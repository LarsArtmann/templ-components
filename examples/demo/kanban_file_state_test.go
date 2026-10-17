package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/larsartmann/templ-components/display"
)

// TestKanbanFileBackedState pins the #189 persistence contract: an
// opt-in file-backed board loads a previously persisted layout across
// instances (the demo restart case), mutations persist, and a corrupt or
// missing file degrades to the initial layout instead of failing the
// server. Memory-only (empty path) stays the default so goldens stay
// deterministic.
func TestKanbanFileBackedState(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "kanban.json")

	board := &kanbanDemoState{columns: kanbanDemoInitialColumns()}
	loadKanbanStateFile(board, path)

	if got := board.moveErr(display.KanbanMove{Card: "kb-1", Column: "done", Index: 0}); got != nil {
		t.Fatalf("move on fresh board: %v", got)
	}

	board.add("backlog")

	if got := board.moveErr(display.KanbanMove{Card: "kb-1", Column: "nope", Index: 0}); got != errKanbanStaleView {
		t.Errorf("move to unknown column = %v, want errKanbanStaleView", got)
	}

	reloaded := &kanbanDemoState{columns: kanbanDemoInitialColumns()}
	loadKanbanStateFile(reloaded, path)

	if got := reloaded.moveErr(display.KanbanMove{Card: "kb-1", Column: "backlog", Index: 0}); got != nil {
		t.Errorf("reloaded board lost the persisted move (kb-1 not found): %v", got)
	}

	for _, column := range reloaded.columns {
		if column.ID != "backlog" {
			continue
		}

		if len(column.Cards) == 0 || column.Cards[len(column.Cards)-1].ID != "kb-new-0" {
			t.Errorf("reloaded board lost the persisted add; backlog tail = %+v", column.Cards)
		}
	}

	// A corrupt file must not take the server down.
	if err := writeFile(path, []byte("{not json")); err != nil {
		t.Fatalf("seed corrupt file: %v", err)
	}

	corrupt := &kanbanDemoState{columns: kanbanDemoInitialColumns()}
	loadKanbanStateFile(corrupt, path)

	if len(corrupt.columns) == 0 {
		t.Error("corrupt file must degrade to the initial layout, not empty columns")
	}
}

// writeFile is os.WriteFile with a test-appropriate signature.
func writeFile(path string, blob []byte) error {
	return os.WriteFile(path, blob, 0o600)
}
