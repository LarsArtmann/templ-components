package navigation

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// navMobileMenuID returns the id shared by a Nav's mobile menu and its toggle.
//
// With a caller-supplied BaseProps.ID the menu id is DERIVED from it
// ("<id>-mobile-menu") rather than reusing it: reusing it puts the same id on
// the <nav> element and the menu div, which is invalid HTML and leaves
// getElementById()/aria-controls resolving the wrong element.
//
// Without one the id is derived from the nav's own content, so rendering the
// same page twice produces byte-identical HTML. A random per-render id (what
// utils.EnsureID mints) makes every response body unique, so ETag
// revalidation can never answer 304 and each navigation or htmx poll
// re-downloads a full document. Two Navs with identical links, current path,
// and label in one page share the derived id — give each Nav its own
// BaseProps.ID when that distinction matters.
func navMobileMenuID(p NavProps) string {
	if p.ID != "" {
		return p.ID + "-mobile-menu"
	}

	var seed strings.Builder
	seed.WriteString(p.CurrentPath)
	seed.WriteString("\x00")
	seed.WriteString(p.AriaLabel)

	if p.ContainerAware {
		seed.WriteString("\x01container")
	}

	for _, link := range p.Links {
		seed.WriteString("\x00")
		seed.WriteString(link.Href)
		seed.WriteString("\x00")
		seed.WriteString(link.Text)
	}

	sum := sha256.Sum256([]byte(seed.String()))

	return "tc-mobile-menu-" + hex.EncodeToString(sum[:8])
}
