package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		c   config
		err error
	}{
		{
			c:   config{},
			err: errors.New("Must specify a number greater then 0"),
		},
		{
			c:   config{numTimes: -1},
			err: errors.New("Must specify a number greater then 0"),
		},
		{
			c:   config{numTimes: 10},
			err: nil,
		},
	}

	for _, tc := range tests {
		err := validateArgs(tc.c)
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Expected error to be: %v, got: %v\n", tc.err, err)
		}

		if tc.err == nil && err != nil {
			t.Errorf("Expected nill error, got: %v\n", err)
		}

	}
}

func TestRun_ValidateArgsError(t *testing.T) {
	stdin := strings.NewReader("Reo\n")
	stdout := &bytes.Buffer{}

	exitCode := run(stdin, stdout, []string{"0"})
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	got := stdout.String()

	if !strings.Contains(got, "Must specify a number greater then 0") {
		t.Fatalf("expected validation error in ouput, got %q", got)
	}

	if !strings.Contains(got, "Usage:") {
		t.Fatalf("expected usage in output, got %q", got)
	}
}
