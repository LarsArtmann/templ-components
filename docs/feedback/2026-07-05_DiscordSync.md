<!-- AUTO-UPDATED 2026-07-10: Retrospective status overlay -->

> ## 🔔 Update Notice — 2026-07-10
>
> This report is **historical**. Many items listed as "open", "todo", or "broken" below
> have since been **fixed and verified**. Do not act on open items without first checking
> [TODO_LIST.md](../../TODO_LIST.md) for current status.
>
> **Key fixes completed since this report:**
>
> - ✅ All 7 P0 bugs fixed (InlineLoadingOverlay a11y, SanitizeID mismatch, FromError fallback,
>   Footer BaseProps, ErrorPage/NotFound404 `<main>` landmark, CSRFTokenName, grid-rows verified)
> - ✅ `encoding/json/v2` purged from all production code + pre-commit guard added
> - ✅ Motion constants centralized in `utils/motion.go`, wired into 13 components
> - ✅ `FamilyFromErrorFamily` → `FromErrorFamily` (old name kept as deprecated alias)
> - ✅ `icons.IconRTL()` + CSS for directional icon RTL mirroring
> - ✅ 33 regression tests added (htmx, errorpage, layout, navigation, feedback, display)
> - ✅ Dark golden test infrastructure (badge/card/button)
> - ✅ CHANGELOG consolidated, ROADMAP updated, migration guide created
> - ✅ All 14 packages pass, 0 lint issues
>
> **Canonical source of truth:** [TODO_LIST.md](../../TODO_LIST.md) (52 items, 37 ✅ done, 12 deferred/blocked)

---

# templ-components — Consumer Feedback (DiscordSync)

**Consumer:** [DiscordSync](https://github.com/LarsArtmann/DiscordSync) — Discord backup bot
**Version used:** v0.6.1
**Usage depth:** Moderate — 7 packages used: display, feedback, htmx, icons, layout, navigation, errorpage. 11 pages rendered with templ + Tailwind v4.
**Date:** 2026-07-05

---

## What Works Superbly

### 1. `utils.BaseProps` embed is the right composition model

Every props struct embeds `utils.BaseProps` → auto-satisfies `utils.ComponentProps` interface → `GetBaseProps()`/`SetBaseProps()` promoted via pointer receivers. This means ID, Class, Attrs, AriaLabel, Nonce are available on every component without code duplication. Consumers override via `@theme` CSS variables.

### 2. `display.StatCard` — the dashboard workhorse

```go
@display.StatCard(display.StatCardProps{
    Label: "Messages", Value: "1,234", Icon: icons.ChatBubbleLeft,
    Description: "Total synced", Href: "/messages",
})
```

One component, responsive grid, icon, description, link. Powers our entire 13-card dashboard.

### 3. `display.Avatar` — clean and flexible

```go
@display.Avatar(display.AvatarProps{
    Src: user.AvatarURL, Alt: user.DisplayName, Size: display.AvatarSizeSM,
})
```

Fallback initials, size variants, clean rounded styling. Used in messages, members, voice states, search, message detail.

### 4. `errorpage.WriteError` — one-call error pages

```go
errorpage.WriteError(w, http.StatusInternalServerError, "Failed to load page", logger)
```

Renders a full styled error page with consistent branding. Used in every web handler error path.

### 5. `htmx.GlobalErrorHandling` + `htmx.LoadingIndicator` + `feedback.ToastContainer`

Layout-level wiring that gives every page HTMX error handling, loading spinners, and toast notifications for free. One-time setup in `layout.templ`, never touched again.

### 6. `navigation.Breadcrumbs` for hierarchical navigation

Message detail → channel → guild. Clean, typed props, active-path highlighting.

### 7. CSP-safe by construction

Every inline `<script>` carries `nonce={ props.Nonce }`. This lets us ship a strict Content-Security-Policy without forking the library. The `nonceFromContext(ctx)` pattern in our layout extracts the per-request nonce and threads it to every component.

### 8. Tailwind v4 CSS-first theming

Consumers override colors via `@theme { --color-blue-600: #...; }` in their CSS — no Go code change. Dark mode via class strategy with `layout.ThemeToggle`.

---

## What's Painful

### 1. No data table / sortable table component

DiscordSync's message browser, member list, and voice states all use hand-built HTML `<table>` elements with manual column headers, sorting, and styling. This is the single biggest gap.

```html
<!-- What we write today: 40 lines of raw HTML per table -->
<table class="min-w-full divide-y divide-gray-200">
  <thead>
    <tr>
      <th class="px-4 py-2 text-left text-xs font-medium...">Author</th>
      ...
    </tr>
  </thead>
  <tbody>
    for _, msg := range messages {
    <tr>
      <td class="px-4 py-2 text-sm...">...</td>
    </tr>
    }
  </tbody>
</table>
```

**Suggestion:** Add `display.Table` or `display.DataTable` with:

- Typed column definitions (`display.Column{Header, Field, Sortable}`)
- Row rendering via `templ.Component` slot
- Responsive overflow wrapper
- Zebra striping option
- Empty state integration

### 2. No modal/drawer for detail views — **✅ EXISTS since v0.8.0**

Message detail is a full page (`GET /messages/{id}`). A slide-out drawer or modal would be better UX — see the detail without losing the list context.

> **✅ RESOLVED:** `display.Modal` and `display.Drawer` both exist with focus trap, aria sync, backdrop, keyboard nav.

### 3. No pagination component that works with cursor-based pagination — **✅ EXISTS since v0.10.0**

> **✅ RESOLVED:** `navigation.LoadMore` exists with cursor pagination, `hx-get`/`hx-swap`, and optional `InfiniteScroll` (`hx-trigger="revealed"`).

### 4. The `_templ.go` gitignore problem

`_templ.go` files are gitignored in DiscordSync (generated at build time). But the templ-components skill says "commit `*_templ.go`" because it's a library and the Go module proxy doesn't run `templ generate`. This is correct for the library but DIFFERENT for consumers.

The BuildFlow pre-commit hook re-appends `*_templ.go` to `.gitignore` in consumer repos, which is harmless for already-tracked files but makes NEW generated files invisible to `git status` until `git add -f`.

**Suggestion:** Document the consumer-side pattern clearly: "Consumers should gitignore `*_templ.go` and generate at build time. The library commits them because the module proxy needs them."

### 5. No accordion/toggle for collapsible sections

Message detail has edit history, reactions, embeds, attachments — all stacked vertically. An accordion (`<details>`-based) would improve scannability. The library may have this (the skill mentions it), but we haven't adopted it.

### 6. No badge with count/number overlay

We show reaction counts and poll answer counts. Currently using `display.Badge` with text like "👍 5". A count badge (small number overlay on an icon) would be cleaner.

---

## What's Missing

### 1. A generic list/ul component with consistent spacing

```go
@display.List(items, renderItem func(item T) templ.Component)
```

Every page has `for _, x := range items { @template(x) }`. A generic list with consistent spacing, dividers, and empty-state integration would reduce boilerplate.

### 2. A description list / key-value component

We use `<dl>` manually for attachment metadata (filename, size, content type, dimensions). `display.DefinitionList` + `display.DefinitionItem` exist (we use them) — but they need a `display.DefinitionGrid` wrapper for responsive 2-column layout.

### 3. A copy-to-clipboard button — **✅ EXISTS since v0.10.0**

> **✅ RESOLVED:** `display.CopyButton` exists with clipboard write, "Copied!" feedback, `<a>` and `<button>` variants, CSP-safe singleton script.

Message IDs, attachment hashes, and error codes are frequently copied. A `display.CopyButton` with clipboard integration + "Copied!" feedback would be useful.

### 4. A relative time component

"2 hours ago", "5 minutes ago", "yesterday". Discord messages have timestamps that are more readable as relative time. We format these manually.

---

## What DiscordSync Should Adopt (Self-Note)

These components exist in the library but DiscordSync doesn't use yet:

| Component                     | Where to use                                                                   |
| ----------------------------- | ------------------------------------------------------------------------------ |
| `display.Tabs`                | Message detail: Content / Edits / Attachments / Reactions tabs                 |
| `feedback.ProgressBar`        | Already used for backfill progress; could use for attachment download progress |
| `display.Badge` with variants | Thread archive status, download status, projection health                      |
| `layout.ThemeToggle`          | Already wired in layout.templ                                                  |

---

## Skill Feedback

The skill file (`/skill/SKILL.md`) is **excellent for library authors** but could be better for consumers.

### Good

- The component anatomy skeleton is the authoritative reference
- Decision trees (map vs if-branch, when to add icon, theming) are practical
- Mandatory conventions (dark mode gray-\*, motion safety, CSP nonce) catch real issues
- Anti-patterns to refuse on review is a great review checklist
- Testing expectations table helps understand what to test

### Could Improve

- **No consumer-side component catalogue** — the skill tells you how to BUILD components, not WHICH components exist. Consumers need a "what's available" reference.
- **No `_templ.go` consumer guidance** — the skill says "commit \_templ.go" but consumers should gitignore them. This distinction isn't clear.
- **No Tailwind v4 consumer setup guide** — the skill references `docs/tailwind-v4-adoption-guide.md` but doesn't include the key setup steps inline.
- **No "adopting existing components" section** — everything is about creating new components. Consumers want to know how to find, import, and use existing ones.

---

## Summary Scorecard

| Dimension             | Score    | Notes                                                     |
| --------------------- | -------- | --------------------------------------------------------- |
| Component quality     | 8/10     | What exists is well-built; gaps in table/pagination/modal |
| API design            | 9/10     | BaseProps embed, typed enums, map lookups — excellent     |
| Accessibility         | 9/10     | ARIA, keyboard nav, motion-reduce on everything           |
| CSP safety            | 10/10    | Nonce on every script, no eval, no inline handlers        |
| Documentation (skill) | 7/10     | Great for authors; thin for consumers                     |
| Component coverage    | 6/10     | Missing data table, modal, cursor pagination, accordion   |
| Theming               | 9/10     | Tailwind v4 CSS variables are the right approach          |
| Overall               | **8/10** | Solid foundation; needs more data-heavy components        |

---

## Appendix: Resolution Status (2026-07-05)

| Pain point                                  | Status                           | Resolution                                                                                                                                                                                                                                                                                                                                                                       |
| ------------------------------------------- | -------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1. No data table / sortable table component | **EXISTS** (discoverability gap) | `display.Table` already exists with typed `TableProps{Headers, Rows}`, striping, hover, caption, bordered, responsive overflow. **Not sortable** — sorting is deferred. The consumer didn't know it existed.                                                                                                                                                                     |
| 2. No modal/drawer for detail views         | **EXISTS** (discoverability gap) | `display.Modal` (focus trap, Escape, backdrop, 5 sizes) and `display.Drawer` (left/right slide, focus trap) both exist. Consumer didn't know.                                                                                                                                                                                                                                    |
| 3. No cursor-based pagination               | **PARTIALLY EXISTS**             | `navigation.Pagination` exists but is page-number-based (CurrentPage, TotalPages, MaxVisible). Cursor-based pagination (`?cursor=...` + "Load More" HTMX pattern) is not implemented. Documenting the recommended pattern is a good next step.                                                                                                                                   |
| 4. `_templ.go` gitignore problem            | **KNOWN ISSUE**                  | Documented in AGENTS.md as the "BuildFlow gotcha." The `.gitignore` has `!*_templ.go` (line 2, unignore) overridden by `*_templ.go` (line 32, re-ignore from BuildFlow). Workaround: `git add -f` for new components. Consumer-side guidance: consumers should gitignore `*_templ.go` and generate at build time — the library commits them because the module proxy needs them. |
| 5. No accordion                             | **EXISTS** (discoverability gap) | `display.Accordion` already exists with typed `AccordionProps{Items}`, open/closed state, keyboard nav. Consumer didn't know.                                                                                                                                                                                                                                                    |
| 6. No badge with count/number overlay       | **OPEN**                         | Not implemented. `display.Badge` with text ("👍 5") is the current workaround. A count-badge overlay on an icon would be a new component.                                                                                                                                                                                                                                        |
| 7. Generic list component                   | **OPEN**                         | Not implemented. The `{ children... }` + `for` pattern is the current approach. A `display.List[T]` with typed items would reduce boilerplate.                                                                                                                                                                                                                                   |
| 8. Description grid wrapper                 | **OPEN**                         | `display.DefinitionList` exists but has no responsive 2-column grid wrapper. A `display.DefinitionGrid` composing `DefinitionList` with `display.Grid` could solve this.                                                                                                                                                                                                         |
| 9. Copy-to-clipboard button                 | **OPEN**                         | Not implemented. Would need minimal JS for clipboard API + "Copied!" feedback.                                                                                                                                                                                                                                                                                                   |
| 10. Relative time component                 | **OPEN**                         | Not implemented. Would need a `display.RelativeTime(timestamp)` component with optional JS auto-refresh.                                                                                                                                                                                                                                                                         |

### Discoverability gap summary

4 of 6 "missing" components already exist (Table, Modal, Drawer, Accordion). This is a **documentation problem**, not a missing-feature problem. The README component catalogue and the SKILL.md need to better surface existing components to consumers.
