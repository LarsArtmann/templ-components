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
| FeedbackStyle     | Visual properties (color, icon, border) for a feedback variant                     | `feedback.feedbackStyleSet`      |
| FillIcon          | SVG rendered with `fill="currentColor"`; used for small 20x20 indicators           | `internal/svg.FillIcon`          |
| StrokeIcon        | SVG rendered with `stroke="currentColor"`; standard 24x24 UI icons                 | `icons.Icon`                     |
| IconPath          | SVG path data string; multi-path icons use a pipe separator                        | `icons.iconPathData`             |
| CardShell         | Shared CSS class for consistent card appearance (border, shadow, radius)           | `display.cardShellClass`         |
| ThemeColor        | CSS custom property for light/dark mode theming                                    | `layout.DefaultThemeColor`       |
| CSP Nonce         | Cryptographic nonce for Content Security Policy compliance                         | All `<script>` tags              |
| Event Delegation  | JS pattern: listeners on `document` for HTMX DOM swap compatibility                | Accordion, Dropdown, ThemeToggle |
| HTMX Error Family | Structured error classification for family-aware toast rendering                   | `htmx.ErrorHandlingConfig`       |

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

## Bounded Contexts

| Context    | Description                                                       | Key Types                                         |
| ---------- | ----------------------------------------------------------------- | ------------------------------------------------- |
| Display    | Visual data presentation (cards, tables, badges, avatars, tabs)   | Card, Table, Badge, Avatar, Tabs                  |
| Feedback   | User-facing notifications (alerts, toasts, progress, spinners)    | Alert, Toast, ProgressBar, Spinner                |
| Forms      | User input controls (text, select, checkbox, radio, toggle, file) | Input, Select, Textarea, Radio, Toggle, FileInput |
| Navigation | Page navigation (nav bars, pagination, breadcrumbs)               | Nav, Pagination, Breadcrumbs                      |
| Layout     | Page-level structure (base HTML, theme, minimal)                  | Base, Minimal, ThemeScript                        |
| HTMX       | HTMX integration (loading, error handling, CSRF)                  | LoadingIndicator, GlobalErrorHandling, SwapOOB    |
| Icons      | SVG icon rendering with typed names                               | Icon, IconWithStrokeWidth                         |
| ErrorPage  | Structured error presentation (page, detail, alert)               | ErrorPage, ErrorDetail, ErrorAlert                |
