# Compiler comparison

2026-10-11T22:27:32.918833+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.320 / 154.0 | 0.318 / 150.9 | 0.317 / 150.9 | 0.121 / 56.9 | 1.01x | 1.00x | yes | yes |
| Compiler | checkers4 | 0.116 / 177.9 | 0.126 / 171.3 | 0.121 / 173.7 | 0.055 / 79.8 | 0.96x | 1.04x | yes | yes |
| Compiler-Unions | single | 0.603 / 161.0 | 0.697 / 167.2 | 0.656 / 158.5 | 0.193 / 59.6 | 0.92x | 1.06x | yes | yes |
| Compiler-Unions | checkers4 | 0.263 / 201.0 | 0.268 / 184.4 | 0.259 / 193.1 | 0.101 / 84.5 | 1.01x | 1.03x | yes | yes |
| xstate-main | single | 1.887 / 675.9 | 1.598 / 579.4 | 1.691 / 628.3 | 0.454 / 192.1 | 1.12x | 0.95x | yes | yes |
| xstate-main | checkers4 | 0.484 / 819.9 | 0.443 / 690.0 | 0.444 / 678.4 | 0.144 / 254.9 | 1.09x | 1.00x | yes | yes |
| webpack | single | 3.575 / 1197.1 | 3.424 / 1163.3 | 3.432 / 1148.7 | 0.894 / 439.5 | 1.04x | 1.00x | yes | yes |
| webpack | checkers4 | 1.207 / 1366.2 | 1.162 / 1183.2 | 1.187 / 1248.5 | 0.279 / 522.4 | 1.02x | 0.98x | yes | yes |
