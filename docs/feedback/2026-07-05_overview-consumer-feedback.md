# templ-components — SDK Feedback from Overview

**Consumer:** [Overview](https://github.com/larsartmann/overview) — local project dashboard (read-only, HTMX + SSE)
**Date:** 2026-07-05
**Version used:** v0.6.1
**Session:** Full page layout, StatCard, StatusBadge, EmptyState, forms, errorpage, ThemeToggle, Pagination

---

## What worked superbly

### 1. `errorpage.ErrorPage` + `errorpage.FromError` — family-aware error rendering

The integration with `go-error-family` is the best part of this library for Overview. `errorpage.FromError(err)` reads the error family (Rejection/Transient/Infrastructure/Corruption) and renders a styled page with the right color, title, message, and fix suggestion. API paths get structured JSON, browser paths get family-colored HTML. One function call replaces what would be a switch statement with 5 branches.

```go
writeErrorPage(w, r, 0, errorpage.FromError(snap.LastError))
```

### 2. `layout.Base(props)` + `DefaultPageProps()` — sensible defaults

The page shell handles `<head>`, meta tags, theme script, footer. `DefaultPageProps()` gives meaningful non-zero values for every field. Override only what you need (Title, Description, HeadContent, Footer). The `HeadContent templ.Component` slot for injecting extra `<head>` content (Tailwind CDN, custom scripts) is the right extension point.

### 3. `StatCard`, `StatusBadge`, `EmptyState` — well-designed primitives

Clean props structs, typed enums (no raw strings for variants), sensible defaults. `StatusBadge(string(p.Activity))` just works. `EmptyState` with Title + Description handles the "no results" state elegantly.

### 4. `forms.Input` with typed `InputType` enum

`forms.InputSearch`, `forms.InputText`, etc. — typed enums over raw strings. The `BaseProps` embed gives ID/Class/Attrs/AriaLabel for free. The forms package is the right level of abstraction for server-rendered inputs.

### 5. `layout.ThemeToggle` + theme script — dark mode done right

Class-based dark mode strategy (`@custom-variant dark (&:where(.dark, .dark *))`), persisted to localStorage + server-side cookie. One component, works end-to-end. No flash of wrong theme on load.

### 6. `navigation.Pagination` — correct, accessible pagination

`CurrentPage`, `TotalPages`, `BaseURL`, `MaxVisible` — handles ellipsis, prev/next, active state. Used directly for project pagination.

---

## Pain points

### 1. `HTMXVersion = ""` to disable CDN injection — undocumented and non-obvious

**Severity:** High (integration friction)

`DefaultPageProps()` sets `HTMXVersion: HTMXVersion2_0_10` and `HTMXUseSRI: true`, which auto-injects:

```html
<script
  src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.10"
  integrity="..."
  crossorigin="anonymous"
></script>
```

When using `cqrs-htmx`'s embedded `/htmx.js`, you must **explicitly blank `HTMXVersion`** to suppress this. I had to read the `base.templ` source to discover this behavior — it's not documented anywhere.

A typed enum would be far clearer:

```go
type HTMXSource string
const (
    HTMXSourceNone   HTMXSource = "none"     // don't inject anything
    HTMXSourceCDN    HTMXSource = "cdn"      // current default behavior
    HTMXSourceLocal  HTMXSource = "local"    // expects caller to serve /htmx.js
)
```

Or at minimum, document in `DefaultPageProps()` godoc: "Set `HTMXVersion = ""` to disable auto-injection."

### 2. `CSSPath` defaults to `"/app.css"` — another silent auto-inject

**Severity:** Medium (integration friction)

`DefaultPageProps()` sets `CSSPath: "/app.css"`, which injects `<link rel="stylesheet" href="/app.css">`. Apps using Tailwind via browser CDN don't have an `app.css`. I had to explicitly set `CSSPath = ""` to suppress a 404.

Same suggestion: document this default, or make `CSSPath` default to `""` (opt-in rather than opt-out).

### 3. Skeleton components exist but are undiscoverable

**Severity:** Medium (discoverability)

Overview hand-rolled a `skeletonGrid()` component (6 placeholder cards with `animate-pulse`). I later discovered that `feedback.Skeleton(variant)` and `feedback.SkeletonGroup(variants)` **already exist** with `SkeletonCard`, `SkeletonText`, `SkeletonAvatar` variants.

This is a documentation problem, not a missing feature. The README component catalogue doesn't prominently feature the loading/skeleton family. Consider:

- Adding `SkeletonCard` to the README "Components" table
- A `feedback.SkeletonCardGrid(count int)` convenience function for the common grid pattern
- Cross-referencing from `EmptyState` ("for loading state, see Skeleton")

### 4. `StatCard` has no link support

**Severity:** Medium (ergonomic)

Overview's stats row makes each StatCard a clickable filter link:

```go
<a href="/?activity=active" class="block" hx-get="/projects?activity=active" hx-target="#content-area" ...>
    @tc.StatCard(tc.StatCardProps{Value: "42", Label: "Active"})
</a>
```

The raw `<a>` wrapper works but bypasses `utils.Class` merging and accessibility patterns. A `StatCardProps.Href string` field (or a `StatCardLink` variant) would handle this natively, with proper `role`, keyboard nav, and class merging.

### 5. `SimpleNavProps` has no right-side slot

**Severity:** Low-medium (ergonomic)

Overview renders `layout.ThemeToggle` outside `SimpleNav` in a separate flex container:

```go
@navigation.SimpleNav(navigation.SimpleNavProps{BrandText: "Overview", BrandHref: "/", CurrentPath: "/"})
// ... then separately:
@layout.ThemeToggle("Toggle theme", "")
```

The parent `NavProps` has `RightItems templ.Component`, but `SimpleNavProps` doesn't expose it — `SimpleNav` never forwards `RightItems` to `Nav`. Adding a `RightItems templ.Component` field to `SimpleNavProps` would solve this:

```go
type SimpleNavProps struct {
    utils.BaseProps
    BrandText   string
    BrandHref   string
    Links       []NavLinkProps
    RightItems  templ.Component  // forwarded to Nav.RightItems
    CurrentPath string
    Sticky      bool
}
```

### 6. No responsive grid helper

**Severity:** Low (DRY)

The pattern `<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">` appears in multiple Overview templates. A `layout.Grid` or `display.CardGrid` with configurable column counts would reduce repetition:

```go
@display.CardGrid(display.CardGridProps{Columns: display.GridCols3, Children: cards})
```

---

## Summary

templ-components is a **solid foundation** for server-rendered Go UIs. The component primitives (StatCard, StatusBadge, EmptyState, forms) are well-designed with typed enums and BaseProps composition. The errorpage + error-family integration is exceptional.

The main friction is **undocumented auto-injection defaults** (`HTMXVersion`, `CSSPath`) — these silently add CDN links and 404s until you discover and blank them. Making the defaults explicit (or opt-in) would save every new consumer from reading source.

The component catalogue is comprehensive but discoverability of the loading/skeleton family needs improvement — I hand-rolled something that already existed.
