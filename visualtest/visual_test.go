package visualtest_test

import (
	"testing"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
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
	visualtest.AssertScreenshot(t, "button/primary_hover", display.Button(primary), visualtest.Options{State: visualtest.StateHover})
	visualtest.AssertScreenshot(t, "button/primary_focus", display.Button(primary), visualtest.Options{State: visualtest.StateFocus})

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
