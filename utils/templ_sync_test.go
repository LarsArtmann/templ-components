package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestTemplGeneratedInSync verifies that every committed *_templ.go file has
// import paths matching its .templ source. This catches the drift class where
// a source file's imports change (e.g., encoding/json → encoding/json/v2) but
// the generated file is not regenerated, producing a stale artifact.
//
// The breadcrumbs_templ.go bug (2026-07-28): source imported encoding/json (v1)
// but the generated file imported encoding/json/v2 — functionally inert under
// GOEXPERIMENT=jsonv2 but breaks consumers who run templ generate themselves.
func TestTemplGeneratedInSync(t *testing.T) {
	t.Parallel()

	root := ".."

	packages := []string{
		"display", "errorpage", "feedback", "forms",
		"htmx", "layout", "navigation", "recipes",
	}

	importRe := regexp.MustCompile(`"([^"]+)"`)

	for _, pkg := range packages {
		pkgDir := filepath.Join(root, pkg)

		templFiles, err := filepath.Glob(filepath.Join(pkgDir, "*.templ"))
		if err != nil {
			t.Fatalf("glob %s: %v", pkgDir, err)
		}

		for _, templFile := range templFiles {
			genFile := strings.TrimSuffix(templFile, ".templ") + "_templ.go"

			t.Run(pkg+"/"+filepath.Base(templFile), func(t *testing.T) {
				t.Parallel()

				src, err := os.ReadFile(templFile)
				if err != nil {
					t.Fatalf("read %s: %v", templFile, err)
				}

				gen, err := os.ReadFile(genFile)
				if err != nil {
					t.Fatalf("read %s: %v", genFile, err)
				}

				srcImports := extractImports(string(src), importRe)
				genImports := extractImports(string(gen), importRe)

				for imp := range srcImports {
					if genImports[imp] {
						continue
					}

					t.Errorf(
						"%s imports %q but %s does not — run `templ generate` to sync",
						filepath.Base(templFile),
						imp,
						filepath.Base(genFile),
					)
				}
			})
		}
	}
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

// TestTemplGeneratedInSyncCoverage confirms the test actually visits files,
// preventing a silent pass if the glob patterns break after a refactor.
func TestTemplGeneratedInSyncCoverage(t *testing.T) {
	t.Parallel()

	root := ".."

	packages := []string{
		"display", "errorpage", "feedback", "forms",
		"htmx", "layout", "navigation", "recipes",
	}

	count := 0

	for _, pkg := range packages {
		files, _ := filepath.Glob(filepath.Join(root, pkg, "*.templ"))

		count += len(files)
	}

	if count < 50 {
		t.Errorf("expected >=50 .templ files across packages, found %d", count)
	}
}
