# Features — templ-components

**Updated:** 2026-05-03 | **Version:** 0.x (pre-release)

A Go component library built on [templ](https://templ.guide) and [Tailwind CSS](https://tailwindcss.com) for building server-rendered web applications.

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

**Totals:** 53 templ components, 42 types, 13 exported Go functions, 30 `.templ` files

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
| `Ptr`         | `[T any](v T) *T`                          | Pointer wrapper                               |
| `Deref`       | `[T any](p *T) T`                          | Nil-safe deref                                |
| `DerefOr`     | `[T any](p *T, fallback T) T`              | Nil-safe deref with fallback                  |

---

## Package: `display`

### Components

| Component          | Description                     | Key Features                                             |
| ------------------ | ------------------------------- | -------------------------------------------------------- |
| `Accordion`        | Collapsible accordion panels    | JS toggle, accessible `aria-expanded`, multiple items    |
| `Avatar`           | User avatar with image/initials | Online/offline indicator, 5 sizes, circle/square         |
| `Badge`            | Status label                    | 7 color types, 3 sizes, pill shape, dot indicator        |
| `StatusBadge`      | Auto-mapped status badge        | Maps ~20 status strings to badge types                   |
| `Card`             | Bordered card container         | Header, subtitle, footer, header action, 4 padding sizes |
| `SimpleCard`       | Minimal card                    | Children-only, no header/footer                          |
| `StatCard`         | Dashboard metric card           | Value, label, change indicator with up/down arrow        |
| `Dropdown`         | Button-triggered menu           | External/internal links, buttons, left/right position    |
| `EmptyState`       | Centered empty state            | Icon, title, description, action link/button             |
| `SimpleEmptyState` | Minimal empty state             | Text-only                                                |
| `Modal`            | Accessible modal dialog         | Backdrop blur, open/close animation, 5 sizes             |
| `Table`            | Responsive data table           | Headers, rows, striping, hover, caption, bordered        |
| `Tabs`             | Tabbed interface                | Default underline or pills style, active state           |
| `Tooltip`          | Hover tooltip                   | 4 positions, arrow indicator, wraps children             |

### Enums

| Type               | Values                                                   |
| ------------------ | -------------------------------------------------------- |
| `AvatarSize`       | XS, SM, MD, LG, XL                                       |
| `AvatarShape`      | Circle, Square                                           |
| `BadgeType`        | Default, Primary, Success, Warning, Error, Info, Neutral |
| `BadgeSize`        | Sm, Md, Lg                                               |
| `CardPadding`      | None, SM, MD, LG                                         |
| `DropdownPosition` | Left, Right                                              |
| `ModalSize`        | SM, MD, LG, XL, Full                                     |
| `TabsStyle`        | Default, Pills                                           |
| `TooltipPosition`  | Top, Bottom, Left, Right                                 |

---

## Package: `feedback`

### Components

| Component        | Description               | Key Features                                            |
| ---------------- | ------------------------- | ------------------------------------------------------- |
| `Alert`          | Full-width alert banner   | 4 types, dismissible, icon, CSP nonce                   |
| `InlineError`    | Compact error message     | Red icon + text                                         |
| `InlineSuccess`  | Compact success message   | Green checkmark + text                                  |
| `Spinner`        | Animated SVG spinner      | 3 sizes, custom color class                             |
| `LoadingOverlay` | Full-screen loading       | Spinner, message, optional progress bar                 |
| `InlineLoading`  | Compact inline loading    | Spinner + message                                       |
| `Skeleton`       | Pulsing placeholder       | 7 variants: text, title, avatar, image, card, table-row |
| `SkeletonGroup`  | Multiple skeletons        | Group with `animate-pulse`                              |
| `ProgressBar`    | Progress indicator        | 3 sizes, custom color, label, percentage                |
| `StepIndicator`  | Horizontal step indicator | Completed/current/future states, labels                 |
| `ToastContainer` | Fixed toast container     | JS `tcShowToast()` for dynamic toasts                   |
| `Toast`          | Server-rendered toast     | 4 types, dismissible, title, message, duration          |

### Enums

| Type              | Values                                                |
| ----------------- | ----------------------------------------------------- |
| `AlertType`       | Success, Error, Warning, Info                         |
| `SpinnerSize`     | Small, Medium, Large                                  |
| `SkeletonVariant` | Text, TextShort, Title, Avatar, Image, Card, TableRow |
| `ProgressBarSize` | SM, MD, LG                                            |
| `ToastType`       | Success, Error, Warning, Info                         |

### Constants

| Name                  | Value  | Purpose             |
| --------------------- | ------ | ------------------- |
| `ToastDurationShort`  | 3000ms | Short auto-dismiss  |
| `ToastDurationMedium` | 5000ms | Medium auto-dismiss |
| `ToastDurationLong`   | 8000ms | Long auto-dismiss   |

---

## Package: `forms`

### Components

| Component    | Description            | Key Features                                                        |
| ------------ | ---------------------- | ------------------------------------------------------------------- |
| `Input`      | Text input with label  | 11 types, error, help text, required, disabled, readonly, autofocus |
| `Checkbox`   | Checkbox with label    | Error, help text, required                                          |
| `Select`     | Select dropdown        | Options, disabled options, pre-selected                             |
| `Textarea`   | Textarea with label    | Configurable rows, error, help text                                 |
| `Label`      | Form label             | Optional `for` attribute, required indicator                        |
| `FieldError` | Field validation error | Optional field ID for `aria-describedby`                            |

### Enums

| Type        | Values                                                                        |
| ----------- | ----------------------------------------------------------------------------- |
| `InputType` | Text, Email, Password, Number, Tel, URL, Date, Time, Datetime, Search, Hidden |

### Functions

| Function     | Purpose                                 |
| ------------ | --------------------------------------- |
| `SanitizeID` | Converts arbitrary text to safe HTML ID |

---

## Package: `htmx`

### Components

| Component              | Description                | Key Features                                                   |
| ---------------------- | -------------------------- | -------------------------------------------------------------- |
| `LoadingIndicator`     | Fixed full-screen loading  | Uses `htmx-indicator`, blur backdrop                           |
| `InlineLoadingOverlay` | Localized loading overlay  | Absolute positioned, for form targets                          |
| `LoadingButton`        | Button with loading state  | Text swaps to spinner during HTMX requests                     |
| `ConfirmDelete`        | Delete button with confirm | `hx-delete`, `hx-target`, `hx-confirm`, `hx-swap`              |
| `SwapOOB`              | Out-of-band swap wrapper   | For updating multiple elements per response                    |
| `CSRFToken`            | Hidden CSRF input          | Standard `csrf_token` name                                     |
| `GlobalErrorHandling`  | HTMX error handler         | Network errors, response errors, auto-retry, toast integration |

---

## Package: `icons`

### Components

| Component | Description      | Key Features                                       |
| --------- | ---------------- | -------------------------------------------------- |
| `Icon`    | SVG icon by name | 42 named icons, custom class, currentColor theming |

### Icon Names (42)

Home, Users, Folder, Document, Search, Settings, Chart, Inbox, Check, X, Plus, Minus, ChevronRight, ChevronLeft, ChevronDown, ChevronUp, ArrowRight, ArrowLeft, Refresh, ExternalLink, Download, Upload, Trash, Edit, Eye, EyeOff, Lock, Unlock, Menu, Bell, Calendar, Clock, Location, Phone, Mail, Globe, Sun, Moon, Spinner, Exclamation, Information, Question

### Functions

| Function    | Purpose                                                   |
| ----------- | --------------------------------------------------------- |
| `IconAttrs` | Returns accessible attributes (aria-label or aria-hidden) |

---

## Package: `layout`

### Components

| Component     | Description           | Key Features                                                                                           |
| ------------- | --------------------- | ------------------------------------------------------------------------------------------------------ |
| `Base`        | Complete HTML5 layout | Meta tags, OG tags, Twitter cards, CSS, HTMX with SRI, theme script, skip-to-content, security headers |
| `Minimal`     | Minimal HTML document | No dependencies, for static pages/emails/PDFs                                                          |
| `ThemeScript` | Dark mode script      | localStorage-based, prevents FOUC                                                                      |
| `ThemeToggle` | Theme switch button   | Sun/moon icons, JS toggle, CSP nonce                                                                   |

### Functions

| Function           | Purpose                                                        |
| ------------------ | -------------------------------------------------------------- |
| `DefaultPageProps` | Returns sensible defaults (locale=en, HTMX 2.0.6, SRI enabled) |

---

## Package: `navigation`

### Components

| Component          | Description               | Key Features                                                      |
| ------------------ | ------------------------- | ----------------------------------------------------------------- |
| `Nav`              | Responsive navigation bar | Brand, links, right items, sticky option, mobile menu             |
| `SimpleNav`        | Simplified nav            | Text brand, sticky by default                                     |
| `NavLink`          | Desktop nav link          | Active state styling via currentPath                              |
| `MobileNavLink`    | Mobile nav link           | Border-left active indicator                                      |
| `Breadcrumbs`      | Breadcrumb navigation     | Chevron separators, active state                                  |
| `Pagination`       | Page navigation           | Mobile/desktop layouts, prev/next arrows, page range, query param |
| `MobileMenu`       | Collapsible mobile menu   | JS toggle, nonce-based CSP                                        |
| `MobileMenuToggle` | Hamburger button          | Conditional visibility                                            |
| `Footer`           | Simple footer             | Copyright with dynamic year                                       |

---

## Cross-Cutting Features

- **CSP Compliance:** All inline scripts use `nonce` attribute
- **Dark Mode:** Full Tailwind `dark:` variant support via `layout.ThemeScript`
- **Tailwind Class Merging:** `utils.Class()` uses tailwind-merge-go for conflict resolution
- **Accessibility:** `aria-*` attributes, `role` attributes, screen-reader text, keyboard navigation
- **Responsive:** Mobile-first designs with `sm:` breakpoints
- **Type Safety:** 16+ typed string enums, `utils.BaseProps` embedded in all Props structs
