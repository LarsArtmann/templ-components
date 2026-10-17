package visualtest

import (
	"net/http"
	"time"
)

// The flaky-interaction contract for optimistic-UI e2e proofs (ADR-0041).
//
// A stateful component's test board arms two special interactions:
//
//   - FlakySlowID: the server stalls before applying the move, opening a
//     deterministic pending window in which the optimistic register is
//     observable (FlakySlowDelay in the e2e suite; stall-until-disconnect
//     in the freeze-frame capture suite).
//   - FlakyFailID: the server always answers 500, exercising the honest
//     revert path (card returns, pending register clears, failure flash).
//
// Server sites share FlakyFail for the doomed branch; each site owns its
// stall shape. Client sites assert on cards carrying these IDs via the
// data-tc-kanban-card hooks.
const (
	FlakySlowID = "e-slow"
	FlakyFailID = "e-fail"

	// FlakySlowDelay is the e2e pending window — long enough that the
	// optimistic pending state is deterministically observable before the
	// response lands, short enough not to dominate the suite.
	FlakySlowDelay = 1200 * time.Millisecond
)

// FlakyFail writes the doomed interaction's 500 response and reports
// whether the caller must stop (id is FlakyFailID); the revert-path branch
// shared by every flaky-board server.
func FlakyFail(w http.ResponseWriter, id string) bool {
	if id != FlakyFailID {
		return false
	}

	http.Error(w, "boom", http.StatusInternalServerError)

	return true
}
