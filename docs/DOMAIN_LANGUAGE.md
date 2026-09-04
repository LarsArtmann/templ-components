# Domain Language

A **Unified Language** for **templ-components** — shared across developers and consumers.
Inspired by Domain-Driven Design (DDD) Ubiquitous Language.

## Glossary

| Term              | Definition                                                                         | Context                          |
| ----------------- | ---------------------------------------------------------------------------------- | -------------------------------- |
| Component         | A reusable UI building block with typed Go props                                   | `display.Card`, `feedback.Alert` |
| Props             | Typed configuration struct for a component                                         | `CardProps`, `AlertProps`        |
| BaseProps         | Shared fields (ID, Class, Attrs, AriaLabel, Nonce) embedded in all component props | `utils.BaseProps`                |
| ComponentProps    | Interface satisfied by all props structs via BaseProps promotion                   | `utils.ComponentProps`           |
| FeedbackType      | Enum for visual feedback severity: Success, Error, Warning, Info                   | `feedback.FeedbackType`          |
| TrendDirection    | Enum for stat change direction: Up, Down, None                                     | `display.TrendDirection`         |
| Eyebrow           | Small uppercase overline label above a title; reads as status, not decoration      | `display.Eyebrow`                |
| Scrollback        | Server-rendered terminal-style log block with a CSS-only staggered line entrance   | `display.Scrollback`             |
| ScrollbackTone    | Enum coloring a scrollback tag column: neutral, info, success, warning, danger     | `display.ScrollbackTone`         |
| FeedbackStyle     | Visual properties (color, icon, border) for a feedback variant                     | `feedback.feedbackStyleSet`      |
| FillIcon          | SVG rendered with `fill="currentColor"`; used for small 20x20 indicators           | `utils/svg.FillIcon`             |
| StrokeIcon        | SVG rendered with `stroke="currentColor"`; standard 24x24 UI icons                 | `icons.Icon`                     |
| IconPath          | SVG path data string; multi-path icons use a pipe separator                        | `icons.iconPathData`             |
| CardShell         | Shared CSS class for consistent card appearance (border, shadow, radius)           | `display.cardShellClass`         |
| ThemeColor        | CSS custom property for light/dark mode theming                                    | `layout.DefaultThemeColor`       |
| CSP Nonce         | Cryptographic nonce for Content Security Policy compliance                         | All `<script>` tags              |
| Event Delegation  | JS pattern: listeners on `document` for HTMX DOM swap compatibility                | Accordion, Dropdown, ThemeToggle |
| HTMX Error Family | Structured error classification for family-aware toast rendering                   | `htmx.ErrorHandlingConfig`       |
| Transport         | The client runtime that executes a wiring: `htmx` (default) or `datastar`           | `wire.Transport`                 |
| Action            | One hypermedia exchange (transport, method, URL, event, target) rendered in either dialect | `wire.Action`              |
| Wiring            | The act of connecting a component to a backend endpoint via a transport's attributes | `display.ButtonProps.Wire`       |
| Dialect           | The attribute syntax a transport speaks: `hx-*` or `data-on:*` + `@action()`        | `wire.Attributes()`              |
| Patch Mode        | How Datastar merges a fragment response into the region: inner, outer, append, ...  | `wire.PatchMode`                 |

## Entities

Objects with identity and lifecycle within the component tree.

| Term     | Definition                                                       | Context                 |
| -------- | ---------------------------------------------------------------- | ----------------------- |
| Page     | Full HTML document rendered by `layout.Base` or `layout.Minimal` | `layout.PageProps`      |
| Nav      | Top-level navigation bar with brand, links, mobile menu          | `navigation.NavProps`   |
| Modal    | Overlay dialog with focus trap and keyboard navigation           | `display.ModalProps`    |
| Dropdown | Button-triggered action menu with keyboard navigation            | `display.DropdownProps` |
| Table    | Data table with headers, rows, sortable columns, caption         | `display.TableProps`    |

## Value Objects

Immutable configuration objects.

| Term            | Definition                                               | Context                      |
| --------------- | -------------------------------------------------------- | ---------------------------- |
| PaginationProps | Page navigation state (current, total, URL construction) | `navigation.PaginationProps` |
| BreadcrumbItem  | Single segment in a breadcrumb trail                     | `navigation.BreadcrumbItem`  |
| DropdownItem    | Single action in a dropdown menu (link or button)        | `display.DropdownItem`       |
| SelectOption    | Single option in a select dropdown                       | `forms.SelectOption`         |
| TableCell       | Single cell with text or component content               | `display.TableCell`          |
| BadgeType       | Visual variant: neutral, success, warning, error, info   | `display.BadgeType`          |
| AvatarStatus    | Online state: online, offline, none                      | `display.AvatarStatus`       |
| CardPadding     | Internal spacing: none, sm, md, lg                       | `display.CardPadding`        |

## Platform terms (v0.20.0+)

| Term                  | Description                                                                                                                                                                                                                                                                                                                               | Where defined                                                                                                                                                                                                                                                   |
| --------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ContainerAware        | Opt-in flag making a component respond to its parent container width via `@container` queries (ADR-0018). 8 components total: `Grid.ContainerResponsive` (the original, differently named) + `Card`/`Nav`/`Split`/`Form`/`Pagination`/`DefinitionGrid`/`SkeletonCardGrid` `.ContainerAware`. Default flip to opt-out is a v2.0 candidate. | `NavProps.ContainerAware`, `CardProps.ContainerAware`, `SplitProps.ContainerAware`, `FormProps.ContainerAware`, `PaginationProps.ContainerAware`, `DefinitionGridProps.ContainerAware`, `SkeletonCardGridProps.ContainerAware`, `GridProps.ContainerResponsive` |
| Recipe                | Pre-composed screen layout built from primitives (Dashboard, SettingsLayout, LoginCard)                                                                                                                                                                                                                                                   | `recipes/` package                                                                                                                                                                                                                                              |
| Semantic Token        | CSS variable alias mapping a Tailwind palette color to a role (`--color-tc-primary` etc.)                                                                                                                                                                                                                                                 | `templates/templ-components-theme.css`                                                                                                                                                                                                                          |
| Theme Preset          | Drop-in CSS file with a complete `@theme` palette override (default, minimal, glass)                                                                                                                                                                                                                                                      | `templates/presets/`                                                                                                                                                                                                                                            |
| HTMXSrc               | `PageProps` field for self-hosting htmx; suppresses CDN preconnect, SRI, and response-targets ext                                                                                                                                                                                                                                         | `layout.PageProps.HTMXSrc`                                                                                                                                                                                                                                      |
| Popover API           | Native HTML `popover="auto"` attribute for click-toggle panels with light-dismiss + top-layer                                                                                                                                                                                                                                             | `display/Popover`, ADR-0017                                                                                                                                                                                                                                     |
| tc CLI                | Scaffolding tool: `tc init`, `tc ls`, `tc add <component>` — copies `.templ` to a consumer dir                                                                                                                                                                                                                                            | `cmd/tc/`                                                                                                                                                                                                                                                       |
| Fluid Typography      | CSS utility classes that scale font size relative to the nearest `@container` ancestor via `clamp(min, Ncqi + base, max)`. Six classes: `.tc-fluid-display`, `.tc-fluid-h1`–`h4`, `.tc-fluid-lead`. Zero JS, Baseline 2023.                                                                                                               | `templates/custom.css`, `docs/recipes/fluid-typography.md`                                                                                                                                                                                                      |
| Container Query Units | CSS length units (`cqi`, `cqw`, `cqh`, `cqmin`, `cqmax`) that resolve relative to the nearest `@container` ancestor's dimensions instead of the viewport. `cqi` = 1% of the container's inline size. Baseline 2023.                                                                                                                       | `templates/custom.css` (`.tc-fluid-*`)                                                                                                                                                                                                                          |
| Web Components        | Custom Elements + Shadow DOM + HTML Templates. **Permanently rejected** (ADR-0033): Shadow DOM breaks Tailwind utility-class theming, Custom Elements require JS, and the distribution problem does not exist for a Go-source library. The library achieves component encapsulation via native APIs instead.                              | ADR-0033                                                                                                                                                                                                                                                        |

## Visual testing & regression (v1.2.0+)

A separate Go module (`visualtest/`) renders components in headless Chromium and diffs pixels against committed golden images. It never enters the consumer dependency graph (separate `go.mod` + `replace`). Run via `nix run .#visual`. See `docs/visual-testing.md`.

| Term               | Definition                                                                                                                                                                                              | Where defined           |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------- |
| Visual Regression  | Detecting unintended pixel-level changes in rendered components (layout shifts, color regressions, dark-mode/RTL breakage) that HTML-string tests miss                                                  | `visualtest/`           |
| Golden (PNG)       | Committed baseline screenshot a test compares against; regenerated with `-update`                                                                                                                       | `visualtest/testdata/`  |
| chromedp           | Go bindings for the Chrome DevTools Protocol; drives headless Chromium for capture                                                                                                                      | `visualtest/harness.go` |
| pixelmatch         | Perceptual pixel-diff algorithm (YIQ color distance + anti-alias skip) via `orisano/pixelmatch`; decides pass/fail against `MaxMismatch` (default 0.1%)                                                 | `visualtest/compare.go` |
| AssertScreenshot   | Public test API: `AssertScreenshot(t, name, component, opts...)` with `Options{Dark, RTL, Viewport, MaxMismatch, Threshold, State}` (`StateRest`/`StateHover`/`StateFocus`/`StateClick`/`StateContext`) | `visualtest/harness.go` |
| `.fail/` artifacts | On mismatch, actual + diff PNGs are written to `.fail/` for review and auto-cleaned on pass                                                                                                             | `visualtest/golden.go`  |

## Bounded Contexts

| Context    | Description                                                         | Key Types                                         |
| ---------- | ------------------------------------------------------------------- | ------------------------------------------------- |
| Display    | Visual data presentation (cards, tables, badges, avatars, tabs)     | Card, Table, Badge, Avatar, Tabs                  |
| Feedback   | User-facing notifications (alerts, toasts, progress, spinners)      | Alert, Toast, ProgressBar, Spinner                |
| Forms      | User input controls (text, select, checkbox, radio, toggle, file)   | Input, Select, Textarea, Radio, Toggle, FileInput |
| Navigation | Page navigation (nav bars, pagination, breadcrumbs)                 | Nav, Pagination, Breadcrumbs                      |
| Layout     | Page-level structure (base HTML, theme, minimal)                    | Base, Minimal, ThemeScript                        |
| HTMX       | HTMX integration (loading, error handling, CSRF)                    | LoadingIndicator, GlobalErrorHandling, SwapOOB    |
| Icons      | SVG icon rendering with typed names                                 | Icon, IconWithStrokeWidth                         |
| ErrorPage  | Structured error presentation (page, detail, alert)                 | ErrorPage, ErrorDetail, ErrorAlert                |
| Recipes    | Top-of-DAG composition layer (Dashboard, SettingsLayout, LoginCard) | Dashboard, SettingsLayout, LoginCard              |
| CLI        | Scaffolding tool for copying components into consumer projects      | `tc` binary (`cmd/tc/`)                           |
