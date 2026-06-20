# AGENTS.md — templ-components



## Build & Test Commands

```bash
# Full build (required before go build after .templ changes)
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./...

# Tests
go test ./...

# Lint
golangci-lint run ./...

# All-in-one verification
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./... && go test ./... && golangci-lint run ./...
```

## CRITICAL: Generated `*_templ.go` Files MUST Be Committed

This is a **templ library**, not an application. The Go module proxy (proxy.golang.org) fetches
source from the Git tag — it does **not** run `templ generate`. Without committed `*_templ.go`
files, consumers get uncompilable code (`undefined` errors on every component function).

- The `.gitignore` uses `!*_templ.go` to override the global gitignore's `*_templ.go` entry
- After editing any `.templ` file, always run `templ generate ./...` and commit the updated `*_templ.go` files alongside the source
- Never add `*_templ.go` back to `.gitignore` — this is the standard pattern for publishable templ packages
- 46 generated files across 10 packages + examples/demo (display, errorpage, feedback, forms, htmx, icons, internal/golden, internal/svg, layout, navigation)

**Why this matters:** The Go module proxy serves source as-is. Consumers who `go get` this package
will have their Go toolchain download the tagged commit. If `*_templ.go` is missing from that
commit, the package won't compile. Unlike applications (where you generate at build time), a
**library's generated code is part of its distributable artifact**.

## Architecture

- **Module:** `github.com/larsartmann/templ-components`
- **Go:** 1.26, **templ:** v0.3.x
- **No framework deps** — pure Go + templ + Tailwind v4 class strings
- **Theming:** Components emit standard Tailwind classes (`bg-blue-600`). Consumers override via `@theme { --color-blue-600: #custom; }` in their CSS. No Go code changes needed. See `templ-components-theme.css` for semantic alias examples.
- **ComponentProps interface:** `utils.ComponentProps` with `GetBaseProps()`/`SetBaseProps()` on `*BaseProps` (pointer receivers for `recvcheck`). All 25+ props structs auto-satisfy via method promotion.
- **Accessibility — motion-reduce:** `motion-reduce:transition-none motion-reduce:duration-0` on all transitions, `motion-reduce:animate-none` on all animations (spinner, skeletons, toast enter/exit, modal, accordion)
- **Dark mode colors:** All components use `gray-*` exclusively (no mixed `slate-*`/`gray-*`). Dark mode via class strategy: `@custom-variant dark (&:where(.dark, .dark *))` toggled by `layout.ThemeScript()` + `layout.ThemeToggle()`
- **CI:** `.github/workflows/ci.yaml` — lint (golangci-lint), build+test with `templ generate`, coverage artifact. Pre-commit: `.git/hooks/pre-commit` → `scripts/pre-commit.sh`
- **Import graph:** `utils ← all`, `internal/svg ← display,feedback,icons`, `icons ← display,feedback,errorpage`, `feedback ← none (htmx decoupled)`, `errorpage ← icons,utils,go-error-family (no feedback dep)`
- **No circular imports** allowed
- **AriaLabel propagation:** All components with `BaseProps` propagate `AriaLabel` to root element. Components with hardcoded aria-labels (Nav, Pagination, Breadcrumbs, StepIndicator) allow AriaLabel override via `utils.Ternary`
- **SVG paths:** Shared constants in `internal/svg` (PathChevronDown, PathChevronSmall, PathArrowUp/Down/Left/Right, PathAvatarFill) — single source of truth


## Code Conventions

- All component props embed `utils.BaseProps` (exception: `layout.PageProps`) — all auto-satisfy `utils.ComponentProps` interface
- All root elements propagate `props.Class`, `props.Attrs`, `props.ID`, and `props.AriaLabel` from BaseProps (25/25 components, including NavLink/MobileNavLink)
- Class attributes use `utils.Class()` for Tailwind conflict resolution (exception: `templ.KV` conditionals where comma-join is required)
- Style lookups use maps/structs, not switches (e.g., `badgeStyleMap`, `badgeSizeLookup`, `cardPaddingLookup`, `iconPathData`, `alertIconMap`, `toastIconMap`, `spinnerSizeLookup`, `progressHeightLookup`, `avatarSizeLookup`, `avatarDotSizeLookup`)
- String enums: `type XxxType string` + `const XxxDefault XxxType = "default"`
- Size constants: uppercase suffix pattern `[Component]Size[SM|MD|LG]` (e.g., `AvatarSizeSM`, `BadgeSizeSM`, `SpinnerSM`)
- Default constructors: `DefaultXxxProps()` for every component with non-zero defaults
- Private helpers: `xxxClass()` for Tailwind class mapping
- CSP: all inline scripts use `nonce={ props.Nonce }`
- Sub-templates: extract shared rendering to private `templ` functions
- Feedback styles: shared `feedbackStyleSet` struct + `lookupFeedbackStyle[T]()` generic + `feedbackIconName()` + `dismissScript()` in `feedback/styles.go`
- FeedbackType: canonical `FeedbackType` enum (`FeedbackSuccess/Error/Warning/Info`); `AlertType` and `ToastType` are type aliases for backward compat
- SVG paths: constants in `internal/svg` (`PathChevronDown`, `PathChevronSmall`, `PathArrowUp/Down/Left/Right`, `PathAvatarFill`) — single source of truth for inline SVG paths
- Icons: `iconPathData` map with `|` separator for multi-path icons. `iconPaths()` validates no empty segments (panics on stray `|`). `allIconNames()` auto-generated from `iconPathData` + Spinner — no manual list to maintain.
- Form errors: `ErrorAttrs(id, errMsg, helpTextID)` helper returns `templ.Attributes` for aria-invalid/aria-describedby
- Card shell CSS: shared `cardShellClass` constant for consistent card styling. `SimpleCard` composes through `Card` internally.
- HTMX loading: accepts `templ.Component` spinner parameter (decoupled from feedback package)
- Toast icons: generated from Go `iconPathData` via `icons.IconPathJS()` (single source of truth)
- TrendDirection: `TrendNone = "none"` (non-empty sentinel, not "")
- Layout: `Minimal(MinimalProps)` uses props struct like `Base(PageProps)`
- Interactive: `cursor-pointer` on buttons, `caret-blue-600 dark:caret-blue-400` on inputs, `scroll-smooth` + selection colors on body, `shadow-sm` on card shell
- Modal: focus save/restore via `data-tc-prev-focus` attribute on open, restored on close
- NavLink/MobileNavLink: `NavLink` uses `utils.Class()` for merge; `MobileNavLink` appends `props.Class` to `templ.KV` chain
- InputType: validates via `inputType()` with `validInputTypes` map; panics on unknown, defaults empty to `"text"`
- Structural variants (TabsVariant, DropdownPosition, TrendDirection): use `if`-branch for DOM structure, not map lookup — map pattern is for pure class lookups only
- `forms.SanitizeID`: exported utility for library consumers, not used internally
- Enum validation: 0 panic-on-unknown, 12 map+fallback, structural variants use if-branch. InputType falls back to "text", icons.Name falls back to Question icon. Only remaining panic: icon path data integrity check (stray `|` separators).
- ID auto-generation: `utils.EnsureID(prefix, id)` generates unique IDs via crypto/rand when consumer omits props.ID. Used by Modal, Drawer, Dropdown, Accordion, Combobox.
- SwapOOB: invalid swap styles fall back to `outerHTML` instead of panicking.
- Zero runtime panics in component code (only 1 developer data integrity check in icons package).
- Modal/Dropdown/Accordion/Combobox: IDs auto-generated via `utils.EnsureID()` when consumer omits props.ID — no panics.
- Combobox JS: global singleton `tcComboboxAttached` handler for input filtering, click-to-select, focus/blur dropdown management, Escape key dismissal. CSP-safe with `nonce={ props.Nonce }`.

- Toast JS: dismiss icon from `icons.IconPathJS()` via `tcToastIcons.dismiss`
- Table: row cells auto-padded/truncated to match header count
- **Error handler:** `errorpage/handler.go` provides `ErrorHandler(err, cfg)` returning `http.Handler`, `FromError(err)` for type-safe conversion from go-error-family errors, 6 pre-built constructors (`NotFound`, `Forbidden`, `BadRequest`, `Conflict`, `ServiceUnavailable`, `InternalError`), `WriteError`/`WriteErrorPage` convenience wrappers, `HTMLShell` mode for valid HTML documents, `JSON` mode for API/HTMX responses. Uses `errors.AsType[errorfamily.Classified]()` for go-error-family integration.
- **Error families:** `errorpage` package integrates with go-error-family via `FamilyFromErrorFamily()` converter + `ParseFamily()` for string-based lookup. `FromError()` extracts Why/Fix defaults from go-error-family's `Family.DefaultWhy()`/`DefaultFix()` methods.
- **Error components:** `ErrorPage` (full-page), `ErrorDetail` (inline card), `ErrorAlert` (family-aware alert) in `errorpage/`
- **Drawer:** `display.Drawer` — accessible side panel with left/right slide, focus trap, Escape key, backdrop click. Follows same pattern as Modal but with translate-x transforms.
- **ValidationSummary:** `forms.ValidationSummary` — accessible error summary with icon, error count, linked field errors, `role="alert"`.
- **Golden testing:** `internal/golden.Assert(t, name, got)` — golden file comparison with CSS class normalization. Supports `-update` flag.
- **Error sub-templates:** 6 shared private sub-templates in `errorpage/shared.templ` (familyIcon, fixCard, causeList, contextTable, timestampFooter, familyBadge)
- HTMX retry: per-element `data-tc-retry` attribute (no shared counter)
- HTMX error handling: family-aware — when server returns structured JSON with `family` field, toast type is mapped. `ErrorHandlerConfig{JSON: true}` produces the JSON format that HTMX consumes.
- GlobalErrorHandling: configurable via `ErrorHandlingConfig` struct (MaxErrorHistory, MaxRetries, RetryDelayMS). Includes `tc-error-announcer` div with `aria-live="polite"` for screen reader announcements.
- Pagination: `rel="prev"`/`rel="next"` on arrow links for SEO. Ellipsis rendering when visible range is truncated. Uses `net/url` for URL construction.
- Breadcrumbs: optional `Separator` field for custom separators. `JSONLD` field enables JSON-LD structured data (`application/ld+json`).
- Theme colors: `DefaultThemeColor` and `DefaultDarkThemeColor` constants in layout package.
- Icon stroke-width: `IconWithStrokeWidth(name, class, strokeWidth)` for custom stroke widths (default Icon uses 1.5).
- Select validation: `normalizeSelectOptions()` resolves Disabled+Selected contradiction (clears Selected).
- Badge href: `BadgeProps.Href` renders `<a>` instead of `<span>` when set, enabling badge-based navigation.
- ProgressBar indeterminate: `ProgressBarProps.Indeterminate bool` renders animated bar with `aria-busy="true"` instead of percentage-based width.
- StepIndicator orientation: `StepIndicatorProps.Orientation` with `StepHorizontal`/`StepVertical` constants for vertical progress tracking.
- Tabs client-side: `TabsProps.ClientSide bool` adds `data-tc-tabs` attribute and inline JS for click-to-activate and keyboard nav (ArrowLeft/Right, Home, End). Uses global singleton guard (`tcTabsAttached`).
- Form component: `forms.Form(FormProps)` with `Action`, `Method` (GET/POST), `CSRFToken` hidden input, and `Content` for composing form fields.
- Icons: 99 total (98 path icons + Spinner).
- Thread safety: `utils.Class()` uses `sync.Mutex` to protect tailwind-merge-go's shared LRU cache from concurrent access. Required even though the LRU has internal mutexes — they don't protect the full Merge() call sequence.
- DropdownItemKind: typed enum (`DropdownItemLink`, `DropdownItemButton`) with backward compat via `IsLink()` fallback to Href-based discrimination.

## Lint Command

```bash
# Must lint specific packages — examples/ excluded via .golangci.yml
golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...
```
