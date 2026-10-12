package main

import (
	"os"
	"runtime/debug"
	"strings"

	osmemory "github.com/mackerelio/go-osstat/memory"
)

const batchMemoryBudget = 8 << 30

func batchCompilation(args []string) bool {
	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			return false
		}
		switch strings.ToLower(arg) {
		case "-w", "--w", "-watch", "--watch", "-b", "--b", "-build", "--build", "--lsp", "--api":
			return false
		}
	}
	return true
}

func applyBatchGCPolicy(args []string) func() {
	if !batchCompilation(args) {
		return func() {}
	}
	if _, configured := os.LookupEnv("GOGC"); configured {
		return func() {}
	}
	if _, configured := os.LookupEnv("GOMEMLIMIT"); configured {
		return func() {}
	}
	memory, err := osmemory.Get()
	if err != nil {
		return func() {}
	}
	return setBatchGCPolicy(effectiveMemory(memory.Total))
}

func setBatchGCPolicy(memory uint64) func() {
	if memory < 12<<30 {
		return func() {}
	}
	// Larger one-shot checks benefit from fewer scans of their retained graphs.
	// This remains a soft Go heap budget; it is not a bound on process RSS.
	oldLimit := debug.SetMemoryLimit(batchMemoryBudget)
	oldPercent := debug.SetGCPercent(400)
	return func() {
		debug.SetGCPercent(oldPercent)
		debug.SetMemoryLimit(oldLimit)
	}
}
