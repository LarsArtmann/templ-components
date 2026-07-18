package htmx

import (
	"context"
	"fmt"
	"html"
	"io"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// ViewTransitionsProps configures native View Transitions API integration
// for HTMX content swaps.
//
// The View Transitions API provides smooth, animated transitions between
// DOM states — cross-fades, shared-element morphing, and custom animations
// — with zero JavaScript beyond enabling the feature.
//
// HTMX 2.0 supports View Transitions natively. This component:
//   - Optionally enables View Transitions globally for all HTMX swaps
//   - Renders CSS for default and customizable transition animations
//
// Graceful degradation: browsers without View Transitions support (older
// Firefox/Safari) perform instant swaps — identical to current behavior.
type ViewTransitionsProps struct {
	utils.BaseProps

	// Global enables View Transitions for ALL HTMX swaps on the page.
	// When false, consumers opt in per-element via hx-swap="... transition:true".
	Global bool
}

// DefaultViewTransitionsProps returns sensible defaults.
func DefaultViewTransitionsProps() ViewTransitionsProps {
	return ViewTransitionsProps{ //nolint:exhaustruct // intentionally minimal defaults
		Global: true,
	}
}

// viewTransitionsCSS provides a smooth cross-fade by default and respects
// prefers-reduced-motion. Consumers can override these pseudo-element
// styles in their own CSS for custom transition animations.
const viewTransitionsCSS = `::view-transition-old(root){animation:tc-vt-fade-out 200ms ease forwards}
::view-transition-new(root){animation:tc-vt-fade-in 200ms ease forwards}
@keyframes tc-vt-fade-out{to{opacity:0}}
@keyframes tc-vt-fade-in{from{opacity:0}}
@media (prefers-reduced-motion:reduce){::view-transition-old(root),::view-transition-new(root){animation:none}}`

// viewTransitionsScript enables global View Transitions for HTMX 2.0.
const viewTransitionsScript = `(function(){if(typeof htmx!=='undefined'&&document.startViewTransition){htmx.config.globalViewTransitions=true;}})();`

// styleComponent renders a CSP-safe <style nonce="..."> tag wrapping the
// given CSS string. Uses the same pattern as display.scriptComponent to
// bypass templ's raw-text element handling.
func styleComponent(nonce, css string) templ.Component {
	escapedNonce := html.EscapeString(nonce)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, "<style nonce=\"%s\">\n%s\n</style>\n", escapedNonce, css); err != nil {
			return fmt.Errorf("write view transitions style: %w", err)
		}

		return nil
	})
}

// scriptComponent renders a CSP-safe <script nonce="..."> tag.
func scriptComponent(nonce, js string) templ.Component {
	escapedNonce := html.EscapeString(nonce)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, "<script nonce=\"%s\">\n%s\n</script>\n", escapedNonce, js); err != nil {
			return fmt.Errorf("write view transitions script: %w", err)
		}

		return nil
	})
}
