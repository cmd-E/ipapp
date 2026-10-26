# ipapp
Offline IPv4 subnet analysis: masks, network boundaries, host ranges and grouping
addresses by network. The application uses only the Go standard library.

![изображение](https://user-images.githubusercontent.com/49339376/136646568-cbd6768b-8097-4697-bdad-62562c3b5e7c.png)


## Build and usage

Install Go 1.26 or later. From the repository root:

```console
go build -o ipapp .
./ipapp -m 10 192.152.10.36 192.152.10.38
```

On Windows, use `go build -o ipapp.exe .` and `./ipapp.exe`.

`-m` accepts a prefix from 0 through 32, defaulting to 24. Both `-m 24` and `-m=24`
work. Options must precede addresses; `-h` prints help. IPv4 input must contain four
decimal octets without signs, leading zeroes, CIDR suffixes or surrounding whitespace.
IPv6 is outside the existing application's scope.

Conventional subnet ranges exclude network and directed-broadcast addresses. `/31`
uses two point-to-point endpoints ([RFC 3021](https://www.rfc-editor.org/rfc/rfc3021));
`/32` contains one address. Neither has a directed broadcast. The `/0` arithmetic
range contains 4,294,967,294 conventional host addresses. Counts describe subnet
arithmetic, not allocation policy or actual reachability of reserved addresses.

The command preserves first-seen network order and the original order of addresses
within groups, including duplicate addresses. Reports go to stdout and diagnostics
to stderr. Invalid addresses are quoted so terminal control characters are escaped.
If input mixes valid and invalid addresses, the valid addresses are still analyzed.

| Exit status | Meaning |
| --- | --- |
| 0 | Successful analysis or help |
| 1 | Missing/invalid addresses or report output failure |
| 2 | Invalid options or prefix length |

## Review and Go API changes

[REVIEW.md](REVIEW.md) records findings by severity and the fixes. The IPv4 type keeps
validated address/prefix state private. Decimal/binary constructors and `ConvertMask`
return `(IP, error)`; callers must handle errors. Derived addresses preserve their
prefix. `GetIPInDec`, `GetIPInBin` and `PrefixLength` provide read access. The legacy
`GetBroadcastAddress` name returns the inclusive upper boundary, including endpoint
boundaries for `/31` and `/32`.

## Tests and coverage

```sh
go test -count=1 ./...
go vet ./...
mkdir -p .artifacts
go test -coverpkg=./... -coverprofile=.artifacts/coverage.out ./...
go tool cover -func=.artifacts/coverage.out
python tests/check_coverage.py
python -m unittest discover -s tests -p 'test_*.py' -v
```

Python 3 is needed only for the coverage gate and its self-tests. Linux CI additionally
runs the Go suite with `-race`. The initial expanded suite measured 114/115 production
statements (99.13%) and 54/54 core statements (100%). The uncovered statement is the
process-exit wrapper; compiled CLI tests exercise it separately. Coverage gates require
99% overall and 100% for `packages/ip`; missing reports or required packages fail CI.

Tests follow FIRST: the pure core cases require no I/O; each command invocation owns
its flags and buffers; the compiled CLI test builds into a temporary directory; inputs
and random seeds are fixed; assertions validate results and exit status. The original
range fixtures remain, together with strict parsing, all 33 prefixes checked against
`net.IP.Mask`, golden host ranges, retained prefixes, grouping order/ownership, escaped
diagnostics, independent invocations and broken-output checks. Fuzz seeds run as part
of normal tests, and a short fuzz campaign runs in CI. Longer local campaigns:

```sh
go test -run='^$' -fuzz='^FuzzDecimalParsing$' -fuzztime=30s ./packages/ip
go test -run='^$' -fuzz='^FuzzBinaryParsing$' -fuzztime=30s ./packages/ip
```

## GitHub Actions and dependency maintenance

CI runs for pull requests, pushes to `master`, manual dispatches and a weekly schedule.
It checks Linux on Go 1.26 and 1.27 with the race detector, and Windows on Go 1.27.
Limiting push runs to the default branch prevents duplicate push/PR checks. Test JSON,
text/JSON/Markdown/HTML coverage and fuzz logs are retained as artifacts for 14 days.

CodeQL performs a Go build and security-extended analysis in a separate job. Actions
are pinned to verified commit hashes. Dependabot checks Actions weekly once its config
is merged into the default branch. There are no third-party Go modules to update;
add a gomod entry if external modules are introduced.

`.github/branch-ruleset.json` proposes default-branch protection requiring PRs, resolved
conversations, an up-to-date branch, all three test jobs and `CodeQL (go)`. Force pushes
and deletion are blocked, with no bypass actors. It requires zero human approvals so
a sole maintainer can merge their own passing PR. This file does not activate rules:
the proposal must be applied separately through GitHub settings or its ruleset API.
