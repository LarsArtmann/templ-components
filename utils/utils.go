package utils

import (
	"maps"
	"sync"
	"time"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/a-h/templ"
)

//nolint:gochecknoglobals // Required because tailwind-merge-go is not thread-safe
var mergeMutex sync.Mutex

// BaseProps provides common configurable attributes for all components
type BaseProps struct {
	ID        string
	Class     string
	Attrs     templ.Attributes
	AriaLabel string
	Nonce     string
}

// Class merges Tailwind classes intelligently using tailwind-merge-go.
// Conflicting classes are resolved with later arguments overriding earlier ones.
func Class(classes ...string) string {
	mergeMutex.Lock()
	defer mergeMutex.Unlock()

	return twmerge.Merge(classes...)
}

// MergeAttrs merges multiple attribute maps, with later maps overriding earlier ones
func MergeAttrs(m ...templ.Attributes) templ.Attributes {
	result := make(templ.Attributes)
	for _, attrs := range m {
		maps.Copy(result, attrs)
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
//
//go:fix inline
func Ptr[T any](v T) *T {
	return new(v)
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

// BoolString returns "true" or "false" for a boolean value
func BoolString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// MapEnum looks up a string key in a map and returns the corresponding enum value,
// or the fallback if the key is not found.
func MapEnum[T ~string](m map[string]T, fallback T, key string) T {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}
