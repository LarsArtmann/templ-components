package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestCoarsePointerCompliance verifies that hover-revealed FUNCTIONALITY in
// .templ source files has a coarse-pointer (touch) fallback: there is no
// hover on touch devices, so controls hidden until group-hover/peer-hover
// are unreachable there.
//
// The compliant pattern (shipped by KanbanBoard): the hover-revealed element
// carries a `tc-*` hook class AND templates/custom.css covers that hook in
// an `@media (pointer: coarse)` block that keeps it visible. Anything else
// needs an entry in coarsePointerExemptions with a reason.
//
// Pure color accents (group-hover:text-*) and request-state reveals
// (htmx-indicator) are not hover-gated functionality and never match.
//
// This is a FAILING test — violations block CI.
//
// Run via: go test ./utils/... -run TestCoarsePointer.
func TestCoarsePointerCompliance(t *testing.T) {
	t.Parallel()

	// Elements revealed only by hover (opacity or display based).
	revealRe := regexp.MustCompile(
		`(opacity-0.*(group|peer)-hover:opacity-100|(group|peer)-hover:opacity-100.*opacity-0)` +
			`|(hidden.*(group|peer)-hover:block|(group|peer)-hover:block.*hidden)`,
	)

	hookRe := regexp.MustCompile(`(tc-[a-z0-9-]+)`)

	coarseHooks := loadCoarsePointerHooks(t)

	exempt := coarsePointerExemptions()

	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx"}

	for _, dir := range dirs {
		err := filepath.Walk(filepath.Join("..", dir), func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}

			if !strings.HasSuffix(path, ".templ") {
				return nil
			}

			rel, relErr := filepath.Rel("..", path)
			if relErr != nil {
				return relErr
			}

			data, readErr := os.ReadFile(path) //nolint:gosec // test scans templ files
			if readErr != nil {
				return fmt.Errorf("read file: %w", readErr)
			}

			for line := range strings.SplitSeq(string(data), "\n") {
				if !revealRe.MatchString(line) {
					continue
				}

				if reason, ok := exempt[rel]; ok {
					t.Logf("coarse-pointer exemption (%s): %s", rel, reason)

					continue
				}

				covered := false

				for _, hook := range hookRe.FindAllString(line, -1) {
					if coarseHooks[hook] {
						covered = true

						break
					}
				}

				if !covered {
					t.Errorf("coarse-pointer violation in %s: hover-revealed element has no "+
						"touch fallback — add a tc-* hook class covered by an @media (pointer: coarse) "+
						"rule in templates/custom.css, or document an exemption:\n  %s",
						rel, strings.TrimSpace(line))
				}
			}

			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", dir, err)
		}
	}
}

// coarsePointerExemptions lists files allowed to carry hover-only reveals.
// Every entry needs a justification future sessions can re-verify.
func coarsePointerExemptions() map[string]string {
	return map[string]string{
		"display/tooltip.templ": "ADR-0017: tooltips are hover/focus progressive " +
			"enhancement; touch devices deliberately do not toggle them",
		"display/hover_card.templ": "focus-within reveals the card when the trigger " +
			"contains focusable content; always-visible would duplicate content on touch " +
			"(progressive enhancement, same family as ADR-0017)",
	}
}

// loadCoarsePointerHooks collects every tc-* class kept visible by an
// @media (pointer: coarse) block in templates/custom.css.
func loadCoarsePointerHooks(t *testing.T) map[string]bool {
	t.Helper()

	cssPath := filepath.Join("..", "templates", "custom.css")

	data, err := os.ReadFile(cssPath) //nolint:gosec // reads repo CSS
	if err != nil {
		t.Fatalf("cannot read templates/custom.css: %v", err)
	}

	hooks := make(map[string]bool)

	selectorRe := regexp.MustCompile(`\.((tc|-)[a-z0-9-]+)`)
	mediaRe := regexp.MustCompile(`@media\s*\(pointer:\s*coarse\)`)

	inBlock := false
	depth := 0

	for line := range strings.SplitSeq(string(data), "\n") {
		if !inBlock && mediaRe.MatchString(line) {
			inBlock = true
		}

		if !inBlock {
			continue
		}

		depth += strings.Count(line, "{") - strings.Count(line, "}")

		for _, sel := range selectorRe.FindAllStringSubmatch(line, -1) {
			hooks[sel[1]] = true
		}

		if depth <= 0 && strings.Contains(line, "}") {
			inBlock = false
		}
	}

	if len(hooks) == 0 {
		t.Fatal("templates/custom.css has no @media (pointer: coarse) block with tc-* hooks — " +
			"the coarse-pointer fallback convention is gone")
	}

	return hooks
}
