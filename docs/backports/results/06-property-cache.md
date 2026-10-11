# Compiler comparison

2026-10-11T21:51:01.704460+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.329 / 149.8 | 0.328 / 148.7 | 0.326 / 147.7 | 0.120 / 56.9 | 1.01x | 1.01x | yes | yes |
| Compiler | checkers4 | 0.117 / 179.3 | 0.116 / 176.3 | 0.117 / 175.6 | 0.055 / 79.4 | 1.00x | 1.00x | yes | yes |
| Compiler-Unions | single | 0.564 / 164.5 | 0.566 / 165.2 | 0.581 / 169.6 | 0.178 / 59.6 | 0.97x | 0.97x | yes | yes |
| Compiler-Unions | checkers4 | 0.218 / 197.4 | 0.221 / 188.8 | 0.220 / 187.0 | 0.083 / 84.7 | 0.99x | 1.00x | yes | yes |
| xstate-main | single | 1.638 / 653.0 | 1.687 / 631.3 | 1.672 / 632.6 | 0.433 / 192.1 | 0.98x | 1.01x | yes | yes |
| xstate-main | checkers4 | 0.487 / 812.1 | 0.460 / 689.0 | 0.457 / 689.3 | 0.143 / 256.8 | 1.07x | 1.01x | yes | yes |
| webpack | single | 3.597 / 1238.0 | 3.510 / 1163.7 | 3.557 / 1176.4 | 0.868 / 439.5 | 1.01x | 0.99x | yes | yes |
| webpack | checkers4 | 1.249 / 1352.3 | 1.226 / 1264.1 | 1.231 / 1223.2 | 0.276 / 523.9 | 1.02x | 1.00x | yes | yes |
