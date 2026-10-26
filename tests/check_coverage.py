"""Enforce measured Go statement coverage, including required package floors."""
import json
import os
from pathlib import Path
import re


def parse_report(text):
    lines = text.splitlines()
    if not lines or lines[0] not in ('mode: set', 'mode: count', 'mode: atomic'):
        raise ValueError('Missing or invalid Go coverage report')
    blocks = {}
    for line in lines[1:]:
        location, statements, hits = line.rsplit(' ', 2)
        match = re.fullmatch(r'(.+):\d+\.\d+,\d+\.\d+', location)
        count, hits = int(statements), int(hits)
        if match is None or count < 0 or hits < 0:
            raise ValueError('Invalid coverage block')
        covered = hits > 0
        if location in blocks:
            previous = blocks[location]
            if previous[1] != count:
                raise ValueError('Inconsistent duplicate coverage block')
            covered |= previous[2]
        # -coverpkg may report a block from multiple test packages. Count its
        # statements once and accept coverage recorded by any test package.
        blocks[location] = match.group(1), count, covered
    return list(blocks.values())


def evaluate(text, policy):
    blocks = parse_report(text)
    floors = {'overall': policy['minimumStatementPercent'], **policy['packages']}
    results = {}
    for name, minimum in floors.items():
        if isinstance(minimum, bool) or not isinstance(minimum, (int, float)) or not 0 <= minimum <= 100:
            raise ValueError('Invalid coverage threshold')
        selected = blocks if name == 'overall' else [block for block in blocks if block[0].startswith(name)]
        total = sum(count for _, count, _ in selected)
        if total == 0:
            raise ValueError(f'No statement coverage for {name}')
        covered = sum(count for _, count, hit in selected if hit)
        percent = covered * 100 / total
        results[name] = {'covered': covered, 'total': total, 'percent': percent, 'minimum': minimum}
        if percent < minimum:
            raise ValueError(f'Coverage regression for {name}: {percent:.2f}% < {minimum}%')
    return results


def main():
    root = Path(__file__).resolve().parents[1]
    policy = json.loads((root / '.github/coverage-baseline.json').read_text(encoding='utf-8'))
    results = evaluate((root / '.artifacts/coverage.out').read_text(encoding='utf-8'), policy)
    summary = '## Go statement coverage\n\n| Scope | Covered | Minimum |\n| --- | --- | --- |\n'
    for name, result in results.items():
        summary += f"| {name} | {result['covered']}/{result['total']} ({result['percent']:.2f}%) | {result['minimum']}% |\n"
    destination = root / '.artifacts'
    (destination / 'coverage-summary.md').write_text(summary, encoding='utf-8')
    (destination / 'coverage-summary.json').write_text(json.dumps({
        'commit': os.environ.get('GITHUB_SHA', 'local'), 'coverage': results,
    }, indent=2) + '\n', encoding='utf-8')
    print(summary)


if __name__ == '__main__':
    try:
        main()
    except (ValueError, KeyError, OSError) as error:
        raise SystemExit(str(error)) from error
