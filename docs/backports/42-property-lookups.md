# Reuse property lookups within structural relation checks

Source: tsrs `crates/tsrs_checker/src/relater_1.rs:get_unmatched_property_recording; crates/tsrs_checker/src/relater_2.rs:properties_related_to_rest`.

Retain only source symbols already queried while checking unmatched properties, then reuse them at the corresponding target position in the relation loop. Go limits reuse to resolved plain objects, checks target-symbol identity and the member epoch, and keeps lazy member queries on their original path. Recursive checks lease independent scratch slices and clear retained pointers. Full shadow baselines recompute every reused lookup and check its symbol and checker effects.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
