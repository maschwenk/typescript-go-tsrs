# Compiler comparison

2026-10-11T21:36:39.631255+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.318 / 149.8 | 0.321 / 143.5 | 0.119 / 56.9 | 0.99x | yes | yes |
| Compiler | checkers4 | 0.117 / 179.7 | 0.116 / 176.3 | 0.055 / 79.6 | 1.00x | yes | yes |
| Compiler-Unions | single | 0.554 / 162.4 | 0.575 / 168.4 | 0.178 / 59.6 | 0.96x | yes | yes |
| Compiler-Unions | checkers4 | 0.217 / 201.0 | 0.217 / 196.0 | 0.083 / 84.8 | 1.00x | yes | yes |
| xstate-main | single | 1.743 / 666.1 | 1.657 / 611.8 | 0.441 / 192.1 | 1.05x | yes | yes |
| xstate-main | checkers4 | 0.488 / 830.0 | 0.450 / 692.0 | 0.144 / 255.1 | 1.09x | yes | yes |
| webpack | single | 3.568 / 1237.7 | 3.442 / 1178.2 | 0.875 / 439.5 | 1.04x | yes | yes |
| webpack | checkers4 | 1.193 / 1364.3 | 1.228 / 1268.2 | 0.265 / 522.4 | 0.97x | yes | yes |
