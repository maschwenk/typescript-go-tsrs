# Scan path segments in eight-byte words

Source: tsrs `notes/perf-checker-cpu4.md; crates/tsrs_tspath/src/path.rs: hasRelativePathSegment`.

The common clean pnpm path takes about 12.4 ns instead of 66.7 ns in paired function benchmarks on this machine, with zero allocations in both versions. Use portable byte loads and verify slash/dot candidates exactly. Exhaustive path strings, Unicode and raw bytes agree with the original semantics; whole compiler results are measured separately.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.008x | 146.5 → 145.8 | 67327 → 67327 | 446849 → 446852 |
| Compiler | checkers4 | 0.996x | 171.6 → 171.7 | 91011 → 91011 | 596327 → 596339 |
| webpack | single | 1.013x | 1162.5 → 1162.3 | 733733 → 733733 | 8103298 → 8103304 |
| webpack | checkers4 | 0.994x | 1145.2 → 1207.1 | 921944 → 921944 | 9452292 → 9457374 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/20-word-path-scan.json).
