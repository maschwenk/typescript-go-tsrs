# Compare sibling declarations without source-file walks

Source: tsrs `crates/tsrs_checker/src/utilities.rs:compare_nodes`.

Nodes sharing a non-null parent share their source file, so compare their positions directly. Preserve nil ordering and retain the file-index path for source-file nodes and nodes with different parents.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
