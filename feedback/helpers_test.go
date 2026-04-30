// Package feedback provides tests for feedback components like Alert, Toast, Spinner, and Loading.
package feedback

import (
	"fmt"
	"testing"
)

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
	tests := []struct {
		typ           ToastType
		wantBorder    string
		wantIconColor string
	}{
		{
			ToastSuccess,
			"border-green-200 dark:border-green-800",
			"text-green-600 dark:text-green-400",
		},
		{ToastError, "border-red-200 dark:border-red-800", "text-red-600 dark:text-red-400"},
		{
			ToastWarning,
			"border-yellow-200 dark:border-yellow-800",
			"text-yellow-600 dark:text-yellow-400",
		},
		{ToastInfo, "border-blue-200 dark:border-blue-800", "text-blue-600 dark:text-blue-400"},
	}
	for _, tt := range tests {
		t.Run(string(tt.typ), func(t *testing.T) {
			t.Parallel()
			border, bg, text, icon := toastStyles(tt.typ)
			if border == "" || bg == "" || text == "" || icon == "" {
				t.Errorf(
					"toastStyles(%q) returned empty value: border=%q bg=%q text=%q icon=%q",
					tt.typ,
					border,
					bg,
					text,
					icon,
				)
			}
			if border != tt.wantBorder {
				t.Errorf("toastStyles(%q) border = %q, want %q", tt.typ, border, tt.wantBorder)
			}
			if icon != tt.wantIconColor {
				t.Errorf("toastStyles(%q) icon = %q, want %q", tt.typ, icon, tt.wantIconColor)
			}
		})
	}
}

func TestToastStylesDefault(t *testing.T) {
	t.Parallel()
	border, bg, text, icon := toastStyles("unknown")
	if border == "" || bg == "" || text == "" || icon == "" {
		t.Errorf(
			"toastStyles(unknown) returned empty value: border=%q bg=%q text=%q icon=%q",
			border,
			bg,
			text,
			icon,
		)
	}
}

func TestAlertStyles(t *testing.T) {
	t.Parallel()
	tests := []struct {
		typ           AlertType
		wantBorder    string
		wantIconColor string
	}{
		{AlertSuccess, "border-green-200 dark:border-green-800", "text-green-400"},
		{AlertError, "border-red-200 dark:border-red-800", "text-red-400"},
		{AlertWarning, "border-yellow-200 dark:border-yellow-800", "text-yellow-400"},
		{AlertInfo, "border-blue-200 dark:border-blue-800", "text-blue-400"},
	}
	for _, tt := range tests {
		t.Run(string(tt.typ), func(t *testing.T) {
			t.Parallel()
			border, bg, text, icon := alertStyles(tt.typ)
			if border == "" || bg == "" || text == "" || icon == "" {
				t.Errorf(
					"alertStyles(%q) returned empty value: border=%q bg=%q text=%q icon=%q",
					tt.typ,
					border,
					bg,
					text,
					icon,
				)
			}
			if border != tt.wantBorder {
				t.Errorf("alertStyles(%q) border = %q, want %q", tt.typ, border, tt.wantBorder)
			}
			if icon != tt.wantIconColor {
				t.Errorf("alertStyles(%q) icon = %q, want %q", tt.typ, icon, tt.wantIconColor)
			}
		})
	}
}

func TestProgressHeightClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		size ProgressBarSize
		want string
	}{
		{ProgressBarSizeSM, "h-1.5"},
		{ProgressBarSizeMD, "h-2.5"},
		{ProgressBarSizeLG, "h-4"},
	}
	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
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
		want          string
	}{
		{0, 2, "bg-blue-600 text-white dark:bg-blue-500"},
		{
			2,
			2,
			"bg-white border-2 border-blue-600 text-blue-600 dark:border-blue-400 dark:text-blue-400",
		},
		{3, 2, "bg-gray-100 text-gray-500 dark:bg-gray-800 dark:text-gray-400"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("step-%d-current-%d", tt.step, tt.current), func(t *testing.T) {
			t.Parallel()
			got := stepCircleClass(tt.step, tt.current)
			if got != tt.want {
				t.Errorf("stepCircleClass(%d, %d) = %q, want %q", tt.step, tt.current, got, tt.want)
			}
		})
	}
}
