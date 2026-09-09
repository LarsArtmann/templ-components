package visualtest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
// of the baseline file) — an accepted violation is a documented a11y debt,
// not a pass.

// axeBaselinePath points at the accepted-violations ledger.
const axeBaselinePath = "testdata/axe_baseline.json"

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

	newTabForTest := newTab // shared-Chromium allocator
	_ = newTabForTest

	for _, route := range axeSweepRoutes {
		t.Run(route.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := newTab(t)
			defer cancel()

			var loadErr error
			actions := []chromedp.Action{
				chromedp.Navigate(server.BaseURL() + route.path),
				chromedp.WaitReady("body"),
			}
			if route.dark {
				actions = append(actions,
					chromedp.Evaluate(`document.documentElement.classList.add('dark'); true`, nil),
					chromedp.Sleep(haltDelay), // let the dark class re-paint before auditing
				)
			}

			if loadErr = chromedp.Run(ctx, actions...); loadErr != nil {
				t.Fatalf("load %s%s: %v", server.BaseURL(), route.path, loadErr)
			}

			results, err := RunAxe(ctx)
			if err != nil {
				t.Fatalf("axe run on %s: %v", route.path, err)
			}

			blocking := results.BlockingViolations()
			if len(blocking) == 0 {
				return
			}

			accepted, ok := baseline[route.name]
			var failures []string

			for _, violation := range blocking {
				key := AxeViolationKey(violation)

				if ok {
					if allowed, seen := accepted[key]; seen && len(violation.Nodes) <= allowed {
						t.Logf("accepted violation %s on %s (%d nodes): %s", key, route.name, len(violation.Nodes), violation.Help)

						continue
					}
				}

				failures = append(failures, describeAxeViolation(route.name, violation))
			}

			if len(failures) > 0 {
				for _, failure := range failures {
					t.Error(failure)
				}

				t.Errorf("axe sweep found %d unaccepted critical/serious violation(s) on %s. "+
					"Fix the markup, or extend testdata/axe_baseline.json with a justification if the finding is a false positive.",
					len(failures), route.name)
			}
		})
	}
}

// describeAxeViolation renders one violation as a self-contained failure line:
// rule, impact, and every offending selector with axe's failure explanation.
func describeAxeViolation(route string, violation AxeViolation) string {
	summary := fmt.Sprintf("axe[%s/%s]: %s (%s)\n  %s\n  docs: %s",
		route, violation.ID, violation.Help, violation.Impact, violation.HelpURL)

	for i, node := range violation.Nodes {
		if i >= axeMaxReportedNodes {
			summary += fmt.Sprintf("\n  ... and %d more nodes", len(violation.Nodes)-axeMaxReportedNodes)

			break
		}

		summary += fmt.Sprintf("\n  target %v\n    %s", node.Target, node.FailureSummary)
	}

	return summary
}

// axeMaxReportedNodes caps the per-violation node report so one widespread
// rule failure cannot flood the test output.
const axeMaxReportedNodes = 5

// readAxeBaseline loads the accepted-violations ledger. A missing file is the
// strictest possible baseline (nothing accepted); a malformed one is a setup
// bug and fails loud.
func readAxeBaseline(t *testing.T) map[string]map[string]int {
	t.Helper()

	raw, err := os.ReadFile(filepath.FromSlash(axeBaselinePath)) //nolint:gosec // path is a package-level constant
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
