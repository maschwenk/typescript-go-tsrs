package checker

import (
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"testing"
)

func indexedAccessTestChecker() (*Checker, *Type, *Type) {
	c, _ := longUnionTestChecker()
	c.initializeClosures()
	c.compilerOptions = &core.CompilerOptions{}
	c.stringLiteralTypes = make(map[string]*Type)
	c.stringType = c.newIntrinsicType(TypeFlagsString, "string")
	c.numberType = c.newIntrinsicType(TypeFlagsNumber, "number")
	c.globalObjectType = c.newAnonymousType(nil, nil, nil, nil, nil)
	members := make(ast.SymbolTable)
	var keys []*Type
	for i := 0; i < 32; i++ {
		name := fmt.Sprintf("p%d", i)
		symbol := &ast.Symbol{Name: name, Flags: ast.SymbolFlagsProperty}
		members[name] = symbol
		c.valueSymbolLinks.Get(symbol).resolvedType = c.numberType
		keys = append(keys, c.getStringLiteralType(name))
	}
	object := c.newAnonymousType(nil, members, nil, nil, nil)
	index := c.getUnionType(keys)
	return c, object, index
}

func TestIndexedAccessMemoIdentityFlagsAndMissingResult(t *testing.T) {
	c, object, index := indexedAccessTestChecker()
	result := c.getIndexedAccessTypeOrUndefined(object, index, AccessFlagsNone, nil, nil)
	if result != c.numberType {
		t.Fatal("property union selected the wrong type")
	}
	before := c.getMemoCounts()
	for range 10 {
		if c.getIndexedAccessTypeOrUndefined(object, index, AccessFlagsNone, nil, nil) != result || c.getMemoCounts() != before {
			t.Fatal("memo changed identity or effects")
		}
	}
	c.getIndexedAccessTypeOrUndefined(object, index, AccessFlagsWriting, nil, nil)
	if len(c.indexedAccessUnionMemo) != 2 {
		t.Fatal("read and write requests were conflated")
	}
	missing := c.getUnionType([]*Type{c.getStringLiteralType("p0"), c.getStringLiteralType("missing")})
	if c.getIndexedAccessTypeOrUndefined(object, missing, AccessFlagsNone, nil, nil) != nil {
		t.Fatal("missing property was accepted")
	}
	key := indexedAccessUnionKey{object.id, missing.id, AccessFlagsNone}
	if value, ok := c.indexedAccessUnionMemo[key]; !ok || value != nil {
		t.Fatal("negative answer was not distinguished from a miss")
	}
}

func BenchmarkIndexedAccessUnionMemo(b *testing.B) {
	for _, cached := range []bool{false, true} {
		name := "walk"
		if cached {
			name = "memo"
		}
		b.Run(name, func(b *testing.B) {
			c, object, index := indexedAccessTestChecker()
			c.getIndexedAccessUnionMemoized(object, index, AccessFlagsNone)
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if cached {
					c.getIndexedAccessUnionMemoized(object, index, AccessFlagsNone)
				} else {
					c.getIndexedAccessTypeUncached(object, index, AccessFlagsNone, nil, nil)
				}
			}
		})
	}
}
