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
		current    int
		total      int
		maxVisible int
		wantStart  int
		wantEnd    int
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
			BaseURL:     "/users",
		}))
		utils.AssertContains(t, output, `href="/users?page=1"`)
		utils.AssertContains(t, output, `href="/users?page=3"`)
		utils.AssertContains(t, output, `aria-current="page"`)
		utils.AssertContains(t, output, "2")
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
		output := renderPagination(t, 1, 1, "/users")
		if output != "" {
			t.Errorf("expected empty output for single page, got: %s", output)
		}
	})

	t.Run("custom query param", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Pagination(PaginationProps{
			CurrentPage: 2,
			TotalPages:  3,
			BaseURL:     "/search?q=test",
			QueryParam:  "p",
		}))
		utils.AssertContains(t, output, `href="/search?q=test&amp;p=1"`)
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
}

func renderPagination(t *testing.T, currentPage, totalPages int, baseURL string) string {
	t.Helper()
	return utils.Render(t, Pagination(PaginationProps{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		BaseURL:     baseURL,
	}))
}
