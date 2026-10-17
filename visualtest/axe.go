package visualtest

import (
	"context"
	_ "embed"
	"errors"
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

// errAxeRuntime wraps in-page axe failures; the payload carries axe's message.
var errAxeRuntime = errors.New("axe runtime error")

// AxeImpact is the axe-core severity of a violation. Only Critical and Serious
// block the sweep (see AxeImpact.Blocking); Moderate and Minor are recorded
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

const axeRunTimeout = 45 * time.Second

// RunAxe injects the vendored axe-core runtime into the current page and runs
// a violations-only audit against the whole document. The context must be a
// chromedp tab context with the page already loaded.
//
// axe.run's promise is awaited natively via chromedp's EvalAwaitPromise: the
// resolved value decodes by-value straight into AxeResults, and a rejected
// promise surfaces as a *chromedp.ExceptionError (wrapped as errAxeRuntime).
// No window-parked payload, no poll-sleep handshake.
func RunAxe(ctx context.Context) (AxeResults, error) {
	runCtx, cancel := context.WithTimeout(ctx, axeRunTimeout)
	defer cancel()

	// chromedp.Run, never action.Do: Do on a tab context hits the documented
	// "invalid context" flakiness (see visualtest e2e lessons).
	if _, err := chromedp.Run(runCtx, chromedp.Evaluate[chromedp.Void](axeSource)); err != nil {
		return AxeResults{}, fmt.Errorf("inject axe runtime: %w", err)
	}

	results, err := chromedp.Run(runCtx, chromedp.Evaluate[AxeResults](
		`axe.run(document, {resultTypes: ['violations']})`,
		chromedp.EvalAwaitPromise,
	))
	if err != nil {
		var jsErr *chromedp.ExceptionError
		if errors.As(err, &jsErr) {
			return AxeResults{}, fmt.Errorf("%w: %w", errAxeRuntime, err)
		}
		return AxeResults{}, fmt.Errorf("run axe audit: %w", err)
	}

	return results, nil
}

// AxeViolationKey identifies a violation class within one route: rule ID +
// impact. The baseline guard compares keys and node counts per key.
func AxeViolationKey(v AxeViolation) string { return v.ID + "|" + string(v.Impact) }
