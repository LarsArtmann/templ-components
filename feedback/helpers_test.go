// Package feedback provides tests for feedback components like Alert, Toast, Spinner, and Loading.
package feedback

import "testing"

func TestSpinnerSizeClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		size SpinnerSize
		want string
	}{
		{SpinnerSmall, "h-4 w-4"},
		{SpinnerMedium, "h-6 w-6"},
		{SpinnerLarge, "h-8 w-8"},
		{SpinnerSize("unknown"), "h-6 w-6"},
	}
	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			t.Parallel()
			got := spinnerSizeClass(tt.size)
			if got != tt.want {
				t.Errorf("spinnerSizeClass(%q) = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}

func TestToastStyles(t *testing.T) {
	t.Parallel()
	types := []ToastType{ToastSuccess, ToastError, ToastWarning, ToastInfo}
	for _, tt := range types {
		t.Run(string(tt), func(t *testing.T) {
			t.Parallel()
			border, bg, text, icon := toastStyles(tt)
			if border == "" || bg == "" || text == "" || icon == "" {
				t.Errorf("toastStyles(%q) returned empty value: border=%q bg=%q text=%q icon=%q", tt, border, bg, text, icon)
			}
		})
	}
}

func TestToastStylesDefault(t *testing.T) {
	t.Parallel()
	border, bg, text, icon := toastStyles("unknown")
	if border == "" || bg == "" || text == "" || icon == "" {
		t.Errorf("toastStyles(unknown) returned empty value: border=%q bg=%q text=%q icon=%q", border, bg, text, icon)
	}
}

func TestAlertStyles(t *testing.T) {
	t.Parallel()
	types := []AlertType{AlertSuccess, AlertError, AlertWarning, AlertInfo}
	for _, tt := range types {
		t.Run(string(tt), func(t *testing.T) {
			t.Parallel()
			border, bg, text, icon := alertStyles(tt)
			if border == "" || bg == "" || text == "" || icon == "" {
				t.Errorf("alertStyles(%q) returned empty value: border=%q bg=%q text=%q icon=%q", tt, border, bg, text, icon)
			}
		})
	}
}

func TestProgressHeightClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		size string
		want string
	}{
		{"sm", "h-1.5"},
		{"md", "h-2.5"},
		{"lg", "h-4"},
		{"unknown", "h-2.5"},
	}
	for _, tt := range tests {
		t.Run(tt.size, func(t *testing.T) {
			t.Parallel()
			got := progressHeightClass(tt.size)
			if got != tt.want {
				t.Errorf("progressHeightClass(%q) = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}

func TestStepLineClass(t *testing.T) {
	t.Parallel()
	if got := stepLineClass(0, 2); got != "bg-blue-600 dark:bg-blue-500" {
		t.Errorf("stepLineClass(0, 2) = %q", got)
	}
	if got := stepLineClass(3, 2); got != "bg-gray-200 dark:bg-gray-700" {
		t.Errorf("stepLineClass(3, 2) = %q", got)
	}
}

func TestStepCircleClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		step, current int
		wantContains  string
	}{
		{0, 2, "bg-blue-600"},
		{2, 2, "border-blue-600"},
		{3, 2, "bg-gray-100"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got := stepCircleClass(tt.step, tt.current)
			if got == "" {
				t.Errorf("stepCircleClass(%d, %d) returned empty", tt.step, tt.current)
			}
		})
	}
}
