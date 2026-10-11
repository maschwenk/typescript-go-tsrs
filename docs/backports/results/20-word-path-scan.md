# Compiler comparison

2026-10-11T22:40:12.671924+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.325 / 149.2 | 0.321 / 146.5 | 0.319 / 145.8 | 0.120 / 56.9 | 1.02x | 1.01x | yes | yes |
| Compiler | checkers4 | 0.118 / 178.6 | 0.114 / 171.6 | 0.114 / 171.7 | 0.055 / 79.5 | 1.03x | 1.00x | yes | yes |
| webpack | single | 3.494 / 1208.8 | 3.461 / 1162.5 | 3.417 / 1162.3 | 0.895 / 439.5 | 1.02x | 1.01x | yes | yes |
| webpack | checkers4 | 1.200 / 1381.5 | 1.149 / 1145.2 | 1.157 / 1207.1 | 0.264 / 523.4 | 1.04x | 0.99x | yes | yes |
