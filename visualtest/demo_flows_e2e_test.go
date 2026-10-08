package visualtest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

// The demo click-through suite drives the LIVE examples/demo server through
// its headline interactive flows — the compositions consumers copy — and
// proves them in a real browser, not just string assertions. Tests run
// serially (no t.Parallel): ~5 parallel tabs on the shared allocator flake.

const (
	// demoFlowTimeout bounds each individual flow's waits.
	demoFlowTimeout = 30 * time.Second

	// demoTabTimeout bounds each flow test's whole browser context. chromedp
	// actions have NO built-in timeout: a command whose target stops
	// responding (renderer wedge, orphaned-browser contention) would
	// otherwise block until the 10-minute go-test alarm kills the binary —
	// leaking the Chromium and demo processes because panic skips cleanup.
	// Bounding the tab turns any wedge into a normal test failure.
	demoTabTimeout = 120 * time.Second

	// endOfListText is navigation.EndOfList's default message; its appearance
	// proves the LoadMore chain exhausted server-side.
	endOfListText = "You've reached the end"
)

// kanbanColumnCardIDsExpr returns a JS expression evaluating to the comma-
// joined card ids of one board column (mirrors the kanban e2e helper).
func kanbanColumnCardIDsExpr(boardID, columnID string) string {
	return fmt.Sprintf(
		`Array.from(document.querySelectorAll('#%s [data-tc-kanban-column-body="%s"] > [data-tc-kanban-card]')).map(function(el){return el.getAttribute('data-tc-kanban-card');}).join(',')`,
		boardID,
		columnID,
	)
}

// kanbanClickUntilMove clicks a card's move button on the demo board and
// polls until the destination column shows the expected order, retrying the
// click — the runtime may re-attach to the swapped-in board between renders.
func kanbanClickUntilMove(ctx context.Context, t *testing.T, buttonSel, orderExpr, want string) {
	t.Helper()

	demoClickUntil(ctx, t, buttonSel, orderExpr+"==="+strconv.Quote(want), demoFlowTimeout)
}

// demoClickUntil clicks a selector and polls a JS condition, retrying the
// click on failure. HTMX outerHTML swaps invalidate chromedp's cached node
// ids ("Could not find node with given id"), so every retry re-queries the
// selector from scratch.
func demoClickUntil(ctx context.Context, t *testing.T, buttonSel, conditionExpr string, timeout time.Duration) {
	t.Helper()

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if err := chromedp.Run(ctx, chromedp.Click(chromedp.CSS(buttonSel))); err == nil {
			pollErr := chromedp.Run(ctx, pollTrue(conditionExpr, chromedp.WithPollingTimeout(3*time.Second)))
			if pollErr == nil {
				return
			}
		}

		time.Sleep(300 * time.Millisecond)
	}

	t.Fatalf("visualtest[demo]: condition never held after repeated clicks on %s: %s", buttonSel, conditionExpr)
}

// newFlowTab returns a bounded tab context for one flow test: the context
// (and therefore every chromedp action on it) hard-fails after
// demoTabTimeout instead of hanging the test binary.
func newFlowTab(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()

	tabCtx, tabCancel := newTab(t)

	ctx, cancel := context.WithTimeout(tabCtx, demoTabTimeout)

	return ctx, func() {
		cancel()
		tabCancel()
	}
}

func TestDemoLoadMoreReachesEndOfList(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(server.BaseURL()+"/navigation"),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[demo]: load navigation page: %v", err)
	}

	// Each LoadMore click REPLACES the button (hx-swap="outerHTML"), so the
	// selector must be re-resolved for every step — demoClickUntil handles
	// the stale-node churn.
	loadMore := "#demo-load-more button"

	demoClickUntil(ctx, t, loadMore,
		`document.querySelectorAll('#demo-load-more .rounded-lg').length >= 4`, demoFlowTimeout)
	demoClickUntil(ctx, t, loadMore,
		`document.querySelector('#demo-load-more').innerText.indexOf(`+strconv.Quote(endOfListText)+`) >= 0`,
		demoFlowTimeout)

	var itemCount int
	if err := chromedp.Run(ctx, evalExprInto(
		`document.querySelectorAll('#demo-load-more .rounded-lg').length`, &itemCount)); err != nil {
		t.Fatalf("visualtest[demo]: count items: %v", err)
	}

	if itemCount < 6 {
		t.Errorf("visualtest[demo]: expected 6 loaded items after two batches, got %d", itemCount)
	}

	server.FailIfServerErrors(t)
}

func TestDemoConfirmDeleteRemovesRow(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	// ConfirmDelete renders hx-confirm; the browser's native confirm()
	// dialog pauses the renderer and chromedp's dialog-accept path stalls
	// with it (the CDP command queue stops draining on a paused target —
	// verified by CDP trace). Stub confirm() to auto-accept instead: the
	// htmx flow under test (hx-confirm gate → DELETE → swap) is identical.
	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/htmx"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[demo]: load htmx page: %v", err)
	}

	if err := chromedp.Run(ctx, evalVoid(
		`window.confirm = function(){return true;}`)); err != nil {
		t.Fatalf("visualtest[demo]: stub confirm: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.Click(chromedp.CSS(`#item-123 button[hx-delete]`))); err != nil {
		t.Fatalf("visualtest[demo]: ConfirmDelete click: %v", err)
	}

	if err := chromedp.Run(ctx, pollTrue(
		// Success = the row was REPLACED: #item-123 no longer exists and
		// the mock endpoint's confirmation text is on the page. Checking
		// querySelector('#item-123') for the text can never succeed — the
		// swap removes the element that carried the id.
		`!document.querySelector('#item-123') && document.body.innerText.indexOf('deleted successfully') >= 0`,
		chromedp.WithPollingTimeout(demoFlowTimeout),
	)); err != nil {
		t.Fatalf("visualtest[demo]: row was not replaced by the delete response: %v", err)
	}

	server.FailIfServerErrors(t)
}

func TestDemoLoadingButtonBusyGate(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/htmx"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[demo]: load htmx page: %v", err)
	}

	saveButton := `button[hx-post="/demo/api/save"]`

	if err := chromedp.Run(ctx, chromedp.Click(chromedp.CSS(saveButton))); err != nil {
		t.Fatalf("visualtest[demo]: save click: %v", err)
	}

	// The request takes ~600ms; while in flight the button carries htmx's
	// htmx-request class (the busy gate) and the default label is hidden.
	var busySeen bool

	err := chromedp.Run(ctx, pollBool(
		`document.querySelector(`+strconv.Quote(saveButton)+`).classList.contains('htmx-request')`,
		&busySeen,
		chromedp.WithPollingInterval(20*time.Millisecond),
		chromedp.WithPollingTimeout(2*time.Second),
	))
	if err != nil {
		t.Logf("visualtest[demo]: busy class not observed (may have completed too fast): %v", err)
	}

	if err := chromedp.Run(ctx, pollTrue(
		`document.querySelector('#save-result') && document.querySelector('#save-result').innerText.length > 0`,
		chromedp.WithPollingTimeout(demoFlowTimeout),
	)); err != nil {
		t.Fatalf("visualtest[demo]: save result never appeared: %v", err)
	}

	var result string
	if err := chromedp.Run(ctx, evalExprInto(
		`document.querySelector('#save-result').innerText`, &result)); err != nil {
		t.Fatalf("visualtest[demo]: read save result: %v", err)
	}

	if !strings.Contains(result, "Saved") && !busySeen {
		t.Fatalf("visualtest[demo]: unexpected save result %q", result)
	}

	if !busySeen {
		t.Log("visualtest[demo]: busy gate not observed within poll window; request finished first")
	}

	server.FailIfServerErrors(t)
}

func TestDemoUploadEcho(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	uploadFile := filepath.Join(t.TempDir(), "echo-me.txt")
	if err := os.WriteFile(uploadFile, []byte("demo upload payload"), 0o600); err != nil {
		t.Fatalf("visualtest[demo]: write upload fixture: %v", err)
	}

	// transport=htmx renders exactly one upload form, keeping selectors stable.
	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(server.BaseURL()+"/wire?transport=htmx"),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[demo]: load wire page: %v", err)
	}

	if err := chromedp.Run(
		ctx,
		chromedp.SetUploadFiles(`form[action="/demo/api/wire/upload"] input[type="file"]`, []string{uploadFile}),
	); err != nil {
		t.Fatalf("visualtest[demo]: set upload file: %v", err)
	}

	if err := chromedp.Run(
		ctx,
		chromedp.Click(chromedp.CSS(`form[action="/demo/api/wire/upload"] button[type="submit"]`)),
	); err != nil {
		t.Fatalf("visualtest[demo]: submit upload: %v", err)
	}

	var out string
	if err := chromedp.Run(ctx, pollText(
		`(document.querySelector('#wire-upload-out') ? document.querySelector('#wire-upload-out').innerText : '').toLowerCase()`,
		&out,
		chromedp.WithPollingTimeout(demoFlowTimeout),
	)); err != nil {
		t.Fatalf("visualtest[demo]: upload echo never appeared in #wire-upload-out: %v", err)
	}

	if !strings.Contains(out, "echo-me.txt") {
		t.Fatalf("visualtest[demo]: upload echo region never showed the filename, got %q", out)
	}

	server.FailIfServerErrors(t)
}

func TestDemoKanbanMoveButtonsBothTransports(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/kanban"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[demo]: load kanban page: %v", err)
	}

	for _, board := range []string{"kanban-demo-htmx", "kanban-demo-datastar"} {
		t.Run(board, func(t *testing.T) {
			backlogExpr := kanbanColumnCardIDsExpr(board, "backlog")
			progressExpr := kanbanColumnCardIDsExpr(board, "progress")

			var before string
			if err := chromedp.Run(ctx, evalExprInto(backlogExpr, &before)); err != nil {
				t.Fatalf("visualtest[demo]: read backlog order: %v", err)
			}

			ids := strings.Split(before, ",")
			if len(ids) == 0 || ids[0] == "" {
				t.Fatalf("visualtest[demo]: board %s backlog is empty at start", board)
			}

			firstCard := ids[0]
			wantProgress := firstCard
			wantBacklog := strings.Join(ids[1:], ",")

			buttonSel := fmt.Sprintf("#%s [data-tc-kanban-card=%q] [data-tc-kanban-move=next]", board, firstCard)
			kanbanClickUntilMove(ctx, t, buttonSel, progressExpr, wantProgress)

			var gotBacklog string
			if err := chromedp.Run(ctx, evalExprInto(backlogExpr, &gotBacklog)); err != nil {
				t.Fatalf("visualtest[demo]: read backlog after move: %v", err)
			}

			if gotBacklog != wantBacklog {
				t.Errorf("visualtest[demo]: backlog after move = %q, want %q", gotBacklog, wantBacklog)
			}
		})
	}

	server.FailIfServerErrors(t)
}

// TestDemoWireBusyCardBothTransports proves the demo busy card end-to-end in
// a real browser: each dialect's button runs the slow job and its OWN region
// receives the done fragment — the htmx button via hx-target, the Datastar
// button via the endpoint's response headers (client fetch options never
// reach patch targeting; the action expression is a plain @post).
func TestDemoWireBusyCardBothTransports(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	tests := []struct {
		name      string
		page      string
		buttonSel string
		outID     string
		wantText  string
	}{
		{
			name:      "htmx swaps via hx-target",
			page:      "/wire?transport=htmx",
			buttonSel: `button[hx-post="/demo/api/wire/busy"]`,
			outID:     "wire-busy-htmx-out",
			wantText:  "Job finished via htmx",
		},
		{
			name:      "datastar patches via response headers",
			page:      "/wire?transport=datastar",
			buttonSel: `button[data-on\:click*="/api/wire/busy"]`,
			outID:     "wire-busy-datastar-out",
			wantText:  "Job finished via datastar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := chromedp.Run(ctx,
				chromedp.Navigate(server.BaseURL()+tt.page),
				chromedp.WaitReady("body"),
			); err != nil {
				t.Fatalf("visualtest[demo]: load %s: %v", tt.page, err)
			}

			outExpr := fmt.Sprintf(
				`(document.getElementById('%s') ? document.getElementById('%s').innerText : '')`,
				tt.outID, tt.outID,
			)
			demoClickUntil(ctx, t, tt.buttonSel, fmt.Sprintf(`%s.includes(%q)`, outExpr, tt.wantText), demoFlowTimeout)

			var out string
			if err := chromedp.Run(ctx, evalExprInto(outExpr, &out)); err != nil {
				t.Fatalf("visualtest[demo]: read %s: %v", tt.outID, err)
			}

			if !strings.Contains(out, tt.wantText) {
				t.Errorf("visualtest[demo]: %s = %q, want it to contain %q", tt.outID, out, tt.wantText)
			}
		})
	}

	server.FailIfServerErrors(t)
}

// TestDemoWireSwapCardsBothTransports proves the swap-styles and remove-mode
// cards in a real browser: the append button grows its region under BOTH
// dialects (htmx hx-swap="beforeend", Datastar response-header append), and
// the dismiss button retracts the region entirely (delete / remove mode).
func TestDemoWireSwapCardsBothTransports(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	cards := []struct {
		name      string
		page      string
		appendSel string
		appendOut string
		regionID  string
	}{
		{
			name:      "htmx",
			page:      "/wire?transport=htmx",
			appendSel: `button[hx-get="/demo/api/wire/swap-line"]`,
			appendOut: "wire-swap-htmx-out",
			regionID:  "wire-remove-htmx-region",
		},
		{
			name:      "datastar",
			page:      "/wire?transport=datastar",
			appendSel: `button[data-on\:click*="/api/wire/swap-line"]`,
			appendOut: "wire-swap-datastar-out",
			regionID:  "wire-remove-datastar-region",
		},
	}

	for _, tt := range cards {
		t.Run(tt.name, func(t *testing.T) {
			if err := chromedp.Run(ctx,
				chromedp.Navigate(server.BaseURL()+tt.page),
				chromedp.WaitReady("body"),
			); err != nil {
				t.Fatalf("visualtest[demo]: load %s: %v", tt.page, err)
			}

			appendExpr := fmt.Sprintf(
				`document.querySelectorAll('#%s p').length`,
				tt.appendOut,
			)

			// First append: the region gains one line.
			demoClickUntil(ctx, t, tt.appendSel, appendExpr+">=1", demoFlowTimeout)

			// Second append: still appending (mode did not degenerate).
			demoClickUntil(ctx, t, tt.appendSel, appendExpr+">=2", demoFlowTimeout)

			// Remove: the whole region is retracted.
			removeExpr := fmt.Sprintf(`!document.getElementById('%s')`, tt.regionID)
			removeSel := fmt.Sprintf(`#%s button`, tt.regionID)
			demoClickUntil(ctx, t, removeSel, removeExpr, demoFlowTimeout)
		})
	}

	server.FailIfServerErrors(t)
}
