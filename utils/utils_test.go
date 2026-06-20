// Package utils provides testing utilities for the templ-components library.
package utils

import (
	"fmt"
	"slices"
	"strings"
	"sync"
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

// Verify that *BaseProps satisfies the ComponentProps interface.
func TestBasePropsImplementsComponentProps(t *testing.T) {
	t.Parallel()

	var _ ComponentProps = &BaseProps{}

	// Verify GetBaseProps returns the correct values.
	bp := &BaseProps{ID: "test-id", Class: "test-class", AriaLabel: "test-label"}
	got := bp.GetBaseProps()
	if got.ID != "test-id" {
		t.Errorf("GetBaseProps().ID = %q, want %q", got.ID, "test-id")
	}
	if got.Class != "test-class" {
		t.Errorf("GetBaseProps().Class = %q, want %q", got.Class, "test-class")
	}
	if got.AriaLabel != "test-label" {
		t.Errorf("GetBaseProps().AriaLabel = %q, want %q", got.AriaLabel, "test-label")
	}

	// Verify SetBaseProps updates the struct.
	bp.SetBaseProps(BaseProps{ID: "new-id"})
	if bp.ID != "new-id" {
		t.Errorf("after SetBaseProps, ID = %q, want %q", bp.ID, "new-id")
	}
}

// ComponentProps interface must be satisfied by any struct embedding BaseProps.
// This test verifies the method-promotion mechanism.
func TestComponentPropsInterfacePromoted(t *testing.T) {
	t.Parallel()

	type testComponentProps struct {
		BaseProps
		Text string
	}

	var cp ComponentProps = &testComponentProps{Text: "hello"}

	// GetBaseProps returns the embedded BaseProps.
	bp := cp.GetBaseProps()
	if bp.Class != "" {
		t.Errorf("expected empty class, got %q", bp.Class)
	}

	// SetBaseProps updates the embedded BaseProps.
	cp.SetBaseProps(BaseProps{Class: "updated"})
	updated := cp.GetBaseProps()
	if updated.Class != "updated" {
		t.Errorf("expected class %q, got %q", "updated", updated.Class)
	}
}

func TestClassConcurrentAccess(t *testing.T) {
	t.Parallel()
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			result := Class("px-4 py-2", "px-6", fmt.Sprintf("bg-%d", i))
			if !strings.Contains(result, "px-6") {
				t.Errorf("Class merge failed: %s", result)
			}
		})
	}
	wg.Wait()
}

func TestValidateID(t *testing.T) {
	t.Parallel()
	t.Run("empty ID returns error", func(t *testing.T) {
		t.Parallel()
		err := ValidateID("modal", "")
		if err == nil {
			t.Error("expected error for empty ID")
		}
	})
	t.Run("non-empty ID returns nil", func(t *testing.T) {
		t.Parallel()
		err := ValidateID("modal", "my-modal")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestEnsureID(t *testing.T) {
	t.Parallel()
	t.Run("returns existing ID unchanged", func(t *testing.T) {
		t.Parallel()
		got := EnsureID("modal", "my-modal")
		if got != "my-modal" {
			t.Errorf("EnsureID() = %q, want %q", got, "my-modal")
		}
	})
	t.Run("generates prefixed ID when empty", func(t *testing.T) {
		t.Parallel()
		got := EnsureID("drawer", "")
		if !strings.HasPrefix(got, "tc-drawer-") {
			t.Errorf("EnsureID() = %q, want prefix %q", got, "tc-drawer-")
		}
		if len(got) <= len("tc-drawer-") {
			t.Errorf("EnsureID() = %q, expected non-empty hex suffix", got)
		}
	})
	t.Run("generates unique IDs", func(t *testing.T) {
		t.Parallel()
		a := EnsureID("modal", "")
		b := EnsureID("modal", "")
		if a == b {
			t.Errorf("EnsureID() generated duplicate IDs: %q == %q", a, b)
		}
	})
}
