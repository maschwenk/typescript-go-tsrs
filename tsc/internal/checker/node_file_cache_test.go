package checker

import (
	"testing"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/parser"
)

func TestNodeFileCacheProgramOrderAndReparenting(t *testing.T) {
	first := parser.ParseSourceFile(ast.SourceFileParseOptions{FileName: "/z.ts"}, "type Z = number;", core.ScriptKindTS)
	second := parser.ParseSourceFile(ast.SourceFileParseOptions{FileName: "/a.ts"}, "type A = string;", core.ScriptKindTS)
	ast.SetParentInChildren(first.AsNode())
	ast.SetParentInChildren(second.AsNode())
	c := &Checker{sourceFilesComplete: true, fileIndexMap: map[*ast.SourceFile]int{first: 0, second: 1}}
	a, b := first.Statements.Nodes[0], second.Statements.Nodes[0]
	for range 3 {
		if c.compareNodes(a, b) >= 0 || c.compareNodes(b, a) <= 0 {
			t.Fatal("node order ignored containing program file order")
		}
	}
	if len(c.nodeFileCache) != 2 {
		t.Fatal("published declarations were not cached")
	}
	a.Parent = second.AsNode()
	if info := c.sourceFileAndIndexOfNode(a); info.file != second || info.index != 1 {
		t.Fatal("direct reparenting reused a stale source file")
	}
	unknown := parser.ParseSourceFile(ast.SourceFileParseOptions{FileName: "/unknown.ts"}, "type U = boolean;", core.ScriptKindTS)
	ast.SetParentInChildren(unknown.AsNode())
	n := unknown.Statements.Nodes[0]
	if c.sourceFileAndIndexOfNode(n).index != 0 {
		t.Fatal("unknown file index changed")
	}
	if _, ok := c.nodeFileCache[n]; ok {
		t.Fatal("unknown program file was cached")
	}
}
