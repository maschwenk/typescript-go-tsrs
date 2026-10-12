# Finite GC budget for one-shot compilation

The remaining gap is partly GC work over Go's pointer-rich compiler graph. The completed Linux profile at four threads spends 17.9% of sampled CPU in GC and write barriers, while tsrs owns its checker allocations without a tracing collector. The full profile and rejected representation prototypes are documented in [the gap investigation](53-profile-vscode-gap.md).

This is a Go-specific bridge to reducing that overhead, not a port of Rust's allocator. Ordinary one-shot CLI compilation uses `GOGC=400` with a soft 8 GiB Go memory budget only when physical memory and Linux cgroup ancestor limits permit at least 12 GiB. Smaller or unknown memory budgets retain Go's defaults. Explicit `GOGC` or `GOMEMLIMIT` variables take precedence, including an explicitly empty variable. The policy never disables GC and does not bound process RSS.

The setting lives in the executable, outside the compiler API. Watch mode, project builds, the language server, API mode and response-file commands are excluded. Keys, symbols, diagnostic order and TypeScript instantiation budgets are unchanged. Setting `GOGC=100` retains the previous default behavior for a workload whose live graph exceeds the automatic budget.

Three interleaved ordinary-CLI native VS Code repetitions at four runtime threads and checkers match the complete diagnostic text, exit status and 10,465 files. Median wall time changes from 5.419s to 4.836s (10.8% reduction), CPU from 19.591s to 17.440s (11.0%), and peak RSS grows from 6276.5 to 6502.6 MiB (3.6%). [Raw repetitions and binary hashes](results/54-batch-gc-native.json). This is the same pinned workload as the README, with Go 1.27, PGO off, on the M5 Max machine described in the investigation.

The automatic implementation also passes three interleaved ordinary-CLI repetitions on one public four-CPU Linux runner. The exact reference is `965d47afc2438d1fbc1f6e57ebe22e17a29c0a6a`, with the same compiler code before this policy; both Go binaries use Go 1.27.2 and PGO off. All validation passes and timed repetitions match the complete diagnostic text, exit status 2, 371 errors and 10,427 checked files. The native corpus has 10,465 files; comparisons are made within each platform's matching dependency installation.

| compiler | median wall | median CPU | median peak RSS |
| --- | ---: | ---: | ---: |
| prior fork | 20.096s | 77.642s | 6304.4 MiB |
| automatic finite GC policy | 16.589s | 61.447s | 6651.4 MiB |
| pinned tsrs | 8.546s | 33.782s | 1962.3 MiB |

This is a 17.5% wall-time reduction and 20.9% CPU reduction for 5.5% more peak RSS. The fork/tsrs wall ratio moves from 2.35x to 1.94x. Absolute seconds must not be compared with other runners. The reference repetitions range from 18.303s to 20.217s and the candidate from 16.499s to 16.758s; these are three repetitions on one runner, not confidence bounds across machines. [Raw repetitions and binary provenance](results/54-batch-gc-linux.json), [completed run](https://github.com/maschwenk/typescript-go-tsrs/actions/runs/38210420494).

CPU sampling remains separate. The profiled automatic implementation performs five collections, with heap goals approximately five times the preceding live heap, confirming that the eligible-machine policy was active. [GC trace](results/54-batch-gc-linux-trace.txt). The manual profiling job accepts an exact `parent_sha` to repeat this comparison. The complete upstream baseline matrix, race tests, build/vet and harness tests passed for source `5b1aa042574578e241e3f33260da0f8e31710009` in [CI run 38210418398](https://github.com/maschwenk/typescript-go-tsrs/actions/runs/38210418398). Rebase and documentation changes preserve the tested compiler and harness trees.

The command package tests use the actual runtime controls to check budget/percent application, restoration in tests and preservation of explicit overrides. Linux filesystem fixtures verify both cgroup v1/v2 parent caps, physical-memory caps and conservative handling of malformed or unavailable limits.
