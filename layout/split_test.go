package layout

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestSplit(t *testing.T) {
	t.Parallel()

	t.Run("default props produce grid + md:grid-cols-3 + gap-6", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Main:  templ.Raw(`<div data-test="main">M</div>`),
			Aside: templ.Raw(`<div data-test="aside">A</div>`),
		}))
		utils.AssertContainsAll(t, output,
			"grid", "grid-cols-1", "md:grid-cols-3", "gap-6",
			`data-test="main"`, `data-test="aside"`,
		)
	})

	t.Run("aside at end (default) — main comes first in source order", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Main:  templ.Raw(`<div id="m">M</div>`),
			Aside: templ.Raw(`<div id="a">A</div>`),
		}))
		mainIdx := strings.Index(output, `id="m"`)

		asideIdx := strings.Index(output, `id="a"`)
		if mainIdx < 0 || asideIdx < 0 {
			t.Fatalf("missing main or aside in output: %q", output)
		}

		if mainIdx > asideIdx {
			t.Errorf("AsidePositionEnd: main should come before aside in source order")
		}
	})

	t.Run("aside at start — aside comes first in source order", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Main:          templ.Raw(`<div id="m">M</div>`),
			Aside:         templ.Raw(`<div id="a">A</div>`),
			AsidePosition: AsidePositionStart,
		}))
		mainIdx := strings.Index(output, `id="m"`)

		asideIdx := strings.Index(output, `id="a"`)
		if mainIdx < 0 || asideIdx < 0 {
			t.Fatalf("missing main or aside in output: %q", output)
		}

		if asideIdx > mainIdx {
			t.Errorf("AsidePositionStart: aside should come before main in source order")
		}
	})

	t.Run("ratio 1To4 produces md:grid-cols-4 + main span-3", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Main:  templ.Raw(`<div>M</div>`),
			Aside: templ.Raw(`<div>A</div>`),
			Ratio: SplitRatio1To4,
		}))
		utils.AssertContainsAll(t, output, "md:grid-cols-4", "md:col-span-3")
	})

	t.Run("ratio 1To2 produces md:grid-cols-2 + main span-1", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Main:  templ.Raw(`<div>M</div>`),
			Aside: templ.Raw(`<div>A</div>`),
			Ratio: SplitRatio1To2,
		}))
		utils.AssertContainsAll(t, output, "md:grid-cols-2", "md:col-span-1")
	})

	t.Run("unknown ratio falls back to default (1To3)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Main:  templ.Raw(`<div>M</div>`),
			Aside: templ.Raw(`<div>A</div>`),
			Ratio: SplitRatio("bogus"),
		}))
		utils.AssertContainsAll(t, output, "md:grid-cols-3", "md:col-span-2")
	})

	t.Run("custom Gap propagates", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Main:  templ.Raw(`<div>M</div>`),
			Aside: templ.Raw(`<div>A</div>`),
			Gap:   "gap-8",
		}))
		utils.AssertContains(t, output, "gap-8")
	})

	t.Run("BaseProps propagate", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			BaseProps: utils.BaseProps{
				ID:        "split-1",
				Class:     "data-tc-test",
				AriaLabel: "Article and TOC",
				Attrs:     templ.Attributes{"data-testid": "split"},
			},
			Main:  templ.Raw(`<div>M</div>`),
			Aside: templ.Raw(`<div>A</div>`),
		}))
		utils.AssertContainsAll(t, output,
			`id="split-1"`, "data-tc-test", `aria-label="Article and TOC"`, `data-testid="split"`,
		)
	})

	t.Run("nil Aside does not panic", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Main: templ.Raw(`<div data-test="main">M</div>`),
		}))
		utils.AssertContains(t, output, `data-test="main"`)

		if strings.Contains(output, "<aside") {
			t.Errorf("nil Aside should not render <aside>")
		}
	})

	t.Run("nil Main does not panic", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Aside: templ.Raw(`<div>A</div>`),
		}))
		utils.AssertContainsAll(t, output, "<aside", "grid")
	})

	t.Run("aside uses <aside> element; main uses <div> (Base owns <main>)", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Split(SplitProps{
			Main:  templ.Raw(`<div>M</div>`),
			Aside: templ.Raw(`<div>A</div>`),
		}))
		if !strings.Contains(output, "<aside") {
			t.Errorf("Split Aside column must use <aside> element")
		}
		// Must NOT emit <main> — Base already provides the main landmark.
		// Nested <main> is invalid HTML (only one per page).
		if strings.Contains(output, "<main") {
			t.Errorf("Split must not emit <main> — layout.Base owns the main landmark")
		}
	})

	t.Run("min-w-0 present on both columns (grid-blowout guard)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Split(SplitProps{
			Main:  templ.Raw(`<div>M</div>`),
			Aside: templ.Raw(`<div>A</div>`),
		}))
		// min-w-0 complements minmax(0,1fr) to prevent flex/grid overflow
		// from wide children. At least two occurrences (one per column).
		count := strings.Count(output, "min-w-0")
		if count < 2 {
			t.Errorf("expected min-w-0 on both columns; got %d occurrences", count)
		}
	})
}

func TestSplitRatioIsValid(t *testing.T) {
	t.Parallel()

	for _, r := range []SplitRatio{SplitRatio1To2, SplitRatio1To3, SplitRatio1To4} {
		if !SplitRatioIsValid(r) {
			t.Errorf("SplitRatioIsValid(%q) = false; want true", r)
		}
	}

	if SplitRatioIsValid(SplitRatio("bogus")) {
		t.Errorf("SplitRatioIsValid(\"bogus\") = true; want false")
	}
}

func TestAsidePositionIsValid(t *testing.T) {
	t.Parallel()

	if !AsidePositionIsValid(AsidePositionStart) {
		t.Errorf("AsidePositionIsValid(Start) = false; want true")
	}

	if !AsidePositionIsValid(AsidePositionEnd) {
		t.Errorf("AsidePositionIsValid(End) = false; want true")
	}

	if AsidePositionIsValid(AsidePosition("middle")) {
		t.Errorf("AsidePositionIsValid(\"middle\") = true; want false")
	}
}
