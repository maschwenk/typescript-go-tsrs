# Reuse stable constituent flags for large union kind queries

Source: tsrs `crates/tsrs_checker/src/checker_13.rs:maybe_type_of_kind/constituent_flags`.

Skip impossible kind queries for primitive unions and retain the transitive flags of large unions and intersections in a bounded per-checker map. Unlike the Rust assumption, Go enum and boolean union normalization can add flags, so exclude every tree containing those kinds. Only complete programs retain settled type graphs. Real nested-type and enum tests and full shadow baselines compare answers with the original recursive walk.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
