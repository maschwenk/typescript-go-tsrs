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

## Bounded experiments

- Increasing the managed arena chunk cap from 256 to 1024 did not help: four-thread medians 7.430s before and 7.466s after; CPU 26.530s and 26.988s. Live reported memory grew roughly 170 MiB. The code was reverted. [Repetitions](results/profile-vscode-20261012/rejected-arena-probe.json).
- A finite `GOGC=400` / `GOMEMLIMIT=8GiB` ordinary-CLI probe saved 16.3% wall time with one thread, but only 6.6% with four, increasing four-thread peak RSS about 6.6%. It neither disables GC nor reproduces Rust's arena ownership. It is not adopted as a default on this evidence. [Repetitions](results/profile-vscode-20261012/finite-gc-probe.json).

## Reproduce the investigation

`bench/profile_vscode.py` measures ordinary CLI runs at one and four threads, validates complete diagnostics, and records Go CPU/allocation profiles and GC traces separately. Rust sampling uses native `samply` on macOS or `perf` on Linux. The manual `Benchmark` workflow's `profile=true` input records both CPU profiles and repeats the finite GC probe on a standard public GitHub runner. The Linux profile is required before selecting a default memory policy from the Darwin data.

The next compiler change should target measured allocation volume, pointer-rich live state, or actual checker imbalance. A larger shared-checker design is not presumed to be a win: tsrs's own `notes/perf-shared-checker.md` rejected it, and VS Code's shared duplication was comparatively small. Scheduler changes also require the history/order-independence fixes, not merely a copied work queue.
