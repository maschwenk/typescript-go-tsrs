package checker

import (
	"os"

	"github.com/microsoft/TypeScript/tsc/internal/ast"
)

var contextualReuseMode = os.Getenv("TSGO_CONTEXTUAL_REUSE")

// Keep the contextual answer for this node, including a computed nil answer.
func (c *Checker) isConstContextWithType(node *ast.Node) (bool, *Type, bool) {
	parent := node.Parent
	if ast.IsConstAssertion(parent) || c.isInlineImportAttributes(node) {
		return true, nil, false
	}
	var contextual *Type
	computed := false
	before := c.getMemoCounts()
	if c.isValidConstAssertionArgument(node) {
		contextual = c.getContextualType(node, ContextFlagsNone)
		computed = len(c.instantiationStack) == 0 && len(c.typeResolutions) == 0 && len(c.inferenceContextInfos) == 0
		if c.isConstTypeVariable(contextual, 0) {
			return true, contextual, computed
		}
	}
	result := (ast.IsParenthesizedExpression(parent) || ast.IsArrayLiteralExpression(parent) || ast.IsSpreadElement(parent)) && c.isConstContext(parent) ||
		(ast.IsPropertyAssignment(parent) || ast.IsShorthandPropertyAssignment(parent) || ast.IsTemplateSpan(parent)) && c.isConstContext(parent.Parent)
	return result, contextual, computed && before == c.getMemoCounts()
}
