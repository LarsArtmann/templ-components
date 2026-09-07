package visualtest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/utils/wire"
)

func TestZZDatastarWizardReplica(t *testing.T) {
	srv := packE2EServer(t)

	var (
		mu     sync.Mutex
		bodies []string
	)
	orig := srv.Config.Handler
	srv.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "wizard") && r.Method == http.MethodPost {
			b, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(b))
			mu.Lock()
			bodies = append(bodies, string(b))
			mu.Unlock()
		}
		orig.ServeHTTP(w, r)
	})
	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()
		t.Logf("WIZARD BODIES: %q", bodies)
	})

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(os.Getenv("CHROMEDP_CHROME_PATH")),
		chromedp.NoSandbox,
		chromedp.DisableGPU,
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, tabCancel := chromedp.NewContext(allocCtx)
	defer tabCancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 90*time.Second)
	defer cancelTimeout()

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if e, ok := ev.(*network.EventResponseReceived); ok {
			t.Logf("NET %d %s", e.Response.Status, e.Response.URL)
		}
	})

	region := packWizardDatastarRegion
	dialect := wire.TransportDatastar

	var ok bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(packGate(dialect), &ok),
	); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := packSubmitUntil(ctx, region, region, packWizardEmailBad); err != nil {
		t.Fatalf("step 0 invalid: %v", err)
	}
	t.Log("phase step0-error OK")

	if err := chromedp.Run(ctx,
		chromedp.Poll(regionExistsExpr(region, `input[name="email"]`), &ok),
		waitSwapSettled(),
		setFieldValue(ctx, region, `input[name="email"]`, "ada@example.com"),
	); err != nil {
		t.Fatalf("step 0 fill: %v", err)
	}

	if err := packSubmitUntil(ctx, region, region, "Full name"); err != nil {
		t.Fatalf("step 0 advance: %v", err)
	}
	t.Log("phase step1 OK")

	if err := chromedp.Run(ctx, waitSwapSettled()); err != nil {
		t.Fatalf("settle: %v", err)
	}

	for attempt := 1; attempt <= 4; attempt++ {
		if err := chromedp.Run(ctx,
			chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
		); err != nil {
			t.Fatalf("click attempt %d: %v", attempt, err)
		}

		if packPollOnce(ctx, regionHasText(region, packWizardNameBad)) {
			t.Logf("attempt %d: needle OK", attempt)

			break
		}

		var probe string
		if err := chromedp.Run(ctx, chromedp.Evaluate(`JSON.stringify({
			txt: document.querySelector('#pack-wizard-datastar-region').innerText.slice(0, 120),
			forms: [...document.querySelectorAll('form')].filter(f => f.closest('#pack-wizard-datastar-region')).map(f => f.querySelector('input[name="step"]') ? f.querySelector('input[name="step"]').value : 'none'),
			dsAttr: (document.querySelector('#pack-wizard-datastar-region form') || {hasAttribute: function(){return false;}}).hasAttribute('data-on:submit'),
			url: location.href
		})`, &probe)); err != nil {
			t.Fatalf("probe attempt %d: %v", attempt, err)
		}

		t.Logf("attempt %d NO NEEDLE: %s", attempt, probe)
	}
}
