# Native VS Code comparison

On pinned VS Code, the final fork default completed in **3.829s**, compared with **5.857s** for stock Go: 34.6% less elapsed time and 11.3% less peak RSS. Defaults used four checkers for stock Go, eight for this fork and sixteen for tsrs on this machine. Equal-checker results are separate below.

The corpus is `microsoft/vscode@9cf0128b9822dcd8a533f08ef7312956a4390301`, project `src`. Stock Go is `microsoft/TypeScript@b85298b6a81f772d080b0455de0ca9d744cd6fd6`; tsrs is `maschwenk/tsrs@e4c7d83a5907d1b8b2dc2f444cef40e5dc2a3751`. Every warmup and measured run completed with identical full diagnostic text, exit status and 10,465-file count. The corpus reports 371 diagnostics at this common compiler revision; the comparison has matching errors.

## Final default after step 52

Fork source: `e290a44cb1d2b04cde6a6f7e92c738813dbcc0ad`. [Raw repetitions, warmups and build metadata](adaptive/compare.json).

| Compiler | Default checkers | Median wall time | Median peak RSS |
| --- | ---: | ---: | ---: |
| Stock Go | 4 | 5.857s | 7,106.0 MiB |
| This fork | 8 | 3.829s | 6,299.9 MiB |
| tsrs | 16 | 0.857s | 2,145.8 MiB |

The fork is 1.53 times as fast as stock at their defaults. tsrs still completes this workload about 4.47 times as fast as the fork. Step 52 changes the ordinary CLI's automatic checker count; it does not imply a 35% improvement at equal checker counts.

## Equal checker counts before the default-policy change

Fork source: `e0c7ad31fbc289dd9c62b848dd13222696402169`, containing the preceding implementation and team tooling. [Full matrix with raw repetitions, warmups and build metadata](matrix/compare.json).

| Mode | Stock Go time / RSS | Fork time / RSS | tsrs time / RSS | Stock time divided by fork |
| --- | ---: | ---: | ---: | ---: |
| Single checker | 23.852s / 6,676.8 MiB | 20.857s / 5,700.6 MiB | 6.679s / 1,616.7 MiB | 1.14 |
| Four checkers | 5.405s / 7,263.2 MiB | 4.792s / 5,815.5 MiB | 1.723s / 1,829.1 MiB | 1.13 |
| Eight checkers | 4.308s / 8,040.1 MiB | 3.676s / 6,524.0 MiB | 1.180s / 1,974.1 MiB | 1.17 |
| Earlier defaults (Go 4, fork 4, tsrs 16) | 5.886s / 7,194.0 MiB | 5.075s / 5,850.7 MiB | 0.814s / 2,150.3 MiB | 1.16 |

At equal checker counts the fork uses 11–15% less wall time and 15–20% less peak RSS than stock. These rows and the final default are two separate runs at two recorded source revisions. The final default row is not part of this earlier binary's matrix.

## Method and correctness

The executables ran directly on an Apple M5 Max Mac with 18 logical CPUs, 128 GiB RAM and macOS 26.6.2 ARM64. Idle CI containers were paused during timing. Each mode warmed all compilers once, then collected five interleaved repetitions, rotating and reversing execution order. Filesystem caches were warm. Go stock and fork used Go 1.27 and `-pgo=off`; tsrs used plain Rust release with no Rust flags. Single mode capped both runtime thread pools at one; other modes cleared runtime pool overrides.

The compiler baselines passed throughout the series. The combined all-package shadow run had one unchanged Darwin fswatch timeout; that package passed its isolated retry. Final checker/compiler race tests passed. Seeded differential checks of diagnostics and emitted JS/declaration bytes passed for 64 seeds in both normal and all-shadow configurations, with single and four-checker modes. All-shadow VS Code checks also matched stock, including the final eight-checker default.

Go PGO training and a separate PGO build were validated, including compiler baselines. This report measures plain builds and claims no PGO application speedup. Instruction counts are unavailable on this Mac. These results describe this pinned workload and do not reproduce tsrs's published PGO/BOLT release figures or establish a universal speedup.

All raw compiler logs are compressed alongside their report. [Publication metadata](publication.json) records the original and public JSON hashes. Local checkout prefixes were replaced with `<repo>` and `<tsrs>`; measurements and binary hashes are unchanged. Each JSON entry retains its original diagnostic hash and also records the hash of its published log after path normalization. All 90 published logs were parsed again to verify their hashes, file counts, exits and complete diagnostic equality. See [reproduction commands](../../../../bench/README.md).
