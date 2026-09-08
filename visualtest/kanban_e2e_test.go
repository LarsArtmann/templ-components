package visualtest

import (
	"context"
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

// kanbanE2EBoard is the mutex-guarded server state one dialect endpoint owns.
type kanbanE2EBoard struct {
	mu      sync.Mutex
	columns []display.KanbanColumn
}

func newKanbanE2EBoard() *kanbanE2EBoard {
	b := &kanbanE2EBoard{}
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

// reset restores the initial layout (tests share the package-global boards
// and must not observe each other's moves).
func (b *kanbanE2EBoard) reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.columns = initialKanbanE2EColumns()
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

		htmxProps := kanbanHTMXBoard.kanbanE2EBoardProps("kb-htmx", htmxAction)
		if err := display.KanbanBoard(htmxProps).Render(ctx, w); err != nil {
			return err
		}

		datastarProps := kanbanDatastarBoard.kanbanE2EBoardProps("kb-ds", datastarAction)

		return display.KanbanBoard(datastarProps).Render(ctx, w)
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

			move, err := display.ParseKanbanMove(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

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
)

// kanbanE2EBoardProps snapshots the state into render props with a stable
// board id and the transport's wire action.
func (b *kanbanE2EBoard) kanbanE2EBoardProps(id string, action wire.Action) display.KanbanBoardProps {
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

	return props
}

// kanbanColumnOrderExpr returns a JS expression evaluating to the ordered
// comma-separated card ids of one column, e.g. "e1,e2".
func kanbanColumnOrderExpr(boardID, columnID string) string {
	return fmt.Sprintf(
		`Array.from(document.querySelectorAll('#%s [data-tc-kanban-column-body="%s"] > [data-tc-kanban-card]')).map(function(el){return el.getAttribute('data-tc-kanban-card');}).join(',')`,
		boardID, columnID,
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

	for attempt := 0; attempt < 3; attempt++ {
		if err := chromedp.Run(ctx,
			chromedp.Click(buttonSel, chromedp.NodeVisible),
		); err != nil {
			t.Fatalf("click %s: %v", buttonSel, err)
		}

		var got string

		err := chromedp.Run(ctx, chromedp.Poll(poll, &got))
		if err == nil && got == "true" {
			return
		}

		time.Sleep(300 * time.Millisecond)
	}

	t.Fatalf("card never reached the expected column order %q via %s", want, buttonSel)
}

// TestKanbanE2EKeyboardMovesBothTransports proves the accessibility path:
// clicking a card's "move to next column" button submits the hidden form
// and the re-rendered board carries the card in the target column — under
// htmx AND the real Datastar runtime.
func TestKanbanE2EKeyboardMovesBothTransports(t *testing.T) {
	kanbanHTMXBoard.reset()
	kanbanDatastarBoard.reset()

	srv := kanbanE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 45*time.Second)
	defer cancelTimeout()

	var ready bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(kanbanE2EReady, &ready),
	); err != nil {
		t.Fatalf("navigate + readiness: %v", err)
	}

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
	kanbanHTMXBoard.reset()
	kanbanDatastarBoard.reset()

	srv := kanbanE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 45*time.Second)
	defer cancelTimeout()

	var ready string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(kanbanE2EReady, &ready),
	); err != nil {
		t.Fatalf("navigate + readiness: %v", err)
	}

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

		var ok string

		poll := kanbanColumnOrderExpr(drop.boardID, drop.column) + "===" + fmt.Sprintf("%q", drop.want)
		if err := chromedp.Run(ctx, chromedp.Poll(poll, &ok)); err != nil || ok != "true" {
			t.Fatalf("%s: card %s never landed at the top of %q (want order %q): err=%v got=%q",
				drop.boardID, drop.card, drop.column, drop.want, err, ok)
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
