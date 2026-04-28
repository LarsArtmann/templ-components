// Package display provides tests for display components.
package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestTabsRender(t *testing.T) {
	t.Parallel()
	t.Run("basic tabs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			Tabs: []Tab{
				{ID: "tab1", Label: "First", Active: true},
				{ID: "tab2", Label: "Second"},
			},
		}))
		utils.AssertContains(t, output, "First")
		utils.AssertContains(t, output, "Second")
		utils.AssertContains(t, output, `role="tablist"`)
		utils.AssertContains(t, output, `aria-selected="true"`)
		utils.AssertContains(t, output, `aria-selected="false"`)
	})

	t.Run("tab with content", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			Tabs: []Tab{
				{
					ID:      "panel1",
					Label:   "Panel",
					Active:  true,
					Content: templ.Raw("<p>Hello</p>"),
				},
			},
		}))
		utils.AssertContains(t, output, `role="tabpanel"`)
		utils.AssertContains(t, output, "Hello")
	})

	t.Run("pills style", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			TabStyle: TabsStylePills,
			Tabs: []Tab{
				{ID: "a", Label: "A", Active: true},
			},
		}))
		utils.AssertContains(t, output, "bg-blue-600")
	})
}
