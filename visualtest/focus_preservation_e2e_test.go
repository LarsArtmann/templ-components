package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/htmx"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// TestFocusPreservationE2E (M20/F093): after HTMX swaps, keyboard focus must
// not silently drop to <body> — the failure mode that strands keyboard and
// screen-reader users mid-interaction.
//
// Contract 1 — LoadMore self-replaces its own trigger: the response button
// (FocusOnSwap) carries autofocus, and the embedded htmx 2.0.10 runtime
// focuses swapped-in [autofocus] elements (verified against
// layout/static/htmx.min.js), so focus lands on the REPLACEMENT button.
//
// Contract 2 — SwapOOB patches content ELSEWHERE in the page while the
// trigger keeps its DOM position: focus on the trigger must survive.
//
// Serial by design (shared tab, ordered interactions); no t.Parallel — the
// pack-e2e allocator lesson.
func TestFocusPreservationE2E(t *testing.T) {
	srv := focusPreservationServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	var ok bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(packGate(wire.TransportHTMX), &ok),
	); err != nil {
		t.Fatalf("focus e2e setup: %v", err)
	}

	t.Run("loadmore replacement button receives focus", func(t *testing.T) {
		var done string

		if err := chromedp.Run(ctx,
			chromedp.Evaluate(`(document.querySelector('#items-load-more button').click(),'')`, &done),
			chromedp.Poll(
				`document.activeElement !== null && document.activeElement.closest('#items-load-more') !== null && document.activeElement.tagName === 'BUTTON' && document.activeElement.textContent.trim() === 'Load more' ? 'ok' : ''`,
				&done,
			),
		); err != nil {
			dumpFocusState(t, ctx, "loadmore")

			t.Fatalf("loadmore focus restoration: %v", err)
		}

		// The swap also delivered the next batch (focus proof is not a
		// no-op render). The response replaces the button AT ITS POSITION —
		// outside #items — so count page-wide.
		if err := chromedp.Run(ctx,
			chromedp.Poll(`document.querySelectorAll('.fp-card').length >= 4 ? 'ok' : ''`, &done),
		); err != nil {
			t.Fatalf("loadmore batch delivery: %v", err)
		}
	})

	t.Run("swapoob keeps trigger focus", func(t *testing.T) {
		var done, active string

		if err := chromedp.Run(ctx,
			chromedp.Evaluate(`(document.getElementById('oob-trigger').focus(),'')`, &done),
			chromedp.Evaluate(`(document.getElementById('oob-trigger').click(),'')`, &done),
			chromedp.Poll(
				`document.getElementById('counter').textContent.includes('1') && document.getElementById('status').textContent.includes('updated') ? 'ok' : ''`,
				&done,
			),
		); err != nil {
			dumpFocusState(t, ctx, "swapoob")

			t.Fatalf("oob swap application: %v", err)
		}

		if err := chromedp.Run(ctx,
			chromedp.Evaluate(`document.activeElement && document.activeElement.id`, &active),
		); err != nil {
			t.Fatalf("activeElement read: %v", err)
		}

		if active != "oob-trigger" {
			t.Errorf("focus lost to %q after OOB swap — want the trigger #oob-trigger", active)
		}
	})
}

// dumpFocusState prints where focus is and the live DOM regions on failure.
func dumpFocusState(t *testing.T, ctx context.Context, flow string) {
	t.Helper()

	var body string

	if err := chromedp.Run(
		ctx,
		chromedp.Evaluate(
			`JSON.stringify({active: document.activeElement && (document.activeElement.id || document.activeElement.tagName), items: document.getElementById('items') && document.getElementById('items').outerHTML.slice(0, 400), status: document.getElementById('status') && document.getElementById('status').outerHTML, counter: document.getElementById('counter') && document.getElementById('counter').outerHTML})`,
			&body,
		),
	); err != nil {
		t.Logf("%s: dom dump failed: %v", flow, err)

		return
	}

	t.Logf("%s: state: %s", flow, body)
}

// focusPreservationServer serves one page with a LoadMore list and a
// SwapOOB trigger, plus the two htmx endpoints they request.
func focusPreservationServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		cursor := r.URL.Query().Get("cursor")
		if cursor != "1" {
			http.Error(w, "unexpected cursor", http.StatusBadRequest)

			return
		}

		props := navigation.LoadMoreProps{
			BaseProps:   utils.BaseProps{ID: "items-load-more"},
			Endpoint:    "/api/items",
			Cursor:      "2",
			FocusOnSwap: true,
		}

		if err := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, focusCard(3)+focusCard(4)); err != nil {
				return err
			}

			return navigation.LoadMore(props).Render(r.Context(), w)
		}).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("GET /api/oob", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, `<p id="status" class="text-sm">updated</p>`); err != nil {
				return err
			}

			return htmx.SwapOOB(htmx.SwapOOBProps{
				BaseProps: utils.BaseProps{ID: "counter"},
				Selector:  "#counter",
				SwapStyle: htmx.SwapOuterHTML,
			}).Render(templ.WithChildren(r.Context(),
				templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
					_, err := io.WriteString(w, `1`)

					return err
				}),
			), w)
		}).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		pageProps := layout.DefaultPageProps()
		pageProps.Title = "Focus preservation e2e"

		body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if _, err := io.WriteString(
				w,
				`<div id="items" class="space-y-2">`+focusCard(1)+focusCard(2)+`</div>`+
					`<p id="status" class="text-sm">idle</p>`+
					`<p>count: <span id="counter">0</span></p>`+
					`<button id="oob-trigger" type="button" hx-get="/api/oob" hx-target="#status" hx-swap="outerHTML">OOB swap</button>`,
			); err != nil {
				return err
			}

			return navigation.LoadMore(navigation.LoadMoreProps{
				BaseProps: utils.BaseProps{ID: "items-load-more"},
				Endpoint:  "/api/items",
				Cursor:    "1",
			}).Render(ctx, w)
		})

		if err := layout.Base(pageProps).Render(templ.WithChildren(r.Context(), body), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// focusCard is one list item card in the focus e2e batch.
func focusCard(n int) string {
	return `<div class="fp-card rounded-lg border border-gray-200 p-3"><p class="text-sm font-medium">Item ` +
		string(rune('0'+n)) + `</p></div>`
}
