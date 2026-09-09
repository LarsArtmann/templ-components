package visualtest

import (
	"testing"

	"github.com/chromedp/chromedp"
)

// demoMobileViewport is a small-phone viewport (iPhone 12/13/14 class).
const demoMobileViewport = 375

// TestDemoMobile375NoHorizontalOverflow proves the demo pages keep
// phone-width layout honest: no route may push the document wider than the
// viewport. Wide content must live in its own scroll container (kanban
// board, data table), never widen the page — the #1 untested form-factor
// defect class.
func TestDemoMobile375NoHorizontalOverflow(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.EmulateViewport(demoMobileViewport, 812)); err != nil {
		t.Fatalf("visualtest[demo]: emulate %dpx viewport: %v", demoMobileViewport, err)
	}

	for _, route := range []string{"/", "/forms", "/users", "/recipes/dashboard"} {
		t.Run(route, func(t *testing.T) {
			if err := chromedp.Run(ctx,
				chromedp.Navigate(server.BaseURL()+route),
				chromedp.WaitReady("body"),
			); err != nil {
				t.Fatalf("visualtest[demo]: navigate %s: %v", route, err)
			}

			var overflow int64

			if err := chromedp.Run(ctx, chromedp.Evaluate(
				`document.documentElement.scrollWidth - document.documentElement.clientWidth`, &overflow,
			)); err != nil {
				t.Fatalf("visualtest[demo]: measure %s overflow: %v", route, err)
			}

			if overflow > 1 {
				t.Errorf("visualtest[demo]: %s overflows a %dpx viewport by %dpx — find the element widening the page, not a scroll container", route, demoMobileViewport, overflow)
			}
		})
	}

	server.FailIfServerErrors(t)
}

// TestDemoMobile375KanbanReachable proves the kanban board stays usable at
// phone width: cards render and the board scrolls horizontally (content
// reachable by panning), instead of being clipped or widening the page.
func TestDemoMobile375KanbanReachable(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.EmulateViewport(demoMobileViewport, 812)); err != nil {
		t.Fatalf("visualtest[demo]: emulate %dpx viewport: %v", demoMobileViewport, err)
	}

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.BaseURL()+"/"),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[demo]: navigate: %v", err)
	}

	var cardCount int64

	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`document.querySelectorAll('#kanban-demo-htmx [data-tc-kanban-card]').length`, &cardCount,
	)); err != nil {
		t.Fatalf("visualtest[demo]: count kanban cards: %v", err)
	}

	if cardCount == 0 {
		t.Fatal("visualtest[demo]: htmx kanban board rendered no cards at phone width")
	}

	var boardScroll struct {
		Scroll int64
		Client int64
	}

	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`(() => { const el = document.querySelector('#kanban-demo-htmx .overflow-x-auto'); return {Scroll: el ? el.scrollWidth : 0, Client: el ? el.clientWidth : 0}; })()`,
		&boardScroll,
	)); err != nil {
		t.Fatalf("visualtest[demo]: measure kanban scroll container: %v", err)
	}

	if boardScroll.Client == 0 {
		t.Fatal("visualtest[demo]: kanban horizontal scroll container missing at phone width")
	}

	var pageOverflow int64

	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`document.documentElement.scrollWidth - document.documentElement.clientWidth`, &pageOverflow,
	)); err != nil {
		t.Fatalf("visualtest[demo]: measure page overflow: %v", err)
	}

	if pageOverflow > 1 {
		t.Errorf("visualtest[demo]: kanban board widens the page by %dpx at %dpx width — the board must scroll internally", pageOverflow, demoMobileViewport)
	}

	server.FailIfServerErrors(t)
}
