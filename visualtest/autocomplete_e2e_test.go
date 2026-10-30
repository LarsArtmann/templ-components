package visualtest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
)

// autocompleteE2EPage renders a real <form> with four Inputs carrying autofill
// tokens (email, current-password, cc-number) plus one input whose token is
// invalid. The browser, not a string test, is the judge: password managers and
// autofill key off the emitted attribute, and a bogus token must not render
// one at all.
func autocompleteE2EPage() templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		email := forms.InputProps{
			Nonce: "e2e-nonce", ID: "e2e-ac-email", Name: "email",
			Type: forms.InputEmail, Autocomplete: forms.AutocompleteEmail,
		}
		password := forms.InputProps{
			Nonce: "e2e-nonce", ID: "e2e-ac-password", Name: "password",
			Type: forms.InputPassword, Autocomplete: forms.AutocompleteCurrentPassword,
		}
		card := forms.InputProps{
			Nonce: "e2e-nonce", ID: "e2e-ac-card", Name: "cc",
			Type: forms.InputTel, Autocomplete: forms.AutocompleteCardNumber,
		}
		bogus := forms.InputProps{
			Nonce: "e2e-nonce", ID: "e2e-ac-bogus", Name: "bogus",
			Type:         forms.InputText,
			Autocomplete: forms.AutocompleteType("bogus"),
		}

		if _, err := io.WriteString(w, `<form method="POST" action="/submit">`); err != nil {
			return err
		}

		for _, input := range []forms.InputProps{email, password, card, bogus} {
			if err := forms.Input(input).Render(ctx, w); err != nil {
				return err
			}
		}

		_, err := io.WriteString(w, `<button type="submit" id="e2e-ac-submit">Submit</button></form>`)

		return err
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(layout.DefaultPageProps()).Render(templ.WithChildren(ctx, body), w)
	})
}

func autocompleteE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = autocompleteE2EPage().Render(r.Context(), w)
	})
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		keys := make([]string, 0, len(r.PostForm))
		for k := range r.PostForm {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		pairs := make([]string, 0, len(keys))
		for _, k := range keys {
			pairs = append(pairs, k+"="+strings.Join(r.PostForm[k], ","))
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprint(w, "RECEIVED:"+strings.Join(pairs, "|"))
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// TestAutocompleteE2ETokensAndSubmission proves two things only a real browser
// can judge: (1) the autofill tokens land in the DOM exactly as emitted —
// including the invalid-token case, which must render NO autocomplete
// attribute at all — and (2) a real form submission round-trips the field
// values to the server. Serial — shared browser allocator.
func TestAutocompleteE2ETokensAndSubmission(t *testing.T) {
	if !browserConfigured() {
		t.Skip("no browser (set CHROMEDP_CHROME_PATH)")
	}

	srv := autocompleteE2EServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Do(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.WaitReady("body"),
		pollTrue(
			`(function(){
				var email=document.querySelector('input#e2e-ac-email');
				var password=document.querySelector('input#e2e-ac-password');
				var card=document.querySelector('input#e2e-ac-card');
				var bogus=document.querySelector('input#e2e-ac-bogus');
				if(!email||!password||!card||!bogus)return false;
				if(email.getAttribute('autocomplete')!=='email')return false;
				if(password.getAttribute('autocomplete')!=='current-password')return false;
				if(card.getAttribute('autocomplete')!=='cc-number')return false;
				return true;})()`,
			pollFast()...,
		),
	); err != nil {
		t.Fatalf("visualtest[autocomplete]: autofill tokens in DOM: %v", err)
	}

	// The invalid token must not render an autocomplete attribute — a typo in
	// the vocabulary would otherwise silently break the browser's heuristics.
	var bogusToken string

	if err := chromedp.Do(ctx, evalExprInto(
		`document.querySelector('input#e2e-ac-bogus').getAttribute('autocomplete')||''`, &bogusToken,
	)); err != nil {
		t.Fatalf("visualtest[autocomplete]: read bogus token attr: %v", err)
	}

	if bogusToken != "" {
		t.Fatalf("visualtest[autocomplete]: invalid AutocompleteType emitted %q, want NO attribute", bogusToken)
	}

	// Fill the form and really submit it — the server echoes the POST body.
	// The submit navigates away mid-Do, which tears down the JS execution
	// context the click ack lives on; that error is expected here (a genuinely
	// broken click fails the poll below instead).
	fillErr := chromedp.Do(ctx,
		chromedp.SetValue("input#e2e-ac-email", "ada@example.com"),
		chromedp.SetValue("input#e2e-ac-password", "s3cr3t"),
		chromedp.SetValue("input#e2e-ac-card", "4242424242424242"),
		chromedp.Click("button#e2e-ac-submit"),
	)
	if fillErr != nil {
		t.Logf("visualtest[autocomplete]: submit click ack: %v (expected during navigation)", fillErr)
	}

	// The full-page navigation tears down the JS execution context earlier
	// actions ran on, so the first poll attempts can race it with -32000;
	// retry with short polls until the new document answers.
	pollDeadline := time.Now().Add(10 * time.Second)

	for {
		pollErr := chromedp.Do(ctx, pollTrue(
			`document.body.innerText.startsWith('RECEIVED:')`,
			chromedp.WithPollingTimeout(750*time.Millisecond),
			chromedp.WithPollingInterval(100*time.Millisecond),
		))
		if pollErr == nil {
			break
		}

		if time.Now().After(pollDeadline) {
			t.Fatalf("visualtest[autocomplete]: form submission echo never appeared: %v", pollErr)
		}

		time.Sleep(100 * time.Millisecond)
	}

	var received string

	if err := chromedp.Do(ctx, evalExprInto(`document.body.innerText`, &received)); err != nil {
		t.Fatalf("visualtest[autocomplete]: read submission echo: %v", err)
	}

	want := "RECEIVED:bogus=|cc=4242424242424242|email=ada@example.com|password=s3cr3t"
	if received != want {
		t.Fatalf("visualtest[autocomplete]: server received %q, want %q", received, want)
	}
}
