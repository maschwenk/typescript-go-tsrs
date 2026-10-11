# Memoize weak-type checks after member resolution

Source: tsrs `notes/perf-icount-hillclimb.md; crates/tsrs_checker/src/relater_1.rs is_weak_type`.

Store the pure weak-type predicate in two spare object flag bits only after members are resolved. Clear the bits whenever members reset or are assigned, and when cloning a reference or padding an object literal. Lazy unresolved member tables still run the original query. New tests exercise optional/required/signature transitions and unresolved clones.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.986x | 147.7 → 146.9 | 67327 → 67327 | 459171 → 459171 |
| Compiler | checkers4 | 0.995x | 165.1 → 174.8 | 91011 → 91011 | 631912 → 631557 |
| Compiler-Unions | single | 1.013x | 171.1 → 165.2 | 63690 → 63690 | 484521 → 484529 |
| Compiler-Unions | checkers4 | 1.012x | 195.9 → 193.0 | 88982 → 88982 | 671815 → 671801 |
| xstate-main | single | 0.998x | 620.2 → 634.6 | 569112 → 569112 | 4500041 → 4500026 |
| xstate-main | checkers4 | 0.993x | 681.7 → 690.6 | 711882 → 711882 | 5847310 → 5851331 |
| webpack | single | 1.063x | 1183.9 → 1160.2 | 733733 → 733733 | 8318556 → 8318585 |
| webpack | checkers4 | 1.014x | 1253.1 → 1268.4 | 921944 → 921944 | 9733312 → 9737199 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/14-weak-memo.json).
