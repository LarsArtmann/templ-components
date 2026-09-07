# Recipe: Server-Side Validation with Dual-Transport Forms

**Audience:** Consumers building forms that submit via `forms.FormProps.Wire`
(htmx or Datastar) and need submit → server validation → inline field errors
on either runtime.

**Problem:** Whole-form submission is now transport-symmetric, but validation
round-trips have transport-specific traps: which region the re-rendered form
lands in, which HTTP status the runtimes actually swap, and how to keep the
form interactive after an error re-render.

**Outcome:** A copy-pasteable pattern where one handler validates, re-renders
the form with `ValidationSummary` + inline field errors (values preserved),
and works unchanged under htmx AND Datastar. Proven end-to-end in Chromium by
`visualtest/wire_form_e2e_test.go`.

---

## The flow

```text
browser                        one endpoint (wire.Handler)
────────                       ───────────────────────────
fill form
click Submit
  ├─ htmx:    hx-post ────────► ParseForm + validate
  └─ Datastar: @post(...,
      {contentType:'form'})
          │
          │                     invalid ──► 200 OK
          │                                 re-rendered form
          │                                 (summary + inline errors,
          │                                  values preserved)
          │                     valid ────► 200 OK
          │                                 verdict + fresh form
          ◄────────────────── fragment
  htmx swaps region (hx-target)
  Datastar patches region (response headers)
```

## 1. Render the form inside its swap region

The re-render replaces the region's content, so the form must live INSIDE it.
On error the user gets the same form back (with errors); on success a verdict
plus a fresh form.

```templ
<div id="form-region" aria-live="polite">
	@forms.Form(forms.FormProps{
		Method: forms.FormPost,
		Wire: &wire.Action{
			Transport: wire.TransportHTMX, // or TransportDatastar
			Method:    wire.MethodPost,
			URL:       "/api/subscribe",
			Target:    "#form-region", // htmx-only; Datastar targets via the response
		},
	}) {
		// children: see step 3
	}
</div>
```

## 2. One handler, wrapped in wire.Handler

```go
mux.Handle("/api/subscribe", wire.Handler(wire.PatchTarget{
	Selector: "#form-region", // Datastar callers get this via response headers
	Mode:     wire.PatchModeInner,
}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form body", http.StatusBadRequest)
		return
	}

	st := formState{
		Name:  strings.TrimSpace(r.PostFormValue("name")),
		Email: strings.TrimSpace(r.PostFormValue("email")),
	}
	if st.Name == "" {
		st.NameErr = "Name is required."
	}
	if !strings.Contains(st.Email, "@") || !strings.Contains(st.Email, ".") {
		st.EmailErr = "Enter an email address with a domain."
	}

	dialect := wire.TransportHTMX
	if wire.IsDatastar(r) {
		dialect = wire.TransportDatastar
	}

	if st.invalid() {
		componentOr500(w, r, subscribeForm(dialect, st)) // errors inline
		return
	}
	componentOr500(w, r, verdictAndFreshForm(st)) // success path
})))
```

Both dialects serialize the fields into a standard urlencoded body (htmx
natively; Datastar via `{contentType: 'form'}` — the `Form` component applies
that default), so `r.ParseForm`/`r.PostFormValue` serve both.

## 3. The error fragment: summary + inline errors + preserved values

```templ
templ subscribeForm(dialect wire.Transport, st formState) {
	if errs := st.summaryErrors(); len(errs) > 0 {
		@forms.ValidationSummary(forms.ValidationSummaryProps{Errors: errs})
	}
	@forms.Form(forms.FormProps{
		Method: forms.FormPost,
		Wire:   &wire.Action{Transport: dialect, Method: wire.MethodPost, URL: "/api/subscribe", Target: "#form-region"},
	}) {
		@forms.Input(forms.InputProps{Name: "name", Label: "Name", Value: st.Name, Error: st.NameErr})
		@forms.Input(forms.InputProps{Name: "email", Type: forms.InputEmail, Label: "Email address", Value: st.Email, Error: st.EmailErr})
		@display.Button(display.ButtonProps{Text: "Subscribe", Type: display.ButtonHTMLSubmit})
	}
}
```

`Input.Error` renders `aria-invalid`, the error id wiring, and the
`role="alert"` field message; `ValidationSummary` links each error to its
field. The library-level markup is pinned by the `form_validation_errors`
golden.

## Why errors MUST return 200 OK (verified runtime facts)

Both runtimes punish 4xx by default, so the error fragment travels as 200:

- **htmx 2.0.10** ships `responseHandling` defaults of
  `{code: "204", swap: false}, {code: "[23]..", swap: true}, {code: "[45]..", swap: false, error: true}`
  — a 422 does **not** swap (it fires `htmx:responseError`) unless you change
  `htmx.config.responseHandling` or handle `htmx:beforeSwap`.
- **Datastar** (pinned v1.0.3 bundle) dispatches a `datastar-fetch` lifecycle
  error event when a fetch resolves with `status >= 400`.

If you need RESTful 422 semantics anyway: return 422, render the same
fragment, and configure the client — htmx via
`htmx.config.responseHandling = [{code:"422", swap:true}]` (or a
`htmx:beforeSwap` listener setting `detail.shouldSwap`), Datastar by treating
the fetch error event as a signal to re-request. The zero-config path above
is what the library demo and e2e tests use.

## Testing the round-trip

Two levels, both in the repo:

1. **String/HTTP level** (demo tests): POST the urlencoded body with and
   without the `Datastar-Request` header; assert the summary text, inline
   error, preserved `value="..."`, and — for Datastar calls — the
   `Datastar-Selector`/`Datastar-Mode` response headers.
2. **Browser level** (`visualtest/wire_form_e2e_test.go`): fill the form,
   submit, assert the errors appear, fix the field, submit again. Note the
   settle wait: htmx wires swapped-in nodes during its ~20ms settle phase, so
   test clients (and only test clients — humans cannot click that fast) must
   wait before re-interacting with swapped content.

## Related

- `docs/transport-wiring.md` — the wire package, dialect mapping, parity table
- `docs/datastar-runtime-facts.md` — bundle-verified Datastar behavior
- `examples/demo/wire_demo.templ` — the live "Dual-transport form" section
- ADR-0036 — the common-subset contract behind `wire.Action`
