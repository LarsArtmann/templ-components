package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunLint proves the `tc lint` engine: exit-code contract (0 clean,
// 1 findings/failure, 2 usage), finding output format, and multi-dir scans.
// The scanners themselves are proven at the source (utils TestScanRTL et al);
// this pins the CLI wiring.
func TestRunLint(t *testing.T) {
	t.Parallel()

	cleanDir := t.TempDir()
	writeLintCLIFixture(t, cleanDir, "ok.templ", "<div class=\"ms-2 text-start\">\n</div>\n")

	dirtyDir := t.TempDir()
	writeLintCLIFixture(t, dirtyDir, "bad.templ", "<div class=\"ml-2\">\n<span class=\"start-0\">\n</div>\n")

	t.Run("clean dir exits 0", func(t *testing.T) {
		t.Parallel()

		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}

		if code := runLint([]string{"--rtl", cleanDir}, stdout, stderr); code != 0 {
			t.Errorf("exit = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
		}

		if !strings.Contains(stdout.String(), "clean") {
			t.Errorf("stdout should report clean:\n%s", stdout.String())
		}
	})

	t.Run("findings exit 1 with file:line output", func(t *testing.T) {
		t.Parallel()

		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}

		if code := runLint([]string{"--rtl", dirtyDir}, stdout, stderr); code != 1 {
			t.Fatalf("exit = %d, want 1\nstdout:\n%s\nstderr:\n%s", code, stdout.String(), stderr.String())
		}

		out := stdout.String()

		for _, want := range []string{"bad.templ:1:", "rtl-physical:", "bad.templ:2:", "rtl-inset-alias:", "ml-2", "2 finding(s)"} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout missing %q:\n%s", want, out)
			}
		}

		if stderr.Len() != 0 {
			t.Errorf("stderr should be empty on findings:\n%s", stderr.String())
		}
	})

	t.Run("no scanner flag is a usage error", func(t *testing.T) {
		t.Parallel()

		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}

		if code := runLint([]string{cleanDir}, stdout, stderr); code != 2 {
			t.Errorf("exit = %d, want 2", code)
		}

		if !strings.Contains(stderr.String(), "select at least one scanner") {
			t.Errorf("stderr should explain the missing scanner:\n%s", stderr.String())
		}
	})

	t.Run("unknown flag is a usage error", func(t *testing.T) {
		t.Parallel()

		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}

		if code := runLint([]string{"--rtl", "--nope", cleanDir}, stdout, stderr); code != 2 {
			t.Errorf("exit = %d, want 2", code)
		}
	})

	t.Run("help exits with usage", func(t *testing.T) {
		t.Parallel()

		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}

		if code := runLint([]string{"--help"}, stdout, stderr); code != 2 {
			t.Errorf("exit = %d, want 2", code)
		}

		if !strings.Contains(stderr.String(), "--rtl") {
			t.Errorf("help should list scanners:\n%s", stderr.String())
		}
	})

	t.Run("missing directory exits 1", func(t *testing.T) {
		t.Parallel()

		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}

		if code := runLint([]string{"--rtl", filepath.Join(t.TempDir(), "gone")}, stdout, stderr); code != 1 {
			t.Errorf("exit = %d, want 1", code)
		}

		if !strings.Contains(stderr.String(), "no such directory") {
			t.Errorf("stderr should name the missing directory:\n%s", stderr.String())
		}
	})

	t.Run("multiple dirs all scanned", func(t *testing.T) {
		t.Parallel()

		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}

		if code := runLint([]string{"--rtl", cleanDir, dirtyDir}, stdout, stderr); code != 1 {
			t.Fatalf("exit = %d, want 1\nstderr:\n%s", code, stderr.String())
		}

		if !strings.Contains(stdout.String(), "in 2 directories") {
			t.Errorf("summary should count 2 directories:\n%s", stdout.String())
		}
	})
}

func writeLintCLIFixture(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}
}
