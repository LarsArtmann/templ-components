package visualtest

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// TestDemoCSPNonceIntegrity pins the demo's CSP nonce contract end-to-end:
// the demo stamps a nonce-requiring CSP header (demoCSP, 'nonce-demo-nonce'),
// so every inline <script> MUST carry that nonce. A component built without
// its nonce prop renders <script nonce=""> — well-formed HTML that the
// browser silently refuses to execute. That exact regression shipped live on
// 2026-10-01 (the MPA rewrite dropped demoNonceConst at 13 call sites and
// killed CopyButton/Kanban/Tooltip/Popover/menu-nav/image-fallback JS on
// every page) because the library's integration/csp_nonce_test.go only
// asserts attribute PRESENCE — an empty value passes it. This guard asserts
// the absence of the broken shape on every rendered demo page.
func TestDemoCSPNonceIntegrity(t *testing.T) {
	server := StartDemoServer(t)

	client := &http.Client{}
	emptyNonce := regexp.MustCompile(`<script nonce="">`)

	pages := []string{
		"/", "/layout", "/display", "/feedback", "/forms", "/navigation",
		"/icons", "/htmx", "/datastar", "/wire", "/kanban", "/echarts",
		"/recipes", "/users", "/error-pages", "/recipes/dashboard",
		"/errors/404",
	}

	for _, page := range pages {
		t.Run(page, func(t *testing.T) {
			t.Parallel()

			resp, err := client.Get(server.BaseURL() + page)
			if err != nil {
				t.Fatalf("visualtest[csp]: GET %s: %v", page, err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("visualtest[csp]: read %s: %v", page, err)
			}

			csp := resp.Header.Get("Content-Security-Policy")
			if !strings.Contains(csp, "'nonce-demo-nonce'") {
				t.Errorf("visualtest[csp]: %s CSP lacks the demo nonce source: %q", page, csp)
			}

			if hits := emptyNonce.FindAllString(string(body), -1); len(hits) > 0 {
				t.Errorf(
					"visualtest[csp]: %s renders %d script(s) with an EMPTY nonce — the demo CSP blocks them, killing their JS. The component was built without its nonce prop (demoNonceConst).",
					page,
					len(hits),
				)
			}
		})
	}
}

// TestDemoCSPJSExecutes is the behavioral twin of TestDemoCSPNonceIntegrity:
// the HTTP-level sweep above catches the nonce="" shape, but a browser-level
// assertion proves scripts actually RUN (CSP headers could regress in other
// ways — wrong nonce value, a second intersecting policy). It loads one
// shell page and checks that the Kanban singleton script attached, which
// only happens when inline execution is allowed.
func TestDemoCSPJSExecutes(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/kanban"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[csp]: load kanban: %v", err)
	}

	var attached bool
	if err := chromedp.Run(ctx, chromedp.Evaluate(`window.tcKanbanAttached===true`, &attached)); err != nil {
		t.Fatalf("visualtest[csp]: evaluate: %v", err)
	}

	if !attached {
		t.Error(
			"visualtest[csp]: window.tcKanbanAttached never became true — inline scripts are not executing under the demo CSP",
		)
	}
}
