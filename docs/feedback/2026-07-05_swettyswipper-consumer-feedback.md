# Consumer Feedback: templ-components

**From:** SwettySwipperWeb integration session (2026-07-05)
**Perspective:** AI agent consuming the library in a server-rendered HTMX app
**Tone:** Honest, direct, grateful but critical where warranted

---

## What Works Superbly

### 1. Layout System (`layout.Base` + `PageProps`)

The `layout.Base` component with `PageProps` is excellent — it handles the full HTML shell (head, body, script tags, theme toggle) with a clean struct-based API:

```go
layout.Base(layout.PageProps{
    Title:        "Dashboard",
    HTMXVersion:  layout.HTMXVersion2_0_10,
    CurrentPath:  r.URL.Path,
    CSRFToken:    csrfToken,
    Nonce:        nonce,
}, content)
```

This eliminated 100+ lines of boilerplate HTML that we previously hand-rolled.

### 2. Navigation (`SimpleNav`)

`SimpleNav` with typed `NavLinkProps` is clean and composable. The active-link styling via `CurrentPath` matching is the right abstraction.

### 3. Display Components (`Badge`, `Table`, `EmptyState`, `StatusBadge`)

These are well-designed:

- `Badge` with typed `BadgeType` enum + map-based style lookup is the correct pattern
- `Table` with `TableRow` / `TableCell` is flexible enough for real-world data
- `EmptyState` handles the "no data" case elegantly

### 4. Icons System

The single-source-of-truth `iconPathData` map with `|`-separated multi-path icons is clean. `allIconNames()` being auto-generated prevents stale lists.

### 5. CSP Safety

Every inline `<script>` carries `nonce={ props.Nonce }`. This is the correct CSP pattern and makes the library usable in strict-CSP environments.

---

## What's Confusing or Hard to Discover

### 1. CDN Dependency is Hidden

**Problem:** `layout.Base` generates `<script src="https://cdn.jsdelivr.net/npm/htmx.org@version">`. This CDN dependency is not obvious from the API surface. We discovered it when our CSP blocked htmx and nothing worked.

**Impact:** Every consumer MUST include `https://cdn.jsdelivr.net` in their CSP `script-src`, but this isn't documented in the component props or the SKILL.md prominently enough.

**Ask:** (a) Add a CSP warning to `PageProps` godoc. (b) Make `HTMXCDN` prop more discoverable — it's the escape hatch for self-hosting but buried in generated code. (c) Consider making self-hosting the default and CDN the opt-in.

### 2. `PageProps` Doesn't Embed `BaseProps`

**Problem:** `layout.PageProps` is the ONE exception that doesn't embed `utils.BaseProps`. This means no `ID`, `Class`, `Attrs`, `AriaLabel` propagation on the layout component.

**Impact:** Can't add custom classes or data attributes to the `<body>` or root `<div>`. Minor but surprising given the universal pattern.

**Ask:** Document why this exception exists (it's noted in the skill but not in the godoc).

### 3. Component Discovery — Hard to Know What Exists

**Problem:** The library has 53 components across 9 packages. As a consumer, it's hard to discover what's available without reading the source. The README helps but doesn't cover all components.

**Impact:** We hand-rolled form inputs, buttons, cards, modals, and toasts because we didn't know the library had them. Only after a deep audit did we discover `forms`, `feedback.Toast`, `display.Button`, `display.Card`, `display.Modal`.

**Ask:** Generate a component catalog page (or at least a complete table in the README) with: component name, package, one-line description, thumbnail/screenshot. Consider a demo site.

### 4. `utils.Class()` vs `utils.Lookup()` — When to Use Which

**Problem:** The skill says "always go through `utils.Class(...)`" for class merging, but also "use `utils.Lookup(map, key, fallback)`" for style lookups. The relationship between these isn't clear.

**Ask:** Add a decision tree: "If merging multiple class strings → `utils.Class()`. If selecting a style variant from a map → `utils.Lookup()`. If both → `utils.Lookup()` first, then pass result to `utils.Class()`."

---

## What's Missing

### 1. Form Components — Exist But Undiscoverable

The `forms` package (`Input`, `Select`, `Textarea`, `Toggle`, `Radio`, `Combobox`, `Label`, `Form`) exists but our project hand-rolls 44+ raw form elements because we didn't know about it.

**Impact:** 44+ `<input>` / `<select>` / `<textarea>` tags with hand-rolled Tailwind classes instead of typed components with validation, a11y, and consistent styling.

**Ask:** This is the #1 missed opportunity. Make the forms package the flagship feature — it's what every server-rendered app needs. Add a forms demo page and put it front-and-center in the README.

### 2. Feedback Components — Toast Undiscovered

We hand-rolled a 30-line JavaScript toast system in `layout.templ`. The library has `feedback.Toast` which is CSP-safe, accessible, and consistent.

**Ask:** Make `feedback.Toast` discoverable. Consider adding it to `layout.Base` as an optional slot so consumers don't need to wire it separately.

### 3. No Pagination Component

We hand-roll pagination with `PaginationData` struct + raw HTML. A `pagination.Pagination` component with typed props (CurrentPage, TotalPages, OnChange HTMX attributes) would be valuable.

### 4. No Image Component

For a media comparison app, an `Image` component with lazy loading, aspect ratio, fallback src, and loading spinner would be valuable. We hand-roll `<img>` tags with Tailwind classes.

### 5. No Aspect Ratio Utility

We hand-roll `aspectRatioStyle(width, height)` returning `padding-top` CSS. A utility or component for aspect-ratio-preserving containers would help.

---

## What's Over-Engineered

### Nothing

The library is well-scoped. Every component serves a real need. The typed enum + map lookup pattern is the right level of abstraction. No over-engineering observed.

---

## Summary Scorecard

| Area                     | Rating | Notes                                                               |
| ------------------------ | ------ | ------------------------------------------------------------------- |
| Layout system            | ★★★★★  | Excellent, clean PageProps API                                      |
| Display components       | ★★★★☆  | Good coverage, missing Button/Card adoption                         |
| Forms package            | ★★★☆☆  | Exists but undiscoverable — #1 gap                                  |
| Feedback package         | ★★★☆☆  | Toast exists but unknown to consumers                               |
| Icons                    | ★★★★★  | Single source of truth, clean                                       |
| CSP safety               | ★★★★★  | Nonce on every script — excellent                                   |
| Component discovery      | ★★☆☆☆  | Hard to know what exists without reading source                     |
| CDN vs self-hosting      | ★★☆☆☆  | CDN dependency hidden, self-hosting path unclear                    |
| Documentation (SKILL.md) | ★★★★☆  | Great for authors, less useful for consumers discovering components |

---

## Top 3 Requests

1. **Make the forms package the flagship.** Generate a forms demo page, put it in the README, make it the first thing consumers see.
2. **Add a component catalog** — a table or page listing all 53 components with descriptions.
3. **Make self-hosting htmx the default** — CDN should be opt-in, not the silent default.

---

_This feedback is given with gratitude for a well-crafted component library. The critique is offered to make it even more adoptable._
