package display

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// kanbanTestColumns is a small realistic board shared by the kanban tests.
func kanbanTestColumns() []KanbanColumn {
	return []KanbanColumn{
		{
			ID:    "todo",
			Title: "To do",
			Cards: []KanbanCard{
				{ID: "c1", Title: "Write docs"},
				{ID: "c2", Title: "Cut release", Href: "/issues/2"},
			},
		},
		{ID: "doing", Title: "In progress", Cards: []KanbanCard{{ID: "c3", Title: "Kanban board"}}},
		{ID: "done", Title: "Done"},
	}
}

// TestKanbanBoardWire verifies the hidden move form renders the right
// attribute dialect per transport, with the board's documented defaults
// applied (POST, submit event, form encoding, self-targeting outerHTML swap
// under htmx) and stays inert without wiring.
func TestKanbanBoardWire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		props       KanbanBoardProps
		contains    []string
		notContains []string
	}{
		{
			name: "nil wire renders a read-only board",
			props: KanbanBoardProps{
				BaseProps: utils.BaseProps{ID: "kb-ro"},
				Columns:   kanbanTestColumns(),
			},
			notContains: []string{
				"data-tc-kanban-form", "hx-post", "data-on:", "draggable",
				"data-tc-kanban-move", "tcKanbanAttached",
			},
		},
		{
			name: "empty wire URL is inert",
			props: KanbanBoardProps{
				BaseProps: utils.BaseProps{ID: "kb-inert"},
				Columns:   kanbanTestColumns(),
				Wire:      &wire.Action{URL: ""},
			},
			notContains: []string{"data-tc-kanban-form", "hx-post", "data-on:"},
		},
		{
			name: "htmx dialect applies move defaults",
			props: KanbanBoardProps{
				BaseProps: utils.BaseProps{ID: "kb-htmx"},
				Columns:   kanbanTestColumns(),
				Wire:      &wire.Action{URL: "/api/kanban/move"},
			},
			contains: []string{
				`hx-post="/api/kanban/move"`,
				`hx-trigger="submit"`,
				`hx-target="#kb-htmx"`,
				`hx-swap="outerHTML"`,
				`data-tc-kanban-form`,
				`name="card"`,
				`name="column"`,
				`name="index"`,
			},
		},
		{
			name: "datastar dialect applies move defaults",
			props: KanbanBoardProps{
				BaseProps: utils.BaseProps{ID: "kb-ds"},
				Columns:   kanbanTestColumns(),
				Wire:      &wire.Action{Transport: wire.TransportDatastar, URL: "/api/kanban/move"},
			},
			// templ escapes the apostrophes inside the attribute value.
			contains: []string{
				`data-on:submit="@post(&#39;/api/kanban/move&#39;, {contentType: &#39;form&#39;})"`,
			},
			notContains: []string{"hx-post", "hx-target", "hx-swap"},
		},
		{
			name: "consumer overrides win over defaults",
			props: KanbanBoardProps{
				BaseProps: utils.BaseProps{ID: "kb-ov"},
				Columns:   kanbanTestColumns(),
				Wire: &wire.Action{
					Method:      wire.MethodPatch,
					URL:         "/api/kanban/move",
					Target:      "#elsewhere",
					ContentType: wire.ContentTypeJSON,
				},
			},
			contains: []string{
				`hx-patch="/api/kanban/move"`,
				`hx-target="#elsewhere"`,
			},
			notContains: []string{`{contentType: 'form'}`},
		},
		{
			name: "csrf token travels in the move form",
			props: KanbanBoardProps{
				BaseProps: utils.BaseProps{ID: "kb-csrf"},
				Columns:   kanbanTestColumns(),
				Wire:      &wire.Action{URL: "/api/kanban/move"},
				CSRFToken: "tok-123",
			},
			contains: []string{
				`name="csrf_token"`,
				`value="tok-123"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			html := utils.Render(t, KanbanBoard(tt.props))
			for _, want := range tt.contains {
				utils.AssertContains(t, html, want)
			}

			for _, notWant := range tt.notContains {
				utils.AssertNotContains(t, html, notWant)
			}
		})
	}
}

// TestKanbanBoardWireNonMutation verifies the wire helper never mutates the
// consumer's action (the defaults are applied to a copy).
func TestKanbanBoardWireNonMutation(t *testing.T) {
	t.Parallel()

	action := &wire.Action{URL: "/api/kanban/move"}
	props := KanbanBoardProps{Columns: kanbanTestColumns(), Wire: action}
	_ = utils.Render(t, KanbanBoard(props))

	if action.Method != wire.MethodUnspecified {
		t.Errorf("consumer action mutated: Method = %q", action.Method)
	}

	if action.Event != wire.EventUnspecified {
		t.Errorf("consumer action mutated: Event = %q", action.Event)
	}

	if action.ContentType != wire.ContentTypeUnspecified {
		t.Errorf("consumer action mutated: ContentType = %q", action.ContentType)
	}

	if action.Target != "" {
		t.Errorf("consumer action mutated: Target = %q", action.Target)
	}
}

// TestParseKanbanMove pins the server-side half of the move contract.
func TestParseKanbanMove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    string
		want    KanbanMove
		wantErr bool
	}{
		{
			name: "full move",
			body: "card=c1&column=done&index=2",
			want: KanbanMove{Card: "c1", Column: "done", Index: 2},
		},
		{
			name: "missing index defaults to top",
			body: "card=c1&column=todo",
			want: KanbanMove{Card: "c1", Column: "todo", Index: 0},
		},
		{
			name:    "empty index defaults to top",
			body:    "card=c1&column=todo&index=",
			want:    KanbanMove{Card: "c1", Column: "todo", Index: 0},
			wantErr: false,
		},
		{name: "missing card", body: "column=todo&index=0", wantErr: true},
		{name: "missing column", body: "card=c1&index=0", wantErr: true},
		{name: "non-numeric index", body: "card=c1&column=todo&index=abc", wantErr: true},
		{name: "negative index", body: "card=c1&column=todo&index=-1", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequestWithContext(
				t.Context(),
				http.MethodPost,
				"/api/kanban/move",
				strings.NewReader(tt.body),
			)
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			got, err := ParseKanbanMove(r)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseKanbanMove(%q) = %+v, want error", tt.body, got)
				}

				return
			}

			if err != nil {
				t.Fatalf("ParseKanbanMove(%q) error: %v", tt.body, err)
			}

			if got != tt.want {
				t.Errorf("ParseKanbanMove(%q) = %+v, want %+v", tt.body, got, tt.want)
			}
		})
	}
}

// TestKanbanHelpers covers the small private helpers.
func TestKanbanHelpers(t *testing.T) {
	t.Parallel()

	if got := kanbanColumnID("kb", 0, KanbanColumn{}); got != "kb-col-0" {
		t.Errorf("kanbanColumnID without column ID = %q", got)
	}

	if got := kanbanColumnID("kb", 1, KanbanColumn{ID: "in progress"}); got != "kb-col-in-progress" {
		t.Errorf("kanbanColumnID with spaces = %q", got)
	}

	if got := kanbanCountLabel(0); got != "no cards" {
		t.Errorf("kanbanCountLabel(0) = %q", got)
	}

	if got := kanbanCountLabel(1); got != "1 card" {
		t.Errorf("kanbanCountLabel(1) = %q", got)
	}

	if got := kanbanCountLabel(3); got != "3 cards" {
		t.Errorf("kanbanCountLabel(3) = %q", got)
	}
}

// TestKanbanJSSingletonGuard verifies the script is a guarded singleton so
// re-renders (the board swaps itself after every move) stay idempotent.
func TestKanbanJSSingletonGuard(t *testing.T) {
	t.Parallel()

	js := kanbanJS()
	if !strings.Contains(js, "if(!window.tcKanbanAttached){window.tcKanbanAttached=true;") {
		t.Error("kanban script lacks the tcKanbanAttached singleton guard")
	}
}

// TestKanbanJSCrossBoardGuard pins the two guards that make each wired board
// ignore drops of cards owned by a DIFFERENT board: dragover must not claim
// the drop (no preventDefault) and drop must not submit. Behavior is proven
// browser-level by TestKanbanE2ECrossBoardDropIgnored; this catches silent
// JS regressions in unit time.
func TestKanbanJSCrossBoardGuard(t *testing.T) {
	t.Parallel()

	js := kanbanJS()

	for _, token := range []string{
		"if(tcKbBoard(tcKbSrc)!==b)return;",
		"if(tcKbBoard(src)!==b)return;",
	} {
		if !strings.Contains(js, token) {
			t.Errorf("kanban script lacks cross-board guard %q", token)
		}
	}

	dragOver := js[strings.Index(js, "dragover"):strings.Index(js, "drop',")]
	if dropGuard := strings.Index(dragOver, "if(tcKbBoard(tcKbSrc)!==b)return;"); dropGuard == -1 {
		t.Error("cross-board guard missing from the dragover listener")
	} else if prevent := strings.Index(dragOver, "e.preventDefault();"); prevent < dropGuard {
		t.Error("dragover preventDefault must come AFTER the cross-board guard")
	}
}

// TestKanbanJSPostSwapAnnouncement pins the post-swap confirmation pipeline:
// every submit stages a "Moved ..." message and arms a swap-agnostic poll
// that writes it into the REPLACED board's live region (the old region dies
// with the outerHTML swap). Browser-level proof: TestKanbanE2EAnnouncesMove.
func TestKanbanJSPostSwapAnnouncement(t *testing.T) {
	t.Parallel()

	js := kanbanJS()

	for _, token := range []string{
		"tcKbAnnounce='Moved '+title+' to '+(cn||colId)+'.';",
		"tcKbAnnounceAfter(b);",
		"function tcKbAnnounceIn(b){",
		"function tcKbAnnounceAfter(old){",
		"var live=b.querySelector('[data-tc-kanban-live]');",
		"if(!old.isConnected&&nb&&nb!==old){clearInterval(timer);tcKbAnnounceIn(nb);}",
	} {
		if !strings.Contains(js, token) {
			t.Errorf("kanban script lacks post-swap announcement token %q", token)
		}
	}
}
