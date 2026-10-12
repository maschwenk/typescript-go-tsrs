# GitHub Actions validation

CI runs all upstream package and compiler baselines with both test-program concurrency settings, all supported cache shadow modes, checker/compiler race checks, vet, a plain Go build and the benchmark harness tests. Existing compiler baselines must not be rewritten to make an optimization pass.

Performance workflows are manual and use VS Code only. They are deferred until source implementation is complete. The primary comparison runs native executables on the Mac with idle CI containers paused. Optional Linux Benchmark, Native Go PGO and CodSpeed workflows provide reproducible team runs; each records its own platform and build mode.

Hosted CI defaults to Ubuntu 24.04. The private repository's hosted allowance was exhausted during this task, so reviewed workflows were dispatched to two temporary Linux ARM64 containers on the Mac. They have no host source, credentials or Docker socket mounts. Self-hosted jobs use local compute; they do not spend GitHub-hosted runner minutes. Their registrations and containers are removed when the validation campaign finishes. Subsequent hosted runs require available allowance; alternatively the team can register its own runner using `.github/runners/Dockerfile` and label `tsgo-backport-local`.

For a stacked PR, dispatch the reviewed workflow ref with `source_sha` set to the exact full source commit. Checkout and tests use that source, while the workflow definition stays at the selected reviewed ref. The last step records the actual job conclusion on the tested commit and links its GitHub Actions run. A workflow SHA is distinct from the source SHA; both are retained. Failed, canceled, unexecuted or stale-source checks cannot authorize a merge. Ports merge in source order once their exact commits pass.

Dependabot updates stay disabled while compiler and benchmark controls are pinned. Reused runner builds prune missing worktree registrations, preserving existing checkouts. Self-hosted caches consume disk even when hosted cache actions are enabled; monitor free space and remove reproducible build caches when needed. Never remove corpus revisions, source changes or diagnostic evidence as cache cleanup.
