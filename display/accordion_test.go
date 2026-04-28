// Package display provides tests for display components.
package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestAccordionRender(t *testing.T) {
	t.Parallel()

	t.Run("basic accordion", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{ID: "faq1", Title: "What is this?", Open: true, Content: templ.Raw("\u003cp\u003eAnswer 1\u003c/p\u003e")},
				{ID: "faq2", Title: "How does it work?", Content: templ.Raw("\u003cp\u003eAnswer 2\u003c/p\u003e")},
			},
		}))
		utils.AssertContains(t, output, "What is this?")
		utils.AssertContains(t, output, "How does it work?")
		utils.AssertContains(t, output, "Answer 1")
		utils.AssertContains(t, output, `role="region"`)
		utils.AssertContains(t, output, `aria-expanded="true"`)
		utils.AssertContains(t, output, `aria-expanded="false"`)
		utils.AssertContains(t, output, `data-accordion-trigger="faq1"`)
	})

	t.Run("closed accordion item", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{ID: "item1", Title: "Closed", Open: false, Content: templ.Raw("Content")},
			},
		}))
		utils.AssertContains(t, output, `class="overflow-hidden transition-all duration-200 max-h-0"`)
	})
}
