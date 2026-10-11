# Memoize pure repeated indexed access over union keys

Source: tsrs `checker_13.rs get_indexed_access_type_by_union_memoized; unioncache.rs indexed_access_memo_mode`.

Reuse exact object/index/flag results, including negative answers, only outside contextual expression access, active instantiation and type resolution. Keep node/alias requests on the original path, require unchanged identity/budget/diagnostic-attempt counters, cap entries at 8192 and recompute every shadow hit.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.037x | 142.0 → 132.6 | 67327 → 67327 | 444995 → 444992 |
| Compiler | checkers4 | 1.076x | 169.4 → 170.6 | 91011 → 91011 | 592134 → 591993 |
| Compiler-Unions | single | 1.025x | 150.6 → 154.6 | 63690 → 63690 | 467940 → 467939 |
| Compiler-Unions | checkers4 | 0.995x | 180.8 → 180.2 | 88982 → 88982 | 631778 → 631699 |
| xstate-main | single | 1.004x | 572.6 → 581.1 | 569112 → 569112 | 4329814 → 4329806 |
| xstate-main | checkers4 | 0.989x | 662.9 → 659.0 | 711882 → 711882 | 5646941 → 5640386 |
| webpack | single | 1.004x | 1145.4 → 1130.4 | 733733 → 733733 | 8065088 → 8064543 |
| webpack | checkers4 | 0.985x | 1153.6 → 1191.9 | 921944 → 921944 | 9420140 → 9412715 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/28-indexed-access-union-memo.json).
