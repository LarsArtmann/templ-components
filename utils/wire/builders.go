package wire

// Constructors and fluent builders for the Action exchange spec. They exist
// so call sites stop hand-rolling the same struct literals and component
// packages stop duplicating the same default-application logic: every With*
// method is a VALUE receiver that returns a modified COPY — a consumer's
// *Action is never mutated (the same non-mutation contract
// kanbanWireAttributes/formWireAttributes pinned by test), and chains read
// left to right:
//
//	wire.Post("/api/save").WithTarget("#out").WithDebounce(300)
//
// The zero value stays the source of truth: an unset field keeps rendering
// as "dialect default", so builders only ever SET what the chain says.

// Get wires a GET exchange — the read default (fragment loads, filters,
// searches).
func Get(url string) *Action {
	return &Action{URL: url} //nolint:exhaustruct_v5 // the zero value IS the read default
}

// Post wires a POST exchange. Mutations MUST set an explicit method — the
// zero-value Action renders GET in BOTH dialects.
func Post(url string) *Action {
	return &Action{URL: url, Method: MethodPost} //nolint:exhaustruct_v5 // zero value stays dialect-default
}

// Put wires a PUT exchange.
func Put(url string) *Action {
	return &Action{URL: url, Method: MethodPut} //nolint:exhaustruct_v5 // zero value stays dialect-default
}

// Patch wires a PATCH exchange.
func Patch(url string) *Action {
	return &Action{URL: url, Method: MethodPatch} //nolint:exhaustruct_v5 // zero value stays dialect-default
}

// Delete wires a DELETE exchange.
func Delete(url string) *Action {
	return &Action{URL: url, Method: MethodDelete} //nolint:exhaustruct_v5 // zero value stays dialect-default
}

// WithEvent sets the triggering DOM event. Under htmx it renders as the
// hx-trigger event token; under Datastar as the data-on:<event> key.
func (a Action) WithEvent(e Event) Action {
	a.Event = e

	return a
}

// WithTransport selects the attribute dialect (zero value stays htmx).
func (a Action) WithTransport(t Transport) Action {
	a.Transport = t

	return a
}

// WithMethod sets the HTTP verb explicitly (the constructors cover the
// common cases; this escapes for dynamic verbs).
func (a Action) WithMethod(m Method) Action {
	a.Method = m

	return a
}

// WithTarget sets the htmx patch region (hx-target). Datastar targeting is
// response-driven (see PatchTarget); use WithSelector for the client-side
// Datastar twin.
func (a Action) WithTarget(target string) Action {
	a.Target = target

	return a
}

// WithSelector selects which form serializes under Datastar ContentTypeForm
// ({selector: '…'}) — the runtime resolves the enclosing form via that
// selector instead of the action element's closest form. It NEVER targets
// patches: the pinned runtime reads patch targeting exclusively from the
// datastar-selector response header (or the fragment root's id) — use
// wire.Handler(PatchTarget) server-side instead.
func (a Action) WithSelector(selector string) Action {
	a.Selector = selector

	return a
}

// WithSwap sets the region-merge style (the shared PatchMode vocabulary).
// Transport-asymmetric like Target: htmx renders hx-swap client-side; under
// Datastar the merge mode is response-driven (wire.Handler's PatchTarget.Mode)
// and Swap renders nothing there.
func (a Action) WithSwap(mode PatchMode) Action {
	a.Swap = mode

	return a
}

// WithContentType sets how field values travel (ContentTypeForm makes a
// Datastar action serialize the enclosing form).
func (a Action) WithContentType(c ContentType) Action {
	a.ContentType = c

	return a
}

// WithDebounce sets the debounce window (delay:<n>ms trigger modifier under
// htmx, __debounce.<n>ms under Datastar). Needs an explicit event.
func (a Action) WithDebounce(ms int) Action {
	a.DebounceMS = ms

	return a
}

// WithThrottle sets the throttle window (throttle:<n>ms under htmx,
// __throttle.<n>ms under Datastar). Needs an explicit event. Throttle and
// debounce answer different questions — throttle bounds the RATE (fire at
// most once per window, trailing call last), debounce waits for QUIET — so
// setting both is allowed but pointless under Datastar (delay → throttle
// pipeline) and unsupported by htmx (delay+throttle together). Prefer one.
func (a Action) WithThrottle(ms int) Action {
	a.ThrottleMS = ms

	return a
}

// WithInterval sets the polling duration (ADR-0043).
func (a Action) WithInterval(every string) Action {
	a.Interval = every

	return a
}

// WithReveal sets the viewport-reveal trigger (ADR-0043).
func (a Action) WithReveal(reveal Reveal) Action {
	a.Reveal = &reveal

	return a
}

// WithPreventDefault renders the Datastar __prevent modifier (htmx ignores
// it — it intercepts hx-* clicks natively). See Action.PreventDefault.
func (a Action) WithPreventDefault() Action {
	a.PreventDefault = true

	return a
}

// WithFormDefaults applies the host-component defaults shared by every
// component that wires INSIDE a form (forms.Form, FilterInput,
// FilterDropdown, KanbanBoard's hidden move form): an unspecified Event
// becomes the host's natural submit-family event, and an unspecified
// ContentType becomes form encoding so the form's fields travel under
// Datastar too. It replaces the per-package copy-and-default helpers — the
// copy semantics (never mutate the consumer's action) are the method's whole
// point.
func (a Action) WithFormDefaults(event Event) Action {
	if a.Event == EventUnspecified {
		a.Event = event
	}

	if a.ContentType == ContentTypeUnspecified {
		a.ContentType = ContentTypeForm
	}

	return a
}
