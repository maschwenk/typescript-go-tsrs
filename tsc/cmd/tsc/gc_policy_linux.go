//go:build linux

package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func effectiveMemory(physical uint64) uint64 {
	return effectiveMemoryFrom(physical, "/proc/self/cgroup", "/sys/fs/cgroup")
}

func effectiveMemoryFrom(physical uint64, groupFile, root string) uint64 {
	groups, err := os.ReadFile(groupFile)
	if err != nil {
		return 0
	}
	for _, group := range strings.Split(string(groups), "\n") {
		parts := strings.SplitN(group, ":", 3)
		if len(parts) != 3 {
			continue
		}
		base, file := root, "memory.max"
		if parts[0] != "0" || parts[1] != "" {
			if !strings.Contains(","+parts[1]+",", ",memory,") {
				continue
			}
			base, file = filepath.Join(root, "memory"), "memory.limit_in_bytes"
		}
		seenLimit := false
		directory := filepath.Join(base, strings.TrimPrefix(filepath.Clean("/"+parts[2]), "/"))
		for {
			value, err := os.ReadFile(filepath.Join(directory, file))
			if err == nil {
				seenLimit = true
				text := strings.TrimSpace(string(value))
				if text != "max" {
					limit, err := strconv.ParseUint(text, 10, 64)
					if err != nil {
						return 0
					}
					physical = min(physical, limit)
				}
			} else if !os.IsNotExist(err) {
				return 0
			}
			if directory == base {
				break
			}
			directory = filepath.Dir(directory)
		}
		if !seenLimit {
			return 0
		}
	}
	return physical
}
