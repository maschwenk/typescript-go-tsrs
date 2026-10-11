# Reuse one membership set for large named unions

Source: tsrs `notes/perf-icount-hillclimb.md; crates/tsrs_checker/src/checker_13.rs get_union_type_worker`.

Large overlapping named unions currently binary-search every named union for every flattened constituent. Build one operation-scoped identity set and retain the existing reduced-type and origin order. Tiny unions and incomplete auto-import hosts keep the original path, avoiding map overhead and preserving partial-comparator behavior.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.996x | 147.9 → 144.1 | 67327 → 67327 | 459443 → 459439 |
| Compiler | checkers4 | 1.003x | 171.9 → 175.7 | 91011 → 91011 | 632092 → 632156 |
| Compiler-Unions | single | 0.991x | 164.8 → 170.2 | 63690 → 63690 | 484502 → 484646 |
| Compiler-Unions | checkers4 | 1.001x | 194.5 → 185.2 | 88982 → 88982 | 672516 → 673156 |
| xstate-main | single | 0.997x | 623.0 → 625.5 | 569112 → 569112 | 4508526 → 4508502 |
| xstate-main | checkers4 | 0.994x | 693.2 → 691.3 | 711882 → 711882 | 5867113 → 5873831 |
| webpack | single | 1.011x | 1160.7 → 1166.6 | 733733 → 733733 | 8343818 → 8343813 |
| webpack | checkers4 | 1.004x | 1269.7 → 1272.4 | 921944 → 921944 | 9768426 → 9796440 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/10-named-membership.json).
