package visualtest

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// TestKeyboardTraversalFocusVisibility is the automated half of the
// keyboard-only audit (TODO #197, residue of #175): Tab through each route
// and assert focus NEVER lands on an invisible or zero-size element — the
// classic keyboard-user trap (sr-only triggers, collapsed panels, off-screen
// drawers) that axe's DOM-level rules cannot see.
//
// What it deliberately does NOT automate: the *ordering* quality of the tab
// sequence (DOM order vs visual order) and visible focus-indicator styling
// on every widget — those need human traversal per
// docs/testing/a11y-gate-policy.md; this test narrows where humans look.
func TestKeyboardTraversalFocusVisibility(t *testing.T) {
	t.Parallel()

	server := StartDemoServer(t)

	routes := []struct {
		path     string
		minFocus int
	}{
		{"/", 5},
		{"/forms", 5},
		{"/users", 5},
		{"/errors/full", 2}, // a minimal error page: go-home + go-back links only
	}

	const maxTabs = 60

	for _, route := range routes {
		t.Run(route.path, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := newTab(t)
			defer cancel()

			timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 90*time.Second)
			defer cancelTimeout()

			chromedp.Run(timeoutCtx,
				chromedp.Navigate(server.BaseURL()+route.path),
				chromedp.WaitVisible("body", chromedp.ByQuery),
				chromedp.Evaluate(`document.activeElement && document.activeElement.blur();`, nil),
			)

			visited := 0

			for i := range maxTabs {
				var info struct {
					Tag      string  `json:"tag"`
					ID       string  `json:"id"`
					Classes  string  `json:"classes"`
					Visible  bool    `json:"visible"`
					RectArea float64 `json:"area"`
				}

				if err := chromedp.Run(timeoutCtx,
					chromedp.KeyEvent("\t"),
					// Let any scroll-into-view settle before measuring.
					chromedp.Sleep(30*time.Millisecond),
					chromedp.Evaluate(`(() => {
						const el = document.activeElement;
						if (!el || el === document.body) return {tag: ''};
						const r = el.getBoundingClientRect();
						const style = getComputedStyle(el);
						const visible = style.visibility !== 'hidden'
							&& style.display !== 'none'
							&& Number(style.opacity) > 0
							&& r.width > 0 && r.height > 0;
						return {
							tag: el.tagName,
							id: el.id || '',
							classes: (el.getAttribute && el.getAttribute('class')) || '',
							visible: visible,
							area: r.width * r.height
						};
					})()`, &info),
				); err != nil {
					t.Fatalf("tab #%d on %s: %v", i+1, route.path, err)
				}

				if info.Tag == "" {
					if visited == 0 {
						t.Fatal("Tab never moved focus off document.body — the key dispatch is broken")
					}

					// Focus wrapped past the last element — cycle complete.
					break
				}

				visited++

				if !info.Visible {
					t.Errorf("tab #%d on %s landed focus on an INVISIBLE element: <%s id=%q class=%q>",
						i+1, route.path, strings.ToLower(info.Tag), info.ID, info.Classes)
				}
			}

			if visited < route.minFocus {
				t.Errorf(
					"only %d focusable element(s) visited on %s — traversal too short to be a real audit",
					visited,
					route.path,
				)
			}
		})
	}
}
