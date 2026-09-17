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
// KanbanMove is one card move as submitted by a wired KanbanBoard. The form
// tags carry the exact field names the board's hidden move form sends
// (kanbanFieldCard/Column/Index); wire.DecodeForm fills them on the server.
type KanbanMove struct {
	Card   string `form:"card"`
	Column string `form:"column"`
	Index  int    `form:"index"`
}

// errKanbanMoveMissingFields reports a move whose card or column field was
// empty — the two values the board can never omit when wired.
var errKanbanMoveMissingFields = errors.New("kanban move: card and column are required")

// errKanbanMoveBadIndex reports a move whose index field was negative.
var errKanbanMoveBadIndex = errors.New("kanban move: index must be a non-negative integer")

// ParseKanbanMove decodes a KanbanBoard move submission (form body or query
// parameters) into a KanbanMove via wire.DecodeForm. A missing index defaults
// to 0 (insert at top); a negative or non-numeric index is an error. Empty
// card or column fields are an error — a wired board always sends both.
//
// The returned Index is advisory: the server owns the real ordering. When
// the board view is server-sorted, reject same-column moves (Card's current
// column == Column) — the sort re-asserts on re-render, so accepting them
// makes the drop look broken. Cross-column moves stay meaningful on sorted
// views. See docs/recipes/kanban-card-anatomy.md, "Moves on a sorted view".
func ParseKanbanMove(r *http.Request) (KanbanMove, error) {
	move, err := wire.DecodeForm[KanbanMove](r)
	if err != nil {
		return KanbanMove{}, fmt.Errorf("parse kanban move form: %w", err)
	}

	if move.Card == "" || move.Column == "" {
		return KanbanMove{}, errKanbanMoveMissingFields
	}

	if move.Index < 0 {
		return KanbanMove{}, fmt.Errorf("%w: got %d", errKanbanMoveBadIndex, move.Index)
	}

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
// DIFFERENT board is ignored (each wired board owns its cards).
//
// Moves apply optimistically (ADR-0041): before the request fires, the card
// is moved in the DOM and marked pending (tc-kanban-pending + aria-busy);
// the pending state is cleared when the re-rendered board lands and REVERTED
// with a visible flash + role="alert" announcement when the transport reports
// failure (htmx:afterRequest without success, datastar-fetch error or
// retries-failed). Every submitted move announces its completion into the
// board's live region once the re-rendered board lands (kanbanAnnounceJS).
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
		`var card=b.querySelector('[data-tc-kanban-card="'+tcKbEsc(cardId)+'"]');` +
		`var t=card?card.querySelector('[data-tc-kanban-title]'):null;` +
		`var c=b.querySelector('[data-tc-kanban-column-body="'+tcKbEsc(colId)+'"]');` +
		`var cn=c?c.getAttribute('data-tc-kanban-col-title'):colId;` +
		`var title=t?t.textContent:cardId;` +
		`if(live){live.textContent='Moving '+title+' to '+(cn||colId)+'.';}` +
		`tcKbAnnounce='Moved '+title+' to '+(cn||colId)+'.';` +
		`tcKbAnnounceAfter(b);` +
		`if(card&&c)tcKbOptimistic(b,card,c,index,title,cn||colId);` +
		`if(typeof f.requestSubmit==='function'){f.requestSubmit();}else{f.submit();}` +
		`return true;}` +
		`var tcKbSrc=null,tcKbZone=null,tcKbIdx=0,tcKbAnnounce=null;` +
		`function tcKbClear(zone){` +
		`if(tcKbZone&&tcKbZone!==zone){tcKbZone.removeAttribute('data-tc-kanban-over');tcKbZone.classList.remove('tc-kanban-drop-end');}` +
		`if(tcKbZone){var prev=tcKbZone.querySelector(':scope > .tc-kanban-drop-before');if(prev)prev.classList.remove('tc-kanban-drop-before');}` +
		`}` +
		kanbanDragJS() +
		kanbanAnnounceJS() +
		kanbanPendingJS() +
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

// kanbanAnnounceJS returns the post-swap confirmation mechanism. tcKbSubmit
// already announces "Moving X to Y." immediately; this adds "Moved X to Y."
// once the re-rendered board actually lands in the DOM, so screen-reader
// users hear the completed state change, not just the intent.
//
// It deliberately polls the board's live region instead of listening for
// transport events: htmx's afterSwap detail.target is invalidated by the
// outerHTML swap the board mandates, and Datastar has no global swap event
// at all. The signal is swap-agnostic: tcKbSubmit writes "Moving ..." into
// the live region before submitting, and the server always re-renders the
// region empty — so an EMPTY region after submit means the response landed,
// whether the runtime replaced the board element (htmx outerHTML: resolve
// the new node by id) or morphed it in place (Datastar outer patches keep
// the old element connected). The 30s budget (up from the original 5s) also
// covers slow servers whose swap would otherwise strand the pending register
// (kanbanPendingJS clears it via transport events regardless). Boards whose
// patch never touches the live region degrade to no post-swap announcement —
// the immediate "Moving" message still fires.
func kanbanAnnounceJS() string {
	return `function tcKbAnnounceIn(b){` +
		`if(!tcKbAnnounce||!b)return;` +
		`var msg=tcKbAnnounce;tcKbAnnounce=null;` +
		`var live=b.querySelector('[data-tc-kanban-live]');` +
		`if(live){live.textContent=msg;}` +
		`}` +
		`function tcKbAnnounceAfter(old){` +
		`var tries=0;` +
		`var timer=setInterval(function(){` +
		`tries++;` +
		`var b=old.isConnected?old:(old.id?document.getElementById(old.id):null);` +
		`if(b){` +
		`var live=b.querySelector('[data-tc-kanban-live]');` +
		`if(live&&live.textContent===''){clearInterval(timer);tcKbAnnounceIn(b);tcKbSucceed(b.id);return;}` +
		`}` +
		`if(tries>300){clearInterval(timer);}` +
		`},100);` +
		`}`
}

// kanbanPendingJS returns the optimistic-move pipeline (ADR-0041): the
// registry of in-flight moves, the DOM placement + count sync that make the
// move instant, the pending register (tc-kanban-pending + aria-busy), and
// the transport event listeners that clear it on success and revert it —
// with a visible flash and a role="alert" announcement — on failure.
//
// Success has three independent clearers so pending can never stick: htmx's
// afterRequest with successful===true (verified in the vendored htmx 2.0.10
// source: fires from onload after the swap, with successful = !isError), the
// Datastar bundle's datastar-fetch "finished" lifecycle type (verified in
// the pinned v1.0.3 bundle: dispatched in finally after every action — after
// an "error" revert it is a harmless no-op), and the announce poll's
// swap-landed signal above.
//
// Failure is honest: htmx's afterRequest without success covers 4xx/5xx AND
// network errors, aborts and timeouts (successful is undefined — never true —
// for those), with responseError/sendError as belt-and-braces; Datastar
// dispatches "error" (HTTP >= 400, immediate) and "retries-failed" (network
// errors retry under the runtime's default backoff first — up to ~2 minutes
// of honest "still trying" before the revert). Revert restores the card to
// its exact original position, unhides placeholders, recomputes counts and
// aria-labels from the DOM, flashes tc-kanban-move-failed for 4s, and writes
// the failure into the board's role="alert" region (the repo's urgency
// policy: role="alert", never aria-live="assertive").
//
// Every listener filters on data-tc-kanban-form so unrelated wired elements
// (the column Action add button, the page's reset button) can never trigger
// a revert; entries are keyed by board id so cross-board events stay inert.
func kanbanPendingJS() string {
	return kanbanPendingRegistryJS() + kanbanPendingPipelineJS() + kanbanPendingEventsJS()
}

// kanbanPendingRegistryJS returns the in-flight move registry (keyed by
// board id) plus the DOM helpers it drives: count-badge/aria-label syncing,
// optimistic placement, and empty-placeholder unhiding.
func kanbanPendingRegistryJS() string {
	return `var tcKbPending=[];` +
		`function tcKbPendingFor(bid){` +
		`for(var i=0;i<tcKbPending.length;i++){if(tcKbPending[i].id===bid)return true;}` +
		`return false;}` +
		`function tcKbForget(bid){tcKbPending=tcKbPending.filter(function(e){return e.id!==bid;});}` +
		`function tcKbCountLabel(n){return n===0?'no cards':(n===1?'1 card':n+' cards');}` +
		`function tcKbSyncCounts(zone){` +
		`if(!zone)return;` +
		`var col=zone.closest('[data-tc-kanban-column]');if(!col)return;` +
		`var n=tcKbCards(zone).length;` +
		`var badge=col.querySelector('[data-tc-kanban-count]');` +
		`if(badge)badge.textContent=String(n);` +
		`zone.setAttribute('aria-label',(zone.getAttribute('data-tc-kanban-col-title')||'')+': '+tcKbCountLabel(n));` +
		`}` +
		`function tcKbPlace(zone,card,idx){` +
		`var cards=tcKbCards(zone);` +
		`var self=cards.indexOf(card);` +
		`if(self>-1)cards.splice(self,1);` +
		`if(idx>=cards.length){zone.appendChild(card);}` +
		`else{zone.insertBefore(card,cards[idx]);}` +
		`}` +
		`function tcKbUnhideEmpties(zone){` +
		`if(!zone)return;` +
		`var phs=zone.querySelectorAll('[data-tc-kanban-empty]');` +
		`for(var i=0;i<phs.length;i++){phs[i].hidden=false;}` +
		`}` +
}

// kanbanPendingPipelineJS returns the register's state transitions:
// tcKbOptimistic (capture origin BEFORE placing — the revert restores that
// exact position), tcKbSucceed, tcKbRevert, and tcKbFailFrom.
func kanbanPendingPipelineJS() string {
	return `function tcKbOptimistic(b,card,zone,idx,title,dest){` +
		`var ph=zone.querySelector(':scope > [data-tc-kanban-empty]');` +
		`if(ph)ph.hidden=true;` +
		`var src=card.parentNode;` +
		`var next=card.nextSibling;` +
		`tcKbPlace(zone,card,idx);` +
		`tcKbSyncCounts(src);` +
		`tcKbSyncCounts(zone);` +
		`card.classList.add('tc-kanban-pending');` +
		`card.setAttribute('aria-busy','true');` +
		`tcKbPending.push({id:b.id,card:card,parent:src,next:next,msg:'Moving '+title+' to '+dest});` +
		`}` +
		`function tcKbSucceed(bid){` +
		`var b=document.getElementById(bid);` +
		`if(b){` +
		`var cards=b.querySelectorAll('[data-tc-kanban-card].tc-kanban-pending');` +
		`for(var i=0;i<cards.length;i++){cards[i].classList.remove('tc-kanban-pending');cards[i].removeAttribute('aria-busy');}` +
		`}` +
		`tcKbForget(bid);` +
		`}` +
		`function tcKbRevert(bid){` +
		`var b=document.getElementById(bid);` +
		`var rest=[];` +
		`for(var i=0;i<tcKbPending.length;i++){` +
		`var e=tcKbPending[i];` +
		`if(e.id!==bid){rest.push(e);continue;}` +
		`var cur=e.card.parentNode;` +
		`if(e.next&&e.next.parentNode===e.parent){e.parent.insertBefore(e.card,e.next);}` +
		`else{e.parent.appendChild(e.card);}` +
		`tcKbUnhideEmpties(e.parent);` +
		`tcKbUnhideEmpties(cur);` +
		`tcKbSyncCounts(e.parent);` +
		`tcKbSyncCounts(cur);` +
		`e.card.classList.remove('tc-kanban-pending');` +
		`e.card.removeAttribute('aria-busy');` +
		`e.card.classList.add('tc-kanban-move-failed');` +
		`(function(c){setTimeout(function(){c.classList.remove('tc-kanban-move-failed');},4000);})(e.card);` +
		`if(b){` +
		`var alertEl=b.querySelector('[data-tc-kanban-alert]');` +
		`if(alertEl)alertEl.textContent=e.msg+' failed. The board was restored.';` +
		`}` +
		`}` +
		`tcKbPending=rest;` +
		`}` +
		`function tcKbFailFrom(el){` +
		`if(!el||!el.closest)return;` +
		`var b=el.closest('[data-tc-kanban]');` +
		`if(!b||!b.id||!tcKbPendingFor(b.id))return;` +
		`tcKbAnnounce=null;` +
		`tcKbRevert(b.id);` +
		`}` +
}

// kanbanPendingEventsJS returns the transport event listeners. Success and
// failure are decided per runtime per the facts pinned on kanbanPendingJS.
func kanbanPendingEventsJS() string {
	return `document.addEventListener('htmx:afterRequest',function(e){` +
		`var d=e.detail||{};` +
		`var f=d.elt;` +
		`if(!f||!f.hasAttribute||!f.hasAttribute('data-tc-kanban-form'))return;` +
		`var b=f.closest('[data-tc-kanban]');` +
		`if(!b||!b.id)return;` +
		`if(d.successful===true){tcKbSucceed(b.id);}` +
		`else{tcKbFailFrom(f);}` +
		`});` +
		`document.addEventListener('htmx:responseError',function(e){tcKbFailFrom((e.detail||{}).elt);});` +
		`document.addEventListener('htmx:sendError',function(e){tcKbFailFrom((e.detail||{}).elt);});` +
		`document.addEventListener('datastar-fetch',function(e){` +
		`var d=e.detail||{};` +
		`var f=d.el;` +
		`if(!f||!f.hasAttribute||!f.hasAttribute('data-tc-kanban-form'))return;` +
		`var b=f.closest('[data-tc-kanban]');` +
		`if(!b||!b.id)return;` +
		`if(d.type==='error'||d.type==='retries-failed'){tcKbFailFrom(f);}` +
		`else if(d.type==='finished'){tcKbSucceed(b.id);}` +
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
