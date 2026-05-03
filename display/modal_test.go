package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestModalRender(t *testing.T) {
	t.Parallel()
	t.Run("closed modal", func(t *testing.T) {
		t.Parallel()
		props := ModalProps{
			BaseProps: utils.BaseProps{
				ID:     "test-modal",
				Nonce:  "test-nonce",
			},
			Title: "Confirm",
			Open:  false,
			Size:  ModalSizeMD,
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
		t.Parallel()
		props := ModalProps{
			BaseProps: utils.BaseProps{
				ID: "open-modal",
			},
			Title: "Hello",
			Open:  true,
			Size:  ModalSizeMD,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, "opacity-100")
		utils.AssertContains(t, output, "pointer-events-auto")
	})

	t.Run("size variants", func(t *testing.T) {
		sizes := map[ModalSize]string{
			ModalSizeSM:   "max-w-sm",
			ModalSizeMD:   "max-w-md",
			ModalSizeLG:   "max-w-lg",
			ModalSizeXL:   "max-w-xl",
			ModalSizeFull: "max-w-4xl",
		}
		for size, want := range sizes {
			t.Run(string(size), func(t *testing.T) {
				t.Parallel()
				props := ModalProps{
					BaseProps: utils.BaseProps{
						ID: "sz-modal",
					},
					Title: "Test Modal",
					Open:  false,
					Size:  size,
				}
				output := utils.Render(t, Modal(props))
				utils.AssertContains(t, output, want)
			})
		}
	})
}
