# Build a separate Go PGO compiler from pinned VS Code profiles

Source strategy: tsrs `notes/perf-pgo.md`; Go translation uses CPU profiles and Go's native `-pgo` build flag instead of Rust LTO/BOLT.

`bench/profile.py` collected ten completed single-threaded profiles from pinned VS Code on the native Mac, verified stable full diagnostics/exits/file counts, merged them with pprof and built a separate optimized compiler. The PGO build passed the full checker and compiler baseline suite with no baseline changes. Raw profiles and source/binary/profile hashes stay in `.bench/pgo`; the manual workflow preserves them as artifacts.

Normal builds remain explicitly `-pgo=off`. Training and evaluation use the same VS Code workload, so this is workload specialization. No application speedup or held-out generalization is claimed from successful training alone. The independent native comparison records plain builds separately.
