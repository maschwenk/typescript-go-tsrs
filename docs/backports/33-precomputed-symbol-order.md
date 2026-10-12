# Compute symbol declaration sort keys once

Source: tsrs `crates/tsrs_checker/src/utilities.rs:sort_symbols_by_key; notes/perf-tsgo-can1357-speedups.md`.

Precompute each symbol first declaration file index and position once, then sort the same declared/file/position/name/ID order. Only complete indexed programs take the path; unknown files and tracing retain the original comparator. Symbol IDs are assigned only on actual ties, as before, and reusable Go scratch is cleared after each operation. A two-file test compares against the original comparator and verifies the missing-index fallback.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
