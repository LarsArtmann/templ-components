// Package display provides tests for display components.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

const dropdownItemDelete = "Delete"

func TestDropdownRender(t *testing.T) {
	t.Parallel()

	t.Run("basic dropdown", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "actions"},
			Label:     dropdownLabelActions,
			Items: []DropdownItem{
				{Text: dropdownItemEdit, Href: dropdownHrefEdit},
				{Text: dropdownItemDelete, Href: "/delete"},
			},
		}))
		utils.AssertContains(t, output, dropdownLabelActions)
		utils.AssertContains(t, output, `id="actions-button"`)
		utils.AssertContains(t, output, `id="actions-menu"`)
		utils.AssertContains(t, output, dropdownItemEdit)
		utils.AssertContains(t, output, dropdownItemDelete)
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

	t.Run("empty ID auto-generates on render", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{Label: "No ID"}))
		utils.AssertContains(t, output, `id="tc-dropdown-`)
	})

	t.Run("right position", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "rmenu"},
			Label:     dropdownLabelMenu,
			Position:  DropdownPositionRight,
			Items:     []DropdownItem{{Text: "Item", Href: "/x"}},
		}))
		utils.AssertContains(t, output, `data-tc-align="end"`)
	})

	t.Run("with icon items", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "imenu"},
			Label:     "Actions",
			Items: []DropdownItem{
				{Text: dropdownItemEdit, Href: "/edit", Icon: icons.Edit},
				{Text: dropdownItemDelete, Href: "/del", Icon: icons.Trash},
			},
		}))
		utils.AssertContains(t, output, "Edit")
		utils.AssertContains(t, output, "Delete")
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()

		props := DefaultDropdownProps()
		if props.Label != "" {
			t.Errorf("DefaultDropdownProps().Label = %q, want empty", props.Label)
		}
	})

	t.Run("with divider items", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "dmenu"},
			Label:     dropdownLabelMenu,
			Items: []DropdownItem{
				{Text: dropdownItemEdit, Href: "/edit"},
				{
					Text:  dropdownItemDelete,
					Href:  "/del",
					Attrs: map[string]any{"class": "text-red-600"},
				},
			},
		}))
		utils.AssertContains(t, output, dropdownItemEdit)
		utils.AssertContains(t, output, dropdownItemDelete)
		utils.AssertContains(t, output, "text-red-600")
	})

	t.Run("RTL keyboard nav computes keys correctly", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "rtl-menu"},
			Label:     dropdownLabelMenu,
			Items:     []DropdownItem{{Text: "Item", Href: "/x"}},
		}))
		// The RTL ternary must NOT be inside a JS string literal (dead code).
		// It should be computed as a variable assignment.
		utils.AssertContains(t, output, "var nextKey")
		utils.AssertNotContains(t, output, "=== '(document.documentElement")
	})
}

func TestDropdownKeyboardEnhancements(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, Dropdown(DropdownProps{
		BaseProps: utils.BaseProps{ID: "kbd-menu", Nonce: "n"},
		Label:     dropdownLabelMenu,
		Items: []DropdownItem{
			{Text: "One", Href: "/1"},
			{Text: "Two", Href: "/2"},
			{Text: "Three", Href: "/3"},
			{Text: "Four", Href: "/4"},
		},
	}))

	// Home/End and PageUp/PageDown move focus inside the menu.
	utils.AssertContains(t, output, "e.key === 'Home'")
	utils.AssertContains(t, output, "e.key === 'End'")
	utils.AssertContains(t, output, "e.key === 'PageDown'")
	utils.AssertContains(t, output, "e.key === 'PageUp'")

	// Disabled items are skipped during keyboard navigation.
	utils.AssertContains(t, output, `querySelectorAll('[role="menuitem"]:not([disabled])')`)

	// Enter/Space activation is intentionally left to native browser behavior
	// so that HTMX-powered links and buttons work correctly.
	utils.AssertNotContains(t, output, "window.location.href")
}
