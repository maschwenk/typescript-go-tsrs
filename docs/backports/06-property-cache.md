# Avoid copying union properties into the augmented cache

Source: tsrs `upstream/patches/0004; notes/mem-lazy.md`.

The ordinary union/intersection property cache is already the authoritative answer. Augmented lookups can consult it instead of duplicating the same property symbols in a second map, preserving absence and access semantics.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.006x | 148.7 → 147.7 | 67327 → 67327 | 461122 → 460258 |
| Compiler | checkers4 | 0.995x | 176.3 → 175.6 | 91011 → 91011 | 634405 → 633183 |
| Compiler-Unions | single | 0.974x | 165.2 → 169.6 | 63690 → 63690 | 486201 → 485059 |
| Compiler-Unions | checkers4 | 1.003x | 188.8 → 187.0 | 88982 → 88982 | 674486 → 673211 |
| xstate-main | single | 1.009x | 631.3 → 632.6 | 569112 → 569112 | 4651207 → 4635893 |
| xstate-main | checkers4 | 1.007x | 689.0 → 689.3 | 711882 → 711882 | 6031825 → 6020569 |
| webpack | single | 0.987x | 1163.7 → 1176.4 | 733733 → 733733 | 8382028 → 8362568 |
| webpack | checkers4 | 0.996x | 1264.1 → 1223.2 | 921944 → 921944 | 9837552 → 9804400 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/06-property-cache.json).
