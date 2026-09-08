package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestComponentSourcesPopulated(t *testing.T) {
	t.Parallel()

	r := newRegistry()
	if len(r.files) == 0 {
		t.Fatal("registry files is empty — embed failed")
	}

	if _, ok := r.files["button"]; !ok {
		t.Error("expected 'button' in registry files")
	}
}

// TestSourcesMatchPackageFiles is the _sources/ drift guard: every embedded
// source file must be byte-identical to the real package file it was copied
// from. The sources were copied by hand; without this guard the scaffolder
// silently serves stale components after a package edit.
func TestSourcesMatchPackageFiles(t *testing.T) {
	t.Parallel()

	repoRoot := filepath.Join("..", "..")

	err := fs.WalkDir(sourcesFS, "_sources", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return walkErr
		}

		rel, _ := filepath.Rel("_sources", path)
		parts := strings.Split(rel, string(filepath.Separator))

		if parts[0] == "starter" {
			// Starter CSS is a curated template, not a copy of a package file.
			return nil
		}

		if rel == filepath.Join("datastar", datastarBumpProtocolDoc) {
			// Scaffolder-owned checklist — no package-file counterpart.
			return nil
		}

		embedded, err := sourcesFS.ReadFile(path)
		if err != nil {
			t.Errorf("read embedded %s: %v", rel, err)

			return nil
		}

		original, err := os.ReadFile(filepath.Join(repoRoot, rel))
		if err != nil {
			t.Errorf("package file for embedded source missing: %s: %v", rel, err)

			return nil
		}

		if string(embedded) != string(original) {
			t.Errorf("embedded %s drifted from the package file — re-copy it (and re-run this test)", rel)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk _sources: %v", err)
	}
}

// TestPackageImportsMatchSources guards the packageImports checklist printed
// by 'tc add --list-deps': it must equal the set of module-level imports the
// package's non-test, non-generated sources actually use.
func TestPackageImportsMatchSources(t *testing.T) {
	t.Parallel()

	for pkg := range packageImports {
		got := collectPackageImports(t, pkg)
		want := map[string]bool{}

		for _, imp := range packageImports[pkg] {
			want[imp] = true
		}

		for imp := range want {
			if !got[imp] {
				t.Errorf("packageImports[%q] lists %q but %s does not import it — prune the checklist", pkg, imp, pkg)
			}
		}

		for imp := range got {
			if !want[imp] {
				t.Errorf("%s imports %q but packageImports[%q] misses it — add it", pkg, imp, pkg)
			}
		}
	}
}

// collectPackageImports scans a package's non-test, non-generated .go and
// .templ sources for module-level imports (the go.mod surface a vendoring
// consumer pulls in). Test-only helpers (/golden) and the package's own path
// are excluded.
func collectPackageImports(t *testing.T, pkg string) map[string]bool {
	t.Helper()

	const selfPrefix = "github.com/larsartmann/templ-components/"

	got := map[string]bool{}

	for _, pattern := range []string{"*.go", "*.templ"} {
		files, err := filepath.Glob(filepath.Join("..", "..", pkg, pattern))
		if err != nil {
			t.Fatalf("glob %s %s: %v", pkg, pattern, err)
		}

		for _, file := range files {
			base := filepath.Base(file)
			if strings.HasSuffix(base, "_test.go") || strings.HasSuffix(base, "_templ.go") {
				continue
			}

			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("read %s: %v", file, err)
			}

			for _, imp := range importPattern.FindAllString(string(content), -1) {
				imp = strings.Trim(imp, `"`)
				if strings.Contains(imp, "/golden") || imp == selfPrefix+pkg {
					continue
				}

				got[imp] = true
			}
		}
	}

	return got
}

var importPattern = regexp.MustCompile(`"github\.com/[^"]+"`)

func TestCmdAddDatastarIncludesBumpProtocol(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	cmdAdd(newRegistry(), []string{"live_region", "--out", tmp})

	if _, err := os.Stat(filepath.Join(tmp, datastarBumpProtocolDoc)); err != nil {
		t.Errorf("%s not copied next to datastar sources: %v", datastarBumpProtocolDoc, err)
	}

	// Non-datastar adds must NOT get the checklist.
	tmp2 := t.TempDir()
	cmdAdd(newRegistry(), []string{"button", "--out", tmp2})

	if _, err := os.Stat(filepath.Join(tmp2, datastarBumpProtocolDoc)); err == nil {
		t.Errorf("%s copied for a non-datastar component", datastarBumpProtocolDoc)
	}
}

func TestCmdAddCopiesTemplAndTypes(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	cmdAdd(newRegistry(), []string{"button", "--out", tmp})

	if _, err := os.Stat(filepath.Join(tmp, "button.templ")); err != nil {
		t.Errorf("button.templ not copied: %v", err)
	}
}

func TestCmdAddUnknownComponent(t *testing.T) {
	t.Parallel()

	r := newRegistry()
	if _, ok := r.files["totally-bogus-component"]; ok {
		t.Error("did not expect bogus component in sources")
	}
}

func TestCmdListDoesNotCrash(t *testing.T) {
	t.Parallel()

	cmdList(newRegistry(), nil)
}
