# Features — templ-components

**Updated:** 2026-07-05 | **Version:** 0.6.1

A Go component library built on [templ](https://templ.guide) and [Tailwind CSS v4](https://tailwindcss.com) for building server-rendered web applications.

---

## Overview

| Package      | Components    | Description                                                                                                                                                                                               |
| ------------ | ------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `utils`      | 0             | Shared types, Tailwind class merging, generic helpers                                                                                                                                                     |
| `display`    | 25            | UI display: cards, badges, buttons, modals, drawers, tables, tabs, avatars, tooltips, accordions, dropdowns, empty states, page headers, definition lists, copy button, relative time, count badge, image |
| `errorpage`  | 3             | Error presentation: full-page errors, error detail cards, family-aware alerts                                                                                                                             |
| `feedback`   | 13            | User feedback: alerts, toasts, spinners, progress bars, skeletons                                                                                                                                         |
| `forms`      | 16            | Form controls: inputs, selects, textareas, checkboxes, radios, toggles, file inputs, date pickers, comboboxes, validation                                                                                 |
| `htmx`       | 7             | HTMX integration: loading indicators, error handling, helpers                                                                                                                                             |
| `icons`      | 2 (101 icons) | SVG icon system with typed name constants                                                                                                                                                                 |
| `layout`     | 5             | Page layout: base HTML, theme toggle, dark mode, CSP-safe script tag                                                                                                                                      |
| `navigation` | 11            | Navigation: nav bars, breadcrumbs, pagination, mobile menus, sidebar nav, load more                                                                                                                       |

**Totals:** 80 templ components, 101 icon names, 26+ typed enums, 59 generated `*_templ.go` files, ~24,000 lines of Go/templ source

---

## Package: `utils`

### `BaseProps` (shared by all components)

```go
type BaseProps struct {
    ID        string
    Class     string
    Attrs     templ.Attributes
    AriaLabel string
    Nonce     string
}
```

### Functions

| Function        | Signature                                               | Purpose                                           |
| --------------- | ------------------------------------------------------- | ------------------------------------------------- |
| `Class`         | `(classes ...string) string`                            | Merges Tailwind classes via tailwind-merge-go     |
| `CurrentYear`   | `() string`                                             | Current year string                               |
| `Ternary`       | `[T any](bool, a, b T) T`                               | Generic ternary                                   |
| `Lookup`        | `[K comparable, V any](m map[K]V, key K, fallback V) V` | Generic map lookup with fallback                  |
| `EnsureID`      | `(prefix, id string) string`                            | Auto-generates unique ID via crypto/rand if empty |
| `DismissScript` | `() string`                                             | Shared JS for [data-dismiss] click delegation     |

### Test Helpers (exported)

| Function            | Purpose                                            |
| ------------------- | -------------------------------------------------- |
| `Render`            | Renders a templ.Component to string                |
| `RenderAll`         | Renders multiple components to concatenated string |
| `AssertContains`    | Asserts substring in rendered output               |
| `AssertNotContains` | Asserts substring absent from output               |
| `AssertContainsAll` | Asserts output contains every substring            |
| `AssertEqual`       | Asserts two values are equal                       |

---

## Package: `display`

### Components

| Component          | Status           | Description                     | Key Features                                                                                                       |
| ------------------ | ---------------- | ------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `Accordion`        | FULLY_FUNCTIONAL | Collapsible accordion panels    | JS toggle, `aria-expanded`, `aria-controls`, chevron rotation                                                      |
| `Avatar`           | FULLY_FUNCTIONAL | User avatar with image/initials | AvatarStatus enum, 5 sizes, circle/square, online/offline dot                                                      |
| `Button`           | FULLY_FUNCTIONAL | Action button                   | 5 variants, 3 sizes, href (link mode), loading state                                                               |
| `Badge`            | FULLY_FUNCTIONAL | Status label                    | 7 color types, 3 sizes, pill shape, dot indicator                                                                  |
| `StatusBadge`      | FULLY_FUNCTIONAL | Auto-mapped status badge        | Maps ~20 status strings to badge types                                                                             |
| `Card`             | FULLY_FUNCTIONAL | Bordered card container         | Header, subtitle, footer, header action, 4 padding sizes                                                           |
| `SimpleCard`       | FULLY_FUNCTIONAL | Minimal card                    | Children-only, no header/footer                                                                                    |
| `StatCard`         | FULLY_FUNCTIONAL | Dashboard metric card           | StatCardProps struct, TrendDirection enum (Up/Down/None), optional `Icon` field, `Href` renders as clickable `<a>` |
| `Grid`             | FULLY_FUNCTIONAL | Responsive grid container       | Typed `GridCols` enum (1–6), responsive breakpoints (1→2→3→N), children slot                                       |
| `Dropdown`         | FULLY_FUNCTIONAL | Button-triggered menu           | External/internal links, buttons, keyboard nav, ARIA menu                                                          |
| `Drawer`           | FULLY_FUNCTIONAL | Side panel                      | Left/right slide, focus trap, Escape, backdrop, 5 sizes                                                            |
| `EmptyState`       | FULLY_FUNCTIONAL | Centered empty state            | Icon, title, description, action link/button                                                                       |
| `SimpleEmptyState` | FULLY_FUNCTIONAL | Minimal empty state             | Text-only                                                                                                          |
| `Modal`            | FULLY_FUNCTIONAL | Accessible modal dialog         | Focus trap, Escape close, backdrop, 5 sizes, open/close API                                                        |
| `Table`            | FULLY_FUNCTIONAL | Responsive data table           | Headers, rows, striping, hover, caption, bordered, cell content                                                    |
| `Tabs`             | FULLY_FUNCTIONAL | Tabbed interface                | Default underline or pills variant, ARIA tablist/tab/tabpanel                                                      |
| `Tooltip`          | FULLY_FUNCTIONAL | Hover tooltip                   | 4 positions, arrow, `role="tooltip"`, CSS-only                                                                     |
| `PageHeader`       | FULLY_FUNCTIONAL | Page title block                | Title, subtitle, breadcrumb + action component slots                                                               |
| `DefinitionList`   | FULLY_FUNCTIONAL | Two-column key/value list       | Typed `DefinitionItem` entries, semantic `<dl>` markup                                                             |
| `ListNote`         | FULLY_FUNCTIONAL | Truncation notice               | "Showing N of M" when a list is truncated                                                                          |
| `CopyButton`       | FULLY_FUNCTIONAL | Clipboard copy button           | CSP-safe singleton JS, "Copied!" feedback, clipboard icon                                                          |
| `RelativeTime`     | FULLY_FUNCTIONAL | Relative timestamp              | `<time datetime>` with "2 hours ago", absolute title on hover                                                      |
| `CountBadge`       | FULLY_FUNCTIONAL | Notification count overlay      | Absolute-positioned badge, overflow "N+", aria-hidden decorative                                                   |
| `DefinitionGrid`   | FULLY_FUNCTIONAL | Responsive key/value grid       | Term-detail pairs in SimpleCard tiles, composes through Grid                                                       |
| `Image`            | FULLY_FUNCTIONAL | Lazy-loaded image               | `loading=lazy` default, width/height for CLS, CSP-safe fallback via event capture                                  |

### Enums

| Type               | Values                                                   |
| ------------------ | -------------------------------------------------------- |
| `AvatarSize`       | XS, SM, MD, LG, XL                                       |
| `AvatarShape`      | Circle, Square                                           |
| `AvatarStatus`     | Online, Offline, None                                    |
| `BadgeType`        | Default, Primary, Success, Warning, Error, Info, Neutral |
| `BadgeSize`        | SM, MD, LG                                               |
| `CardPadding`      | None, SM, MD, LG                                         |
| `DropdownPosition` | Left, Right                                              |
| `ModalSize`        | SM, MD, LG, XL, 2XL, Full                                |
| `DrawerSide`       | Left, Right                                              |
| `DrawerSize`       | SM, MD, LG, XL, 2XL, Full                                |
| `TabsVariant`      | Default, Pills                                           |
| `TrendDirection`   | Up, Down, None                                           |
| `GridCols`         | 1, 2, 3 (default), 4, 5, 6                               |
| `TooltipPosition`  | Top, Bottom, Left, Right                                 |

### Known Issues

- **Tooltip** calls `tooltipLookupPosition()` twice per render instead of caching

---

## Package: `errorpage`

### Components

| Component     | Status           | Description              | Key Features                                                            |
| ------------- | ---------------- | ------------------------ | ----------------------------------------------------------------------- |
| `ErrorPage`   | FULLY_FUNCTIONAL | Full-page error view     | Wix-style What/Why/Fix/WayOut, 5 families, context, cause chain, action |
| `NotFound404` | FULLY_FUNCTIONAL | Dedicated 404 page       | Gradient numeral hero, search form, quick-links grid, go-home/go-back   |
| `ErrorDetail` | FULLY_FUNCTIONAL | Inline error detail card | Code badge, family badge, context table, cause chain, suggested fix     |
| `ErrorAlert`  | FULLY_FUNCTIONAL | Family-aware alert       | 5 distinct color schemes, dismiss, fix suggestion, family badge         |

### Enums

| Type     | Values                                                     |
| -------- | ---------------------------------------------------------- |
| `Family` | Rejection, Conflict, Transient, Corruption, Infrastructure |

### Bridge Helpers

| Function            | Purpose                                                     |
| ------------------- | ----------------------------------------------------------- |
| `FamilyStatusCode`  | Maps Family → HTTP status code (400/409/503/500/503)        |
| `ContextMap`        | Converts map[string]string → []ContextPair                  |
| `ExtractCauseChain` | Walks Unwrap() chain → []CauseItem with ErrorCode() support |
| `FromError`         | Converts any `error` → `ErrorPageProps` (family/code/cause) |

### HTTP Handler

|                      | Function                       | Signature                                  | Purpose |
| -------------------- | ------------------------------ | ------------------------------------------ | ------- |
| `ErrorHandler`       | `(err, cfg) http.Handler`      | Serves error page with correct HTTP status |
| `WriteError`         | `(w, r, err, nonce)`           | Convenience wrapper for ErrorHandler       |
| `WriteErrorPage`     | `(w, r, status, props, nonce)` | Renders pre-configured ErrorPageProps      |
| `NotFound`           | `() ErrorPageProps`            | Pre-built 404 (Rejection)                  |
| `Forbidden`          | `() ErrorPageProps`            | Pre-built 403 (Rejection)                  |
| `BadRequest`         | `(msg) ErrorPageProps`         | Pre-built 400 (Rejection)                  |
| `Conflict`           | `(msg) ErrorPageProps`         | Pre-built 409 (Conflict)                   |
| `ServiceUnavailable` | `() ErrorPageProps`            | Pre-built 503 (Transient)                  |
| `InternalError`      | `() ErrorPageProps`            | Pre-built 500 (Infrastructure)             |

### Family Visual Mapping

| Family         | Color  | Icon                | Tone          | HTTP Status |
| -------------- | ------ | ------------------- | ------------- | ----------- |
| Rejection      | Amber  | ExclamationTriangle | Instructional | 400         |
| Conflict       | Orange | ExclamationCircle   | Explanatory   | 409         |
| Transient      | Blue   | Refresh             | Reassuring    | 503         |
| Corruption     | Red    | ExclamationTriangle | Urgent        | 500         |
| Infrastructure | Gray   | Globe               | Apologetic    | 503         |

---

## Package: `feedback`

### Components

| Component          | Status           | Description               | Key Features                                                      |
| ------------------ | ---------------- | ------------------------- | ----------------------------------------------------------------- |
| `Alert`            | FULLY_FUNCTIONAL | Full-width alert banner   | 4 types, dismissible, icon, CSP nonce                             |
| `InlineError`      | FULLY_FUNCTIONAL | Compact error message     | Red icon + text                                                   |
| `InlineSuccess`    | FULLY_FUNCTIONAL | Compact success message   | Green checkmark + text                                            |
| `Spinner`          | FULLY_FUNCTIONAL | Animated SVG spinner      | 3 sizes, custom color class                                       |
| `LoadingOverlay`   | FULLY_FUNCTIONAL | Full-screen loading       | Spinner, message, optional progress bar                           |
| `InlineLoading`    | FULLY_FUNCTIONAL | Compact inline loading    | Spinner + message                                                 |
| `Skeleton`         | FULLY_FUNCTIONAL | Pulsing placeholder       | 7 variants: text, title, avatar, image, card, table-row           |
| `SkeletonGroup`    | FULLY_FUNCTIONAL | Multiple skeletons        | Group with `animate-pulse`                                        |
| `SkeletonCardGrid` | FULLY_FUNCTIONAL | Loading card grid         | N skeleton cards in responsive 3-col grid, single `role="status"` |
| `ProgressBar`      | FULLY_FUNCTIONAL | Progress indicator        | 3 sizes, custom color, label, percentage                          |
| `StepIndicator`    | FULLY_FUNCTIONAL | Horizontal step indicator | Completed/current/future states, labels                           |
| `ToastContainer`   | FULLY_FUNCTIONAL | Fixed toast container     | JS `tcShowToast()` for dynamic toasts                             |
| `Toast`            | FULLY_FUNCTIONAL | Server-rendered toast     | 4 types, dismissible, title, message, duration                    |

### Enums

| Type              | Values                                                |
| ----------------- | ----------------------------------------------------- |
| `AlertType`       | Success, Error, Warning, Info                         |
| `SpinnerSize`     | SM, MD, LG                                            |
| `SkeletonVariant` | Text, TextShort, Title, Avatar, Image, Card, TableRow |
| `ProgressBarSize` | SM, MD, LG                                            |
| `ToastType`       | Success, Error, Warning, Info                         |

### Constants

| Name                  | Value  | Purpose             |
| --------------------- | ------ | ------------------- |
| `ToastDurationShort`  | 3000ms | Short auto-dismiss  |
| `ToastDurationMedium` | 5000ms | Medium auto-dismiss |
| `ToastDurationLong`   | 8000ms | Long auto-dismiss   |

### Shared Style System

```go
type feedbackStyleSet struct { Border, BG, Text, Icon string }
func lookupFeedbackStyle[T ~string](m map[T]feedbackStyleSet, def feedbackStyleSet, t T) feedbackStyleSet
```

Used by both Alert and Toast for consistent visual styling.

### Known Issues

- _None known_

---

## Package: `forms`

### Components

| Component           | Status           | Description              | Key Features                                                        |
| ------------------- | ---------------- | ------------------------ | ------------------------------------------------------------------- |
| `Input`             | FULLY_FUNCTIONAL | Text input with label    | 11 types, error, help text, required, disabled, readonly, autofocus |
| `Checkbox`          | FULLY_FUNCTIONAL | Checkbox with label      | Error, help text, required                                          |
| `RadioGroup`        | FULLY_FUNCTIONAL | Radio button group       | Inline/stacked, AriaLabel propagation on fieldset                   |
| `Select`            | FULLY_FUNCTIONAL | Select dropdown          | Options, disabled options, pre-selected, normalize contradiction    |
| `Textarea`          | FULLY_FUNCTIONAL | Textarea with label      | Configurable rows, error, help text                                 |
| `Toggle`            | FULLY_FUNCTIONAL | Toggle switch            | 3 sizes (SM/MD/LG), label, error, help text                         |
| `FileInput`         | FULLY_FUNCTIONAL | File upload input        | Multiple, accept filter, error, help text                           |
| `DatePicker`        | FULLY_FUNCTIONAL | Date input               | Native `<input type="date">`, min/max constraints                   |
| `Combobox`          | FULLY_FUNCTIONAL | Autocomplete with filter | `role="combobox"`, client-side filtering, auto-generated IDs        |
| `Label`             | FULLY_FUNCTIONAL | Form label               | Optional `for` attribute, required indicator                        |
| `FieldError`        | FULLY_FUNCTIONAL | Field validation error   | Accessible with ID linking for aria-describedby                     |
| `ValidationSummary` | FULLY_FUNCTIONAL | Accessible error summary | Icon, error count, linked fields, `role="alert"`                    |
| `Form`              | FULLY_FUNCTIONAL | Form wrapper             | Action, Method (GET/POST), CSRF token, children pattern             |
| `InputGroup`        | FULLY_FUNCTIONAL | Input group container    | Groups multiple inputs with shared styling                          |
| `FormFieldWrapper`  | FULLY_FUNCTIONAL | Shared field chrome      | Label + FieldError + helpText, used by Input/Select/Textarea        |
| `Radio`             | FULLY_FUNCTIONAL | Single radio button      | Sub-component of RadioGroup                                         |

### Enums

| Type        | Values                                                                        |
| ----------- | ----------------------------------------------------------------------------- |
| `InputType` | Text, Email, Password, Number, Tel, URL, Date, Time, Datetime, Search, Hidden |

### Functions

| Function     | Purpose                                 |
| ------------ | --------------------------------------- |
| `SanitizeID` | Converts arbitrary text to safe HTML ID |
| `ErrorAttrs` | Returns aria-invalid + aria-describedby |

### Known Issues

- **Combobox** lacks full WAI-ARIA keyboard pattern (no ArrowDown/Up navigation within listbox, no `aria-activedescendant` tracking)

---

## Package: `htmx`

### Components

| Component              | Status           | Description                | Key Features                                                                               |
| ---------------------- | ---------------- | -------------------------- | ------------------------------------------------------------------------------------------ |
| `LoadingIndicator`     | FULLY_FUNCTIONAL | Fixed full-screen loading  | Uses `htmx-indicator`, blur backdrop                                                       |
| `InlineLoadingOverlay` | FULLY_FUNCTIONAL | Localized loading overlay  | Absolute positioned, for form targets                                                      |
| `LoadingButton`        | FULLY_FUNCTIONAL | Button with loading state  | Text swaps to spinner during HTMX requests                                                 |
| `ConfirmDelete`        | FULLY_FUNCTIONAL | Delete button with confirm | `hx-delete`, `hx-target`, `hx-confirm`, `hx-swap`                                          |
| `SwapOOB`              | FULLY_FUNCTIONAL | Out-of-band swap wrapper   | For updating multiple elements per response                                                |
| `CSRFToken`            | FULLY_FUNCTIONAL | Hidden CSRF input          | Standard `csrf_token` name                                                                 |
| `GlobalErrorHandling`  | FULLY_FUNCTIONAL | HTMX error handler         | Network errors, response errors, auto-retry, toast integration, family-aware error parsing |

### Known Issues

- **Hidden coupling**: `GlobalErrorHandling` calls `tcShowToast()` — requires `ToastContainer` on page, silently fails otherwise

---

## Package: `icons`

### Components

| Component             | Status           | Description      | Key Features                                        |
| --------------------- | ---------------- | ---------------- | --------------------------------------------------- |
| `Icon`                | FULLY_FUNCTIONAL | SVG icon by name | 101 named icons, custom class, currentColor theming |
| `IconWithStrokeWidth` | FULLY_FUNCTIONAL | Icon variant     | Custom stroke-width (default Icon uses 1.5)         |

### Icon Names (101)

100 path icons + Spinner covering navigation, UI actions, chevrons/arrows, communication, media, and status. See `icons/icon_names.go` for the complete list.

### Functions

| Function              | Purpose                                                                |
| --------------------- | ---------------------------------------------------------------------- |
| `IconWithStrokeWidth` | Icon with custom stroke-width                                          |
| `IconPathData`        | Returns raw path data for a named icon (full `<svg>` wrapper control)  |
| `allIconNames`        | Auto-generated list of all icon names from `iconPathData` (unexported) |

### Known Issues

- Unknown icon names fall back to the Question icon instead of crashing render

---

## Package: `layout`

### Components

| Component     | Status           | Description           | Key Features                                                                                           |
| ------------- | ---------------- | --------------------- | ------------------------------------------------------------------------------------------------------ |
| `Base`        | FULLY_FUNCTIONAL | Complete HTML5 layout | Meta tags, OG tags, Twitter cards, CSS, HTMX with SRI, theme script, skip-to-content, security headers |
| `Minimal`     | FULLY_FUNCTIONAL | Minimal HTML document | No dependencies, for static pages/emails/PDFs                                                          |
| `ThemeScript` | FULLY_FUNCTIONAL | Dark mode script      | localStorage-based, prevents FOUC                                                                      |
| `ThemeToggle` | FULLY_FUNCTIONAL | Theme switch button   | Sun/moon icons, JS toggle, CSP nonce                                                                   |
| `Script`      | FULLY_FUNCTIONAL | CSP-safe script tag   | Auto-injects nonce, optional attrs (async, defer, type)                                                |

### Functions

| Function           | Purpose                                                         |
| ------------------ | --------------------------------------------------------------- |
| `DefaultPageProps` | Returns sensible defaults (locale=en, HTMX 2.0.10, SRI enabled) |

### Known Issues

- _None known_

---

## Package: `navigation`

### Components

| Component          | Status           | Description               | Key Features                                                       |
| ------------------ | ---------------- | ------------------------- | ------------------------------------------------------------------ |
| `Nav`              | FULLY_FUNCTIONAL | Responsive navigation bar | Brand, links, right items, sticky option, mobile menu              |
| `SimpleNav`        | FULLY_FUNCTIONAL | Simplified nav            | Text brand, sticky by default                                      |
| `NavLink`          | FULLY_FUNCTIONAL | Desktop nav link          | Active state styling via currentPath                               |
| `MobileNavLink`    | FULLY_FUNCTIONAL | Mobile nav link           | Border-left active indicator                                       |
| `Breadcrumbs`      | FULLY_FUNCTIONAL | Breadcrumb navigation     | Chevron separators, active state                                   |
| `Pagination`       | FULLY_FUNCTIONAL | Page navigation           | Mobile/desktop layouts, prev/next arrows, page range, query param  |
| `MobileMenu`       | FULLY_FUNCTIONAL | Collapsible mobile menu   | JS toggle, nonce-based CSP                                         |
| `MobileMenuToggle` | FULLY_FUNCTIONAL | Hamburger button          | Conditional visibility                                             |
| `Footer`           | FULLY_FUNCTIONAL | Simple footer             | Copyright with dynamic year                                        |
| `SidebarNav`       | FULLY_FUNCTIONAL | Vertical sidebar nav      | Brand/footer slots, icons, active-route detection                  |
| `LoadMore`         | FULLY_FUNCTIONAL | Cursor pagination button  | hx-get + hx-swap outerHTML, cursor as query param, centered layout |

### Known Issues

- `Footer` could arguably live in `layout/` instead

---

## Cross-Cutting Features

- **CSP Compliance:** All inline scripts use `nonce` attribute
- **Dark Mode:** Full Tailwind `dark:` variant support via `layout.ThemeScript`
- **Tailwind Class Merging:** `utils.Class()` uses tailwind-merge-go for conflict resolution
- **Accessibility:** `aria-*` attributes, `role` attributes, screen-reader text, keyboard navigation (modal focus trap, dropdown arrows, tabs)
- **Responsive:** Mobile-first designs with `sm:` breakpoints
- **Type Safety:** 26 typed string enums, `utils.BaseProps` embedded in all Props structs
- **Test Coverage:** 74% average across packages, BDD + snapshot + a11y + benchmark + integration tests
- **Theming:** Tailwind v4 `@theme` override support via `templ-components-theme.css`. Components emit standard utility classes (`bg-blue-600`, `text-gray-900`) — consumers override `--color-*` variables to theme globally without touching component code.

---

## Planned / Not Yet Implemented

| Component       | Package | Notes                      |
| --------------- | ------- | -------------------------- |
| Docs site       | —       | Auto-generated from source |
| Release tooling | —       | goreleaser, tag-based      |
