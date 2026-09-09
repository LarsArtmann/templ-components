package display

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils/wire"
)

// Field names of the hidden inputs the kanban move form submits. They are the
// server-side contract of a wired KanbanBoard: on every move (drag-and-drop
// or keyboard button) the component POSTs these three values to Wire.URL.
// Use the constants (or ParseKanbanMove) in the handler so both sides stay
// pinned to one spelling.
const (
	kanbanFieldCard   = "card"
	kanbanFieldColumn = "column"
	kanbanFieldIndex  = "index"
)

// kanbanColumnID derives the DOM id of a column's heading (the
// aria-labelledby target). Consumer column IDs win; the index keeps the id
// unique for columns without one. Whitespace is collapsed so the value
// stays a valid HTML id.
func kanbanColumnID(boardID string, index int, col KanbanColumn) string {
	if col.ID == "" {
		return boardID + "-col-" + strconv.Itoa(index)
	}

	return boardID + "-col-" + strings.Join(strings.Fields(col.ID), "-")
}

// kanbanCountLabel returns the screen-reader text for a column's card
// count badge ("3 cards" / "1 card" / "no cards").
func kanbanCountLabel(count int) string {
	switch count {
	case 0:
		return "no cards"
	case 1:
		return "1 card"
	default:
		return strconv.Itoa(count) + " cards"
	}
}

// KanbanMove is one decoded card move: the card (its consumer-supplied ID),
// the target column ID, and the 0-based insertion position within that
// column's cards (0 = top).
type KanbanMove struct {
	Card   string
	Column string
	Index  int
}

// errKanbanMoveMissingFields reports a move whose card or column field was
// empty — the two values the board can never omit when wired.
var errKanbanMoveMissingFields = errors.New("kanban move: card and column are required")

// errKanbanMoveBadIndex reports a move whose index field was neither a
// non-negative integer nor empty.
var errKanbanMoveBadIndex = errors.New("kanban move: index must be a non-negative integer")

// ParseKanbanMove decodes a KanbanBoard move submission (form body or query
// parameters) into a KanbanMove. A missing or non-numeric index defaults to
// 0 (insert at top); a negative index is an error. Empty card or column
// fields are an error — a wired board always sends both.
func ParseKanbanMove(r *http.Request) (KanbanMove, error) {
	if err := r.ParseForm(); err != nil {
		return KanbanMove{}, fmt.Errorf("parse kanban move form: %w", err)
	}

	move := KanbanMove{
		Card:   r.PostForm.Get(kanbanFieldCard),
		Column: r.PostForm.Get(kanbanFieldColumn),
		Index:  0,
	}

	if move.Card == "" || move.Column == "" {
		return KanbanMove{}, errKanbanMoveMissingFields
	}

	raw := r.PostForm.Get(kanbanFieldIndex)
	if raw == "" {
		return move, nil
	}

	index, err := strconv.Atoi(raw)
	if err != nil || index < 0 {
		return KanbanMove{}, fmt.Errorf("%w: got %q", errKanbanMoveBadIndex, raw)
	}

	move.Index = index

	return move, nil
}

// kanbanWireAttributes renders the wiring for a KanbanBoard's hidden move
// form. It copies the action (never mutates the consumer's) and applies the
// board defaults: an unspecified Method becomes POST (a card move is a
// mutation), an unspecified Event becomes submit, and an unspecified
// ContentType becomes form encoding so the hidden inputs travel under
// Datastar too (the same defaults forms.Form applies). Under htmx an empty
// Target resolves to the board root and hx-swap="outerHTML" is added — the
// documented response contract is "re-render the whole board". An empty URL
// wires nothing (read-only board).
func kanbanWireAttributes(w *wire.Action, boardID string) templ.Attributes {
	if w == nil {
		return nil
	}

	action := *w

	if action.Method == wire.MethodUnspecified {
		action.Method = wire.MethodPost
	}

	if action.Event == wire.EventUnspecified {
		action.Event = wire.EventSubmit
	}

	if action.ContentType == wire.ContentTypeUnspecified {
		action.ContentType = wire.ContentTypeForm
	}

	if action.Target == "" && action.Transport != wire.TransportDatastar {
		action.Target = "#" + boardID
	}

	attrs := action.Attributes()
	if attrs == nil {
		return nil
	}

	if action.Transport != wire.TransportDatastar {
		attrs["hx-swap"] = "outerHTML"
	}

	return attrs
}

// kanbanJS returns the singleton JavaScript for kanban card moves. All
// affordances funnel into one pipeline: fill the hidden move form's fields
// (card, column, index) and call requestSubmit() — the wire attributes on
// that form (hx-post or data-on:submit) perform the exchange, so the script
// itself is transport-agnostic. Drag-and-drop computes the insertion index
// from the pointer position between cards (adjusting for the dragged card's
// own removal when moving within a column); the keyboard buttons move a card
// to the end of the adjacent column. A drop whose source card lives on a
// DIFFERENT board is ignored (each wired board owns its cards), and every
// submitted move announces its completion into the board's live region once
// the re-rendered board lands (kanbanAnnounceJS).
func kanbanJS() string {
	return `if(!window.tcKanbanAttached){window.tcKanbanAttached=true;` +
		`function tcKbBoard(el){return el.closest('[data-tc-kanban]');}` +
		`function tcKbEsc(v){return (window.CSS&&CSS.escape)?CSS.escape(v):v;}` +
		`function tcKbForm(b){return b.querySelector('[data-tc-kanban-form]');}` +
		`function tcKbCards(zone){return Array.prototype.slice.call(zone.querySelectorAll(':scope > [data-tc-kanban-card]'));}` +
		`function tcKbSubmit(b,cardId,colId,index){` +
		`var f=tcKbForm(b);if(!f)return false;` +
		`f.querySelector('[data-tc-kanban-f-card]').value=cardId;` +
		`f.querySelector('[data-tc-kanban-f-column]').value=colId;` +
		`f.querySelector('[data-tc-kanban-f-index]').value=String(index);` +
		`var live=b.querySelector('[data-tc-kanban-live]');` +
		`var t=b.querySelector('[data-tc-kanban-card="'+tcKbEsc(cardId)+'"] [data-tc-kanban-title]');` +
		`var c=b.querySelector('[data-tc-kanban-column-body="'+tcKbEsc(colId)+'"]');` +
		`var cn=c?c.getAttribute('data-tc-kanban-col-title'):colId;` +
		`var title=t?t.textContent:cardId;` +
		`if(live){live.textContent='Moving '+title+' to '+(cn||colId)+'.';}` +
		`tcKbAnnounce='Moved '+title+' to '+(cn||colId)+'.';` +
		`tcKbAnnounceAfter(b);` +
		`if(typeof f.requestSubmit==='function'){f.requestSubmit();}else{f.submit();}` +
		`return true;}` +
		`var tcKbSrc=null,tcKbZone=null,tcKbIdx=0,tcKbAnnounce=null;` +
		`function tcKbClear(zone){` +
		`if(tcKbZone&&tcKbZone!==zone){tcKbZone.removeAttribute('data-tc-kanban-over');tcKbZone.classList.remove('tc-kanban-drop-end');}` +
		`if(tcKbZone){var prev=tcKbZone.querySelector(':scope > .tc-kanban-drop-before');if(prev)prev.classList.remove('tc-kanban-drop-before');}` +
		`}` +
		kanbanDragJS() +
		kanbanAnnounceJS() +
		kanbanClickJS() +
		`}`
}

// kanbanDragJS returns the HTML5 drag-and-drop listeners: dragstart marks
// the source card, dragover computes the insertion index from the pointer
// position between cards (showing the drop indicator), drop submits the move
// (adjusting the index for the dragged card's own removal when reordering
// within a column, and skipping no-op drops), dragend/dragleave clean up.
func kanbanDragJS() string {
	return `document.addEventListener('dragstart',function(e){` +
		`var card=e.target.closest('[data-tc-kanban-card]');if(!card)return;` +
		`var b=tcKbBoard(card);if(!b||!tcKbForm(b))return;` +
		`tcKbSrc=card;` +
		`e.dataTransfer.setData('text/plain',card.getAttribute('data-tc-kanban-card'));` +
		`e.dataTransfer.effectAllowed='move';` +
		`card.classList.add('tc-kanban-dragging');` +
		`});` +
		`document.addEventListener('dragend',function(e){` +
		`var card=e.target.closest('[data-tc-kanban-card]');` +
		`if(card)card.classList.remove('tc-kanban-dragging');` +
		`tcKbClear(null);tcKbZone=null;tcKbSrc=null;tcKbIdx=0;` +
		`});` +
		`document.addEventListener('dragover',function(e){` +
		`var zone=e.target.closest('[data-tc-kanban-column-body]');if(!zone)return;` +
		`var b=tcKbBoard(zone);if(!b||!tcKbForm(b)||!tcKbSrc)return;` +
		`if(tcKbBoard(tcKbSrc)!==b)return;` +
		`e.preventDefault();` +
		`e.dataTransfer.dropEffect='move';` +
		`tcKbClear(zone);` +
		`tcKbZone=zone;` +
		`var cards=tcKbCards(zone);` +
		`var idx=cards.length;` +
		`for(var i=0;i<cards.length;i++){` +
		`var r=cards[i].getBoundingClientRect();` +
		`if(e.clientY<r.top+r.height/2){idx=i;break;}` +
		`}` +
		`zone.setAttribute('data-tc-kanban-over',String(idx));` +
		`zone.classList.remove('tc-kanban-drop-end');` +
		`if(idx<cards.length){cards[idx].classList.add('tc-kanban-drop-before');}` +
		`else{zone.classList.add('tc-kanban-drop-end');}` +
		`tcKbIdx=idx;` +
		`});` +
		`document.addEventListener('drop',function(e){` +
		`var zone=e.target.closest('[data-tc-kanban-column-body]');if(!zone)return;` +
		`var b=tcKbBoard(zone);if(!b||!tcKbForm(b))return;` +
		`e.preventDefault();` +
		`var cardId=e.dataTransfer.getData('text/plain');` +
		`var idx=tcKbZone===zone?tcKbIdx:0;` +
		`tcKbClear(null);tcKbZone=null;` +
		`var src=tcKbSrc;tcKbSrc=null;` +
		`if(!cardId||!src)return;` +
		`if(tcKbBoard(src)!==b)return;` +
		`if(src.closest('[data-tc-kanban-column-body]')===zone){` +
		`var cards=tcKbCards(zone);` +
		`var srcIdx=cards.indexOf(src);` +
		`if(srcIdx>-1){` +
		`if(srcIdx<idx)idx=idx-1;` +
		`if(idx===srcIdx)return;` +
		`}` +
		`}` +
		`tcKbSubmit(b,cardId,zone.getAttribute('data-tc-kanban-column-body'),idx);` +
		`});` +
		`document.addEventListener('dragleave',function(e){` +
		`var zone=e.target.closest('[data-tc-kanban-column-body]');if(!zone)return;` +
		`if(!zone.contains(e.relatedTarget))tcKbClear(null);` +
		`});`
}

// kanbanClickJS returns the click listener behind the per-card keyboard
// move buttons: it resolves the adjacent column, appends the card to its
// end, and submits the same hidden form as a drop.
func kanbanClickJS() string {
	return `document.addEventListener('click',function(e){` +
		`var btn=e.target.closest('[data-tc-kanban-move]');if(!btn)return;` +
		`var card=btn.closest('[data-tc-kanban-card]');if(!card)return;` +
		`var b=tcKbBoard(card);if(!b)return;` +
		`e.preventDefault();` +
		`var dir=btn.getAttribute('data-tc-kanban-move')==='next'?1:-1;` +
		`var zones=Array.prototype.slice.call(b.querySelectorAll('[data-tc-kanban-column-body]'));` +
		`var ci=zones.indexOf(card.closest('[data-tc-kanban-column-body]'));` +
		`var target=zones[ci+dir];if(!target)return;` +
		`var index=tcKbCards(target).length;` +
		`tcKbSubmit(b,card.getAttribute('data-tc-kanban-card'),target.getAttribute('data-tc-kanban-column-body'),index);` +
		`});`
}

// kanbanScriptComponent renders the kanban singleton script CSP-safe.
func kanbanScriptComponent(nonce string) templ.Component {
	return scriptComponent(nonce, kanbanJS(), "kanban script")
}
