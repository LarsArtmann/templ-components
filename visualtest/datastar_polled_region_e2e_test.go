package visualtest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/go-datastar/static"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/layout"
)

// The Datastar PolledRegion browser E2E proves the polling loop end to end
// against the PINNED runtime bundle: the data-on-interval attribute fires on
// schedule, the endpoint's outer self-patch replaces the region, the runtime
// re-initializes the patched-in element (so polling continues), and — the
// anti-regression core — the loop stays INTERVAL-BOUND (no eager re-arm
// storm; this component deliberately never emits the leading flag, which
// under outer self-patch would refetch immediately forever).
//
// Run via `nix run .#visual`; tests skip gracefully without a browser.

// datastarPolledRegionPartial renders the PolledRegion with a tick counter
// as its server-side initial content. Both the page body and the
// /partials/stats endpoint render THROUGH this helper, so the endpoint
// response is the same component the page rendered (the self-patch contract).
func datastarPolledRegionPartial(tick int) templ.Component {
	body := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<p id="tick-value">tick %d</p>`, tick)

		return err
	})

	props := datastar.PolledRegionProps{
		ID:            "polled-region",
		URL:           "/partials/stats",
		Every:         "500ms",
		ShowTimestamp: true,
	}

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return datastar.PolledRegion(props).Render(templ.WithChildren(ctx, body), w)
	})
}

func datastarPolledRegionE2EServer(t *testing.T) (*httptest.Server, *atomic.Int64) {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	var ticks atomic.Int64

	props := layout.DefaultPageProps()
	props.Title = "Datastar PolledRegion E2E — templ-components"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(
			w,
			`<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`,
		)

		return err
	})

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, component := range []templ.Component{
			datastar.SDKScript(datastar.SDKScriptProps{
				Nonce: "polled-e2e-nonce",
				Src:   "/datastar.js",
			}),
			datastarPolledRegionPartial(0),
		} {
			if err := component.Render(ctx, w); err != nil {
				return err
			}
		}

		return nil
	})

	page := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(css)
	})

	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(static.Bytes())
	})

	mux.HandleFunc("/partials/stats", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		var sb strings.Builder
		if err := datastarPolledRegionPartial(int(ticks.Add(1))).Render(context.Background(), &sb); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		_, _ = io.WriteString(w, sb.String())
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := page.Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv, &ticks
}

// TestDatastarPolledRegionBrowser proves: the interval fires (server-side
// counter advances past 2), the self-patch lands (tick text updated inside
// the region), the region never duplicates (exactly one element), and the
// request rate stays interval-bound (no eager re-arm loop).
func TestDatastarPolledRegionBrowser(t *testing.T) {
	t.Parallel()

	srv, ticks := datastarPolledRegionE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollTrue(`window.__dsReady===true && document.querySelector('#polled-region')!==null`),
		// At least two interval ticks must have fetched AND patched: the
		// endpoint increments per request, so "tick 2" proves two real
		// round-trips through the runtime's interval plugin.
		pollTrue(`(() => {
			const el = document.getElementById('tick-value');
			return el !== null && /tick ([2-9]|\d{2,})/.test(el.textContent);
		})()`),
	); err != nil {
		t.Fatalf("PolledRegion interval/self-patch proof: %v", err)
	}

	// Exactly ONE region element: the outer self-patch replaces, never
	// duplicates (a broken re-arm strategy would stack regions).
	var regionCount int
	if err := chromedp.Do(ctx,
		evalExprInto(`document.querySelectorAll('#polled-region').length`, &regionCount),
	); err != nil {
		t.Fatalf("region count: %v", err)
	}

	if regionCount != 1 {
		t.Errorf("region must not duplicate under self-patch; got %d #polled-region elements", regionCount)
	}

	// Interval-bound request rate: settle ~2.5s more and read the server
	// counter. At a 500ms interval the expected band is single digits; an
	// eager re-arm loop (leading + outer self-patch) would be hundreds+
	// within the first second. The generous ceiling keeps the assertion
	// load-insensitive while still failing the loop class by two orders.
	before := ticks.Load()

	time.Sleep(2500 * time.Millisecond)

	after := ticks.Load()

	if delta := after - before; delta > 20 {
		t.Errorf("request rate not interval-bound: %d requests in 2.5s at a 500ms interval — eager re-arm loop", delta)
	}

	if after < 3 {
		t.Errorf("polling stopped early: %d total ticks after ~4s at a 500ms interval", after)
	}
}
