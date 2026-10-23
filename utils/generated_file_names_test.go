package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// generatedFileNameRe captures the FileName literal templ writes into every
// generated error path.
var generatedFileNameRe = regexp.MustCompile("FileName: `([^`]+)`")

// findOwningModuleRoot walks up from dir until it finds a go.mod.
func findOwningModuleRoot(t *testing.T, dir string) string {
	t.Helper()

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}

		dir = parent
	}
}

// TestGeneratedFileNameRoots guards the FileName literals templ writes into
// generated files. Each literal MUST resolve relative to the module root that
// owns the generated file:
//
//   - root-module files carry repo-root-relative paths (`display/x.templ`);
//   - submodules carry module-relative paths (`indicator.templ`,
//     `internal/pages/x.templ`).
//
// The 2026-10-23 incident is the reason this guard exists: running
// `templ generate` from the repo root with go.work active rewrote 28
// submodule paths workspace-relative (`datastar/indicator.templ` inside the
// datastar module), and the auto-commit daemon committed the corruption before
// anyone noticed. A literal that only resolves from the repo root is exactly
// that failure class.
func TestGeneratedFileNameRoots(t *testing.T) {
	t.Parallel()

	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	checked := 0

	err = filepath.WalkDir(repoRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".cache", "node_modules", "dist", "testdata", "artifacts", "results":
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(path, "_templ.go") {
			return nil
		}

		moduleRoot := findOwningModuleRoot(t, filepath.Dir(path))
		if moduleRoot == "" {
			t.Errorf("%s: no go.mod found above this generated file", path)

			return nil
		}

		source, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Errorf("%s: read: %v", path, readErr)

			return nil
		}

		relFile, relErr := filepath.Rel(repoRoot, path)
		if relErr != nil {
			relFile = path
		}

		relModule, relModuleErr := filepath.Rel(repoRoot, moduleRoot)
		if relModuleErr != nil {
			relModule = moduleRoot
		}

		for _, match := range generatedFileNameRe.FindAllStringSubmatch(string(source), -1) {
			literal := match[1]
			target := filepath.Join(moduleRoot, filepath.FromSlash(literal))

			if _, statErr := os.Stat(target); statErr != nil {
				t.Errorf("%s: FileName %q does not resolve from module root %q — "+
					"a repo-root-relative path in a submodule is the root-generate corruption class",
					relFile, literal, relModule)
			}

			checked++
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk repo: %v", err)
	}

	if checked == 0 {
		t.Fatal("no FileName literals found — the guard would be vacuously green")
	}
}
