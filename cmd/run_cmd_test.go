package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		c      config
		input  string
		output string
		err    error
	}{
		{
			c:      config{numTimes: 5},
			input:  "\n",
			output: strings.Repeat("Your Name Please? Press the Enter Key when done.\n", 1),
			err:    errors.New("You didn't enter your name"),
		},
		{
			c:      config{numTimes: 5},
			input:  "Bill Bryson",
			output: "Your Name Please? Press the Enter Key when done.\n" + strings.Repeat("Nice to meet you Bill Bryson\n", 5),
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		rd := strings.NewReader(tc.input)
		err := runCmd(rd, byteBuf, tc.c)
		if err != nil && tc.err == nil {
			t.Fatalf("Expected nil error,  got : %v\n", err)
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error: %v, Got error: %v\n", tc.err.Error(), err.Error())
		}

		gotMsg := byteBuf.String()
		if gotMsg != tc.output {
			t.Errorf("Expected stdout message to be: %v, Got: %v\n", tc.output, gotMsg)
		}
		byteBuf.Reset()
	}
}

func TestRun_RunCmdError(t *testing.T) {
	stdin := strings.NewReader("\n")
	stdout := new(bytes.Buffer)

	exitCode := run(stdin, stdout, []string{"2"})
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	got := stdout.String()
	if !strings.Contains(got, "You didn't enter your name") {
		t.Fatalf("expected name error in output, got %q", got)
	}
}

type errReader struct{}

func (e *errReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func TestGetName_ScannerError(t *testing.T) {
	stdin := &errReader{}
	stdout := new(bytes.Buffer)

	_, err := getName(stdin, stdout)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "read error" {
		t.Fatalf("expected read error, got %v", err)
	}
}
