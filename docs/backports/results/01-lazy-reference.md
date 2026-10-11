# Compiler comparison

2026-10-11T21:31:46.744104+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | --- | --- |
| references | single | 0.377 / 203.2 | 0.078 / 50.2 | 0.019 / 22.3 | 4.84x | yes | yes |
| references | checkers4 | 0.155 / 208.8 | 0.042 / 55.4 | 0.021 / 35.5 | 3.71x | yes | yes |
| members | single | 0.821 / 430.3 | 0.818 / 429.2 | 0.065 / 78.0 | 1.00x | yes | yes |
| members | checkers4 | 0.445 / 418.0 | 0.451 / 412.0 | 0.071 / 91.2 | 0.99x | yes | yes |
| unions | single | 0.367 / 115.7 | 0.366 / 129.0 | 0.096 / 58.0 | 1.00x | yes | yes |
| unions | checkers4 | 0.292 / 117.3 | 0.299 / 116.9 | 0.099 / 70.3 | 0.98x | yes | yes |
| flow | single | 0.081 / 59.4 | 0.082 / 59.6 | 0.020 / 25.1 | 0.98x | yes | yes |
| flow | checkers4 | 0.048 / 64.1 | 0.049 / 64.0 | 0.023 / 38.5 | 0.98x | yes | yes |
| modules | single | 0.056 / 51.2 | 0.058 / 51.5 | 0.017 / 23.0 | 0.98x | yes | yes |
| modules | checkers4 | 0.033 / 56.9 | 0.033 / 57.1 | 0.019 / 38.3 | 1.00x | yes | yes |
| errors | single | 0.042 / 46.0 | 0.041 / 46.2 | 0.013 / 21.0 | 1.02x | yes | yes |
| errors | checkers4 | 0.027 / 51.6 | 0.027 / 51.7 | 0.016 / 34.3 | 1.00x | yes | yes |
