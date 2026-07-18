// ARIA attribute helpers for form field accessibility.
package forms

import "github.com/a-h/templ"

const (
	ariaDescribedBy = "aria-describedby"
	ariaInvalid     = "aria-invalid"
	ariaTrue        = "true"
	errorIDSuffix   = "-error"
)

// ErrorAttrs returns templ attributes for aria-invalid and aria-describedby
// when an error message or help text is present. Returns nil only when both
// errorMessage and helpTextID are empty. When helpTextID is non-empty it is always
// emitted via aria-describedby (independent of id), so id-less fields with help
// text remain accessible to assistive technology.
func ErrorAttrs(id, errorMessage, helpTextID string) templ.Attributes {
	if errorMessage == "" {
		if helpTextID != "" {
			return templ.Attributes{ariaDescribedBy: helpTextID}
		}

		return nil
	}

	attrs := templ.Attributes{
		ariaInvalid: ariaTrue,
	}

	switch {
	case id != "" && helpTextID != "":
		attrs[ariaDescribedBy] = id + errorIDSuffix + " " + helpTextID
	case id != "":
		attrs[ariaDescribedBy] = id + errorIDSuffix
	case helpTextID != "":
		attrs[ariaDescribedBy] = helpTextID
	}

	return attrs
}
