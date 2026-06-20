package forms

import (
	"testing"
)

func TestInputTypeValidation(t *testing.T) {
	t.Parallel()

	t.Run("valid InputType passes through", func(t *testing.T) {
		t.Parallel()
		for _, tt := range []InputType{InputText, InputEmail, InputPassword, InputNumber, InputTel, InputURL, InputDate, InputTime, InputDatetime, InputSearch, InputHidden} {
			got := inputType(tt)
			if got != string(tt) {
				t.Errorf("inputType(%q) = %q, want %q", tt, got, string(tt))
			}
		}
	})

	t.Run("invalid InputType falls back to text", func(t *testing.T) {
		t.Parallel()
		for _, in := range []InputType{"", "javascript:alert(1)"} {
			if got := inputType(in); got != "text" {
				t.Errorf("inputType(%q) = %q, want %q", in, got, "text")
			}
		}
	})
}
