package visualtest

import (
	"context"
	"testing"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

// setPrintMedia emulates CSS media type "print" for the tab (raw CDP call —
// this chromedp major has no typed action for it). An empty media string
// disables the override, restoring screen media.
func setPrintMedia(ctx context.Context, media string) error {
	return chromedp.Do(ctx, chromedp.Func(func(ctx context.Context, tgt *chromedp.Target) error {
		_, err := cdp.Call(
			ctx,
			tgt,
			emulation.SetEmulatedMedia,
			emulation.SetEmulatedMediaParams{Media: media},
		)

		return err
	}))
}

// TestPrintMediaHidesInteractiveChrome proves the global print layer
// (templates/custom.css @media print) actually applies in a real browser
// under print emulation: copy buttons and the theme toggle are display:none,
// and shadows are stripped. String tests cannot see media-query application —
// this is the browser proof (TODO #424 spot-check).
func TestPrintMediaHidesInteractiveChrome(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := setPrintMedia(ctx, "print"); err != nil {
		t.Fatalf("visualtest[print]: emulate print media: %v", err)
	}

	if err := chromedp.Do(ctx,
		chromedp.Navigate(server.BaseURL()+"/"),
		chromedp.WaitReady("body"),
		pollTrue(`!!document.querySelector('[data-tc-copy]')`),
	); err != nil {
		t.Fatalf("visualtest[print]: navigate home: %v", err)
	}

	assertions := []struct {
		name     string
		selector string
	}{
		{"copy button hidden on paper", "[data-tc-copy]"},
		{"theme toggle hidden on paper", "[data-theme-toggle]"},
	}

	for _, a := range assertions {
		t.Run(a.name, func(t *testing.T) {
			var display string

			if err := chromedp.Do(ctx, evalExprInto(
				`getComputedStyle(document.querySelector('`+a.selector+`')).display`, &display)); err != nil {
				t.Fatalf("visualtest[print]: computed display of %s: %v", a.selector, err)
			}

			if display != "none" {
				t.Errorf(
					"visualtest[print]: %s computes to display:%s under print emulation — the @media print layer in templates/custom.css is not applying (stale compiled demo CSS?)",
					a.selector,
					display,
				)
			}
		})
	}

	t.Run("shadows stripped on paper", func(t *testing.T) {
		var boxShadow string

		if err := chromedp.Do(ctx, evalExprInto(
			`(() => { const el = document.querySelector('[class*="shadow"]'); return el ? getComputedStyle(el).boxShadow : 'none'; })()`,
			&boxShadow)); err != nil {
			t.Fatalf("visualtest[print]: computed boxShadow: %v", err)
		}

		if boxShadow != "none" {
			t.Errorf(
				"visualtest[print]: first shadow-bearing element keeps box-shadow %q on paper — the print layer's shadow reset is not applying",
				boxShadow,
			)
		}
	})

	server.FailIfServerErrors(t)
}
