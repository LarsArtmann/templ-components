package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestPaginationEllipsisCoverage(t *testing.T) {
	t.Parallel()

	t.Run("many pages shows ellipsis", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 5,
			TotalPages:  20,
			BaseURL:     "/items", QueryParam: "page",
		}))

		utils.AssertContains(t, output, "/items?page=1")
		utils.AssertContains(t, output, "/items?page=20")
	})

	t.Run("first page with many pages", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 1,
			TotalPages:  15,
			BaseURL:     "/p",
		}))
		utils.AssertContains(t, output, "/p?page=2")
	})

	t.Run("last page with many pages", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 15,
			TotalPages:  15,
			BaseURL:     "/p",
		}))

		utils.AssertContains(t, output, "15")
	})

	t.Run("prev next arrows on first page", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 1,
			TotalPages:  5,
			BaseURL:     "/p",
		}))
		utils.AssertContains(t, output, `rel="next"`)
	})

	t.Run("prev next arrows on last page", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 5,
			TotalPages:  5,
			BaseURL:     "/p",
		}))
		utils.AssertContains(t, output, `rel="prev"`)
	})

	t.Run("with ID and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  5,
			BaseURL:     "/p",
			BaseProps:   utils.BaseProps{ID: "pager", Class: "mt-4"},
		}))
		utils.AssertContains(t, output, `id="pager"`)
		utils.AssertContains(t, output, "mt-4")
	})
}

func TestBreadcrumbsExtraCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with custom separator", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
			Items: []BreadcrumbItem{
				{Text: "Home", Href: "/"},
				{Text: "Users", Href: "/users"},
				{Text: "Profile", Active: true},
			},
			Separator: "→",
		}))
		utils.AssertContains(t, output, "→")
		utils.AssertContains(t, output, "Profile")
	})

	t.Run("with JSON-LD", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
			Items: []BreadcrumbItem{
				{Text: "Home", Href: "https://example.com/"},
				{Text: "Page", Href: "https://example.com/page"},
			},
			JSONLD: true,
		}))
		utils.AssertContains(t, output, "application/ld+json")
		utils.AssertContains(t, output, "BreadcrumbList")
	})

	t.Run("with ID", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Breadcrumbs(BreadcrumbsProps{
			Items:     []BreadcrumbItem{{Text: "Home", Href: "/"}},
			BaseProps: utils.BaseProps{ID: "bc"},
		}))
		utils.AssertContains(t, output, `id="bc"`)
	})
}

func TestNavLinkExtraCoverage(t *testing.T) {
	t.Parallel()

	t.Run("external link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NavLink(NavLinkProps{
			Href:     "https://external.com",
			Text:     "External",
			External: true,
		}, ""))
		utils.AssertContains(t, output, "target=\"_blank\"")
		utils.AssertContains(t, output, "rel=\"noopener noreferrer\"")
	})

	t.Run("active link with aria-current", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NavLink(NavLinkProps{
			Href: "/about",
			Text: "About",
		}, "/about"))
		utils.AssertContains(t, output, `aria-current="page"`)
	})
}
