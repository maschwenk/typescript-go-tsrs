# Compiler comparison

2026-10-11T22:08:28.951063+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.326 / 148.2 | 0.320 / 145.8 | 0.320 / 147.1 | 0.120 / 56.9 | 1.02x | 1.00x | yes | yes |
| Compiler | checkers4 | 0.116 / 180.6 | 0.118 / 172.1 | 0.118 / 176.1 | 0.054 / 79.5 | 0.99x | 1.00x | yes | yes |
| Compiler-Unions | single | 0.572 / 173.1 | 0.571 / 163.2 | 0.575 / 169.3 | 0.192 / 59.6 | 0.99x | 0.99x | yes | yes |
| Compiler-Unions | checkers4 | 0.216 / 193.7 | 0.222 / 193.9 | 0.216 / 195.8 | 0.082 / 85.2 | 1.00x | 1.03x | yes | yes |
| xstate-main | single | 1.820 / 673.7 | 1.671 / 623.4 | 1.681 / 609.9 | 0.442 / 192.1 | 1.08x | 0.99x | yes | yes |
| xstate-main | checkers4 | 0.478 / 832.2 | 0.446 / 690.5 | 0.453 / 692.0 | 0.144 / 255.9 | 1.05x | 0.98x | yes | yes |
| webpack | single | 3.602 / 1240.8 | 3.484 / 1174.7 | 3.362 / 1161.3 | 0.896 / 439.5 | 1.07x | 1.04x | yes | yes |
| webpack | checkers4 | 1.204 / 1400.5 | 1.158 / 1253.8 | 1.168 / 1215.4 | 0.266 / 523.8 | 1.03x | 0.99x | yes | yes |
