// Package utils provides shared utilities for templ components including Tailwind class
// merging, attribute helpers, and test rendering utilities.
package utils

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// Render renders a templ component to a string for testing
func Render(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("failed to render component: %v", err)
	}
	return strings.TrimSpace(buf.String())
}

// AssertContains checks that the rendered output contains a substring
func AssertContains(t *testing.T, output, want string) {
	t.Helper()
	if !strings.Contains(output, want) {
		t.Errorf("output does not contain %q:\n%s", want, output)
	}
}

// AssertNotContains checks that the rendered output does not contain a substring
func AssertNotContains(t *testing.T, output, notWant string) {
	t.Helper()
	if strings.Contains(output, notWant) {
		t.Errorf("output should not contain %q:\n%s", notWant, output)
	}
}

// AssertEqual checks that got equals want, reporting a test error with context if not
func AssertEqual[T comparable](t *testing.T, context string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", context, got, want)
	}
}

// AssertContainsClass checks that the rendered output contains a CSS class
// in a class attribute. Handles classes embedded in larger class strings.
func AssertContainsClass(t *testing.T, output, class string) {
	t.Helper()
	if !strings.Contains(output, class) {
		t.Errorf("output does not contain class %q:\n%s", class, output)
	}
}
