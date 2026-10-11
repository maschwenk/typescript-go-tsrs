# Compiler comparison

2026-10-11T22:01:42.460191+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.326 / 150.1 | 0.328 / 147.1 | 0.325 / 147.1 | 0.123 / 56.9 | 1.00x | 1.01x | yes | yes |
| Compiler | checkers4 | 0.120 / 179.8 | 0.119 / 165.2 | 0.121 / 171.5 | 0.056 / 79.5 | 0.99x | 0.98x | yes | yes |
| Compiler-Unions | single | 0.595 / 172.0 | 0.596 / 172.0 | 0.615 / 165.7 | 0.186 / 59.6 | 0.97x | 0.97x | yes | yes |
| Compiler-Unions | checkers4 | 0.231 / 199.5 | 0.236 / 183.6 | 0.231 / 186.4 | 0.090 / 85.0 | 1.00x | 1.02x | yes | yes |
| xstate-main | single | 1.833 / 679.3 | 1.695 / 617.4 | 1.687 / 627.3 | 0.445 / 192.1 | 1.09x | 1.00x | yes | yes |
| xstate-main | checkers4 | 0.480 / 831.7 | 0.449 / 682.4 | 0.449 / 691.9 | 0.144 / 255.6 | 1.07x | 1.00x | yes | yes |
| webpack | single | 3.573 / 1236.9 | 3.481 / 1166.3 | 3.499 / 1184.9 | 0.891 / 439.5 | 1.02x | 0.99x | yes | yes |
| webpack | checkers4 | 1.232 / 1386.8 | 1.207 / 1177.6 | 1.199 / 1239.1 | 0.276 / 523.8 | 1.03x | 1.01x | yes | yes |
