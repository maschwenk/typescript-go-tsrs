# Use managed Go arenas for retained reference and literal types

Source: tsrs `notes/mem-layout.md; notes/perf-tsgo-can1357-speedups.md; can1357/TypeScript@7dd635e6ff`.

Type references and literal types are retained for the checker lifetime. Allocate these records in the existing GC-managed arena to reduce individual allocations while preserving pointers and IDs. Temporary anonymous objects remain individually collected.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.009x | 147.5 → 147.0 | 67327 → 67327 | 459168 → 443245 |
| Compiler | checkers4 | 1.012x | 175.6 → 168.7 | 91011 → 91011 | 631677 → 589571 |
| Compiler-Unions | single | 1.008x | 174.4 → 166.0 | 63690 → 63690 | 484529 → 469057 |
| Compiler-Unions | checkers4 | 1.000x | 195.5 → 195.5 | 88982 → 88982 | 671822 → 631369 |
| xstate-main | single | 0.906x | 605.8 → 618.3 | 569112 → 569112 | 4500074 → 4429747 |
| xstate-main | checkers4 | 1.009x | 682.5 → 684.2 | 711882 → 711882 | 5849801 → 5739725 |
| webpack | single | 0.964x | 1166.2 → 1163.1 | 733733 → 733733 | 8318590 → 8180233 |
| webpack | checkers4 | 0.996x | 1185.8 → 1253.9 | 921944 → 921944 | 9752178 → 9539298 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/16-managed-arenas.json).
