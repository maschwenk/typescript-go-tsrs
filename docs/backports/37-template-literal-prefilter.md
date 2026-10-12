# Prefilter template literals by fixed text

Source: tsrs `crates/tsrs_checker/src/checker_13.rs:remove_string_literals_matched_by_template_literals`.

Apply the existing template-inference length, prefix and suffix test before calling the full matcher. Impossible string-literal matches avoid inference setup. The predicate uses the same byte-string operations as the original Go inference path; string mappings and potentially matching templates retain the full matcher.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
