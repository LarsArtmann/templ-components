# templ-components

[![CI](https://img.shields.io/github/actions/workflow/status/larsartmann/templ-components/ci.yaml?branch=master&style=flat-square)](https://github.com/larsartmann/templ-components/actions)
[![Go Reference](https://img.shields.io/badge/go-pkg.go.dev-blue?style=flat-square)](https://pkg.go.dev/github.com/larsartmann/templ-components)
[![Go Report Card](https://goreportcard.com/badge/github.com/larsartmann/templ-components?style=flat-square)](https://goreportcard.com/report/github.com/larsartmann/templ-components)
[![License: MIT](https://img.shields.io/badge/License-MIT-green?style=flat-square)](https://github.com/larsartmann/templ-components/blob/master/LICENSE)
[![Version](https://img.shields.io/badge/version-v0.17.0-blue?style=flat-square)](https://github.com/larsartmann/templ-components/releases)

**Server-rendered UI components for Go web apps — built on [templ](https://templ.guide), [HTMX](https://htmx.org), and [Tailwind CSS v4](https://tailwindcss.com).**

[Documentation](https://templcomponents.lars.software) · [Quick Start](#quick-start) · [Component Catalog](#component-catalog)

No DaisyUI. No Node.js. No framework lock-in.

---

## Why templ-components?

97 server-rendered components. 34 typed string enums. 102 SVG icons. Zero client-side framework.

templ-components follows [HATEOAS](https://htmx.org/essays/hateoas/) — the server renders HTML, JavaScript enhances it rather than replacing it. Every component uses Tailwind CSS v4 utility classes with built-in dark mode, CSP nonce support, and ARIA accessibility.

| Feature                | templ-components           | [templUI](https://templui.io) | [goshipit](https://github.com/haatos/goshipit) |
| ---------------------- | -------------------------- | ----------------------------- | ---------------------------------------------- |
| **CSS approach**       | Tailwind v4 (CSS-first)    | Tailwind + CSS vars           | Tailwind + DaisyUI                             |
| **JavaScript**         | HATEOAS (enhances HTML)    | Alpine.js                     | DaisyUI JS                                     |
| **Requires Node.js**   | No                         | No                            | Yes                                            |
| **Components**         | 97                         | 40+                           | —                                              |
| **Typed props**        | 34 enums                   | —                             | —                                              |
| **Dark mode**          | Built-in (tested)          | CSS custom properties         | Via DaisyUI                                    |
| **CSP compliant**      | Yes (nonce on all scripts) | Yes                           | —                                              |
| **HTMX integration**   | Built-in package           | —                             | —                                              |
| **Standalone library** | Yes                        | No                            | No                                             |

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
    "github.com/larsartmann/templ-components/display"
    "github.com/larsartmann/templ-components/icons"
)

templ Dashboard() {
    @layout.Base(layout.DefaultPageProps()) {
        @layout.ThemeScript("")
        @display.PageHeader(display.PageHeaderProps{Title: "Dashboard"})
        @display.Grid(display.GridProps{Cols: display.GridCols3}) {
            @display.StatCard(display.StatCardProps{
                Label: "Revenue", Value: "$42,189", Trend: display.TrendUp,
            })
            @display.StatCard(display.StatCardProps{
                Label: "Users", Value: "1,204", Icon: icons.Users,
            })
        }
    }
}
```

**3. Generate and run**

```bash
templ generate && go run .
```

**Full guide:** [Installation](https://templcomponents.lars.software/getting-started/installation/) · [Quick Start](https://templcomponents.lars.software/getting-started/quick-start/)

---

## Component Catalog

### `display` — Data Display (30 components)

Cards, tables (Table + DataTable), tabs, modals, badges, buttons, avatars, tooltips, accordions, dropdowns, stat cards, page headers, definition lists, responsive grid, carousel, and more.

```templ
@display.Card(display.CardProps{Title: "Users", Subtitle: "Manage users"}) {
    <p>Card content</p>
}

@display.StatCard(display.StatCardProps{Label: "Users", Value: "1,204", Icon: icons.Users, Change: "12%", Trend: display.TrendUp})

@display.Grid(display.GridProps{Cols: display.GridCols3, Gap: display.GridGapLG}) {
    for _, u := range users {
        @display.Card(display.CardProps{Title: u.Name}) { <p>{ u.Email }</p> }
    }
}

@display.Table(display.TableProps{
    Headers: []string{"Name", "Email", "Role"},
    Rows: []display.TableRow{
        display.SimpleTableRow("Alice", "alice@example.com", "Admin"),
    },
    Striped: true,
})

@display.Modal(display.ModalProps{Title: "Confirm", Size: display.ModalSizeSM}) {
    <p>Are you sure?</p>
}
```

### `feedback` — User Feedback (13 components)

Alerts, toasts, spinners, progress bars, skeletons, step indicators, loading states.

```templ
@feedback.ToastContainer("")
@feedback.Toast(feedback.ToastProps{Message: "Saved!", Type: feedback.ToastSuccess})
@feedback.Alert(feedback.AlertProps{Title: "Warning", Type: feedback.AlertWarning})
@feedback.ProgressBar(feedback.ProgressBarProps{Current: 45, Total: 100})
@feedback.SkeletonCardGrid(6)
```

### `forms` — Form Controls (21 components)

Inputs, selects, textareas, checkboxes, radios, toggles, file inputs, date pickers, comboboxes, sliders, ratings, tags input, validation.

```templ
@forms.Input(forms.InputProps{Name: "email", Type: forms.InputEmail, Label: "Email"})
@forms.Select(forms.SelectProps{Name: "country", Label: "Country",
    Options: []forms.SelectOption{{Value: "de", Label: "Germany"}}})
@forms.Toggle(forms.ToggleProps{Name: "notifications", Label: "Enable notifications"})
@forms.Combobox(forms.ComboboxProps{Name: "country", Label: "Country",
    Options: []forms.ComboboxOption{{Value: "de", Label: "Germany"}}})
```

### `navigation` — Navigation (12 components)

Nav bars, breadcrumbs, pagination, mobile menus, sidebar, load-more.

```templ
@navigation.SimpleNav(navigation.SimpleNavProps{BrandText: "MyApp", CurrentPath: "/"})
@navigation.Breadcrumbs(navigation.BreadcrumbsProps{Items: []navigation.BreadcrumbItem{
    {Text: "Home", Href: "/"}, {Text: "Users", Active: true},
}})
@navigation.Pagination(navigation.PaginationProps{CurrentPage: 2, TotalPages: 10})
@navigation.SidebarNav(navigation.SidebarNavProps{CurrentPath: "/users"})
```

### `icons` — SVG Icons (102 icons)

Typed icon constants, no icon library dependency.

```templ
@icons.Icon(icons.Home, "h-5 w-5 text-gray-500")
@icons.Icon(icons.Check, "h-6 w-6 text-green-500")
```

### `htmx` — HTMX Integration (8 components)

Loading indicators, error handling, CSRF protection, out-of-band swaps, View Transitions.

```templ
@htmx.GlobalErrorHandling(htmx.DefaultErrorHandlingConfig())
@htmx.LoadingIndicator(feedback.Spinner(feedback.SpinnerMD, "text-blue-600"))
@htmx.ViewTransitions(htmx.ViewTransitionsProps{Global: true})
```

### `errorpage` — Error Pages (4 components)

Structured error pages with family-aware styling, HTTP handler integration, dedicated 404.

```templ
@errorpage.NotFound404(errorpage.DefaultNotFound404Props())
@errorpage.ErrorPage(errorpage.ErrorPageProps{Family: errorpage.FamilyNotFound, Why: "Not found"})
```

---

## Design Principles

**Type-safe.** 34 typed string enums make invalid states unrepresentable. Props structs embed `utils.BaseProps` for consistent ID, class, attributes, ARIA label, and CSP nonce propagation.

**Accessible.** ARIA attributes, roles, keyboard navigation, and screen-reader text across all interactive components. Native `<dialog>` for modals, `<details>` for accordions, `<search>` landmark for search inputs.

**CSP-ready.** All inline scripts use `nonce` attributes. No `eval()`, no inline event handlers. Integration test suite verifies compliance on every component.

**Dark mode.** Every component has proper `dark:` variants — enforced by `TestDarkModeCompliance` + `TestDarkModeSemanticColors` regression tests. `ThemeScript` prevents FOUC.

**Server-rendered.** Zero client-side JavaScript by default. Interactive features use minimal vanilla JS with nonce-based CSP.

**Pay for what you use.** Import only the packages you need. No monolithic bundle.

---

## Tailwind CSS Setup

Tailwind v4 uses CSS-first configuration. Vendor the dependency so Tailwind can scan the `.templ` source files:

```bash
go mod vendor
```

Then in your CSS:

```css
@import "tailwindcss";
@source "../vendor/github.com/larsartmann/templ-components";
@custom-variant dark (&:where(.dark, .dark *));
```

```bash
tailwindcss -i app.css -o styles.css --minify
```

If your project uses [BuildFlow](https://github.com/larsartmann/buildflow), the `tailwind-build` provider handles this automatically.

## Theming

Components emit standard Tailwind classes (`bg-blue-600`, `text-gray-900`). Override colors without touching component code:

```css
@theme {
  --color-blue-600: #4f46e5;
  --color-blue-500: #6366f1;
}
```

For semantic tokens (`bg-tc-primary`, `text-tc-danger`), copy the included [`templ-components-theme.css`](templ-components-theme.css).

See the [Theming guide](https://templcomponents.lars.software/guides/theming/) for details.

---

## By the Numbers

| Metric       | Value                                               |
| ------------ | --------------------------------------------------- |
| Components   | 97                                                  |
| SVG icons    | 102                                                 |
| Typed enums  | 34                                                  |
| Packages     | 9                                                   |
| Tests        | ~890 functions + ~1,650 subtests                    |
| Dependencies | 3 (`templ`, `tailwind-merge-go`, `go-error-family`) |

---

## Requirements

- **Go** 1.26+ (`GOEXPERIMENT=jsonv2`)
- **templ** CLI ([install](https://templ.guide/quick-start/installation))
- **Tailwind CSS** 4.x+
- **HTMX** 2.x (optional, for `htmx` package)

---

## Ecosystem

This library is part of the **GOTH stack** (Go + Templ + HTMX):

| Project                                                           | What it does                                                            |
| ----------------------------------------------------------------- | ----------------------------------------------------------------------- |
| [cqrs-htmx](https://github.com/LarsArtmann/cqrs-htmx)             | Production CQRS+ES framework with WebAuthn, RBAC, multi-tenancy, SSE.   |
| [go-cqrs-lite](https://github.com/larsartmann/go-cqrs-lite)       | Minimal CQRS/ES building blocks.                                        |
| [go-error-family](https://github.com/larsartmann/go-error-family) | Structured error families. Used by templ-components' errorpage package. |

---

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for setup, conventions, and workflow.

---

## License

[MIT](https://github.com/larsartmann/templ-components/blob/master/LICENSE)
