# ADR-0030: Datastar Integration Strategy

**Status:** Accepted (2026-08-02)
**Decider:** Lars Artmann

## Context

[Datastar](https://data-star.dev/) (v1.0.2) is a ~12 KiB zero-dependency frontend
framework that unifies HTMX's backend reactivity with Alpine.js's frontend
reactivity via Server-Sent Events (SSE) and reactive signals. It is the natural
evolution of the hypermedia approach this library champions.

The question: should templ-components adopt Datastar, and if so, how?

The library currently has three interactivity tiers:
1. Native HTML (`<dialog>`, `<details>`, Popover API) — zero JS
2. HTMX (`hx-*` attributes) — zero JS runtime dependency
3. Singleton-guard inline JS (15 scripts across interactive components) — zero
   dependency, but ~500 lines of hand-maintained event-delegation JS

Datastar offers a fourth tier: a single ~12 KiB runtime that replaces both HTMX
and the singleton-guard JS with declarative reactive signals and SSE streaming.

The tension: this library's core selling point is **zero JS runtime
dependencies** (ADR 0005, ADR 0007). Making Datastar a hard dependency would
destroy that guarantee and break every consumer.

See `docs/research/datastar-integration-analysis.md` for the full deep-research
analysis.

## Decision

**Adopt Datastar as an opt-in complement — not a replacement for HTMX.**

1. **HTMX remains the default** interactivity layer. Zero-dependency, HATEOAS-first.

2. **A new `datastar` package** provides first-class Datastar support with **zero
   new library dependencies**. It mirrors the existing `htmx` package pattern:
   emit `data-*` attributes and inject the runtime `<script>`, without importing
   the `datastar-go` SDK. Consumers who want SSE streaming add `datastar-go` to
   **their** go.mod, not ours.

3. **SSE-powered components** (`LiveRegion`) unlock real-time capabilities that
   HTMX polling cannot match: push-based updates, zero idle requests,
   sub-second latency, signal-only patches (no HTML round-trip).

4. **The existing 104 components remain framework-agnostic.** ~80% are pure
   server-rendered HTML + Tailwind and work unchanged in any interactivity model.
   Only the 8 `htmx`-package components have HTMX-specific coupling; the recipe
   documents the Datastar equivalents.

5. **Consumers choose** per-page or app-wide: HTMX (default), Datastar (opt-in),
   both, or neither (native HTML only).

## Why Not Full Migration (Tier D — Rejected)

- Breaks the zero-dependency guarantee (the library's #1 selling point)
- Breaks every `hx-*` consumer and all 8 `htmx` package components
- Adds ~12 KiB JS to every page even for static content
- Requires SSE infrastructure not available in all deployments (serverless)
- Datastar's `data-*` expressions are JavaScript, not pure hypermedia controls

## Consequences

- **New package:** `datastar/` with `SDKScript`, `LiveRegion`, `Indicator`, and
  typed `data-*` attribute helpers. Same conventions as all other packages
  (BaseProps, typed enums with IsValid, golden tests, dark-mode compliance).
- **No go.mod changes.** The library stays at 3 production dependencies (templ,
  tailwind-merge-go, go-error-family).
- **Documentation burden:** the JS decision ladder (rung 7) and a new recipe
  (`docs/recipes/datastar-integration.md`) guide consumers through the choice.
- **Two interactivity models** coexist: consumers who use both HTMX and Datastar
  on different pages face no conflict (distinct attribute namespaces: `hx-*` vs
  `data-*`).
- **Future reactive variants:** Combobox, TagsInput, and multi-step forms may get
  Datastar-native alternatives in the `datastar` package — not replacements for
  the `forms.*` equivalents.
