package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCSSFreshness warns (or, in CI, FAILS) when the committed demo CSS might
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
		// In CI a stale committed CSS ships missing classes to consumers; fail
		// hard. Locally, just warn so `go test ./...` stays green during edits
		// before a CSS recompile.
		if os.Getenv("CI") != "" {
			t.Error(msg)
		} else {
			t.Logf("WARNING: %s", msg)
		}
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
