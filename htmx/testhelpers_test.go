package htmx

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// testSpinner renders a minimal spinner SVG for tests. Avoids importing the
// feedback package (root module) so the htmx sub-module stays standalone.
// The colorClasses parameter is included in the SVG class attribute so tests
// that assert on color classes (dark:, htmx-indicator) continue to pass.
func testSpinner(colorClasses string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		svg := `<svg class="animate-spin h-5 w-5 ` + colorClasses +
			`" viewBox="0 0 24 24" fill="none" aria-hidden="true"></svg>`

		_, err := io.WriteString(w, svg)

		return err //nolint:wrapcheck // test helper
	})
}
