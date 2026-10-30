package forms

import "testing"

func TestIsValidEnums(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func() bool
		want bool
	}{
		{"ToggleSize SM", func() bool { return ToggleSizeIsValid(ToggleSizeSM) }, true},
		{"ToggleSize invalid", func() bool { return ToggleSizeIsValid(ToggleSize("bogus")) }, false},
		{"InputType Text", func() bool { return InputTypeIsValid(InputText) }, true},
		{"InputType Email", func() bool { return InputTypeIsValid(InputEmail) }, true},
		{"InputType invalid", func() bool { return InputTypeIsValid(InputType("bogus")) }, false},
		{"FormMethod GET", func() bool { return FormMethodIsValid(FormGet) }, true},
		{"FormMethod POST", func() bool { return FormMethodIsValid(FormPost) }, true},
		{"FormMethod invalid", func() bool { return FormMethodIsValid(FormMethod("bogus")) }, false},
		{"RatingSize SM", func() bool { return RatingSizeIsValid(RatingSizeSM) }, true},
		{"RatingSize MD", func() bool { return RatingSizeIsValid(RatingSizeMD) }, true},
		{"RatingSize LG", func() bool { return RatingSizeIsValid(RatingSizeLG) }, true},
		{"RatingSize invalid", func() bool { return RatingSizeIsValid(RatingSize("bogus")) }, false},
		{"EnterKeyHint Send", func() bool { return EnterKeyHintTypeIsValid(EnterKeyHintSend) }, true},
		{"EnterKeyHint Search", func() bool { return EnterKeyHintTypeIsValid(EnterKeyHintSearch) }, true},
		{"EnterKeyHint invalid", func() bool { return EnterKeyHintTypeIsValid(EnterKeyHintType("bogus")) }, false},
		{"Autocomplete current-password", func() bool { return AutocompleteTypeIsValid(AutocompleteCurrentPassword) }, true},
		{"Autocomplete one-time-code", func() bool { return AutocompleteTypeIsValid(AutocompleteOneTimeCode) }, true},
		{"Autocomplete cc-number", func() bool { return AutocompleteTypeIsValid(AutocompleteCCNumber) }, true},
		{"Autocomplete invalid", func() bool { return AutocompleteTypeIsValid(AutocompleteType("bogus")) }, false},
		{"Autocomplete empty", func() bool { return AutocompleteTypeIsValid(AutocompleteType("")) }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.fn(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
