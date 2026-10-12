# TypeScript Go with tsrs strategies

An experimental Go TypeScript compiler that brings transferable strategies from [tsrs](https://github.com/maschwenk/tsrs) back to [Microsoft TypeScript](https://github.com/microsoft/TypeScript). It preserves upstream history and the Apache-2.0 license. This fork is maintained independently of Microsoft.

All 52 individually documented implementation steps are merged. The [strategy audit](docs/backports/README.md) links each PR and explains the Go translation and strategies that require more architectural work. Every published benchmark identifies the exact compiler commit it measured.

- [VS Code comparison and raw results](docs/backports/results/native-vscode/README.md)
- [Build and benchmark instructions](bench/README.md)
- [GitHub Actions](docs/ci.md)
- [Contributing](CONTRIBUTING.md)

```sh
# Go automatically selects the version required by tsc/go.mod.
GOTOOLCHAIN=auto go -C tsc build -pgo=off -o ../.bench/bin/tsgo ./cmd/tsc
.bench/bin/tsgo --version
```

CI runs upstream compiler baselines, both test-program concurrency settings, cache shadow checks and checker/compiler race tests on standard GitHub-hosted Ubuntu runners. Every push to `main` benchmarks VS Code and updates the tables below. Implementation PRs run correctness checks. The native Mac comparison uses plain Go and Rust release builds with matching complete diagnostics. Runtime and memory results describe that workload and machine. No packages are published automatically.

<!-- vscode-benchmarks:start -->
## VS Code benchmarks

<!-- measured-commit: 51e09e1eaf323f8bb23fbee3221323f9c693acf9 -->
Measured [51e09e1eaf32](https://github.com/maschwenk/typescript-go-tsrs/commit/51e09e1eaf323f8bb23fbee3221323f9c693acf9) on 2026-10-12T04:26:37.344171+00:00. [Actions run](https://github.com/maschwenk/typescript-go-tsrs/actions/runs/38203593935) · [raw measurements](bench/results/51e09e1eaf323f8bb23fbee3221323f9c693acf9.json)

Linux-6.17.0-1022-azure-x86_64-with-glibc2.39 · x86_64 · 4 logical CPUs · AMD EPYC 9V45 96-Core Processor · MemTotal:       16373452 kB

VS Code `9cf0128b9822` (`src`), warm caches, one warmup and 5 interleaved repetitions per compiler and mode. Tables show medians. Hosted runner hardware can vary between commits.

### Wall time (seconds)

| mode | tsgo | tsgo nightly | tsrs |
| --- | ---: | ---: | ---: |
| default settings | 14.778 | 15.503 | 4.592 |
| single thread | 31.971 | 34.860 | 10.977 |
| 4 checkers | 13.433 | 15.374 | 4.344 |
| 8 checkers | 15.479 | 15.852 | 4.722 |

### Peak memory (MiB)

| mode | tsgo | tsgo nightly | tsrs |
| --- | ---: | ---: | ---: |
| default settings | 6321.2 | 6939.8 | 1978.6 |
| single thread | 5333.5 | 6291.8 | 1659.8 |
| 4 checkers | 6253.9 | 6630.7 | 1958.7 |
| 8 checkers | 6849.7 | 7136.5 | 2173.0 |

### Compiler versions and diagnostics

| compiler | version | files | errors | matches tsgo diagnostics, exit and files |
| --- | --- | ---: | ---: | --- |
| tsgo | Version 7.1.0-dev | 10427 | 371 | yes |
| tsgo nightly | Version 7.1.0-dev.20261011.1 | 10427 | 371 | yes |
| tsrs | Version 7.1.0-dev (tsrs 0.14.0-ts7.1.0-dev.20260929, microsoft/TypeScript@b85298b6a81f) | 10427 | 371 | yes |

`tsgo` is this fork. `tsgo nightly` is `@typescript/typescript-linux-x64@7.1.0-dev.20261011.1`, resolved once from `typescript@next`. `tsrs` is pinned to `e4c7d83a5907`.

fork: Go -pgo=off; tsrs: plain Rust release; nightly: official npm binary. Default settings use each compiler’s own thread policy; the other rows compare equal checker counts. Single thread also limits runtime pools to one.

All warmups and repetitions completed and produced stable full diagnostics within each compiler. Different compiler versions can produce different diagnostics; those rows are descriptive timings, with no speedup ratio or correctness claim. Compiler correctness is gated separately by CI.

[Method and reproduction](bench/README.md#automatic-readme-benchmarks) · [Native Mac comparison](docs/backports/results/native-vscode/README.md)
<!-- vscode-benchmarks:end -->

[LICENSE.txt](LICENSE.txt) and [NOTICE.txt](NOTICE.txt) retain upstream licensing and third-party notices. [NOTICE.tsrs.txt](NOTICE.tsrs.txt) credits the tsrs source of the translated strategies.
