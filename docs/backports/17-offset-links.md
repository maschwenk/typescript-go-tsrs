# Replace sparse symbol-link pointers with exact value offsets

Source: tsrs `notes/mem-dense-link-tables.md; notes/mem-layout.md; can1357/TypeScript@cf4a9c69cc`.

The existing value-symbol store already assigns symbol IDs. Preserve that order while replacing sparse pointer pages with uint32 offsets into stable Go-managed blocks. Empty entries remain absent and growth preserves every returned pointer; overflow is checked.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.004x | 150.9 → 150.9 | 67327 → 67327 | 443249 → 443103 |
| Compiler | checkers4 | 1.041x | 171.3 → 173.7 | 91011 → 91011 | 589325 → 589288 |
| Compiler-Unions | single | 1.062x | 167.2 → 158.5 | 63690 → 63690 | 469051 → 468918 |
| Compiler-Unions | checkers4 | 1.032x | 184.4 → 193.1 | 88982 → 88982 | 631111 → 630829 |
| xstate-main | single | 0.945x | 579.4 → 628.3 | 569112 → 569112 | 4429751 → 4428829 |
| xstate-main | checkers4 | 0.996x | 690.0 → 678.4 | 711882 → 711882 | 5741884 → 5740844 |
| webpack | single | 0.998x | 1163.3 → 1148.7 | 733733 → 733733 | 8180196 → 8178754 |
| webpack | checkers4 | 0.979x | 1183.2 → 1248.5 | 921944 → 921944 | 9532975 → 9529523 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/17-offset-links.json).
