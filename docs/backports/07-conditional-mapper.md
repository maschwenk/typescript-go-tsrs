# Instantiate conditional types without a throwaway composite mapper

Source: tsrs `upstream/patches/0005; notes/mem-lazy.md`.

Map conditional type parameters directly through the existing mapper pair instead of allocating a temporary combined mapper and probing its cache. Preserve mapper precedence, alias keys and all instantiation limits.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.011x | 146.1 → 146.7 | 67327 → 67327 | 460261 → 460257 |
| Compiler | checkers4 | 1.007x | 165.5 → 174.1 | 91011 → 91011 | 633097 → 632997 |
| Compiler-Unions | single | 0.990x | 163.1 → 169.4 | 63690 → 63690 | 485059 → 485061 |
| Compiler-Unions | checkers4 | 1.004x | 195.4 → 195.5 | 88982 → 88982 | 673192 → 672941 |
| xstate-main | single | 1.084x | 609.0 → 593.3 | 569112 → 569112 | 4635877 → 4498200 |
| xstate-main | checkers4 | 1.018x | 687.4 → 690.6 | 711882 → 711882 | 6010241 → 5860457 |
| webpack | single | 0.982x | 1156.9 → 1173.7 | 733733 → 733733 | 8362594 → 8338199 |
| webpack | checkers4 | 0.977x | 1184.8 → 1177.2 | 921944 → 921944 | 9815330 → 9762932 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/07-conditional-mapper.json).
