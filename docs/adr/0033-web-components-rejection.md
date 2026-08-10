# ADR 0033: Web Components Rejection — Use the Platform, Not the Stack

## Date

2026-08-10

## Status

Accepted

## Context

"Web Components" is invoked in three separate research documents (`docs/research/what-we-are-missing.md`
§2.4, `docs/planning/2026-07-21_03-54_NEXT-LEVEL-PLATFORM-FIRST-ROADMAP.md`,
`docs/planning/2026-07-12_19-00_MODERN-BROWSER-INTEGRATION.md`) as a candidate capability.
Each time it was deferred or skipped with a one-line rationale. Because the question
keeps returning, this ADR records the decision permanently so future sessions stop
re-litigating it.

The term "Web Components" conflates two distinct ideas that must be separated:

1. **The Web Components technology stack** — Custom Elements (`customElements.define`),
   Shadow DOM (`attachShadow`, Declarative Shadow DOM), and HTML Templates (`<template>`).
2. **The "use the platform" philosophy** — reaching for native browser APIs instead of
   custom JavaScript or framework abstractions.

templ-components wholeheartedly adopts **(2)** and explicitly rejects **(1)**. This ADR
explains why the technology stack is the wrong tool for this library's thesis, while the
philosophy is already the library's governing principle.

## Decision

**Do NOT adopt the Web Components technology stack (Custom Elements, Shadow DOM, or
Declarative Shadow DOM) in templ-components, opt-in or otherwise, for the foreseeable
lifecycle of the v1.x / v2.x line.**

This is a permanent architectural boundary, not a "deferred until Baseline" item. It is
added to the ROADMAP's "Explicitly NOT Planned" table alongside the React/Vue/Svelte
wrapper rejection.

## Rationale

### 1. Shadow DOM breaks the theming model

The library's entire styling system is built on **Tailwind utility classes that consumers
override via `@theme` custom properties** (`@theme { --color-blue-600: #custom; }`). This
is a light-DOM, global-stylesheet model. Shadow DOM creates a hard encapsulation boundary:
stylesheets and utility classes from the host document **do not cross into the shadow
tree**. A `<tc-card>` with a shadow root would not receive the consumer's Tailwind classes.

The workarounds all re-architect the theming model:

| Workaround                          | Why it fails here                                                                 |
| ----------------------------------- | --------------------------------------------------------------------------------- |
| CSS custom properties (do pierce)   | Forces a total rewrite: every utility class becomes a token reference. Abandons Tailwind. |
| `::part()` / `::theme()`            | Per-shadow-root, awkward, and each part must be manually exported. Anti-ergonomic for consumers. |
| Constructable Stylesheets           | Requires JavaScript to inject the Tailwind output into each shadow root — violates the zero-JS principle. |
| `@adoptedstylesheets` on the server | No SSR path exists; DSD `<style>` blocks can't reference the consumer's compiled Tailwind output. |

None of these preserve the library's core value: **write Go, get styled HTML, override
colors without touching Go or JS.**

### 2. The distribution problem Web Components solve does not exist here

Web Components exist to distribute **framework-agnostic widgets** that work in React, Vue,
Angular, or plain HTML. templ-components distributes **Go source code** to **Go developers**
who vendor it and compile with their own Tailwind. There is no cross-framework consumer to
serve. A Web Component wrapper would solve a problem nobody in the target audience has.

### 3. Custom Elements require JavaScript

Custom Elements upgrade via `connectedCallback` — a JavaScript lifecycle hook. The library
prides itself on a large surface of **zero-JavaScript components**: native `<dialog>`
(Modal, Drawer), `<details>` (Accordion, CollapsibleSection), scroll-snap (Carousel),
CSS-only Tooltip, native Popover API (Dropdown, ContextMenu), `field-sizing` (Textarea),
`accent-color` (form controls). A Custom Element trades this principle for marginal
encapsulation that the library does not need.

### 4. HTMX interaction is hostile to shadow trees

HTMX swaps HTML fragments over the wire into the light DOM. A Custom Element renders its
shadow tree client-side on upgrade. After an HTMX swap that replaces a custom element's
light-DOM children, the shadow tree and the server-rendered content **desync** — the
shadow tree must be re-rendered or invalidated, reintroducing exactly the client-side
hydration complexity the library was designed to eliminate. This is a fundamental model
conflict, not an integration detail.

### 5. Declarative Shadow DOM (DSD) does not rescue the case

DSD (Baseline 2024) server-renders a shadow root via `<template shadowrootmode="open">`,
removing the JS-upgrade requirement. But it does not remove boundaries 1, 2, or 4 above:
the shadow tree still can't receive the consumer's Tailwind utility classes, the
distribution problem is unchanged, and HTMX swaps still desync. DSD solves the *hydration*
problem that the library doesn't have, while leaving the *theming* problem intact.

## What the library DOES adopt: the "use the platform" philosophy

The library achieves every goal that motivates teams to reach for Web Components —
encapsulated, accessible, dependency-light components — through **native platform APIs in
the light DOM**:

| Goal                         | Web Components approach          | templ-components approach                              |
| ---------------------------- | -------------------------------- | ------------------------------------------------------ |
| Modal/dialog                 | Custom Element + shadow root     | Native `<dialog>` + `showModal()` (ADR-0014)           |
| Floating menus / tooltips    | Custom Element positioning       | Native Popover API + `popovertarget` (ADR-0017)        |
| Collapsible regions          | Custom Element open/close        | Native `<details>`/`<summary>` (ADR-0027)              |
| Slideshows                   | Custom Element transform logic   | Native CSS scroll-snap (zero JS for touch/drag)        |
| Component-scoped reflow      | Shadow DOM containment           | CSS Container Queries `@container` (ADR-0018)          |
| Theme encapsulation          | Shadow DOM style boundary        | Tailwind `@theme` tokens + semantic alias layer (ADR-0008) |
| Style encapsulation          | Shadow boundary                  | BEM-free utility classes + tailwind-merge conflict resolution |

The result: the library ships ~690 lines of JavaScript total across 112 components — and
most components ship **zero** JS. This is the "use the platform" outcome that Web
Components promise, achieved without the shadow-boundary cost.

## The one narrow exception (explicitly out of scope)

A future **opt-in, separate Go module** could expose templ-components output as Web
Components for non-Go consumers (React/Vue/Astro apps). This would be a **wrapper project**,
not part of this library, and would own the theming-model compromises (token-only theming,
JS requirement, HTMX desync) on its own terms. It is not planned and would only begin if
concrete consumer demand materializes. It is explicitly **not** a v1.x or v2.x item.

## Consequences

**Positive:**

- The question is closed. Future sessions and contributors find a binding decision instead
  of re-evaluating from scratch.
- The ROADMAP's "Explicitly NOT Planned" table gains its missing Web Components entry.
- The library's "use the platform" identity is reinforced as a deliberate stance, not an
  oversight.

**Negative:**

- Consumers who specifically need framework-agnostic widget distribution get no path from
  this library. They must build their own wrapper (per "The one narrow exception" above).

## References

- [ADR-0008: Semantic tokens](0008-semantic-tokens.md) — the theming model WC would break
- [ADR-0014: `<dialog>` migration](0014-dialog-migration.md) — "use the platform" exemplar
- [ADR-0017: Popover API migration](0017-popover-api-migration.md) — "use the platform" exemplar
- [ADR-0018: Container-Query-Native Contract](0018-container-query-native-contract.md) — light-DOM reflow
- `docs/research/modern-browser-capabilities.md` — the native-API inventory
- `docs/research/what-we-are-missing.md` §2.4 — the prior DSD analysis this ADR supersedes
