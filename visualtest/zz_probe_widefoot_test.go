package visualtest

// TEMPORARY manual-verification probe for the sticky-footer + stacked-kanban
// work. Captures demo routes at custom viewports (1920x1080 wide desktop,
// 375x667 mobile) to /tmp for human review. DELETE after review — not part
// of the suite.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func TestZZProbeCustomViewports(t *testing.T) {
	base := demoRouteBase()
	if len(base) < 4 || base[:4] != "http" {
		t.Fatalf("demo server did not start: %s", base)
	}

	wide := Viewport{Width: 1920, Height: 1080}

	cases := []struct {
		name     string
		path     string
		viewport Viewport
		fullPage bool
	}{
		{"probe_users_wide_full", "/users", wide, true},
		{"probe_users_wide_fold", "/users", wide, false},
		{"probe_kanban_wide_fold", "/kanban", wide, false},
		{"probe_users_mobile_fold", "/users", ViewportMobile, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := newTab(t)
			defer cancel()

			timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 40*time.Second)
			defer cancelTimeout()

			tasks := []chromedp.Action[chromedp.Void]{
				chromedp.EmulateViewport(int64(tc.viewport.Width), int64(tc.viewport.Height)),
				chromedp.Navigate(base + tc.path),
				evalVoid(`try { localStorage.setItem('theme', "light"); } catch (e) {}`),
				chromedp.Reload(),
				chromedp.WaitVisible(chromedp.CSS("body")),
				waitAnimationsSettled(),
				chromedp.Sleep(settleDelay),
			}

			if err := chromedp.Do(timeoutCtx, tasks...); err != nil {
				t.Fatalf("capture: %v", err)
			}

			// quality 100 = PNG (any other value yields JPEG).
			shot, err := chromedp.Run(timeoutCtx, chromedp.FullScreenshot(100))
			if err != nil {
				t.Fatalf("shot: %v", err)
			}

			out := filepath.Join("/tmp", tc.name+".png")
			if err := os.WriteFile(out, shot, 0o644); err != nil {
				t.Fatalf("write: %v", err)
			}

			fmt.Printf("probe wrote %s (%d bytes)\n", out, len(shot))
		})
	}
}
