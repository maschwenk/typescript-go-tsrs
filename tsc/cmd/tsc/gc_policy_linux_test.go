//go:build linux

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEffectiveMemoryRespectsCgroupAncestors(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		t.Run(map[bool]string{false: "unified", true: "legacy"}[legacy], func(t *testing.T) {
			root := t.TempDir()
			base, file, group := root, "memory.max", "0::/parent/child\n"
			if legacy {
				base, file, group = filepath.Join(root, "memory"), "memory.limit_in_bytes", "4:cpu,memory:/parent/child\n"
			}
			if err := os.MkdirAll(filepath.Join(base, "parent/child"), 0o755); err != nil {
				t.Fatal(err)
			}
			write := func(path, value string) {
				t.Helper()
				if err := os.WriteFile(path, []byte(value), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			groups := filepath.Join(root, "groups")
			write(groups, group)
			write(filepath.Join(base, file), "17179869184")
			write(filepath.Join(base, "parent", file), "4294967296")
			write(filepath.Join(base, "parent/child", file), "8589934592")
			if got := effectiveMemoryFrom(128<<30, groups, root); got != 4<<30 {
				t.Fatalf("ancestor limit ignored: %d", got)
			}
			if got := effectiveMemoryFrom(2<<30, groups, root); got != 2<<30 {
				t.Fatalf("physical memory ignored: %d", got)
			}
			write(filepath.Join(base, "parent", file), "invalid")
			if got := effectiveMemoryFrom(128<<30, groups, root); got != 0 {
				t.Fatalf("unknown budget enabled policy: %d", got)
			}
		})
	}
}
