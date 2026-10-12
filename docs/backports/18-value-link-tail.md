# Keep rare value-symbol link fields in a lazy tail

Source: tsrs `notes/mem-layout.md: value links hot and cold split`.

Every value-symbol link used to reserve all rare fields. Keep resolved type, target and mapper in a 32-byte header and allocate the rare tail only when needed. Typed accessors preserve independent write, name, containing-type and checked values; empty writes allocate nothing.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.992x | 150.5 → 149.8 | 67327 → 67327 | 443108 → 454476 |
| Compiler | checkers4 | 1.011x | 173.3 → 171.6 | 91011 → 91011 | 588954 → 604352 |
| Compiler-Unions | single | 0.984x | 162.7 → 158.2 | 63690 → 63690 | 468930 → 480906 |
| Compiler-Unions | checkers4 | 0.989x | 193.6 → 192.2 | 88982 → 88982 | 630920 → 647576 |
| xstate-main | single | 0.962x | 594.7 → 601.8 | 569112 → 569112 | 4428839 → 4477220 |
| xstate-main | checkers4 | 0.991x | 670.2 → 668.4 | 711882 → 711882 | 5743832 → 5818078 |
| webpack | single | 0.973x | 1149.3 → 1162.5 | 733733 → 733733 | 8178773 → 8248148 |
| webpack | checkers4 | 1.002x | 1234.4 → 1155.7 | 921944 → 921944 | 9535415 → 9612836 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/18-value-link-tail.json).
