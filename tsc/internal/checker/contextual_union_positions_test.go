package checker

import (
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

func TestContextualUnionPositionsEpochAndObjectAugmentation(t *testing.T) {
	c, _, _ := indexedAccessTestChecker()
	property := &ast.Symbol{Name: "wanted", Flags: ast.SymbolFlagsProperty}
	c.valueSymbolLinks.Get(property).resolvedType = c.numberType
	var types []*Type
	for range 32 {
		types = append(types, c.newAnonymousType(nil, nil, nil, nil, nil))
	}
	c.setStructuredTypeMembers(types[5], ast.SymbolTable{"wanted": property}, nil, nil, nil)
	union := c.getUnionType(types)
	if c.getTypeOfPropertyOfContextualType(union, "wanted") != c.numberType {
		t.Fatal("filtered union lost its property")
	}
	entry := c.contextualUnionPositions[contextualUnionKey{union.id, "wanted"}]
	if len(entry.positions) != 1 {
		t.Fatal("resolved absent constituents were not skipped")
	}
	c.setStructuredTypeMembers(types[11], ast.SymbolTable{"wanted": property}, nil, nil, nil)
	c.getTypeOfPropertyOfContextualType(union, "wanted")
	entry = c.contextualUnionPositions[contextualUnionKey{union.id, "wanted"}]
	if len(entry.positions) != 2 || entry.epoch != c.memberResolutionEpoch {
		t.Fatal("member changes reused stale positions")
	}
	augmentation := &ast.Symbol{Name: "toString", Flags: ast.SymbolFlagsProperty}
	c.valueSymbolLinks.Get(augmentation).resolvedType = c.stringType
	c.setStructuredTypeMembers(c.globalObjectType, ast.SymbolTable{"toString": augmentation}, nil, nil, nil)
	if c.getTypeOfPropertyOfContextualType(union, "toString") != c.stringType {
		t.Fatal("filter skipped an inherited Object member")
	}
	if c.getTypeOfPropertyOfContextualType(union, "missing") != nil {
		t.Fatal("absent name acquired a contextual type")
	}
}
