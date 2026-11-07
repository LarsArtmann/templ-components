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

// TestBarChartVerticalHeightContract pins the #48 vertical layout contract
// in a real browser: with ONLY Height set — the documented usage — every
// bar measures > 0px and scales with its value. The historical markup
// (items-end columns against an indefinite height chain) resolved every
// percentage height to 0px; goldens are string-level and never saw it, so
// this guard measures actual layout.
func TestBarChartVerticalHeightContract(t *testing.T) {
	page := func() templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				return display.BarChart(display.BarChartProps{
					Orient: display.BarVertical,
					Height: "8rem",
					Bars: []display.BarChartBar{
						{Label: "low", Value: 200},
						{Label: "mid", Value: 450},
						{Label: "high", Value: 1200},
					},
				}).Render(ctx, w)
			})

			props := layout.DefaultPageProps()
			props.Title = "Vertical bar height contract"

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
		Count int64   `json:"count"`
		Low   float64 `json:"low"`
		Mid   float64 `json:"mid"`
		High  float64 `json:"high"`
	}

	if err := chromedp.Do(ctx, evalExprInto(
		`(() => {
			const chart = document.querySelector('[role="img"]');
			const bars = chart ? chart.querySelectorAll('.rounded-t') : [];
			const h = el => (el ? el.getBoundingClientRect().height : 0);
			return { count: bars.length, low: h(bars[0]), mid: h(bars[1]), high: h(bars[2]) };
		})()`, &m)); err != nil {
		t.Fatalf("visualtest[display]: measure bar heights: %v", err)
	}

	if m.Count != 3 {
		t.Fatalf("visualtest[display]: expected 3 vertical bars, found %d — markup moved", m.Count)
	}

	if m.Low <= 0 || m.Mid <= 0 || m.High <= 0 {
		t.Fatalf(
			"visualtest[display]: vertical bars collapsed to 0px with only Height set (low=%v mid=%v high=%v) — the #48 regression is back",
			m.Low,
			m.Mid,
			m.High,
		)
	}

	// The bars must scale with value, not just render: 200 < 450 < 1200 and
	// the tallest is a multiple of the shortest (6x nominal; bound loosely
	// against rounding and the plot's baseline padding).
	if !(m.Low < m.Mid && m.Mid < m.High) {
		t.Errorf("visualtest[display]: bar heights not ordered by value: low=%v mid=%v high=%v", m.Low, m.Mid, m.High)
	}

	if m.High < 4*m.Low {
		t.Errorf(
			"visualtest[display]: height ratio high/low = %.2f, want >= 4 (values 1200 vs 200) — bars are not scaling with data",
			m.High/m.Low,
		)
	}
}

// TestBarChartVerticalSignedLayout pins the #59 signed rendering in a real
// browser: the zero line sits between the positive and negative bars, and
// each bar stays on its own side of it.
func TestBarChartVerticalSignedLayout(t *testing.T) {
	page := func() templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				return display.BarChart(display.BarChartProps{
					Orient: display.BarVertical,
					Height: "12rem",
					Bars: []display.BarChartBar{
						{Label: "in", Value: 420},
						{Label: "out", Value: -180},
					},
				}).Render(ctx, w)
			})

			props := layout.DefaultPageProps()
			props.Title = "Signed bar layout"

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
		PositiveTop float64 `json:"positiveTop"`
		PositiveBot float64 `json:"positiveBot"`
		NegativeTop float64 `json:"negativeTop"`
		NegativeBot float64 `json:"negativeBot"`
		LineTop     float64 `json:"lineTop"`
		LineBot     float64 `json:"lineBot"`
		Count       int64   `json:"count"`
	}

	if err := chromedp.Do(ctx, evalExprInto(
		`(() => {
			const chart = document.querySelector('[role="img"]');
			const pos = chart ? chart.querySelector('.rounded-t') : null;
			const neg = chart ? chart.querySelector('.rounded-b') : null;
			const line = chart ? chart.querySelector('.h-px') : null;
			const r = el => (el ? el.getBoundingClientRect() : null);
			const pr = r(pos), nr = r(neg), lr = r(line);
			return {
				positiveTop: pr ? pr.top : -1,
				positiveBot: pr ? pr.bottom : -1,
				negativeTop: nr ? nr.top : -1,
				negativeBot: nr ? nr.bottom : -1,
				lineTop: lr ? lr.top : -1,
				lineBot: lr ? lr.bottom : -1,
				count: (pos?1:0) + (neg?1:0) + (line?1:0),
			};
		})()`, &m)); err != nil {
		t.Fatalf("visualtest[display]: measure signed layout: %v", err)
	}

	if m.Count != 3 {
		t.Fatalf("visualtest[display]: expected positive bar + negative bar + zero line, found %d", m.Count)
	}

	// Vertical order by element CENTER: positive bar, then the zero line,
	// then the hanging negative bar — the sign is readable from geometry
	// alone. Centers (not edges) absorb the line's own 1px thickness.
	positiveCenter := (m.PositiveTop + m.PositiveBot) / 2
	lineCenter := (m.LineTop + m.LineBot) / 2
	negativeCenter := (m.NegativeTop + m.NegativeBot) / 2

	if !(positiveCenter < lineCenter && lineCenter < negativeCenter) {
		t.Errorf(
			"visualtest[display]: signed layout broken: positiveCenter=%v lineCenter=%v negativeCenter=%v (want positive < line < negative)",
			positiveCenter,
			lineCenter,
			negativeCenter,
		)
	}

	// The positive bar must grow UP from the line, not straddle it (1px
	// tolerance for the hairline's own extent).
	if m.PositiveBot > m.LineTop+1.5 {
		t.Errorf(
			"visualtest[display]: positive bar bottom (%v) crosses the zero line (%v) by more than a pixel",
			m.PositiveBot, m.LineTop,
		)
	}
}
