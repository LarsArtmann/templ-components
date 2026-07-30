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
