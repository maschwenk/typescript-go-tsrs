package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"os"
)

var propertyLookupsMode = os.Getenv("TSGO_PROPERTY_LOOKUPS")

type sourcePropertyLookup struct {
	target, source *ast.Symbol
	epoch          uint64
}
