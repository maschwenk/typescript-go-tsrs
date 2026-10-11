# Compiler comparison

2026-10-11T21:29:00.473851+00:00 · local · macOS-26.6.2-arm64-arm-64bit · 18 logical CPUs

Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.

| project | mode | stock s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | Go equality | tsrs equality |
| --- | --- | ---: | ---: | ---: | ---: | --- | --- |
| members | single | 0.824 / 430.9 | 0.832 / 428.4 | 0.066 / 77.9 | 0.99x | yes | yes |
| members | checkers4 | 0.546 / 423.5 | 0.579 / 425.9 | 0.072 / 91.1 | 0.94x | yes | yes |
| members | default | 0.473 / 412.1 | 0.484 / 422.7 | 0.072 / 91.1 | 0.98x | yes | yes |
| unions | single | 0.387 / 123.6 | 0.398 / 117.2 | 0.099 / 58.0 | 0.97x | yes | yes |
| unions | checkers4 | 0.407 / 122.4 | 0.399 / 119.4 | 0.126 / 70.3 | 1.02x | yes | yes |
| unions | default | 0.322 / 118.5 | 0.323 / 121.0 | 0.102 / 70.3 | 1.00x | yes | yes |
| flow | single | 0.084 / 59.5 | 0.083 / 59.8 | 0.020 / 25.1 | 1.01x | yes | yes |
| flow | checkers4 | 0.050 / 64.0 | 0.050 / 64.3 | 0.024 / 38.5 | 0.99x | yes | yes |
| flow | default | 0.051 / 63.9 | 0.051 / 64.3 | 0.027 / 38.4 | 1.01x | yes | yes |
| modules | single | 0.062 / 51.4 | 0.061 / 51.4 | 0.018 / 23.0 | 1.01x | yes | yes |
| modules | checkers4 | 0.035 / 57.4 | 0.035 / 57.3 | 0.022 / 38.0 | 1.00x | yes | yes |
| modules | default | 0.033 / 57.4 | 0.035 / 57.2 | 0.020 / 39.3 | 0.94x | yes | yes |
| errors | single | 0.045 / 45.8 | 0.044 / 45.7 | 0.014 / 21.0 | 1.02x | yes | yes |
| errors | checkers4 | 0.029 / 51.9 | 0.029 / 51.4 | 0.020 / 34.1 | 1.02x | yes | yes |
| errors | default | 0.028 / 52.2 | 0.029 / 51.8 | 0.026 / 34.1 | 0.97x | yes | yes |
