# ADR 0047: Context-Injected Defaults — props > context > zero

## Date

2026-10-14

## Status

Accepted. Executes idea #1 of the 2026-10-14 big-picture review
([docs/reviews/2026-10-14_50-big-picture-ideas.md](../reviews/2026-10-14_50-big-picture-ideas.md));
plan lane M3/M4 in
[docs/planning/2026-10-14_18-52_FIFTY-IDEAS-PARETO-MASTER-PLAN.md](../planning/2026-10-14_18-52_FIFTY-IDEAS-PARETO-MASTER-PLAN.md).

## Context

Every script-emitting component carries a `Nonce` prop (CSP), and future axes
(locale, direction) want the same shape. The value is almost always constant
per REQUEST — the middleware that minted it has it — but props structs are
filled deep in the render tree, far from that middleware. Today every call
site must thread the nonce through by hand:

```go
@display.Button(display.ButtonProps{Nonce: nonce, ...})   // at every level
```

That is the library's single largest per-props papercut: mechanical, easy to
forget (producing `nonce=""` — the dead-script class pinned by
`TestNoEmptyNonceAttribute`), and it multiplies with every new
request-scoped prop.

One repo rule constrains any solution: **a struct-literal zero value MUST
equal the documented default** (consumer audit lesson — consumers use
`XxxProps{...}` literals). Whatever carries request-scoped defaults may only
FILL values the caller left empty. It must never reinterpret a literal.

## Decision

**Precedence: explicit props > context > zero value.** A component resolves a
request-scoped value through ONE helper that applies that order:

```go
nonce := utils.ResolveNonce(ctx, props.Nonce)
```

`utils/context.go` ships:

- `ContextWithNonce(ctx, nonce) context.Context` — middleware-side setter.
- `NonceFromContext(ctx) string` — typed getter, zero-value fallback `""`
  (never panics, never a sentinel error).
- `ResolveNonce(ctx, explicit string) string` — the precedence rule as code:
  explicit wins; context fills empty; both empty yields `""`, and the
  component's existing omit-empty behavior applies unchanged.

Keys are unexported zero-size structs (`nonceKey struct{}`); values are plain
strings. Future request-scoped axes (locale, direction) get their own key +
`Resolve*` helper — the precedence rule generalizes, the per-axis helper keeps
types honest.

Adoption is per-component and opt-in: a component starts resolving through
the helper in the same commit that re-baselines its goldens (a context-free
render is byte-identical, so goldens only gain NEW context-nonce cases, never
change existing ones).

## Alternatives considered

- **Global default variable** (`utils.DefaultNonce = …`): process-global
  mutable state; breaks multi-tenant servers and tests that render with
  different tenants in parallel. Rejected.
- **Context overrides props**: breaks the struct-literal rule — the same
  literal would render differently depending on invisible request state, the
  un-reasonable-about trap this library exists to prevent. Rejected.
- **Do nothing** (keep threading): taxes every call site forever and keeps
  producing the empty-nonce dead-script class. Rejected.

## Consequences

- Middleware wraps the handler once: `utils.ContextWithNonce(r.Context(), nonce)`;
  every adopted component picks the nonce up automatically, and an explicit
  prop still wins (test-proven: `TestNoncePrecedence`).
- The empty-nonce contract is untouched: `ResolveNonce` returning `""` renders
  NO nonce attribute (the omit-empty rule), never `nonce=""`.
- Golden safety: `ResolveNonce(ctx, "")` with a bare `context.Background()` is
  byte-identical to the old `props.Nonce == ""` path.
- Each adoption is a visible, golden-backed diff — no silent behavior change
  can enter through this door.
