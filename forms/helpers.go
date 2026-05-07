// Package forms provides utility functions for form components.
package forms

import (
	"regexp"

	"github.com/a-h/templ"
)

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9-]`)

// SanitizeID returns a safe HTML ID from arbitrary text
func SanitizeID(s string) string {
	return nonAlphaNum.ReplaceAllString(s, "-")
}

// ErrorAttrs returns templ attributes for aria-invalid and aria-describedby
// when an error message is present. Returns nil if errMsg is empty.
func ErrorAttrs(id, errMsg string) templ.Attributes {
	if errMsg == "" {
		return nil
	}
	attrs := templ.Attributes{
		"aria-invalid": "true",
	}
	if id != "" {
		attrs["aria-describedby"] = id + "-error"
	}
	return attrs
}
