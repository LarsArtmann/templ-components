package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// panicAllowlist is every sanctioned `panic(` site in library (non-test,
// non-generated) source. Adding an entry requires a comment here.
var panicAllowlist = map[string]string{
	"icons/icon_paths.go": "developer data-integrity check (stray '|' separator in path data) — panics at init on a malformed icon table, caught by any test run",
	"cmd/tc/":             "CLI fail-fast on unusable embedded assets (developer-facing tool, not library runtime)",
}

// panicCallRe matches direct panic( calls (heuristic: the identifier at call
// position — enough to catch the convention drift; indirect panics via helper
// funcs are out of scope).
var panicCallRe = regexp.MustCompile(`(?m)^\s*panic\(`)

// TestZeroRuntimePanics pins the library guarantee "zero runtime panics in
// component code" (AGENTS.md): no library source file may call panic()
// outside the allowlist. A panic in a render path is a consumer crash by
// definition.
func TestZeroRuntimePanics(t *testing.T) {
	t.Parallel()

	packages := []string{
		"charts/echarts", "datastar", "display", "errorpage",
		"feedback", "forms", "htmx", "icons", "layout", "navigation",
		"utils", "utils/wire",
	}

	violations := 0

	for _, pkg := range packages {
		files, err := filepath.Glob(filepath.Join("..", pkg, "*.go"))
		if err != nil {
			t.Fatalf("glob %s: %v", pkg, err)
		}

		for _, file := range files {
			name := filepath.Base(file)
			if strings.HasSuffix(name, "_test.go") {
				continue // test files may panic to fail
			}

			content, readErr := os.ReadFile(file)
			if readErr != nil {
				t.Fatalf("read %s: %v", file, readErr)
			}

			if !panicCallRe.Match(content) {
				continue
			}

			rel := pkg + "/" + name
			if allowlisted(rel) {
				continue
			}

			violations++

			t.Errorf(
				"%s calls panic() — the library guarantees zero runtime panics (fail soft: render a fallback, return an error, or validate at construction). If this site is sanctioned, add it to panicAllowlist with a reason.",
				rel,
			)
		}
	}

	if violations == 0 {
		t.Log("zero panic() calls in library runtime code outside the allowlist")
	}
}

// allowlisted reports whether a file path is covered by a prefix allowlist
// entry.
func allowlisted(rel string) bool {
	for prefix := range panicAllowlist {
		if strings.HasPrefix(rel, prefix) || strings.HasPrefix(rel, "../"+prefix) {
			return true
		}
	}

	return false
}
