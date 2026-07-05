package htmx

import "testing"

func TestIsValidEnums(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func() bool
		want bool
	}{
		{"SwapStyle InnerHTML", func() bool { return SwapStyleIsValid(SwapInnerHTML) }, true},
		{"SwapStyle OuterHTML", func() bool { return SwapStyleIsValid(SwapOuterHTML) }, true},
		{"SwapStyle invalid", func() bool { return SwapStyleIsValid(SwapStyle("bogus")) }, false},
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
