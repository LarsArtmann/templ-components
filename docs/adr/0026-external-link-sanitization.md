# ADR-0026: ExternalLink URL Sanitization

**Status:** Accepted (2026-07-30)
**Decider:** Lars Artmann

## Context

External links need `target="_blank" rel="noopener noreferrer"` to prevent
tabnabbing and back-history manipulation. But the href itself is also a
security concern — `javascript:`, `data:`, and `vbscript:` schemes can
execute arbitrary code.

templ provides `templ.SafeURL` as a type for trusted URLs. The question is
whether to use it for the ExternalLink href.

## Decision

**Use plain `string` for the href, NOT `templ.SafeURL`.**

templ's `SafeURL` is a type assertion ("I promise this URL is safe"), not a
validator. When you pass a string to `href={ ... }` in templ, the built-in
URL sanitizer runs and blocks dangerous schemes by rewriting them to
`about:invalid`.

By using a plain string, the sanitizer runs on every href. By using `SafeURL`,
the sanitizer is bypassed.

```go
// CORRECT: sanitizer runs, javascript: is blocked
href={ props.Href }        // props.Href is string

// DANGEROUS: sanitizer bypassed, any scheme allowed
href={ templ.SafeURL(props.Href) }
```

## Consequences

- Every href passed to ExternalLink is sanitized — no way to bypass.
- Consumers cannot use `SafeURL` to pre-validate and skip sanitization.
- The `↗` arrow icon is `aria-hidden` to avoid screen-reader noise.
- `rel="noopener noreferrer"` is always set, even for `mailto:` links (harmless).
