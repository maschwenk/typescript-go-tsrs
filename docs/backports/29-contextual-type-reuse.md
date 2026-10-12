# Reuse pure contextual types during const-context checks

Source: tsrs `crates/tsrs_checker/src/checker_07.rs:is_const_context_ex; notes/perf-icount-hillclimb.md`.

Reuse the current node contextual type when the first query and recursive const checks create no identities, consume no instantiation budget and run outside active inference or resolution. Keep the original repeated query for effectful contexts. Shadow validation first caught repeated generic queries that incremented instantiation counters; the guard excludes them.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
