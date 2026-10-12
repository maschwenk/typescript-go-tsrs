# Memoize intersection normalization predicates

Source: tsrs `crates/tsrs_checker/src/checker_14.rs: shouldNormalizeIntersection; notes/perf-icount-hillclimb.md`.

Repeated normalization asks the same question of an immutable intersection. Cache the pure predicate in object flags, and clear the memo when constructing another type so propagated flags cannot carry another intersection’s answer.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.992x | 145.0 → 149.7 | 67327 → 67327 | 459175 → 459167 |
| Compiler | checkers4 | 0.981x | 170.3 → 175.8 | 91011 → 91011 | 631725 → 632078 |
| xstate-main | single | 1.012x | 630.2 → 623.0 | 569112 → 569112 | 4500058 → 4500072 |
| xstate-main | checkers4 | 0.999x | 691.4 → 682.3 | 711882 → 711882 | 5845021 → 5858825 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/15-intersection-memo.json).
