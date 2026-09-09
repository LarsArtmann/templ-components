package visualtest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
)

// The axe sweep audits the LIVE demo routes — the library's shop window and
// the composition patterns consumers copy — with the vendored axe-core
// runtime. The guard fails on any critical/serious violation that is not
// explicitly accepted in testdata/axe_baseline.json; moderate/minor findings
// are logged but never gate.
//
// Baseline format: {"<route>": {"<rule>|<impact>": <accepted-node-count>, ...}}
// Accepted entries must carry a justification in the audit trail (git history
// of the baseline file) — an accepted violation is documented a11y debt,
// not a pass.

// axeBaselinePath points at the accepted-violations ledger.
const axeBaselinePath = "testdata/axe_baseline.json"

// axeMaxReportedNodes caps the per-violation node report so one widespread
// rule failure cannot flood the test output.
const axeMaxReportedNodes = 5

// axeSweepRoutes enumerates the demo pages under audit. index/forms get light
// + dark passes (the dark palette is where contrast regressions hide); the
// kanban section lives on the index route.
var axeSweepRoutes = []struct {
	name string
	path string
	dark bool
}{
	{name: "index", path: "/", dark: false},
	{name: "index_dark", path: "/", dark: true},
	{name: "forms", path: "/forms", dark: false},
	{name: "forms_dark", path: "/forms", dark: true},
	{name: "recipes_dashboard", path: "/recipes/dashboard", dark: false},
	{name: "recipes_settings", path: "/recipes/settings", dark: false},
	{name: "recipes_login", path: "/recipes/login", dark: false},
	{name: "recipes_auth", path: "/recipes/auth", dark: false},
	{name: "users", path: "/users", dark: false},
}

func TestAxeSweepDemoRoutes(t *testing.T) {
	t.Parallel()

	baseline := readAxeBaseline(t)
	server := StartDemoServer(t)

	for _, route := range axeSweepRoutes {
		t.Run(route.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := newTab(t)
			defer cancel()

			results := axeAuditRoute(t, ctx, server.BaseURL(), route.path, route.dark)

			t.Logf("axe[%s]: %d total violation rule(s), %d blocking", route.name, len(results.Violations), len(results.BlockingViolations()))
			assertNoUnacceptedViolations(t, route.name, baseline, results)
		})
	}
}

// TestAxeHarnessDetectsViolations is the sweep's positive control: a page with
// a textbook violation (image without alt text = critical "image-alt") must
// be caught, proving the harness can fail — a guard that cannot fail guards
// nothing.
func TestAxeHarnessDetectsViolations(t *testing.T) {
	t.Parallel()

	const badPage = `<!DOCTYPE html><html><head><title>axe positive control</title></head>` +
		`<body><img src="x.png"><a href="#"></a></body></html>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, badPage)
	}))
	defer srv.Close()

	ctx, cancel := newTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(srv.URL), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[axe]: load control page: %v", err)
	}

	results, err := RunAxe(ctx)
	if err != nil {
		t.Fatalf("visualtest[axe]: control audit: %v", err)
	}

	blocking := results.BlockingViolations()
	if len(blocking) == 0 {
		t.Fatal("visualtest[axe]: positive control FAILED — axe reported no critical/serious violations on a page with an unlabeled image and an empty link; the sweep harness is broken")
	}

	found := map[string]bool{}
	for _, violation := range blocking {
		found[violation.ID] = true
	}

	if !found["image-alt"] {
		t.Errorf("visualtest[axe]: expected image-alt among blocking rules, got %v", found)
	}
}

// axeAuditRoute loads one demo page (optionally forced to dark mode) and runs
// the axe audit against it.
func axeAuditRoute(t *testing.T, ctx context.Context, baseURL, path string, dark bool) AxeResults {
	t.Helper()

	actions := []chromedp.Action{
		chromedp.Navigate(baseURL + path),
		chromedp.WaitReady("body"),
	}

	if dark {
		actions = append(actions,
			chromedp.Evaluate(`document.documentElement.classList.add('dark'); true`, nil),
			chromedp.Sleep(settleDelay),
		)
	}

	if err := chromedp.Run(ctx, actions...); err != nil {
		t.Fatalf("visualtest[axe]: load %s%s: %v", baseURL, path, err)
	}

	results, err := RunAxe(ctx)
	if err != nil {
		t.Fatalf("visualtest[axe]: audit %s%s: %v", baseURL, path, err)
	}

	return results
}

// assertNoUnacceptedViolations fails the test for every blocking violation
// that the baseline does not explicitly accept (same rule + impact, node
// count within the accepted budget).
func assertNoUnacceptedViolations(t *testing.T, route string, baseline map[string]map[string]int, results AxeResults) {
	t.Helper()

	blocking := results.BlockingViolations()
	if len(blocking) == 0 {
		return
	}

	accepted := baseline[route]

	failures := make([]string, 0, len(blocking))

	for _, violation := range blocking {
		key := AxeViolationKey(violation)

		if budget, ok := accepted[key]; ok && len(violation.Nodes) <= budget {
			t.Logf("accepted violation %s on %s (%d nodes): %s", key, route, len(violation.Nodes), violation.Help)

			continue
		}

		failures = append(failures, describeAxeViolation(route, violation))
	}

	if len(failures) == 0 {
		return
	}

	for _, failure := range failures {
		t.Error(failure)
	}

	t.Errorf("axe sweep found %d unaccepted critical/serious violation(s) on %s. "+
		"Fix the markup, or extend testdata/axe_baseline.json with a justification if the finding is a false positive.",
		len(failures), route)
}

// describeAxeViolation renders one violation as a self-contained failure line:
// rule, impact, and every offending selector with axe's failure explanation.
func describeAxeViolation(route string, violation AxeViolation) string {
	var report strings.Builder

	fmt.Fprintf(&report, "axe[%s/%s]: %s (%s)\n  docs: %s",
		route, violation.ID, violation.Help, violation.Impact, violation.HelpURL)

	for i, node := range violation.Nodes {
		if i >= axeMaxReportedNodes {
			fmt.Fprintf(&report, "\n  ... and %d more nodes", len(violation.Nodes)-axeMaxReportedNodes)

			break
		}

		fmt.Fprintf(&report, "\n  target %v\n    %s", node.Target, node.FailureSummary)
	}

	return report.String()
}

// readAxeBaseline loads the accepted-violations ledger. A missing file is the
// strictest possible baseline (nothing accepted); a malformed one is a setup
// defect and fails loud.
func readAxeBaseline(t *testing.T) map[string]map[string]int {
	t.Helper()

	raw, err := os.ReadFile(filepath.FromSlash(axeBaselinePath))
	if os.IsNotExist(err) {
		return map[string]map[string]int{}
	}

	if err != nil {
		t.Fatalf("visualtest[axe]: read baseline: %v", err)
	}

	var baseline map[string]map[string]int
	if err := json.Unmarshal(raw, &baseline); err != nil {
		t.Fatalf("visualtest[axe]: decode baseline %s: %v", axeBaselinePath, err)
	}

	return baseline
}
