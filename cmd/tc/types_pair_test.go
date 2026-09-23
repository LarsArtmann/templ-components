package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTypesFilesHaveTemplTwin catches ghost type files: a *_types.go whose
// component's .templ was deleted or renamed. The types file then ships
// dead types in the published module AND — worse — the scaffolder mirror
// keeps embedding it as a companion for a component that no longer exists
// (the 2026-09-22 pass found nine *_types.go copies whose components had
// never been embedded; the inverse ghost had zero coverage until this test,
// backlog item #283b). The reverse direction (a .templ without a *_types.go) is
// legal: small components inline their props next to the component.
func TestTypesFilesHaveTemplTwin(t *testing.T) {
	t.Parallel()

	for _, pkg := range mirroredPackages {
		entries, err := os.ReadDir(filepath.Join("..", "..", pkg))
		if err != nil {
			t.Fatalf("read package dir %s: %v", pkg, err)
		}

		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() || !strings.HasSuffix(name, "_types.go") {
				continue
			}

			twin := strings.TrimSuffix(name, "_types.go") + ".templ"
			if _, err := os.Stat(filepath.Join("..", "..", pkg, twin)); err != nil {
				t.Errorf("%s/%s has no .templ twin (%s missing) — delete or rename the types file", pkg, name, twin)
			}
		}
	}
}
