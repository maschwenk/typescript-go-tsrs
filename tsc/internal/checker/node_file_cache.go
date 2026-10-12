package checker

import (
	"os"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

var nodeFileCacheMode = os.Getenv("TSGO_NODE_FILE_CACHE")

type nodeFileInfo struct {
	parent *ast.Node
	file   *ast.SourceFile
	index  int
}

func (c *Checker) sourceFileAndIndexOfNode(n *ast.Node) nodeFileInfo {
	eligible := c.sourceFilesComplete && !ast.NodeIsSynthesized(n) && nodeFileCacheMode != "off" && nodeFileCacheMode != "0"
	if eligible {
		if info, ok := c.nodeFileCache[n]; ok && info.parent == n.Parent {
			if nodeFileCacheMode == "shadow" {
				file := ast.GetSourceFileOfNode(n)
				if file != info.file || c.fileIndexMap[file] != info.index {
					panic("node file cache changed declaration order")
				}
			}
			return info
		}
	}
	file := ast.GetSourceFileOfNode(n)
	info := nodeFileInfo{n.Parent, file, c.fileIndexMap[file]}
	if _, indexed := c.fileIndexMap[file]; eligible && indexed && (len(c.nodeFileCache) < 4096 || c.nodeFileCache[n].file != nil) {
		if c.nodeFileCache == nil {
			c.nodeFileCache = make(map[*ast.Node]nodeFileInfo)
		}
		c.nodeFileCache[n] = info
	}
	return info
}
