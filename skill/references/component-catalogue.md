# Component catalogue — templ-components

Full per-package reference: exact Go signatures and one-line behaviour notes for all
136 exported templ components, the 102-icon set, and the shared helper packages.

This file is the deep reference for `skill/SKILL.md` Part 1. The counts printed in the
package headings are machine-guarded against the `.templ` sources by
`utils.TestDocsCountDrift` in `skill/SKILL.md`; keep them in sync in the same commit as
any add/remove. Anything with a trailing `(docs)` pointer lives in the repo under `docs/`.

### By package (import path reference)

#### `display` — 52 components

| Component            | Signature                                           | One-liner                                                                                                                                                               |
| -------------------- | --------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Card`               | `Card(props CardProps)`                             | Bordered card with title, subtitle, footer, header action, 4 padding sizes, `ContainerAware`, sharp corners by default (`rounded-*` via `Class`, `.tc-squircle` opt-in) |
| `SimpleCard`         | `SimpleCard(props SimpleCardProps)`                 | Minimal card — children only, no header/footer                                                                                                                          |
| `StatCard`           | `StatCard(props StatCardProps)`                     | Dashboard metric card with value, label, change, trend, icon, semantic `Tone` (`StatTone`: blue/green/yellow/red/purple), optional `Href` link                          |
| `Grid`               | `Grid(props GridProps)`                             | Responsive grid — typed `GridCols` enum, `GridGap` enum, `ContainerAware` (default `true` since v2.0)                                                                   |
| `Badge`              | `Badge(props BadgeProps)`                           | Compact status label — 7 types, 3 sizes, pill, dot, optional `Href`                                                                                                     |
| `StatusBadge`        | `StatusBadge(status string)`                        | Auto-maps ~20 status strings to badge types                                                                                                                             |
| `Button`             | `Button(props ButtonProps)`                         | Button or link button — variants, sizes, icons, HTMX attrs                                                                                                              |
| `Avatar`             | `Avatar(props AvatarProps)`                         | Image avatar with fallback initials, sizes, status dot                                                                                                                  |
| `Modal`              | `Modal(props ModalProps)`                           | Native `<dialog>` overlay — showModal(), focus trap, Escape, backdrop, 5 sizes                                                                                          |
| `Drawer`             | `Drawer(props DrawerProps)`                         | Native `<dialog>` side panel — slide animation, focus trap, Escape, backdrop                                                                                            |
| `Dropdown`           | `Dropdown(props DropdownProps)`                     | Button-triggered menu — links, buttons, keyboard nav                                                                                                                    |
| `Tooltip`            | `Tooltip(props TooltipProps)`                       | Hover tooltip — 6 positions (logical start/end + legacy physical left/right), arrow, touch support                                                                      |
| `Popover`            | `Popover(props PopoverProps)`                       | Button-triggered floating panel — 6 positions (logical start/end + legacy physical left/right), arbitrary content, click-outside dismiss                                |
| `Accordion`          | `Accordion(props AccordionProps)`                   | Collapsible sections — open/closed state, keyboard nav                                                                                                                  |
| `Table`              | `Table(props TableProps)`                           | Responsive data table — striping, hover, caption, bordered, `Body` slot, `Row.Href`, `Flush` (card nesting), `CellPadding` (comfortable/compact)                        |
| `DataTable`          | `DataTable(props DataTableProps)`                   | Data table with integrated sort management, pagination, empty-state — composes `Table` internally                                                                       |
| `Tabs`               | `Tabs(props TabsProps)`                             | Tabbed interface — underline/pills, optional client-side JS, `Wire` for server-side tab switching (`{tab}` URL placeholder)                                             |
| `EmptyState`         | `EmptyState(props EmptyStateProps)`                 | Empty-data placeholder — icon, title, description, action                                                                                                               |
| `SimpleEmptyState`   | `SimpleEmptyState(message string)`                  | Minimal empty state — text only                                                                                                                                         |
| `PageHeader`         | `PageHeader(props PageHeaderProps)`                 | Page title block — title, subtitle, breadcrumb, action slots                                                                                                            |
| `DefinitionList`     | `DefinitionList(props DefinitionListProps)`         | Two-column `<dl>` key/value list                                                                                                                                        |
| `ListNote`           | `ListNote(props ListNoteProps)`                     | "Showing N of M" truncation notice; `Variant: ListNoteCount` = count-only "Showing N items." (always renders — range/list semantics)                                    |
| `Eyebrow`            | `Eyebrow(props EyebrowProps)`                       | Small uppercase monospace overline above a title; accent color via `Class`                                                                                              |
| `Scrollback`         | `Scrollback(props ScrollbackProps)`                 | Terminal-style log block — timestamp + tone-colored tags, CSS-only staggered entrance, reduced-motion safe, zero JS                                                     |
| `CopyButton`         | `CopyButton(props CopyButtonProps)`                 | Clipboard copy button or link — CSP-safe, "Copied!" feedback, optional `Href` variant                                                                                   |
| `RelativeTime`       | `RelativeTime(props RelativeTimeProps)`             | `<time datetime>` with relative text ("2 hours ago")                                                                                                                    |
| `CountBadge`         | `CountBadge(props CountBadgeProps)`                 | Icon + notification count overlay — overflow "N+"                                                                                                                       |
| `DefinitionGrid`     | `DefinitionGrid(props DefinitionGridProps)`         | Responsive grid of term-detail cards, `ContainerAware`                                                                                                                  |
| `Image`              | `Image(props ImageProps)`                           | Lazy-loaded `<img>` with CSP-safe fallback, `SrcSet`/`Sizes` for responsive delivery, `Placeholder` (zero-JS blur-up), optional `Rounded`                               |
| `HoverCard`          | `HoverCard(props HoverCardProps)`                   | CSS-only hover-activated card — 4 positions, focus-within support                                                                                                       |
| `ContextMenu`        | `ContextMenu(props ContextMenuProps)`               | Right-click context menu — CSP-safe JS, role=menu, Escape/click-outside dismiss                                                                                         |
| `Carousel`           | `Carousel(props CarouselProps)`                     | Slide carousel with prev/next arrows and dot indicators                                                                                                                 |
| `KanbanBoard`        | `KanbanBoard(props KanbanBoardProps)`               | Dual-transport kanban — drag-and-drop + keyboard moves via one hidden form, per-column `Action` add-card slot, `Tone` status dot (`KanbanTone`), read-only without Wire |
| `Sparkline`          | `Sparkline(props SparklineProps)`                   | Pure SVG mini line chart — trend visualization, `Filled`, `AriaLabel`, customizable dims                                                                                |
| `BarChart`           | `BarChart(props BarChartProps)`                     | CSS bar chart (horizontal/vertical) — per-bar colors, links, value formatting, empty state                                                                              |
| `ExternalLink`       | `ExternalLink(props ExternalLinkProps)`             | Safe off-site link — `target="_blank" rel="noopener"`, URL-sanitized href, optional arrow icon                                                                          |
| `CollapsibleSection` | `CollapsibleSection(props CollapsibleSectionProps)` | Native `<details>/<summary>` — configurable heading, open/closed default, optional `StorageKey` for localStorage persistence, chevron rotation                          |
| `CommandPalette`     | `CommandPalette(props CommandPaletteProps)`         | ⌘K command palette — native `<dialog>`, client-side filter over server-rendered groups, link or `Wire` selection, arrow/Enter navigation                                |
| `Heatmap`            | `Heatmap(props HeatmapProps)`                       | CSS grid heatmap — row/column labels, opacity-based coloring, peak highlighting, tooltips, empty state                                                                  |
| `LineChart`          | `LineChart(props LineChartProps)`                   | Pure SVG line chart — axes, gridlines, multi-series, dots, legend, linear/smooth styles, ARIA. Zero JS (ADR-0031)                                                       |
| `PieChart`           | `PieChart(props PieChartProps)`                     | Pure SVG pie/donut chart — arc paths, external labels, legend, center label, custom colors, ARIA. Zero JS                                                               |
| `AreaChart`          | `AreaChart(props AreaChartProps)`                   | Pure SVG area chart — filled areas, multi-series, fill opacity, smooth curves, ARIA. Zero JS                                                                            |
| `CodeBlock`          | `CodeBlock(props CodeBlockProps)`                   | Copyable code display — block (figure+pre/code, language/label header) or compact-ID variant, integrated CopyButton, HTML-escaped                                       |
| `DataState`          | `DataState(props DataStateProps)`                   | Honesty ladder — Content (children) / Empty / Unavailable / Disabled; composes EmptyState, per-state default icons                                                      |
| `StatusDot`          | `StatusDot(props StatusDotProps)`                   | Colored status dot — 5 tones, optional pulse ring (motion-reduce safe), sr-only Label for a11y                                                                          |
| `LivePill`           | `LivePill(props LivePillProps)`                     | Dot+text status pill — the Live / Reconnecting / Paused connection indicator shape                                                                                      |
| `SegmentBar`         | `SegmentBar(props SegmentBarProps)`                 | Stacked proportion bar + legend — flex-grow proportions, deterministic label-keyed palette, `Color` override, role=img on the bar + real legend list                    |
| `Lightbox`           | `Lightbox(props LightboxProps)`                     | Native-dialog image viewer — thumbnails or trigger button, prev/next, rotate, zoom, captions; CSP-safe singleton script                                                 |
| `Prose`              | `Prose(props ProseProps)`                           | Longform-content container — `.tc-prose` descendant styling for rendered markdown (headings/lists/quotes/code/tables/images), typed `ProseSize`, dark+RTL safe          |

#### `forms` — 26 components

| Component           | Signature                                         | One-liner                                                                                                                            |
| ------------------- | ------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `Input`             | `Input(props InputProps)`                         | Typed input (text, email, password, search, url, tel, number) with label, error, help text                                           |
| `Checkbox`          | `Checkbox(props CheckboxProps)`                   | Checkbox with label, help text, error                                                                                                |
| `Select`            | `Select(props SelectProps)`                       | Select dropdown with typed options, disabled, selected, optgroups; `Stylable` for customizable `<select>` API                        |
| `Textarea`          | `Textarea(props TextareaProps)`                   | Multi-line text input with label, rows, max length; `AutoGrow` (field-sizing), `EnterKeyHint`                                        |
| `Toggle`            | `Toggle(props ToggleProps)`                       | Toggle switch with label, required, error, help text                                                                                 |
| `Radio`             | `Radio(props RadioProps)`                         | Radio button with label                                                                                                              |
| `RadioGroup`        | `RadioGroup(props RadioGroupProps)`               | Group of radios with options                                                                                                         |
| `Combobox`          | `Combobox(props ComboboxProps)`                   | Autocomplete input with filterable dropdown                                                                                          |
| `DatePicker`        | `DatePicker(props DatePickerProps)`               | Date input with label                                                                                                                |
| `FileInput`         | `FileInput(props FileInputProps)`                 | File upload input with label, accept types                                                                                           |
| `Form`              | `Form(props FormProps)`                           | Form wrapper with action, method, CSRF token, `ContainerAware`, `Wire` for dual-transport submit (fields serialize in both dialects) |
| `Label`             | `Label(forID, text string, required bool)`        | Form label element                                                                                                                   |
| `FieldError`        | `FieldError(fieldID, message string)`             | Inline field error message                                                                                                           |
| `FormFieldWrapper`  | `FormFieldWrapper(props FormFieldProps)`          | Wraps inputs with label, error, help text                                                                                            |
| `InputGroup`        | `InputGroup(props InputGroupProps)`               | Input with prefix/suffix                                                                                                             |
| `ValidationSummary` | `ValidationSummary(props ValidationSummaryProps)` | Accessible error summary with icon, count, linked fields                                                                             |
| `FilterDropdown`    | `FilterDropdown(props FilterDropdownProps)`       | Auto-submit select for filter bars — HTMX-native or dual-transport via `Wire`                                                        |
| `FilterInput`       | `FilterInput(props FilterInputProps)`             | Debounced dual-transport search input — `<search>` form, auto-submits as you type (`Wire` + `DebounceMS`)                            |
| `FilterBar`         | `FilterBar(props FilterBarProps)`                 | Horizontal auto-submit filter row — GET form, composite change trigger, Reset, noscript Apply, sticky opt-in, `Wire` dual-transport  |
| `FilterChips`       | `FilterChips(props FilterChipsProps)`             | Zero-JS filter chips — link-based, `aria-current`, `FilterToggleHref` query builder; crawler-friendly                                |
| `SegmentedControl`  | `SegmentedControl(props SegmentedControlProps)`   | Segmented single-select — radio mode (radiogroup, submits with the form) or link mode (aria-current); zero JS                        |
| `Slider`            | `Slider(props SliderProps)`                       | Range input slider with label, value display, help text                                                                              |
| `Rating`            | `Rating(props RatingProps)`                       | Accessible star rating — radio inputs, read-only mode, 3 sizes                                                                       |
| `TagsInput`         | `TagsInput(props TagsInputProps)`                 | Tag input with add/remove — hidden inputs for submission, MaxTags/AllowDuplicate                                                     |
| `Calendar`          | `Calendar(props CalendarProps)`                   | Month-view calendar — server-side navigation, day links, MinDate/MaxDate disabling                                                   |

#### `feedback` — 14 components

| Component          | Signature                                       | One-liner                                                                                                   |
| ------------------ | ----------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `Alert`            | `Alert(props AlertProps)`                       | Inline alert — 4 types, dismissible, icon                                                                   |
| `InlineError`      | `InlineError(message string)`                   | Compact inline error                                                                                        |
| `InlineSuccess`    | `InlineSuccess(message string)`                 | Compact inline success                                                                                      |
| `Toast`            | `Toast(props ToastProps)`                       | Toast notification — 4 types, dismissible, duration                                                         |
| `ToastContainer`   | `ToastContainer(nonce string)`                  | Fixed container for toasts (include once per page)                                                          |
| `Spinner`          | `Spinner(props SpinnerProps)`                   | Loading spinner — 3 sizes, custom color                                                                     |
| `LoadingOverlay`   | `LoadingOverlay(props LoadingOverlayProps)`     | Full-screen loading overlay with spinner                                                                    |
| `InlineLoading`    | `InlineLoading(message string)`                 | Inline "Saving..." indicator                                                                                |
| `Skeleton`         | `Skeleton(variant SkeletonVariant)`             | Pulsing placeholder — 7 variants                                                                            |
| `SkeletonGroup`    | `SkeletonGroup(variants []SkeletonVariant)`     | Multiple skeletons under single `role="status"`                                                             |
| `SkeletonCardGrid` | `SkeletonCardGrid(props SkeletonCardGridProps)` | N skeleton cards in responsive grid — loading state, optional `ContainerAware`                              |
| `ProgressBar`      | `ProgressBar(props ProgressBarProps)`           | Progress indicator — 3 sizes, indeterminate, label                                                          |
| `StepIndicator`    | `StepIndicator(props StepIndicatorProps)`       | Horizontal/vertical step progress — typed `StepState` (complete/current/pending) derived from `CurrentStep` |

#### `layout` — 11 components

| Component     | Signature                              | One-liner                                                                                                                                    |
| ------------- | -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| `Base`        | `Base(props PageProps)`                | Full HTML5 shell — head, meta, theme, HTMX, CSS auto-inject, SEO (noindex/canonical/hreflang/JSON-LD)                                        |
| `Minimal`     | `Minimal(props MinimalProps)`          | Minimal HTML doc — no dependencies, for static/PDF                                                                                           |
| `ThemeScript` | `ThemeScript(nonce string)`            | Dark mode script — prevents FOUC, include in `<head>`                                                                                        |
| `ThemeToggle` | `ThemeToggle(ariaLabel, nonce string)` | Dark/light toggle button with sun/moon icons                                                                                                 |
| `Script`      | `Script(nonce, src string, attrs)`     | CSP-safe `<script src>` — auto-injects nonce                                                                                                 |
| `Stylesheet`  | `Stylesheet(href, attrs)`              | CSP-safe `<link rel="stylesheet">`                                                                                                           |
| `MetaRefresh` | `MetaRefresh(delaySeconds, url)`       | Zero-JS reload/redirect via `<meta http-equiv=refresh>` — retry/stale-session pattern; empty url reloads                                     |
| `AppShell`    | `AppShell(props AppShellProps)`        | Sidebar + header + main app shell — the #1 admin dashboard pattern                                                                           |
| `Container`   | `Container(props ContainerProps)`      | Centered max-width wrapper with responsive gutter ON by default (`NoPad: true` opts out) — replaces `max-w-Nxl mx-auto px-4 sm:px-6 lg:px-8` |
| `Split`       | `Split(props SplitProps)`              | 2-col content+aside — article+sidebar, detail+metadata (RTL-aware), `ContainerAware`                                                         |
| `Stack`       | `Stack(props StackProps)`              | Vertical rhythm — typed Gap enum, replaces repeated `space-y-N`                                                                              |

#### `navigation` — 12 components

| Component          | Signature                                               | One-liner                                                                                                   |
| ------------------ | ------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `Nav`              | `Nav(props NavProps)`                                   | Full nav bar — brand, links, right items slot, mobile menu, `ContainerAware` (hamburger by container width) |
| `SimpleNav`        | `SimpleNav(props SimpleNavProps)`                       | Simplified nav — text brand, links, `RightItems` slot                                                       |
| `NavLink`          | `NavLink(props NavLinkProps, currentPath string)`       | Single nav link with active detection                                                                       |
| `MobileNavLink`    | `MobileNavLink(props NavLinkProps, currentPath string)` | Mobile menu link                                                                                            |
| `MobileMenu`       | `MobileMenu(links, currentPath, nonce, menuID)`         | Responsive mobile menu                                                                                      |
| `MobileMenuToggle` | `MobileMenuToggle(show bool, menuID string)`            | Hamburger toggle button                                                                                     |
| `Pagination`       | `Pagination(props PaginationProps)`                     | Page-number pagination — ellipsis, prev/next, SEO rel, `ContainerAware`                                     |
| `Breadcrumbs`      | `Breadcrumbs(props BreadcrumbsProps)`                   | Breadcrumb trail — separator, JSON-LD                                                                       |
| `SidebarNav`       | `SidebarNav(props SidebarNavProps)`                     | Vertical sidebar — brand, icon+label items, footer slot                                                     |
| `LoadMore`         | `LoadMore(props LoadMoreProps)`                         | Cursor-based "Load more" button — hx-get, hx-swap, optional `InfiniteScroll`                                |
| `EndOfList`        | `EndOfList(props EndOfListProps)`                       | "You've reached the end" indicator — companion to LoadMore/Pagination                                       |
| `Footer`           | `Footer(props FooterProps)`                             | Simple footer with copyright                                                                                |

#### `htmx` — 9 components

| Component              | Signature                                                 | One-liner                                                                                 |
| ---------------------- | --------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| `LoadingIndicator`     | `LoadingIndicator(spinner templ.Component)`               | Fixed HTMX loading overlay                                                                |
| `InlineLoadingOverlay` | `InlineLoadingOverlay(id string, spinner)`                | Per-target loading overlay                                                                |
| `LoadingButton`        | `LoadingButton(defaultText, loadingText string, spinner)` | Button with loading state                                                                 |
| `ConfirmDelete`        | `ConfirmDelete(props ConfirmDeleteProps)`                 | HTMX delete confirmation                                                                  |
| `SwapOOB`              | `SwapOOB(props SwapOOBProps)`                             | HTMX out-of-band swap                                                                     |
| `CSRFToken`            | `CSRFToken(token string)`                                 | Hidden CSRF input                                                                         |
| `GlobalErrorHandling`  | `GlobalErrorHandling(cfg ErrorHandlingConfig)`            | HTMX error → toast pipeline                                                               |
| `PolledRegion`         | `PolledRegion(props PolledRegionProps)`                   | Auto-refreshing HTMX region — interval polling, eager-load, aria-live, optional timestamp |

#### `datastar` — 6 components + action helpers

| Component          | Signature                                      | One-liner                                                                                                                                                    |
| ------------------ | ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `SDKScript`        | `SDKScript(props SDKScriptProps)`              | Datastar runtime injection — versioned CDN, self-hostable via `Src`, CSP nonce                                                                               |
| `LiveRegion`       | `LiveRegion(props LiveRegionProps)`            | SSE-powered live region — `data-init` auto-connect, `RetryMode` self-healing, `LivePoliteness`, busy cue                                                     |
| `PolledRegion`     | `PolledRegion(props PolledRegionProps)`        | Interval-polling region (Datastar twin of htmx's) — `data-on-interval__duration.<n>`, m/h normalized to seconds, empty URL inert, timestamp, NO leading flag |
| `Indicator`        | `Indicator(props IndicatorProps)`              | Loading indicator driven by a Datastar indicator signal                                                                                                      |
| `LoadingButton`    | `LoadingButton(props LoadingButtonProps)`      | Busy-state button label (Datastar twin of htmx's) — `data-indicator` signal + `data-show` swap, zero component JS, empty signal degrades inert               |
| `SSEErrorHandling` | `SSEErrorHandling(cfg SSEErrorHandlingConfig)` | SSE error feedback on the real `datastar-fetch` lifecycle events                                                                                             |

#### `errorpage` — 4 components + 6 constructors + handler

| Component / Function           | Signature                             | One-liner                                                                                                                                                                                                                             |
| ------------------------------ | ------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `ErrorPage`                    | `ErrorPage(props ErrorPageProps)`     | Full-page error display (neutral card, family accent bar, `HTTP {StatusCode}`/code/family chip row, primary + ghost secondary action, `MaxWidth`, timestamp + trace footer — 2026-09-17 redesign; 4-capture mobile/RTL visual shield) |
| `NotFound404`                  | `NotFound404(props NotFound404Props)` | Dedicated 404 page — hero numeral, search, links, configurable `LinksTitle`                                                                                                                                                           |
| `ErrorDetail`                  | `ErrorDetail(props ErrorDetailProps)` | Inline error card — `Variant` picks Tinted (default) or Neutral accent-bar shell                                                                                                                                                      |
| `ErrorAlert`                   | `ErrorAlert(props ErrorAlertProps)`   | Family-aware alert                                                                                                                                                                                                                    |
| `ErrorHandler`                 | `ErrorHandler(err, cfg) http.Handler` | go-error-family aware HTTP handler                                                                                                                                                                                                    |
| `WriteError`                   | `WriteError(w, r, err, nonce)`        | One-call error page from any handler                                                                                                                                                                                                  |
| `FromError`                    | `FromError(err) ErrorPageProps`       | Extract family/title/code/why/fix/trace/status from error; prefers `Public()` user-safe message, derives family title when none set                                                                                                   |
| `NotFound` ... `InternalError` | 6 constructors                        | Pre-built error page props by HTTP family                                                                                                                                                                                             |

#### `icons` — 102 icons + 3 functions

| Function                                        | One-liner                                                      |
| ----------------------------------------------- | -------------------------------------------------------------- |
| `Icon(name, class)`                             | Inline SVG icon by name (101 path icons + Spinner)             |
| `IconWithStrokeWidth(name, class, strokeWidth)` | Same, custom stroke width (default 1.5)                        |
| `IconPathData(name) []string`                   | Raw path data for custom SVG wrapper                           |
| `IconPathJS(name) string`                       | Path data formatted for JS injection                           |
| `Render(icon CustomIcon, class)`                | Render a consumer-provided icon (arbitrary viewBox/paths/fill) |
| `AnimatedIcon(name, class)`                     | Icon with hover animation (default per-icon)                   |
| `AnimatedIconWithAnimation(name, anim, class)`  | Icon with explicit animation preset                            |
| `DefaultAnimation(name) Animation`              | Get default animation for an icon                              |
| `AllAnimations() []Animation`                   | All 11 animation presets, sorted                               |

#### `utils` — shared types + helpers

| Function                           | One-liner                                                        |
| ---------------------------------- | ---------------------------------------------------------------- |
| `BaseProps`                        | Embed in every props struct — ID, Class, Attrs, AriaLabel, Nonce |
| `ComponentProps`                   | Interface: `GetBaseProps()` / `SetBaseProps()`                   |
| `Class(classes ...string)`         | Tailwind class merge (tailwind-merge-go, mutex-protected)        |
| `Lookup(m, key, fallback)`         | Map lookup with fallback for style variants                      |
| `Ternary(cond, a, b)`              | Functional ternary for templ expressions                         |
| `EnsureID(prefix, id)`             | Auto-generate DOM-safe IDs (crypto/rand)                         |
| `ScriptAttrs(nonce)`               | Script attrs with the omit-empty nonce rule (never `nonce=""`)   |
| `ScriptComponent(nonce, js, name)` | The canonical CSP-safe inline-script writer for custom scripts   |
| `Version`                          | Library version string (matches CHANGELOG)                       |

#### `utils/format` — display formatting (pure functions, no deps)

| Function                    | One-liner                                                             |
| --------------------------- | --------------------------------------------------------------------- |
| `format.Bytes(n)`           | IEC bytes: `1.5 KiB`, `42 MiB` — sign preserved, no panic on overflow |
| `format.CompactDuration(d)` | Operator style: `42s`, `1m30s`, `12h30m` (no Go `String` padding)     |
| `format.ClockDuration(d)`   | Display style: `3:45`, `1:02:03`, `2d 5h`; negatives clamp to `0:00`  |
| `format.Percent(ratio)`     | 0–1 → `87.3%`; >100 honest; NaN/Inf → `—`                             |
| `format.StringOrDash(s)`    | Empty → `—` placeholder for tables/detail views                       |
| `format.CompactCount(n)`    | `999`, `1.2k`, `3.4M`, `7.8T` — sign preserved                        |

#### `utils/wire` — transport-agnostic wiring (ADR-0036)

| Function                                                                                                                               | One-liner                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| -------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Action{Transport, Method, URL, Event, Target, ContentType, Selector, Swap, DebounceMS, ThrottleMS, PreventDefault, Interval, Reveal}` | One typed wiring spec; `.Attributes()` renders the htmx or Datastar dialect; `ContentTypeForm` serializes the enclosing form's fields under Datastar (`{contentType: 'form'}`, verified v1.0.3 runtime), ignored by htmx; `Swap` (a `PatchMode`) renders `hx-swap` under htmx and NOTHING under Datastar (the mode is response-header driven via `wire.Handler/PatchTarget` — corrected 2026-10-02, browser-proven); `Selector` renders only under ContentTypeForm (form lookup; never targets patches) |
| `Handler(target PatchTarget, next)`                                                                                                    | Both-transports endpoint middleware: Datastar callers get response-header targeting                                                                                                                                                                                                                                                                                                                                                                                                                     |
| `IsDatastar(r)` / `IsHTMX(r)`                                                                                                          | Request predicates on `Datastar-Request` / `HX-Request` headers                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `PatchTarget{Selector, Mode}`                                                                                                          | Server-side patch region for Datastar callers (`PatchMode` typed enum, 8 modes)                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `HeaderDatastarRequest` etc.                                                                                                           | Typed constants for the request/response marker headers                                                                                                                                                                                                                                                                                                                                                                                                                                                 |

Scope: common subset only (ADR-0036) — OOB/confirm/indicators/SSE stay in the `htmx`/`datastar` modules. `Action.Target` renders only for htmx (Datastar targeting is response-DRIVEN ONLY — `wire.Handler(PatchTarget)`; client fetch options never reach patch targeting, browser-proven; `Action.Selector` is form-lookup only). Whole-form submission IS symmetric (`forms.FormProps.Wire` defaults to `ContentTypeForm`). Per-field value binding stays asymmetric under Datastar (bound signals escape hatch). Components take wiring via `BaseProps.Attrs` or a typed `Wire` field (`display.Button`, `navigation.LoadMore`, `forms.Form`). Components take wiring via `BaseProps.Attrs` or a typed `Wire` field (`display.Button`, `navigation.LoadMore`, `forms.Form`, `display.Tabs` — `{tab}` URL placeholder). Confirm dialogs are htmx-ONLY (the pinned runtime ships zero confirm machinery, bundle-verified); per-patch View Transitions under Datastar ride `PatchTarget.UseViewTransitions` (the `Datastar-Use-View-Transition` response header). Constructors/builders: `wire.Get/Post/Put/Patch/Delete(url)` + value-receiver `With*` chainers + `WithFormDefaults(event)` — never mutate a consumer's `*Action`. All Datastar fetch options merge into ONE object literal. Typed triggers (ADR-0043): `Action.Interval` (polling; m/h normalized to seconds — both runtimes misparse them as ms) and `Action.Reveal *wire.Reveal` (viewport reveal; zero value fires once; `Exit` is Datastar-only) render in BOTH dialects — htmx merges all trigger sources into one comma-separated `hx-trigger`, Datastar one attribute per source; a trigger-only Action never also fires on click.
