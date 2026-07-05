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
