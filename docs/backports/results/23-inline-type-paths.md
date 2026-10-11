# Compiler comparison

2026-10-11T22:56:17.877348+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.317 / 147.3 | 0.313 / 146.2 | 0.317 / 144.2 | 0.118 / 56.9 | 1.00x | 0.99x | yes | yes |
| Compiler | checkers4 | 0.117 / 179.9 | 0.113 / 171.2 | 0.111 / 171.6 | 0.054 / 79.6 | 1.05x | 1.02x | yes | yes |
| Compiler-Unions | single | 0.571 / 173.8 | 0.569 / 165.3 | 0.558 / 166.7 | 0.177 / 59.6 | 1.02x | 1.02x | yes | yes |
| Compiler-Unions | checkers4 | 0.217 / 201.6 | 0.214 / 190.5 | 0.211 / 189.7 | 0.083 / 84.6 | 1.03x | 1.02x | yes | yes |
| xstate-main | single | 1.832 / 684.6 | 1.633 / 604.8 | 1.605 / 606.8 | 0.435 / 192.1 | 1.14x | 1.02x | yes | yes |
| xstate-main | checkers4 | 0.480 / 831.9 | 0.436 / 657.8 | 0.431 / 668.7 | 0.145 / 255.5 | 1.11x | 1.01x | yes | yes |
| webpack | single | 3.565 / 1232.9 | 3.339 / 1141.1 | 3.458 / 1161.0 | 0.891 / 439.5 | 1.03x | 0.97x | yes | yes |
| webpack | checkers4 | 1.188 / 1383.6 | 1.141 / 1231.0 | 1.143 / 1229.6 | 0.271 / 522.8 | 1.04x | 1.00x | yes | yes |
