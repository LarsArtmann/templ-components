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
	"github.com/larsartmann/go-datastar/static"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/layout"
)

// The Datastar LoadingButton browser E2E proves the indicator-signal label
// swap against the PINNED runtime: clicking a button wired with
// data-indicator:<signal> + a slow endpoint must show the busy label while
// the request is in flight and restore the rest label when it lands.
// The swap is runtime-driven (data-show) — this test pins that the runtime
// actually toggles the signal the component binds to.
//
// Run via `nix run .#visual`; tests skip gracefully without a browser.

func datastarLoadingButtonE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	props := layout.DefaultPageProps()
	props.Title = "Datastar LoadingButton E2E — templ-components"
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
				Nonce: "lbtn-e2e-nonce",
				Src:   "/datastar.js",
			}),
		} {
			if err := component.Render(ctx, w); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(
			w,
			`<button id="save-btn" type="button" data-indicator:saving data-on:click="@post('/api/slow-save')">`,
		); err != nil {
			return err
		}

		if err := datastar.LoadingButton(datastar.LoadingButtonProps{
			Signal:      "saving",
			DefaultText: "Save",
			LoadingText: "Saving…",
		}).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</button>`)

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

	mux.HandleFunc("/api/slow-save", func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(1200 * time.Millisecond)
		w.WriteHeader(http.StatusNoContent)
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

// TestDatastarLoadingButtonBrowser proves the busy/rest label swap under the
// pinned runtime: at rest only "Save" is visible; clicking fires the slow
// endpoint and "Saving…" + spinner become visible while "Save" hides; the
// response restores "Save".
func TestDatastarLoadingButtonBrowser(t *testing.T) {
	t.Parallel()

	srv := datastarLoadingButtonE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	// visibleTextIn returns a predicate asserting the given text is inside a
	// non-display:none subtree of the button.
	visibleTextIn := func(text string) string {
		return `(() => {
			const spans = document.querySelectorAll('#save-btn span[data-show]');
			for (const s of spans) {
				if (s.textContent.includes('` + text + `') && s.offsetParent !== null) return true;
			}
			return false;
		})()`
	}

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollTrue(`window.__dsReady===true && document.querySelector('#save-btn')!==null`),
		// At rest: "Save" visible, "Saving…" hidden.
		pollTrue(visibleTextIn("Save")),
		pollTrue("!("+visibleTextIn("Saving…")+")"),
		// Click → the 1.2s endpoint keeps the busy state observable.
		chromedp.Click("#save-btn", chromedp.NodeVisible),
		pollTrue(visibleTextIn("Saving…")),
		pollTrue("!("+visibleTextIn("Save")+")"),
		// Response lands → rest state restored.
		pollTrue(visibleTextIn("Save")),
		pollTrue("!("+visibleTextIn("Saving…")+")"),
	); err != nil {
		t.Fatalf("LoadingButton busy/rest swap proof: %v", err)
	}
}
