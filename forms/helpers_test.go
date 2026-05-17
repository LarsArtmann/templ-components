// Package forms provides tests for form components like Input, Select, Textarea, and helpers.
package forms

import (
	"testing"

	"github.com/a-h/templ"
)

func TestSanitizeID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple text", input: "This field is required", want: "This-field-is-required"},
		{name: "already clean", input: "email-error", want: "email-error"},
		{name: "special chars", input: "foo@bar.baz!", want: "foo-bar-baz-"},
		{name: "empty", input: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := SanitizeID(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestErrorAttrs(t *testing.T) {
	t.Parallel()
	t.Run("returns nil when no error", func(t *testing.T) {
		t.Parallel()
		got := ErrorAttrs("email", "", "")
		if got != nil {
			t.Errorf("ErrorAttrs(%q, %q, %q) = %v, want nil", "email", "", "", got)
		}
	})
	t.Run("returns aria attrs with id", func(t *testing.T) {
		t.Parallel()
		got := ErrorAttrs("email", "required", "")
		want := templ.Attributes{
			"aria-invalid":     "true",
			"aria-describedby": "email-error",
		}
		if len(got) != len(want) {
			t.Fatalf("ErrorAttrs returned %d attrs, want %d", len(got), len(want))
		}
		for k, v := range want {
			if got[k] != v {
				t.Errorf("attrs[%q] = %v, want %v", k, got[k], v)
			}
		}
	})
	t.Run("returns aria attrs without id", func(t *testing.T) {
		t.Parallel()
		got := ErrorAttrs("", "required", "")
		if _, ok := got["aria-describedby"]; ok {
			t.Error("should not have aria-describedby when id is empty")
		}
		if got["aria-invalid"] != "true" {
			t.Error("should have aria-invalid=true")
		}
	})
	t.Run("includes helpTextID in aria-describedby", func(t *testing.T) {
		t.Parallel()
		got := ErrorAttrs("email", "required", "email-help")
		want := templ.Attributes{
			"aria-invalid":     "true",
			"aria-describedby": "email-error email-help",
		}
		if len(got) != len(want) {
			t.Fatalf("ErrorAttrs returned %d attrs, want %d", len(got), len(want))
		}
		for k, v := range want {
			if got[k] != v {
				t.Errorf("attrs[%q] = %v, want %v", k, got[k], v)
			}
		}
	})
}
