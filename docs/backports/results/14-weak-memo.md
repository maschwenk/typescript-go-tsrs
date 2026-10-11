# Compiler comparison

2026-10-11T22:16:55.308045+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | parent s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | parent/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Compiler | single | 0.322 / 149.7 | 0.316 / 147.7 | 0.320 / 146.9 | 0.120 / 56.9 | 1.01x | 0.99x | yes | yes |
| Compiler | checkers4 | 0.116 / 178.8 | 0.115 / 165.1 | 0.116 / 174.8 | 0.055 / 79.6 | 1.00x | 0.99x | yes | yes |
| Compiler-Unions | single | 0.574 / 172.9 | 0.567 / 171.1 | 0.559 / 165.2 | 0.180 / 59.6 | 1.03x | 1.01x | yes | yes |
| Compiler-Unions | checkers4 | 0.217 / 201.4 | 0.219 / 195.9 | 0.216 / 193.0 | 0.082 / 84.9 | 1.00x | 1.01x | yes | yes |
| xstate-main | single | 1.680 / 666.7 | 1.671 / 620.2 | 1.674 / 634.6 | 0.445 / 192.1 | 1.00x | 1.00x | yes | yes |
| xstate-main | checkers4 | 0.485 / 810.2 | 0.448 / 681.7 | 0.451 / 690.6 | 0.146 / 254.5 | 1.08x | 0.99x | yes | yes |
| webpack | single | 3.578 / 1238.0 | 3.481 / 1183.9 | 3.273 / 1160.2 | 0.880 / 439.5 | 1.09x | 1.06x | yes | yes |
| webpack | checkers4 | 1.207 / 1369.0 | 1.174 / 1253.1 | 1.158 / 1268.4 | 0.266 / 522.2 | 1.04x | 1.01x | yes | yes |
