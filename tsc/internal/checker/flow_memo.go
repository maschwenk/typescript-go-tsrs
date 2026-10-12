package checker

import (
	"os"
	"slices"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

var flowMemoMode = os.Getenv("TSGO_FLOW_MEMO")

type flowMemoKey struct {
	flow              *ast.FlowNode
	symbol            *ast.Symbol
	declared, initial *Type
	container         *ast.Node
	kind              ast.Kind
	reliability       RelationComparisonResult
}

type flowMemoEntry struct {
	result       FlowType
	shared       []SharedFlow
	height       int
	epoch, steps uint64
}

func (c *Checker) initializeFlowMemoKey(f *FlowState) {
	if !c.sourceFilesComplete || flowMemoMode == "off" || flowMemoMode == "0" {
		return
	}
	reference := f.reference
	if reference.Kind != ast.KindIdentifier && reference.Kind != ast.KindThisKeyword {
		return
	}
	var symbol *ast.Symbol
	if reference.Kind == ast.KindIdentifier {
		if ast.IsThisInTypeQuery(reference) {
			return
		}
		// Identifier checking resolved this link before starting its flow walk.
		symbol = c.getResolvedSymbolOrNil(reference)
		if symbol == nil || symbol == c.unknownSymbol {
			return
		}
	}
	f.memoKey = flowMemoKey{symbol: symbol, declared: f.declaredType, initial: f.initialType, container: f.flowContainer, kind: reference.Kind}
	f.memoEnabled = true
}

func (c *Checker) getTypeAtFlowNode(f *FlowState, flow *ast.FlowNode) FlowType {
	if !f.memoEnabled || f.memoDisabled {
		return c.getTypeAtFlowNodeUncached(f, flow)
	}
	depth := f.depth
	f.memoDepths = append(f.memoDepths, depth+1)
	defer func() {
		last := len(f.memoDepths) - 1
		height := f.memoDepths[last]
		f.memoDepths = f.memoDepths[:last]
		if last > 0 {
			f.memoDepths[last-1] = max(f.memoDepths[last-1], height)
		}
	}()
	eligible := len(c.instantiationStack) == 0 && len(c.typeResolutions) == 0 && len(c.flowLoopStack) == 0 &&
		c.inlineLevel == 0 && len(f.reduceLabels) == 0 && len(c.sharedFlows) == f.sharedFlowStart && !c.flowAnalysisDisabled
	if !eligible {
		return c.getTypeAtFlowNodeUncached(f, flow)
	}
	key := f.memoKey
	key.flow = flow
	key.reliability = c.reliabilityFlags
	if entry := c.flowMemoEntries[key]; entry != nil && entry.epoch == c.memberResolutionEpoch && depth+entry.height < 2000 {
		if flowMemoMode == "shadow" {
			before := c.inferenceEffects()
			start := len(c.sharedFlows)
			f.memoDisabled = true
			fresh := c.getTypeAtFlowNodeUncached(f, flow)
			f.memoDisabled = false
			if fresh != entry.result || !slices.Equal(c.sharedFlows[start:], entry.shared) || before != c.inferenceEffects() {
				panic("flow memo changed type, shared-flow history or checker effects")
			}
			c.sharedFlows = c.sharedFlows[:start]
		}
		// Reproduce all shared nodes, in order, so later depth checks see the same history.
		c.sharedFlows = append(c.sharedFlows, entry.shared...)
		f.memoDepths[len(f.memoDepths)-1] = depth + entry.height
		c.flowMemoHits++
		c.flowMemoSkipped += entry.steps
		return entry.result
	}
	before := c.inferenceEffects()
	steps := c.flowMemoSteps
	start := len(c.sharedFlows)
	result := c.getTypeAtFlowNodeUncached(f, flow)
	steps = c.flowMemoSteps - steps
	shared := c.sharedFlows[start:]
	bytes := 128 + len(shared)*24
	if result.incomplete || result.t == c.errorType || before != c.inferenceEffects() || len(shared) > 64 || len(c.flowMemoEntries) >= 8192 || c.flowMemoBytes+bytes > 8<<20 {
		return result
	}
	if c.flowMemoEntries == nil {
		c.flowMemoEntries = make(map[flowMemoKey]*flowMemoEntry)
	}
	c.flowMemoEntries[key] = &flowMemoEntry{result: result, shared: slices.Clone(shared), height: f.memoDepths[len(f.memoDepths)-1] - depth, epoch: c.memberResolutionEpoch, steps: steps}
	c.flowMemoBytes += bytes
	return result
}
