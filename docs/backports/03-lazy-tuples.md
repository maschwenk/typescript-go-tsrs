# Lazy tuple reference members

Source: tsrs `upstream/patches/0001`, prototype L1. Remove the tuple exclusion from lazy reference eligibility: tuple element access, length, optional/rest/named elements, readonly behavior and array methods still resolve through the same mapper and instantiated base types. Full enumeration uses the same eager table.

Fresh local real-project comparison, three interleaved repetitions at one and four checkers: XState is 3.1% / 3.7% faster than the previous step with 1.6% / 4.0% less peak RSS. TypeScript's two compiler projects are essentially neutral, with about 0.4–1.9% less RSS. Webpack is 0.7% / 5.8% faster but four-checker RSS rises 3.7% in this noisy three-run sample; no unconditional memory improvement is claimed.

The upstream checker/compiler suite passes, including the added tuple regression whose baseline comes from the original Go compiler. All eight real-project cells agree across stock, parent, fork and tsrs on complete diagnostics, exit status and file count. See [raw results](results/03-tuple.json).
