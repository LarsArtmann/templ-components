package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// goDirectivePattern extracts the `go <version>` toolchain directive line
// from a go.mod / go.work file.
var goDirectivePattern = regexp.MustCompile(`(?m)^go (\d+\.\d+(?:\.\d+)?)\s*$`)

// normalizeGoDirective reduces a go directive to major.minor: the repo
// canonical form is `go 1.26.0` (the go1.26 toolchain's own tidy
// canonicalization — the 2026-10-03 e2e outage proved the toolchain sorts
// `1.26` BELOW `1.26.0` and rejects mixed directives), but historical trees
// still carry bare `1.26` from before the normalization. Comparing
// normalized forms keeps the historical spellings passing instead of
// failing CI.
func normalizeGoDirective(v string) string {
	if dot := strings.IndexByte(v, '.'); dot >= 0 {
		if second := strings.IndexByte(v[dot+1:], '.'); second >= 0 {
			return v[:dot+1+second]
		}
	}

	return v
}

// TestGoWorkDirectiveMatchesRootGoMod guards the go.work ↔ go.mod version
// sync: the workspace directive was left unguarded after the 1.26.7 bump
// (2026-09-02 f15). A go.work pinned to an older toolchain makes workspace
// builds download (or reject) a different Go than per-module CI runs use,
// splitting local vs CI behavior. Comparison is major.minor-normalized
// (see normalizeGoDirective).
func TestGoWorkDirectiveMatchesRootGoMod(t *testing.T) {
	t.Parallel()

	goWork, err := os.ReadFile(filepath.Join("..", "go.work"))
	if err != nil {
		t.Skipf("go.work not present (release tags strip it? expected in dev): %v", err)
	}

	goMod, err := os.ReadFile(filepath.Join("..", "go.mod"))
	if err != nil {
		t.Fatalf("read root go.mod: %v", err)
	}

	workMatch := goDirectivePattern.FindSubmatch(goWork)
	if workMatch == nil {
		t.Fatalf("could not parse `go <version>` directive from go.work:\n%s", goWork)
	}

	modMatch := goDirectivePattern.FindSubmatch(goMod)
	if modMatch == nil {
		t.Fatalf("could not parse `go <version>` directive from root go.mod")
	}

	if normalizeGoDirective(string(workMatch[1])) != normalizeGoDirective(string(modMatch[1])) {
		t.Errorf(
			"go.work go directive (%s) != root go.mod go directive (%s) — align them (canonical form: `go 1.26.0`; bump go.work when the toolchain input moves)",
			workMatch[1],
			modMatch[1],
		)
	}
}

// TestGoDirectivesAlignAcrossWorkspace extends the go.work ↔ root guard to
// EVERY module in the repo (root, the 7 published sub-modules, visualtest,
// website): all go directives must agree at major.minor resolution. This is
// the normalized-compare guard (2026-10-02): historical `go 1.26` spellings
// pass next to the canonical `go 1.26.0`, but a module left on a genuinely
// different toolchain (the 2026-09-17 1.27.1 incident) fails.
func TestGoDirectivesAlignAcrossWorkspace(t *testing.T) {
	t.Parallel()

	root, err := os.ReadFile(filepath.Join("..", "go.mod"))
	if err != nil {
		t.Fatalf("read root go.mod: %v", err)
	}

	rootMatch := goDirectivePattern.FindSubmatch(root)
	if rootMatch == nil {
		t.Fatalf("could not parse `go <version>` directive from root go.mod")
	}

	want := normalizeGoDirective(string(rootMatch[1]))

	var misaligned []string

	err = filepath.WalkDir("..", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			switch filepath.Base(path) {
			case ".git", "node_modules", "result", "testdata":
				return filepath.SkipDir
			}

			return nil
		}

		if filepath.Base(path) != "go.mod" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		match := goDirectivePattern.FindSubmatch(content)
		if match == nil {
			t.Errorf("%s: no `go <version>` directive", path)

			return nil
		}

		if normalizeGoDirective(string(match[1])) != want {
			misaligned = append(misaligned, fmt.Sprintf("%s: %s (want %s.x)", path, match[1], want))
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk modules: %v", err)
	}

	for _, m := range misaligned {
		t.Error(m)
	}
}

// TestGoDatastarStaticPinSurface pins the dependency SURFACE of
// go-datastar/static: exactly 3 go.mod files may reference it — datastar
// (direct require), root + visualtest (indirect siblings). A fourth mention
// means a module silently gained a Datastar dependency (consumer graph
// pollution); a missing one means the release script's lockstep bump will
// desync that module's go.sum.
func TestGoDatastarStaticPinSurface(t *testing.T) {
	t.Parallel()

	const dep = "github.com/larsartmann/go-datastar/static"

	want := []string{"datastar/go.mod", "go.mod", "visualtest/go.mod"}
	sort.Strings(want)

	mods, err := filepath.Glob(filepath.Join("..", "*", "go.mod"))
	if err != nil {
		t.Fatalf("glob go.mods: %v", err)
	}

	mods = append(mods, filepath.Join("..", "go.mod"))

	var got []string

	for _, mod := range mods {
		content, err := os.ReadFile(mod)
		if err != nil {
			t.Fatalf("read %s: %v", mod, err)
		}

		if strings.Contains(string(content), dep) {
			rel, err := filepath.Rel("..", mod)
			if err != nil {
				t.Fatalf("rel %s: %v", mod, err)
			}

			got = append(got, filepath.ToSlash(rel))
		}
	}

	sort.Strings(got)

	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf(
			"go-datastar/static pin surface drifted:\n got: %v\nwant: %v\n(direct: datastar; indirect: root, visualtest — a change here is a release-script + consumer-graph decision, not a mechanical bump)",
			got,
			want,
		)
	}
}
