# Compiler comparison

2026-10-11T22:42:45.069783+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.324 / 147.5 | 0.317 / 142.8 | 0.330 / 141.6 | 0.119 / 56.9 | 0.98x | 0.96x | yes | yes |
| Compiler | checkers4 | 0.116 / 180.3 | 0.113 / 171.7 | 0.114 / 163.7 | 0.054 / 79.6 | 1.01x | 0.99x | yes | yes |
| xstate-main | single | 1.770 / 663.1 | 1.656 / 600.4 | 1.612 / 582.1 | 0.442 / 192.1 | 1.10x | 1.03x | yes | yes |
| xstate-main | checkers4 | 0.485 / 817.9 | 0.435 / 668.3 | 0.443 / 660.8 | 0.145 / 255.6 | 1.09x | 0.98x | yes | yes |
