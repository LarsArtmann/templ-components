package utils

import (
	"sync"
	"time"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/a-h/templ"
)

// BaseProps provides common configurable attributes for all components
type BaseProps struct {
	ID        string
	Class     string
	Attrs     templ.Attributes
	AriaLabel string
	Nonce     string
}

// ComponentProps is the interface implemented by all component props structs
// that embed BaseProps. It enables generic component composition, validation,
// and wrapper functions across the entire component library.
type ComponentProps interface {
	GetBaseProps() BaseProps
	SetBaseProps(BaseProps)
}

// GetBaseProps returns the embedded BaseProps. This method is promoted to all
// structs that embed BaseProps, satisfying the ComponentProps interface.
func (p *BaseProps) GetBaseProps() BaseProps {
	return *p
}

// SetBaseProps updates the embedded BaseProps. This method is promoted to all
// structs that embed BaseProps, satisfying the ComponentProps interface.
func (p *BaseProps) SetBaseProps(bp BaseProps) {
	*p = bp
}

//nolint:gochecknoglobals // Package-level merge cache + mutex for thread safety
var classMu sync.Mutex

// Class merges Tailwind classes intelligently using tailwind-merge-go.
// Conflicting classes are resolved with later arguments overriding earlier ones.
// Thread-safe via sync.Mutex to protect the shared LRU cache.
func Class(classes ...string) string {
	classMu.Lock()
	defer classMu.Unlock()
	return twmerge.Merge(classes...)
}

// CurrentYear returns the current year as a string
func CurrentYear() string {
	return time.Now().Format("2006")
}

// Ternary returns a if condition is true, otherwise b.
// Note: both a and b are eagerly evaluated (this is a function, not a macro),
// so callers passing function calls or side-effecting expressions should use
// inline if/else in templ instead: { if cond }...{ else }...{ end }
func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// Lookup returns the map value for key, or fallback if not found.
// Replaces the repetitive if-ok-return pattern used across all enum lookups.
func Lookup[K comparable, V any](m map[K]V, key K, fallback V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}

// DismissScript returns shared JavaScript for dismissing elements via [data-dismiss]
// click delegation. Used by Alert, Toast, and ErrorAlert components.
// Handles both role="alert" (Alert) and role="status" (Toast) containers.
func DismissScript() string {
	return `if(!window.tcDismissAttached){window.tcDismissAttached=true;document.addEventListener('click',function(e){var btn=e.target.closest('[data-dismiss]');if(btn){var el=btn.closest('[role="alert"],[role="status"]');if(el)el.remove();}});}`
}
