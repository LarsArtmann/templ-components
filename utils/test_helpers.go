// Test rendering utilities: Render, RenderAll, Assert* helpers for component testing.
package utils

import (
	"bytes"
	"context"
	"strings"

	"github.com/a-h/templ"
)

// TestReporter is the subset of *testing.T the render helpers use.
// *testing.T satisfies it. The interface exists so this PRODUCTION package
// never imports testing: importing testing here linked the Go test
// framework (testing.init et al.) into every consumer binary that pulls
// utils (found in dnsblockd via `go tool nm`, 2026-09-19). Test-only
// helpers stay importable from _test.go files with zero changes while the
// production import graph stays test-free.
type TestReporter interface {
	Helper()
	Fatalf(format string, args ...any)
	Errorf(format string, args ...any)
}

// Render renders a templ component to a string for testing.
func Render(t TestReporter, c templ.Component) string {
	t.Helper()

	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	return strings.TrimSpace(buf.String())
}

// RenderAll renders multiple templ components into a single concatenated string.
// Useful for integration tests that verify component composition.
func RenderAll(t TestReporter, components ...templ.Component) string {
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

// AssertContains checks that the rendered output contains a substring.
func AssertContains(t TestReporter, output, want string) {
	t.Helper()

	if !strings.Contains(output, want) {
		t.Errorf("output does not contain %q:\n%s", want, output)
	}
}

// AssertNotContains checks that the rendered output does not contain a substring.
func AssertNotContains(t TestReporter, output, notWant string) {
	t.Helper()

	if strings.Contains(output, notWant) {
		t.Errorf("output should not contain %q:\n%s", notWant, output)
	}
}

// AssertEqual checks that got equals want, reporting a test error with context if not.
func AssertEqual[T comparable](t TestReporter, context string, got, want T) {
	t.Helper()

	if got != want {
		t.Errorf("%s = %v, want %v", context, got, want)
	}
}

// AssertContainsAll checks that the rendered output contains every substring
// in wants. Reports a single test error per missing substring.
func AssertContainsAll(t TestReporter, output string, wants ...string) {
	t.Helper()

	for _, want := range wants {
		if !strings.Contains(output, want) {
			t.Errorf("output does not contain %q:\n%s", want, output)
		}
	}
}
