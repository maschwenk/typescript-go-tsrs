package checker

import (
	"encoding/binary"
	"os"
	"slices"
)

var inferMemoMode = os.Getenv("TSGO_INFER_MEMO")

type inferMemoOutcome struct {
	candidates, contraCandidates []*Type
	priority                     InferencePriority
	topLevel                     bool
}

type inferMemoEntry struct {
	infos        []inferMemoOutcome
	cleared      bool
	epoch, steps uint64
}

type inferenceEffects struct {
	counts                          memoCounts
	expressions, transient, members uint64
	flow                            int
	reliability                     RelationComparisonResult
}

func (c *Checker) inferenceEffects() inferenceEffects {
	return inferenceEffects{c.getMemoCounts(), c.memoExpressionChecks, c.memoTransient, c.memberResolutionEpoch, c.flowInvocationCount, c.reliabilityFlags}
}

func (c *Checker) inferTypes(infos []*InferenceInfo, source, target *Type, priority InferencePriority, contra bool) {
	eligible := c.sourceFilesComplete && inferMemoMode != "off" && inferMemoMode != "0" && len(infos) <= 8 &&
		len(c.instantiationStack) == 0 && len(c.typeResolutions) == 0 && len(c.flowLoopStack) == 0 && c.skipDirectInferenceNodes.Len() == 0
	if !eligible || !c.inferMemoLongTargets[target.id] {
		steps := c.inferMemoSteps
		c.inferTypesUncached(infos, source, target, priority, contra)
		if eligible && c.inferMemoSteps-steps >= 16 && len(c.inferMemoLongTargets) < 4096 {
			if c.inferMemoLongTargets == nil {
				c.inferMemoLongTargets = make(map[TypeId]bool)
			}
			c.inferMemoLongTargets[target.id] = true
		}
		return
	}
	key := c.inferMemoScratch[:0]
	c.inferMemoScratch = nil
	defer func() { c.inferMemoScratch = key[:0] }()
	word := func(n uint32) { key = binary.LittleEndian.AppendUint32(key, n) }
	word(uint32(source.id))
	word(uint32(target.id))
	word(uint32(priority))
	word(uint32(c.reliabilityFlags))
	if contra {
		word(1)
	} else {
		word(0)
	}
	word(uint32(len(infos)))
	for _, info := range infos {
		word(uint32(info.typeParameter.id))
		word(uint32(info.priority))
		flags := uint32(0)
		if info.topLevel {
			flags |= 1
		}
		if info.isFixed {
			flags |= 2
		}
		word(flags)
		key = binary.LittleEndian.AppendUint64(key, uint64(info.impliedArity))
		word(uint32(len(info.candidates)))
		for _, t := range info.candidates {
			word(uint32(t.id))
		}
		word(uint32(len(info.contraCandidates)))
		for _, t := range info.contraCandidates {
			word(uint32(t.id))
		}
	}
	if len(key) > 16384 {
		c.inferTypesUncached(infos, source, target, priority, contra)
		return
	}
	if entry := c.inferMemoEntries[string(key)]; entry != nil && entry.epoch == c.memberResolutionEpoch {
		if inferMemoMode == "shadow" {
			fresh := make([]*InferenceInfo, len(infos))
			for i := range infos {
				fresh[i] = cloneInferenceInfo(infos[i])
			}
			before := c.inferenceEffects()
			cleared := c.inferTypesUncached(fresh, source, target, priority, contra)
			if before != c.inferenceEffects() || cleared != entry.cleared {
				panic("inference memo skipped checker effects")
			}
			for i, info := range fresh {
				kept := entry.infos[i]
				if !slices.Equal(info.candidates, kept.candidates) || !slices.Equal(info.contraCandidates, kept.contraCandidates) || info.priority != kept.priority || info.topLevel != kept.topLevel {
					panic("inference memo changed candidate order or priority")
				}
			}
		}
		if entry.cleared {
			clearCachedInferences(infos)
		}
		for i, info := range infos {
			kept := entry.infos[i]
			info.candidates = append(info.candidates[:0], kept.candidates...)
			info.contraCandidates = append(info.contraCandidates[:0], kept.contraCandidates...)
			info.priority = kept.priority
			info.topLevel = kept.topLevel
		}
		c.inferMemoHits++
		c.inferMemoSkipped += entry.steps
		return
	}
	before := c.inferenceEffects()
	steps := c.inferMemoSteps
	cleared := c.inferTypesUncached(infos, source, target, priority, contra)
	steps = c.inferMemoSteps - steps
	if before != c.inferenceEffects() || steps < 16 || len(c.inferMemoEntries) >= 2048 {
		return
	}
	bytes := len(key) + len(infos)*64
	for _, info := range infos {
		bytes += (len(info.candidates) + len(info.contraCandidates)) * 8
	}
	if c.inferMemoBytes+bytes > 8<<20 {
		return
	}
	entry := &inferMemoEntry{cleared: cleared, epoch: c.memberResolutionEpoch, steps: steps, infos: make([]inferMemoOutcome, len(infos))}
	for i, info := range infos {
		entry.infos[i] = inferMemoOutcome{slices.Clone(info.candidates), slices.Clone(info.contraCandidates), info.priority, info.topLevel}
	}
	if c.inferMemoEntries == nil {
		c.inferMemoEntries = make(map[string]*inferMemoEntry)
	}
	c.inferMemoEntries[string(key)] = entry
	c.inferMemoBytes += bytes
}
