# Cross-Project Analysis: cqrs-htmx/adminui ↔ templ-components

**Date:** 2026-06-27
**Scope:** How templ-components can improve cqrs-htmx/adminui (Part A) and how
adminui's needs can improve templ-components (Part B).
**Architecture decision:** Tailwind CSS v4+ is the standard for all projects
(see `docs/adr-001-tailwind-v4-standard.md`).

---

## Executive Summary

`cqrs-htmx/adminui` is a **drop-in embeddable admin panel** with its own
**706-line hand-rolled CSS design system** (CSS custom properties, no Tailwind).
`templ-components` is a **Tailwind v4+ component library**.

**Decision:** Tailwind v4+ everywhere. adminui should migrate to Tailwind v4+
(CSS-first, no Node.js) rather than templ-components accommodating non-Tailwind
projects. This eliminates adminui's 706-line custom design system and unlocks
every templ-components component for free.

**Completed work:** 4 new components, 3 new icons, StatCard enhancement, and
`IconPathData()` API — all filling gaps adminui revealed. The color bridge CSS
experiment was built and **rejected** (see ADR-001).

---

## Part A: How templ-components improves cqrs-htmx/adminui

### A1. Component gaps filled — DONE

adminui hand-builds several patterns that templ-components now provides:

| adminui hand-built | templ-components replacement | Status |
|--------------------|------------------------------|--------|
| `.admin-content__head` (title+action) | `display.PageHeader` | ✅ Added |
| `.kv` (definition list) | `display.DefinitionList` | ✅ Added |
| `statCardView` (icon+value tile) | `display.StatCard` (+ `Icon` field) | ✅ Enhanced |
| `listNote()` ("Showing N of M") | `display.ListNote` | ✅ Added |
| `.admin-sidebar` + `.admin-nav` | `navigation.SidebarNav` | ✅ Added |
| `@icon()` (15 hardcoded SVGs) | `icons.Icon` (101 icons) | ✅ Available |

adminui can adopt all of these **once it migrates to Tailwind v4+**.

### A2. The path forward: migrate to Tailwind v4+

adminui's `admin.css` is a well-designed but isolated system. The migration path
is documented in `docs/tailwind-v4-adoption-guide.md`:

1. **Phase 1 (1hr):** Wrap existing CSS tokens in a Tailwind `@theme` block
2. **Phase 2 (ongoing):** Replace custom classes with templ-components or raw utilities
3. **Phase 3 (when done):** Delete the 706-line `admin.css`

Tailwind v4 is CSS-first — no Node.js, no `package.json`. The standalone CLI
binary builds CSS in one command. This is not a framework lock-in; it's adopting
the standard tool.

---

## Part B: How adminui's needs improved templ-components

### B1. Missing icons — DONE

| Icon | adminui use case | Added as |
|------|-----------------|----------|
| Building/office | tenants nav + dashboard stat | `BuildingOffice2` |
| Key | credentials/API keys | `Key` |
| Logout (arrow-exit) | sign-out button | `ArrowRightOnRectangle` |

Total: 99 → **101 icons** (official Heroicons v2 outline paths).

### B2. StatCard with icon — DONE

adminui hand-builds `statCardView` with a leading icon tile. tc's `StatCard`
previously lacked icon support. Now enhanced with optional `Icon icons.Name`
field — backward-compatible.

### B3. PageHeader — DONE

adminui hand-builds `.admin-content__head` on every page. Added
`display.PageHeader` with Title, Subtitle, Breadcrumb slot, Action slot.

### B4. DefinitionList — DONE

adminui hand-builds `.kv` for user detail pages. Added `display.DefinitionList`
with typed `DefinitionItem{Term, Detail, DetailComponent}`.

### B5. SidebarNav — DONE

adminui builds a complete vertical sidebar. Added `navigation.SidebarNav` with
brand/footer slots, icon+label items, `CurrentPath` auto-active detection.

### B6. ListNote — DONE

adminui hand-builds `listNote()`. Added `display.ListNote` with `role="status"`.

### B7. IconPathData API — DONE

Added `icons.IconPathData(name) []string` — raw SVG path d-strings for consumers
needing full `<svg>` wrapper control. Useful for any project, not just adminui.

---

## What changed (files)

### templ-components repo

| File | Change |
|------|--------|
| `icons/icon_names.go` | +3 icon constants (BuildingOffice2, Key, ArrowRightOnRectangle) |
| `icons/icon_paths.go` | +3 path entries, +`IconPathData()` exported function |
| `icons/snapshot_test.go` | +`TestIconPathData` (3 subtests) |
| `display/card.templ` | StatCard enhanced with optional `Icon` field |
| `display/card_test.go` | +2 StatCard icon tests |
| `display/page_header.templ` | **NEW** — PageHeader component |
| `display/page_header_test.go` | **NEW** — 7 tests |
| `display/definition_list.templ` | **NEW** — DefinitionList component |
| `display/definition_list_test.go` | **NEW** — 5 tests |
| `display/list_note.templ` | **NEW** — ListNote component |
| `display/list_note_test.go` | **NEW** — 5 tests |
| `navigation/sidebar_nav.templ` | **NEW** — SidebarNav component |
| `navigation/sidebar_nav_test.go` | **NEW** — 9 tests |
| `docs/adr-001-tailwind-v4-standard.md` | **NEW** — Tailwind v4+ decision |
| `docs/tailwind-v4-adoption-guide.md` | **NEW** — migration guide |
| `docs/icons-only-adoption.md` | **NEW** — icons package docs |
| `docs/cross-project-analysis.md` | **NEW** — this report |
| `README.md` | Updated component counts, catalog, theming docs |
| `AGENTS.md` | New components documented, Tailwind v4+ stance added |

### cqrs-htmx/adminui repo

No changes. adminui's icon system was handled by its own commit (`8091422`)
with self-contained inline Heroicons paths. Migration to Tailwind v4+ is a
future session.

---

## Metrics

| Metric | Before | After |
|--------|--------|-------|
| tc icons | 99 | **101** |
| tc components | 69 | **73** (+PageHeader, DefinitionList, ListNote, SidebarNav) |
| tc new tests | — | **+34** |
| tc test packages passing | 12 | 12 |
| tc lint issues | 0 | 0 |
