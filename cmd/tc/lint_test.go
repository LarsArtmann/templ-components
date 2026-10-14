package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runLintCase runs `tc lint` arg parsing + scanning and returns the exit code
// plus both output streams.
func runLintCase(args []string) (int, string, string) {
	outBuf, errBuf := &bytes.Buffer{}, &bytes.Buffer{}

	code := runLint(args, outBuf, errBuf)

	return code, outBuf.String(), errBuf.String()
}

func writeLintCLIFixture(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}
}

// TestRunLintClean proves the exit-code contract on a violation-free
// directory: 0, with a summary naming the scanner.
func TestRunLintClean(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeLintCLIFixture(t, dir, "ok.templ", "<div class=\"ms-2 text-start\">\n</div>\n")

	code, stdout, stderr := runLintCase([]string{"--rtl", dir})
	if code != 0 {
		t.Errorf("exit = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}

	if !strings.Contains(stdout, "clean") || !strings.Contains(stdout, "--rtl") {
		t.Errorf("summary should name the scanner and report clean:\n%s", stdout)
	}
}

// TestRunLintFindings proves findings exit 1 with file:line attribution and
// a finding-count summary.
func TestRunLintFindings(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeLintCLIFixture(t, dir, "bad.templ", "<div class=\"ml-2\">\n<span class=\"start-0\">\n</div>\n")

	code, stdout, stderr := runLintCase([]string{"--rtl", dir})
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}

	for _, want := range []string{"bad.templ:1:", "rtl-physical:", "bad.templ:2:", "rtl-inset-alias:", "ml-2", "2 finding(s)"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing %q:\n%s", want, stdout)
		}
	}

	if stderr != "" {
		t.Errorf("stderr should be empty on findings:\n%s", stderr)
	}
}

// TestRunLintUsageErrors pins the usage-error contract: exit 2 with a hint on
// stderr for every malformed invocation.
func TestRunLintUsageErrors(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "no scanner flag", args: []string{dir}, wantErr: "select at least one scanner"},
		{name: "unknown flag", args: []string{"--rtl", "--nope", dir}, wantErr: "unknown flag"},
		{name: "help", args: []string{"--help"}, wantErr: "--rtl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := runLintCase(tt.args)
			if code != 2 {
				t.Errorf("exit = %d, want 2\nstderr:\n%s", code, stderr)
			}

			if !strings.Contains(stderr, tt.wantErr) {
				t.Errorf("stderr missing %q:\n%s", tt.wantErr, stderr)
			}

			if stdout != "" {
				t.Errorf("usage errors should not write findings to stdout:\n%s", stdout)
			}
		})
	}
}

// TestRunLintMissingDir proves a missing target directory is a hard failure
// (exit 1) that names the directory — a typo'd path must never read as
// "clean".
func TestRunLintMissingDir(t *testing.T) {
	t.Parallel()

	code, _, stderr := runLintCase([]string{"--rtl", filepath.Join(t.TempDir(), "gone")})
	if code != 1 {
		t.Errorf("exit = %d, want 1", code)
	}

	if !strings.Contains(stderr, "no such directory") {
		t.Errorf("stderr should name the missing directory:\n%s", stderr)
	}
}

// TestRunLintMultipleDirs proves the summary counts every scanned directory.
func TestRunLintMultipleDirs(t *testing.T) {
	t.Parallel()

	clean := t.TempDir()
	writeLintCLIFixture(t, clean, "ok.templ", "<p class=\"pe-1\">\n</p>\n")

	dirty := t.TempDir()
	writeLintCLIFixture(t, dirty, "bad.templ", "<div class=\"mr-1\">\n</div>\n")

	code, stdout, stderr := runLintCase([]string{"--rtl", clean, dirty})
	if code != 1 {
		t.Fatalf("exit = %d, want 1\nstderr:\n%s", code, stderr)
	}

	if !strings.Contains(stdout, "in 2 directories") {
		t.Errorf("summary should count 2 directories:\n%s", stdout)
	}
}
