// ID sanitization for form field elements.
package forms

import "regexp"

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9-]`)

// SanitizeID returns a safe HTML ID from arbitrary text.
func SanitizeID(s string) string {
	return nonAlphaNum.ReplaceAllString(s, "-")
}

// HelpTextID returns the HTML ID for a help text element, or "" if id is empty.
func HelpTextID(id string) string {
	if id == "" {
		return ""
	}
	return id + "-help"
}
