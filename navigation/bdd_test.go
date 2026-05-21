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

	t.Run("user sees brand and navigation links", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Links: []NavLinkProps{
				{Href: "/", Text: navItemHome},
				{Href: navItemAbout, Text: navItemAbout},
			},
			CurrentPath: "/",
		}))
		utils.AssertContains(t, output, navItemHome)
		utils.AssertContains(t, output, navItemAbout)
	})

	t.Run("user sees active link styling on current page", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Links: []NavLinkProps{
				{Href: "/", Text: navItemHome},
				{Href: navItemAbout, Text: navItemAbout},
			},
			CurrentPath: navItemAbout,
		}))
		utils.AssertContains(t, output, "border-blue-500")
	})

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
		output := utils.Render(t, SimpleNav("MyApp", "/", []NavLinkProps{
			{Href: "/", Text: navItemHome},
			{Href: "/settings", Text: "Settings"},
		}, "/"))
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

	t.Run("user sees current page as aria-current", func(t *testing.T) {
		t.Parallel()
		output := renderBreadcrumbs(t, breadcrumbHomeAndCurrent())
		utils.AssertContains(t, output, `aria-current="page"`)
	})

	t.Run("user sees breadcrumb navigation landmark", func(t *testing.T) {
		t.Parallel()
		output := renderBreadcrumbs(t, breadcrumbHomeAndCurrent())
		utils.AssertContains(t, output, `aria-label="Breadcrumb"`)
	})

	t.Run("user sees clickable parent breadcrumbs", func(t *testing.T) {
		t.Parallel()
		output := renderBreadcrumbs(t, breadcrumbHomeAndCurrent())
		utils.AssertContains(t, output, `href="/"`)
	})
}

func renderDefaultPagination(t *testing.T, currentPage, totalPages int) string {
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

	t.Run("user sees page numbers for multiple pages", func(t *testing.T) {
		t.Parallel()
		output := renderDefaultPagination(t, 1, 5)
		utils.AssertContains(t, output, testItemsPath)
	})

	t.Run("user sees previous and next navigation", func(t *testing.T) {
		t.Parallel()
		output := renderDefaultPagination(t, 2, 5)
		utils.AssertContains(t, output, `aria-label="`)
	})

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

	t.Run("user sees brand name and copyright year", func(t *testing.T) {
		t.Parallel()
		assertFooterContainsAll(t, "MyApp", "All rights reserved", "<footer")
	})
}
