package visualtest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
)

// scrollContainmentE2EPage renders a tall page with an open Modal containing
// its own scrollable region. The native showModal inertness is what must keep
// the page behind the backdrop put: a real (CDP) wheel over the backdrop must
// not move the page, while a wheel over the modal's scrollable body must
// scroll the body.
func scrollContainmentE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Fixture CSS is self-contained on purpose: the compiled demo CSS only
		// carries utilities the demo sources use, so arbitrary Tailwind
		// classes in a test string would silently no-op.
		if _, err := io.WriteString(w,
			`<style>.e2e-tall{height:3000px}.e2e-scrollbox{max-height:10rem;overflow-y:auto}</style>`+
				`<div class="e2e-tall"><p>scroll containment fixture</p></div>`); err != nil {
			return err
		}

		var rows strings.Builder
		for i := 1; i <= 40; i++ {
			fmt.Fprintf(&rows, `<p>Line %d</p>`, i)
		}

		inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w,
				`<div id="e2e-scroll-body" class="e2e-scrollbox">`+rows.String()+`</div>`)

			return err
		})

		return display.Modal(display.ModalProps{
			Nonce: "e2e-nonce", ID: "e2e-scroll-modal",
			Title: "Scroll containment", Open: true,
		}).Render(templ.WithChildren(ctx, inner), w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(layout.DefaultPageProps()).Render(templ.WithChildren(ctx, body), w)
	})
}

func scrollContainmentE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = scrollContainmentE2EPage().Render(r.Context(), w)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// wheelAt dispatches a real CDP mouse-wheel tick at viewport coordinates, so
// the browser's own hit-testing decides what scrolls — the same path a user's
// trackpad takes.
func wheelAt(x, y, deltaY float64) chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, tab *chromedp.Target) error {
		// Raw JSON params: Chromium requires BOTH deltaX and deltaY keys
		// present for mouseWheel events, and cdproto's omitempty drops a
		// zero deltaX — which Chrome then rejects with -32602.
		params := map[string]any{
			"type":   "mouseWheel",
			"x":      x,
			"y":      y,
			"deltaX": 0,
			"deltaY": deltaY,
		}

		err := tab.Call(ctx, "Input.dispatchMouseEvent", params, &cdp.Empty{})

		return err
	})
}

// TestModalE2EScrollContainment proves the modal's scroll containment with
// real input events: a wheel over the ::backdrop region must NOT scroll the
// page behind the modal dialog (native showModal inertness is the guarantee),
// while a wheel over the dialog's scrollable body must scroll the body and
// still leave the page untouched. Serial — shared browser allocator.
func TestModalE2EScrollContainment(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := scrollContainmentE2EServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitReady("body"),
		pollTrue(`document.getElementById('e2e-scroll-modal').open===true`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[scroll-containment]: open modal: %v", err)
	}

	// Put the page at a known non-zero scroll offset; any backdrop leak moves
	// it from here.
	if err := chromedp.Do(ctx,
		evalVoid(`window.scrollTo(0, 800);`),
		pollTrue(`document.scrollingElement.scrollTop===800`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[scroll-containment]: scroll page: %v", err)
	}

	// Wheel over the backdrop (top-left corner is guaranteed outside the
	// centered panel): the page must not move.
	if err := chromedp.Do(ctx,
		wheelAt(40, 40, 600),
		chromedp.Sleep(150*time.Millisecond),
		pollTrue(`document.scrollingElement.scrollTop===800`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[scroll-containment]: backdrop wheel scrolled the page: %v", err)
	}

	// Locate the scrollable body's centre and wheel over it: the BODY scrolls.
	var center struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	if err := chromedp.Do(ctx, evalExprInto(
		`(function(){var r=document.getElementById('e2e-scroll-body').getBoundingClientRect();return {x:r.left+r.width/2,y:r.top+r.height/2};})()`,
		&center,
	)); err != nil {
		t.Fatalf("visualtest[scroll-containment]: locate scroll body: %v", err)
	}

	if err := chromedp.Do(ctx,
		wheelAt(center.X, center.Y, 200),
		chromedp.Sleep(150*time.Millisecond),
		pollTrue(`document.getElementById('e2e-scroll-body').scrollTop>0`, pollFast()...),
		pollTrue(`document.scrollingElement.scrollTop===800`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[scroll-containment]: body wheel did not scroll body (or leaked): %v", err)
	}

	var bodyTop, pageTop float64

	if err := chromedp.Do(ctx, evalExprInto(
		`document.getElementById('e2e-scroll-body').scrollTop`, &bodyTop,
	), evalExprInto(`document.scrollingElement.scrollTop`, &pageTop)); err != nil {
		t.Fatalf("visualtest[scroll-containment]: read scroll offsets: %v", err)
	}

	t.Logf("body scrollTop=%.0f, page scrollTop=%.0f (want >0 / ==800)", bodyTop, pageTop)
}
