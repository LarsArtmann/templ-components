# ADR 0021: Headless Variants (Decision: Defer)

## Date

2026-07-21

## Status

**Deferred indefinitely.** Spike concluded that the `Unstyled bool` flag
approach does not work cleanly for this library's design.

## Context

Phase 5.4 of the platform-first roadmap asked: should templ-components ship
"headless" variants — components that emit semantics + ARIA but NO Tailwind
classes, leaving styling entirely to the consumer? The Radix UI / Headless UI
pattern is well-known in React. The plan suggested evaluating a `Unstyled
bool` flag per component.

## Decision

**Defer.** Three options were considered; none were a clear win.

### Option A: `Unstyled bool` on every component (spike target)

Add `Unstyled bool` to `BaseProps`. When true, the component emits its HTML
structure + ARIA but skips all Tailwind class strings.

**Why rejected:**

1. Every component already accepts `Class` + `Attrs` for class override. The
   `utils.Class()` merge means consumer classes already win against the
   defaults. The "unstyled" use case is largely served by passing
   `Class: "my-classes"` and overriding.
2. Components that compose other components (e.g. `Dropdown` uses `dropdownItemLink`
   sub-template) would need the flag threaded through every layer — invasive.
3. The library's value proposition is the Tailwind styling. Stripping it
   leaves bare HTML that consumers could write themselves in 30 seconds.
4. CSP-safe inline JS (Tooltip, Dropdown keyboard nav) is styled by class —
   stripping classes breaks the JS visual feedback loop.
5. Golden files assume the styling is present — every test would need an
   "unstyled golden" variant.

### Option B: Separate `*-headless` package

A parallel package `display/headless/` with unstyled variants. Each headless
component wraps the styled one and strips classes.

**Why rejected:**

- Doubles the API surface (98 → 196 components)
- Maintainers must keep two render paths in sync on every change
- The styled components already have slots — consumers can compose their own
  DOM and pass it via the slot

### Option C: Document the existing override pattern

Document that consumers wanting a "headless" experience should pass
`Class: ""` (which suppresses the library classes via tailwind-merge-go's
"empty string wins" behavior) plus their own classes via `Attrs` or `Class`.

**Accepted.** This is the existing pattern, costs zero, and works today. The
recipe is now documented in `docs/theming.md` under "Component-level Class
override."

## Consequences

- No code changes. No new APIs.
- The "headless variants" item in TODO #41 is closed as "won't do" with this
  ADR as the rationale.
- Consumers who need Radix-style headless primitives should reach for a
  different library. templ-components is opinionated about Tailwind.

## References

- TODO #41 — closed by this ADR
- [`docs/theming.md`](../theming.md) — documents the existing override pattern
- [ADR-0019: recipes package](0019-recipes-package.md) — slot-based composition
  as an alternative to headless variants
