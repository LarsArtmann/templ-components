package utils

import (
	"sync"
	"time"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/a-h/templ"
)

//nolint:gochecknoglobals // Required to protect tailwind-merge-go's lazy init on first call to twmerge.Merge
var classMu sync.Mutex

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

// Class merges Tailwind classes intelligently using tailwind-merge-go.
// Conflicting classes are resolved with later arguments overriding earlier ones.
func Class(classes ...string) string {
	classMu.Lock()
	defer classMu.Unlock()
	return twmerge.Merge(classes...)
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

// MapEnum looks up a string key in a map and returns the corresponding enum value,
// or the fallback if the key is not found.
func MapEnum[T ~string](m map[string]T, fallback T, key string) T {
	if v, ok := m[key]; ok {
		return v
	}
	return fallback
}

// DismissScript returns shared JavaScript for dismissing elements via [data-dismiss]
// click delegation. Used by Alert, Toast, and ErrorAlert components.
func DismissScript() string {
	return `if(!window.tcDismissAttached){window.tcDismissAttached=true;document.addEventListener('click',function(e){var btn=e.target.closest('[data-dismiss]');if(btn){var el=btn.closest('[role="alert"]');if(el)el.remove();}});}`
}
