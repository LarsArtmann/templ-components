package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// concatPrefixRe matches quoted literals in concatenation position.
var concatPrefixRe = regexp.MustCompile(`"([^"]*)"\s*\+`)

// goldenSweepPackages are the library packages with committed HTML golden
// files under testdata/. The orphan detector walks each one.
var goldenSweepPackages = []string{
	"charts/echarts", "datastar", "display", "errorpage",
	"feedback", "forms", "htmx", "icons", "layout", "navigation",
}

// TestNoOrphanGoldens fails when a committed .golden file is no longer
// referenced by any test in its package — an orphaned snapshot pins output
// nobody produces and silently rots.
//
// Reference rule: the golden name (file stem) is referenced when it appears
// as a quoted string in a non-generated test file of the same package —
// either verbatim (golden.Assert(t, "<name>", ...)) or as the literal PREFIX
// of a concatenation ("circular_progress_"+tt.name, where the table's case
// names supply the rest). Both forms are how the repo's sweep tests name
// snapshots; a name matching neither is an orphan.
func TestNoOrphanGoldens(t *testing.T) {
	t.Parallel()

	for _, pkg := range goldenSweepPackages {
		pkgDir := filepath.Join("..", pkg)

		testSource := readPackageTestSources(t, pkgDir)
		prefixes := concatPrefixLiterals(testSource)

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
			if goldenReferenced(testSource, prefixes, name) {
				continue
			}

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

// goldenReferenced reports whether a golden name is wired into a test:
// verbatim quoted, or built from a concatenated literal prefix the name
// starts with.
func goldenReferenced(testSource string, prefixes []string, name string) bool {
	if strings.Contains(testSource, `"`+name+`"`) {
		return true
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

// concatPrefixLiterals extracts string literals directly followed by a +
// concatenation (the `"prefix"+computed` sweep pattern).
func concatPrefixLiterals(testSource string) []string {
	var prefixes []string

	for _, match := range concatPrefixRe.FindAllStringSubmatch(testSource, -1) {
		literal := match[1]
		if literal != "" && strings.Contains(literal, "_") {
			prefixes = append(prefixes, strings.TrimSuffix(literal, "_")+"_")
		}
	}

	return prefixes
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
