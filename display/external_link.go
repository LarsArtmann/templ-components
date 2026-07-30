package display

import "github.com/larsartmann/templ-components/utils"

// ExternalLinkProps configures a safe external link.
type ExternalLinkProps struct {
	utils.BaseProps

	// Href is the target URL. Passed as a plain string (NOT templ.SafeURL)
	// so that templ's built-in URL sanitizer runs — it blocks javascript:,
	// data:, vbscript: and other dangerous schemes by rewriting them.
	Href string

	// Text is the visible link text. When empty, children are rendered instead.
	Text string

	// ShowIcon controls whether the external-arrow icon (↗) is rendered.
	// Default: true.
	ShowIcon bool
}

// DefaultExternalLinkProps returns sensible defaults for an external link.
func DefaultExternalLinkProps() ExternalLinkProps {
	return ExternalLinkProps{ //nolint:exhaustruct // intentionally minimal defaults
		ShowIcon: true,
	}
}
