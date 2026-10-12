package checker

import "os"

var recursionStackMode = os.Getenv("TSGO_DEEP_NESTED")

type relaterStackEntry struct {
	t        *Type
	target   *Type
	identity RecursionId
}

// Each active relation entry computes its identity at the first original lookup.
func (c *Checker) isDeeplyNestedOnStack(t *Type, stack []relaterStackEntry, maxDepth int) bool {
	if len(stack) < maxDepth {
		return false
	}
	if recursionStackMode == "off" || recursionStackMode == "0" {
		types := make([]*Type, len(stack))
		for i := range stack {
			types[i] = stack[i].t
		}
		return c.isDeeplyNestedType(t, types, maxDepth)
	}
	result := c.isDeeplyNestedOnStackKept(t, stack, maxDepth)
	if recursionStackMode == "shadow" {
		before := c.getMemoCounts()
		types := make([]*Type, len(stack))
		for i := range stack {
			types[i] = stack[i].t
		}
		if result != c.isDeeplyNestedType(t, types, maxDepth) {
			panic("kept recursion identities changed depth result")
		}
		for i := range stack {
			e := &stack[i]
			if e.target != nil {
				target := getRecursionIdentityTarget(e.t)
				if target != e.target || target.flags&TypeFlagsIntersection == 0 && getRecursionIdentityFromTarget(target) != e.identity {
					panic("active recursion identity changed")
				}
			}
		}
		if before != c.getMemoCounts() {
			panic("kept recursion identities skipped checker effects")
		}
	}
	return result
}

func (c *Checker) isDeeplyNestedOnStackKept(t *Type, stack []relaterStackEntry, maxDepth int) bool {
	target := getRecursionIdentityTarget(t)
	if target.flags&TypeFlagsIntersection != 0 {
		for _, constituent := range target.Types() {
			if c.isDeeplyNestedOnStackKept(constituent, stack, maxDepth) {
				return true
			}
		}
		return false
	}
	identity := getRecursionIdentityFromTarget(target)
	count := 0
	lastTypeId := TypeId(0)
	for i := range stack {
		e := &stack[i]
		if e.target == nil {
			e.target = getRecursionIdentityTarget(e.t)
			if e.target.flags&TypeFlagsIntersection == 0 {
				e.identity = getRecursionIdentityFromTarget(e.target)
			}
		}
		matches := e.identity == identity
		if e.target.flags&TypeFlagsIntersection != 0 {
			matches = false
			for _, constituent := range e.target.Types() {
				if hasMatchingRecursionIdentity(constituent, identity) {
					matches = true
					break
				}
			}
		}
		if matches {
			if e.t.id >= lastTypeId {
				count++
				if count >= maxDepth {
					return true
				}
			}
			lastTypeId = e.t.id
		}
	}
	return false
}
