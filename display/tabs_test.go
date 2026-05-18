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
			ActiveTabID: "tab1",
			Tabs: []Tab{
				{ID: "tab1", Label: "First"},
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
			ActiveTabID: "panel1",
			Tabs: []Tab{
				{
					ID:      "panel1",
					Label:   "Panel",
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
			ActiveTabID: "a",
			Variant:     TabsPills,
			Tabs: []Tab{
				{ID: "a", Label: "A"},
			},
		}))
		utils.AssertContains(t, output, "bg-blue-600")
	})

	t.Run("no active tab renders all inactive", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			Tabs: []Tab{
				{ID: "x", Label: "X"},
				{ID: "y", Label: "Y"},
			},
		}))
		utils.AssertNotContains(t, output, `aria-selected="true"`)
	})

	t.Run("default props has TabsDefault variant", func(t *testing.T) {
		t.Parallel()
		props := DefaultTabsProps()
		if props.Variant != TabsDefault {
			t.Errorf("DefaultTabsProps().Variant = %q, want %q", props.Variant, TabsDefault)
		}
	})

	t.Run("tabs with empty tabs list", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{}))
		utils.AssertContains(t, output, `role="tablist"`)
	})
}
