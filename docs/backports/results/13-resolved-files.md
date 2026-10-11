# Compiler comparison

2026-10-11T22:14:27.040669+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.322 / 149.5 | 0.339 / 147.6 | 0.343 / 147.7 | 0.127 / 56.9 | 0.94x | 0.99x | yes | yes |
| Compiler | checkers4 | 0.116 / 180.3 | 0.116 / 164.1 | 0.116 / 175.3 | 0.055 / 79.5 | 1.00x | 1.00x | yes | yes |
| xstate-main | single | 1.829 / 682.1 | 1.662 / 616.0 | 1.597 / 580.5 | 0.445 / 192.1 | 1.15x | 1.04x | yes | yes |
| xstate-main | checkers4 | 0.482 / 832.3 | 0.448 / 701.7 | 0.449 / 683.9 | 0.144 / 256.1 | 1.07x | 1.00x | yes | yes |
