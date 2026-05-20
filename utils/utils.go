package utils

import (
	"sync"
	"time"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/a-h/templ"
)

//nolint:gochecknoglobals // Required to protect twmerge.Merge — the library's internal cache is not thread-safe
var classMu sync.Mutex

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
	classMu.Lock()
	defer classMu.Unlock()
	return twmerge.Merge(classes...)
}

// MapEnum returns the current year as a string
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

// MapEnum looks up a string key in a map and returns the corresponding enum value,
// or the fallback if the key is not found.
func MapEnum[T ~string](m map[string]T, fallback T, key string) T {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}
