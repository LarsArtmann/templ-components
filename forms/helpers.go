package forms

import (
	"fmt"
	"regexp"
)

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9-]`)

// SanitizeID returns a safe HTML ID from arbitrary text
func SanitizeID(s string) string {
	return nonAlphaNum.ReplaceAllString(s, "-")
}

// FormatFloat formats a float64 for use in form inputs
func FormatFloat(v float64, decimals int) string {
	if v == 0 {
		return ""
	}
	format := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(format, v)
}

// IsSelected returns "selected" if value matches option (for templ boolean attributes)
func IsSelected(value, option string) bool {
	return value == option
}

// IfNotNil returns the string representation of a field value if the object is not nil
func IfNotNil[T any](obj *T, getter func(T) string) string {
	if obj == nil {
		return ""
	}
	return getter(*obj)
}

// IfNotNilString returns the string value if the pointer is not nil
func IfNotNilString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BoolToString returns "true" or "false" for a boolean
func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
