# Keep recursion identities on active relation stacks

Source: tsrs `crates/tsrs_checker/src/relater_1.rs:is_deeply_nested_on_stack; notes/perf-icount-round4.md`.

Cache a recursion identity at the first original lookup of each active relation entry. Preserve exact pointer identities, intersection constituent checks, nondecreasing type-ID counting and the existing 3/10/100 depth thresholds. Use managed Go pointers rather than packed Rust addresses. The full shadow baseline suite compared every answer and retained identity with the original walk.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
