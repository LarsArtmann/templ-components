# ADR 0007: Self-Host HTMX as Default (CDN Opt-In)

**Date:** 2026-07-05
**Status:** Accepted (deferred to v1.0)
**Supersedes:** None

## Context

Currently, `layout.Base` injects htmx from a CDN (`cdn.jsdelivr.net`) by default.
Consumers who want a strict Content-Security-Policy must either:

1. Add the CDN to their `script-src` CSP allow-list
2. Set `HTMXVersion = ""` to suppress the CDN tag and self-host manually

Two consumers (SwettySwipper, Overview) reported CSP friction with the CDN
default. The CDN approach also adds a network dependency and a third-party
trust requirement that not all deployments can satisfy.

## Decision

**In v1.0**, change the default from CDN to self-hosted:

- `PageProps.HTMXVersion = ""` → no htmx script (current behavior)
- `PageProps.HTMXVersion = "2.0.10"` → CDN script (current behavior)
- `PageProps.HTMXSrc = "/static/htmx.org@2.0.10.min.js"` → **NEW default**: self-hosted

The library will **not** bundle the htmx JS file. Consumers download it once
and serve it from their own static directory. The `layout.Script` helper makes
this a one-liner:

```templ
@layout.Script(nonce, "/static/htmx.min.js", nil)
```

### Migration path (v0.x → v1.0)

1. Download htmx: `curl -o static/htmx.min.js https://cdn.jsdelivr.net/npm/htmx.org@2.0.10/dist/htmx.min.js`
2. Set `PageProps.HTMXVersion = ""` (suppress CDN)
3. Add `@layout.Script(nonce, "/static/htmx.min.js", nil)` to your page

### Why not bundle?

Bundling a JS file in a Go library is unconventional and adds binary size for
consumers who don't use HTMX. The library's principle is "zero JavaScript by
default" — self-hosting makes the JS dependency explicit.

## Consequences

- **Breaking change** — consumers upgrading to v1.0 must add the self-hosted script
- `PageProps.HTMXCDN` field is removed in v1.0 (no CDN default)
- `layout.sri.go` SRI hash validation is simplified (no CDN comparison needed)
- The Tailwind CSS CDN path (`CSSPath`) is already consumer-controlled — this
  brings htmx in line with the same pattern

## Timeline

Deferred to v1.0 (after `Validate() error` on all props and `internal/testutil/`
move). Target: after at least two external consumers are stable on v0.x.
