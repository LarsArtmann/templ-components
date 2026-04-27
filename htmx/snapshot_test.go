package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestLoadingIndicatorRender(t *testing.T) {
	output := utils.Render(t, LoadingIndicator())
	utils.AssertContains(t, output, "tc-loading-indicator")
	utils.AssertContains(t, output, "animate-spin")
}

func TestGlobalErrorHandlingRender(t *testing.T) {
	output := utils.Render(t, GlobalErrorHandling("test-nonce"))
	utils.AssertContains(t, output, `nonce="test-nonce"`)
	utils.AssertContains(t, output, "htmx:responseError")
	utils.AssertContains(t, output, "htmx:sendError")
}
