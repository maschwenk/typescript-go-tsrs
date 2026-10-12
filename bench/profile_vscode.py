#!/usr/bin/env python3
"""Compare ordinary VS Code compilation, then collect CPU/allocation profiles separately."""
import argparse
import datetime
import importlib.util
import json
import os
import platform
import statistics
import subprocess
from pathlib import Path

from run import ROOT, digest, equal_outputs, run, signature


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--fork', type=Path, required=True)
    parser.add_argument('--tsrs', type=Path, required=True)
    parser.add_argument('--tsrs-source', type=Path, required=True)
    parser.add_argument('--tsrs-profile', type=Path, help='same pinned release source built with debug symbols')
    parser.add_argument('--perf', type=Path, help='Linux perf binary, with perf_event access configured')
    parser.add_argument('--samply', type=Path, help='native samply binary, as an alternative to perf')
    parser.add_argument('--work-dir', type=Path, default=ROOT / '.bench/work')
    parser.add_argument('--out', type=Path, default=ROOT / '.bench/profile')
    parser.add_argument('--reps', type=int, default=3)
    parser.add_argument('--gc-probe', action='store_true', help='also measure finite GOGC=400, GOMEMLIMIT=8GiB; not a default change')
    args = parser.parse_args()
    if not 1 <= args.reps <= 10 or (args.perf and args.samply):
        parser.error('use 1-10 repetitions and at most one native sampling tool')
    if any(os.environ.get(k) for k in ('GOGC', 'GOMEMLIMIT', 'GODEBUG')):
        parser.error('clear GC overrides for the control measurements')
    source = args.tsrs_source.resolve()
    pins = json.loads((ROOT / 'bench/pins.json').read_text())
    actual = subprocess.check_output(['git', '-C', str(source), 'rev-parse', 'HEAD'], text=True).strip()
    if actual != pins['tsrs']['commit']:
        parser.error('tsrs must match bench/pins.json')
    spec = importlib.util.spec_from_file_location('tsrs_profile_projects', source / 'bench/run.py')
    upstream = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(upstream)
    cfg = json.loads((source / 'bench/projects.json').read_text())
    project = next(p for p in cfg['projects'] if p['name'] == 'vscode')
    work = args.work_dir.resolve()
    if work.is_relative_to(ROOT / '.bench'):
        upstream.setup_project(cfg, project, work)
    actual = subprocess.check_output(['git', '-C', str(work / 'solutions/vscode'), 'rev-parse', 'HEAD'], text=True).strip()
    if actual != project['commit']:
        parser.error('VS Code source differs from the pin')
    _, config = upstream.project_path(cfg, project, work)
    executables = {'tsgo': args.fork.resolve(), 'tsrs': args.tsrs.resolve()}
    args.out = args.out.resolve()
    if (args.out / 'result.json').exists():
        parser.error('choose a new output directory for each profiling experiment')
    args.out.mkdir(parents=True, exist_ok=True)
    result = {'schema': 'profile-v1', 'measurement': 'ordinary-cli-v1',
              'date': datetime.datetime.now(datetime.timezone.utc).isoformat(),
              'head': subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=ROOT, text=True).strip(),
              'machine': {'platform': platform.platform(), 'cpus': os.cpu_count()},
              'project': project, 'pins': pins, 'reps': args.reps, 'tools': {}, 'rows': []}
    for name, binary in executables.items():
        result['tools'][name] = {'sha256': digest(binary), 'version': subprocess.check_output([str(binary), '--version'], text=True).strip()}
    result['tools']['tsgo']['go_build'] = subprocess.check_output(['go', 'version', '-m', str(args.fork.resolve())], text=True)
    for mode in ('single', 'checkers4'):
        pools = {'GOMAXPROCS': '1', 'RAYON_NUM_THREADS': '1'} if mode == 'single' else {'GOMAXPROCS': '4', 'RAYON_NUM_THREADS': '4'}
        validations = {name: run(binary, config, mode, args.out / f'{name}-{mode}-validation.log', 600, extra_env=pools)
                       for name, binary in executables.items()}
        if not equal_outputs(validations['tsgo'], list(validations.values())):
            raise RuntimeError('full diagnostic text, exit status or file count differs')
        if result['rows'] and any(signature(v) != signature(result['rows'][0]['validations'][name]) for name, v in validations.items()):
            raise RuntimeError('diagnostics differ between thread counts')
        measurements = {name: [] for name in executables}
        if args.gc_probe:
            measurements['tsgo-gc-probe'] = []
        for repetition in range(args.reps):
            order = list(measurements)
            if repetition % 2:
                order.reverse()
            for name in order:
                tool = 'tsgo' if name == 'tsgo-gc-probe' else name
                env = {**pools, **({'GOGC': '400', 'GOMEMLIMIT': '8GiB'} if name == 'tsgo-gc-probe' else {})}
                measurements[name].append(run(executables[tool], config, mode, args.out / f'{name}-{mode}-{repetition}.log', 600,
                                              validation=validations[tool], extra_env=env))
        if any(not equal_outputs(validations['tsgo'], samples) for samples in measurements.values()):
            raise RuntimeError('timed ordinary CLI output differs from the validated compile')
        summary = {name: {k: statistics.median(r[k] for r in rs) for k in ('wall_s', 'cpu_s', 'peak_rss_bytes')}
                   for name, rs in measurements.items()}
        profiles = args.out / mode
        profiles.mkdir(exist_ok=True)
        for name in ('tsgo', 'tsgo-gc-probe') if args.gc_probe else ('tsgo',):
            directory = profiles / name
            directory.mkdir(exist_ok=True)
            env = {**pools, 'GODEBUG': 'gctrace=1', **({'GOGC': '400', 'GOMEMLIMIT': '8GiB'} if name == 'tsgo-gc-probe' else {})}
            sample = run(executables['tsgo'], config, mode, directory / 'compile.log', 600, validation=validations['tsgo'],
                         extra_args=['--pprofDir', str(directory)], extra_env=env)
            if not equal_outputs(validations['tsgo'], [sample]) or len(list(directory.glob('*-cpuprofile.pb.gz'))) != 1:
                raise RuntimeError('Go profile did not cover the validated compilation')
        if args.perf or args.samply:
            binary = (args.tsrs_profile or args.tsrs).resolve()
            if subprocess.check_output([str(binary), '--version'], text=True) != subprocess.check_output([str(args.tsrs.resolve()), '--version'], text=True):
                raise RuntimeError('profiling tsrs version differs from the timing binary')
            result['tools']['tsrs-profile'] = {'sha256': digest(binary), 'purpose': 'profiling only; excluded from timing summaries'}
            if args.perf:
                prefix = [args.perf.resolve(), 'record', '-e', 'cpu-clock', '-F', '499', '--call-graph', 'dwarf,8192', '-o', profiles / 'tsrs.perf.data', '--']
            else:
                prefix = [args.samply.resolve(), 'record', '--save-only', '--unstable-presymbolicate', '--rate', '2000', '-o', profiles / 'tsrs.json.gz', '--']
            sample = run(binary, config, mode, profiles / 'tsrs-profile.log', 600, validation=validations['tsrs'], extra_env=pools, command_prefix=prefix)
            if not equal_outputs(validations['tsrs'], [sample]):
                raise RuntimeError('tsrs profile did not cover the validated compilation')
            if args.perf:
                with (profiles / 'tsrs-perf-report.txt').open('w') as report:
                    subprocess.run([str(args.perf.resolve()), 'report', '--stdio', '--no-children', '-i', str(profiles / 'tsrs.perf.data')], stdout=report, check=True)
        result['rows'].append({'mode': mode, 'validations': validations, 'runs': measurements, 'summary': summary})
        (args.out / 'result.json').write_text(json.dumps(result, indent=2) + '\n')
        print(mode, summary, flush=True)
    result['completed'] = True
    (args.out / 'result.json').write_text(json.dumps(result, indent=2) + '\n')


if __name__ == '__main__':
    main()
