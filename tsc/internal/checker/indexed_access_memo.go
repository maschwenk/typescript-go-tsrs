package checker

import "os"

var indexedAccessMemoMode = os.Getenv("TSGO_IA_MEMO")

type indexedAccessUnionKey struct {
	object, index TypeId
	flags         AccessFlags
}

type memoCounts struct {
	types, symbols, signatures, total, inst uint32
	impure, diagnostics                     uint64
}

func (c *Checker) getMemoCounts() memoCounts {
	return memoCounts{c.TypeCount, c.SymbolCount, c.SignatureCount, c.TotalInstantiationCount, c.instantiationCount, c.unionCacheImpure, c.memoDiagnosticAttempts}
}

func (c *Checker) getIndexedAccessUnionMemoized(objectType, indexType *Type, flags AccessFlags) *Type {
	key := indexedAccessUnionKey{objectType.id, indexType.id, flags}
	if result, ok := c.indexedAccessUnionMemo[key]; ok {
		if indexedAccessMemoMode == "shadow" {
			before := c.getMemoCounts()
			fresh := c.getIndexedAccessTypeUncached(objectType, indexType, flags, nil, nil)
			if fresh != result || before != c.getMemoCounts() {
				panic("indexed access union memo changed identity, budgets or diagnostic attempts")
			}
		}
		return result
	}
	before := c.getMemoCounts()
	result := c.getIndexedAccessTypeUncached(objectType, indexType, flags, nil, nil)
	after := c.getMemoCounts()
	createdOnlyResult := after.types == before.types || result != nil && after.types == before.types+1 && uint32(result.id) == after.types && result.flags&(TypeFlagsUnionOrIntersection|TypeFlagsIndexedAccess) != 0
	after.types = before.types
	if createdOnlyResult && after == before && result != c.errorType && len(c.indexedAccessUnionMemo) < 8192 {
		if c.indexedAccessUnionMemo == nil {
			c.indexedAccessUnionMemo = make(map[indexedAccessUnionKey]*Type)
		}
		c.indexedAccessUnionMemo[key] = result
	}
	return result
}
