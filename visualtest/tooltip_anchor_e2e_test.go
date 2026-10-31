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
// that gate; this e2e forces the bubbles visible and asserts the CLAMP and
// the declarative popover-panel layer instead.
const tooltipAnchorForceStyle = `<style>[data-tc-tooltip] > [role="tooltip"]{display:block !important}</style>`

// tooltipAnchorE2EPage renders: a Tooltip on a trigger pinned against the
// RIGHT viewport edge (the clamp must pull the bubble back inside), and a
// Dropdown near the right edge (the CSS anchor layer must keep its open
// panel inside the viewport).
func tooltipAnchorE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		edgeTooltip := display.TooltipProps{
			ID:       "tt-edge",
			Text:     "Edge tooltip with quite a lot of text that extends far",
			Position: display.PositionTop,
		}
		edgeWrap := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if _, err := io.WriteString(w, `<div class="flex justify-end px-1 pt-4">`); err != nil {
				return err
			}
			inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				_, err := io.WriteString(w, `<button id="tt-edge-btn" type="button" class="rounded border border-gray-400 px-2 py-1">Edge trigger</button>`)

				return err
			})
			if err := display.Tooltip(edgeTooltip).Render(templ.WithChildren(ctx, inner), w); err != nil {
				return err
			}

			_, err := io.WriteString(w, `</div>`)

			return err
		})
		if err := edgeWrap.Render(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `<div class="flex justify-end px-1 pt-24">`); err != nil {
			return err
		}
		edgeDropdown := display.DropdownProps{
			ID:       "dd-edge",
			Label:    "Edge actions",
			Position: display.DropdownPositionEnd,
			Items: []display.DropdownItem{
				{Text: "Rename", Href: "#rename"},
				{Text: "Duplicate with a fairly long label", Href: "#dup"},
			},
		}
		if err := display.Dropdown(edgeDropdown).Render(ctx, w); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `</div>`); err != nil {
			return err
		}

		return nil
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

// TestTooltipAnchorPositioning proves the C23 positioning layers:
//
//  1. The tooltip viewport clamp (tooltipAriaJS singleton) pulls a bubble
//     on a right-edge trigger back inside the viewport on focus.
//  2. On anchor-supporting engines the shared popover JS fallback stays
//     dormant (tcPositionPopover undefined) and the declarative CSS layer
//     (templates/custom.css @supports block, implicit anchor via the
//     invoker) keeps an open Dropdown panel inside the viewport.
//
// Serial — shared browser allocator.
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

	var vw float64
	if err := chromedp.Do(ctx, evalExprInto(`window.innerWidth`, &vw)); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: read viewport width: %v", err)
	}

	// 1. Tooltip clamp: focusing the trigger fires focusin -> the singleton
	// clamps the (forced-visible) bubble into the viewport.
	if err := chromedp.Do(ctx,
		evalVoid(`document.getElementById("tt-edge-btn").focus();`),
	); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: focus edge trigger: %v", err)
	}

	// The clamp runs on double-rAF after focusin; poll until it settles.
	if err := chromedp.Do(ctx, pollTrue(
		`(function(){var b=document.querySelector('#tt-edge-tooltip');if(!b)return false;var r=b.getBoundingClientRect();return r.width>0&&r.right<=window.innerWidth+0.5&&r.left>=-0.5;})()`,
		pollFast()...,
	)); err != nil {
		var edge struct {
			Left, Right float64
		}
		_ = chromedp.Do(ctx, evalExprInto(
			`(function(){var b=document.querySelector('#tt-edge-tooltip');if(!b)return null;var r=b.getBoundingClientRect();return {left:r.left,right:r.right};})()`,
			&edge,
		))
		var diag string
		_ = chromedp.Do(ctx, evalExprInto(
			`(function(){var b=document.querySelector('#tt-edge-tooltip');return JSON.stringify({ml:b.style.marginLeft,disp:getComputedStyle(b).display,active:document.activeElement&&document.activeElement.id,attached:typeof window.tcTooltipAriaAttached})})()`,
			&diag,
		))
		t.Fatalf("tooltip clamp: edge bubble stayed outside the viewport: left=%.1f right=%.1f vw=%.1f diag=%s", edge.Left, edge.Right, vw, diag)
	}

	// 2. Anchor-layer posture: on supporting engines the JS fallback must
	// be dormant — one positioning brain per engine.
	var anchorSupported bool
	if err := chromedp.Do(ctx, evalExprInto(
		`window.CSS&&CSS.supports&&CSS.supports("anchor-name: --x")`,
		&anchorSupported,
	)); err != nil {
		t.Fatalf("visualtest[tooltip-anchor]: probe anchor support: %v", err)
	}

	if anchorSupported {
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

		// Open the edge dropdown (real click on the popovertarget button).
		if err := chromedp.Do(ctx,
			chromedp.WaitVisible(chromedp.CSS("#dd-edge-button")),
			chromedp.Click(chromedp.CSS("#dd-edge-button")),
		); err != nil {
			t.Fatalf("visualtest[tooltip-anchor]: open edge dropdown: %v", err)
		}

		var panelBox struct {
			Left, Right float64
		}
		if err := chromedp.Do(ctx, evalExprInto(
			`(function(){var p=document.getElementById('dd-edge-menu');if(!p||!p.matches(':popover-open'))return null;var r=p.getBoundingClientRect();return {left:r.left,right:r.right};})()`,
			&panelBox,
		)); err != nil {
			t.Fatalf("visualtest[tooltip-anchor]: decode dropdown panel: %v", err)
		}

		if panelBox.Left == 0 && panelBox.Right == 0 {
			t.Fatal("visualtest[tooltip-anchor]: dropdown panel not open")
		}

		if panelBox.Right > vw+0.5 {
			t.Errorf("anchor layer: dropdown panel overflows right: right=%.1f vw=%.1f", panelBox.Right, vw)
		}

		if panelBox.Left < -0.5 {
			t.Errorf("anchor layer: dropdown panel overflows left: left=%.1f", panelBox.Left)
		}
	}
}
