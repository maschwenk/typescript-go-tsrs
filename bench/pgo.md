# Native Go PGO experiment

This is a separate build using Go's native profile-guided optimization. Ordinary stock and fork builds use `-pgo=off`; no `default.pgo` is committed. Rust LTO/BOLT flags do not transfer to Go's toolchain.

Build the plain fork and prepare the pinned VS Code corpus, then run on the same native machine as the final comparison:

```sh
python3 bench/profile.py --compiler .bench/bin/fork --tsrs-source /path/to/tsrs \
  --projects vscode --reps 10
python3 bench/run.py --stock .bench/bin/stock --parent .bench/bin/fork \
  --fork .bench/pgo/tsgo-pgo --tsrs /path/to/tsrs/target/release/tsrs \
  --tsrs-source /path/to/tsrs --projects vscode --modes single,checkers4,checkers8,default \
  --reps 5 --label native-mac-pgo --out .bench/results/vscode-pgo.json
```

`profile.py` collects ten single-threaded CPU profiles from completed VS Code checks, verifies stable diagnostics/exits/file counts, merges the profiles with Go's pprof tool and builds `.bench/pgo/tsgo-pgo`. Raw profiles, logs, profile hashes, compiler hashes, source revision, platform and build flags are retained in metadata. The PGO workflow also runs compiler baselines using the collected profile.

Training and evaluation use the same pinned VS Code workload. Results measure that workload's specialization, not held-out performance or generalization to other projects. A profile depends on source and platform; retrain after relevant source changes and compare plain versus PGO before selecting a release build. The ordinary harness still checks diagnostic equality on every warmup and measured repetition.

The manual `Native Go PGO` Actions workflow reproduces training, conformance checking and the VS Code comparison on Linux. Its timings are separate from native Mac results. See the [Go PGO guide](https://go.dev/doc/pgo).
