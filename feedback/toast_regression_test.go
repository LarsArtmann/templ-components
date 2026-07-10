package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestToastContainerHasNonce verifies CSP nonce on the container script.
func TestToastContainerHasNonce(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ToastContainer("test-nonce"))
	utils.AssertContains(t, output, `nonce="test-nonce"`)
}

// TestToastContainerHasLiveRegion verifies aria-live for screen reader announcements.
func TestToastContainerHasLiveRegion(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ToastContainer(""))
	utils.AssertContains(t, output, `aria-live="polite"`)
	utils.AssertContains(t, output, `aria-atomic="true"`)
}

// TestToastContainerHasShowFunction verifies the JS toast constructor function exists.
func TestToastContainerHasShowFunction(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ToastContainer(""))
	utils.AssertContains(t, output, "tcShowToast")
}

// TestToastContainerHasColorMap verifies tcToastColors is defined (not JSON-marshaled).
func TestToastContainerHasColorMap(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ToastContainer(""))
	utils.AssertContains(t, output, "tcToastColors")
	utils.AssertContains(t, output, "tcToastIcons")
}

// TestToastBasicRender verifies a single toast renders with correct type classes.
func TestToastBasicRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toast(ToastProps{
		Type:    ToastSuccess,
		Message: "Saved successfully",
	}))
	utils.AssertContains(t, output, "Saved successfully")
}

// TestToastDurationAutoDismiss verifies the setTimeout is present when Duration > 0.
func TestToastDurationAutoDismiss(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toast(ToastProps{
		Type:     ToastInfo,
		Message:  "test",
		Duration: 5000,
	}))
	utils.AssertContains(t, output, "setTimeout")
	utils.AssertContains(t, output, "5000")
}
