package visualtest

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/utils/wire"
)

func TestZZDatastarWizardReplica(t *testing.T) {
	srv := packE2EServer(t)

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

	if err := packSubmitUntil(ctx, region, region, packWizardNameBad); err != nil {
		t.Fatalf("step 1 invalid: %v", err)
	}
	t.Log("phase step1-error OK")

	if err := chromedp.Run(ctx,
		waitSwapSettled(),
		setFieldValue(ctx, region, `input[name="name"]`, "Ada Lovelace"),
	); err != nil {
		t.Fatalf("step 1 fill: %v", err)
	}

	if err := packSubmitUntil(ctx, region, region, "Wizard complete"); err != nil {
		t.Fatalf("step 1 complete: %v", err)
	}
	t.Log("phase complete OK")
}
