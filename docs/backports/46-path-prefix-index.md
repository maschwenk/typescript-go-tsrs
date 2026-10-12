# Index path patterns by their ordered wildcard prefixes

Source: tsrs `crates/tsrs_module/src/resolver.rs:ParsedPatterns/match_pattern_or_exact`.

Large path-mapping tables index patterns by literal prefix and search prefix lengths from longest to shortest. Keep exact matches first and original insertion order among equal prefixes and suffix ties. Small tables keep the original scan. A real ordered mapping test compares indexed answers with the original longest-prefix matcher, including Unicode, exact matches, invalid stars and suffix ties; module package tests passed.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
