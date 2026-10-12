# Reject distinct plain literals before relation dispatch

Source: tsrs `crates/tsrs_checker/src/relater_1.rs:is_type_related_to; notes/perf-icount-hillclimb.md`.

After preserving fresh-literal regularization and pointer equality, distinct regular string, number and bigint literal types cannot match any relation. Skip relation dispatch for those exact flags. Enum literals, boolean literals, unions and all other flags retain the original path.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
