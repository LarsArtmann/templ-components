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

// TestWaitAnimationsSettled exercises the code paths in
// waitAnimationsSettled: (1) page with no animations returns after the
// registration window, (2) finished animations return, (3) infinite running
// animations are filtered out (never block), (4) finite running animations
// are waited out. All paths must return without error.
func TestWaitAnimationsSettled(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		page   string
		assert func(t *testing.T, elapsed time.Duration)
	}{
		{
			name: "no_animations",
			page: `<!DOCTYPE html><html><body>
<div id="test" style="width:100px;height:100px;background:blue"></div>
</body></html>`,
			assert: func(t *testing.T, elapsed time.Duration) {
				// Waits through the registration window (~300ms) since no
				// animations appear to confirm the element has no transition.
				assertElapsedUnder(t, "no_animations", elapsed, 500*time.Millisecond,
					"should return after registration window")
			},
		},
		{
			name: "finished_animations",
			page: `<!DOCTYPE html><html><body>
<div id="test" style="width:100px;height:100px;background:red;
  animation: instant 0.01s forwards;"></div>
<style>@keyframes instant { to { opacity: 1; } }</style>
</body></html>`,
			assert: func(t *testing.T, elapsed time.Duration) {
				// May wait through the registration window if the short animation
				// has already been cleaned up by the time getAnimations() runs.
				assertElapsedUnder(t, "finished_animations", elapsed, 500*time.Millisecond,
					"should return quickly")
			},
		},
		{
			name: "long_running",
			page: `<!DOCTYPE html><html><body>
<div id="test" style="width:100px;height:100px;background:green;
  animation: spin 10s infinite;"></div>
<style>@keyframes spin { to { transform: rotate(360deg); } }</style>
</body></html>`,
			assert: func(t *testing.T, elapsed time.Duration) {
				// Infinite animations are filtered out — they never finish, so
				// the helper must not block on them. Returns after the
				// registration window.
				assertElapsedUnder(t, "long_running", elapsed, 500*time.Millisecond,
					"infinite animations must not block")
			},
		},
		{
			name: "finite_running",
			page: `<!DOCTYPE html><html><body>
<div id="test" style="width:100px;height:100px;background:orange;
  animation: slide 2s forwards;"></div>
<style>@keyframes slide { to { transform: translateX(50px); } }</style>
</body></html>`,
			assert: func(t *testing.T, elapsed time.Duration) {
				// Finite animations are waited out. The 2s animation may have
				// partially elapsed during navigation, so only bound away from
				// the failure modes: returning at the registration window
				// (~300ms) or burning the full 8s deadline.
				if elapsed < 1*time.Second {
					t.Errorf("finite_running took %v — finite animations must be waited out", elapsed)
				}

				if elapsed > 4*time.Second {
					t.Errorf("finite_running took %v — should return soon after the animation finishes", elapsed)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			elapsed := navigateAndSettleAnimations(t, tc.name, tc.page)
			tc.assert(t, elapsed)
		})
	}
}

// navigateAndSettleAnimations loads one fixture page, waits for its #test
// element, runs waitAnimationsSettled, and returns how long settling took.
func navigateAndSettleAnimations(t *testing.T, name, page string) time.Duration {
	t.Helper()

	ctx, cancel := newTab(t)
	defer cancel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, page)
	}))
	defer srv.Close()

	taskCtx, taskCancel := context.WithTimeout(ctx, 5*time.Second)
	defer taskCancel()

	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(srv.URL),
		chromedp.WaitVisible("#test", chromedp.ByQuery),
	); err != nil {
		t.Fatalf("navigate(%s): %v", name, err)
	}

	start := time.Now()

	if err := chromedp.Run(taskCtx, waitAnimationsSettled()); err != nil {
		t.Fatalf("waitAnimationsSettled(%s): %v", name, err)
	}

	return time.Since(start)
}

// assertElapsedUnder fails the test when a settle path exceeded its deadline.
func assertElapsedUnder(t *testing.T, label string, elapsed, limit time.Duration, why string) {
	t.Helper()

	if elapsed > limit {
		t.Errorf("%s took %v — %s", label, elapsed, why)
	}
}
