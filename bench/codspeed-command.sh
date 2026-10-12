#!/usr/bin/env bash
set -euo pipefail

repo=$(cd -- "$(dirname -- "$0")/.." && pwd)
state=${TSGO_CODSPEED_STATE:-"$repo/.bench/results/codspeed-vscode"}
source "$state/paths.env"

case ${1:-} in
  stock) compiler=$STOCK ;;
  fork) compiler=$FORK ;;
  tsrs) compiler=$TSRS ;;
  *) echo 'expected stock, fork, or tsrs' >&2; exit 64 ;;
esac
mode=${2:-}
flags=()
unset GOMAXPROCS RAYON_NUM_THREADS
case "$mode" in
  single)
    flags+=(--singleThreaded)
    export GOMAXPROCS=1 RAYON_NUM_THREADS=1
    ;;
  checkers4) flags+=(--checkers 4) ;;
  default) ;;
  *) echo 'expected single, checkers4, or default' >&2; exit 64 ;;
esac

cd -- "$PROJECT_CWD"
log=$(mktemp "$LOGS/$1-$mode.XXXXXX.log")
set +e
timeout --signal=KILL 600 "$compiler" -p "$PROJECT" \
  --noEmit --incremental false --pretty false --extendedDiagnostics \
  "${flags[@]}" > "$log" 2>&1
status=$?
set -e
printf '%s\n' "$status" > "$log.exit"

# Diagnostic exits remain in the sidecar and are checked after all repetitions.
case "$status" in
  0|1|2) exit 0 ;;
  *) cat "$log" >&2; exit "$status" ;;
esac
