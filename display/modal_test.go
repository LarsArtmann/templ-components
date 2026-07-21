package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const modalTestNonce = "test-nonce"

func TestModalRender(t *testing.T) {
	t.Parallel()
	t.Run("closed modal renders dialog element", func(t *testing.T) {
		t.Parallel()

		props := ModalProps{
			BaseProps: utils.BaseProps{
				ID:    "test-modal",
				Nonce: modalTestNonce,
			},
			Title: "Confirm",
			Open:  false,
			Size:  ModalSizeMD,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, `<dialog`)
		utils.AssertContains(t, output, `id="test-modal"`)
		utils.AssertContains(t, output, "Confirm")
		utils.AssertContains(t, output, `nonce="`+modalTestNonce+`"`)
		utils.AssertContains(t, output, "tcCloseModal")
		utils.AssertContains(t, output, "tcOpenModal")
	})

	t.Run("open modal has data-tc-open attribute", func(t *testing.T) {
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
		utils.AssertContains(t, output, `data-tc-open="true"`)
	})

	t.Run("size variants", func(t *testing.T) {
		sizes := map[ModalSize]string{
			ModalSizeSM:  "max-w-sm",
			ModalSizeMD:  "max-w-md",
			ModalSizeLG:  "max-w-lg",
			ModalSizeXL:  "max-w-xl",
			ModalSize2XL: "max-w-4xl",
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

	t.Run("special chars in ID are JS-escaped", func(t *testing.T) {
		t.Parallel()

		props := ModalProps{
			BaseProps: utils.BaseProps{
				ID:    "modal-with-'quotes'",
				Nonce: modalTestNonce,
			},
			Title: "Escape Test",
			Open:  false,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertNotContains(t, output, "('modal-with-'quotes'')")
		utils.AssertContains(t, output, `"modal-with-'quotes'"`)
	})

	t.Run("JS uses native showModal API", func(t *testing.T) {
		t.Parallel()

		props := ModalProps{
			BaseProps: utils.BaseProps{
				ID:    "dialog-modal",
				Nonce: modalTestNonce,
			},
			Title: "Dialog Test",
			Open:  false,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, "showModal")
		utils.AssertContains(t, output, "d.close()")
		utils.AssertContains(t, output, "data-tc-close")
	})

	t.Run("empty ID auto-generates on render", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Modal(ModalProps{Title: "No ID"}))
		utils.AssertContains(t, output, `id="tc-modal-`)
	})

	t.Run("no title omits aria-labelledby (no dangling ref)", func(t *testing.T) {
		t.Parallel()

		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "notitle-modal"},
			Title:     "",
			Open:      false,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertNotContains(t, output, "aria-labelledby")
	})

	t.Run("with title emits aria-labelledby", func(t *testing.T) {
		t.Parallel()

		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "titled-modal"},
			Title:     "Titled",
			Open:      false,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, `aria-labelledby="titled-modal-title"`)
	})

	t.Run("closed dialog has no data-tc-open", func(t *testing.T) {
		t.Parallel()

		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "closed-modal"},
			Title:     "Closed",
			Open:      false,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertNotContains(t, output, `data-tc-open="true"`)
	})

	t.Run("open dialog has data-tc-open for auto-open", func(t *testing.T) {
		t.Parallel()

		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "open-modal-a11y"},
			Title:     "Open",
			Open:      true,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, `data-tc-open="true"`)
	})

	t.Run("JS singleton guard prevents double-binding", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Modal(ModalProps{
			BaseProps: utils.BaseProps{ID: "js-guard-modal"},
			Title:     "Guard",
			Open:      false,
		}))
		utils.AssertContains(t, output, "window.tcOverlayModalAttached")
	})

	t.Run("JS backdrop click handler detects dialog target", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Modal(ModalProps{
			BaseProps: utils.BaseProps{ID: "backdrop-modal"},
			Title:     "Backdrop",
			Open:      false,
		}))
		utils.AssertContains(t, output, "e.target===d")
	})
}
