# Memoize pure flow frames and replay shared-node history

Source: tsrs `crates/tsrs_checker/src/flowmemo.rs; upstream/flow-memo.md`.

Memoize complete flow frames for resolved identifier and this references without assigning new symbol IDs. Retain every shared-flow write in its original order and each subtree height, so the original 2000-frame limit stays observable. Go guards are intentionally stricter: no active instantiation, resolution or loop; no reduce label, shared-flow read, depth or inline limit, position-sensitive generic narrowing, expression check, nested flow invocation, new identities, diagnostics, budgets, reliability or member-resolution changes. Cached frames with up to 64 shared writes and bounded retained memory replay the original history. Shadow validation and a real nested flow graph check identity, history and the near-limit fallback.

The full upstream checker and compiler baseline suite passed locally. No existing baseline was changed.

Performance measurements are deferred to the final stock Go / fork / tsrs comparison on the pinned VS Code repository. This change has no claimed application speedup yet.
