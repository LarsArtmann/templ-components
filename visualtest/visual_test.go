package visualtest_test

import (
	"testing"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/visualtest"
)

// TestButtons covers the button color system in both light and dark mode —
// the most common source of dark-mode regressions in this library.
func TestButtons(t *testing.T) {
	t.Parallel()

	primary := display.DefaultButtonProps()
	primary.Text = "Save changes"
	visualtest.AssertScreenshot(t, "button/primary_light", display.Button(primary))
	visualtest.AssertScreenshot(t, "button/primary_dark", display.Button(primary), visualtest.Options{Dark: true})

	danger := display.DefaultButtonProps()
	danger.Text = "Delete account"
	danger.Variant = display.ButtonDanger
	visualtest.AssertScreenshot(t, "button/danger_light", display.Button(danger))

	secondary := display.DefaultButtonProps()
	secondary.Text = "Cancel"
	secondary.Variant = display.ButtonSecondary
	visualtest.AssertScreenshot(t, "button/secondary_dark", display.Button(secondary), visualtest.Options{Dark: true})
}

// TestButtonStates covers the interactive states that are invisible to HTML
// golden tests: hover (color shift) and focus-visible (focus ring), plus the
// disabled (greyed-out) variant.
func TestButtonStates(t *testing.T) {
	t.Parallel()

	primary := display.DefaultButtonProps()
	primary.Text = "Save changes"
	visualtest.AssertScreenshot(
		t,
		"button/primary_hover",
		display.Button(primary),
		visualtest.Options{State: visualtest.StateHover},
	)
	visualtest.AssertScreenshot(
		t,
		"button/primary_focus",
		display.Button(primary),
		visualtest.Options{State: visualtest.StateFocus},
	)

	disabled := display.DefaultButtonProps()
	disabled.Text = "Save changes"
	disabled.Disabled = true
	visualtest.AssertScreenshot(t, "button/primary_disabled", display.Button(disabled))
}

// TestAlerts covers all four feedback types with icons and colors.
func TestAlerts(t *testing.T) {
	t.Parallel()

	success := feedback.DefaultAlertProps()
	success.Title = "Payment received"
	success.Message = "Your subscription is now active."
	success.Type = feedback.AlertSuccess
	visualtest.AssertScreenshot(t, "alert/success_light", feedback.Alert(success))

	errAlert := feedback.DefaultAlertProps()
	errAlert.Title = "Could not save"
	errAlert.Message = "The server rejected the request. Try again in a moment."
	errAlert.Type = feedback.AlertError
	visualtest.AssertScreenshot(t, "alert/error_dark", feedback.Alert(errAlert), visualtest.Options{Dark: true})

	warn := feedback.DefaultAlertProps()
	warn.Title = "Storage almost full"
	warn.Message = "You have used 92% of your quota."
	warn.Type = feedback.AlertWarning
	visualtest.AssertScreenshot(t, "alert/warning_light", feedback.Alert(warn))
}

// TestBadges covers the badge color + dot/pill variants.
func TestBadges(t *testing.T) {
	t.Parallel()

	success := display.DefaultBadgeProps()
	success.Text = "Active"
	success.Type = display.BadgeSuccess
	success.Dot = true
	visualtest.AssertScreenshot(t, "badge/success_dot_light", display.Badge(success))

	errorBadge := display.DefaultBadgeProps()
	errorBadge.Text = "Failed"
	errorBadge.Type = display.BadgeError
	errorBadge.Pill = true
	visualtest.AssertScreenshot(t, "badge/error_pill_dark", display.Badge(errorBadge), visualtest.Options{Dark: true})
}

// TestCard covers the structural card with header + body.
func TestCard(t *testing.T) {
	t.Parallel()

	card := display.DefaultCardProps()
	card.Title = "Monthly revenue"
	visualtest.AssertScreenshot(t, "card/basic_light", display.Card(card))
	visualtest.AssertScreenshot(t, "card/basic_dark", display.Card(card), visualtest.Options{Dark: true})
}

// TestResponsiveViewport captures a card at a mobile viewport width to catch
// responsive breakpoint regressions (padding collapses, text wrapping).
func TestResponsiveViewport(t *testing.T) {
	t.Parallel()

	card := display.DefaultCardProps()
	card.Title = "Monthly revenue"
	card.Subtitle = "Net of taxes and fees"
	visualtest.AssertScreenshot(t, "card/mobile", display.Card(card),
		visualtest.Options{Viewport: visualtest.Viewport{Width: 375, Height: 667}}) // iPhone SE width
}

// TestModal covers the native <dialog> modal in both light and dark mode.
// The modal is rendered with Open=true so the dialog is visible.
func TestModal(t *testing.T) {
	t.Parallel()

	modal := display.DefaultModalProps()
	modal.Title = "Delete project"
	modal.Open = true
	visualtest.AssertScreenshot(t, "modal/open_light", display.Modal(modal))
	visualtest.AssertScreenshot(t, "modal/open_dark", display.Modal(modal), visualtest.Options{Dark: true})
}

// TestDrawer covers the native <dialog> drawer on both sides.
func TestDrawer(t *testing.T) {
	t.Parallel()

	right := display.DefaultDrawerProps()
	right.Title = "Settings"
	right.Open = true
	visualtest.AssertScreenshot(t, "drawer/right_light", display.Drawer(right))

	left := display.DefaultDrawerProps()
	left.Title = "Filters"
	left.Open = true
	left.Side = display.DrawerLeft
	visualtest.AssertScreenshot(
		t,
		"drawer/left_dark",
		display.Drawer(left),
		visualtest.Options{Dark: true},
	)
}

// TestInput covers text, error, and disabled states of the most-used form input.
func TestInput(t *testing.T) {
	t.Parallel()

	basic := forms.DefaultInputProps()
	basic.Label = "Email address"
	basic.Placeholder = "you@example.com"
	visualtest.AssertScreenshot(t, "input/text_light", forms.Input(basic))
	visualtest.AssertScreenshot(t, "input/text_dark", forms.Input(basic), visualtest.Options{Dark: true})

	withError := forms.DefaultInputProps()
	withError.Label = "Email address"
	withError.Value = "not-an-email"
	withError.Error = "Please enter a valid email address"
	visualtest.AssertScreenshot(t, "input/error_light", forms.Input(withError))

	disabled := forms.DefaultInputProps()
	disabled.Label = "API key"
	disabled.Value = "sk-••••••••"
	disabled.Disabled = true
	visualtest.AssertScreenshot(t, "input/disabled_light", forms.Input(disabled))
}

// TestSelect covers the select component with options.
func TestSelect(t *testing.T) {
	t.Parallel()

	sel := forms.DefaultSelectProps()
	sel.Label = "Country"
	sel.Options = []forms.SelectOption{
		{Value: "us", Label: "United States"},
		{Value: "de", Label: "Germany"},
		{Value: "jp", Label: "Japan"},
	}
	visualtest.AssertScreenshot(t, "select/basic_light", forms.Select(sel))
	visualtest.AssertScreenshot(t, "select/basic_dark", forms.Select(sel), visualtest.Options{Dark: true})
}

// TestRTL verifies that logical CSS properties (ms-, me-, ps-, pe-, start-,
// end-) correctly mirror in right-to-left mode. Uses Button and Card — the
// most common components where RTL breakage is user-visible.
func TestRTL(t *testing.T) {
	t.Parallel()

	primary := display.DefaultButtonProps()
	primary.Text = "Save changes"
	visualtest.AssertScreenshot(
		t,
		"button/primary_rtl",
		display.Button(primary),
		visualtest.Options{RTL: true},
	)

	card := display.DefaultCardProps()
	card.Title = "Monthly revenue"
	card.Subtitle = "Net of taxes and fees"
	visualtest.AssertScreenshot(
		t,
		"card/basic_rtl",
		display.Card(card),
		visualtest.Options{RTL: true},
	)
}
