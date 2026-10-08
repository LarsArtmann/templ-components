package display

import (
	"os"
	"strings"
	"testing"
)

// TestRTLDirectionReadsAreSubtreeScoped pins every RTL-aware JS handler
// (menu keyboard nav in shared.go, Tabs, Carousel) to resolve direction from
// the nearest [dir] ancestor, falling back to <html> — so RTL widgets inside
// LTR pages (and vice versa) map ArrowLeft/ArrowRight correctly. A bare
// documentElement read regresses mixed-direction subtrees and must fail here.
func TestRTLDirectionReadsAreSubtreeScoped(t *testing.T) {
	t.Parallel()

	sources := []string{"shared.go", "tabs.templ", "carousel.templ"}

	for _, name := range sources {
		content, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}

		src := string(content)

		if !strings.Contains(src, "closest('[dir]')") {
			t.Errorf("%s: no subtree-scoped direction read found (missing closest('[dir]'))", name)

			continue
		}

		scoped := strings.ReplaceAll(src, "||document.documentElement", "")
		scoped = strings.ReplaceAll(scoped, "|| document.documentElement", "")

		if strings.Contains(scoped, "documentElement.getAttribute('dir')") {
			t.Errorf(
				"%s: page-scoped direction read found; use (el.closest('[dir]')||document.documentElement).getAttribute('dir')",
				name,
			)
		}
	}
}
