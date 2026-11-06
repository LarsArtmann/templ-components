package utils

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestSemanticColorTokensCompiled extends TestCustomClassTokensCompiled from
// tc-* hooks to semantic color classes: every literal `<utility>-<hue>-<shade>`
// token written into a .templ class attribute must be RESOLVED by the compiled
// demo CSS (examples/demo/static/app.css, which @source-scans every repo
// .templ file). Tailwind v4 emits a color utility only when its theme color
// resolves, so even a freshly compiled stylesheet silently produces nothing
// for a family the consumer theme does not define — the exact mechanism
// behind the invisible accent focus ring the token-axis sweep shipped
// (ADR-0049): the class string was present, the rule never was.
//
// Substring matching against the minified CSS is deliberate, the same choice
// the tc-* guard makes: variant selectors (dark:, focus-visible:,
// peer-checked:*, @sm:…) embed the bare token, and an opacity stem (a token
// with a slash and an alpha percentage, like the accent-600 bar tint) still
// matches the unslashed selector fragment. False
// negatives (dynamically composed class strings) stay out of scope — the
// authoritative full check remains `scripts/ci-repro.sh --css`.
func TestSemanticColorTokensCompiled(t *testing.T) {
	t.Parallel()

	// unresolvedByDesign lists color tokens that are deliberately absent from
	// the compiled CSS, each with the reason. Everything else must resolve.
	unresolvedByDesign := map[string]string{}

	// colorTokenRe anchors on the canonical hue families (tone.Hue + the
	// accent token family) and a shade: `[a-z-]+-(hue)-(shade)`. Anchoring on
	// families kills the numeric false-positive class (z-50, duration-700,
	// opacity-60) because those never carry a hue name.
	colorTokenRe := regexp.MustCompile(
		`\b[a-z]+[a-z-]*-(?:accent|gray|green|red|yellow|amber|orange|blue|emerald|rose|violet|cyan|pink|purple)-` +
			`(?:50|100|200|300|400|500|600|700|800|900|950)\b`)

	css, err := os.ReadFile("../examples/demo/static/app.css")
	if err != nil {
		t.Skipf("compiled demo CSS not found: %v", err)
	}

	cssText := string(css)

	// class="..." (HTML attribute) and class={ "..." } (templ expression),
	// same extraction contract as TestCustomClassTokensCompiled.
	classAttrRe := regexp.MustCompile(`class=(?:\{[[:space:]]*)?"([^"]*)"?`)

	templFiles := walkTemplFiles(t)

	seen := make(map[string]bool)

	var missing []string

	for _, file := range templFiles {
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}

		for _, attr := range classAttrRe.FindAllStringSubmatch(string(src), -1) {
			for _, token := range colorTokenRe.FindAllString(attr[1], -1) {
				if seen[token] {
					continue
				}

				seen[token] = true

				if strings.Contains(cssText, token) {
					continue
				}

				if _, ok := unresolvedByDesign[token]; ok {
					continue
				}

				missing = append(missing, token+" ("+file+")")
			}
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)

		t.Errorf(
			"%d semantic color class(es) referenced in .templ class attributes but not resolved by examples/demo/static/app.css — the rule was never emitted (undefined theme family or stale compile): recompile with `nix run .#css`, define the family in templates/custom.css, or fix the typo:\n%s",
			len(missing),
			strings.Join(missing, "\n"),
		)
	}

	if len(seen) < 20 {
		t.Errorf("expected >=20 distinct color tokens across sources, found %d — the scan is likely broken", len(seen))
	}
}
