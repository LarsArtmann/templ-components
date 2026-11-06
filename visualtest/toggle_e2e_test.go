package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
)

// toggleE2EPage renders two Toggles — one server-rendered checked, one
// unchecked — so the switch contract is provable in both directions:
// role="switch" + native checked matching Checked, flipped by trusted label
// clicks alone. aria-checked is deliberately absent from the contract: ARIA
// in HTML forbids it on checkbox inputs, and HTML-AAM derives the switch
// state from the native checked property, so no JS sync (and no bfcache
// staleness) exists to test.
func toggleE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		on := forms.ToggleProps{
			Nonce: "e2e-nonce", ID: "e2e-toggle-on", Name: "notify-on",
			Label: "Notifications on", Checked: true,
			Attrs: templ.Attributes{"data-testid": "e2e-toggle-on-label"},
		}
		if err := forms.Toggle(on).Render(ctx, w); err != nil {
			return err
		}

		off := forms.ToggleProps{
			Nonce: "e2e-nonce", ID: "e2e-toggle-off", Name: "notify-off",
			Label: "Notifications off",
			Attrs: templ.Attributes{"data-testid": "e2e-toggle-off-label"},
		}
		if err := forms.Toggle(off).Render(ctx, w); err != nil {
			return err
		}

		return nil
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(layout.DefaultPageProps()).Render(templ.WithChildren(ctx, body), w)
	})
}

func toggleE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = toggleE2EPage().Render(r.Context(), w)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestToggleE2ESwitchSemantics proves the switch contract in a real browser:
// the server-rendered role/aria-checked pair matches Checked in both
// directions, a trusted click flips both .checked and aria-checked via the
// delegated change listener, and a pageshow dispatch repairs a stale
// aria-checked after .checked was restored behind the page's back (the
// bfcache/form-history path). Serial — shared browser allocator.
func TestToggleE2ESwitchSemantics(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := toggleE2EServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitReady("body"),
	); err != nil {
		t.Fatalf("visualtest[toggle]: navigate: %v", err)
	}

	// Contract pre-check: both inputs render role="switch" with the checked
	// property matching the server-rendered Checked value, plus type and name.
	if err := chromedp.Do(ctx, pollTrue(
		`(function(){
			var on=document.querySelector('input#e2e-toggle-on');
			var off=document.querySelector('input#e2e-toggle-off');
			if(!on||!off)return false;
			if(on.getAttribute('role')!=='switch'||off.getAttribute('role')!=='switch')return false;
			if(on.checked!==true)return false;
			if(off.checked!==false)return false;
			if(on.type!=='checkbox'||off.type!=='checkbox')return false;
			if(on.name!=='notify-on'||off.name!=='notify-off')return false;
			return true;})()`,
		pollFast()...,
	)); err != nil {
		t.Fatalf("visualtest[toggle]: server-rendered switch contract: %v", err)
	}

	// Trusted click on the label flips the unchecked toggle on. The input is
	// sr-only, so the click lands on the wrapping label (native label
	// activation toggles the input and fires a real change event).
	click := `document.querySelector('[data-testid="e2e-toggle-off-label"]').click()`
	if err := chromedp.Do(ctx,
		evalVoid(click),
		pollTrue(`document.querySelector('input#e2e-toggle-off').checked===true`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[toggle]: click flip to checked: %v", err)
	}

	if err := chromedp.Do(ctx,
		evalVoid(click),
		pollTrue(`document.querySelector('input#e2e-toggle-off').checked===false`, pollFast()...),
	); err != nil {
		t.Fatalf("visualtest[toggle]: click flip to unchecked: %v", err)
	}
}
