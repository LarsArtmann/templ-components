package visualtest

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Scratch probe 2: isolates the ConfirmDelete click wedge. DELETE AFTER DIAGNOSIS.
func TestScratchConfirmDeleteProbe(t *testing.T) {
	server := StartDemoServer(t)

	chromePath := os.Getenv("CHROMEDP_CHROME_PATH")
	if chromePath == "" {
		t.Skip("no chrome")
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.NoSandbox,
		chromedp.DisableGPU,
	)

	debugCtx, debugCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer debugCancel()

	ctx, cancel := chromedp.NewContext(debugCtx, chromedp.WithDebugf(t.Logf))
	defer cancel()

	ctx, timeoutCancel := context.WithTimeout(ctx, 40*time.Second)
	defer timeoutCancel()

	chromedp.ListenTarget(ctx, func(ev any) {
		if _, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			go func() {
				_ = page.HandleJavaScriptDialog(true).Do(ctx)
			}()
		}
	})

	t.Log("phase 1: navigate")

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("navigate: %v", err)
	}

	t.Log("phase 2: sanity evaluate")

	var one int

	if err := chromedp.Run(ctx, chromedp.Evaluate(`1`, &one)); err != nil {
		t.Fatalf("evaluate: %v", err)
	}

	t.Log("phase 3: click #item-123 button[hx-delete]")

	if err := chromedp.Run(ctx, chromedp.Click(`#item-123 button[hx-delete]`, chromedp.ByQuery)); err != nil {
		t.Fatalf("click: %v", err)
	}

	t.Log("phase 4: click returned — poll for replacement")

	var removed bool

	if err := chromedp.Run(ctx, chromedp.Poll(
		`Boolean(document.querySelector('#item-123') && document.querySelector('#item-123').innerText.indexOf('deleted successfully') >= 0)`,
		&removed,
		chromedp.WithPollingTimeout(15*time.Second),
	)); err != nil {
		t.Fatalf("poll: %v", err)
	}

	t.Log("phase 4 ok: delete flow complete")
}
