# Lazy keyof-mapped members

Source: tsrs `upstream/patches/base/0002`, originating from microsoft/TypeScript#64526. Keep the keyof-constrained mapped type's declaration and mapper; synthesize the requested member, preserving optional/read-only modifiers, signatures, index information and declaration order. Full enumeration retains the original eager path. Depends on lazy reference member tables.

Fresh local three-repetition comparison: the mapped-member stress case improves from roughly 0.82 to 0.39 s single-threaded and peak RSS from 429 to 223 MiB. Four-checker time improves roughly 0.44 to 0.20 s. The new compiler regression covers homomorphic mapping, modifiers, key remapping, index signatures, inheritance and discriminated intersections. No existing baseline was changed.

Upstream checker/compiler baselines pass. Stock, parent, fork and tsrs all agree on diagnostic text, exits and files in every measured cell. These are workload-specific measurements; [raw results](results/02-mapped.json) include the independent parent contribution.
