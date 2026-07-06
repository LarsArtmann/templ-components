package navigation

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

// ---------------------------------------------------------------------------
// Breadcrumbs: custom separator, JSON-LD, BaseProps
// ---------------------------------------------------------------------------

func TestBreadcrumbsCustomSeparator(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "/"},
			{Text: "Products", Href: "/products"},
			{Text: "Detail", Active: true},
		},
		Separator: ">",
	}))
	utils.AssertContainsAll(t, output, "Home", "Products", "Detail")
}

func TestBreadcrumbsWithJSONLD(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "https://example.com/"},
			{Text: "Blog", Href: "https://example.com/blog"},
		},
		JSONLD: true,
	}))
	utils.AssertContains(t, output, "application/ld+json")
}

func TestBreadcrumbsWithBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "/"},
		},
		BaseProps: utils.BaseProps{
			ID:        "bc-1",
			AriaLabel: "Site breadcrumbs",
		},
	}))
	utils.AssertContainsAll(t, output, `id="bc-1"`, `aria-label="Site breadcrumbs"`)
}

func TestBreadcrumbsSingleItem(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Only", Active: true},
		},
	}))
	utils.AssertContains(t, output, "Only")
}

// ---------------------------------------------------------------------------
// Pagination: various page counts, BaseURL, query param
// ---------------------------------------------------------------------------

func TestPaginationFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Pagination(PaginationProps{
		CurrentPage: 3,
		TotalPages:  10,
		BaseURL:     "/items",
		QueryParam:  "p",
		MaxVisible:  5,
		BaseProps: utils.BaseProps{
			ID:        "pager",
			AriaLabel: "Item pagination",
		},
	}))
	utils.AssertContainsAll(t, output, `id="pager"`, "/items")
}

func TestPaginationFirstPage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Pagination(PaginationProps{
		CurrentPage: 1,
		TotalPages:  5,
		BaseURL:     "/list",
	}))
	utils.AssertContains(t, output, "/list")
}

func TestPaginationLastPage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Pagination(PaginationProps{
		CurrentPage: 5,
		TotalPages:  5,
		BaseURL:     "/list",
	}))
	utils.AssertContains(t, output, "/list")
}

func TestPaginationWithEllipsis(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Pagination(PaginationProps{
		CurrentPage: 5,
		TotalPages:  20,
		BaseURL:     "/page",
		MaxVisible:  5,
	}))
	// Should have ellipsis when many pages
	utils.AssertContains(t, output, "&hellip;")
}

func TestPaginationSinglePage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Pagination(PaginationProps{
		CurrentPage: 1,
		TotalPages:  1,
		BaseURL:     "/single",
	}))
	_ = output // should not panic
}

// ---------------------------------------------------------------------------
// NavLink: active/inactive, external link, BaseProps
// ---------------------------------------------------------------------------

func TestNavLinkActiveState(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NavLink(NavLinkProps{
		Href: "/dashboard",
		Text: "Dashboard",
		BaseProps: utils.BaseProps{
			ID:        "nl-dashboard",
			AriaLabel: "Dashboard link",
		},
	}, "/dashboard"))
	utils.AssertContainsAll(t, output,
		`id="nl-dashboard"`,
		`aria-label="Dashboard link"`,
		"Dashboard",
		`aria-current="page"`,
	)
}

func TestNavLinkInactive(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NavLink(NavLinkProps{
		Href: "/settings",
		Text: "Settings",
	}, "/dashboard"))
	utils.AssertContainsAll(t, output, `href="/settings"`, "Settings")
}

func TestNavLinkExternalLink(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NavLink(NavLinkProps{
		Href:     "https://external.com",
		Text:     "External",
		External: true,
	}, ""))
	utils.AssertContainsAll(t, output,
		`href="https://external.com"`,
		`target="_blank"`,
		`rel="noopener noreferrer"`,
	)
}

func TestMobileNavLinkActiveRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileNavLink(NavLinkProps{
		Href: "/profile",
		Text: "Profile",
	}, "/profile"))
	utils.AssertContainsAll(t, output, "Profile", `aria-current="page"`)
}

// ---------------------------------------------------------------------------
// Nav + SimpleNav: full props, RightItems slot
// ---------------------------------------------------------------------------

func TestNavFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Nav(NavProps{
		Brand: templ.Raw("<span>MyApp</span>"),
		Links: []NavLinkProps{
			{Href: "/", Text: "Home"},
			{Href: "/about", Text: "About"},
		},
		CurrentPath: "/about",
		Sticky:      true,
		BaseProps: utils.BaseProps{
			ID:        "main-nav",
			AriaLabel: "Main navigation",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="main-nav"`,
		"MyApp",
		"Home",
		"About",
	)
}

func TestSimpleNavFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleNav(SimpleNavProps{
		BrandText: "MyApp",
		BrandHref: "/",
		Links: []NavLinkProps{
			{Href: "/", Text: "Home"},
			{Href: "/docs", Text: "Docs"},
		},
		CurrentPath: "/docs",
		Sticky:      false,
		RightItems:  templ.Raw("<button>Sign in</button>"),
		BaseProps: utils.BaseProps{
			ID:        "snav",
			AriaLabel: "Top nav",
		},
	}))
	utils.AssertContainsAll(t, output,
		"MyApp",
		"Home",
		"Docs",
		"Sign in",
	)
}

func TestDefaultSimpleNavPropsRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SimpleNav(DefaultSimpleNavProps()))
	_ = output // should not panic
}

// ---------------------------------------------------------------------------
// SidebarNav: full props, icon items, footer slot, CurrentPath active
// ---------------------------------------------------------------------------

func TestSidebarNavFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SidebarNav(SidebarNavProps{
		Brand: templ.Raw("<div>Logo</div>"),
		Items: []SidebarNavItem{
			{Label: "Dashboard", Href: "/dashboard", Icon: icons.Home},
			{Label: "Settings", Href: "/settings", Icon: icons.Settings},
			{Label: "Users", Href: "/users", Icon: icons.Users},
		},
		Footer:      templ.Raw("<div>v1.0</div>"),
		CurrentPath: "/dashboard",
		BaseProps: utils.BaseProps{
			ID:        "sidebar",
			AriaLabel: "Main sidebar",
		},
	}))
	utils.AssertContainsAll(t, output,
		`id="sidebar"`,
		"Dashboard",
		"Settings",
		"Users",
		"v1.0",
	)
}

func TestSidebarNavExplicitActive(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SidebarNav(SidebarNavProps{
		Items: []SidebarNavItem{
			{Label: "A", Href: "/a", Active: true},
			{Label: "B", Href: "/b"},
		},
	}))
	utils.AssertContainsAll(t, output, "A", "B")
}

// ---------------------------------------------------------------------------
// LoadMore: cursor encoding, existing query params
// ---------------------------------------------------------------------------

func TestLoadMoreWithExistingQuery(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint: "/api/items?filter=active",
		Cursor:   "abc123",
		Label:    "Load more",
	}))
	utils.AssertContainsAll(t, output, "Load more", "abc123")
}

func TestLoadMoreDefaultLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint: "/api/items",
		Cursor:   "cur1",
	}))
	utils.AssertContains(t, output, "Load more")
}

func TestLoadMoreWithBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint: "/more",
		Cursor:   "c",
		BaseProps: utils.BaseProps{
			ID:        "lm-1",
			AriaLabel: "Load more items",
		},
	}))
	utils.AssertContainsAll(t, output, `id="lm-1"`, `aria-label="Load more items"`)
}

// ---------------------------------------------------------------------------
// Footer
// ---------------------------------------------------------------------------

func TestFooterRenderOutput(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Footer("Acme Inc."))
	utils.AssertContains(t, output, "Acme Inc.")
}

// ---------------------------------------------------------------------------
// MobileMenu + Toggle
// ---------------------------------------------------------------------------

func TestMobileMenuFullRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileMenu(
		[]NavLinkProps{
			{Href: "/", Text: "Home"},
			{Href: "/about", Text: "About"},
		},
		"/about",
		"nonce-abc",
		"mobile-menu-1",
	))
	utils.AssertContainsAll(t, output, "Home", "About", "mobile-menu-1")
}

func TestMobileMenuToggleShown(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileMenuToggle(true, "mm-1"))
	utils.AssertContains(t, output, "mm-1")
}

func TestMobileMenuToggleHidden(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileMenuToggle(false, "mm-2"))
	_ = output // should not panic when hidden
}
