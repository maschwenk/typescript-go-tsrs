#!/usr/bin/env python3
"""Resolve typescript@next once, install that exact native binary and record provenance."""
import argparse
import json
import platform
import subprocess
import sys
from pathlib import Path


def npm_json(*args):
    return json.loads(subprocess.check_output(['npm', 'view', *args, '--json'], text=True))


def install(directory):
    directory.mkdir(parents=True, exist_ok=True)
    version = npm_json('typescript@next', 'version')
    arch = {'x86_64': 'x64', 'aarch64': 'arm64', 'arm64': 'arm64'}[platform.machine().lower()]
    package = f'@typescript/typescript-{sys.platform}-{arch}'
    metadata = npm_json(f'{package}@{version}')
    if metadata['version'] != version:
        raise ValueError('nightly platform package version differs')
    (directory / 'package.json').write_text('{"name":"tsgo-nightly-benchmark","private":true}\n')
    subprocess.run(['npm', 'install', '--save-exact', '--ignore-scripts', '--no-audit', '--no-fund',
                    f'--min-release-age-exclude={package}', f'{package}@{version}'], cwd=directory, check=True)
    binary = directory / 'node_modules' / package / 'lib/tsc'
    if binary.read_bytes()[:4] not in (b'\x7fELF', b'\xcf\xfa\xed\xfe', b'\xca\xfe\xba\xbe'):
        raise ValueError('nightly must be a native executable')
    actual = subprocess.check_output([str(binary.resolve()), '--version'], text=True).strip()
    if actual != f'Version {version}':
        raise ValueError(f'nightly version differs: {actual}')
    record = {'package': package, 'version': version, 'gitHead': metadata.get('gitHead'),
              'integrity': metadata['dist']['integrity'], 'resolved_from': 'typescript@next'}
    (directory / 'provenance.json').write_text(json.dumps(record, indent=2) + '\n')
    print(binary.resolve())


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--out', type=Path, default=Path('.bench/nightly'))
    install(parser.parse_args().out.resolve())
