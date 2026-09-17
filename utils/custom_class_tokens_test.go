package utils

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestCustomClassTokensCompiled is the content-based half of the CSS
// freshness story (TODO #263): every `tc-*` hook class referenced in a
// .templ source MUST appear in the compiled demo CSS. These classes are
// hand-written in templates/custom.css (component hooks like tc-modal,
// tc-kanban-pending), so a missing token means either a stale compiled CSS
// or a typo'd hook — both shipped silently when the old check relied only
// on file mtimes (git checkouts and daemon commits reset mtimes, so the
// signal flipped randomly).
//
// Arbitrary Tailwind utilities are NOT checked here (constructing them
// without a real Tailwind build over-produces false positives); the
// authoritative full check remains `scripts/ci-repro.sh --css`, which
// recompiles and byte-diffs.
func TestCustomClassTokensCompiled(t *testing.T) {
	t.Parallel()

	css, err := os.ReadFile("../examples/demo/static/app.css")
	if err != nil {
		t.Skipf("compiled demo CSS not found: %v", err)
	}

	cssText := string(css)

	tokenRe := regexp.MustCompile(`\btc-[a-z0-9]+(?:-[a-z0-9]+)*\b`)
	// data-tc-* are HOOK ATTRIBUTE names, not classes — strip them before
	// tokenizing so they don't false-positive as missing class selectors.
	dataAttrRe := regexp.MustCompile(`data-tc-[a-z0-9-]+`)

	templFiles := walkTemplFiles(t, "..")

	seen := make(map[string]bool)

	var missing []string

	for _, file := range templFiles {
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}

		for _, token := range tokenRe.FindAllString(string(src), -1) {
			if seen[token] {
				continue
			}

			seen[token] = true

			// The class must appear in the compiled CSS either as its
			// escaped selector (.tc-kanban-pending) or inside a custom.css
			// at-rule block. Match on the bare token with a non-word
			// boundary handled by the token regex already excluding
			// longer names.
			if !strings.Contains(cssText, token) {
				missing = append(missing, token+" ("+file+")")
			}
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)

		t.Errorf("%d tc-* class token(s) referenced in .templ sources but absent from examples/demo/static/app.css — recompile the CSS (nix run .#css) or fix the typo:\n%s",
			len(missing), strings.Join(missing, "\n"))
	}

	if len(seen) < 10 {
		t.Errorf("expected >=10 distinct tc-* tokens across sources, found %d — the scan is likely broken", len(seen))
	}
}
