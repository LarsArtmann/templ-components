package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/htmx"
	"github.com/larsartmann/templ-components/layout"
)

// The PolledRegion busy-cue E2E pins the browser-level contract the string
// tests cannot: the singleton script must actually LISTEN for
// htmx:afterRequest in a real browser and clear aria-busy + the
// data-tc-polled-busy marker on the requesting element. SwapNone makes the
// script the ONLY possible clearer — a swap can never replace the region's
// markup — so a green run proves the listener, not a side effect.
func polledRegionE2EPage() templ.Component {
	props := layout.DefaultPageProps()
	props.Title = "PolledRegion busy cue E2E — templ-components"
	props.CSSPath = "/app.css"

	regionA := htmx.DefaultPolledRegionProps()
	regionA.ID = "region-a"
	regionA.URL = "/api/region"
	regionA.Every = "10s"
	regionA.Eager = true
	regionA.Swap = htmx.SwapNone
	regionA.ShowTimestamp = false

	regionB := regionA
	regionB.ID = "region-b"

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, region := range []htmx.PolledRegionProps{regionA, regionB} {
			if err := htmx.PolledRegion(region).
				Render(templ.WithChildren(ctx, templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
					_, err := io.WriteString(w, `<span>initial content</span>`)

					return err
				})), w); err != nil {
				return err
			}
		}

		return nil
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

func polledRegionE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	return e2ePageServer(t, polledRegionE2EPage, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/region", func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(800 * time.Millisecond)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(`<span>fresh content</span>`))
		})
	})
}

func TestPolledRegionBusyCueClearsBrowser(t *testing.T) {
	t.Parallel()

	srv := polledRegionE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	regionBusyJS := func(id string) string {
		return `(() => {
			const el = document.querySelector('#` + id + `');
			return el !== null && el.getAttribute('aria-busy') === 'true' && el.hasAttribute('data-tc-polled-busy');
		})()`
	}

	var (
		initialBusyA, initialBusyB bool
		clearedA, clearedB         bool
		contentIntact              string
		syntheticCleared           bool
	)

	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollTrue(
			`document.readyState==='complete' && window.htmx!==undefined && document.querySelector('#region-a')!==null`,
		),
		chromedp.Evaluate(regionBusyJS("region-a"), &initialBusyA),
		chromedp.Evaluate(regionBusyJS("region-b"), &initialBusyB),
		// The eager one-shot script fires the real request via htmx.ajax;
		// htmx dispatches the real htmx:afterRequest on completion — the
		// script must clear both regions (document-level delegation, one
		// listener for many regions).
		pollBool(`!`+regionBusyJS("region-a"), &clearedA),
		pollBool(`!`+regionBusyJS("region-b"), &clearedB),
		chromedp.Evaluate(`document.querySelector('#region-a span').textContent`, &contentIntact),
		// Synthetic re-arm: put the cue back on region-b and dispatch
		// htmx:afterRequest by hand — proves the listener directly, decoupled
		// from htmx's own request timing.
		chromedp.Evaluate(`(() => {
			const el = document.querySelector('#region-b');
			el.setAttribute('aria-busy', 'true');
			el.setAttribute('data-tc-polled-busy', '');
			el.dispatchEvent(new CustomEvent('htmx:afterRequest', {bubbles: true, detail: {elt: el}}));
			return el.getAttribute('aria-busy') === null && !el.hasAttribute('data-tc-polled-busy');
		})()`, &syntheticCleared),
	); err != nil {
		t.Fatalf("PolledRegion busy-cue E2E: %v", err)
	}

	if !initialBusyA || !initialBusyB {
		t.Errorf(
			"initial render: both eager regions must carry aria-busy + marker (a=%v b=%v)",
			initialBusyA,
			initialBusyB,
		)
	}

	if !clearedA || !clearedB {
		t.Errorf(
			"after real htmx:afterRequest: cue not cleared on both regions (a=%v b=%v) — SwapNone means only the script could clear it",
			clearedA,
			clearedB,
		)
	}

	if contentIntact != "initial content" {
		t.Errorf("region content must be untouched under SwapNone; got %q", contentIntact)
	}

	if !syntheticCleared {
		t.Error(
			"synthetic htmx:afterRequest dispatch did not clear the re-armed cue — listener not attached or wrong element targeting",
		)
	}
}

// TestPolledRegionEagerOuterHTMLNoRefetchLoop is the browser-level guard for
// the 2026-09-04 DiscordSync projection-health outage: an eager PolledRegion
// with the default hx-swap="outerHTML" must issue EXACTLY ONE eager request,
// even though the endpoint re-renders the same eager region on every
// response. Under the old hx-trigger="load, every Ns" implementation htmx
// re-fired "load" for every swapped-in element, producing an unbounded
// self-refetch loop. The interval is set to 1h so it can never fire during
// the test — every hit after the first is an illegal refetch.
func TestPolledRegionEagerOuterHTMLNoRefetchLoop(t *testing.T) {
	t.Parallel()

	var hits atomic.Int64

	regionProps := htmx.DefaultPolledRegionProps()
	regionProps.ID = "loop-region"
	regionProps.URL = "/api/loop"
	regionProps.Every = "1h"
	regionProps.Eager = true
	regionProps.ShowTimestamp = false

	renderRegion := func(w io.Writer) {
		_ = htmx.PolledRegion(regionProps).Render(
			templ.WithChildren(context.Background(), templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<span>loop content</span>`)

				return err
			})), w)
	}

	page := func() templ.Component {
		pageProps := layout.DefaultPageProps()
		pageProps.Title = "PolledRegion loop-guard E2E — templ-components"
		pageProps.CSSPath = "/app.css"

		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			return layout.Base(pageProps).Render(templ.WithChildren(ctx, templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				renderRegion(w)

				return nil
			})), w)
		})
	}

	srv := e2ePageServer(t, page, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/loop", func(w http.ResponseWriter, _ *http.Request) {
			hits.Add(1)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			renderRegion(w)
		})
	})

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollTrue(
			`document.readyState==='complete' && window.htmx!==undefined && document.querySelector('#loop-region')!==null`,
		),
		// Give any illegal refetch loop ample time to rack up hits.
		chromedp.Sleep(3*time.Second),
	); err != nil {
		t.Fatalf("PolledRegion loop-guard E2E: %v", err)
	}

	if got := hits.Load(); got != 1 {
		t.Errorf(
			"eager outerHTML region must fetch exactly once (the one-shot htmx.ajax); got %d hits — the pre-fix 'load' trigger looped unboundedly here",
			got,
		)
	}
}
