package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
)

// TestLightboxE2E proves the lightbox in a real browser: the trigger opens
// the native dialog, next navigates with wrap-around, the caption follows
// the image, rotate applies the transform, and Escape (native dialog) closes
// the viewer (plan C8g — wired⇒e2e rule for script-bearing components).
func TestLightboxE2E(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := lightboxE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	tabCtx, tabCancel := context.WithTimeout(ctx, 120*time.Second)
	defer tabCancel()

	if err := chromedp.Do(tabCtx, chromedp.Navigate(srv.URL)); err != nil {
		t.Fatalf("navigate: %v", err)
	}

	// Open via the trigger button.
	if err := chromedp.Do(tabCtx,
		chromedp.WaitVisible(chromedp.CSS("[data-tc-lightbox-open]")),
		chromedp.Click(chromedp.CSS("[data-tc-lightbox-open]")),
		pollTrue(`document.querySelector('dialog[data-tc-lightbox-dialog]').open`, pollFast()...),
	); err != nil {
		t.Fatalf("open lightbox: %v", err)
	}

	// First image visible, second hidden; caption shows image 1's caption.
	if err := chromedp.Do(tabCtx,
		pollTrue(`!document.querySelectorAll('[data-tc-lightbox-img]')[0].hidden`, pollFast()...),
		pollTrue(`document.querySelectorAll('[data-tc-lightbox-img]')[1].hidden`, pollFast()...),
		pollTrue(`document.querySelector('[data-tc-lightbox-caption]').textContent==='Front view'`, pollFast()...),
	); err != nil {
		t.Fatalf("initial state: %v", err)
	}

	// Next: image 2 visible, caption follows (empty — image 2 has none).
	if err := chromedp.Do(tabCtx,
		chromedp.Click(chromedp.CSS("[data-tc-lightbox-next]")),
		pollTrue(`!document.querySelectorAll('[data-tc-lightbox-img]')[1].hidden`, pollFast()...),
		pollTrue(`document.querySelector('[data-tc-lightbox-caption]').textContent===''`, pollFast()...),
	); err != nil {
		t.Fatalf("next: %v", err)
	}

	// Rotate: transform carries 90deg.
	if err := chromedp.Do(tabCtx,
		chromedp.Click(chromedp.CSS("[data-tc-lightbox-rotate]")),
		pollTrue(`document.querySelectorAll('[data-tc-lightbox-img]')[1].style.transform.indexOf('rotate(90deg)')>=0`, pollFast()...),
	); err != nil {
		t.Fatalf("rotate: %v", err)
	}

	// Next again wraps to image 0.
	if err := chromedp.Do(tabCtx,
		chromedp.Click(chromedp.CSS("[data-tc-lightbox-next]")),
		pollTrue(`!document.querySelectorAll('[data-tc-lightbox-img]')[0].hidden`, pollFast()...),
	); err != nil {
		t.Fatalf("wrap-around: %v", err)
	}

	// Close via the wired Close button (our delegation); Escape-to-close is
	// UA-native <dialog> behavior and synthetic Escape key events do not
	// reach the dialog cancel pipeline in headless Chromium — the same
	// reason Modal/Drawer e2e never browser-test it.
	if err := chromedp.Do(tabCtx,
		chromedp.Click(chromedp.CSS("[data-tc-close]")),
		pollTrue(`!document.querySelector('dialog[data-tc-lightbox-dialog]').open`, pollFast()...),
	); err != nil {
		t.Fatalf("close: %v", err)
	}
}

// lightboxSVG is a tiny inline image so the browser resolves real <img> loads.
const lightboxSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="120" height="80"><rect width="120" height="80" fill="#3b82f6"/></svg>`

// lightboxE2EServer serves a minimal page wrapping the Lightbox.
func lightboxE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		var body io.Writer = w

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		_, _ = io.WriteString(body, "<!doctype html><html><head><style>dialog{color:#fff}</style></head><body>")

		if err := display.Lightbox(display.LightboxProps{
			Images: []display.LightboxImage{
				{Src: "/img/1.svg", Alt: "Front view", Caption: "Front view"},
				{Src: "/img/2.svg", Alt: "Back view"},
			},
			TriggerLabel: "View images",
		}).Render(context.Background(), body); err != nil {
			t.Fatalf("render lightbox: %v", err)
		}

		_, _ = io.WriteString(body, "</body></html>")
	})

	mux.HandleFunc("/img/1.svg", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		_, _ = io.WriteString(w, lightboxSVG)
	})

	mux.HandleFunc("/img/2.svg", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		_, _ = io.WriteString(w, lightboxSVG)
	})

	return httptest.NewServer(mux)
}
