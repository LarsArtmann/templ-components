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
)

// overlayMatrixE2EPage renders the three remaining <dialog> overlays — Drawer,
// CommandPalette, Lightbox — each with its real trigger, so the shared
// backdrop-dismissal contract (position-based, never target-based) can be
// driven per component the way the Modal fixture does.
func overlayMatrixE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(
			w,
			`<button type="button" id="e2e-mx-drawer-open" data-tc-overlay-trigger="e2e-mx-drawer">Open drawer</button>`,
		)
		if err != nil {
			return err
		}

		drawer := display.DefaultDrawerProps()
		drawer.ID = "e2e-mx-drawer"
		drawer.Title = "Drawer"
		drawer.Nonce = "e2e-nonce"

		drawerChildren := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<p id="e2e-mx-drawer-body">Drawer content</p>`)

			return err
		})

		if err := display.Drawer(drawer).Render(templ.WithChildren(ctx, drawerChildren), w); err != nil {
			return err
		}

		palette := display.CommandPaletteProps{
			Nonce:        "e2e-nonce",
			ID:           "e2e-mx-palette",
			TriggerLabel: "Commands",
			Groups: []display.CommandGroup{
				{Label: "Actions", Items: []display.CommandItem{{Label: "Deploy"}}},
			},
		}
		if err := display.CommandPalette(palette).Render(ctx, w); err != nil {
			return err
		}

		lightbox := display.LightboxProps{
			Nonce: "e2e-nonce", ID: "e2e-mx-lightbox",
			TriggerLabel: "Open gallery",
			Images:       []display.LightboxImage{{Src: "/static/fixture.svg", Alt: "Gallery fixture"}},
		}

		return display.Lightbox(lightbox).Render(ctx, w)
	})

	props := layout.DefaultPageProps()
	props.CSSPath = "/app.css"

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

func overlayMatrixE2EServer(t *testing.T) *httptest.Server {
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
	mux.HandleFunc("/static/fixture.svg", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		_, _ = io.WriteString(
			w,
			`<svg xmlns="http://www.w3.org/2000/svg" width="320" height="180"><rect width="100%" height="100%" fill="#3b82f6"/></svg>`,
		)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = overlayMatrixE2EPage().Render(r.Context(), w)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestOverlayBackdropMatrix extends the Modal backdrop contract to the other
// dialog overlays: for Drawer, CommandPalette, and Lightbox — (1) a real
// click inside the dialog box never closes it, (2) a press+release outside
// dismisses, (3) a press inside released outside (selection drag) does not.
// The three overlays wire their backdrop logic through the same
// dialogBackdropOutsideJS helper, but each ships its own script, so the
// contract is proven per component, not per helper. Serial — shared allocator.
func TestOverlayBackdropMatrix(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := overlayMatrixE2EServer(t)

	overlays := []struct {
		name        string
		dialogID    string
		openTrigger string
	}{
		{name: "drawer", dialogID: "e2e-mx-drawer", openTrigger: "#e2e-mx-drawer-open"},
		{name: "command-palette", dialogID: "e2e-mx-palette", openTrigger: "[data-tc-palette-open]"},
		{name: "lightbox", dialogID: "e2e-mx-lightbox-dialog", openTrigger: "[data-tc-lightbox-open]"},
	}

	for _, tc := range overlays {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := newTab(t)
			defer cancel()

			tabCtx, tabCancel := context.WithTimeout(ctx, 120*time.Second)
			defer tabCancel()

			if err := chromedp.Do(tabCtx,
				chromedp.EmulateViewport(640, 480),
				chromedp.Navigate(srv.URL+"/"),
				chromedp.WaitReady("body"),
				chromedp.WaitVisible(chromedp.CSS(tc.openTrigger)),
				chromedp.Click(chromedp.CSS(tc.openTrigger)),
				pollTrue(`document.getElementById('`+tc.dialogID+`').open`, pollFast()...),
				// Settle the @starting-style open animation (drawer slides,
				// lightbox zooms) BEFORE reading the box — a mid-animation
				// rect sends the "inside" click outside or the "outside"
				// click inside.
				waitAnimationsSettled(),
				chromedp.Sleep(settleDelay),
			); err != nil {
				t.Fatalf("open via trigger: %v", err)
			}

			box := dialogBox(t, tabCtx, tc.dialogID)
			outsideX, outsideY := box.outsidePoint()

			// Inside click must keep the overlay open.
			if err := chromedp.Do(tabCtx, clickAt(box.Left+4, box.Top+4)); err != nil {
				t.Fatalf("click inside dialog: %v", err)
			}

			time.Sleep(300 * time.Millisecond)

			if !modalOpen(t, tabCtx, tc.dialogID) {
				t.Fatal("clicking inside the overlay closed it — dismissal is target-based again")
			}

			// Press inside, release outside (selection drag) must not dismiss.
			if err := chromedp.Do(tabCtx,
				dragFromTo(box.Left+8, box.Top+8, outsideX, outsideY),
			); err != nil {
				t.Fatalf("drag out of dialog: %v", err)
			}

			time.Sleep(300 * time.Millisecond)

			if !modalOpen(t, tabCtx, tc.dialogID) {
				t.Fatal("a press inside released outside dismissed the overlay")
			}

			// Press+release outside dismisses.
			if err := chromedp.Do(tabCtx,
				clickAt(outsideX, outsideY),
				pollTrue(`!document.getElementById('`+tc.dialogID+`').open`, pollFast()...),
			); err != nil {
				t.Fatalf("backdrop press+release did not dismiss: %v", err)
			}
		})
	}
}
