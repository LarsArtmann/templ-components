package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/go-datastar/static"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// The dual-transport FORM e2e suite: the browser-level proof that a
// forms.Form with Wire submits its fields, re-renders server-side validation
// errors inline, and stays interactive — under BOTH the self-hosted htmx
// runtime and the pinned Datastar bundle, against ONE endpoint wrapped in the
// library's own wire.Handler middleware.
//
// This is the trust pin for forms.FormProps.Wire: string tests prove what we
// emit; only a real runtime proves the runtime accepts it.

const (
	wireFormHTMXRegion     = "#wire-form-htmx-region"
	wireFormDatastarRegion = "#wire-form-datastar-region"

	wireFormNameMissing = "Name is required."
	wireFormEmailBad    = "Enter an email address with a domain."

	// wireFormSettleWait gives htmx's settle phase (default 20ms) time to
	// process swapped-in nodes before the test interacts with them again —
	// a click inside that window hits an unwired form and falls through to a
	// native submit. Real users cannot click this fast; e2e clients can.
	wireFormSettleWait = 250 * time.Millisecond
)

// formSel scopes a CSS selector to one dialect's form region.
func formSel(region, rest string) string {
	return region + " " + rest
}

// regionHasText builds a JS predicate: the region's text includes needle.
func regionHasText(region, needle string) string {
	return `document.querySelector('` + region + `').innerText.includes(` + jsString(needle) + `)`
}

// formValueExpr builds a JS expression reading a field's current value.
func formValueExpr(region, field string) string {
	return `document.querySelector('` + formSel(region, field) + `').value`
}

// jsString quotes a Go string as a single-quoted JS literal.
func jsString(s string) string {
	return "'" + strings.ReplaceAll(s, `\`, `\\`) + "'"
}

// setFieldValue sets an input's value programmatically and fires the input
// event, so runtimes tracking the field observe the change. Deterministic
// replacement — unlike chromedp.SendKeys, which types at the cursor (position
// 0 for values set via the value attribute) and prepends.
func setFieldValue(region, field, value string) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(cctx context.Context) error {
		expr := `var i=document.querySelector('` + formSel(region, field) + `'); i.value=` + jsString(value) +
			`; i.dispatchEvent(new Event('input',{bubbles:true})); i.value`

		var out string

		return chromedp.Evaluate(expr, &out).Do(cctx)
	})
}

// waitSwapSettled blocks until a freshly swapped region is safe to interact
// with again (see wireFormSettleWait).
func waitSwapSettled() chromedp.Action {
	return chromedp.Sleep(wireFormSettleWait)
}

// wireFormState carries the submitted values and server-side validation
// verdicts re-rendered into the form fragment — the canonical round-trip.
type wireFormState struct {
	Name     string
	Email    string
	NameErr  string
	EmailErr string
}

func (s wireFormState) invalid() bool {
	return s.NameErr != "" || s.EmailErr != ""
}

func (s wireFormState) summaryErrors() []forms.ValidationError {
	errs := make([]forms.ValidationError, 0, 2)
	if s.NameErr != "" {
		errs = append(errs, forms.ValidationError{Field: "name", Message: s.NameErr})
	}

	if s.EmailErr != "" {
		errs = append(errs, forms.ValidationError{Field: "email", Message: s.EmailErr})
	}

	return errs
}

// wireFormFragment renders the form exactly as both the initial page and the
// endpoint response do: optional ValidationSummary on top, inputs that
// preserve values and carry inline errors, and the submit button. The htmx
// dialect carries hx-target on the form; the Datastar dialect relies on
// wire.Handler's response headers.
func wireFormFragment(dialect wire.Transport, st wireFormState) templ.Component {
	props := forms.FormProps{
		Method: forms.FormPost,
		Wire: &wire.Action{
			Transport: dialect,
			Method:    wire.MethodPost,
			URL:       "/api/wire/form",
		},
	}
	if dialect == wire.TransportHTMX {
		props.Wire.Target = wireFormHTMXRegion
	}

	children := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if errs := st.summaryErrors(); len(errs) > 0 {
			if err := forms.ValidationSummary(forms.ValidationSummaryProps{Errors: errs}).Render(ctx, w); err != nil {
				return err
			}
		}

		if err := forms.Input(forms.InputProps{
			Name:        "name",
			Label:       "Name",
			Placeholder: "Ada Lovelace",
			Value:       st.Name,
			Error:       st.NameErr,
		}).Render(ctx, w); err != nil {
			return err
		}

		if err := forms.Input(forms.InputProps{
			Name:  "email",
			Type:  forms.InputEmail,
			Label: "Email address",
			Value: st.Email,
			Error: st.EmailErr,
		}).Render(ctx, w); err != nil {
			return err
		}

		return display.Button(display.ButtonProps{
			Text:    "Subscribe",
			Type:    display.ButtonHTMLSubmit,
			Size:    display.ButtonSizeSM,
			Variant: display.ButtonPrimary,
		}).Render(ctx, w)
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return forms.Form(props).Render(templ.WithChildren(ctx, children), w)
	})
}

// wireFormVerdict echoes the submitted fields and the transport that carried
// them — the assertion needle that proves the fields traveled.
func wireFormVerdict(name, email, transport string) templ.Component {
	return feedback.InlineSuccess(
		"Subscribed " + name + " (" + email + ") via " + transport + ".",
	)
}

// wireFormEndpoint validates server-side and ALWAYS answers 200 with a
// fragment: the re-rendered form (errors inline, values preserved) on
// invalid input, or the success verdict plus a fresh form. 200 is the
// zero-config parity pattern: htmx 2.0.10's default responseHandling is
// {code:"[45]..", swap:false, error:true} (verified in the pinned bundle),
// and the pinned Datastar bundle dispatches a fetch lifecycle error event at
// status >= 400 — a 422 needs runtime configuration to swap in either.
func wireFormEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form body", http.StatusBadRequest)

		return
	}

	st := wireFormState{
		Name:  strings.TrimSpace(r.PostFormValue("name")),
		Email: strings.TrimSpace(r.PostFormValue("email")),
	}
	if st.Name == "" {
		st.NameErr = wireFormNameMissing
	}

	if !strings.Contains(st.Email, "@") || !strings.Contains(st.Email, ".") {
		st.EmailErr = wireFormEmailBad
	}

	dialect := wire.TransportHTMX
	transport := "htmx"

	if wire.IsDatastar(r) {
		dialect = wire.TransportDatastar
		transport = "datastar"
	}

	ctx := context.Background()
	if !st.invalid() {
		if err := wireFormVerdict(st.Name, st.Email, transport).Render(ctx, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
		// Fresh form after success.
		st = wireFormState{}
	}

	if err := wireFormFragment(dialect, st).Render(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// wireFormE2EServer serves the form E2E page, the pinned Datastar bundle,
// the compiled CSS, and the shared /api/wire/form endpoint. The endpoint is
// wrapped in wire.Handler so Datastar callers get response-header targeting
// for the datastar region — the production server-side contract.
func wireFormE2EServer(t *testing.T) *httptest.Server {
	t.Helper()

	css, err := loadCSS()
	if err != nil {
		t.Fatalf("load compiled CSS: %v", err)
	}

	props := layout.DefaultPageProps()
	props.Title = "Wire Form E2E — templ-components"
	props.CSSPath = "/app.css"
	props.HeadContent = templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, werr := io.WriteString(
			w,
			`<script>window.__dsReady=false;document.addEventListener('datastar-ready',function(){window.__dsReady=true;},{once:true});</script>`,
		)

		return werr
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(css)
	})

	mux.HandleFunc("/datastar.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = w.Write(static.Bytes())
	})

	mux.Handle("/api/wire/form", wire.Handler(wire.PatchTarget{
		Selector: wireFormDatastarRegion,
		Mode:     wire.PatchModeInner,
	}, http.HandlerFunc(wireFormEndpoint)))

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := wireFormE2EPage(props).Render(context.Background(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	return srv
}

// wireFormE2EPage composes the page shell with both dialect forms rendered
// into their region wrappers.
func wireFormE2EPage(props layout.PageProps) templ.Component {
	body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := datastar.SDKScript(datastar.SDKScriptProps{
			BaseProps: utils.BaseProps{Nonce: "wire-form-e2e-nonce"},
			Src:       "/datastar.js",
		}).Render(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `<div class="p-4 grid grid-cols-2 gap-8 items-start">`); err != nil {
			return err
		}

		if _, err := io.WriteString(
			w,
			`<div><h2 class="text-sm font-semibold mb-2">htmx dialect</h2><div id="wire-form-htmx-region">`,
		); err != nil {
			return err
		}

		if err := wireFormFragment(wire.TransportHTMX, wireFormState{}).Render(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(
			w,
			`</div></div><div><h2 class="text-sm font-semibold mb-2">Datastar dialect</h2><div id="wire-form-datastar-region">`,
		); err != nil {
			return err
		}

		if err := wireFormFragment(wire.TransportDatastar, wireFormState{}).Render(ctx, w); err != nil {
			return err
		}

		_, werr := io.WriteString(w, `</div></div></div>`)

		return werr
	})

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return layout.Base(props).Render(templ.WithChildren(ctx, body), w)
	})
}

// TestWireE2EHTMXFormSubmitsFields proves the headline feature under htmx:
// fill the wired form, submit, and the region receives the verdict echoing
// the fields — carried natively by hx-post.
func TestWireE2EHTMXFormSubmitsFields(t *testing.T) {
	t.Parallel()

	srv := wireFormE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var ok bool

	if err := chromedp.Run(ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(`document.readyState==='complete' && window.htmx!==undefined`, &ok),
		chromedp.SendKeys(formSel(wireFormHTMXRegion, `input[name="name"]`), "Ada Lovelace", chromedp.NodeVisible),
		chromedp.SendKeys(formSel(wireFormHTMXRegion, `input[name="email"]`), "ada@example.com", chromedp.NodeVisible),
		chromedp.Click(formSel(wireFormHTMXRegion, `button[type="submit"]`), chromedp.NodeVisible),
		chromedp.Poll(regionHasText(wireFormHTMXRegion, "Subscribed Ada Lovelace (ada@example.com) via htmx."), &ok),
	); err != nil {
		t.Fatalf("htmx form E2E: %v", err)
	}
}

// TestWireE2EDatastarFormSubmitsFields proves the same form under the pinned
// Datastar runtime: the fields serialize via {contentType:'form'}, the
// response-header patch lands in the region, and the verdict names datastar.
func TestWireE2EDatastarFormSubmitsFields(t *testing.T) {
	t.Parallel()

	srv := wireFormE2EServer(t)

	ctx, cancel := newTab(t)
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
	defer cancelTimeout()

	var ok bool

	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(srv.URL+"/"),
		chromedp.Poll(`window.__dsReady===true`, &ok),
		chromedp.SendKeys(formSel(wireFormDatastarRegion, `input[name="name"]`), "Grace Hopper", chromedp.NodeVisible),
		chromedp.SendKeys(
			formSel(wireFormDatastarRegion, `input[name="email"]`),
			"grace@example.com",
			chromedp.NodeVisible,
		),
		chromedp.Click(formSel(wireFormDatastarRegion, `button[type="submit"]`), chromedp.NodeVisible),
		chromedp.Poll(
			regionHasText(wireFormDatastarRegion, "Subscribed Grace Hopper (grace@example.com) via datastar."),
			&ok,
		),
	); err != nil {
		t.Fatalf("datastar form E2E: %v", err)
	}
}

// TestWireE2EFormValidationRoundTrip proves the #1 forms use case under both
// runtimes: submit invalid input (an email the HTML5 client gate accepts but
// the server rejects), see the ValidationSummary + inline field error,
// confirm the submitted value survived, fix it, and submit again — the
// re-rendered form must stay fully interactive.
func TestWireE2EFormValidationRoundTrip(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		region string
		gate   string
	}{
		{name: "htmx", region: wireFormHTMXRegion, gate: `document.readyState==='complete' && window.htmx!==undefined`},
		{name: "datastar", region: wireFormDatastarRegion, gate: `window.__dsReady===true`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := wireFormE2EServer(t)

			ctx, cancel := newTab(t)
			defer cancel()

			ctx, cancelTimeout := context.WithTimeout(ctx, 30*time.Second)
			defer cancelTimeout()

			var (
				ok        bool
				preserved string
			)

			if err := chromedp.Run(ctx,
				chromedp.Navigate(srv.URL+"/"),
				chromedp.Poll(tc.gate, &ok),
				// "ada@example" passes the browser's HTML5 email gate but
				// fails the server's domain rule.
				chromedp.SendKeys(formSel(tc.region, `input[name="name"]`), "Ada Lovelace", chromedp.NodeVisible),
				chromedp.SendKeys(formSel(tc.region, `input[name="email"]`), "ada@example", chromedp.NodeVisible),
				chromedp.Click(formSel(tc.region, `button[type="submit"]`), chromedp.NodeVisible),
				// Error round-trip: summary + inline error visible, then let
				// the runtime settle before re-interacting (htmx wires swapped
				// nodes during its 20ms settle phase; Datastar initializes
				// observed mutations asynchronously too).
				chromedp.Poll(regionHasText(tc.region, "1 error found"), &ok),
				chromedp.Poll(regionHasText(tc.region, wireFormEmailBad), &ok),
				waitSwapSettled(),
				// The submitted value survived the re-render.
				chromedp.Evaluate(formValueExpr(tc.region, `input[name="email"]`), &preserved),
				// Fix the email in the re-rendered form and resubmit.
				setFieldValue(tc.region, `input[name="email"]`, "ada@example.com"),
				chromedp.Click(formSel(tc.region, `button[type="submit"]`), chromedp.NodeVisible),
				chromedp.Poll(regionHasText(tc.region, "Subscribed Ada Lovelace (ada@example.com)"), &ok),
			); err != nil {
				t.Fatalf("%s validation round-trip: %v", tc.name, err)
			}

			if preserved != "ada@example" {
				t.Fatalf("%s: submitted email not preserved across re-render; got %q", tc.name, preserved)
			}
		})
	}
}
