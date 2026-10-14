package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestScanRTL seeds violation fixtures in a temp dir and proves the reusable
// scanner core (shared by TestRTLLogicalProperties and `tc lint --rtl`) finds
// exactly the banned shapes with true file/line attribution.
func TestScanRTL(t *testing.T) {
	t.Parallel()

	t.Run("clean dir yields no findings", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeLintFixture(t, dir, "clean.templ", "<div class=\"ms-2 me-2 ps-1 pe-1 text-start inset-s-0\">\n</div>\n")

		findings, err := ScanRTL(dir)
		if err != nil {
			t.Fatalf("ScanRTL: %v", err)
		}

		if len(findings) != 0 {
			t.Errorf("expected 0 findings, got %d: %+v", len(findings), findings)
		}
	})

	t.Run("long classes and safe tokens never match", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeLintFixture(t, dir, "safe.templ",
			strings.Join([]string{
				"<div class=\"inset-s-0 rounded-s-lg border-s-2\">",
				"  <span class=\"text-start items-start justify-end col-start-2\">",
				"    <span class=\"ms-auto me-auto\">",
				"  </span>",
				"</div>",
			}, "\n")+"\n")

		findings, err := ScanRTL(dir)
		if err != nil {
			t.Fatalf("ScanRTL: %v", err)
		}

		if len(findings) != 0 {
			t.Errorf("expected 0 findings, got %d: %+v", len(findings), findings)
		}
	})

	t.Run("physical and inset-alias violations with line attribution", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeLintFixture(t, dir, "bad.templ", strings.Join([]string{
			"<div class=\"flex\">",
			"  <!-- prose mentioning start-0 and ml-2 must NOT trip the scanner -->",
			"  <span class=\"ml-2\">",
			"  <span class=\"pl-3 text-left\">",
			"  <span class=\"start-0\">",
			"</div>",
		}, "\n")+"\n")

		findings, err := ScanRTL(dir)
		if err != nil {
			t.Fatalf("ScanRTL: %v", err)
		}

		if len(findings) != 3 {
			t.Fatalf("expected 3 findings (ml-, pl-3+text-left, start-0), got %d:\n%+v", len(findings), findings)
		}

		wantLines := []int{3, 4, 5}
		for i, want := range wantLines {
			if findings[i].Line != want {
				t.Errorf("finding %d: Line = %d, want %d (comments must not shift line numbers)", i, findings[i].Line, want)
			}
		}

		wantRules := []string{"rtl-physical", "rtl-physical", "rtl-inset-alias"}
		for i, want := range wantRules {
			if findings[i].Rule != want {
				t.Errorf("finding %d: Rule = %q, want %q", i, findings[i].Rule, want)
			}
		}

		if !strings.Contains(findings[0].Snippet, "ml-2") {
			t.Errorf("finding 0: Snippet = %q, want it to contain ml-2", findings[0].Snippet)
		}
	})

	t.Run("missing dir is an error", func(t *testing.T) {
		t.Parallel()

		if _, err := ScanRTL(filepath.Join(t.TempDir(), "does-not-exist")); err == nil {
			t.Error("expected an error for a missing directory, got nil")
		}
	})
}

func writeLintFixture(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}
}
