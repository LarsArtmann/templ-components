# templ-components — SDK Feedback from cqrs-htmx

**Consumer:** [cqrs-htmx](https://github.com/larsartmann/cqrs-htmx) — Go CQRS+HTMX library with an admin dashboard module (`adminui/v4`)
**Date:** 2026-07-10
**Version used:** v0.13.0
**Session:** Full adoption pass — replaced 10 hand-rolled components with library equivalents in the adminui module (dashboard, users, tenants, members, audit pages)

---

## Context: What cqrs-htmx adminui is

`adminui` is a ready-made admin dashboard panel — one-call mount behind session middleware, rendering users/tenants/members/audit pages with HTMX partial swaps, toast notifications, and an offline-aware sync indicator. Built with templ + HTMX + a custom CSS-variable design system (not Tailwind palette). It ships as a Go module that consumers import.

**Before this session:** adminui imported templ-components v0.13.0 but used only 1 of 84 components (`icons.IconPathData`). Everything else was hand-rolled.

**After this session:** 10 library components adopted, 1,475 lines of hand-rolled code deleted.

---

## What worked superbly

### 1. `icons.IconPathData` — the icons-only adoption path is excellent

We adopted the icons package first (months before the rest). The `IconPathData(name) []string` API returns raw SVG path data, letting us wrap it in our own `<svg>` element with custom dimensions and stroke width. No CSS dependency, no component coupling. This is the gateway adoption path and it's flawless.

### 2. Typed BadgeType + graceful fallback

`display.BadgeType` (BadgePrimary, BadgeSuccess, etc.) is a typed enum with a graceful fallback to `BadgeNeutral` on unknown values. The mapping from our internal string-based badge kinds (`"green"`, `"blue"`, `"amber"`) to typed BadgeType was a clean one-function migration. No panics, no missing styles.

### 3. `display.EmptyState` — replaced 6 hand-rolled copies

Every page had its own `@empty(iconName, title, hint)` with the same div structure. `display.EmptyState` with `Icon`, `Title`, `Description` is the exact same shape, typed. One of the highest-ROI adoptions — eliminated duplication across 5 templates.

### 4. `display.RelativeTime` — replaced a hand-rolled time formatter

Our `relTime()` function (19 lines of switch/case formatting "just now", "3m ago", etc.) is now `display.RelativeTime` which renders the same text AND adds a `<time datetime>` element for accessibility + auto-refresh via `Intl.RelativeTimeFormat`. We got accessibility and live-updating for free.

### 5. `display.ListNote` — exact behavioral match

Our `listNote(shown, total)` rendered "Showing N of M" only when `total > shown`. The library version has the exact same conditional logic and the same semantics. Drop-in replacement, zero behavior change.

### 6. `display.StatCard` — clean props, icon integration

Our `statCardView` was a custom div with hardcoded Tailwind classes. `display.StatCard` accepts `Value`, `Label`, `Icon` (typed `icons.Name`) — same shape, cleaner API, and the icon tile styling is consistent with the library's design language.

---

## Pain points

### 1. `bg-white` hardcoded — breaks CSS-variable-based theming

**Severity:** Medium (required a fragile workaround)

`display.Card` uses `cardShellClass = "bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-xs"`. Same in `display.Table` (`bg-white dark:bg-gray-900`), `display.Accordion`, `display.Modal`, etc.

adminui uses CSS variables (`var(--surface)`) for all surface colors, not the Tailwind palette. When a `display.Card` renders `bg-white`, it stays white in dark mode because:

- adminui uses `prefers-color-scheme` to swap CSS variables (not the `.dark` class the library expects)
- The `dark:bg-gray-800` variant never activates (no `.dark` class in the DOM)

**Workaround I used:** Added a CSS bridge in `tailwind.css`:

```css
.bg-white {
  background-color: var(--surface);
}
```

This works because adminui never uses `bg-white` directly (it uses `style="background:var(--surface)"`). But it's fragile — if a future adminui component uses `bg-white` intentionally (rare but possible), the bridge silently overrides it.

**Suggestion:** Consider using CSS custom properties as the surface token instead of hardcoding `bg-white`. E.g., `--tc-surface: #fff` (light) / `--tc-surface: oklch(...)` (dark) via `@theme`, then components use `bg-[var(--tc-surface)]` or a custom utility. Consumers override `--tc-surface` in their `@theme`. This is the same pattern adminui already uses internally.

### 2. Promoted `BaseProps.Class` field — struct literal gotcha

**Severity:** Low (Go language limitation, but undocumented)

Every component's props embed `utils.BaseProps`. In Go, promoted fields cannot be set in struct literals:

```go
// This compiles:
display.GridProps{Cols: display.GridCols3, Class: "..."}
// Wait — Class is promoted from BaseProps. This does NOT compile:
```

The workaround:

```go
display.GridProps{
    BaseProps: utils.BaseProps{Class: "..."},
    Cols:      display.GridCols3,
}
```

This bit me during the Grid adoption — I wrote `Class: "..."` in the struct literal and got a compile error. The fix is simple once you know it, but it's not documented in the consumer guide.

**Suggestion:** Add a one-liner to the consumer guide under "Quick start" or "Theming":

> **Setting `Class` on components:** Because `Class` is promoted from the embedded `utils.BaseProps`, you must set it via the embedded struct in literal initializers:
>
> ```go
> display.GridProps{BaseProps: utils.BaseProps{Class: "..."}, Cols: display.GridCols3}
> ```

### 3. `display.Grid` lacks `auto-fit`/`minmax` patterns

**Severity:** Low (workaround via `BaseProps.Class`)

Dashboard stat cards commonly use `grid-template-columns: repeat(auto-fit, minmax(190px, 1fr))` — the card grid responds to container width, not viewport breakpoints. The library's `GridCols` enum offers fixed breakpoint-based columns (`grid-cols-1 sm:grid-cols-2 lg:grid-cols-3`). There's no `GridColsAutoFit` or equivalent.

**Workaround:** Passed the custom grid template via `BaseProps.Class`:

```go
display.GridProps{
    BaseProps: utils.BaseProps{Class: "[grid-template-columns:repeat(auto-fit,minmax(190px,1fr))]"},
    Cols:      display.GridCols3, // ignored when Class overrides grid-template-columns
}
```

Note: setting `Cols` AND `Class` with a grid-template override is slightly misleading — `Cols` is effectively ignored but still required for the component to emit the grid wrapper class.

**Suggestion:** Add a `GridColsAutoFit` variant or a `MinColWidth string` field that generates the `auto-fit, minmax(...)` template. Alternatively, document the `Class` override pattern in the Grid section of the consumer guide.

### 4. Tailwind v4 `@source` path is GOPATH-dependent and fragile

**Severity:** Medium (silent CSS breakage if path is wrong)

To get templ-components' utility classes included in the compiled CSS, the consumer must add:

```css
@source "/home/lars/go/pkg/mod/github.com/larsartmann/templ-components@v0.13.0";
```

This path is:

- GOPATH-dependent (changes per developer/machine)
- Version-dependent (must update when bumping the dependency)
- Silent — if the path is wrong, classes simply don't appear in the CSS, and components render unstyled. No error, no warning.

The `tailwind-v4-adoption-guide.md` documents the `@source` pattern but uses a `vendor/` example. For Go module cache consumers, the path is non-obvious.

**Suggestion:**

- Document the Go module cache path pattern explicitly: `@source "$(go env GOMODCACHE)/github.com/larsartmann/templ-components@v0.13.0"`
- Or better: provide a nix flake app / make target that generates the `@source` line from `go list -m -f '{{.Dir}}'`
- Consider whether Tailwind v4's `@source` can accept a Go-import-path-based resolution (unlikely, but worth investigating)

### 5. `Card` title hardcodes `<h3>` and header styling

**Severity:** Low (worked around by using `CardPaddingNone` for table cards)

`display.Card` renders `Title` as:

```html
<h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">{ props.Title }</h3>
```

adminui's card headers use `text-[15px]` with a different padding pattern (`px-[18px] py-[15px]`). There's no way to customize the title element type, its classes, or the header padding without using the `Body` slot and building the entire header yourself.

This wasn't blocking — I used `CardPaddingNone` for table-inside-card layouts (where the title is rendered separately). But for titled cards with custom header styling, the component is less flexible than the hand-rolled version.

**Suggestion:** Consider adding `TitleClass string` and `HeaderClass string` fields to `CardProps`, or a `Header templ.Component` slot (like `Footer`) for fully custom headers.

### 6. Dark mode strategy mismatch (`.dark` class vs `prefers-color-scheme`)

**Severity:** Low for adminui (bridge worked), but conceptual mismatch worth noting

The library uses class-based dark mode: `@custom-variant dark (&:where(.dark, .dark *))`. adminui uses `@media (prefers-color-scheme: dark)`. These are different strategies:

- Library: `.dark` class on `<html>` (or any ancestor), toggled by `ThemeToggle` + `ThemeScript`
- adminui: `prefers-color-scheme` media query swaps CSS variable values

Because adminui never adds a `.dark` class, ALL `dark:` variants from templ-components are inert. The `bg-white` CSS bridge handles the surface color, but any component that relies on `dark:` variants for non-surface styling (e.g., `dark:text-gray-300` in BadgeNeutral) won't dark-mode-correctly.

In practice this was acceptable — the badge type colors (green/blue/red/amber) are vivid enough in both modes. But it's a conceptual gap: consumers with `prefers-color-scheme` dark mode get partial dark mode support from the library.

**Suggestion:** This is a known tradeoff (class-based dark mode enables user-toggleable themes). Document it explicitly for `prefers-color-scheme` consumers:

- Surface colors: bridge via CSS variable override (as I did with `.bg-white`)
- Component-internal `dark:` variants: will not activate without `.dark` class
- Full dark mode support requires adopting the `.dark` class strategy (or adding a media-query-to-class polyfill)

---

## What didn't block but could be better

### 7. `DefaultSpinnerProps()` is nice — wish all components had one

`feedback.DefaultSpinnerProps()`, `display.DefaultEmptyStateProps()`, etc. are great for the "I just want the defaults with one field changed" pattern. But not all components have them (e.g., `display.Grid` has no `DefaultGridProps()`). Consistency would help.

### 8. `CardPaddingNone` still wraps children in a padding div

When using `display.Card` with `CardPaddingNone` (for table-inside-card), children are still wrapped in `<div class={cardPaddingClass("none")}>`. For `CardPaddingNone`, this div presumably has no padding class — but the extra wrapping `<div>` is unnecessary for table layouts where you want the `<table>` to be a direct child of the card for `overflow-x-auto` to work correctly.

Not blocking (my table-in-card works), but slightly inelegant.

### 9. Discoverability — the skill says "84 components" but finding the right one is manual

The consumer guide's "By use case" table is helpful, but the adoption process still requires reading the component's `.templ` source to understand:

- Exact prop field names (which are promoted from `BaseProps`)
- Whether the component's hardcoded classes match your design system
- What the dark mode strategy requires

A `go doc`-style generated API reference (pkg.go.dev works for this) would speed up the evaluation phase.

---

## Adoption summary

| Library component      | Replaced hand-rolled                                               | Call sites | Verdict                        |
| ---------------------- | ------------------------------------------------------------------ | ---------- | ------------------------------ |
| `feedback.Spinner`     | `spinner()` inline SVG                                             | 1          | Clean                          |
| `display.ListNote`     | `listNote()`                                                       | 3          | Exact match                    |
| `display.RelativeTime` | `relTime()` 19-line switch                                         | 3          | Upgrade (a11y + auto-refresh)  |
| `display.EmptyState`   | `empty()`                                                          | 6          | Clean                          |
| `display.Badge`        | `badge()` + `badgeColor()` + `badgeColors` map + `roleBadgeKind()` | ~20        | Clean (typed enum)             |
| `display.StatCard`     | `statCardView()`                                                   | 1          | Clean                          |
| `display.Card`         | `<div class="rounded-[10px]...">`                                  | 2          | Worked (CSS bridge needed)     |
| `display.Grid`         | `<div class="grid gap-4...">`                                      | 1          | Worked (Class override needed) |
| `icons.IconPathData`   | N/A (pre-existing)                                                 | ~15        | Excellent                      |

**Net:** -1,475 lines, +10 components adopted, build + tests + lint green.

---

## Recommendations for the library (prioritized)

1. **Document the `BaseProps.Class` struct literal gotcha** — 5-minute fix, saves every new consumer a compile error
2. **Document the Go module cache `@source` path** — saves 30 minutes of CSS-not-applying debugging
3. **Add `GridColsAutoFit` or `MinColWidth` to GridProps** — common dashboard pattern, currently requires `Class` escape hatch
4. **Consider CSS-variable-based surface tokens** instead of hardcoded `bg-white` — eliminates the fragile `.bg-white` bridge workaround
5. **Add `TitleClass`/`HeaderClass` or `Header templ.Component` slot to Card** — unlocks custom card headers without full Body-slot replacement
6. **Document dark mode strategy implications for `prefers-color-scheme` consumers** — set expectations about partial dark mode support

---

_Analysis: See [Dark Mode & Theming Research](../dark-mode-research.md) for a
first-principles analysis of pain points #1 and #6 (they share a root cause) and
the recommended `@theme` palette override pattern that replaces the fragile `.bg-white`
bridge without redesigning the library's token architecture._

_Reported by: Crush (AI engineering partner) during cqrs-htmx adminui templ-components adoption session._
