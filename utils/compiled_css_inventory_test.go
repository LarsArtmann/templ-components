package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
)

// compiledCSSTargetsFile is the single source of truth for compiled CSS
// distribution targets — shared with scripts/release.sh (recompile) and
// scripts/check-css-minified.sh (minification guard). Format: one
// "<input> <output>" pair per repo-relative line, # comments allowed.
const compiledCSSTargetsFile = "../scripts/compiled-css-targets.txt"

// readCompiledCSSTargets parses the shared target list and returns the
// compiled output paths. It fails the test loudly when the fixture is
// missing — a guard that owns its fixture never skips (AGENTS.md rule).
func readCompiledCSSTargets(t *testing.T) []string {
	t.Helper()

	data, err := os.ReadFile(compiledCSSTargetsFile)
	if err != nil {
		t.Fatalf("read %s: %v — the inventory guard owns this fixture (fail loud, never skip)", compiledCSSTargetsFile, err)
	}

	var outputs []string

	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			t.Fatalf("%s: malformed line %q — want exactly \"<input> <output>\"", compiledCSSTargetsFile, line)
		}

		outputs = append(outputs, filepath.ToSlash(parts[1]))
	}

	if len(outputs) == 0 {
		t.Fatalf("%s lists no targets — the compiled-CSS guards are dead without entries", compiledCSSTargetsFile)
	}

	return outputs
}

// TestCompiledCSSInventory pins the tracked compiled-CSS artifact set.
//
// Everything else with an .out.css suffix was dead weight: compiled output of
// a `tc new` preset scaffolder that never shipped, byte-identical duplicates
// of static/app.css, and website output the Astro/Vite build regenerates.
// They were removed 2026-09-13 after audit (see AGENTS.md "CSS artifact
// inventory"). This guard keeps the tracked set at exactly the targets that
// scripts/compiled-css-targets.txt lists — the BuildFlow tailwind-build
// daemon repeatedly resurrects stray compiled CSS in the working tree
// (proven 2026-09-13: it re-added six deleted artifacts within hours), so
// the invariant is asserted against `git ls-files` (what consumers actually
// receive), not the worktree, where ignored litter reappears on every daemon
// cycle.
func TestCompiledCSSInventory(t *testing.T) {
	t.Parallel()

	repoRoot := ".."

	targets := readCompiledCSSTargets(t)

	var wantOut []string
	for _, target := range targets {
		if strings.HasSuffix(target, ".out.css") {
			wantOut = append(wantOut, target)
		}
	}

	gitOut, err := exec.CommandContext(t.Context(), "git", "-C", repoRoot, "ls-files", "--", "*.out.css").Output()
	if err != nil {
		t.Fatalf("git ls-files *.out.css: %v — the tracked-artifact guard requires git (fail loud, never skip)", err)
	}

	var gotOut []string

	for line := range strings.SplitSeq(strings.TrimSpace(string(gitOut)), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			gotOut = append(gotOut, filepath.ToSlash(line))
		}
	}

	sort.Strings(gotOut)
	sort.Strings(wantOut)

	if strings.Join(gotOut, ",") != strings.Join(wantOut, ",") {
		t.Fatalf(
			"tracked .out.css set drifted.\n want: %v\n got:  %v\n\nIf you added a compiled artifact intentionally, add it to scripts/compiled-css-targets.txt AND give it an in-repo consumer or a release.sh compile line — artifacts nobody consumes are daemon-recompile churn bait.",
			wantOut,
			gotOut,
		)
	}

	var litter []string

	err = filepath.Walk(repoRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == "dist" {
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(path, ".out.css") {
			rel, relErr := filepath.Rel(repoRoot, path)
			if relErr != nil {
				return relErr
			}

			if !slices.Contains(wantOut, filepath.ToSlash(rel)) {
				litter = append(litter, rel)
			}
		}

		return nil
	})
	if err != nil {
		t.Logf("walk repo for untracked .out.css litter: %v", err)
	}

	for _, path := range litter {
		t.Logf(
			"untracked .out.css litter in worktree (gitignored, daemon-recompile residue — safe to delete): %s",
			path,
		)
	}

	for _, mustExist := range targets {
		if _, statErr := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(mustExist))); statErr != nil {
			t.Errorf(
				"distribution target missing: %s (%v) — restore it or update scripts/compiled-css-targets.txt (single source, shared with release.sh and check-css-minified.sh)",
				mustExist,
				statErr,
			)
		}
	}
}
