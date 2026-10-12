package main

import (
	"math"
	"runtime/debug"
	"testing"
)

func TestBatchGCPolicyScope(t *testing.T) {
	for _, args := range [][]string{{"-w"}, {"--watch", "false"}, {"--WATCH"}, {"-build"}, {"--api"}, {"--lsp"}, {"@args.txt"}} {
		if batchCompilation(args) {
			t.Fatalf("persistent or deferred mode eligible: %v", args)
		}
	}
	if !batchCompilation([]string{"-p", "src", "--noEmit", "--incremental", "false"}) {
		t.Fatal("ordinary compile excluded")
	}
}

func TestBatchGCPolicyRestoresRuntime(t *testing.T) {
	oldPercent := debug.SetGCPercent(100)
	oldLimit := debug.SetMemoryLimit(math.MaxInt64)
	t.Cleanup(func() { debug.SetGCPercent(oldPercent); debug.SetMemoryLimit(oldLimit) })
	setBatchGCPolicy(8 << 30)()
	if got := debug.SetGCPercent(100); got != 100 {
		t.Fatal("small-machine policy changed GC")
	}
	if debug.SetMemoryLimit(-1) != math.MaxInt64 {
		t.Fatal("small-machine policy changed memory budget")
	}
	restore := setBatchGCPolicy(16 << 30)
	if got := debug.SetGCPercent(400); got != 400 {
		t.Fatalf("GC target: %d", got)
	}
	if debug.SetMemoryLimit(-1) != batchMemoryBudget {
		t.Fatal("finite memory budget not applied")
	}
	restore()
	if got := debug.SetGCPercent(100); got != 100 {
		t.Fatal("GC policy not restored")
	}
	if debug.SetMemoryLimit(-1) != math.MaxInt64 {
		t.Fatal("memory budget not restored")
	}
}

func TestBatchGCPolicyPreservesUserOverrides(t *testing.T) {
	for _, variable := range []string{"GOGC", "GOMEMLIMIT"} {
		t.Run(variable, func(t *testing.T) {
			t.Setenv(variable, "")
			old := debug.SetGCPercent(175)
			t.Cleanup(func() { debug.SetGCPercent(old) })
			applyBatchGCPolicy([]string{"-p", "src"})()
			if got := debug.SetGCPercent(175); got != 175 {
				t.Fatal("explicit environment override changed")
			}
		})
	}
}
