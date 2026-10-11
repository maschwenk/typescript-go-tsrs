# Find unmatched properties without materializing a lazy target

Source: tsrs `upstream/patches/0003; notes/mem-lazy.md`.

Structural comparison can walk target names and instantiate only those needed to report an unmatched property. Preserve required/optional, private and discriminant behavior while avoiding full target materialization.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.004x | 147.1 → 147.8 | 67912 → 67327 | 458775 → 461117 |
| Compiler | checkers4 | 0.981x | 175.7 → 170.7 | 92133 → 91011 | 630923 → 634451 |
| Compiler-Unions | single | 0.998x | 169.8 → 169.6 | 64228 → 63690 | 483516 → 486216 |
| Compiler-Unions | checkers4 | 1.010x | 195.4 → 196.0 | 90108 → 88982 | 670330 → 674579 |
| xstate-main | single | 0.988x | 615.0 → 628.1 | 582108 → 569112 | 4613451 → 4651261 |
| xstate-main | checkers4 | 0.981x | 698.5 → 692.8 | 733962 → 711882 | 5971514 → 6029204 |
| webpack | single | 0.983x | 1150.5 → 1158.3 | 740119 → 733733 | 8355888 → 8382027 |
| webpack | checkers4 | 0.993x | 1185.2 → 1252.6 | 932857 → 921944 | 9780883 → 9841518 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/05-unmatched.json).
