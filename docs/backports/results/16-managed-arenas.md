# Compiler comparison

2026-10-11T22:23:52.747172+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.325 / 149.8 | 0.327 / 147.5 | 0.324 / 147.0 | 0.121 / 56.9 | 1.00x | 1.01x | yes | yes |
| Compiler | checkers4 | 0.116 / 180.8 | 0.115 / 175.6 | 0.114 / 168.7 | 0.055 / 79.5 | 1.02x | 1.01x | yes | yes |
| Compiler-Unions | single | 0.572 / 168.6 | 0.577 / 174.4 | 0.572 / 166.0 | 0.179 / 59.6 | 1.00x | 1.01x | yes | yes |
| Compiler-Unions | checkers4 | 0.222 / 198.2 | 0.218 / 195.5 | 0.218 / 195.5 | 0.083 / 84.5 | 1.02x | 1.00x | yes | yes |
| xstate-main | single | 1.835 / 686.5 | 1.513 / 605.8 | 1.669 / 618.3 | 0.443 / 192.1 | 1.10x | 0.91x | yes | yes |
| xstate-main | checkers4 | 0.482 / 830.5 | 0.449 / 682.5 | 0.445 / 684.2 | 0.145 / 256.5 | 1.08x | 1.01x | yes | yes |
| webpack | single | 3.583 / 1233.9 | 3.350 / 1166.2 | 3.477 / 1163.1 | 0.885 / 439.5 | 1.03x | 0.96x | yes | yes |
| webpack | checkers4 | 1.192 / 1395.1 | 1.155 / 1185.8 | 1.160 / 1253.9 | 0.266 / 522.1 | 1.03x | 1.00x | yes | yes |
