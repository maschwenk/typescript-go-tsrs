#!/usr/bin/env python3
"""Prepare VS Code commands and check existing CodSpeed logs without rerunning compilers."""
from __future__ import annotations

import argparse
import hashlib
import importlib.util
import json
import math
import os
import platform
import re
import shlex
import subprocess
import sys
from pathlib import Path

sys.dont_write_bytecode = True
ROOT = Path(os.environ.get('TSGO_CODSPEED_REPO', Path(__file__).resolve().parents[1])).resolve()
VSCODE_SHA = '9cf0128b9822dcd8a533f08ef7312956a4390301'
TSRS_SHA = 'e4c7d83a5907d1b8b2dc2f444cef40e5dc2a3751'
TOOLS = ('stock', 'fork', 'tsrs')
MODES = ('single', 'checkers4', 'default')


def load_module(name, path):
    spec = importlib.util.spec_from_file_location(name, path)
    module = importlib.util.module_from_spec(spec)
    sys.path.insert(0, str(path.parent))
    try:
        spec.loader.exec_module(module)
    finally:
        sys.path.pop(0)
    return module


def git_head(path):
    return subprocess.check_output(['git', '-C', str(path), 'rev-parse', 'HEAD'], text=True).strip()


def sha256(path):
    digest = hashlib.sha256()
    with path.open('rb') as stream:
        for block in iter(lambda: stream.read(1024 * 1024), b''):
            digest.update(block)
    return digest.hexdigest()


def prepare(args):
    pins = json.loads((ROOT / 'bench/pins.json').read_text())
    if pins['tsrs']['commit'] != TSRS_SHA or git_head(args.tsrs_source) != TSRS_SHA:
        raise ValueError('tsrs source must match the inspected bench/pins.json revision')
    cfg = json.loads((args.tsrs_source / 'bench/projects.json').read_text())
    project = next(p for p in cfg['projects'] if p['name'] == 'vscode')
    expected = {'repo': 'https://github.com/microsoft/vscode.git', 'commit': VSCODE_SHA,
                'project': 'src', 'install': 'npm ci --ignore-scripts --no-audit --no-fund'}
    if any(project.get(key) != value for key, value in expected.items()):
        raise ValueError('pinned tsrs VS Code setup differs from the reviewed workload')
    upstream = load_module('codspeed_tsrs_setup', args.tsrs_source / 'bench/run.py')
    work = (args.prepared_work_dir or ROOT / '.bench/codspeed/work').resolve()
    if args.prepared_work_dir is None:
        upstream.setup_project(cfg, project, work)
    cwd, target = upstream.project_path(cfg, project, work)
    if git_head(cwd) != VSCODE_SHA or not (cwd / 'node_modules').is_dir():
        raise ValueError('prepared VS Code checkout must have the pinned revision and installed dependencies')
    if not (target / 'tsconfig.json').is_file():
        raise ValueError('missing VS Code src/tsconfig.json')
    state = args.state.resolve()
    if state.exists():
        raise ValueError(f'{state} already exists; choose a fresh state directory to preserve previous results')
    binaries = {'stock': ROOT / '.bench/bin/stock', 'fork': ROOT / '.bench/bin/fork',
                'tsrs': args.tsrs_source / 'target/release/tsrs'}
    metadata = {}
    for name, exe in binaries.items():
        exe = exe.resolve()
        if not exe.is_file():
            raise ValueError(f'missing native compiler: {exe}')
        item = {'path': str(exe), 'sha256': sha256(exe),
                'version': subprocess.check_output([str(exe), '--version'], text=True).strip()}
        if name != 'tsrs':
            item['go_build'] = subprocess.check_output(['go', 'version', '-m', str(exe)], text=True)
            if re.search(r'\tbuild\t-pgo=(?!off(?:\n|$))', item['go_build']):
                raise ValueError(f'{name} must be built with bench/build.py (-pgo=off)')
        metadata[name] = item
    state.mkdir(parents=True)
    logs = state / 'logs'
    logs.mkdir()
    values = {'PROJECT': target.resolve(), 'PROJECT_CWD': target.parent.resolve(), 'LOGS': logs,
              **{name.upper(): path.resolve() for name, path in binaries.items()}}
    (state / 'paths.env').write_text(''.join(f'{key}={shlex.quote(str(value))}\n' for key, value in values.items()))
    manifest = {'pins': pins, 'head': git_head(ROOT), 'project': expected, 'modes': args.modes.split(','),
                'repetitions': args.repetitions, 'tools': metadata, 'repo': str(ROOT),
                'machine': {'platform': platform.platform(), 'cpus': os.cpu_count(),
                            'affinity': sorted(os.sched_getaffinity(0)) if hasattr(os, 'sched_getaffinity') else None},
                'harness_sha256': sha256(ROOT / 'bench/run.py'),
                'build': 'Go -pgo=off; tsrs cargo +1.99.0 build --locked --release -p tsrs_cli'}
    (state / 'manifest.json').write_text(json.dumps(manifest, indent=2) + '\n')
    config = ['options:', f'  working-directory: {json.dumps(str(ROOT))}', '  warmup-time: 1ms',
              f'  min-rounds: {args.repetitions}', f'  max-rounds: {args.repetitions}', '', 'benchmarks:']
    for mode in manifest['modes']:
        for tool in TOOLS:
            config += [f'  - name: vscode/src/{mode}/{tool}',
                       f'    exec: bash bench/codspeed-command.sh {tool} {mode}']
    (state / 'run.yml').write_text('\n'.join(config) + '\n')
    print(state)


def parse_log(log, harness):
    lines = log.read_text(errors='replace').splitlines()
    status = int(log.with_suffix(log.suffix + '.exit').read_text().strip())
    return {**harness.parse_output(lines, status), 'log':str(log)}


def gate(state, strict_tsrs=False):
    manifest = json.loads((state / 'manifest.json').read_text())
    harness_path = Path(manifest['repo']) / 'bench/run.py'
    if sha256(harness_path) != manifest['harness_sha256']:
        raise ValueError('bench/run.py changed since preparation; refusing to change the equality gate mid-run')
    harness = load_module('codspeed_equality_gate', harness_path)
    for item in manifest['tools'].values():
        if sha256(Path(item['path'])) != item['sha256']:
            raise ValueError('compiler binary changed during measurement')
    timing = {}
    for path in (state / 'profile/results').glob('*.json'):
        data = json.loads(path.read_text())
        if data.get('instrument', {}).get('type') != 'walltime':
            continue
        for benchmark in data['benchmarks']:
            if benchmark['name'] in timing:
                raise ValueError('duplicate CodSpeed benchmark results')
            timing[benchmark['name']] = benchmark['stats']
    passed, rows = True, []
    summary = ['# VS Code src: stock / fork / tsrs', '',
               f'Corpus `{VSCODE_SHA}`. Plain native builds; grouped CodSpeed repetitions; profiling disabled.', '',
               '| mode | stock s | fork s | tsrs s | stock/fork | stock/tsrs | Go equality | tsrs equality |',
               '| --- | ---: | ---: | ---: | ---: | ---: | --- | --- |']
    for mode in manifest['modes']:
        runs, medians = {}, {}
        complete = True
        for tool in TOOLS:
            runs[tool] = [parse_log(log, harness) for log in sorted((state / 'logs').glob(f'{tool}-{mode}.*.log'))]
            stats = timing.get(f'vscode/src/{mode}/{tool}', {})
            count_ok = len(runs[tool]) == manifest['repetitions'] + 1
            timing_ok = stats.get('rounds') == manifest['repetitions']
            median = stats.get('median_ns', 0) / 1e9
            complete &= count_ok and timing_ok and median > 0 and math.isfinite(median)
            complete &= all(r['completed'] for r in runs[tool])
            medians[tool] = median
        reference = harness.signature(runs['stock'][0]) if runs['stock'] else None
        equality = {tool: bool(rs) and all(r['completed'] and harness.signature(r) == reference for r in rs)
                    for tool, rs in runs.items()}
        go_equal = equality['stock'] and equality['fork']
        tsrs_equal = equality['stock'] and equality['tsrs']
        passed &= complete and go_equal and (equality['tsrs'] or not strict_tsrs)
        rows.append({'mode': mode, 'complete': complete, 'equality': equality, 'median_wall_s': medians, 'runs': runs})
        cells = [f'{medians[tool]:.3f}' if complete else 'incomplete' for tool in TOOLS]
        fork_ratio = f"{medians['stock']/medians['fork']:.2f}x" if complete and go_equal else 'invalid'
        tsrs_ratio = f"{medians['stock']/medians['tsrs']:.2f}x" if complete and tsrs_equal else 'invalid'
        summary.append(f"| {mode} | {' | '.join(cells)} | {fork_ratio} | {tsrs_ratio} | "
                       f"{'yes' if go_equal else 'MISMATCH'} | {'yes' if tsrs_equal else 'MISMATCH'} |")
    summary += ['', 'tsrs diagnostic/exit/file-count equality is advisory unless --gate-tsrs is set.',
                'Every compiler must complete. Speed ratios are suppressed for unequal or incomplete results.', '']
    result = {'manifest': manifest, 'gate_passed': bool(passed), 'gate_tsrs': strict_tsrs, 'rows': rows}
    (state / 'compare.json').write_text(json.dumps(result, indent=2) + '\n')
    (state / 'compare.md').write_text('\n'.join(summary))
    print('\n'.join(summary))
    return bool(passed)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest='command', required=True)
    prep = commands.add_parser('prepare')
    prep.add_argument('--tsrs-source', type=Path, default=ROOT / '.bench/tsrs')
    prep.add_argument('--prepared-work-dir', type=Path)
    prep.add_argument('--repetitions', type=int, default=7)
    prep.add_argument('--modes', default=','.join(MODES))
    check = commands.add_parser('gate')
    check.add_argument('--gate-tsrs', action='store_true')
    for command in (prep, check):
        command.add_argument('--state', type=Path, default=ROOT / '.bench/results/codspeed-vscode')
    args = parser.parse_args()
    try:
        if args.command == 'prepare':
            if not 1 <= args.repetitions <= 50 or not set(args.modes.split(',')) <= set(MODES):
                parser.error('use 1-50 repetitions and only single,checkers4,default modes')
            if len(set(args.modes.split(','))) != len(args.modes.split(',')):
                parser.error('modes must not repeat')
            prepare(args)
        elif not gate(args.state.resolve(), args.gate_tsrs):
            sys.exit(1)
    except (ValueError, OSError, KeyError, StopIteration) as error:
        parser.exit(1, f'{error}\n')


if __name__ == '__main__':
    main()
