# Answer empty-object checks from lazy members

Source: tsrs `upstream/patches/0002 and notes/mem-lazy.md L11`.

Emptiness checks inspect lazy signatures, index infos and member names without instantiating every property. The original enumeration path remains for ineligible or already materialized types; no relation or instantiation budget is bypassed.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.995x | 145.9 → 143.8 | 68407 → 67912 | 458782 → 458767 |
| Compiler | checkers4 | 1.018x | 170.8 → 169.9 | 92660 → 92133 | 631045 → 630914 |
| Compiler-Unions | single | 0.992x | 166.6 → 170.8 | 64823 → 64228 | 483700 → 483506 |
| Compiler-Unions | checkers4 | 0.975x | 195.6 → 189.3 | 90799 → 90108 | 669620 → 670325 |
| xstate-main | single | 1.016x | 630.7 → 602.9 | 582976 → 582108 | 4610654 → 4613488 |
| xstate-main | checkers4 | 0.993x | 695.7 → 696.2 | 734928 → 733962 | 5969852 → 5970684 |
| webpack | single | 1.004x | 1177.3 → 1159.4 | 743886 → 740119 | 8352540 → 8355932 |
| webpack | checkers4 | 1.001x | 1170.6 → 1249.9 | 937356 → 932857 | 9779996 → 9807769 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/04-empty-object.json).
