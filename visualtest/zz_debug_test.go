package visualtest

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/utils/wire"
)

func TestZZOwnBrowserDatastarWizard(t *testing.T) {
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

	ctx, cancelTimeout := context.WithTimeout(ctx, 45*time.Second)
	defer cancelTimeout()

	region := packWizardDatastarRegion
	dialect := wire.TransportDatastar

	var ok bool
	var probe string

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if e, isCT := ev.(*runtime.EventConsoleAPICalled); isCT {
			for _, a := range e.Args {
				t.Logf("CONSOLE %s: %s", e.Type, a.Value)
			}
		}

		if e, isEX := ev.(*runtime.EventExceptionThrown); isEX {
			t.Logf("EXCEPTION: %s", e.ExceptionDetails.Text)
		}
	})

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(packGate(dialect), &ok),
		chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
		chromedp.Poll(regionHasText(region, packWizardEmailBad), &ok),
		chromedp.Evaluate(`JSON.stringify({
			stepForms: [...document.querySelectorAll('form')].map(f => ({
				step: f.querySelector('input[name="step"]') ? f.querySelector('input[name="step"]').value : 'none',
				ds: f.hasAttribute('data-on:submit'),
				inDsRegion: !!f.closest('#pack-wizard-datastar-region'),
			})),
			regionHTML: document.querySelector('#pack-wizard-datastar-region').innerHTML.slice(0, 400),
		})`, &probe),
	); err != nil {
		t.Fatalf("probe after step0 error: %v", err)
	}
	t.Logf("probe-after-req1: %s", probe)

	if err := chromedp.Run(ctx,
		waitSwapSettled(),
		setFieldValue(ctx, region, `input[name="email"]`, "ada@example.com"),
		chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
		chromedp.Sleep(3 * time.Second),
		chromedp.Evaluate(`JSON.stringify({
			stepForms: [...document.querySelectorAll('form')].map(f => ({
				step: f.querySelector('input[name="step"]') ? f.querySelector('input[name="step"]').value : 'none',
				inDsRegion: !!f.closest('#pack-wizard-datastar-region'),
			})),
			regionText2: document.querySelector('#pack-wizard-datastar-region').innerText,
			nameInput: !!document.querySelector('#pack-wizard-datastar-region input[name="name"]'),
		})`, &probe),
	); err != nil {
		t.Fatalf("probe after step1: %v", err)
	}
	t.Logf("probe-after-req2: %s", probe)
}
