# Datastar Runtime Facts

Extracted from the pinned runtime bundle (`go-datastar/static` v0.2.0, Datastar
v1.0.2) during the 2026-08-21 SSE integration audit. Re-verified 2026-09-02 at
pin v0.4.0: the embedded `datastar.js` is byte-identical (same sha256) across
static v0.2.0–v0.4.0 — every fact below still applies. These contradict several
plausible-but-wrong assumptions that caused shipped bugs — re-verify against
the bundle if the version pin ever bumps. Enforcement lives in tests:
`examples/demo/sse_test.go` (wire format) and
`datastar/sse_error_handling_test.go` (lifecycle event names).

## v1.0.3 re-audit (2026-09-05, `go-datastar/static` v0.5.0)

The bundle CHANGED (sha256 `4df1f98a…` → `5d6b7794…`, 56330 → 33538 bytes —
upstream minification refactor). Every fact below was re-verified against the
new bundle unless marked otherwise:

- **Unchanged**: only `datastar-patch-elements`/`datastar-patch-signals` have
  registered handlers; `datastar-sse-error` and the pre-v1.0 `datastar-merge-*`
  names remain absent; `PatchElementsExpectedSelector` /
  `PatchElementsNoTargetsFound` enforcement intact; attribute syntax unchanged.
- **Unchanged (machinery verified)**: retry defaults in the bundle's fetch
  options destructure (`retryInterval=1000`, `retryScaler=2`, `retryMaxWait=30000`,
  `retryMaxCount=10`, `retry='auto'`). The reconnect matrix itself
  (which mode reconnects on clean EOF) was verified behaviorally on v1.0.2 and
  the machinery is byte-present in v1.0.3; no behavioral counter-evidence found.
- **ADOPTED (2026-09-07) — fetch actions accept a client-side `selector`
  option**: the v1.0.3 fetch options destructure includes `selector`
  (v1.0.2 had none). The v1.0.2-era fact "fetch actions accept no target
  option" is therefore OUTDATED. Response dispatch decodes (from the
  bundle's response handler): for a non-SSE `text/html` response the
  runtime reads `datastar-selector` / `datastar-mode` / … response
  headers, then **the fetch options override them field-by-field** — a
  string `selector` option replaces the header value (`mode`,
  `namespace`, `useViewTransition` behave the same). Under
  `contentType: 'form'` the `selector` option additionally picks which
  form serializes (`querySelector(sel)` over `closest("form")`).
  Consumed by `wire.Action.Selector` (ADR-0038): renders Datastar-only,
  empty keeps response-header targeting authoritative.
- **NEW (2026-09-07) — fetch actions accept `contentType: 'form'`** for
  whole-form serialization (consumed by `wire.ContentTypeForm` /
  `forms.FormProps.Wire`):
  - Default is `contentType: 'json'` (signals object as JSON body; on
    body-less GET requests the JSON rides in a `?datastar=<json>` query
    param). Any other value throws `FetchInvalidContentType`.
  - Under `'form'` the action serializes the action element's
    `closest("form")` (or the form matched by the `selector` option); no
    enclosing form throws `FetchFormNotFound`.
  - **HTML5 constraint validation gates the request**: unless the form has
    `novalidate`, an invalid form calls `checkValidity` + `reportValidity`
    and the fetch never fires.
  - The **submitter button's** `name`/`value` is appended (the
    `SubmitEvent.submitter` when the action element is the form, else an
    `input[type=submit]`/`button[type=submit]` the browser would submit).
  - `enctype="multipart/form-data"` sends a `FormData` body (file uploads
    work); anything else sends `application/x-www-form-urlencoded`.
  - GET requests carry the form fields as query parameters instead of a body.
  - The `data-on` plugin **auto-calls `preventDefault()` when the element is
    a `HTMLFormElement` and the event is `submit`** — so
    `data-on:submit="@post('/x', {contentType: 'form'})"` on a `<form>`
    suppresses the native full-page submission with no explicit modifier.
  - **NEW (2026-09-07, decoded from the pinned bundle) — `data-on` modifier
    spelling**: the attribute-name parser splits the name on `__` (the first
    segment is `plugin:event`, each further segment is one modifier group),
    and each group splits on `.` into `name.arg1.arg2…`. A debounced search
    input is therefore `data-on:input__debounce.300ms`. Durations parse as
    `500ms` / `2s` / bare `300` (first arg wins). Verified flags: `debounce`
    (trailing by default; `leading` opt-in, `notrailing` opt-out),
    `throttle` (leading by default; `trailing` opt-in, `noleading` opt-out),
    `delay.<dur>`, `window`, `document`, `capture`, `passive`, `once`,
    `outside`, `prevent`, `stop`, `viewtransition`, and `case.<x>`
    (event-name casing; `on` defaults to kebab). Trigger syntax stays
    dialect-specific per ADR-0036's scope rule regardless.
  - **NEW (2026-09-07, e2e-proven) — fetch error status dispatches a
    lifecycle error event**: the fetch options include an error hook that
    fires when `response.status >= 400` (dispatches the `datastar-fetch`
    error with the status). Combined with htmx 2.0.10's default
    `responseHandling` (`{code: "[45]..", swap: false, error: true}` — 4xx/5xx
    responses DO NOT swap), server-side validation errors must travel as
    **200 OK + error fragment** for zero-config parity on both runtimes. The
    422 variant needs client configuration in each runtime; see
    `docs/recipes/server-side-validation.md`.
  - **NEW (2026-09-07, e2e-proven) — patched-in `data-on:submit` forms keep
    working**: the runtime's MutationObserver initializes `data-*` attributes
    on nodes patched in by `datastar-patch-elements`, so a re-rendered form
    (validation errors round-trip) submits again without any manual
    re-initialization. Browser-proven by
    `visualtest/wire_form_e2e_test.go` (`TestWireE2EFormValidationRoundTrip`).

Full audit context: `docs/research/2026-08-21_go-sse-go-datastar-deep-dive.html`.

## Wire format (server → client SSE)

- **SSE event types**: only `datastar-patch-elements` and
  `datastar-patch-signals` have registered handlers. The pre-v1.0 names
  (`datastar-merge-fragments`, `datastar-merge-signals`) are **silently
  ignored**.
- **Datalines are keyed per line**: the client splits each `data:` line on the
  FIRST space (`selector #x`, `mode inner`, `elements <html>`). Values for a
  repeated key are joined with `\n`, so multi-line HTML must repeat the
  `elements` prefix on every line.
- **Event termination**: real newlines per line + a blank line terminate the
  event. A literal `\n` (backslash-n text) in the format string means the
  browser **never dispatches** the event — this exact bug shipped.
- **Default `outer` mode matches by id**: without a `selector` dataline,
  incoming root elements are looked up via `getElementById`; id-less fragments
  are dropped with a `PatchElementsNoTargetsFound` console warning.
- **Non-default modes REQUIRE a selector** — otherwise the runtime throws
  `PatchElementsExpectedSelector`.
- **Non-SSE HTML responses** are patched via response headers: the body is the
  elements payload, `Datastar-Selector` picks the target, `Datastar-Mode` the
  merge mode.

## Lifecycle observability

- Everything flows through the document-level `datastar-fetch` CustomEvent:
  `detail = {type, el, argsRaw}`, type ∈ `started`, `finished`, `error`
  (HTTP ≥ 400; status in `argsRaw.status`), `retrying`, `retries-failed`
  (reconnects exhausted; defaults: 10 retries, 1s ×2 exponential backoff,
  30s cap).
- **There is NO `datastar-sse-error` event** — listening for it is dead code.
- Reconnection matrix for `@get`/`@post`/... actions (verified in the bundle's
  fetch plugin): clean stream EOF reconnects **only** under `retry: 'always'`;
  HTTP ≥ 400 under `'always'` or `'error'`; thrown network errors under every
  mode. The failure counter resets on every successful (200) connect, so
  `'always'` self-heals indefinitely across individual server restarts
  (defaults: 10 retries, 1s ×2 exponential backoff, 30s cap).
  `datastar.LiveRegionProps.Retry` maps to this argument
  (`RetryAlways` → `@get(url, {retry: 'always'})`).
- Enforced by `datastar.TestPinnedRuntimeBundleContract` (bundle byte-content
  guard) — a pin bump that renames any of these tokens fails CI.

## CSP

- `data-*` expressions compile via the `Function()` constructor → CSP
  classifies it as eval → `script-src` needs `'unsafe-eval'`.
- `'unsafe-inline'` is NOT needed (external module script + nonced inline
  scripts).

## Attribute syntax

- `data-<plugin>:<key>` (split on the first colon) is correct:
  `data-init` (no key), `data-on:click`, `data-indicator:<signal>`,
  `data-show="$<signal>"`.
- Bare `$` is the signals object itself and is **always truthy** — never emit
  `data-show="$"` (the empty-signal Indicator bug).
