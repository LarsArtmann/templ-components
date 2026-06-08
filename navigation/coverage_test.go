package navigation

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// --- Pagination coverage (was 58.7%) ---

func TestPaginationEllipsis(t *testing.T) {
	t.Parallel()
	t.Run("many pages shows ellipsis at start", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 8,
			TotalPages:  20,
			BaseURL:     "/items",
			MaxVisible:  5,
		}))
		utils.AssertContains(t, output, "&hellip;")
		utils.AssertContains(t, output, `href="/items?page=1"`)
	})

	t.Run("many pages shows ellipsis at end", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  20,
			BaseURL:     "/items",
			MaxVisible:  5,
		}))
		utils.AssertContains(t, output, "&hellip;")
		utils.AssertContains(t, output, `href="/items?page=20"`)
	})

	t.Run("previous arrow has rel=prev", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 3,
			TotalPages:  5,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, `rel="prev"`)
	})

	t.Run("next arrow has rel=next", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  5,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, `rel="next"`)
	})

	t.Run("custom aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  5,
			BaseURL:     "/items",
			BaseProps:   utils.BaseProps{AriaLabel: "Results navigation"},
		}))
		utils.AssertContains(t, output, `aria-label="Results navigation"`)
	})

	t.Run("custom class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  5,
			BaseURL:     "/items",
			BaseProps:   utils.BaseProps{Class: "my-pager"},
		}))
		utils.AssertContains(t, output, "my-pager")
	})

	t.Run("custom attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  5,
			BaseURL:     "/items",
			BaseProps:   utils.BaseProps{Attrs: templ.Attributes{"data-testid": "pager"}},
		}))
		utils.AssertContains(t, output, `data-testid="pager"`)
	})

	t.Run("mobile pagination shows previous and next", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  5,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, "Previous")
		utils.AssertContains(t, output, "Next")
	})

	t.Run("first page mobile shows disabled previous", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 1,
			TotalPages:  5,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, "cursor-not-allowed")
	})

	t.Run("last page mobile shows disabled next", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 5,
			TotalPages:  5,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, "cursor-not-allowed")
	})

	t.Run("last page desktop disables next arrow", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 5,
			TotalPages:  5,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, "Next")
	})

	t.Run("page count text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  5,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, "Showing page")
		utils.AssertContains(t, output, "of")
	})
}

// --- NavLink / MobileNavLink coverage ---

func TestNavLinkExternal(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NavLink(NavLinkProps{
		Href:     "https://example.com",
		Text:     "External",
		External: true,
	}, "/"))
	utils.AssertContains(t, output, `target="_blank"`)
	utils.AssertContains(t, output, `rel="noopener noreferrer"`)
}

func TestNavLinkWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NavLink(NavLinkProps{
		Href:      "/about",
		Text:      "About",
		BaseProps: utils.BaseProps{AriaLabel: "About us"},
	}, "/home"))
	utils.AssertContains(t, output, `aria-label="About us"`)
}

func TestNavLinkWithClass(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NavLink(NavLinkProps{
		Href:      "/about",
		Text:      "About",
		BaseProps: utils.BaseProps{Class: "custom-link"},
	}, "/home"))
	utils.AssertContains(t, output, "custom-link")
}

func TestNavLinkActive(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NavLink(NavLinkProps{
		Href: "/about",
		Text: "About",
	}, "/about"))
	utils.AssertContains(t, output, `aria-current="page"`)
}

func TestMobileNavLinkActive(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileNavLink(NavLinkProps{
		Href: "/about",
		Text: "About",
	}, "/about"))
	utils.AssertContains(t, output, `aria-current="page"`)
}

func TestMobileNavLinkInactive(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileNavLink(NavLinkProps{
		Href: "/settings",
		Text: "Settings",
	}, "/home"))
	utils.AssertContains(t, output, "Settings")
	utils.AssertNotContains(t, output, `aria-current="page"`)
}

func TestMobileNavLinkWithClass(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileNavLink(NavLinkProps{
		Href:      "/settings",
		Text:      "Settings",
		BaseProps: utils.BaseProps{Class: "my-mobile-link"},
	}, "/home"))
	utils.AssertContains(t, output, "my-mobile-link")
}

func TestMobileNavLinkWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileNavLink(NavLinkProps{
		Href:      "/settings",
		Text:      "Settings",
		BaseProps: utils.BaseProps{AriaLabel: "Settings page"},
	}, "/home"))
	utils.AssertContains(t, output, `aria-label="Settings page"`)
}

// --- Breadcrumbs coverage ---

func TestBreadcrumbsWithCustomSeparator(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "/"},
			{Text: "Products", Href: "/products"},
		},
		Separator: ">",
	}))
	utils.AssertContains(t, output, ">")
}

// --- Nav coverage ---

func TestSimpleNavMinimal(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleNav("App", "/", nil, "/"))
	utils.AssertContains(t, output, "App")
}

func TestSimpleNavWithLinks(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleNav("App", "/", []NavLinkProps{
		{Href: "/about", Text: "About"},
		{Href: "/contact", Text: "Contact"},
	}, "/about"))
	utils.AssertContains(t, output, "About")
	utils.AssertContains(t, output, "Contact")
}

func TestFooterMinimal(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Footer("2024 Acme"))
	utils.AssertContains(t, output, "2024 Acme")
}

func TestDefaultPaginationProps(t *testing.T) {
	t.Parallel()
	props := DefaultPaginationProps()
	if props.QueryParam != "page" {
		t.Errorf("QueryParam = %q, want %q", props.QueryParam, "page")
	}
	if props.MaxVisible != 5 {
		t.Errorf("MaxVisible = %d, want 5", props.MaxVisible)
	}
}
