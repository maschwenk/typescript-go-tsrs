import copy
import json
import subprocess
import tempfile
import unittest
from pathlib import Path

from publish_readme import START, END, VSCODE, publish
from main_benchmark import MODES_MAIN

ROOT = Path(__file__).resolve().parent.parent


class PublicationTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.repo = Path(self.tmp.name)
        self.git('init', '-q')
        self.git('config', 'user.name', 'benchmark test')
        self.git('config', 'user.email', 'benchmark@example.invalid')
        (self.repo / 'README.md').write_text(f'prose before\n{START}\npending\n{END}\nprose after\n')
        self.git('add', '.')
        self.git('commit', '-qm', 'first source')
        self.first = self.git('rev-parse', 'HEAD')
        (self.repo / 'source.txt').write_text('second source')
        self.git('add', '.')
        self.git('commit', '-qm', 'second source')
        self.second = self.git('rev-parse', 'HEAD')
        self.result_path = self.repo / 'incoming.json'

    def git(self, *args):
        return subprocess.check_output(['git', *args], cwd=self.repo, text=True).strip()

    def result(self, sha):
        # Exercise publication with the recorded native VS Code measurements and diagnostics.
        recorded = json.loads((ROOT / 'docs/backports/results/native-vscode/matrix/compare.json').read_text())
        mapping = {'tsgo': 'fork', 'tsgo nightly': 'stock', 'tsrs': 'tsrs'}
        result = {'schema': 'readme-v1', 'completed': True, 'head': sha, 'date': recorded['date'],
                  'reps': recorded['reps'], 'machine': recorded['machine'], 'run_url': '',
                  'project': {'repo': 'https://github.com/microsoft/vscode.git', 'commit': VSCODE, 'project': 'src'},
                  'build': recorded['build'], 'tools': {}, 'rows': []}
        for name, old in mapping.items():
            result['tools'][name] = copy.deepcopy(recorded['tools'][old])
        result['tools']['tsgo']['commit'] = sha
        result['machine']['arch'] = 'arm64'
        result['tools']['tsgo nightly']['npm'] = {'package': '@typescript/typescript-darwin-arm64', 'version': 'recorded-stock'}
        result['tools']['tsrs']['commit'] = recorded['pins']['tsrs']['commit']
        for mode in MODES_MAIN:
            old = next(row for row in recorded['rows'] if row['mode'] == mode)
            result['rows'].append({'mode': mode, **{key: {name: copy.deepcopy(old[key][mapped]) for name, mapped in mapping.items()} for key in ('summary', 'runs', 'warmups', 'equality')}})
        return result

    def send(self, result):
        self.result_path.write_text(json.dumps(result))
        return publish(self.result_path, self.repo)

    def test_out_of_order_runs_keep_the_newer_table_and_archive_both(self):
        self.assertTrue(self.send(self.result(self.second)))
        newer = (self.repo / 'README.md').read_text()
        self.assertFalse(self.send(self.result(self.first)))
        self.assertEqual(newer, (self.repo / 'README.md').read_text())
        self.assertTrue((self.repo / f'bench/results/{self.first}.json').is_file())
        self.assertTrue((self.repo / f'bench/results/{self.second}.json').is_file())
        self.assertTrue(newer.startswith('prose before\n'))
        self.assertTrue(newer.endswith('\nprose after\n'))
        self.assertIn('| mode | tsgo | tsgo nightly | tsrs |', newer)

    def test_crashes_and_unstable_warmups_cannot_publish(self):
        for mutation in ('crash', 'warmup'):
            result = self.result(self.second)
            sample = result['rows'][0]['warmups']['tsgo']
            if mutation == 'crash':
                sample['completed'] = False
            else:
                sample['diagnostics_sha256'] = '0' * 64
            before = (self.repo / 'README.md').read_text()
            with self.assertRaises(ValueError):
                self.send(result)
            self.assertEqual(before, (self.repo / 'README.md').read_text())
            self.assertFalse((self.repo / 'bench/results').exists())

    def test_cross_version_disagreements_are_visible_without_speedup_claims(self):
        result = self.result(self.second)
        for row in result['rows']:
            for sample in [row['warmups']['tsgo nightly'], *row['runs']['tsgo nightly']]:
                sample['diagnostics_sha256'] = '1' * 64
        self.assertTrue(self.send(result))
        text = (self.repo / 'README.md').read_text()
        self.assertIn('| differs |', text)
        self.assertIn('no speedup ratio', text)

    def test_thread_dependent_diagnostics_and_invalid_metrics_cannot_publish(self):
        for mutation in ('thread-output', 'nan'):
            result = self.result(self.second)
            if mutation == 'nan':
                result['rows'][0]['runs']['tsrs'][0]['wall_s'] = float('nan')
            else:
                row = result['rows'][1]
                for sample in [row['warmups']['tsrs'], *row['runs']['tsrs']]:
                    sample['diagnostics_sha256'] = '2' * 64
            with self.assertRaises(ValueError):
                self.send(result)

    def test_older_rerun_cannot_replace_same_source_newer_measurements(self):
        result = self.result(self.second)
        self.assertTrue(self.send(result))
        before = (self.repo / 'README.md').read_text()
        result['date'] = '2020-01-01T00:00:00+00:00'
        self.assertFalse(self.send(result))
        self.assertEqual(before, (self.repo / 'README.md').read_text())


if __name__ == '__main__':
    unittest.main()
