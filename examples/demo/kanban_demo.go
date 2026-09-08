// Kanban demo: server-owned board state + the move endpoint every consumer
// of display.KanbanBoard needs. Two boards share this file — one wired for
// htmx, one for Datastar — each with its own endpoint and state so the two
// runtimes never fight over one board.
package main

import (
	"sync"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// kanbanDemoState is the mutex-guarded board a demo endpoint owns. Real apps
// swap this for a database; the move semantics stay identical.
type kanbanDemoState struct {
	mu      sync.Mutex
	columns []display.KanbanColumn
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
			Cards: []display.KanbanCard{
				{ID: "kb-1", Title: "Audit CSP headers"},
				{ID: "kb-2", Title: "Dark mode pass"},
				{ID: "kb-3", Title: "Write migration guide"},
			},
		},
		{
			ID:    "progress",
			Title: "In progress",
			Cards: []display.KanbanCard{},
		},
		{
			ID:    "review",
			Title: "In review",
			Cards: []display.KanbanCard{},
		},
		{
			ID:    "done",
			Title: "Done",
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

// kanbanDemoBoardProps builds the render props for one demo board: stable
// board ID (the response re-renders it), the transport's wire.Action, and a
// snapshot of the current state.
func (s *kanbanDemoState) kanbanDemoBoardProps(id string, action wire.Action) display.KanbanBoardProps {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshot := make([]display.KanbanColumn, len(s.columns))
	for ci, col := range s.columns {
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
