package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/core"
)

// nodeLinkStore is a links store keyed by node references. Values are stored directly
// in the pages of the store which is suitable for values where sizeof(V) is small.
type nodeLinkStore[V any] struct {
	store core.PagedLinkStore[V]
}

func (s *nodeLinkStore[V]) Get(node *ast.Node) *V {
	return s.store.Get(uint64(ast.GetNodeId(node)))
}

func (s *nodeLinkStore[V]) Has(node *ast.Node) bool {
	return s.store.Has(uint64(ast.GetNodeId(node)))
}

func (s *nodeLinkStore[V]) TryGet(node *ast.Node) *V {
	return s.store.TryGet(uint64(ast.GetNodeId(node)))
}

// symbolArenaLinkStore retains existing symbol ID assignment and stable values.
type symbolArenaLinkStore[V any] struct {
	store core.PagedArenaLinkStore[V]
}

func (s *symbolArenaLinkStore[V]) Get(symbol *ast.Symbol) *V {
	return s.store.Get(uint64(ast.GetSymbolId(symbol)))
}

func (s *symbolArenaLinkStore[V]) Has(symbol *ast.Symbol) bool {
	return s.store.Has(uint64(ast.GetSymbolId(symbol)))
}

func (s *symbolArenaLinkStore[V]) TryGet(symbol *ast.Symbol) *V {
	return s.store.TryGet(uint64(ast.GetSymbolId(symbol)))
}
