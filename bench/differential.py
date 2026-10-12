#!/usr/bin/env python3
"""Seeded native-compiler differential testing of diagnostics and emitted bytes."""
import argparse
import hashlib
import json
import os
import random
import subprocess
from pathlib import Path

MODES = {'single': ['--singleThreaded'], 'checkers4': ['--checkers', '4']}

def fixture(seed):
    r=random.Random(seed)
    n=r.randrange(3,40)
    lines=[f"type Event = "+' | '.join(f"{{ kind: 'k{i}'; value: {i}; extra?: string }}" for i in range(n))+';',
           "type Copy<T> = { readonly [P in keyof T]: T[P] };",
           "type Select<K extends Event['kind']> = Extract<Event, {kind: K}>;",
           "declare function event<K extends Event['kind']>(kind: K, data: Select<K>['value']): Select<K>;",
           "export const events: Event[] = [" + ','.join(f"event('k{i}', {i if r.random()>.15 else n+1})" for i in range(n))+"];",
           "export interface Box<T> { value: T; nested: Copy<T>; tuple: readonly [T, T?] }",
           "export class Base<T> { value!: T; optional?: T; method(v: T): T {return v} }",
           "export class Derived extends Base<Event> { readonly tag = 'derived' as const }",
           "export type Keys = keyof Copy<Box<Event>>;",
           "export type Values = Box<Event>[Keys];",
           "export function narrow(v: Event | undefined): number {",
           "if (!v) return -1;",
           "let result: number = v.value;",
           "for (let i=0; i<2; i++) { if (v.kind === 'k0') result += v.value; }",
           "try { if (v.kind === 'k1') return v.value; } finally { result += 1; }",
           "return result; }",
           f"export const expected: {r.randrange(n)} = events[0].value;",
           "import { identity } from './module.js';",
           "export const shared = identity(events);"]
    return {'main.ts':'\n'.join(lines)+'\n', 'module.ts':'export function identity<T>(v:T):T {return v}\n'}

def emitted(directory):
    return {str(p.relative_to(directory)):hashlib.sha256(p.read_bytes()).hexdigest()
            for p in sorted(directory.rglob('*')) if p.is_file()}

def compile(exe, project, out, mode, log, timeout, shadow):
    env=dict(os.environ)
    env.pop('GOMAXPROCS',None)
    if mode=='single': env['GOMAXPROCS']='1'
    if shadow:
        for name in ('UNION_CACHE','IA_MEMO','FLOW_MEMO','INFER_MEMO','LAZY_UNAFFECTED','CONTEXTUAL_REUSE','CONTEXTUAL_POSITIONS','DEEP_NESTED','NODE_FILE_CACHE','DISCRIMINANT_PROPERTIES','PROPERTY_LOOKUPS','SUBTYPE_PROPERTIES','CONSTITUENT_FLAGS'): env['TSGO_'+name]='shadow'
    cmd=[str(exe),'-p',str(project),'--pretty','false','--incremental','false',
         '--declaration','--outDir',str(out),*MODES[mode]]
    try:
        result=subprocess.run(cmd,cwd=project.parent,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=timeout)
    except subprocess.TimeoutExpired as exc:
        log.write_bytes(exc.stdout or b'')
        return {'exit':'timeout','diagnostics':'','emit':{}}
    log.write_bytes(result.stdout)
    text=result.stdout.replace(str(out).encode(),b'<out>')
    return {'exit':result.returncode,'diagnostics':hashlib.sha256(text).hexdigest(),'emit':emitted(out)}

def main():
    ap=argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--stock',type=Path,required=True)
    ap.add_argument('--fork',type=Path,required=True)
    ap.add_argument('--seeds',type=int,default=64)
    ap.add_argument('--start',type=int,default=0)
    ap.add_argument('--modes',default='single,checkers4')
    ap.add_argument('--timeout',type=int,default=30)
    ap.add_argument('--work-dir',type=Path,default=Path('.bench/differential'))
    ap.add_argument('--shadow',action='store_true')
    args=ap.parse_args()
    args.stock=args.stock.resolve();args.fork=args.fork.resolve()
    args.work_dir=args.work_dir.resolve()
    if args.seeds<1 or not set(args.modes.split(','))<=MODES.keys(): ap.error('positive seed count and supported modes required')
    report={'schema':1,'source':subprocess.check_output(['git','rev-parse','HEAD'],text=True).strip(),'tools':{name:hashlib.sha256(exe.read_bytes()).hexdigest() for name,exe in [('stock',args.stock),('fork',args.fork)]},'seeds':args.seeds,'start':args.start,'modes':args.modes,'shadow':args.shadow,'failures':[]}
    for seed in range(args.start,args.start+args.seeds):
        d=args.work_dir/f'seed-{seed}'
        src=d/'src';src.mkdir(parents=True,exist_ok=True)
        for name,text in fixture(seed).items(): (src/name).write_text(text)
        config={'compilerOptions':{'strict':True,'target':'esnext','module':'nodenext','skipLibCheck':True,
                'declaration':True,'noEmitOnError':False},'include':['*.ts']}
        project=src/'tsconfig.json';project.write_text(json.dumps(config)+'\n')
        for mode in args.modes.split(','):
            results={}
            for name,exe in [('stock',args.stock),('fork',args.fork)]:
                out=d/f'{mode}-{name}-out';out.mkdir(exist_ok=True)
                for old in out.rglob('*'):
                    if old.is_file(): old.unlink()
                results[name]=compile(exe,project,out,mode,d/f'{mode}-{name}.log',args.timeout,args.shadow and name=='fork')
            if results['stock']!=results['fork'] or results['stock']['exit'] not in (0,1,2) or not {'main.js','main.d.ts','module.js','module.d.ts'}<=results['stock']['emit'].keys():
                report['failures'].append({'seed':seed,'mode':mode,'results':results})
    args.work_dir.mkdir(parents=True,exist_ok=True)
    (args.work_dir/'report.json').write_text(json.dumps(report,indent=2)+'\n')
    print(json.dumps({'seeds':args.seeds,'modes':args.modes,'failures':len(report['failures']),'report':str(args.work_dir/'report.json')}))
    return bool(report['failures'])

if __name__=='__main__': raise SystemExit(main())
