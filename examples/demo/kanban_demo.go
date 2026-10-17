// Kanban demo: server-owned board state + the move endpoint every consumer
// of display.KanbanBoard needs. Two boards share this file — one wired for
// htmx, one for Datastar — each with its own endpoint and state so the two
// runtimes never fight over one board.
//
// ANTI-DRIFT TIE (TODO #227): visualtest's kanbanE2EServer re-implements
// this same contract (CSRF-403 / same-origin-403 / unknown-column-404 /
// reset) because visualtest boots this demo as an external binary and cannot
// import it. The two implementations are pinned together by ONE probe table:
// runKanbanContractProbes in visualtest/demo_kanban_http_test.go runs
// TestDemoKanbanHTTPContracts (this binary) and
// TestKanbanE2EHTTPContractParity (the harness) through identical probes —
// change a status code or guard here, and the shared table names the side
// that drifted. Keep both green.
package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// errKanbanStaleView is the sorted-view rejection: the submitted move
// references state the server does not have (unknown card or column).
var errKanbanStaleView = errors.New("stale board view: card or column not found")

// kanbanDemoMoveDelay simulates real network latency on the move endpoints.
// The board applies moves optimistically (ADR-0041): the card moves the
// moment it is dropped and wears a spinner until the server's re-render
// lands — localhost would make that state invisible, so the demo delays on
// purpose to show what consumers' users actually experience.
const kanbanDemoMoveDelay = 800 * time.Millisecond

// kanbanDemoState is the mutex-guarded board a demo endpoint owns. Real apps
// swap this for a database; the move semantics stay identical.
type kanbanDemoState struct {
	mu      sync.Mutex
	columns []display.KanbanColumn
	next    int64
	// path, when set, file-backs the state: loaded at wiring time, saved
	// atomically after every mutation (backlog #189). Empty means
	// memory-only, which keeps tests and route goldens deterministic.
	path string
}

// The two demo boards, one per transport, each with independent state so
// moves on one never mutate the other mid-render.
var (
	kanbanHTMXState     = &kanbanDemoState{columns: kanbanDemoInitialColumns()}
	kanbanDatastarState = &kanbanDemoState{columns: kanbanDemoInitialColumns()}
)

// kanbanDemoInitialColumns returns the demo board's starting layout.
func kanbanDemoInitialColumns() []display.KanbanColumn {
	return []display.KanbanColumn{
		{
			ID:    "backlog",
			Title: "Backlog",
			Tone:  display.KanbanToneGray,
			Cards: []display.KanbanCard{
				{ID: "kb-1", Title: "Audit CSP headers"},
				{ID: "kb-2", Title: "Dark mode pass"},
				{ID: "kb-3", Title: "Write migration guide"},
			},
		},
		{
			ID:    "progress",
			Title: "In progress",
			Tone:  display.KanbanToneBlue,
			Cards: []display.KanbanCard{},
		},
		{
			ID:    "review",
			Title: "In review",
			Tone:  display.KanbanToneYellow,
			Cards: []display.KanbanCard{},
		},
		{
			ID:    "done",
			Title: "Done",
			Tone:  display.KanbanToneGreen,
			Cards: []display.KanbanCard{{ID: "kb-4", Title: "Set up CI"}},
		},
	}
}

// move applies a decoded KanbanMove: remove the card from wherever it lives,
// insert it at the requested position (clamped to the column bounds).
// Unknown cards or columns are no-ops — the board is the source of truth.
// moveErr applies one move and reports the sorted-view rejections
// (backlog #224): a card the server state does not know, or a target
// column the board does not have, means the submitter operated on a stale
// view — the endpoint answers 422 instead of silently no-oping. The old
// behavior not only hid the stale view, it DELETED the card from its
// source column when the target column was unknown. Index overflow stays
// a clamp (append) — a legitimate race the board self-heals.
func (s *kanbanDemoState) moveErr(move display.KanbanMove) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var card display.KanbanCard
	found := false

	for ci := range s.columns {
		for pi := range s.columns[ci].Cards {
			if s.columns[ci].Cards[pi].ID == move.Card {
				card = s.columns[ci].Cards[pi]
				s.columns[ci].Cards = append(s.columns[ci].Cards[:pi], s.columns[ci].Cards[pi+1:]...)
				found = true

				break
			}
		}
	}

	if !found {
		return errKanbanStaleView
	}

	target := -1

	for ci := range s.columns {
		if s.columns[ci].ID == move.Column {
			target = ci

			break
		}
	}

	if target < 0 {
		return errKanbanStaleView
	}

	if move.Index > len(s.columns[target].Cards) {
		move.Index = len(s.columns[target].Cards)
	}

	cards := make([]display.KanbanCard, 0, len(s.columns[target].Cards)+1)
	cards = append(cards, s.columns[target].Cards[:move.Index]...)
	cards = append(cards, card)
	cards = append(cards, s.columns[target].Cards[move.Index:]...)
	s.columns[target].Cards = cards
	s.persistLocked()

	return nil
}

// move keeps the silent legacy shape for in-memory tests that drive moves
// directly and do not care about rejections.
func (s *kanbanDemoState) move(move display.KanbanMove) {
	_ = s.moveErr(move)
}

// add appends a fresh demo card to the named column and reports whether
// the column exists — the handler turns false into a 404 so a broken link
// surfaces instead of silently no-oping.
func (s *kanbanDemoState) add(columnID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ci := range s.columns {
		if s.columns[ci].ID != columnID {
			continue
		}

		n := s.next
		s.next++
		s.columns[ci].Cards = append(s.columns[ci].Cards, display.KanbanCard{
			ID:    "kb-new-" + strconv.FormatInt(n, 10),
			Title: "New card " + strconv.FormatInt(n+1, 10),
		})
		s.persistLocked()

		return true
	}

	return false
}

// reset restores the demo board's starting layout. The new-card counter
// keeps counting so a stale open tab can never resubmit an id a fresh
// card already owns.
func (s *kanbanDemoState) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.columns = kanbanDemoInitialColumns()
	s.persistLocked()
}

// persistLocked writes the state to the backing file (atomic temp+rename).
// Caller must hold s.mu. No-op for memory-only boards.
func (s *kanbanDemoState) persistLocked() {
	if s.path == "" {
		return
	}

	blob, err := json.Marshal(kanbanPersistedState{Columns: s.columns, Next: s.next})
	if err != nil {
		return
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, blob, 0o600); err != nil {
		return
	}
	_ = os.Rename(tmp, s.path)
}

// loadKanbanStateFile opts one board into file-backed state (backlog
// #189): the KANBAN_STATE_FILE env var names the backing file shared by
// both demo boards. Called at wiring time, before the server accepts
// requests. A missing file starts from the initial layout; a corrupt one
// is ignored the same way (the demo holds no real data).
func loadKanbanStateFile(s *kanbanDemoState, path string) {
	s.path = path

	blob, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var persisted kanbanPersistedState
	if err := json.Unmarshal(blob, &persisted); err != nil || len(persisted.Columns) == 0 {
		return
	}

	s.columns = persisted.Columns
	s.next = persisted.Next
}

// kanbanPersistedState is the file shape — columns plus the id counter, so
// reloaded adds cannot collide with pre-restart ids.
type kanbanPersistedState struct {
	Columns []display.KanbanColumn `json:"columns"`
	Next    int64                  `json:"next"`
}

// kanbanDemoBoardProps builds the render props for one demo board: stable
// board ID (the response re-renders it), the transport's wire.Action, a
// per-column add-card button (posting to <move URL>/add/<column> in the
// SAME transport dialect), and a snapshot of the current state.
func (s *kanbanDemoState) kanbanDemoBoardProps(id string, action wire.Action, csrf string) display.KanbanBoardProps {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshot := make([]display.KanbanColumn, len(s.columns))
	for ci, col := range s.columns {
		add := action
		add.URL += "/add/" + col.ID
		add.Method = wire.MethodPost
		add.Event = wire.EventClick
		add.Target = "#" + id
		snapshot[ci] = display.KanbanColumn{
			ID:     col.ID,
			Title:  col.Title,
			Tone:   col.Tone,
			Cards:  append([]display.KanbanCard(nil), col.Cards...),
			Action: kanbanAddButton(add, col.Title),
		}
	}

	props := display.DefaultKanbanBoardProps()
	boardLabel := "Kanban board (HTMX)"
	if action.Transport == wire.TransportDatastar {
		boardLabel = "Kanban board (Datastar)"
	}
	props.BaseProps = utils.BaseProps{ID: id, Nonce: demoNonceConst, AriaLabel: boardLabel}
	props.Columns = snapshot
	props.Wire = &action
	props.CSRFToken = csrf

	return props
}

// kanbanAddButton builds the column header's add-card affordance. The Wire
// attributes carry the transport dialect (hx-post vs data-on:click), so the
// same builder serves both boards. htmx targets the board root (Target is
// set by the caller) and must swap outerHTML — otherwise the whole-board
// response lands inside this button.
func kanbanAddButton(action wire.Action, columnTitle string) templ.Component {
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

// kanbanResetButton builds the per-board reset affordance in the board's
// dialect. htmx swaps the board root via Target + hx-swap outerHTML (the
// move form's response contract); Datastar targeting is response-driven —
// the endpoint is wrapped in wire.Handler — so no target is carried.
func kanbanResetButton(action wire.Action) templ.Component {
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

// kanbanCSRFFieldName matches KanbanBoardProps's default hidden-input name.
const kanbanCSRFFieldName = "csrf_token"

// kanbanSameOriginRequest reports whether r comes from the demo's own
// origin. The move form carries the CSRF token in a hidden input; the
// bodyless add/reset POSTs have no form to carry one, so they enforce
// same-origin instead — a token in the URL would leak into logs.
func kanbanSameOriginRequest(r *http.Request) bool {
	if r.Header.Get("Sec-Fetch-Site") == "same-origin" {
		return true
	}

	origin := r.Header.Get("Origin")

	return origin == "http://"+r.Host || origin == "https://"+r.Host
}
