package checker

import (
	"slices"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

func TestFlowMemoSharedHistoryAndDepthBoundary(t *testing.T) {
	c, _, _ := indexedAccessTestChecker()
	factory := ast.NodeFactory{}
	reference := factory.NewIdentifier("x")
	c.symbolNodeLinks.Get(reference).resolvedSymbol = &ast.Symbol{Name: "x", Flags: ast.SymbolFlagsFunctionScopedVariable}
	flow := &ast.FlowNode{Flags: ast.FlowFlagsStart | ast.FlowFlagsShared}
	for range 3 {
		expression := factory.NewToken(ast.KindTrueKeyword)
		expression.Parent = factory.NewIfStatement(expression, factory.NewEmptyStatement(), nil)
		flow = &ast.FlowNode{Flags: ast.FlowFlagsTrueCondition | ast.FlowFlagsShared, Node: expression, Antecedent: flow}
	}
	f := &FlowState{reference: reference, declaredType: c.numberType, initialType: c.numberType}
	c.initializeFlowMemoKey(f)
	result := c.getTypeAtFlowNode(f, flow)
	shared := slices.Clone(c.sharedFlows)
	if result.t != c.numberType || len(shared) != 4 {
		t.Fatal("real flow walk lost shared history")
	}
	c.sharedFlows = c.sharedFlows[:0]
	c.getTypeAtFlowNode(f, flow)
	if c.flowMemoHits == 0 || !slices.Equal(c.sharedFlows, shared) || f.depth != 0 {
		t.Fatal("flow hit changed history or depth")
	}
	hits := c.flowMemoHits
	c.sharedFlows = c.sharedFlows[:0]
	f.depth = 1996
	if c.getTypeAtFlowNode(f, flow) != result || c.flowMemoHits != hits || f.depth != 1996 {
		t.Fatal("memo bypassed the 2000-depth boundary")
	}
	f.depth = 0
	c.sharedFlows = c.sharedFlows[:0]
	f.initialType = c.stringType
	c.initializeFlowMemoKey(f)
	if c.getTypeAtFlowNode(f, flow).t != c.stringType {
		t.Fatal("different initial type shared a flow memo entry")
	}
}
