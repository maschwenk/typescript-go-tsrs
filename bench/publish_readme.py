#!/usr/bin/env python3
"""Archive a complete VS Code run and update the README without regressing its source commit."""
import argparse
import datetime
import json
import math
import re
import statistics
import subprocess
from pathlib import Path

from main_benchmark import TOOLS, MODES_MAIN, stable
from run import signature

START = '<!-- vscode-benchmarks:start -->'
END = '<!-- vscode-benchmarks:end -->'
REPO = 'maschwenk/typescript-go-tsrs'
VSCODE = '9cf0128b9822dcd8a533f08ef7312956a4390301'


def validate(result):
    if result.get('schema') != 'readme-v1' or result.get('completed') is not True:
        raise ValueError('a complete README benchmark is required')
    if not re.fullmatch('[0-9a-f]{40}', result['head']):
        raise ValueError('a full source commit is required')
    if not isinstance(result['reps'], int) or not 1 <= result['reps'] <= 20:
        raise ValueError('invalid repetition count')
    date = datetime.datetime.fromisoformat(result['date'])
    if date.tzinfo is None:
        raise ValueError('measurement date must include its timezone')
    project = result['project']
    if (project['repo'], project['commit'], project['project']) != ('https://github.com/microsoft/vscode.git', VSCODE, 'src'):
        raise ValueError('only the pinned VS Code src workload can be published')
    if result['run_url'] and not re.fullmatch(f'https://github.com/{REPO}/actions/runs/[0-9]+', result['run_url']):
        raise ValueError('unexpected Actions run link')
    if set(result['tools']) != set(TOOLS) or result['tools']['tsgo']['commit'] != result['head']:
        raise ValueError('unexpected compilers or source provenance')
    for metadata in result['tools'].values():
        if not re.fullmatch('[0-9a-f]{64}', metadata['sha256']):
            raise ValueError('missing binary hash')
    if [row['mode'] for row in result['rows']] != list(MODES_MAIN):
        raise ValueError('all four thread modes are required')
    for row in result['rows']:
        if set(row['runs']) != set(TOOLS) or set(row['warmups']) != set(TOOLS):
            raise ValueError('missing compiler output')
        for name in TOOLS:
            runs = row['runs'][name]
            warmup = row['warmups'][name]
            if len(runs) != result['reps'] or not stable(warmup, runs):
                raise ValueError('incomplete or unstable compiler output')
            if signature(warmup) != signature(result['rows'][0]['warmups'][name]):
                raise ValueError('compiler output changed with its thread count')
            for sample in [warmup, *runs]:
                if sample['counters'].get('Files', 0) < 1000:
                    raise ValueError('VS Code checking did not cover the workload')
                if not re.fullmatch('[0-9a-f]{64}', sample['diagnostics_sha256']):
                    raise ValueError('invalid diagnostic hash')
                for key in ('wall_s', 'peak_rss_bytes', 'cpu_s', 'errors'):
                    value = sample[key]
                    if not isinstance(value, (int, float)) or not math.isfinite(value) or value < 0:
                        raise ValueError('invalid measurement')
                    if key in ('wall_s', 'peak_rss_bytes') and value == 0:
                        raise ValueError('empty measurement')
            # Derive the published numbers and parity from raw runs, not an uploaded summary.
            row['summary'][name] = {key: statistics.median(s[key] for s in runs) for key in ('wall_s', 'peak_rss_bytes', 'cpu_s', 'errors')}
            row['equality'][name] = stable(row['warmups']['tsgo'], [warmup, *runs])


def text(value):
    return str(value).replace('|', '\\|').replace('\n', ' ').replace('\r', ' ').replace('`', "'").replace('<', '&lt;').replace('>', '&gt;')


def table(result):
    sha = result['head']
    link = f'https://github.com/{REPO}/commit/{sha}'
    machine = result['machine']
    lines = ['## VS Code benchmarks', '', f'<!-- measured-commit: {sha} -->',
             f"Measured [{sha[:12]}]({link}) on {text(result['date'])}. " + (f"[Actions run]({result['run_url']}) · " if result['run_url'] else '') + f'[raw measurements](bench/results/{sha}.json)', '',
             f"{text(machine['platform'])} · {text(machine['arch'])} · {text(machine['cpus'])} logical CPUs · {text(machine.get('cpu_model', 'CPU model unavailable'))} · {text(machine.get('mem_total', 'RAM unavailable'))}", '',
             f"VS Code `{VSCODE[:12]}` (`src`), warm caches, one warmup and {result['reps']} interleaved repetitions per compiler and mode. Tables show medians. Hosted runner hardware can vary between commits.", '',
             '### Wall time (seconds)', '', '| mode | tsgo | tsgo nightly | tsrs |', '| --- | ---: | ---: | ---: |']
    labels = {'default': 'default settings', 'single': 'single thread', 'checkers4': '4 checkers', 'checkers8': '8 checkers'}
    for row in result['rows']:
        lines.append('| ' + labels[row['mode']] + ' | ' + ' | '.join(f"{row['summary'][name]['wall_s']:.3f}" for name in TOOLS) + ' |')
    lines += ['', '### Peak memory (MiB)', '', '| mode | tsgo | tsgo nightly | tsrs |', '| --- | ---: | ---: | ---: |']
    for row in result['rows']:
        lines.append('| ' + labels[row['mode']] + ' | ' + ' | '.join(f"{row['summary'][name]['peak_rss_bytes'] / 2**20:.1f}" for name in TOOLS) + ' |')
    default = result['rows'][0]
    lines += ['', '### Compiler versions and diagnostics', '', '| compiler | version | files | errors | matches tsgo diagnostics, exit and files |', '| --- | --- | ---: | ---: | --- |']
    for name in TOOLS:
        sample = default['warmups'][name]
        equal = all(row['equality'][name] for row in result['rows'])
        lines.append(f"| {name} | {text(result['tools'][name]['version'])} | {sample['counters']['Files']} | {sample['errors']} | {'yes' if equal else 'differs'} |")
    nightly = result['tools']['tsgo nightly']['npm']
    lines += ['', f"`tsgo` is this fork. `tsgo nightly` is `{text(nightly['package'])}@{text(nightly['version'])}`, resolved once from `typescript@next`. `tsrs` is pinned to `{text(result['tools']['tsrs']['commit'])[:12]}`.", '',
              text(result['build']) + '. Default settings use each compiler’s own thread policy; the other rows compare equal checker counts. Single thread also limits runtime pools to one.', '',
              'All warmups and repetitions completed and produced stable full diagnostics within each compiler. Different compiler versions can produce different diagnostics; those rows are descriptive timings, with no speedup ratio or correctness claim. Compiler correctness is gated separately by CI.', '',
              '[Method and reproduction](bench/README.md#automatic-readme-benchmarks) · [Native Mac comparison](docs/backports/results/native-vscode/README.md)']
    return '\n'.join(lines)


def ancestor(repo, older, newer):
    return subprocess.run(['git', 'merge-base', '--is-ancestor', older, newer], cwd=repo, check=False).returncode == 0


def publish(result_path, repo):
    result = json.loads(result_path.read_text())
    validate(result)
    head = result['head']
    if not ancestor(repo, head, 'HEAD'):
        raise ValueError('measured commit is not an ancestor of the publication checkout')
    readme = repo / 'README.md'
    content = readme.read_text()
    if content.count(START) != 1 or content.count(END) != 1:
        raise ValueError('README must contain exactly one benchmark marker pair')
    before, rest = content.split(START, 1)
    current, after = rest.split(END, 1)
    displayed = re.search(r'<!-- measured-commit: ([0-9a-f]{40}) -->', current)
    archive = repo / 'bench/results' / f'{head}.json'
    if archive.exists() and json.loads(archive.read_text())['date'] > result['date']:
        print('keeping the newer repetition of this commit')
        return False
    block = table(result)
    keep_newer = False
    if displayed and displayed[1] != head:
        if ancestor(repo, head, displayed[1]):
            keep_newer = True
        elif not ancestor(repo, displayed[1], head):
            raise ValueError('benchmark sources are on unrelated histories')
    archive.parent.mkdir(parents=True, exist_ok=True)
    archive.write_text(json.dumps(result, indent=2) + '\n')
    if keep_newer:
        print('archived older run; keeping the newer README table')
        return False
    readme.write_text(before + START + '\n' + block + '\n' + END + after)
    print('updated README for ' + head)
    return True


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('result', type=Path)
    parser.add_argument('--repo', type=Path, default=Path.cwd())
    args = parser.parse_args()
    publish(args.result, args.repo.resolve())
