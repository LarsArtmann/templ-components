package visualtest

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

// Viewport sets the browser window size for a screenshot.
type Viewport struct {
	Width  int
	Height int
}

// Device viewport dimensions in CSS pixels, matching real devices. The viewport
// presets compose these so the magic-number linter stays happy without burying
// the values.
const (
	viewportMobileWidth   = 375
	viewportMobileHeight  = 667
	viewportTabletWidth   = 768
	viewportTabletHeight  = 1024
	viewportDesktopWidth  = 1280
	viewportDesktopHeight = 800
)

// Common viewport presets matching real device CSS widths. Use them as
// Options{Viewport: visualtest.ViewportMobile} to test responsive breakpoints.
//
//nolint:gochecknoglobals // viewport presets are intentional package-level constants consumed by callers
var (
	ViewportMobile  = Viewport{Width: viewportMobileWidth, Height: viewportMobileHeight}   // iPhone SE / small phone
	ViewportTablet  = Viewport{Width: viewportTabletWidth, Height: viewportTabletHeight}   // iPad portrait
	ViewportDesktop = Viewport{Width: viewportDesktopWidth, Height: viewportDesktopHeight} // default laptop
)

// InteractionState selects an interactive state to apply to #tc-root before
// capturing, so hover/focus styles are covered by golden screenshots.
type InteractionState int

const (
	// StateRest is the default: no interaction.
	StateRest InteractionState = iota
	// StateHover moves the mouse over #tc-root, triggering :hover styles.
	StateHover
	// StateFocus focuses #tc-root, triggering :focus-visible styles.
	StateFocus
	// StateClick clicks the first interactive descendant of #tc-root (a
	// [popovertarget] trigger, button, or link). Use it with FullViewport
	// and WaitSelector to capture components whose open state renders in the
	// top layer — Dropdown, Popover, ContextMenu (native Popover API), Modal,
	// and Drawer (native <dialog>).
	StateClick
	// StateContext dispatches a contextmenu (right-click) event on the first
	// [data-tc-ctxmenu-trigger] (or interactive descendant) of #tc-root. This
	// opens components that activate on right-click, e.g. ContextMenu.
	StateContext
)

// String returns a human-readable name for the interaction state, used in
// failure messages and test output for quick diagnosis.
func (s InteractionState) String() string {
	switch s {
	case StateRest:
		return "rest"
	case StateHover:
		return "hover"
	case StateFocus:
		return "focus"
	case StateClick:
		return "click"
	case StateContext:
		return "context"
	default:
		return fmt.Sprintf("unknown(%d)", int(s))
	}
}

// Bool is a convenience helper for setting *bool option fields. It returns a
// pointer to b so callers can distinguish "unset" (nil) from an explicit true
// or false — the tri-state that bare bool fields cannot express.
//
//	visualtest.Options{Dark: visualtest.Bool(true)}  // explicit dark mode
//	visualtest.Options{Dark: visualtest.Bool(false)} // explicit light mode
//	visualtest.Options{}                             // unset → default (light)

//nolint:modernize // wrapper exists for API stability + tri-state clarity, not micro-optimization
func Bool(b bool) *bool { return &b }

// Options configures how a component is rendered and captured.
type Options struct {
	// Dark renders with the .dark class on <html> (the library's dark-mode
	// strategy). nil or false renders light mode; use Bool(true) for dark.
	// The pointer (not bare bool) lets callers distinguish "unset" from
	// "explicitly light" in test metadata.
	Dark *bool
	// RTL sets dir="rtl" on <html> to test logical-property mirroring.
	// nil or false renders LTR; use Bool(true) for RTL.
	RTL *bool
	// Viewport sets the emulated window size. Defaults to 1280x800
	// (ViewportDesktop). See ViewportMobile / ViewportTablet presets.
	Viewport Viewport
	// MaxMismatch is the largest fraction of mismatched pixels (0–1) that
	// still passes. Defaults to 0.001 (0.1%). Anti-aliasing noise is filtered
	// out by pixelmatch; a real regression blows past this.
	MaxMismatch float64
	// Threshold is the pixelmatch perceptual color-distance threshold (0–1).
	// Pixels whose YIQ distance is below it count as identical. Default 0.1;
	// raise it to tolerate more rendering noise, lower it for stricter checks.
	Threshold float64
	// State applies an interaction (hover/focus/click) to #tc-root before capture.
	State InteractionState
	// FullViewport captures the full browser viewport instead of a tightly
	// cropped #tc-root element. Required for components whose open state
	// renders in the top layer (Popover API menus, <dialog>), because that
	// content is painted outside #tc-root's bounding box and would be cropped
	// by an element screenshot.
	FullViewport bool
	// WaitSelector, when set, is waited for (chromedp.WaitVisible) AFTER the
	// interaction state is applied and BEFORE capture. Use it with StateClick
	// to wait for an overlay menu to appear, e.g. WaitSelector: "[popover]".
	WaitSelector string
}

// defaultOptions fills zero values with sensible defaults.
func defaultOptions(options Options) Options {
	if options.Viewport.Width == 0 {
		options.Viewport.Width = ViewportDesktop.Width
	}

	if options.Viewport.Height == 0 {
		options.Viewport.Height = ViewportDesktop.Height
	}

	if options.MaxMismatch == 0 {
		options.MaxMismatch = 0.001
	}

	if options.Threshold == 0 {
		options.Threshold = 0.1
	}

	return options
}

// isDark reports whether the Dark option is explicitly true.
func isDark(o Options) bool { return o.Dark != nil && *o.Dark }

// isRTL reports whether the RTL option is explicitly true.
func isRTL(o Options) bool { return o.RTL != nil && *o.RTL }

// renderHTML wraps a component in a standalone HTML document with the compiled
// Tailwind CSS inline. The component is placed inside #tc-root so it can be
// screenshot as a tightly-cropped element.
func renderHTML(component templ.Component, opts Options) (string, error) {
	css, err := loadCSS()
	if err != nil {
		return "", err
	}

	var body strings.Builder
	if err := component.Render(context.Background(), &body); err != nil {
		return "", fmt.Errorf("render component: %w", err)
	}

	htmlClass := ""
	bodyClass := "bg-white text-gray-900"

	if isDark(opts) {
		htmlClass = ` class="dark"`
		bodyClass = "bg-gray-900 text-gray-100"
	}

	dir := ""
	if isRTL(opts) {
		dir = ` dir="rtl"`
	}

	return fmt.Sprintf(pageTemplate, htmlClass, dir, css, bodyClass, body.String()), nil
}

// pageTemplate is the isolated document shell. The CSS is inlined so the page
// is fully self-contained (no HTTP server needed for static renders).
//

const pageTemplate = `<!DOCTYPE html>
<html lang="en"%s%s>
<head>
<meta charset="utf-8">
<meta name="color-scheme" content="light dark">
<style>%s</style>
<style>
  /* Neutralize the page chrome so the element screenshot is just the component. */
  html, body { margin: 0; }
  #tc-root { display: inline-block; padding: 16px; }
</style>
</head>
<body class="%s">
  <div id="tc-root">%s</div>
</body>
</html>`
