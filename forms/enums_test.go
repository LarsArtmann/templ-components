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
