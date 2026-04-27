# templ-components

A reusable component library for Go projects using [templ](https://templ.guide/) and [Tailwind CSS](https://tailwindcss.com/).

Designed to be shared across all your Go projects. Import only the packages you need.

## Features

- **Layout** - Base HTML layouts, theme toggle, dark mode support
- **Feedback** - Toasts, alerts, loading spinners, skeletons, progress bars
- **Display** - Badges, cards, stat cards, empty states
- **Forms** - Inputs, textareas, selects, labels, validation errors
- **Navigation** - Navbars, breadcrumbs, mobile menus
- **Icons** - 40+ common SVG icons
- **HTMX** - Error handling, loading indicators, helpers
- **Utils** - Common Go helpers

## Installation

```bash
go get github.com/larsartmann/templ-components
```

Requires [templ](https://templ.guide/) CLI for code generation:

```bash
go install github.com/a-h/templ/cmd/templ@latest
```

## Quick Start

```go
package main

import (
    "github.com/a-h/templ"
    "github.com/larsartmann/templ-components/layout"
    "github.com/larsartmann/templ-components/feedback"
    "github.com/larsartmann/templ-components/display"
)

func MyPage() templ.Component {
    return layout.Base(layout.BaseProps{
        Title:       "Dashboard",
        Description: "My awesome dashboard",
    }) {
        @feedback.ToastContainer("")
        @display.Card(display.CardProps{Title: "Welcome"}) {
            <p>Hello, world!</p>
        }
    }
}
```

## Package Overview

### `layout` - Page Structure

```templ
@layout.Base(layout.BaseProps{Title: "My Page"}) {
    <h1>Content</h1>
}

@layout.ThemeToggle("Toggle theme", "")

@layout.Minimal("Title", "en") {
    <p>Static content</p>
}
```

### `feedback` - User Feedback

```templ
// Toast notifications
@feedback.ToastContainer("")

@feedback.Toast(feedback.ToastProps{
    Message: "Saved successfully!",
    Type:    feedback.ToastSuccess,
    Title:   "Success",
})

// Alerts
@feedback.Alert(feedback.AlertProps{
    Title:   "Warning",
    Message: "This action cannot be undone.",
    Type:    feedback.AlertWarning,
})

// Loading states
@feedback.Spinner(feedback.SpinnerMedium, "text-blue-600")
@feedback.InlineLoading("Saving...")
@feedback.Skeleton("card")

// Progress
@feedback.ProgressBar(feedback.ProgressBarProps{Current: 45, Total: 100})
@feedback.StepIndicator(feedback.StepIndicatorProps{Steps: []string{"Details", "Review", "Confirm"}, CurrentStep: 1})
```

### `display` - Data Display

```templ
@display.Badge(display.BadgeProps{Text: "Active", Type: display.BadgeSuccess, Dot: true})
@display.StatusBadge("healthy")

@display.Card(display.CardProps{Title: "Users", Subtitle: "Manage users"}) {
    <p>Card content</p>
}

@display.StatCard("1,234", "Total Users", "+12%", true)

@display.EmptyState(display.EmptyStateProps{
    Title:       "No repositories",
    Description: "Connect your first repository.",
    Icon:        "folder",
    ActionText:  "Connect",
    ActionHref:  "/connect",
})
```

### `forms` - Form Primitives

```templ
@forms.Input(forms.InputProps{
    Name:  "email",
    Type:  forms.InputEmail,
    Label: "Email address",
})

@forms.Textarea(forms.TextareaProps{
    Name:  "bio",
    Label: "Bio",
    Rows:  4,
})

@forms.Select(forms.SelectProps{
    Name:    "country",
    Label:   "Country",
    Options: []forms.SelectOption{
        {Value: "de", Label: "Germany"},
        {Value: "at", Label: "Austria"},
    },
})

@forms.Checkbox(forms.CheckboxProps{
    Name:  "terms",
    Label: "I agree to the terms",
})
```

### `navigation` - Navigation

```templ
@navigation.SimpleNav("MyApp", "/", []navigation.NavLinkProps{
    {Href: "/", Text: "Home"},
    {Href: "/about", Text: "About"},
}, "/")

@navigation.Breadcrumbs([]navigation.BreadcrumbItem{
    {Text: "Home", Href: "/"},
    {Text: "Users", Active: true},
})

@navigation.Footer("MyApp")
```

### `icons` - SVG Icons

```templ
@icons.Icon(icons.Home, "h-5 w-5 text-gray-500")
@icons.Icon(icons.Check, "h-6 w-6 text-green-500")
```

### `htmx` - HTMX Utilities

```templ
@htmx.GlobalErrorHandling("")
@htmx.LoadingIndicator()
@htmx.InlineLoadingOverlay("form-loading")
```

## Tailwind CSS Setup

This library assumes you have Tailwind CSS configured with the `dark` class strategy:

```js
// tailwind.config.js
module.exports = {
  darkMode: 'class',
  content: [
    './templates/**/*.templ',
    './node_modules/github.com/larsartmann/templ-components/**/*.templ',
  ],
  // ...
}
```

## Dark Mode

All components support dark mode via Tailwind's `dark:` prefix. Include the theme script in your base layout:

```templ
@layout.ThemeScript("")
```

And add the toggle button:

```templ
@layout.ThemeToggle("Toggle theme", "")
```

## Browser Support

- Modern evergreen browsers (Chrome, Firefox, Safari, Edge)
- Requires CSS custom properties support
- HTMX 2.x recommended

## License

MIT
