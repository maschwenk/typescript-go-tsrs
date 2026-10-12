# tsrs strategy audit

Source: `maschwenk/tsrs@e4c7d83a5907d1b8b2dc2f444cef40e5dc2a3751`, sharing the compiler revision `microsoft/TypeScript@b85298b6a81f772d080b0455de0ca9d744cd6fd6`. This independent fork preserves upstream history. It began as a private history-preserving copy because GitHub public forks cannot be private.

The implementation series contains 52 separate steps, including 50 source ports, optional Go PGO tooling and an adaptive checker default. Each note identifies the Rust source, the Go invariant and correctness evidence. PR links below show current merge state. [Performance was measured after implementation](results/native-vscode/README.md), using pinned VS Code on the native Mac. Earlier per-port experiments remain historical evidence; later ports have no claimed individual application speedup. No compiler baseline is rewritten.

## Translation decisions

| Family | Go translation |
| --- | --- |
| Lazy class/interface, mapped and tuple members; empty and unmatched queries | Ports 01–05 retain the original eager path where resolution is required. |
| Union properties and conditional mapper churn | Ports 06–08 avoid throwaway copies and use exact mapper identities. |
| Sorted runs, galloping merge and named-union membership | Ports 09–10 preserve constituent order and the total type comparator. |
| Name literals, program/file lookups and modifier walks | Ports 11–13, 33–50 skip repeated pure work or reuse immutable answers within their operation. |
| Weak/intersection flags, union fronts and indexed access | Ports 14–15, 19, 27–28, 49 retain bounded answers with the original budget and side-effect rules. Enum/boolean normalization is excluded from the kind-flags memo. |
| Layouts, dense links, hot/cold splits and scratch reset | Ports 16–18, 21–22 use GC-managed records and typed accessors. Temporary objects stay collectible. |
| Fast paths, relation keys and stack assumptions | Ports 23–25, 36–39 and 50 preserve nil, alias, enum, freshness and identity distinctions. |
| Unaffected members, contextual queries, recursion, inference and flow | Ports 26, 29–32, 35 and 40–43 include off/shadow modes and conservative fallback guards. Cached property symbols still read their current types. Flow hits replay shared writes and preserve depth limits; inference hits preserve candidate order and budgets. |
| Go-native PGO | Step 51 provides a separate optional build using actual Go CPU profiles. VS Code training and evaluation are the same pinned workload; no held-out generalization claim. |
| Borrowing versus cloning, explicit Rust inline annotations, Fx/PackedMap hashes and allocator entry points | Go already borrows slices and uses its own compiler inliner and Swiss maps. Do not copy Rust-specific register-frame and custom-map machinery without Go evidence. Packed argument keys use Go's normal map hash. |
| Rust leak arenas, pointer compression, leaf retirement, mimalloc, madvise/THP, LLVM LTO/BOLT | Depend on Rust's allocator/compiler and arena ownership. They do not transfer as safe substitutions for the Go collector, pointer rules and linker. Go-managed arenas and native PGO are the corresponding experiments here. |
| Lazy declaration-file member parsing/binding | Not shipped in this series. Go exposes `NodeList.Nodes` and `Symbol.Members` directly. The Rust implementation requires forcing accessors, parser checkpoints, reentrant binding, shared-list publication and an audit of external references. Discarding these Go lists under today's API would lose observable AST/symbol data. This needs a separate frontend/API design; tsrs reports zero VS Code benefit in `notes/mem-lazy-dts-members.md`. |
| Checker defaults | Step 52 selects a bounded worker count from the available CPUs and checked-file workload for ordinary CLI compilation; explicit counts and existing API/build/watch behavior are retained. |
| Dynamic stealing/module affinity and shared checker graphs | Keep upstream Go assignment. Assignment order can change history-dependent diagnostics; extra stealing needs the order-independence fixes and dedicated concurrency validation, not only a copied scheduler. See tsrs `notes/open-history-dependence.md` and `notes/perf-order-independence.md`. |
| Type-print memo, persisted frontend, derived variance shortcuts and speculative budget skipping | Rejected or parked tsrs experiments are not presumed wins. Type-print memo is absent from the pinned implementation. Existing instantiation limits stay unchanged. |

These decisions distinguish shipped ports from architecture work and techniques specific to Rust. They do not claim that every conceivable Go compiler redesign is complete. Source notes include `notes/perf-icount-hillclimb.md`, `notes/perf-icount-round4.md`, `notes/perf-tsgo-can1357-speedups.md`, `notes/mem-layout.md` and `docs/PORTING.md`.

## Source PRs

| Port | Change and validation | PR |
| --- | --- | --- |
| 01 | [Lazy class and interface reference members](01-lazy-references.md) | [#2](https://github.com/maschwenk/typescript-go-tsrs/pull/2) |
| 02 | [Lazy keyof-mapped members](02-lazy-mapped.md) | [#3](https://github.com/maschwenk/typescript-go-tsrs/pull/3) |
| 03 | [Lazy tuple reference members](03-lazy-tuples.md) | [#4](https://github.com/maschwenk/typescript-go-tsrs/pull/4) |
| 04 | [Answer empty-object checks from lazy members](04-empty-object.md) | [#5](https://github.com/maschwenk/typescript-go-tsrs/pull/5) |
| 05 | [Find unmatched properties without materializing a lazy target](05-unmatched.md) | [#6](https://github.com/maschwenk/typescript-go-tsrs/pull/6) |
| 06 | [Avoid copying union properties into the augmented cache](06-property-cache.md) | [#7](https://github.com/maschwenk/typescript-go-tsrs/pull/7) |
| 07 | [Instantiate conditional types without a throwaway composite mapper](07-conditional-mapper.md) | [#8](https://github.com/maschwenk/typescript-go-tsrs/pull/8) |
| 08 | [Use exact type IDs in alias-free mapper caches](08-mapper-cache.md) | [#9](https://github.com/maschwenk/typescript-go-tsrs/pull/9) |
| 09 | [Build unions by merging existing sorted runs](09-union-runs.md) | [#10](https://github.com/maschwenk/typescript-go-tsrs/pull/10) |
| 10 | [Reuse one membership set for large named unions](10-named-membership.md) | [#11](https://github.com/maschwenk/typescript-go-tsrs/pull/11) |
| 11 | [Avoid name literal types for objects without index signatures](11-index-literals.md) | [#12](https://github.com/maschwenk/typescript-go-tsrs/pull/12) |
| 12 | [Skip redirect map and lock probes when no project references exist](12-no-redirect.md) | [#13](https://github.com/maschwenk/typescript-go-tsrs/pull/13) |
| 13 | [Reuse resolved-module file answers within each checker](13-resolved-files.md) | [#14](https://github.com/maschwenk/typescript-go-tsrs/pull/14) |
| 14 | [Memoize weak-type checks after member resolution](14-weak-memo.md) | [#15](https://github.com/maschwenk/typescript-go-tsrs/pull/15) |
| 15 | [Memoize intersection normalization predicates](15-intersection-memo.md) | [#16](https://github.com/maschwenk/typescript-go-tsrs/pull/16) |
| 16 | [Use managed Go arenas for retained reference and literal types](16-managed-arenas.md) | [#17](https://github.com/maschwenk/typescript-go-tsrs/pull/17) |
| 17 | [Replace sparse symbol-link pointers with exact value offsets](17-offset-links.md) | [#18](https://github.com/maschwenk/typescript-go-tsrs/pull/18) |
| 18 | [Keep rare value-symbol link fields in a lazy tail](18-value-link-tail.md) | [#19](https://github.com/maschwenk/typescript-go-tsrs/pull/19) |
| 19 | [Cache exact small alias-free union constructions](19-union-front-cache.md) | [#20](https://github.com/maschwenk/typescript-go-tsrs/pull/20) |
| 20 | [Scan path segments in eight-byte words](20-word-path-scan.md) | [#21](https://github.com/maschwenk/typescript-go-tsrs/pull/21) |
| 21 | [Presize inherited member tables without changing creation order](21-inherited-capacity.md) | [#22](https://github.com/maschwenk/typescript-go-tsrs/pull/22) |
| 22 | [Reset pooled checker state without copying whole records](22-pool-reset.md) | [#23](https://github.com/maschwenk/typescript-go-tsrs/pull/23) |
| 23 | [Inline the unchanged plain-type normalization and apparent paths](23-inline-type-paths.md) | [#36](https://github.com/maschwenk/typescript-go-tsrs/pull/36) |
| 24 | [Pack simple relation keys without mixing cache domains](24-packed-relation-keys.md) | [#37](https://github.com/maschwenk/typescript-go-tsrs/pull/37) |
| 25 | [Keep shallow recursive relation assumptions in an inline scan](25-inline-relation-assumptions.md) | [#39](https://github.com/maschwenk/typescript-go-tsrs/pull/39) |
| 26 | [Reuse settled unaffected members across lazy generic instances](26-lazy-unaffected-memo.md) | [#40](https://github.com/maschwenk/typescript-go-tsrs/pull/40) |
| 27 | [Reuse repeated large unions before flattening sorting and interning](27-long-union-front-cache.md) | [#41](https://github.com/maschwenk/typescript-go-tsrs/pull/41) |
| 28 | [Memoize pure repeated indexed access over union keys](28-indexed-access-union-memo.md) | [#42](https://github.com/maschwenk/typescript-go-tsrs/pull/42) |
| 29 | [Reuse pure contextual types during const-context checks](29-contextual-type-reuse.md) | [#43](https://github.com/maschwenk/typescript-go-tsrs/pull/43) |
| 30 | [Keep recursion identities on active relation stacks](30-kept-recursion-identities.md) | [#44](https://github.com/maschwenk/typescript-go-tsrs/pull/44) |
| 31 | [Memoize pure repeated inference walks](31-inference-walk-memo.md) | [#45](https://github.com/maschwenk/typescript-go-tsrs/pull/45) |
| 32 | [Memoize pure flow frames and replay shared-node history](32-flow-frame-memo.md) | [#46](https://github.com/maschwenk/typescript-go-tsrs/pull/46) |
| 33 | [Compute symbol declaration sort keys once](33-precomputed-symbol-order.md) | [#47](https://github.com/maschwenk/typescript-go-tsrs/pull/47) |
| 34 | [Avoid contextual index-name literals without index signatures](34-contextual-index-literals.md) | [#48](https://github.com/maschwenk/typescript-go-tsrs/pull/48) |
| 35 | [Reuse contextual union property positions](35-contextual-union-positions.md) | [#49](https://github.com/maschwenk/typescript-go-tsrs/pull/49) |
| 36 | [Reject distinct plain literals before relation dispatch](36-plain-literal-relations.md) | [#50](https://github.com/maschwenk/typescript-go-tsrs/pull/50) |
| 37 | [Prefilter template literals by fixed text](37-template-literal-prefilter.md) | [#51](https://github.com/maschwenk/typescript-go-tsrs/pull/51) |
| 38 | [Compare sibling declarations without source-file walks](38-sibling-node-order.md) | [#52](https://github.com/maschwenk/typescript-go-tsrs/pull/52) |
| 39 | [Split common type-order checks and plain string comparisons](39-type-order-fastpaths.md) | [#53](https://github.com/maschwenk/typescript-go-tsrs/pull/53) |
| 40 | [Cache declaration source files during type ordering](40-node-file-cache.md) | [#54](https://github.com/maschwenk/typescript-go-tsrs/pull/54) |
| 41 | [Reuse settled property lookups in large union discrimination](41-discriminant-properties.md) | [#55](https://github.com/maschwenk/typescript-go-tsrs/pull/55) |
| 42 | [Reuse property lookups within structural relation checks](42-property-lookups.md) | [#56](https://github.com/maschwenk/typescript-go-tsrs/pull/56) |
| 43 | [Reuse discriminator property lookups during subtype reduction](43-subtype-property-reuse.md) | [#57](https://github.com/maschwenk/typescript-go-tsrs/pull/57) |
| 44 | [Skip irrelevant class accessibility modifier walks](44-property-accessibility.md) | [#58](https://github.com/maschwenk/typescript-go-tsrs/pull/58) |
| 45 | [Reuse declaration modifiers during union property synthesis](45-property-modifiers.md) | [#59](https://github.com/maschwenk/typescript-go-tsrs/pull/59) |
| 46 | [Index path patterns by their ordered wildcard prefixes](46-path-prefix-index.md) | [#60](https://github.com/maschwenk/typescript-go-tsrs/pull/60) |
| 47 | [Reuse the current file emit module format](47-emit-module-format.md) | [#61](https://github.com/maschwenk/typescript-go-tsrs/pull/61) |
| 48 | [Carry relevant exclude patterns through directory walks](48-live-directory-excludes.md) | [#62](https://github.com/maschwenk/typescript-go-tsrs/pull/62) |
| 49 | [Reuse stable constituent flags for large union kind queries](49-constituent-kind-flags.md) | [#63](https://github.com/maschwenk/typescript-go-tsrs/pull/63) |
| 50 | [Pack short type argument keys without the byte builder](50-packed-type-arguments.md) | [#64](https://github.com/maschwenk/typescript-go-tsrs/pull/64) |
| 51 | [Train an optional native Go PGO build](51-native-go-pgo.md) | [#66](https://github.com/maschwenk/typescript-go-tsrs/pull/66) |
| 52 | [Choose a bounded checker default from the workload](52-adaptive-default-checkers.md) | [#67](https://github.com/maschwenk/typescript-go-tsrs/pull/67) |
