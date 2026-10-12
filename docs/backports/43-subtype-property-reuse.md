# Reuse discriminator property lookups during subtype reduction

Source: tsrs `crates/tsrs_checker/src/checker_13.rs:remove_subtypes`.

Keep bounded operation-scoped target property answers for the current discriminator name and member epoch. Go only retains settled property symbols and still reads their types normally, preserving type resolution and contextual symbol updates. A name or member change clears the operation cache; unknown cases keep the original getter. Full shadow baselines check skipped lookups and checker effects.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
