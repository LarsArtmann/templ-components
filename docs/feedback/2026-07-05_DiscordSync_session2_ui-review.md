# templ-components — Consumer Feedback Update (DiscordSync)

**Consumer:** [DiscordSync](https://github.com/LarsArtmann/DiscordSync) — Discord backup bot
**Version used:** v0.6.1
**Date:** 2026-07-05 (session 2 — full UI/UX integration review)
**Session scope:** Reviewed all 14 dashboard pages for component consistency, adopted shared helpers, investigated `forms` package adoption

---

## What This Session Covered

A full audit of every `.templ` file in DiscordSync's `internal/web/` package (14 pages, 9 supporting `.go` files). The goal was to find every place where the project hand-rolls something the library already provides, or where pages diverge from each other in inconsistent ways.

**8 integration fixes shipped.** All verified by build + full test suite + lint.

---

## New Findings

### 1. `forms` Package: investigated deeply, decided NOT to adopt — and here's exactly why

**Severity:** Design decision (not a bug)

DiscordSync uses 7 filter forms across pages (messages, search, members, threads, voice-states, attachments, events). Each uses hand-rolled `filterForm` + `filterSelect` + `selectOption` helpers (168 lines total in `filters.templ`).

The `forms` package (`forms.Form` + `forms.Select`) was investigated as a replacement. I read the source of both `forms.Form`, `forms.Select`, `FormFieldWrapper`, `Label`, and `baseInputClass`. The conclusion:

**The HTMX selector works.** `forms.Select` renders `<label>` → `<select>` → optional error/help with NO wrapper div. The `<select>` is a direct descendant of `<form>`. htmx's `from:find select` finds it at any depth.

**But the layout is incompatible for our use case:**

| Concern         | Our `filterForm`                                     | `forms.Form` + `forms.Select`                       |
| --------------- | ---------------------------------------------------- | --------------------------------------------------- |
| **Layout**      | `flex flex-wrap gap-3` → horizontal row of dropdowns | `space-y-6` default → vertical stack                |
| **Input style** | `border-gray-300` (traditional border)               | `ring-1 ring-inset shadow-xs` (inner shadow ring)   |
| **Label style** | `text-xs text-gray-500` (muted, compact)             | `text-sm text-gray-900 dark:text-white` (prominent) |

Our filter bar is a compact horizontal row — the defining UX pattern of every list page. `forms.Form` defaults to a vertical form layout designed for settings pages and create/edit forms. Overriding `forms.Form`'s default class with `flex flex-wrap gap-3` and reconciling the input/label styling differences would eat most of the simplification.

**Verdict:** The custom helpers are purpose-built for the horizontal filter-bar pattern and fully tested. The `forms` package targets a different use case (vertical forms with validation). Both earn their keep.

**Library suggestion:** Consider documenting this distinction: "`forms.Form` is for vertical data-entry forms. For horizontal filter bars with HTMX auto-submit, consumers typically build a thin `filterForm` helper." Or add a `forms.InlineForm` variant with flex layout + compact labels.

### 2. `display.Grid` exists but DiscordSync doesn't use it

**Severity:** Medium (consistency)

DiscordSync hand-rolls grid classes everywhere:

```html
<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4"></div>
```

The library has `display.Grid(display.GridProps{Cols: display.GridCols3})` with typed column enums and responsive breakpoints. This was flagged as resolved in the Overview feedback, but DiscordSync never adopted it. Same discoverability gap.

**Status:** DiscordSync should adopt `display.Grid` for the dashboard stat-card rows.

### 3. `display.Table` exists but DiscordSync built `listTableWithHeader`

**Severity:** Low (already addressed in this session)

DiscordSync has `listTableWithHeader(ariaLabel, headers)` in `filters.templ` — a shared table wrapper with `<thead>` from a `[]string`. This session standardized all pages to use it (guilds card was the last holdout with a raw `<table>`).

The library's `display.Table` already exists with `TableProps{Headers, Rows, Caption}`. DiscordSync's custom helper is thinner (no `Rows` type — children render `<tr>` directly), which is actually more flexible for the templ pattern.

**Verdict:** The custom helper is the right call for consumer templates. The library component's `Rows []TableRow` abstraction is less ergonomic than `for _, msg := range messages { <tr>...</tr> }` when each row has custom cell rendering.

### 4. `display.Tabs` — adopted since last feedback

The previous feedback said "DiscordSync should adopt `display.Tabs` for message detail." Verified: `message_detail.templ` now uses `display.Tabs` with `messageTabs()` builder that dynamically includes only tabs with content (Edits, Attachments, Embeds, Extras). Working correctly.

### 5. `feedback.ProgressBar` — adopted since last feedback

Verified: `backfill.templ` uses `feedback.ProgressBar` for channel completion progress. Working correctly.

---

## Persistent Discoverability Gap (unchanged from session 1)

The #1 problem remains: **consumers don't know what exists.** In this session alone:

| Component                     | Status | Discovered?                                     |
| ----------------------------- | ------ | ----------------------------------------------- |
| `forms.Form` / `forms.Select` | Exists | Investigated, decided against (layout mismatch) |
| `display.Grid`                | Exists | Found but not yet adopted                       |
| `display.Table`               | Exists | Custom helper preferred (more flexible)         |
| `display.Tabs`                | Exists | Already adopted since last session              |
| `feedback.ProgressBar`        | Exists | Already adopted since last session              |

The SKILL.md now has a consumer-facing component catalogue (76 components across 9 packages). This is a major improvement. But the DiscordSync project's AGENTS.md doesn't cross-reference which library components are actually used vs available-but-not-adopted.

---

## Scorecard Update

| Dimension             | Session 1 | Session 2 | Change                                    |
| --------------------- | --------- | --------- | ----------------------------------------- |
| Component quality     | 8/10      | 8/10      | —                                         |
| API design            | 9/10      | 9/10      | —                                         |
| Accessibility         | 9/10      | 9/10      | —                                         |
| CSP safety            | 10/10     | 10/10     | —                                         |
| Documentation (skill) | 7/10      | 8/10      | +1 (consumer catalogue added to SKILL.md) |
| Component coverage    | 6/10      | 6/10      | —                                         |
| Theming               | 9/10      | 9/10      | —                                         |
| **Overall**           | **8/10**  | **8/10**  | —                                         |

---

## Recommendations for templ-components

1. **Document the `forms` vs `filterForm` distinction.** `forms.Form` is for vertical data-entry forms. Horizontal HTMX filter bars are a different pattern. Either document this or add `forms.InlineForm`.

2. **Reconsider `display.Table`'s `Rows` abstraction.** Consumer templates need `for _, item := range items { <tr>...</tr> }` with custom cells. The `Rows []TableRow` type forces an intermediate conversion. Consider a `display.TableSlots` variant where children render `<tr>` directly (like DiscordSync's `listTableWithHeader`).

3. **Add a "Which components does this project use?" cross-reference.** Consumers should be able to quickly audit which library components they've adopted vs hand-rolled. A grep-able table in the consumer's AGENTS.md would help.

4. **The component catalogue in SKILL.md is great but 76 entries is a lot to scan.** Consider grouping by use case: "Dashboard pages" (StatCard, Grid, Card, ProgressBar), "List pages" (Table, Badge, Avatar, Pagination), "Forms" (Input, Select, Form), "Feedback" (Toast, Alert, Skeleton, Spinner).
