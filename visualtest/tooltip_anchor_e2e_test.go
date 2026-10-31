package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
)

// tooltipAnchorForceStyle keeps the tooltip bubbles on-screen unconditionally.
// Headless Chrome reports (hover: none) — CDP cannot emulate the hover media
// feature — so Tailwind's @media (hover: hover) gate never opens and the
// CSS hover-show path stays dark. Real desktop engines show tooltips via
// that gate; this e2e exists to prove the ANCHOR-POSITIONING layer (C23),
// so the show/hide path is forced and geometry is asserted instead.
const tooltipAnchorForceStyle = `<style>[data-tc-tooltip] > [role="tooltip"]{display:block !important}</style>`

// tooltipAnchorE2EPage renders two tooltips: one on a trigger pinned against
// the RIGHT viewport edge (PositionTop — a base-centered bubble would
// overflow), one in normal flow (control).
func tooltipAnchorE2EPage() templ.Component {
	tooltipAt := func(id string, wrapperClass string, label string) templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, `<div class="`+wrapperClass+`">`); err != nil {
				return err
			}

			props := display.TooltipProps{ID: id, Text: label, Position: display.PositionTop}
			inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<button id="`+id+`-btn" type="button" class="rounded border border-gray-400 px-2 py-1">`+label+` trigger</button>`)

				return err
			})
			if err := display.Tooltip(props).Render(templ.WithChildren(ctx, inner), w); err != nil {
				return err
			}

			_, err := io.WriteString(w, `</div>`)

			return err
		})
	}

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := tooltipAt("tt-edge", "flex justify-end px-1 pt-4", "Edge").Render(ctx, w); err != nil {
			return err
		}

		return tooltipAt("tt-center", "flex justify-center pt-64", "Center").Render(ctx, w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		props := layout.DefaultPageProps()
		props.Title = "Tooltip Anchor E2E — templ-components"
		props.CSSPath = "/app.css"
		props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
			_, err := io.WriteString(w, tooltipAnchorForceStyle)

			return err
		})

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

// TestTooltipAnchorPositioning proves the CSS Anchor Positioning layer
// (templates/custom.css @supports block) keeps tooltips inside the viewport
// and over their triggers. On anchor-supporting engines the shared JS
// positioning fallback must stay dormant (tcPositionPopover undefined); on
// non-supporting engines the assertions degrade to the legacy geometry
// checks that do not depend on the layer. Serial — shared browser allocator.
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

	var anchorSupported bool
	if err := chromedp.Do(ctx, evalExprInto(
		`window.CSS&&CSS.supports&&CSS.supports("anchor-name: --x")`,
		&anchorSupported,
	)); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: probe anchor support: %v", err)
	}

	tooltipBox := func(id string) overlayBox {
		t.Helper()

		var box overlayBox
		if err := chromedp.Do(ctx, evalExprInto(
			`(() => { const b=document.querySelector('#`+id+`-tooltip'); if(!b) return null; const r=b.getBoundingClientRect(); return {left:r.left, top:r.top, right:r.right, bottom:r.bottom}; })()`,
			&box,
		)); err != nil {
			t.Fatalf("visualtest[tooltip-anchor]: read %s box: %v", id, err)
		}

		return box
	}

	triggerBox := func(id string) overlayBox {
		t.Helper()

		var box overlayBox
		if err := chromedp.Do(ctx, evalExprInto(
			`(() => { const r=document.querySelector('#`+id+`-btn').getBoundingClientRect(); return {left:r.left, top:r.top, right:r.right, bottom:r.bottom}; })()`,
			&box,
		)); err != nil {
			t.Fatalf("visualtest[tooltip-anchor]: read %s trigger box: %v", id, err)
		}

		return box
	}

	var vw float64
	if err := chromedp.Do(ctx, evalExprInto(`window.innerWidth`, &vw)); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: read viewport width: %v", err)
	}

	if anchorSupported {
		// The JS fallback must NOT have attached its positioning machinery
		// when the CSS layer is active — one positioning brain per engine.
		var jsFallbackAttached bool
		if err := chromedp.Do(ctx, evalExprInto(
			`typeof window.tcPositionPopover==="function"`,
			&jsFallbackAttached,
		)); err != nil {
			t.Fatalf("visualtest[tooltip-anchor]: probe JS fallback state: %v", err)
		}

		if jsFallbackAttached {
			t.Error("anchor layer active but the JS positioning fallback attached anyway")
		}
	}

	edge := tooltipBox("tt-edge")
	if edge.Right == 0 && edge.Left == 0 {
		t.Fatal("visualtest[tooltip-anchor]: edge bubble not rendered")
	}

	if anchorSupported {
		if edge.Right > vw+0.5 {
			t.Errorf("anchor layer: edge tooltip overflows right: right=%.1f vw=%.1f", edge.Right, vw)
		}

		if edge.Left < -0.5 {
			t.Errorf("anchor layer: edge tooltip overflows left: left=%.1f", edge.Left)
		}

		edgeTrigger := triggerBox("tt-edge")
		if edge.Bottom > edgeTrigger.Top+0.5 {
			t.Errorf("anchor layer: PositionTop tooltip sits below its trigger: bubbleBottom=%.1f triggerTop=%.1f", edge.Bottom, edgeTrigger.Top)
		}
	}

	center := tooltipBox("tt-center")
	centerTrigger := triggerBox("tt-center")

	if center.Left < -0.5 || center.Right > vw+0.5 {
		t.Errorf("center tooltip overflows viewport: left=%.1f right=%.1f vw=%.1f", center.Left, center.Right, vw)
	}

	centerDelta := (center.Left+center.Right)/2 - (centerTrigger.Left+centerTrigger.Right)/2
	if centerDelta < -24 || centerDelta > 24 {
		t.Errorf("center tooltip not over its trigger: bubble center off by %.1fpx", centerDelta)
	}
}
