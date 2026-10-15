package visualtest

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

// M11 CDP cheap wins (plan μ41–μ45): vision-deficiency passes, reduced-motion
// emulation, and a print-media probe — all on top of the existing chromedp
// plumbing (newFlowTab + StartDemoServer), extending the primary pipeline
// rather than duplicating it. Emulation applies via raw CDP calls because
// chromedp has no typed action for these domains.

func emulateVisionDeficiency(t *testing.T, ctx context.Context, deficiency emulation.SetEmulatedVisionDeficiencyType) {
	t.Helper()

	if err := chromedp.Do(ctx, chromedp.Func(func(ctx context.Context, tgt *chromedp.Target) error {
		_, err := cdp.Call(
			ctx,
			tgt,
			emulation.SetEmulatedVisionDeficiency,
			emulation.SetEmulatedVisionDeficiencyParams{Type: deficiency},
		)

		return err
	})); err != nil {
		t.Fatalf("visualtest[emulation]: set vision deficiency %q: %v", deficiency, err)
	}
}

func emulateMedia(t *testing.T, ctx context.Context, media string, features ...*emulation.MediaFeature) {
	t.Helper()

	if err := chromedp.Do(ctx, chromedp.Func(func(ctx context.Context, tgt *chromedp.Target) error {
		_, err := cdp.Call(
			ctx,
			tgt,
			emulation.SetEmulatedMedia,
			emulation.SetEmulatedMediaParams{Media: media, Features: features},
		)

		return err
	})); err != nil {
		t.Fatalf("visualtest[emulation]: set emulated media %q: %v", media, err)
	}
}

func fullShot(t *testing.T, ctx context.Context) []byte {
	t.Helper()

	shot, err := chromedp.Run(ctx, chromedp.FullScreenshot(100))
	if err != nil {
		t.Fatalf("visualtest[emulation]: full screenshot: %v", err)
	}

	return shot
}

func saveProbeArtifact(t *testing.T, dir, name string, png []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), dir)
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("visualtest[emulation]: mkdir: %v", err)
	}

	p := filepath.Join(path, name)
	if err := os.WriteFile(p, png, 0o644); err != nil {
		t.Fatalf("visualtest[emulation]: write artifact: %v", err)
	}

	return p
}

// TestVisionDeficiencyProbe proves the CDP vision emulation applies to the
// demo routes and captures review artifacts: under achromatopsia and
// deuteranopia the rendering MUST differ from the trichromatic baseline.
// The screenshots land in a temp dir for human review; the test gates the
// mechanism, color-safety judgment stays a human finding (μ41–μ42).
func TestVisionDeficiencyProbe(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	for _, route := range []string{"/display", "/feedback", "/users"} {
		t.Run(route, func(t *testing.T) {
			if err := chromedp.Do(ctx,
				chromedp.Navigate(server.BaseURL()+route),
				chromedp.WaitReady("body"),
			); err != nil {
				t.Fatalf("navigate %s: %v", route, err)
			}

			baseline := fullShot(t, ctx)

			for _, deficiency := range []emulation.SetEmulatedVisionDeficiencyType{
				emulation.SetEmulatedVisionDeficiencyTypeAchromatopsia,
				emulation.SetEmulatedVisionDeficiencyTypeDeuteranopia,
			} {
				emulateVisionDeficiency(t, ctx, deficiency)
				shot := fullShot(t, ctx)
				emulateVisionDeficiency(t, ctx, emulation.SetEmulatedVisionDeficiencyTypeNone)

				if bytes.Equal(shot, baseline) {
					t.Errorf(
						"visualtest[emulation]: %s under %s is byte-identical to baseline — vision emulation did not apply",
						route,
						deficiency,
					)
				}

				path := saveProbeArtifact(t, "vision", string(deficiency)+routeFilename(route)+".png", shot)
				t.Logf("captured %s", path)
			}
		})
	}

	server.FailIfServerErrors(t)
}

// TestReducedMotionEmulationKillsSpinnerAnimation proves the demo's animated
// components honor prefers-reduced-motion through the CDP media emulation:
// the spinner's CSS animation must collapse to `none` under the reduce
// setting (the motion-reduce:animate-none convention), and be active
// without it.
func TestReducedMotionEmulationKillsSpinnerAnimation(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(server.BaseURL()+"/feedback"),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("navigate /feedback: %v", err)
	}

	const animationNameExpr = `(() => { const el = document.querySelector('.animate-spin'); return el ? getComputedStyle(el).animationName : 'NO-SPINNER'; })()`

	var active string
	if err := chromedp.Do(ctx, evalExprInto(animationNameExpr, &active)); err != nil {
		t.Fatalf("read spinner animation: %v", err)
	}

	if active == "none" || active == "NO-SPINNER" {
		t.Fatalf("spinner animation without emulation = %q, want an active animation name", active)
	}

	emulateMedia(t, ctx, "", &emulation.MediaFeature{Name: "prefers-reduced-motion", Value: "reduce"})
	defer emulateMedia(t, ctx, "")

	var reduced string
	if err := chromedp.Do(ctx, evalExprInto(animationNameExpr, &reduced)); err != nil {
		t.Fatalf("read spinner animation under reduce: %v", err)
	}

	if reduced != "none" {
		t.Errorf(
			"spinner animation under prefers-reduced-motion: reduce = %q, want \"none\" (motion-reduce:animate-none missing?)",
			reduced,
		)
	}
}

// TestPrintMediaProbe renders Table (/users) and Card-heavy (/recipes) pages
// under print media and asserts they still produce non-degenerate output.
// Print styling is NOT gated here (μ45: findings go to the evaluation doc /
// TODO seeds) — the probe exists to keep the print path observable.
func TestPrintMediaProbe(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	emulateMedia(t, ctx, "print")
	defer emulateMedia(t, ctx, "")

	for _, route := range []string{"/users", "/recipes"} {
		t.Run(route, func(t *testing.T) {
			if err := chromedp.Do(ctx,
				chromedp.Navigate(server.BaseURL()+route),
				chromedp.WaitReady("body"),
			); err != nil {
				t.Fatalf("navigate %s: %v", route, err)
			}

			var textLen int64
			if err := chromedp.Do(ctx, evalExprInto(`document.body.innerText.length`, &textLen)); err != nil {
				t.Fatalf("measure print text: %v", err)
			}

			if textLen < 200 {
				t.Errorf(
					"visualtest[emulation]: %s under print media renders only %d chars of text — layout may be collapsing",
					route,
					textLen,
				)
			}

			shot := fullShot(t, ctx)
			path := saveProbeArtifact(t, "print", "print"+routeFilename(route)+".png", shot)
			t.Logf("print capture for %s: %d bytes → %s", route, len(shot), path)
		})
	}

	server.FailIfServerErrors(t)
}

func routeFilename(route string) string {
	if route == "/" {
		return "-home"
	}

	return "-" + string(bytes.ReplaceAll([]byte(route[1:]), []byte("/"), []byte("-")))
}
