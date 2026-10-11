package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"slices"
	"testing"
)

func TestLazyUnaffectedPendingSetterAndSnapshots(t *testing.T) {
	c := &Checker{sourceFilesComplete: true}
	c.couldContainTypeVariables = c.couldContainTypeVariablesWorker
	source := c.newObjectType(ObjectFlagsClass, nil)
	parameter := c.newTypeParameter(nil)
	number := c.newIntrinsicType(TypeFlagsNumber, "number")
	mapper := newSimpleTypeMapper(parameter, number)
	property := &ast.Symbol{Name: "a", Flags: ast.SymbolFlagsProperty}
	setter := &ast.Symbol{Name: "b", Flags: ast.SymbolFlagsSetAccessor}
	generic := &ast.Symbol{Name: "c", Flags: ast.SymbolFlagsProperty}
	members := ast.SymbolTable{"a": property, "b": setter, "c": generic}
	c.valueSymbolLinks.Get(generic).resolvedType = parameter
	first := c.getLazyUnaffectedMembers(source, members, mapper)
	if len(first) != 0 {
		t.Fatal("unresolved and generic members became unaffected")
	}
	c.valueSymbolLinks.Get(property).resolvedType = number
	c.valueSymbolLinks.Get(setter).resolvedType = number
	second := c.getLazyUnaffectedMembers(source, members, mapper)
	if !slices.Equal(second, []string{"a"}) {
		t.Fatal("unresolved setter write type was skipped")
	}
	c.valueSymbolLinks.Get(setter).setWriteType(number)
	third := c.getLazyUnaffectedMembers(source, members, mapper)
	if !slices.Equal(third, []string{"a", "b"}) || !slices.Equal(second, []string{"a"}) || len(first) != 0 {
		t.Fatal("settling a member changed an existing lazy table snapshot")
	}
	if len(c.lazyUnaffected[lazyUnaffectedKey{target: source}].pending) != 0 {
		t.Fatal("settled generic member was rechecked forever")
	}
	added := &ast.Symbol{Name: "d", Flags: ast.SymbolFlagsProperty}
	c.valueSymbolLinks.Get(added).resolvedType = number
	members["d"] = added
	fourth := c.getLazyUnaffectedMembers(source, members, mapper)
	if !slices.Equal(fourth, []string{"a", "b", "d"}) || !slices.Equal(third, []string{"a", "b"}) {
		t.Fatal("table growth changed an earlier snapshot")
	}
}

func BenchmarkLazyUnaffectedMembers(b *testing.B) {
	for _, cached := range []bool{false, true} {
		name := "scan"
		if cached {
			name = "memo"
		}
		b.Run(name, func(b *testing.B) {
			c := &Checker{sourceFilesComplete: true}
			c.couldContainTypeVariables = c.couldContainTypeVariablesWorker
			source := c.newObjectType(ObjectFlagsClass, nil)
			parameter := c.newTypeParameter(nil)
			number := c.newIntrinsicType(TypeFlagsNumber, "number")
			mapper := newSimpleTypeMapper(parameter, number)
			members := make(ast.SymbolTable)
			for i := 0; i < 512; i++ {
				name := string(rune(0x100 + i))
				symbol := &ast.Symbol{Name: name, Flags: ast.SymbolFlagsProperty}
				members[name] = symbol
				c.valueSymbolLinks.Get(symbol).resolvedType = number
			}
			c.getLazyUnaffectedMembers(source, members, mapper)
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if cached {
					c.getLazyUnaffectedMembers(source, members, mapper)
				} else {
					c.computeLazyUnaffectedMembers(members, mapper)
				}
			}
		})
	}
}
