# Datastar v1.0.2 Runtime Facts

Extracted from the pinned runtime bundle (`go-datastar/static` v0.2.0, Datastar
v1.0.2) during the 2026-08-21 SSE integration audit. Re-verified 2026-09-02 at
pin v0.4.0: the embedded `datastar.js` is byte-identical (same sha256) across
static v0.2.0–v0.4.0 — every fact below still applies. These contradict several
plausible-but-wrong assumptions that caused shipped bugs — re-verify against
the bundle if the version pin ever bumps. Enforcement lives in tests:
`examples/demo/sse_test.go` (wire format) and
`datastar/sse_error_handling_test.go` (lifecycle event names).

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

## Attribute syntax (v1.0.2)

- `data-<plugin>:<key>` (split on the first colon) is correct:
  `data-init` (no key), `data-on:click`, `data-indicator:<signal>`,
  `data-show="$<signal>"`.
- Bare `$` is the signals object itself and is **always truthy** — never emit
  `data-show="$"` (the empty-signal Indicator bug).
