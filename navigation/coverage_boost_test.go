package navigation

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestBreadcrumbsFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with JSONLD and separator", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
			BaseProps: utils.BaseProps{ID: "bc-1", Class: "py-2", AriaLabel: "Trail"},
			Items: []BreadcrumbItem{
				{Text: "Home", Href: "/"},
				{Text: "Users", Href: "/users"},
				{Text: "Profile", Active: true},
			},
			Separator: "/",
			JSONLD:    true,
		}))
		utils.AssertContains(t, output, "Home")
		utils.AssertContains(t, output, "application/ld+json")
	})
}

func TestPaginationFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("many pages with ellipsis", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			BaseProps:   utils.BaseProps{ID: "pager", AriaLabel: "Pagination"},
			CurrentPage: 5,
			TotalPages:  20,
			BaseURL:     "/items",
			QueryParam:  "page",
			MaxVisible:  7,
		}))
		utils.AssertContains(t, output, "/items")
	})
	t.Run("single page", func(t *testing.T) {
		t.Parallel()
		utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 1,
			TotalPages:  1,
			BaseURL:     "/x",
			QueryParam:  "p",
		}))
	})
}

func TestNavFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with all props", func(t *testing.T) {
		t.Parallel()
		brand := templ.Raw(`<span>Brand</span>`)
		output := utils.Render(t, Nav(NavProps{
			BaseProps:   utils.BaseProps{ID: "main-nav", Class: "shadow", AriaLabel: "Main"},
			Brand:       brand,
			Links:       []NavLinkProps{{Href: "/", Text: "Home"}, {Href: "/about", Text: "About"}},
			CurrentPath: "/about",
			Sticky:      true,
		}))
		utils.AssertContains(t, output, "Brand")
		utils.AssertContains(t, output, "About")
	})
}

func TestSimpleNavFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SimpleNav(SimpleNavProps{
			BaseProps:   utils.BaseProps{ID: "simple-nav", AriaLabel: "Navigation"},
			BrandText:   "MyApp",
			BrandHref:   "/home",
			Links:       []NavLinkProps{{Href: "/dashboard", Text: "Dashboard"}},
			CurrentPath: "/dashboard",
		}))
		utils.AssertContains(t, output, "MyApp")
		utils.AssertContains(t, output, "Dashboard")
	})
}

func TestMobileMenuFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("renders with links", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, MobileMenu([]NavLinkProps{
			{Href: "/a", Text: "Page A"},
			{Href: "/b", Text: "Page B"},
		}, "/a", "nonce-m", "cov-menu"))
		utils.AssertContains(t, output, "Page A")
	})
}

func TestFooterFullCoverage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Footer("MyApp"))
	utils.AssertContains(t, output, "MyApp")
}
