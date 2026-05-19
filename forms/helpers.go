// Package forms provides utility functions for form components.
package forms

import (
	"regexp"

	"github.com/a-h/templ"
)

const ariaDescribedBy = "aria-describedby"
const ariaInvalid    = "aria-invalid"

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9-]`)

// SanitizeID returns a safe HTML ID from arbitrary text
func SanitizeID(s string) string {
	return nonAlphaNum.ReplaceAllString(s, "-")
}

// ErrorAttrs returns templ attributes for aria-invalid and aria-describedby
// when an error message is present. Returns nil if errMsg is empty and no helpTextID.
// When helpTextID is non-empty, it is appended to aria-describedby alongside the error ID.
func ErrorAttrs(id, errMsg, helpTextID string) templ.Attributes {
	if errMsg == "" {
		if helpTextID != "" && id != "" {
			return templ.Attributes{
				ariaDescribedBy: helpTextID,
			}
		}
		return nil
	}
	attrs := templ.Attributes{
		ariaInvalid: "true",
	}
	if id != "" {
		describedBy := id + "-error"
		if helpTextID != "" {
			describedBy += " " + helpTextID
		}
		attrs[ariaDescribedBy] = describedBy
	}
	return attrs
}

// HelpTextID returns the HTML ID for a help text element, or "" if id is empty.
func HelpTextID(id string) string {
	if id == "" {
		return ""
	}
	return id + "-help"
}
