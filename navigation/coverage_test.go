package navigation

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// --- Pagination coverage (was 58.7%) ---

func TestPaginationEllipsis(t *testing.T) {
	t.Parallel()
	t.Run("ellipsis at start and end", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			name        string
			currentPage uint
			wantHref    string
		}{
			{"start", 8, `href="/items?page=1"`},
			{"end", 2, `href="/items?page=20"`},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				output := utils.Render(t, Pagination(PaginationProps{
					CurrentPage: tc.currentPage,
					TotalPages:  20,
					BaseURL:     "/items",
					MaxVisible:  5,
				}))
				utils.AssertContains(t, output, "&hellip;")
				utils.AssertContains(t, output, tc.wantHref)
			})
		}
	})

	t.Run("rel attributes on arrows", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			name        string
			currentPage uint
			wantRel     string
		}{
			{"previous", 3, `rel="prev"`},
			{"next", 2, `rel="next"`},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				output := utils.Render(t, Pagination(PaginationProps{
					CurrentPage: tc.currentPage,
					TotalPages:  5,
					BaseURL:     "/items",
				}))
				utils.AssertContains(t, output, tc.wantRel)
			})
		}
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

	t.Run("boundary pages disable arrows", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			name        string
			currentPage uint
			totalPages  uint
			wantText    string
		}{
			{"first page disables previous", 1, 5, "cursor-not-allowed"},
			{"last page disables next", 5, 5, "cursor-not-allowed"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				output := utils.Render(t, Pagination(PaginationProps{
					CurrentPage: tc.currentPage,
					TotalPages:  tc.totalPages,
					BaseURL:     "/items",
				}))
				utils.AssertContains(t, output, tc.wantText)
			})
		}
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
	output := utils.Render(t, SimpleNav(SimpleNavProps{BrandText: "App", BrandHref: "/", CurrentPath: "/"}))
	utils.AssertContains(t, output, "App")
}

func TestSimpleNavWithLinks(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleNav(SimpleNavProps{BrandText: "App", BrandHref: "/", Links: []NavLinkProps{
		{Href: "/about", Text: "About"},
		{Href: "/contact", Text: "Contact"},
	}, CurrentPath: "/about"}))
	utils.AssertContains(t, output, "About")
	utils.AssertContains(t, output, "Contact")
}

func TestSimpleNavForwardsRightItems(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleNav(SimpleNavProps{
		BrandText:  "App",
		BrandHref:  "/",
		RightItems: simpleBrand("Sign in", "/login"),
	}))
	utils.AssertContains(t, output, "Sign in")
	utils.AssertContains(t, output, `/login`)
}

func TestSimpleNavOmitsRightItemsWhenNil(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleNav(SimpleNavProps{BrandText: "App", BrandHref: "/"}))
	// RightItems slot is empty: no "flex items-center gap-2" wrapper rendered.
	if got := strings.Count(output, "Sign in"); got != 0 {
		t.Errorf("expected no RightItems text, got %d occurrence(s)", got)
	}
}

func TestFooterMinimal(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Footer(FooterProps{BrandText: "2024 Acme"}))
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
