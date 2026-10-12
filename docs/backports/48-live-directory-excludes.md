# Carry relevant exclude patterns through directory walks

Source: tsrs `crates/tsrs_vfs/src/vfsmatch/vfsmatch.rs:may_match_below/live_excludes`.

Once a leading literal path component differs from a directory prefix, that exclude cannot match any entry below it. Carry only still-possible excludes into descendants, keeping wildcard and ancestor cases conservative. File and directory matchers share immutable exclude patterns and one per-directory selection. Case-sensitive and case-insensitive tests compare every selected answer with the original matcher, including Unicode and carried descendant selections; filesystem matcher package tests passed.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
