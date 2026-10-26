import unittest

from check_coverage import evaluate


class CoverageGateTests(unittest.TestCase):
    def test_exact_threshold_passes_and_lower_coverage_fails(self):
        report = 'mode: atomic\nipapp/main.go:1.1,2.2 1 1\nipapp/main.go:3.1,4.2 1 0\n'
        self.assertEqual(evaluate(report, {'minimumStatementPercent': 50, 'packages': {}})['overall']['percent'], 50)
        with self.assertRaisesRegex(ValueError, 'Coverage regression'):
            evaluate(report, {'minimumStatementPercent': 51, 'packages': {}})

    def test_missing_required_package_fails(self):
        report = 'mode: count\nipapp/main.go:1.1,2.2 1 1\n'
        with self.assertRaisesRegex(ValueError, 'No statement coverage'):
            evaluate(report, {'minimumStatementPercent': 0, 'packages': {'ipapp/packages/ip/': 100}})

    def test_package_floor_cannot_be_hidden_by_overall_coverage(self):
        report = 'mode: set\nipapp/main.go:1.1,2.2 100 1\nipapp/packages/ip/ip.go:1.1,2.2 1 0\n'
        with self.assertRaisesRegex(ValueError, 'Coverage regression for ipapp/packages/ip/'):
            evaluate(report, {'minimumStatementPercent': 90, 'packages': {'ipapp/packages/ip/': 100}})

    def test_duplicate_blocks_merge_hits_without_inflating_statement_count(self):
        report = 'mode: count\nipapp/main.go:1.1,2.2 2 0\nipapp/main.go:1.1,2.2 2 3\n'
        result = evaluate(report, {'minimumStatementPercent': 100, 'packages': {}})['overall']
        self.assertEqual((result['covered'], result['total']), (2, 2))

    def test_inconsistent_duplicates_fail(self):
        report = 'mode: count\nipapp/main.go:1.1,2.2 1 1\nipapp/main.go:1.1,2.2 2 1\n'
        with self.assertRaisesRegex(ValueError, 'Inconsistent'):
            evaluate(report, {'minimumStatementPercent': 0, 'packages': {}})

    def test_invalid_empty_or_negative_reports_fail(self):
        for report in ('', 'mode: unknown', 'mode: count\n', 'mode: count\ninvalid 1 1',
                       'mode: count\nipapp/main.go:1.1,2.2 -1 1', 'mode: count\nipapp/main.go:1.1,2.2 1 -1'):
            with self.subTest(report=report), self.assertRaises(ValueError):
                evaluate(report, {'minimumStatementPercent': 0, 'packages': {}})

    def test_invalid_thresholds_fail(self):
        report = 'mode: count\nipapp/main.go:1.1,2.2 1 1'
        for minimum in (-1, 101, '100', True, float('nan')):
            with self.subTest(minimum=minimum), self.assertRaises(ValueError):
                evaluate(report, {'minimumStatementPercent': minimum, 'packages': {}})
