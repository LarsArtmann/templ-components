package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestModalEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("modal without title renders no header", func(t *testing.T) {
		t.Parallel()
		props := ModalProps{
			BaseProps: utils.BaseProps{ID: "no-title-modal"},
			Open:      true,
			Size:      ModalSizeMD,
		}
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, `id="no-title-modal"`)
		utils.AssertNotContains(t, output, `id="no-title-modal-title"`)
		utils.AssertNotContains(t, output, "aria-label=\"Close\"")
	})

	t.Run("dropdown with empty items list renders button only", func(t *testing.T) {
		t.Parallel()
		props := DropdownProps{
			BaseProps: utils.BaseProps{ID: "empty-dd"},
			Items:     []DropdownItem{},
		}
		output := utils.Render(t, Dropdown(props))
		utils.AssertContains(t, output, `id="empty-dd"`)
		utils.AssertContains(t, output, `id="empty-dd-button"`)
	})

	t.Run("dropdown item with both Href and action renders link", func(t *testing.T) {
		t.Parallel()
		props := DropdownProps{
			BaseProps: utils.BaseProps{ID: "both-dd"},
			Items: []DropdownItem{
				{Text: "Link", Href: "/link"},
			},
		}
		output := utils.Render(t, Dropdown(props))
		utils.AssertContains(t, output, "/link")
	})
}

func TestAccordionEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty item ID panics", func(t *testing.T) {
		t.Parallel()
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic for Accordion item with empty ID")
			}
		}()
		props := AccordionProps{
			Items: []AccordionItem{
				{ID: "", Title: "Missing"},
			},
		}
		utils.Render(t, Accordion(props))
	})

	t.Run("empty items list renders container only", func(t *testing.T) {
		t.Parallel()
		props := AccordionProps{
			Items: []AccordionItem{},
		}
		output := utils.Render(t, Accordion(props))
		utils.AssertContains(t, output, "divide-y")
	})
}
