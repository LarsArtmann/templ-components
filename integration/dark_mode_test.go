package integration

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/errorpage"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
)

// TestDarkModeVariantsPresent renders key components and asserts that their
// HTML output contains at least one dark: class. This is a smoke test — the
// comprehensive scanning is done by utils.TestDarkModeCompliance.
func TestDarkModeVariantsPresent(t *testing.T) {
	t.Parallel()

	renderings := []struct {
		name string
		html string
	}{
		{"Button", utils.Render(t, display.Button(display.ButtonProps{Variant: display.ButtonPrimary}))},
		{"Card", utils.Render(t, display.Card(display.CardProps{Title: "T"}))},
		{"Badge", utils.Render(t, display.Badge(display.BadgeProps{Text: "B", Type: display.BadgeSuccess}))},
		{"Avatar", utils.Render(t, display.Avatar(display.AvatarProps{Initials: "AB"}))},
		{"Alert", utils.Render(t, feedback.Alert(feedback.AlertProps{Type: feedback.AlertInfo, Message: "M"}))},
		{
			"Spinner",
			utils.Render(t, feedback.Spinner(feedback.SpinnerProps{Color: "text-blue-600 dark:text-blue-400"})),
		},
		{"ProgressBar", utils.Render(t, feedback.ProgressBar(feedback.ProgressBarProps{Current: 50, Total: 100}))},
		{"Input", utils.Render(t, forms.Input(forms.InputProps{Label: "L", Name: "n"}))},
		{"Toggle", utils.Render(t, forms.Toggle(forms.ToggleProps{Label: "L"}))},
		{
			"Select",
			utils.Render(
				t,
				forms.Select(
					forms.SelectProps{Label: "L", Name: "n", Options: []forms.SelectOption{{Value: "a", Label: "A"}}},
				),
			),
		},
		{"Nav", utils.Render(t, navigation.Nav(navigation.DefaultNavProps()))},
		{
			"Breadcrumbs",
			utils.Render(
				t,
				navigation.Breadcrumbs(
					navigation.BreadcrumbsProps{Items: []navigation.BreadcrumbItem{{Text: "H", Href: "/"}}},
				),
			),
		},
		{
			"Pagination",
			utils.Render(t, navigation.Pagination(navigation.PaginationProps{CurrentPage: 1, TotalPages: 3})),
		},
		{"ErrorPage", utils.Render(t, errorpage.ErrorPage(errorpage.DefaultErrorPageProps()))},
		{"NotFound404", utils.Render(t, errorpage.NotFound404(errorpage.DefaultNotFound404Props()))},
		{"ThemeToggle", utils.Render(t, layout.ThemeToggle("Toggle", "nonce"))},
		{
			"StepIndicator",
			utils.Render(
				t,
				feedback.StepIndicator(feedback.StepIndicatorProps{Steps: []string{"A", "B"}, CurrentStep: 0}),
			),
		},
	}

	for _, r := range renderings {
		if !strings.Contains(r.html, "dark:") {
			t.Errorf("%s: output contains no dark: classes", r.name)
		}
	}
}

// TestBasePageHasColorScheme verifies that layout.Base includes color-scheme
// CSS for native form control rendering in both light and dark modes.
func TestBasePageHasColorScheme(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, layout.Base(layout.DefaultPageProps()))
	if !strings.Contains(html, "color-scheme") {
		t.Error("layout.Base output does not contain color-scheme CSS")
	}
}
