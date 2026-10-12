package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/core"
)

func (c *Checker) emitModuleFormatOfFile(file *ast.SourceFile) core.ModuleKind {
	if !c.sourceFilesComplete {
		return c.program.GetEmitModuleFormatOfFile(file)
	}
	if file != nil && file == c.lastEmitModuleFile {
		return c.lastEmitModuleKind
	}
	kind := c.program.GetEmitModuleFormatOfFile(file)
	c.lastEmitModuleFile, c.lastEmitModuleKind = file, kind
	return kind
}
