#!/usr/bin/env python3
"""Interleaved native compiler benchmark with diagnostic equality gates."""
from __future__ import annotations
import argparse
import datetime
import hashlib
import importlib.util
import json
import os
import platform
import re
import statistics
import subprocess
import sys
import threading
import time
from pathlib import Path
from generate import generate

ROOT = Path(__file__).resolve().parent.parent
MODES = {'single': ['--singleThreaded'], 'checkers4': ['--checkers', '4'], 'checkers8': ['--checkers', '8'], 'default': []}
COUNTER = re.compile(r'^(Files|Lines|Identifiers|Symbols|Types|Instantiations|Memory allocs):\s+(\d+)\s*$')
TIMING = re.compile(r'^(Config|Parse|Bind|Check|Emit|Total) time:')
METRIC = re.compile(r'^(Memory used|Memory allocs|Memory frees|GC count|GC time|Checker count|Threads):')


def digest(path):
    h = hashlib.sha256()
    with path.open('rb') as f:
        for block in iter(lambda: f.read(1024 * 1024), b''):
            h.update(block)
    return h.hexdigest()


def run(exe, project, mode, log, timeout):
    argv = [str(exe), '-p', str(project), '--noEmit', '--incremental', 'false', '--pretty', 'false', '--extendedDiagnostics', *MODES[mode]]
    env = dict(os.environ)
    if mode == 'single':
        env['GOMAXPROCS'] = '1'
        env['RAYON_NUM_THREADS'] = '1'
    with log.open('wb') as out:
        t = time.perf_counter()
        p = subprocess.Popen(argv, stdout=out, stderr=subprocess.STDOUT, cwd=project.parent, env=env)
        timer = threading.Timer(timeout, p.kill)
        timer.start()
        try:
            _, status, ru = os.wait4(p.pid, 0)
        finally:
            timer.cancel()
        p.returncode = os.waitstatus_to_exitcode(status)
        wall = time.perf_counter() - t
    lines = log.read_text(errors='replace').splitlines()
    counters = {m[1]: int(m[2]) for line in lines if (m := COUNTER.match(line))}
    # Preserve the entire diagnostic including continuation lines. Strip only explicit metrics.
    diagnostic_lines = []
    in_diagnostic = False
    for line in lines:
        if re.search(r'error TS\d+:', line):
            in_diagnostic = True
            diagnostic_lines.append(line)
        elif COUNTER.match(line) or TIMING.match(line) or METRIC.match(line):
            in_diagnostic = False
        elif in_diagnostic and (not line or line[0].isspace()):
            diagnostic_lines.append(line)
        else:
            in_diagnostic = False
    diagnostics = '\n'.join(diagnostic_lines).strip()
    completed = any(line.startswith('Check time:') for line in lines) and p.returncode in (0, 1, 2)
    return {'exit': p.returncode, 'wall_s': wall, 'peak_rss_bytes': ru.ru_maxrss * (1 if sys.platform == 'darwin' else 1024),
            'cpu_s': ru.ru_utime + ru.ru_stime, 'counters': counters, 'completed': completed,
            'diagnostics_sha256': hashlib.sha256(diagnostics.encode()).hexdigest(),
            'errors': sum(bool(re.search(r'error TS\d+:', line)) for line in lines), 'log': str(log)}


def signature(r):
    return r['exit'], r['diagnostics_sha256'], r['counters'].get('Files')


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    for tool in ('stock', 'parent', 'fork', 'tsrs'):
        ap.add_argument('--' + tool, type=Path, required=tool in ('stock', 'fork', 'tsrs'))
    ap.add_argument('--tsrs-source', type=Path, required=True)
    ap.add_argument('--work-dir', type=Path, default=ROOT / '.bench/work')
    ap.add_argument('--out', type=Path, default=ROOT / '.bench/results/latest.json')
    ap.add_argument('--projects', default='synthetic', help='synthetic and/or names from tsrs bench/projects.json')
    ap.add_argument('--modes', default='single,checkers4,default')
    ap.add_argument('--reps', type=int, default=5)
    ap.add_argument('--timeout', type=int, default=600)
    ap.add_argument('--label', default='local')
    ap.add_argument('--gate-tsrs', action='store_true', help='also fail on known tsrs disagreements')
    args = ap.parse_args()
    if args.reps < 1 or not set(args.modes.split(',')) <= MODES.keys():
        ap.error('positive repetitions and supported modes required')
    tools = {name: path.resolve() for name in ('stock','parent','fork','tsrs') if (path := getattr(args,name))}
    for exe in tools.values():
        if not exe.is_file():
            ap.error(f'missing compiler: {exe}')
    args.work_dir = args.work_dir.resolve()
    names = args.projects.split(',')
    projects = generate(args.work_dir / 'synthetic') if 'synthetic' in names else []
    if any(name != 'synthetic' for name in names):
        spec = importlib.util.spec_from_file_location('tsrs_bench', args.tsrs_source / 'bench/run.py')
        upstream = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(upstream)
        cfg = json.loads((args.tsrs_source / 'bench/projects.json').read_text())
        available = {p['name']: p for p in cfg['projects']}
        for name in names:
            if name == 'synthetic':
                continue
            if name not in available:
                ap.error(f'unknown project {name}')
            upstream.setup_project(cfg, available[name], args.work_dir)
            _, project = upstream.project_path(cfg, available[name], args.work_dir)
            projects.append((name, project))
    args.out.parent.mkdir(parents=True, exist_ok=True)
    logs = args.out.parent / (args.out.stem + '-logs')
    logs.mkdir(exist_ok=True)
    pins = json.loads((ROOT / 'bench/pins.json').read_text())
    actual_tsrs = subprocess.check_output(['git','-C',str(args.tsrs_source),'rev-parse','HEAD'],text=True).strip()
    if actual_tsrs != pins['tsrs']['commit']:
        ap.error('tsrs source does not match bench/pins.json')
    result = {'schema': 1, 'date': datetime.datetime.now(datetime.timezone.utc).isoformat(), 'pins': pins,
              'head': subprocess.check_output(['git','rev-parse','HEAD'],cwd=ROOT,text=True).strip(),
              'machine': {'label':args.label,'platform':platform.platform(),'cpus':os.cpu_count()},
              'reps':args.reps, 'build':'plain release, Go -pgo=off unless binary metadata says otherwise',
              'tools': {}, 'rows': [], 'instructions': 'not collected; Go counts vary with GC; use bench/count.py on Linux separately'}
    for name,exe in tools.items():
        metadata = {'sha256':digest(exe),'version':subprocess.check_output([str(exe),'--version'],text=True).strip()}
        if name != 'tsrs':
            metadata['go_build'] = subprocess.check_output(['go','version','-m',str(exe)],text=True)
        result['tools'][name] = metadata
    failed = False
    report = ['# Compiler comparison', '', f"{result['date']} · {args.label} · {platform.platform()} · {os.cpu_count()} logical CPUs", '',
              'Native binaries; warm filesystem caches; medians of interleaved runs. Single mode caps both compiler and runtime pools. Default mode uses each compiler’s settings.', '',
              '| project | mode | stock s / MiB | fork s / MiB | tsrs s / MiB | stock/fork | Go equality | tsrs equality |',
              '| --- | --- | ---: | ---: | ---: | ---: | --- | --- |']
    for project_name, project in projects:
        for mode in args.modes.split(','):
            runs = {name:[] for name in tools}
            for name, exe in tools.items():
                warm = run(exe, project, mode, logs / f'{project_name}-{mode}-{name}-warm.log', args.timeout)
                if not warm['completed']:
                    failed = True
            order = list(tools)
            for rep in range(args.reps):
                rotated = order[rep % len(order):] + order[:rep % len(order)]
                if rep % 2:
                    rotated.reverse()
                for name in rotated:
                    r = run(tools[name], project, mode, logs / f'{project_name}-{mode}-{name}-{rep}.log', args.timeout)
                    runs[name].append(r)
                    if not r['completed']:
                        failed = True
            reference = signature(runs['stock'][0])
            equality = {name: all(r['completed'] and signature(r) == reference for r in rs) for name,rs in runs.items()}
            go_equal = all(equality[n] for n in tools if n != 'tsrs')
            failed |= not go_equal or (args.gate_tsrs and not equality['tsrs'])
            summary = {name:{key:statistics.median(r[key] for r in rs) for key in ('wall_s','peak_rss_bytes','cpu_s','errors')} for name,rs in runs.items()}
            row = {'project':project_name,'mode':mode,'equality':equality,'summary':summary,'runs':runs}
            result['rows'].append(row)
            def cell(name):
                r=summary[name]
                return f"{r['wall_s']:.3f} / {r['peak_rss_bytes']/2**20:.1f}"
            report.append(f"| {project_name} | {mode} | {cell('stock')} | {cell('fork')} | {cell('tsrs')} | {summary['stock']['wall_s']/summary['fork']['wall_s']:.2f}x | {'yes' if go_equal else 'MISMATCH'} | {'yes' if equality['tsrs'] else 'MISMATCH'} |")
            print(report[-1],flush=True)
    result['go_gate_passed'] = not failed
    args.out.write_text(json.dumps(result,indent=2)+'\n')
    args.out.with_suffix('.md').write_text('\n'.join(report)+'\n')
    sys.exit(1 if failed else 0)

if __name__ == '__main__':
    main()
