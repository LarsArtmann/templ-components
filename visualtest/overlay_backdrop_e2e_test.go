package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
)

// The backdrop-dismissal contract of every <dialog>-based overlay (Modal,
// Drawer, CommandPalette, Lightbox) is POSITION-based, never target-based:
//
//   - a press AND release outside the dialog closes it;
//   - a press inside the dialog never closes it, even when the event target is
//     the <dialog> element itself (which is what a click on the dialog's own
//     padding or on the gaps between its children targets);
//   - a press inside followed by a release outside (text selection drag) does
//     not dismiss.
//
// The first rule is the one that regressed once already: `event.target ===
// dialog` closed the modal when the user clicked inside it. These tests render
// the real component in the real browser and drive it with real mouse input,
// so the predicate is exercised the way Chromium actually delivers it.

// overlayBox is the decoded result of getBoundingClientRect for the dialog.
type overlayBox struct {
	Left   float64 `json:"left"`
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
}

func (b overlayBox) insideX() float64 { return (b.Left + b.Right) / 2 }
func (b overlayBox) insideY() float64 { return (b.Top + b.Bottom) / 2 }

// outsidePoint returns a viewport point guaranteed to lie outside the dialog
// box with a 20px margin, clamped away from the viewport edge so the press
// lands inside the page.
func (b overlayBox) outsidePoint() (float64, float64) {
	x := b.Left - 20
	if x < 2 {
		x = 2
	}

	y := b.Top - 20
	if y < 2 {
		y = 2
	}

	return x, y
}

// mouseEvent builds a single CDP mouse dispatch at viewport coordinates. The
// real input pipeline is what turns these into pointerdown/pointerup, so the
// component's pointer listeners see exactly what a user produces.
func mouseEvent(typ input.DispatchMouseEventType, x, y float64) chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, tab *chromedp.Target) error {
		//nolint:exhaustruct_v5 // omitted CDP fields take protocol defaults
		params := input.DispatchMouseEventParams{
			Type:   typ,
			X:      x,
			Y:      y,
			Button: input.MouseButtonLeft,
		}
		if typ != input.DispatchMouseEventTypeMouseMoved {
			params.ClickCount = 1
		}

		_, err := cdp.Call(ctx, tab, input.DispatchMouseEvent, params)

		return err
	})
}

// clickAt dispatches a full press+release pair at viewport coordinates.
func clickAt(x, y float64) chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, tab *chromedp.Target) error {
		if _, err := mouseEvent(input.DispatchMouseEventTypeMousePressed, x, y)(ctx, tab); err != nil {
			return err
		}

		_, err := mouseEvent(input.DispatchMouseEventTypeMouseReleased, x, y)(ctx, tab)

		return err
	})
}

// dragFromTo presses at (x1,y1), moves through the midpoint, and releases at
// (x2,y2) — the shape of a text-selection drag that happens to end outside.
func dragFromTo(x1, y1, x2, y2 float64) chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, tab *chromedp.Target) error {
		if _, err := mouseEvent(input.DispatchMouseEventTypeMousePressed, x1, y1)(ctx, tab); err != nil {
			return err
		}

		if _, err := mouseEvent(
			input.DispatchMouseEventTypeMouseMoved,
			(x1+x2)/2,
			(y1+y2)/2,
		)(ctx, tab); err != nil {
			return err
		}

		_, err := mouseEvent(input.DispatchMouseEventTypeMouseReleased, x2, y2)(ctx, tab)

		return err
	})
}

// modalOpen reports whether the dialog is currently in the top layer.
func modalOpen(t *testing.T, ctx context.Context, id string) bool {
	t.Helper()

	var open bool
	if err := chromedp.Do(ctx, evalExprInto(
		`(() => { const d = document.getElementById(`+jsString(id)+`); return !!d && d.open; })()`,
		&open,
	)); err != nil {
		t.Fatalf("read dialog open state: %v", err)
	}

	return open
}

// dialogBox reads the dialog's viewport rectangle.
func dialogBox(t *testing.T, ctx context.Context, id string) overlayBox {
	t.Helper()

	var box overlayBox
	if err := chromedp.Do(ctx, evalExprInto(
		`(() => { const r = document.getElementById(`+jsString(id)+`).getBoundingClientRect();`+
			`return {left: r.left, top: r.top, right: r.right, bottom: r.bottom}; })()`,
		&box,
	)); err != nil {
		t.Fatalf("read dialog box: %v", err)
	}

	return box
}

// overlayBackdropE2EPage renders a trigger plus a titled Modal, so the close
// button (rendered only when Title is set) is present too.
func overlayBackdropE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w,
			`<button type="button" id="e2e-open" data-tc-overlay-trigger="e2e-modal">Open modal</button>`)
		if err != nil {
			return err
		}

		props := display.DefaultModalProps()
		props.ID = "e2e-modal"
		props.Title = "Confirm"
		props.Nonce = "e2e-nonce"

		children := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<p id="e2e-modal-body">Are you sure you want to continue?</p>`)

			return err
		})

		return display.Modal(props).Render(templ.WithChildren(ctx, children), w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(layout.DefaultPageProps()).Render(templ.WithChildren(ctx, body), w)
	})
}

func overlayBackdropE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = overlayBackdropE2EPage().Render(r.Context(), w)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// openOverlay navigates to the fixture page and opens the modal through the
// declarative [data-tc-overlay-trigger] path (a real click, not showModal()).
func openOverlay(t *testing.T, ctx context.Context, url string) {
	t.Helper()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		pollTrue(`!!document.getElementById('e2e-open')`, pollFast()...),
		chromedp.Click(chromedp.CSS("#e2e-open")),
		pollTrue(`document.getElementById('e2e-modal').open`, pollFast()...),
	); err != nil {
		t.Fatalf("open modal via trigger: %v", err)
	}
}

// TestModalInsideClickKeepsOverlayOpen is the regression test for the
// `event.target === dialog` dismissal bug: a click inside the modal must never
// close it. The pointer-position predicate is checked by clicking the dialog's
// own box — the region that targets the <dialog> element in Chromium.
func TestModalInsideClickKeepsOverlayOpen(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := overlayBackdropE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	tabCtx, tabCancel := context.WithTimeout(ctx, 120*time.Second)
	defer tabCancel()

	if err := chromedp.Do(tabCtx, chromedp.EmulateViewport(640, 480)); err != nil {
		t.Fatalf("set viewport: %v", err)
	}

	openOverlay(t, tabCtx, srv.URL)

	box := dialogBox(t, tabCtx, "e2e-modal")

	// Click the dialog's top-left inside corner: 4px in from the border, in
	// the header's own box — a point whose event target is the dialog (or its
	// header), never the backdrop.
	if err := chromedp.Do(tabCtx, clickAt(box.Left+4, box.Top+4)); err != nil {
		t.Fatalf("click inside dialog: %v", err)
	}

	// Wait past the press/release pair before asserting: a wrongly-wired
	// dismissal runs on pointerup, not later.
	time.Sleep(300 * time.Millisecond)

	if !modalOpen(t, tabCtx, "e2e-modal") {
		t.Fatal("clicking inside the modal closed it — backdrop dismissal is target-based again")
	}

	// A synthetic pointerdown/pointerup pair dispatched directly on the dialog
	// with INSIDE coordinates reproduces the exact event shape a backdrop click
	// has (target === dialog). The position predicate must keep it open.
	if err := chromedp.Do(tabCtx, evalVoid(
		`(() => { const d = document.getElementById('e2e-modal');`+
			`const r = d.getBoundingClientRect();`+
			`const opts = {bubbles: true, cancelable: true, clientX: r.left + 8, clientY: r.top + 8};`+
			`d.dispatchEvent(new PointerEvent('pointerdown', opts));`+
			`d.dispatchEvent(new PointerEvent('pointerup', opts)); })()`,
	)); err != nil {
		t.Fatalf("dispatch inside pointer events: %v", err)
	}

	time.Sleep(300 * time.Millisecond)

	if !modalOpen(t, tabCtx, "e2e-modal") {
		t.Fatal("pointer events with inside coordinates on the dialog closed it — predicate is not position-based")
	}
}

// TestModalBackdropClickDismisses proves the positive half: press and release
// outside the dialog closes it.
func TestModalBackdropClickDismisses(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := overlayBackdropE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	tabCtx, tabCancel := context.WithTimeout(ctx, 120*time.Second)
	defer tabCancel()

	if err := chromedp.Do(tabCtx, chromedp.EmulateViewport(640, 480)); err != nil {
		t.Fatalf("set viewport: %v", err)
	}

	openOverlay(t, tabCtx, srv.URL)

	box := dialogBox(t, tabCtx, "e2e-modal")
	outsideX, outsideY := box.outsidePoint()

	if err := chromedp.Do(tabCtx,
		clickAt(outsideX, outsideY),
		pollTrue(`!document.getElementById('e2e-modal').open`, pollFast()...),
	); err != nil {
		t.Fatalf("backdrop click did not dismiss the modal: %v", err)
	}
}

// TestModalCloseButtonDismisses covers the wired [data-tc-close] path.
func TestModalCloseButtonDismisses(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := overlayBackdropE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	tabCtx, tabCancel := context.WithTimeout(ctx, 120*time.Second)
	defer tabCancel()

	openOverlay(t, tabCtx, srv.URL)

	if err := chromedp.Do(tabCtx,
		chromedp.Click(chromedp.CSS("#e2e-modal [data-tc-close]")),
		pollTrue(`!document.getElementById('e2e-modal').open`, pollFast()...),
	); err != nil {
		t.Fatalf("close button did not dismiss the modal: %v", err)
	}
}

// TestModalDragOutDoesNotDismiss proves a press inside released outside (the
// text-selection drag shape) is not a dismissal, and its mirror: a press
// outside released inside is not one either.
func TestModalDragOutDoesNotDismiss(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := overlayBackdropE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	tabCtx, tabCancel := context.WithTimeout(ctx, 120*time.Second)
	defer tabCancel()

	if err := chromedp.Do(tabCtx, chromedp.EmulateViewport(640, 480)); err != nil {
		t.Fatalf("set viewport: %v", err)
	}

	openOverlay(t, tabCtx, srv.URL)

	box := dialogBox(t, tabCtx, "e2e-modal")
	outsideX, outsideY := box.outsidePoint()

	if err := chromedp.Do(tabCtx,
		dragFromTo(box.insideX(), box.insideY(), outsideX, outsideY),
	); err != nil {
		t.Fatalf("drag inside to outside: %v", err)
	}

	time.Sleep(300 * time.Millisecond)

	if !modalOpen(t, tabCtx, "e2e-modal") {
		t.Fatal("press inside + release outside dismissed the modal — press and release must BOTH land outside")
	}

	// Mirror case: the release is inside the dialog, so the press cannot have
	// started on the backdrop either.
	if err := chromedp.Do(tabCtx,
		dragFromTo(outsideX, outsideY, box.insideX(), box.insideY()),
	); err != nil {
		t.Fatalf("drag outside to inside: %v", err)
	}

	time.Sleep(300 * time.Millisecond)

	if !modalOpen(t, tabCtx, "e2e-modal") {
		t.Fatal("press outside + release inside dismissed the modal")
	}
}
