package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestSSEErrorHandlingDefaults(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SSEErrorHandling(DefaultSSEErrorHandlingConfig()))

	utils.AssertContainsAll(t, output,
		"datastar-sse-error",
		"tcShowToast",
		`id="tc-datastar-announcer"`,
		`aria-live="polite"`,
	)
	utils.AssertContains(t, output, "var DURATION = 6000;")
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
