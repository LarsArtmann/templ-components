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
				{
					ID:      accordionIDFAQ1,
					Title:   "What is this?",
					Open:    true,
					Content: templ.Raw("\u003cp\u003eAnswer 1\u003c/p\u003e"),
				},
				{
					ID:      "faq2",
					Title:   "How does it work?",
					Content: templ.Raw("\u003cp\u003eAnswer 2\u003c/p\u003e"),
				},
			},
		}))

		utils.AssertContains(t, output, "What is this?")
		utils.AssertContains(t, output, "How does it work?")
		utils.AssertContains(t, output, "Answer 1")
		utils.AssertContains(t, output, "<details")
		utils.AssertContains(t, output, "<summary")
		utils.AssertContains(t, output, " open")
		utils.AssertContains(t, output, `id="`+accordionIDFAQ1+`"`)
		utils.AssertContains(t, output, `data-tc-chevron`)
	})

	t.Run("closed accordion item", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{ID: "item1", Title: "Closed", Open: false, Content: templ.Raw("Content")},
			},
		}))
		utils.AssertContains(t, output, "<details")
		utils.AssertContains(t, output, "<summary")
		utils.AssertNotContains(t, output, "<details open")
	})

	t.Run("no script tag - uses native details", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			BaseProps: utils.BaseProps{Nonce: testNonce},
			Items: []AccordionItem{
				{ID: "n1", Title: "Q", Content: templ.Raw("A")},
			},
		}))
		utils.AssertNotContains(t, output, "<script")
		utils.AssertNotContains(t, output, "tcAccordionAttached")
	})

	t.Run("custom class and id", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			BaseProps: utils.BaseProps{ID: "faq", Class: cssClassMt4},
			Items: []AccordionItem{
				{ID: "n1", Title: "Q", Content: templ.Raw("A")},
			},
		}))
		utils.AssertContains(t, output, `id="faq"`)
		utils.AssertContains(t, output, cssClassMt4)
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()

		props := DefaultAccordionProps()
		if props.Items != nil {
			t.Error("DefaultAccordionProps().Items should be nil")
		}
	})
}
