# Pack simple relation keys without mixing cache domains

Source: tsrs `notes/perf-icount-hillclimb.md; notes/perf-tsgo-can1357-speedups.md; can1357/TypeScript@76097835cd`.

Two 31-bit type IDs and two intersection-state bits fit in a uint64. Give those simple pairs a separate exact map and retain the original generic constraint hash path, wider-ID fallback, identity symmetry and cache-size budget. Explicit tags prevent packed and hashed keys from colliding.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.002x | 145.4 → 147.7 | 67327 → 67327 | 446100 → 446128 |
| Compiler | checkers4 | 1.013x | 168.6 → 159.1 | 91011 → 91011 | 593523 → 593709 |
| Compiler-Unions | single | 1.000x | 158.2 → 149.5 | 63690 → 63690 | 472979 → 472990 |
| Compiler-Unions | checkers4 | 0.963x | 185.3 → 182.9 | 88982 → 88982 | 637416 → 637480 |
| xstate-main | single | 0.993x | 621.6 → 612.8 | 569112 → 569112 | 4379222 → 4379283 |
| xstate-main | checkers4 | 1.002x | 668.9 → 658.7 | 711882 → 711882 | 5706634 → 5710303 |
| webpack | single | 1.030x | 1159.5 → 1137.8 | 733733 → 733733 | 8096169 → 8096156 |
| webpack | checkers4 | 0.994x | 1231.2 → 1136.3 | 921944 → 921944 | 9446216 → 9434319 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/24-packed-relation-keys.json).
