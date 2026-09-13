package visualtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/utils/wire"
)

// TestCalendarNavDiag is a throwaway diagnostic for the htmx MonthNav
// second-click failure: probes whether htmx processed the swapped-in anchors.
func TestCalendarNavDiag(t *testing.T) {
	dialect := wire.TransportHTMX
	srv := calendarNavServer(t, dialect, &requestLog{})

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 40*time.Second)
	defer cancelTimeout()

	var consoleBuf []string

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if e, ok := ev.(*runtime.ConsoleAPICalledEvent); ok {
			consoleBuf = append(consoleBuf, fmt.Sprintf("%v %v", e.Type, e.Args))
		}
	})

	var ok bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(packGate(dialect), &ok),
	); err != nil {
		t.Fatalf("setup: %v", err)
	}

	var done, h3 string

	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector('a[aria-label="Next month"]').dispatchEvent(new MouseEvent('click',{bubbles:true,cancelable:true}))&&''`, &done),
		chromedp.Poll(`document.querySelector('#cal-nav h3')?.textContent.includes('August')?'ok':''`, &done),
	); err != nil {
		t.Fatalf("next click: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector('#cal-nav h3')?.textContent`, &h3)); err != nil {
		t.Fatalf("h3 read: %v", err)
	}

	t.Logf("h3 after next: %q", h3)

	// Raw prev click, no manual processing.
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`document.querySelector('a[aria-label="Previous month"]').dispatchEvent(new MouseEvent('click',{bubbles:true,cancelable:true}))&&''`, &done),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`document.querySelector('#cal-nav h3')?.textContent`, &h3),
	); err != nil {
		t.Fatalf("raw prev: %v", err)
	}

	t.Logf("h3 after raw prev click (no process): %q", h3)

	// Manual htmx.process then prev click again.
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`htmx.process(document.getElementById('cal-nav'))&&''`, &done),
		chromedp.Evaluate(`document.querySelector('a[aria-label="Previous month"]').dispatchEvent(new MouseEvent('click',{bubbles:true,cancelable:true}))&&''`, &done),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`document.querySelector('#cal-nav h3')?.textContent`, &h3),
	); err != nil {
		t.Fatalf("processed prev: %v", err)
	}

	t.Logf("h3 after htmx.process + prev click: %q", h3)
	t.Logf("console: %v", consoleBuf)
}
