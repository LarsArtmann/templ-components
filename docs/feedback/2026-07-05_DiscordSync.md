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

### 2. No modal/drawer for detail views

Message detail is a full page (`GET /messages/{id}`). A slide-out drawer or modal would be better UX — see the detail without losing the list context.

**Suggestion:** If `display.Modal` or `layout.Drawer` doesn't exist yet, it's the most-requested missing component for data-heavy apps.

### 3. No pagination component that works with cursor-based pagination

`navigation.Pagination` (if it exists) is page-number-based. DiscordSync uses cursor pagination (`?cursor=...`) with a "Load More" button + HTMX infinite scroll. We build this manually with `display.Button`.

**Suggestion:** Add `navigation.LoadMore` or `navigation.CursorPagination` that handles the HTMX swap pattern. Even just documenting the recommended pattern would help.

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

### 3. A copy-to-clipboard button

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
