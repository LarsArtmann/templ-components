package datastar

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// IndicatorProps configures a loading indicator that shows/hides based on a
// Datastar indicator signal. Place data-indicator:<name> on the element that
// triggers a backend action, then render an Indicator wherever you want the
// spinner to appear.
//
// This is the Datastar equivalent of htmx.LoadingIndicator / InlineLoadingOverlay.
//
//	<button data-on:click="@post('/save')" data-indicator:saving>Save</button>
//	@datastar.Indicator(datastar.IndicatorProps{Signal: "saving"})
type IndicatorProps struct {
	utils.BaseProps

	// Signal is the indicator signal name (without the $ prefix).
	// Must match the data-indicator:<name> on the triggering element.
	Signal string

	// Spinner is the loading animation to display while the signal is true.
	// Typically feedback.Spinner(...). When nil, a CSS pulse is shown.
	Spinner templ.Component
}

// DefaultIndicatorProps returns sensible defaults for an indicator.
func DefaultIndicatorProps() IndicatorProps {
	return IndicatorProps{} //nolint:exhaustruct // intentionally minimal defaults
}

// indicatorSignalExpr returns the Datastar show expression for an indicator.
// Returns "$<signal>" which is true during in-flight requests.
// An empty signal degrades to "false" so the indicator stays hidden — the
// bare "$" object is always truthy and would pin the spinner visible.
func indicatorSignalExpr(signal string) string {
	if signal == "" {
		return "false"
	}

	return "$" + signal
}

// Get returns a Datastar @get('url') action expression for use in data-on:*
// or data-init attributes. The response should be an SSE stream.
//
//	data-on:click={ datastar.Get("/api/search") }
func Get(url string) string {
	return actionExpr("get", url)
}

// getActionExpr builds an @get expression with retry and requestCancellation
// options. RetryAuto and CancellationNone are the runtime defaults and are
// omitted; any other valid combination emits only the non-default options.
func getActionExpr(url string, retry RetryMode, cancellation RequestCancellation) string {
	escaped := strings.ReplaceAll(url, "'", "\\'")

	var opts []string

	if mode := retryModeValue(retry); mode != RetryAuto {
		opts = append(opts, fmt.Sprintf("retry: '%s'", mode))
	}

	if requestCancellationValue(cancellation) == CancellationCleanup {
		opts = append(opts, "requestCancellation: 'cleanup'")
	}

	if len(opts) == 0 {
		return fmt.Sprintf("@get('%s')", escaped)
	}

	return fmt.Sprintf("@get('%s', {%s})", escaped, strings.Join(opts, ", "))
}

// Post returns a Datastar @post('url') action expression.
// Sends all signals as the request body.
func Post(url string) string {
	return actionExpr("post", url)
}

// Put returns a Datastar @put('url') action expression.
func Put(url string) string {
	return actionExpr("put", url)
}

// Patch returns a Datastar @patch('url') action expression.
func Patch(url string) string {
	return actionExpr("patch", url)
}

// Delete returns a Datastar @delete('url') action expression.
func Delete(url string) string {
	return actionExpr("delete", url)
}

// actionExpr builds a Datastar backend action expression string.
// Escapes single quotes in the URL to prevent expression injection.
func actionExpr(method, url string) string {
	escaped := strings.ReplaceAll(url, "'", "\\'")

	return fmt.Sprintf("@%s('%s')", method, escaped)
}
