package icons

import "github.com/a-h/templ"

// IconAttrs returns attributes for an accessible icon.
// When ariaLabel is non-empty, returns aria-label. Otherwise returns aria-hidden="true".
func IconAttrs(ariaLabel string) templ.Attributes {
	if ariaLabel != "" {
		return templ.Attributes{
			"aria-label": ariaLabel,
		}
	}
	return templ.Attributes{
		"aria-hidden": "true",
	}
}
