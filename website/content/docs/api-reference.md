---
title: API Reference
description: Package-level overview and links to full godoc.
---

## Packages

| Package      | Components | Purpose                                                                                                        |
| ------------ | ---------- | -------------------------------------------------------------------------------------------------------------- |
| `display`    | 43         | Cards, tables, modals, badges, buttons, avatars, carousel, tabs, accordion                                     |
| `feedback`   | 14         | Alerts, toasts, spinners, skeletons, progress bars                                                             |
| `forms`      | 23         | Inputs, selects, toggles, combobox, slider, rating, tags input                                                 |
| `layout`     | 10         | Page shell, theme toggle, CSP-safe script/style tags                                                           |
| `navigation` | 12         | Nav bars, pagination, breadcrumbs, sidebar, load-more                                                          |
| `htmx`       | 9          | Loading, error handling, OOB swaps, View Transitions                                                           |
| `datastar`   | 4          | Datastar runtime injection, SSE LiveRegion, loading Indicator                                                  |
| `icons`      | 105        | Heroicons v2 outline + Spinner                                                                                 |
| `errorpage`  | 4          | Error pages, 404, go-error-family integration                                                                  |
| `recipes`    | 4          | Screen-level compositions: Dashboard, SettingsLayout, LoginCard, AuthLayout                                    |
| `utils`      | —          | BaseProps, Class(), EnsureID, test helpers                                                                     |
| `utils/wire` | —          | Transport-agnostic wiring contract: one `wire.Action`, both HTMX and Datastar — see the Transport Wiring guide |

## Full Godoc

Complete API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/larsartmann/templ-components).

## Key Types

### BaseProps

Every props struct embeds `utils.BaseProps`:

```go
type BaseProps struct {
    ID        string
    Class     string
    Attrs     templ.Attributes
    AriaLabel string
    Nonce     string
}
```

### wire.Action

One typed wiring spec, rendered as htmx or Datastar attributes (see the
Transport Wiring guide):

```go
type Action struct {
    Transport      Transport   // "" (htmx default) | "htmx" | "datastar"
    Method         Method      // "" (GET — set explicitly on writes!) | get | post | put | patch | delete
    URL            string      // empty = inert (no wiring)
    Event          Event       // dialect default when empty
    Target         string      // htmx only (hx-target); Datastar targets via response headers
    ContentType    ContentType // "" (json) | json | form (form encoding)
    Selector       string      // Datastar form-encoding only: picks which <form> serializes; never targets patches
    Swap           PatchMode   // htmx hx-swap; Datastar's mode is response-driven (wire.Handler)
    DebounceMS     int         // htmx delay:<n>ms / datastar __debounce.<n>ms
    ThrottleMS     int         // htmx throttle:<n>ms / datastar __throttle.<n>ms
    PreventDefault bool        // datastar __prevent (runtime only auto-prevents form+submit)
    Interval       string      // polling: "10s" | "500ms" | "2m" | "1h" (m/h normalized to seconds)
    Reveal         *Reveal     // lazy-load on scroll into view
}
```

### Override Props

Some components expose additional typed class overrides so you can tweak their
internal layout without replacing the whole component. For example, `CardProps`
includes `TitleClass` and `HeaderClass`:

```templ
@display.Card(display.CardProps{
    Title:       "Users",
    TitleClass:  "text-indigo-600",
    HeaderClass: "bg-gray-50 dark:bg-gray-900/50",
}) {
    <p>Card content</p>
}
```

### Typed Enums

63 typed string enums make invalid states unrepresentable. Each ships with an `IsValid()` method.

```go
type BadgeType string
const (
    BadgePrimary BadgeType = "primary"
    BadgeNeutral BadgeType = "neutral"
    BadgeSuccess BadgeType = "success"
    BadgeWarning BadgeType = "warning"
    BadgeError   BadgeType = "error"
    BadgeInfo    BadgeType = "info"
)
```

### Lookup Maps

All style lookups use typed maps with `utils.Lookup()` fallback — no switches, no panics on unknown values.

## Import Paths

```go
import (
    "github.com/larsartmann/templ-components/display"
    "github.com/larsartmann/templ-components/feedback"
    "github.com/larsartmann/templ-components/forms"
    "github.com/larsartmann/templ-components/layout"
    "github.com/larsartmann/templ-components/navigation"
    "github.com/larsartmann/templ-components/htmx"
    "github.com/larsartmann/templ-components/icons"
    "github.com/larsartmann/templ-components/errorpage"
    "github.com/larsartmann/templ-components/utils"
)
```
