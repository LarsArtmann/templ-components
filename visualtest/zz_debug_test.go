package visualtest

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/utils/wire"
)

func TestZZWizardParallelPair(t *testing.T) {
	t.Parallel()
	for _, dialect := range []wire.Transport{wire.TransportHTMX, wire.TransportDatastar} {
		t.Run(string(dialect), func(t *testing.T) {
			t.Parallel()
			srv := packE2EServer(t)
			ctx, cancel := newTab(t)
			defer cancel()
			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			chromedp.ListenTarget(ctx, func(ev interface{}) {
				if e, ok := ev.(*runtime.EventConsoleAPICalled); ok {
					for _, a := range e.Args {
						t.Logf("[%s] console %s: %s", dialect, e.Type, a.Value)
					}
				}
				if e, ok := ev.(*runtime.EventExceptionThrown); ok {
					t.Logf("[%s] exception: %s", dialect, e.ExceptionDetails.Text)
				}
			})

			region := packWizardDatastarRegion
			if dialect == wire.TransportHTMX {
				region = packWizardHTMXRegion
			}
			var ok bool
			if err := chromedp.Run(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.Poll(packGate(dialect), &ok),
				chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionHasText(region, packWizardEmailBad), &ok),
			); err != nil {
				text, textErr := regionText(ctx, region)
				t.Logf("[%s] wizard FAILED: %v", dialect, err)
				text, textErr := regionText(ctx, region)
				t.Logf("[%s] region: %.300s read-err=%v", dialect, text, textErr)
			}
		})
	}
}

func TestZZWizardSerial(t *testing.T) {
	// EXACT copy of the suite subtest, minus t.Parallel.
	srv := packE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	region := packWizardDatastarRegion
	dialect := wire.TransportDatastar

	var ok bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(packGate(dialect), &ok),
		chromedp.Click(formSel(region, `button[type="submit"]`), chromedp.NodeVisible),
		chromedp.Poll(regionHasText(region, packWizardEmailBad), &ok),
	); err != nil {
		text, textErr := regionText(ctx, region)
		t.Logf("serial wizard E2E: %v\nregion text: %.600s (read err: %v)", err, text, textErr)
	}
}
