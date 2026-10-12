# Reuse declaration modifiers during union property synthesis

Source: tsrs `crates/tsrs_checker/src/checker_09.rs:create_union_or_intersection_property`.

Reuse the read modifier flags already computed while combining a property for readonly and ordinary write checks. Synthetic symbols and symbols with setters retain the original write-specific declaration search, preserving mixed accessor accessibility and readonly behavior.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
