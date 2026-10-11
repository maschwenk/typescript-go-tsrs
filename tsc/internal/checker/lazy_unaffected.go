package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"os"
	"slices"
)

var lazyUnaffectedMode = os.Getenv("TSGO_LAZY_UNAFFECTED")

type lazyUnaffectedKey struct {
	target       *Type
	mapsThisOnly bool
}
type lazyUnaffectedMemo struct {
	length  int
	settled []string
	pending []*ast.Symbol
}

func (c *Checker) getLazyUnaffectedMembers(source *Type, members ast.SymbolTable, mapper *TypeMapper) []string {
	if !c.sourceFilesComplete || lazyUnaffectedMode == "0" || lazyUnaffectedMode == "off" {
		return c.computeLazyUnaffectedMembers(members, mapper)
	}
	key := lazyUnaffectedKey{target: source, mapsThisOnly: mapper.MapsThisOnly()}
	memo := c.lazyUnaffected[key]
	// Declared tables are immutable once resolved in a complete Program. The length
	// guard also handles a target whose table was not complete on the earlier query.
	if memo == nil || memo.length != len(members) {
		memo = &lazyUnaffectedMemo{length: len(members)}
		for name, symbol := range members {
			if c.isNamedMember(symbol, name) {
				unaffected, settled := c.settledUnaffectedMember(symbol, mapper)
				if unaffected {
					memo.settled = append(memo.settled, name)
				} else if !settled {
					memo.pending = append(memo.pending, symbol)
				}
			}
		}
		slices.Sort(memo.settled)
		if c.lazyUnaffected == nil {
			c.lazyUnaffected = make(map[lazyUnaffectedKey]*lazyUnaffectedMemo)
		}
		c.lazyUnaffected[key] = memo
	} else if len(memo.pending) != 0 {
		var added []string
		pending := memo.pending[:0]
		for _, symbol := range memo.pending {
			unaffected, settled := c.settledUnaffectedMember(symbol, mapper)
			if unaffected {
				added = append(added, symbol.Name)
			} else if !settled {
				pending = append(pending, symbol)
			}
		}
		clear(memo.pending[len(pending):])
		memo.pending = pending
		if len(added) != 0 {
			// Existing lazy tables retain their earlier snapshot.
			added = append(added, memo.settled...)
			slices.Sort(added)
			memo.settled = added
		}
	}
	if lazyUnaffectedMode == "shadow" && !slices.Equal(memo.settled, c.computeLazyUnaffectedMembers(members, mapper)) {
		panic("lazy unaffected memo changed the declared-member snapshot")
	}
	return memo.settled
}

func (c *Checker) computeLazyUnaffectedMembers(members ast.SymbolTable, mapper *TypeMapper) []string {
	var names []string
	for name, symbol := range members {
		if c.isNamedMember(symbol, name) && c.isSymbolUnaffectedByInstantiation(symbol, mapper) {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	return names
}

func (c *Checker) settledUnaffectedMember(symbol *ast.Symbol, mapper *TypeMapper) (unaffected, settled bool) {
	if c.isSymbolUnaffectedByInstantiation(symbol, mapper) {
		return true, true
	}
	links := c.valueSymbolLinks.Get(symbol)
	return false, links.resolvedType != nil && (symbol.Flags&ast.SymbolFlagsSetAccessor == 0 || links.getWriteType() != nil)
}
