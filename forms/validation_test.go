package forms

import (
	"testing"
)

func TestInputTypeValidation(t *testing.T) {
	t.Parallel()

	t.Run("empty InputType defaults to text", func(t *testing.T) {
		t.Parallel()
		got := inputType("")
		if got != "text" {
			t.Errorf("inputType('') = %q, want %q", got, "text")
		}
	})

	t.Run("valid InputType passes through", func(t *testing.T) {
		t.Parallel()
		for _, tt := range []InputType{InputText, InputEmail, InputPassword, InputNumber, InputTel, InputURL, InputDate, InputTime, InputDatetime, InputSearch, InputHidden} {
			got := inputType(tt)
			if got != string(tt) {
				t.Errorf("inputType(%q) = %q, want %q", tt, got, string(tt))
			}
		}
	})

	t.Run("unknown InputType panics", func(t *testing.T) {
		t.Parallel()
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic for unknown InputType")
			}
		}()
		inputType("javascript:alert(1)")
	})
}
