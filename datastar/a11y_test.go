package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestSSEErrorHandlingA11y is the accessibility lens for the global SSE error
// handler: every failure path must reach a non-sighted user even when the
// toast container is missing.
func TestSSEErrorHandlingA11y(t *testing.T) {
	t.Parallel()

	t.Run("announcer region exists and is visually hidden", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SSEErrorHandling(DefaultSSEErrorHandlingConfig()))

		utils.AssertContains(t, output, `<div id="tc-datastar-announcer" aria-live="polite" class="sr-only"></div>`)
	})

	t.Run("announcements feed the aria-live region via textContent", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SSEErrorHandling(DefaultSSEErrorHandlingConfig()))

		// Screen readers announce textContent changes on aria-live regions.
		// The handler must not use innerHTML (injection) or create/remove
		// nodes (some screen readers miss region creation).
		utils.AssertContains(t, output, "announcer.textContent = text")
		utils.AssertNotContains(t, output, "announcer.innerHTML")
	})

	t.Run("announcer is polite, not assertive", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SSEErrorHandling(DefaultSSEErrorHandlingConfig()))

		// Stream errors are important but not emergencies: polite queueing
		// avoids interrupting the user mid-sentence. The region is declared
		// once at render time (a11y best practice).
		utils.AssertContains(t, output, `aria-live="polite"`)
	})

	t.Run("both failure paths announce through the same region", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SSEErrorHandling(DefaultSSEErrorHandlingConfig()))

		utils.AssertContains(t, output, `announce('Stream error' + status)`)
		utils.AssertContains(t, output, "announce('Live stream lost. Automatic reconnection failed.')")
	})
}

// TestSSEErrorHandlingCSP pins the security lens: the inline script must be
// nonced (the integration CSP test covers the cross-package invariant; this
// keeps the package-local guarantee explicit).
func TestSSEErrorHandlingCSP(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SSEErrorHandling(SSEErrorHandlingConfig{Nonce: "a11y-nonce"}))

	utils.AssertContains(t, output, `<script nonce="a11y-nonce">`)
}
