# ADR 0009: Accepted Code Duplication

## Status

Accepted

## Context

The `templ-components` library is checked for code duplication via
`art-dupl --semantic --sort total-tokens -t 8`. After multiple dedup passes
(commits `4c2b00e` and earlier), the detector still flags 4 clone groups.
Each is examined below and the rationale for accepting it is documented here
so future maintainers do not waste time re-investigating.

The **goal** is zero harmful duplication — not zero report lines. Each
remaining clone is either idiomatic UI layout, required by the templ DSL, or
in the demo binary which is not production code.

## Accepted Clones

### 1. `display/card.templ:246-254`, `errorpage/errordetail.templ:24-32`, `errorpage/errorpage.templ:23-31` — flex row with icon + content (3 occurrences)

Each implementation is a **visually distinct design choice** that happens to
share the Tailwind `flex` row idiom:

| File                | Alignment      | Gap | Icon container                       | Content element        |
| ------------------- | -------------- | --- | ------------------------------------ | ---------------------- |
| `card.templ`        | `items-center` | `4` | 12×12 colored box (`bg-blue-50`)     | `<dl>` (semantic HTML) |
| `errordetail.templ` | `items-start`  | `3` | bare icon (`mt-0.5`)                 | `<div>`                |
| `errorpage.templ`   | `items-start`  | `4` | 10×10 colored box (`style.AccentBG`) | `<div>`                |

The card pattern uses `<dl>` because `statCardFigures` renders `<dt>` and
`<dd>` (semantic definition-list pairs). An `IconRow` helper would either need
8+ parameters (alignment, gap, icon container class, icon body, content
class, content body, content element) or break the `<dl>` semantics.

**Verdict**: accept as idiomatic UI primitive. Extracting would require more
parameter noise than the duplicated lines save.

### 2. `display/copy_button.templ:63-78` and `:79-94` — `<a>` vs `<button>` CopyButton branches

Both branches render the same DOM contract (`data-tc-copy`, `data-tc-copy-label`,
optional `id`, `class`, optional `aria-label`, `props.Attrs`). The element
type (`<a>` vs `<button>`) and the two element-specific attributes (`href`
vs `type="button"`) are **required by templ**: attributes must live on the
interactive element itself because the singleton copy script listens for
clicks via `closest('[data-tc-copy]')` on the document.

The `copyButtonContent(icon, label)` sub-template already factors out the
inner icon + status span; the outer attribute block must remain inline per
branch.

**Verdict**: accept as required DOM duplication. Templ has no mechanism to
push attributes onto a wrapper element.

### 3. `examples/demo/demo.templ:322-343` — three demo table rows (3 occurrences)

Three rows of `@demoTableBody()` rendering distinct badge types
(Success / Failed / Passed) for visual demonstration. The structural
similarity is intentional demo content.

`examples/demo` is excluded from `golangci-lint` per `.golangci.yml` paths
filter and is the consumer-facing demo binary, not production code.

**Verdict**: accept as demo content.

### 4. `examples/demo/demo.templ:129-140` and `:141-171` — Buttons and Form Controls sections

Both blocks follow the demo's standard pattern:

```
@demoSection("Section Name")
<div class="layout-classes-for-this-section">
  ... distinct component calls ...
</div>
```

The `@demoSection` wrapper plus `<div>` container is the project's demo
shell convention. The content (9 button variants vs 8 form controls) is
completely different. The structural similarity is intentional demo
rhythm.

**Verdict**: accept as demo structure.

## Decision

These 4 clone groups are **intentional** and remain in the codebase.
Future dedup passes should not attempt to extract them unless the underlying
design changes (e.g. consolidating the error header and StatCard into a
single visual primitive — which would require a major redesign, not just
templ refactoring).

## Consequences

- New contributors running `art-dupl` will see 4 "false positive" clone groups
- This ADR is the canonical reference for why they remain
- Any new component that falls into one of these categories should not be
  forced into an existing extraction just to reduce the report count
- The dedup code review can stop at "matches an entry in ADR 0009"
