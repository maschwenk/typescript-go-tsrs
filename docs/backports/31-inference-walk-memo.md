# Memoize pure repeated inference walks

Source: tsrs `crates/tsrs_checker/src/infermemo.rs; notes/perf-heavy-files-infer-memo.md`.

Keep exact candidate and contravariant candidate order, priority, top-level flags and inferred-type cache clearing for repeated long walks. Only warmed targets are keyed, with at most eight inference infos and bounded retained memory. Go uses stricter guards than tsrs: no flow work or expression checks, no active instantiation or resolution, no creation, diagnostics, budget or reliability changes, and member-resolution epoch invalidation. Shadow baselines and a real 32-property inference test check outcomes, cache clearing, mutation isolation and invalidation.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
