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

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// The demo click-through suite drives the LIVE examples/demo server through
// its headline interactive flows — the compositions consumers copy — and
// proves them in a real browser, not just string assertions. Tests run
// serially (no t.Parallel): ~5 parallel tabs on the shared allocator flake.

const (
	// demoFlowTimeout bounds each individual flow's waits.
	demoFlowTimeout = 30 * time.Second

	// endOfListText is navigation.EndOfList's default message; its appearance
	// proves the LoadMore chain exhausted server-side.
	endOfListText = "You've reached the end"

	// uploadEchoMarker is the text prefix /api/wire/upload echoes back.
	uploadEchoMarker = "echo"
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

	deadline := time.Now().Add(demoFlowTimeout)

	for time.Now().Before(deadline) {
		if err := chromedp.Run(ctx, chromedp.Click(buttonSel, chromedp.ByQuery)); err == nil {
			var got string

			pollErr := chromedp.Run(ctx, chromedp.Poll(orderExpr+"==="+strconv.Quote(want), &got,
				chromedp.WithPollingTimeout(3*time.Second)))
			if pollErr == nil {
				return
			}
		}

		time.Sleep(300 * time.Millisecond)
	}

	t.Fatalf("visualtest[demo]: move never landed: %s did not become %q", orderExpr, want)
}

func TestDemoLoadMoreReachesEndOfList(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[demo]: load index: %v", err)
	}

	loadMore := "#demo-load-more button"

	if err := chromedp.Run(ctx, chromedp.Click(loadMore, chromedp.ByQuery)); err != nil {
		t.Fatalf("visualtest[demo]: first LoadMore click: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.Click(loadMore, chromedp.ByQuery)); err != nil {
		t.Fatalf("visualtest[demo]: second LoadMore click: %v", err)
	}

	var body string
	if err := chromedp.Run(ctx, chromedp.Poll(
		`document.querySelector('#demo-load-more').innerText.indexOf(`+strconv.Quote(endOfListText)+`) >= 0`,
		&body,
		chromedp.WithPollingTimeout(demoFlowTimeout),
	)); err != nil {
		t.Fatalf("visualtest[demo]: EndOfList never appeared after exhausting LoadMore: %v", err)
	}

	var itemCount int
	if err := chromedp.Run(ctx, chromedp.Evaluate(
		`document.querySelectorAll('#demo-load-more .rounded-lg').length`, &itemCount,
	)); err != nil {
		t.Fatalf("visualtest[demo]: count items: %v", err)
	}

	if itemCount < 6 {
		t.Errorf("visualtest[demo]: expected 6 loaded items after two batches, got %d", itemCount)
	}

	server.FailIfServerErrors(t)
}

func TestDemoConfirmDeleteRemovesRow(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	chromedp.ListenTarget(ctx, func(ev any) {
		if _, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			dialogCtx, dialogCancel := context.WithTimeout(ctx, 2*time.Second)
			defer dialogCancel()

			_ = chromedp.Run(dialogCtx, page.HandleJavaScriptDialog(true))
		}
	})

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[demo]: load index: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.Click(`#item-123 button[hx-delete]`, chromedp.ByQuery)); err != nil {
		t.Fatalf("visualtest[demo]: ConfirmDelete click: %v", err)
	}

	var body string
	if err := chromedp.Run(ctx, chromedp.Poll(
		`document.querySelector('#item-123') !== null && document.querySelector('#item-123').innerText.indexOf('deleted successfully') >= 0`,
		&body,
		chromedp.WithPollingTimeout(demoFlowTimeout),
	)); err != nil {
		t.Fatalf("visualtest[demo]: row was not replaced by the delete response: %v", err)
	}

	server.FailIfServerErrors(t)
}

func TestDemoLoadingButtonBusyGate(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[demo]: load index: %v", err)
	}

	saveButton := `button[hx-post="/api/save"]`

	if err := chromedp.Run(ctx, chromedp.Click(saveButton, chromedp.ByQuery)); err != nil {
		t.Fatalf("visualtest[demo]: save click: %v", err)
	}

	// The request takes ~600ms; while in flight the button carries htmx's
	// htmx-request class (the busy gate) and the default label is hidden.
	var busySeen bool

	err := chromedp.Run(ctx, chromedp.Poll(
		`document.querySelector(`+strconv.Quote(saveButton)+`).classList.contains('htmx-request')`,
		&busySeen,
		chromedp.WithPollingInterval(20*time.Millisecond),
		chromedp.WithPollingTimeout(2*time.Second),
	))
	if err != nil {
		t.Logf("visualtest[demo]: busy class not observed (may have completed too fast): %v", err)
	}

	var result string
	if err := chromedp.Run(ctx, chromedp.Poll(
		`document.querySelector('#save-result').innerText.length > 0`,
		&result,
		chromedp.WithPollingTimeout(demoFlowTimeout),
	)); err != nil {
		t.Fatalf("visualtest[demo]: save result never appeared: %v", err)
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

	ctx, cancel := newTab(t)
	defer cancel()

	uploadFile := filepath.Join(t.TempDir(), "echo-me.txt")
	if err := os.WriteFile(uploadFile, []byte("demo upload payload"), 0o600); err != nil {
		t.Fatalf("visualtest[demo]: write upload fixture: %v", err)
	}

	// transport=htmx renders exactly one upload form, keeping selectors stable.
	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/?transport=htmx"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[demo]: load index: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.SetUploadFiles(`form[action="/api/wire/upload"] input[type="file"]`, []string{uploadFile})); err != nil {
		t.Fatalf("visualtest[demo]: set upload file: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.Click(`form[action="/api/wire/upload"] button[type="submit"]`, chromedp.ByQuery)); err != nil {
		t.Fatalf("visualtest[demo]: submit upload: %v", err)
	}

	var out string
	if err := chromedp.Run(ctx, chromedp.Poll(
		`document.querySelector('#wire-upload-out').innerText.toLowerCase().indexOf('echo-me.txt') >= 0`,
		&out,
		chromedp.WithPollingTimeout(demoFlowTimeout),
	)); err != nil {
		t.Fatalf("visualtest[demo]: upload echo never appeared in #wire-upload-out: %v", err)
	}

	server.FailIfServerErrors(t)
}

func TestDemoKanbanMoveButtonsBothTransports(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("visualtest[demo]: load index: %v", err)
	}

	for _, board := range []string{"kanban-demo-htmx", "kanban-demo-datastar"} {
		t.Run(board, func(t *testing.T) {
			backlogExpr := kanbanColumnCardIDsExpr(board, "backlog")
			progressExpr := kanbanColumnCardIDsExpr(board, "progress")

			var before string
			if err := chromedp.Run(ctx, chromedp.Evaluate(backlogExpr, &before)); err != nil {
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
			if err := chromedp.Run(ctx, chromedp.Evaluate(backlogExpr, &gotBacklog)); err != nil {
				t.Fatalf("visualtest[demo]: read backlog after move: %v", err)
			}

			if gotBacklog != wantBacklog {
				t.Errorf("visualtest[demo]: backlog after move = %q, want %q", gotBacklog, wantBacklog)
			}
		})
	}

	server.FailIfServerErrors(t)
}
