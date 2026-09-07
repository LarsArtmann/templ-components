package visualtest

import (
	"context"
	"net/http"
	"os"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Temporary control experiment — delete after use.
//
// Minimal pure-htmx page: a button whose swap target is its own ANCESTOR
// region; the endpoint returns a NEW button with the same hx- attributes.
// Mirrors the wire form structure (form swapped into ancestor region).
func TestDebugHtmxSwapControl(t *testing.T) {
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
		_, _ = w.Write([]byte(`<button id="btn2" hx-post="/api/x" hx-target="#region" hx-swap="innerHTML">again `))
		_, _ = w.Write([]byte("hit " + string(rune('0'+hits)) + `</button>`))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html><html><head><script src="/embed.js"></script></head><body>
<button id="btn1" hx-post="/api/x" hx-target="#region" hx-swap="innerHTML">go</button>
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
		ok   bool
		evts any
	)

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(`document.readyState==='complete' && window.htmx!==undefined`, &ok),
		chromedp.Evaluate(`window.__events=[]; ['htmx:afterSwap','htmx:afterSettle','htmx:load','htmx:afterProcessNode','htmx:error','htmx:sendError'].forEach(function(n){ document.addEventListener(n, function(e){ window.__events.push(n + (e.target && e.target.id ? ':'+e.target.id : '')); }); }); 'armed'`, &evts),
		chromedp.Click("#btn1", chromedp.NodeVisible),
		chromedp.Poll(`document.querySelector('#btn2')!==null`, &ok),
		chromedp.Evaluate(`new Promise(function(res){ setTimeout(function(){ res(JSON.stringify(window.__events)); }, 300); })`, &evts),
		chromedp.Poll(`document.querySelector('#region').innerText.includes('hit 2')`, &ok),
	); err != nil {
		t.Fatalf("control: %v (hits=%d) events=%v", err, hits, evts)
	}
	t.Logf("EVENTS: %v", evts)
	var evts2 string
	if err := chromedp.Run(ctx, chromedp.Evaluate(`JSON.stringify(window.__events)`, &evts2)); err != nil {
		t.Fatalf("evts2: %v", err)
	}
	t.Logf("EVENTS2: %s", evts2)

	t.Logf("control passed: hits=%d", hits)
}
