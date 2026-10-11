# Compiler comparison

2026-10-11T21:54:02.565930+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.448 / 152.0 | 0.351 / 146.1 | 0.347 / 146.7 | 0.146 / 56.9 | 1.29x | 1.01x | yes | yes |
| Compiler | checkers4 | 0.125 / 173.6 | 0.126 / 165.5 | 0.125 / 174.1 | 0.060 / 79.6 | 1.00x | 1.01x | yes | yes |
| Compiler-Unions | single | 0.580 / 172.8 | 0.578 / 163.1 | 0.584 / 169.4 | 0.181 / 59.6 | 0.99x | 0.99x | yes | yes |
| Compiler-Unions | checkers4 | 0.218 / 201.9 | 0.221 / 195.4 | 0.220 / 195.5 | 0.083 / 85.0 | 0.99x | 1.00x | yes | yes |
| xstate-main | single | 1.860 / 678.4 | 1.693 / 609.0 | 1.562 / 593.3 | 0.459 / 192.1 | 1.19x | 1.08x | yes | yes |
| xstate-main | checkers4 | 0.484 / 811.6 | 0.459 / 687.4 | 0.451 / 690.6 | 0.144 / 254.9 | 1.07x | 1.02x | yes | yes |
| webpack | single | 3.524 / 1195.3 | 3.437 / 1156.9 | 3.501 / 1173.7 | 0.874 / 439.5 | 1.01x | 0.98x | yes | yes |
| webpack | checkers4 | 1.312 / 1328.1 | 1.221 / 1184.8 | 1.250 / 1177.2 | 0.274 / 522.8 | 1.05x | 0.98x | yes | yes |
