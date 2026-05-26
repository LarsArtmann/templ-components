// Package forms provides utility functions for form components.
package forms

import (
	"regexp"

	"github.com/a-h/templ"
)

const (
	ariaDescribedBy = "aria-describedby"
	ariaInvalid     = "aria-invalid"
	ariaTrue        = "true"
)

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
		ariaInvalid: ariaTrue,
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

// baseInputClass returns the shared Tailwind classes for text inputs, selects, and textareas
func baseInputClass(hasError bool) string {
	base := "block w-full rounded-md border-0 py-1.5 text-gray-900 dark:text-white shadow-xs ring-1 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-inset sm:text-sm sm:leading-6 transition-colors dark:bg-gray-800 dark:ring-gray-700 dark:placeholder:text-gray-500"
	if hasError {
		return base + " ring-red-300 focus:ring-red-500 dark:ring-red-700 dark:focus:ring-red-500"
	}
	return base + " ring-gray-300 focus:ring-blue-600 dark:focus:ring-blue-500"
}
