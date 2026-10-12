# Reuse the current file emit module format

Source: tsrs `crates/tsrs_checker/src/checker_09.rs:emit_module_format_of_file`.

A checker repeatedly asks for the module format of its current source file. Cache one file and its format after the program is complete, retaining the original program lookup for incomplete programs and file changes. Program options and file module metadata remain immutable for this checker lifetime; language-service programs create their own checkers.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
