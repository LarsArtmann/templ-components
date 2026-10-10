package forms

import (
	"net/url"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestFilterChipsRender(t *testing.T) {
	t.Parallel()

	chips := []FilterChip{
		{Label: "All", Href: "/history"},
		{Label: "Uploads", Href: "/history?source=upload", Active: true},
		{Label: "Fetches", Href: "/history?source=fetch"},
	}

	t.Run("renders links with zero JS", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterChips(FilterChipsProps{Chips: chips}))
		utils.AssertContains(t, output, `href="/history?source=upload"`)
		utils.AssertContains(t, output, ">Uploads</a>")
		utils.AssertNotContains(t, output, "<script")
	})

	t.Run("active chip carries aria-current and filled style", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterChips(FilterChipsProps{Chips: chips}))
		utils.AssertContains(t, output, `aria-current="true"`)
		utils.AssertContains(t, output, "bg-blue-600")
	})

	t.Run("inactive chips have no aria-current", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, FilterChips(FilterChipsProps{Chips: chips}))
		if got, want := strings.Count(output, `aria-current="true"`), 1; got != want {
			t.Errorf("aria-current count = %d, want %d", got, want)
		}
	})

	t.Run("group is labelled with fallback", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterChips(FilterChipsProps{Chips: chips}))
		utils.AssertContains(t, output, `role="group"`)
		utils.AssertContains(t, output, `aria-label="Filters"`)
	})

	t.Run("aria label override", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterChips(FilterChipsProps{
			AriaLabel: "Source filter",
			Chips:     chips,
		}))
		utils.AssertContains(t, output, `aria-label="Source filter"`)
	})

	t.Run("propagates base props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterChips(FilterChipsProps{
			ID: "chips", Class: "mb-2", Attrs: templ.Attributes{"data-x": "1"},
			Chips: chips,
		}))
		utils.AssertContains(t, output, `id="chips"`)
		utils.AssertContains(t, output, "mb-2")
		utils.AssertContains(t, output, `data-x="1"`)
	})

	t.Run("empty chip list renders empty group", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterChips(FilterChipsProps{}))
		utils.AssertContains(t, output, `role="group"`)
		utils.AssertNotContains(t, output, "<a")
	})

	t.Run("motion-reduce fallback present", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterChips(FilterChipsProps{Chips: chips}))
		utils.AssertContains(t, output, "motion-reduce:transition-none")
	})
}

func TestFilterToggleHref(t *testing.T) {
	t.Parallel()

	t.Run("sets the key", func(t *testing.T) {
		t.Parallel()

		got := FilterToggleHref("/history", url.Values{}, "source", "upload")
		if got != "/history?source=upload" {
			t.Errorf("FilterToggleHref = %q, want /history?source=upload", got)
		}
	})

	t.Run("empty value removes the key", func(t *testing.T) {
		t.Parallel()

		got := FilterToggleHref("/history", url.Values{"source": {"upload"}}, "source", "")
		if got != "/history" {
			t.Errorf("FilterToggleHref = %q, want /history", got)
		}
	})

	t.Run("preserves sibling params", func(t *testing.T) {
		t.Parallel()

		got := FilterToggleHref("/history", url.Values{"lane": {"main"}, "source": {"upload"}}, "source", "fetch")

		want := "/history?lane=main&source=fetch"
		if got != want {
			t.Errorf("FilterToggleHref = %q, want %q", got, want)
		}
	})

	t.Run("escapes values", func(t *testing.T) {
		t.Parallel()

		got := FilterToggleHref("/history", url.Values{}, "q", "a b&c=d")
		if got != "/history?q=a+b%26c%3Dd" {
			t.Errorf("FilterToggleHref = %q, want escaped", got)
		}
	})

	t.Run("does not mutate the caller's values", func(t *testing.T) {
		t.Parallel()

		params := url.Values{"source": {"upload"}}

		_ = FilterToggleHref("/history", params, "source", "")
		if params.Get("source") != "upload" {
			t.Error("caller's url.Values mutated")
		}
	})
}

func TestGoldenSweepFilterChips(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "filter_chips_row", HTML: utils.Render(t, FilterChips(FilterChipsProps{
			Chips: []FilterChip{
				{Label: "All", Href: "/history"},
				{Label: "Uploads", Href: "/history?source=upload", Active: true},
				{Label: "Fetches", Href: "/history?source=fetch"},
				{Label: "Body", Href: "/history?source=body"},
			},
		}))},
		{Name: "filter_chips_empty", HTML: utils.Render(t, FilterChips(FilterChipsProps{}))},
	})
}
