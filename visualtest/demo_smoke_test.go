package visualtest

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// TestDemoSmokeAllRoutes is CI's demo-health signal: every demo route a
// consumer might copy must render its page (unique <title>) and the server
// must log zero 500s while serving them. Runs in the existing Visual
// Regression job — no separate workflow needed.
func TestDemoSmokeAllRoutes(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	routes := []struct {
		path  string
		title string
	}{
		{"/", "templ-components Demo"},
		{"/forms", "Forms Demo"},
		{"/users", "Users"},
		{"/recipes/dashboard", "Dashboard Recipe"},
		{"/recipes/settings", "Settings Recipe"},
		{"/recipes/login", "Login Recipe"},
		{"/recipes/auth", "Auth Layout Recipe"},
	}

	for _, route := range routes {
		t.Run(route.path, func(t *testing.T) {
			if err := chromedp.Run(ctx,
				chromedp.Navigate(server.BaseURL()+route.path),
				chromedp.WaitReady("body"),
			); err != nil {
				t.Fatalf("visualtest[demo]: navigate %s: %v", route.path, err)
			}

			var title string

			if err := chromedp.Run(ctx, chromedp.Evaluate(`document.title`, &title)); err != nil {
				t.Fatalf("visualtest[demo]: read %s title: %v", route.path, err)
			}

			if !strings.Contains(title, route.title) {
				t.Fatalf("visualtest[demo]: %s rendered title %q, want prefix %q", route.path, title, route.title)
			}
		})
	}

	server.FailIfServerErrors(t)
}
