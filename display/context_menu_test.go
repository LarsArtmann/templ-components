package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestContextMenuRender(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, ContextMenu(ContextMenuProps{
		BaseProps: utils.BaseProps{Nonce: "n"},
		Items: []ContextMenuItem{
			{Text: "Edit", Href: "/edit"},
			{Text: "Delete", Href: "/delete"},
		},
	}))

	utils.AssertContains(t, output, `data-tc-ctxmenu-trigger`)
	utils.AssertContains(t, output, `role="menu"`)
	utils.AssertContains(t, output, `role="menuitem"`)
	utils.AssertContains(t, output, `popover="auto"`)
}

func TestContextMenuKeyboardTrigger(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, ContextMenu(ContextMenuProps{
		BaseProps: utils.BaseProps{Nonce: "n"},
		Items: []ContextMenuItem{
			{Text: "Edit", Href: "/edit"},
			{Text: "Delete", Href: "/delete"},
		},
	}))

	// Keyboard users cannot right-click: Shift+F10 and the dedicated ContextMenu
	// key must open the menu and position it at the trigger element.
	utils.AssertContains(t, output, "e.key === 'F10'")
	utils.AssertContains(t, output, "e.shiftKey")
	utils.AssertContains(t, output, "e.key === 'ContextMenu'")
	utils.AssertContains(t, output, "getBoundingClientRect()")

	// Menuitems use roving tabindex so arrow navigation can focus them.
	utils.AssertContains(t, output, `tabindex="-1"`)
}

func TestContextMenuSharedMenuNav(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, ContextMenu(ContextMenuProps{
		BaseProps: utils.BaseProps{Nonce: "n"},
		Items: []ContextMenuItem{
			{Text: "Edit", Href: "/edit"},
			{Text: "Delete", Href: "/delete"},
		},
	}))

	// The shared menu keyboard-nav helper (Arrow/Home/End + focus-first on open)
	// is injected by the context menu, the same one Dropdown uses.
	utils.AssertContains(t, output, "tcMenuKeyNavAttached")
}

func TestContextMenuDisabledItem(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, ContextMenu(ContextMenuProps{
		BaseProps: utils.BaseProps{Nonce: "n"},
		Items: []ContextMenuItem{
			{Text: "Edit", Href: "/edit"},
			{Text: "Delete", Disabled: true},
		},
	}))

	// Disabled items expose aria-disabled so the shared nav selector skips them.
	utils.AssertContains(t, output, `aria-disabled="true"`)
	// Class order is not asserted (tailwind-merge reorders); check tokens instead.
	utils.AssertContainsAll(t, output, "opacity-50", "pointer-events-none")
}
