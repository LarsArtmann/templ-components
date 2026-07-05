package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenLoadMore(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, LoadMore(LoadMoreProps{
		Endpoint: "/api/items",
		Cursor:   "abc123",
		BaseProps: utils.BaseProps{
			ID: "test-load-more",
		},
	}))
	golden.Assert(t, "loadmore", output)
}

func TestLoadMoreBehavior(t *testing.T) {
	t.Parallel()

	t.Run("renders button with hx-get", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			Endpoint: "/api/items",
			Cursor:   "next",
		}))
		utils.AssertContains(t, output, `hx-get="/api/items?cursor=next"`)
		utils.AssertContains(t, output, `hx-swap="outerHTML"`)
		utils.AssertContains(t, output, "Load more")
	})

	t.Run("appends cursor with ampersand when endpoint has query", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			Endpoint: "/api/items?filter=active",
			Cursor:   "abc",
		}))
		utils.AssertContainsAll(t, output, `hx-get="/api/items?cursor=abc&amp;filter=active"`)
	})

	t.Run("uses custom label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			Endpoint: "/x",
			Label:    "Show more results",
		}))
		utils.AssertContains(t, output, "Show more results")
	})

	t.Run("renders without cursor", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			Endpoint: "/api/items",
		}))
		utils.AssertContains(t, output, `hx-get="/api/items"`)
		utils.AssertNotContains(t, output, "cursor=")
	})
}

func TestLoadMoreA11y(t *testing.T) {
	t.Parallel()

	t.Run("button has type=button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{Endpoint: "/x"}))
		utils.AssertContains(t, output, `type="button"`)
	})

	t.Run("has focus-visible ring", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{Endpoint: "/x"}))
		utils.AssertContains(t, output, "focus-visible:ring-2")
	})

	t.Run("has motion-reduce classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{Endpoint: "/x"}))
		utils.AssertContains(t, output, "motion-reduce:transition-none")
		utils.AssertContains(t, output, "motion-reduce:duration-0")
	})

	t.Run("propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			Endpoint:  "/x",
			BaseProps: utils.BaseProps{AriaLabel: "Load more results"},
		}))
		utils.AssertContains(t, output, `aria-label="Load more results"`)
	})
}
