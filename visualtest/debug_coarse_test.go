package visualtest

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

// TestDebugCoarseEmulation prints whether CDP media emulation flips
// matchMedia and the computed opacity. Scratch test — delete after use.
func TestDebugCoarseEmulation(t *testing.T) {
	srv := kanbanE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var ready bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(kanbanE2EReady, &ready),
	); err != nil {
		t.Fatalf("navigate: %v", err)
	}

	var matchBefore bool
	if err := chromedp.Run(ctx, chromedp.Evaluate(`matchMedia('(pointer:coarse)').matches`, &matchBefore)); err != nil {
		t.Fatalf("eval before: %v", err)
	}
	t.Logf("matchMedia coarse BEFORE: %v", matchBefore)

	err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return emulation.SetEmulatedMedia().
			WithFeatures([]*emulation.MediaFeature{
				{Name: "pointer", Value: "coarse"},
			}).
			Do(ctx)
	}))
	if err != nil {
		t.Fatalf("emulate: %v", err)
	}

	var matchAfter, matchHover bool
	if err := chromedp.Run(ctx, chromedp.Evaluate(`matchMedia('(pointer:coarse)').matches`, &matchAfter)); err != nil {
		t.Fatalf("eval after: %v", err)
	}
	if err := chromedp.Run(ctx, chromedp.Evaluate(`matchMedia('(hover:none)').matches`, &matchHover)); err != nil {
		t.Fatalf("eval hover: %v", err)
	}
	t.Logf("matchMedia coarse AFTER: %v, hover:none: %v", matchAfter, matchHover)

	var opacity string
	if err := chromedp.Run(ctx, chromedp.Evaluate(kanbanButtonsOpacityExpr, &opacity)); err != nil {
		t.Fatalf("eval opacity: %v", err)
	}
	t.Logf("computed opacity AFTER: %q", opacity)
}
