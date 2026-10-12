# Reuse settled property lookups in large union discrimination

Source: tsrs `crates/tsrs_checker/src/relater_1.rs:discriminant_property_types`.

Keep bounded per-name property answers for large unions with resolved plain object members. Go stores symbols rather than permanently copying their types, so contextual symbol type updates still use the original getter. Unknown, module, tuple, mapped and augmented/indexed cases use the original query. Member epochs invalidate entries, including within a discrimination loop. Full shadow baselines and an actual union test check lookup effects, symbol-type changes and removed members.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
