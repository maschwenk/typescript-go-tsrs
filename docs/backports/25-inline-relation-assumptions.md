# Keep shallow recursive relation assumptions in an inline scan

Source: tsrs `notes/perf-icount-hillclimb.md; can1357/TypeScript@57a6a9a62f`.

Use the existing maybe-key stack for its first eight entries and hash only deeper entries. Preserve identity and constraint keys, tail unwind, propagated reliability flags and every recursion and complexity budget.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.056x | 148.1 → 143.5 | 67327 → 67327 | 446122 → 446117 |
| Compiler | checkers4 | 0.965x | 170.4 → 158.9 | 91011 → 91011 | 593921 → 593777 |
| Compiler-Unions | single | 0.998x | 153.2 → 167.2 | 63690 → 63690 | 473006 → 472995 |
| Compiler-Unions | checkers4 | 0.939x | 181.0 → 183.0 | 88982 → 88982 | 637392 → 637388 |
| xstate-main | single | 0.961x | 587.5 → 594.7 | 569112 → 569112 | 4379303 → 4379294 |
| xstate-main | checkers4 | 1.056x | 665.5 → 661.8 | 711882 → 711882 | 5733797 → 5711190 |
| webpack | single | 1.013x | 1140.6 → 1125.6 | 733733 → 733733 | 8096201 → 8096148 |
| webpack | checkers4 | 0.989x | 1225.6 → 1225.1 | 921944 → 921944 | 9445049 → 9447090 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/25-inline-relation-assumptions.json).
