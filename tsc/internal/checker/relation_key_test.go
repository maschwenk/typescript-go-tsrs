package checker

import "testing"

func TestRelationKeyDomainsAndBoundaries(t *testing.T) {
	relation := &Relation{}
	for _, sourceID := range []TypeId{1, 1<<31 - 1, 1 << 31, ^TypeId(0)} {
		for _, targetID := range []TypeId{1, 1<<31 - 1, 1 << 31, ^TypeId(0)} {
			for _, state := range []IntersectionState{0, 1, 2, 3, 4} {
				source, target := &Type{id: sourceID}, &Type{id: targetID}
				key, constrained := getRelationKey(source, target, state, false, false)
				if constrained {
					t.Fatal("simple pair became constrained")
				}
				if key.simple != (sourceID < 1<<31 && targetID < 1<<31 && state < 4) {
					t.Fatal("packing guard lost boundary bits")
				}
				relation.set(key, RelationComparisonResultSucceeded)
				if relation.get(key) != RelationComparisonResultSucceeded {
					t.Fatal("relation cache lost the pair")
				}
				reverse, _ := getRelationKey(target, source, state, true, false)
				ordered, _ := getRelationKey(source, target, state, true, false)
				if reverse != ordered {
					t.Fatal("identity relation lost symmetry")
				}
			}
		}
	}
	simple := relationKey{hash: CacheHashKey{Lo: 42}, simple: true}
	hashed := relationKey{hash: CacheHashKey{Lo: 42}}
	relation.set(simple, RelationComparisonResultSucceeded)
	relation.set(hashed, RelationComparisonResultFailed)
	if relation.get(simple) != RelationComparisonResultSucceeded || relation.get(hashed) != RelationComparisonResultFailed {
		t.Fatal("packed and hashed domains collided")
	}
}
