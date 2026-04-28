// Package display provides tests for display components.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestDropdownRender(t *testing.T) {
	t.Parallel()

	t.Run("basic dropdown", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "actions"},
			Label:     "Actions",
			Items: []DropdownItem{
				{Text: "Edit", Href: "/edit"},
				{Text: "Delete", Href: "/delete"},
			},
		}))
		utils.AssertContains(t, output, "Actions")
		utils.AssertContains(t, output, `id="actions-button"`)
		utils.AssertContains(t, output, `id="actions-menu"`)
		utils.AssertContains(t, output, "Edit")
		utils.AssertContains(t, output, "Delete")
		utils.AssertContains(t, output, `role="menu"`)
		utils.AssertContains(t, output, `aria-haspopup="true"`)
		utils.AssertContains(t, output, `data-dropdown-trigger="actions"`)
	})

	t.Run("external link item", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "links"},
			Label:     "Links",
			Items: []DropdownItem{
				{Text: "Docs", Href: "https://example.com", External: true},
			},
		}))
		utils.AssertContains(t, output, `target="_blank"`)
		utils.AssertContains(t, output, `rel="noopener noreferrer"`)
	})

	t.Run("button-only item", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "cmds"},
			Label:     "Commands",
			Items: []DropdownItem{
				{Text: "Copy"},
			},
		}))
		utils.AssertContains(t, output, "Copy")
		utils.AssertContains(t, output, `type="button"`)
	})
}
