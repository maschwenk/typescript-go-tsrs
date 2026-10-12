package checker

import (
	"cmp"
	"slices"
	"strings"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

type symbolSortKey struct {
	symbol    *ast.Symbol
	file, pos int
	declared  bool
}

func (c *Checker) sortSymbolsByKey(symbols []*ast.Symbol) bool {
	keyed := c.symbolSortScratch[:0]
	defer func() { clear(keyed); c.symbolSortScratch = keyed[:0] }()
	for _, symbol := range symbols {
		if symbol == nil {
			return false
		}
		key := symbolSortKey{symbol: symbol}
		if len(symbol.Declarations) != 0 {
			declaration := symbol.Declarations[0]
			file := ast.GetSourceFileOfNode(declaration)
			index, ok := c.fileIndexMap[file]
			if file == nil || !ok {
				return false
			}
			key.file, key.pos, key.declared = index, declaration.Pos(), true
		}
		keyed = append(keyed, key)
	}
	slices.SortFunc(keyed, func(a, b symbolSortKey) int {
		if a.symbol == b.symbol {
			return 0
		}
		if a.declared != b.declared {
			if a.declared {
				return -1
			}
			return 1
		}
		if a.declared {
			if result := cmp.Compare(a.file, b.file); result != 0 {
				return result
			}
			if result := cmp.Compare(a.pos, b.pos); result != 0 {
				return result
			}
		}
		if result := strings.Compare(a.symbol.Name, b.symbol.Name); result != 0 {
			return result
		}
		// Assign IDs only on the same ties reached by the original comparator.
		return int(ast.GetSymbolId(a.symbol)) - int(ast.GetSymbolId(b.symbol))
	})
	for i := range symbols {
		symbols[i] = keyed[i].symbol
	}
	return true
}
