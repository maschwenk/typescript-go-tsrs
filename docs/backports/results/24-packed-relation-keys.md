# Compiler comparison

2026-10-11T23:05:02.726204+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.364 / 152.1 | 0.352 / 145.4 | 0.351 / 147.7 | 0.132 / 56.9 | 1.04x | 1.00x | yes | yes |
| Compiler | checkers4 | 0.130 / 179.3 | 0.118 / 168.6 | 0.116 / 159.1 | 0.057 / 79.7 | 1.12x | 1.01x | yes | yes |
| Compiler-Unions | single | 0.625 / 171.8 | 0.582 / 158.2 | 0.582 / 149.5 | 0.195 / 59.6 | 1.07x | 1.00x | yes | yes |
| Compiler-Unions | checkers4 | 0.229 / 201.2 | 0.219 / 185.3 | 0.228 / 182.9 | 0.089 / 85.1 | 1.00x | 0.96x | yes | yes |
| xstate-main | single | 1.831 / 679.7 | 1.635 / 621.6 | 1.646 / 612.8 | 0.448 / 192.1 | 1.11x | 0.99x | yes | yes |
| xstate-main | checkers4 | 0.485 / 821.4 | 0.438 / 668.9 | 0.437 / 658.7 | 0.145 / 255.5 | 1.11x | 1.00x | yes | yes |
| webpack | single | 3.415 / 1211.9 | 3.319 / 1159.5 | 3.224 / 1137.8 | 0.889 / 439.5 | 1.06x | 1.03x | yes | yes |
| webpack | checkers4 | 1.194 / 1384.5 | 1.132 / 1231.2 | 1.139 / 1136.3 | 0.268 / 523.5 | 1.05x | 0.99x | yes | yes |
