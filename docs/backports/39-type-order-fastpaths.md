# Split common type-order checks and plain string comparisons

Source: tsrs `crates/tsrs_checker/src/utilities.rs:compare_types; notes/perf-icount-hillclimb.md`.

Keep identity, nil, checker ownership and flag ordering ahead of the structural comparator. Alias-free plain string literals compare values and IDs directly; other equal-flag types retain the original comparator. Enum normalization and cross-checker rejection keep their original order. Go inlining and application speedup are not assumed.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
