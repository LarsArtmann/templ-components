package recipes

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// AuthLayoutProps configures a recipes.AuthLayout screen — the canonical
// split-screen authentication layout: a centered auth card on one side and a
// branding/feature panel on the other. On mobile, the branding panel collapses
// to show only the card.
//
// The card content is consumer-supplied (typically a LoginCard, a signup form,
// or a password-reset form). The branding panel shows an optional title,
// description, feature bullets, and footer.
type AuthLayoutProps struct {
	utils.BaseProps

	// Card is the main auth content (typically a LoginCard or forms.Form
	// wrapped in a Card). Required.
	Card templ.Component
	// PanelTitle renders as a large heading in the branding panel. Optional.
	PanelTitle string
	// PanelText renders as body text under the PanelTitle. Optional.
	PanelText string
	// PanelFeatures is a list of bullet points rendered with checkmark icons
	// in the branding panel. Optional.
	PanelFeatures []string
	// PanelFooter renders at the bottom of the branding panel (e.g. a
	// testimonial, copyright, or security badge). Optional.
	PanelFooter templ.Component
	// Reverse flips the layout so the card appears on the right and the
	// branding panel on the left. Default: card left, panel right.
	// (On mobile the panel is always hidden.)
	Reverse bool
}

// DefaultAuthLayoutProps returns sensible defaults (empty).
//
//nolint:exhaustruct_v5 // constructor intentionally sets only non-zero defaults
func DefaultAuthLayoutProps() AuthLayoutProps {
	return AuthLayoutProps{}
}
