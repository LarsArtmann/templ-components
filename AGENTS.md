# AGENTS.md — templ-components

## Multi-Module Structure (go.work + 6 modules)

This repo uses a **Go multi-module workspace** with 6 modules coordinated by `go.work`:

| Module            | Path                   | Contains                                                      | Purpose                                               |
| ----------------- | ---------------------- | ------------------------------------------------------------- | ----------------------------------------------------- |
| **root**          | `go.mod`               | display, feedback, forms, layout, navigation, htmx, internal/ | Core UI components                                    |
| **svg**           | `svg/go.mod`           | svg (promoted from `internal/svg`)                            | SVG rendering primitives                              |
| **utils**         | `utils/go.mod`         | utils                                                         | BaseProps, Class(), EnsureID, test helpers            |
| **icons**         | `icons/go.mod`         | icons                                                         | Named SVG icons (lightweight: no tailwind-merge-go)   |
| **errorpage**     | `errorpage/go.mod`     | errorpage                                                     | Error pages + HTTP handler (isolates go-error-family) |
| **examples/demo** | `examples/demo/go.mod` | demo                                                          | Demo binary                                           |

**Each sub-module has replace directives** for its siblings (e.g., `icons/go.mod` has
`replace github.com/larsartmann/templ-components/svg => ../svg`). This ensures
`GOWORK=off go build` works per-module — essential for CI and external consumers.

**go.work is committed** (un-ignored via `!go.work` in `.gitignore`). BuildFlow may
re-add `go.work` to `.gitignore` on pre-commit — the `!` negation after it keeps it tracked.

## Build & Test Commands

```bash
# Full build (required before go build after .templ changes)
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./...

# Tests (root module)
go test ./...

# Tests (all modules including sub-modules)
for mod in svg utils icons errorpage examples/demo; do (cd $mod && go test ./...); done && go test ./...

# Lint (must include ./svg/...)
golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./svg/... ./internal/...

# GOWORK=off isolation test (verify replace directives work standalone)
for mod in svg utils icons errorpage examples/demo; do (cd $mod && GOWORK=off go build ./... && GOWORK=off go test ./...); done && GOWORK=off go build ./...

# All-in-one verification
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./... && go test ./... && golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./svg/... ./internal/...
```

## CRITICAL: Generated `*_templ.go` Files MUST Be Committed

This is a **templ library**, not an application. The Go module proxy (proxy.golang.org) fetches
source from the Git tag — it does **not** run `templ generate`. Without committed `*_templ.go`
files, consumers get uncompilable code (`undefined` errors on every component function).

- The `.gitignore` uses `!*_templ.go` to override the global gitignore's `*_templ.go` entry
- After editing any `.templ` file, always run `templ generate ./...` and commit the updated `*_templ.go` files alongside the source
- Never add `*_templ.go` back to `.gitignore` — this is the standard pattern for publishable templ packages
- 59 generated files across 6 modules: root (display, errorpage, feedback, forms, htmx, layout, navigation, internal/golden), svg, utils, icons, examples/demo
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
- **Theming:** Components emit standard Tailwind classes (`bg-blue-600`). Consumers override via `@theme { --color-blue-600: #custom; }` in their CSS. No Go code changes needed. See `templ-components-theme.css` for semantic alias examples.
- **ComponentProps interface:** `utils.ComponentProps` with `GetBaseProps()`/`SetBaseProps()` on `*BaseProps` (pointer receivers for `recvcheck`). All 26+ props structs auto-satisfy via method promotion.
- **Accessibility — motion-reduce:** `motion-reduce:transition-none motion-reduce:duration-0` on all transitions, `motion-reduce:animate-none` on all animations (spinner, skeletons, toast enter/exit, modal, accordion)
- **Dark mode colors:** All components use `gray-*` exclusively (no mixed `slate-*`/`gray-*`). Dark mode via class strategy: `@custom-variant dark (&:where(.dark, .dark *))` toggled by `layout.ThemeScript()` + `layout.ThemeToggle()`
- **CI:** `.github/workflows/ci.yaml` — lint (golangci-lint), build+test with `templ generate`, coverage artifact. Pre-commit: `.git/hooks/pre-commit` → `scripts/pre-commit.sh`
- **Import graph (multi-module):** Module DAG: `svg` ← `icons`; `utils` (leaf); `icons,svg,utils` ← root (display,feedback,forms,layout,navigation,htmx); `icons,utils` ← errorpage; all ← examples/demo. Production deps: `icons → svg`, `display → icons,svg,utils`, `feedback → icons,svg,utils`, `forms → icons,utils`, `layout → icons,utils`, `navigation → icons,svg,utils`, `htmx → utils`, `errorpage → icons,utils`
- **No circular imports** allowed
- **AriaLabel propagation:** All components with `BaseProps` propagate `AriaLabel` to root element. Components with hardcoded aria-labels (Nav, Pagination, Breadcrumbs, StepIndicator) allow AriaLabel override via `utils.Ternary`
- **SVG paths:** Shared constants in `internal/svg` (PathChevronDown, PathChevronSmall, PathArrowUp/Down/Left/Right, PathAvatarFill) — single source of truth

## Code Conventions

- All component props embed `utils.BaseProps` (exception: `layout.PageProps`) — all auto-satisfy `utils.ComponentProps` interface
- All root elements propagate `props.Class`, `props.Attrs`, `props.ID`, and `props.AriaLabel` from BaseProps (26/26 components, including NavLink/MobileNavLink)
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
- Muted text: shared `mutedTextClass` constant (in `display/shared.go`) for the standard secondary-text pattern (`text-sm text-gray-500 dark:text-gray-400`); callers combine it with a margin via `utils.Class(mutedTextClass, "mt-N")`. Used by Card subtitle, PageHeader subtitle, EmptyState description.
- Pagination: `paginationPageItem`/`paginationEllipsisItem` sub-templates wrap `activeSpanOrLink`/`paginationEllipsis` in `<li>` so the page-range loop body stays flat.
- EmptyState: single `emptyStateAction(text, href, attrs)` helper renders an anchor when `href != ""`, else a button (replaces the former link/button pair).
- HTMX loading: accepts `templ.Component` spinner parameter (decoupled from feedback package)
- Toast icons: generated from Go `iconPathData` via `icons.IconPathJS()` (single source of truth)
- TrendDirection: `TrendNone = "none"` (non-empty sentinel, not "")
- Layout: `Minimal(MinimalProps)` uses props struct like `Base(PageProps)`
- Layout: `PageProps.HTMXCDN` overrides the CDN base URL for htmx scripts. Empty defaults to `https://cdn.jsdelivr.net/npm`. Both `htmx.org` and `htmx-ext-response-targets` URLs derive from this value. Consumers with a different CSP can set e.g. `HTMXCDN: "https://unpkg.com"` without forking the library.
- Interactive: `cursor-pointer` on buttons, `caret-blue-600 dark:caret-blue-400` on inputs, `scroll-smooth` + selection colors on body, `shadow-sm` on card shell
- Modal: focus save/restore via `data-tc-prev-focus` attribute on open, restored on close
- NavLink/MobileNavLink: both render through the shared `navLinkAnchor` sub-template; each supplies an active/inactive base-class builder (`navLinkClasses`, `mobileNavLinkClass`) and `navLinkAnchor` merges `props.Class` via `utils.Class()` so consumer Tailwind overrides resolve correctly. Do NOT assert ordered class substrings in tests — `utils.Class`/tailwind-merge reorders classes; use `utils.AssertContainsAll` for multi-token checks.
- InputType: validates via `inputType()` with `validInputTypes` map; panics on unknown, defaults empty to `"text"`
- Structural variants (TabsVariant, DropdownPosition, TrendDirection): use `if`-branch for DOM structure, not map lookup — map pattern is for pure class lookups only
- `forms.SanitizeID`: exported utility for library consumers; also used internally by `forms.RadioGroup` to derive per-option IDs from option values
- Enum validation: 0 panic-on-unknown, 14 map+fallback (InputType, ButtonHTMLType, FormMethod now included), structural variants use if-branch. InputType falls back to "text", ButtonHTMLType/FormMethod fall back to HTML-spec defaults ("button"/"GET"), icons.Name falls back to Question icon. Only remaining panic: icon path data integrity check (stray `|` separators).
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
- LoadMore: cursor appended as `?cursor=` query param (detects existing `?` for `&`). `hx-swap="outerHTML"` + `hx-target="this"` for self-replacement.
- StatCard HTMX: `HxGet`/`HxTarget`/`HxSwap` typed fields on both `<a>` and `<div>` variants. When empty, attributes are omitted.
- Card.Body: `Body templ.Component` slot — when set, overrides children. Backward compatible.

- Toast JS: dismiss icon from `icons.IconPathJS()` via `tcToastIcons.dismiss`
- Table: row cells auto-padded/truncated to match header count
- **Error handler:** `errorpage/handler.go` provides `ErrorHandler(err, cfg)` returning `http.Handler`, `FromError(err)` for type-safe conversion from go-error-family errors, 6 pre-built constructors (`NotFound`, `Forbidden`, `BadRequest`, `Conflict`, `ServiceUnavailable`, `InternalError`), `WriteError`/`WriteErrorPage` convenience wrappers, `HTMLShell` mode for valid HTML documents, `JSON` mode for API/HTMX responses. Uses `errors.AsType[errorfamily.Classified]()` for go-error-family integration.
- **Error families:** `errorpage` package integrates with go-error-family via `FamilyFromErrorFamily()` converter + `ParseFamily()` for string-based lookup. `FromError()` extracts Why/Fix defaults from go-error-family's `Family.DefaultWhy()`/`DefaultFix()` methods.
- **Error components:** `ErrorPage` (full-page), `NotFound404` (dedicated 404 with hero numeral + search + links), `ErrorDetail` (inline card), `ErrorAlert` (family-aware alert) in `errorpage/`
- **NotFound404:** `errorpage.NotFound404(props NotFound404Props)` — dedicated 404 page with large gradient numeral (`text-[8rem]`), optional search form (`SearchAction`), quick-links card grid (`[]NotFoundLink`), and "Go home" / "Go back" buttons. Unlike `ErrorPage` (family-colored error card), NotFound404 is a welcoming navigation aid using neutral blue/indigo palette. `DefaultNotFound404Props()` returns full defaults. `DefaultNotFoundLinks()` returns starter links (Home, Documentation). Types and constructors live in `notfound404_types.go`. All string constants (titles, messages, labels) defined as package-private `notFound404*` constants for goconst compliance. Shares the `tcGoBackAttached` singleton with `ErrorPage`.
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
- Icons: 101 total (100 path icons + Spinner). Added BuildingOffice2 (tenants/orgs), Key (credentials), ArrowRightOnRectangle (logout) — gaps surfaced by cqrs-htmx/adminui.
- Thread safety: `utils.Class()` uses `sync.Mutex` to protect tailwind-merge-go's shared LRU cache from concurrent access. Required even though the LRU has internal mutexes — they don't protect the full Merge() call sequence.
- DropdownItemKind: typed enum (`DropdownItemLink`, `DropdownItemButton`) with backward compat via `IsLink()` fallback to Href-based discrimination.
- PageHeader: `display.PageHeader(PageHeaderProps)` with Title, Subtitle, Breadcrumb slot, Action slot. No navigation import — uses `templ.Component` slots for compositional freedom.
- DefinitionList: `display.DefinitionList(DefinitionListProps)` with `[]DefinitionItem{Term, Detail, DetailComponent}`. Two-column `<dl>` grid for metadata/key-value display.
- ListNote: `display.ListNote(ListNoteProps{Shown, Total})` — renders "Showing N of M" truncation notice when Total > Shown. `role="status"` for a11y.
- SidebarNav: `navigation.SidebarNav(SidebarNavProps)` — vertical sidebar with Brand/Footer slots, icon+label nav items, CurrentPath auto-active detection, `aria-current="page"` on active item.
- icons.IconPathData: exported function returning raw SVG path d-strings for consumers needing full `<svg>` wrapper control (used by cqrs-htmx/adminui for icons-only adoption without Tailwind).
- icons-only adoption: the `icons` sub-module depends only on `svg` + `templ` — no tailwind-merge-go, no CSS framework. This is a natural property of icons, not a portability strategy. See `docs/icons-only-adoption.md`.
- Grid: `display.Grid(GridProps)` — responsive grid container with typed `GridCols` enum (1–6) and `gridColsLookup` map+fallback. Replaces the repeated `grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4` pattern. Children passed via `{ children... }`.
- StatCard.Href: `StatCardProps.Href` renders the whole card as an `<a>` with hover shadow, focus ring, cursor-pointer when set. Mirrors `Badge.Href` pattern. Shared body extracted to `statCardInner` sub-template so linked/unlinked layouts don't diverge.
- layout.Script: `layout.Script(nonce, src, attrs)` — CSP-safe `<script src>` tag that auto-injects the nonce. Prevents forgetting `nonce={...}` under strict CSP. Use instead of raw `<script>` tags for external scripts.
- SkeletonCardGrid: `feedback.SkeletonCardGrid(count)` — renders N skeleton cards in a responsive 3-col grid under a single `role="status"`. Pairs with `display.Grid` for the common dashboard loading pattern. Non-positive count falls back to 1.
- SimpleNav.RightItems: `SimpleNavProps.RightItems` forwards a `templ.Component` (e.g. ThemeToggle, sign-in button) to `Nav.RightItems`. Consumers no longer render right-side nav items in a separate flex container.
- PageProps auto-injects: `DefaultPageProps()` auto-injects `<link rel="stylesheet" href="/app.css">` (CSSPath) and `<script src="...htmx.org@...">` (HTMXVersion). Set either to `""` to suppress. Documented in godoc on both fields + DefaultPageProps.
- ProgressBar clamp: uses `max(0, min(100, props.Progress))` (Go 1.21+ builtin) instead of manual if-branch clamping.
- OverlayKind: typed enum (`OverlayModal`, `OverlayDrawer`) replaces untyped `closeKind`/`componentName` strings on `overlayShellProps`. The `componentName()` method derives JS function names from the kind.
- CopyButton: `execCommand('copy')` fallback for non-secure HTTP contexts. `role="status"` + `aria-live="polite"` on label span for screen reader feedback.
- Table.Body: `TableProps.Body templ.Component` slot — when set, overrides Rows for custom `<tr>` rendering. Follows Card.Body pattern.

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
3. Bumps `utils.Version` via in-place sed
4. Inserts a new `## [<version>] — YYYY-MM-DD` heading into CHANGELOG, replacing
   `[Unreleased]` and adding a fresh empty `[Unreleased]` block above
5. Prompts for release notes on stdin (Ctrl-D on an empty line to finish)
6. Regenerates `*_templ.go` and runs the full verify suite (build + test + lint)
7. Asserts the version drift-guard test passes (CHANGELOG heading matches `utils.Version`)
8. Stages and commits as `release: <version> — <summary>` (one-commit convention)
9. Creates an annotated, SSH-signed tag `v<version>: <summary>`

The script does **not** push. House rule: "NEVER PUSH TO REMOTE". Push manually
after reviewing the release commit and tag with `git show v<version>` and
`git show <commit>`.

## Lint Command

```bash
# Must lint specific packages — examples/ excluded via .golangci.yml
golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./svg/... ./internal/...
```
