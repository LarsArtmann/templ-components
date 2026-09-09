# Bug Report: Table-in-Card Double Border

**Consumer:** cqrs-htmx adminui (v0.15.0)
**Date:** 2026-07-12
**Severity:** Medium — visible visual defect on every table page
**Components affected:** `display.Table` inside `display.Card(CardPaddingNone)`

---

## The Problem

When `display.Table` is nested inside `display.Card(CardPaddingNone)`, two concentric 1px borders render — one from each component. This is visible as a double-line border on every table-in-card layout.

### Reproduction

```go
@display.Card(display.CardProps{
    Title:   "Users",
    Padding: display.CardPaddingNone,
}) {
    @display.Table(display.TableProps{
        Headers: []string{"Name", "Email"},
        Hover:   true,
        Body:    userRows,
    })
}
```

### Root Cause

Both components hardcode identical border classes:

**Card** (`card.templ:10`):

```go
const cardShellClass = "bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-xs"
```

Card always emits `cardShellClass + " overflow-hidden"` (line 58).

**Table** (`table.templ:163`):

```html
<div class="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700"></div>
```

Table always emits this wrapper div — no option to disable the border.

When Table is a direct child of `CardPaddingNone` (which renders children without a padding wrapper, lines 85-90), both borders are adjacent:

```
Card div:   border border-gray-200 rounded-lg overflow-hidden
  └─ Table div:  border border-gray-200 rounded-lg overflow-x-auto
```

The result: two stacked 1px gray borders with slightly mismatched `rounded-lg` corners.

### Why `props.Class` can't fix it

`utils.Class()` uses `tailwind-merge-go` to merge classes, but it can only **append** — it cannot strip hardcoded classes from the component's own template. Passing `Class: "!border-0"` on `TableProps` adds the class to the inner `<table>` element, not the wrapper div where the border lives. There's no way for the consumer to suppress the wrapper div's border.

---

## Current Workaround

We added a CSS rule targeting the exact nesting pattern:

```css
/* Card always emits overflow-hidden; Table wrapper always emits overflow-x-auto.
   This selector matches only Table-inside-Card and suppresses the inner border. */
.overflow-hidden > .overflow-x-auto {
  border: 0 !important;
  border-radius: 0 !important;
}
```

This works but is **fragile** — it depends on Tailwind utility class names that could change in a future version. If Table's wrapper ever drops `overflow-x-auto` or Card drops `overflow-hidden`, the selector silently stops matching and the double border returns.

---

## Suggested Fixes (pick one)

### Option A: Add `Flush` option to TableProps (preferred)

```go
type TableProps struct {
    // ... existing fields ...
    Flush bool // When true, suppresses the wrapper div's border + rounded corners.
               // Use when Table is nested inside a Card(CardPaddingNone).
}
```

In the template:

```templ
<div class={ utils.Class(
    "overflow-x-auto",
    tableFlushClass(props.Flush), // "rounded-lg border border-gray-200 dark:border-gray-700" when !Flush, "" when Flush
)}>
```

This is the cleanest API: the consumer opts in with `Flush: true` and the component handles it.

### Option B: Add `Borderless` option to TableProps

Same pattern, different name. `Borderless: true` removes all borders (wrapper + `Bordered` table element). More aggressive than `Flush` but covers more use cases.

### Option C: Make Card auto-detect Table children

Card could check if its child is a Table and suppress its own border. This is magic — not recommended. Components shouldn't inspect their children's type.

### Option D: Split the wrapper into a separate concern

Extract the `overflow-x-auto rounded-lg border` wrapper into a `TableScrollWrapper` that consumers compose explicitly. Table renders just the `<table>` element. This is the most flexible but breaks the simple API for standalone tables.

---

## Impact

This pattern — Table inside Card(CardPaddingNone) — is the **standard layout for admin dashboards**. It appears on:

- Dashboard recent activity tables
- Audit log tables
- User list tables
- Tenant list tables
- Member list tables

In cqrs-htmx adminui alone, this nesting appears in 5 templates. The double border is visible on every page.

---

## Related: Table cell padding

Separately, Table hardcodes `px-4 py-3` on all cells (`<th>` and `<td>`). Our previous adminui tables used `px-4 py-2` (more compact, standard for data-heavy admin panels). The library's `py-3` makes tables feel spacious for a dashboard context. A `Compact` option or a `CellPadding` field would help, but this is lower priority than the border issue.
