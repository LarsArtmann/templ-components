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
	"github.com/larsartmann/templ-components/utils"
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
			BaseProps: utils.BaseProps{Nonce: "wire-e2e-nonce"},
			Src:       "/datastar.js",
		}).Render(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `<div class="p-4 flex flex-wrap gap-3 items-start">`); err != nil {
			return err
		}

		htmxButton := display.ButtonProps{
			BaseProps: utils.BaseProps{ID: "btn-wire-htmx", Class: "mr-4"},
			Text:      "Load via htmx",
			Variant:   display.ButtonSecondary,
			Size:      display.ButtonSizeSM,
			Wire: &wire.Action{
				URL:    "/api/wire/fragment",
				Target: "#wire-htmx-out",
			},
		}
		if err := display.Button(htmxButton).Render(ctx, w); err != nil {
			return err
		}

		datastarButton := display.ButtonProps{
			BaseProps: utils.BaseProps{ID: "btn-wire-datastar"},
			Text:      "Load via Datastar",
			Variant:   display.ButtonSecondary,
			Size:      display.ButtonSizeSM,
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

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		// htmx loads inline (self-host) and processes hx-* nodes on
		// DOMContentLoaded; "complete" readyState guarantees both.
		chromedp.Poll(`document.readyState==='complete' && window.htmx!==undefined`, &htmxDefined),
		chromedp.Click("#btn-wire-htmx", chromedp.NodeVisible),
		chromedp.Poll(`document.querySelector('#wire-htmx-out').innerHTML.length>0`, &fragment),
		chromedp.InnerHTML("#wire-htmx-out", &out, chromedp.NodeVisible),
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

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		// The pinned runtime dispatches datastar-ready on document when its
		// engine booted; before that, data-on:* clicks are inert.
		chromedp.Poll(`window.__dsReady===true`, &dsReady),
		chromedp.Click("#btn-wire-datastar", chromedp.NodeVisible),
		chromedp.Poll(`document.querySelector('#wire-datastar-out').innerHTML.length>0`, &fragment),
		chromedp.InnerHTML("#wire-datastar-out", &out, chromedp.NodeVisible),
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
