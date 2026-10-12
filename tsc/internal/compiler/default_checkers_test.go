package compiler

import "testing"

func TestDefaultCheckersRespectMachineAndCheckedWork(t *testing.T) {
	for _, test := range []struct{ cores, files, want int }{
		{1, 4096, 4}, {8, 4096, 8}, {18, 4096, 8}, {33, 4096, 8}, {34, 4096, 8}, {64, 4096, 8}, {128, 4096, 8},
		{18, 31, 4}, {18, 127, 4}, {18, 159, 4}, {18, 160, 5}, {18, 255, 7}, {18, 256, 8}, {18, 511, 8}, {18, 512, 8}, {64, 511, 8},
	} {
		if got := defaultCheckersForWorkload(test.cores, test.files); got != test.want {
			t.Fatalf("%d cores/%d files: got %d want %d", test.cores, test.files, got, test.want)
		}
	}
}
