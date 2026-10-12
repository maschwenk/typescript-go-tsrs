//go:build !linux

package main

func effectiveMemory(physical uint64) uint64 { return physical }
