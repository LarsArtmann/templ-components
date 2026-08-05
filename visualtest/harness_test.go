package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestWaitAnimationSettled exercises the three code paths in
// waitAnimationSettled: (1) element with no animations returns immediately,
// (2) element with finished animations returns, (3) element with running
// animations times out gracefully. All paths must return without error.
func TestWaitAnimationSettled(t *testing.T) {
	t.Parallel()

	ctx, cancel := newTab(t)
	defer cancel()

	pages := map[string]string{
		"no_animations": `<!DOCTYPE html><html><body>
<div id="test" style="width:100px;height:100px;background:blue"></div>
</body></html>`,
		"finished_animations": `<!DOCTYPE html><html><body>
<div id="test" style="width:100px;height:100px;background:red;
  animation: instant 0.01s forwards;"></div>
<style>@keyframes instant { to { opacity: 1; } }</style>
</body></html>`,
		"long_running": `<!DOCTYPE html><html><body>
<div id="test" style="width:100px;height:100px;background:green;
  animation: spin 10s infinite;"></div>
<style>@keyframes spin { to { transform: rotate(360deg); } }</style>
</body></html>`,
	}

	for name, page := range pages {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, page)
			}))
			defer srv.Close()

			taskCtx, taskCancel := context.WithTimeout(ctx, 5*time.Second)
			defer taskCancel()

			start := time.Now()

			err := chromedp.Run(taskCtx,
				chromedp.Navigate(srv.URL),
				chromedp.WaitVisible("#test", chromedp.ByQuery),
				waitAnimationSettled("#test"),
			)

			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("waitAnimationSettled(%s): %v", name, err)
			}

			switch name {
			case "no_animations":
				// Should return after the initial 80ms sleep.
				if elapsed > 300*time.Millisecond {
					t.Errorf("no_animations took %v — should return almost immediately after initial sleep", elapsed)
				}
			case "finished_animations":
				// Should return after the initial 80ms sleep + one poll.
				if elapsed > 400*time.Millisecond {
					t.Errorf("finished_animations took %v — should return quickly", elapsed)
				}
			case "long_running":
				// Should poll for the full 800ms deadline before giving up.
				if elapsed < 800*time.Millisecond {
					t.Errorf("long_running took %v — should wait at least 800ms before timeout", elapsed)
				}

				if elapsed > 1200*time.Millisecond {
					t.Errorf("long_running took %v — should not exceed ~900ms total (80ms sleep + 800ms poll)", elapsed)
				}
			}
		})
	}
}
