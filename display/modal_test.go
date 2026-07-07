package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const modalTestNonce = "test-nonce"

func TestModalRender(t *testing.T) {
	t.Parallel()
	t.Run("closed modal", func(t *testing.T) {
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
		utils.AssertContains(t, output, `id="test-modal"`)
		utils.AssertContains(t, output, "Confirm")
		utils.AssertContains(t, output, `role="dialog"`)
		utils.AssertContains(t, output, `aria-modal="true"`)
		utils.AssertContains(t, output, `nonce="`+modalTestNonce+`"`)
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
			ModalSize2XL:  "max-w-4xl",
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

	t.Run("focus restore JS saves activeElement", func(t *testing.T) {
		t.Parallel()
		props := ModalProps{
			BaseProps: utils.BaseProps{
				ID:    "focus-modal",
				Nonce: modalTestNonce,
			},
			Title: "Focus Test",
			Open:  false,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, "data-tc-prev-focus")
		utils.AssertContains(t, output, "document.activeElement")
		utils.AssertContains(t, output, "removeAttribute")
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

	t.Run("closed overlay has aria-hidden and inert", func(t *testing.T) {
		t.Parallel()
		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "closed-modal"},
			Title:     "Closed",
			Open:      false,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, `aria-hidden="true"`)
		utils.AssertContains(t, output, " inert")
	})

	t.Run("open overlay has aria-hidden=false, no inert, open-on-load hook", func(t *testing.T) {
		t.Parallel()
		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "open-modal-a11y"},
			Title:     "Open",
			Open:      true,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, `aria-hidden="false"`)
		utils.AssertNotContains(t, output, " inert")
		utils.AssertContains(t, output, `data-tc-open-on-load="true"`)
	})

	t.Run("overlay JS syncs aria-hidden and inert on open/close", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Modal(ModalProps{
			BaseProps: utils.BaseProps{ID: "js-sync-modal"},
			Title:     "Sync",
			Open:      false,
		}))
		// Open function must remove inert and set aria-hidden=false so the
		// dialog is focusable and in the a11y tree when opened via JS.
		utils.AssertContains(t, output, "setAttribute('aria-hidden', 'false')")
		utils.AssertContains(t, output, "removeAttribute('inert')")
		// Close function must add inert and set aria-hidden=true.
		utils.AssertContains(t, output, "setAttribute('aria-hidden', 'true')")
		utils.AssertContains(t, output, "setAttribute('inert', '')")
	})
}
