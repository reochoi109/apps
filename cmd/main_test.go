package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := &bytes.Buffer{}

	exitCode := run(stdin, stdout, []string{})
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}

	got := stdout.String()
	if got == "" {
		t.Fatalf("expected output, got empty string")
	}
}

func TestRun_Success(t *testing.T) {
	stdin := strings.NewReader("홍길동\n")
	stdout := &bytes.Buffer{}

	exitCode := run(stdin, stdout, []string{"5"})
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	got := stdout.String()
	if !strings.Contains(got, "홍길동") {
		t.Fatalf("expected output to contain name, got %q", got)
	}
}

func TestRun_Help(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := &bytes.Buffer{}

	exitCode := run(stdin, stdout, []string{"-h"})
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	got := stdout.String()
	if got == "" {
		t.Fatalf("expected output, got empty string")
	}
}
