# Container Query Leveraging Strategy

> **How to get the most out of CSS Container Queries in templ-components.**
> This document maps the full landscape: what's shipped, what's new, and what's next.
>
> **Date:** 2026-08-10 · **Author:** Architecture review

---

## Executive summary

templ-components has **the deepest container-query integration of any Go UI library**:
8 opt-in `ContainerAware` components, a binding contract (ADR-0018), an enforcement
scanner (`TestContainerQueryCompliance`), and now **fluid typography via container query
units** (`cqi`). This document consolidates the strategy and identifies the remaining
high-value, on-thesis work.

Container Queries are the **single best "use the platform" capability** for this library's
thesis (see ADR-0033): they are pure CSS, zero JavaScript, zero dependencies, Baseline
2023, and they solve the exact problem server-rendered component libraries face — a
component's correct layout depends on where it's *placed*, not on the *viewport*.

---

## Part 1: What's shipped (the foundation)

### The 8 container-aware components

| Component                   | Flag                  | What adapts                                              |
| --------------------------- | --------------------- | -------------------------------------------------------- |
| `display.Grid`              | `ContainerResponsive` | Column count (1→2→3→N)                                   |
| `display.Card`              | `ContainerAware`      | Padding (compact below `@sm:`)                           |
| `display.DefinitionGrid`    | `ContainerAware`      | Term-detail card grid column count                       |
| `navigation.Nav`            | `ContainerAware`      | Collapse to hamburger below `@sm:`                       |
| `navigation.Pagination`     | `ContainerAware`      | Mobile prev/next vs full page numbers                    |
| `layout.Split`              | `ContainerAware`      | 2-col main+aside collapses to stacked below `@md:`       |
| `forms.Form`                | `ContainerAware`      | Grid layout label/value columns below `@sm:`             |
| `feedback.SkeletonCardGrid` | `ContainerAware`      | Loading skeleton grid matches `Grid.ContainerResponsive` |

All follow the ADR-0018 contract: opt-in (default off), emit `<div class="@container">`
wrapper (or `@container` on root for Pagination), swap `sm:`/`md:`/`lg:` for
`@sm:`/`@md:`/`@lg:`. Byte-identical when off.

### Fluid typography via container query units (NEW)

Six `.tc-fluid-*` utility classes in `templates/custom.css` size text with `cqi` (container
inline-size units) + `clamp()`. Text scales smoothly with its container — the typography
analog of the container-aware components. See the [Fluid Typography recipe](../recipes/fluid-typography.md).

### Enforcement

- `utils.TestContainerQueryCompliance` — scans every `.templ` file for viewport breakpoints
  (`sm:`/`md:`/`lg:`) on components without a `ContainerAware` escape hatch. Catches
  regressions where someone adds responsive behavior without the container variant.

---

## Part 2: Container Query Length Units (`cqi`, `cqw`, `cqh`, ...)

**Baseline 2023** — safe everywhere. These units resolve against the nearest `@container`
ancestor instead of the viewport. The library now uses them for fluid typography, but they
have broader application:

| Use case                      | Pattern                              | Status       |
| ----------------------------- | ------------------------------------ | ------------ |
| Fluid headings / metrics      | `clamp(min, Ncqi + base, max)`       | **Shipped** (`.tc-fluid-*`) |
| Fluid spacing (gap, padding)  | `clamp(0.5rem, 2cqi, 1.5rem)`        | Consumer CSS |
| Fluid icon sizing             | `width: clamp(1rem, 3cqi, 2rem)`     | Consumer CSS |
| Fluid chart dimensions        | SVG `viewBox` + container `cqw` sizing | Research     |

**Why not bake fluid spacing into components?** Spacing is opinionated and consumers
override it constantly via Tailwind utilities. The library's stance: provide fluid
**typography** (universal, low-controversy) as utility classes, and document the `cqi`
pattern so consumers apply it to spacing/icons where they see fit. This respects the
"CSS-first, consumer-overridable" principle.

---

## Part 3: Container Style Queries (`@container style(...)`)

**Status: Chrome-only (not yet Baseline).** Document but do NOT implement in components.

Style queries let a component react to a **computed custom property** on an ancestor,
not just size:

```css
@container style(--tc-density: compact) {
  .card { padding: 0.5rem; }
}
@container style(--tc-density: comfortable) {
  .card { padding: 1.5rem; }
}
```

This would let a parent set `--tc-density: compact` and have ALL descendants react — a
powerful "context styling" mechanism that's cleaner than prop-drilling.

**Why defer:**
- Firefox and Safari support is incomplete (Baseline not reached as of mid-2026).
- The library's deprecation policy is "Baseline 2024+ only" for new CSS features.
- When it hits Baseline, the natural application is a `Density` prop on `AppShell` that
  sets `--tc-density` and lets child Cards/Forms react without explicit per-component flags.

**Action:** Revisit when `@container style()` reaches Baseline. Tracked in the v2.0+
research section of the ROADMAP.

---

## Part 4: The v2.0 default flip (ADR-0022)

ADR-0022 proposes flipping 3 container-aware components to **default-on** at v2.0:

- `Grid.ContainerResponsive` → default `true`
- `Card.ContainerAware` → default `true`
- `Split.ContainerAware` → default `true`

These three are chosen because they are the components **most commonly placed in
constrained containers** (grid cells, sidebars, splits). The other five (`Nav`,
`Pagination`, `Form`, `DefinitionGrid`, `SkeletonCardGrid`) stay opt-in because viewport
behavior is usually correct for them.

**Status:** Draft. This is a **major-version decision** — it changes default rendering for
every consumer on upgrade. It must ship with a migration guide and the other v2.0
default flips (self-host HTMX, semantic tokens). It is NOT a v1.x item and should not be
executed piecemeal.

**What this strategy document adds:** the recommendation that, at v2.0, the three
flipped components should also gain `container-name` so that nested containers don't
conflict (see Part 6).

---

## Part 5: Named containers (`container-name`)

When multiple `@container` wrappers nest (e.g., a `Card` with `ContainerAware: true` inside
an `AppShell` sidebar that's also a container), the innermost container wins for size
queries. This is usually correct. But **style queries** (Part 3) and explicit `@container
[name]` queries benefit from named containers:

```css
.parent { container-type: inline-size; container-name: sidebar; }
@container sidebar (min-width: 20rem) { /* only reacts to sidebar width */ }
```

Tailwind v4 supports `@container/sidebar` naming syntax. This is a **v2.0 consideration**
when style queries land — not needed for the current size-query-only world.

---

## Part 6: New component candidates — evaluated honestly

The ROADMAP lists five container-aware candidates. Here is the rigorous evaluation against
ADR-0018's three criteria (viewport-responsive today, plausibly placed in a constrained
container, clear named behavior change):

| Candidate       | Responsive today? | In constrained container? | Named behavior change?        | Verdict   |
| --------------- | ----------------- | ------------------------- | ----------------------------- | --------- |
| `Container`     | No (it IS the width limiter) | N/A (it sets the width) | None | **Reject** — circular |
| `Breadcrumbs`   | Minimal (separator) | Sometimes (header in split) | Weak — truncate vs wrap | **Defer** — no strong behavior change |
| `EmptyState`    | Minimal            | Yes (card, drawer)        | Weak — icon/text size only    | **Reject** — cosmetic only |
| `NotFound404`   | Yes (hero numeral) | No (always full page)     | N/A                           | **Reject** — never constrained |
| `Footer`        | Yes (columns)      | Rarely (usually full-width) | Column collapse             | **Defer** — rarely placed in a container |

**Conclusion:** None of the five candidates meet all three ADR-0018 criteria convincingly.
Adding `ContainerAware` to marginal candidates **dilutes the pattern** and adds props
noise for no real consumer benefit. The library's container-query coverage is already
strong at 8 components. **Do not expand for the sake of a higher count.**

---

## Part 7: The `containerAwareWrapper` consolidation — assessed and declined

The ROADMAP suggested extracting a shared `containerAwareWrapper(containerAware bool,
content templ.Component)` sub-template to remove the 8× hand-written `if/else` wrapper.

**Assessment after review:** the "boilerplate" is already minimal. Each component delegates
to an `xxxInner` sub-template (e.g., `cardInner`, the Grid renders inline), so the
per-component wrapper is 3–5 lines:

```templ
if props.ContainerAware {
    <div class="@container">
        @cardInner(props) { { children... } }
    </div>
} else {
    @cardInner(props) { { children... } }
}
```

A shared wrapper would require passing the entire inner rendering as a `templ.Component`
closure, adding indirection without meaningfully reducing lines. ADR-0009 (Accepted Clones)
explicitly blesses this kind of idiomatic, low-line-count duplication. **Declined — not
worth the refactor risk for 8 components × ~3 lines.**

---

## Recommendation summary

| Action                                                        | Priority | Risk   | Status     |
| ------------------------------------------------------------- | -------- | ------ | ---------- |
| Fluid typography `.tc-fluid-*` classes                        | High     | None   | **Done**   |
| Bind WC rejection as ADR-0033                                 | High     | None   | **Done**   |
| Document container style queries (`@container style()`)       | Medium   | None   | **Done** (this doc) |
| Hold the line on marginal `ContainerAware` candidates         | Ongoing  | None   | **Decision recorded** |
| Decline `containerAwareWrapper` consolidation                 | —        | —      | **Declined** |
| v2.0 default flip (Grid/Card/Split)                           | v2.0     | High   | ADR-0022 draft — do NOT execute pre-v2.0 |
| Named containers + style queries                              | v2.0+    | —      | Blocked on Baseline |

---

## References

- [ADR-0018: Container-Query-Native Contract](../adr/0018-container-query-native-contract.md)
- [ADR-0022: v2.0 Default-Flip Migration Plan](../adr/0022-v2-default-flip-migration.md)
- [ADR-0033: Web Components Rejection](../adr/0033-web-components-rejection.md)
- [Recipe: Container Queries](../recipes/container-queries.md)
- [Recipe: Fluid Typography](../recipes/fluid-typography.md)
- [MDN: Container Query Length Units](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_containment/Container_queries#container_query_length_units)
