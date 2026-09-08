package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/forms"
)

// Temporary diagnostic for the combobox open-state timeout.
func TestZZComboboxFocusDebug(t *testing.T) {
	cb := forms.DefaultComboboxProps()
	cb.Label = "Country"
	cb.Placeholder = "Select a country..."
	cb.Options = []forms.ComboboxOption{
		{Label: "United States", Value: "us"},
		{Label: "Canada", Value: "ca"},
		{Label: "Germany", Value: "de"},
	}
	cb.Nonce = "test-nonce"

	page, err := renderHTML(forms.Combobox(cb), defaultOptions(Options{}))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("page has script: %v", strings.Contains(page, "tcComboboxAttached"))
	t.Logf("page has listbox: %v", strings.Contains(page, `role="listbox"`))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, page)
	}))
	t.Cleanup(srv.Close)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 10*time.Second)
	defer cancelTimeout()

	var (
		focused       bool
		listHidden    bool
		expanded      string
		active        string
		syntheticOpen bool
	)

	err = chromedp.Run(ctx,
		chromedp.Navigate(srv.URL),
		chromedp.WaitVisible("#tc-root", chromedp.ByQuery),
		focusAction("#tc-root"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Evaluate(`document.activeElement.getAttribute('data-combobox-input') !== null`, &focused),
		chromedp.Evaluate(`document.querySelector('[role="listbox"]').classList.contains('hidden')`, &listHidden),
		chromedp.Evaluate(`document.querySelector('[role="combobox"]').getAttribute('aria-expanded')`, &expanded),
		chromedp.Evaluate(`document.activeElement.outerHTML.slice(0, 80)`, &active),
		chromedp.Evaluate(`document.hasFocus()`, &syntheticOpen),
	)
	if err != nil {
		t.Fatalf("chromedp: %v", err)
	}

	t.Logf("focused=%v listStillHidden=%v expanded=%q active=%q docHasFocus=%v", focused, listHidden, expanded, active, syntheticOpen)

	var syntheticListHidden bool

	err = chromedp.Run(ctx,
		chromedp.Evaluate(`document.activeElement.dispatchEvent(new FocusEvent('focusin', {bubbles: true}))`, nil),
		chromedp.Sleep(200*time.Millisecond),
		chromedp.Evaluate(`document.querySelector('[role="listbox"]').classList.contains('hidden')`, &syntheticListHidden),
	)
	if err != nil {
		t.Fatalf("chromedp synthetic: %v", err)
	}

	t.Logf("after synthetic focusin: listStillHidden=%v", syntheticListHidden)
}
