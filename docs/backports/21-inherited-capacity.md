# Presize inherited member tables without changing creation order

Source: tsrs `notes/mem-layout.md; notes/perf-tsgo-can1357-speedups.md; can1357/TypeScript@7dd635e6ff`.

Use the settled declared table and the first base property list to size one owned copy before adding inheritance. Keep the original base-resolution order, value-member precedence and private-static exclusions. Later bases use the existing table; lazy-member forcing follows its original path.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.961x | 142.8 → 141.6 | 67327 → 67327 | 446847 → 446105 |
| Compiler | checkers4 | 0.992x | 171.7 → 163.7 | 91011 → 91011 | 596302 → 594050 |
| xstate-main | single | 1.028x | 600.4 → 582.1 | 569112 → 569112 | 4379358 → 4379206 |
| xstate-main | checkers4 | 0.980x | 668.3 → 660.8 | 711882 → 711882 | 5701983 → 5722593 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/21-inherited-capacity.json).
