# Compiler comparison

2026-10-11T21:47:59.040272+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.320 / 150.2 | 0.319 / 147.1 | 0.318 / 147.8 | 0.120 / 56.9 | 1.01x | 1.00x | yes | yes |
| Compiler | checkers4 | 0.115 / 179.8 | 0.113 / 175.7 | 0.116 / 170.7 | 0.054 / 79.5 | 0.99x | 0.98x | yes | yes |
| Compiler-Unions | single | 0.578 / 177.3 | 0.570 / 169.8 | 0.570 / 169.6 | 0.179 / 59.6 | 1.01x | 1.00x | yes | yes |
| Compiler-Unions | checkers4 | 0.217 / 202.0 | 0.222 / 195.4 | 0.220 / 196.0 | 0.083 / 84.7 | 0.99x | 1.01x | yes | yes |
| xstate-main | single | 1.825 / 691.1 | 1.674 / 615.0 | 1.694 / 628.1 | 0.444 / 192.1 | 1.08x | 0.99x | yes | yes |
| xstate-main | checkers4 | 0.493 / 830.2 | 0.453 / 698.5 | 0.462 / 692.8 | 0.143 / 256.6 | 1.07x | 0.98x | yes | yes |
| webpack | single | 3.572 / 1239.0 | 3.398 / 1150.5 | 3.457 / 1158.3 | 0.902 / 439.5 | 1.03x | 0.98x | yes | yes |
| webpack | checkers4 | 1.189 / 1368.2 | 1.170 / 1185.2 | 1.179 / 1252.6 | 0.268 / 523.9 | 1.01x | 0.99x | yes | yes |
