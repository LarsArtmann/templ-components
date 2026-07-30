package display

import (
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

// CollapsibleSectionProps configures a collapsible section using native
// <details>/<summary> elements.
type CollapsibleSectionProps struct {
	utils.BaseProps

	// Title is the section heading text.
	Title string

	// TitleTag is the heading element (h1–h6). Default: "h3".
	TitleTag string

	// Collapsed controls whether the section starts collapsed on initial
	// render. Default: false (section is expanded).
	Collapsed bool

	// StorageKey, when non-empty, persists the open/closed state to
	// localStorage under this key. A consumer-side script reads the
	// data-collapsible attribute and toggles accordingly.
	StorageKey string

	// Icon overrides the default chevron. Default: icons.ChevronDown.
	Icon icons.Name
}

// DefaultCollapsibleSectionProps returns sensible defaults.
func DefaultCollapsibleSectionProps() CollapsibleSectionProps {
	return CollapsibleSectionProps{ //nolint:exhaustruct // intentionally minimal
		TitleTag:  "h3",
		Collapsed: false,
		Icon:      icons.ChevronDown,
	}
}

// resolveCollapsibleDefaults merges user props with defaults.
func resolveCollapsibleDefaults(p CollapsibleSectionProps) CollapsibleSectionProps {
	if p.TitleTag == "" {
		p.TitleTag = "h3"
	}

	if p.Icon == "" {
		p.Icon = icons.ChevronDown
	}

	return p
}

// isValidHeadingTag reports whether tag is h1–h6.
func isValidHeadingTag(tag string) bool {
	switch tag {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return true
	default:
		return false
	}
}
