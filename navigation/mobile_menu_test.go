package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestMobileMenuKeyboardNavigation(t *testing.T) {
	t.Parallel()

	links := []NavLinkProps{
		{Text: "Home", Href: "/"},
		{Text: "About", Href: "/about"},
	}
	output := utils.Render(t, MobileMenu(links, "/", "n", "mobile-kbd", false))

	// Escape closes the menu when focus is inside it.
	utils.AssertContains(t, output, "e.key !== 'Escape'")

	// Opening moves focus to the first focusable menu item.
	utils.AssertContains(t, output, "first = menu.querySelector")
	utils.AssertContains(t, output, "first.focus()")

	// Closing returns focus to the toggle button.
	utils.AssertContains(t, output, "btn.focus()")
}

// TestMobileMenuScriptMarkupCoRender pins the script/markup co-rendering
// invariant (#365 decision): the singleton script and the menu markup render
// under the SAME condition (links exist). A script without markup would be
// dead weight; markup without the script is the dead-hamburger trap — a
// toggle whose clicks do nothing (the #402/#403 fleet class). The Nav
// composition gates MobileMenuToggle on the SAME len(links) condition, so
// the composed path can never ship a toggle without its script.
func TestMobileMenuScriptMarkupCoRender(t *testing.T) {
	t.Parallel()

	links := []NavLinkProps{{Text: "Home", Href: "/"}}

	t.Run("with links: menu and script co-render", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileMenu(links, "/", "", "mm-x", false))
		utils.AssertContains(t, output, `data-mobile-menu-id="mm-x"`)
		utils.AssertContains(t, output, "tcMobileMenuAttached")
	})

	t.Run("without links: neither menu nor script renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileMenu(nil, "/", "", "mm-y", false))
		utils.AssertNotContains(t, output, `data-mobile-menu-id="mm-y"`)
		utils.AssertNotContains(t, output, "tcMobileMenuAttached")
		utils.AssertNotContains(t, output, "<script")
	})

	t.Run("Nav with no links renders no dead toggle", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{Links: nil, CurrentPath: "/"}))
		utils.AssertNotContains(t, output, "data-mobile-menu-toggle")
	})
}
