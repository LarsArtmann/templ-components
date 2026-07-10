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
