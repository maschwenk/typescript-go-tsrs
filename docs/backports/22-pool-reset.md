# Reset pooled checker state without copying whole records

Source: tsrs `notes/mem-recycle.md; notes/mem-layout.md; can1357/TypeScript@ddc71b880b`.

The checker already retains these states and buffers. Clear each transient field in place and retain the same buffers and empty sets, avoiding full-record copies through Go write barriers. Every original depth, budget and recursion counter is reset in the same way.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.019x | 146.4 → 142.6 | 67327 → 67327 | 446104 → 446100 |
| Compiler | checkers4 | 1.006x | 164.5 → 170.8 | 91011 → 91011 | 593888 → 593627 |
| Compiler-Unions | single | 0.995x | 162.0 → 162.7 | 63690 → 63690 | 472985 → 472990 |
| Compiler-Unions | checkers4 | 1.020x | 190.4 → 191.3 | 88982 → 88982 | 637271 → 637335 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/22-pool-reset.json).
