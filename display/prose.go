// Prose component: a longform-content container for rendered markdown and
// article HTML. One wrapper class (.tc-prose, templates/custom.css) styles
// headings, lists, blockquotes, code, tables, and images via descendant
// selectors — the consumer renders raw content HTML inside and gets the
// library's typographic contract (theme-token colors, dark-mode overrides,
// logical/RTL-safe spacing) without per-element classes.
package display

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// ProseSize controls the container's font scale. Utilities beat the
// .tc-prose descendant rules, so the size classes override the base font.
type ProseSize string

const (
	// ProseSizeSM renders compact longform text (0.875rem).
	ProseSizeSM ProseSize = "sm"
	// ProseSizeMD is the default reading scale (1rem, the .tc-prose base).
	ProseSizeMD ProseSize = "md"
	// ProseSizeLG renders large text (1.125rem) for lead articles.
	ProseSizeLG ProseSize = "lg"
)

// ProseSizeDefault is the zero-value-compatible default size.
const ProseSizeDefault = ProseSizeMD

// ProseSizeIsValid reports whether v is one of the defined ProseSize constants.
func ProseSizeIsValid(v ProseSize) bool {
	switch v {
	case ProseSizeSM, ProseSizeMD, ProseSizeLG:
		return true
	default:
		return false
	}
}

// proseSizeLookup maps each size to the classes layered over .tc-prose.
// MD is intentionally empty: the .tc-prose base IS the medium scale.
var proseSizeLookup = map[ProseSize]string{ //nolint:gochecknoglobals // typed-key style lookup (repo convention)
	ProseSizeSM: "text-sm leading-relaxed",
	ProseSizeMD: "",
	ProseSizeLG: "text-lg leading-relaxed",
}

// proseSizeClass resolves the size classes with the zero value = default.
func proseSizeClass(size ProseSize) string {
	if size == "" {
		return proseSizeLookup[ProseSizeMD]
	}

	return utils.Lookup(proseSizeLookup, size, "")
}

// ProseProps configures a longform-content container.
type ProseProps struct {
	utils.BaseProps

	// Children is the longform content (typically rendered markdown HTML).
	Children templ.Component
	// Size sets the font scale. Zero value falls back to ProseSizeMD.
	Size ProseSize
}

// DefaultProseProps returns sensible defaults.
func DefaultProseProps() ProseProps {
	return ProseProps{Size: ProseSizeMD}
}
