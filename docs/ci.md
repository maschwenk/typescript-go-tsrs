# GitHub Actions validation

CI uses standard GitHub-hosted Ubuntu 24.04 runners. The optional CodSpeed workflow uses standard Ubuntu 24.04 ARM64. [GitHub documents free standard hosted runner use for public repositories](https://docs.github.com/en/actions/reference/runners/github-hosted-runners#standard-github-hosted-runners-for-public-repositories). Larger runners and optional external service subscriptions have separate billing. Benchmark artifacts have seven-day retention.

CI runs all upstream package and compiler baselines with both test-program concurrency settings, supported cache shadow modes, checker/compiler race tests, vet, a plain Go build and the benchmark harness tests. Existing compiler baselines must not be rewritten to make an optimization pass. The duplicated union-shadow YAML key from the implementation branch has been removed.

Tests receive a read-only token, and checkout does not persist credentials. Workflows use pinned action commits and `pull_request`, without `pull_request_target`. An exact-source replay has a separate reporting job with check-write permission. It never checks out or executes the tested source; it reads the actual matrix conclusions from the current run attempt. Ordinary PR jobs do not get check-write permission. Workflow approval is required for all outside contributors.

For a stacked PR, dispatch the reviewed CI workflow ref with `source_sha` set to the exact full source commit. Checkout and tests use that source, while the workflow definition stays at the reviewed ref. Reporting records the actual job conclusions on that commit and links the run, workflow commit and attempt. Failed, canceled, unexecuted or stale-source checks cannot authorize a merge. Ports merge in source order once their exact commits pass.

Failed baseline checks retain the generated compiler and language-service baselines as a seven-day artifact for that matrix configuration. Compare them with the checked-in reference files before changing code; do not accept a new baseline to clear a performance PR. Rerun the full workflow after investigating a suspected flake so both configurations are present in the current attempt.

README benchmarks run on every main push and use VS Code only. The workflow measures this fork, the official nightly and pinned tsrs, then updates time, memory and diagnostic tables in the README. Each commit has its own concurrency group; older runs archive their measurements without replacing a newer source's table. Publication runs in a separate job with content-write permission; measurement receives a read-only token. Bot results commits carry `[skip ci]` and use `GITHUB_TOKEN` to prevent loops. See [the method](../bench/README.md#automatic-readme-benchmarks).

The native Mac comparison remains separate. Manual Linux Benchmark, Native Go PGO and optional CodSpeed provide team runs, each recording its platform and build mode. They do not start for implementation PRs or on a schedule. CodSpeed upload defaults off and requires an existing repository connection.

The repository is public and hosted checks have passed. The temporary Linux runners used during the private validation campaign have been removed. No local or self-hosted runner is registered, and the active workflows do not accept private runner labels. Historical branches and logs retain the earlier validation provenance. Pending backports are rerun and merged in source order. Native Mac performance evidence remains distinct from hosted CI correctness evidence.

Dependabot updates stay disabled while compiler and benchmark controls are pinned. Compiler source, corpus pins and diagnostic evidence must be retained when cleaning reproducible caches.
