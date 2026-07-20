# AGENTS.md — templ-components

## Module Structure (single module)

This repo is a **single Go module** (`github.com/larsartmann/templ-components`) with 14 packages:

| Package             | Contains                                   | Purpose                                                                                                             |
| ------------------- | ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------- |
| `display`           | 30 UI components                           | Cards, tables (Table + DataTable), modals, badges, buttons, avatars, carousel, context menu, hover card             |
| `feedback`          | 13 components                              | Alerts, toasts, spinners, skeletons, progress bars                                                                  |
| `forms`             | 21 components                              | Inputs, selects, toggles, combobox, slider, rating, tags input, calendar, validation                                |
| `layout`            | 10 components                              | Page shell, theme toggle, CSP-safe script/style tags, **body-layout primitives**: AppShell, Container, Split, Stack |
| `navigation`        | 12 components                              | Nav bars, pagination, breadcrumbs, sidebar, EndOfList                                                               |
| `htmx`              | 8 components                               | HTMX loading, error handling, OOB swaps, View Transitions                                                           |
| `icons`             | 102 named SVG icons                        | Heroicons v2 outline + Spinner                                                                                      |
| `errorpage`         | 4 components + handler                     | Error pages, 404, go-error-family integration                                                                       |
| `utils`             | BaseProps, Class(), EnsureID, test helpers | Shared utilities                                                                                                    |
| `internal/svg`      | SVG path constants                         | Single source of truth for inline SVG paths                                                                         |
| `internal/golden`   | Golden file testing                        | CSS-normalized HTML snapshot comparison                                                                             |
| `internal/contract` | Contract tests                             | Cross-package interface verification                                                                                |
| `integration`       | CSP nonce tests                            | Asserts nonce on all inline scripts                                                                                 |
| `examples/demo`     | Demo binary                                | Showcases components                                                                                                |

> **Note:** A multi-module workspace split was prototyped on the `modularize/strategic-split`
> branch but was never merged to `master`. The split may be re-attempted post-v1.0 if the
> package graph warrants it.

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

### Nix flake commands

```bash
# Format all .nix (nixfmt) + .go (gofmt + goimports) files.
# Generated *_templ.go, website/, and demo static/ are excluded.
nix fmt

# Run all flake checks (currently: treefmt format verification).
nix flake check

# Flake apps (export GOEXPERIMENT=jsonv2 automatically)
nix run .#build     # templ generate + go build
nix run .#test      # go test -race
nix run .#lint      # golangci-lint
nix run .#verify    # generate + build + test + lint
nix run .#coverage  # go test -coverprofile
```

The flake uses `flake-parts` + `treefmt-nix` (mirrors `website/flake.nix`). The
`formatter` output is provided by treefmt-nix's flakeModule (replaces the former
bare `formatter = pkgs.nixfmt;`). BuildFlow still owns the pre-commit hook.

## CRITICAL: Generated `*_templ.go` Files MUST Be Committed

This is a **templ library**, not an application. The Go module proxy (proxy.golang.org) fetches
source from the Git tag — it does **not** run `templ generate`. Without committed `*_templ.go`
files, consumers get uncompilable code (`undefined` errors on every component function).

- The `.gitignore` uses `!*_templ.go` to override the global gitignore's `*_templ.go` entry
- After editing any `.templ` file, always run `templ generate ./...` and commit the updated `*_templ.go` files alongside the source
- Never add `*_templ.go` back to `.gitignore` — this is the standard pattern for publishable templ packages
- 87 generated files across all packages
- **BuildFlow gotcha:** the BuildFlow pre-commit `templ-generate` step re-appends `*_templ.go` to `.gitignore` on every run, which (being the last pattern) overrides the `!*_templ.go` unignore and hides generated files from `git status`. This is harmless for already-tracked files (gitignore cannot untrack), but any NEW component's `*_templ.go` will be invisible until `git add -f`. After each commit, check `git status` for a re-added `*_templ.go` line and remove it. Consider fixing this in BuildFlow itself (it is `larsartmann/buildflow`).

**Why this matters:** The Go module proxy serves source as-is. Consumers who `go get` this package
will have their Go toolchain download the tagged commit. If `*_templ.go` is missing from that
commit, the package won't compile. Unlike applications (where you generate at build time), a
**library's generated code is part of its distributable artifact**.

## templ Version Pin: go.mod v0.3.1020, system binary may be v0.3.1036

`go.mod` pins `github.com/a-h/templ v0.3.1020` — the latest **published** version on the Go module
proxy (https://proxy.golang.org/github.com/a-h/templ/@v/list). The system `templ` binary in
`~/.nix-profile/bin/templ` may be a local Nix build of unreleased upstream master
(`github:a-h/templ` flake), reporting `v0.3.1036`. This causes a cosmetic import-block style diff
across all 51 `*_templ.go` files on every regen:

- v0.3.1020 emits `import "github.com/a-h/templ"` on its own line, then a separate
  `import (...)` block for project imports
- v0.3.1036 collapses both into a single `import (...)` block

**Rule:** always use `nix develop` to enter the dev shell before running `templ generate`. The dev
shell provides `pkgs.templ` (v0.3.1020) which matches `go.mod` and produces zero diff. If you run
`templ generate` with the system binary, expect 51 files to change cosmetically — these are no-op
import-style changes; the generated code is semantically identical.

**Do not bump `go.mod` to v0.3.1036** — that version is not yet on the module proxy, so consumers
who `go get` this package would fail. Wait for the official upstream release, then bump in lockstep.

## Architecture

- **Module:** `github.com/larsartmann/templ-components`
- **Go:** 1.26, **templ:** v0.3.x
- **No framework deps** — pure Go + templ + Tailwind v4 class strings
- **CSS standard:** Tailwind CSS v4+ (latest) for ALL LarsArtmann projects. CSS-first config, no Node.js runtime, no DaisyUI. Small custom CSS only where Tailwind doesn't cover something. See `docs/adr-001-tailwind-v4-standard.md` and `docs/tailwind-v4-adoption-guide.md`.
- **CSS setup:** Consumers vendor the library and copy `templates/app.css` + `templates/custom.css` as a starter entry point, then compile with `tailwindcss`. `app.css` imports `custom.css` via `@import "./custom.css"`. BuildFlow's `tailwind-build` provider automates this in its DAG. See `docs/tailwind-v4-adoption-guide.md` for details.
- **JavaScript patterns:** see `docs/javascript-guide.md` for the complete decision ladder (native HTML → HTMX → singleton-guard → Alpine → Datastar → islands), CSP compliance, and templ's built-in JS features. See ADR 0005 for the singleton-guard pattern used by all interactive components in this repo.
- **Accepted code duplication:** see `docs/adr/0009-accepted-clones.md` for the 4 clone groups that art-dupl flags but are intentional (idiomatic UI layout, required by templ DSL, or demo content). New dedup passes should not force extraction beyond what's documented in that ADR.
- **Theming:** Components emit standard Tailwind classes (`bg-blue-600`). Consumers override via `@theme { --color-blue-600: #custom; }` in their CSS. No Go code changes needed. See `templ-components-theme.css` for semantic alias examples.
- **ComponentProps interface:** `utils.ComponentProps` with `GetBaseProps()`/`SetBaseProps()` on `*BaseProps` (pointer receivers for `recvcheck`). All 26+ props structs auto-satisfy via method promotion.
- **Accessibility — motion-reduce:** `motion-reduce:transition-none motion-reduce:duration-0` on all transitions, `motion-reduce:animate-none` on all animations (spinner, skeletons, toast enter/exit, modal, accordion)
- **Dark mode colors:** All components use `gray-*` exclusively (no mixed `slate-*`/`gray-*`). Dark mode via class strategy: `@custom-variant dark (&:where(.dark, .dark *))` toggled by `layout.ThemeScript()` + `layout.ThemeToggle()`. `color-scheme: light` on `:root`, `color-scheme: dark` on `.dark` (native form control rendering).
- **Dark mode color convention:** Light mode uses `-600` shade for backgrounds (`bg-blue-600`), dark mode uses `-500` (`dark:bg-blue-500`). Light mode uses `-600` for text (`text-blue-600`), dark mode uses `-400` (`dark:text-blue-400`). Neutral text: `text-gray-500` → `dark:text-gray-400`, `text-gray-400` → `dark:text-gray-500`. Every neutral and semantic color class MUST have a `dark:` variant — enforced by `utils.TestDarkModeCompliance` + `utils.TestDarkModeSemanticColors` (failing tests, block CI). Exceptions: Toggle thumb (`bg-white` both modes), SidebarNav (permanently dark sidebar), avatar silhouette icon (`text-blue-200` decorative).
- **Dark mode compliance tests:** `utils.TestDarkModeCompliance` scans all `.templ`/`.go` source files for neutral colors (`text-gray-*`, `bg-white`, `bg-gray-*`, `border-gray-*`, `ring-gray-*`) without `dark:` variants. `utils.TestDarkModeSemanticColors` scans for semantic colors (`bg-blue-600`, `text-red-600`, etc.) without `dark:` variants. Both are FAILING tests — they block CI. Run via `go test ./utils/... -run TestDarkMode`. For the full dark mode strategy analysis (Tailwind v4 default is `prefers-color-scheme`, three consumer paths, `@theme` palette override pattern), see `docs/dark-mode-research.md`.
- **CI:** `.github/workflows/ci.yaml` — lint (golangci-lint), build+test with `templ generate`, coverage artifact. Pre-commit: `.git/hooks/pre-commit` → `scripts/pre-commit.sh`
- **Import graph:** `internal/svg` ← `icons`; `utils` (leaf); `icons,internal/svg,utils` ← root packages (display,feedback,forms,layout,navigation,htmx); `icons,utils` ← errorpage; all ← examples/demo. Production deps: `icons → internal/svg`, `display → icons,internal/svg,utils`, `feedback → icons,internal/svg,utils`, `forms → icons,utils`, `layout → icons,utils`, `navigation → icons,internal/svg,utils`, `htmx → utils`, `errorpage → icons,utils`
- **No circular imports** allowed
- **AriaLabel propagation:** All components with `BaseProps` propagate `AriaLabel` to root element. Components with hardcoded aria-labels (Nav, Pagination, Breadcrumbs, StepIndicator) allow AriaLabel override via `utils.Ternary`
- **SVG paths:** Shared constants in `internal/svg` (PathChevronDown, PathChevronSmall, PathArrowUp/Down/Left/Right, PathAvatarFill) — single source of truth

## Code Conventions

- All component props embed `utils.BaseProps` (exception: `layout.PageProps`) — all auto-satisfy `utils.ComponentProps` interface
- All root elements propagate `props.Class`, `props.Attrs`, `props.ID`, and `props.AriaLabel` from BaseProps (26/26 components, including NavLink/MobileNavLink)
- Class attributes use `utils.Class()` for Tailwind conflict resolution (exception: `templ.KV` conditionals where comma-join is required)
- **RTL/i18n: use logical CSS properties exclusively.** Never use `ml-`/`mr-`/`pl-`/`pr-`/`left-`/`right-`/`text-left`/`border-l-`/`border-r-` — use `ms-`/`me-`/`ps-`/`pe-`/`start-`/`end-`/`text-start`/`border-s-`/`border-e-` instead. These are CSS logical properties that automatically mirror in RTL (`dir="rtl"`). Exception: `left-1/2 -translate-x-1/2` for centering (not directional).
- **Motion: use shared transition constants.** Use `transitionFast` (150ms), `transitionNormal` (200ms), `transitionColors`, `transitionTransform` from `display/shared.go` instead of inline timing strings. All include `motion-reduce:*` fallbacks. Wire into CopyButton, Accordion, Modal, Drawer — do NOT leave inline `transition-colors motion-reduce:...` strings when a constant matches.
- **Container queries: use `@container` for context-responsive grids.** Set `GridProps.ContainerResponsive = true` to make the Grid respond to its parent container width instead of the viewport. Uses Tailwind v4 container-query variants (`@sm:`, `@lg:`).
- Style lookups use maps/structs, not switches (e.g., `badgeStyleMap`, `badgeSizeLookup`, `cardPaddingLookup`, `iconPathData`, `alertIconMap`, `toastIconMap`, `spinnerSizeLookup`, `progressHeightLookup`, `avatarSizeLookup`, `avatarDotSizeLookup`)
- **Lookup maps MUST use typed enum keys** (never `map[string]X`). If a typed enum exists, its lookup map uses it as the key type — `badgeSizeLookup[BadgeSizeMD]`, not `badgeSizeLookup[string(v)]`. All map lookups go through `utils.Lookup(m, key, fallback)` (generic, no per-call `if ok` boilerplate). `ButtonHTMLType` uses `map[ButtonHTMLType]string` + `utils.Lookup` (not `map[X]bool`).
- **Every closed-set enum MUST ship an `IsValid()` method + a test in the same commit.** 31 enums have IsValid (e.g., `SortDirectionIsValid`, `ButtonHTMLTypeIsValid`, `TableCellPaddingIsValid`). Test in the package's `enums_test.go` table-driven `TestIsValidEnums`. No IsValid without a test — this prevents the dead-code ghost system.
- **Drift-guard tests:** `utils.TestVersionMatchesChangelog` (CHANGELOG heading == `utils.Version`) and `utils.TestVersionMatchesFeatures` (FEATURES.md `**Version:**` == `utils.Version`). Bump version, CHANGELOG heading, and FEATURES.md version together at release time.
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
- Muted text: shared `mutedTextClass` constant (in `display/shared.go`) for the standard secondary-text pattern (`text-sm text-gray-500 dark:text-gray-400`); callers combine it with a margin via `utils.Class(mutedTextClass, "mt-N")`. Used by Card subtitle, PageHeader subtitle, EmptyState description.
- Pagination: `paginationPageItem`/`paginationEllipsisItem` sub-templates wrap `activeSpanOrLink`/`paginationEllipsis` in `<li>` so the page-range loop body stays flat.
- EmptyState: single `emptyStateAction(text, href, attrs)` helper renders an anchor when `href != ""`, else a button (replaces the former link/button pair).
- HTMX loading: accepts `templ.Component` spinner parameter (decoupled from feedback package)
- Toast icons: generated from Go `iconPathData` via `icons.IconPathJS()` (single source of truth)
- TrendDirection: `TrendUp`/`TrendDown`/`TrendWarn`/`TrendNone`. `TrendWarn` uses amber (`text-amber-600 dark:text-amber-400`) + right-pointing arrow + sr-only "Holding at". `TrendNone = "none"` (non-empty sentinel, not "")
- Layout: `Minimal(MinimalProps)` uses props struct like `Base(PageProps)`
- Layout: `PageProps.HTMXCDN` overrides the CDN base URL for htmx scripts. Empty defaults to `https://cdn.jsdelivr.net/npm`. Both `htmx.org` and `htmx-ext-response-targets` URLs derive from this value. Consumers with a different CSP can set e.g. `HTMXCDN: "https://unpkg.com"` without forking the library.
- Interactive: `cursor-pointer` on buttons, `caret-blue-600 dark:caret-blue-400` on inputs, `scroll-smooth` + selection colors on body, `shadow-sm` on card shell
- Modal/Drawer: native `<dialog>` element provides focus trapping, Escape-to-close, focus restore, top-layer rendering, and `::backdrop` — zero JS for those behaviors. CSS `@starting-style` + `allow-discrete` handle open/close animations (defined in `templates/custom.css` under `.tc-modal` / `.tc-drawer`). `tcOpenModal(id)` / `tcCloseModal(id)` are thin wrappers around `dialog.showModal()` / `dialog.close()` for backward compat. Backdrop click detection: `e.target === dialog` (the backdrop is a pseudo-element of dialog, so clicks register on the dialog itself).
- NavLink/MobileNavLink: both render through the shared `navLinkAnchor` sub-template; each supplies an active/inactive base-class builder (`navLinkClasses`, `mobileNavLinkClass`) and `navLinkAnchor` merges `props.Class` via `utils.Class()` so consumer Tailwind overrides resolve correctly. Do NOT assert ordered class substrings in tests — `utils.Class`/tailwind-merge reorders classes; use `utils.AssertContainsAll` for multi-token checks.
- InputType: validates via `inputType()` with `validInputTypes` map; panics on unknown, defaults empty to `"text"`
- Structural variants (TabsVariant, DropdownPosition, TrendDirection): use `if`-branch for DOM structure, not map lookup — map pattern is for pure class lookups only
- `forms.SanitizeID`: exported utility for library consumers; also used internally by `forms.RadioGroup` to derive per-option IDs from option values
- Enum validation: 0 panic-on-unknown, 15 map+fallback (InputType, ButtonHTMLType, FormMethod, TableCellPadding now included), structural variants use if-branch. InputType falls back to "text", ButtonHTMLType/FormMethod fall back to HTML-spec defaults ("button"/"GET"), icons.Name falls back to Question icon. Only remaining panic: icon path data integrity check (stray `|` separators).
- AvatarStatus: only `online`/`offline` render a colored status dot; unknown values render no dot (graceful degradation, no invisible element).
- ID auto-generation: `utils.EnsureID(prefix, id)` generates unique IDs via crypto/rand when consumer omits props.ID. Used by Modal, Drawer, Dropdown, Accordion, Combobox.
- SwapOOB: invalid swap styles fall back to `outerHTML` instead of panicking.
- Zero runtime panics in component code (only 1 developer data integrity check in icons package).
- Version: `utils.Version` is the single source of truth for the library release; a drift-guard test (`utils.TestVersionMatchesChangelog`) asserts it matches the latest CHANGELOG heading. Bump in lockstep with the Git tag.
- Modal/Dropdown/Accordion/Combobox: IDs auto-generated via `utils.EnsureID()` when consumer omits props.ID — no panics.
- Combobox JS: global singleton `tcComboboxAttached` handler for input filtering, click-to-select, focus/blur dropdown management, Escape key dismissal. CSP-safe with `nonce={ props.Nonce }`.
- CopyButton JS: global singleton `tcCopyAttached` handler — click delegation on `[data-tc-copy]`, clipboard write via `navigator.clipboard.writeText`, temporary label swap via `[data-tc-copy-label]` for 2s.
- Image fallback JS: global singleton `tcImageFallbackAttached` handler — error event capture (true) on `[data-tc-img-fallback]`, swaps src to fallback and removes attribute. Uses capture phase because error events don't bubble.
- CountBadge: zero count hides badge (aria-hidden decorative), overflow shows "N+" (default max 99). `formatInt` helper is shared with CountBadge.
- RelativeTime: pure Go formatting (`formatRelativeTime`), no JS. `<time datetime>` for a11y/SEO, `title` for absolute time on hover.
- LoadMore: cursor appended as `?cursor=` query param (detects existing `?` for `&`). `hx-swap="outerHTML"` + `hx-target="this"` for self-replacement. `InfiniteScroll: true` adds `hx-trigger="revealed"`.
- EndOfList: `navigation.EndOfList(EndOfListProps)` — "You've reached the end" indicator for the bottom of a list. Companion to LoadMore/Pagination. `role="status"`, `text-gray-500 dark:text-gray-400`. Customizable `Message`.
- StatCard HTMX: `HxGet`/`HxTarget`/`HxSwap` typed fields on both `<a>` and `<div>` variants. When empty, attributes are omitted.
- Card.Body: `Body templ.Component` slot — when set, overrides children. Backward compatible.
- Card.Header: `Header templ.Component` slot — when set, replaces the entire default header section (title, subtitle, header action). Use for custom header layouts. When nil, default header renders if Title or HeaderAction is set.
- Card.CardPaddingNone: when `Padding == CardPaddingNone`, children/body render without the wrapping padding `<div>` — directly inside the card shell. Enables table-in-card layouts where `<table>` must be a direct child for `overflow-x-auto`.
- Table.Flush: `TableProps.Flush bool` suppresses the wrapper div's border + rounded corners. Set when nesting a Table inside `Card(CardPaddingNone)` to avoid the double-border defect — the card provides the outer border. The `overflow-x-auto` scroll wrapper is always retained.
- Table.CellPadding: `TableProps.CellPadding` typed enum (`TableCellPaddingComfortable` px-4 py-3 default / `TableCellPaddingCompact` px-4 py-2). Controls vertical density of `<th>`/`<td>` cells via `tableCellPaddingLookup` map + `tableCellPaddingClass()` + `utils.Lookup` fallback. `TableCellPaddingIsValid` included.

- Toast JS: dismiss icon from `icons.IconPathJS()` via `tcToastIcons.dismiss`
- Table: row cells auto-padded/truncated to match header count. `TypedHeaders []TableHeader` takes precedence over `Headers []string` for sortable columns: each `TableHeader` has `Sortable bool`, `SortDirection` (`SortNone`/`SortAsc`/`SortDesc`), and `Href` for server-side sort links. Renders `aria-sort="ascending/descending/none"` + ↑/↓ indicators. `ariaSortValue()` maps the enum; `tableHeaderCount()` aligns cell padding to whichever header type is used. `TableRow.Href` makes rows clickable: sets `data-tc-row-href` + `role="link"` + `tabindex="0"` + `cursor-pointer`. A CSP-safe singleton script (`tableRowHrefJS` in `shared.go`) handles click + keyboard navigation. Clicks on interactive child elements (links, buttons) are not hijacked.
- **Error handler:** `errorpage/handler.go` provides `ErrorHandler(err, cfg)` returning `http.Handler`, `FromError(err)` for type-safe conversion from go-error-family errors, 6 pre-built constructors (`NotFound`, `Forbidden`, `BadRequest`, `Conflict`, `ServiceUnavailable`, `InternalError`), `WriteError`/`WriteErrorPage` convenience wrappers, `HTMLShell` mode for valid HTML documents, `JSON` mode for API/HTMX responses. Uses `errors.AsType[errorfamily.Classified]()` for go-error-family integration.
- **Error families:** `errorpage` package integrates with go-error-family via `FromErrorFamily()` converter + `ParseFamily()` for string-based lookup. `FromError()` extracts Why/Fix defaults from go-error-family's `Family.DefaultWhy()`/`DefaultFix()` methods. (`FamilyFromErrorFamily` is a deprecated alias, will be removed in v1.0.)
- **Error components:** `ErrorPage` (full-page), `NotFound404` (dedicated 404 with hero numeral + search + links), `ErrorDetail` (inline card), `ErrorAlert` (family-aware alert) in `errorpage/`
- **NotFound404:** `errorpage.NotFound404(props NotFound404Props)` — dedicated 404 page with large gradient numeral (`text-[8rem]`), optional search form (`SearchAction`), quick-links card grid (`[]NotFoundLink`), and "Go home" / "Go back" buttons. Unlike `ErrorPage` (family-colored error card), NotFound404 is a welcoming navigation aid using neutral blue/indigo palette. `DefaultNotFound404Props()` returns full defaults. `DefaultNotFoundLinks()` returns starter links (Home, Documentation). Types and constructors live in `notfound404_types.go`. All string constants (titles, messages, labels) defined as package-private `notFound404*` constants for goconst compliance. Shares the `tcGoBackAttached` singleton with `ErrorPage`.
- **Drawer:** `display.Drawer` — accessible side panel rendered as a native `<dialog>` with `data-side="left"`/`"right"`. CSS positions the dialog via `margin-inline-*` (auto-mirrors in RTL) and animates via `translateX`. Side positioning is in `templates/custom.css` under `dialog.tc-drawer[data-side=...]`.
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
- Select validation: `normalizeSelectOptions()` resolves Disabled+Selected contradiction (clears Selected). `SelectProps.Groups []SelectGroup` renders `<optgroup>` elements when set (Options is ignored). Each group's options go through the same normalization.
- Badge href: `BadgeProps.Href` renders `<a>` instead of `<span>` when set, enabling badge-based navigation.
- ProgressBar indeterminate: `ProgressBarProps.Indeterminate bool` renders animated bar with `aria-busy="true"` instead of percentage-based width.
- StepIndicator orientation: `StepIndicatorProps.Orientation` with `StepHorizontal`/`StepVertical` constants for vertical progress tracking.
- Tabs client-side: `TabsProps.ClientSide bool` adds `data-tc-tabs` attribute and inline JS for click-to-activate and keyboard nav (ArrowLeft/Right, Home, End). Uses global singleton guard (`tcTabsAttached`).
- Form component: `forms.Form(FormProps)` with `Action`, `Method` (GET/POST), `CSRFToken` hidden input, and `Content` for composing form fields. `Inline bool` switches to horizontal `flex flex-wrap items-end gap-3` layout (for filter bars / toolbars) instead of the default vertical `space-y-6` stack. Follows the `RadioGroup.Inline` precedent.
- StatCard HTMX: `HxGet`/`HxTarget`/`HxSwap` typed fields on both `<a>` and `<div>` variants. `HxSwap` is typed `htmx.SwapStyle` (not raw `string`) — consumers pass `htmx.SwapInnerHTML`/`htmx.SwapOuterHTML`. When empty, attributes are omitted.
- Icons: 102 total (101 path-icon Name consts + Spinner). Includes single-letter `X` and its discoverability alias `Close` (both `"x"`), plus 3 cross-API aliases (`Menu`/`Bars3`, `Refresh`/`ArrowPath`, `Location`/`MapPin`, `ThumbUp`/`HandThumbUp`) resolved via `iconAliases`. Added BuildingOffice2 (tenants/orgs), Key (credentials), ArrowRightOnRectangle (logout) — gaps surfaced by cqrs-htmx/adminui. `AllIconNames()` returns 97 (96 unique iconPathData keys + Spinner).
- Thread safety: `utils.Class()` uses `sync.Mutex` to protect tailwind-merge-go's shared LRU cache from concurrent access. Required even though the LRU has internal mutexes — they don't protect the full Merge() call sequence.
- DropdownItemKind: typed enum (`DropdownItemLink`, `DropdownItemButton`) with backward compat via `IsLink()` fallback to Href-based discrimination.
- PageHeader: `display.PageHeader(PageHeaderProps)` with Title, Subtitle, Breadcrumb slot, Action slot. No navigation import — uses `templ.Component` slots for compositional freedom.
- DefinitionList: `display.DefinitionList(DefinitionListProps)` with `[]DefinitionItem{Term, Detail, DetailComponent}`. Two-column `<dl>` grid for metadata/key-value display.
- ListNote: `display.ListNote(ListNoteProps{Shown, Total})` — renders "Showing N of M" truncation notice when Total > Shown. `role="status"` for a11y.
- SidebarNav: `navigation.SidebarNav(SidebarNavProps)` — vertical sidebar with Brand/Footer slots, icon+label nav items, CurrentPath auto-active detection, `aria-current="page"` on active item.
- icons.IconPathData: exported function returning raw SVG path d-strings for consumers needing full `<svg>` wrapper control (used by cqrs-htmx/adminui for icons-only adoption without Tailwind).
- icons-only adoption: the `icons` package depends only on `internal/svg` + `templ` + `utils` — no tailwind-merge-go, no CSS framework. This is a natural property of icons, not a portability strategy. See `docs/icons-only-adoption.md`.
- Grid: `display.Grid(GridProps)` — responsive grid container with typed `GridCols` enum (1–6, plus `GridColsAutoFit`) and `gridColsLookup` map+fallback. Replaces the repeated `grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4` pattern. Children passed via `{ children... }`. `GridColsAutoFit` + `MinColWidth` generates CSS `repeat(auto-fit, minmax(Xpx, 1fr))` for container-width-responsive grids (common dashboard pattern).
- StatCard.Href: `StatCardProps.Href` renders the whole card as an `<a>` with hover shadow, focus ring, cursor-pointer when set. Mirrors `Badge.Href` pattern. Shared body extracted to `statCardInner` sub-template so linked/unlinked layouts don't diverge.
- layout.Script: `layout.Script(nonce, src, attrs)` — CSP-safe `<script src>` tag that auto-injects the nonce. Prevents forgetting `nonce={...}` under strict CSP. Use instead of raw `<script>` tags for external scripts.
- SkeletonCardGrid: `feedback.SkeletonCardGrid(count)` — renders N skeleton cards in a responsive 3-col grid under a single `role="status"`. Pairs with `display.Grid` for the common dashboard loading pattern. Non-positive count falls back to 1.
- SimpleNav.RightItems: `SimpleNavProps.RightItems` forwards a `templ.Component` (e.g. ThemeToggle, sign-in button) to `Nav.RightItems`. Consumers no longer render right-side nav items in a separate flex container.
- PageProps auto-injects: `DefaultPageProps()` auto-injects `<link rel="stylesheet" href="/app.css">` (CSSPath) and `<script src="...htmx.org@...">` (HTMXVersion). Set either to `""` to suppress. Documented in godoc on both fields + DefaultPageProps.
- ProgressBar clamp: uses `max(0, min(100, props.Progress))` (Go 1.21+ builtin) instead of manual if-branch clamping.
- OverlayKind: typed enum (`OverlayModal`, `OverlayDrawer`) replaces untyped `closeKind`/`componentName` strings on `overlayShellProps`. The `componentName()` method derives JS function names from the kind.
- CopyButton: `execCommand('copy')` fallback for non-secure HTTP contexts. `role="status"` + `aria-live="polite"` on label span for screen reader feedback.
- Table.Body: `TableProps.Body templ.Component` slot — when set, overrides Rows for custom `<tr>` rendering. Follows Card.Body pattern.
- SimpleCard.Body: `SimpleCardProps.Body templ.Component` slot — same pattern, forwarded to Card.Body internally.
- layout.Stylesheet: `layout.Stylesheet(href, attrs)` — CSP-safe `<link rel="stylesheet">` companion to `layout.Script`.
- LoadMore URL: uses `net/url` for cursor encoding — base64 cursors with `=`/`+` are properly escaped. Physical-property `containsChar` helper deleted.
- RTL/i18n: all Tailwind classes use logical properties (`ms-`/`me-`/`start-0`/`end-0`/`ps-`/`pe-`/`text-start`/`border-s-`/`border-e-`). Components automatically mirror in RTL contexts when consumer sets `dir="rtl"`.
- CSP nonce test: `integration/csp_nonce_test.go` renders every inline-script component and asserts every `<script>` tag has `nonce=`. Prevents CSP regressions.
- GridProps.Gap: `GridGap` typed enum (`GridGapSM`/`MD`/`LG`/`XL` → `gap-2`/`4`/`6`/`8`) with `gridGapLookup` map + `GridGapIsValid`. Gap is a separate field from `Cols` — the lookup maps no longer include hardcoded `gap-4`. `DefinitionGrid` passes `gridGapClass(GridGapDefault)` explicitly.
- CopyButton.Href: `CopyButtonProps.Href` renders `<a>` instead of `<button>` when set. Link still copies to clipboard via the singleton script. Same `data-tc-copy` attribute on both variants.
- Image.Rounded: `ImageProps.Rounded bool` adds `rounded-full` when true (for avatars/icons), defaults to `rounded-md`. Convenience shortcut — consumers wanting other rounding use `props.Class`.
- LoadMore.InfiniteScroll: `LoadMoreProps.InfiniteScroll bool` adds `hx-trigger="revealed"` for auto-load on scroll. Defaults to false (click-to-load).
- NotFound404.LinksTitle: `NotFound404Props.LinksTitle string` — configurable heading for the quick-links section. Defaults to "Popular pages".
- WriteNotFound404: `errorpage.WriteNotFound404(w, r, props, nonce)` — convenience HTTP handler that renders `NotFound404` with a 404 status code. Mirrors `WriteErrorPage` pattern.
- ROADMAP.md and CONTRIBUTING.md: both created. ROADMAP tracks v0.x/v1.0/v2.0+ direction. CONTRIBUTING has Nix setup, conventions table, release flow.
- goBackScript promotion: stays in `errorpage/shared.templ` — only `errorpage` uses `history.back()` (ErrorPage + NotFound404). Promote to `utils` when a 2nd package needs it, following the `DismissButton` precedent (promoted to `utils/dismiss.templ` only when both `feedback` and `errorpage` shared the same markup).
- overlayShellProps field count: 8 fields, 2 callers (Modal + Drawer). Fields: `id`, `open`, `title`, `ariaLabel`, `kind`, `nonce`, `dialogClass`, `side`, `attrs`. Previous version had 11 fields with `outerClass`, `cfg`, `panelKVs` — eliminated by the `<dialog>` migration (CSS handles positioning, animations, and open/closed state).
- **Tailwind variant-prefix classes MUST be complete literals.** Never dynamically concatenate variant prefixes (`"peer-checked:" + translateClass`) — Tailwind's content scanner cannot find the resulting token and the CSS is never generated. Store complete literals (`peer-checked:translate-x-5`) in lookup maps. The `forms.Toggle` component follows this pattern.
- **`<dialog>` migration eliminated overlay JS complexity.** Modal and Drawer now use the native `<dialog>` element. Previously ~200 lines of JS per overlay instance handled focus trapping, Tab cycling, Escape, aria-hidden/inert toggling, focus save/restore, and animation class toggling. Now: `dialog.showModal()` / `dialog.close()` handles all of that natively. The JS (`overlayDialogJS` in `display/shared.go`) is ~15 lines: singleton guard + thin `tcOpen`/`tcClose` wrappers + per-instance IIFE for auto-open and click delegation. CSS `@starting-style` + `allow-discrete` replace the JS class toggling for animations.
- **`normalizeSelectOptions` returns a defensive copy.** Never mutate a caller's slice in a normalize function — Go slices share backing arrays, so in-place modification corrupts the caller's data across re-renders. Always `make` + `copy` first.
- **Combobox hidden input mirrors disabled state.** When `Disabled: true`, BOTH the visible text input and the hidden submission input must get `disabled` — otherwise the disabled field's value is still submitted (HTML spec violation).
- **Combobox `input` event clears hidden value.** When the user types without selecting an option, the hidden value must be cleared so a stale server-provided value isn't silently submitted.
- **Combobox Enter does NOT preventDefault unless an option is highlighted.** Unconditional `preventDefault()` on Enter blocks form submission when no option is active. Enter should fall through to let the form submit naturally.
- **Checkbox without ID renders `<span>` not `<label for="">`.** An empty `for=""` is invalid HTML and breaks label-input association. Guard like the Radio component.
- **Toast auto-generates ID via `EnsureID`** so `Duration > 0` always works (the auto-dismiss `setTimeout` references the toast by ID). `DefaultToastProps()` sets Duration: 5000 but previously omitted ID, silently disabling auto-dismiss.
- **ProgressBar clamps `aria-valuenow` to `[0, Total]`** matching the visual width clamp. Raw `props.Current` can be negative or exceed Total, violating the ARIA spec.
- **Dropdown RTL keys computed as variables** (`var nextKey = isRtl ? 'ArrowLeft' : 'ArrowRight'`), never inside string-literal comparisons. The Tabs component already does this correctly — Dropdown was fixed to match.
- **CopyButton calls `e.preventDefault()`** in the click handler so the `<a>` variant doesn't navigate away before the "Copied!" label swap fires.
- **Tabs `ensureTabIDs` + `resolveActiveTabID`**: auto-generate IDs for tabs that omit them (prevents `id="-tab"` invalid HTML + JS `querySelector('#')` crash), and default `ActiveTabID` to the first tab so exactly one tab has `tabindex="0"` (WAI-ARIA requirement).
- **Tooltip JS propagates `aria-describedby`** from the non-focusable wrapper `<div>` to the first focusable child (the trigger). Without this, screen readers never announce the tooltip text because `aria-describedby` must be on the focusable element.
- **Accordion uses native `<details>`/`<summary>`**: zero JavaScript, native keyboard support, built-in accessibility (implicit aria-expanded, role=group). Chevron rotation via CSS `details[open] [data-tc-chevron]` in `templates/custom.css`. The old JS toggle with grid-rows animation has been removed — `<details>` provides native open/close behavior.
- **LoadingButton hides default text via Tailwind arbitrary variant**: `[.htmx-request_&]:hidden` compiles to `.htmx-request .element { display: none }`. Never use fictional CSS classes like `htmx-hide-during-request` — it was never defined by HTMX or any CSS file, so the default text never hid during loading.
- **InlineLoadingOverlay uses `role="status"` not `aria-hidden="true"`**: HTMX indicators show/hide via CSS opacity, not DOM insertion. A static `aria-hidden="true"` is never toggled, so screen readers never announce the loading state. Use `role="status"` + `aria-live="polite"` instead.
- **HTMX retry counter: set and clear on the same element**: `data-tc-retry` is set on `event.detail.elt` (triggering element) but was cleared on `event.detail.target` (swap target) — different DOM nodes when `hx-target` points elsewhere. Always use `event.detail.elt` in both the set and clear paths.
- **ErrorPageProps.StatusCode overrides family-derived status**: The family-to-status-code map is too coarse (FamilyRejection → 400 for 400/403/404). Constructors that need a specific code set `StatusCode` explicitly; the handler checks `props.StatusCode` first, falling back to `FamilyStatusCode()` when unset (0).
- **ThemeToggle uses `querySelectorAll`**: `querySelector` only initializes the first toggle. Multiple `ThemeToggle` instances on a page need `querySelectorAll` for init + click handler must sync all instances' `aria-checked`.
- **localStorage wrapped in try/catch**: `setItem`/`getItem` throw `QuotaExceededError` in Safari private mode. ThemeScript and ThemeToggle both guard with try/catch.
- **ThemeScript before HTMX CDN scripts**: The FOUC-prevention script must run before the page paints. Placing it after synchronous CDN `<script>` tags delays first paint if the CDN is slow.
- **RadioGroup Required propagates to individual radios**: `aria-required` on `<fieldset>` is an incomplete substitute — browsers ignore it for constraint validation. Per HTML spec, `required` on any one radio in a group makes the group required.
- **FieldError has `role="alert"`**: bare `<p>` errors are invisible to aria-live. Add `role="alert"` so screen readers announce dynamic errors immediately. Guard against empty message to prevent rendering an empty red paragraph.
- **InputGroup both addons get `pointer-events-none`**: the right addon div was missing this class, intercepting clicks over the right ~40px of the input. Interactive addons (buttons) can override with `pointer-events-auto`.
- **ConfirmDelete `hx-confirm` is conditional**: render `hx-confirm` only when `Confirm != ""` — an empty `hx-confirm=""` shows a browser confirmation dialog with no text, confusing UX.
- **SwapOOB empty Selector omits colon**: `hx-swap-oob` format is `style:selector`. Empty selector was producing `outerHTML:` (trailing colon). HTMX resolves the element's own ID when no selector is given — omit the colon entirely.
- **Breadcrumb URL resolver uses `net/url.Parse`**: `strings.Contains(href, "://")` misses protocol-relative URLs (`//cdn.example.com`). Use `url.Parse` + `IsAbs()` + `strings.HasPrefix(href, "//")`.
- **Carousel uses CSS scroll-snap**: native touch/drag support via `snap-x snap-mandatory scroll-smooth` on the track. Slides use `snap-center`. JS simplified to `scrollBy`/`scrollTo` for prev/next/dot navigation + `scrollend` for dot sync (with scroll+debounce fallback). Replaced the old `translateX` + manual transform approach.
- **ViewTransitions API**: `htmx.ViewTransitions(ViewTransitionsProps{Global: true})` enables native View Transitions for HTMX swaps via `htmx.config.globalViewTransitions = true` (HTMX 2.0 built-in). Renders default cross-fade CSS with `prefers-reduced-motion` support. Graceful degradation — browsers without View Transitions do instant swaps.
- **Modern browser capabilities**: see `docs/research/modern-browser-capabilities.md` for the comprehensive analysis of native APIs (`<dialog>`, Popover API, `@starting-style`, `<details>`, scroll-snap, View Transitions, `content-visibility`, `:has()`) and the phased migration roadmap. CSS foundation for these APIs is in `templates/custom.css`.
- **Modal/Drawer use native `<dialog>`**: `showModal()`/`close()` replace ~200 lines of custom focus-trap JS. Browser handles focus trap, Tab cycling, Escape, focus restore, top-layer, backdrop, inert. CSS `.tc-overlay`/`.tc-modal`/`.tc-drawer` classes handle animations via `@starting-style` + `allow-discrete`. JS wrappers `tcOpenOverlay(id)`/`tcCloseOverlay(id)` kept for HTMX compat. Backdrop click detection: `e.target === dialog` (the `::backdrop` pseudo-element registers as the dialog).
- **Stylable Select (`appearance: base-select`)**: `SelectProps.Stylable: true` opts into the customizable `<select>` API. Emits `<button><selectedcontent></selectedcontent></button>` inside the select. CSS `.tc-select` in `templates/custom.css` styles button, picker (`::picker(select)`), options, arrow (`::picker-icon`), checkmark. Progressive enhancement — non-supporting browsers (Firefox, iOS Safari) ignore the structure and render native `<select>`.
- **Textarea AutoGrow**: `TextareaProps.AutoGrow` (default `true`) uses CSS `field-sizing: content` via `.tc-auto-grow` class. No JavaScript. `field-sizing` is Baseline 2024.
- **EnterKeyHint (unified API)**: Both `InputProps.EnterKeyHint` and `TextareaProps.EnterKeyHint` use the same typed `EnterKeyHintType` enum. Constants: `EnterKeyHintSend`, `EnterKeyHintDone`, `EnterKeyHintGo`, `EnterKeyHintNext`, `EnterKeyHintPrevious`, `EnterKeyHintSearch`, `EnterKeyHintEnter`. `EnterKeyHintTypeIsValid` included. Input also auto-derives a smart default from `InputType` (email→next, search→search, etc.) via `enterKeyHintValue()`; explicit `EnterKeyHint` overrides the auto-derived value.
- **Input search semantic landmark**: `Input` with `Type: InputSearch` auto-wraps in `<search>` element (Baseline 2023). Screen readers announce it as a search landmark. No API change — auto-detected from the InputType.
- **Form hx-validate**: `FormProps.Validate: true` emits `hx-validate="true"` for HTML5 constraint validation before HTMX submit. Pair with native `required`, `pattern`, `type="email"` etc. for client-side validation.
- **Image responsive delivery**: `ImageProps.SrcSet` and `Sizes` are typed string fields — no more `Attrs` workaround. Example: `SrcSet: "/img-480w.jpg 480w, /img-800w.jpg 800w", Sizes: "(max-width: 600px) 480px, 800px"`.
- **Table content-visibility**: `TableProps.LazyRows: true` applies `content-visibility: auto` to body rows via `.tc-content-auto` class (48px intrinsic height). When `CellPadding: TableCellPaddingCompact`, uses `.tc-content-auto-compact` (40px intrinsic height) to avoid scrollbar jitter. Browser skips rendering off-screen rows. Recommended for tables with 100+ rows.
- **Global accent-color CSS**: `templates/custom.css` sets `accent-color: blue-600` (light) / `blue-400` (dark) on checkboxes, radios, range inputs, and progress bars. Consumers override via `@theme { --color-blue-600: #custom; }`. No Go code changes needed — same theming model as all other components.

## Demo Infrastructure

- **Demo CSS**: `examples/demo/demo.css` compiles Tailwind CSS scanning ALL `.templ` files in the repo (`@source "../../**/*.templ"`). Path is relative to the CSS file location (`examples/demo/`), so `../../` reaches the repo root. The compiled CSS is embedded via `//go:embed static/app.css` and served at `/css/app.css`.
- **@source path gotcha**: Tailwind v4's `@source` directive resolves relative to the CSS file, not the CWD. Using `@source "./**/*.templ"` from `examples/demo/demo.css` only scans `examples/demo/*.templ` — missing ALL component classes from `display/`, `forms/`, `feedback/`, etc. Must use `@source "../../**/*.templ"` to scan the entire repo.
- **Custom CSS**: Component-specific CSS (dialog animations, stylable select, auto-grow textarea, scroll-snap, accordion chevron, accent-color) lives in `templates/custom.css`. Both `templates/app.css` (consumer template) and `examples/demo/demo.css` (demo) import it via `@import "./custom.css"` / `@import "../../templates/custom.css"`. Single source of truth.
- **Dockerfile 3-stage pipeline**: CSS (Node 22, compiles Tailwind) → Go binary (templ generate + go build, overwrites committed CSS with fresh compile from Stage 1) → Distroless runtime. CSS is always freshly compiled during Docker build — the committed `static/app.css` is never stale because it's overwritten. `.dockerignore` excludes `.git/`, `website/`, `docs/` from Docker context (~653MB → ~15MB).
- **Demo endpoints**: `/health` returns `{"status":"ok"}` for Cloud Run health checks. `/css/app.css` serves embedded CSS with `Cache-Control: public, max-age=31536000, immutable`. `/api/load-more` and `/api/delete` are mock HTMX endpoints for interactive demo components.

## Release Convention: One-Commit Release

Established with v0.4.0 → v0.5.0 → v0.6.0. Each version is cut with a **single
release commit** at the tip of `master`, even if many feature/fix commits preceded it.
The release commit message is the canonical user-facing description of what changed.

**`[Unreleased]` must be warm at all times.** Every feature/fix commit that lands on
`master` must add its changelog entry to the `[Unreleased]` section immediately — not
deferred to release time. The release script (`scripts/release.sh`) enforces this by
failing if `[Unreleased]` has no body.

**Release commit message structure:**

```
release: <version> — <one-line summary>

<one-paragraph "why this version" / headline summary>
<one-paragraph "what's in it" / feature highlights>
<one-paragraph "notes" / breaking changes, deprecations, migration paths>

💘 Generated with Crush
Assisted-by: Crush:MiniMax-M3
```

**Release commit body must include:**

- The version bump in `utils/version.go`
- The CHANGELOG heading (e.g., `## [0.6.0] — YYYY-MM-DD`) replacing `[Unreleased]`,
  with a fresh empty `## [Unreleased]` inserted above it
- The release notes in the commit body **and** the CHANGELOG (both kept in sync)

**Tag format:** annotated + SSH-signed, message `<version>: <one-line summary>`,
e.g., `v0.6.0: typed Props structs, tooltip a11y, error handler hardening`.
Sign with the same key used for v0.5.0.

**Post-release commits** (e.g., backfilling tests, fixing docs, regenerating
after a templ patch release) are committed normally on `master` and roll up
into the next minor/patch release. Do not retag the same version.

**Two-commit alternative was considered and rejected** for v0.6.0: it
doubles the commit noise without improving the reviewability of the
release. The one-commit convention keeps the release timeline obvious
in `git log` and matches every prior release in this repo.

**To cut a release:** use `scripts/release.sh` (see "Release Script" below).

## Release Script

`scripts/release.sh` automates the full release cut in one command:

```bash
scripts/release.sh <new-version> "<release-summary>"
# Example: scripts/release.sh 0.7.0 "typed HTMX retry, Drawer motion-reduce"
```

What it does:

1. Validates the working tree is clean and on `master`
2. Confirms the new version is greater than the current one (via `sort -V`)
3. Collects release notes (`--notes-file FILE`, or auto-extracted from CHANGELOG `[Unreleased]`)
4. Installs an `EXIT`-trap rollback (`release_rollback`) that restores `utils/version.go`, `CHANGELOG.md`, and `FEATURES.md` if any later step fails — so a failed verify never leaves a dirty tree
5. Bumps `utils.Version` via in-place sed
6. Moves the `[Unreleased]` body under a new `## [<version>] — YYYY-MM-DD` heading (inserts a fresh empty `[Unreleased]` above)
7. Bumps `FEATURES.md` `**Version:**` + `**Updated:**` date (the three version files must move together; `utils.TestVersionMatchesFeatures` enforces it)
8. Regenerates `*_templ.go` and runs the full verify suite (build + test + lint)
9. Asserts the version drift-guard (`TestVersionMatches(Changelog|Features)`)
10. Stages and commits as `release: <version> — <summary>` (one-commit convention; body carries the release notes, `Assisted-by: Crush:${CRUSH_MODEL}`)
11. Creates an annotated, SSH-signed tag `v<version>: <summary>`

The script does **not** push. House rule: "NEVER PUSH TO REMOTE". Push manually
after reviewing the release commit and tag with `git show v<version>` and
`git show <commit>`.

## Lint Command

```bash
# examples/ excluded via .golangci.yml paths exclusion
golangci-lint run ./...
```

**Disabled linters (do NOT re-enable — fundamentally incompatible with this codebase):**

- `ireturn` — every component returns `templ.Component` (an interface) by design; the linter's premise is antithetical to templ.
- `godoclint` — demands exactly one `// Package` godoc per package, but the repo intentionally documents per-file.
- `testableexamples` — `Example*` funcs render HTML that is verbose and version-dependent; output isn't asserted.

**Reconciled at v0.18.1:** commit 73395d9 expanded to 67 linters but left `golangci-lint run` failing (187 findings, CI would have gone red on first push). The config now uses an explicit depguard allow-list (the `$module` token did not resolve — use literal `github.com/larsartmann/templ-components` + the three runtime deps), extends `varnamelen`/`mnd` ignore lists, and excludes test files from `err113`/`makezero`/`varnamelen`/`gocheckcompilerdirectives`. If you add a linter, run `golangci-lint run` to 0 findings before committing.

## `encoding/json/v2` Adoption

This library uses `encoding/json/v2` + `encoding/json/jsontext` (Go 1.26+ with
`GOEXPERIMENT=jsonv2`). The pre-commit hook (`scripts/pre-commit.sh`) sets
`GOEXPERIMENT=jsonv2` automatically. The `.golangci.yml` enables the
`goexperiment.jsonv2` build tag. BuildFlow also auto-detects and sets it.

**Consumers** must set `GOEXPERIMENT=jsonv2` when building (or wait for Go 1.27
where it becomes stable). The `errorpage` package uses `json.MarshalEncode` +
`jsontext.NewEncoder` for JSON error responses. Other packages (breadcrumbs,
tests) still use `encoding/json` v1 — both coexist fine under the experiment flag.

## Conventions

- **Naming hygiene:** `forms/radio_go.go` renamed to `forms/radio.go` (the `_go.go` suffix falsely implied generated code). `icons.Close` added as alias for `icons.X` (prefer `Close` in new code — `X` is a single-letter identifier with poor discoverability). `errMsg` → `errorMessage` (no abbreviations). `cleanMessage` → `sanitizeErrorMessage` (precise verb). `htmxMainSRIDefault` → `sriHTMXMainDefault` (consistent word order with `sriHTMXMainByVersion`).
- **RTL keyboard mapping:** `display.Tabs` and `display.Dropdown` JS handlers check `document.documentElement.getAttribute('dir') === 'rtl'` and swap ArrowLeft/Right mappings per WAI-ARIA APG. In LTR: ArrowRight=next, ArrowLeft=prev. In RTL: ArrowLeft=next, ArrowRight=prev.
- **Sub-template extraction pattern (ADR 0010):** Extract when 2+ callers, 5+ lines, clear domain name. Do NOT extract for single caller, demo code, no clean name, 8+ parameters, or same-file callers. See `docs/adr/0010-sub-template-extraction-pattern.md`.
- **Motion-reduce compliance test:** `utils.TestMotionReduceCompliance` scans all `.templ` files for `transition-*`/`animate-*` classes without `motion-reduce:` fallback. Run via `go test ./utils/... -run TestMotionReduce`.
- **SKILL.md drift-guard:** `utils.TestSkillComponentCount` logs actual vs documented component count. Informational (not failing) — intended to surface drift during code review.
- **Fuzz tests:** `forms.FuzzInputType`, `forms.FuzzFormMethod`, `display.FuzzButtonHTMLType` verify enum validation never panics on arbitrary input. Run via `go test -fuzz=. -run=Fuzz ./...`.
- **Benchmark suites:** Now in 7 packages (display, feedback, navigation, forms, layout, htmx, icons, utils). Run via `go test -bench=. -benchmem ./...`.
- **goconst zero issues:** All repeated string literals extracted to named constants. `msgGoBack` in constructors.go uses `notFound404GoBackText` (single source of truth across errorpage package).
- **Golden package coverage:** 81.8% (was 70.5%). New tests cover: `-update` flag, `MkdirAll`, normalization edge cases, diff identical/multi-line, `lineAt` out-of-range.
- **FooterProps:** `navigation.Footer` now takes `FooterProps` (embeds `BaseProps`) instead of a raw `brandText string`. All components in the library now accept `BaseProps`.
- **CSRFTokenName:** `forms.FormProps` has a `CSRFTokenName` field (defaults to `"csrf_token"`) for framework compatibility.
- **ErrorPage/NotFound404 landmark:** Both use `<main>` (not `<div role="region">`) for WCAG 2.4.1 Bypass Blocks compliance.
- **FromError fallback:** Unknown errors return `FamilyCorruption` (→500), not `FamilyInfrastructure` (→503). An unrecognized error is a bug, not a transient outage.
