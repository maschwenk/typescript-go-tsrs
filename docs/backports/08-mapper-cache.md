# Use exact type IDs in alias-free mapper caches

Source: tsrs `notes/perf-checker-cpu3.md; crates/tsrs_checker/src/mapper.rs; Go reference can1357/TypeScript a8c4fd4456`.

Skip key building and cache probes for a newly active mapper, allocate recursive cache maps only when storing, and key alias-free results directly by TypeId. Aliased keys remain in a separate domain. Keep the original instantiation budget checks before lookup, refresh the frame pointer after recursion, and clear both domains at the original reset sites.

Fresh local macOS arm64 / Go 1.27, plain-release binaries, three interleaved repetitions. Results against the previous step:

| project | mode | parent/fork wall | parent → fork peak MiB | parent → fork symbols | parent → fork allocations |
| --- | --- | ---: | ---: | ---: | ---: |
| Compiler | single | 0.992x | 148.5 → 147.2 | 67327 → 67327 | 460249 → 460256 |
| Compiler | checkers4 | 1.058x | 176.0 → 167.6 | 91011 → 91011 | 633094 → 633242 |
| Compiler-Unions | single | 0.997x | 165.6 → 171.1 | 63690 → 63690 | 485059 → 485050 |
| Compiler-Unions | checkers4 | 1.000x | 194.2 → 195.7 | 88982 → 88982 | 673243 → 673164 |
| xstate-main | single | 1.009x | 618.4 → 608.8 | 569112 → 569112 | 4498204 → 4498155 |
| xstate-main | checkers4 | 1.020x | 691.2 → 691.2 | 711882 → 711882 | 5874112 → 5858530 |
| webpack | single | 0.989x | 1164.8 → 1179.8 | 733733 → 733733 | 8338251 → 8338218 |
| webpack | checkers4 | 1.002x | 1225.6 → 1178.8 | 921944 → 921944 | 9761970 → 9752745 |

The full upstream checker and compiler baseline suite passed at this step. No existing baseline was changed. Every measured Go run matched stock diagnostics, exit status and file count. tsrs agreement is shown independently in the raw results; stress-test and shared-machine timings do not establish a universal speedup.

[Raw repetitions](results/08-mapper-cache.json).
