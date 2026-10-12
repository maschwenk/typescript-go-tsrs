# Reuse contextual union property positions

Source: tsrs `crates/tsrs_checker/src/checker_15.rs:get_type_of_property_of_contextual_type; notes/perf-flow-union-inference.md`.

Keep bounded property-position lists for large contextual unions. Only settled objects with no own or global Object property, signatures, index signatures or tuple/mapped handling are skipped. Unknown and contextual constituents retain their original order and original queries. Member-resolution epochs invalidate cached positions. The shadow suite and a real union test cover skipped-query effects, mutation and Object augmentation.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
