package visualtest

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
)

// The Datastar synthetics browser-prove the two JS paths that string tests
// can only pin textually: the global SSEErrorHandling toast pipeline and the
// LiveRegion busy-cue clearing. The DataStar runtime's document-level
// `datastar-fetch` event is synthesized with CustomEvent — no real SSE
// endpoint is needed.

func datastarSyntheticsServer(t *testing.T) *httptest.Server {
	t.Helper()

	page := func() templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				if err := feedback.ToastContainer("synthetics-nonce").Render(ctx, w); err != nil {
					return err
				}

				if err := datastar.SSEErrorHandling(datastar.DefaultSSEErrorHandlingConfig()).
					Render(ctx, w); err != nil {
					return err
				}

				return datastar.LiveRegion(datastar.LiveRegionProps{
					BaseProps: utils.BaseProps{ID: "synth-live", Nonce: "synthetics-nonce"},
					URL:       "/api/stream",
					AutoStart: true,
				}).Render(ctx, w)
			})

			props := layout.DefaultPageProps()
			props.Title = "Datastar synthetics"

			return layout.Base(props).Render(templ.WithChildren(ctx, content), w)
		})
	}

	return e2ePageServer(t, page, nil)
}

// dispatchDatastarFetch synthesizes the document-level datastar-fetch event
// whose detail mirrors the DataStar runtime's shape.
func dispatchDatastarFetch(detailJSON string) chromedp.Action {
	return chromedp.Evaluate(
		`document.dispatchEvent(new CustomEvent('datastar-fetch', {detail: `+detailJSON+`}));`, nil,
	)
}

func TestDemoDatastarSSEErrorToast(t *testing.T) {
	server := datastarSyntheticsServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[datastar]: navigate: %v", err)
	}

	if err := chromedp.Run(ctx, dispatchDatastarFetch(
		`{type: 'error', argsRaw: {status: 500}}`)); err != nil {
		t.Fatalf("visualtest[datastar]: dispatch error event: %v", err)
	}

	var announcer string

	if err := chromedp.Run(ctx, chromedp.Poll(
		`(document.getElementById('tc-datastar-announcer')||{innerText:''}).innerText`,
		&announcer,
		chromedp.WithPollingTimeout(5*time.Second),
	)); err != nil {
		t.Fatalf("visualtest[datastar]: announcer never spoke: %v", err)
	}

	if announcer == "" {
		t.Fatal("visualtest[datastar]: SSE error event did not set the announcer text")
	}

	var toastCount int64

	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`document.querySelectorAll('#tc-toast-container .pointer-events-auto').length`,
		&toastCount,
	)); err != nil {
		t.Fatalf("visualtest[datastar]: count toasts: %v", err)
	}

	if toastCount == 0 {
		t.Error("visualtest[datastar]: SSE error event produced no toast (tcShowToast path dead)")
	}
}

func TestDemoDatastarRetriesFailedToast(t *testing.T) {
	server := datastarSyntheticsServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[datastar]: navigate: %v", err)
	}

	if err := chromedp.Run(ctx, dispatchDatastarFetch(`{type: 'retries-failed'}`)); err != nil {
		t.Fatalf("visualtest[datastar]: dispatch retries-failed event: %v", err)
	}

	var held bool

	if err := chromedp.Run(ctx, chromedp.Poll(
		`Boolean((document.getElementById('tc-datastar-announcer')||{innerText:''}).innerText.indexOf('Live stream lost') >= 0)`,
		&held,
		chromedp.WithPollingTimeout(5*time.Second),
	)); err != nil {
		t.Fatalf("visualtest[datastar]: retries-failed announcement never appeared: %v", err)
	}
}

func TestDemoDatastarBusyClearsOnFirstPatch(t *testing.T) {
	server := datastarSyntheticsServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[datastar]: navigate: %v", err)
	}

	var busyBefore string

	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`(document.getElementById('synth-live')||{getAttribute:function(){return 'MISSING'}}).getAttribute('aria-busy')`,
		&busyBefore,
	)); err != nil {
		t.Fatalf("visualtest[datastar]: read aria-busy: %v", err)
	}

	if busyBefore != "true" {
		t.Fatalf("visualtest[datastar]: live region should start aria-busy=true, got %q", busyBefore)
	}

	if err := chromedp.Run(ctx, dispatchDatastarFetch(`{type: 'datastar-patch-elements'}`)); err != nil {
		t.Fatalf("visualtest[datastar]: dispatch patch event: %v", err)
	}

	var cleared bool

	if err := chromedp.Run(ctx, chromedp.Poll(
		`Boolean(!(document.getElementById('synth-live').hasAttribute('aria-busy')))`,
		&cleared,
		chromedp.WithPollingTimeout(5*time.Second),
	)); err != nil {
		t.Fatalf("visualtest[datastar]: aria-busy was never cleared after first patch: %v", err)
	}
}
