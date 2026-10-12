# Compiler comparison

2026-10-11T23:46:28.798427+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.414 / 141.9 | 0.427 / 142.0 | 0.412 / 132.6 | 0.152 / 56.9 | 1.01x | 1.04x | yes | yes |
| Compiler | checkers4 | 0.157 / 179.4 | 0.166 / 169.4 | 0.154 / 170.6 | 0.073 / 79.6 | 1.02x | 1.08x | yes | yes |
| Compiler-Unions | single | 0.697 / 168.3 | 0.576 / 150.6 | 0.562 / 154.6 | 0.207 / 59.6 | 1.24x | 1.03x | yes | yes |
| Compiler-Unions | checkers4 | 0.235 / 190.2 | 0.217 / 180.8 | 0.218 / 180.2 | 0.090 / 84.6 | 1.08x | 0.99x | yes | yes |
| xstate-main | single | 1.837 / 647.0 | 1.604 / 572.6 | 1.598 / 581.1 | 0.449 / 192.1 | 1.15x | 1.00x | yes | yes |
| xstate-main | checkers4 | 0.502 / 819.9 | 0.433 / 662.9 | 0.438 / 659.0 | 0.145 / 256.0 | 1.15x | 0.99x | yes | yes |
| webpack | single | 4.448 / 1199.8 | 4.236 / 1145.4 | 4.221 / 1130.4 | 1.146 / 439.5 | 1.05x | 1.00x | yes | yes |
| webpack | checkers4 | 1.256 / 1399.2 | 1.161 / 1153.6 | 1.178 / 1191.9 | 0.288 / 522.2 | 1.07x | 0.99x | yes | yes |
