// Package forms provides utility functions for form components.
package forms

import "regexp"

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9-]`)

// SanitizeID returns a safe HTML ID from arbitrary text
func SanitizeID(s string) string {
	return nonAlphaNum.ReplaceAllString(s, "-")
}
