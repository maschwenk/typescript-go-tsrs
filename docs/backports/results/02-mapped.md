# Compiler comparison

2026-10-11T21:34:14.743510+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | --- | --- |
| references | single | 0.352 / 203.4 | 0.076 / 50.2 | 0.018 / 22.3 | 4.62x | yes | yes |
| references | checkers4 | 0.155 / 208.2 | 0.042 / 55.5 | 0.021 / 35.5 | 3.71x | yes | yes |
| members | single | 0.823 / 429.1 | 0.387 / 223.1 | 0.066 / 78.0 | 2.13x | yes | yes |
| members | checkers4 | 0.439 / 418.8 | 0.201 / 224.9 | 0.070 / 91.1 | 2.18x | yes | yes |
| unions | single | 0.365 / 128.0 | 0.369 / 129.7 | 0.095 / 58.0 | 0.99x | yes | yes |
| unions | checkers4 | 0.369 / 120.7 | 0.330 / 117.5 | 0.105 / 70.3 | 1.12x | yes | yes |
| flow | single | 0.082 / 59.5 | 0.083 / 59.4 | 0.020 / 25.1 | 0.99x | yes | yes |
| flow | checkers4 | 0.050 / 63.3 | 0.049 / 63.9 | 0.024 / 38.4 | 1.00x | yes | yes |
| modules | single | 0.058 / 51.3 | 0.058 / 51.4 | 0.018 / 23.0 | 1.00x | yes | yes |
| modules | checkers4 | 0.032 / 57.1 | 0.031 / 57.4 | 0.019 / 38.1 | 1.01x | yes | yes |
| errors | single | 0.042 / 45.7 | 0.043 / 46.0 | 0.013 / 21.0 | 0.97x | yes | yes |
| errors | checkers4 | 0.028 / 52.3 | 0.027 / 52.2 | 0.017 / 34.2 | 1.05x | yes | yes |
