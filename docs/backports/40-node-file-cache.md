# Cache declaration source files during type ordering

Source: tsrs `crates/tsrs_checker/src/utilities.rs:source_file_and_index_of_node`.

Reuse source-file and program-index answers for published AST nodes in complete programs. Go retains managed pointers in a bounded 4096-entry map. Incomplete programs, synthesized nodes and missing file indexes use the original walk; direct parent changes refresh entries. The shadow suite checks every hit against the parent walk, and parsed-file tests cover file order, reparenting and missing indexes. This relies on published source trees retaining their ancestor structure within a checker.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
