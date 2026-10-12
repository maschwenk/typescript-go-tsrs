package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"testing"
)

func TestConstituentFlagsNestedKindsAndEnumExclusion(t *testing.T) {
	c, object, keys := indexedAccessTestChecker()
	types := append([]*Type(nil), keys.Types()...)
	types = append(types, object)
	union := c.getUnionType(types)
	for _, kind := range []TypeFlags{TypeFlagsObject, TypeFlagsStringLiteral, TypeFlagsNumber, TypeFlagsUnion, TypeFlagsAny, TypeFlagsInstantiable, TypeFlagsObject | TypeFlagsVoid} {
		if c.maybeTypeOfKind(union, kind) != c.maybeTypeOfKindUncached(union, kind) {
			t.Fatalf("kind %v changed", kind)
		}
	}
	if _, ok := c.constituentFlagMemo[union.id]; !ok {
		t.Fatal("settled large union was not retained")
	}
	intersection := c.newIntersectionType(ObjectFlagsNone, []*Type{object, union})
	for _, kind := range []TypeFlags{TypeFlagsObject, TypeFlagsStringLiteral, TypeFlagsUnion, TypeFlagsIntersection} {
		if c.maybeTypeOfKind(intersection, kind) != c.maybeTypeOfKindUncached(intersection, kind) {
			t.Fatal("nested constituent was lost")
		}
	}
	c.enumLiteralTypes = make(map[EnumLiteralKey]*Type)
	enumSymbol := &ast.Symbol{Name: "E", Flags: ast.SymbolFlagsRegularEnum}
	member := &ast.Symbol{Name: "M", Flags: ast.SymbolFlagsEnumMember}
	enum := c.getEnumLiteralType("member", enumSymbol, member)
	enumUnion := c.getUnionType(append(types, enum))
	c.maybeTypeOfKind(enumUnion, TypeFlagsObject)
	if _, ok := c.constituentFlagMemo[enumUnion.id]; ok {
		t.Fatal("enum normalization can change flags")
	}
	for _, kind := range []TypeFlags{TypeFlagsVoid, TypeFlagsObject, TypeFlagsAny, TypeFlagsIntersection, TypeFlagsInstantiable} {
		if c.maybeTypeOfKind(keys, kind) != c.maybeTypeOfKindUncached(keys, kind) {
			t.Fatal("primitive union exclusion changed a query")
		}
	}
}
