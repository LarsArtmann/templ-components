// Package utils provides testing utilities for the templ-components library.
package utils

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
)

func TestCurrentYear(t *testing.T) {
	t.Parallel()
	got := CurrentYear()
	want := time.Now().Format("2006")
	if got != want {
		t.Errorf("CurrentYear() = %q, want %q", got, want)
	}
}

func TestBoolString(t *testing.T) {
	t.Parallel()
	t.Run("true", func(t *testing.T) {
		t.Parallel()
		got := BoolString(true)
		AssertEqual(t, "BoolString(true)", got, "true")
	})
	t.Run("false", func(t *testing.T) {
		t.Parallel()
		got := BoolString(false)
		AssertEqual(t, "BoolString(false)", got, "false")
	})
}

func TestPtr(t *testing.T) {
	t.Parallel()
	v := "hello"
	p := new(v)
	AssertEqual(t, "*new(v)", *p, v)
}

func TestTernary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		condition bool
		a, b      string
		want      string
	}{
		{name: "true returns a", condition: true, a: "yes", b: "no", want: "yes"},
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

func TestDeref(t *testing.T) {
	t.Parallel()
	t.Run("non-nil", func(t *testing.T) {
		t.Parallel()
		v := 42
		got := Deref(&v)
		if got != 42 {
			t.Errorf("Deref() = %d, want 42", got)
		}
	})
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		var p *int
		got := Deref(p)
		if got != 0 {
			t.Errorf("Deref() = %d, want 0", got)
		}
	})
}

func TestDerefOr(t *testing.T) {
	t.Parallel()
	t.Run("non-nil", func(t *testing.T) {
		t.Parallel()
		v := 42
		got := DerefOr(&v, 99)
		if got != 42 {
			t.Errorf("DerefOr() = %d, want 42", got)
		}
	})
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		var p *int
		got := DerefOr(p, 99)
		if got != 99 {
			t.Errorf("DerefOr() = %d, want 99", got)
		}
	})
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
	return slices.Contains(splitSpace(s), word)
}

func splitSpace(s string) []string {
	return strings.Fields(s)
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

func TestMergeAttrs(t *testing.T) {
	t.Parallel()
	t.Run("merge two maps", func(t *testing.T) {
		t.Parallel()
		a := templ.Attributes{"data-foo": "1"}
		b := templ.Attributes{"data-bar": "2"}
		got := MergeAttrs(a, b)
		if got["data-foo"] != "1" || got["data-bar"] != "2" {
			t.Errorf("MergeAttrs() = %v", got)
		}
	})
	t.Run("later overrides earlier", func(t *testing.T) {
		t.Parallel()
		a := templ.Attributes{"data-foo": "1"}
		b := templ.Attributes{"data-foo": "2"}
		got := MergeAttrs(a, b)
		if got["data-foo"] != "2" {
			t.Errorf("MergeAttrs() data-foo = %q, want 2", got["data-foo"])
		}
	})
	t.Run("empty maps", func(t *testing.T) {
		t.Parallel()
		got := MergeAttrs()
		if len(got) != 0 {
			t.Errorf("MergeAttrs() len = %d, want 0", len(got))
		}
	})
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
