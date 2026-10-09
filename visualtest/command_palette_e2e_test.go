package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// commandPaletteE2EPage is the palette test page: trigger button + ⌘K hotkey,
// two link groups, and one wire item whose action outer-patches #output.
func commandPaletteE2EPage() templ.Component {
	wired := wire.Post("/api/run").WithTarget("#output")
	action := &wired

	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		palette := display.CommandPaletteProps{
			BaseProps:    utils.BaseProps{ID: "e2e-palette", Nonce: "e2e-nonce"},
			TriggerLabel: "Commands",
			Hotkey:       true,
			Groups: []display.CommandGroup{
				{Label: "Navigation", Items: []display.CommandItem{
					{Label: "Open target page", Href: "/target", Hint: "link flow"},
					{Label: "External docs", Href: "https://example.com"},
				}},
				{Label: "Actions", Items: []display.CommandItem{
					{Label: "Run job", Hint: "wire flow", Wire: action},
				}},
			},
		}
		if err := display.CommandPalette(palette).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `<div id="output" class="p-4">idle</div>`)

		return err
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(layout.DefaultPageProps()).Render(templ.WithChildren(ctx, body), w)
	})
}

func commandPaletteE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = commandPaletteE2EPage().Render(r.Context(), w)
	})

	mux.HandleFunc("/target", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><body><h1 id='target-marker'>Target reached</h1></body></html>"))
	})

	mux.Handle("POST /api/run", wire.Handler(wire.PatchTarget{
		Selector: "#output",
		Mode:     wire.PatchModeOuter,
	}, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<div id=\"output\">job ran</div>"))
	})))

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestCommandPaletteE2E proves the palette in a real browser: open via
// trigger, filter client-side (non-matching rows hide), navigate via Enter on
// a link row. Serial — shared browser allocator.
func TestCommandPaletteE2E(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := commandPaletteE2EServer(t)
	ctx, cancel := newTab(t)
	defer cancel()

	tabCtx, tabCancel := context.WithTimeout(ctx, 120*time.Second)
	defer tabCancel()

	if err := chromedp.Run(tabCtx, chromedp.Navigate(srv.URL)); err != nil {
		t.Fatalf("navigate: %v", err)
	}

	// Open via the trigger button.
	if err := chromedp.Run(tabCtx,
		chromedp.WaitVisible(`[data-tc-palette-open]`, chromedp.ByQuery),
		chromedp.Click(`[data-tc-palette-open]`, chromedp.ByQuery),
		pollTrue(`document.getElementById('e2e-palette').open`),
	); err != nil {
		t.Fatalf("open palette: %v", err)
	}

	// Filter: "target" keeps the link row visible and hides the others.
	if err := chromedp.Run(tabCtx,
		chromedp.SendKeys(`#e2e-palette-input`, "target"),
		pollTrue(`document.getElementById('e2e-palette-item-1').getAttribute('data-tc-palette-hidden')==='true'`),
		pollTrue(`document.getElementById('e2e-palette-item-0').getAttribute('data-tc-palette-hidden')==='false'`),
	); err != nil {
		t.Fatalf("filter: %v", err)
	}

	// Enter on the highlighted (first visible) link row navigates.
	if err := chromedp.Run(tabCtx,
		chromedp.SendKeys(`#e2e-palette-input`, "\r"),
		pollTrue(`location.pathname==='/target' && !!document.getElementById('target-marker')`),
	); err != nil {
		t.Fatalf("enter navigate: %v", err)
	}
}

// TestCommandPaletteWireE2E proves the wire flow: the ⌘K/Ctrl+K hotkey opens
// the palette, selecting a wire item closes it, and the htmx response patches
// the target region.
func TestCommandPaletteWireE2E(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := commandPaletteE2EServer(t)
	ctx, cancel := newTab(t)
	defer cancel()

	tabCtx, tabCancel := context.WithTimeout(ctx, 120*time.Second)
	defer tabCancel()

	if err := chromedp.Run(tabCtx, chromedp.Navigate(srv.URL)); err != nil {
		t.Fatalf("navigate: %v", err)
	}

	// Open via the hotkey (synthetic keydown — the JS listens on document).
	if err := chromedp.Run(tabCtx,
		chromedp.WaitVisible(`#e2e-palette`, chromedp.ByQuery),
		chromedp.Evaluate(`document.dispatchEvent(new KeyboardEvent('keydown',{key:'k',ctrlKey:true,bubbles:true,cancelable:true}));`, nil),
		pollTrue(`document.getElementById('e2e-palette').open`),
	); err != nil {
		t.Fatalf("hotkey open: %v", err)
	}

	// Click the wire row: palette closes, outer patch lands.
	if err := chromedp.Run(tabCtx,
		chromedp.Click(`#e2e-palette-item-2`, chromedp.ByQuery),
		pollTrue(`!document.getElementById('e2e-palette').open`),
		pollTrue(`document.getElementById('output').textContent.trim()==='job ran'`),
	); err != nil {
		t.Fatalf("wire select: %v", err)
	}
}
