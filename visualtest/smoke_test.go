//go:build smoke

package visualtest_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestSmoke verifies chromedp can find and drive a headless Chromium.
func TestSmoke(t *testing.T) {
	chromePath := os.Getenv("CHROMEDP_CHROME_PATH")
	if chromePath == "" {
		t.Skip("set CHROMEDP_CHROME_PATH to run smoke test")
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.NoSandbox,
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var title string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,<html><head><title>tc</title></head><body><h1 id=h>hello</h1></body></html>"),
		chromedp.Text("#h", &title, chromedp.ByQuery),
	); err != nil {
		t.Fatalf("chromedp run: %v", err)
	}
	if title != "hello" {
		t.Fatalf("got %q want hello", title)
	}
	fmt.Println("SMOKE OK: chromedp drove Chromium at", chromePath)
}
