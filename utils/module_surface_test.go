package utils

import (
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

// TestGoWorkDirectiveMatchesRootGoMod guards the go.work ↔ go.mod version
// sync: the workspace directive was left unguarded after the 1.26.7 bump
// (2026-09-02 f15). A go.work pinned to an older toolchain makes workspace
// builds download (or reject) a different Go than per-module CI runs use,
// splitting local vs CI behavior.
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

	if string(workMatch[1]) != string(modMatch[1]) {
		t.Errorf(
			"go.work go directive (%s) != root go.mod go directive (%s) — align them (bump go.work when the toolchain input moves)",
			workMatch[1], modMatch[1],
		)
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
			got, want,
		)
	}
}
