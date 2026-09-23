package visualtest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/go-datastar/static"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// The kanban E2E suite proves display.KanbanBoard's move pipeline at the
// browser level under BOTH real runtimes: the keyboard move buttons (the
// accessibility path) and synthetic HTML5 drag-and-drop (the pointer path)
// must both land the card in the target column through the same hidden form
// — htmx via hx-post + hx-target/hx-swap, Datastar via @post with form
// encoding and response-header targeting from wire.Handler.
//
// Run via `nix run .#visual`; tests skip gracefully without a browser.
//
// ANTI-DRIFT TIE (TODO #227): this server's add/reset/CSRF/same-origin
// contract MIRRORS examples/demo/kanban_demo.go (which cannot be imported —
// the demo runs as an external binary). The two are pinned together by ONE
// probe table: runKanbanContractProbes in demo_kanban_http_test.go drives
// TestDemoKanbanHTTPContracts (real demo) and TestKanbanE2EHTTPContractParity
// (this server) through identical probes. Diverge from the demo's guards and
// the parity test names this file.

// kanbanE2EBoard is the mutex-guarded server state one dialect endpoint owns.
type kanbanE2EBoard struct {
	mu      sync.Mutex
	columns []display.KanbanColumn
	next    int64
	initial func() []display.KanbanColumn
}

func newKanbanE2EBoard() *kanbanE2EBoard {
	return newKanbanE2EBoardWith(initialKanbanE2EColumns)
}

// newKanbanE2EBoardWith builds a board that always resets to the given
// layout — the flaky boards use their own slow/failing card mix.
func newKanbanE2EBoardWith(initial func() []display.KanbanColumn) *kanbanE2EBoard {
	b := &kanbanE2EBoard{initial: initial}
	b.reset()

	return b
}

// initialKanbanE2EColumns returns the board layout every test starts from.
func initialKanbanE2EColumns() []display.KanbanColumn {
	return []display.KanbanColumn{
		{
			ID:    "todo",
			Title: "To do",
			Cards: []display.KanbanCard{{ID: "e1", Title: "First"}, {ID: "e2", Title: "Second"}},
		},
		{ID: "doing", Title: "In progress"},
		{ID: "done", Title: "Done", Cards: []display.KanbanCard{{ID: "e3", Title: "Third"}}},
	}
}

// kanbanFlakyE2EColumns arms the optimistic-pending proof: one card whose
// move endpoint stalls 1.2s (the pending window) and one whose move endpoint
// always 500s (the revert path).
func kanbanFlakyE2EColumns() []display.KanbanColumn {
	return []display.KanbanColumn{
		{
			ID:    "todo",
			Title: "To do",
			Cards: []display.KanbanCard{
				{ID: "e-slow", Title: "Slow move"},
				{ID: "e-fail", Title: "Doomed move"},
			},
		},
		{ID: "doing", Title: "In progress"},
	}
}

// kanbanFlakyMoveDelay is the stall the e-slow move endpoint applies — long
// enough that the optimistic pending state is deterministically observable
// before the response lands.
const kanbanFlakyMoveDelay = 1200 * time.Millisecond

// reset restores the initial layout (tests share the package-global boards
// and must not observe each other's moves).
func (b *kanbanE2EBoard) reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	initial := b.initial
	if initial == nil {
		initial = initialKanbanE2EColumns
	}

	b.columns = initial()
	b.next = 0
}

// apply mutates the board by one decoded move: remove the card from
// wherever it lives, insert it at the (clamped) requested position.
func (b *kanbanE2EBoard) apply(move display.KanbanMove) {
	b.mu.Lock()
	defer b.mu.Unlock()

	var card display.KanbanCard

	found := false

	for ci := range b.columns {
		for pi := range b.columns[ci].Cards {
			if b.columns[ci].Cards[pi].ID == move.Card {
				card = b.columns[ci].Cards[pi]
				b.columns[ci].Cards = append(b.columns[ci].Cards[:pi], b.columns[ci].Cards[pi+1:]...)
				found = true

				break
			}
		}
	}

	if !found {
		return
	}

	for ci := range b.columns {
		if b.columns[ci].ID != move.Column {
			continue
		}

		if move.Index > len(b.columns[ci].Cards) {
			move.Index = len(b.columns[ci].Cards)
		}

		cards := make([]display.KanbanCard, 0, len(b.columns[ci].Cards)+1)
		cards = append(cards, b.columns[ci].Cards[:move.Index]...)
		cards = append(cards, card)
		cards = append(cards, b.columns[ci].Cards[move.Index:]...)
		b.columns[ci].Cards = cards

		return
	}
}

// add appends a fresh card to the named column and reports whether the
// column exists — the handler turns false into a 404 so a broken link
// surfaces instead of silently no-oping.
func (b *kanbanE2EBoard) add(columnID string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	for ci := range b.columns {
		if b.columns[ci].ID != columnID {
			continue
		}

		n := b.next
		b.next++
		b.columns[ci].Cards = append(b.columns[ci].Cards, display.KanbanCard{
			ID:    "ea" + strconv.FormatInt(n, 10),
			Title: "Added " + strconv.FormatInt(n+1, 10),
		})

		return true
	}

	return false
}

// kanbanE2EPage renders both boards on one consumer page shell: layout.Base
// injects the self-hosted htmx runtime; the Datastar runtime loads from
// /datastar.js (the pinned bundle) with a datastar-ready catcher registered
// before the module executes.
func kanbanE2EPage() templ.Component {
	props := layout.DefaultPageProps()
	props.Title = "Kanban E2E — templ-components"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(
			w,
			`<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`,
		)

		return err
	})

	htmxAction := wire.Action{URL: "/api/kanban/htmx"}
	datastarAction := wire.Action{Transport: wire.TransportDatastar, URL: "/api/kanban/datastar"}

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := datastar.SDKScript(datastar.SDKScriptProps{
			BaseProps: utils.BaseProps{Nonce: "kanban-e2e-nonce"},
			Src:       "/datastar.js",
		}).Render(ctx, w); err != nil {
			return err
		}

		// Each board sits in a wrapper div carrying its reset button,
		// mirroring the demo's per-board header layout.
		if _, err := io.WriteString(w, `<div id="kb-htmx-wrap">`); err != nil {
			return err
		}

		if err := kanbanE2EResetButton(wire.Action{
			URL:    "/api/kanban/htmx/reset",
			Method: wire.MethodPost,
			Target: "#kb-htmx",
		}).Render(ctx, w); err != nil {
			return err
		}

		htmxProps := kanbanHTMXBoard.kanbanE2EBoardProps("kb-htmx", htmxAction)
		if err := display.KanbanBoard(htmxProps).Render(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `</div><div id="kb-ds-wrap">`); err != nil {
			return err
		}

		if err := kanbanE2EResetButton(wire.Action{
			Transport: wire.TransportDatastar,
			Method:    wire.MethodPost,
			URL:       "/api/kanban/datastar/reset",
		}).Render(ctx, w); err != nil {
			return err
		}

		datastarProps := kanbanDatastarBoard.kanbanE2EBoardProps("kb-ds", datastarAction)
		if err := display.KanbanBoard(datastarProps).Render(ctx, w); err != nil {
			return err
		}

		// The flaky boards prove the optimistic-pending register: e-slow's
		// move stalls 1.2s, e-fail's move always 500s (see the flaky handlers
		// in kanbanE2EServer).
		if _, err := io.WriteString(
			w,
			`</div><div id="kb-flaky-wrap"><h3 class="text-sm font-semibold text-gray-900 dark:text-white">Optimistic pending register (slow + failing moves)</h3>`,
		); err != nil {
			return err
		}

		flakyHTMXProps := kanbanFlakyHTMXBoard.kanbanFlakyBoardProps(
			"kb-htmx-flaky",
			wire.Action{URL: "/api/kanban/htmx-flaky"},
		)
		if err := display.KanbanBoard(flakyHTMXProps).Render(ctx, w); err != nil {
			return err
		}

		if err := display.KanbanBoard(kanbanFlakyDatastarBoard.kanbanFlakyBoardProps(
			"kb-ds-flaky",
			wire.Action{Transport: wire.TransportDatastar, URL: "/api/kanban/datastar-flaky"},
		)).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</div>`)

		return err
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

// kanbanE2EServer serves the E2E page, its assets, and the two move
// endpoints — the Datastar one wrapped in the production wire.Handler so
// its response owns the patch targeting (outer mode, the board's own id).
func kanbanE2EServer(t *testing.T) *httptest.Server {
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

	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(static.Bytes())
	})

	moveHandler := func(board *kanbanE2EBoard, id string, action wire.Action) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")

			move, ok := kanbanParseMoveWithCSRF(w, r)
			if !ok {
				return
			}

			board.apply(move)

			if err := display.KanbanBoard(board.kanbanE2EBoardProps(id, action)).Render(r.Context(), w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}

	mux.Handle("POST /api/kanban/htmx", moveHandler(kanbanHTMXBoard, "kb-htmx", wire.Action{URL: "/api/kanban/htmx"}))
	mux.Handle("POST /api/kanban/datastar", wire.Handler(wire.PatchTarget{
		Selector: "#kb-ds",
		Mode:     wire.PatchModeOuter,
	}, moveHandler(kanbanDatastarBoard, "kb-ds", wire.Action{Transport: wire.TransportDatastar, URL: "/api/kanban/datastar"})))

	// Flaky move endpoints for the optimistic-pending proof: e-fail always
	// 500s (revert path), e-slow stalls kanbanFlakyMoveDelay (pending window)
	// before the normal apply + re-render.
	flakyMoveHandler := func(board *kanbanE2EBoard, id string, action wire.Action) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")

			move, ok := kanbanParseMoveWithCSRF(w, r)
			if !ok {
				return
			}

			switch move.Card {
			case "e-fail":
				http.Error(w, "boom", http.StatusInternalServerError)

				return
			case "e-slow":
				time.Sleep(kanbanFlakyMoveDelay)
			}

			board.apply(move)

			if err := display.KanbanBoard(board.kanbanFlakyBoardProps(id, action)).Render(r.Context(), w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
	mux.Handle("POST /api/kanban/htmx-flaky", flakyMoveHandler(
		kanbanFlakyHTMXBoard,
		"kb-htmx-flaky",
		wire.Action{URL: "/api/kanban/htmx-flaky"},
	))
	mux.Handle("POST /api/kanban/datastar-flaky", wire.Handler(wire.PatchTarget{
		Selector: "#kb-ds-flaky",
		Mode:     wire.PatchModeOuter,
	}, flakyMoveHandler(
		kanbanFlakyDatastarBoard,
		"kb-ds-flaky",
		wire.Action{Transport: wire.TransportDatastar, URL: "/api/kanban/datastar-flaky"},
	)))

	// The column Action slot's add endpoint: column id in the path, no
	// request body, same-origin enforced — mirrors the demo.
	addHandler := func(board *kanbanE2EBoard, id string, action wire.Action) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")

			if !kanbanE2ESameOrigin(r) {
				http.Error(w, "cross-origin request rejected", http.StatusForbidden)

				return
			}

			column := r.PathValue("column")
			if column == "" {
				http.Error(w, "missing column path parameter", http.StatusBadRequest)

				return
			}

			if !board.add(column) {
				http.Error(w, "unknown column: "+column, http.StatusNotFound)

				return
			}

			if err := display.KanbanBoard(board.kanbanE2EBoardProps(id, action)).Render(r.Context(), w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}

	// Board reset: restores the starting layout (mirrors the demo).
	resetHandler := func(board *kanbanE2EBoard, id string, action wire.Action) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")

			if !kanbanE2ESameOrigin(r) {
				http.Error(w, "cross-origin request rejected", http.StatusForbidden)

				return
			}

			board.reset()

			if err := display.KanbanBoard(board.kanbanE2EBoardProps(id, action)).Render(r.Context(), w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}

	mux.Handle(
		"POST /api/kanban/htmx/add/{column}",
		addHandler(kanbanHTMXBoard, "kb-htmx", wire.Action{URL: "/api/kanban/htmx"}),
	)
	mux.Handle("POST /api/kanban/datastar/add/{column}", wire.Handler(wire.PatchTarget{
		Selector: "#kb-ds",
		Mode:     wire.PatchModeOuter,
	}, addHandler(kanbanDatastarBoard, "kb-ds", wire.Action{Transport: wire.TransportDatastar, URL: "/api/kanban/datastar"})))
	mux.Handle(
		"POST /api/kanban/htmx/reset",
		resetHandler(kanbanHTMXBoard, "kb-htmx", wire.Action{URL: "/api/kanban/htmx"}),
	)
	mux.Handle("POST /api/kanban/datastar/reset", wire.Handler(wire.PatchTarget{
		Selector: "#kb-ds",
		Mode:     wire.PatchModeOuter,
	}, resetHandler(kanbanDatastarBoard, "kb-ds", wire.Action{Transport: wire.TransportDatastar, URL: "/api/kanban/datastar"})))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := kanbanE2EPage().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// The two boards, one per transport, each with independent state.
var (
	kanbanHTMXBoard     = newKanbanE2EBoard()
	kanbanDatastarBoard = newKanbanE2EBoard()

	// The flaky boards prove the optimistic-pending register (slow + failing
	// moves) under both runtimes.
	kanbanFlakyHTMXBoard     = newKanbanE2EBoardWith(kanbanFlakyE2EColumns)
	kanbanFlakyDatastarBoard = newKanbanE2EBoardWith(kanbanFlakyE2EColumns)
)

// kanbanE2ECSRFToken is the e2e server's CSRF token; every move handler
// rejects form posts without it, so the browser tests prove the token
// round-trips through the move form's hidden input under both runtimes.
var kanbanE2ECSRFToken = newKanbanE2ECSRFToken()

// kanbanE2ECSRFField matches KanbanBoardProps's default hidden-input name.
const kanbanE2ECSRFField = "csrf_token"

func newKanbanE2ECSRFToken() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic("kanban e2e: generate csrf token: " + err.Error())
	}

	return hex.EncodeToString(b[:])
}

// kanbanE2ESameOrigin mirrors the demo's same-origin enforcement for the
// bodyless add/reset POSTs: Sec-Fetch-Site, then an Origin/host match.
func kanbanE2ESameOrigin(r *http.Request) bool {
	if r.Header.Get("Sec-Fetch-Site") == "same-origin" {
		return true
	}

	origin := r.Header.Get("Origin")

	return origin == "http://"+r.Host || origin == "https://"+r.Host
}

// kanbanE2EAddButton mirrors the demo's add affordance: a ghost Button
// whose Wire attributes carry the transport dialect. htmx targets the board
// root (Target is set by the caller) and must swap outerHTML — otherwise
// the whole-board response lands inside this button.
func kanbanE2EAddButton(action wire.Action, columnTitle string) templ.Component {
	base := utils.BaseProps{AriaLabel: "Add card to " + columnTitle}
	if action.Transport != wire.TransportDatastar {
		base.Attrs = templ.Attributes{"hx-swap": "outerHTML"}
	}

	return display.Button(display.ButtonProps{
		BaseProps: base,
		Text:      "+ Add",
		Variant:   display.ButtonGhost,
		Size:      display.ButtonSizeSM,
		Wire:      &action,
	})
}

// kanbanE2EResetButton mirrors the demo's reset affordance.
func kanbanE2EResetButton(action wire.Action) templ.Component {
	base := utils.BaseProps{AriaLabel: "Reset demo board"}
	if action.Transport != wire.TransportDatastar {
		base.Attrs = templ.Attributes{"hx-swap": "outerHTML"}
	}

	return display.Button(display.ButtonProps{
		BaseProps: base,
		Text:      "Reset",
		Variant:   display.ButtonGhost,
		Size:      display.ButtonSizeSM,
		Wire:      &action,
	})
}

// kanbanE2EBoardProps snapshots the state into render props with a stable
// board id, the transport's wire action, a CSRF token on the move form, and
// a per-column add-card button in the Action slot (mirroring the demo).
func (b *kanbanE2EBoard) kanbanE2EBoardProps(id string, action wire.Action) display.KanbanBoardProps {
	b.mu.Lock()
	defer b.mu.Unlock()

	snapshot := make([]display.KanbanColumn, len(b.columns))
	for ci, col := range b.columns {
		add := action
		add.URL += "/add/" + col.ID
		add.Method = wire.MethodPost
		add.Event = wire.EventClick
		add.Target = "#" + id
		snapshot[ci] = display.KanbanColumn{
			ID:     col.ID,
			Title:  col.Title,
			Cards:  append([]display.KanbanCard(nil), col.Cards...),
			Action: kanbanE2EAddButton(add, col.Title),
		}
	}

	props := display.DefaultKanbanBoardProps()
	props.BaseProps = utils.BaseProps{ID: id}
	props.Columns = snapshot
	props.Wire = &action
	props.CSRFToken = kanbanE2ECSRFToken

	return props
}

// kanbanFlakyBoardProps snapshots the flaky board state into render props —
// like kanbanE2EBoardProps but without the add-card Action slots, which the
// pending/revert tests never use.
func (b *kanbanE2EBoard) kanbanFlakyBoardProps(id string, action wire.Action) display.KanbanBoardProps {
	b.mu.Lock()
	defer b.mu.Unlock()

	snapshot := make([]display.KanbanColumn, len(b.columns))
	for ci, col := range b.columns {
		snapshot[ci] = display.KanbanColumn{
			ID:    col.ID,
			Title: col.Title,
			Cards: append([]display.KanbanCard(nil), col.Cards...),
		}
	}

	props := display.DefaultKanbanBoardProps()
	props.BaseProps = utils.BaseProps{ID: id}
	props.Columns = snapshot
	props.Wire = &action
	props.CSRFToken = kanbanE2ECSRFToken

	return props
}

// kanbanColumnOrderExpr returns a JS expression evaluating to the ordered
// comma-separated card ids of one column, e.g. "e1,e2".
func kanbanColumnOrderExpr(boardID, columnID string) string {
	return fmt.Sprintf(
		`Array.from(document.querySelectorAll('#%s [data-tc-kanban-column-body="%s"] > [data-tc-kanban-card]')).map(function(el){return el.getAttribute('data-tc-kanban-card');}).join(',')`,
		boardID,
		columnID,
	)
}

// kanbanE2EReady returns the readiness poll expression: htmx is inline
// (self-hosted), Datastar announces itself via datastar-ready.
const kanbanE2EReady = `document.readyState==='complete' && window.htmx!==undefined && window.__dsReady===true`

// kanbanClickMoveUntil clicks a move button and polls for the expected
// column order, retrying the click a few times — a click issued while a
// runtime re-attaches to swapped-in markup can be swallowed (the
// wire-forms-pack lesson).
func kanbanClickMoveUntil(ctx context.Context, t *testing.T, buttonSel, orderExpr, want string) {
	t.Helper()

	poll := orderExpr + "===" + strconv.Quote(want)

	for range 3 {
		// No NodeVisible: the move buttons are opacity-0 until hover/focus,
		// which chromedp's visibility check treats as hidden. Clicking the
		// coordinates still delivers the event.
		if err := chromedp.Run(ctx,
			chromedp.Click(buttonSel),
		); err != nil {
			t.Fatalf("click %s: %v", buttonSel, err)
		}

		var got bool

		err := chromedp.Run(ctx, pollBool(poll, &got))
		if err == nil && got {
			return
		}

		time.Sleep(300 * time.Millisecond)
	}

	t.Fatalf("card never reached the expected column order %q via %s", want, buttonSel)
}

// newKanbanReadyTab bootstraps one kanban e2e test: fresh boards, a demo
// server, and a browser tab (45s budget) parked on "/" with both boards
// reporting ready. The returned context is cancelled via t.Cleanup.
func newKanbanReadyTab(t *testing.T) context.Context {
	t.Helper()

	kanbanHTMXBoard.reset()
	kanbanDatastarBoard.reset()

	srv := kanbanE2EServer(t)

	ctx, cancel := newTab(t)
	t.Cleanup(cancel)

	ctx, cancelTimeout := context.WithTimeout(ctx, 45*time.Second)
	t.Cleanup(cancelTimeout)

	var ready bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollBool(kanbanE2EReady, &ready),
	); err != nil {
		t.Fatalf("navigate + readiness: %v", err)
	}

	return ctx
}

// TestKanbanE2EKeyboardMovesBothTransports proves the accessibility path:
// clicking a card's "move to next column" button submits the hidden form
// and the re-rendered board carries the card in the target column — under
// htmx AND the real Datastar runtime.
func TestKanbanE2EKeyboardMovesBothTransports(t *testing.T) {
	ctx := newKanbanReadyTab(t)

	// htmx board: e1 moves to the end of "doing" (empty column).
	kanbanClickMoveUntil(ctx, t,
		`#kb-htmx [data-tc-kanban-card="e1"] [data-tc-kanban-move="next"]`,
		kanbanColumnOrderExpr("kb-htmx", "doing"), "e1")

	// Datastar board: same move through the same markup.
	kanbanClickMoveUntil(ctx, t,
		`#kb-ds [data-tc-kanban-card="e1"] [data-tc-kanban-move="next"]`,
		kanbanColumnOrderExpr("kb-ds", "doing"), "e1")

	// The moved card must keep its move identity (the buttons are on the
	// re-rendered card, so it can move on).
	kanbanClickMoveUntil(ctx, t,
		`#kb-htmx [data-tc-kanban-card="e1"] [data-tc-kanban-move="next"]`,
		kanbanColumnOrderExpr("kb-htmx", "done"), "e3,e1")
}

// kanbanDropScript returns the JS that simulates one HTML5 drag-and-drop:
// dragstart on the card, dragover + drop on the target zone at clientY 0
// (above every card → insert at index 0), dragend for cleanup.
func kanbanDropScript(boardID, cardID, columnID string) string {
	return fmt.Sprintf(`(function(){
		var b=document.getElementById(%q);
		var dt=new DataTransfer();
		var card=b.querySelector('[data-tc-kanban-card="%s"]');
		var zone=b.querySelector('[data-tc-kanban-column-body="%s"]');
		if(!card||!zone){return 'missing';}
		card.dispatchEvent(new DragEvent('dragstart',{bubbles:true,cancelable:true,dataTransfer:dt}));
		zone.dispatchEvent(new DragEvent('dragover',{bubbles:true,cancelable:true,clientY:0,dataTransfer:dt}));
		zone.dispatchEvent(new DragEvent('drop',{bubbles:true,cancelable:true,dataTransfer:dt}));
		card.dispatchEvent(new DragEvent('dragend',{bubbles:true,cancelable:true,dataTransfer:dt}));
		return 'dropped';
	})()`, boardID, cardID, columnID)
}

// TestKanbanE2EDropMovesBothTransports proves the pointer path: synthetic
// drag events route through the same hidden form (index 0 = insert at top
// of the target column) under both runtimes.
func TestKanbanE2EDropMovesBothTransports(t *testing.T) {
	ctx := newKanbanReadyTab(t)

	drops := []struct {
		boardID string
		card    string
		column  string
		want    string
	}{
		{boardID: "kb-htmx", card: "e2", column: "done", want: "e2,e3"},
		{boardID: "kb-ds", card: "e2", column: "done", want: "e2,e3"},
	}

	for _, drop := range drops {
		var dropped string

		if err := chromedp.Run(ctx,
			chromedp.Evaluate(kanbanDropScript(drop.boardID, drop.card, drop.column), &dropped),
		); err != nil {
			t.Fatalf("%s drop dispatch: %v", drop.boardID, err)
		}

		if dropped != "dropped" {
			t.Fatalf("%s drop dispatch returned %q", drop.boardID, dropped)
		}

		var ok bool

		poll := kanbanColumnOrderExpr(drop.boardID, drop.column) + "===" + strconv.Quote(drop.want)
		if err := chromedp.Run(ctx, pollBool(poll, &ok)); err != nil || !ok {
			t.Fatalf("%s: card %s never landed at the top of %q (want order %q): err=%v",
				drop.boardID, drop.card, drop.column, drop.want, err)
		}
	}

	// Sanity: the source column lost the dragged card.
	var todo string

	if err := chromedp.Run(ctx,
		chromedp.Evaluate(kanbanColumnOrderExpr("kb-htmx", "todo"), &todo),
	); err != nil || strings.TrimSpace(todo) != "e1" {
		t.Fatalf("htmx todo column after drop = %q, want e1 (err %v)", todo, err)
	}
}

// kanbanButtonsOpacityExpr evaluates to the computed opacity of one card's
// hover-revealed move-button wrapper.
const kanbanButtonsOpacityExpr = `getComputedStyle(document.querySelector('#kb-htmx [data-tc-kanban-card="e1"] .tc-kanban-buttons')).opacity`

// TestKanbanE2ECoarsePointerButtonsVisible proves the touch fallback in a
// real browser: emulating a coarse pointer flips the wrapper's computed
// opacity from 0 to 1 via the unlayered @media (pointer: coarse) rule in
// templates/custom.css, so touch users always see the move buttons. The
// fine-pointer control leg asserts the hidden starting state first, so the
// flip can never pass vacuously.
//
// Emulation notes (verified empirically): Emulation.setEmulatedMedia IGNORES
// the pointer/hover features in this Chromium build — only
// Emulation.setTouchEmulationEnabled flips matchMedia('(pointer:coarse)').
// And because the wrapper carries utils.TransitionFast (150ms), the computed
// opacity must be POLLED to its settled value, not read immediately.
func TestKanbanE2ECoarsePointerButtonsVisible(t *testing.T) {
	srv := kanbanE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var ready bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollBool(kanbanE2EReady, &ready),
	); err != nil {
		t.Fatalf("navigate + readiness: %v", err)
	}

	var fine string

	if err := chromedp.Run(ctx, chromedp.Evaluate(kanbanButtonsOpacityExpr, &fine)); err != nil {
		t.Fatalf("fine-pointer opacity eval: %v", err)
	}

	if fine != "0" {
		t.Fatalf("fine pointer computed opacity = %q, want 0 (control leg)", fine)
	}

	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return emulation.SetTouchEmulationEnabled(true).WithMaxTouchPoints(5).Do(ctx)
	})); err != nil {
		t.Fatalf("emulate coarse pointer: %v", err)
	}

	var coarseMatches bool

	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`matchMedia('(pointer:coarse)').matches`, &coarseMatches),
	); err != nil || !coarseMatches {
		t.Fatalf("pointer:coarse media feature not emulated (matches=%v, err=%v)", coarseMatches, err)
	}

	var visible bool

	if err := chromedp.Run(ctx, pollBool(kanbanButtonsOpacityExpr+`==="1"`, &visible)); err != nil || !visible {
		t.Fatalf(
			"coarse pointer computed opacity never settled at 1 (visible=%v, err=%v) — touch fallback broken",
			visible,
			err,
		)
	}
}

// kanbanCrossBoardDropScript drags a card from one board and drops it on a
// DIFFERENT board's column — a move the component must ignore (each wired
// board owns its cards).
func kanbanCrossBoardDropScript(srcBoardID, cardID, dstBoardID, columnID string) string {
	return fmt.Sprintf(`(function(){
		var sb=document.getElementById(%q);
		var db=document.getElementById(%q);
		var dt=new DataTransfer();
		var card=sb.querySelector('[data-tc-kanban-card="%s"]');
		var zone=db.querySelector('[data-tc-kanban-column-body="%s"]');
		if(!card||!zone){return 'missing';}
		card.dispatchEvent(new DragEvent('dragstart',{bubbles:true,cancelable:true,dataTransfer:dt}));
		var over=zone.dispatchEvent(new DragEvent('dragover',{bubbles:true,cancelable:true,clientY:0,dataTransfer:dt}));
		zone.dispatchEvent(new DragEvent('drop',{bubbles:true,cancelable:true,dataTransfer:dt}));
		card.dispatchEvent(new DragEvent('dragend',{bubbles:true,cancelable:true,dataTransfer:dt}));
		return over?'accepted':'rejected';
	})()`, srcBoardID, dstBoardID, cardID, columnID)
}

// TestKanbanE2ECrossBoardDropIgnored proves the cross-board guard under a
// real drag sequence: dropping a card from the htmx board onto the Datastar
// board must not submit anything — both boards keep their exact card order
// and the dragover is not even accepted (no preventDefault).
func TestKanbanE2ECrossBoardDropIgnored(t *testing.T) {
	ctx := newKanbanReadyTab(t)

	var accepted string

	if err := chromedp.Run(ctx,
		chromedp.Evaluate(
			kanbanCrossBoardDropScript("kb-htmx", "e2", "kb-ds", "done"),
			&accepted,
		),
	); err != nil {
		t.Fatalf("cross-board drop dispatch: %v", err)
	}

	if accepted == "missing" {
		t.Fatal("cross-board drop script could not find the card or zone")
	}

	if accepted == "claimed" {
		t.Fatal("cross-board dragover was claimed (preventDefault ran) — guard missing from dragover listener")
	}

	// Both boards must still be in their initial order a moment later: no
	// request was submitted, no swap happened.
	time.Sleep(500 * time.Millisecond)

	for _, check := range []struct {
		boardID string
		column  string
		want    string
	}{
		{boardID: "kb-htmx", column: "todo", want: "e1,e2"},
		{boardID: "kb-htmx", column: "done", want: "e3"},
		{boardID: "kb-ds", column: "todo", want: "e1,e2"},
		{boardID: "kb-ds", column: "done", want: "e3"},
	} {
		var got string

		if err := chromedp.Run(ctx,
			chromedp.Evaluate(kanbanColumnOrderExpr(check.boardID, check.column), &got),
		); err != nil || got != check.want {
			t.Fatalf("%s %s after cross-board drop = %q, want %q (err %v)",
				check.boardID, check.column, got, check.want, err)
		}
	}
}

// TestKanbanE2EAnnouncesMove proves the post-swap live-region confirmation
// under both runtimes: after a keyboard move lands, the REPLACED board's
// live region reads "Moved <card> to <column>." — the poll-based mechanism
// from kanbanAnnounceJS, exercised through htmx's outerHTML swap AND the
// Datastar patch flow.
func TestKanbanE2EAnnouncesMove(t *testing.T) {
	ctx := newKanbanReadyTab(t)

	const want = "Moved First to In progress."

	for _, boardID := range []string{"kb-htmx", "kb-ds"} {
		kanbanClickMoveUntil(ctx, t,
			fmt.Sprintf(`#%s [data-tc-kanban-card="e1"] [data-tc-kanban-move="next"]`, boardID),
			kanbanColumnOrderExpr(boardID, "doing"), "e1",
		)

		var announced bool

		livePoll := fmt.Sprintf(
			`document.getElementById(%q).querySelector('[data-tc-kanban-live]').textContent===%q`,
			boardID, want,
		)

		if err := chromedp.Run(ctx, pollBool(livePoll, &announced)); err != nil || !announced {
			t.Fatalf("%s: post-swap live region never announced %q (announced=%v, err=%v)",
				boardID, want, announced, err)
		}
	}
}

// kanbanColumnCountExpr returns a JS expression evaluating to the number of
// cards in one column.
func kanbanColumnCountExpr(boardID, columnID string) string {
	return fmt.Sprintf(
		`document.querySelectorAll('#%s [data-tc-kanban-column-body="%s"] > [data-tc-kanban-card]').length`,
		boardID,
		columnID,
	)
}

// kanbanClickCountUntil clicks an always-visible affordance (the add/reset
// buttons live in the column/board header — no opacity-0 dance like the
// hover-revealed move buttons) and polls for the expected card count,
// retrying the click a few times: a click issued while the runtime
// re-attaches to swapped-in markup can be swallowed (the wire-forms-pack
// lesson).
func kanbanClickCountUntil(ctx context.Context, t *testing.T, buttonSel, countExpr string, want int) {
	t.Helper()

	poll := countExpr + "===" + strconv.Itoa(want)

	for range 3 {
		if err := chromedp.Run(ctx,
			chromedp.Click(buttonSel),
		); err != nil {
			t.Fatalf("click %s: %v", buttonSel, err)
		}

		var got bool

		err := chromedp.Run(ctx, pollBool(poll, &got))
		if err == nil && got {
			return
		}

		time.Sleep(300 * time.Millisecond)
	}

	t.Fatalf("column count never reached %d via %s", want, buttonSel)
}

// TestKanbanE2EAddAndResetBothTransports proves the KanbanColumn.Action slot
// end-to-end: clicking the add-card button appends a card to the target
// column through the transport dialect (htmx hx-post + outerHTML swap of
// the board root; Datastar @post + response-header targeting), and the
// reset button restores the empty column — both under the real runtimes.
// This is the "wired affordance ⇒ browser proof" doctrine test for Action:
// without it, a missing hx-target/hx-swap pair silently swaps the whole
// board into the button itself.
func TestKanbanE2EAddAndResetBothTransports(t *testing.T) {
	ctx := newKanbanReadyTab(t)

	// htmx board: one click on the empty "In progress" column's add button.
	kanbanClickCountUntil(ctx, t,
		`#kb-htmx [data-tc-kanban-column="doing"] [aria-label="Add card to In progress"]`,
		kanbanColumnCountExpr("kb-htmx", "doing"), 1)

	// Datastar board: the same affordance in the other dialect.
	kanbanClickCountUntil(ctx, t,
		`#kb-ds [data-tc-kanban-column="doing"] [aria-label="Add card to In progress"]`,
		kanbanColumnCountExpr("kb-ds", "doing"), 1)

	// Reset restores the starting layout on both boards.
	kanbanClickCountUntil(ctx, t,
		`#kb-htmx-wrap [aria-label="Reset demo board"]`,
		kanbanColumnCountExpr("kb-htmx", "doing"), 0)
	kanbanClickCountUntil(ctx, t,
		`#kb-ds-wrap [aria-label="Reset demo board"]`,
		kanbanColumnCountExpr("kb-ds", "doing"), 0)
}

// kanbanClickOnce clicks an affordance exactly once (the pending/revert tests
// need a single move in flight, unlike the retry-until helpers).
func kanbanClickOnce(ctx context.Context, t *testing.T, buttonSel string) {
	t.Helper()

	if err := chromedp.Run(ctx, chromedp.Click(buttonSel)); err != nil {
		t.Fatalf("click %s: %v", buttonSel, err)
	}
}

// kanbanFlakyPendingExpr returns the expression asserting e-slow sits in the
// target column ALREADY (optimistic placement), wearing the full pending
// register: tc-kanban-pending, aria-busy, and both count badges synced.
func kanbanFlakyPendingExpr(boardID string) string {
	return fmt.Sprintf(`(function(){
		var b=document.getElementById(%q);
		var card=b.querySelector('[data-tc-kanban-card="e-slow"]');
		if(!card)return false;
		return !!card.closest('[data-tc-kanban-column-body="doing"]')
			&& card.classList.contains('tc-kanban-pending')
			&& card.getAttribute('aria-busy')==='true'
			&& b.querySelector('[data-tc-kanban-column="doing"] [data-tc-kanban-count]').textContent==='1'
			&& b.querySelector('[data-tc-kanban-column="todo"] [data-tc-kanban-count]').textContent==='1';
	})()`, boardID)
}

// kanbanFlakyDoneExpr returns the expression asserting the move CONFIRMED:
// e-slow in place, pending register cleared, live region announcing the
// completed move.
func kanbanFlakyDoneExpr(boardID string) string {
	return fmt.Sprintf(`(function(){
		var b=document.getElementById(%q);
		var card=b.querySelector('[data-tc-kanban-card="e-slow"]');
		if(!card)return false;
		return !!card.closest('[data-tc-kanban-column-body="doing"]')
			&& !card.classList.contains('tc-kanban-pending')
			&& card.getAttribute('aria-busy')===null
			&& b.querySelector('[data-tc-kanban-live]').textContent==='Moved Slow move to In progress.';
	})()`, boardID)
}

// kanbanFlakyRevertedExpr returns the expression asserting the FAILED move
// was reverted honestly: e-fail back in its original column at the end,
// pending cleared, the tc-kanban-move-failed flash applied, counts restored,
// and the role="alert" region announcing the revert.
func kanbanFlakyRevertedExpr(boardID string) string {
	return fmt.Sprintf(`(function(){
		var b=document.getElementById(%q);
		var card=b.querySelector('[data-tc-kanban-card="e-fail"]');
		if(!card)return false;
		return !!card.closest('[data-tc-kanban-column-body="todo"]')
			&& !card.classList.contains('tc-kanban-pending')
			&& card.classList.contains('tc-kanban-move-failed')
			&& b.querySelector('[data-tc-kanban-column="todo"] [data-tc-kanban-count]').textContent==='2'
			&& b.querySelector('[data-tc-kanban-column="doing"] [data-tc-kanban-count]').textContent==='0'
			&& b.querySelector('[data-tc-kanban-alert]').textContent==='Moving Doomed move to In progress failed. The board was restored.';
	})()`, boardID)
}

// TestKanbanE2EPendingStateBothTransports proves the optimistic register
// (ADR-0041) at the browser level: the moment e-slow's move is submitted, the
// card sits in the target column wearing tc-kanban-pending + aria-busy with
// both count badges synced — long BEFORE the 1.2s-stalled response lands —
// and once the re-rendered board arrives the pending state is cleared and
// the live region confirms. Under htmx AND the real Datastar runtime.
func TestKanbanE2EPendingStateBothTransports(t *testing.T) {
	kanbanFlakyHTMXBoard.reset()
	kanbanFlakyDatastarBoard.reset()

	srv := kanbanE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	var ready bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollBool(kanbanE2EReady, &ready),
	); err != nil {
		t.Fatalf("navigate + readiness: %v", err)
	}

	for _, boardID := range []string{"kb-htmx-flaky", "kb-ds-flaky"} {
		buttonSel := fmt.Sprintf(`#%s [data-tc-kanban-card="e-slow"] [data-tc-kanban-move="next"]`, boardID)
		kanbanClickOnce(ctx, t, buttonSel)

		// The optimistic state must appear within 500ms — well inside the
		// 1.2s stall, so this poll can only pass BEFORE the response lands.
		var pending bool

		if err := chromedp.Run(ctx, pollBool(
			kanbanFlakyPendingExpr(boardID), &pending,
			chromedp.WithPollingTimeout(500*time.Millisecond), chromedp.WithPollingInterval(50*time.Millisecond),
		)); err != nil || !pending {
			t.Fatalf(
				"%s: optimistic pending state never appeared within 500ms (pending=%v, err=%v)",
				boardID,
				pending,
				err,
			)
		}

		// Then the response lands: pending cleared, confirmation announced.
		var done bool

		if err := chromedp.Run(ctx, pollBool(kanbanFlakyDoneExpr(boardID), &done)); err != nil || !done {
			t.Fatalf("%s: pending state never cleared after the swap (done=%v, err=%v)", boardID, done, err)
		}
	}
}

// TestKanbanE2EFailureRevertsBothTransports proves the honest failure path:
// e-fail's move 500s, the card snaps back to its original column with the
// tc-kanban-move-failed flash, counts are restored, and the role="alert"
// region announces the revert — under htmx AND the real Datastar runtime.
func TestKanbanE2EFailureRevertsBothTransports(t *testing.T) {
	kanbanFlakyHTMXBoard.reset()
	kanbanFlakyDatastarBoard.reset()

	srv := kanbanE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	var ready bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		pollBool(kanbanE2EReady, &ready),
	); err != nil {
		t.Fatalf("navigate + readiness: %v", err)
	}

	for _, boardID := range []string{"kb-htmx-flaky", "kb-ds-flaky"} {
		buttonSel := fmt.Sprintf(`#%s [data-tc-kanban-card="e-fail"] [data-tc-kanban-move="next"]`, boardID)
		kanbanClickOnce(ctx, t, buttonSel)

		var reverted bool

		if err := chromedp.Run(
			ctx,
			pollBool(kanbanFlakyRevertedExpr(boardID), &reverted),
		); err != nil ||
			!reverted {
			t.Fatalf("%s: failed move never reverted cleanly (reverted=%v, err=%v)", boardID, reverted, err)
		}

		// The failed-state flash self-clears after 4s.
		var flashGone bool

		flashExpr := fmt.Sprintf(
			`!document.querySelector('#%s [data-tc-kanban-card="e-fail"]').classList.contains('tc-kanban-move-failed')`,
			boardID,
		)

		if err := chromedp.Run(ctx, pollBool(flashExpr, &flashGone)); err != nil || !flashGone {
			t.Fatalf("%s: tc-kanban-move-failed flash never self-cleared (gone=%v, err=%v)", boardID, flashGone, err)
		}
	}
}

// kanbanParseMoveWithCSRF is the shared prelude of every kanban test move
// endpoint: parse the move form and enforce the test CSRF token. On failure
// it writes the error response and returns ok=false (backlog #262: one home
// for the prelude the e2e and pending-visual servers both need).
func kanbanParseMoveWithCSRF(w http.ResponseWriter, r *http.Request) (display.KanbanMove, bool) {
	move, err := display.ParseKanbanMove(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return display.KanbanMove{}, false
	}

	if r.FormValue(kanbanE2ECSRFField) != kanbanE2ECSRFToken {
		http.Error(w, "invalid CSRF token", http.StatusForbidden)

		return display.KanbanMove{}, false
	}

	return move, true
}

