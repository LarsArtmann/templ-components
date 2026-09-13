package utils

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// compiledCSSDistributionTargets lists the ONLY compiled CSS artifacts the
// repo tracks. Each has a consumer: static/app.css is go:embed'ed by the
// demo binary; templates/styles.css and templates/templ-components-theme.out.css
// are pre-built consumer artifacts refreshed by scripts/release.sh.
//
// Everything else with an .out.css suffix was dead weight: compiled output of
// a `tc new` preset scaffolder that never shipped, byte-identical duplicates
// of static/app.css, and website output the Astro/Vite build regenerates.
// They were removed 2026-09-13 after audit (see AGENTS.md "CSS artifact
// inventory"). This guard keeps the set at exactly these three — the
// BuildFlow tailwind-build daemon has repeatedly resurrected stray compiled
// CSS in the working tree, and without this test the zombies would silently
// re-enter master.
func TestCompiledCSSInventory(t *testing.T) {
	t.Parallel()

	repoRoot := ".."

	wantOut := []string{
		"templates/templ-components-theme.out.css",
	}

	var gotOut []string

	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() {
			// Website builds its CSS via @tailwindcss/vite; dist/ is
			// gitignored build output.
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

			gotOut = append(gotOut, rel)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk repo for .out.css files: %v", err)
	}

	sort.Strings(gotOut)
	sort.Strings(wantOut)

	if strings.Join(gotOut, ",") != strings.Join(wantOut, ",") {
		t.Fatalf(
			"committed .out.css set drifted.\n want: %v\n got:  %v\n\nIf you added a compiled artifact intentionally, add it to compiledCSSDistributionTargets AND give it an in-repo consumer or a release.sh compile line — artifacts nobody consumes are daemon-recompile churn bait.",
			wantOut,
			gotOut,
		)
	}

	for _, mustExist := range []string{
		filepath.Join(repoRoot, "templates", "styles.css"),
		filepath.Join(repoRoot, "examples", "demo", "static", "app.css"),
	} {
		if _, statErr := os.Stat(mustExist); statErr != nil {
			t.Errorf(
				"distribution target missing: %s (%v) — restore it or update compiledCSSDistributionTargets, scripts/release.sh, and AGENTS.md together",
				mustExist,
				statErr,
			)
		}
	}
}
