// Package utils provides testing utilities for the templ-components library.
package utils

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
)

const (
	textYes = "yes"
)

func TestCurrentYear(t *testing.T) {
	t.Parallel()
	got := CurrentYear()
	want := time.Now().Format("2006")
	if got != want {
		t.Errorf("CurrentYear() = %q, want %q", got, want)
	}
}

func TestTernary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		condition bool
		a, b      string
		want      string
	}{
		{name: "true returns a", condition: true, a: textYes, b: "no", want: textYes},
		{name: "false returns b", condition: false, a: "yes", b: "no", want: "no"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Ternary(tt.condition, tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Ternary() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		classes        []string
		wantContain    []string
		wantNotContain []string
	}{
		{
			name:           "single string",
			classes:        []string{"a b c"},
			wantContain:    []string{"a", "b", "c"},
			wantNotContain: nil,
		},
		{
			name:           "tailwind merge",
			classes:        []string{"bg-red-500 hover:bg-blue-500", "bg-green-500"},
			wantContain:    []string{"bg-green-500", "hover:bg-blue-500"},
			wantNotContain: []string{"bg-red-500"},
		},
		{
			name:           "empty ignored",
			classes:        []string{"a b", ""},
			wantContain:    []string{"a", "b"},
			wantNotContain: nil,
		},
		{
			name:           "multiple overrides",
			classes:        []string{"p-4", "p-6", "px-8"},
			wantContain:    []string{"p-6", "px-8"},
			wantNotContain: []string{"p-4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Class(tt.classes...)
			for _, want := range tt.wantContain {
				if !containsWord(got, want) {
					t.Errorf("Class(%v) = %q, want to contain %q", tt.classes, got, want)
				}
			}
			for _, notWant := range tt.wantNotContain {
				if containsWord(got, notWant) {
					t.Errorf("Class(%v) = %q, should not contain %q", tt.classes, got, notWant)
				}
			}
		})
	}
}

func containsWord(s, word string) bool {
	return slices.Contains(strings.Fields(s), word)
}

type testEnum string

const (
	testEnumA testEnum = "a"
	testEnumB testEnum = "b"
	testEnumC testEnum = "c"
)

func TestMapEnum(t *testing.T) {
	t.Parallel()

	lookup := map[string]testEnum{
		"alpha": testEnumA,
		"beta":  testEnumB,
	}

	t.Run("found key returns mapped value", func(t *testing.T) {
		t.Parallel()
		got := MapEnum(lookup, testEnumC, "alpha")
		if got != testEnumA {
			t.Errorf("MapEnum() = %q, want %q", got, testEnumA)
		}
	})

	for _, key := range []string{"unknown", ""} {
		t.Run(key+" key returns fallback", func(t *testing.T) {
			t.Parallel()
			got := MapEnum(lookup, testEnumC, key)
			if got != testEnumC {
				t.Errorf("MapEnum() = %q, want %q (fallback)", got, testEnumC)
			}
		})
	}
}

func TestRender(t *testing.T) {
	t.Parallel()
	t.Run("renders component to string", func(t *testing.T) {
		t.Parallel()
		c := templ.Raw("<div>hello</div>")
		got := Render(t, c)
		AssertContains(t, got, "<div>hello</div>")
	})
}

func TestAssertContains(t *testing.T) {
	t.Parallel()
	AssertContains(t, "hello world", "hello")
}

func TestAssertNotContains(t *testing.T) {
	t.Parallel()
	AssertNotContains(t, "hello world", "xyz")
}

func TestAssertEqual(t *testing.T) {
	t.Parallel()
	AssertEqual(t, "matching values", 42, 42)
}
