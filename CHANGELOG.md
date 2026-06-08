# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [0.2.0] — 2026-06-08

### Added

- `errorpage` package: 3 components for presenting structured errors on the web
  - `ErrorPage`: full-page error view with Wix-style What/Why/Fix/WayOut layout
  - `ErrorDetail`: inline error detail card with context table, cause chain, and suggested fix
  - `ErrorAlert`: family-aware alert banner with dismiss support
  - 5 error families (Rejection, Conflict, Transient, Corruption, Infrastructure) with distinct color/icon/tone
  - `FamilyStatusCode()`: maps Family → HTTP status code (400/409/503/500/503)
  - `ContextMap()`: converts map[string]string → []ContextPair
  - `ExtractCauseChain()`: walks Unwrap() chain to build CauseItem slice
  - `ParseFamily(string) Family`: case-insensitive string→Family conversion
  - `FamilyFromErrorFamily(errorfamily.Family) Family`: converts go-error-family int enum to errorpage string
  - `FamilyIsValid(Family) bool`, `FamilyIcon(Family) icons.Name`: validation and icon lookup
- `utils.DismissScript()`: shared dismiss JS extracted from feedback package (single source of truth)
- DismissScript call pattern unified: both feedback and errorpage call `utils.DismissScript()` directly (removed `feedback.dismissScript()` private wrapper)
- `errorpage/handler.go`: `http.Handler` integration for serving error pages
  - `ErrorHandler(err, cfg)`: returns `http.Handler` with correct HTTP status, Content-Type, and family-aware rendering
  - `FromError(err)`: type-safe conversion — uses `errors.AsType[errorfamily.Classified]()` for go-error-family, falls back to string-based interface, extracts Why/Fix from `Family.DefaultWhy()`/`DefaultFix()`
  - 6 pre-built constructors: `NotFound()`, `Forbidden()`, `BadRequest(msg)`, `Conflict(msg)`, `ServiceUnavailable()`, `InternalError()` with code constants
  - `WriteError()` and `WriteErrorPage()` convenience wrappers for `http.ResponseWriter`
  - `ErrorHandlerConfig.Override` callback for per-error customization
  - `ErrorHandlerConfig.HTMLShell`: wraps error page in valid HTML document (DOCTYPE/html/head/title/body)
  - `ErrorHandlerConfig.JSON`: renders JSON error response (family/code/message/title/why/fix) for API/HTMX endpoints
  - Render errors logged via `slog.Error` instead of silently discarded
- `errorpage/shared.templ`: 6 shared sub-templates extracted (familyIcon, fixCard, causeList, contextTable, timestampFooter, familyBadge) — eliminated 9 duplicated HTML patterns
- HTMX `GlobalErrorHandling`: family-aware error parsing — structured JSON responses with `family` field now map to appropriate toast types instead of generic status-code logic
- HTMX `ErrorHandlingConfig`: configurable error handling — `MaxErrorHistory`, `MaxRetries`, `RetryDelayMS` with `DefaultErrorHandlingConfig()`. Includes `tc-error-announcer` div with `aria-live="polite"` for screen reader announcements
- `icons.IconWithStrokeWidth(name, class, strokeWidth)`: custom stroke-width variant of Icon (default uses 1.5)
- `icons.allIconNames()`: auto-generated from `iconPathData` map — no manual list to maintain
- `icons.iconPaths()`: validates no empty segments from stray `|` separators (panics on malformed data)
- `navigation.Pagination`: `rel="prev"`/`rel="next"` on arrow links for SEO, ellipsis rendering when visible range is truncated
- `navigation.Breadcrumbs`: `Separator` field for custom separators, `JSONLD` field enables JSON-LD structured data
- `display.DropdownItemKind`: typed enum (`DropdownItemLink`, `DropdownItemButton`) with backward compat via `IsLink()` fallback
- `layout.DefaultThemeColor` / `layout.DefaultDarkThemeColor`: named constants replacing magic hex values
- `forms.normalizeSelectOptions()`: resolves Disabled+Selected contradiction (clears Selected when both are true)
- `display.SimpleCard`: composes through `Card` internally instead of duplicating shell CSS

### Changed

- **BREAKING**: `ConflictError(msg)` renamed to `Conflict(msg)` for naming consistency with other constructors
- **BREAKING**: `GlobalErrorHandling(nonce string)` → `GlobalErrorHandling(cfg ErrorHandlingConfig)` — configurable error handling with struct
- **BREAKING**: `DropdownItem` now has `Kind DropdownItemKind` field; backward compat via `IsLink()` fallback to Href discrimination
- **BREAKING**: `FromError()` now uses `errors.AsType[errorfamily.Classified]()` — requires `github.com/larsartmann/go-error-family` v0.2.0
- Added `github.com/larsartmann/go-error-family` v0.2.0 as dependency for type-safe error family bridging
- Render errors in `ErrorHandler`/`WriteErrorPage` now logged via `slog.Error` instead of silently discarded
- DismissScript call pattern unified: removed `feedback.dismissScript()` wrapper, all callers use `utils.DismissScript()` directly
- **BREAKING**: `Tab.Active bool` removed from `Tab` struct → `TabsProps.ActiveTabID string` on parent. Prevents zero/multiple active tabs
- Test deduplication: eliminated all 19 clone groups across 7 packages using extracted helpers, table-driven tests, and merged duplicates
- Coverage improvements: display 71.8%→72.7%, forms 70.8%→72.0%, navigation 72.2%→73.2%
- Added comprehensive edge case tests for error boundaries (nil/empty inputs, invalid enum fallbacks)
- Added benchmarks for hot render paths: Class merge, Badge, Card, Table, Modal, Dropdown
- Standardized error messages in `validateDropdownID()` and `validateModalID()` to use `fmt.Errorf()`
- Fixed 5 pre-existing goconst lint warnings in `forms/bdd_test.go` by extracting test constants
- Removed stale `MergeAttrs`, `Deref`, `DerefOr`, `BoolString` references from FEATURES.md (removed in v0.2)
- **BREAKING**: `BadgeDefault` constant removed → use `BadgeNeutral`. `DefaultBadgeProps()` now returns `BadgeNeutral`
- **BREAKING**: `ErrorAttrs(id, errMsg)` → `ErrorAttrs(id, errMsg, helpTextID)` — now links both error and help text IDs in `aria-describedby`
- **BREAKING**: `Minimal(title, locale string)` → `Minimal(MinimalProps)` for consistency with `Base`
- **BREAKING**: `LoadingIndicator()` → `LoadingIndicator(spinner templ.Component)` — decoupled from feedback package
- **BREAKING**: `InlineLoadingOverlay(id)` → `InlineLoadingOverlay(id, spinner templ.Component)`
- **BREAKING**: `LoadingButton(default, loading)` → `LoadingButton(default, loading, spinner templ.Component)`
- Badge color/dot maps consolidated into single `badgeStyleMap` with `badgeStyle` struct
- Tooltip position functions consolidated into `tooltipPositionMap` with `tooltipPositionStyles` struct
- Card shell CSS (`bg-white dark:bg-slate-800 border...`) extracted to `cardShellClass` constant
- HTMX CDN URL construction extracted to `htmxCDNURL()` helper
- Error handling JS magic numbers extracted to named constants (`MAX_ERROR_HISTORY`, `MAX_RETRIES`, `RETRY_DELAY_MS`)
- Toast icon paths now generated from Go `iconPathData` via `icons.IconPathJS()` — fixes copy-paste bug where error and warning had identical paths
- Avatar status dot now scales with avatar size (XS→1.5, SM→2, MD→2.5, LG→3, XL→3.5)
- `Exclamation` icon constant deprecated — use `ExclamationCircle` instead
- `icons.IconAttrs()` removed (was dead code — never called outside tests)
- ProgressBar a11y test moved from display to feedback package
- `TestIconCount` now dynamically checks `allIconNames` count matches `iconPathData` (+1 Spinner)

### Fixed

- NavLinkProps `Attrs` field shadowing `BaseProps.Attrs` — consumer attrs were silently dropped
- Dropdown JS XSS vulnerability — raw `props.ID` interpolated into JS. Now uses `strconv.Quote()`
- Accordion state coupling — `hidden` attribute prevented JS toggle from working on server-closed items. Now uses `data-open` attribute
- Modal/Dropdown empty ID produces broken ARIA attributes — now panics with clear error message
- Dropdown `sanitizeJSIdent` and `dropdownInitScript` unused functions removed
- Toast JS `error` and `warning` had identical SVG path data (copy-paste bug)

### Added

- `validateDropdownID()` and `validateModalID()` for required ID validation at render time
- Pre-commit hook replaced with project's own script
- `.golangci.yml` excludes examples from lint
- `icons.IconPathJS()` exported helper for JS icon path generation
- `toastJSIconPaths()` generates toast icon map from Go icon data (single source of truth)
- `htmxCDNURL()` helper for HTMX CDN URL construction
- `MinimalProps` struct and `DefaultMinimalProps()` for minimal layout
- ADR 0001: Two Icon Systems documentation
- `ErrorAttrs` now supports `helpTextID` parameter for dual `aria-describedby` references
- `avatarDotSizeClass()` for proportional status dot sizing

## [0.1.0] - 2026-01-01

### Added

- Initial release with 56 components, 44 types, 42 icons
- Display: Card, Badge, Modal, Table, Tabs, Avatar, Tooltip, Accordion, Dropdown, Empty State
- Feedback: Alert, Toast, Spinner, Progress Bar, Step Indicator, Skeleton, Loading
- Forms: Input, Select, Textarea, Checkbox, Label
- Navigation: Nav, Breadcrumbs, Pagination, Mobile Menu
- HTMX: Loading indicators, error handling, CSRF, OOB swap, confirm delete
- Layout: Base HTML, theme toggle, dark mode support
- Icons: 42 named SVG icons
