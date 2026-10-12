#!/usr/bin/env python3
"""Measure VS Code with the fork, official nightly and pinned tsrs for the README."""
import argparse
import datetime
import importlib.util
import json
import os
import platform
import statistics
import subprocess
from pathlib import Path

from run import ROOT, MODES, digest, equal_outputs, run, signature

TOOLS = ('tsgo', 'tsgo nightly', 'tsrs')
MODES_MAIN = ('default', 'single', 'checkers4', 'checkers8')


def stable(warmup, repetitions):
    return bool(repetitions) and equal_outputs(warmup, [warmup, *repetitions])


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--fork', type=Path, required=True)
    parser.add_argument('--nightly', type=Path, required=True)
    parser.add_argument('--nightly-provenance', type=Path, required=True)
    parser.add_argument('--tsrs-source', type=Path, required=True)
    parser.add_argument('--tsrs', type=Path, required=True)
    parser.add_argument('--work-dir', type=Path, default=ROOT / '.bench/work')
    parser.add_argument('--out', type=Path, default=ROOT / '.bench/readme/result.json')
    parser.add_argument('--reps', type=int, default=5)
    args = parser.parse_args()
    if not 1 <= args.reps <= 20:
        parser.error('repetitions must be between 1 and 20')
    pins = json.loads((ROOT / 'bench/pins.json').read_text())
    tsrs_source = args.tsrs_source.resolve()
    actual = subprocess.check_output(['git', '-C', str(tsrs_source), 'rev-parse', 'HEAD'], text=True).strip()
    if actual != pins['tsrs']['commit']:
        parser.error('tsrs must match bench/pins.json')
    spec = importlib.util.spec_from_file_location('tsrs_projects', tsrs_source / 'bench/run.py')
    upstream = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(upstream)
    cfg = json.loads((tsrs_source / 'bench/projects.json').read_text())
    project = next(p for p in cfg['projects'] if p['name'] == 'vscode')
    upstream.setup_project(cfg, project, args.work_dir.resolve())
    _, config = upstream.project_path(cfg, project, args.work_dir.resolve())
    executables = dict(zip(TOOLS, (args.fork.resolve(), args.nightly.resolve(), args.tsrs.resolve())))
    machine = {'platform': platform.platform(), 'cpus': os.cpu_count(),
               'runner': os.environ.get('RUNNER_NAME', 'local'), 'arch': platform.machine()}
    if Path('/proc/cpuinfo').is_file():
        machine['cpu_model'] = next((line.split(':', 1)[1].strip() for line in Path('/proc/cpuinfo').read_text().splitlines() if line.startswith('model name')), 'unknown')
        machine['mem_total'] = next(line for line in Path('/proc/meminfo').read_text().splitlines() if line.startswith('MemTotal:'))
    head = subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=ROOT, text=True).strip()
    result = {'schema': 'readme-v1', 'head': head, 'date': datetime.datetime.now(datetime.timezone.utc).isoformat(),
              'machine': machine, 'reps': args.reps, 'project': project, 'tools': {}, 'rows': [],
              'run_url': os.environ.get('BENCH_RUN_URL', ''), 'build': 'fork: Go -pgo=off; tsrs: plain Rust release; nightly: official npm binary'}
    for name, exe in executables.items():
        metadata = {'version': subprocess.check_output([str(exe), '--version'], text=True).strip(), 'sha256': digest(exe)}
        if name != 'tsrs':
            metadata['go_build'] = subprocess.check_output(['go', 'version', '-m', str(exe)], text=True)
        if name == 'tsgo nightly':
            metadata['npm'] = json.loads(args.nightly_provenance.read_text())
        else:
            metadata['commit'] = head if name == 'tsgo' else actual
        result['tools'][name] = metadata
    args.out.parent.mkdir(parents=True, exist_ok=True)
    logs = args.out.parent / 'logs'
    logs.mkdir(exist_ok=True)
    failed = False
    for mode in MODES_MAIN:
        warmups = {name: run(exe, config, mode, logs / f'{mode}-{name.replace(" ", "-")}-warm.log', 600) for name, exe in executables.items()}
        runs = {name: [] for name in TOOLS}
        for rep in range(args.reps):
            order = list(TOOLS[rep % 3:] + TOOLS[:rep % 3])
            if rep % 2:
                order.reverse()
            for name in order:
                runs[name].append(run(executables[name], config, mode, logs / f'{mode}-{name.replace(" ", "-")}-{rep}.log', 600))
        stability = {name: stable(warmups[name], runs[name]) for name in TOOLS}
        failed |= not all(stability.values())
        if result['rows']:
            failed |= any(signature(warmups[name]) != signature(result['rows'][0]['warmups'][name]) for name in TOOLS)
        summaries = {name: {key: statistics.median(r[key] for r in rs) for key in ('wall_s', 'peak_rss_bytes', 'cpu_s', 'errors')} for name, rs in runs.items()}
        equality = {name: equal_outputs(warmups['tsgo'], [warmups[name], *runs[name]]) for name in TOOLS}
        result['rows'].append({'mode': mode, 'summary': summaries, 'stable': stability, 'equality': equality, 'runs': runs, 'warmups': warmups})
        print(mode, summaries, flush=True)
    result['completed'] = not failed
    # Published measurements use portable log paths. Full logs stay in the Actions artifact.
    for row in result['rows']:
        for name in TOOLS:
            for r in [row['warmups'][name], *row['runs'][name]]:
                r['log'] = 'logs/' + Path(r['log']).name
    args.out.write_text(json.dumps(result, indent=2) + '\n')
    if failed:
        raise SystemExit('incomplete or unstable compiler output; refusing to publish')


if __name__ == '__main__':
    main()
