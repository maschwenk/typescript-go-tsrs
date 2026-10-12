package checker

import "testing"

func TestNormalizationMemoKeepsIntersectionIdentity(t *testing.T) {
	c := &Checker{}
	c.compareSymbols = c.compareSymbolsWorker
	tp := c.newTypeParameter(nil)
	nilType := c.newIntrinsicType(TypeFlagsNull, "null")
	object := c.newAnonymousType(nil, nil, nil, nil, nil)
	for _, constituents := range [][]*Type{{tp, nilType}, {tp, object}, {nilType, object}, {tp}} {
		typ := c.newIntersectionType(ObjectFlagsNone, constituents)
		expected := c.shouldNormalizeIntersectionWorker(typ)
		count := c.TypeCount
		for range 5 {
			if c.shouldNormalizeIntersection(typ) != expected {
				t.Fatal("memo changed normalization decision")
			}
		}
		if count != c.TypeCount {
			t.Fatal("normalization predicate allocated a type")
		}
		clone := c.newIntersectionType(typ.objectFlags, constituents)
		if clone.objectFlags&ObjectFlagsShouldNormalizeComputed != 0 {
			t.Fatal("new intersection inherited stale memo")
		}
	}
}
