# Cache exact small alias-free union constructions

Source: tsrs `crates/tsrs_checker/src/unioncache.rs; notes/perf-union-inference.md`.

Repeated small unions reuse an exact input-order and reduction key before flattening and sorting. Store only results that create nothing except the returned union, instantiate nothing, take no impure reduction and are not error types. Alias-bearing calls keep the original ID assignment path. The complete compiler suite and four real projects pass shadow mode; CI now recomputes cache hits to verify identity and counters.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.023x | 144.3 → 143.0 | 67327 → 67327 | 454477 → 446853 |
| Compiler | checkers4 | 0.999x | 158.8 → 163.2 | 91011 → 91011 | 604672 → 596170 |
| Compiler-Unions | single | 0.982x | 165.9 → 170.1 | 63690 → 63690 | 480899 → 473733 |
| Compiler-Unions | checkers4 | 1.015x | 191.5 → 185.6 | 88982 → 88982 | 647888 → 640161 |
| xstate-main | single | 0.994x | 613.9 → 609.9 | 569112 → 569112 | 4477248 → 4379329 |
| xstate-main | checkers4 | 0.995x | 665.2 → 656.3 | 711882 → 711882 | 5803871 → 5711802 |
| webpack | single | 0.993x | 1143.9 → 1147.7 | 733733 → 733733 | 8248159 → 8103251 |
| webpack | checkers4 | 0.991x | 1137.3 → 1144.7 | 921944 → 921944 | 9619773 → 9452275 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/19-union-front-cache.json).
