package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/jsnum"
	"slices"
	"testing"
)

func longUnionTestChecker() (*Checker, []*Type) {
	c := &Checker{sourceFilesComplete: true, strictNullChecks: true}
	c.compareSymbols = c.compareSymbolsWorker
	c.unionTypes = make(map[CacheHashKey]*Type)
	c.unionOfUnionTypes = make(map[UnionOfUnionKey]*Type)
	c.numberLiteralTypes = make(map[jsnum.Number]*Type)
	c.neverType = c.newIntrinsicType(TypeFlagsNever, "never")
	c.errorType = c.newIntrinsicType(TypeFlagsAny, "error")
	c.wildcardType = c.newIntrinsicType(TypeFlagsAny, "wildcard")
	var inputs []*Type
	for i := 0; i < 128; i++ {
		inputs = append(inputs, c.getNumberLiteralType(jsnum.Number(i)))
	}
	return c, inputs
}

func TestLongUnionCacheExactOrderReductionAndMemoryBound(t *testing.T) {
	c, inputs := longUnionTestChecker()
	result := c.getUnionTypeEx(inputs, UnionReductionNone, nil, nil)
	count := c.TypeCount
	for range 10 {
		if c.getUnionTypeEx(inputs, UnionReductionNone, nil, nil) != result || c.TypeCount != count {
			t.Fatal("long-key hit changed identity or creation count")
		}
	}
	reversed := slices.Clone(inputs)
	slices.Reverse(reversed)
	c.getUnionTypeEx(reversed, UnionReductionNone, nil, nil)
	c.getUnionTypeEx(inputs, UnionReductionLiteral, nil, nil)
	if len(c.unionFrontCache.long) != 3 {
		t.Fatal("input order or reduction keys were conflated")
	}
	c.unionFrontCache.longBytes = maxLongUnionBytes
	altered := append(slices.Clone(inputs), c.neverType)
	before := len(c.unionFrontCache.long)
	if c.getUnionTypeEx(altered, UnionReductionNone, nil, nil) != result || len(c.unionFrontCache.long) != before {
		t.Fatal("memory cap changed the answer or kept storing")
	}
	if c.getUnionTypeEx(inputs, UnionReductionNone, nil, nil) != result {
		t.Fatal("memory cap disabled existing hits")
	}
}

func BenchmarkLongUnionCache(b *testing.B) {
	for _, cached := range []bool{false, true} {
		name := "uncached"
		if cached {
			name = "front"
		}
		b.Run(name, func(b *testing.B) {
			c, inputs := longUnionTestChecker()
			c.getUnionTypeFrontCached(inputs, UnionReductionNone)
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if cached {
					c.getUnionTypeFrontCached(inputs, UnionReductionNone)
				} else {
					c.getUnionTypeUncached(inputs, UnionReductionNone, nil, nil)
				}
			}
		})
	}
}
