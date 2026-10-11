# Comparing stock Go, this fork and tsrs

The TypeScript base and tsrs source are pinned in `pins.json`; tsrs ports that exact TypeScript revision. All tools receive the same tsconfig and files. Run native executables directly, with `--noEmit --incremental false --pretty false --extendedDiagnostics`. This avoids Node launcher overhead. Warm every compiler, rotate the order in each repetition, report medians of wall time, CPU time and per-process maximum RSS (`wait4`). Logs and raw repetitions stay in the JSON artifact. Single mode also sets `GOMAXPROCS=1` and `RAYON_NUM_THREADS=1`; fixed checker counts and default settings are separate comparisons.

```sh
python3 bench/build.py
# Clone the tsrs revision from pins.json and build its pinned Rust toolchain:
cargo build --locked --release -p tsrs_cli
python3 bench/run.py --stock .bench/bin/stock --fork .bench/bin/fork \
  --tsrs /path/to/tsrs/target/release/tsrs --tsrs-source /path/to/tsrs \
  --projects synthetic,Compiler,Compiler-Unions,xstate-main,webpack --reps 7
```

`--work-dir` can point at an existing tsrs `bench/.work` directory. The real-project setup and config overlays come from the pinned tsrs harness, with the exact suite and project revisions in its `bench/projects.json`. Synthetic workloads require no dependency install. They cover wide generic and mapped members, overlapping named unions and conditional distribution, long flow walks, module chains and deliberate type errors. They are stress tests, not representative headline applications.

The Go gate requires every repeated stock, parent and fork run to finish and produce identical complete diagnostics, exit code and file count. Full diagnostic continuation text is compared. Types and symbol counts are recorded but may decrease under lazy resolution. tsrs equality is always shown; `--gate-tsrs` additionally makes its known discrepancies fail the command. An incomplete run never counts as a performance win.

PR Actions runs build all Go binaries with the same toolchain and `-pgo=off`, plus plain Rust release. Their parent column isolates the current PR's contribution. The manual workflow includes XState and webpack by default and accepts any pinned tsrs project. Artifacts contain binaries' hashes and Go build metadata. These results do not claim parity with tsrs's published PGO/BOLT builds. Linux `count.py out.json -- compiler ...` optionally records user-space instructions; absent permissions are reported as unavailable. Go counts fluctuate with GC and cannot use tsrs's deterministic 1% instruction threshold unchanged.

No wall-time threshold fails a PR on a shared hosted runner. Examine multiple paired runs before calling a small delta a win. Compare RSS alongside speed: wider checker pools duplicate caches and can increase memory. Compiler conformance baselines, types/symbols tests, package tests, concurrent test programs and race checks remain separate correctness gates.
