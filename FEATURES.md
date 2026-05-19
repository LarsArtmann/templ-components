# Features — templ-components

**Updated:** 2026-05-19 | **Version:** 0.x (pre-release)

A Go component library built on [templ](https://templ.guide) and [Tailwind CSS v4](https://tailwindcss.com) for building server-rendered web applications.

---

## Overview

| Package      | Components   | Description                                                                               |
| ------------ | ------------ | ----------------------------------------------------------------------------------------- |
| `utils`      | 0            | Shared types, Tailwind class merging, generic helpers                                     |
| `display`    | 14           | UI display: cards, badges, modals, tables, tabs, avatars, tooltips, accordions, dropdowns |
| `feedback`   | 12           | User feedback: alerts, toasts, spinners, progress bars, skeletons                         |
| `forms`      | 6            | Form controls: inputs, selects, textareas, checkboxes, labels, errors                     |
| `htmx`       | 7            | HTMX integration: loading indicators, error handling, helpers                             |
| `icons`      | 1 (42 icons) | SVG icon system with typed name constants                                                 |
| `layout`     | 4            | Page layout: base HTML, theme toggle, dark mode                                           |
| `navigation` | 9            | Navigation: nav bars, breadcrumbs, pagination, mobile menus                               |

**Totals:** 53 templ components, 44 icon names, 17 typed enums, 30 `.templ` files, ~3,400 lines of Go/templ source

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

| Function      | Signature                                  | Purpose                                       |
| ------------- | ------------------------------------------ | --------------------------------------------- |
| `Class`       | `(classes ...string) string`               | Merges Tailwind classes via tailwind-merge-go |
| `MergeAttrs`  | `(m ...templ.Attributes) templ.Attributes` | Merges attribute maps                         |
| `CurrentYear` | `() string`                                | Current year string                           |
| `Ternary`     | `[T any](bool, a, b T) T`                  | Generic ternary                               |
| `Deref`       | `[T any](p *T) T`                          | Nil-safe deref                                |
| `DerefOr`     | `[T any](p *T, fallback T) T`              | Nil-safe deref with fallback                  |
| `BoolString`  | `(b bool) string`                          | Returns `"true"` or `"false"`                 |
| `MapEnum`     | `[T ~string](m map[T]U, key T) U`          | Generic map lookup with fallback              |

### Test Helpers (exported)

| Function            | Purpose                              |
| ------------------- | ------------------------------------ |
| `Render`            | Renders a templ.Component to string  |
| `AssertContains`    | Asserts substring in rendered output |
| `AssertNotContains` | Asserts substring absent from output |
| `AssertEqual`       | Asserts two values are equal         |

---

## Package: `display`

### Components

| Component          | Status           | Description                     | Key Features                                                    |
| ------------------ | ---------------- | ------------------------------- | --------------------------------------------------------------- |
| `Accordion`        | FULLY_FUNCTIONAL | Collapsible accordion panels    | JS toggle, `aria-expanded`, `aria-controls`, chevron rotation   |
| `Avatar`           | FULLY_FUNCTIONAL | User avatar with image/initials | AvatarStatus enum, 5 sizes, circle/square, online/offline dot   |
| `Badge`            | FULLY_FUNCTIONAL | Status label                    | 7 color types, 3 sizes, pill shape, dot indicator               |
| `StatusBadge`      | FULLY_FUNCTIONAL | Auto-mapped status badge        | Maps ~20 status strings to badge types                          |
| `Card`             | FULLY_FUNCTIONAL | Bordered card container         | Header, subtitle, footer, header action, 4 padding sizes        |
| `SimpleCard`       | FULLY_FUNCTIONAL | Minimal card                    | Children-only, no header/footer                                 |
| `StatCard`         | FULLY_FUNCTIONAL | Dashboard metric card           | StatCardProps struct, TrendDirection enum (Up/Down/None)        |
| `Dropdown`         | FULLY_FUNCTIONAL | Button-triggered menu           | External/internal links, buttons, keyboard nav, ARIA menu       |
| `EmptyState`       | FULLY_FUNCTIONAL | Centered empty state            | Icon, title, description, action link/button                    |
| `SimpleEmptyState` | FULLY_FUNCTIONAL | Minimal empty state             | Text-only                                                       |
| `Modal`            | FULLY_FUNCTIONAL | Accessible modal dialog         | Focus trap, Escape close, backdrop, 5 sizes, open/close API     |
| `Table`            | FULLY_FUNCTIONAL | Responsive data table           | Headers, rows, striping, hover, caption, bordered, cell content |
| `Tabs`             | FULLY_FUNCTIONAL | Tabbed interface                | Default underline or pills variant, ARIA tablist/tab/tabpanel   |
| `Tooltip`          | FULLY_FUNCTIONAL | Hover tooltip                   | 4 positions, arrow, `role="tooltip"`, CSS-only                  |

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
| `ModalSize`        | SM, MD, LG, XL, Full                                     |
| `TabsVariant`      | Default, Pills                                           |
| `TrendDirection`   | Up, Down, None                                           |
| `TooltipPosition`  | Top, Bottom, Left, Right                                 |

### Known Issues

- **Tooltip** calls `tooltipLookupPosition()` twice per render instead of caching

---

## Package: `feedback`

### Components

| Component        | Status           | Description               | Key Features                                            |
| ---------------- | ---------------- | ------------------------- | ------------------------------------------------------- |
| `Alert`          | FULLY_FUNCTIONAL | Full-width alert banner   | 4 types, dismissible, icon, CSP nonce                   |
| `InlineError`    | FULLY_FUNCTIONAL | Compact error message     | Red icon + text                                         |
| `InlineSuccess`  | FULLY_FUNCTIONAL | Compact success message   | Green checkmark + text                                  |
| `Spinner`        | FULLY_FUNCTIONAL | Animated SVG spinner      | 3 sizes, custom color class                             |
| `LoadingOverlay` | FULLY_FUNCTIONAL | Full-screen loading       | Spinner, message, optional progress bar                 |
| `InlineLoading`  | FULLY_FUNCTIONAL | Compact inline loading    | Spinner + message                                       |
| `Skeleton`       | FULLY_FUNCTIONAL | Pulsing placeholder       | 7 variants: text, title, avatar, image, card, table-row |
| `SkeletonGroup`  | FULLY_FUNCTIONAL | Multiple skeletons        | Group with `animate-pulse`                              |
| `ProgressBar`    | FULLY_FUNCTIONAL | Progress indicator        | 3 sizes, custom color, label, percentage                |
| `StepIndicator`  | FULLY_FUNCTIONAL | Horizontal step indicator | Completed/current/future states, labels                 |
| `ToastContainer` | FULLY_FUNCTIONAL | Fixed toast container     | JS `tcShowToast()` for dynamic toasts                   |
| `Toast`          | FULLY_FUNCTIONAL | Server-rendered toast     | 4 types, dismissible, title, message, duration          |

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

- **Toast icon SVG paths duplicated** across Go and JS — single source of truth missing
- **Dismiss JS** pattern duplicated between Alert and Toast
- **Spinner SVG** rendered 3 different ways across packages

---

## Package: `forms`

### Components

| Component    | Status           | Description            | Key Features                                                        |
| ------------ | ---------------- | ---------------------- | ------------------------------------------------------------------- |
| `Input`      | FULLY_FUNCTIONAL | Text input with label  | 11 types, error, help text, required, disabled, readonly, autofocus |
| `Checkbox`   | FULLY_FUNCTIONAL | Checkbox with label    | Error, help text, required                                          |
| `Select`     | FULLY_FUNCTIONAL | Select dropdown        | Options, disabled options, pre-selected                             |
| `Textarea`   | FULLY_FUNCTIONAL | Textarea with label    | Configurable rows, error, help text                                 |
| `Label`      | FULLY_FUNCTIONAL | Form label             | Optional `for` attribute, required indicator                        |
| `FieldError` | FULLY_FUNCTIONAL | Field validation error | Accessible with ID linking for aria-describedby                     |

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

- `ErrorAttrs` doesn't link help text ID when both error and help text are present
- No Radio, File input, or Toggle/Switch components yet

---

## Package: `htmx`

### Components

| Component              | Status           | Description                | Key Features                                                   |
| ---------------------- | ---------------- | -------------------------- | -------------------------------------------------------------- |
| `LoadingIndicator`     | FULLY_FUNCTIONAL | Fixed full-screen loading  | Uses `htmx-indicator`, blur backdrop                           |
| `InlineLoadingOverlay` | FULLY_FUNCTIONAL | Localized loading overlay  | Absolute positioned, for form targets                          |
| `LoadingButton`        | FULLY_FUNCTIONAL | Button with loading state  | Text swaps to spinner during HTMX requests                     |
| `ConfirmDelete`        | FULLY_FUNCTIONAL | Delete button with confirm | `hx-delete`, `hx-target`, `hx-confirm`, `hx-swap`              |
| `SwapOOB`              | FULLY_FUNCTIONAL | Out-of-band swap wrapper   | For updating multiple elements per response                    |
| `CSRFToken`            | FULLY_FUNCTIONAL | Hidden CSRF input          | Standard `csrf_token` name                                     |
| `GlobalErrorHandling`  | FULLY_FUNCTIONAL | HTMX error handler         | Network errors, response errors, auto-retry, toast integration |

### Known Issues

- **Hidden coupling**: `GlobalErrorHandling` calls `tcShowToast()` — requires `ToastContainer` on page, silently fails otherwise
- **Magic numbers**: Retry count (2), delay (1000ms), error history (10) are hardcoded
- **Package coupling**: `htmx/loading.templ` accepts `templ.Component` for spinner (decoupled)

---

## Package: `icons`

### Components

| Component | Status           | Description      | Key Features                                       |
| --------- | ---------------- | ---------------- | -------------------------------------------------- |
| `Icon`    | FULLY_FUNCTIONAL | SVG icon by name | 44 named icons, custom class, currentColor theming |

### Icon Names (44)

Home, Users, Folder, Document, Search, Settings, Chart, Inbox, Check, X, Plus, Minus, ChevronRight, ChevronLeft, ChevronDown, ChevronUp, ArrowRight, ArrowLeft, Refresh, ExternalLink, Download, Upload, Trash, Edit, Eye, EyeOff, Lock, Unlock, Menu, Bell, Calendar, Clock, Location, Phone, Mail, Globe, Sun, Moon, Spinner, Exclamation _(deprecated)_, ExclamationTriangle, ExclamationCircle, CheckCircle, Information, Question

### Functions

| Function    | Status          | Purpose                                                   |
| ----------- | --------------- | --------------------------------------------------------- |
| `IconAttrs` | EXPORTED_UNUSED | Returns accessible attributes (aria-label or aria-hidden) |

### Known Issues

- **`IconAttrs` is dead code** — exported but never called anywhere
- Unknown icon names silently fall back to a clock icon — no warning

---

## Package: `layout`

### Components

| Component     | Status           | Description           | Key Features                                                                                           |
| ------------- | ---------------- | --------------------- | ------------------------------------------------------------------------------------------------------ |
| `Base`        | FULLY_FUNCTIONAL | Complete HTML5 layout | Meta tags, OG tags, Twitter cards, CSS, HTMX with SRI, theme script, skip-to-content, security headers |
| `Minimal`     | FULLY_FUNCTIONAL | Minimal HTML document | No dependencies, for static pages/emails/PDFs                                                          |
| `ThemeScript` | FULLY_FUNCTIONAL | Dark mode script      | localStorage-based, prevents FOUC                                                                      |
| `ThemeToggle` | FULLY_FUNCTIONAL | Theme switch button   | Sun/moon icons, JS toggle, CSP nonce                                                                   |

### Functions

| Function           | Purpose                                                        |
| ------------------ | -------------------------------------------------------------- |
| `DefaultPageProps` | Returns sensible defaults (locale=en, HTMX 2.0.6, SRI enabled) |

### Known Issues

- **HTMX CDN URL** repeated 4 times — should be a constant
- **`Minimal` uses positional params** while `Base` uses props struct — inconsistent API
- **Magic theme colors** `#4f46e5` and `#1e1b4b` hardcoded in PageProps defaults

---

## Package: `navigation`

### Components

| Component          | Status           | Description               | Key Features                                                      |
| ------------------ | ---------------- | ------------------------- | ----------------------------------------------------------------- |
| `Nav`              | FULLY_FUNCTIONAL | Responsive navigation bar | Brand, links, right items, sticky option, mobile menu             |
| `SimpleNav`        | FULLY_FUNCTIONAL | Simplified nav            | Text brand, sticky by default                                     |
| `NavLink`          | FULLY_FUNCTIONAL | Desktop nav link          | Active state styling via currentPath                              |
| `MobileNavLink`    | FULLY_FUNCTIONAL | Mobile nav link           | Border-left active indicator                                      |
| `Breadcrumbs`      | FULLY_FUNCTIONAL | Breadcrumb navigation     | Chevron separators, active state                                  |
| `Pagination`       | FULLY_FUNCTIONAL | Page navigation           | Mobile/desktop layouts, prev/next arrows, page range, query param |
| `MobileMenu`       | FULLY_FUNCTIONAL | Collapsible mobile menu   | JS toggle, nonce-based CSP                                        |
| `MobileMenuToggle` | FULLY_FUNCTIONAL | Hamburger button          | Conditional visibility                                            |
| `Footer`           | FULLY_FUNCTIONAL | Simple footer             | Copyright with dynamic year                                       |

### Known Issues

- **`NavLinkProps.Attrs` shadows `BaseProps.Attrs`** — split brain bug. Consumer attributes on BaseProps are silently ignored
- `Footer` could arguably live in `layout/` instead

---

## Cross-Cutting Features

- **CSP Compliance:** All inline scripts use `nonce` attribute
- **Dark Mode:** Full Tailwind `dark:` variant support via `layout.ThemeScript`
- **Tailwind Class Merging:** `utils.Class()` uses tailwind-merge-go for conflict resolution
- **Accessibility:** `aria-*` attributes, `role` attributes, screen-reader text, keyboard navigation (modal focus trap, dropdown arrows, tabs)
- **Responsive:** Mobile-first designs with `sm:` breakpoints
- **Type Safety:** 16+ typed string enums, `utils.BaseProps` embedded in all Props structs
- **Test Coverage:** 71.8% average across packages, BDD + snapshot + a11y + benchmark tests

---

## Planned / Not Yet Implemented

| Component       | Package     | Notes                                 |
| --------------- | ----------- | ------------------------------------- |
| Radio button    | `forms`     | Not present                           |
| File input      | `forms`     | Not present                           |
| Toggle/Switch   | `forms`     | Not present                           |
| Docs site       | —           | Auto-generated from source            |
| Demo app        | `examples/` | Working — Nav, Alert, StatCard, Icons |
| Release tooling | —           | goreleaser, tag-based                 |
