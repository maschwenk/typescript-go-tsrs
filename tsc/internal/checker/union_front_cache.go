package checker

import "os"

// Alias-bearing calls keep the original path so this cache assigns no symbol IDs.
var unionFrontMode = os.Getenv("TSGO_UNION_CACHE")

type unionFrontKey struct {
	ids       [8]TypeId
	count     uint32
	reduction UnionReduction
}

type unionFrontEntry struct {
	key    unionFrontKey
	result *Type
}

type unionFrontCache struct{ slots [1024]unionFrontEntry }

func (c *Checker) getUnionTypeFrontCached(types []*Type, reduction UnionReduction) *Type {
	key := unionFrontKey{count: uint32(len(types)), reduction: reduction}
	hash := uint32(reduction)
	for i, t := range types {
		key.ids[i] = t.id
		hash = (hash ^ uint32(t.id)) * 16777619
	}
	if c.unionFrontCache == nil {
		c.unionFrontCache = &unionFrontCache{}
	}
	slot := hash & 1023
	entry := c.unionFrontCache.slots[slot]
	if entry.result != nil && entry.key == key {
		if unionFrontMode == "shadow" {
			count, total, inst, symbols, sigs := c.TypeCount, c.TotalInstantiationCount, c.instantiationCount, c.SymbolCount, c.SignatureCount
			got := c.getUnionTypeUncached(types, reduction, nil, nil)
			if got != entry.result || count != c.TypeCount || total != c.TotalInstantiationCount || inst != c.instantiationCount || symbols != c.SymbolCount || sigs != c.SignatureCount {
				panic("union front cache changed identity or counters")
			}
		}
		return entry.result
	}
	count, total, inst, impure, symbols, sigs := c.TypeCount, c.TotalInstantiationCount, c.instantiationCount, c.unionCacheImpure, c.SymbolCount, c.SignatureCount
	result := c.getUnionTypeUncached(types, reduction, nil, nil)
	createdOnlyResult := c.TypeCount == count || c.TypeCount == count+1 && result.flags&TypeFlagsUnion != 0 && uint32(result.id) == c.TypeCount
	if createdOnlyResult && symbols == c.SymbolCount && sigs == c.SignatureCount && total == c.TotalInstantiationCount && inst == c.instantiationCount && impure == c.unionCacheImpure && result != c.errorType {
		c.unionFrontCache.slots[slot] = unionFrontEntry{key: key, result: result}
	}
	return result
}
