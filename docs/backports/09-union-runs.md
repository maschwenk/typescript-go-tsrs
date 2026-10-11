# Build unions by merging existing sorted runs

Source: tsrs `notes/perf-tsgo-can1357-speedups.md; crates/tsrs_checker/src/checker_13.rs merge_ascending_runs`.

Filtering a flattened union preserves its order. Track ascending runs at input boundaries, binary-insert a single extra constituent, and merge other runs pairwise, with binary searches for highly unequal sizes. Preserve stable sorting for hosts whose reachable source files are not complete: the auto-import resolver can load additional files, where the comparator is not a total order. The merge tests compare type identity and order against the original stable sort over 500 generated cases.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.009x | 147.1 → 147.1 | 67327 → 67327 | 460252 → 459443 |
| Compiler | checkers4 | 0.985x | 165.2 → 171.5 | 91011 → 91011 | 633068 → 632257 |
| Compiler-Unions | single | 0.970x | 172.0 → 165.7 | 63690 → 63690 | 485053 → 484488 |
| Compiler-Unions | checkers4 | 1.022x | 183.6 → 186.4 | 88982 → 88982 | 672971 → 672520 |
| xstate-main | single | 1.004x | 617.4 → 627.3 | 569112 → 569112 | 4498194 → 4508505 |
| xstate-main | checkers4 | 0.999x | 682.4 → 691.9 | 711882 → 711882 | 5852456 → 5867453 |
| webpack | single | 0.995x | 1166.3 → 1184.9 | 733733 → 733733 | 8338207 → 8343791 |
| webpack | checkers4 | 1.006x | 1177.6 → 1239.1 | 921944 → 921944 | 9761718 → 9767465 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/09-union-runs.json).
