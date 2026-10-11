# Skip redirect map and lock probes when no project references exist

Source: tsrs `notes/perf-checker-cpu4.md section 1; crates/tsrs_compiler/src/projectreferencefilemapper.rs`.

With no configured project references the builder cannot populate the redirect maps or symlink redirect cache. Return the original filename directly, avoiding the shared map lock and hash for every module-format question. Programs with references retain their complete original redirect path. The compiler package tests also passed.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 1.011x | 148.4 → 148.3 | 67327 → 67327 | 459169 → 459169 |
| Compiler | checkers4 | 1.023x | 175.8 → 174.6 | 91011 → 91011 | 631733 → 631769 |
| xstate-main | single | 0.992x | 616.2 → 630.9 | 569112 → 569112 | 4508338 → 4508335 |
| xstate-main | checkers4 | 1.017x | 679.8 → 692.2 | 711882 → 711882 | 5864802 → 5855291 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/12-no-redirect.json).
