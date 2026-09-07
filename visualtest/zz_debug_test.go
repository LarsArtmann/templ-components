package visualtest

import (
	"context"
	"os"
	"testing"
	"time"

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
	var one string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(packGate(dialect), &ok),
		chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
		chromedp.Poll(regionHasText(region, packWizardEmailBad), &ok),
		waitSwapSettled(),
		setFieldValue(ctx, region, `input[name="email"]`, "ada@example.com"),
		chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
		chromedp.Poll(regionExistsExpr(region, `input[name="name"]`), &ok),
		waitSwapSettled(),
		chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
		chromedp.Poll(regionHasText(region, packWizardNameBad), &ok),
		chromedp.Evaluate(`'alive-before'`, &one),
	); err != nil {
		t.Fatalf("flow up to step6: %v", err)
	}
	t.Logf("before final click: %s", one)

	err := chromedp.Run(ctx,
		setFieldValue(ctx, region, `input[name="name"]`, "Ada Lovelace"),
		chromedp.Evaluate(`document.querySelector('#pack-wizard-datastar-region button[type="submit"]').click(); 'clicked'`, &one),
		chromedp.Sleep(1 * time.Second),
		chromedp.Evaluate(`'alive-after'`, &one),
	)
	t.Logf("final click run: err=%v value=%q", err, one)

	err2 := chromedp.Run(ctx, chromedp.Evaluate(`'zombie-check'`, &one))
	t.Logf("zombie check: err=%v value=%q", err2, one)
}
