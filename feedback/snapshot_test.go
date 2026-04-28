// Package feedback provides tests for feedback components like Alert, Toast, Spinner, and Loading.
package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestAlertRender(t *testing.T) {
	t.Parallel()
	t.Run("error alert", func(t *testing.T) {
		t.Parallel()
		props := AlertProps{
			BaseProps: utils.BaseProps{
				ID:        "",
				Class:     "",
				Attrs:     nil,
				AriaLabel: "",
			},
			Title:       "Error",
			Message:     "Something failed",
			Type:        AlertError,
			Dismissible: false,
		}
		output := utils.Render(t, Alert(props))
		utils.AssertContains(t, output, "Error")
		utils.AssertContains(t, output, "Something failed")
		utils.AssertContains(t, output, "bg-red-50")
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("dismissible alert", func(t *testing.T) {
		t.Parallel()
		props := AlertProps{
			BaseProps: utils.BaseProps{
				ID:        "",
				Class:     "",
				Attrs:     nil,
				AriaLabel: "",
			},
			Title:       "Warning",
			Message:     "",
			Type:        AlertWarning,
			Dismissible: true,
		}
		output := utils.Render(t, Alert(props))
		utils.AssertContains(t, output, `aria-label="Dismiss"`)
		utils.AssertContains(t, output, `data-dismiss="alert"`)
		utils.AssertNotContains(t, output, `onclick=`)
	})
}

func TestToastContainerRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ToastContainer("test-nonce"))
	utils.AssertContains(t, output, "tc-toast-container")
	utils.AssertContains(t, output, `nonce="test-nonce"`)
}

func TestToastRender(t *testing.T) {
	t.Parallel()
	props := ToastProps{
		BaseProps: utils.BaseProps{
			ID:        "",
			Class:     "",
			Attrs:     nil,
			AriaLabel: "",
		},
		Message:  "Saved!",
		Title:    "",
		Type:     ToastSuccess,
		Duration: 0,
	}
	output := utils.Render(t, Toast(props))
	utils.AssertContains(t, output, "Saved!")
	utils.AssertContains(t, output, `role="alert"`)
	utils.AssertContains(t, output, `data-dismiss="toast"`)
	utils.AssertNotContains(t, output, `onclick=`)
}

func TestInlineErrorRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineError("Required"))
	utils.AssertContains(t, output, "Required")
	utils.AssertContains(t, output, "text-red-600")
}

func TestInlineSuccessRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InlineSuccess("Done"))
	utils.AssertContains(t, output, "Done")
	utils.AssertContains(t, output, "text-green-600")
}
