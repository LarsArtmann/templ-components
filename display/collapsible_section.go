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

	// TitleTag is the heading element (h1–h6). Default: HeadingTagH3.
	TitleTag HeadingTagType

	// Collapsed controls whether the section starts collapsed on initial
	// render. Default: false (section is expanded).
	Collapsed bool

	// StorageKey, when non-empty, persists the open/closed state to
	// localStorage under this key. A consumer-side script reads the
	// data-collapsible attribute and toggles accordingly.
	StorageKey string

	// PersistState, when true, ships the persistence script with the
	// component: a CSP-safe, nonce'd singleton that restores the saved
	// open/closed state for every details[data-collapsible] on the page
	// and writes back on toggle (capture phase — toggle does not bubble).
	// Requires StorageKey AND BaseProps.Nonce; without a nonce the script
	// is skipped so strict-CSP pages never receive un-nonced inline JS.
	// Consumers who bring their own persistence script leave this false
	// (the default) and keep reading data-collapsible themselves.
	PersistState bool

	// Icon overrides the default chevron. Default: icons.ChevronDown.
	Icon icons.Name
}

// DefaultCollapsibleSectionProps returns sensible defaults.
func DefaultCollapsibleSectionProps() CollapsibleSectionProps {
	return CollapsibleSectionProps{ //nolint:exhaustruct_v5 // intentionally minimal
		TitleTag:  HeadingTagH3,
		Collapsed: false,
		Icon:      icons.ChevronDown,
	}
}

// resolveCollapsibleDefaults merges user props with defaults.
func resolveCollapsibleDefaults(p CollapsibleSectionProps) CollapsibleSectionProps {
	if p.TitleTag == "" {
		p.TitleTag = HeadingTagH3
	}

	if p.Icon == "" {
		p.Icon = icons.ChevronDown
	}

	return p
}
