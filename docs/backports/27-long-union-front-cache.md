# Reuse repeated large unions before flattening sorting and interning

Source: tsrs `notes/perf-icount-hillclimb.md; unioncache.rs LongKeys`.

Extend the union front cache beyond eight inputs with exact ordered type-ID and reduction keys. Take scratch buffers around recursion, keep alias and origin calls on the original path, require unchanged budgets/identities/impurity, cap retained key bytes at one MiB and entries at 4096, and shadow-check hits.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.995x | 149.8 → 138.8 | 67327 → 67327 | 445137 → 444996 |
| Compiler | checkers4 | 0.938x | 170.3 → 165.1 | 91011 → 91011 | 591875 → 591562 |
| Compiler-Unions | single | 0.994x | 156.8 → 153.5 | 63690 → 63690 | 472032 → 467956 |
| Compiler-Unions | checkers4 | 1.028x | 183.3 → 180.2 | 88982 → 88982 | 635757 → 631905 |
| xstate-main | single | 1.067x | 569.5 → 573.7 | 569112 → 569112 | 4330019 → 4329822 |
| xstate-main | checkers4 | 1.258x | 660.5 → 659.3 | 711882 → 711882 | 5631923 → 5637796 |
| webpack | single | 1.049x | 1142.5 → 1132.6 | 733733 → 733733 | 8072075 → 8065123 |
| webpack | checkers4 | 1.028x | 1140.8 → 1214.5 | 921944 → 921944 | 9417235 → 9418424 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/27-long-union-front-cache.json).

For an already interned 128-input literal union, the real uncached worker takes median 1,402 ns and 2,432 bytes/4 allocations; the long-key hit takes 72.51 ns and zero allocations (19.3x in this microbenchmark). [Raw measurements](results/27-long-union-micro.txt).
