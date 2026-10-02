package main

import (
	"bytes"
	"strings"
	"testing"
)

func parcelArgs(weight, zone string) []string {
	return []string{"-weight", weight, "-length", "30", "-width", "20", "-height", "10", "-zone", zone}
}

func TestRunPrintsQuote(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(parcelArgs("3", "eu"), &stdout, &stderr)

	if code != 0 {
		t.Fatalf("exit status = %d, want 0 (stderr: %s)", code, stderr.String())
	}
	want := "" +
		"weight tier        10.00\n" +
		"zone surcharge      8.00\n" +
		"total              18.00 EUR\n"
	if stdout.String() != want {
		t.Errorf("stdout =\n%s\nwant\n%s", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want it empty", stderr.String())
	}
}

func TestRunAcceptsZoneInAnyCase(t *testing.T) {
	for _, zone := range []string{"EU", "Eu", "eU", " eu ", "\tEU\n"} {
		t.Run(zone, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(parcelArgs("3", zone), &stdout, &stderr)

			if code != 0 {
				t.Fatalf("exit status = %d, want 0 (stderr: %s)", code, stderr.String())
			}
			if want := "18.00 EUR"; !strings.Contains(stdout.String(), want) {
				t.Errorf("stdout = %q, want it to contain %q", stdout.String(), want)
			}
		})
	}
}

func TestRunPrintsHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-h"}, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit status = %d, want 0", code)
	}
	if want := "Usage of shiprate"; !strings.Contains(stderr.String(), want) {
		t.Errorf("stderr = %q, want it to contain %q", stderr.String(), want)
	}
}

func TestRunRejectsBadInput(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantStderr string
	}{
		{"no flags", nil, "missing required flags: -weight, -length, -width, -height, -zone"},
		{"missing zone", []string{"-weight", "3", "-length", "30", "-width", "20", "-height", "10"}, "missing required flags: -zone"},
		{"unknown zone", parcelArgs("3", "mars"), `unknown zone "mars", valid zones are [domestic eu international]`},
		{"misspelt zone", parcelArgs("3", "internationl"), `unknown zone "internationl"`},
		{"zone with a space inside", parcelArgs("3", "inter national"), `unknown zone "inter national"`},
		{"negative weight", parcelArgs("-3", "eu"), "weight must be a positive number"},
		{"negative height", []string{"-weight", "3", "-length", "30", "-width", "20", "-height", "-10", "-zone", "eu"}, "height must be a positive number"},
		{"weight is not a number", parcelArgs("heavy", "eu"), "invalid value"},
		{"unknown flag", append(parcelArgs("3", "eu"), "-express"), "flag provided but not defined"},
		{"positional argument", append(parcelArgs("3", "eu"), "extra"), `unexpected argument "extra"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(tt.args, &stdout, &stderr)

			if code != 1 {
				t.Errorf("exit status = %d, want 1", code)
			}
			if stdout.Len() != 0 {
				t.Errorf("stdout = %q, want it empty", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}
