# Public repository preparation

Prepared 2026-10-12 for `maschwenk/typescript-go-tsrs`. The owner authorized publication and the repository is now public. All 52 implementation steps merged individually after correctness validation; publication did not mark pending CI successful.

## Completed preparation

- Standard GitHub-hosted Ubuntu 24.04 runners for CI, Benchmark and optional Go PGO; standard Ubuntu 24.04 ARM64 for optional CodSpeed. Private runner labels and the temporary runner image definition were removed. Both old runner registrations were deleted; no local runner process remains.
- Read-only test tokens, credentials not persisted by checkout, pinned action revisions and no `pull_request_target` trigger. Exact-source checks are reported by a separate job that never executes the tested source.
- Every main push compares this fork, the official nightly and tsrs on pinned VS Code and updates the README tables. Separate manual comparison, PGO and CodSpeed workflows remain available, with seven-day artifact retention. Implementation PRs run correctness checks. Native Mac results retain source, build, hardware and diagnostic evidence.
- Public README, contributor instructions, fork-specific issue/PR templates and a fork security-reporting policy. Upstream guidance is retained separately. The Apache-2.0 license and upstream notices are unchanged; tsrs strategy attribution is explicit.
- Local checkout prefixes were normalized in current published benchmark JSON and logs. Original diagnostic hashes remain available alongside hashes of the normalized public logs. All 90 native benchmark logs were parsed again and their full diagnostic equality, exits and file counts verified.

## Audit evidence

TruffleHog 3.94.0 scanned the fork changes across fetched branch histories after the common upstream pin, including merge-parent diffs (135 commits in the snapshot). It also scanned the current project docs/tooling and downloaded GitHub records. Credential verification was disabled; no credential candidates were detected. Inherited Microsoft history was already public and was preserved rather than rewritten.

The GitHub audit covered 67 issues/PRs, 12 issue and review comments, all 27 retained benchmark artifacts and 274 available Actions log archives. Three failed hosted runs had no jobs and returned no log archive; their metadata was checked. No private company material was found in those records. The earlier queued private-runner replay was canceled. The preparation checks were rerun successfully on public hosted runners; earlier failures from the private allowance remain historical records.

The repository has no Actions secrets, Actions variables or environments. Older commits retain earlier machine paths and public commit authorship; history and PR Session links were not rewritten. No private corpus or credentials were uploaded. Native results use public, pinned VS Code source.

## Publication

The repository is public. Workflow approval is required for all outside contributors. Private vulnerability reporting, GitHub secret scanning and push protection are enabled. No self-hosted runners are registered.

The hosted infrastructure PR passed both test configurations and merged. Compiler PRs were replayed against the reviewed hosted CI workflow with exact source SHAs and merged in source order. History-only conflict resolutions preserved the validated compiler, tests and dependency inputs and recorded both revisions in the PR. Failed checks caused by the former private allowance are retained as history; no artificial green checks were created.

GitHub documents [free standard hosted runner use](https://docs.github.com/en/actions/reference/runners/github-hosted-runners#standard-github-hosted-runners-for-public-repositories) and that [Actions history and logs become visible](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/setting-repository-visibility#making-a-repository-public) when a repository is made public.
