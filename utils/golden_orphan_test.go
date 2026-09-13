package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// goldenSweepPackages are the library packages with committed HTML golden
// files under testdata/. The orphan detector walks each one.
var goldenSweepPackages = []string{
	"charts/echarts", "datastar", "display", "errorpage",
	"feedback", "forms", "htmx", "icons", "layout", "navigation",
}

// TestNoOrphanGoldens fails when a committed .golden file is no longer
// referenced by any test in its package — an orphaned snapshot pins output
// nobody produces and silently rots (found 3 orphans in the 2026-09-13
// sweep that inspired this guard... none existed, which is exactly when a
// guard belongs: BEFORE the first orphan ships).
//
// Reference rule: the golden name (file stem) must appear as a quoted string
// in at least one non-generated test file of the same package — that is how
// golden.Assert(t, "<name>", ...) and AssertSnapshots name their snapshots.
// A name referenced only from a *_templ.go file does not count (generated
// code is not a test author).
func TestNoOrphanGoldens(t *testing.T) {
	t.Parallel()

	for _, pkg := range goldenSweepPackages {
		pkgDir := filepath.Join("..", pkg)

		testSource := readPackageTestSources(t, pkgDir)

		entries, err := os.ReadDir(filepath.Join(pkgDir, "testdata"))
		if err != nil {
			if os.IsNotExist(err) {
				continue // package has no goldens at all
			}

			t.Fatalf("read %s/testdata: %v — the orphan guard owns its fixture (fail loud, never skip)", pkgDir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".golden") {
				continue
			}

			name := strings.TrimSuffix(entry.Name(), ".golden")
			if !strings.Contains(testSource, `"`+name+`"`) {
				t.Errorf(
					"orphan golden: %s/testdata/%s — no test in %s references the name %q. Delete the file or re-wire its test (golden.Assert / golden.Snapshot{Name: ...}).",
					pkg,
					entry.Name(),
					pkg,
					name,
				)
			}
		}
	}
}

// readPackageTestSources concatenates every *_test.go file of a package so
// golden names can be searched as quoted strings.
func readPackageTestSources(t *testing.T, pkgDir string) string {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(pkgDir, "*_test.go"))
	if err != nil {
		t.Fatalf("glob %s tests: %v", pkgDir, err)
	}

	var builder strings.Builder

	for _, file := range files {
		content, readErr := os.ReadFile(file)
		if readErr != nil {
			t.Fatalf("read %s: %v", file, readErr)
		}

		builder.Write(content)
		builder.WriteByte('\n')
	}

	return builder.String()
}
