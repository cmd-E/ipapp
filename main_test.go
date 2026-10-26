package main

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestCommandArgumentsAndExitStatus(t *testing.T) {
	cases := []struct {
		name               string
		args               []string
		code               int
		output, diagnostic string
	}{
		{"default", []string{"192.0.2.1"}, 0, "Number of hosts               254", ""},
		{"separate-mask", []string{"-m", "10", "192.152.10.36"}, 0, "192.128.0.0", ""},
		{"equals-mask", []string{"-m=10", "192.152.10.36"}, 0, "192.128.0.0", ""},
		{"double-dash-mask", []string{"--m=24", "192.0.2.1"}, 0, "192.0.2.0", ""},
		{"repeated-mask", []string{"-m=8", "-m=24", "192.0.2.1"}, 0, "192.0.2.0", ""},
		{"terminator", []string{"--", "192.0.2.1"}, 0, "192.0.2.0", ""},
		{"zero-prefix", []string{"-m=0", "192.0.2.1"}, 0, "Number of hosts               4294967294", ""},
		{"point-to-point", []string{"-m=31", "192.0.2.11"}, 0, "Number of hosts               2", ""},
		{"host-route", []string{"-m=32", "192.0.2.11"}, 0, "Number of hosts               1", ""},
		{"no-address", nil, 1, "", "weren't provided"},
		{"all-invalid", []string{"bad"}, 1, "", "No IPs to examine"},
		{"mixed", []string{"bad", "192.0.2.1"}, 1, "Examining 192.0.2.1", "invalid IPv4"},
		{"ipv6", []string{"::1"}, 1, "", "invalid IPv4"},
		{"negative-prefix", []string{"-m=-1", "192.0.2.1"}, 2, "", "between 0 and 32"},
		{"large-prefix", []string{"-m=33", "192.0.2.1"}, 2, "", "between 0 and 32"},
		{"noninteger-prefix", []string{"-m=no", "192.0.2.1"}, 2, "", "Invalid arguments"},
		{"missing-mask", []string{"-m"}, 2, "", "Invalid arguments"},
		{"unknown-option", []string{"-unknown"}, 2, "", "Invalid arguments"},
		{"short-help", []string{"-h"}, 0, "", "Usage: ipapp"},
		{"long-help", []string{"--help"}, 0, "", "Usage: ipapp"},
	}
	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			before := append([]string(nil), item.args...)
			var out, diagnostics bytes.Buffer
			code := run(item.args, &out, &diagnostics)
			if code != item.code || !strings.Contains(out.String(), item.output) || !strings.Contains(diagnostics.String(), item.diagnostic) {
				t.Fatalf("code=%d\nstdout=%s\nstderr=%s", code, &out, &diagnostics)
			}
			if item.diagnostic == "" && diagnostics.Len() != 0 {
				t.Fatalf("unexpected diagnostic: %s", &diagnostics)
			}
			if item.output == "" && out.Len() != 0 {
				t.Fatalf("unexpected report: %s", &out)
			}
			if !reflect.DeepEqual(before, item.args) {
				t.Fatal("argument slice mutated")
			}
		})
	}
}

func TestCommandReportsSubnetBoundariesAndGrouping(t *testing.T) {
	var output, diagnostics bytes.Buffer
	code := run([]string{"-m=24", "192.0.2.129", "10.0.0.1", "192.0.2.130"}, &output, &diagnostics)
	for _, text := range []string{"192.0.2.0 | 11000000.00000000.00000010.00000000", "192.0.2.1 |", "192.0.2.254 |", "192.0.2.255 |", "Ips in 192.0.2.0 network:\n192.0.2.129\n192.0.2.130\nIps in 10.0.0.0 network:\n10.0.0.1\n"} {
		if code != 0 || !strings.Contains(output.String(), text) {
			t.Fatalf("missing %q in %s (code %d)", text, &output, code)
		}
	}
	for _, prefix := range []string{"-m=31", "-m=32"} {
		output.Reset()
		diagnostics.Reset()
		if run([]string{prefix, "192.0.2.129"}, &output, &diagnostics) != 0 || !strings.Contains(output.String(), "N/A (no directed broadcast") {
			t.Fatal("endpoint misreported as broadcast")
		}
	}
}

func TestDiagnosticsEscapeTerminalControlCharacters(t *testing.T) {
	for _, args := range [][]string{{"bad\nFORGED\x1b[2J"}, {"-unknown\x1b[2J"}} {
		var output, diagnostics bytes.Buffer
		if run(args, &output, &diagnostics) == 0 || strings.Contains(diagnostics.String(), "\x1b") || strings.Contains(diagnostics.String(), "\nFORGED") {
			t.Fatalf("unsafe diagnostics: %q", diagnostics.String())
		}
	}
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, errors.New("closed output") }

func TestOutputFailureReturnsNonzero(t *testing.T) {
	var diagnostics bytes.Buffer
	if run([]string{"192.0.2.1"}, failedWriter{}, &diagnostics) != 1 || !strings.Contains(diagnostics.String(), "Cannot write report") {
		t.Fatal(diagnostics.String())
	}
}

func TestCompiledCLI(t *testing.T) {
	if testing.Short() {
		t.Skip("compiled CLI integration check")
	}
	binary := filepath.Join(t.TempDir(), "ipapp")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	for _, item := range []struct {
		name string
		args []string
		code int
	}{
		{"success", []string{"-m", "24", "192.0.2.1"}, 0},
		{"invalid", []string{"invalid"}, 1},
		{"mask", []string{"-m=33", "192.0.2.1"}, 2},
		{"help", []string{"-h"}, 0},
	} {
		t.Run(item.name, func(t *testing.T) {
			command := exec.Command(binary, item.args...)
			command.Dir = t.TempDir()
			output, err := command.CombinedOutput()
			code := 0
			if err != nil {
				var exit *exec.ExitError
				if !errors.As(err, &exit) {
					t.Fatal(err)
				}
				code = exit.ExitCode()
			}
			if code != item.code {
				t.Fatalf("code=%d want=%d output=%s", code, item.code, output)
			}
		})
	}
}
