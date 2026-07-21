package recipes

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// SettingsLayoutProps configures a recipes.SettingsLayout screen — the
// canonical "section nav on the left, form content on the right" pattern
// used by almost every settings page in every admin panel.
//
// The aside is a vertical list of section links; the main is a stack of
// Card-wrapped forms (one per section). Both are consumer-supplied
// templ.Component slots.
type SettingsLayoutProps struct {
	utils.BaseProps

	// Title and Subtitle render at the top of the page.
	Title    string
	Subtitle string
	// Aside is the section navigation slot (typically a Card-wrapped list
	// of links to #section-id anchors). Required.
	Aside templ.Component
	// Sections is the list of form sections rendered in the main column.
	// Each section is rendered as a Card with Title + a form body. The
	// Section.ID becomes the card's id (for anchor navigation from Aside).
	Sections []SettingsSection
}

// SettingsSection is one titled section of a settings page.
type SettingsSection struct {
	// ID is used as the anchor target (id attribute) so the Aside links
	// can jump to it. Optional but recommended.
	ID string
	// Title renders as the Card header.
	Title string
	// Subtitle renders as muted text under the Title. Optional.
	Subtitle string
	// Body is the form content of the section. Typically a forms.Form
	// or a Card.Body slot. Required.
	Body templ.Component
}

// DefaultSettingsLayoutProps returns sensible defaults (empty).
func DefaultSettingsLayoutProps() SettingsLayoutProps {
	return SettingsLayoutProps{} //nolint:exhaustruct // intentionally minimal
}
