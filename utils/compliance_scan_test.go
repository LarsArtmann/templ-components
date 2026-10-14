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
}

// TestScanRTLSeededViolations proves both RTL rules fire with true file/line
// attribution on a seeded fixture (HTML comments must not shift line numbers
// nor survive as prose).
func TestScanRTLSeededViolations(t *testing.T) {
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
}

// TestScanRTLMissingDir proves a missing target directory is an error, never
// a silent "clean".
func TestScanRTLMissingDir(t *testing.T) {
	t.Parallel()

	if _, err := ScanRTL(filepath.Join(t.TempDir(), "does-not-exist")); err == nil {
		t.Error("expected an error for a missing directory, got nil")
	}
}

// TestScanDarkMode seeds dark-mode gap fixtures and proves the reusable
// scanner core reports both rule families with same-line dark: coverage
// respected and generated/test files excluded.
func TestScanDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("gaps in templ and go sources, both rules", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeLintFixture(t, dir, "gaps.templ", "<span class=\"text-gray-500\">x</span>\n")
		writeLintFixture(t, dir, "gaps_styles.go", "package p\n\nvar c = \"bg-blue-600\"\n")
		writeLintFixture(t, dir, "gen_templ.go", "package p\n\nvar c = \"text-gray-900\"\n") // generated: excluded
		writeLintFixture(t, dir, "gaps_test.go", "package p\n\nvar c = \"bg-red-600\"\n")    // test: excluded

		findings, err := ScanDarkMode(dir)
		if err != nil {
			t.Fatalf("ScanDarkMode: %v", err)
		}

		if len(findings) != 2 {
			t.Fatalf("expected 2 findings, got %d:\n%+v", len(findings), findings)
		}

		if findings[0].Rule != "darkmode-neutral" || findings[0].Line != 1 {
			t.Errorf("finding 0 = %+v, want darkmode-neutral line 1", findings[0])
		}

		if findings[1].Rule != "darkmode-semantic" || findings[1].Line != 3 {
			t.Errorf("finding 1 = %+v, want darkmode-semantic line 3", findings[1])
		}
	})

	t.Run("same-line dark variant passes", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeLintFixture(t, dir, "ok.templ",
			"<span class=\"text-gray-500 dark:text-gray-400 bg-blue-600 dark:bg-blue-500\">x</span>\n"+
				"// text-gray-900 inside a Go comment is not markup\n")

		findings, err := ScanDarkMode(dir)
		if err != nil {
			t.Fatalf("ScanDarkMode: %v", err)
		}

		if len(findings) != 0 {
			t.Errorf("expected 0 findings, got %d:\n%+v", len(findings), findings)
		}
	})

	t.Run("missing dir is an error", func(t *testing.T) {
		t.Parallel()

		if _, err := ScanDarkMode(filepath.Join(t.TempDir(), "gone")); err == nil {
			t.Error("expected an error for a missing directory, got nil")
		}
	})
}

// TestScanMotionReduce seeds transition fixtures and proves the rule: a
// transition/animate class needs a same-line motion-reduce fallback or a
// shared motion constant.
func TestScanMotionReduce(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeLintFixture(t, dir, "gaps.templ", "<span class=\"transition-colors\">x</span>\n")
	writeLintFixture(t, dir, "ok.templ",
		"<span class=\"transition-colors motion-reduce:transition-none\">x</span>\n"+
			"<span class={ \"transitionFast\" }>x</span>\n"+
			"<!-- prose about transition-colors stays clean -->\n")
	writeLintFixture(t, dir, "sub", "ignored") // non-.templ files are skipped

	findings, err := ScanMotionReduce(dir)
	if err != nil {
		t.Fatalf("ScanMotionReduce: %v", err)
	}

	if len(findings) != 1 || findings[0].Line != 1 || findings[0].File != "gaps.templ" {
		t.Errorf("expected exactly the gaps.templ line-1 finding, got %+v", findings)
	}

	if findings[0].Rule != "motion-reduce" {
		t.Errorf("Rule = %q, want motion-reduce", findings[0].Rule)
	}
}

// TestScanIconPaths seeds inlined-path fixtures and proves both literal forms
// fire while component-sourced paths and canonical names stay clean.
func TestScanIconPaths(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeLintFixture(t, dir, "bad.templ", "<path d=\"M6 18L18 6M6 6l12 12\"></path>\n")
	writeLintFixture(t, dir, "bad_js.go", "package p\n\nvar s = \"<path d=\\\"M6 18\\\">\"\n")
	writeLintFixture(t, dir, "clean.templ", "<path d={ svg.PathXMark }></path>\n")

	findings, err := ScanIconPaths(dir)
	if err != nil {
		t.Fatalf("ScanIconPaths: %v", err)
	}

	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d:\n%+v", len(findings), findings)
	}

	if findings[0].File != "bad.templ" || findings[0].Rule != "icon-path-inline" {
		t.Errorf("finding 0 = %+v, want bad.templ icon-path-inline", findings[0])
	}

	if findings[1].File != "bad_js.go" {
		t.Errorf("finding 1 = %+v, want bad_js.go (embedded-JS form)", findings[1])
	}
}

func writeLintFixture(t *testing.T, dir, name, content string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}
}
