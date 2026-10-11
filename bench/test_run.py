import tempfile
import unittest
from pathlib import Path
from run import run, signature

class HarnessTest(unittest.TestCase):
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
