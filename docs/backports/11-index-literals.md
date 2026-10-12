# Avoid name literal types for objects without index signatures

Source: tsrs `notes/perf-icount-hillclimb.md; crates/tsrs_checker/src/checker_09.rs`.

Query index infos once and return immediately when none exist. Only a type with index signatures needs the property name literal type or late-bound symbol key. Keep index applicability and synthetic intersection handling on the same original path.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.001x | 145.8 → 147.1 | 67327 → 67327 | 459434 → 459170 |
| Compiler | checkers4 | 1.000x | 172.1 → 176.1 | 91011 → 91011 | 632089 → 631575 |
| Compiler-Unions | single | 0.992x | 163.2 → 169.3 | 63690 → 63690 | 484641 → 484529 |
| Compiler-Unions | checkers4 | 1.026x | 193.9 → 195.8 | 88982 → 88982 | 672496 → 671798 |
| xstate-main | single | 0.994x | 623.4 → 609.9 | 569112 → 569112 | 4508491 → 4508330 |
| xstate-main | checkers4 | 0.983x | 690.5 → 692.0 | 711882 → 711882 | 5856947 → 5850642 |
| webpack | single | 1.036x | 1174.7 → 1161.3 | 733733 → 733733 | 8343789 → 8342506 |
| webpack | checkers4 | 0.991x | 1253.8 → 1215.4 | 921944 → 921944 | 9756294 → 9774471 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/11-index-literals.json).
