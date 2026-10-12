package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"testing"
)

func TestPackedTypeArgumentsSeparateLengthOrderAndFlags(t *testing.T) {
	values := []*Type{{id: 0}, {id: 1}, {id: 1 << 31}, {id: TypeId(^uint32(0))}}
	seen := make(map[CacheHashKey]bool)
	var visit func([]*Type)
	visit = func(args []*Type) {
		for _, flag := range []bool{false, true} {
			for _, tag := range []uint64{0x1a5e, 0xc0d7} {
				key := packedTypeArgumentKey(tag, args, flag)
				if seen[key] {
					t.Fatalf("packed key collided at length %d", len(args))
				}
				seen[key] = true
			}
		}
		if len(args) < 3 {
			for _, v := range values {
				visit(append(args, v))
			}
		}
	}
	visit(nil)
	alias := &TypeAlias{symbol: &ast.Symbol{Name: "A", Flags: ast.SymbolFlagsTypeAlias}, typeArguments: values[:1]}
	for _, args := range [][]*Type{values[:1], values} {
		for _, a := range []*TypeAlias{nil, alias} {
			if a == nil && len(args) <= 3 {
				continue
			}
			for _, flag := range []bool{false, true} {
				var b keyBuilder
				b.writeTypes(args)
				b.writeAlias(a)
				if flag {
					b.writeByte('!')
				}
				if getTypeInstantiationKey(args, a, flag) != b.hash() || getConditionalTypeKey(args, a, flag) != b.hash() {
					t.Fatal("long or aliased fallback changed its key")
				}
			}
		}
	}
}
