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

// copyButtonE2EPage renders BOTH CopyButton variants — the default <button>
// and the Href <a>. The announcement contract under test is identical for
// both: live semantics live in the SIBLING [data-tc-copy-status] region
// (role=status, sr-only), never inside the interactive element.
func copyButtonE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		button := display.CopyButtonProps{
			Nonce: "e2e-nonce", ID: "e2e-copy-button",
			Text: "go get github.com/larsartmann/templ-components",
		}
		if err := display.CopyButton(button).Render(ctx, w); err != nil {
			return err
		}

		link := display.CopyButtonProps{
			Nonce: "e2e-nonce", ID: "e2e-copy-link",
			Text: "https://github.com/larsartmann/templ-components",
			Href: "https://github.com/larsartmann/templ-components",
		}
		if err := display.CopyButton(link).Render(ctx, w); err != nil {
			return err
		}

		return nil
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(layout.DefaultPageProps()).Render(templ.WithChildren(ctx, body), w)
	})
}

func copyButtonE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = copyButtonE2EPage().Render(r.Context(), w)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestCopyButtonE2EAnnouncesViaSiblingRegion proves the copy announcement in
// a real browser for both variants: the interactive element carries no live
// semantics, a click fills the sibling role=status region with the copied
// label, and the region clears after 2s so repeat copies re-announce.
// Serial — shared browser allocator.
func TestCopyButtonE2EAnnouncesViaSiblingRegion(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := copyButtonE2EServer(t)

	for _, variant := range []struct{ name, sel string }{
		{name: "button", sel: "#e2e-copy-button"},
		{name: "link", sel: "#e2e-copy-link"},
	} {
		t.Run(variant.name, func(t *testing.T) {
			ctx, cancel := newFlowTab(t)
			defer cancel()

			if err := chromedp.Do(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.WaitReady("body"),
			); err != nil {
				t.Fatalf("visualtest[copy-button]: navigate: %v", err)
			}

			// Contract pre-check: no aria-live/role=status on the interactive
			// element; its next sibling IS the live region.
			if err := chromedp.Do(ctx, pollTrue(
				`(function(){var el=document.querySelector('`+variant.sel+`');`+
					`if(!el)return false;`+
					`if(el.hasAttribute('aria-live'))return false;`+
					`if(el.getAttribute('role')==='status')return false;`+
					`var st=el.nextElementSibling;`+
					`return !!st&&st.hasAttribute('data-tc-copy-status')&&st.getAttribute('role')==='status';})()`,
				pollFast()...,
			)); err != nil {
				t.Fatalf("visualtest[copy-button] %s: announcement contract pre-check: %v", variant.name, err)
			}

			// A bubbling MouseEvent hits the document-level delegated listener
			// the same way a trusted click does (kanban/calendar-e2e-proven
			// pattern; trusted clicks are unreliable in headless).
			var clickNop string

			click := evalExprInto(
				`(document.querySelector('`+variant.sel+`').dispatchEvent(new MouseEvent('click',{bubbles:true,cancelable:true})),'')`,
				&clickNop,
			)

			var statusText string
			if err := chromedp.Do(
				ctx,
				click,
				pollText(
					`document.querySelector('`+variant.sel+`').nextElementSibling.textContent`,
					&statusText,
					pollFast()...),
			); err != nil {
				t.Fatalf("visualtest[copy-button] %s: copy announcement: %v", variant.name, err)
			}

			if statusText != "Copied!" {
				t.Fatalf(
					"visualtest[copy-button] %s: live region announced %q, want %q",
					variant.name,
					statusText,
					"Copied!",
				)
			}

			// The region clears after 2s — without the clear, a repeat copy
			// would not re-announce (screen readers skip unchanged text).
			if err := chromedp.Do(
				ctx,
				chromedp.Sleep(2300*time.Millisecond),
				pollTrue(
					`document.querySelector('`+variant.sel+`').nextElementSibling.textContent===''`,
					pollFast()...),
			); err != nil {
				t.Fatalf("visualtest[copy-button] %s: live region did not clear after 2s: %v", variant.name, err)
			}
		})
	}
}
