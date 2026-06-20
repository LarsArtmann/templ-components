// Package navigation provides behavior-driven tests for navigation components.
// These tests verify end-user-facing behavior: navigating pages, seeing active states, browsing breadcrumbs.
package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const (
	testItemsPath = "/items"
	navItemAbout  = "/about"
)

// --- Nav Behavior ---

func TestNavUserCanNavigateBetweenPages(t *testing.T) {
	t.Parallel()

	links := []NavLinkProps{
		{Href: "/", Text: navItemHome},
		{Href: navItemAbout, Text: navItemAbout},
	}
	for _, tt := range []struct {
		name        string
		currentPath string
		want        string
	}{
		{"user sees brand and navigation links", "/", navItemHome},
		{"user sees active link styling on current page", navItemAbout, "border-blue-500"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Nav(NavProps{
				Links:       links,
				CurrentPath: tt.currentPath,
			}))
			utils.AssertContains(t, output, tt.want)
			if tt.currentPath == "/" {
				utils.AssertContains(t, output, navItemAbout)
			}
		})
	}

	t.Run("user sees sticky navigation when enabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Sticky: true,
			Links:  []NavLinkProps{{Href: "/", Text: navItemHome}},
		}))
		utils.AssertContains(t, output, "sticky")
	})

	t.Run("user sees nav with semantic landmark", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(DefaultNavProps()))
		utils.AssertContains(t, output, "<nav")
	})
}

// --- SimpleNav Behavior ---

func TestSimpleNavUserGetsQuickNavigation(t *testing.T) {
	t.Parallel()

	t.Run("user sees brand text and links", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SimpleNav(SimpleNavProps{BrandText: "MyApp", BrandHref: "/", Links: []NavLinkProps{
			{Href: "/", Text: navItemHome},
			{Href: "/settings", Text: "Settings"},
		}, CurrentPath: "/"}))
		utils.AssertContains(t, output, "MyApp")
		utils.AssertContains(t, output, "Settings")
	})
}

// --- NavLink Behavior ---

func TestNavLinkUserSeesLinkStates(t *testing.T) {
	t.Parallel()

	t.Run("user sees active link for current page", func(t *testing.T) {
		t.Parallel()
		props := NavLinkProps{Href: "/dashboard", Text: "Dashboard"}
		output := utils.Render(t, NavLink(props, "/dashboard"))
		utils.AssertContains(t, output, "Dashboard")
		utils.AssertContains(t, output, "border-blue-500")
	})

	t.Run("user sees external link with security attributes", func(t *testing.T) {
		t.Parallel()
		props := NavLinkProps{Href: navHrefExternal, Text: navItemExternal, External: true}
		output := utils.Render(t, NavLink(props, ""))
		utils.AssertContains(t, output, `target="_blank"`)
		utils.AssertContains(t, output, `rel="noopener noreferrer"`)
	})
}

func breadcrumbHomeAndCurrent() []BreadcrumbItem {
	return []BreadcrumbItem{
		{Text: navItemHome, Href: "/"},
		{Text: "Current", Active: true},
	}
}

// --- Breadcrumbs Behavior ---

func renderBreadcrumbs(t *testing.T, items []BreadcrumbItem) string {
	t.Helper()
	return utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: items}))
}

func TestBreadcrumbsUserCanSeeWhereTheyAre(t *testing.T) {
	t.Parallel()

	t.Run("user sees breadcrumb trail from home to current page", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{
			{Text: navItemHome, Href: "/"},
			{Text: "Users", Href: "/users"},
			{Text: breadcrumbItemEdit, Active: true},
		}}))
		utils.AssertContains(t, output, navItemHome)
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, "Edit")
	})

	t.Run("breadcrumb attributes: aria-current, landmark, and clickable parents", func(t *testing.T) {
		t.Parallel()
		output := renderBreadcrumbs(t, breadcrumbHomeAndCurrent())
		utils.AssertContains(t, output, `aria-current="page"`)
		utils.AssertContains(t, output, `aria-label="Breadcrumb"`)
		utils.AssertContains(t, output, `href="/"`)
	})

	t.Run("user sees JSON-LD structured data when enabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
			Items: []BreadcrumbItem{
				{Text: navItemHome, Href: "https://example.com/"},
				{Text: "Users", Href: "https://example.com/users"},
			},
			JSONLD: true,
		}))
		utils.AssertContains(t, output, `<script type="application/ld+json">`)
		utils.AssertContains(t, output, `"@context":"https://schema.org"`)
		utils.AssertContains(t, output, `"@type":"BreadcrumbList"`)
		utils.AssertContains(t, output, `"name":"Home"`)
		utils.AssertContains(t, output, `"name":"Users"`)
	})

	t.Run("user sees custom separator when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
			Items:     breadcrumbHomeAndCurrent(),
			Separator: "/",
		}))
		utils.AssertContains(t, output, "/")
	})
}

func renderDefaultPagination(t *testing.T, currentPage, totalPages uint) string {
	t.Helper()
	props := DefaultPaginationProps()
	props.CurrentPage = currentPage
	props.TotalPages = totalPages
	props.BaseURL = testItemsPath
	return utils.Render(t, Pagination(props))
}

// --- Pagination Behavior ---

func TestPaginationUserCanBrowsePages(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name        string
		currentPage uint
		totalPages  uint
		want        string
	}{
		{"user sees page numbers for multiple pages", 1, 5, testItemsPath},
		{"user sees previous and next navigation", 2, 5, `aria-label="`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := renderDefaultPagination(t, tt.currentPage, tt.totalPages)
			utils.AssertContains(t, output, tt.want)
		})
	}

	t.Run("user sees nothing when only one page exists", func(t *testing.T) {
		t.Parallel()
		output := renderDefaultPagination(t, 1, 1)
		utils.AssertNotContains(t, output, "<nav")
	})

	t.Run("user can customize query parameter name", func(t *testing.T) {
		t.Parallel()
		props := DefaultPaginationProps()
		props.CurrentPage = 1
		props.TotalPages = 3
		props.BaseURL = "/search"
		props.QueryParam = "p"
		output := utils.Render(t, Pagination(props))
		utils.AssertContains(t, output, "p=")
	})
}

// --- Footer Behavior ---

func TestFooterUserSeesCopyright(t *testing.T) {
	t.Parallel()

	t.Run("user sees brand name, copyright year, and dark mode classes", func(t *testing.T) {
		t.Parallel()
		assertFooterContainsAll(t, "MyApp", "All rights reserved", "<footer",
			"dark:border-gray-800", "dark:bg-gray-900", "dark:text-gray-400")
	})
}

func TestNavEmptyLinks(t *testing.T) {
	t.Parallel()
	t.Run("nav with no links renders brand only", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SimpleNav(SimpleNavProps{BrandText: "App", BrandHref: "/", CurrentPath: "/"}))
		utils.AssertContains(t, output, "App")
		utils.AssertNotContains(t, output, `aria-current="page"`)
	})
}
