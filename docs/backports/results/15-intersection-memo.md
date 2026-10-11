# Compiler comparison

2026-10-11T22:21:53.744186+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.319 / 149.8 | 0.317 / 145.0 | 0.320 / 149.7 | 0.120 / 56.9 | 1.00x | 0.99x | yes | yes |
| Compiler | checkers4 | 0.116 / 169.3 | 0.115 / 170.3 | 0.117 / 175.8 | 0.054 / 79.5 | 0.99x | 0.98x | yes | yes |
| xstate-main | single | 1.832 / 684.3 | 1.791 / 630.2 | 1.770 / 623.0 | 0.441 / 192.1 | 1.04x | 1.01x | yes | yes |
| xstate-main | checkers4 | 0.521 / 813.5 | 0.491 / 691.4 | 0.491 / 682.3 | 0.155 / 254.4 | 1.06x | 1.00x | yes | yes |
