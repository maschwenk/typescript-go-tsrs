package checker

import "testing"

func TestValueLinkTailRetainsIndependentRareValues(t *testing.T) {
	links := ValueSymbolLinks{}
	if links.getWriteType() != nil || links.getNameType() != nil || links.getContainingType() != nil || links.getFunctionOrConstructorChecked() {
		t.Fatal("cold rare values are not empty")
	}
	links.setWriteType(nil)
	if links.rare != nil {
		t.Fatal("empty write allocated a tail")
	}
	first, second := &Type{}, &Type{}
	links.setWriteType(first)
	links.setNameType(second)
	links.setContainingType(first)
	links.setFunctionOrConstructorChecked(true)
	if links.getWriteType() != first || links.getNameType() != second || links.getContainingType() != first || !links.getFunctionOrConstructorChecked() {
		t.Fatal("rare fields overlap")
	}
	links.setNameType(nil)
	if links.getNameType() != nil || links.getWriteType() != first {
		t.Fatal("clearing a name changed the write type")
	}
}
