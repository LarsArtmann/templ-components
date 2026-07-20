# ADR 0016: Grid-first for 2D layouts

## Date

2026-07-20

## Status

Accepted

## Context

Before this ADR, the `layout` package shipped six head/theme primitives
(`Base`, `Minimal`, `ThemeScript`, `ThemeToggle`, `Script`, `Stylesheet`) and
**zero body-layout primitives**. Every consumer rebuilt the same 2D patterns
by hand:

- **Sidebar + main app shell** — hand-rolled `lg:grid-cols-[16rem_1fr]` (or
  fragile flex with hardcoded widths). The single most-rebuilt pattern.
- **Page content container** — the snippet `max-w-6xl mx-auto px-4 sm:px-6
lg:px-8` was literally copy-pasted in `examples/demo/demo.templ` and in
  every consumer codebase.
- **Content + aside split** — hand-rolled `grid md:grid-cols-3` with manual
  column spans. No RTL awareness.
- **Multi-column footer** — `Footer` rendered only a single brand+copyright
  row. Multi-column required forking the component.
- **Settings forms with aligned labels** — `Form.Inline` used
  `flex flex-wrap`, which broke label alignment when fields wrapped to a
  second row.

Grid was already used **correctly** in 9 places (`display.Grid`,
`DefinitionList`'s `grid-cols-[auto_1fr]`, `Calendar`'s 7-col day grid,
demo responsive grids). The gap was not "we use grid wrong" — it was "we
don't ship grid-based layout primitives at all."

The 48 flex usages were mostly correct (1D layouts where flex is the right
tool). A blanket "grid-ify all flex" migration would have regressed quality,
not improved it.

## Decision

Codify **grid = 2D, flex = 1D** as the library's layout-selection rule, and
ship the missing 2D layout primitives.

### Rule

| Layout shape                      | Tool | Rationale                                                             |
| --------------------------------- | ---- | --------------------------------------------------------------------- |
| Both rows AND columns matter (2D) | Grid | Sidebar+main, content+aside, multi-col footer, card dashboards        |
| Only one axis matters (1D)        | Flex | Nav bars, button rows, chip wraps, vertical stacks, inline forms      |
| Vertical rhythm only              | Flex | `Stack` uses `flex flex-col` — grid would be overkill for a 1D column |

### Primitives added

| Component                | Purpose                                                              | Grid pattern                                                  |
| ------------------------ | -------------------------------------------------------------------- | ------------------------------------------------------------- |
| `layout.AppShell`        | Sidebar + header + main app shell                                    | `lg:grid-cols-[var(--tc-sidebar-w)_minmax(0,1fr)]`            |
| `layout.Container`       | Centered max-width wrapper with responsive padding                   | Single column (no grid — wrapper only)                        |
| `layout.Split`           | Content + aside (article+sidebar / detail+metadata)                  | `grid-cols-1 md:grid-cols-{N}` with `md:col-span-{M}` on Main |
| `layout.Stack`           | Vertical rhythm (1D, flex not grid)                                  | `flex flex-col` (deliberately NOT grid)                       |
| `navigation.Footer`      | Multi-column footer (backward-compat: empty Columns = legacy single) | `grid grid-cols-2 md:grid-cols-4 gap-8`                       |
| `forms.Form` Layout enum | Stack (default), Inline (legacy flex), Grid (settings)               | `sm:grid-cols-[auto_minmax(0,1fr)]` for Grid variant          |

### Mandatory: `minmax(0, 1fr)`, never bare `1fr`

Every grid column that should shrink-to-fit MUST use `minmax(0, 1fr)`, never
bare `1fr`. Bare `1fr` allows a wide child (a `<table>`, a long URL, a
`<pre>` block) to force the column wider than its container, causing
page-wide horizontal scroll. This is the single most common grid footgun.

Enforced in:

- `AppShell`: main column is `minmax(0, 1fr)`
- `Split`: both columns get `min-w-0` (the flex/grid complement)
- `Form` Grid layout: value column is `minmax(0, 1fr)`

See the recipe `docs/recipes/grid-blowout-minmax.md` for the before/after
reproduction.

### Logical positioning for RTL

`Split.AsidePosition` uses logical names (`Start`/`End`), never physical
(`Left`/`Right`). Source order + CSS grid auto-placement handles the mirror
in `dir="rtl"`. No physical `left-`/`right-` utilities are emitted by any
new primitive.

## Consequences

### Positive

- The #1 hand-rolled consumer pattern (sidebar+main app shell) becomes a
  one-liner: `@layout.AppShell(props)`.
- The repeated `max-w-6xl mx-auto px-4 sm:px-6 lg:px-8` snippet has a single
  source of truth: `layout.Container`.
- Grid blowout is prevented by construction (`minmax(0, 1fr)` baked in).
- Future contributors have a documented rule preventing the "grid-ify
  everything" Verschlimmbesserung trap.
- `Footer` and `Form` extensions are backward compatible — zero breakage
  for existing consumers.

### Negative

- `layout` package grows from 6 to 10 components. Acceptable — the new
  primitives are cohesive (all are body-layout concerns).
- `Form.Inline` is soft-deprecated. Removed target is v1.0; until then
  both `Inline bool` and `Layout FormLayout` coexist (Layout wins).
- Mobile sidebar drawer is NOT built into `AppShell` (would require a
  `layout → display` import, breaking the import graph). Consumers wire
  their own mobile nav via the `MobileNav` slot (typically a
  `display.Drawer`). Documented in the AppShell recipe.

### Non-goals (explicitly deferred)

- Migrating existing 48 flex usages to grid — most are correctly 1D. M15
  audits and documents keep-decisions per usage.
- Container-query-as-default — viewport-default remains the safer baseline.
  `Grid.ContainerResponsive` is the opt-in path; may revisit post-v1.0.
- CSS subgrid — tracked in `docs/research/css-subgrid.md`. Would unlock
  true 2D alignment in Card/DefinitionList once Baseline 2025.

## Appendix: New-component checklist (apply before adding any grid primitive)

- [ ] Uses `minmax(0, 1fr)` for any flexible column (never bare `1fr`)
- [ ] Uses logical CSS properties (`ms-`, `me-`, `start-`, `end-`, `ps-`, `pe-`)
- [ ] Has `dark:` variants on every neutral and semantic color
- [ ] Embeds `utils.BaseProps` + propagates `ID`/`Class`/`Attrs`/`AriaLabel`
- [ ] Typed enums have `IsValid()` + map+fallback (no panics on unknown values)
- [ ] Constructor `DefaultXxxProps()` provided
- [ ] Tests cover: defaults, each enum value, unknown-value fallback,
      BaseProps propagation, dark-mode variants, grid-blowout guard
- [ ] Documented in SKILL.md + FEATURES.md

## Appendix: Flex-usage audit (M15, 2026-07-20)

Audited all 48 `flex` usages across `.templ` source files to classify each as
1D-keep or 2D-migrate. **Result: 48 keep, 0 migrate.** Every existing flex
usage is a correct 1D layout where grid would be the wrong tool.

### Categories observed

| Pattern                             | Count | Verdict | Why flex is correct                                  |
| ----------------------------------- | ----- | ------- | ---------------------------------------------------- |
| `flex items-center justify-center`  | ~12   | keep    | 1D centering (modals, overlays, empty states)        |
| `flex items-center justify-between` | ~8    | keep    | 1D space-between (nav bars, headers, card footers)   |
| `inline-flex items-center gap-N`    | ~14   | keep    | 1D button/icon rows (buttons, badges, action groups) |
| `flex flex-col`                     | ~6    | keep    | 1D vertical stack (use `layout.Stack` for new ones)  |
| `flex-wrap` / `sm:flex-wrap`        | ~4    | keep    | 1D wrapping row (nav links, filter chips)            |
| `flex-shrink-0` / `flex-1`          | ~4    | keep    | Flex child sizing (prevent blowout, fill remaining)  |

### Decision

No existing flex usage was migrated. The rule "grid = 2D, flex = 1D" is
retroactively confirmed: every flex usage in the codebase has exactly one
axis that matters. Future 2D layouts should use the new grid primitives
(`AppShell`, `Split`, multi-col `Footer`, `Form` Grid layout). Future 1D
layouts should use `Stack` (vertical) or raw `flex` utilities (horizontal).

This appendix exists to prevent future "why isn't this flex a grid?" churn.
If a contributor questions a flex usage, point them here.
