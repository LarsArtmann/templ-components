package visualtest

import (
	"context"
	"io"
	"testing"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
)

// TestSparklineResponsiveStretches pins the #60 responsive contract in a real
// browser: a Responsive sparkline fills its container's width (the fixed
// 120px viewBox coordinate system scales up) while the fixed-mode control
// stays at its attribute width. Markup-level goldens can assert the viewBox
// attribute but never that the svg actually stretches.
func TestSparklineResponsiveStretches(t *testing.T) {
	page := func() templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				if _, err := io.WriteString(w, `<div style="width:600px" id="wide-parent">`); err != nil {
					return err
				}

				if err := display.Sparkline(display.SparklineProps{
					Values:     []float64{1, 3, 2, 5, 4, 6, 3, 7},
					Responsive: true,
				}).Render(ctx, w); err != nil {
					return err
				}

				if err := display.Sparkline(display.SparklineProps{
					Values: []float64{1, 3, 2, 5, 4, 6, 3, 7},
				}).Render(ctx, w); err != nil {
					return err
				}

				_, err := io.WriteString(w, `</div>`)

				return err
			})

			props := layout.DefaultPageProps()
			props.Title = "Sparkline responsive contract"

			return layout.Base(props).Render(templ.WithChildren(ctx, children), w)
		})
	}

	server := e2ePageServer(t, page, nil)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[display]: navigate: %v", err)
	}

	var m struct {
		Responsive float64 `json:"responsive"`
		Fixed      float64 `json:"fixed"`
	}

	if err := chromedp.Do(ctx, evalExprInto(
		`(() => {
			const svgs = document.querySelectorAll('#wide-parent svg');
			return {
				responsive: svgs[0] ? svgs[0].getBoundingClientRect().width : 0,
				fixed: svgs[1] ? svgs[1].getBoundingClientRect().width : 0,
			};
		})()`, &m)); err != nil {
		t.Fatalf("visualtest[display]: measure svg widths: %v", err)
	}

	if m.Fixed < 115 || m.Fixed > 125 {
		t.Errorf("visualtest[display]: fixed-mode svg drifted from its 120px attr: %v", m.Fixed)
	}

	if m.Responsive < 550 {
		t.Fatalf(
			"visualtest[display]: Responsive sparkline did not stretch into its 600px parent (width=%v) — the #60 responsive contract is broken",
			m.Responsive,
		)
	}
}

// TestSparklineZeroLinePlacement pins the #60 sign-visibility contract in a
// real browser: on an all-negative series the zero line sits ABOVE every
// plotted point (negatives hang below it); on an all-positive series it sits
// BELOW every point. Min/max normalization alone renders both series with
// identical shapes and no zero reference — exactly the confusion #60 files.
func TestSparklineZeroLinePlacement(t *testing.T) {
	page := func() templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				if err := display.Sparkline(display.SparklineProps{
					Values:       []float64{-5, -3, -8, -4, -6},
					ShowZeroLine: true,
					AriaLabel:    "negative trend",
				}).Render(ctx, w); err != nil {
					return err
				}

				return display.Sparkline(display.SparklineProps{
					Values:       []float64{5, 3, 8, 4, 6},
					ShowZeroLine: true,
					AriaLabel:    "positive trend",
				}).Render(ctx, w)
			})

			props := layout.DefaultPageProps()
			props.Title = "Sparkline zero line placement"

			return layout.Base(props).Render(templ.WithChildren(ctx, children), w)
		})
	}

	server := e2ePageServer(t, page, nil)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(server.URL),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[display]: navigate: %v", err)
	}

	var m struct {
		LineCount       int64   `json:"lineCount"`
		NegLineTop      float64 `json:"negLineTop"`
		NegLineBottom   float64 `json:"negLineBottom"`
		NegPointsTop    float64 `json:"negPointsTop"`
		PosLineTop      float64 `json:"posLineTop"`
		PosPointsBottom float64 `json:"posPointsBottom"`
	}

	if err := chromedp.Do(ctx, evalExprInto(
		`(() => {
			const svgs = document.querySelectorAll('svg[aria-label]');
			const lines = document.querySelectorAll('svg[aria-label] line');
			const negLine = lines[0] ? lines[0].getBoundingClientRect() : null;
			const posLine = lines[1] ? lines[1].getBoundingClientRect() : null;
			const negPts = svgs[0] ? svgs[0].querySelector('polyline').getBoundingClientRect() : null;
			const posPts = svgs[1] ? svgs[1].querySelector('polyline').getBoundingClientRect() : null;
			return {
				lineCount: lines.length,
				negLineTop: negLine ? negLine.top : -1,
				negLineBottom: negLine ? negLine.bottom : -1,
				negPointsTop: negPts ? negPts.top : -1,
				posLineTop: posLine ? posLine.top : -1,
				posPointsBottom: posPts ? posPts.bottom : -1,
			};
		})()`, &m)); err != nil {
		t.Fatalf("visualtest[display]: measure zero line geometry: %v", err)
	}

	if m.LineCount != 2 {
		t.Fatalf("visualtest[display]: expected 2 zero lines, found %d — markup moved", m.LineCount)
	}

	// All-negative: every plotted point renders BELOW the zero line.
	if m.NegPointsTop <= m.NegLineBottom {
		t.Fatalf(
			"visualtest[display]: negative series overlaps the zero line (points top=%v, line bottom=%v) — sign is not visible",
			m.NegPointsTop,
			m.NegLineBottom,
		)
	}

	// All-positive: every plotted point renders ABOVE the zero line.
	if m.PosPointsBottom >= m.PosLineTop {
		t.Fatalf(
			"visualtest[display]: positive series overlaps the zero line (points bottom=%v, line top=%v) — sign is not visible",
			m.PosPointsBottom,
			m.PosLineTop,
		)
	}
}
