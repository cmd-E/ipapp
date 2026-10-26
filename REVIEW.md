# Implementation and security review

Reviewed 26 October 2026. No new application features were added. No critical or
high-severity remote security issue was identified: this is a local, offline IPv4
calculator with no third-party Go dependencies, network listeners, or file loading.
That scope does not establish that the application is free of vulnerabilities.

## Findings, ordered by severity

| Severity | Original issue and evidence | Implemented approach |
| --- | --- | --- |
| Medium | Decimal and binary constructors indexed four components without checking their count and discarded numeric conversion errors. Malformed inputs could panic or become a different address; public constructors also accepted invalid prefixes. | Return explicit errors, use `net/netip` for strict IPv4 parsing, validate binary octets and prefixes, and keep validated address/prefix state private. |
| Medium | The CLI parsed flags, then reparsed `os.Args` with a manual flag-removal loop. `-m 10 192.152.10.36` incorrectly diagnosed `10` as an IP; adjacent options could be skipped. | Use a command-local `flag.FlagSet` and its positional arguments. Document that options precede addresses. |
| Medium | Invalid-only input printed `No IPs to examine` and exited with status 0. Invalid mixed inputs also returned success. | Report errors to stderr, retain analysis of valid addresses, and return a nonzero status whenever an address is rejected. |
| Medium | Minimum/maximum/broadcast calculations discarded prefix metadata, breaking subsequent network operations. `/31` and `/32` reported zero hosts despite returning endpoints. | Use fixed-width integer masks, preserve the prefix on all derived addresses, count two `/31` point-to-point endpoints and one `/32` host, and avoid labeling these as directed-broadcast networks. |
| Medium | Grouping performed repeated linear searches and recalculated network masks, giving quadratic behavior as distinct networks increased. | Use a map for lookup and an ordered result slice to retain first-seen network and input-address order. |
| Low | Invalid address strings were echoed verbatim, allowing terminal control sequences in error output. This is local-output hardening, not evidence of a remote exploit. | Quote rejected input and flag diagnostics so control characters are escaped. |
| Low | Global flags, direct stdout writes, startup logger initialization and unused string/mask helper packages coupled command behavior to process state. Only three address-range tests existed; there were no CI or coverage gates. | Separate command execution from process exit with injected writers; remove helpers superseded by standard flag parsing and integer operations; retain the original regression fixtures and add deterministic unit/integration/fuzz tests and CI. |

## Architecture and compatibility

The application now has a small CLI layer and a validated IPv4 value type. The core
is independent of terminal I/O and global configuration. The zero IP value remains
valid as `0.0.0.0/0`. Grouping preserves the original first-seen order and duplicates.

The existing constructor names now return `(IP, error)` instead of silently accepting
bad input. Address fields are private; use `GetIPInDec`, `GetIPInBin` and `PrefixLength`.
This is an intentional Go API change for callers; all repository callers are updated.
The command name, `-m` option, default `/24` and valid-address report remain available.
Go 1.26 or later is required; CI checks both supported Go 1.26 and 1.27 toolchains.

`GetBroadcastAddress` retains its name for compatibility and returns the inclusive
upper subnet boundary. For `/31` and `/32`, this boundary is an endpoint rather than
a directed broadcast; the CLI says so explicitly.

The `/31` interpretation follows [RFC 3021](https://www.rfc-editor.org/rfc/rfc3021).
Strict parsing uses the standard [net/netip package](https://pkg.go.dev/net/netip).
Actual allocation, routing, multicast/reserved-use policy and network reachability
are outside this calculator's existing scope; tests exercise arithmetic and CLI behavior.
