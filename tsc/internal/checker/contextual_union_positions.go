package checker

import (
	"fmt"
	"os"
)

var contextualPositionsMode = os.Getenv("TSGO_CONTEXTUAL_POSITIONS")

type contextualUnionKey struct {
	union TypeId
	name  string
}
type contextualUnionEntry struct {
	epoch     uint64
	positions []uint32
}

func (c *Checker) canSkipContextualConstituent(t *Type, name string) bool {
	if t.flags != TypeFlagsObject || t.objectFlags&(ObjectFlagsMapped|ObjectFlagsMembersResolved) != ObjectFlagsMembersResolved || isTupleType(t) {
		return false
	}
	data := t.AsStructuredType()
	// Callable types can be augmented by several global function interfaces.
	if data.members[name] != nil || len(data.indexInfos) != 0 || len(data.signatures) != 0 {
		return false
	}
	object := c.globalObjectType
	return object != nil && object.flags == TypeFlagsObject && object.objectFlags&ObjectFlagsMembersResolved != 0 && object.AsStructuredType().members[name] == nil
}

func (c *Checker) getContextualUnionPropertyFiltered(t *Type, name string, nameType *Type) (*Type, bool) {
	if !c.sourceFilesComplete || t.flags&TypeFlagsUnion == 0 || len(c.typeResolutions) != 0 || len(c.instantiationStack) != 0 || contextualPositionsMode == "off" || contextualPositionsMode == "0" {
		return nil, false
	}
	u := t.AsUnionType()
	types := u.types
	if u.origin != nil && u.origin.flags&TypeFlagsUnion != 0 {
		types = u.origin.Types()
	}
	if len(types) < 8 || len(types) > 4096 {
		return nil, false
	}
	key := contextualUnionKey{t.id, name}
	entry, ok := c.contextualUnionPositions[key]
	if !ok || entry.epoch != c.memberResolutionEpoch {
		positions := make([]uint32, 0, len(types))
		for i, s := range types {
			if !c.canSkipContextualConstituent(s, name) {
				positions = append(positions, uint32(i))
			}
		}
		entry = contextualUnionEntry{c.memberResolutionEpoch, positions}
		bytes := len(positions)*4 + len(name) + 64
		if len(c.contextualUnionPositions) < 1024 && c.contextualUnionBytes+bytes <= 1<<20 {
			if c.contextualUnionPositions == nil {
				c.contextualUnionPositions = make(map[contextualUnionKey]contextualUnionEntry)
			}
			c.contextualUnionPositions[key] = entry
			c.contextualUnionBytes += bytes
		}
	}
	if len(entry.positions) == len(types) {
		return nil, false
	}
	if contextualPositionsMode == "shadow" {
		position := 0
		for i, s := range types {
			if position < len(entry.positions) && entry.positions[position] == uint32(i) {
				position++
				continue
			}
			before := c.inferenceEffects()
			fresh := c.getContextualPropertyConstituent(s, name, nameType)
			if fresh != nil || before != c.inferenceEffects() {
				panic(fmt.Sprintf("contextual union skipped a property answer or checker effects: result=%p before=%+v after=%+v", fresh, before, c.inferenceEffects()))
			}
		}
	}
	mapped := make([]*Type, 0, len(entry.positions))
	for _, position := range entry.positions {
		s := types[position]
		var result *Type
		if s.flags&TypeFlagsUnion != 0 {
			result = c.mapTypeEx(s, func(t *Type) *Type { return c.getContextualPropertyConstituent(t, name, nameType) }, true)
		} else {
			result = c.getContextualPropertyConstituent(s, name, nameType)
		}
		if result != nil {
			mapped = append(mapped, result)
		}
	}
	if len(mapped) == 0 {
		return nil, true
	}
	// At least one original constituent returned nil, so mapTypeEx's changed path applies.
	return c.getUnionTypeEx(mapped, UnionReductionNone, nil, nil), true
}
