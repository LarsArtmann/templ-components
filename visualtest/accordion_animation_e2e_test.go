package visualtest

import (
	"context"
	"testing"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

// The Accordion animates its panel through the ::details-content pseudo-element
// (Modern Web Guidance: animate-to-intrinsic-sizes), gated on
// @supports selector(::details-content) and prefers-reduced-motion:
// no-preference. These tests prove BOTH halves in a real browser against the
// compiled demo CSS — an intermediate height while opening, and an instant
// snap when the user asks for reduced motion. String assertions cannot see
// either: a missing transition and a permanently-collapsed panel look the same
// in HTML.

// accordionSample is the decoded result of one animation probe: the panel's
// height before opening, then one height per animation frame for 500ms.
type accordionSample struct {
	Closed  float64   `json:"closed"`
	Samples []float64 `json:"samples"`
}

// accordionHeightSamples opens the <details> element with the given id and
// records its border-box height on every animation frame for 500ms. Sampling
// the <details> height (not the pseudo-element's computed block-size) keeps the
// probe on stable CSSOM ground and measures the user-visible panel growth
// directly. The pre-open height is captured first: the first sample after
// `open = true` is already moving, and under reduced motion it is the settled
// height, so "did it grow" must be judged against the closed state.
func accordionHeightSamples(id string) string {
	return `(async () => {
  const d = document.getElementById('` + id + `');
  if (!d) return null;
  const closed = d.getBoundingClientRect().height;
  d.open = true;
  const t0 = performance.now();
  const hs = [];
  while (performance.now() - t0 < 500) {
    hs.push(Math.round(d.getBoundingClientRect().height * 100) / 100);
    await new Promise(r => requestAnimationFrame(r));
  }
  return {closed: closed, samples: hs};
})()`
}

// openDisplayPage navigates to the demo's display route, which carries the
// real Accordion markup and the compiled app.css that animates it.
func openDisplayPage(t *testing.T, ctx context.Context, baseURL string) {
	t.Helper()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(baseURL+"/display"),
		chromedp.WaitReady("body"),
		pollTrue(`!!document.getElementById('acc-2')`, pollFast()...),
	); err != nil {
		t.Fatalf("open /display: %v", err)
	}
}

// sampleAccordionHeights runs the frame loop for the given accordion item.
func sampleAccordionHeights(t *testing.T, ctx context.Context, id, label string) accordionSample {
	t.Helper()

	var probe accordionSample
	// EvalAwaitPromise is REQUIRED: the expression is an async IIFE, and
	// without the option the protocol returns the Promise object itself,
	// which cannot unmarshal into the result struct.
	if err := chromedp.Do(ctx, evalInto(
		chromedp.Evaluate[accordionSample](accordionHeightSamples(id), chromedp.EvalAwaitPromise),
		&probe,
	)); err != nil {
		t.Fatalf("%s: sample accordion heights: %v", label, err)
	}

	if len(probe.Samples) < 5 {
		t.Fatalf("%s: only %d height samples collected, want >= 5 (animation frame loop did not run)",
			label, len(probe.Samples))
	}

	settled := probe.Samples[len(probe.Samples)-1]
	if settled <= probe.Closed+1 {
		t.Fatalf("%s: panel never grew: closed %.2f, settled %.2f", label, probe.Closed, settled)
	}

	return probe
}

// TestAccordionPanelAnimatesOpen asserts a mid-transition height exists while
// the panel opens: some sample is strictly between the closed height and the
// settled height. Without the transition every sample equals the final height
// and the test fails.
func TestAccordionPanelAnimatesOpen(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	emulateMedia(t, ctx, "", &emulation.MediaFeature{Name: "prefers-reduced-motion", Value: "no-preference"})

	openDisplayPage(t, ctx, server.BaseURL())

	// acc-2 starts closed in the demo data, so the loop covers the whole
	// opening transition.
	probe := sampleAccordionHeights(t, ctx, "acc-2", "animated open")

	closed := probe.Closed
	settled := probe.Samples[len(probe.Samples)-1]

	intermediate := 0

	for _, h := range probe.Samples {
		if h > closed+2 && h < settled-2 {
			intermediate++
		}
	}

	if intermediate == 0 {
		t.Fatalf(
			"no intermediate panel height in %d samples (%.2f -> %.2f): the ::details-content transition is missing",
			len(probe.Samples), closed, settled,
		)
	}
}

// TestAccordionSnapsUnderReducedMotion asserts the panel jumps straight to its
// open height when prefers-reduced-motion: reduce is set — the media query
// removes the animation entirely rather than shortening it.
func TestAccordionSnapsUnderReducedMotion(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	emulateMedia(t, ctx, "", &emulation.MediaFeature{Name: "prefers-reduced-motion", Value: "reduce"})
	defer emulateMedia(t, ctx, "")

	openDisplayPage(t, ctx, server.BaseURL())

	// acc-3 also starts closed; under reduce its first frame must already be
	// the settled height.
	probe := sampleAccordionHeights(t, ctx, "acc-3", "reduced motion")

	settled := probe.Samples[len(probe.Samples)-1]

	for i, h := range probe.Samples {
		if h < settled-1 {
			t.Fatalf(
				"sample %d = %.2f is below the settled height %.2f under prefers-reduced-motion: reduce — the panel still animates",
				i,
				h,
				settled,
			)
		}
	}
}
