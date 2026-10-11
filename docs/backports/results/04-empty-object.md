# Compiler comparison

2026-10-11T21:41:24.620180+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.320 / 148.0 | 0.317 / 145.9 | 0.319 / 143.8 | 0.119 / 56.9 | 1.00x | 1.00x | yes | yes |
| Compiler | checkers4 | 0.117 / 178.8 | 0.115 / 170.8 | 0.113 / 169.9 | 0.054 / 79.2 | 1.03x | 1.02x | yes | yes |
| Compiler-Unions | single | 0.565 / 168.5 | 0.572 / 166.6 | 0.577 / 170.8 | 0.178 / 59.6 | 0.98x | 0.99x | yes | yes |
| Compiler-Unions | checkers4 | 0.217 / 201.3 | 0.220 / 195.6 | 0.226 / 189.3 | 0.083 / 84.4 | 0.96x | 0.98x | yes | yes |
| xstate-main | single | 1.882 / 689.2 | 1.734 / 630.7 | 1.707 / 602.9 | 0.457 / 192.1 | 1.10x | 1.02x | yes | yes |
| xstate-main | checkers4 | 0.479 / 815.2 | 0.450 / 695.7 | 0.453 / 696.2 | 0.143 / 254.5 | 1.06x | 0.99x | yes | yes |
| webpack | single | 3.562 / 1216.1 | 3.491 / 1177.3 | 3.478 / 1159.4 | 0.880 / 439.5 | 1.02x | 1.00x | yes | yes |
| webpack | checkers4 | 1.187 / 1338.7 | 1.157 / 1170.6 | 1.155 / 1249.9 | 0.265 / 522.3 | 1.03x | 1.00x | yes | yes |
