"""File-based regressions for the correctness gate; no compiler benchmark is run."""
import contextlib
import io
import json
import tempfile
import unittest
from pathlib import Path

import codspeed_vscode as integration


class EqualityGateTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.state = Path(self.temp.name)
        (self.state / 'logs').mkdir()
        (self.state / 'profile/results').mkdir(parents=True)
        tools = {}
        for tool in integration.TOOLS:
            binary = self.state / tool
            binary.write_bytes(b'fixture binary identity')
            tools[tool] = {'path': str(binary), 'sha256': integration.sha256(binary)}
            for rep in range(2):
                self.log(tool, rep).write_text('src/example.ts(1,1): error TS2322: Type mismatch.\n'
                                              '  Complete continuation text.\n\nFiles: 8\n'
                                              'Types: 99\nCheck time: 0.20s\nTotal time: 0.25s\n')
                self.log(tool, rep).with_suffix('.log.exit').write_text('1\n')
        self.manifest = {'repo': str(integration.ROOT), 'modes': ['single'], 'repetitions': 1,
                         'tools': tools, 'harness_sha256': integration.sha256(integration.ROOT / 'bench/run.py')}
        (self.state / 'manifest.json').write_text(json.dumps(self.manifest))
        results = {'instrument': {'type': 'walltime'}, 'benchmarks': [
            {'name': f'vscode/src/single/{tool}', 'stats': {'rounds': 1, 'warmup_iters': 0,
                                                        'median_ns': (index + 1) * 1000000}}
            for index, tool in enumerate(integration.TOOLS)]}
        (self.state / 'profile/results/fixture.json').write_text(json.dumps(results))

    def log(self, tool, rep=0):
        return self.state / 'logs' / f'{tool}-single.{rep}.log'

    def gate(self, strict=False):
        with contextlib.redirect_stdout(io.StringIO()):
            return integration.gate(self.state, strict)

    def test_equal_nonzero_diagnostic_exits_pass(self):
        self.assertTrue(self.gate())

    def test_continuation_difference_fails_and_suppresses_speed_ratio(self):
        path = self.log('fork', 1)
        path.write_text(path.read_text().replace('Complete continuation text.', 'Different continuation text.'))
        self.assertFalse(self.gate())
        self.assertIn('invalid', (self.state / 'compare.md').read_text())

    def test_equal_diagnostics_with_different_exit_fail(self):
        self.log('fork', 1).with_suffix('.log.exit').write_text('2\n')
        self.assertFalse(self.gate())

    def test_unstable_stock_suppresses_both_speed_ratios(self):
        path = self.log('stock', 1)
        path.write_text(path.read_text().replace('TS2322', 'TS9999'))
        self.assertFalse(self.gate())
        row = next(line for line in (self.state / 'compare.md').read_text().splitlines()
                   if line.startswith('| single |'))
        self.assertEqual(row.count('invalid'), 2)

    def test_file_count_difference_fails(self):
        path = self.log('fork', 1)
        path.write_text(path.read_text().replace('Files: 8', 'Files: 9'))
        self.assertFalse(self.gate())

    def test_types_counter_difference_does_not_fail(self):
        path = self.log('fork', 1)
        path.write_text(path.read_text().replace('Types: 99', 'Types: 10'))
        self.assertTrue(self.gate())

    def test_tsrs_disagreement_is_advisory_or_explicitly_strict(self):
        path = self.log('tsrs', 1)
        path.write_text(path.read_text().replace('TS2322', 'TS9999'))
        self.assertTrue(self.gate())
        self.assertFalse(self.gate(strict=True))

    def test_incomplete_tsrs_fails_even_when_advisory(self):
        self.log('tsrs', 1).write_text('Files: 8\n')
        self.assertFalse(self.gate())

    def test_missing_repetition_fails(self):
        self.log('stock', 1).unlink()
        self.assertFalse(self.gate())

    def test_binary_mutation_refuses_results(self):
        (self.state / 'fork').write_bytes(b'changed binary')
        with self.assertRaisesRegex(ValueError, 'binary changed'):
            self.gate()

    def test_missing_exit_sidecar_refuses_results(self):
        self.log('fork', 1).with_suffix('.log.exit').unlink()
        with self.assertRaises(FileNotFoundError):
            self.gate()


if __name__ == '__main__':
    unittest.main()
