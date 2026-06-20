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

	t.Run("edge cases", func(t *testing.T) {
		t.Parallel()
		for _, tt := range []struct {
			name        string
			currentPage uint
			totalPages  uint
			baseURL     string
			want        string
		}{
			{"first page with many pages", 1, 15, "/p", "/p?page=2"},
			{"last page with many pages", 15, 15, "/p", "15"},
			{"prev next on first page", 1, 5, "/p", `rel="next"`},
			{"prev next on last page", 5, 5, "/p", `rel="prev"`},
		} {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				output := utils.Render(t, Pagination(PaginationProps{
					CurrentPage: tt.currentPage,
					TotalPages:  tt.totalPages,
					BaseURL:     tt.baseURL,
				}))
				utils.AssertContains(t, output, tt.want)
			})
		}
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
