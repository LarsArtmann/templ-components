package visualtest

import (
	"context"
	"encoding/json"
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
)

// Signals-channel E2E (2026-11-07 Datastar plan T10, gate G3): a REAL
// datastar-patch-signals SSE stream must update data-text bindings without
// any element patch — the browser-proof twin of
// docs/recipes/datastar-signals-live-data.md. The two-tick scenario also
// proves MERGE-PATCH semantics (RFC 7396): tick one patches only
// metrics.users, tick two patches only metrics.rps, and each sibling value
// must survive the other's patch untouched.
//
// Run via `nix run .#visual`; tests skip gracefully without a browser.

func writeSSESignals(w io.Writer, patch any) error {
	buf, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("event: datastar-patch-signals\n")
	b.WriteString("data: signals ")
	b.Write(buf)
	b.WriteString("\n")
	b.WriteString("\n")

	_, err = io.WriteString(w, b.String())

	return err
}

func datastarSignalsE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	props := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w,
			`<!doctype html><html lang="en"><head><meta charset="utf-8">`+
				`<title>Datastar signals e2e — templ-components</title><link rel="stylesheet" href="/app.css">`+
				`<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`+
				`</head><body>`); err != nil {
			return err
		}

		if err := datastar.SDKScript(datastar.SDKScriptProps{
			Nonce: "ds-e2e-nonce",
			Src:   "/datastar.js",
		}).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w,
			`<div data-init="@get('/api/metrics')" data-signals="{'metrics': {users: 0, rps: 0}}">`+
				`<dl><dt>Users online</dt><dd id="sig-users" data-text="$metrics.users">0</dd>`+
				`<dt>Requests / s</dt><dd id="sig-rps" data-text="$metrics.rps">0</dd></dl></div></body></html>`)

		return err
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

	// Two ticks, then hold: users-only patch at 300ms, rps-only patch at
	// 900ms — partial patches are the merge-semantics proof.
	mux.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)

			return
		}

		_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		write := func(patch map[string]any) bool {
			if err := writeSSESignals(w, patch); err != nil {
				return false
			}

			flusher.Flush()

			return true
		}

		select {
		case <-r.Context().Done():
			return
		case <-time.After(300 * time.Millisecond):
		}

		if !write(map[string]any{"metrics": map[string]any{"users": 1237}}) {
			return
		}

		select {
		case <-r.Context().Done():
			return
		case <-time.After(600 * time.Millisecond):
		}

		if !write(map[string]any{"metrics": map[string]any{"rps": 44}}) {
			return
		}

		<-r.Context().Done()
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := props.Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestDatastarSignalsBrowser is the G3 gate: signal patches flow server →
// client without element patches, bindings re-render, and merge-patch keeps
// sibling values intact.
func TestDatastarSignalsBrowser(t *testing.T) {
	t.Parallel()

	srv := datastarSignalsE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL),
		pollTrue(`window.__dsReady===true`, chromedp.WithPollingTimeout(15*time.Second)),

		// Tick 1: users merged in, rps untouched by the partial patch.
		pollTrue(`document.getElementById('sig-users').textContent==='1237' && document.getElementById('sig-rps').textContent==='0'`, chromedp.WithPollingTimeout(10*time.Second)),

		// Tick 2: rps merged in, users SURVIVED the second partial patch.
		pollTrue(`document.getElementById('sig-rps').textContent==='44' && document.getElementById('sig-users').textContent==='1237'`, chromedp.WithPollingTimeout(10*time.Second)),
	); err != nil {
		t.Fatalf("signals channel browser proof: %v", err)
	}
}
