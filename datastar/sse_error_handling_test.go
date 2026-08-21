package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestSSEErrorHandlingDefaults(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SSEErrorHandling(DefaultSSEErrorHandlingConfig()))

	utils.AssertContainsAll(t, output,
		"datastar-fetch",
		"tcShowToast",
		`id="tc-datastar-announcer"`,
		`aria-live="polite"`,
	)
	utils.AssertContains(t, output, "var DURATION = 6000;")
}

// TestSSEErrorHandlingReportsTerminalStates pins the two failure types the
// DataStar v1.x runtime actually dispatches on the datastar-fetch event.
// "error" fires on HTTP-level failures (status >= 400), "retries-failed"
// when automatic reconnection gave up. There is no "datastar-sse-error"
// event — guarding against that regression here.
func TestSSEErrorHandlingReportsTerminalStates(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SSEErrorHandling(DefaultSSEErrorHandlingConfig()))

	utils.AssertContainsAll(t, output,
		"detail.type === 'error'",
		"detail.type === 'retries-failed'",
		"detail.type === 'retrying'",
		"args.status",
	)
	utils.AssertNotContains(t, output, "datastar-sse-error")
}

func TestSSEErrorHandlingCustomDuration(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SSEErrorHandling(SSEErrorHandlingConfig{DurationMS: 9000}))

	utils.AssertContains(t, output, "var DURATION = 9000;")
}

func TestSSEErrorHandlingZeroDurationNormalizes(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SSEErrorHandling(SSEErrorHandlingConfig{}))

	utils.AssertContains(t, output, "var DURATION = 6000;")
}

func TestSSEErrorHandlingNonce(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SSEErrorHandling(SSEErrorHandlingConfig{Nonce: "abc123"}))

	utils.AssertContains(t, output, `nonce="abc123"`)
}

func TestSSEErrorHandlingAnnouncesToLiveRegion(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SSEErrorHandling(DefaultSSEErrorHandlingConfig()))

	utils.AssertContains(t, output, "tc-datastar-announcer")
}
