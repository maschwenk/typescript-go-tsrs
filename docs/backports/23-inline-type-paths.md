# Inline the unchanged plain-type normalization and apparent paths

Source: tsrs `notes/perf-checker-cpu2.md; notes/perf-icount-hillclimb.md`.

Plain object types already return themselves from these operations. Put those pure guards in small wrappers that Go can inline and keep all fresh-literal, tuple, reference, union, intersection, substitution and simplification work in the original workers. Reduction and mapped-type paths retain their order and limits.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.985x | 146.2 → 144.2 | 67327 → 67327 | 446102 → 446106 |
| Compiler | checkers4 | 1.018x | 171.2 → 171.6 | 91011 → 91011 | 593600 → 593699 |
| Compiler-Unions | single | 1.020x | 165.3 → 166.7 | 63690 → 63690 | 472975 → 472989 |
| Compiler-Unions | checkers4 | 1.018x | 190.5 → 189.7 | 88982 → 88982 | 637323 → 637293 |
| xstate-main | single | 1.018x | 604.8 → 606.8 | 569112 → 569112 | 4379179 → 4379256 |
| xstate-main | checkers4 | 1.012x | 657.8 → 668.7 | 711882 → 711882 | 5706119 → 5713681 |
| webpack | single | 0.966x | 1141.1 → 1161.0 | 733733 → 733733 | 8096161 → 8096164 |
| webpack | checkers4 | 0.998x | 1231.0 → 1229.6 | 921944 → 921944 | 9457789 → 9435079 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/23-inline-type-paths.json).
