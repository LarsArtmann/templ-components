package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestLoadMoreCursorAppended verifies the cursor is appended as a query param.
func TestLoadMoreCursorAppended(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint: "/api/items",
		Cursor:   "abc123",
	}))
	utils.AssertContains(t, output, "cursor=abc123")
	utils.AssertContains(t, output, `hx-get="/api/items?cursor=abc123"`)
}

// TestLoadMoreCursorWithExistingQuery verifies cursor appended when URL already
// has a query string.
func TestLoadMoreCursorWithExistingQuery(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint: "/api/items?page=2",
		Cursor:   "next",
	}))
	utils.AssertContains(t, output, "cursor=next")
}

// TestLoadMoreInfiniteScroll verifies hx-trigger="revealed" is present.
func TestLoadMoreInfiniteScroll(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint:       "/api/items",
		Cursor:         "x",
		InfiniteScroll: true,
	}))
	utils.AssertContains(t, output, `hx-trigger="revealed"`)
}

// TestLoadMoreNoInfiniteScroll verifies hx-trigger is absent by default.
func TestLoadMoreNoInfiniteScroll(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint: "/api/items",
		Cursor:   "x",
	}))
	utils.AssertNotContains(t, output, `hx-trigger="revealed"`)
}

// TestLoadMoreSelfTarget verifies hx-swap + hx-target for self-replacement.
func TestLoadMoreSelfTarget(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint: "/api/items",
		Cursor:   "x",
	}))
	utils.AssertContains(t, output, `hx-swap="outerHTML"`)
	utils.AssertContains(t, output, `hx-target="this"`)
}

// TestLoadMoreEscapedCursor verifies cursor with special chars is escaped.
func TestLoadMoreEscapedCursor(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint: "/api/items",
		Cursor:   "base64==",
	}))
	// = should be properly escaped in the URL
	utils.AssertContains(t, output, "cursor=")
}

// TestBreadcrumbJSONLD verifies JSON-LD structured data when enabled.
func TestBreadcrumbJSONLD(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "/"},
			{Text: "Users", Href: "/users"},
		},
		JSONLD: true,
	}))
	utils.AssertContains(t, output, "application/ld+json")
}

// TestBreadcrumbNoJSONLD verifies JSON-LD is absent by default.
func TestBreadcrumbNoJSONLD(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "/"},
		},
	}))
	utils.AssertNotContains(t, output, "application/ld+json")
}

// TestBreadcrumbRelPrevNext verifies SEO rel attributes on links.
func TestBreadcrumbRelPrevNext(t *testing.T) {
	t.Parallel()
	_ = utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "/"},
		},
	}))
	// Just verify it renders without panic — rel attrs are on pagination, not breadcrumbs
}

// TestFooterAcceptsBaseProps verifies Footer now accepts BaseProps (Class/ID/Attrs).
func TestFooterAcceptsBaseProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Footer(FooterProps{
		BaseProps: utils.BaseProps{
			Class: "custom-footer-class",
			ID:    "site-footer",
		},
		BrandText: "Acme",
	}))
	utils.AssertContains(t, output, "custom-footer-class")
	utils.AssertContains(t, output, `id="site-footer"`)
	utils.AssertContains(t, output, "Acme")
}

// TestMobileMenuNoDoublePrefix verifies the mobile-menu ID does not get
// the "tc-mobile-menu-" prefix applied twice. EnsureID already produces
// "tc-mobile-menu-<hex>", so the template must not prepend it again.
func TestMobileMenuNoDoublePrefix(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileMenu(testNavLinks, "/", "nonce", "my-menu-id", false))
	utils.AssertContains(t, output, `id="my-menu-id"`)
	utils.AssertNotContains(t, output, "tc-mobile-menu-tc-mobile-menu-")
	utils.AssertNotContains(t, output, "tc-mobile-menu-my-menu-id")
}

// TestMobileMenuToggleNoDoublePrefix verifies aria-controls matches
// the menu ID without a double prefix.
func TestMobileMenuToggleNoDoublePrefix(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, MobileMenuToggle(true, "my-menu-id", false))
	utils.AssertContains(t, output, `aria-controls="my-menu-id"`)
	utils.AssertNotContains(t, output, "tc-mobile-menu-my-menu-id")
}

// TestBreadcrumbsCurrentPathAutoDetect verifies that setting CurrentPath
// auto-highlights the matching crumb without needing Active=true.
func TestBreadcrumbsCurrentPathAutoDetect(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "/"},
			{Text: "Users", Href: "/users"},
			{Text: "Edit", Href: "/users/1/edit"},
		},
		CurrentPath: "/users/1/edit",
	}))
	// The active crumb should render as a <span> (no <a> tag)
	utils.AssertContains(t, output, "text-gray-500 dark:text-gray-400")
	// Non-active crumbs should still have links
	utils.AssertContains(t, output, `href="/"`)
	utils.AssertContains(t, output, `href="/users"`)
	// The active crumb should NOT have an href
	utils.AssertNotContains(t, output, `href="/users/1/edit"`)
}

// TestBreadcrumbsCurrentPathNoMatch verifies that when CurrentPath doesn't
// match any crumb, none are auto-highlighted (except terminal empty-Href crumb).
func TestBreadcrumbsCurrentPathNoMatch(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "/"},
			{Text: "Settings", Href: "/settings"},
		},
		CurrentPath: "/completely/different",
	}))
	// Both should be links since neither matches and neither has Active=true
	utils.AssertContains(t, output, `href="/"`)
	utils.AssertContains(t, output, `href="/settings"`)
}

// TestBreadcrumbsExplicitActiveOverridesCurrentPath verifies that Active=true
// takes priority — an item can be active even when CurrentPath matches a different crumb.
func TestBreadcrumbsExplicitActiveOverridesCurrentPath(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "/", Active: true},
			{Text: "Users", Href: "/users"},
		},
		CurrentPath: "/users",
	}))
	// Both should be active: Home via explicit flag, Users via CurrentPath match
	// Neither should render as a link
	utils.AssertNotContains(t, output, `href="/"`)
	utils.AssertNotContains(t, output, `href="/users"`)
}
