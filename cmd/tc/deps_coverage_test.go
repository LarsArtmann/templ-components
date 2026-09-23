package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPackageDepsCoverPackageFiles keeps the --list-deps output honest
// (backlog #285, 2026-09-23): every non-test, non-generated .go file in a
// mirrored package must be either listed in packageDeps[pkg] (a sibling
// `tc add` does not copy) or be a *_types.go companion (which `tc add`
// copies automatically next to its .templ). The kanban, chart,
// calendar-nav, and heading-tag helper files all predate this test and were
// silently missing from `tc add <component> --list-deps`.
func TestPackageDepsCoverPackageFiles(t *testing.T) {
	t.Parallel()

	for _, pkg := range mirroredPackages {
		entries, err := os.ReadDir(filepath.Join("..", "..", pkg))
		if err != nil {
			t.Fatalf("read package dir %s: %v", pkg, err)
		}

		listed := make(map[string]bool)
		for _, dep := range packageDeps[pkg] {
			listed[dep] = true
		}

		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() || !strings.HasSuffix(name, ".go") {
				continue
			}

			switch {
			case strings.HasSuffix(name, "_test.go"),
				strings.HasSuffix(name, "_templ.go"),
				strings.HasSuffix(name, "_types.go"),
				name == "doc.go":
				continue
			}

			if !listed[name] {
				t.Errorf(
					"%s/%s is not listed in packageDeps[%q] — add it (the deps listing must match what a vendored component needs)",
					pkg,
					name,
					pkg,
				)
			}
		}

		// Listed entries must also exist on disk — a stale dep name is as
		// misleading as a missing one.
		for _, dep := range packageDeps[pkg] {
			if _, err := os.Stat(filepath.Join("..", "..", pkg, dep)); err != nil {
				t.Errorf("packageDeps[%q] lists %q which does not exist", pkg, dep)
			}
		}
	}
}
