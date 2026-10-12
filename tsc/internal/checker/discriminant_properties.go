package checker

import (
	"os"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

var discriminantPropertiesMode = os.Getenv("TSGO_DISCRIMINANT_PROPERTIES")

type settledProperty struct {
	symbol *ast.Symbol
	known  bool
}

type discriminantPropertyEntry struct {
	epoch      uint64
	properties []settledProperty
}

func (c *Checker) settledPropertyOfType(t *Type, name string) settledProperty {
	if t.flags != TypeFlagsObject || t.objectFlags&(ObjectFlagsMembersResolved|ObjectFlagsMapped) != ObjectFlagsMembersResolved ||
		isTupleType(t) || t.symbol != nil && t.symbol.Flags&ast.SymbolFlagsValueModule != 0 {
		return settledProperty{}
	}
	if symbol := t.AsStructuredType().members[name]; symbol != nil {
		return settledProperty{symbol, symbol.Flags&ast.SymbolFlagsValue != 0}
	}
	return settledProperty{known: c.canSkipContextualConstituent(t, name)}
}

func (c *Checker) discriminantProperties(target *Type, types []*Type, name string) *discriminantPropertyEntry {
	if !c.sourceFilesComplete || len(types) < 8 || len(types) > 4096 || len(c.instantiationStack) != 0 || len(c.typeResolutions) != 0 ||
		discriminantPropertiesMode == "off" || discriminantPropertiesMode == "0" {
		return nil
	}
	key := contextualUnionKey{target.id, name}
	if entry := c.discriminantPropertyCache[key]; entry != nil && entry.epoch == c.memberResolutionEpoch {
		return entry
	}
	if len(c.discriminantPropertyCache) >= 1024 || c.discriminantPropertyBytes+len(types)*16+len(name)+64 > 1<<20 {
		return nil
	}
	entry := &discriminantPropertyEntry{epoch: c.memberResolutionEpoch, properties: make([]settledProperty, len(types))}
	for i, t := range types {
		entry.properties[i] = c.settledPropertyOfType(t, name)
	}
	if c.discriminantPropertyCache == nil {
		c.discriminantPropertyCache = make(map[contextualUnionKey]*discriminantPropertyEntry)
	}
	c.discriminantPropertyCache[key] = entry
	c.discriminantPropertyBytes += len(types)*16 + len(name) + 64
	return entry
}

func (c *Checker) typeOfDiscriminantProperty(entry *discriminantPropertyEntry, i int, t *Type, name string) *Type {
	if entry == nil || entry.epoch != c.memberResolutionEpoch || !entry.properties[i].known {
		return c.getTypeOfPropertyOrIndexSignatureOfType(t, name)
	}
	property := entry.properties[i].symbol
	if discriminantPropertiesMode == "shadow" {
		before := c.inferenceEffects()
		if fresh := c.getPropertyOfType(t, name); fresh != property || before != c.inferenceEffects() {
			panic("discriminant property reuse changed symbol or checker effects")
		}
		if property == nil && c.getTypeOfPropertyOrIndexSignatureOfType(t, name) != nil {
			panic("discriminant property reuse skipped an index signature")
		}
		if before != c.inferenceEffects() {
			panic("discriminant property reuse skipped lookup effects")
		}
	}
	if property == nil {
		return nil
	}
	// Read the symbol's current type normally: contextual symbol types need not be immutable.
	return c.getTypeOfSymbol(property)
}
