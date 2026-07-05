// Test rendering utilities: Render, RenderAll, Assert* helpers for component testing.
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

// RenderAll renders multiple templ components into a single concatenated string.
// Useful for integration tests that verify component composition.
func RenderAll(t *testing.T, components ...templ.Component) string {
	t.Helper()
	var sb strings.Builder
	for _, c := range components {
		var buf bytes.Buffer
		if err := c.Render(context.Background(), &buf); err != nil {
			t.Fatalf("failed to render component: %v", err)
		}
		sb.WriteString(buf.String())
	}
	return strings.TrimSpace(sb.String())
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

// AssertContainsAll checks that the rendered output contains every substring
// in wants. Reports a single test error per missing substring.
func AssertContainsAll(t *testing.T, output string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(output, want) {
			t.Errorf("output does not contain %q:\n%s", want, output)
		}
	}
}
