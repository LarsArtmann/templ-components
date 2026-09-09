package visualtest

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// Scratch probe: isolates the LoadMore flow hang. DELETE AFTER DIAGNOSIS.
func TestScratchLoadMoreProbe(t *testing.T) {
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

	ctx, timeoutCancel := context.WithTimeout(ctx, 45*time.Second)
	defer timeoutCancel()

	t.Log("phase 1: navigate")

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("navigate: %v", err)
	}

	t.Log("phase 2: sanity evaluate 1+1")

	var sum int

	if err := chromedp.Run(ctx, chromedp.Evaluate(`1+1`, &sum)); err != nil {
		t.Fatalf("evaluate 1+1: %v", err)
	}

	t.Logf("phase 2 ok: sum=%d", sum)

	t.Log("phase 3: count #demo-load-more buttons via Evaluate")

	var count int

	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`document.querySelectorAll('#demo-load-more button').length`, &count)); err != nil {
		t.Fatalf("evaluate count: %v", err)
	}

	t.Logf("phase 3 ok: buttons=%d", count)

	t.Log("phase 4: chromedp.Click on the button")

	if err := chromedp.Run(ctx, chromedp.Click(`#demo-load-more button`, chromedp.ByQuery)); err != nil {
		t.Fatalf("click: %v", err)
	}

	t.Log("phase 4 ok: click returned")

	var html string

	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`document.documentElement.outerHTML`, &html)); err != nil {
		t.Fatalf("dump html: %v", err)
	}

	if werr := os.WriteFile("/tmp/tc-demo-index.html", []byte(html), 0o600); werr != nil {
		t.Fatalf("write html: %v", werr)
	}

	t.Logf("phase 4b: dumped %d bytes to /tmp/tc-demo-index.html", len(html))

	t.Log("phase 5: poll for >= 4 cards")

	var got string

	if err := chromedp.Run(ctx, chromedp.Poll(
		`document.querySelectorAll('#demo-load-more .rounded-lg').length >= 4`, &got,
		chromedp.WithPollingTimeout(10*time.Second))); err != nil {
		t.Fatalf("poll cards: %v", err)
	}

	t.Log("phase 5 ok: poll satisfied — LoadMore flow is healthy")
}
