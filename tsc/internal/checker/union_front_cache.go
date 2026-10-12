package checker

import (
	"encoding/binary"
	"os"
)

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

type unionFrontCache struct {
	slots     [1024]unionFrontEntry
	long      map[string]*Type
	longBytes int
	probe     []byte
}

const maxLongUnionWords = 4096
const maxLongUnionBytes = 1 << 20

func (c *Checker) getUnionTypeFrontCached(types []*Type, reduction UnionReduction) *Type {
	if len(types) > 8 {
		return c.getLongUnionFrontCached(types, reduction)
	}

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

func (c *Checker) getLongUnionFrontCached(types []*Type, reduction UnionReduction) *Type {
	if c.unionFrontCache == nil {
		c.unionFrontCache = &unionFrontCache{}
	}
	cache := c.unionFrontCache
	// Take the scratch buffer: recursive union construction needs a separate probe.
	probe := cache.probe
	cache.probe = nil
	size := 4 * (len(types) + 1)
	if cap(probe) < size {
		probe = make([]byte, size)
	} else {
		probe = probe[:size]
	}
	binary.LittleEndian.PutUint32(probe, uint32(reduction))
	for i, t := range types {
		binary.LittleEndian.PutUint32(probe[4*(i+1):], uint32(t.id))
	}
	defer func() {
		if cap(probe) >= cap(cache.probe) {
			cache.probe = probe
		}
	}()
	if result := cache.long[string(probe)]; result != nil {
		if unionFrontMode == "shadow" {
			count, total, inst, impure, symbols, sigs := c.TypeCount, c.TotalInstantiationCount, c.instantiationCount, c.unionCacheImpure, c.SymbolCount, c.SignatureCount
			got := c.getUnionTypeUncached(types, reduction, nil, nil)
			if got != result || count != c.TypeCount || total != c.TotalInstantiationCount || inst != c.instantiationCount || impure != c.unionCacheImpure || symbols != c.SymbolCount || sigs != c.SignatureCount {
				panic("long union front cache changed identity or counters")
			}
		}
		return result
	}
	count, total, inst, impure, symbols, sigs := c.TypeCount, c.TotalInstantiationCount, c.instantiationCount, c.unionCacheImpure, c.SymbolCount, c.SignatureCount
	result := c.getUnionTypeUncached(types, reduction, nil, nil)
	createdOnlyResult := c.TypeCount == count || c.TypeCount == count+1 && result.flags&TypeFlagsUnion != 0 && uint32(result.id) == c.TypeCount
	if createdOnlyResult && symbols == c.SymbolCount && sigs == c.SignatureCount && total == c.TotalInstantiationCount && inst == c.instantiationCount && impure == c.unionCacheImpure && result != c.errorType && cache.longBytes+size <= maxLongUnionBytes && len(cache.long) < 4096 {
		if cache.long == nil {
			cache.long = make(map[string]*Type)
		}
		if cache.long[string(probe)] == nil {
			cache.longBytes += size
		}
		cache.long[string(probe)] = result
	}
	return result
}
