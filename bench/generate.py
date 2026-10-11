#!/usr/bin/env python3
"""Deterministic compiler workloads; no dependencies or network installs."""
import json
from pathlib import Path


def generate(root):
    configs = []
    for name in ("references", "members", "unions", "flow", "modules", "errors"):
        d = root / name
        d.mkdir(parents=True, exist_ok=True)
        config = {"compilerOptions": {"strict": True, "noEmit": True, "skipLibCheck": True,
                  "target": "esnext", "module": "nodenext"}, "include": ["*.ts"]}
        (d / "tsconfig.json").write_text(json.dumps(config) + "\n")
        configs.append((name, d / "tsconfig.json"))
    (root / "references/main.ts").write_text(
        "interface Wide<T> {\n" + "".join(f"p{i}: T;\n" for i in range(2500)) + "}\n" +
        "".join(f"declare const v{i}: Wide<{{ n: {i} }}>; const r{i}: {i} = v{i}.p{i}.n;\n" for i in range(300)))
    (root / "members/main.ts").write_text(
        "interface Wide<T> {\n" + "".join(f"p{i}: T;\n" for i in range(2500)) + "}\n" +
        "type Copy<T> = { [K in keyof T]: T[K] };\n" +
        "".join(f"declare const v{i}: Copy<Wide<{{ n: {i} }}>>; const r{i}: {i} = v{i}.p{i}.n;\n" for i in range(300)))
    (root / "unions/main.ts").write_text(
        "".join(f"type U{i} = " + " | ".join(f"{{kind: 'k{j}'; value: {j}}}" for j in range(i, i+100)) + ";\n" for i in range(40)) +
        "type All = " + " | ".join(f"U{i}" for i in range(40)) + ";\n" +
        "type Picks = { [K in All['kind']]: Extract<All, {kind: K}> };\n" +
        "".join(f"declare const p{i}: Picks['k{i}']; const k{i}: 'k{i}' = p{i}.kind;\n" for i in range(100)))
    (root / "flow/main.ts").write_text("".join(
        f"function f{i}(v: string | number | null) {{\n" +
        "".join(f"if (v === null) v = {j}; else if (typeof v === 'string') v = v.length;\n" for j in range(80)) +
        "return v; }\n" for i in range(30)))
    for i in range(180):
        (root / f"modules/m{i}.ts").write_text(
            (f"import {{ make as prev }} from './m{i-1}.js';\n" if i else "") +
            "export function make<T>(v: T): { value: T } { return { value: v }; }\n" +
            (f"export const value = prev({{ id: {i}, label: 'm{i}' }});\n" if i else ""))
    (root / "errors/main.ts").write_text("const a: string = 123;\nconst b: number = 'wrong';\n")
    return configs

if __name__ == "__main__":
    import sys
    generate(Path(sys.argv[1]))
