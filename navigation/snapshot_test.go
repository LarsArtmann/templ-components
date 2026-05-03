// Package navigation provides rendering tests for navigation components.
package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestNavLinkRender(t *testing.T) {
	t.Parallel()
	t.Run("active link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NavLink(NavLinkProps{Href: "/", Text: "Home"}, "/"))
		utils.AssertContains(t, output, "Home")
		utils.AssertContains(t, output, `href="/"`)
		utils.AssertContains(t, output, `aria-current="page"`)
	})

	t.Run("inactive link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NavLink(NavLinkProps{Href: "/about", Text: "About"}, "/"))
		utils.AssertContains(t, output, "About")
		utils.AssertContains(t, output, `href="/about"`)
		utils.AssertNotContains(t, output, `aria-current="page"`)
	})

	t.Run("external link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			NavLink(
				NavLinkProps{Href: "https://example.com", Text: "External", External: true},
				"/",
			),
		)
		utils.AssertContains(t, output, `target="_blank"`)
		utils.AssertContains(t, output, `rel="noopener noreferrer"`)
	})
}

func TestMobileNavLinkRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileNavLink(NavLinkProps{Href: "/", Text: "Home"}, "/"))
	utils.AssertContains(t, output, "Home")
	utils.AssertContains(t, output, `aria-current="page"`)
}

func TestBreadcrumbsRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs([]BreadcrumbItem{
		{Text: "Home", Href: "/"},
		{Text: "Users", Active: true},
	}))
	utils.AssertContains(t, output, "Home")
	utils.AssertContains(t, output, "Users")
	utils.AssertContains(t, output, `href="/"`)
}

func TestFooterRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Footer("MyApp"))
	utils.AssertContains(t, output, "MyApp")
	utils.AssertContains(t, output, "All rights reserved")
}

func TestMobileMenuToggleRender(t *testing.T) {
	t.Parallel()
	t.Run("shown", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileMenuToggle(true))
		utils.AssertContains(t, output, `data-mobile-menu-toggle`)
		utils.AssertContains(t, output, `aria-controls="tc-mobile-menu"`)
		utils.AssertContains(t, output, `aria-expanded="false"`)
	})
	t.Run("hidden", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileMenuToggle(false))
		utils.AssertNotContains(t, output, "button")
	})
}

var testNavLinks = []NavLinkProps{{Href: "/", Text: "Home"}, {Href: "/about", Text: "About"}}

func TestMobileMenuRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileMenu(testNavLinks, "/", "test-nonce"))
	utils.AssertContains(t, output, "Home")
	utils.AssertContains(t, output, "About")
	utils.AssertContains(t, output, `id="tc-mobile-menu"`)
	utils.AssertContains(t, output, `nonce="test-nonce"`)
}

func TestSimpleNavRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleNav("MyApp", "/", testNavLinks, "/"))
	utils.AssertContains(t, output, "MyApp")
	utils.AssertContains(t, output, "Home")
	utils.AssertContains(t, output, `aria-label="Main navigation"`)
}
