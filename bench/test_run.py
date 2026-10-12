import tempfile
import unittest
from pathlib import Path
from run import run, signature, parse_output, equal_outputs

class HarnessTest(unittest.TestCase):
    def test_ordinary_cli_requires_successful_coverage_and_exact_diagnostics(self):
        validation = parse_output(['a.ts(1,1): error TS2322: wrong', '  detail', 'Files: 10465', 'Check time: 1s'], 2)
        lines = ['a.ts(1,1): error TS2322: wrong', '  detail']
        ordinary = parse_output(lines, 2, validation=validation)
        self.assertTrue(ordinary['completed'])
        self.assertEqual(signature(ordinary), signature(validation))
        self.assertEqual(ordinary['counters_source'], 'validation')
        self.assertFalse(parse_output(lines, -9, validation=validation)['completed'])
        self.assertFalse(parse_output(lines + ['panic: crash'], 2, validation=validation)['completed'])
        self.assertFalse(parse_output(lines[:1], 2, validation=validation)['completed'])
        self.assertFalse(parse_output(lines, 2, validation=parse_output([], 2))['completed'])

    def test_warmup_mismatch_is_part_of_equality(self):
        reference=parse_output(['Files: 2', 'Check time: 0.01s'], 0)
        warmup=parse_output(['file.ts(1,1): error TS2322: wrong', 'Files: 2', 'Check time: 0.01s'], 2)
        self.assertFalse(equal_outputs(reference, [warmup, reference]))

    def test_panic_after_metrics_and_missing_file_count_are_incomplete(self):
        self.assertFalse(parse_output(['Files: 2', 'Check time: 0.01s', 'panic: late failure'], 2)['completed'])
        self.assertFalse(parse_output(['Check time: 0.01s'], 0)['completed'])

    def test_errors_and_continuations_are_part_of_equality(self):
        with tempfile.TemporaryDirectory() as tmp:
            root=Path(tmp)
            exe=root/'compiler'
            exe.write_text('#!/bin/sh\nprintf "main.ts(1,1): error TS2322: wrong type\\n  detail A\\nFiles: 2\\nMemory used: 25K\\nCheck time: 0.010s\\n"\nexit 2\n')
            exe.chmod(0o755)
            project=root/'tsconfig.json'
            project.write_text('{}')
            a=run(exe,project,'single',root/'a.log',5)
            self.assertTrue(a['completed'])
            self.assertEqual(a['errors'],1)
            exe.write_text(exe.read_text().replace('detail A','detail B'))
            b=run(exe,project,'single',root/'b.log',5)
            self.assertNotEqual(signature(a),signature(b))

    def test_timeout_and_crash_are_incomplete(self):
        with tempfile.TemporaryDirectory() as tmp:
            root=Path(tmp)
            exe=root/'compiler'
            exe.write_text('#!/bin/sh\nexec sleep 5\n')
            exe.chmod(0o755)
            project=root/'tsconfig.json'
            project.write_text('{}')
            r=run(exe,project,'single',root/'a.log',0.1)
            self.assertFalse(r['completed'])
            self.assertLess(r['exit'],0)
