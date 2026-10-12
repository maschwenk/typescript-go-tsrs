# Finite GC budget for one-shot compilation

The remaining gap is partly GC work over Go's pointer-rich compiler graph. The completed Linux profile at four threads spends 17.9% of sampled CPU in GC and write barriers, while tsrs owns its checker allocations without a tracing collector. The full profile and rejected representation prototypes are documented in [the gap investigation](53-profile-vscode-gap.md).

This is a Go-specific bridge to reducing that overhead, not a port of Rust's allocator. Ordinary one-shot CLI compilation uses `GOGC=400` with a soft 8 GiB Go memory budget only when physical memory and Linux cgroup ancestor limits permit at least 12 GiB. Smaller or unknown memory budgets retain Go's defaults. Explicit `GOGC` or `GOMEMLIMIT` variables take precedence, including an explicitly empty variable. The policy never disables GC and does not bound process RSS.

The setting lives in the executable, outside the compiler API. Watch mode, project builds, the language server, API mode and response-file commands are excluded. Keys, symbols, diagnostic order and TypeScript instantiation budgets are unchanged. Setting `GOGC=100` retains the previous default behavior for a workload whose live graph exceeds the automatic budget.

Three interleaved ordinary-CLI native VS Code repetitions at four runtime threads and checkers match the complete diagnostic text, exit status and 10,465 files. Median wall time changes from 5.419s to 4.836s (10.8% reduction), CPU from 19.591s to 17.440s (11.0%), and peak RSS grows from 6276.5 to 6502.6 MiB (3.6%). [Raw repetitions and binary hashes](results/54-batch-gc-native.json). This is the same pinned workload as the README, with Go 1.27, PGO off, on the M5 Max machine described in the investigation.

The earlier public Linux finite-GC probe saved 18.3% wall time with 4.7% more RSS in one timed execution. That is motivation, not a repeated A/B claim for this automatic implementation. The manual profiling job now accepts an exact `parent_sha` and compares that plain parent binary, the candidate and pinned tsrs in interleaved ordinary four-thread runs on one runner. CPU sampling remains separate. The automatic policy requires that comparison and upstream correctness CI before merging.

The command package tests use the actual runtime controls to check budget/percent application, restoration in tests and preservation of explicit overrides. Linux filesystem fixtures verify both cgroup v1/v2 parent caps, physical-memory caps and conservative handling of malformed or unavailable limits.
