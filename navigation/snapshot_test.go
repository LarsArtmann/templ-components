// Package navigation provides rendering tests for navigation components.
package navigation

import (
	"regexp"
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
					Class: "my-custom-link",
					Href:  "/",
					Text:  navItemHome,
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
					ID:   "nav-home",
					Href: "/",
					Text: navItemHome,
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
				Class: "my-mobile-link",
				Href:  "/",
				Text:  navItemHome,
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
			ID: "trail", Class: "my-crumbs",
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
		output := utils.Render(t, MobileMenuToggle(true, "test-menu", false))
		utils.AssertContains(t, output, `data-mobile-menu-toggle="test-menu"`)
		utils.AssertContains(t, output, `aria-controls="test-menu"`)
		utils.AssertContains(t, output, `aria-expanded="false"`)
	})
	t.Run("hidden", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileMenuToggle(false, "test-menu", false))
		utils.AssertNotContains(t, output, "button")
	})
}

var testNavLinks = []NavLinkProps{
	{Href: "/", Text: navItemHome},
	{Href: navItemAbout, Text: navItemAbout},
}

func TestMobileMenuRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileMenu(testNavLinks, "/", "test-nonce", "test-menu", false))
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

// navMenuIDRe matches the content-derived mobile-menu id: the same
// tc-<prefix>-<16 hex> shape utils.EnsureID mints, so goldens normalize it.
var navMenuIDRe = regexp.MustCompile(`tc-mobile-menu-[a-f0-9]{16}`)

// domIDRe matches id="" attributes only (not data-*-id=""), so a duplicate-id
// assertion sees exactly the ids getElementById() resolves.
var domIDRe = regexp.MustCompile(`(?:^|\s)id="([^"]*)"`)

func domIDs(html string) []string {
	matches := domIDRe.FindAllStringSubmatch(html, -1)
	ids := make([]string, 0, len(matches))

	for _, match := range matches {
		ids = append(ids, match[1])
	}

	return ids
}

// TestNavMenuIDIsDeterministic pins the caching contract: the mobile-menu id
// is derived from the nav's content, so rendering the same props twice is
// byte-identical and ETag revalidation can answer 304 instead of re-sending
// a freshly randomized page on every request.
func TestNavMenuIDIsDeterministic(t *testing.T) {
	t.Parallel()

	props := NavProps{Links: testNavLinks, CurrentPath: "/"}
	first := utils.Render(t, Nav(props))
	second := utils.Render(t, Nav(props))

	if first != second {
		t.Error("Nav rendered different HTML for identical props — the mobile-menu id must not be render-random")
	}

	if got := navMenuIDRe.FindString(first); got == "" {
		t.Errorf("mobile-menu id missing from the rendered nav:\n%s", first)
	}
}

// TestNavMenuIDTracksContent proves the derived id is content-scoped: two
// navs that differ in links or current path get different menu ids, so their
// toggles cannot resolve to each other's menu.
func TestNavMenuIDTracksContent(t *testing.T) {
	t.Parallel()
	base := navMenuIDRe.FindString(utils.Render(t, Nav(NavProps{Links: testNavLinks, CurrentPath: "/"})))
	otherPath := navMenuIDRe.FindString(utils.Render(t, Nav(NavProps{Links: testNavLinks, CurrentPath: "/about"})))

	if base == "" || otherPath == "" {
		t.Fatalf("mobile-menu id missing (base=%q otherPath=%q)", base, otherPath)
	}

	if base == otherPath {
		t.Errorf("navs with different current paths share the menu id %q", base)
	}
}

// TestNavMenuIDDerivesFromCallerID pins the duplicate-id fix: the <nav> keeps
// the consumer's id and the mobile menu gets its own derived id, so a document
// never carries the same id twice (invalid HTML, and aria-controls/
// getElementById resolve to the wrong element).
func TestNavMenuIDDerivesFromCallerID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Nav(NavProps{
		ID:          "site-nav",
		Links:       testNavLinks,
		CurrentPath: "/",
	}))

	utils.AssertContains(t, output, `id="site-nav-mobile-menu"`)
	utils.AssertContains(t, output, `aria-controls="site-nav-mobile-menu"`)

	seen := map[string]int{}
	for _, id := range domIDs(output) {
		seen[id]++
	}

	if seen["site-nav"] != 1 {
		t.Errorf("expected exactly one id=%q, got %d", "site-nav", seen["site-nav"])
	}

	for id, count := range seen {
		if count > 1 {
			t.Errorf("duplicate id %q rendered %d times", id, count)
		}
	}
}

func TestNavContainerAware(t *testing.T) {
	t.Parallel()

	t.Run("viewport breakpoints by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			ID:    "vnav",
			Links: []NavLinkProps{{Href: "/", Text: "Home"}},
		}))
		utils.AssertNotContains(t, output, "@container")
		utils.AssertContains(t, output, "sm:flex")
		utils.AssertNotContains(t, output, "@sm:flex")
	})

	t.Run("container breakpoints when flag set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			ID:             "cnav",
			Links:          []NavLinkProps{{Href: "/", Text: "Home"}},
			ContainerAware: true,
		}))
		utils.AssertContains(t, output, "@container")
		utils.AssertContains(t, output, "@sm:flex")
		utils.AssertContains(t, output, "@sm:hidden")
		utils.AssertNotContains(t, output, " sm:flex")
	})
}

func TestPaginationContainerAware(t *testing.T) {
	t.Parallel()

	t.Run("viewport breakpoints by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  5,
			BaseURL:     "/items",
		}))
		utils.AssertNotContains(t, output, "@container")
		utils.AssertContains(t, output, "sm:hidden")
		utils.AssertContains(t, output, "sm:flex")
		utils.AssertNotContains(t, output, "@sm:")
	})

	t.Run("container breakpoints when flag set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage:    2,
			TotalPages:     5,
			BaseURL:        "/items",
			ContainerAware: true,
		}))
		utils.AssertContains(t, output, "@container")
		utils.AssertContains(t, output, "@sm:hidden")
		utils.AssertContains(t, output, "@sm:flex")
		utils.AssertNotContains(t, output, " sm:")
	})
}
