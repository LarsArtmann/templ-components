package visualtest

import (
	"errors"
	"fmt"
)

// Sentinel errors for the visualtest package.
var (
	// errNoBrowser is returned when no Chromium binary is available. Tests
	// that encounter it skip rather than fail.
	errNoBrowser = errors.New("no Chromium binary found")

	// errHoverElementNotFound is wrapped with the selector when hoverAction
	// cannot find a target element. Use fmt.Errorf("...: %w", errHoverElementNotFound, ...)
	// at the call site to attach the selector.
	errHoverElementNotFound = errors.New("hover: element not found")

	// errFocusNoElement is wrapped with the selector when focusAction cannot
	// find any focusable descendant.
	errFocusNoElement = errors.New("focus: no focusable element under selector")

	// errClickNoElement is wrapped with the selector when clickAction cannot
	// find any clickable descendant.
	errClickNoElement = errors.New("click: no clickable element under selector")

	// errContextNoTrigger is wrapped with the selector when contextAction
	// cannot find any trigger element.
	errContextNoTrigger = errors.New("context: no trigger element under selector")
)

// wrapSelector attaches the selector string to a sentinel error so callers
// see both the cause and the affected selector.
func wrapSelector(sentinel error, sel string) error {
	return fmt.Errorf("%w: %q", sentinel, sel)
}
