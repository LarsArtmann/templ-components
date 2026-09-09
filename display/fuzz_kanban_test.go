package display

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// FuzzParseKanbanMove verifies the move-contract decoder never panics on
// arbitrary form bodies and, when it does succeed, echoes the submitted
// card/column verbatim (the consumer contract: IDs are opaque strings).
func FuzzParseKanbanMove(f *testing.F) {
	f.Add("c1", "todo", "0")
	f.Add("c1", "todo", "")
	f.Add("", "todo", "2")
	f.Add("c1", "", "-1")
	f.Add("c  ard", "col&=x", "9999999999999999999999")
	f.Add("c1", "todo", "abc")
	f.Add("%s", "%d", "\x00")
	f.Fuzz(func(t *testing.T, card, column, index string) {
		form := url.Values{}
		form.Set(kanbanFieldCard, card)
		form.Set(kanbanFieldColumn, column)
		form.Set(kanbanFieldIndex, index)

		r := httptest.NewRequest(http.MethodPost, "/api/kanban", strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		move, err := ParseKanbanMove(r)
		if err != nil {
			return
		}

		if move.Card != card || move.Column != column {
			t.Fatalf("ParseKanbanMove echoed (%q, %q), want (%q, %q)", move.Card, move.Column, card, column)
		}

		if index == "" && move.Index != 0 {
			t.Fatalf("ParseKanbanMove defaulted empty index to %d, want 0", move.Index)
		}
	})
}
