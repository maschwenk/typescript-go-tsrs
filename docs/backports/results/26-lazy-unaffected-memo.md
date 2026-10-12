# Compiler comparison

2026-10-11T23:22:13.840598+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.358 / 153.4 | 0.314 / 143.8 | 0.317 / 142.9 | 0.122 / 56.9 | 1.13x | 0.99x | yes | yes |
| Compiler | checkers4 | 0.117 / 179.1 | 0.115 / 169.9 | 0.113 / 164.5 | 0.055 / 79.4 | 1.03x | 1.02x | yes | yes |
| Compiler-Unions | single | 0.579 / 167.0 | 0.561 / 166.8 | 0.557 / 160.9 | 0.182 / 59.6 | 1.04x | 1.01x | yes | yes |
| Compiler-Unions | checkers4 | 0.253 / 199.1 | 0.236 / 181.1 | 0.237 / 178.4 | 0.097 / 84.9 | 1.07x | 1.00x | yes | yes |
| xstate-main | single | 1.857 / 658.9 | 1.694 / 562.5 | 1.776 / 613.0 | 0.467 / 192.1 | 1.05x | 0.95x | yes | yes |
| xstate-main | checkers4 | 0.491 / 812.2 | 0.451 / 667.7 | 0.429 / 663.6 | 0.150 / 254.4 | 1.14x | 1.05x | yes | yes |
| webpack | single | 4.027 / 1212.6 | 3.665 / 1152.0 | 4.037 / 1132.9 | 1.112 / 439.5 | 1.00x | 0.91x | yes | yes |
| webpack | checkers4 | 1.732 / 1356.9 | 1.578 / 1131.4 | 1.574 / 1227.6 | 0.395 / 523.5 | 1.10x | 1.00x | yes | yes |
