# Compiler comparison

2026-10-11T22:11:55.226000+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.322 / 147.7 | 0.322 / 148.4 | 0.319 / 148.3 | 0.119 / 56.9 | 1.01x | 1.01x | yes | yes |
| Compiler | checkers4 | 0.117 / 178.5 | 0.117 / 175.8 | 0.114 / 174.6 | 0.054 / 79.5 | 1.02x | 1.02x | yes | yes |
| xstate-main | single | 1.820 / 687.8 | 1.661 / 616.2 | 1.673 / 630.9 | 0.444 / 192.1 | 1.09x | 0.99x | yes | yes |
| xstate-main | checkers4 | 0.481 / 817.0 | 0.456 / 679.8 | 0.448 / 692.2 | 0.143 / 255.0 | 1.07x | 1.02x | yes | yes |
