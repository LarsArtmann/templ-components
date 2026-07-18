// Package navigation provides rendering tests for navigation components.
package navigation

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestNavLinkRender(t *testing.T) {
	t.Parallel()
	t.Run("active link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NavLink(NavLinkProps{Href: "/", Text: navItemHome}, "/"))
		utils.AssertContains(t, output, "Home")
		utils.AssertContains(t, output, `href="/"`)
		utils.AssertContains(t, output, `aria-current="page"`)
	})

	t.Run("inactive link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			NavLink(NavLinkProps{Href: navItemAbout, Text: navItemAbout}, "/"),
		)
		utils.AssertContains(t, output, navItemAbout)
		utils.AssertContains(t, output, `href="/about"`)
		utils.AssertNotContains(t, output, `aria-current="page"`)
	})

	t.Run("external link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			NavLink(
				NavLinkProps{Href: navHrefExternal, Text: navItemExternal, External: true},
				"/",
			),
		)
		utils.AssertContains(t, output, `target="_blank"`)
		utils.AssertContains(t, output, `rel="noopener noreferrer"`)
	})

	t.Run("custom class propagated via utils.Class merge", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			NavLink(
				NavLinkProps{
					BaseProps: utils.BaseProps{Class: "my-custom-link"},
					Href:      "/",
					Text:      navItemHome,
				},
				"/",
			),
		)
		utils.AssertContains(t, output, "my-custom-link")
		utils.AssertContains(t, output, "inline-flex")
	})

	t.Run("custom ID propagated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(
			t,
			NavLink(
				NavLinkProps{
					BaseProps: utils.BaseProps{ID: "nav-home"},
					Href:      "/",
					Text:      navItemHome,
				},
				"/",
			),
		)
		utils.AssertContains(t, output, `id="nav-home"`)
	})
}

func TestMobileNavLinkRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileNavLink(NavLinkProps{Href: "/", Text: navItemHome}, "/"))
	utils.AssertContains(t, output, "Home")
	utils.AssertContains(t, output, `aria-current="page"`)
}

func TestMobileNavLinkClassPropagation(t *testing.T) {
	t.Parallel()
	output := utils.Render(
		t,
		MobileNavLink(
			NavLinkProps{
				BaseProps: utils.BaseProps{Class: "my-mobile-link"},
				Href:      "/",
				Text:      navItemHome,
			},
			"/",
		),
	)
	utils.AssertContains(t, output, "my-mobile-link")
	utils.AssertContainsAll(t, output, "block", "border-s-4")
}

func TestMobileNavLinkVariants(t *testing.T) {
	t.Parallel()

	t.Run("inactive link has no aria-current", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileNavLink(
			NavLinkProps{Href: "/about", Text: navItemAbout},
			"/",
		))
		utils.AssertContains(t, output, navItemAbout)
		utils.AssertNotContains(t, output, `aria-current="page"`)
		utils.AssertContains(t, output, "border-transparent")
	})

	t.Run("link with custom href", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileNavLink(
			NavLinkProps{
				Href: navHrefExternal,
				Text: navItemExternal,
			},
			"/",
		))
		utils.AssertContains(t, output, `href="https://example.com"`)
		utils.AssertContains(t, output, "External")
	})
}

func TestBreadcrumbsRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{
		{Text: navItemHome, Href: "/"},
		{Text: "Users", Active: true},
	}}))
	utils.AssertContains(t, output, "Home")
	utils.AssertContains(t, output, "Users")
	utils.AssertContains(t, output, `href="/"`)
}

func TestBreadcrumbsCoverage(t *testing.T) {
	t.Parallel()

	t.Run("custom ID and class propagated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
			BaseProps: utils.BaseProps{ID: "trail", Class: "my-crumbs"},
			Items: []BreadcrumbItem{
				{Text: navItemHome, Href: "/"},
			},
		}))
		utils.AssertContains(t, output, `id="trail"`)
		utils.AssertContains(t, output, "my-crumbs")
	})

	t.Run("separator icon between items", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{
			{Text: navItemHome, Href: "/"},
			{Text: navItemUsers, Href: navPathUsers},
			{Text: breadcrumbItemEdit, Active: true},
		}}))
		utils.AssertContains(t, output, "stroke-linecap")
	})

	t.Run("empty items renders nav", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{}}))
		utils.AssertContains(t, output, `aria-label="Breadcrumb"`)
	})

	t.Run("item with empty Href renders as span", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{
			{Text: "Root"},
		}}))
		utils.AssertContains(t, output, "Root")
		utils.AssertNotContains(t, output, "<a ")
	})
}

func TestFooterRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Footer(FooterProps{BrandText: "MyApp"}))
	utils.AssertContains(t, output, "MyApp")
	utils.AssertContains(t, output, "All rights reserved")
}

func TestMobileMenuToggleRender(t *testing.T) {
	t.Parallel()
	t.Run("shown", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileMenuToggle(true, "test-menu"))
		utils.AssertContains(t, output, `data-mobile-menu-toggle="test-menu"`)
		utils.AssertContains(t, output, `aria-controls="test-menu"`)
		utils.AssertContains(t, output, `aria-expanded="false"`)
	})
	t.Run("hidden", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileMenuToggle(false, "test-menu"))
		utils.AssertNotContains(t, output, "button")
	})
}

var testNavLinks = []NavLinkProps{
	{Href: "/", Text: navItemHome},
	{Href: navItemAbout, Text: navItemAbout},
}

func TestMobileMenuRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileMenu(testNavLinks, "/", "test-nonce", "test-menu"))
	utils.AssertContains(t, output, "Home")
	utils.AssertContains(t, output, navItemAbout)
	utils.AssertContains(t, output, `id="test-menu"`)
	utils.AssertNotContains(t, output, `tc-mobile-menu-tc-mobile-menu-`)
	utils.AssertContains(t, output, `nonce="test-nonce"`)
}

func TestSimpleNavRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleNav(SimpleNavProps{
		BrandText: "MyApp", BrandHref: "/", Links: testNavLinks, CurrentPath: "/",
	}))
	utils.AssertContains(t, output, "MyApp")
	utils.AssertContains(t, output, "Home")
	utils.AssertContains(t, output, `aria-label="Main navigation"`)
}

func TestMobileNavLinkExternalAddsSecurityAttrs(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileNavLink(NavLinkProps{
		Href: "https://example.com", Text: "External", External: true,
	}, "/"))
	utils.AssertContains(t, output, `target="_blank"`)
	utils.AssertContains(t, output, `rel="noopener noreferrer"`)
}

func TestTwoNavsProduceUniqueMenuIDs(t *testing.T) {
	t.Parallel()
	out1 := utils.Render(t, Nav(NavProps{Links: testNavLinks, CurrentPath: "/"}))
	out2 := utils.Render(t, Nav(NavProps{Links: testNavLinks, CurrentPath: "/"}))
	idCount1 := strings.Count(out1, `id="tc-mobile-menu-`)
	idCount2 := strings.Count(out2, `id="tc-mobile-menu-`)

	if idCount1 != 2 {
		t.Errorf("expected 2 mobile-menu ids (toggle + menu) in first Nav, got %d", idCount1)
	}

	if idCount2 != 2 {
		t.Errorf("expected 2 mobile-menu ids (toggle + menu) in second Nav, got %d", idCount2)
	}
	// Both should have valid (non-empty, non-hardcoded) IDs
	utils.AssertContains(t, out1, `data-mobile-menu-id="tc-mobile-menu-`)
	utils.AssertContains(t, out2, `data-mobile-menu-id="tc-mobile-menu-`)
}
