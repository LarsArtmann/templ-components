package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestAlertRender(t *testing.T) {
	t.Run("error alert", func(t *testing.T) {
		props := AlertProps{Title: "Error", Message: "Something failed", Type: AlertError}
		output := utils.Render(t, Alert(props))
		utils.AssertContains(t, output, "Error")
		utils.AssertContains(t, output, "Something failed")
		utils.AssertContains(t, output, "bg-red-50")
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("dismissible alert", func(t *testing.T) {
		props := AlertProps{Title: "Warning", Type: AlertWarning, Dismissible: true}
		output := utils.Render(t, Alert(props))
		utils.AssertContains(t, output, `aria-label="Dismiss"`)
	})
}

func TestToastContainerRender(t *testing.T) {
	output := utils.Render(t, ToastContainer("test-nonce"))
	utils.AssertContains(t, output, "tc-toast-container")
	utils.AssertContains(t, output, `nonce="test-nonce"`)
}

func TestToastRender(t *testing.T) {
	props := ToastProps{Message: "Saved!", Type: ToastSuccess}
	output := utils.Render(t, Toast(props))
	utils.AssertContains(t, output, "Saved!")
	utils.AssertContains(t, output, `role="alert"`)
}

func TestInlineErrorRender(t *testing.T) {
	output := utils.Render(t, InlineError("Required"))
	utils.AssertContains(t, output, "Required")
	utils.AssertContains(t, output, "text-red-600")
}

func TestInlineSuccessRender(t *testing.T) {
	output := utils.Render(t, InlineSuccess("Done"))
	utils.AssertContains(t, output, "Done")
	utils.AssertContains(t, output, "text-green-600")
}
