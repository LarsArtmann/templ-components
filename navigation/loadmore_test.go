package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils/wire"
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

func TestGoldenLoadMoreWired(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "loadmore_wired_htmx", HTML: utils.Render(t, LoadMore(LoadMoreProps{
			BaseProps: utils.BaseProps{ID: "items-more"},
			Wire:      &wire.Action{URL: "/api/items"},
			Cursor:    "next",
		}))},
		{Name: "loadmore_wired_datastar", HTML: utils.Render(t, LoadMore(LoadMoreProps{
			BaseProps: utils.BaseProps{ID: "items-more"},
			Wire:      &wire.Action{Transport: wire.TransportDatastar, URL: "/api/items"},
			Cursor:    "next",
		}))},
		{Name: "loadmore_wired_datastar_infinite_scroll_ignored", HTML: utils.Render(t, LoadMore(LoadMoreProps{
			BaseProps:      utils.BaseProps{ID: "items-more"},
			Wire:           &wire.Action{Transport: wire.TransportDatastar, URL: "/api/items"},
			InfiniteScroll: true,
		}))},
		{Name: "loadmore_wired_empty_url_inert", HTML: utils.Render(t, LoadMore(LoadMoreProps{
			BaseProps: utils.BaseProps{ID: "items-more"},
			Wire:      &wire.Action{},
		}))},
	})
}

func TestLoadMoreWire(t *testing.T) {
	t.Parallel()

	t.Run("htmx wire keeps self-replacement and appends cursor to Wire.URL", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			BaseProps: utils.BaseProps{ID: "items-more"},
			Wire:      &wire.Action{URL: "/api/items"},
			Cursor:    "next",
		}))
		utils.AssertContains(t, output, `hx-get="/api/items?cursor=next"`)
		utils.AssertContains(t, output, `hx-swap="outerHTML"`)
		utils.AssertContains(t, output, `hx-target="this"`)
		utils.AssertNotContains(t, output, "data-on:")
	})

	t.Run("datastar wire renders the expression and never a target", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			BaseProps: utils.BaseProps{ID: "items-more"},
			Wire:      &wire.Action{Transport: wire.TransportDatastar, URL: "/api/items"},
			Cursor:    "next",
		}))
		utils.AssertContains(t, output, `data-on:click="@get(&#39;/api/items?cursor=next&#39;)"`)
		utils.AssertNotContains(t, output, "hx-")
	})

	t.Run("InfiniteScroll is ignored under datastar", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			BaseProps:      utils.BaseProps{ID: "items-more"},
			Wire:           &wire.Action{Transport: wire.TransportDatastar, URL: "/api/items"},
			InfiniteScroll: true,
		}))
		utils.AssertNotContains(t, output, "revealed")
		utils.AssertNotContains(t, output, "hx-trigger")
	})

	t.Run("InfiniteScroll survives the htmx wire path", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			BaseProps:      utils.BaseProps{ID: "items-more"},
			Wire:           &wire.Action{URL: "/api/items"},
			InfiniteScroll: true,
		}))
		utils.AssertContains(t, output, `hx-trigger="revealed"`)
	})

	t.Run("empty Wire URL wires nothing", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			BaseProps: utils.BaseProps{ID: "items-more"},
			Wire:      &wire.Action{},
		}))
		utils.AssertNotContains(t, output, "hx-get")
		utils.AssertNotContains(t, output, "data-on:")
		utils.AssertNotContains(t, output, "hx-swap")
	})
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

	t.Run("adds hx-trigger=revealed when InfiniteScroll is true", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			Endpoint:       "/api/items",
			Cursor:         "next",
			InfiniteScroll: true,
		}))
		utils.AssertContains(t, output, `hx-trigger="revealed"`)
	})

	t.Run("omits hx-trigger by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadMore(LoadMoreProps{
			Endpoint: "/api/items",
		}))
		utils.AssertNotContains(t, output, "hx-trigger")
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
