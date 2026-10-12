# Reuse settled unaffected members across lazy generic instances

Source: tsrs `notes/perf-icount-round4.md; checker_09.rs lazy_unaffected_members`.

Keep the sorted immutable snapshot of resolved unaffected declared members for each completed-program generic target and mapper mode. Recheck only pending symbols and setter write types, retain old lazy-table snapshots, and recompute every hit in shadow mode.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.992x | 143.8 → 142.9 | 67327 → 67327 | 446122 → 445127 |
| Compiler | checkers4 | 1.016x | 169.9 → 164.5 | 91011 → 91011 | 593670 → 592197 |
| Compiler-Unions | single | 1.007x | 166.8 → 160.9 | 63690 → 63690 | 472981 → 472012 |
| Compiler-Unions | checkers4 | 0.998x | 181.1 → 178.4 | 88982 → 88982 | 637343 → 635743 |
| xstate-main | single | 0.954x | 562.5 → 613.0 | 569112 → 569112 | 4379243 → 4329957 |
| xstate-main | checkers4 | 1.049x | 667.7 → 663.6 | 711882 → 711882 | 5711285 → 5649452 |
| webpack | single | 0.908x | 1152.0 → 1132.9 | 733733 → 733733 | 8096112 → 8072131 |
| webpack | checkers4 | 1.002x | 1131.4 → 1227.6 | 921944 → 921944 | 9448475 → 9424446 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/26-lazy-unaffected-memo.json).

A 512-member settled-target microbenchmark compares the real original scan with the memo: median 30,297 ns and 18,800 bytes/10 allocations per scan versus 7.993 ns and zero allocations for a memo hit. This isolates the repeated work removed; it does not predict application-wide speedup. [Raw microbenchmark](results/26-lazy-unaffected-micro.txt).

A seven-repetition single-checker follow-up, interleaved while the Linux runner shared this machine, found xstate parent/fork 1.038x and webpack 0.984x. The broader application effect is small and mixed; this port does not establish a universal wall-time gain. [Raw follow-up](results/26-lazy-repeat.json).
