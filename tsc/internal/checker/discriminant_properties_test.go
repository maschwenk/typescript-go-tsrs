package checker

import (
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

func TestDiscriminantPropertiesReadCurrentTypesAndInvalidateMembers(t *testing.T) {
	c, _, _ := indexedAccessTestChecker()
	var types []*Type
	var symbols []*ast.Symbol
	for range 32 {
		symbol := &ast.Symbol{Name: "kind", Flags: ast.SymbolFlagsProperty}
		c.valueSymbolLinks.Get(symbol).resolvedType = c.numberType
		symbols = append(symbols, symbol)
		types = append(types, c.newAnonymousType(nil, ast.SymbolTable{"kind": symbol}, nil, nil, nil))
	}
	union := c.getUnionType(types)
	types = union.Types()
	entry := c.discriminantProperties(union, types, "kind")
	if entry == nil || !entry.properties[0].known || c.typeOfDiscriminantProperty(entry, 0, types[0], "kind") != c.numberType {
		t.Fatal("settled property was not reused")
	}
	property := entry.properties[0].symbol
	c.valueSymbolLinks.Get(property).resolvedType = c.stringType
	if c.typeOfDiscriminantProperty(entry, 0, types[0], "kind") != c.stringType {
		t.Fatal("reuse returned a stale contextual symbol type")
	}
	c.setStructuredTypeMembers(types[0], nil, nil, nil, nil)
	if c.typeOfDiscriminantProperty(entry, 0, types[0], "kind") != nil {
		t.Fatal("member epoch change retained a removed property")
	}
	fresh := c.discriminantProperties(union, types, "kind")
	if fresh == nil || fresh == entry || !fresh.properties[0].known || fresh.properties[0].symbol != nil {
		t.Fatal("removed property did not become a known absence")
	}
}
