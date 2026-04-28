// Package forms provides tests for form components like Input, Select, Textarea, and helpers.
package forms

import "testing"

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
