#!/usr/bin/env python3
"""Train native Go PGO on pinned projects and build a separate optimized binary."""
import argparse
import hashlib
import importlib.util
import json
import os
import platform
import subprocess
from pathlib import Path
from run import parse_output, signature

ROOT=Path(__file__).resolve().parent.parent

def sha(p): return hashlib.sha256(p.read_bytes()).hexdigest()

def main():
    ap=argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--compiler',type=Path,required=True)
    ap.add_argument('--tsrs-source',type=Path,required=True)
    ap.add_argument('--work-dir',type=Path,default=ROOT/'.bench/work')
    ap.add_argument('--out',type=Path,default=ROOT/'.bench/pgo')
    ap.add_argument('--projects',default='vscode')
    ap.add_argument('--reps',type=int,default=10)
    args=ap.parse_args()
    args.compiler=args.compiler.resolve();args.out=args.out.resolve();args.work_dir=args.work_dir.resolve()
    if args.reps<1: ap.error('positive repetitions required')
    if args.projects != 'vscode': ap.error('this experiment uses VS Code only')
    pins=json.loads((ROOT/'bench/pins.json').read_text())
    actual=subprocess.check_output(['git','-C',str(args.tsrs_source),'rev-parse','HEAD'],text=True).strip()
    if actual!=pins['tsrs']['commit']: ap.error('tsrs project definitions do not match the source pin')
    cfg=json.loads((args.tsrs_source/'bench/projects.json').read_text())
    available={p['name']:p for p in cfg['projects']}
    spec=importlib.util.spec_from_file_location('tsrs_bench',args.tsrs_source/'bench/run.py')
    upstream=importlib.util.module_from_spec(spec);spec.loader.exec_module(upstream)
    args.out.mkdir(parents=True,exist_ok=True)
    profiles=[];training=[];outputs=[];reference=None
    for name in args.projects.split(','):
        if name not in available: ap.error(f'unknown training project {name}')
        p=available[name]
        if args.work_dir.is_relative_to(ROOT/'.bench'): upstream.setup_project(cfg,p,args.work_dir)
        checkout=args.work_dir/'suite' if 'suite_dir' in p else args.work_dir/'solutions'/name
        expected=cfg['suite']['commit'] if 'suite_dir' in p else p['commit']
        actual=subprocess.check_output(['git','-C',str(checkout),'rev-parse','HEAD'],text=True).strip()
        if actual!=expected: ap.error(f'{name}: corpus revision differs from pin')
        _,project=upstream.project_path(cfg,p,args.work_dir)
        training.append({'name':name,'commit':actual,'project':str(project)})
        for i in range(args.reps):
            d=args.out/'raw'/name/str(i);d.mkdir(parents=True,exist_ok=True)
            for old in d.glob('*-cpuprofile.pb.gz'): old.unlink()
            cmd=[str(args.compiler),'-p',str(project),'--noEmit','--incremental','false','--pretty','false','--singleThreaded','--extendedDiagnostics','--pprofDir',str(d)]
            result=subprocess.run(cmd,cwd=project.parent,env=dict(os.environ,GOMAXPROCS='1'),stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=600)
            (d/'compile.log').write_bytes(result.stdout)
            cpus=list(d.glob('*-cpuprofile.pb.gz'))
            if result.returncode not in (0,1,2) or b'Check time:' not in result.stdout or len(cpus)!=1 or cpus[0].stat().st_size<100:
                raise RuntimeError(f'{name} repetition {i}: no completed compile and CPU profile; see {d}')
            parsed=parse_output(result.stdout.decode(errors='replace').splitlines(),result.returncode)
            if not parsed['completed']: raise RuntimeError(f'{name} repetition {i}: incomplete compiler output')
            if reference is None: reference=signature(parsed)
            if signature(parsed)!=reference: raise RuntimeError(f'{name} repetition {i}: unstable diagnostics or file count')
            outputs.append({'project':name,'repetition':i,**parsed,'profile_sha256':sha(cpus[0])})
            profiles.extend(cpus)
    profile=args.out/'compiler.pprof'
    env=dict(os.environ,GOTOOLCHAIN='auto')
    subprocess.run(['go','-C',str(ROOT/'tsc'),'tool','pprof','-proto','-output',str(profile),*map(str,profiles)],env=env,check=True)
    summary=subprocess.check_output(['go','-C',str(ROOT/'tsc'),'tool','pprof','-top',str(profile)],env=env,text=True)
    if 'Total samples = 0' in summary: raise RuntimeError('empty merged CPU profile')
    (args.out/'profile-top.txt').write_text(summary)
    binary=args.out/'tsgo-pgo'
    subprocess.run(['go','-C',str(ROOT/'tsc'),'build','-pgo='+str(profile),'-o',str(binary),'./cmd/tsc'],env=env,check=True)
    metadata={'schema':1,'source':subprocess.check_output(['git','rev-parse','HEAD'],cwd=ROOT,text=True).strip(),
              'platform':platform.platform(),'architecture':platform.machine(),'profile_sha256':sha(profile),
              'training_compiler_sha256':sha(args.compiler),'pgo_binary_sha256':sha(binary),
              'training':training,'outputs':outputs,'repetitions':args.reps,'profiles':len(profiles),
              'evaluation_scope':'VS Code only; training and evaluation use the same pinned workload, not held-out generalization',
              'go_build':subprocess.check_output(['go','version','-m',str(binary)],text=True)}
    (args.out/'metadata.json').write_text(json.dumps(metadata,indent=2)+'\n')
    print(json.dumps({'binary':str(binary),'profile':str(profile),'metadata':str(args.out/'metadata.json')}))

if __name__=='__main__': main()
