# Skip irrelevant class accessibility modifier walks

Source: tsrs `crates/tsrs_checker/src/utilities.rs:get_declaration_non_public_flags_from_symbol`.

Object property compatibility only needs private and protected modifiers. Ordinary symbols outside a class cannot carry those accessibility modifiers, so bypass their declaration walk. Class members and synthetic symbols retain the original modifier lookup. Other modifier consumers keep their full flags.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
