package checker

import "os"

var constituentFlagsMode = os.Getenv("TSGO_CONSTITUENT_FLAGS")

func (c *Checker) maybeTypeOfKindFast(t *Type, kind TypeFlags) bool {
	if t.flags&kind != 0 {
		return true
	}
	if t.flags&TypeFlagsUnionOrIntersection == 0 {
		return false
	}
	const notInPrimitiveUnion = TypeFlagsAny | TypeFlagsUnknown | TypeFlagsVoid | TypeFlagsNever | TypeFlagsObject | TypeFlagsIntersection | TypeFlagsInstantiable
	if kind & ^notInPrimitiveUnion == 0 && t.objectFlags&ObjectFlagsPrimitiveUnion != 0 {
		return false
	}
	if len(t.Types()) >= 8 && c.sourceFilesComplete {
		if flags, ok := c.constituentFlagMemo[t.id]; ok {
			return flags&kind != 0
		}
		flags, stable := constituentFlags(t)
		if stable {
			if c.constituentFlagMemo == nil {
				c.constituentFlagMemo = make(map[TypeId]TypeFlags)
			}
			if len(c.constituentFlagMemo) < 4096 {
				c.constituentFlagMemo[t.id] = flags
			}
			return flags&kind != 0
		}
	}
	for _, child := range t.Types() {
		if child.flags&kind != 0 || child.flags&TypeFlagsUnionOrIntersection != 0 && c.maybeTypeOfKindFast(child, kind) {
			return true
		}
	}
	return false
}

func constituentFlags(t *Type) (TypeFlags, bool) {
	var flags TypeFlags
	stable := t.flags&(TypeFlagsEnumLike|TypeFlagsBooleanLike) == 0
	for _, child := range t.Types() {
		flags |= child.flags
		stable = stable && child.flags&(TypeFlagsEnumLike|TypeFlagsBooleanLike) == 0
		if child.flags&TypeFlagsUnionOrIntersection != 0 {
			nested, nestedStable := constituentFlags(child)
			flags |= nested
			stable = stable && nestedStable
		}
	}
	return flags, stable
}
