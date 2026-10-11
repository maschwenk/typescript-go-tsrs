# Reuse resolved-module file answers within each checker

Source: tsrs `notes/perf-checker-cpu4.md section 2; crates/tsrs_checker/src/checker_08.rs`.

Cache the immutable completed program answer by its existing resolved filename, including negative answers. This avoids repeated path normalization and file/redirect lookups for repeated imports. Hosts that load files on demand, including the auto-import resolver, keep the original lookup. Auto-import package tests passed in addition to compiler baselines.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.988x | 147.6 → 147.7 | 67327 → 67327 | 459172 → 459170 |
| Compiler | checkers4 | 1.003x | 164.1 → 175.3 | 91011 → 91011 | 631470 → 631486 |
| xstate-main | single | 1.041x | 616.0 → 580.5 | 569112 → 569112 | 4508335 → 4500071 |
| xstate-main | checkers4 | 0.999x | 701.7 → 683.9 | 711882 → 711882 | 5871537 → 5846489 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/13-resolved-files.json).
