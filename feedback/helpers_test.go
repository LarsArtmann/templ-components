// Package feedback provides tests for feedback components like Alert, Toast, Spinner, and Loading.
package feedback

import (
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

const (
	spinnerMDClass      = "h-6 w-6"
	borderGreenClass    = "border-green-200 dark:border-green-800"
	borderRedClass      = "border-red-200 dark:border-red-800"
	borderYellowClass   = "border-yellow-200 dark:border-yellow-800"
	borderBlueClass     = "border-blue-200 dark:border-blue-800"
	stepLineActiveClass = "bg-blue-600 dark:bg-blue-500"
	textBlueInfoClass   = "text-blue-600 dark:text-blue-400"
)

func TestSpinnerSizeClass(t *testing.T) {
	t.Parallel()

	tests := []struct {
		size SpinnerSize
		want string
	}{
		{SpinnerSM, "h-4 w-4"},
		{SpinnerMD, spinnerMDClass},
		{SpinnerLG, "h-8 w-8"},
		{SpinnerSize("unknown"), spinnerMDClass},
	}
	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				fmt.Sprintf("spinnerSizeClass(%q)", tt.size),
				spinnerSizeClass(tt.size),
				tt.want,
			)
		})
	}
}

func assertStyleFunc4(
	t *testing.T,
	funcName, border, bg, text, icon, wantBorder, wantIconColor string,
) {
	t.Helper()

	if border == "" || bg == "" || text == "" || icon == "" {
		t.Errorf(
			"%s returned empty value: border=%q bg=%q text=%q icon=%q",
			funcName,
			border,
			bg,
			text,
			icon,
		)
	}

	if wantBorder != "" && border != wantBorder {
		t.Errorf("%s border = %q, want %q", funcName, border, wantBorder)
	}

	if wantIconColor != "" && icon != wantIconColor {
		t.Errorf("%s icon = %q, want %q", funcName, icon, wantIconColor)
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
			borderGreenClass,
			"text-green-600 dark:text-green-400",
		},
		{ToastError, borderRedClass, "text-red-600 dark:text-red-400"},
		{
			ToastWarning,
			borderYellowClass,
			"text-yellow-600 dark:text-yellow-400",
		},
		{ToastInfo, borderBlueClass, textBlueInfoClass},
	}
	for _, tt := range tests {
		t.Run(string(tt.typ), func(t *testing.T) {
			t.Parallel()

			s := feedbackStyle(tt.typ)
			assertStyleFunc4(
				t,
				fmt.Sprintf("feedbackStyle(%q)", tt.typ),
				s.Border,
				s.BG,
				s.Text,
				s.Icon,
				tt.wantBorder,
				tt.wantIconColor,
			)
		})
	}
}

func TestToastStylesDefault(t *testing.T) {
	t.Parallel()

	s := feedbackStyle("unknown")
	assertStyleFunc4(t, "feedbackStyle(unknown)", s.Border, s.BG, s.Text, s.Icon, "", "")
}

func TestAlertStyles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		typ           AlertType
		wantBorder    string
		wantIconColor string
	}{
		{AlertSuccess, borderGreenClass, "text-green-600 dark:text-green-400"},
		{AlertError, borderRedClass, "text-red-600 dark:text-red-400"},
		{AlertWarning, borderYellowClass, "text-yellow-600 dark:text-yellow-400"},
		{AlertInfo, borderBlueClass, textBlueInfoClass},
	}
	for _, tt := range tests {
		t.Run(string(tt.typ), func(t *testing.T) {
			t.Parallel()

			s := feedbackStyle(tt.typ)
			assertStyleFunc4(
				t,
				fmt.Sprintf("feedbackStyle(%q)", tt.typ),
				s.Border,
				s.BG,
				s.Text,
				s.Icon,
				tt.wantBorder,
				tt.wantIconColor,
			)
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
			utils.AssertEqual(
				t,
				fmt.Sprintf("progressHeightClass(%q)", tt.size),
				progressHeightClass(tt.size),
				tt.want,
			)
		})
	}
}

func TestStepLineClass(t *testing.T) {
	t.Parallel()

	tests := []struct {
		step, current int
		want          string
	}{
		{0, 2, stepLineActiveClass},
		{3, 2, "bg-gray-200 dark:bg-gray-700"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("step-%d-current-%d", tt.step, tt.current), func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				fmt.Sprintf("stepLineClass(%d, %d)", tt.step, tt.current),
				stepLineClass(tt.step, tt.current),
				tt.want,
			)
		})
	}
}

func TestDismissScript(t *testing.T) {
	t.Parallel()

	s := utils.DismissScript()
	if s == "" {
		t.Error("utils.DismissScript() returned empty string")
	}

	if !strings.Contains(s, "tcDismissAttached") {
		t.Error("utils.DismissScript() missing tcDismissAttached guard")
	}

	if !strings.Contains(s, "data-dismiss") {
		t.Error("utils.DismissScript() missing data-dismiss selector")
	}
}

func TestFeedbackIcon(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		typ  FeedbackType
		want icons.Name
	}{
		{"success maps to CheckCircle", FeedbackSuccess, icons.CheckCircle},
		{"error maps to XCircle", FeedbackError, icons.XCircle},
		{"unknown falls back to Information", "unknown", icons.Information},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := feedbackIcon(tt.typ)
			if got != tt.want {
				t.Errorf("feedbackIcon(%q) = %q, want %q", tt.typ, got, tt.want)
			}
		})
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
			"bg-white dark:bg-gray-800 border-2 border-blue-600 text-blue-600 dark:border-blue-400 dark:text-blue-400",
		},
		{3, 2, "bg-gray-100 text-gray-500 dark:bg-gray-800 dark:text-gray-400"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("step-%d-current-%d", tt.step, tt.current), func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				fmt.Sprintf("stepCircleClass(%d, %d)", tt.step, tt.current),
				stepCircleClass(tt.step, tt.current),
				tt.want,
			)
		})
	}
}
