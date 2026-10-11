package checker

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestMergeTypeRunsMatchesStableSort(t *testing.T) {
	c := &Checker{}
	pool := make([]*Type, 80)
	for i := range pool {
		pool[i] = c.newIntrinsicType(TypeFlagsNumber, "number")
	}
	rng := rand.New(rand.NewPCG(17, 29))
	for trial := 0; trial < 500; trial++ {
		var input []*Type
		var starts []int
		for run := 0; run < 1+rng.IntN(20); run++ {
			if len(input) > 0 {
				starts = append(starts, len(input))
			}
			chunk := make([]*Type, 1+rng.IntN(150))
			for i := range chunk {
				chunk[i] = pool[rng.IntN(len(pool))]
			}
			slices.SortStableFunc(chunk, CompareTypes)
			input = append(input, chunk...)
		}
		want := slices.Clone(input)
		slices.SortStableFunc(want, CompareTypes)
		got := mergeAscendingTypeRuns(input, starts)
		if !slices.Equal(got, want) {
			t.Fatalf("trial %d: merged type identities differ", trial)
		}
	}
	for _, first := range []bool{true, false} {
		input := slices.Clone(pool)
		if first {
			input = append([]*Type{pool[70]}, input...)
		} else {
			input = append(input, pool[10])
		}
		want := slices.Clone(input)
		slices.SortStableFunc(want, CompareTypes)
		mid := len(input) - 1
		if first {
			mid = 1
		}
		if !slices.Equal(mergeAscendingTypeRuns(input, []int{mid}), want) {
			t.Fatal("single constituent insertion differs")
		}
	}
}
