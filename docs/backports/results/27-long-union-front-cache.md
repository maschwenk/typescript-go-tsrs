# Compiler comparison

2026-10-11T23:34:17.460006+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.446 / 137.1 | 0.444 / 149.8 | 0.446 / 138.8 | 0.160 / 56.9 | 1.00x | 0.99x | yes | yes |
| Compiler | checkers4 | 0.178 / 179.9 | 0.154 / 170.3 | 0.165 / 165.1 | 0.093 / 79.4 | 1.08x | 0.94x | yes | yes |
| Compiler-Unions | single | 0.728 / 162.8 | 0.670 / 156.8 | 0.674 / 153.5 | 0.242 / 59.6 | 1.08x | 0.99x | yes | yes |
| Compiler-Unions | checkers4 | 0.309 / 197.9 | 0.299 / 183.3 | 0.291 / 180.2 | 0.118 / 84.2 | 1.06x | 1.03x | yes | yes |
| xstate-main | single | 2.354 / 674.7 | 2.146 / 569.5 | 2.011 / 573.7 | 0.616 / 192.1 | 1.17x | 1.07x | yes | yes |
| xstate-main | checkers4 | 0.722 / 814.8 | 0.795 / 660.5 | 0.632 / 659.3 | 0.214 / 256.8 | 1.14x | 1.26x | yes | yes |
| webpack | single | 3.925 / 1203.3 | 3.885 / 1142.5 | 3.704 / 1132.6 | 0.975 / 439.5 | 1.06x | 1.05x | yes | yes |
| webpack | checkers4 | 1.599 / 1373.5 | 1.487 / 1140.8 | 1.447 / 1214.5 | 0.330 / 524.5 | 1.11x | 1.03x | yes | yes |
