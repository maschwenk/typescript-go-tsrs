package checker

import (
	"fmt"
	"slices"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

func TestInferenceMemoCandidatesAndIsolation(t *testing.T) {
	c, source, _ := indexedAccessTestChecker()
	parameter := c.newTypeParameter(nil)
	members := make(ast.SymbolTable)
	for i := 0; i < 32; i++ {
		name := fmt.Sprintf("p%d", i)
		symbol := &ast.Symbol{Name: name, Flags: ast.SymbolFlagsProperty}
		members[name] = symbol
		c.valueSymbolLinks.Get(symbol).resolvedType = parameter
	}
	target := c.newAnonymousType(nil, members, nil, nil, nil)
	target.objectFlags |= ObjectFlagsCouldContainTypeVariablesComputed | ObjectFlagsCouldContainTypeVariables
	for range 4 {
		info := newInferenceInfo(parameter)
		info.inferredType = c.stringType
		c.inferTypes([]*InferenceInfo{info}, source, target, InferencePriorityNone, false)
		if !slices.Equal(info.candidates, []*Type{c.numberType}) || info.inferredType != nil {
			t.Fatal("inference lost candidates or failed to clear the inferred type")
		}
		info.candidates[0] = c.stringType
	}
	if c.inferMemoHits == 0 || c.inferMemoSkipped < 16 {
		t.Fatal("repeated real inference did not hit the memo")
	}
	hits := c.inferMemoHits
	c.memberResolutionEpoch++
	c.inferTypes([]*InferenceInfo{newInferenceInfo(parameter)}, source, target, InferencePriorityNone, false)
	if c.inferMemoHits != hits {
		t.Fatal("member resolution change reused an old inference")
	}
	info := newInferenceInfo(parameter)
	info.isFixed = true
	c.inferTypes([]*InferenceInfo{info}, source, target, InferencePriorityNone, false)
	if len(info.candidates) != 0 {
		t.Fatal("fixed inference reused an unfixed result")
	}
}
