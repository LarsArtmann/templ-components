# templ-components

Reusable UI components for Go web apps — built on [templ](https://templ.guide), [HTMX](https://htmx.org), and [Tailwind CSS](https://tailwindcss.com).

[![CI](https://img.shields.io/github/actions/workflow/status/larsartmann/templ-components/ci.yaml?branch=master&style=flat-square)](https://github.com/larsartmann/templ-components/actions)
[![Go Reference](https://img.shields.io/badge/go-pkg.go.dev-blue?style=flat-square)](https://pkg.go.dev/github.com/larsartmann/templ-components)
[![Go Report Card](https://goreportcard.com/badge/github.com/larsartmann/templ-components?style=flat-square)](https://goreportcard.com/report/github.com/larsartmann/templ-components)
[![License: MIT](https://img.shields.io/badge/License-MIT-green?style=flat-square)](https://github.com/larsartmann/templ-components/blob/master/LICENSE)
[![Pre-release](https://img.shields.io/badge/status-pre--release-orange?style=flat-square)](https://github.com/larsartmann/templ-components)

**No DaisyUI. No Node.js. No framework lock-in.**

templ-components is a pure Tailwind CSS v4 component library for Go's templ engine, with first-class HTMX integration. Every component renders server-side, ships zero JavaScript by default, and uses only Tailwind utility classes — so you stay in full control of your CSS and build pipeline.

> **Status:** Work in progress. No stability is guaranteed — APIs may change without notice at any time.

---

## Why templ-components?

| Feature              | templ-components            | [templUI](https://templui.io) | [goshipit](https://github.com/haatos/goshipit) |
| -------------------- | --------------------------- | ----------------------------- | ---------------------------------------------- |
| **CSS approach**     | Tailwind v4+ (CSS-first)    | Tailwind + CSS vars           | Tailwind + DaisyUI                             |
| **JavaScript**       | None (pure server-rendered) | Alpine.js                     | DaisyUI JS                                     |
| **Requires Node.js** | No                          | No                            | Yes                                            |
| **Go module**        | Yes                         | Yes                           | Yes                                            |
| **Components**       | 76                          | 40+                           | —                                              |
| **Dark mode**        | Built-in (Tailwind `dark:`) | CSS custom properties         | Via DaisyUI                                    |
| **CSP compliant**    | Yes (nonce support)         | Yes                           | —                                              |
| **Typed props**      | 27 string enums             | —                             | —                                              |
| **HTMX integration** | Built-in package            | —                             | —                                              |

**templ-components is for developers who want server-rendered HTML with zero client-side JavaScript — no Alpine.js, no DaisyUI, no build pipeline beyond `templ generate`.**

---

## Quick Start

**1. Install**

```bash
go get github.com/larsartmann/templ-components
```

**2. Build a page**

```templ
package main

import (
    "github.com/larsartmann/templ-components/layout"
    "github.com/larsartmann/templ-components/feedback"
    "github.com/larsartmann/templ-components/display"
)

templ Page() {
    @layout.Base(layout.DefaultPageProps()) {
        @layout.ThemeScript("")
        @feedback.ToastContainer("")
        @display.Card(display.DefaultCardProps()) {
            <h1>Hello, world</h1>
        }
    }
}
```

**3. Generate and run**

```bash
templ generate && go run .
```

---

## Component Catalog

### `utils` — Shared Types & Helpers

All component props embed `utils.BaseProps` for consistent `ID`, `Class`, `Attrs`, `AriaLabel`, and `Nonce` propagation. Key helpers:

```go
utils.Class("bg-white dark:bg-gray-800", props.Class)  // Tailwind class merging
utils.Ternary(condition, "yes", "no")                    // Generic ternary
```

### `layout` — Page Structure

Base HTML layouts, theme toggle, and dark mode support.

```templ
@layout.Base(layout.DefaultPageProps()) { <main>Content</main> }
@layout.Minimal(layout.DefaultMinimalProps()) { <p>Static content</p> }
@layout.ThemeScript("")
@layout.ThemeToggle("Toggle theme", "")
```

#### Suppressing auto-injected `<head>` tags

`DefaultPageProps()` auto-injects two tags that new consumers often want to
suppress. Override them in your `PageProps` literal:

```templ
props := layout.DefaultPageProps()
// 1. You bundle htmx yourself (e.g. via embedded /htmx.js), or don't use htmx:
props.HTMXVersion = ""               // suppresses <script src="...htmx.org...">
// 2. You serve styles another way (Tailwind Play CDN, inline, HeadContent):
props.CSSPath = ""                   // suppresses <link rel="stylesheet" href="/app.css">
```

If you load htmx from a different CDN or self-host, see `PageProps.HTMXCDN`.

### `display` — Data Display (20 components)

Cards, badges, modals, tables, tabs, avatars, tooltips, accordions, dropdowns, stat cards, page headers, definition lists, responsive grid, and more.

```templ
@display.Card(display.CardProps{Title: "Users", Subtitle: "Manage users"}) {
    <p>Card content</p>
}

@display.StatCard(display.StatCardProps{Label: "Users", Value: "1,204", Icon: icons.Users, Change: "12%", Trend: display.TrendUp})

@display.StatCard(display.StatCardProps{Label: "Active", Value: "42", Href: "/?activity=active"})
@display.Badge(display.BadgeProps{Text: "Active", Type: display.BadgeSuccess, Dot: true})
@display.StatusBadge("healthy")

@display.Grid(display.GridProps{Cols: display.GridCols3}) {
    for _, u := range users {
        @display.Card(display.CardProps{Title: u.Name}) { <p>{ u.Email }</p> }
    }
}

@display.Modal(display.ModalProps{Title: "Confirm", Size: display.ModalSizeSM}) {
    <p>Are you sure?</p>
}

@display.Table(display.TableProps{
    Headers: []string{"Name", "Email", "Role"},
    Rows: []display.TableRow{
        display.SimpleTableRow("Alice", "alice@example.com", "Admin"),
        display.SimpleTableRow("Bob", "bob@example.com", "User"),
    },
    Striped: true,
})

@display.Tabs(display.TabsProps{
    Tabs: []display.Tab{
        {ID: "users", Label: "Users", Active: true},
        {ID: "settings", Label: "Settings"},
    },
})

@display.Accordion(display.AccordionProps{
    Items: []display.AccordionItem{
        {ID: "faq1", Title: "What is this?", Open: true},
        {ID: "faq2", Title: "How does it work?"},
    },
})

@display.Avatar(display.AvatarProps{Src: "/avatar.jpg", Alt: "Alice"})
@display.Tooltip(display.TooltipProps{Text: "More info"}) { <button>Hover me</button> }
@display.Dropdown(display.DropdownProps{Label: "Actions", Items: []display.DropdownItem{
    {Text: "Edit", Href: "/edit"},
    {Text: "Delete", Href: "/delete"},
}})

@display.PageHeader(display.PageHeaderProps{Title: "Users", Subtitle: "Manage accounts"})
@display.DefinitionList(display.DefinitionListProps{Items: []display.DefinitionItem{
    {Term: "Email", Detail: "alice@example.com"},
    {Term: "Status", Detail: "Active"},
}})
@display.ListNote(display.ListNoteProps{Shown: 50, Total: 127})
```

### `feedback` — User Feedback (12 components)

Alerts, toasts, spinners, progress bars, skeletons, and loading states.

```templ
@feedback.ToastContainer("")
@feedback.Toast(feedback.ToastProps{Message: "Saved!", Type: feedback.ToastSuccess})

@feedback.Alert(feedback.AlertProps{
    Title: "Warning", Message: "This cannot be undone.", Type: feedback.AlertWarning,
})

@feedback.Spinner(feedback.SpinnerMD, "text-blue-600")
@feedback.Skeleton(feedback.SkeletonCard)

{{ /* Loading state for a card grid (pairs with display.Grid) */ }}
@feedback.SkeletonCardGrid(6)

@feedback.ProgressBar(feedback.ProgressBarProps{Current: 45, Total: 100})
@feedback.StepIndicator(feedback.StepIndicatorProps{
    Steps: []string{"Details", "Review", "Confirm"}, CurrentStep: 1,
})
```

### `forms` — Form Controls (16 components)

Inputs, selects, textareas, checkboxes, radios, toggles, file inputs, date pickers, comboboxes, labels, and validation errors.

```templ
@forms.Input(forms.InputProps{Name: "email", Type: forms.InputEmail, Label: "Email address"})
@forms.Textarea(forms.TextareaProps{Name: "bio", Label: "Bio", Rows: 4})
@forms.Select(forms.SelectProps{
    Name: "country", Label: "Country",
    Options: []forms.SelectOption{{Value: "de", Label: "Germany"}, {Value: "at", Label: "Austria"}},
})
@forms.Checkbox(forms.CheckboxProps{Name: "terms", Label: "I agree"})
@forms.RadioGroup(forms.RadioGroupProps{Name: "plan", Options: []forms.RadioOption{
    {Value: "free", Label: "Free"}, {Value: "pro", Label: "Pro"},
}})
@forms.Toggle(forms.ToggleProps{Name: "notifications", Label: "Enable notifications"})
@forms.FileInput(forms.FileInputProps{Name: "avatar", Label: "Upload photo", Accept: "image/*"})
@forms.DatePicker(forms.DatePickerProps{Name: "dob", Label: "Date of birth"})
@forms.Combobox(forms.ComboboxProps{Name: "country", Label: "Country", Options: []forms.ComboboxOption{
    {Value: "de", Label: "Germany"}, {Value: "at", Label: "Austria"},
}})
```

### `navigation` — Navigation (10 components)

Nav bars, breadcrumbs, pagination, mobile menus, and sidebar navigation.

```templ
@navigation.SimpleNav(navigation.SimpleNavProps{
    BrandText: "MyApp", BrandHref: "/",
    Links: []navigation.NavLinkProps{
        {Href: "/", Text: "Home"}, {Href: "/about", Text: "About"},
    },
    RightItems: @layout.ThemeToggle("Toggle theme", ""),
    CurrentPath: "/",
})

@navigation.Breadcrumbs(navigation.BreadcrumbsProps{
    Items: []navigation.BreadcrumbItem{
        {Text: "Home", Href: "/"}, {Text: "Users", Active: true},
    },
})

@navigation.Pagination(navigation.PaginationProps{CurrentPage: 2, TotalPages: 10, BaseURL: "/users"})
@navigation.Footer("MyApp")

@navigation.SidebarNav(navigation.SidebarNavProps{
    Brand: templ.Raw(`<span class="font-bold text-white">MyApp</span>`),
    Items: []navigation.SidebarNavItem{
        {Label: "Dashboard", Href: "/", Icon: icons.Squares2x2},
        {Label: "Users", Href: "/users", Icon: icons.Users},
    },
    CurrentPath: "/users",
})
```

### `icons` — SVG Icons (101 icons)

Typed icon constants, no icon library dependency.

```templ
@icons.Icon(icons.Home, "h-5 w-5 text-gray-500")
@icons.Icon(icons.Check, "h-6 w-6 text-green-500")
```

101 named icons (100 path icons + Spinner) covering navigation, UI actions, communication, media, status, and admin/security. See `icons/icon_names.go` for the full list.

### `htmx` — HTMX Integration (7 components)

Loading indicators, error handling, CSRF protection, and out-of-band swaps.

```templ
@htmx.GlobalErrorHandling(htmx.DefaultErrorHandlingConfig())
@htmx.LoadingIndicator(feedback.Spinner(feedback.SpinnerMD, "text-blue-600"))
@htmx.InlineLoadingOverlay("form-loading", feedback.Spinner(feedback.SpinnerSM, "text-white"))
@htmx.CSRFToken("token-value")
```

---

## Design Principles

**Type-safe.** 26 typed string enums make invalid states unrepresentable. Props structs embed `utils.BaseProps` for consistent ID, class, attributes, ARIA label, and CSP nonce propagation.

**Accessible.** ARIA attributes, roles, keyboard navigation (modal focus trap, dropdown arrows, tabs), and screen-reader text across all interactive components.

**CSP-ready.** All inline scripts use `nonce` attributes. No `eval()`, no inline event handlers.

**Dark mode.** Every component supports Tailwind's `dark:` variant. Include `@layout.ThemeScript("")` to prevent flash of unstyled content.

**Server-rendered.** Zero client-side JavaScript by default. Interactive features (accordion, dropdown, modal, theme toggle) use minimal vanilla JS with nonce-based CSP.

**Pay for what you use.** Import only the packages you need. No monolithic bundle.

```go
import (
    "github.com/larsartmann/templ-components/display"  // only if you use display components
    "github.com/larsartmann/templ-components/feedback" // only if you use feedback components
)
```

---

## Tailwind CSS Setup

Tailwind v4 uses CSS-first configuration. Since this is a Go module (not an npm package), you need to vendor the dependency so Tailwind can scan the `.templ` source files for class names:

```bash
go mod vendor
```

Then in your CSS:

```css
/* app.css */
@import "tailwindcss";

/* Scan templ-components for class names */
@source "../vendor/github.com/larsartmann/templ-components";

/* Enable class-based dark mode toggle */
@custom-variant dark (&:where(.dark, .dark *));
```

Tailwind extracts class names from any file content — including `.templ` files — so no additional configuration is needed beyond the `@source` path.

---

## Theming

Components emit standard Tailwind utility classes (`bg-blue-600`, `text-gray-900`). To customize colors **without touching component code**, override the underlying CSS variables in your `@theme` block:

```css
/* app.css */
@import "tailwindcss";
@source "../vendor/github.com/larsartmann/templ-components";

/* Override the primary brand color globally */
@theme {
  --color-blue-600: #4f46e5; /* changes ALL bg-blue-600, text-blue-600, etc. */
  --color-blue-500: #6366f1;
}
```

For a ready-made configuration with semantic aliases (`bg-tc-primary`, `text-tc-danger`, etc.), copy the included theme file into your project:

```css
@import "./templ-components-theme.css";
```

Or grab it from the repo root and customize:

```bash
curl -O https://raw.githubusercontent.com/larsartmann/templ-components/master/templ-components-theme.css
```

The theme file provides `@custom-variant dark` (required for `ThemeScript`/`ThemeToggle`) and semantic tokens like `--color-tc-primary`, `--color-tc-surface`, `--color-tc-danger` that you can use in your own CSS or Tailwind arbitrary values.

### Tailwind v4+ is the standard

**All LarsArtmann projects use Tailwind CSS v4+ (latest).** CSS-first config, no Node.js runtime, no DaisyUI. Small custom CSS only where Tailwind doesn't cover something. New to Tailwind v4? See the [`docs/tailwind-v4-adoption-guide.md`](docs/tailwind-v4-adoption-guide.md) for setup, theming, and migration from custom CSS design systems. Migrating from the Play CDN? See [`docs/migration/play-cdn-to-tailwind-v4.md`](docs/migration/play-cdn-to-tailwind-v4.md).

The `icons` package is the only CSS-agnostic package (pure SVG path data). See [`docs/icons-only-adoption.md`](docs/icons-only-adoption.md).

Wiring up HTMX error feedback in a server-rendered app? See [`docs/recipes/server-rendered-htmx-error-feedback.md`](docs/recipes/server-rendered-htmx-error-feedback.md).

---

## By the Numbers

| Metric       | Value                                               |
| ------------ | --------------------------------------------------- |
| Components   | 76                                                  |
| SVG icons    | 101                                                 |
| Typed enums  | 27                                                  |
| Packages     | 11                                                  |
| Tests        | 456 functions + ~660 subtests                       |
| Coverage     | ~75% (72–79% per package)                           |
| Dependencies | 3 (`templ`, `tailwind-merge-go`, `go-error-family`) |

---

## Requirements

- **Go** 1.26+
- **templ** CLI ([install](https://templ.guide/quick-start/installation))
- **Tailwind CSS** 4.x+
- **HTMX** 2.x (optional, for `htmx` package)

---

### `errorpage` — Error Pages (3 components)

Structured error pages with family-aware styling, HTTP handler integration, and pre-built constructors.

```templ
@errorpage.ErrorPage(errorpage.ErrorPageProps{Family: errorpage.FamilyNotFound, Why: "Page not found"})
@errorpage.ErrorDetail(errorpage.ErrorDetailProps{Family: errorpage.FamilyConflict, Why: "Conflict", Fix: "Refresh and retry"})
@errorpage.ErrorAlert(errorpage.ErrorAlertProps{Family: errorpage.FamilyTransient, Why: "Temporary issue"})
```

HTTP handler integration:

```go
http.Handle("/error", errorpage.ErrorHandler(err, errorpage.ErrorHandlerConfig{HTMLShell: true}))
http.Error(w, "not found", http.StatusNotFound) // use errorpage.NotFound() instead
```

---

## Ecosystem

This library is part of the **GOTH stack** (Go + Templ + HTMX):

| Project                                                           | What it does                                                                                                    |
| ----------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| [cqrs-htmx](https://github.com/LarsArtmann/cqrs-htmx)             | Production CQRS+ES framework with WebAuthn, RBAC, multi-tenancy, SSE. Uses templ-components in its admin panel. |
| [go-cqrs-lite](https://github.com/LarsArtmann/go-cqrs-lite)       | Minimal CQRS/ES building blocks (command bus, event store, projections, snapshots).                             |
| [go-error-family](https://github.com/LarsArtmann/go-error-family) | Structured error families (classified, contextual, actionable). Used by templ-components' errorpage package.    |

---

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for setup, conventions, and workflow.

---

## License

[MIT](https://github.com/larsartmann/templ-components/blob/master/LICENSE)
