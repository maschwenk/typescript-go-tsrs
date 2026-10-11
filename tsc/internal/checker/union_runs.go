package checker

import (
	"math/bits"
	"slices"
)

// Flattened unions are sorted already. Only input boundaries can begin another
// run. This requires all reachable source files to have checker file indices.
func mergeAscendingTypeRuns(types []*Type, starts []int) []*Type {
	n := len(types)
	if len(starts) == 1 && (starts[0] == 1 || starts[0] == n-1) {
		var t *Type
		if starts[0] == 1 {
			t = types[0]
			copy(types, types[1:])
		} else {
			t = types[n-1]
		}
		i, _ := slices.BinarySearchFunc(types[:n-1], t, CompareTypes)
		copy(types[i+1:], types[i:n-1])
		types[i] = t
		return types
	}
	bounds := make([]int, 0, len(starts)+2)
	bounds = append(bounds, 0)
	bounds = append(bounds, starts...)
	bounds = append(bounds, n)
	src, dst := types, make([]*Type, n)
	for len(bounds) > 2 {
		next := bounds[:1]
		runs := len(bounds) - 1
		for j := 0; j < runs; j += 2 {
			start, mid := bounds[j], bounds[j+1]
			end := mid
			if j+1 < runs {
				end = bounds[j+2]
				mergeSortedTypeRuns(dst[start:end], src[start:mid], src[mid:end])
			} else {
				copy(dst[start:end], src[start:end])
			}
			next = append(next, end)
		}
		src, dst = dst, src
		bounds = next
	}
	return src
}

func mergeSortedTypeRuns(out, left, right []*Type) {
	short, long := left, right
	if len(short) > len(long) {
		short, long = long, short
	}
	if len(short) > 0 && len(short)*bits.Len(uint(len(long)))*2 < len(long) {
		i := 0
		for _, t := range short {
			pos, _ := slices.BinarySearchFunc(long, t, CompareTypes)
			i += copy(out[i:], long[:pos])
			long = long[pos:]
			out[i] = t
			i++
			if len(long) > 0 && long[0] == t {
				out[i] = t
				i++
				long = long[1:]
			}
		}
		copy(out[i:], long)
		return
	}
	i := 0
	for len(left) > 0 && len(right) > 0 {
		if left[0] == right[0] || CompareTypes(left[0], right[0]) < 0 {
			out[i] = left[0]
			left = left[1:]
		} else {
			out[i] = right[0]
			right = right[1:]
		}
		i++
	}
	i += copy(out[i:], left)
	copy(out[i:], right)
}
