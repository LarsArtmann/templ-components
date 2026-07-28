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

// Options configures how a component is rendered and captured.
type Options struct {
	// Dark renders with the .dark class on <html> (the library's dark-mode
	// strategy). Default is light mode.
	Dark bool
	// RTL sets dir="rtl" on <html> to test logical-property mirroring.
	RTL bool
	// Viewport sets the emulated window size. Defaults to 1280x800.
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
func defaultOptions(o Options) Options {
	if o.Viewport.Width == 0 {
		o.Viewport.Width = 1280
	}

	if o.Viewport.Height == 0 {
		o.Viewport.Height = 800
	}

	if o.MaxMismatch == 0 {
		o.MaxMismatch = 0.001
	}

	if o.Threshold == 0 {
		o.Threshold = 0.1
	}

	return o
}

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

	if opts.Dark {
		htmlClass = ` class="dark"`
		bodyClass = "bg-gray-900 text-gray-100"
	}

	dir := ""
	if opts.RTL {
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
