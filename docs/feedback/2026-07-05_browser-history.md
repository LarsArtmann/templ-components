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

# Consumer Feedback — templ-components

**From:** browser-history project (github.com/larsartmann/browser-history)
**Date:** 2026-07-05
**Version used:** v0.6
**Consumer:** Crush (AI assistant) + Lars

---

## What Works Great

### `layout.Base` + `layout.PageProps` composition

```templ
@layout.Base(layout.PageProps{
    Title:           "Browser History",
    Description:     "Where your browsing time goes",
    BodyClass:       "bg-gray-950 text-gray-200 min-h-screen antialiased",
    HeadContent:     headContent(data.Nonce),
    Nonce:           data.Nonce,
    SecurityHeaders: true,
}) {
    // page content
}
```

Clean composition. `HeadContent` lets us inject CSP-nonce'd scripts. `SecurityHeaders: true` adds defense-in-depth meta tags. The `Nonce` field flows from server → template → script tags automatically.

### `display.Badge` with typed `BadgeType`

```go
var categoryBadgeTypeMap = map[visit.Category]display.BadgeType{
    "development":   display.BadgeInfo,
    "news":          display.BadgeSuccess,
    "entertainment": display.BadgeError,
    "social":        display.BadgeWarning,
    "shopping":      display.BadgePrimary,
}
```

Type-safe badge variants via map + fallback is the right pattern. Much better than a switch statement. The skill correctly recommends this approach and we followed it.

### `gray-*` color palette for dark mode

The skill says "use gray-_ not slate-_ for dark mode" — we corrected our entire dashboard from slate to gray. The result looks more neutral and professional. Good call.

### `motion-reduce:*` on all transitions

```html
class="hover:bg-gray-800 motion-reduce:transition-none motion-reduce:duration-0"
```

Accessibility-first. The skill mandates it, we applied it everywhere. Correct default.

---

## Pain Points

### 1. Tailwind v4 CSS-first config vs. our `tailwind.min.js` CDN approach

**Problem:** Our dashboard loads Tailwind via `<script src="/static/tailwind.min.js">` (the Play CDN). The skill describes Tailwind v4 CSS-first config (`@import "tailwindcss"`). We're still on the Play CDN because migrating to CSS-first requires a build step we haven't set up.

**Impact:** The Play CDN is slow (downloads ~300KB JS, compiles CSS client-side on every page load). It also requires `'unsafe-inline'` in the CSP for `style-src`, which weakens our security posture.

**Suggestion:** Add a migration recipe to the skill: "How to move from Play CDN to Tailwind v4 CSS-first in a Go + templ project." Show the `flake.nix` or Makefile target that generates the CSS file.

### 2. No loading/error feedback components documented for server-rendered apps

**Problem:** When our `/extract` endpoint fails, the dashboard does `http.Error(w, "query error", 500)` — a blank white page. The skill mentions `feedback.Toast` and `feedback.Alert` but we couldn't find clear examples of wiring them with HTMX swap targets in a server-rendered (non-SPA) context.

**Impact:** Dashboard error UX is poor. We know `feedback.Toast` exists but the integration path from "Huma handler returns error" → "show Toast in dashboard" is unclear.

**Suggestion:** Add a recipe: "Server-rendered HTMX error feedback loop" showing:

1. Huma handler returns error
2. Server renders `feedback.Alert` component as HTML fragment
3. HTMX swaps it into a target div
4. Toast auto-dismisses after 5s

### 3. CSP nonce propagation — works but requires manual wiring

**Problem:** The nonce flows through correctly (`PageProps.Nonce` → template `nonce={ nonce }`), but every `<script>` tag needs it manually applied. If you forget one, the browser blocks the script silently.

**Impact:** We had to audit every script tag after enabling CSP. Not a library bug — just a DX improvement opportunity.

**Suggestion:** Consider a `templ` helper: `@layout.Script(nonce, src, attrs...)` that auto-injects the nonce, so you can't forget it.

### 4. Component coverage feels incomplete for dashboards

**What we have:** Badge, Button, layout.Base, display components.

**What we need for a dashboard:** Cards with stats, tables with sortable columns, pagination, filter dropdowns, loading skeletons, empty states, time-range pickers.

**Impact:** We hand-rolled all dashboard HTML instead of composing components. The result works but is inconsistent and harder to maintain.

**Suggestion:** Prioritize dashboard-oriented components:

- `card.Stat(label, value, icon)` — stat card with label, big number, trend indicator
- `table.Data(columns, rows, opts...)` — sortable data table with pagination
- `state.Empty(title, description, action)` — empty state with CTA
- `state.Loading(skeletonType)` — loading skeleton

---

## Minor Notes

- **The `map + fallback` pattern recommendation** is excellent guidance — it's more extensible than switch and the skill explains when to use each clearly.
- **`BaseProps` embed** — good for composition. We don't use it yet but the pattern is sound.
- **No accessibility issues found** — the `motion-reduce:*` and semantic HTML patterns produce accessible output.

---

## Summary

templ-components provides a **solid foundation** (layout, display, badge) with excellent accessibility and CSP guidance. The main gaps are **dashboard-specific components** (stat cards, data tables, empty/loading states) and **migration recipes** (Play CDN → Tailwind v4 CSS-first, server-rendered HTMX feedback loops). The design philosophy (gray-\* colors, map+fallback, motion-reduce, nonce-everywhere) is correct and we follow it happily.

---

## Appendix: Resolution Status (2026-07-05)

| Pain point                                                       | Status                 | Resolution                                                                                                                                                                                                                                                                                                                    |
| ---------------------------------------------------------------- | ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1. Tailwind v4 CSS-first vs Play CDN approach                    | **RESOLVED**           | Added migration recipe: `docs/migration/play-cdn-to-tailwind-v4.md` — 7-step guide covering CLI install, `@source` scanning of vendored templ-components, `@theme` overrides, flake.nix integration, and CSP tightening from `style-src 'unsafe-inline'` to `style-src 'self'`.                                               |
| 2. No loading/error feedback documented for server-rendered apps | **RESOLVED**           | Added recipe: `docs/recipes/server-rendered-htmx-error-feedback.md` — documents the 3 render modes (full error page via `errorpage.ErrorHandler`, inline alert fragment via `errorpage.ErrorAlert`, toast via `htmx.GlobalErrorHandling` + `feedback.ToastContainer`), with go-error-family integration and CSP nonce wiring. |
| 3. CSP nonce propagation — manual wiring                         | **RESOLVED**           | Added `layout.Script(nonce, src, attrs)` — CSP-safe `<script src>` tag that auto-injects the nonce. Prevents the common bug of forgetting `nonce={...}` on a script tag. Optional `attrs` covers `async`, `defer`, `type="module"`, `crossorigin`.                                                                            |
| 4. Component coverage incomplete for dashboards                  | **PARTIALLY RESOLVED** | Added: `StatCard.Href` (clickable stat cards), `display.Grid` (responsive grid), `feedback.SkeletonCardGrid` (loading card grid). Still missing: sortable data table with typed columns, filter dropdowns, time-range pickers. These are deferred to future work.                                                             |
