// Package forms provides tests for form components like Input, Select, Textarea, and helpers.
package forms

import (
	"testing"

	"github.com/a-h/templ"
)

const emailErrorSuffix = "email-error"

func TestSanitizeID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple text", input: "This field is required", want: "This-field-is-required"},
		{name: "already clean", input: emailErrorSuffix, want: emailErrorSuffix},
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

func TestErrorAttrs_NoErrorNoHelpText(t *testing.T) {
	t.Parallel()
	got := ErrorAttrs("email", "", "")
	if got != nil {
		t.Errorf("ErrorAttrs(%q, %q, %q) = %v, want nil", "email", "", "", got)
	}
}

func TestErrorAttrs_AriaAttrsWithID(t *testing.T) {
	t.Parallel()
	got := ErrorAttrs("email", "required", "")
	want := templ.Attributes{
		ariaInvalid:     "true",
		ariaDescribedBy: "email-error",
	}
	if len(got) != len(want) {
		t.Fatalf("ErrorAttrs returned %d attrs, want %d", len(got), len(want))
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("attrs[%q] = %v, want %v", k, got[k], v)
		}
	}
}

func TestErrorAttrs_AriaAttrsWithoutID(t *testing.T) {
	t.Parallel()
	got := ErrorAttrs("", "required", "")
	if _, ok := got[ariaDescribedBy]; ok {
		t.Error("should not have aria-describedby when id is empty")
	}
	if got[ariaInvalid] != ariaTrue {
		t.Error("should have aria-invalid=true")
	}
}

func TestErrorAttrs_HelpTextIDIncluded(t *testing.T) {
	t.Parallel()
	got := ErrorAttrs("email", "required", "email-help")
	want := templ.Attributes{
		ariaInvalid:     "true",
		ariaDescribedBy: "email-error email-help",
	}
	if len(got) != len(want) {
		t.Fatalf("ErrorAttrs returned %d attrs, want %d", len(got), len(want))
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("attrs[%q] = %v, want %v", k, got[k], v)
		}
	}
}

func TestErrorAttrs_HelpTextOnlyNoError(t *testing.T) {
	t.Parallel()
	got := ErrorAttrs("email", "", "email-help")
	if got == nil {
		t.Fatal("ErrorAttrs with helpTextID should not return nil")
	}
	if got[ariaDescribedBy] != "email-help" {
		t.Errorf("aria-describedby = %v, want %q", got[ariaDescribedBy], "email-help")
	}
	if _, ok := got[ariaInvalid]; ok {
		t.Error("should not have aria-invalid when there is no error")
	}
}
