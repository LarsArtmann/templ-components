package utils

import "time"

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
