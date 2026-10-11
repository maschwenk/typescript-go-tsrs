# Compiler comparison

2026-10-11T21:58:34.499995+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.323 / 147.5 | 0.318 / 148.5 | 0.321 / 147.2 | 0.120 / 56.9 | 1.01x | 0.99x | yes | yes |
| Compiler | checkers4 | 0.115 / 169.1 | 0.124 / 176.0 | 0.117 / 167.6 | 0.054 / 79.7 | 0.98x | 1.06x | yes | yes |
| Compiler-Unions | single | 0.576 / 169.8 | 0.575 / 165.6 | 0.577 / 171.1 | 0.180 / 59.6 | 1.00x | 1.00x | yes | yes |
| Compiler-Unions | checkers4 | 0.218 / 201.2 | 0.221 / 194.2 | 0.221 / 195.7 | 0.083 / 84.8 | 0.99x | 1.00x | yes | yes |
| xstate-main | single | 1.853 / 689.8 | 1.703 / 618.4 | 1.688 / 608.8 | 0.451 / 192.1 | 1.10x | 1.01x | yes | yes |
| xstate-main | checkers4 | 0.488 / 812.4 | 0.454 / 691.2 | 0.445 / 691.2 | 0.141 / 255.3 | 1.10x | 1.02x | yes | yes |
| webpack | single | 3.584 / 1215.7 | 3.481 / 1164.8 | 3.521 / 1179.8 | 0.883 / 439.5 | 1.02x | 0.99x | yes | yes |
| webpack | checkers4 | 1.184 / 1325.7 | 1.156 / 1225.6 | 1.154 / 1178.8 | 0.265 / 523.8 | 1.03x | 1.00x | yes | yes |
