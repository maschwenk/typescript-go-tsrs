# Public repository preparation

Prepared 2026-10-12 for `maschwenk/typescript-go-tsrs`. The repository remains private until its owner changes visibility. Compiler ports continue to merge individually after exact-source validation; publishing the repository does not mark pending CI successful.

## Completed preparation

- Standard GitHub-hosted Ubuntu 24.04 runners for CI, Benchmark and optional Go PGO; standard Ubuntu 24.04 ARM64 for optional CodSpeed. Private runner labels and the temporary runner image definition were removed. Both old runner registrations were deleted; no local runner process remains.
- Read-only test tokens, credentials not persisted by checkout, pinned action revisions and no `pull_request_target` trigger. Exact-source checks are reported by a separate job that never executes the tested source.
- VS Code performance workflows are manual, with seven-day artifact retention. No recurring or per-PR performance matrix runs. Native Mac results retain source, build, hardware and diagnostic evidence.
- Public README, contributor instructions, fork-specific issue/PR templates and a fork security-reporting policy. Upstream guidance is retained separately. The Apache-2.0 license and upstream notices are unchanged; tsrs strategy attribution is explicit.
- Local checkout prefixes were normalized in current published benchmark JSON and logs. Original diagnostic hashes remain available alongside hashes of the normalized public logs. All 90 native benchmark logs were parsed again and their full diagnostic equality, exits and file counts verified.

## Audit evidence

TruffleHog 3.94.0 scanned the fork changes across fetched branch histories after the common upstream pin, including merge-parent diffs (135 commits in the snapshot). It also scanned the current project docs/tooling and downloaded GitHub records. Credential verification was disabled; no credential candidates were detected. Inherited Microsoft history was already public and was preserved rather than rewritten.

The GitHub audit covered 67 issues/PRs, 12 issue and review comments, all 27 retained benchmark artifacts and 274 available Actions log archives. Three failed hosted runs had no jobs and returned no log archive; their metadata was checked. No private company material was found in those records. The earlier queued private-runner replay was canceled. New preparation PRs may show the existing hosted-allowance failure while the repository remains private; they must be rerun after publication.

The repository has no Actions secrets, Actions variables or environments. Older commits retain earlier machine paths and public commit authorship; history and PR Session links were not rewritten. No private corpus or credentials were uploaded. Native results use public, pinned VS Code source.

## At publication

Change the repository visibility only after owner approval. Then require workflow approval for all outside contributors, enable private vulnerability reporting and confirm GitHub secret scanning/push protection availability. The approval setting cannot be applied while the repository is private.

Rerun hosted CI for the public-preparation PR, merge it, then replay pending compiler PRs against the reviewed hosted CI workflow with exact source SHAs. Continue merging source PRs in order. Failed checks caused by the former private allowance are retained as history; no artificial green checks are created.

GitHub documents [free standard hosted runner use](https://docs.github.com/en/actions/reference/runners/github-hosted-runners#standard-github-hosted-runners-for-public-repositories) and that [Actions history and logs become visible](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/setting-repository-visibility#making-a-repository-public) when a repository is made public.
