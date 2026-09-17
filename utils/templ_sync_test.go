package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestTemplGeneratedInSync verifies that every committed *_templ.go file has
// import paths exactly matching its .templ source (both directions, excluding
// templ's auto-injected runtime imports). This catches the drift class where
// a source file's imports change but the generated file is not regenerated —
// or the reverse, where a daemon commits a generated file whose imports no
// longer match the source.
//
// The walk is REPO-WIDE (all modules: root, utils, icons, errorpage,
// charts/echarts, datastar, htmx, website, examples/demo). The 2026-09-17
// incident: website/internal/pages/base_templ.go was flipped to
// encoding/json/v2 while base.templ said v1 — invisible because the guard
// only walked the 8 root-module packages.
//
// Breadcrumbs bug (2026-07-28): source imported encoding/json (v1) but the
// generated file imported encoding/json/v2 — functionally inert under
// GOEXPERIMENT=jsonv2 but breaks consumers who run templ generate themselves.
func TestTemplGeneratedInSync(t *testing.T) {
	t.Parallel()

	templFiles := walkTemplFiles(t, "..")
	if len(templFiles) == 0 {
		t.Fatal("no .templ files found — glob root is wrong")
	}

	importRe := regexp.MustCompile(`"([^"]+)"`)

	for _, templFile := range templFiles {
		genFile := strings.TrimSuffix(templFile, ".templ") + "_templ.go"

		t.Run(templFile, func(t *testing.T) {
			t.Parallel()

			src, err := os.ReadFile(templFile)
			if err != nil {
				t.Fatalf("read %s: %v", templFile, err)
			}

			gen, err := os.ReadFile(genFile)
			if err != nil {
				t.Fatalf("read %s: %v (run `templ generate ./...` and commit)", genFile, err)
			}

			srcImports := extractImports(string(src), importRe)
			genImports := extractImports(string(gen), importRe)

			assertImportsMatch(t, templFile, genFile, srcImports, genImports)
		})
	}
}

// assertImportsMatch fails the test for every import present in one file but
// not the other, in both directions. Missing-in-generated means a stale
// generated file; missing-in-source means the generated file was hand-edited
// or generated from a different source. Extracted to keep the driving test's
// cognitive complexity under the gocognit gate.
func assertImportsMatch(t *testing.T, templFile, genFile string, srcImports, genImports map[string]bool) {
	t.Helper()

	onlyInSrc := missingFrom(srcImports, genImports)
	onlyInGen := missingFrom(genImports, srcImports)

	for _, imp := range onlyInSrc {
		t.Errorf(
			"%s imports %q but %s does not — run `templ generate ./...` to sync",
			filepath.Base(templFile), imp, filepath.Base(genFile),
		)
	}

	for _, imp := range onlyInGen {
		t.Errorf(
			"%s imports %q but its .templ source does not — the generated file is stale or hand-edited; run `templ generate ./...` to sync",
			filepath.Base(genFile), imp,
		)
	}
}

// missingFrom returns the sorted paths present in want but absent from have.
func missingFrom(want, have map[string]bool) []string {
	var missing []string

	for imp := range want {
		if !have[imp] {
			missing = append(missing, imp)
		}
	}

	sort.Strings(missing)

	return missing
}

// walkTemplFiles returns every .templ file under root that is expected to have
// a generated *_templ.go twin. Excluded: .git, node_modules, build output
// (dist), test fixtures, and cmd/tc/_sources (the `tc new` scaffolding
// sources — intentionally standalone templates with no generated twins, the
// same exclusion CI's tracked-files check uses).
func walkTemplFiles(t *testing.T, root string) []string {
	t.Helper()

	var files []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			name := d.Name()

			switch {
			case name == ".git" || name == "node_modules" || name == "dist" || name == "testdata":
				return filepath.SkipDir
			case path == filepath.Join(root, "cmd", "tc", "_sources"):
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(path, ".templ") {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}

	sort.Strings(files)

	return files
}

// extractImports pulls quoted import paths from the import section of a Go or
// templ source file. It scans only import (...) blocks and single-line
// import "..." declarations, excluding templ's own runtime imports
// (github.com/a-h/templ and github.com/a-h/templ/runtime) which are always
// injected by the generator.
func extractImports(src string, re *regexp.Regexp) map[string]bool {
	result := make(map[string]bool)

	importBlockRe := regexp.MustCompile(
		`(?s)import\s*\(([^)]*)\)|import\s*"([^"]+)"`,
	)

	for _, match := range importBlockRe.FindAllStringSubmatch(src, -1) {
		var block string

		if match[1] != "" {
			block = match[1]
		} else {
			block = match[2]
		}

		for _, imp := range re.FindAllStringSubmatch(block, -1) {
			path := imp[1]

			if strings.HasPrefix(path, "github.com/a-h/templ") {
				continue
			}

			result[path] = true
		}
	}

	return result
}

// TestTemplGeneratedInSyncCoverage confirms the repo-wide walk actually visits
// the modules whose drift shipped incidents (website: the 2026-09-17
// base_templ.go import flip; examples/demo: the 20:52 stale generated file),
// preventing a silent pass if the skip rules over-exclude after a refactor.
func TestTemplGeneratedInSyncCoverage(t *testing.T) {
	t.Parallel()

	templFiles := walkTemplFiles(t, "..")

	if len(templFiles) < 100 {
		t.Errorf("expected >=100 .templ files repo-wide (outside _sources), found %d", len(templFiles))
	}

	mustCover := []string{
		"website/internal/pages/",
		"examples/demo/",
		"charts/echarts/",
		"datastar/",
		"icons/",
		"display/",
	}

	for _, dir := range mustCover {
		found := false

		for _, f := range templFiles {
			if strings.Contains(f, dir) {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("walk missed .templ files under %q — skip rules over-exclude", dir)
		}
	}
}
