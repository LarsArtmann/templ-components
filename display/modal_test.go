package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestModalRender(t *testing.T) {
	t.Run("closed modal", func(t *testing.T) {
		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "test-modal"},
			Title:     "Confirm",
			Size:      "md",
			Nonce:     "test-nonce",
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, `id="test-modal"`)
		utils.AssertContains(t, output, "Confirm")
		utils.AssertContains(t, output, `role="dialog"`)
		utils.AssertContains(t, output, `aria-modal="true"`)
		utils.AssertContains(t, output, `nonce="test-nonce"`)
		utils.AssertContains(t, output, "tcCloseModal")
		utils.AssertContains(t, output, "opacity-0")
	})

	t.Run("open modal", func(t *testing.T) {
		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "open-modal"},
			Title:     "Hello",
			Open:      true,
			Nonce:     "n",
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, "opacity-100")
		utils.AssertContains(t, output, "pointer-events-auto")
	})

	t.Run("size variants", func(t *testing.T) {
		sizes := map[string]string{
			"sm":   "max-w-sm",
			"md":   "max-w-md",
			"lg":   "max-w-lg",
			"xl":   "max-w-xl",
			"full": "max-w-4xl",
		}
		for size, want := range sizes {
			t.Run(size, func(t *testing.T) {
				props := ModalProps{
					BaseProps: utils.BaseProps{ID: "sz-modal"},
					Size:      size,
					Nonce:     "n",
				}
				output := utils.Render(t, Modal(props))
				utils.AssertContains(t, output, want)
			})
		}
	})
}
