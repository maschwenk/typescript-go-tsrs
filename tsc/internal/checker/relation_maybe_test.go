package checker

import "testing"

func TestMaybeKeyInlineBoundaryAndUnwind(t *testing.T) {
	r := &Relater{relation: &Relation{}, relationCount: 100}
	for i := 0; i < 12; i++ {
		key := relationKey{hash: CacheHashKey{Lo: uint64(i + 1)}, simple: i%2 == 0}
		r.maybeKeys = append(r.maybeKeys, key)
		if i >= inlineMaybeKeys {
			r.maybeKeysSet.Add(key)
		}
	}
	for _, key := range r.maybeKeys {
		if !r.hasMaybeKey(key) {
			t.Fatal("active assumption missing")
		}
	}
	removed := append([]relationKey(nil), r.maybeKeys[7:]...)
	r.resetMaybeStack(7, RelationComparisonResultReportsUnreliable, true)
	if r.relationCount != 95 || len(r.maybeKeys) != 7 {
		t.Fatal("unwind changed the comparison budget or retained tail")
	}
	for _, key := range removed {
		if r.hasMaybeKey(key) || r.maybeKeysSet.Has(key) {
			t.Fatal("unwound assumption remained active")
		}
		if r.relation.get(key) != RelationComparisonResultSucceeded|RelationComparisonResultReportsUnreliable {
			t.Fatal("assumption result or reliability was lost")
		}
	}
	for _, key := range r.maybeKeys {
		if !r.hasMaybeKey(key) {
			t.Fatal("unwind removed a parent assumption")
		}
	}
	r.resetMaybeStack(0, 0, false)
	if len(r.maybeKeys) != 0 || r.relationCount != 95 {
		t.Fatal("failed unwind changed the budget")
	}
}
