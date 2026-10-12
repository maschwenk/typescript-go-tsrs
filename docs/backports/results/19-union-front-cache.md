# Compiler comparison

2026-10-11T22:35:34.881919+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.323 / 149.7 | 0.329 / 144.3 | 0.322 / 143.0 | 0.121 / 56.9 | 1.00x | 1.02x | yes | yes |
| Compiler | checkers4 | 0.122 / 175.1 | 0.118 / 158.8 | 0.118 / 163.2 | 0.057 / 79.4 | 1.04x | 1.00x | yes | yes |
| Compiler-Unions | single | 0.574 / 169.2 | 0.568 / 165.9 | 0.579 / 170.1 | 0.179 / 59.6 | 0.99x | 0.98x | yes | yes |
| Compiler-Unions | checkers4 | 0.218 / 201.4 | 0.220 / 191.5 | 0.216 / 185.6 | 0.082 / 84.7 | 1.01x | 1.02x | yes | yes |
| xstate-main | single | 1.713 / 659.7 | 1.663 / 613.9 | 1.673 / 609.9 | 0.452 / 192.1 | 1.02x | 0.99x | yes | yes |
| xstate-main | checkers4 | 0.492 / 820.6 | 0.446 / 665.2 | 0.448 / 656.3 | 0.147 / 255.8 | 1.10x | 1.00x | yes | yes |
| webpack | single | 3.601 / 1208.8 | 3.449 / 1143.9 | 3.473 / 1147.7 | 0.894 / 439.5 | 1.04x | 0.99x | yes | yes |
| webpack | checkers4 | 1.210 / 1399.4 | 1.152 / 1137.3 | 1.162 / 1144.7 | 0.267 / 522.6 | 1.04x | 0.99x | yes | yes |
