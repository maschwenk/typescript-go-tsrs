# TypeScript Go with tsrs strategies

Private experiment applying the transferable optimizations from [tsrs](https://github.com/maschwenk/tsrs) to Microsoft's Go TypeScript compiler. Upstream source, history and license are preserved; the common compiler revision is pinned for fair comparison.

- [Strategy audit and per-backport explanations](docs/backports/README.md)
- [Reproducible three-way benchmark](bench/README.md)
- [Upstream TypeScript contributor guide](CONTRIBUTING.md)

GitHub Actions runs the upstream Go tests and compiler baselines, concurrent programs, race checks and a comparison of stock Go, PR parent, this fork and pinned tsrs. Benchmark JSON, logs and Markdown are downloadable from the Actions run; use the manual Benchmark workflow for the real-project matrix. No packages are published automatically.

Build: `cd tsc && go build -pgo=off -o ../.bench/bin/fork ./cmd/tsc`.

This branch is an experiment. Measured runtime gains and memory costs are stated per workload; tsrs's published Rust gains are not attributed to Go without measurement.
