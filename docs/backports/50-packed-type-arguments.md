# Pack short type argument keys without the byte builder

Source: tsrs `crates/tsrs_checker/src/checker_09.rs:packed_type_list_key`.

Alias-free instantiation and conditional keys of up to three 32-bit type IDs fit directly in the existing 128-bit cache key with separate family, length and flag fields. Go maps already hash keys, so do not import the Rust map-specific mixing step. Longer and aliased lists keep the original byte builder and hash. Packed inputs are injective; mixed packed/hash caches retain the original 128-bit probabilistic collision model. Tests exercise zero and maximum IDs, order, lengths, flags, families and original hashed fallbacks.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
