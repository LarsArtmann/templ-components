package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// kanbanSingleTransportPage renders a page with EXACTLY ONE wired board (the
// transport under test) — what a consumer embedding a single-transport board
// sees (backlog #258). The pending-visual page proves states; this one pins
// the single-transport LAYOUTS as PNG evidence.
func kanbanSingleTransportPage(transport wire.Transport) templ.Component {
	props := layout.DefaultPageProps()
	props.Title = "Kanban single transport"
	props.CSSPath = "/app.css"

	board := display.DefaultKanbanBoardProps()
	board.BaseProps = utils.BaseProps{ID: "kb-visual"}
	board.Columns = kanbanFlakyE2EColumns()
	action := wire.Action{URL: "/api/kanban/single"}
	if transport == wire.TransportDatastar {
		action.Transport = wire.TransportDatastar
	}
	board.Wire = &action
	board.CSRFToken = kanbanE2ECSRFToken

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="w-[64rem]">`); err != nil {
			return err
		}

		if err := display.KanbanBoard(board).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</div>`)

		return err
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

// kanbanSingleTransportServer serves the single-board page for one transport.
// The move endpoint is never called by this test — it pins the rendered
// layout, not the move flow — so a 204 sink is enough.
func kanbanSingleTransportServer(t *testing.T, transport wire.Transport) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(css)
	})

	mux.HandleFunc("/api/kanban/single", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := kanbanSingleTransportPage(transport).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return httptest.NewServer(mux)
}

// TestKanbanSingleTransportBoards pins both single-transport board renders
// (htmx-only and Datastar-only) as goldens: a consumer embedding one board —
// not the demo's side-by-side pair — must get the identical, working layout
// in either dialect (backlog #258).
func TestKanbanSingleTransportBoards(t *testing.T) {
	cases := []struct {
		name      string
		transport wire.Transport
	}{
		{name: "kanban/single_htmx", transport: wire.TransportHTMX},
		{name: "kanban/single_datastar", transport: wire.TransportDatastar},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			if !browserConfigured() {
				t.Skipf("visual tests skipped: %v (set CHROMEDP_CHROME_PATH)", errNoBrowser)
			}

			srv := kanbanSingleTransportServer(t, tcase.transport)
			defer srv.Close()

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
			defer cancelTimeout()

			if err := chromedp.Run(ctx,
				chromedp.EmulateViewport(viewportDesktopWidth, viewportDesktopHeight),
				chromedp.Navigate(srv.URL+"/"),
				chromedp.WaitVisible("#kb-visual", chromedp.ByQuery),
				chromedp.Evaluate(`try { localStorage.setItem('theme', 'light'); } catch (e) {}`, nil),
				chromedp.Reload(),
				chromedp.WaitVisible("#kb-visual", chromedp.ByQuery),
			); err != nil {
				t.Fatalf("navigate: %v", err)
			}

			var shot []byte

			if err := chromedp.Run(ctx,
				waitAnimationsSettled(),
				chromedp.Sleep(settleDelay),
				chromedp.Screenshot("#kb-visual", &shot, chromedp.ByQuery, chromedp.NodeVisible),
			); err != nil {
				t.Fatalf("capture: %v", err)
			}

			if *update {
				writeGolden(t, tcase.name, shot)

				return
			}

			golden, exists := readGolden(t, tcase.name)
			if !exists {
				writeGolden(t, tcase.name, shot)
				t.Errorf("visualtest[%s]: no golden yet — wrote %s (re-run without -update to verify)", tcase.name, goldenPath(tcase.name))

				return
			}

			assertGoldenMatch(t, tcase.name, golden, shot)
		})
	}
}
