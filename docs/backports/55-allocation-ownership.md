# Allocation ownership and pointer representation on VS Code

Two bounded experiments failed to produce a meaningful speed improvement: lazy resolved-member records saved 2.4% peak RSS with unchanged CPU, and checker-owned typed arenas reduced median wall time by 0.3% with similar RSS. Both compiler patches were reverted. Their results narrow the next useful direction to compact records with explicit file/checker ownership, rather than more allocation batching in the existing pointer graph.

The sampled live heap is approximately 4.79 GiB after completed final collections. Allocation-stack attribution assigns 50.1% to frontend work and 45.5% to checker work. Optimizing checker allocation alone leaves a large frontend graph. This investigation does not establish how much of the remaining wall-time gap a representation redesign would recover.

## Workload and evidence

The compiler base is `d9839fdac02711f0e0ad2565305df86437e758f8`. The reference binary was built at `3e75fadc11f0a2aca1718486b97898e2da5ebb86`; the compiler tree and benchmark harness code are identical, with only published benchmark results differing. Both candidates used Go 1.27.0, PGO off and `-trimpath`. The Rust source audit uses `maschwenk/tsrs@e4c7d83a5907d1b8b2dc2f444cef40e5dc2a3751`, sharing upstream TypeScript `b85298b6a81f772d080b0455de0ca9d744cd6fd6`.

Measurements ran natively on an Apple M5 Max, 18 logical CPUs, 128 GiB memory, macOS 26.6.2 arm64. All use VS Code `9cf0128b9822dcd8a533f08ef7312956a4390301`, four checkers and `GOMAXPROCS=4`. The repository's automatic finite CLI GC policy was active, with no environment override. No VM was used.

Each candidate has three interleaved ordinary CLI repetitions, ordered reference/candidate, candidate/reference, reference/candidate. Full diagnostic text and exit status match on every repetition, against separate statistics-enabled validations with 10,465 files, 371 errors and exit 2. Timing excludes profiles, statistics and allocation census instrumentation. The native dependency installation has more files than the Linux installation; these are within-machine A/B comparisons, not comparisons of seconds across runners.

[Provenance](results/allocation-ownership-20261012/provenance.json), [lazy-member repetitions](results/allocation-ownership-20261012/lazy-members.json), [arena repetitions](results/allocation-ownership-20261012/typed-arenas.json), [live-heap categories](results/allocation-ownership-20261012/heap-categories.json), [static layouts](results/allocation-ownership-20261012/layout.json).

## Where the graph comes from

The heap profile was collected in a separate statistics-enabled run with `--pprofDir`. Statistics force final collections, making this a sampled live-heap view rather than a cumulative allocation count. It estimates 5,147,933,456 live bytes. The profile uses Go's default sampling rate; it is not an exact object census.

| Allocation category | Live MiB | Share |
| --- | ---: | ---: |
| Frontend: AST, parser, scanner, binder | 2461.25 | 50.13% |
| Checker | 2232.41 | 45.47% |
| AST created during checking | 20.21 | 0.41% |
| Compiler/module work | 41.74 | 0.85% |
| Source/file contents | 135.56 | 2.76% |
| Other | 18.28 | 0.37% |

Classification uses the first semantic compiler function from the allocation leaf, skipping runtime, core allocator and collection helpers. A package name inside a generic type argument does not classify the function. Frontend allocations with a checker caller are separated; file-content allocations are matched first. Each sample is assigned once. These are allocation sites, not a dominator tree or proof of which reference keeps each object alive.

Large sites include identifier records (501.14 MiB), binder/checker symbol arenas (543.32 MiB combined), property-access records (182.35 MiB), call-expression records (125.63 MiB), node lists (104.60 MiB), and flow nodes (99.57 MiB). Checker `newObjectType` is the first semantic allocation frame for 332.08 MiB; its flat allocation attribution is 236.55 MiB. The difference includes allocations through callees. Symbol-link records account for another 224.79 MiB. These sites use the existing pointer representation even where allocations already come from typed arenas.

The [raw heap profile](results/allocation-ownership-20261012/reference-memprofile.pb.gz) and [pprof text](results/allocation-ownership-20261012/heap-inuse-space.txt) are preserved. Diagnostic text from the profile matches validation after removing only the leading relative source-directory prefix from filenames: profiling ran at the repository root, while the timing harness ran at `vscode/src`. Diagnostic continuation text is unchanged.

## Current ownership and representation

Go already embeds a type header in its concrete type record; `newType` takes `data.AsType()`. It does not allocate a separate header for every concrete type. The cost is the size and references in the combined record, including its `TypeData` interface and checker pointer.

| Go record on this arm64 build | Bytes |
| --- | ---: |
| `ast.Node` | 48 |
| `ast.Identifier`, including embedded node | 72 |
| `ast.Symbol` | 96 |
| `ast.NodeList`, excluding its backing list | 32 |
| `checker.Type` header | 56 |
| `checker.StructuredType`, including header | 160 |
| `checker.ObjectType`, including header | 184 |
| `checker.TypeReference`, including header | 216 |
| `checker.UnionType`, including header | 272 |
| `checker.Signature` | 128 |

These are static record sizes, not allocator size classes or totals including referenced objects. The layout JSON also counts pointer-shaped fields, including interface metadata and possibly nil references; those counts do not measure live heap edges or GC work.

`ast.Node` carries a parent pointer and payload interface. Identifiers add flow state and a Go string. Symbols carry declaration lists, member/export maps and several node/symbol references. `Type` points to its checker, symbol, alias and payload. `Program` retains the checker pool and source files; the pool retains its checker list and program. A retained type can therefore retain its checker and associated compiler state. This follows from source references; the heap profile does not quantify the impact of that retaining path.

The existing `core.Arena[T]` uses ordinary typed Go array allocations, growing requested chunks up to 256 elements. Escaped interior pointers keep their backing chunks alive. Go still traces pointer-bearing records in those chunks, and a live record can retain dead neighbors. AST records, symbols, signatures, references and literals already use this approach. Adding more such arenas changes allocation granularity without replacing tracing or compressing references.

Relevant Go sources: [AST header](../../tsc/internal/ast/ast.go), [symbol ownership recovery](../../tsc/internal/ast/symbol.go), [type layouts](../../tsc/internal/checker/types.go), [checker pool](../../tsc/internal/compiler/checkerpool.go), [typed arena](../../tsc/internal/core/arena.go). Types obtained from different checkers cannot be mixed; `Program.GetTypeCheckerForFile` documents that boundary.

At the pinned Rust source, [type layouts](https://github.com/maschwenk/tsrs/blob/e4c7d83a5907d1b8b2dc2f444cef40e5dc2a3751/crates/tsrs_checker/src/types.rs) combine a 24-byte tagged header with its payload. Structured types carry a four-byte resolved-member handle when pointer compression is enabled, and the lazily allocated member record is 32 bytes on 64-bit targets. Thin lists, packed count/index state, sparse base-constraint state and target-only instantiation caches contribute to the smaller representation. Copying only the lazy allocation decision leaves most of those differences intact.

Rust's [arena implementation](https://github.com/maschwenk/tsrs/blob/e4c7d83a5907d1b8b2dc2f444cef40e5dc2a3751/crates/tsrs_core/src/arena.rs) has explicit thread allocations, recycling, parser checkpoints and regions tied to file versions, checkers and programs. Reclamation requires proving no surviving reference enters the discarded storage. Its verification modes check escaped references. Ordinary Go typed arenas provide a different ownership mechanism.

## Rejected experiment: lazy structured members

Moved the five resolved-member fields into a per-type lazy record allocated on first write. Nil getters preserved existing zero values; type identity, creation order, member-resolution epoch and mutable state ownership were preserved. The base constraint and rare abstract-constructor field stayed in the header.

`StructuredType` shrank from 160 to 80 bytes, `ObjectType` from 184 to 104, `TypeReference` from 216 to 136 and `UnionType` from 272 to 192. The new member record was 88 bytes. A separate allocation census recorded 2,334,317 structured types and 1,219,411 member records: 52.2% needed the tail. These are cumulative allocations, not surviving objects. The nominal cumulative payload saving is about 75.8 MiB before allocator rounding, spare arena capacity and lifetimes; halving the header does not halve total compiler memory.

| Median | Reference | Lazy member record | Change |
| --- | ---: | ---: | ---: |
| Wall seconds | 4.6157 | 4.6382 | +0.49% |
| CPU seconds | 16.7361 | 16.7301 | -0.04% |
| Peak RSS MiB | 6537.0 | 6381.3 | -2.38% |

Focused checker tests and full VS Code diagnostic parity passed. The complete upstream baseline suite was not run for this rejected prototype. The small memory improvement and extra indirection did not provide the requested speed gain. [Unapplied prototype](results/allocation-ownership-20261012/lazy-members.patch), [census](results/allocation-ownership-20261012/census.json).

## Rejected experiment: checker-owned type arenas

Added `core.Arena[T]` allocation for object, union, intersection, type-parameter, indexed-access, conditional, index and mapped types. Addresses stayed stable, objects began zeroed, and IDs/creation order were unchanged. Reference/literal/symbol/signature arenas already existed.

| Median | Reference | Additional typed arenas | Change |
| --- | ---: | ---: | ---: |
| Wall seconds | 4.5349 | 4.5219 | -0.29% |
| CPU seconds | 16.3748 | 16.2467 | -0.78% |
| Peak RSS MiB | 6445.4 | 6453.2 | +0.12% |

Full VS Code diagnostic parity passed; the complete upstream baseline suite was not run for this rejected prototype. Three repetitions on one machine cannot establish a sub-percent speedup. The result does not justify shipping this as a large win. [Unapplied prototype](results/allocation-ownership-20261012/typed-arenas.patch). Its reference repetitions belong to a separate A/B campaign; do not compare them with the lazy-member campaign to attribute a change.

## Next substantial prototype

Start with file-owned identifier records, which account for roughly 501 MiB in the sampled heap. Put compact records in an owner table, use integer references for relationships within that owner, and retain stable adapters where the existing `*Node` API requires them. Text references need to preserve source slices, decoded escape text and synthesized names. Measure the combined cost of compact storage, adapters and lookups; another tag-only change to the current pointer layout has already regressed time in [the earlier investigation](53-profile-vscode-gap.md).

The broader design needs an explicit ownership contract:

1. Separate parsed file versions, binder-owned symbols, checker-local types/links and operation scratch. A transient checker symbol is not file-owned merely because it borrows a declaration.
2. Use owner-qualified references for cross-file relationships and keep checker-local type identities isolated. Preserve node/symbol/type IDs and allocation order.
3. Keep owners and externally referenced Go objects reachable through typed Go roots. Integer handles alone do not keep Go objects alive. Avoid a permanent global owner registry that retains every completed program.
4. Preserve API/LSP/watch lifetimes and incremental reuse. Reclaim an owner only after all adapters, caches and surviving programs release it; parser rollback must account for escaped references.
5. Migrate the public fields, accessors and generated AST construction together. Validate upstream diagnostic, type, symbol and emit baselines, concurrency/race behavior, incremental reuse and full VS Code parity before merging.

This is an architectural port with an unproven speed benefit. It is not achieved by hiding live Go pointers in `uintptr` or putting pointer-bearing Go records in untracked byte storage. Go's [GC guide](https://go.dev/doc/gc-guide) explains tracing through references; the [unsafe pointer rules](https://pkg.go.dev/unsafe#Pointer) explain why an integer address is insufficient to retain a Go object.

The [finite-GC result](54-batch-gc-budget.md) already reduced the measured Linux wall-time gap to approximately 1.94x. Its post-change profile attributes about 10% of CPU to GC/write barriers and about 22% to maps. Even eliminating remaining GC work alone would not close that gap. Compact ownership should be evaluated for record size, tracing, locality and access cost together, without assuming Rust's remaining advantage comes from one allocator switch.
