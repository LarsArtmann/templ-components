package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/go-datastar/static"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
)

// The Datastar runtime E2E suite closes the #147 gap: the SSEErrorHandling
// listener and the LiveRegion busy-cue listener were only string-pinned —
// these tests drive the PINNED runtime bundle (served locally, no CDN) with
// REAL fetches (an HTTP 500 and a real datastar-patch-elements SSE stream)
// and assert the DOM outcomes.
//
// Run via `nix run .#visual`; tests skip gracefully without a browser.

// writeSSEPatch mirrors the demo's writeDatastarPatch wire format (pinned by
// examples/demo/sse_test.go): event name, keyed datalines, blank terminator.
func writeSSEPatch(w io.Writer, selector, mode, html string) {
	var b strings.Builder

	b.WriteString("event: datastar-patch-elements\n")
	if selector != "" {
		b.WriteString("data: selector " + selector + "\n")
	}
	if mode != "" {
		b.WriteString("data: mode " + mode + "\n")
	}
	for line := range strings.SplitSeq(strings.TrimSpace(html), "\n") {
		b.WriteString("data: elements " + strings.TrimSuffix(line, "\r") + "\n")
	}
	b.WriteString("\n")

	_, _ = io.WriteString(w, b.String()) //nolint:errcheck // best-effort; client-gone ends the stream
}

// datastarRuntimeE2EServer serves a page with the pinned Datastar runtime,
// ToastContainer, SSEErrorHandling, a failing-fetch trigger, and an
// auto-starting LiveRegion wired to a healthy SSE endpoint. Both #147
// assertions share this server.
func datastarRuntimeE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	props := layout.DefaultPageProps()
	props.Title = "Datastar runtime E2E — templ-components"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(
			w,
			`<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`,
		)

		return err
	})

	liveProps := datastar.DefaultLiveRegionProps()
	liveProps.BaseProps.ID = "live-region"
	liveProps.URL = "/api/good-sse"
	liveProps.AutoStart = true

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, component := range []templ.Component{
			datastar.SDKScript(datastar.SDKScriptProps{
				BaseProps: utils.BaseProps{Nonce: "ds-e2e-nonce"},
				Src:       "/datastar.js",
			}),
			feedback.ToastContainer(""),
			datastar.SSEErrorHandling(datastar.SSEErrorHandlingConfig{Nonce: "ds-e2e-nonce"}),
			datastar.LiveRegion(liveProps),
		} {
			if err := component.Render(ctx, w); err != nil {
				return err
			}
		}

		_, err := io.WriteString(
			w,
			`<div id="live-child">initial</div>`+
				`<button id="bad-fetch-trigger" type="button" data-on:click="@get('/api/bad-sse')">trigger</button>`,
		)

		return err
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

	mux.HandleFunc("/api/bad-sse", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})

	mux.HandleFunc("/api/good-sse", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)

			return
		}

		_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		// First patch delayed so the test can observe the initial aria-busy
		// state before any event arrives.
		select {
		case <-r.Context().Done():
			return
		case <-time.After(1500 * time.Millisecond):
		}

		writeSSEPatch(w, "#live-child", "inner", "<span>fresh via sse</span>")
		flusher.Flush()

		// Hold the stream open until the client goes away.
		<-r.Context().Done()
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := page.Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestDatastarSSEErrorHandlingBrowser proves the SSEErrorHandling JS against
// the pinned runtime: a REAL @get that answers HTTP 500 must make the runtime
// dispatch datastar-fetch {type: 'error'}, which the component turns into an
// announcer message and a toast.
func TestDatastarSSEErrorHandlingBrowser(t *testing.T) {
	t.Parallel()

	srv := datastarRuntimeE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var (
		announcerText string
		toastText     string
	)

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(`window.__dsReady===true && document.querySelector('#tc-datastar-announcer')!==null`, nil),
		// A real click drives the pinned runtime's fetch plugin; the 500
		// response dispatches datastar-fetch {type: 'error'} and the component
		// announces + toasts it.
		chromedp.Click("#bad-fetch-trigger", chromedp.NodeVisible),
		chromedp.Poll(`document.getElementById('tc-datastar-announcer').textContent.includes('Stream error')`, nil),
		chromedp.Evaluate(`document.getElementById('tc-datastar-announcer').textContent`, &announcerText),
		chromedp.Poll(`(() => {
			const toasts = document.querySelectorAll('#tc-toast-container > div');
			for (const t of toasts) {
				if (t.textContent.includes('live stream endpoint returned an error')) return true;
			}
			return false;
		})()`, nil),
		chromedp.Evaluate(`(() => {
			const toasts = document.querySelectorAll('#tc-toast-container > div');
			for (const t of toasts) {
				if (t.textContent.includes('live stream endpoint returned an error')) return t.textContent;
			}
			return '';
		})()`, &toastText),
	); err != nil {
		t.Fatalf("SSEErrorHandling browser proof: %v", err)
	}

	if !strings.Contains(announcerText, "HTTP 500") {
		t.Errorf("announcer must carry the HTTP status; got %q", announcerText)
	}

	if !strings.Contains(toastText, "Stream Error") {
		t.Errorf("toast must carry the error title; got %q", toastText)
	}
}

// TestDatastarLiveRegionBusyClearBrowser proves the LiveRegion busy-cue
// listener against the pinned runtime: the region renders aria-busy + the
// data-tc-live-busy marker, and a REAL datastar-patch-elements SSE event must
// clear the cue (the listener treats the applied patch event name as the
// "first content arrived" signal) while patching the child content.
func TestDatastarLiveRegionBusyClearBrowser(t *testing.T) {
	t.Parallel()

	srv := datastarRuntimeE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var (
		initialBusy bool
		childText   string
	)

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(`window.__dsReady===true && document.querySelector('#live-region')!==null`, nil),
		// Before the delayed first patch: the cue must be present.
		chromedp.Evaluate(`(() => {
			const el = document.querySelector('#live-region');
			return el.getAttribute('aria-busy') === 'true' && el.hasAttribute('data-tc-live-busy');
		})()`, &initialBusy),
		// After the patch: cue cleared on the region AND content landed.
		chromedp.Poll(`(() => {
			const el = document.querySelector('#live-region');
			return el.getAttribute('aria-busy') === null && !el.hasAttribute('data-tc-live-busy')
				&& document.querySelector('#live-child').textContent.includes('fresh via sse');
		})()`, nil),
		chromedp.Evaluate(`document.querySelector('#live-child').textContent`, &childText),
	); err != nil {
		t.Fatalf("LiveRegion busy-clear browser proof: %v", err)
	}

	if !initialBusy {
		t.Error("auto-starting LiveRegion must render aria-busy + data-tc-live-busy before the first patch")
	}

	if !strings.Contains(childText, "fresh via sse") {
		t.Errorf("patched child content missing; got %q", childText)
	}
}
