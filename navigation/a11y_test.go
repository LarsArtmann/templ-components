package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const (
	navItemHome           = "Home"
	navItemUsers          = "Users"
	navPathUsers          = "/users"
	breadcrumbItemCurrent = "Current"
	breadcrumbItemEdit    = "Edit"
	navItemExternal       = "External"
	navHrefExternal       = "https://example.com"
)

func breadcrumbHomeOnly() []BreadcrumbItem {
	return []BreadcrumbItem{{Text: navItemHome, Href: "/"}}
}

func breadcrumbHomeAndActive() []BreadcrumbItem {
	return []BreadcrumbItem{
		{Text: navItemHome, Href: "/"},
		{Text: navItemUsers, Active: true},
	}
}

func TestBreadcrumbsA11y(t *testing.T) {
	t.Parallel()

	t.Run("nav has aria-label", func(t *testing.T) {
		t.Parallel()
		output := renderBreadcrumbs(t, breadcrumbHomeOnly())
		utils.AssertContains(t, output, `aria-label="Breadcrumb"`)
	})

	t.Run("active item has aria-current", func(t *testing.T) {
		t.Parallel()
		output := renderBreadcrumbs(t, breadcrumbHomeAndActive())
		utils.AssertContains(t, output, `aria-current="page"`)
	})

	t.Run("inactive items are links", func(t *testing.T) {
		t.Parallel()
		output := renderBreadcrumbs(t, breadcrumbHomeAndActive())
		utils.AssertContains(t, output, `<a href="/"`)
	})

	t.Run("empty list renders nav", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{}}))
		utils.AssertContains(t, output, "<nav")
		utils.AssertContains(t, output, `aria-label="Breadcrumb"`)
	})

	t.Run("single item with no href renders span", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{
			{Text: navItemHome},
		}}))
		utils.AssertContains(t, output, navItemHome)
		utils.AssertContains(t, output, `aria-current="page"`)
		utils.AssertNotContains(t, output, `<a`)
	})

	t.Run("three levels with separators", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{
			{Text: navItemHome, Href: "/"},
			{Text: navItemUsers, Href: navPathUsers},
			{Text: breadcrumbItemEdit, Active: true},
		}}))
		utils.AssertContains(t, output, navItemHome)
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, "Edit")
		utils.AssertContains(t, output, "text-gray-400")
	})
}

func TestBreadcrumbsDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("breadcrumb items have dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{Items: []BreadcrumbItem{
			{Text: breadcrumbItemCurrent, Active: true},
			{Text: navItemHome, Href: "/"},
		}}))
		utils.AssertContains(t, output, "dark:text-gray-400")
	})
}

func TestNavRender(t *testing.T) {
	t.Parallel()

	t.Run("nav with brand and links", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Brand:       simpleBrand("App", "/"),
			Links:       testNavLinks,
			CurrentPath: "/",
		}))
		utils.AssertContains(t, output, `aria-label="Main navigation"`)
		utils.AssertContains(t, output, "App")
		utils.AssertContains(t, output, navItemHome)
		utils.AssertContains(t, output, navItemAbout)
	})

	t.Run("sticky nav has sticky class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Brand:  simpleBrand("App", "/"),
			Sticky: true,
		}))
		utils.AssertContains(t, output, "sticky top-0")
	})

	t.Run("non-sticky nav has no sticky class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Brand: simpleBrand("App", "/"),
		}))
		utils.AssertNotContains(t, output, "sticky top-0")
	})

	t.Run("nav with right items", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Brand:      simpleBrand("App", "/"),
			RightItems: simpleBrand("Profile", "/profile"),
		}))
		utils.AssertContains(t, output, "Profile")
	})

	t.Run("nav without brand", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Links: testNavLinks,
		}))
		utils.AssertContains(t, output, navItemHome)
		utils.AssertNotContains(t, output, "text-xl font-bold")
	})

	t.Run("nav with custom ID and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			BaseProps: utils.BaseProps{
				ID:    "main-nav",
				Class: "shadow-lg",
			},
			Links: testNavLinks,
		}))
		utils.AssertContains(t, output, `id="main-nav"`)
		utils.AssertContains(t, output, "shadow-lg")
	})

	t.Run("nav with nonce passes to mobile menu script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Links:       testNavLinks,
			CurrentPath: "/",
			BaseProps:   utils.BaseProps{Nonce: "test-nonce-123"},
		}))
		utils.AssertContains(t, output, `nonce="test-nonce-123"`)
	})
}

func TestNavDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("nav has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Brand: simpleBrand("App", "/"),
		}))
		utils.AssertContains(t, output, "dark:border-gray-800")
		utils.AssertContains(t, output, "dark:bg-gray-900")
	})
}

func renderFooter(t *testing.T) string {
	t.Helper()
	return utils.Render(t, Footer("MyApp"))
}

func assertFooterContainsAll(t *testing.T, contains ...string) {
	t.Helper()
	output := renderFooter(t)
	for _, s := range contains {
		utils.AssertContains(t, output, s)
	}
}

func TestFooterDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("footer has dark mode classes", func(t *testing.T) {
		t.Parallel()
		assertFooterContainsAll(t, "dark:border-gray-800", "dark:bg-gray-900", "dark:text-gray-400")
	})
}

func TestAriaLabelOverride(t *testing.T) {
	t.Parallel()
	customLabel := "Custom navigation label"

	t.Run("Nav AriaLabel overrides default", func(t *testing.T) {
		t.Parallel()
		props := NavProps{
			Links:     testNavLinks,
			BaseProps: utils.BaseProps{AriaLabel: customLabel},
		}
		output := utils.Render(t, Nav(props))
		utils.AssertContains(t, output, `aria-label="`+customLabel+`"`)
		utils.AssertNotContains(t, output, `aria-label="Main navigation"`)
	})

	t.Run("Breadcrumbs AriaLabel overrides default", func(t *testing.T) {
		t.Parallel()
		props := BreadcrumbsProps{
			Items:     breadcrumbHomeOnly(),
			BaseProps: utils.BaseProps{AriaLabel: customLabel},
		}
		output := utils.Render(t, Breadcrumbs(props))
		utils.AssertContains(t, output, `aria-label="`+customLabel+`"`)
		utils.AssertNotContains(t, output, `aria-label="Breadcrumb"`)
	})

	t.Run("Pagination AriaLabel overrides default", func(t *testing.T) {
		t.Parallel()
		props := PaginationProps{
			CurrentPage: 1,
			TotalPages:  5,
			QueryParam:  "page",
			BaseProps:   utils.BaseProps{AriaLabel: customLabel},
		}
		output := utils.Render(t, Pagination(props))
		utils.AssertContains(t, output, `aria-label="`+customLabel+`"`)
	})

	t.Run("NavLink AriaLabel renders", func(t *testing.T) {
		t.Parallel()
		props := NavLinkProps{
			Href:      "/",
			Text:      "Home",
			BaseProps: utils.BaseProps{AriaLabel: customLabel},
		}
		output := utils.Render(t, NavLink(props, "/other"))
		utils.AssertContains(t, output, `aria-label="`+customLabel+`"`)
	})
}
