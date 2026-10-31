package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
)

// tooltipAnchorE2EPage renders two tooltips: one near the RIGHT viewport
// edge (PositionTop — the base CSS position would overflow right) and one
// in normal flow (control). The wrapper div reserves space and pushes the
// first trigger against the right edge.
func tooltipAnchorE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Right-edge trigger: margin-left:auto shoves it to the container's
		// end; the container spans the viewport (Base body has no max width).
		edge := display.TooltipProps{
			ID:       "tt-edge",
			Text:     "Edge tooltip with quite a lot of text that extends far",
			Position: display.PositionTop,
		}
		edgeWrap := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, `<div class="flex justify-end px-1 pt-4">`); err != nil {
				return err
			}
			if err := display.Tooltip(edge).Render(templ.WithChildren(ctx, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<button id="tt-edge-btn" type="button" class="rounded border border-gray-400 px-2 py-1">Edge trigger</button>`)

				return err
			})), w); err != nil {
				return err
			}
			_, err := io.WriteString(w, `</div>`)

			return err
		})
		if err := edgeWrap.Render(ctx, w); err != nil {
			return err
		}

		center := display.TooltipProps{
			ID:       "tt-center",
			Text:     "Centered control tooltip",
			Position: display.PositionTop,
		}
		centerWrap := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, `<div class="flex justify-center pt-64">`); err != nil {
				return err
			}
			if err := display.Tooltip(center).Render(templ.WithChildren(ctx, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<button id="tt-center-btn" type="button" class="rounded border border-gray-400 px-2 py-1">Center trigger</button>`)

				return err
			})), w); err != nil {
				return err
			}
			_, err := io.WriteString(w, `</div>`)

			return err
		})

		return centerWrap.Render(ctx, w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		props := layout.DefaultPageProps()
		props.Title = "Tooltip Anchor E2E — templ-components"
		props.CSSPath = "/app.css"

		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

func tooltipAnchorE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(css)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = tooltipAnchorE2EPage().Render(r.Context(), w)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// hoverTooltip moves the REAL mouse (CDP) onto the trigger so the CSS
// :hover path engages (synthetic JS mouseover events do not flip :hover),
// then waits for the bubble to become visible.
func hoverTooltip(t *testing.T, ctx context.Context, sel string) {
	t.Helper()

	var box overlayBox
	if err := chromedp.Do(ctx, evalExprInto(
		`(() => { const r = document.querySelector(`+jsString(sel)+`).getBoundingClientRect();`+
			`return {left: r.left, top: r.top, right: r.right, bottom: r.bottom}; })()`,
		&box,
	)); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: read trigger box %s: %v", sel, err)
	}

	if err := chromedp.Do(ctx, mouseEvent(input.DispatchMouseEventTypeMouseMoved, box.insideX(), box.insideY())); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: hover %s: %v", sel, err)
	}

	visibleExpr := `(function(){var w=document.querySelector(` + jsString(sel) + `).parentElement;` +
		`var b=w.querySelector('[role=tooltip]');return b&&getComputedStyle(b).display!=='none'})()`
	if err := chromedp.Do(ctx, pollTrue(visibleExpr, pollFast()...)); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: bubble visible for %s: %v", sel, err)
	}
}

// TestTooltipAnchorPositioning proves the CSS Anchor Positioning layer
// (templates/custom.css @supports block) keeps tooltips inside the
// viewport: a tooltip on a trigger pinned against the RIGHT viewport edge
// must not overflow right, and a control tooltip in the middle must stay
// centered over its trigger. On engines without anchor positioning the
// static CSS classes ship the legacy behavior (base-centered, possibly
// clipped) — the flip/clamp assertions then only run when the layer is
// active. Serial — shared browser allocator.
func TestTooltipAnchorPositioning(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := tooltipAnchorE2EServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: navigate: %v", err)
	}

	hoverTooltip(t, ctx, "#tt-edge-btn")

	var edgeBox struct {
		Left, Right, Top, Bottom, VW float64
	}

	readEdge := `(function(){var b=document.querySelector('#tt-edge-tooltip');if(!b)return null;var r=b.getBoundingClientRect();return {left:r.left,right:r.right,top:r.top,bottom:r.bottom,vw:window.innerWidth};})()`
	if err := chromedp.Do(ctx, evalExprInto(readEdge, &edgeBox)); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: read edge bubble box: %v", err)
	}

	if edgeBox.VW == 0 {
		t.Fatal("visualtest[tooltip-anchor]: edge bubble not found")
	}

	var anchorSupported bool

	if err := chromedp.Do(ctx, evalExprInto(
		`window.CSS&&CSS.supports&&CSS.supports("anchor-name: --x")`,
		&anchorSupported,
	)); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: probe anchor support: %v", err)
	}

	if anchorSupported {
		if edgeBox.Right > edgeBox.VW+0.5 {
			t.Errorf("anchor layer: edge tooltip overflows right: right=%.1f vw=%.1f", edgeBox.Right, edgeBox.VW)
		}

		if edgeBox.Left < -0.5 {
			t.Errorf("anchor layer: edge tooltip overflows left: left=%.1f", edgeBox.Left)
		}
	}

	hoverTooltip(t, ctx, "#tt-center-btn")

	var centerBox struct {
		Left, Right, Top, Bottom, VW float64
	}

	readCenter := `(function(){var b=document.querySelector('#tt-center-tooltip');if(!b)return null;var r=b.getBoundingClientRect();return {left:r.left,right:r.right,top:r.top,bottom:r.bottom,vw:window.innerWidth};})()`
	if err := chromedp.Do(ctx, evalExprInto(readCenter, &centerBox)); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: read center bubble box: %v", err)
	}

	if centerBox.Left < -0.5 || centerBox.Right > centerBox.VW+0.5 {
		t.Errorf("center tooltip overflows viewport: left=%.1f right=%.1f vw=%.1f", centerBox.Left, centerBox.Right, centerBox.VW)
	}
}
