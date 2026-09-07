# Multi-Step Forms (Stepper) — Both Transports

A wizard is a state machine on the server: each step is a fragment
(`StepIndicator` + that step's fields), and Next/Back are wired requests
that re-render both. The server validates per step and never trusts the
client's step index. This recipe composes existing components — no stepper
component needed.

## The pattern

```templ
// One region holds indicator + step body; the endpoint re-renders both.
<div id="wizard-region">
    @wizardStep(dialect, step, data)
</div>

templ wizardStep(dialect wire.Transport, step int, data WizardData) {
    @feedback.StepIndicator(feedback.StepIndicatorProps{
        Steps:       []string{"Account", "Profile", "Confirm"},
        CurrentStep: step,
    })
    @forms.Form(forms.FormProps{
        Action: "/wizard",
        Method: forms.FormPost,
        CSRFToken: data.CSRF,
        Wire: &wire.Action{
            Transport: dialect,
            Method:    wire.MethodPost,
            URL:       "/wizard/step",
            Target:    utils.Ternary(dialect == wire.TransportDatastar, "", "#wizard-region"),
        },
    }) {
        <input type="hidden" name="step" value={ strconv.Itoa(step) }/>
        if step == 0 {
            @forms.Input(forms.InputProps{Name: "email", Type: forms.InputEmail, Label: "Email"})
        } else if step == 1 {
            @forms.Input(forms.InputProps{Name: "name", Label: "Full name"})
        } else {
            <p class="text-sm text-gray-700 dark:text-gray-300">Review and submit.</p>
        }
        @display.Button(display.ButtonProps{Text: "Next", Type: display.ButtonHTMLSubmit})
    }
}
```

## The endpoint

```go
mux.HandleFunc("POST /wizard/step", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    if err := r.ParseForm(); err != nil {
        http.Error(w, "invalid form", http.StatusBadRequest)
        return
    }

    step, _ := strconv.Atoi(r.PostFormValue("step"))
    data := loadAndValidate(step, r) // server owns state; validate THIS step

    if errCount := data.Errors(); errCount > 0 {
        // Re-render the SAME step with errors — 200 OK fragments
        // (see docs/recipes/server-side-validation.md).
        render(w, wizardStep(dialectOf(r), step, data))
        return
    }

    render(w, wizardStep(dialectOf(r), step+1, data))
})
```

Both transports work unchanged: htmx swaps into `#wizard-region` via
`hx-target`; Datastar needs response-header targeting — wrap the handler in
`wire.Handler(wire.PatchTarget{Selector: "#wizard-region", Mode:
wire.PatchModeInner}, ...)` or set `wire.Action.Selector` to target
client-side.

## Rules that keep wizards honest

- **Server owns the step.** Never advance based on a client-sent step
  number without validating that the current step's required fields passed.
  The hidden `step` input is a hint, not authorization.
- **Validate per step, save at the end** (or persist partial progress
  server-side keyed by a session ID — the pattern is identical).
- **Errors never advance the wizard** and travel as 200 OK fragments —
  both runtimes ignore non-2xx responses by default.
- **Back is a request too** (re-render the previous step from server
  state), not client-side show/hide — otherwise refresh loses the wizard.

## Demo

The wire demo's "Multi-step wizard" card runs this exact pattern against
`/api/wire/wizard` on both transports.
