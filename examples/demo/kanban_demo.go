// Kanban demo: server-owned board state + the move endpoint every consumer
// of display.KanbanBoard needs. Two boards share this file — one wired for
// htmx, one for Datastar — each with its own endpoint and state so the two
// runtimes never fight over one board.
package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"sync"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// kanbanCSRFToken is the demo's server-issued CSRF token. One token per
// process keeps the demo honest but simple; a real app issues one per
// session and validates it server-side on every mutating request.
var kanbanCSRFToken = newKanbanCSRFToken()

func newKanbanCSRFToken() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic("kanban demo: generate csrf token: " + err.Error())
	}

	return hex.EncodeToString(b[:])
}

// kanbanDemoState is the mutex-guarded board a demo endpoint owns. Real apps
// swap this for a database; the move semantics stay identical.
type kanbanDemoState struct {
	mu      sync.Mutex
	columns []display.KanbanColumn
	next    int64
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
func (s *kanbanDemoState) move(move display.KanbanMove) {
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
		return
	}

	for ci := range s.columns {
		if s.columns[ci].ID != move.Column {
			continue
		}

		if move.Index > len(s.columns[ci].Cards) {
			move.Index = len(s.columns[ci].Cards)
		}

		cards := make([]display.KanbanCard, 0, len(s.columns[ci].Cards)+1)
		cards = append(cards, s.columns[ci].Cards[:move.Index]...)
		cards = append(cards, card)
		cards = append(cards, s.columns[ci].Cards[move.Index:]...)
		s.columns[ci].Cards = cards

		return
	}
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
}

// kanbanDemoBoardProps builds the render props for one demo board: stable
// board ID (the response re-renders it), the transport's wire.Action, a
// per-column add-card button (posting to <move URL>/add/<column> in the
// SAME transport dialect), and a snapshot of the current state.
func (s *kanbanDemoState) kanbanDemoBoardProps(id string, action wire.Action) display.KanbanBoardProps {
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
	props.BaseProps = utils.BaseProps{ID: id}
	props.Columns = snapshot
	props.Wire = &action
	props.CSRFToken = kanbanCSRFToken

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
