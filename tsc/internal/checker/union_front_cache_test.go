package checker

import "testing"

func TestUnionFrontCachePreservesIdentityAndReductionDomain(t *testing.T) {
	c := &Checker{}
	c.sourceFilesComplete = true
	c.strictNullChecks = true
	c.compareSymbols = c.compareSymbolsWorker
	c.unionTypes = make(map[CacheHashKey]*Type)
	c.unionOfUnionTypes = make(map[UnionOfUnionKey]*Type)
	c.neverType = c.newIntrinsicType(TypeFlagsNever, "never")
	c.numberType = c.newIntrinsicType(TypeFlagsNumber, "number")
	c.stringType = c.newIntrinsicType(TypeFlagsString, "string")
	c.errorType = c.newIntrinsicType(TypeFlagsAny, "error")
	c.wildcardType = c.newIntrinsicType(TypeFlagsAny, "wildcard")
	inputs := []*Type{c.numberType, c.stringType}
	first := c.getUnionTypeEx(inputs, UnionReductionNone, nil, nil)
	count := c.TypeCount
	for range 10 {
		if c.getUnionTypeEx(inputs, UnionReductionNone, nil, nil) != first || c.TypeCount != count {
			t.Fatal("front hit changed type identity or creation count")
		}
	}
	other := c.getUnionTypeEx(inputs, UnionReductionLiteral, nil, nil)
	if other != first {
		t.Fatal("equivalent reductions did not reach the interned union")
	}
	if got := c.getUnionTypeEx([]*Type{c.neverType, c.numberType}, UnionReductionLiteral, nil, nil); got != c.numberType {
		t.Fatal("never absorption changed")
	}
	if got := c.getUnionTypeEx([]*Type{c.errorType, c.numberType}, UnionReductionLiteral, nil, nil); got != c.errorType {
		t.Fatal("error absorption changed")
	}
	for _, entry := range c.unionFrontCache.slots {
		if entry.result == c.errorType {
			t.Fatal("error result must be recomputed at its reporting site")
		}
	}
}
