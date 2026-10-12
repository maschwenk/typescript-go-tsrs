# Profile the remaining VS Code gap

This is a measurement correction and architecture investigation, not another claimed compiler speedup. Compiler source is unchanged from the completed 52-step series.

## What the original Linux result includes

[Run 38203593935](https://github.com/maschwenk/typescript-go-tsrs/actions/runs/38203593935) measures the complete fork on a four-CPU AMD EPYC 9V45 runner. The default medians are 14.778s wall / 56.537s CPU for this fork and 4.592s / 17.903s for tsrs. The fork's median reported compiler `Total time` is 11.343s. Go's `EmitAndReportStatistics` records that total before forcing two full GCs to report live memory. tsrs has no tracing collector. Consequently approximately 3.435s of Go's measured wall time is outside its reported compilation time. This is a real reporting cost, not type-checking work.

README repetitions now omit `--extendedDiagnostics`. A separate instrumented validation pass proves completed checking and records the file count. Every timed repetition must match that pass's full diagnostic text and exit status. The JSON and generated table distinguish the new method from historical results; changing methods is not presented as a compiler optimization.

## Native profiles

The first native comparison used the same pinned VS Code checkout, plain Go 1.27 with PGO off, and pinned tsrs 0.14.0 with plain Rust release optimization. At four checkers both runtime pools were capped at four. Three interleaved repetitions and all warmups matched the full diagnostic text, exit status, and 10,465-file count. The machine is an Apple M5 Max with 18 logical CPUs and 128 GiB RAM. Profiling runs are separate from timed repetitions.

| mode | fork wall / CPU | tsrs wall / CPU | fork / tsrs peak RSS |
| --- | ---: | ---: | ---: |
| one thread | 20.934s / 19.722s | 6.534s / 6.503s | 5514.8 / 1619.2 MiB |
| four runtime threads and checkers | 6.862s / 25.405s | 1.851s / 7.154s | 6384.8 / 1815.4 MiB |

These initial rows **include statistics reporting**. [Raw repetitions](results/profile-vscode-20261012/instrumented-comparison.json). Ordinary CLI timing and profiles use the new collector described below instead.

In the instrumented four-thread Go CPU profile, `runtime.madvise` consumes 38.44% of CPU samples, entirely through Darwin's `sysUsedOS` while allocating heap spans. This is page recommit/accounting, not the scavenger returning unused memory. `runtime.gcDrain` consumes 28.73% inclusively, including the two forced reporting collections. These categories are disjoint; do not add other inclusive GC rows to them. The forced collections' GC traces show about 4.59 CPU seconds combined. Darwin's recommit syscall cost is platform-specific and must not be used to explain the Linux gap.

The Rust sampler records 7.685 CPU seconds at four checkers. Its checker threads each use approximately 1.51 CPU seconds and finish at 1.58–1.59s; the profile is balanced. Roughly 78.6% of sampled CPU is semantic checking, 10.7% is source-file parsing, and 4.2% is binding. These are sampled categories, not a sum of every inclusive frame. Sleeping-thread samples are excluded by weighting frames with `threadCPUDelta` in microseconds. [Rust single-thread profile summary](results/profile-vscode-20261012/tsrs-single-cpu.json), [four-thread summary](results/profile-vscode-20261012/tsrs-four-cpu.json).

Go's sampled single-thread allocation profile totals 5737 MiB and about 29.15 million allocations. AST/symbol arena records dominate retained bytes; anonymous object types, object-literal checking, and inference/mapper temporaries dominate many individual allocations. The Go process peaks at roughly three times Rust's RSS. Compact representations and allocation ownership remain the major investigation; another small relation cache will not remove these costs.

## Ordinary compilation on the same machine

A clean collector run repeats the comparison without statistics reporting. Every run matches the separate validation pass's full diagnostics and exit status. CPU sampling is performed afterward, outside these timing medians.

| mode | fork wall / CPU | tsrs wall / CPU | fork / tsrs peak RSS | wall ratio |
| --- | ---: | ---: | ---: | ---: |
| one thread | 16.084s / 15.813s | 6.220s / 6.200s | 5508.4 / 1619.2 MiB | 2.586x |
| four runtime threads and checkers | 5.279s / 19.031s | 1.864s / 7.212s | 6123.8 / 1820.7 MiB | 2.832x |

[Raw repetitions and binary provenance](results/profile-vscode-20261012/ordinary-comparison.json). These are native Mac measurements; they do not establish the Linux improvement from correcting the method.

In the ordinary Go CPU samples, Darwin's `sysUsedOS` / `madvise` accounts for 51.68% at one thread and 51.94% at four. The gap remains dominated by heap-management costs on this platform even after excluding the reporting GCs. The Rust profile remains balanced, with approximately 79% semantic checking, 10% parsing and 4% binding. [Go four-thread self profile](results/profile-vscode-20261012/ordinary-go-checkers4-cpu-self.txt), [Rust four-thread CPU profile summary](results/profile-vscode-20261012/ordinary-tsrs-checkers4-cpu.json).

The ordinary Go allocation snapshots are also saved. They can lag the latest GC cycles: `ProfileSession.Stop` writes the allocation profile without a final collection. Use the instrumented allocation profile above for the completed allocation census rather than treating these delayed snapshots as the complete byte total.

The collector's finite GC probe changes four-thread ordinary compilation from 5.279s to 4.709s, with peak RSS growing from 6123.8 to 6486.9 MiB. This repeats the CPU reduction but wall gains vary between runs. It remains an experiment, not an adopted default.

## Bounded experiments

- Increasing the managed arena chunk cap from 256 to 1024 did not help: four-thread medians 7.430s before and 7.466s after; CPU 26.530s and 26.988s. Live reported memory grew roughly 170 MiB. The code was reverted. [Repetitions](results/profile-vscode-20261012/rejected-arena-probe.json).
- A finite `GOGC=400` / `GOMEMLIMIT=8GiB` ordinary-CLI probe saved 16.3% wall time with one thread, but only 6.6% with four, increasing four-thread peak RSS about 6.6%. It neither disables GC nor reproduces Rust's arena ownership. It is not adopted as a default on this evidence. [Repetitions](results/profile-vscode-20261012/finite-gc-probe.json).

## Compact AST payload prototypes: rejected

The pinned tsrs source uses a 24-byte node header and a payload tag instead of a stored interface. A Go prototype replaces the node's 16-byte `nodeData` interface with a tag in existing alignment padding, reducing its header from 48 to 32 bytes. Payloads remain in ordinary Go-managed typed memory: a checked cast uses the embedded header at offset zero, and generated compile-time assertions verify that invariant for every concrete payload. Parent pointers and 64-bit node IDs are preserved.

The original interface-reconstruction version passes AST/parser package tests and VS Code's full diagnostic parity gate. Three interleaved ordinary four-thread repetitions reduce median RSS from 6139.3 to 5569.8 MiB but regress wall time from 4.939s to 5.128s and CPU from 18.133s to 18.944s. [Raw repetitions](results/profile-vscode-20261012/rejected-compact-ast-probe.json).

A second version generates direct accessors and passes concrete header pointers into factories. It also preserves VS Code diagnostics, but does not reverse the regression: 5.057s / 18.570s CPU before, 5.150s / 18.755s after; median RSS 6086.3 -> 5644.1 MiB. [Raw repetitions](results/profile-vscode-20261012/rejected-compact-ast-static-probe.json).

Neither prototype is included in compiler source. Saving 7-9% RSS at a runtime cost does not satisfy the request for major speed wins. The remaining representation work needs to remove substantially more allocated bytes or GC-visible state, rather than simply replacing one interface field with generated dispatch.

## Completed public Linux profile

[Run 38208040052](https://github.com/maschwenk/typescript-go-tsrs/actions/runs/38208040052) succeeds with both compilers' raw CPU profiles, full diagnostic parity and 10,427 checked files. This is a four-CPU public Ubuntu runner. The recorded platform is `Linux-6.17.0-1022-azure-x86_64`; its CPU model was not captured. These absolute times must not be compared with the original run on another runner. One ordinary timed execution per variant is used here for diagnosis, not a repeated performance claim.

| mode | fork wall / CPU | tsrs wall / CPU | finite GC probe wall / CPU | fork / tsrs peak RSS |
| --- | ---: | ---: | ---: | ---: |
| one thread | 45.399s / 45.554s | 20.091s / 20.088s | 38.154s / 38.281s | 5694.6 / 1658.3 MiB |
| four runtime threads and checkers | 20.542s / 80.748s | 8.744s / 34.283s | 16.792s / 62.696s | 6356.5 / 1970.3 MiB |

[Raw results and binary provenance](results/profile-vscode-20261012/linux-result.json). The ordinary wall ratio is 2.26x at one thread and 2.35x at four. The finite four-thread probe reduces wall time 18.3% and CPU 22.4%, with peak RSS growing 4.7%. Its candidate gap is 1.92x. A repeated A/B comparison is required before adopting it.

The earlier Linux sample, recovered from the stalled exporter, classifies CPU samples without double counting: 23.9% map operations (including map allocation), 19.0% GC/write barriers, 4.0% other allocation/zeroing, 53.1% other work. Each sample belongs to the first matching category in that order with GC preceding map allocation; this is sampled cost attribution, not an arithmetic sum of overlapping inclusive frames. [CPU classification](results/profile-vscode-20261012/linux-incomplete-tsgo-cpu-categories.json). Darwin's page-accounting hotspot is absent. Reducing GC helps, but GC alone cannot explain or eliminate a 3x gap.

The Rust Linux sampling reports zero lost samples and exposes symbol-table lookup and flow analysis among the hot symbols, without one dominant individual algorithm. [Four-thread flat report](results/profile-vscode-20261012/linux-checkers4-tsrs-perf-report.txt). Raw mangled Rust v0 symbols are preserved; no unsupported demangling or source-frame percentages are inferred.

## Other rejected or separate probes

A compact pointer-identity link table was tested across 25 checker stores. It preserves actual Go pointer keys, stable value addresses and exact presence, and does not change observable symbol-ID assignment. Core identity-map differential tests, checker/compiler package tests and VS Code diagnostic parity pass. Three interleaved native four-thread runs save only 2.7% wall / 1.6% CPU while increasing median RSS 6.1%, so the prototype is not shipped. [Raw repetitions](results/profile-vscode-20261012/rejected-compact-links-probe.json).

A Go PGO build trained from the recovered Linux VS Code profile saves 7.4% native wall / 10.5% CPU and 3.4% RSS against the plain fork in three interleaved four-thread runs. This is the existing optional PGO build strategy, not a new checker backport or a generalization benchmark: VS Code is both training and evaluation workload. The README comparison continues to use PGO off. [Raw repetitions](results/profile-vscode-20261012/pgo-linux-profile-native-probe.json).

## Reproduce the investigation

`bench/profile_vscode.py` measures ordinary CLI runs at one and four threads, validates complete diagnostics, and records Go CPU/allocation profiles and GC traces separately. Rust sampling uses native `samply` on macOS or `perf` on Linux. The manual `Benchmark` workflow's `profile=true` input records both CPU profiles and repeats the finite GC probe on a standard public GitHub runner. The Linux profile is required before selecting a default memory policy from the Darwin data.

The measured finite GC policy is documented in [the follow-up](54-batch-gc-budget.md). Further compiler changes should target measured allocation volume, pointer-rich live state, or actual checker imbalance. A larger shared-checker design is not presumed to be a win: tsrs's own `notes/perf-shared-checker.md` rejected it, and VS Code's shared duplication was comparatively small. Scheduler changes also require the history/order-independence fixes, not merely a copied work queue.
