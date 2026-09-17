package visualtest

import (
	"bytes"
	"context"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// kanbanPendingVisualServer serves a single-board page whose move endpoint
// never resolves: e-slow stalls until the browser disconnects (the pending
// register can be framed leisurely) and e-fail answers 500 (the revert path).
// The handler is deliberately stateless — neither branch applies the move —
// so both captures render from the same static starting layout.
func kanbanPendingVisualServer(t *testing.T) *httptest.Server {
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

	mux.HandleFunc("POST /api/kanban/visual-state", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		move, err := display.ParseKanbanMove(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if r.FormValue(kanbanE2ECSRFField) != kanbanE2ECSRFToken {
			http.Error(w, "invalid CSRF token", http.StatusForbidden)

			return
		}

		switch move.Card {
		case "e-fail":
			http.Error(w, "boom", http.StatusInternalServerError)
		case "e-slow":
			// Stall until the tab (and its in-flight fetch) goes away.
			<-r.Context().Done()
		default:
			http.Error(w, "unknown card", http.StatusNotFound)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := kanbanPendingVisualPage().Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// kanbanPendingVisualPage renders one wired board with the flaky card pair.
// The pending spinner ring's animation is frozen inline: this documents the
// exact rendering reduced-motion users get (templates/custom.css disables the
// animation under prefers-reduced-motion: reduce) and makes the golden pixel
// deterministic — a mid-rotation ring would never capture identically twice.
func kanbanPendingVisualPage() templ.Component {
	props := layout.DefaultPageProps()
	props.Title = "Kanban pending register — visual states"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(
			w,
			`<style>[data-tc-kanban] [data-tc-kanban-card].tc-kanban-pending::after{animation:none !important}</style>`,
		)

		return err
	})

	board := display.DefaultKanbanBoardProps()
	board.BaseProps = utils.BaseProps{ID: "kb-visual"}
	board.Columns = kanbanFlakyE2EColumns()
	board.Wire = &wire.Action{URL: "/api/kanban/visual-state"}
	board.CSRFToken = kanbanE2ECSRFToken

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="w-[64rem]">`); err != nil {
			return err
		}

		if err := display.KanbanBoard(board).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</div>`)

		return err
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

// kanbanCaptureBoardState drives one optimistic state to visibility and
// screenshots the board for the named golden: click the move button, poll
// until the state expression holds, settle, capture #kb-visual, and compare
// (or write, with -update / on first run) against the golden PNG.
func kanbanCaptureBoardState(
	t *testing.T,
	srv *httptest.Server,
	name, buttonSel, stateExpr string,
	pollTimeout time.Duration,
) {
	t.Helper()

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(viewportDesktopWidth, viewportDesktopHeight),
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitVisible("#kb-visual", chromedp.ByQuery),
		// Pin light mode before capture: ThemeScript resolves un-pinned
		// pages from prefers-color-scheme, which headless Chromium reports
		// as dark by default (the route-golden lesson — unpinned "light"
		// captures silently render dark).
		chromedp.Evaluate(`try { localStorage.setItem('theme', 'light'); } catch (e) {}`, nil),
		chromedp.Reload(),
		chromedp.WaitVisible("#kb-visual", chromedp.ByQuery),
	); err != nil {
		t.Fatalf("navigate: %v", err)
	}

	kanbanClickOnce(ctx, t, buttonSel)

	var state bool

	if err := chromedp.Run(ctx, pollBool(stateExpr, &state,
		chromedp.WithPollingTimeout(pollTimeout), chromedp.WithPollingInterval(50*time.Millisecond)),
	); err != nil || !state {
		t.Fatalf("%s: state never appeared (state=%v, err=%v)", name, state, err)
	}

	var shot []byte

	if err := chromedp.Run(ctx,
		waitAnimationsSettled(),
		chromedp.Sleep(settleDelay),
		chromedp.Screenshot("#kb-visual", &shot, chromedp.ByQuery, chromedp.NodeVisible),
	); err != nil {
		t.Fatalf("%s: capture: %v", name, err)
	}

	if *update {
		writeGolden(t, name, shot)

		return
	}

	golden, exists := readGolden(t, name)
	if !exists {
		writeGolden(t, name, shot)
		t.Errorf("visualtest[%s]: no golden yet — wrote %s (re-run without -update to verify)", name, goldenPath(name))

		return
	}

	actual, err := png.Decode(bytes.NewReader(shot))
	if err != nil {
		t.Fatalf("visualtest[%s]: decode actual: %v", name, err)
	}

	result, diff := comparePixels(golden, actual, 0.1, 0.1*percentMultiplier)
	if !result.Match {
		writeFailureArtifacts(t, name, shot, diff)
		t.Errorf("visualtest[%s]: visual mismatch — %s (max %.4f%%).\n"+
			"Inspect testdata/.fail/%s.{actual,diff}.png, then run `go test -update` if the change is intended.",
			name, result, 0.1, name)
	}
}

// TestKanbanPendingRegisterVisualStates pins the ADR-0041 states as PNG
// evidence — the optimistic pending register (card in the target column,
// dimmed, spinner ring frozen as reduced-motion users see it, counts synced)
// and the honest failure revert (card back home, red flash, counts restored).
// Behavior under both transports is browser-proven by
// TestKanbanE2EPendingStateBothTransports / TestKanbanE2EFailureRevertsBothTransports.
func TestKanbanPendingRegisterVisualStates(t *testing.T) {
	srv := kanbanPendingVisualServer(t)

	kanbanCaptureBoardState(t, srv, "kanban/pending_state",
		`#kb-visual [data-tc-kanban-card="e-slow"] [data-tc-kanban-move="next"]`,
		kanbanFlakyPendingExpr("kb-visual"), 500*time.Millisecond)

	kanbanCaptureBoardState(t, srv, "kanban/failed_state",
		`#kb-visual [data-tc-kanban-card="e-fail"] [data-tc-kanban-move="next"]`,
		kanbanFlakyRevertedExpr("kb-visual"), 8*time.Second)
}
