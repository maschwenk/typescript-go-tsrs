# TypeScript Go with tsrs strategies

Keep each independent backport in its own PR. Run the upstream compiler baseline suite and relevant package tests before merging. Defer performance measurements until the implementation series is complete, then compare only pinned VS Code natively. Never accept a changed existing baseline to make a performance patch pass. Record source, validation, measured results when available and tradeoffs in docs/backports.

Pin both compiler revisions and benchmark inputs. Compare full diagnostic text and exit status, not just error counts. Show default settings and equal checker counts separately. Report build mode, runner, wall time, peak RSS and unavailable instruction counters honestly. Avoid unsafe runtime replacement, GC disabling and changes to the TypeScript instantiation limits.

Use GitHub Actions. Label PRs coder-task-generated and end their description with the current Session link.
