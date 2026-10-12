# Contributing

This repository evaluates tsrs strategies in the Go TypeScript compiler. Open issues and PRs here for behavior introduced by this fork. The [upstream contributor guide](docs/upstream-contributing.md) is retained for reference; upstream submissions follow Microsoft's own policies.

Keep each independent optimization in its own PR. Describe the repeated work or allocation it avoids, the corresponding pinned tsrs implementation, the Go invariants, and any memory or concurrency tradeoffs. Include a small reproducer or behavioral test when it covers a failure that existing compiler baselines miss. Do not rewrite existing baselines to make a performance patch pass.

Install the package dependencies and run correctness checks with the required Go toolchain:

```sh
npm ci --ignore-scripts --no-audit --no-fund
GOTOOLCHAIN=auto go -C tsc test -count=1 ./...
GOTOOLCHAIN=auto go -C tsc test -race ./internal/compiler/... ./internal/checker/...
python3 -m unittest discover -s bench -p 'test_*.py'
```

GitHub Actions also runs both test-program concurrency settings and supported cache shadow modes. Maintainers merge the implementation sequence in source order after the exact source revision passes. Public fork PRs use read-only tokens on disposable hosted runners.

Performance evaluation uses only the pinned VS Code workload at the end of the implementation series. Follow [bench/README.md](bench/README.md), retain raw results and full diagnostic equality, and report equal checker counts separately from defaults. State build mode, source revisions, machine, repetitions and peak RSS. Disclose AI assistance in your PR when used.
