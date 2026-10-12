package checker

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/parser"
)

func TestPrecomputedSymbolOrderAndMissingFileFallback(t *testing.T) {
	c := &Checker{sourceFilesComplete: true, fileIndexMap: make(map[*ast.SourceFile]int)}
	c.compareSymbols = c.compareSymbolsWorker
	var symbols []*ast.Symbol
	text := strings.Builder{}
	for i := range 32 {
		fmt.Fprintf(&text, "type T%d = number;\n", i)
	}
	for i, name := range []string{"/z.ts", "/a.ts"} {
		file := parser.ParseSourceFile(ast.SourceFileParseOptions{FileName: name}, text.String(), core.ScriptKindTS)
		ast.SetParentInChildren(file.AsNode())
		c.fileIndexMap[file] = i
		for j, declaration := range file.Statements.Nodes {
			symbols = append(symbols, &ast.Symbol{Name: fmt.Sprintf("n%d", j%3), Declarations: []*ast.Node{declaration}})
		}
	}
	symbols = append(symbols, &ast.Symbol{Name: "without"})
	for _, s := range symbols {
		ast.GetSymbolId(s)
	}
	slices.Reverse(symbols)
	expected := slices.Clone(symbols)
	slices.SortFunc(expected, c.compareSymbols)
	if !c.sortSymbolsByKey(symbols) || !slices.Equal(symbols, expected) {
		t.Fatal("precomputed sort changed declaration/file/name/ID order")
	}
	for _, key := range c.symbolSortScratch[:cap(c.symbolSortScratch)] {
		if key.symbol != nil {
			t.Fatal("sort scratch retained symbols")
		}
	}
	c.fileIndexMap = nil
	before := slices.Clone(symbols)
	if c.sortSymbolsByKey(symbols) || !slices.Equal(symbols, before) {
		t.Fatal("missing file index did not preserve the original fallback")
	}
}
