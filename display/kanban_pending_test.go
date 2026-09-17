package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// kanbanWiredHTML renders the standard wired test board once for the
// pending-register assertions.
func kanbanWiredHTML(t *testing.T) string {
	t.Helper()

	return utils.Render(t, KanbanBoard(KanbanBoardProps{
		BaseProps: utils.BaseProps{ID: "kb-pending"},
		Columns:   kanbanTestColumns(),
		Wire:      &wire.Action{URL: "/api/kanban/move"},
	}))
}

// TestKanbanPendingRegisterMarkup pins the markup hooks the optimistic-move
// pipeline (ADR-0041) drives at runtime: the count badge and empty-column
// placeholder hooks (present on every board, wired or not) and the wired-only
// sr-only role="alert" region the failure announcement writes into.
func TestKanbanPendingRegisterMarkup(t *testing.T) {
	t.Parallel()

	wired := kanbanWiredHTML(t)

	for _, token := range []string{
		`data-tc-kanban-count`,
		`data-tc-kanban-empty`,
		`data-tc-kanban-alert`,
		`role="alert"`,
		`data-tc-kanban-live`,
	} {
		if !strings.Contains(wired, token) {
			t.Errorf("wired board lacks pending-register hook %q", token)
		}
	}

	readonly := utils.Render(t, KanbanBoard(KanbanBoardProps{
		BaseProps: utils.BaseProps{ID: "kb-ro"},
		Columns:   kanbanTestColumns(),
	}))

	for _, token := range []string{
		`data-tc-kanban-count`,
		`data-tc-kanban-empty`,
	} {
		if !strings.Contains(readonly, token) {
			t.Errorf("read-only board lacks inert hook %q", token)
		}
	}

	for _, token := range []string{
		`data-tc-kanban-alert`,
		`data-tc-kanban-live`,
		`data-tc-kanban-form`,
	} {
		if strings.Contains(readonly, token) {
			t.Errorf("read-only board must not render %q", token)
		}
	}
}

// TestKanbanJSOptimisticPending pins the optimistic-move pipeline in the
// singleton script: instant DOM placement, the pending register, the three
// success clearers, and the failure revert. Browser-level proof:
// TestKanbanE2EPendingStateBothTransports and
// TestKanbanE2EFailureRevertsBothTransports.
func TestKanbanJSOptimisticPending(t *testing.T) {
	t.Parallel()

	js := kanbanJS()

	for _, token := range []string{
		// Optimistic placement before the request fires.
		"if(card&&c)tcKbOptimistic(b,card,c,index,title,cn||colId);",
		"function tcKbOptimistic(b,card,zone,idx,title,dest){",
		"function tcKbPlace(zone,card,idx){",
		"function tcKbSyncCounts(zone){",
		// The pending register itself.
		"card.classList.add('tc-kanban-pending');",
		"card.setAttribute('aria-busy','true');",
		"if(ph)ph.hidden=true;",
		"var tcKbPending=[];",
		// Success: htmx afterRequest success clears; Datastar finished clears;
		// the swap poll clears (kanbanAnnounceJS, pinned there).
		"if(d.successful===true){tcKbSucceed(b.id);}",
		"function tcKbSucceed(bid){",
		"d.type==='error'||d.type==='retries-failed'",
		"else if(d.type==='finished'){tcKbSucceed(b.id);}",
		// Failure: revert with original-position restore + visible flash +
		// role="alert" announcement.
		"function tcKbRevert(bid){",
		"if(e.next&&e.next.parentNode===e.parent){e.parent.insertBefore(e.card,e.next);}",
		"e.card.classList.add('tc-kanban-move-failed');",
		"if(alertEl)alertEl.textContent=e.msg+' failed. The board was restored.';",
		"tcKbAnnounce=null;",
		// Listener guards: only the kanban move form may drive the pipeline.
		"!f.hasAttribute('data-tc-kanban-form')",
		"document.addEventListener('htmx:afterRequest'",
		"document.addEventListener('htmx:responseError'",
		"document.addEventListener('htmx:sendError'",
		"document.addEventListener('datastar-fetch'",
	} {
		if !strings.Contains(js, token) {
			t.Errorf("kanban script lacks optimistic-pending token %q", token)
		}
	}

	// The htmx listener must decide success vs failure on detail.successful
	// AFTER resolving the board, and the Datastar listener must treat error
	// states before the finished no-op.
	at := func(needle string) int {
		idx := strings.Index(js, needle)
		if idx == -1 {
			t.Fatalf("kanban script lacks %q", needle)
		}

		return idx
	}

	htmx := js[at("htmx:afterRequest"):at("htmx:responseError")]
	if success := strings.Index(htmx, "d.successful===true"); success == -1 {
		t.Error("htmx afterRequest listener lacks the success branch")
	} else if fail := strings.Index(htmx, "tcKbFailFrom(f);"); fail != -1 && fail < success {
		t.Error("htmx afterRequest listener must check success BEFORE failing")
	}

	datastar := js[at("datastar-fetch"):]
	if fail := strings.Index(datastar, "tcKbFailFrom(f);"); fail == -1 {
		t.Error("datastar listener lacks the failure branch")
	} else if done := strings.Index(datastar, "tcKbSucceed(b.id);"); done != -1 && done < fail {
		t.Error("datastar listener must handle error/retries-failed BEFORE finished")
	}

	// The revert origin (parent + next sibling) must be captured BEFORE the
	// card is placed, or a failed move would restore to the wrong position.
	opt := js[at("function tcKbOptimistic("):at("function tcKbSucceed(")]
	if origin := strings.Index(opt, "var next=card.nextSibling;"); origin == -1 {
		t.Error("tcKbOptimistic lacks the origin capture")
	} else if place := strings.Index(opt, "tcKbPlace(zone,card,idx);"); place < origin {
		t.Error("tcKbOptimistic must capture the origin BEFORE placing the card")
	}
}
