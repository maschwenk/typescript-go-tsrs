# Compiler comparison

2026-10-11T22:31:07.377704+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.359 / 153.7 | 0.360 / 150.5 | 0.363 / 149.8 | 0.125 / 56.9 | 0.99x | 0.99x | yes | yes |
| Compiler | checkers4 | 0.118 / 178.5 | 0.115 / 173.3 | 0.114 / 171.6 | 0.055 / 79.7 | 1.03x | 1.01x | yes | yes |
| Compiler-Unions | single | 0.574 / 161.8 | 0.566 / 162.7 | 0.574 / 158.2 | 0.183 / 59.6 | 1.00x | 0.98x | yes | yes |
| Compiler-Unions | checkers4 | 0.218 / 201.1 | 0.216 / 193.6 | 0.219 / 192.2 | 0.083 / 84.8 | 1.00x | 0.99x | yes | yes |
| xstate-main | single | 1.788 / 656.4 | 1.605 / 594.7 | 1.669 / 601.8 | 0.445 / 192.1 | 1.07x | 0.96x | yes | yes |
| xstate-main | checkers4 | 0.480 / 820.7 | 0.441 / 670.2 | 0.445 / 668.4 | 0.144 / 255.1 | 1.08x | 0.99x | yes | yes |
| webpack | single | 3.564 / 1217.3 | 3.353 / 1149.3 | 3.447 / 1162.5 | 0.895 / 439.5 | 1.03x | 0.97x | yes | yes |
| webpack | checkers4 | 1.190 / 1369.6 | 1.166 / 1234.4 | 1.164 / 1155.7 | 0.265 / 521.5 | 1.02x | 1.00x | yes | yes |
