package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestBreadcrumbsA11y(t *testing.T) {
	t.Parallel()

	t.Run("nav has aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs([]BreadcrumbItem{
			{Text: "Home", Href: "/"},
		}))
		utils.AssertContains(t, output, `aria-label="Breadcrumb"`)
	})

	t.Run("active item has aria-current", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs([]BreadcrumbItem{
			{Text: "Home", Href: "/"},
			{Text: "Users", Active: true},
		}))
		utils.AssertContains(t, output, `aria-current="page"`)
	})

	t.Run("inactive items are links", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs([]BreadcrumbItem{
			{Text: "Home", Href: "/"},
			{Text: "Users", Active: true},
		}))
		utils.AssertContains(t, output, `<a href="/"`)
	})

	t.Run("empty list renders nav", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs([]BreadcrumbItem{}))
		utils.AssertContains(t, output, "<nav")
		utils.AssertContains(t, output, `aria-label="Breadcrumb"`)
	})

	t.Run("single item with no href renders span", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs([]BreadcrumbItem{
			{Text: "Home"},
		}))
		utils.AssertContains(t, output, "Home")
		utils.AssertContains(t, output, `aria-current="page"`)
		utils.AssertNotContains(t, output, `<a`)
	})

	t.Run("three levels with separators", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs([]BreadcrumbItem{
			{Text: "Home", Href: "/"},
			{Text: "Users", Href: "/users"},
			{Text: "Edit", Active: true},
		}))
		utils.AssertContains(t, output, "Home")
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, "Edit")
		utils.AssertContains(t, output, "text-gray-400")
	})
}

func TestBreadcrumbsDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("active item has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs([]BreadcrumbItem{
			{Text: "Current", Active: true},
		}))
		utils.AssertContains(t, output, "dark:text-gray-400")
	})

	t.Run("link items have dark mode hover", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs([]BreadcrumbItem{
			{Text: "Home", Href: "/"},
		}))
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
		utils.AssertContains(t, output, "Home")
		utils.AssertContains(t, output, "About")
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
		utils.AssertContains(t, output, "Home")
		utils.AssertNotContains(t, output, "text-xl font-bold")
	})
}

func TestNavDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("nav has dark mode border", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Brand: simpleBrand("App", "/"),
		}))
		utils.AssertContains(t, output, "dark:border-gray-800")
	})

	t.Run("nav has dark mode background", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Nav(NavProps{
			Brand: simpleBrand("App", "/"),
		}))
		utils.AssertContains(t, output, "dark:bg-gray-900")
	})
}

func TestFooterDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("footer has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer("MyApp"))
		utils.AssertContains(t, output, "dark:border-gray-800")
		utils.AssertContains(t, output, "dark:bg-gray-900")
		utils.AssertContains(t, output, "dark:text-gray-400")
	})
}
