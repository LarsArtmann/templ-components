package visualtest

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// The 2026-10-30 webfont race (docs/status/2026-10-30_15-40 §d5/d6) proved that
// a golden suite without its determinism helpers rots back into flake silently:
// route goldens raced the Google Fonts CDN for weeks, and nothing failed when
// the fix landed because goldens only prove "unchanged", never "deterministic".
// These guards pin the CONTRACT instead of the pixels:
//
//   1. Every pixel-comparison capture path (component goldens via harness.go,
//      demo routes, site routes) must call both blockExternalWebFonts() and
//      waitFontsSettled(). A future capture path or refactor that drops either
//      helper fails here, not in a flaky PNG diff three sessions later.
//   2. newTab must hand back a context with a deadline (tabLifetime), so a
//      wedged renderer fails ONE test instead of panicking the suite.
//
// This test reads source files, so `go test` caching cannot see edits to them:
// run guards like this with -count=1 when proving a failure or a fix
// (AGENTS.md, the TestGoDirectiveSkew lesson).

// capturePathFiles are the three pixel-comparison capture paths. A new capture
// path must be added here AND call both determinism helpers.
var capturePathFiles = []string{
	"harness.go",           // component goldens (AssertScreenshot chain)
	"route_golden_test.go", // demo route goldens
	"site_routes_test.go",  // website route goldens
}

func capturePathSource(t *testing.T, name string) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed: cannot locate the visualtest package dir")
	}

	path := filepath.Join(filepath.Dir(thisFile), name)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read capture path %s: %v", name, err)
	}

	return string(data)
}

func TestCapturePathsBlockExternalFonts(t *testing.T) {
	t.Parallel()

	// The helpers themselves must exist; a rename should break THIS assertion
	// loudly instead of letting every capture path drift silently.
	fontsSource := capturePathSource(t, "external_fonts.go")
	for _, symbol := range []string{
		"func blockExternalWebFonts()",
		"func waitFontsSettled()",
	} {
		if !strings.Contains(fontsSource, symbol) {
			t.Errorf("external_fonts.go no longer defines %q; update the determinism guards", symbol)
		}
	}

	for _, name := range capturePathFiles {
		source := capturePathSource(t, name)
		for _, call := range []string{"blockExternalWebFonts()", "waitFontsSettled()"} {
			if !strings.Contains(source, call) {
				t.Errorf("%s captures pixels without %s — the webfont race is back (see external_fonts.go)", name, call)
			}
		}
	}
}

func TestNewTabHasDeadline(t *testing.T) {
	if tabLifetime <= 0 {
		t.Fatalf("tabLifetime must be positive, got %v", tabLifetime)
	}

	ctx, cancel := newTab(t)
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal(
			"newTab returned a context WITHOUT a deadline: a wedged renderer would hang until the go-test alarm kills the whole suite",
		)
	}

	remaining := time.Until(deadline)
	if remaining <= 0 || remaining > tabLifetime {
		t.Fatalf("newTab deadline out of range: %v (want (0, %v])", remaining, tabLifetime)
	}
}
