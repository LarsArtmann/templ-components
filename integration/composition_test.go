package integration_test

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
)

func TestFullPageComposition(t *testing.T) {
	t.Parallel()
	output := utils.RenderAll(
		t,
		navigation.SimpleNav(navigation.SimpleNavProps{
			BrandText: "TestApp",
			Links: []navigation.NavLinkProps{
				{Href: "/", Text: "Home"},
				{Href: "/settings", Text: "Settings"},
			},
			CurrentPath: "/",
		}),
		display.Card(display.CardProps{Title: "Welcome"}),
		feedback.Alert(feedback.AlertProps{
			Title:   "Heads up",
			Message: "This is a test alert",
			Type:    feedback.FeedbackInfo,
		}),
		forms.Input(forms.InputProps{
			Name:  "email",
			Label: "Email",
			Type:  forms.InputEmail,
		}),
	)

	if !strings.Contains(output, "TestApp") {
		t.Error("expected nav brand in output")
	}
	if !strings.Contains(output, "Welcome") {
		t.Error("expected card title in output")
	}
	if !strings.Contains(output, "Heads up") {
		t.Error("expected alert title in output")
	}
	if !strings.Contains(output, `name="email"`) {
		t.Error("expected input name in output")
	}
}

func TestFormWithMultipleInputs(t *testing.T) {
	t.Parallel()
	output := utils.RenderAll(
		t,
		forms.Input(forms.InputProps{
			BaseProps: utils.BaseProps{ID: "name"},
			Name:      "name",
			Label:     "Full Name",
			Required:  true,
		}),
		forms.Input(forms.InputProps{
			BaseProps: utils.BaseProps{ID: "email"},
			Name:      "email",
			Label:     "Email",
			Type:      forms.InputEmail,
			Required:  true,
			Error:     "Email already taken",
		}),
		forms.Select(forms.SelectProps{
			BaseProps: utils.BaseProps{ID: "role"},
			Name:      "role",
			Label:     "Role",
			Options: []forms.SelectOption{
				{Value: "admin", Label: "Admin"},
				{Value: "user", Label: "User", Selected: true},
			},
		}),
		forms.Combobox(forms.ComboboxProps{
			BaseProps: utils.BaseProps{ID: "country"},
			Name:      "country",
			Label:     "Country",
			Options: []forms.ComboboxOption{
				{Value: "de", Label: "Germany"},
				{Value: "at", Label: "Austria"},
			},
		}),
		display.Button(display.ButtonProps{
			Text: "Register",
		}),
	)

	utils.AssertContains(t, output, `for="name"`)
	utils.AssertContains(t, output, `for="email"`)
	utils.AssertContains(t, output, "Email already taken")
	utils.AssertContains(t, output, "Germany")
	utils.AssertContains(t, output, "Register")
}

func TestModalWithFormContent(t *testing.T) {
	t.Parallel()
	output := utils.RenderAll(
		t,
		forms.Input(forms.InputProps{
			BaseProps: utils.BaseProps{ID: "modal-name"},
			Name:      "name",
			Label:     "Name",
		}),
		display.Button(display.ButtonProps{
			Text: "Save",
		}),
	)

	utils.AssertContains(t, output, `name="name"`)
	utils.AssertContains(t, output, "Save")
}

func TestFeedbackComponentsCompose(t *testing.T) {
	t.Parallel()
	output := utils.RenderAll(
		t,
		feedback.Spinner(feedback.SpinnerProps{}),
		feedback.LoadingOverlay(feedback.LoadingOverlayProps{
			Message: "Loading data",
		}),
		feedback.ProgressBar(feedback.ProgressBarProps{
			Current: 60,
			Total:   100,
			Label:   "Uploading",
		}),
	)

	utils.AssertContains(t, output, "animate-spin")
	utils.AssertContains(t, output, "Loading data")
	utils.AssertContains(t, output, "60%")
}

func TestNavigationComponentsCompose(t *testing.T) {
	t.Parallel()
	output := utils.RenderAll(
		t,
		navigation.Breadcrumbs(navigation.BreadcrumbsProps{
			Items: []navigation.BreadcrumbItem{
				{Text: "Home", Href: "/"},
				{Text: "Users", Href: "/users"},
				{Text: "Profile", Active: true},
			},
		}),
		navigation.Pagination(navigation.PaginationProps{
			CurrentPage: 3,
			TotalPages:  10,
			BaseURL:     "/users",
			QueryParam:  "page",
		}),
	)

	utils.AssertContains(t, output, "Home")
	utils.AssertContains(t, output, "Profile")
	utils.AssertContains(t, output, "/users")
}

func TestDisplayComponentsCompose(t *testing.T) {
	t.Parallel()
	output := utils.RenderAll(
		t,
		display.Card(display.CardProps{Title: "Dashboard"}),
		display.StatCard(display.StatCardProps{
			Value:  "$12,345",
			Label:  "Revenue",
			Change: "+15%",
			Trend:  display.TrendUp,
		}),
		display.Table(display.TableProps{
			Headers: []string{"Name", "Status"},
			Rows: []display.TableRow{
				{Cells: []display.TableCell{{Text: "Alice"}, {Text: "Active"}}},
			},
		}),
	)

	utils.AssertContains(t, output, "Dashboard")
	utils.AssertContains(t, output, "$12,345")
	utils.AssertContains(t, output, "Alice")
}

func TestBasePropsPropagation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name   string
		render func() string
		checks []string
	}

	cases := []testCase{
		{
			name: "display.Card",
			render: func() string {
				return utils.Render(t, display.Card(display.CardProps{
					BaseProps: utils.BaseProps{ID: "c1", Class: "shadow-lg", AriaLabel: "card"},
				}))
			},
			checks: []string{`id="c1"`, "shadow-lg", `aria-label="card"`},
		},
		{
			name: "forms.Input",
			render: func() string {
				return utils.Render(t, forms.Input(forms.InputProps{
					BaseProps: utils.BaseProps{ID: "i1", AriaLabel: "input"},
					Name:      "test",
				}))
			},
			checks: []string{`id="i1"`, `aria-label="input"`},
		},
		{
			name: "feedback.Alert",
			render: func() string {
				return utils.Render(t, feedback.Alert(feedback.AlertProps{
					BaseProps: utils.BaseProps{ID: "a1", Class: "mt-4", AriaLabel: "alert"},
					Title:     "Test",
				}))
			},
			checks: []string{`id="a1"`, "mt-4", `aria-label="alert"`},
		},
		{
			name: "forms.DatePicker",
			render: func() string {
				return utils.Render(t, forms.DatePicker(forms.DatePickerProps{
					BaseProps: utils.BaseProps{ID: "dp1", AriaLabel: "pick date"},
				}))
			},
			checks: []string{`id="dp1"`, `type="date"`, `aria-label="pick date"`},
		},
		{
			name: "forms.Combobox",
			render: func() string {
				return utils.Render(t, forms.Combobox(forms.ComboboxProps{
					BaseProps: utils.BaseProps{ID: "cb1", AriaLabel: "search"},
					Name:      "x",
				}))
			},
			checks: []string{`role="combobox"`, `aria-label="search"`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			output := tc.render()
			for _, check := range tc.checks {
				if !strings.Contains(output, check) {
					t.Errorf("%s: expected %q in output", tc.name, check)
				}
			}
		})
	}
}

// TestCSPNonceConsistency verifies that components with inline scripts all
// use the nonce from their props, not hardcoded values.
func TestCSPNonceConsistency(t *testing.T) {
	t.Parallel()

	nonce := "test-nonce-12345"

	outputs := []string{
		utils.Render(t, display.Modal(display.ModalProps{
			BaseProps: utils.BaseProps{ID: "m1", Nonce: nonce},
		})),
		utils.Render(t, display.Dropdown(display.DropdownProps{
			BaseProps: utils.BaseProps{ID: "d1", Nonce: nonce},
			Label:     "Actions",
		})),
		utils.Render(t, display.Accordion(display.AccordionProps{
			BaseProps: utils.BaseProps{Nonce: nonce},
			Items:     []display.AccordionItem{{ID: "a1", Title: "Item 1"}},
		})),
		utils.Render(t, forms.Combobox(forms.ComboboxProps{
			BaseProps: utils.BaseProps{ID: "cb1", Nonce: nonce},
			Name:      "x",
		})),
		utils.Render(t, layout.ThemeScript(nonce)),
		utils.Render(t, layout.ThemeToggle("", nonce)),
	}

	for i, output := range outputs {
		if !strings.Contains(output, `nonce="`+nonce+`"`) {
			t.Errorf("component %d: expected nonce %q in script tag", i, nonce)
		}
	}
}
