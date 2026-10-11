package checker

import (
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"testing"
)

func TestActiveMapperCacheKeepsAliasDomainsAndFrameLifetime(t *testing.T) {
	c := &Checker{}
	c.couldContainTypeVariables = c.couldContainTypeVariablesWorker
	source := c.newTypeParameter(nil)
	target := c.newIntrinsicType(TypeFlagsNumber, "number")
	m := newSimpleTypeMapper(source, target)
	alias := &TypeAlias{symbol: &ast.Symbol{Name: "A"}}
	c.pushActiveMapper(m)
	for _, a := range []*TypeAlias{nil, alias, nil, alias} {
		if got := c.instantiateTypeWithAlias(source, m, a); got != target {
			t.Fatal("wrong mapped type")
		}
	}
	if c.TotalInstantiationCount != 2 {
		t.Fatalf("alias and alias-free results must each be computed once: %d", c.TotalInstantiationCount)
	}
	c.clearActiveMapperCaches()
	c.instantiateTypeWithAlias(source, m, nil)
	if c.TotalInstantiationCount != 3 {
		t.Fatal("explicit cache reset retained a result")
	}
	c.popActiveMapper()
	c.pushActiveMapper(m)
	c.instantiateTypeWithAlias(source, m, nil)
	if c.TotalInstantiationCount != 4 {
		t.Fatal("a reused mapper frame retained a result")
	}
	c.popActiveMapper()
	c.instantiateTypeWithAlias(source, m, nil)
	if c.TotalInstantiationCount != 5 {
		t.Fatal("a new mapper frame reused a recursive result")
	}
}
