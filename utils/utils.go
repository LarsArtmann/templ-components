package utils

import (
	"strings"
	"time"

	"github.com/a-h/templ"
)

// BaseProps provides common configurable attributes for all components
type BaseProps struct {
	ID        string
	Class     string
	Attrs     templ.Attributes
	AriaLabel string
}

// Class merges Tailwind classes, with overrides taking precedence over defaults
// Later arguments override earlier ones when they conflict
func Class(defaults string, overrides string) string {
	if overrides == "" {
		return strings.TrimSpace(defaults)
	}
	return strings.TrimSpace(overrides)
}

// MergeAttrs merges multiple attribute maps, with later maps overriding earlier ones
func MergeAttrs(maps ...templ.Attributes) templ.Attributes {
	result := make(templ.Attributes)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// CurrentYear returns the current year as a string
func CurrentYear() string {
	return time.Now().Format("2006")
}

// Ternary returns a if condition is true, otherwise b
func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// Ptr returns a pointer to the given value
func Ptr[T any](v T) *T {
	return &v
}

// Deref returns the value pointed to by p, or the zero value if p is nil
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// DerefOr returns the value pointed to by p, or fallback if p is nil
func DerefOr[T any](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}
