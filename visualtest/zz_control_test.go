package visualtest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/go-datastar/static"
)

// Temporary control experiment — delete after use.
//
// Minimal pure-htmx page: a button swaps new hx-wired content into a target
// region. Mirrors the wire form structure (content swapped into a region).
func TestDebugHtmxSwapControl(t *testing.T) {
	for _, variant := range []struct{ name, formAttrs string }{
		{name: "explicit-trigger", formAttrs: `hx-trigger="submit"`},
		{name: "implicit-trigger", formAttrs: ``},
	} {
		t.Run(variant.name, func(t *testing.T) { testSwapControlVariant(t, variant.formAttrs) })
	}
}

func testSwapControlVariant(t *testing.T, formAttrs string) {
	t.Helper()
	var hits int

	mux := http.NewServeMux()

	mux.HandleFunc("/embed.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		b, err := os.ReadFile("../layout/static/htmx.min.js")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(b)
	})

	mux.HandleFunc("/api/x", func(w http.ResponseWriter, _ *http.Request) {
		hits++
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<form id="f2" hx-post="/api/x" hx-target="#region" hx-swap="innerHTML" ` + formAttrs + `><input name="v" value="x"><button type="submit" id="btn2">again</button></form>`))
	})

	dsBytes := static.Bytes()
	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(dsBytes)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html><html><head><script src="/embed.js"></script><script src="/datastar.js"></script></head><body>
<form id="f1" hx-post="/api/x" hx-target="#region" hx-swap="innerHTML" ` + formAttrs + `><input name="v" value="y"><button type="submit" id="btn1">go</button></form>
<div id="region"></div>
</body></html>`))
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 15*time.Second)
	defer cancelTimeout()

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if e, ok := ev.(*network.EventRequestWillBeSent); ok && e.Request.URL == srv.URL+"/api/x" {
			t.Logf("REQ %s", e.Request.URL)
		}
	})

	var (
		ok    bool
		state string
	)

	arm := `window.__events=[]; ['htmx:afterSwap','htmx:afterSettle','htmx:load','htmx:afterProcessNode','htmx:configRequest','htmx:beforeSwap'].forEach(function(n){ document.addEventListener(n, function(e){ window.__events.push(n + (e.target && e.target.id ? ':'+e.target.id : '')); }); }); 'armed'`

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(`document.readyState==='complete' && window.htmx!==undefined`, &ok),
		chromedp.Evaluate(arm, &state),
		chromedp.Click("#btn1", chromedp.NodeVisible),
		chromedp.Poll(`document.querySelector('#btn2')!==null`, &ok),
		chromedp.Sleep(400*time.Millisecond),
		chromedp.Evaluate(`JSON.stringify({events: window.__events})`, &state),
	); err != nil {
		t.Fatalf("control pre: %v (hits=%d)", err, hits)
	}

	t.Logf("STATE: %s", state)

	if err := chromedp.Run(ctx,
		chromedp.Click("#btn2", chromedp.NodeVisible),
		chromedp.Sleep(400*time.Millisecond),
		chromedp.Evaluate(`JSON.stringify({events2: window.__events})`, &state),
	); err != nil {
		t.Fatalf("control click2: %v (hits=%d)", err, hits)
	}
	t.Logf("STATE2: %s", state)
	if hits < 2 {
		t.Fatalf("second submit did not fire: hits=%d", hits)
	}

	t.Logf("control passed: hits=%d", hits)
}
