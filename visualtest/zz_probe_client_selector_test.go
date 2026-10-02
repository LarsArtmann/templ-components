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
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

func TestProbeClientSelectorTargeting(t *testing.T) {
	props := layout.DefaultPageProps()
	props.Title = "probe"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`)
		return err
	})

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := datastar.SDKScript(datastar.SDKScriptProps{
			BaseProps: utils.BaseProps{Nonce: "probe"},
			Src:       "/datastar.js",
		}).Render(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `<div class="p-4">`); err != nil {
			return err
		}

		btn := display.ButtonProps{
			BaseProps: utils.BaseProps{ID: "btn-probe"},
			Text:      "probe",
			Wire: &wire.Action{
				Transport: wire.TransportDatastar,
				URL:       "/api/probe",
				Selector:  "#probe-region",
			},
		}
		if err := display.Button(btn).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `<div id="probe-region">probe-sentinel</div></div>`)
		return err
	})

	page := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})

	css, errCSS := loadCSS()
	if errCSS != nil {
		t.Fatalf("css: %v", errCSS)
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
	mux.HandleFunc("/api/probe", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<p id="other-root">probe-fragment</p>`))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = page.Render(context.Background(), w)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var ready, regionGone bool
	var regionHTML string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollBool(`window.__dsReady===true`, &ready),
		chromedp.Click("#btn-probe", chromedp.NodeVisible),
		chromedp.Sleep(1500*time.Millisecond),
		pollBool(`!document.querySelector('#probe-region')`, &regionGone),
		chromedp.Evaluate(`document.body.innerHTML`, &regionHTML),
	); err != nil {
		t.Logf("poll error (expected if region survived): %v", err)
	}

	t.Logf("regionGone=%v regionHTML=%s", regionGone, regionHTML)
}
