package visualtest_test

import (
	"context"
	"io"
	"testing"

	"github.com/a-h/templ"
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
	visualtest.AssertScreenshot(t, "button/primary_dark", display.Button(primary), visualtest.Options{Dark: visualtest.Bool(true)})

	danger := display.DefaultButtonProps()
	danger.Text = "Delete account"
	danger.Variant = display.ButtonDanger
	visualtest.AssertScreenshot(t, "button/danger_light", display.Button(danger))

	secondary := display.DefaultButtonProps()
	secondary.Text = "Cancel"
	secondary.Variant = display.ButtonSecondary
	visualtest.AssertScreenshot(t, "button/secondary_dark", display.Button(secondary), visualtest.Options{Dark: visualtest.Bool(true)})
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
	visualtest.AssertScreenshot(t, "alert/error_dark", feedback.Alert(errAlert), visualtest.Options{Dark: visualtest.Bool(true)})

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
	visualtest.AssertScreenshot(t, "badge/error_pill_dark", display.Badge(errorBadge), visualtest.Options{Dark: visualtest.Bool(true)})
}

// TestCard covers the structural card with header + body.
func TestCard(t *testing.T) {
	t.Parallel()

	card := display.DefaultCardProps()
	card.Title = "Monthly revenue"
	visualtest.AssertScreenshot(t, "card/basic_light", display.Card(card))
	visualtest.AssertScreenshot(t, "card/basic_dark", display.Card(card), visualtest.Options{Dark: visualtest.Bool(true)})
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
// The modal is rendered with Open=true; the auto-open JS calls showModal(),
// promoting the dialog to the top layer. FullViewport captures it.
func TestModal(t *testing.T) {
	t.Parallel()

	modal := display.DefaultModalProps()
	modal.Title = "Delete project"
	modal.Open = true
	opts := dialogOpen(visualtest.Viewport{Width: 480, Height: 400})
	visualtest.AssertScreenshot(t, "modal/open_light", display.Modal(modal), opts)

	opts.Dark = visualtest.Bool(true)
	visualtest.AssertScreenshot(t, "modal/open_dark", display.Modal(modal), opts)
}

// TestDrawer covers the native <dialog> drawer on both sides.
func TestDrawer(t *testing.T) {
	t.Parallel()

	opts := dialogOpen(visualtest.Viewport{Width: 480, Height: 400})

	right := display.DefaultDrawerProps()
	right.Title = "Settings"
	right.Open = true
	visualtest.AssertScreenshot(t, "drawer/right_light", display.Drawer(right), opts)

	opts.Dark = visualtest.Bool(true)
	left := display.DefaultDrawerProps()
	left.Title = "Filters"
	left.Open = true
	left.Side = display.DrawerLeft
	visualtest.AssertScreenshot(t, "drawer/left_dark", display.Drawer(left), opts)
}

// TestInput covers text, error, and disabled states of the most-used form input.
func TestInput(t *testing.T) {
	t.Parallel()

	basic := forms.DefaultInputProps()
	basic.Label = "Email address"
	basic.Placeholder = "you@example.com"
	visualtest.AssertScreenshot(t, "input/text_light", forms.Input(basic))
	visualtest.AssertScreenshot(t, "input/text_dark", forms.Input(basic), visualtest.Options{Dark: visualtest.Bool(true)})

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
	visualtest.AssertScreenshot(t, "select/basic_dark", forms.Select(sel), visualtest.Options{Dark: visualtest.Bool(true)})
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
		visualtest.Options{RTL: visualtest.Bool(true)},
	)

	card := display.DefaultCardProps()
	card.Title = "Monthly revenue"
	card.Subtitle = "Net of taxes and fees"
	visualtest.AssertScreenshot(
		t,
		"card/basic_rtl",
		display.Card(card),
		visualtest.Options{RTL: visualtest.Bool(true)},
	)
}

// overlayOpen is the shared option set for components whose open state renders
// in the browser top layer (Popover API menus, <dialog>): click/context to
// open, wait for the [popover] panel to appear, then capture the full viewport
// (the panel paints outside #tc-root's box, so an element screenshot would
// crop it). Nonce is required so the components render their positioning
// scripts and — for ContextMenu — the menu + event handler at all.
//
// MaxMismatch is raised to 2%: these menus are positioned in JS from the
// trigger's getBoundingClientRect(), so a 1px layout-timing shift shows up as
// edge anti-aliasing variance (~0.5-1%). A real regression (missing menu,
// wrong colors, broken layout) blows far past 2%. Pure-CSS components stay at
// the strict 0.1% default.
func overlayOpen(viewport visualtest.Viewport, state visualtest.InteractionState) visualtest.Options {
	return visualtest.Options{
		State:        state,
		WaitSelector: "[popover]",
		FullViewport: true,
		Viewport:     viewport,
		MaxMismatch:  0.02,
	}
}

// dialogOpen returns options for native <dialog>-based overlays (Modal, Drawer)
// that are server-rendered with Open=true. The auto-open JS calls showModal(),
// promoting the dialog to the top layer — outside #tc-root's bounding box — so
// FullViewport is required. WaitSelector: "dialog" ensures the screenshot is
// taken after showModal() fires. MaxMismatch is raised to 2% to absorb
// top-layer positioning and backdrop anti-aliasing variance.
func dialogOpen(viewport visualtest.Viewport) visualtest.Options {
	return visualtest.Options{
		WaitSelector: "dialog",
		FullViewport: true,
		Viewport:     viewport,
		MaxMismatch:  0.02,
	}
}

// TestDropdownOpen covers the Dropdown's open state via the native Popover API
// (popovertarget invoker) in light and dark mode. This was T8 — skipped in the
// prior session as "needs click simulation"; the harness now supports it.
func TestDropdownOpen(t *testing.T) {
	t.Parallel()

	dropdown := func() display.DropdownProps {
		d := display.DefaultDropdownProps()
		d.Label = "Options"
		d.Nonce = "test-nonce"
		d.Items = []display.DropdownItem{
			{Text: "Edit", Href: "/edit"},
			{Text: "Duplicate", Href: "/duplicate"},
			{Text: "Archive", Href: "/archive"},
		}

		return d
	}

	light := overlayOpen(visualtest.Viewport{Width: 480, Height: 360}, visualtest.StateClick)
	visualtest.AssertScreenshot(t, "dropdown/open_light", display.Dropdown(dropdown()), light)

	dark := light
	dark.Dark = visualtest.Bool(true)
	visualtest.AssertScreenshot(t, "dropdown/open_dark", display.Dropdown(dropdown()), dark)
}

// TestPopoverOpen covers the Popover's open state. Popover takes HTML children
// for its panel content, injected via templ.WithChildren.
func TestPopoverOpen(t *testing.T) {
	t.Parallel()

	const content = `<div class="space-y-1">` +
		`<p class="font-medium text-gray-900 dark:text-white">Popover title</p>` +
		`<p class="text-gray-500 dark:text-gray-400">Some helpful detail shown in the floating panel.</p>` +
		`</div>`

	props := display.DefaultPopoverProps()
	props.TriggerText = "Details"
	props.Nonce = "test-nonce"

	visualtest.AssertScreenshot(
		t,
		"popover/open_light",
		withChildren(display.Popover(props), templ.Raw(content)),
		overlayOpen(visualtest.Viewport{Width: 480, Height: 400}, visualtest.StateClick),
	)
}

// TestContextMenuOpen covers the ContextMenu's open state. ContextMenu opens on
// a right-click (contextmenu event), so it uses StateContext rather than
// StateClick.
func TestContextMenuOpen(t *testing.T) {
	t.Parallel()

	const content = `<div class="p-4 rounded-md border border-dashed border-gray-300 dark:border-gray-600 ` +
		`text-sm text-gray-500 dark:text-gray-400">Right-click this region</div>`

	props := display.DefaultContextMenuProps()
	props.Nonce = "test-nonce"
	props.Items = []display.ContextMenuItem{
		{Text: "Edit", Href: "/edit"},
		{Text: "Copy", Href: "/copy"},
		{Text: "Delete", Href: "/delete"},
	}

	visualtest.AssertScreenshot(
		t,
		"contextmenu/open_light",
		withChildren(display.ContextMenu(props), templ.Raw(content)),
		overlayOpen(visualtest.Viewport{Width: 480, Height: 360}, visualtest.StateContext),
	)
}

// withChildren returns a component that renders parent with the given child
// injected into the context, so components that consume { children... }
// (Popover, ContextMenu) show their content in a visual golden.
func withChildren(parent, child templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return parent.Render(templ.WithChildren(ctx, child), w)
	})
}
