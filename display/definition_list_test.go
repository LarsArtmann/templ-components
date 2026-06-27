package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestDefinitionListRender(t *testing.T) {
	t.Parallel()

	t.Run("text items render as dt/dd pairs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionList(DefinitionListProps{
			Items: []DefinitionItem{
				{Term: "Email", Detail: "alice@example.com"},
				{Term: "Plan", Detail: "Pro"},
			},
		}))
		utils.AssertContains(t, output, "Email")
		utils.AssertContains(t, output, "alice@example.com")
		utils.AssertContains(t, output, "Plan")
		utils.AssertContains(t, output, "Pro")
		utils.AssertContains(t, output, "<dt")
		utils.AssertContains(t, output, "<dd")
	})

	t.Run("detail component overrides text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionList(DefinitionListProps{
			Items: []DefinitionItem{
				{Term: "Status", DetailComponent: templ.Raw(`<span class="badge">Active</span>`)},
			},
		}))
		utils.AssertContains(t, output, "Status")
		utils.AssertContains(t, output, `<span class="badge">Active</span>`)
	})

	t.Run("empty items renders empty dl", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionList(DefinitionListProps{}))
		utils.AssertContains(t, output, "<dl")
		utils.AssertContains(t, output, "</dl>")
		utils.AssertNotContains(t, output, "<dt")
	})

	t.Run("propagates BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionList(DefinitionListProps{
			Items: []DefinitionItem{{Term: "Key", Detail: "Val"}},
			BaseProps: utils.BaseProps{
				ID:        "deflist",
				Class:     "extra-class",
				AriaLabel: "Details",
			},
		}))
		utils.AssertContains(t, output, `id="deflist"`)
		utils.AssertContains(t, output, "extra-class")
		utils.AssertContains(t, output, `aria-label="Details"`)
	})

	t.Run("grid layout classes present", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionList(DefinitionListProps{
			Items: []DefinitionItem{{Term: "X", Detail: "Y"}},
		}))
		utils.AssertContains(t, output, "grid")
		utils.AssertContains(t, output, "gap-x-6")
	})
}
