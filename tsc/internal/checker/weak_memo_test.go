package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"testing"
)

func TestWeakTypeMemoIsResetWithMembers(t *testing.T) {
	c := &Checker{}
	c.compareSymbols = c.compareSymbolsWorker
	optional := &ast.Symbol{Name: "x", Flags: ast.SymbolFlagsProperty | ast.SymbolFlagsOptional}
	required := &ast.Symbol{Name: "x", Flags: ast.SymbolFlagsProperty}
	typ := c.newAnonymousType(nil, ast.SymbolTable{"x": optional}, nil, nil, nil)
	if !c.isWeakType(typ) || !c.isWeakType(typ) {
		t.Fatal("optional property type should be weak")
	}
	c.setStructuredTypeMembers(typ, ast.SymbolTable{"x": required}, nil, nil, nil)
	if c.isWeakType(typ) {
		t.Fatal("required member inherited stale weak result")
	}
	c.setStructuredTypeMembers(typ, nil, nil, nil, nil)
	if c.isWeakType(typ) {
		t.Fatal("empty object is not weak")
	}
	c.setStructuredTypeMembers(typ, ast.SymbolTable{"x": optional}, []*Signature{{}}, nil, nil)
	if c.isWeakType(typ) {
		t.Fatal("call signature is not weak")
	}
	clone := c.newObjectType(ObjectFlagsReference, nil)
	clone.objectFlags |= ObjectFlagsMembersResolved | ObjectFlagsIsWeakTypeComputed | ObjectFlagsIsWeakType
	copied := c.cloneTypeReference(clone)
	if copied.objectFlags&(ObjectFlagsWeakTypeMemo|ObjectFlagsMembersResolved) != 0 {
		t.Fatal("unresolved clone inherited computed members")
	}
}
