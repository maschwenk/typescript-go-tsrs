# Lazy class and interface reference members

Source: tsrs `upstream/patches/base/0001` and `notes/lazy-members.md`, originating from Max's microsoft/TypeScript#64475. A generic reference keeps its mapper, instantiated signatures/index information and base types; instantiate each property only when requested. Enumerating members still materializes the table. Private/protected members, callable and constructable types and object augmentation keep their existing behavior.

Fresh local Go 1.27 / macOS arm64, plain release, three interleaved repetitions against the preceding commit: wide generic reference single-threaded wall 0.377 → 0.078 s (4.84x), peak RSS 203 → 50 MiB; four checkers 0.155 → 0.042 s (3.71x), 209 → 55 MiB. The mapped-member workload stays unchanged because mapped tables still enumerate everything; that is the next independent backport. These are stress-test gains, not a claim about application-wide speedup.

The full upstream checker and compiler baseline suite passed, as did compiler/checker race tests and vet. Every measured stock/parent/fork/tsrs diagnostic, exit code and file count matched in 12 cells. Added compiler regression baselines come from the Go prototype; no existing baseline was changed. Raw repetitions: [results](results/01-lazy-reference.json).
