// Package forms provides tests for form components like Input, Select, Textarea, and helpers.
package forms

import "testing"

const testHelloString = "hello"

func TestFormatFloat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		v        float64
		decimals int
		want     string
	}{
		{name: "zero returns empty", v: 0, decimals: 2, want: ""},
		{name: "integer", v: 42, decimals: 0, want: "42"},
		{name: "two decimals", v: 3.14159, decimals: 2, want: "3.14"},
		{name: "negative", v: -2.5, decimals: 1, want: "-2.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := FormatFloat(tt.v, tt.decimals)
			if got != tt.want {
				t.Errorf("FormatFloat(%v, %d) = %q, want %q", tt.v, tt.decimals, got, tt.want)
			}
		})
	}
}

func TestIsSelected(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		value  string
		option string
		want   bool
	}{
		{name: "matches", value: "a", option: "a", want: true},
		{name: "no match", value: "a", option: "b", want: false},
		{name: "empty", value: "", option: "", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := IsSelected(tt.value, tt.option)
			if got != tt.want {
				t.Errorf("IsSelected(%q, %q) = %v, want %v", tt.value, tt.option, got, tt.want)
			}
		})
	}
}

func TestIfNotNil(t *testing.T) {
	t.Parallel()
	type item struct {
		Name string
	}

	t.Run("non-nil", func(t *testing.T) {
		t.Parallel()
		obj := &item{Name: testHelloString}
		got := IfNotNil(obj, func(i item) string { return i.Name })
		if got != testHelloString {
			t.Errorf("IfNotNil() = %q, want %q", got, testHelloString)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		var obj *item
		got := IfNotNil(obj, func(i item) string { return i.Name })
		if got != "" {
			t.Errorf("IfNotNil() = %q, want empty", got)
		}
	})
}

func TestIfNotNilString(t *testing.T) {
	t.Parallel()
	t.Run("non-nil", func(t *testing.T) {
		t.Parallel()
		s := testHelloString
		got := IfNotNilString(&s)
		if got != testHelloString {
			t.Errorf("IfNotNilString() = %q, want %q", got, testHelloString)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		got := IfNotNilString(nil)
		if got != "" {
			t.Errorf("IfNotNilString() = %q, want empty", got)
		}
	})
}

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

func TestBoolToString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		b    bool
		want string
	}{
		{name: "true", b: true, want: "true"},
		{name: "false", b: false, want: "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := BoolToString(tt.b)
			if got != tt.want {
				t.Errorf("BoolToString(%v) = %q, want %q", tt.b, got, tt.want)
			}
		})
	}
}
