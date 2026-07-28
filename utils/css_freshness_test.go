package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCSSFreshness warns when the committed demo CSS might be stale — older
// than the most recently modified .templ or .go source file. A stale CSS
// means new Tailwind classes added in source may not be in the compiled CSS.
//
// This is informational (t.Logf) not failing because:
// 1. The Dockerfile 3-stage pipeline always recompiles CSS during image build
// 2. CI runs templ generate + go build independently of CSS
// 3. Local developers use `nix run .#build` which includes templ generate
//
// But a stale CSS committed to the repo means `go run ./examples/demo` serves
// outdated styles. This test surfaces that gap.
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
		t.Logf(
			"WARNING: demo CSS (%s) is older than newest source file — "+
				"recompile with: tailwindcss -i examples/demo/demo.css -o examples/demo/static/app.css --minify",
			cssInfo.ModTime().Format("2006-01-02 15:04"),
		)
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
