# Compiler comparison

2026-10-11T22:45:55.752368+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.327 / 149.8 | 0.323 / 146.4 | 0.317 / 142.6 | 0.119 / 56.9 | 1.03x | 1.02x | yes | yes |
| Compiler | checkers4 | 0.117 / 180.6 | 0.114 / 164.5 | 0.113 / 170.8 | 0.055 / 79.4 | 1.03x | 1.01x | yes | yes |
| Compiler-Unions | single | 0.574 / 171.9 | 0.567 / 162.0 | 0.570 / 162.7 | 0.179 / 59.6 | 1.01x | 1.00x | yes | yes |
| Compiler-Unions | checkers4 | 0.225 / 201.8 | 0.227 / 190.4 | 0.223 / 191.3 | 0.087 / 84.8 | 1.01x | 1.02x | yes | yes |
