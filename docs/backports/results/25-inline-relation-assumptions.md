# Compiler comparison

2026-10-11T23:17:04.642381+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.336 / 148.4 | 0.330 / 148.1 | 0.313 / 143.5 | 0.122 / 56.9 | 1.07x | 1.06x | yes | yes |
| Compiler | checkers4 | 0.134 / 177.4 | 0.124 / 170.4 | 0.129 / 158.9 | 0.063 / 79.4 | 1.04x | 0.97x | yes | yes |
| Compiler-Unions | single | 0.676 / 164.9 | 0.606 / 153.2 | 0.607 / 167.2 | 0.189 / 59.6 | 1.11x | 1.00x | yes | yes |
| Compiler-Unions | checkers4 | 0.234 / 201.6 | 0.229 / 181.0 | 0.243 / 183.0 | 0.092 / 84.7 | 0.96x | 0.94x | yes | yes |
| xstate-main | single | 2.112 / 645.3 | 1.983 / 587.5 | 2.064 / 594.7 | 0.560 / 192.1 | 1.02x | 0.96x | yes | yes |
| xstate-main | checkers4 | 0.586 / 810.5 | 0.569 / 665.5 | 0.539 / 661.8 | 0.179 / 255.7 | 1.09x | 1.06x | yes | yes |
| webpack | single | 4.452 / 1206.8 | 4.341 / 1140.6 | 4.286 / 1125.6 | 1.121 / 439.5 | 1.04x | 1.01x | yes | yes |
| webpack | checkers4 | 1.297 / 1385.8 | 1.241 / 1225.6 | 1.255 / 1225.1 | 0.293 / 519.8 | 1.03x | 0.99x | yes | yes |
