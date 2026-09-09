package visualtest

import (
	"testing"

	"github.com/chromedp/chromedp"
)

// setRTL returns a chromedp action flipping the document to right-to-left —
// the browser-level trigger for the logical-property mirroring the RTL
// compliance scanner can only approximate.
func setRTL() chromedp.Action {
	return chromedp.Evaluate(`document.documentElement.setAttribute('dir','rtl');`, nil)
}

// TestDemoRTLNoHorizontalOverflow proves the demo routes keep physical
// overflow behaviour under dir="rtl": logical properties must mirror without
// creating page-wide sideways scroll, exactly as in LTR.
func TestDemoRTLNoHorizontalOverflow(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	for _, route := range []string{"/", "/forms", "/users", "/recipes/dashboard"} {
		t.Run(route, func(t *testing.T) {
			if err := chromedp.Run(ctx,
				chromedp.Navigate(server.BaseURL()+route),
				chromedp.WaitReady("body"),
				setRTL(),
			); err != nil {
				t.Fatalf("visualtest[demo]: navigate %s: %v", route, err)
			}

			var overflow int64

			if err := chromedp.Run(ctx, chromedp.Evaluate(
				`document.documentElement.scrollWidth - document.documentElement.clientWidth`, &overflow,
			)); err != nil {
				t.Fatalf("visualtest[demo]: measure %s RTL overflow: %v", route, err)
			}

			if overflow > 1 {
				t.Errorf(
					"visualtest[demo]: %s overflows an RTL viewport by %dpx — a physical property or unmirrored offset is breaking mirroring",
					route,
					overflow,
				)
			}
		})
	}

	server.FailIfServerErrors(t)
}

// TestDemoRTLKanbanMoveWorks proves the kanban board stays operable under
// RTL on both transports: the move buttons still resolve, submit, and the
// server's re-render lands the card in the expected column.
func TestDemoRTLKanbanMoveWorks(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(server.BaseURL()+"/"),
		chromedp.WaitReady("body"),
		setRTL(),
	); err != nil {
		t.Fatalf("visualtest[demo]: navigate: %v", err)
	}

	for _, board := range []string{"kanban-demo-htmx", "kanban-demo-datastar"} {
		t.Run(board, func(t *testing.T) {
			backlogExpr := kanbanColumnCardIDsExpr(board, "backlog")
			progressExpr := kanbanColumnCardIDsExpr(board, "progress")

			var before string

			if err := chromedp.Run(ctx, chromedp.Evaluate(backlogExpr, &before)); err != nil {
				t.Fatalf("visualtest[demo]: read RTL backlog order: %v", err)
			}

			buttonSel := `#` + board + ` [data-tc-kanban-card="` + firstCSV(before) + `"] [data-tc-kanban-move=next]`
			kanbanClickUntilMove(ctx, t, buttonSel, progressExpr, firstCSV(before))
		})
	}

	server.FailIfServerErrors(t)
}

// firstCSV returns the first token of a comma-separated list, or "" when
// the list is empty.
func firstCSV(list string) string {
	for i := range len(list) {
		if list[i] == ',' {
			return list[:i]
		}
	}

	return list
}
