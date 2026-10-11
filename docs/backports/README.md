# tsrs strategy audit

Source: `maschwenk/tsrs@e4c7d83a5907d1b8b2dc2f444cef40e5dc2a3751`. Compiler spec: `microsoft/TypeScript@b85298b6a81f772d080b0455de0ca9d744cd6fd6`. This is a private history-preserving copy: GitHub does not allow a public repository's fork to be private. The local `upstream` remote tracks Microsoft; `main` starts at the common compiler revision.

Each accepted strategy receives its own PR with source attribution, its invariant, upstream tests and measurements against the previous merged step. Existing compiler baselines are never rewritten to conceal a behavior change. The complete-source comparison with tsrs is in `bench/README.md`.

| strategy family | Go applicability | source |
| --- | --- | --- |
| Lazy class/interface and mapped members | Direct Go prototypes, first prerequisites | `upstream/patches/base`, `notes/lazy-members.md` |
| Lazy tuple members, empty-object and unmatched-property queries | Direct Go prototypes, depend on lazy reference tables | `upstream/patches/0001..0003` |
| Lazy union/intersection property caches | Direct Go prototype, independent allocation reduction | `upstream/patches/0004` |
| Conditional instantiation without throwaway composite mapper | Direct Go prototype, independent allocation reduction | `upstream/patches/0005` |
| Sorted union runs and named-union membership sets | Transferable algorithm; maintain exact constituent order | `notes/perf-tsgo-can1357-speedups.md`, `notes/perf-icount-hillclimb.md` |
| Avoid name literal types when no index signature exists | Transferable, skip pure work only | `notes/perf-icount-hillclimb.md` |
| Operation-scoped resolved-module and reference-redirect answers | Transferable immutable program data reuse | `notes/perf-checker-cpu4.md` |
| Weak-type and intersection-normalization memo bits | Transferable only with explicit invalidation and side-effect proof | `notes/perf-icount-hillclimb.md` |
| Flow memo, generic-call memo, union construction front cache | Complex context/budget-sensitive reuse; needs shadow modes and depth-limit tests | `upstream/flow-memo.md`, `notes/perf-flow-union-inference.md`, `notes/perf-heavy-files-infer-memo.md` |
| Worker count, import locality, heavy files and stealing | Locality already exists upstream; dynamic reassignment depends on fixing history-dependent diagnostics | `notes/perf-checker-stealing.md`, `notes/perf-order-independence.md`, `notes/open-history-dependence.md` |
| Lazy declaration-file member lists | Transferable if declared order, merged symbols and declaration/emit checking are preserved | `notes/mem-lazy-dts-members.md` |
| Packed layouts, dense links, scratch reuse | Go-specific translation and allocation profiles required | `notes/mem-layout.md`, `notes/mem-dense-link-tables.md`, `notes/mem-checker-scratch.md` |
| PGO and release training | Go has native PGO; train Go profiles on held-out pinned inputs | `notes/perf-pgo.md` |
| Rust Thin/Fat LTO, LLVM BOLT pipeline | No direct transfer to Go's compiler/linker; use native PGO | `notes/perf-build-level.md`, `notes/perf-pgo.md` |
| mimalloc, leak arenas, compressed pointers, madvise arena THP, manual leaf retirement | Rust memory-model techniques, cannot substitute for Go GC safely as a backport | `docs/PORTING.md`, `notes/mem-pointer-compression.md`, `notes/mem-free-leaf-files.md` |
| Borrow instead of clone, cold splits, Fx hashing, inline register-frame fast paths | Many address Rust-port artifacts; inspect Go assembly/allocations before adding machinery | `docs/RUST.md`, `notes/perf-checker-cpu2.md`, `notes/perf-checker-cpu3.md` |
| Persisted front end, process checkers, shared type graph, variance shortcuts, speculative budget skipping | Rejected or parked in tsrs; no presumed win to backport | `notes/perf-round2-followups.md`, `notes/perf-derived-variance.md` |
| Differential baselines/fuzz, per-PR comparisons, regression artifacts, negative-result notes | Transferable engineering process | `docs/DEBUGGING.md`, `bench/README.md`, `tools/fuzz/parser.py` |

The table is a disposition ledger, not a claim that every possible Go translation has shipped. Per-strategy notes record measured accept/reject decisions.
