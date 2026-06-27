# Cross-Project Analysis: cqrs-htmx/adminui ↔ templ-components

**Date:** 2026-06-27
**Scope:** How templ-components can improve cqrs-htmx/adminui (Part A) and how
adminui's needs can improve templ-components (Part B).

---

## Executive Summary

`cqrs-htmx/adminui` is a **drop-in embeddable admin panel** with its own
**706-line hand-rolled CSS design system** (CSS custom properties, no Tailwind).
`templ-components` is a **Tailwind-coupled component library**. The central
tension is a **deliberate, legitimate divergence** between two deployment models.

The clean win is **icons-only adoption**: adminui deleted its hand-maintained
15-icon SVG map and now delegates to templ-components' 101-icon set. For the
rest, we identified and filled real component gaps in templ-components itself.

---

## Part A: How templ-components improves cqrs-htmx/adminui

### A1. Icons-only adoption — DONE

adminui previously maintained `icons.go`: a 45-line file with 15 hand-written
SVG path strings, a fallback dot, and a lookup map. This is now **deleted**.
adminui delegates to `templ-components/icons`:

```go
// Before: 45 lines of hardcoded SVG path data
var icons = map[string]string{
    "dashboard": `<rect x="3" y="3" width="7" height="9" rx="1.5"/>...`,
    "users":     `<path d="M16 21v-2a4 4 0 0 0-4-4H6..."/>`,
    // ... 15 entries of copy-pasted SVG markup
}

// After: 19-entry name mapping, zero SVG data
var iconNames = map[string]icons.Name{
    "dashboard": icons.Squares2x2,
    "users":     icons.Users,
    "tenants":   icons.BuildingOffice2,
    // ...
}
```

**Impact:**

- adminui no longer maintains SVG path data — it gets 101 icons instead of 15
- 3 icons adminui needed (`BuildingOffice2`, `Key`, `ArrowRightOnRectangle`) were missing from tc and have been added
- Zero Tailwind dependency added — the `icons` package is CSS-agnostic
- `icons.IconPathData()` API was added to tc for consumers needing full `<svg>` wrapper control

### A2. CSS-coupled components — BLOCKED (see Part B, portability layer)

Every other tc component (`Card`, `Badge`, `Button`, `Table`, `Input`, `Toast`,
`Spinner`, `Avatar`, etc.) emits Tailwind utility classes. adminui **cannot**
adopt these without adding a Tailwind build step — which contradicts its
core promise of "No Tailwind, no build step."

**This is the single biggest blocker to broader adoption.**

---

## Part B: How adminui's needs improve templ-components

### B1. Missing icons — DONE

| Icon                | adminui use case             | Added as                |
| ------------------- | ---------------------------- | ----------------------- |
| Building/office     | tenants nav + dashboard stat | `BuildingOffice2`       |
| Key                 | credentials/API keys         | `Key`                   |
| Logout (arrow-exit) | sign-out button              | `ArrowRightOnRectangle` |

All sourced from official Heroicons v2 outline. Total: 99 → **101 icons**.

### B2. StatCard with icon — DONE

adminui hand-builds `statCardView` with a leading icon tile. tc's `StatCard`
previously lacked icon support. Now enhanced with optional `Icon icons.Name`
field — backward-compatible (empty = unchanged layout).

### B3. PageHeader — DONE

adminui hand-builds `.admin-content__head` (title + spacer + action buttons)
on every page. tc had no page-level header component. Added
`display.PageHeader` with Title, Subtitle, Breadcrumb slot, and Action slot.

### B4. DefinitionList — DONE

adminui hand-builds `.kv` (`<dl>` with grid layout) for user detail pages.
tc had no definition list. Added `display.DefinitionList` with typed
`DefinitionItem{Term, Detail, DetailComponent}` — supports both text and
rich content (badges, links).

### B5. SidebarNav — DONE

adminui's layout.templ builds a complete vertical sidebar (brand + icon nav +
footer). tc's navigation package only had horizontal top bars. Added
`navigation.SidebarNav` with brand/footer slots, icon+label items,
`CurrentPath` auto-active detection, and `aria-current="page"`.

### B6. ListNote — DONE

adminui hand-builds `listNote()` for "Showing N of M" truncation notices.
Added `display.ListNote(ListNoteProps{Shown, Total})` with `role="status"`.

### B7. CSS-variable portability layer — PROPOSED (see ADR)

The highest-leverage improvement: let tc components optionally render against
consumer-defined `--tc-*` CSS variables instead of hardcoded Tailwind classes.
This would unlock adoption in non-Tailwind projects like adminui.

---

## What changed (files)

### templ-components repo

| File                              | Change                                                          |
| --------------------------------- | --------------------------------------------------------------- |
| `icons/icon_names.go`             | +3 icon constants (BuildingOffice2, Key, ArrowRightOnRectangle) |
| `icons/icon_paths.go`             | +3 path entries, +`IconPathData()` exported function            |
| `icons/snapshot_test.go`          | +`TestIconPathData` (3 subtests)                                |
| `display/card.templ`              | StatCard enhanced with optional `Icon` field                    |
| `display/card_test.go`            | +2 StatCard icon tests                                          |
| `display/page_header.templ`       | **NEW** — PageHeader component                                  |
| `display/page_header_test.go`     | **NEW** — 7 tests                                               |
| `display/definition_list.templ`   | **NEW** — DefinitionList component                              |
| `display/definition_list_test.go` | **NEW** — 5 tests                                               |
| `display/list_note.templ`         | **NEW** — ListNote component                                    |
| `display/list_note_test.go`       | **NEW** — 5 tests                                               |
| `navigation/sidebar_nav.templ`    | **NEW** — SidebarNav component                                  |
| `navigation/sidebar_nav_test.go`  | **NEW** — 9 tests                                               |
| `docs/icons-only-adoption.md`     | **NEW** — adoption guide for non-Tailwind apps                  |
| `docs/cross-project-analysis.md`  | **NEW** — this report                                           |
| `README.md`                       | Updated component counts (69→73), new components in catalog     |
| `AGENTS.md`                       | New components documented, icon count corrected                 |

### cqrs-htmx/adminui repo

| File                | Change                                                                   |
| ------------------- | ------------------------------------------------------------------------ |
| `icons.go`          | Replaced 45-line hardcoded SVG map with `iconNames` map delegating to tc |
| `icons_test.go`     | Updated to reference `iconNames` instead of `icons`                      |
| `go.mod` / `go.sum` | Added templ-components dependency (local replace for dev)                |

---

## Metrics

| Metric                      | Before | After                                                      |
| --------------------------- | ------ | ---------------------------------------------------------- |
| tc icons                    | 99     | **101**                                                    |
| tc components               | 69     | **73** (+PageHeader, DefinitionList, ListNote, SidebarNav) |
| adminui hardcoded SVG paths | 15     | **0** (delegated to tc)                                    |
| adminui available icons     | 15     | **101** (via tc)                                           |
| tc new tests                | —      | **+34** (3+2+7+5+5+9+3)                                    |
| tc test packages passing    | 12     | 12                                                         |
| tc lint issues              | 0      | 0                                                          |
