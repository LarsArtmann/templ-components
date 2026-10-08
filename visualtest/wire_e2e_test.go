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
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils/wire"
)

// The wire E2E suite proves the transport-agnostic wiring contract at the
// browser level, which string tests cannot: a real htmx trigger engine and a
// real Datastar runtime must both execute the SAME wire.Action shape against
// the SAME endpoint (wrapped in the library's own wire.Handler middleware)
// and land the fragment in the right region.
//
// Run via `nix run .#visual` (sets CHROMEDP_CHROME_PATH); tests skip
// gracefully without a browser, exactly like the screenshot suite.

// wireE2EPage renders the full consumer page shell: layout.Base injects the
// self-hosted (embedded) htmx runtime; the Datastar runtime loads from
// /datastar.js, served from the pinned go-datastar/static bundle so the test
// never touches a CDN. The head registers a datastar-ready catcher BEFORE the
// module executes — the runtime dispatches that event on document when its
// engine has booted, giving a deterministic readiness signal instead of a
// sleep.
func wireE2EPage() templ.Component {
	props := layout.DefaultPageProps()
	props.Title = "Wire E2E — templ-components"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(
			w,
			`<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`,
		)

		return err
	})

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Same runtime injection point a consumer page uses (the demo does
		// exactly this inside its datastar section).
		if err := datastar.SDKScript(datastar.SDKScriptProps{
			Nonce: "wire-e2e-nonce",
			Src:   "/datastar.js",
		}).Render(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `<div class="p-4 flex flex-wrap gap-3 items-start">`); err != nil {
			return err
		}

		htmxButton := display.ButtonProps{
			ID: "btn-wire-htmx", Class: "mr-4",
			Text:    "Load via htmx",
			Variant: display.ButtonSecondary,
			Size:    display.ButtonSizeSM,
			Wire: &wire.Action{
				URL:    "/api/wire/fragment",
				Target: "#wire-htmx-out",
			},
		}
		if err := display.Button(htmxButton).Render(ctx, w); err != nil {
			return err
		}

		datastarButton := display.ButtonProps{
			ID:      "btn-wire-datastar",
			Text:    "Load via Datastar",
			Variant: display.ButtonSecondary,
			Size:    display.ButtonSizeSM,
			Wire: &wire.Action{
				Transport: wire.TransportDatastar,
				URL:       "/api/wire/fragment",
			},
		}
		if err := display.Button(datastarButton).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w,
			`<div id="wire-htmx-out" class="basis-full"></div>`+
				`<div id="wire-datastar-out" class="basis-full"></div>`+
				`</div>`)

		return err
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

// wireE2EServer serves the E2E page plus its assets and the shared fragment
// endpoint. The endpoint is wrapped in the production wire.Handler middleware
// — the E2E exercises the library's own server-side contract, not a copy.
func wireE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(css)
	})

	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(static.Bytes())
	})

	fragment := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := feedback.InlineSuccess(wireFragmentText).Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.Handle("/api/wire/fragment", wire.Handler(wire.PatchTarget{
		Selector: "#wire-datastar-out",
		Mode:     wire.PatchModeInner,
	}, fragment))

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := wireE2EPage().Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// wireFragmentText is the assertion needle inside the served fragment.
const wireFragmentText = "wire contract"

func TestWireE2EHTMXButtonPatchesTarget(t *testing.T) {
	t.Parallel()

	srv := wireE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var (
		htmxDefined, fragment bool
		out                   string
	)

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		// htmx loads inline (self-host) and processes hx-* nodes on
		// DOMContentLoaded; "complete" readyState guarantees both.
		pollBool(`document.readyState==='complete' && window.htmx!==undefined`, &htmxDefined),
		chromedp.Click("#btn-wire-htmx", chromedp.NodeVisible),
		pollBool(`document.querySelector('#wire-htmx-out').innerHTML.length>0`, &fragment),
		evalInto(chromedp.InnerHTML(chromedp.CSS("#wire-htmx-out"), chromedp.NodeVisible), &out),
	); err != nil {
		t.Fatalf("htmx E2E: %v", err)
	}

	assertFragmentLanded(t, "htmx", out)
}

func TestWireE2EDatastarButtonPatchesSelector(t *testing.T) {
	t.Parallel()

	srv := wireE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var (
		dsReady, fragment bool
		out               string
	)

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		// The pinned runtime dispatches datastar-ready on document when its
		// engine booted; before that, data-on:* clicks are inert.
		pollBool(`window.__dsReady===true`, &dsReady),
		chromedp.Click("#btn-wire-datastar", chromedp.NodeVisible),
		pollBool(`document.querySelector('#wire-datastar-out').innerHTML.length>0`, &fragment),
		evalInto(chromedp.InnerHTML(chromedp.CSS("#wire-datastar-out"), chromedp.NodeVisible), &out),
	); err != nil {
		t.Fatalf("datastar E2E: %v", err)
	}

	assertFragmentLanded(t, "datastar", out)
}

func assertFragmentLanded(t *testing.T, transport, out string) {
	t.Helper()

	if !strings.Contains(out, wireFragmentText) {
		t.Fatalf("%s transport: fragment did not land in the target region; got %q", transport, out)
	}
}

// wireSwapE2EPage renders the Swap/mode matrix (L1-06), pinned against the
// corrected runtime model (bundle-decoded + browser-proven 2026-10-02): htmx
// outerHTML swap is client-side; Datastar patch target and mode come
// EXCLUSIVELY from the response headers (wire.Handler/PatchTarget), a client
// {mode}/{selector} option is inert for targeting, and a header-less
// response whose fragment root carries the region's id still patches by
// id-match (outer default).
func wireSwapE2EPage() templ.Component {
	props := layout.DefaultPageProps()
	props.Title = "Wire Swap E2E — templ-components"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(
			w,
			`<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`,
		)

		return err
	})

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := datastar.SDKScript(datastar.SDKScriptProps{
			Nonce: "wire-e2e-nonce",
			Src:   "/datastar.js",
		}).Render(ctx, w); err != nil {
			return err
		}

		buttons := []display.ButtonProps{
			{
				ID:      "btn-swap-htmx",
				Text:    "htmx outer swap",
				Variant: display.ButtonSecondary,
				Size:    display.ButtonSizeSM,
				Wire: &wire.Action{
					URL:    "/api/wire/swap-outer",
					Target: "#htmx-outer-region",
					Swap:   wire.PatchModeOuter,
				},
			},
			{
				ID:      "btn-swap-ds-outer",
				Text:    "Datastar outer (headers)",
				Variant: display.ButtonSecondary,
				Size:    display.ButtonSizeSM,
				Wire: &wire.Action{
					Transport: wire.TransportDatastar,
					URL:       "/api/wire/swap-outer",
				},
			},
			{
				ID:      "btn-swap-ds-append",
				Text:    "Datastar: server append header wins",
				Variant: display.ButtonSecondary,
				Size:    display.ButtonSizeSM,
				Wire: &wire.Action{
					Transport: wire.TransportDatastar,
					URL:       "/api/wire/swap-override",
					Swap:      wire.PatchModeInner,
				},
			},
			{
				ID:      "btn-swap-ds-idmatch",
				Text:    "Datastar: id-matched outer (no headers)",
				Variant: display.ButtonSecondary,
				Size:    display.ButtonSizeSM,
				Wire: &wire.Action{
					Transport: wire.TransportDatastar,
					URL:       "/api/wire/idmatch",
				},
			},
			{
				ID:      "btn-swap-ds-selector",
				Text:    "Datastar: client selector does not target",
				Variant: display.ButtonSecondary,
				Size:    display.ButtonSizeSM,
				Wire: &wire.Action{
					Transport: wire.TransportDatastar,
					URL:       "/api/wire/unheaded",
					Selector:  "#ds-selector-region",
				},
			},
		}

		if _, err := io.WriteString(w, `<div class="p-4 flex flex-wrap gap-3 items-start">`); err != nil {
			return err
		}

		for i := range buttons {
			if err := display.Button(buttons[i]).Render(ctx, w); err != nil {
				return err
			}
		}

		_, err := io.WriteString(w,
			`<div id="htmx-outer-wrap" class="basis-full"><div id="htmx-outer-region">htmx-outer-sentinel</div></div>`+
				`<div id="ds-outer-wrap" class="basis-full"><div id="ds-outer-region">ds-outer-sentinel</div></div>`+
				`<div id="ds-override-wrap" class="basis-full"><div id="ds-override-region">ds-override-sentinel</div></div>`+
				`<div id="ds-idmatch-region">idmatch-sentinel</div>`+
				`<div id="ds-selector-region">ds-selector-sentinel</div>`+
				`</div>`)

		return err
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

func wireSwapE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(css)
	})

	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(static.Bytes())
	})

	fragment := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := feedback.InlineSuccess(wireFragmentText).Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.Handle("/api/wire/swap-outer", wire.Handler(wire.PatchTarget{
		Selector: "#ds-outer-region",
		Mode:     wire.PatchModeOuter,
	}, fragment))

	mux.Handle("/api/wire/swap-override", wire.Handler(wire.PatchTarget{
		Selector: "#ds-override-region",
		Mode:     wire.PatchModeAppend,
	}, fragment))

	// Header-less endpoint whose fragment root carries the region's id: the
	// runtime's outer-default id-match must patch the region anyway.
	mux.HandleFunc("/api/wire/idmatch", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<div id="ds-idmatch-region">id-matched contract</div>`))
	})

	// Header-less endpoint whose fragment root matches NOTHING: with no
	// client-side targeting, nothing may change.
	mux.HandleFunc("/api/wire/unheaded", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<p id="never-in-dom">unheaded probe fragment</p>`))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := wireSwapE2EPage().Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

func TestWireE2ESwapModes(t *testing.T) {
	srv := wireSwapE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 90*time.Second)
	defer cancelTimeout()

	var (
		ready                                            bool
		htmxWrap, dsWrap, override, idmatch, selectorBox string
	)

	if err := chromedp.Do(
		ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollBool(`document.readyState==='complete' && window.htmx!==undefined && window.__dsReady===true`, &ready),
		// htmx: hx-swap="outerHTML" replaces the region element itself.
		chromedp.Click("#btn-swap-htmx", chromedp.NodeVisible),
		pollBool(`!document.querySelector('#htmx-outer-region')`, &ready),
		evalInto(chromedp.InnerHTML(chromedp.CSS("#htmx-outer-wrap"), chromedp.NodeVisible), &htmxWrap),
		// Datastar: response-header targeting (selector + mode from the
		// wire.Handler) patches the region with outer mode.
		chromedp.Click("#btn-swap-ds-outer", chromedp.NodeVisible),
		pollBool(`!document.querySelector('#ds-outer-region')`, &ready),
		evalInto(chromedp.InnerHTML(chromedp.CSS("#ds-outer-wrap"), chromedp.NodeVisible), &dsWrap),
		// The server's Datastar-Mode header WINS: a client {mode:'inner'} is
		// inert, so the append header appends — the sentinel survives.
		chromedp.Click("#btn-swap-ds-append", chromedp.NodeVisible),
		pollBool(
			`document.querySelector('#ds-override-region') && document.querySelector('#ds-override-region').innerHTML.indexOf('`+wireFragmentText+`')>=0`,
			&ready,
		),
		evalInto(chromedp.InnerHTML(chromedp.CSS("#ds-override-region"), chromedp.NodeVisible), &override),
		// Header-less response, fragment root id matches the region: the
		// outer-default id-match replaces it.
		chromedp.Click("#btn-swap-ds-idmatch", chromedp.NodeVisible),
		pollBool(
			`document.querySelector('#ds-idmatch-region') && document.querySelector('#ds-idmatch-region').innerHTML==='id-matched contract'`,
			&ready,
		),
		evalInto(chromedp.InnerHTML(chromedp.CSS("#ds-idmatch-region"), chromedp.NodeVisible), &idmatch),
		// Client {selector} must NOT target: header-less + unmatched fragment
		// root leaves the region untouched.
		chromedp.Click("#btn-swap-ds-selector", chromedp.NodeVisible),
		chromedp.Sleep(1200*time.Millisecond),
		evalInto(chromedp.InnerHTML(chromedp.CSS("#ds-selector-region"), chromedp.NodeVisible), &selectorBox),
	); err != nil {
		t.Fatalf("swap/mode E2E: %v", err)
	}

	assertFragmentLanded(t, "htmx outer", htmxWrap)
	assertFragmentLanded(t, "datastar outer", dsWrap)
	assertFragmentLanded(t, "datastar header-mode append", override)

	if !strings.Contains(override, "ds-override-sentinel") {
		t.Fatalf(
			"server Datastar-Mode: append header did not win over the client {mode:'inner'}: sentinel vanished: %q",
			override,
		)
	}

	if idmatch != "id-matched contract" {
		t.Fatalf("id-matched outer patch failed: got %q", idmatch)
	}

	if selectorBox != "ds-selector-sentinel" {
		t.Fatalf("client {selector} must not target patches: region changed: %q", selectorBox)
	}
}
