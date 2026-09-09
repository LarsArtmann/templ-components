package visualtest

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

// axeSource is the vendored axe-core runtime (Deque Systems, MPL-2.0). Vendored
// — not CDN-fetched — so the a11y sweep runs offline and hermetically: the
// library's accessibility guarantee must not depend on a third party's uptime.
// The embedded file keeps its full copyright header (an MPL-2.0 requirement).
//
//go:embed testdata/axe.min.js
var axeSource string

// AxeImpact is the axe-core severity of a violation. Only Critical and Serious
// block the sweep (see AxeViolation.Blocking); Moderate and Minor are recorded
// for trending but never gate CI.
type AxeImpact string

const (
	AxeImpactCritical AxeImpact = "critical"
	AxeImpactSerious  AxeImpact = "serious"
	AxeImpactModerate AxeImpact = "moderate"
	AxeImpactMinor    AxeImpact = "minor"
)

// Blocking reports whether the impact gates the sweep: the zero-critical guard
// fails on any NEW critical or serious violation.
func (i AxeImpact) Blocking() bool { return i == AxeImpactCritical || i == AxeImpactSerious }

// AxeViolation is one axe-core rule failure: a rule ID, its severity, and the
// DOM nodes that violated it.
type AxeViolation struct {
	ID      string    `json:"id"`
	Impact  AxeImpact `json:"impact"`
	Help    string    `json:"help"`
	HelpURL string    `json:"helpUrl"`
	Nodes   []AxeNode `json:"nodes"`
	Tags    []string  `json:"tags"`
}

// AxeNode is a single offending DOM node: its outer HTML, a CSS selector path,
// and axe's plain-language failure explanation.
type AxeNode struct {
	HTML           string   `json:"html"`
	Target         []string `json:"target"`
	FailureSummary string   `json:"failureSummary"`
}

// AxeResults is the subset of axe.run's result payload the sweep consumes.
type AxeResults struct {
	Violations []AxeViolation `json:"violations"`
}

// BlockingViolations returns only the critical/serious violations.
func (r AxeResults) BlockingViolations() []AxeViolation {
	blocking := make([]AxeViolation, 0, len(r.Violations))

	for _, violation := range r.Violations {
		if violation.Impact.Blocking() {
			blocking = append(blocking, violation)
		}
	}

	return blocking
}

const (
	axeRunTimeout    = 45 * time.Second
	axePollInterval  = 100 * time.Millisecond
	axeResultTimeout = 30 * time.Second
)

// RunAxe injects the vendored axe-core runtime into the current page and runs
// a violations-only audit against the whole document. The context must be a
// chromedp tab context with the page already loaded.
//
// The async axe.run promise is bridged via a polling handshake because
// chromedp v0.16 has no await-promise evaluate option: the run parks its JSON
// payload in window.__tcAxeJSON (or the error in window.__tcAxeErr), and the
// poll waits for either to appear.
func RunAxe(ctx context.Context) (AxeResults, error) {
	runCtx, cancel := context.WithTimeout(ctx, axeRunTimeout)
	defer cancel()

	if err := chromedp.Evaluate(axeSource, nil).Do(runCtx); err != nil {
		return AxeResults{}, fmt.Errorf("inject axe runtime: %w", err)
	}

	const bootScript = `window.__tcAxeJSON = null; window.__tcAxeErr = null;
axe.run(document, {resultTypes: ['violations']})
  .then(r => { window.__tcAxeJSON = JSON.stringify(r); })
  .catch(e => { window.__tcAxeErr = String(e); });
true`

	if err := chromedp.Evaluate(bootScript, nil).Do(runCtx); err != nil {
		return AxeResults{}, fmt.Errorf("start axe run: %w", err)
	}

	if err := chromedp.Poll(
		`window.__tcAxeJSON !== null || window.__tcAxeErr !== null`,
		nil,
		chromedp.WithPollingInterval(axePollInterval),
		chromedp.WithPollingTimeout(axeResultTimeout),
	).Do(runCtx); err != nil {
		return AxeResults{}, fmt.Errorf("axe run did not settle: %w", err)
	}

	var payload string
	if err := chromedp.Evaluate(
		`window.__tcAxeErr !== null ? "ERR:" + window.__tcAxeErr : window.__tcAxeJSON`,
		&payload,
	).Do(runCtx); err != nil {
		return AxeResults{}, fmt.Errorf("fetch axe result: %w", err)
	}

	if len(payload) > 4 && payload[:4] == "ERR:" {
		return AxeResults{}, fmt.Errorf("axe runtime error: %s", payload[4:])
	}

	var results AxeResults
	if err := json.Unmarshal([]byte(payload), &results); err != nil {
		return AxeResults{}, fmt.Errorf("decode axe result (%d bytes): %w", len(payload), err)
	}

	return results, nil
}

// AxeViolationKey identifies a violation class within one route: rule ID +
// impact. The baseline guard compares keys and node counts per key.
func AxeViolationKey(v AxeViolation) string { return v.ID + "|" + string(v.Impact) }
