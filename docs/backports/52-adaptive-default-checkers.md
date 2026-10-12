# Size batch compiler default checkers for CPU and checked work

Source: tsrs `crates/tsrs_compiler/src/checkerpool.rs:default_checker_count`.

Size the ordinary batch CLI default pool by available Go scheduler parallelism and at least 32 checked files per checker. Keep four as the small-workload floor and cap at eight for this Go translation, based on native VS Code scaling and its higher per-checker memory cost. Explicit checker/single-thread flags take precedence. Incremental, build, watch, API and noCheck paths keep the existing defaults. A boundary table covers CPU caps and checked-file thresholds; full compiler and execute package tests passed.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
