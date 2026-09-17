package utils

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestCustomClassTokensCompiled is the content-based half of the CSS
// freshness story (TODO #263): every `tc-*` hook class written in a .templ
// class attribute MUST appear in the compiled demo CSS. These classes are
// hand-written in templates/custom.css (component hooks like tc-modal,
// tc-kanban-pending), so a missing token means either a stale compiled CSS
// or a typo'd hook — both shipped silently when the old check relied only
// on file mtimes (git checkouts and daemon commits reset mtimes, so the
// signal flipped randomly).
//
// Extraction is scoped to class="..." / class={ "..." } attribute values:
// tokens appearing elsewhere are deliberately out of scope — data-tc-* are
// hook ATTRIBUTE names, tc-error-announcer / tc-toast-container are element
// IDs, and JS-applied classes (tc-menu-open) are runtime state styled
// directly in custom.css.
//
// Arbitrary Tailwind utilities are NOT checked here (without a real Tailwind
// build that over-produces false positives); the authoritative full check
// remains `scripts/ci-repro.sh --css`, which recompiles and byte-diffs.
func TestCustomClassTokensCompiled(t *testing.T) {
	t.Parallel()

	// selectorOnlyHooks are emitted classes that are deliberately unstyled:
	// stable public selectors consumers target (the justification lives in
	// the referenced doc). Everything else must exist in the compiled CSS.
	selectorOnlyHooks := map[string]string{
		"tc-btn-loading": "deliberately unstyled hook for hx-indicator scoping — see htmx/doc.go",
	}

	css, err := os.ReadFile("../examples/demo/static/app.css")
	if err != nil {
		t.Skipf("compiled demo CSS not found: %v", err)
	}

	cssText := string(css)

	// class="..." (HTML attribute) and class={ "..." } (templ expression).
	classAttrRe := regexp.MustCompile(`class=(?:\{[[:space:]]*)?"([^"]*)"?`)
	tokenRe := regexp.MustCompile(`\btc-[a-z0-9]+(?:-[a-z0-9]+)*\b`)

	templFiles := walkTemplFiles(t, "..")

	seen := make(map[string]bool)

	var missing []string

	for _, file := range templFiles {
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}

		for _, attr := range classAttrRe.FindAllStringSubmatch(string(src), -1) {
			for _, token := range tokenRe.FindAllString(attr[1], -1) {
				if seen[token] {
					continue
				}

				seen[token] = true

				if !strings.Contains(cssText, token) {
					if _, ok := selectorOnlyHooks[token]; ok {
						continue
					}

					missing = append(missing, token+" ("+file+")")
				}
			}
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)

		t.Errorf(
			"%d tc-* class token(s) referenced in .templ class attributes but absent from examples/demo/static/app.css — recompile the CSS (nix run .#css) or fix the typo:\n%s",
			len(missing),
			strings.Join(missing, "\n"),
		)
	}

	if len(seen) < 5 {
		t.Errorf("expected >=5 distinct tc-* tokens across sources, found %d — the scan is likely broken", len(seen))
	}
}
