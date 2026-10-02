package datastar

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// LoadingButtonProps configures a button label that swaps to a busy state
// while the enclosing button's Datastar action is in flight. This is the
// Datastar equivalent of htmx.LoadingButton.
//
// The component renders CONTENT, not the button itself: place it inside a
// button that carries both the action (data-on:click=...) and the indicator
// marker (data-indicator:<signal>), and pass the SAME signal name here.
type LoadingButtonProps struct {
	utils.BaseProps

	// Signal is the indicator signal name (without the $ prefix). Must
	// match the data-indicator:<name> attribute on the enclosing button.
	// An empty signal degrades gracefully: the default text always shows
	// and the busy state never triggers (mirrors Indicator's degradation).
	Signal string

	// DefaultText is the label shown at rest.
	DefaultText string

	// LoadingText is the label shown while the action is in flight.
	// An empty LoadingText keeps only the spinner in the busy state.
	LoadingText string

	// Spinner is the loading animation shown next to LoadingText. When nil,
	// a small inline CSS spinner is used (the same default as Indicator).
	Spinner templ.Component
}

// DefaultLoadingButtonProps returns sensible defaults for a loading button.
func DefaultLoadingButtonProps() LoadingButtonProps {
	return LoadingButtonProps{} //nolint:exhaustruct_v5 // intentionally minimal defaults
}

// loadingButtonBusyExpr returns the Datastar show expression for the busy
// state. An empty signal degrades to "false" — the bare "$" object is always
// truthy and would pin the busy state visible (the Indicator degradation
// rule).
func loadingButtonBusyExpr(signal string) string {
	if strings.TrimSpace(signal) == "" {
		return "false"
	}

	return "$" + signal
}

// loadingButtonRestExpr returns the Datastar show expression for the rest
// state: the negation of the busy expression. An empty signal degrades to
// "true" so the default label always shows.
func loadingButtonRestExpr(signal string) string {
	if strings.TrimSpace(signal) == "" {
		return "true"
	}

	return "!$" + signal
}
