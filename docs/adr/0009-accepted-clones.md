# ADR 0009: Accepted Code Duplication

## Status

Accepted

## Context

The `templ-components` library is checked via `art-dupl --semantic`. After
systematic extraction passes (commits `4c2b00e`, `7671XXX`, and this session),
the following sub-templates were created to eliminate harmful duplication:

- `utils.DismissButton(bgClass, textClass)` — shared dismiss control for Alert + ErrorAlert
- `errorpage.errorBody(...)` — shared badge + title + message block
- `errorpage.errorHeader(...)` — shared flex-row icon + content header for ErrorDetail + ErrorPage
- `errorpage.goBackScript(nonce)` — shared history.back() JS for ErrorPage + NotFound404
- `errorpage.actionLinkBody(text)` — shared action link body (text + ArrowRight icon)
- `display.overlayPanel(...)` — shared panel body (merged into overlayShell)
- `display.overlayShell(...)` — now owns the complete overlay: shell + backdrop + panel + script
- `display.copyButtonContent(icon, label)` — shared icon + status span for CopyButton branches
- `display.definitionDetailContent(item)` — shared DetailComponent-or-Detail fallback
- `feedback.skeletonContainer(layoutClass, label)` — shared aria shell for SkeletonGroup + SkeletonCardGrid

Despite these extractions, `art-dupl` still reports clones at lower thresholds.
Each is examined below. The goal is **zero harmful duplication** — not zero
report lines.

## Accepted Production Clones

### 1. `display/card.templ:246-254` vs `errorpage/shared.templ:54-62` — flex-row icon + content (t=8)

`errorHeader` already consolidates ErrorDetail and ErrorPage into one template.
StatCard's `statCardInner` remains structurally similar because it uses:

- **`items-center`** (StatCard vertically centers icon vs. ErrorDetail's `items-start`)
- **`<dl class="min-w-0 flex-1">`** (semantic HTML — `statCardFigures` renders `<dt>`/`<dd>` pairs; cannot use `<div>`)
- **`flex h-12 w-12 ... bg-blue-50`** (hardcoded icon box vs. family-styled icon container)

An `IconRow` helper would need 8+ parameters (alignment, gap, icon container
class, icon body, content element tag, content class, content body) — more
indirection than the 8 duplicated lines justify.

**Why not lazy**: The extraction was attempted with `IconRow` in `utils/` and
rejected because the `<dl>` vs `<div>` semantic requirement and different
alignment made the call sites longer than the inline code.

### 2. `display/copy_button.templ:63-78` vs `:79-94` — `<a>` vs `<button>` (t=8)

Both branches must carry `data-tc-copy`, `data-tc-copy-label`, `id`, `class`,
`aria-label`, and `props.Attrs` directly on the interactive element — the
singleton copy script uses `closest('[data-tc-copy]')` on the document.

The `copyButtonContent(icon, label)` sub-template already factors out the
inner icon + status span. The outer attribute declarations must remain inline
because templ requires attributes on the element they decorate.

**Why not lazy**: Templ's DSL has no mechanism to push attributes from a
wrapper element onto a child `<a>` or `<button>`. This is a structural
constraint of the framework, not a missed extraction.

### 3. `errorpage/erroralert.templ:28-34` vs `errorpage/shared.templ:122-130` — badge row (t=7)

Both render `<div class="flex items-center gap-2">` + optional content +
`@familyBadge`. The "optional content" differs completely:

- ErrorAlert: an `<h3>` title with `style.Text` classes
- codeAndFamilyBadge: a `<code>` element with `style.AccentBG` + `style.AccentText`

Extracting a `badgeRow(leftContent, family, style)` template would add a third
template layer for a 1-line flex-row idiom (`flex items-center gap-2`).

**Why not lazy**: The duplicated portion is a single `<div class="flex items-center gap-2">`
line — the most common Tailwind flexbox one-liner. The content inside is
genuinely different (heading vs code element).

### 4. `display/table.templ:199-215` vs `navigation/breadcrumbs.templ:116-133` — for-loop (t=7)

Both iterate over a slice and render list items. Table renders `<tr>` → `<td>`
cells with a content/text fallback. Breadcrumbs renders `<li>` items with a
separator + link/span.

The shared pattern is `for i, item := range props.Items { <element>...</element> }`
— Go's only iteration syntax. The HTML elements, content logic, and CSS
classes are completely different.

**Why not lazy**: These are different HTML list structures (`<tbody>` vs `<ol>`)
serving different domains (data table vs navigation). They share Go's for-loop
syntax, not business logic.

### 5. `errorpage/notfound404.templ:62-70` vs `feedback/alert.templ:81-85` — icon + text (t=7)

NotFound404 renders a `<button data-tc-go-back>` with an icon + text span.
Alert's `inlineMessage` renders a `<div role="..." aria-live="polite">` with an
icon + text span.

The shared tokens are "element with attributes + child icon + child text span"
— the fundamental icon-with-label pattern used everywhere in UI design.

**Why not lazy**: The elements (`<button>` vs `<div>`), accessibility patterns
(`data-tc-go-back` vs `role`/`aria-live`), and purposes (navigation action vs
status announcement) are completely different.

### 6. `display/shared.templ:45-53` / `feedback/loading.templ:147-150` / `icons/icon.templ:23-26` — element with child component (t=6)

At t=6, art-dupl matches "HTML element with attributes + single child
component call". The three occurrences:

- shared.templ: `<button>` close button calling `@icons.Icon(icons.X, ...)`
- loading.templ: `<div>` skeleton wrapper calling `@skeletonBody(variant)`
- icon.templ: `<svg>` spinner wrapper calling `@svg.SpinnerSVG()`

These are different elements (button, div, svg) with different purposes (close
control, loading placeholder, animated spinner).

**Why not lazy**: This is the fundamental template composition primitive — an
element wrapping a child component. It cannot be extracted further without
inventing a generic "element with child" wrapper that takes the tag name as
a string, which would bypass templ's type safety.

## Accepted Demo Clones

### 7. `examples/demo/demo.templ:322-343` — three demo table rows (t=8)

Three `<tr>` rows in `demoTableBody()` showing Success/Failed/Passed badge
types. The demo binary is excluded from lint and is not production code.

### 8. `examples/demo/demo.templ:129-140` vs `:141-171` — section wrappers (t=8)

Both blocks follow the `@demoSection("...") + <div class="...">` convention.
Content is completely different (9 button variants vs 8 form controls).

### 9. `examples/demo/demo.templ:55-62` vs `:172-179` — badge/avatar demo blocks (t=7)

Two demo sections using `<div class="flex flex-wrap gap-3">` to showcase
badges and error alerts respectively.

## Decision

These clones remain because each extraction attempt either:

1. Was already performed (see extraction list above)
2. Would add more indirection than lines saved
3. Is a structural constraint of the templ DSL
4. Is in the demo binary (not production code)

## Consequences

- Running `art-dupl -t 8` shows 4 groups (2 production + 2 demo)
- Running `art-dupl -t 7` shows 8 groups (5 production + 3 demo)
- Running `art-dupl -t 6` shows 9 groups (6 production + 3 demo)
- Each production clone is a different component type sharing a Tailwind idiom
- New components should use existing extractions (errorHeader, overlayShell,
  skeletonContainer, DismissButton, definitionDetailContent) where applicable
