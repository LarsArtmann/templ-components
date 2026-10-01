package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// inlineStyleAllowlist names the library .templ files permitted to emit an
// inline `style=` attribute, each with the reason it exists. These carry
// runtime data (a width, height, colour, or CSS custom property) that has no
// static Tailwind class, so they are documented CSP exemptions — see
// docs/tailwind-v4-adoption-guide.md, "What about CSP?".
//
// Any OTHER .templ file that emits `style=` fails TestInlineStyleCompliance:
// under a strict `style-src` policy (`style-src 'self'` or a nonce) the inline
// attribute is dropped silently, so the library must not grow new inline styles
// by accident. Add a file here only with a matching docs update.
var inlineStyleAllowlist = map[string]string{
	filepath.Join("layout", "appshell.templ"):      "the --tc-sidebar-w custom property (SidebarWidth enum)",
	filepath.Join("display", "bar_chart.templ"):    "runtime bar height/width",
	filepath.Join("display", "heatmap.templ"):      "runtime cell background",
	filepath.Join("feedback", "progressbar.templ"): "runtime bar width",
	filepath.Join("feedback", "loading.templ"):     "runtime overlay progress width",
}

// inlineStyleRe matches a `style=` attribute in either templ form: the Go
// expression `style={ ... }` or the static string `style="..."`.
var inlineStyleRe = regexp.MustCompile(`\bstyle\s*=\s*\{|\bstyle\s*=\s*"`)

// skippedStyleDirs are subtrees that are not part of the published component
// surface: the demo (allows inline styles for layout scaffolding), the website
// generator, vendor/.git, the `cmd/tc/_sources` mirror (guarded separately),
// and build output.
var skippedStyleDirs = map[string]bool{
	"examples":     true,
	"website":      true,
	"vendor":       true,
	"node_modules": true,
	"_sources":     true,
	"dist":         true,
	"testdata":     true,
	".git":         true,
}

// TestInlineStyleCompliance asserts that only the documented exemption files
// emit inline `style=` attributes, and that each exemption still exists and
// still has one (so a removed exemption is pruned rather than left stale).
func TestInlineStyleCompliance(t *testing.T) {
	t.Parallel()

	root := ".."
	seen := make(map[string]bool)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if skippedStyleDirs[info.Name()] {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(path, ".templ") {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}

		if !inlineStyleRe.MatchString(string(data)) {
			return nil
		}

		if _, allowed := inlineStyleAllowlist[rel]; allowed {
			seen[rel] = true

			return nil
		}

		t.Errorf("inline style= in %s: not an allowlisted CSP exemption (document it first)", rel)

		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}

	for rel, reason := range inlineStyleAllowlist {
		if !seen[rel] {
			t.Errorf("inline style= allowlist entry %q (%s) matched no style= attribute; prune it", rel, reason)
		}
	}
}
