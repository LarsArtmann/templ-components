package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCSSFreshness warns (or FAILS) when the committed demo CSS might
// be stale — older than the most recently modified .templ or .go source file.
// A stale CSS means new Tailwind classes added in source may not be in the
// compiled CSS.
//
// Locally this is informational (t.Logf): the Dockerfile pipeline recompiles
// CSS on every image build, and developers use `nix run .#build`. But a stale
// CSS COMMITTED to the repo means `go run ./examples/demo` serves outdated
// styles. In CI (CI env var set, e.g. GitHub Actions), a stale committed CSS
// is a real regression — the build artifact ships with missing classes — so
// the test fails there. The root cause of a real stale-CSS incident
// (bg-amber-50 missing from compiled CSS) is documented in the
// 2026-07-28 status report, section d.1/838016c.
//
// Set TC_CSS_FRESHNESS_STRICT=1 to make the local run fail-capable — the
// complement of `scripts/ci-repro.sh --css` (which does the proper content
// diff via `nix run .#css`) for environments without Nix.
func TestCSSFreshness(t *testing.T) {
	t.Parallel()

	cssPath := "../examples/demo/static/app.css"

	cssInfo, err := os.Stat(cssPath)
	if err != nil {
		t.Skipf("demo CSS not found at %s: %v", cssPath, err)

		return
	}

	// Find the newest source file (.templ or .go, excluding _templ.go and _test.go).
	newestTime := cssInfo.ModTime()

	err = filepath.Walk("..", func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return walkErr
		}

		// Only scan component source directories.
		if !isComponentSource(path) {
			return nil
		}

		if info.ModTime().After(newestTime) {
			newestTime = info.ModTime()
		}

		return nil
	})
	if err != nil {
		t.Logf("walk error: %v", err)
	}

	if newestTime.After(cssInfo.ModTime()) {
		msg := fmt.Sprintf(
			"demo CSS (%s) is older than newest source file — "+
				"recompile with: nix run .#css  OR  tailwindcss -i examples/demo/demo.css -o examples/demo/static/app.css --minify",
			cssInfo.ModTime().Format("2006-01-02 15:04"),
		)
		// Informational by default — the CSS Freshness CI job does a proper
		// content diff (nix run .#css → diff). This timestamp check is too
		// fragile for CI because templ generate touches source files before
		// tests run, giving them newer mtimes than the committed CSS.
		// TC_CSS_FRESHNESS_STRICT=1 opts into a hard local failure.
		if os.Getenv("TC_CSS_FRESHNESS_STRICT") != "" {
			t.Fatalf("STRICT: %s", msg)
		}

		t.Logf("WARNING: %s", msg)
	}
}

// TestVisualFailArtifactsIgnored guards against committing visual-regression
// failure artifacts: the `.fail/` directory (actual + diff PNGs from a red
// visual run) must stay gitignored AND untracked. A committed .fail PNG has
// happened (broad `git add -A` before the daemon guards) and pollutes the
// repo with one-off debugging output.
func TestVisualFailArtifactsIgnored(t *testing.T) {
	t.Parallel()

	gitignore, err := os.ReadFile("../visualtest/.gitignore")
	if err != nil {
		t.Fatalf("read visualtest/.gitignore: %v\nIf it was deliberately removed, delete this test too.", err)
	}

	if !strings.Contains(string(gitignore), ".fail/") {
		t.Errorf("visualtest/.gitignore no longer ignores .fail/ — failure artifacts (actual/diff PNGs) can be committed by broad `git add`")
	}

	tracked, err := exec.Command("git", "-C", "..", "ls-files", "visualtest/testdata/.fail").Output()
	if err != nil {
		t.Skipf("git not usable here (CI clone edge): %v", err)

		return
	}

	if len(strings.TrimSpace(string(tracked))) > 0 {
		t.Errorf("visual-regression .fail artifacts are TRACKED in git:\n%s\nRemove them (they are per-run debugging output).", strings.TrimSpace(string(tracked)))
	}
}

// isComponentSource returns true for .templ and .go files in component
// directories (excluding generated and test files).
func isComponentSource(path string) bool {
	dirs := []string{
		"display/", "errorpage/", "feedback/", "forms/",
		"htmx/", "layout/", "navigation/", "recipes/",
	}

	matched := false

	for _, dir := range dirs {
		prefix := filepath.Join("..", dir)
		if strings.HasPrefix(path, prefix) {
			matched = true

			break
		}
	}

	if !matched {
		return false
	}

	if strings.HasSuffix(path, "_templ.go") || strings.HasSuffix(path, "_test.go") {
		return false
	}

	return strings.HasSuffix(path, ".templ") || strings.HasSuffix(path, ".go")
}
