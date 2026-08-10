package navigation

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestGoldenSidebarNav(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SidebarNav(SidebarNavProps{
		Brand: templ.Raw(`<span class="font-bold text-white">MyApp</span>`),
		Items: []SidebarNavItem{
			{Label: "Dashboard", Href: "/", Icon: icons.Squares2x2, Active: true},
			{Label: "Users", Href: "/users", Icon: icons.Users},
		},
	}))
	golden.Assert(t, "sidebar_nav", output)
}

// TestGoldenPagination snapshots the pagination component in a representative
// middle-of-range state (current 5 of 20, ellipsis + rel=canonical). This is
// the proof-of-concept for converting navigation assertion tests to golden
// files (T14): the snapshot catches structural/visual regressions that the
// behavioural assertions in pagination_test.go (page-range math, edge cases)
// do not. The two complement each other — see the project's testing checklist.
func TestGoldenPagination(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Pagination(PaginationProps{
		CurrentPage: 5,
		TotalPages:  20,
		BaseURL:     "/items",
	}))
	golden.Assert(t, "pagination", output)
}

// TestGoldenBreadcrumbs snapshots a three-level breadcrumb trail with JSON-LD
// structured data enabled — the richest render path.
func TestGoldenBreadcrumbs(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
		Items: []BreadcrumbItem{
			{Text: "Home", Href: "https://example.com/"},
			{Text: "Settings", Href: "https://example.com/settings"},
			{Text: "Profile"},
		},
		JSONLD:  true,
		BaseURL: "https://example.com",
	}))
	golden.Assert(t, "breadcrumbs", output)
}
