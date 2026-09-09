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
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/htmx"
	"github.com/larsartmann/templ-components/layout"
)

// The LoadingButton E2E pins the browser-level contract the HTML goldens
// cannot: htmx itself must flip the .htmx-request class and drive the
// htmx-indicator opacity gates. At rest the spinner/loading text are
// invisible and the default text is shown; during the request they invert.
// The endpoint sleeps 600ms (same as the demo's /api/save) so the in-flight
// assertions are deterministic.
func loadingButtonE2EPage() templ.Component {
	props := layout.DefaultPageProps()
	props.Title = "LoadingButton E2E — templ-components"
	props.CSSPath = "/app.css"

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(
			w,
			`<button id="btn-loading-e2e" hx-post="/api/save-slow" hx-swap="none" class="inline-flex items-center rounded-md bg-blue-600 dark:bg-blue-500 px-4 py-2 text-sm font-medium text-white">`,
		); err != nil {
			return err
		}

		if err := htmx.LoadingButton("Save changes", "Saving...", feedback.Spinner(feedback.SpinnerProps{Size: feedback.SpinnerSM, Color: "text-white"})).
			Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</button>`)

		return err
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

func loadingButtonE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	return e2ePageServer(t, loadingButtonE2EPage, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/save-slow", func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(600 * time.Millisecond)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(`<span class="text-sm text-green-700 dark:text-green-300">Saved.</span>`))
		})
	})
}

func TestLoadingButtonE2EStateGatesDuringRequest(t *testing.T) {
	t.Parallel()

	srv := loadingButtonE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	buttonJS := `document.querySelector('#btn-loading-e2e')`
	spinnerOpacityJS := `parseFloat(getComputedStyle(document.querySelector('#btn-loading-e2e .tc-btn-loading .htmx-indicator')).opacity)`
	defaultTextDisplayJS := `getComputedStyle(document.querySelector('#btn-loading-e2e .tc-btn-loading [class*="htmx-request_"]')).display`

	var (
		htmxReady, requesting, settled, spinnerGatedAgain bool
		spinnerAtRest                                     float64
		defaultTextDuring                                 string
	)

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(`document.readyState==='complete' && window.htmx!==undefined`, &htmxReady),
		chromedp.Evaluate(spinnerOpacityJS, &spinnerAtRest),
		chromedp.Click("#btn-loading-e2e", chromedp.NodeVisible),
		chromedp.Poll(buttonJS+`.classList.contains('htmx-request')`, &requesting),
		chromedp.Poll(spinnerOpacityJS+`>0.99`, nil),
		chromedp.Evaluate(defaultTextDisplayJS, &defaultTextDuring),
		chromedp.Poll(`!(`+buttonJS+`.classList.contains('htmx-request'))`, &settled),
		chromedp.Poll(spinnerOpacityJS+`<0.01`, &spinnerGatedAgain),
	); err != nil {
		t.Fatalf("LoadingButton E2E: %v", err)
	}

	if spinnerAtRest != 0 {
		t.Errorf("at rest: spinner must be gated (opacity 0); got %v", spinnerAtRest)
	}

	if defaultTextDuring != "none" {
		t.Errorf("during request: default text must be hidden; got display=%q", defaultTextDuring)
	}

	if !spinnerGatedAgain {
		t.Error("after response: spinner did not return to opacity 0")
	}
}
