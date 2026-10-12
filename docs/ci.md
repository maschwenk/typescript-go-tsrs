# GitHub Actions validation

CI uses standard GitHub-hosted Ubuntu 24.04 runners. The optional CodSpeed workflow uses standard Ubuntu 24.04 ARM64. [GitHub documents free standard hosted runner use for public repositories](https://docs.github.com/en/actions/reference/runners/github-hosted-runners#standard-github-hosted-runners-for-public-repositories). Larger runners and optional external service subscriptions have separate billing. Benchmark artifacts have seven-day retention.

CI runs all upstream package and compiler baselines with both test-program concurrency settings, supported cache shadow modes, checker/compiler race tests, vet, a plain Go build and the benchmark harness tests. Existing compiler baselines must not be rewritten to make an optimization pass. The duplicated union-shadow YAML key from the implementation branch has been removed.

Tests receive a read-only token, and checkout does not persist credentials. Workflows use pinned action commits and `pull_request`, without `pull_request_target`. An exact-source replay has a separate reporting job with check-write permission. It never checks out or executes the tested source; it reads the actual matrix conclusions from the current run attempt. Ordinary PR jobs do not get check-write permission. Repository owners should require approval for all outside contributors' workflow runs after making the repository public.

For a stacked PR, dispatch the reviewed CI workflow ref with `source_sha` set to the exact full source commit. Checkout and tests use that source, while the workflow definition stays at the reviewed ref. Reporting records the actual job conclusions on that commit and links the run, workflow commit and attempt. Failed, canceled, unexecuted or stale-source checks cannot authorize a merge. Ports merge in source order once their exact commits pass.

Performance workflows are manual and use VS Code only. The primary comparison runs executables directly on the native Mac. Linux Benchmark, Native Go PGO and optional CodSpeed provide separate team runs, each recording its platform and build mode. They do not start for every implementation PR or on a schedule. CodSpeed upload defaults off and requires an existing repository connection.

The temporary Linux runners used during the private validation campaign have been removed. No local or self-hosted runner is registered, and the active workflows do not accept private runner labels. Historical branches and logs retain the earlier validation provenance. Once the repository is public, rerun the blocked CI jobs on the hosted workflow and finish the pending source merges. Native Mac performance evidence remains distinct from hosted CI correctness evidence.

Dependabot updates stay disabled while compiler and benchmark controls are pinned. Compiler source, corpus pins and diagnostic evidence must be retained when cleaning reproducible caches.
