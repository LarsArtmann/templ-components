# Datastar vendored sources — bump & re-audit protocol

You vendored these files from `github.com/larsartmann/templ-components/datastar`
via `tc add`. They integrate with an EXTERNAL runtime: the pinned
`github.com/larsartmann/go-datastar/static` bundle (embedded JS). When you
bump that dependency, these vendored files must be re-audited — a version bump
is a contract re-audit, never a mechanical `go get`.

## Checklist (run on every go-datastar/static bump)

1. **Verify at source.** Hash the old and new bundle:
   `sha256sum "$(go env GOMODCACHE)/github.com/larsartmann/go-datastar/static@<ver>/datastar.js"`
2. **Diff the bundles** and read the changes around: event names
   (`datastar-patch-elements`, `datastar-patch-signals`), the lifecycle event
   (`datastar-fetch` with `detail.type` started/finished/error/retrying/
   retries-failed), retry defaults, and the fetch options destructure
   (`contentType`, `selector`, `retry`, `requestCancellation`).
3. **Re-check these vendored files against the diff:**
   - `live_region.templ` — `data-init` expression, retry/cancellation options,
     the busy-cue script's `datastar-fetch` listener and its `detail.type`
     values.
   - `sse_error_handling.templ` — same `datastar-fetch` lifecycle values.
   - `indicator.templ` — signal expression syntax (`data-indicator:<signal>`).
   - `sdk_script.templ` + `version.go` — CDN URL shape and the pinned version
     constant (the constant name must equal `static.Version`).
4. **Update the pin** in your go.mod, run your tests, and exercise the SSE
   stream end-to-end (wire format: keyed `data:` datalines, blank line
   terminates the event).

## Known runtime facts (quick reference)

- Only `datastar-patch-*` SSE events have registered handlers; there is NO
  `datastar-sse-error` event.
- CSP requires `'unsafe-eval'` (expressions compile via `Function()`).
- Clean stream EOF reconnects ONLY under `retry: 'always'`.
- Non-SSE HTML responses are patched via `Datastar-Selector` /
  `Datastar-Mode` response headers.

Upstream source of truth: `docs/datastar-runtime-facts.md` and
`docs/external-dependency-bumps.md` in the templ-components repo.
