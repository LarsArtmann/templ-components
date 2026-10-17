package visualtest

import (
	"testing"

	"github.com/chromedp/chromedp"
)

// TestDemoNavTimingBudget is the reproducible home of the T18 Lighthouse
// measurement: the same headless Chromium and the same in-process demo
// server grade the demo routes' real navigation metrics — no npx +
// CHROME_PATH dance, no machine-local /tmp JSON that dies with the session.
//
// Budgets are generous on purpose: an httptest server on localhost with
// warm caches should be orders of magnitude under them. The guard exists
// to catch regressions (a bloated stylesheet, a huge DOM, a late-loading
// font causing shifts), not to grade production hosting.
const (
	navBudgetFCPMS = 1500.0
	navBudgetLCPMS = 2500.0
	navBudgetCLS   = 0.1
	navBudgetLoad  = 5000.0
)

// navTimingMetrics is the payload the in-page PerformanceObserver probe
// resolves with (all values milliseconds unless named otherwise).
type navTimingMetrics struct {
	FCPMS  float64 `json:"fcp"`
	LCPMS  float64 `json:"lcp"`
	CLS    float64 `json:"cls"`
	LoadMS float64 `json:"load"`
}

// navTimingJS buffers LCP and layout-shift entries (Chromium-only entry
// types, which is fine — the visual module is Chromium-pinned) and resolves
// after a short settle so late entries land.
const navTimingJS = `new Promise(res => {
  let lcp = 0, cls = 0;
  try {
    new PerformanceObserver(list => {
      const es = list.getEntries();
      if (es.length) { lcp = es[es.length - 1].startTime; }
    }).observe({ type: 'largest-contentful-paint', buffered: true });
    new PerformanceObserver(list => {
      for (const e of list.getEntries()) { if (!e.hadRecentInput) { cls += e.value; } }
    }).observe({ type: 'layout-shift', buffered: true });
  } catch (e) {}
  setTimeout(() => {
    const nav = performance.getEntriesByType('navigation')[0];
    const fcpEntry = performance.getEntriesByType('paint')
      .find(p => p.name === 'first-contentful-paint');
    res({ fcp: fcpEntry ? fcpEntry.startTime : -1, lcp: lcp, cls: cls,
          load: nav ? nav.loadEventEnd : -1 });
  }, 600);
})`

// TestDemoNavTimingBudget audits a representative route set: the home (the
// widest composition), the two heaviest pages (display, forms), and the
// dual-runtime kanban page (htmx + Datastar both loaded).
func TestDemoNavTimingBudget(t *testing.T) {
	t.Parallel()

	server := StartDemoServer(t)

	routes := []struct {
		name string
		path string
	}{
		{name: "index", path: "/"},
		{name: "display", path: "/display"},
		{name: "forms", path: "/forms"},
		{name: "kanban", path: "/kanban"},
	}

	for _, route := range routes {
		t.Run(route.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := newTab(t)
			defer cancel()

			if err := chromedp.Do(ctx,
				chromedp.Navigate(server.BaseURL()+route.path),
				chromedp.WaitReady("body"),
			); err != nil {
				t.Fatalf("visualtest[nav-timing][%s]: navigate: %v", route.name, err)
			}

			metrics, err := chromedp.Run(ctx, chromedp.Evaluate[navTimingMetrics](
				navTimingJS,
				chromedp.EvalAwaitPromise,
			))
			if err != nil {
				t.Fatalf("visualtest[nav-timing][%s]: %v", route.name, err)
			}

			checks := []struct {
				name   string
				value  float64
				budget float64
			}{
				{"first-contentful-paint", metrics.FCPMS, navBudgetFCPMS},
				{"largest-contentful-paint", metrics.LCPMS, navBudgetLCPMS},
				{"cumulative-layout-shift", metrics.CLS, navBudgetCLS},
				{"load-event-end", metrics.LoadMS, navBudgetLoad},
			}

			for _, check := range checks {
				if check.value < 0 {
					t.Errorf(
						"visualtest[nav-timing][%s]: %s missing (entry type unsupported or page never emitted it)",
						route.name,
						check.name,
					)

					continue
				}

				if check.value > check.budget {
					t.Errorf(
						"visualtest[nav-timing][%s]: %s = %.0fms exceeds budget %.0fms — inspect before raising (regressions first, budgets second)",
						route.name,
						check.name,
						check.value,
						check.budget,
					)
				}
			}

			t.Logf(
				"nav-timing[%s]: FCP=%.0fms LCP=%.0fms CLS=%.3f load=%.0fms",
				route.name, metrics.FCPMS, metrics.LCPMS, metrics.CLS, metrics.LoadMS,
			)
		})
	}
}
