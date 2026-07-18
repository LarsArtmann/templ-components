// Package navigation provides tests for navigation components.
package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestPaginationRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		current    uint
		total      uint
		maxVisible uint
		wantStart  uint
		wantEnd    uint
	}{
		{"small total", 2, 3, 5, 1, 3},
		{"current in middle", 5, 10, 5, 3, 7},
		{"current near start", 2, 10, 5, 1, 5},
		{"current near end", 9, 10, 5, 6, 10},
		{"default max visible", 5, 20, 0, 3, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			start, end := paginationRange(tt.current, tt.total, tt.maxVisible)
			if start != tt.wantStart || end != tt.wantEnd {
				t.Errorf("paginationRange(%d, %d, %d) = (%d, %d), want (%d, %d)",
					tt.current, tt.total, tt.maxVisible, start, end, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestPaginationRender(t *testing.T) {
	t.Parallel()
	t.Run("renders page numbers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  5,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, `href="/items?page=1"`)
		utils.AssertContains(t, output, `href="/items?page=3"`)
		utils.AssertContains(t, output, `aria-current="page"`)
		utils.AssertContains(t, output, "2")
		utils.AssertContains(t, output, "Previous")
		utils.AssertContains(t, output, "Next")
		utils.AssertContains(t, output, "Showing page")
		utils.AssertContains(t, output, "of")
	})

	t.Run("first page disables previous", func(t *testing.T) {
		t.Parallel()
		output := renderPagination(t, 1, 3, "/users")
		utils.AssertContains(t, output, `aria-disabled="true"`)
		utils.AssertNotContains(t, output, `href="/users?page=0"`)
	})

	t.Run("last page disables next", func(t *testing.T) {
		t.Parallel()
		output := renderPagination(t, 3, 3, "/users")
		utils.AssertNotContains(t, output, `href="/users?page=4"`)
	})

	t.Run("hidden when single page", func(t *testing.T) {
		t.Parallel()
		assertEmptyPaginationOutput(t, 1, 1, "/users")
	})

	t.Run("custom query param", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  3,
			BaseURL:     "/search?q=test",
			QueryParam:  "p",
		}))
		utils.AssertContains(t, output, `href="/search?p=1&amp;q=test"`)
	})

	t.Run("zero CurrentPage clamped to 1", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 0,
			TotalPages:  3,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, `aria-current="page"`)
		utils.AssertContains(t, output, "1")
	})

	t.Run("zero TotalPages renders nothing", func(t *testing.T) {
		t.Parallel()
		assertEmptyPaginationOutput(t, 1, 0, "/items")
	})

	t.Run("custom ID propagated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			BaseProps:   utils.BaseProps{ID: "my-pager"},
			CurrentPage: 2,
			TotalPages:  3,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, `id="my-pager"`)
	})

	t.Run("rel=canonical on first page when ellipsis shown", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 5,
			TotalPages:  20,
			BaseURL:     "/items",
		}))
		utils.AssertContains(t, output, `rel="canonical"`)
	})
}

func renderPagination(t *testing.T, currentPage, totalPages uint, baseURL string) string {
	t.Helper()

	return utils.Render(t, Pagination(PaginationProps{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		BaseURL:     baseURL,
	}))
}

func assertEmptyPaginationOutput(t *testing.T, currentPage, totalPages uint, baseURL string) {
	t.Helper()

	output := renderPagination(t, currentPage, totalPages, baseURL)
	if output != "" {
		t.Errorf("expected empty output, got: %s", output)
	}
}
