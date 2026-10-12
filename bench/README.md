# Comparing stock Go, this fork and tsrs

The workload is VS Code `src`, pinned to `9cf0128b9822dcd8a533f08ef7312956a4390301` in the pinned tsrs project definitions. Compiler controls are in `pins.json`. Implementation PRs run correctness checks; the README comparison updates on every main push.

Published results identify the exact fork source commit. Check out that commit before building to reproduce a published row; `main` can be an earlier step while the source series merges.

## Automatic README benchmarks

`readme-benchmarks.yml` measures VS Code on a standard GitHub-hosted `ubuntu-24.04` runner after every push to `main`. It compares `tsgo` (this fork, plain Go with `-pgo=off`), `tsgo nightly` (the official native platform binary at the version resolved from `typescript@next`) and `tsrs` (the pinned plain Rust release build). This follows tsrs's main-push measurement and publication pattern; its own published chart uses a different machine and optimized distribution build.

Each run resolves the nightly once, installs that exact version and records its npm integrity, source commit when available, binary hash and Go build settings. It runs one warmup and five interleaved repetitions for all three compilers in default, single-threaded, four-checker and eight-checker modes. The README shows separate time and memory tables, versions, full-diagnostic parity and the machine. All repetitions must finish and produce stable diagnostics, exit status and file counts within each compiler and across thread counts. Nightly versions can disagree with the pinned compiler semantics; disagreement is shown and no speedup ratio is calculated. The original same-source stock/fork correctness gate below stays separate.

Every measured main commit gets `bench/results/<full-commit>.json`. Raw logs remain in the seven-day Actions artifact. Different commits do not share a concurrency group, so rapid merges cannot cancel or replace queued measurements. A separate publisher writes on top of the latest main, retries a racing push, preserves prose outside the README markers and keeps a newer source's table when an older run finishes later. Its bot commits carry `[skip ci]` and use `GITHUB_TOKEN` to prevent a results loop. Branch dispatches collect artifacts without writing main.

Hosted CPU models can vary, so inspect machine metadata before comparing commits. These tables are separate from the existing native Mac results. Optional PGO and CodSpeed workflows stay manual.

To reproduce after building this fork and the pinned tsrs:

```sh
python3 bench/nightly.py --out .bench/nightly
nightly=$(python3 -c 'import json; print(json.load(open(".bench/nightly/provenance.json"))["package"])')
python3 bench/main_benchmark.py --fork .bench/bin/fork \
  --nightly ".bench/nightly/node_modules/$nightly/lib/tsc" \
  --nightly-provenance .bench/nightly/provenance.json \
  --tsrs-source /path/to/tsrs --tsrs /path/to/tsrs/target/release/tsrs
```

Use the exact recorded npm version instead of `@next` to reproduce an older nightly. `publish_readme.py` can regenerate the tables from recorded JSON without rerunning the compilers.

## Native Mac comparison

Build and run the three native executables on the same Mac. Stop other builds and pause other CPU-heavy tasks before timing. Keep ordinary GC settings. Go stock and fork use the same toolchain and `-pgo=off`; tsrs uses plain Rust release.

```sh
python3 bench/build.py
# In the tsrs checkout at the commit in pins.json:
cargo +1.99.0 build --locked --release -p tsrs_cli
# In this repository:
python3 bench/run.py --stock .bench/bin/stock --fork .bench/bin/fork \
  --tsrs /path/to/tsrs/target/release/tsrs --tsrs-source /path/to/tsrs \
  --projects vscode --modes single,checkers4,checkers8,default --reps 5 \
  --label native-mac --out .bench/results/vscode.json
```

Preparation installs only the pinned VS Code corpus with `npm ci --ignore-scripts --no-audit --no-fund`. An external `--work-dir` is read-only and must already contain that exact checkout and its dependencies. Synthetic and other pinned project support remains available for manual diagnosis; the team workflows use VS Code only.

Every compiler receives the same project and `--noEmit --incremental false --pretty false --extendedDiagnostics`. Each mode warms all three compilers and interleaves measured repetitions, rotating and reversing their order. Single mode also sets both runtime pools to one thread; four and eight checkers are separate from each compiler's default. The JSON records all raw runs, warmups, wall and CPU times, peak process RSS, binary hashes, Go build metadata, source pins, machine details and relevant runtime environment. Mac RSS is recorded in bytes; Linux `wait4` values are converted from KiB.

Every stock and fork warmup and measured repetition must finish with identical complete diagnostics, exit status and file count. Panic, timeout, missing metrics and unstable warmups fail the gate. Type, symbol and allocation counters may fall because lazy work is skipped. tsrs equality is reported separately; `--gate-tsrs` also requires it. Ratios involving a compiler with different diagnostics are marked invalid. This avoids presenting skipped or incomplete checking as speed.

The report shows medians; raw runs expose noise. These measurements describe the pinned VS Code workload on the recorded machine, with warm filesystem caches. They do not establish a universal speedup or reproduce tsrs's published PGO/BOLT release numbers. Instruction counts are unavailable on this Mac; `count.py` is optional Linux tooling.

## Correctness and profiles

```sh
python3 -m unittest discover -s bench -p 'test_*.py'
python3 bench/differential.py --stock .bench/bin/stock --fork .bench/bin/fork \
  --seeds 64 --shadow --work-dir .bench/differential-shadow
```

The seeded differential tool compiles generated TypeScript modules with generics, mapped types, unions, flow narrowing and deliberate errors. It compares full diagnostic bytes, exits and every emitted JS/declaration file in single and four-checker modes. Failure fixtures and logs are retained. This complements upstream conformance and race checks; it is bounded test coverage.

CPU profiles and a separate native Go PGO build are described in `pgo.md`. No `default.pgo` changes ordinary builds.

## Optional CodSpeed

CodSpeed supports Go `testing.B` walltime benchmarks and executable walltime comparisons. See [Go support](https://codspeed.io/docs/benchmarks/go) and [executable support](https://codspeed.io/docs/benchmarks/cli-commands). Go does not have CodSpeed CPU simulation support.

The local walltime engine supports macOS, but the pinned executable harness distributes Linux binaries only. Therefore the native Mac result uses `run.py`; `codspeed-vscode.yml` is a separate, manual GitHub-hosted Ubuntu ARM64 workflow. See [CLI release](https://github.com/CodSpeedHQ/codspeed/releases/tag/v5.4.2) and [harness platform list](https://github.com/CodSpeedHQ/codspeed/blob/v5.4.2/crates/exec-harness/Cargo.toml#L31-L32).

The workflow pins CLI 5.4.2 and harness 1.4.1 with asset hashes, disables profiling and process isolation for reproducibility with the validated configuration, saves raw logs and actual exits, and gates the collected results without rerunning compilers. Its grouped repetitions differ from the native interleaved comparison. Harmless executable smoke commands validated the restricted Linux container path; a CodSpeed compiler benchmark was not run for the native comparison.

Upload defaults off. Upload requires an existing CodSpeed repository connection; this task does not install the app or configure a subscription. CodSpeed does not remove the need for compute or GitHub-hosted runner minutes. The optional workflow uses the standard GitHub-hosted `ubuntu-24.04-arm` runner. See [CodSpeed setup](https://codspeed.io/docs/integrations/ci/github-actions/configuration).
