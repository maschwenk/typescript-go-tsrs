# Compiler comparison

2026-10-11T22:04:59.504292+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.322 / 152.2 | 0.320 / 147.9 | 0.322 / 144.1 | 0.120 / 56.9 | 1.00x | 1.00x | yes | yes |
| Compiler | checkers4 | 0.114 / 177.5 | 0.116 / 171.9 | 0.116 / 175.7 | 0.054 / 79.9 | 0.99x | 1.00x | yes | yes |
| Compiler-Unions | single | 0.584 / 171.8 | 0.571 / 164.8 | 0.576 / 170.2 | 0.181 / 59.6 | 1.01x | 0.99x | yes | yes |
| Compiler-Unions | checkers4 | 0.217 / 201.6 | 0.219 / 194.5 | 0.219 / 185.2 | 0.083 / 84.7 | 0.99x | 1.00x | yes | yes |
| xstate-main | single | 1.690 / 660.3 | 1.681 / 623.0 | 1.686 / 625.5 | 0.464 / 192.1 | 1.00x | 1.00x | yes | yes |
| xstate-main | checkers4 | 0.483 / 829.4 | 0.450 / 693.2 | 0.453 / 691.3 | 0.143 / 256.2 | 1.07x | 0.99x | yes | yes |
| webpack | single | 3.698 / 1234.0 | 3.479 / 1160.7 | 3.442 / 1166.6 | 0.895 / 439.5 | 1.07x | 1.01x | yes | yes |
| webpack | checkers4 | 1.284 / 1380.5 | 1.183 / 1269.7 | 1.178 / 1272.4 | 0.284 / 522.3 | 1.09x | 1.00x | yes | yes |
