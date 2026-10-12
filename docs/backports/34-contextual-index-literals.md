# Avoid contextual index-name literals without index signatures

Source: tsrs `crates/tsrs_checker/src/checker_15.rs:get_type_from_index_infos_of_contextual_type; notes/perf-icount-hillclimb.md`.

Resolve contextual index signatures before interning a property-name literal. Types without index signatures now return directly; tuple numeric-rest handling and actual index-signature checks retain their original order.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
