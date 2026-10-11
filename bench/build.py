#!/usr/bin/env python3
"""Build the fork and pinned stock Go compiler with the same toolchain and flags."""
import argparse
import json
import os
import subprocess
from pathlib import Path
ROOT = Path(__file__).resolve().parent.parent

def sh(argv, cwd=ROOT):
    subprocess.run(argv, cwd=cwd, check=True, env=dict(os.environ, GOTOOLCHAIN='auto'))

ap = argparse.ArgumentParser(description=__doc__)
ap.add_argument('--out', type=Path, default=ROOT/'.bench/bin')
ap.add_argument('--reference', help='optional parent revision, built alongside stock')
args=ap.parse_args()
args.out=args.out.resolve()
args.out.mkdir(parents=True,exist_ok=True)
pins=json.loads((ROOT/'bench/pins.json').read_text())
for name, ref in [('stock',pins['upstream']['commit']),('parent',args.reference)]:
    if not ref:
        continue
    directory=ROOT/'.bench'/f'{name}-source'
    if not directory.exists():
        sh(['git','worktree','add','--detach',str(directory),ref])
    elif subprocess.check_output(['git','rev-parse','HEAD'],cwd=directory,text=True).strip()!=ref:
        sh(['git','checkout','--detach',ref],directory)
    sh(['go','build','-pgo=off','-o',str(args.out/name),'./cmd/tsc'],directory/'tsc')
sh(['go','build','-pgo=off','-o',str(args.out/'fork'),'./cmd/tsc'],ROOT/'tsc')
