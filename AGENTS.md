# AGENTS.md — templ-components

## Module Structure (7-module workspace)

This repo is a **7-module Go workspace** (`github.com/larsartmann/templ-components`) coordinated by `go.work` (local dev, gitignored) + `replace` directives (CI/consumers). See ADR-0034 for the split rationale and DAG. The modules: **root** (core UI + recipes + integration + demo + CLI), **utils** (leaf: BaseProps, Class(), EnsureID, svg, cdn, golden), **icons** (106 SVG icons, icons-only adoption), **errorpage** (isolates go-error-family), **charts/echarts** (opt-in adapter), **htmx** (HTMX loading/error/OOB components), **datastar** (Datastar runtime + SSE LiveRegion). Additionally, **visualtest** is a separate module for visual regression tests (not part of the library DAG). See `docs/modularization/README.md` for contributor setup.

### Root-module packages (10)

| Package             | Contains                                            | Purpose                                                                                                                                                                |
| ------------------- | --------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `display`           | 40 UI components                                    | Cards, tables (Table + DataTable), modals, badges, buttons, avatars, carousel, context menu, hover card, **SVG charts** (LineChart, PieChart, AreaChart)               |
| `feedback`          | 14 components                                       | Alerts, toasts, spinners, skeletons, progress bars                                                                                                                     |
| `forms`             | 21 components                                       | Inputs, selects, toggles, combobox, slider, rating, tags input, calendar, validation                                                                                   |
| `layout`            | 10 components                                       | Page shell, theme toggle, CSP-safe script/style tags, **body-layout primitives**: AppShell, Container, Split, Stack                                                    |
| `navigation`        | 12 components                                       | Nav bars, pagination, breadcrumbs, sidebar, EndOfList                                                                                                                  |
| `htmx`              | **Separate module** — 9 components                  | HTMX loading, error handling, OOB swaps, View Transitions, PolledRegion                                                                                                |
| `datastar`          | **Separate module** — 4 components + action helpers | Datastar runtime injection, SSE-powered LiveRegion, loading Indicator, SSE error handling. Pins version via `go-datastar/static`. Opt-in complement to HTMX (ADR-0030) |
| `recipes`           | 4 composition screens                               | Dashboard, SettingsLayout, LoginCard, AuthLayout — screen-level compositions of display/forms/layout/navigation (ADR-0019)                                             |
| `internal/contract` | Contract tests                                      | Cross-package interface verification                                                                                                                                   |
| `integration`       | CSP nonce tests                                     | Asserts nonce on all inline scripts                                                                                                                                    |
| `examples/demo`     | Demo binary                                         | Showcases components                                                                                                                                                   |
| `cmd/tc`            | CLI tool                                            | Component scaffolding (excluded from lint — uses different conventions)                                                                                                |

> **Note:** `go.work` and `go.work.sum` are in `.gitignore` (local dev only). CI and consumers use `replace` directives in each module's `go.mod`. `internal/` packages (`svg`, `cdn`, `golden`) were promoted to `utils/` sub-packages because Go's `internal/` rule blocks cross-module access.

> **Datastar v1.0.2 runtime facts (2026-08-21 audit):** see `docs/datastar-runtime-facts.md` — only `datastar-patch-*` SSE events (keyed datalines, blank line terminates), lifecycle errors only on `datastar-fetch` (no `datastar-sse-error`), CSP needs `'unsafe-eval'`, clean stream EOF only reconnects under `retry: 'always'` (`LiveRegionProps.Retry`). Wire format pinned by `examples/demo/sse_test.go`; bundle contract pinned by `datastar.TestPinnedRuntimeBundleContract`.

## Build & Test Commands

```bash
# Full build (workspace mode via go.work — builds all 7 modules)
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./...

# Tests (all modules via go.work)
go test ./...

# Per-module isolation tests (verify each module builds standalone without go.work)
for mod in utils icons errorpage charts/echarts datastar; do (cd "$mod" && GOWORK=off go test ./...); done

# All-in-one verification
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./... && go test ./... && nix run .#lint
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
nix run .#visual    # pixel-level visual regression tests (headless Chromium). See docs/visual-testing.md
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
- 114 generated files across all packages
- **BuildFlow gotcha:** the BuildFlow pre-commit `templ-generate` step re-appends `*_templ.go` to `.gitignore` on every run, which (being the last pattern) overrides the `!*_templ.go` unignore and hides generated files from `git status`. This is harmless for already-tracked files (gitignore cannot untrack), but any NEW component's `*_templ.go` will be invisible until `git add -f`. After each commit, check `git status` for a re-added `*_templ.go` line and remove it. Consider fixing this in BuildFlow itself (it is `larsartmann/buildflow`).
- **BuildFlow daemon commit messages (identified 2026-07-28, T13):** The auto-commit daemon commits with generic, hallucinated messages (e.g., `"chore: update project configuration and documentation"`) authored as `"Unknown Author <unknown@example.com>"`. These commits are invisible to `git log --grep` for specific features. 5+ sessions have documented this. The daemon also has a 60s budget and does NOT run `go test ./...`, meaning the `TestGolangciDisabledLinters` guard only fires in CI. Root cause: the daemon generates messages from a template, not from `git diff --stat`. Fix requires modifying BuildFlow (`larsartmann/buildflow`). The `.golangci.yml` regression root cause (T1) is directly related — the daemon commits a stale working tree without running tests.

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
- **Dark mode color convention:** Light mode uses `-600` shade for backgrounds (`bg-blue-600`), dark mode uses `-500` (`dark:bg-blue-500`). Light mode uses `-600` for text (`text-blue-600`), dark mode uses `-400` (`dark:text-blue-400`). Neutral text: `text-gray-500` → `dark:text-gray-400`, `text-gray-400` → `dark:text-gray-500`. Every neutral and semantic color class MUST have a `dark:` variant — enforced by `utils.TestDarkModeCompliance` + `utils.TestDarkModeSemanticColors` (both now pass). Exceptions: Toggle thumb (`bg-white` both modes), SidebarNav (permanently dark sidebar), avatar silhouette icon (`text-blue-200` decorative).
- **Dark mode compliance tests:** `utils.TestDarkModeCompliance` scans all `.templ`/`.go` source files for neutral colors (`text-gray-*`, `bg-white`, `bg-gray-*`, `border-gray-*`, `ring-gray-*`) without `dark:` variants. `utils.TestDarkModeSemanticColors` scans for semantic colors (`bg-blue-600`, `text-red-600`, etc.) without `dark:` variants. Both now pass. Run via `go test ./utils/... -run TestDarkMode`. For the full dark mode strategy analysis (Tailwind v4 default is `prefers-color-scheme`, three consumer paths, `@theme` palette override pattern), see `docs/dark-mode-research.md`.
- **CI:** `.github/workflows/ci.yaml` — lint (golangci-lint), build+test with `templ generate`, coverage artifact. Pre-commit: `.git/hooks/pre-commit` → `scripts/pre-commit.sh`
- **Import graph:** `utils` (leaf module with `utils/svg`, `utils/cdn`, `utils/golden`); `icons` → utils; `errorpage` → utils, icons; `charts/echarts` → utils; `htmx` → utils; `datastar` → utils, go-datastar/static; root module (display, feedback, forms, layout, navigation, recipes) → utils, icons, errorpage, charts/echarts, htmx, datastar. All 7 modules form a strict DAG. Production deps: `icons → utils/svg`, `display → icons,utils`, `feedback → icons,utils`, `forms → icons,utils`, `layout → icons,utils`, `navigation → icons,utils`, `htmx → utils` (separate module), `datastar → utils, go-datastar/static` (separate module, zero transitive deps), `errorpage → icons,utils`, `charts/echarts → utils`, `recipes → display,icons,layout,utils`. Root module `require`s all sub-modules for backward compatibility (old import paths resolve to the sub-modules automatically).
- **No circular imports** allowed
- **AriaLabel propagation:** All components with `BaseProps` propagate `AriaLabel` to root element. Components with hardcoded aria-labels (Nav, Pagination, Breadcrumbs, StepIndicator) allow AriaLabel override via `utils.Ternary`
- **SVG paths:** Shared constants in `utils/svg` (PathChevronDown, PathChevronSmall, PathArrowUp/Down/Left/Right, PathAvatarFill) — single source of truth

## Code Conventions

- All component props embed `utils.BaseProps` (exception: `layout.PageProps`) — all auto-satisfy `utils.ComponentProps` interface
- All root elements propagate `props.Class`, `props.Attrs`, `props.ID`, and `props.AriaLabel` from BaseProps (26/26 components, including NavLink/MobileNavLink)
- Class attributes use `utils.Class()` for Tailwind conflict resolution (exception: `templ.KV` conditionals where comma-join is required)
- **RTL/i18n: use logical CSS properties exclusively.** Never use `ml-`/`mr-`/`pl-`/`pr-`/`left-`/`right-`/`text-left`/`border-l-`/`border-r-` — use `ms-`/`me-`/`ps-`/`pe-`/`start-`/`end-`/`text-start`/`border-s-`/`border-e-` instead. These are CSS logical properties that automatically mirror in RTL (`dir="rtl"`). Exception: `left-1/2 -translate-x-1/2` for centering (not directional).
- **Motion: use shared transition constants.** Use `transitionFast` (150ms), `transitionNormal` (200ms), `transitionColors`, `transitionTransform` from `display/shared.go` instead of inline timing strings. All include `motion-reduce:*` fallbacks. Wire into CopyButton, Accordion, Modal, Drawer — do NOT leave inline `transition-colors motion-reduce:...` strings when a constant matches.
- **Container queries: use `@container` for context-responsive components (ADR-0018).** 8 components have a `ContainerAware` bool flag. When true, the component emits a `<div class="@container">` wrapper (or adds `@container` to its root for Pagination) and swaps viewport breakpoints (`sm:`/`md:`/`lg:`) for container variants (`@sm:`/`@md:`/`@lg:`) via parallel lookup maps. **Since v2.0, `Grid` and `Split` default `ContainerAware: true`** (opt-out); **`Card` was reverted to default `false`** (post-v1.8.2, bugfix) — the `@container` wrapper's `container-type: inline-size` containment suppresses intrinsic width, collapsing cards to zero width inside shrink-to-fit parents (flex rows, inline-block, auto grid columns; proven by `visualtest/testdata/card/*`). `Nav`, `DefinitionGrid`, `Form`, `Pagination`, and `SkeletonCardGrid` still default `false` (opt-in). Container-aware components: `Grid.ContainerAware`, `Card.ContainerAware`, `Nav.ContainerAware`, `Split.ContainerAware`, `DefinitionGrid.ContainerAware`, `Form.ContainerAware`, `Pagination.ContainerAware`, `SkeletonCardGrid.ContainerAware`. **Fluid typography** via container query units (`cqi`): six `.tc-fluid-*` utility classes in `templates/custom.css` size text with `clamp(min, Ncqi + base, max)` — see `docs/recipes/fluid-typography.md`. **Do NOT expand `ContainerAware` to marginal candidates** (`Container`, `Breadcrumbs`, `EmptyState`, `NotFound404`, `Footer`) — evaluated and rejected in `docs/container-query-strategy.md`; none meet all three ADR-0018 criteria. **Web Components / Shadow DOM are permanently rejected** (ADR-0033) — Shadow DOM breaks the Tailwind theming model; the library achieves "use the platform" via native APIs instead. **CRITICAL: after adding `@md:` or other container variant classes to `.templ` files, the demo CSS must be recompiled** — Tailwind v4 scans `.templ` files at CSS compile time; the committed `examples/demo/static/app.css` will be stale until recompiled via `nix run .#build` or the Dockerfile pipeline.
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
- FeedbackType: canonical `FeedbackType` enum (`FeedbackSuccess/Error/Warning/Info`). **v2.0 removed** the `AlertType`/`ToastType` aliases and `AlertSuccess`/`ToastSuccess`/etc. constants — use `FeedbackType`/`FeedbackSuccess`/etc. directly. See ADR-0022.
- Icons: `iconPathData` map with `|` separator for multi-path icons. `iconPaths()` validates no empty segments (panics on stray `|`). `allIconNames()` auto-generated from `iconPathData` + Spinner — no manual list to maintain.
- Form errors: `ErrorAttrs(id, errMsg, helpTextID)` helper returns `templ.Attributes` for aria-invalid/aria-describedby
- Shared constants: `cardShellClass`, `mutedTextClass` (in `display/shared.go`) — use for consistent card styling and secondary-text pattern.
- **HTMX loading (v2.0 default: self-host).** `DefaultPageProps()` sets `HTMXSrc: HTMXSelfHost`, which embeds HTMX 2.0.10 inline via `//go:embed` (`layout/embed.go`) and renders it as `<script nonce="...">`. No external CDN request. Consumers who prefer CDN loading set `HTMXSrc: ""` and provide `HTMXVersion` to use the jsDelivr CDN path. `PageProps.HTMXCDN` overrides the CDN base URL (defaults to `https://cdn.jsdelivr.net/npm`).
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
- Combobox JS: global singleton `tcComboboxAttached` handler for input filtering, click-to-select, focus/blur dropdown management, Escape key dismissal. CSP-safe with `nonce={ props.Nonce }`.
- CopyButton JS: global singleton `tcCopyAttached` handler — click delegation on `[data-tc-copy]`, clipboard write via `navigator.clipboard.writeText`, temporary label swap via `[data-tc-copy-label]` for 2s.
- Image fallback JS: global singleton `tcImageFallbackAttached` handler — error event capture (true) on `[data-tc-img-fallback]`, swaps src to fallback and removes attribute. Uses capture phase because error events don't bubble.
- CountBadge: zero count hides badge (aria-hidden decorative), overflow shows "N+" (default max 99). `formatInt` helper is shared with CountBadge.
- RelativeTime: pure Go formatting (`formatRelativeTime`), no JS. `<time datetime>` for a11y/SEO, `title` for absolute time on hover.
- LoadMore: cursor appended as `?cursor=` query param (detects existing `?` for `&`). `hx-swap="outerHTML"` + `hx-target="this"` for self-replacement. `InfiniteScroll: true` adds `hx-trigger="revealed"`.
- EndOfList: `navigation.EndOfList(EndOfListProps)` — "You've reached the end" indicator for the bottom of a list. Companion to LoadMore/Pagination. `role="status"`, `text-gray-500 dark:text-gray-400`. Customizable `Message`.
- Card/Table composition: `CardPaddingNone` + `Table.Flush` together enable table-in-card (no double border). `Card.Body`/`Card.Header`/`Table.Body` are `templ.Component` slots that override default rendering when set.
- Table: row cells auto-padded/truncated to match header count. `TypedHeaders []TableHeader` takes precedence over `Headers []string` for sortable columns: each `TableHeader` has `Sortable bool`, `SortDirection` (`SortNone`/`SortAsc`/`SortDesc`), and `Href` for server-side sort links. Renders `aria-sort="ascending/descending/none"` + ↑/↓ indicators. `ariaSortValue()` maps the enum; `tableHeaderCount()` aligns cell padding to whichever header type is used. `TableRow.Href` makes rows clickable: sets `data-tc-row-href` + `role="link"` + `tabindex="0"` + `cursor-pointer`. A CSP-safe singleton script (`tableRowHrefJS` in `shared.go`) handles click + keyboard navigation. Clicks on interactive child elements (links, buttons) are not hijacked.
- **Error handler:** `errorpage/handler.go` provides `ErrorHandler(err, cfg)` returning `http.Handler`, `FromError(err)` for type-safe conversion from go-error-family errors, 6 pre-built constructors (`NotFound`, `Forbidden`, `BadRequest`, `Conflict`, `ServiceUnavailable`, `InternalError`), `WriteError`/`WriteErrorPage` convenience wrappers, `HTMLShell` mode for valid HTML documents, `JSON` mode for API/HTMX responses. Uses `errors.AsType[errorfamily.Classified]()` for go-error-family integration.
- **Error families:** `errorpage` package integrates with go-error-family via `FromErrorFamily()` converter + `ParseFamily()` for string-based lookup. `FromError()` extracts Why/Fix defaults from go-error-family's `Family.DefaultWhy()`/`DefaultFix()` methods. (`FamilyFromErrorFamily` is a deprecated alias, will be removed in v1.0.)
- **Error components:** `ErrorPage` (full-page), `NotFound404` (dedicated 404 with hero numeral + search + links), `ErrorDetail` (inline card), `ErrorAlert` (family-aware alert) in `errorpage/`
- **NotFound404:** `errorpage.NotFound404(props NotFound404Props)` — dedicated 404 page with large gradient numeral (`text-[8rem]`), optional search form (`SearchAction`), quick-links card grid (`[]NotFoundLink`), and "Go home" / "Go back" buttons. Unlike `ErrorPage` (family-colored error card), NotFound404 is a welcoming navigation aid using neutral blue/indigo palette. `DefaultNotFound404Props()` returns full defaults. `DefaultNotFoundLinks()` returns starter links (Home, Documentation). Types and constructors live in `notfound404_types.go`. All string constants (titles, messages, labels) defined as package-private `notFound404*` constants for goconst compliance. Shares the `tcGoBackAttached` singleton with `ErrorPage`.
- **Drawer:** `display.Drawer` — accessible side panel rendered as a native `<dialog>` with `data-side="left"`/`"right"`. CSS positions the dialog via `margin-inline-*` (auto-mirrors in RTL) and animates via `translateX`. Side positioning is in `templates/custom.css` under `dialog.tc-drawer[data-side=...]`.
- **ValidationSummary:** `forms.ValidationSummary` — accessible error summary with icon, error count, linked field errors, `role="alert"`.
- **Snapshot testing strategy (three tiers):**
  1. **HTML golden tests** (`utils/golden`) — fast, deterministic, the backbone. 102 golden files across all packages. Two normalizations: (a) CSS class tokens sorted alphabetically inside `class="..."`, (b) auto-generated EnsureID values (`tc-<prefix>-<16hex>`) replaced with `tc-<prefix>-NORMALIZED` so EnsureID-using components (Accordion, Tabs, Dropdown, Tooltip, Carousel, ContextMenu, Combobox, Nav) are safe for golden testing without explicit IDs. LCS-based diff with line numbers. Update with `-update` flag.
  2. **Substring assertions** (`utils.AssertContains`/`AssertNotContains`/`AssertContainsAll`) — lightweight checks for individual attributes/structure. Use for targeted invariant checks (e.g., "must contain `aria-label`"), not as a substitute for golden tests.
  3. **Visual regression** (`visualtest.AssertScreenshot`) — pixel-level PNG comparison in headless Chromium. Catches what string tests cannot: layout shifts, color regressions, RTL mirroring. Separate module. Run via `nix run .#visual`.
- **Table-driven golden tests:** Use `golden.AssertSnapshots(t, []golden.Snapshot{{name, html}, ...})` to batch multiple snapshots as parallel subtests. Eliminates `t.Run` + `golden.Assert` boilerplate. Pattern: `golden_sweep_test.go` files in each package cover all component variants.
- **Adding a new component's golden test:** Create a `TestGoldenSweepXxx` function in the package's `golden_sweep_test.go`, render with `utils.Render(t, Component(props))`, assert via `golden.AssertSnapshots`. Run `go test -run TestGoldenSweep -update ./pkg/...` to generate golden files. No need to pass explicit IDs — normalization handles EnsureID automatically.
- **Visual regression testing:** `visualtest.AssertScreenshot(t, name, component, opts...)` — renders a component in headless Chromium (chromedp) and diffs pixels against a committed PNG in `visualtest/testdata/`. Catches what HTML-string golden tests cannot: layout shifts, dark-mode color regressions, RTL mirroring. Supports `Dark`, `RTL`, `Viewport`, `MaxMismatch` (default 0.1%), and `-update`. Run via `nix run .#visual` (Nix provides Chromium; tests skip if no browser). Lives in a **separate Go module** (`visualtest/go.mod` with a local replace) so chromedp never pollutes the library's consumer dependency graph. See `docs/visual-testing.md`.
- **Error sub-templates:** 6 shared private sub-templates in `errorpage/shared.templ` (familyIcon, fixCard, causeList, contextTable, timestampFooter, familyBadge)
- HTMX retry: per-element `data-tc-retry` attribute (no shared counter)
- HTMX error handling: family-aware — when server returns structured JSON with `family` field, toast type is mapped. `ErrorHandlerConfig{JSON: true}` produces the JSON format that HTMX consumes.
- GlobalErrorHandling: configurable via `ErrorHandlingConfig` struct (MaxErrorHistory, MaxRetries, RetryDelayMS). Includes `tc-error-announcer` div with `aria-live="polite"` for screen reader announcements.

- Thread safety: `utils.Class()` uses `sync.Mutex` to protect tailwind-merge-go's shared LRU cache from concurrent access. Required even though the LRU has internal mutexes — they don't protect the full Merge() call sequence.
- CSP nonce test: `integration/csp_nonce_test.go` renders every inline-script component and asserts every `<script>` tag has `nonce=`. Prevents CSP regressions.
- `layout.Script(nonce, src, attrs)` / `layout.Stylesheet(href, attrs)` — CSP-safe external resource helpers that auto-inject nonce.
- `OverlayKind` typed enum (`OverlayModal`/`OverlayDrawer`) replaces untyped strings on overlay internals.
- Icons: 106 total. `icons.IconPathData` exported for consumers needing raw SVG paths (icons-only adoption without Tailwind).
- **Animated Icons:** `icons.AnimatedIcon(name, class)` and `icons.AnimatedIconWithAnimation(name, anim, class)` render any icon with a hover-triggered CSS animation, inspired by [heroicons-animated.com](https://www.heroicons-animated.com/). Pure CSS (zero JavaScript), respects `prefers-reduced-motion`, triggers on both `:hover` and `:focus-within`. 11 animation presets (`AnimPulse`, `AnimBeat`, `AnimBounce`, `AnimWiggle`, `AnimSpin`, `AnimJump`, `AnimNod`, `AnimShake`, `AnimBlink`, `AnimWobble`, `AnimDraw`). Every icon has an explicit default animation via `DefaultAnimation(name)` — no silent `AnimPulse` fallback; aliases (ArrowPath, Bars3, MapPin, HandThumbUp) resolve to their canonical icon's animation; Spinner defaults to `AnimNone` (has its own built-in spin). Per-path animations: `AnimBlink` (Eye) targets `svg path:nth-child(N)` via CSS and requires 2-path icons; `AnimDraw` (Bolt) renders paths with `pathLength="1"` and animates `stroke-dashoffset` for a self-draw effect; `AnimWobble` (Beaker) combines `scale: 0.9` + rotation oscillation. CSS lives in `templates/custom.css` under `.tc-anim-*` classes. `AnimatedIcon` wraps the SVG in `<span class="tc-anim tc-anim-{type} inline-flex">`; `AnimNone` renders plain `Icon()` without a wrapper; `AnimDraw` renders via a specialized `drawIcon` template that adds `pathLength="1"` to each `<path>`.
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
- **Tooltip is pure CSS (no JS).** The `:hover` / `:focus-within` classes show/hide the tooltip. Consumers must set `aria-describedby` directly on the focusable trigger element — the wrapper `<div>`'s `aria-describedby` is no longer auto-propagated to inner focusable children by JS. Touch devices do not toggle on tap (tooltips are progressive enhancement). See ADR-0017 for the migration rationale.
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
- **Carousel uses CSS scroll-snap + keyboard navigation**: native touch/drag support via `snap-x snap-mandatory scroll-smooth` on the track. Slides use `snap-center`. JS simplified to `scrollBy`/`scrollTo` for prev/next/dot navigation + `scrollend` for dot sync (with scroll+debounce fallback). Replaced the old `translateX` + manual transform approach. The carousel region is focusable (`tabindex="0"`) and responds to ArrowLeft/ArrowRight (prev/next), Home (first), End (last). RTL-aware via `document.documentElement.getAttribute('dir')`. The region shows a visible `focus-visible:ring-2` outline so keyboard users see where focus landed.
- **ViewTransitions API**: `htmx.ViewTransitions(ViewTransitionsProps{Global: true})` enables native View Transitions for HTMX swaps via `htmx.config.globalViewTransitions = true` (HTMX 2.0 built-in). Renders default cross-fade CSS with `prefers-reduced-motion` support. Graceful degradation — browsers without View Transitions do instant swaps.
- **Modern browser capabilities**: see `docs/research/modern-browser-capabilities.md` for the comprehensive analysis of native APIs (`<dialog>`, Popover API, `@starting-style`, `<details>`, scroll-snap, View Transitions, `content-visibility`, `:has()`) and the phased migration roadmap. CSS foundation for these APIs is in `templates/custom.css`.
- **Modal/Drawer use native `<dialog>`**: `showModal()`/`close()` replace ~200 lines of custom focus-trap JS. Browser handles focus trap, Tab cycling, Escape, focus restore, top-layer, backdrop, inert. CSS `.tc-overlay`/`.tc-modal`/`.tc-drawer` classes handle animations via `@starting-style` + `allow-discrete`. JS wrappers `tcOpenOverlay(id)`/`tcCloseOverlay(id)` kept for HTMX compat. Backdrop click detection: `e.target === dialog` (the `::backdrop` pseudo-element registers as the dialog).
- **Dropdown/Popover/ContextMenu use native Popover API (ADR-0017)**: `popover="auto"` + declarative `popovertarget` on the trigger button replaces the old show/hide/state singleton JS. Native light-dismiss (click-outside + Escape), top-layer rendering, focus restore. **Top-layer positioning:** the UA stylesheet forces `position: fixed; inset: 0` on `[popover]`, detaching the panel from its trigger's DOM subtree, so CSS classes like `top-full` resolve against the viewport, not the trigger. A shared singleton `popoverPositionJS` (in `display/shared.go`) reads `getBoundingClientRect()` on `toggle` open and sets `style.left/top` with viewport clamping. Mirrors the proven ContextMenu cursor-positioning pattern. `Dropdown` and `ContextMenu` share a single `menuKeyboardNavScriptComponent` (in `display/shared.go`, singleton `tcMenuKeyNavAttached`) for WAI-ARIA menu keyboard nav: ArrowUp/Down (with RTL-aware Left/Right), Home, End, PageUp/PageDown, first-menuitem focus on `toggle` open, skipping both `[disabled]` and `[aria-disabled="true"]` items. Enter/Space activation is intentionally left to native browser behavior so HTMX-powered links/buttons work correctly — never use `window.location.href` in a custom handler (regression guard: `TestDropdownKeyboardEnhancements` asserts its absence). `ContextMenu` is fully keyboard-operable: Shift+F10 and the dedicated ContextMenu key open it (positioned at the trigger via `getBoundingClientRect()`), then the shared menu nav takes over; the `contextmenu` mouse event still positions at the cursor. `Tooltip` uses pure CSS for show/hide and a small singleton (`tooltipAriaJS`) to propagate `aria-describedby` to the focusable trigger. Escape dismisses an active tooltip via `data-tc-tooltip-dismissed` attribute + CSS rule (keeps focus on trigger, preserving tab position). Hover or re-focus clears the dismissed state. `[popover]::backdrop { background-color: transparent; }` in `templates/custom.css` ensures popovers don't dim the page (menus aren't modals).
- **Stylable Select (`appearance: base-select`)**: `SelectProps.Stylable: true` opts into the customizable `<select>` API. Emits `<button><selectedcontent></selectedcontent></button>` inside the select. CSS `.tc-select` in `templates/custom.css` styles button, picker (`::picker(select)`), options, arrow (`::picker-icon`), checkmark. Progressive enhancement — non-supporting browsers (Firefox, iOS Safari) ignore the structure and render native `<select>`.
- **Textarea AutoGrow**: `TextareaProps.AutoGrow` (default `true`) uses CSS `field-sizing: content` via `.tc-auto-grow` class. No JavaScript. `field-sizing` is Baseline 2024.
- **EnterKeyHint (unified API)**: Both `InputProps.EnterKeyHint` and `TextareaProps.EnterKeyHint` use the same typed `EnterKeyHintType` enum. Constants: `EnterKeyHintSend`, `EnterKeyHintDone`, `EnterKeyHintGo`, `EnterKeyHintNext`, `EnterKeyHintPrevious`, `EnterKeyHintSearch`, `EnterKeyHintEnter`. `EnterKeyHintTypeIsValid` included. Input also auto-derives a smart default from `InputType` (email→next, search→search, etc.) via `enterKeyHintValue()`; explicit `EnterKeyHint` overrides the auto-derived value.
- **Input search semantic landmark**: `Input` with `Type: InputSearch` auto-wraps in `<search>` element (Baseline 2023). Screen readers announce it as a search landmark. No API change — auto-detected from the InputType.
- **Form hx-validate**: `FormProps.Validate: true` emits `hx-validate="true"` for HTML5 constraint validation before HTMX submit. Pair with native `required`, `pattern`, `type="email"` etc. for client-side validation.
- **Image responsive delivery**: `ImageProps.SrcSet` and `Sizes` are typed string fields — no more `Attrs` workaround. Example: `SrcSet: "/img-480w.jpg 480w, /img-800w.jpg 800w", Sizes: "(max-width: 600px) 480px, 800px"`.
- **Table content-visibility**: `TableProps.LazyRows: true` applies `content-visibility: auto` to body rows via `.tc-content-auto` class (48px intrinsic height). When `CellPadding: TableCellPaddingCompact`, uses `.tc-content-auto-compact` (40px intrinsic height) to avoid scrollbar jitter. Browser skips rendering off-screen rows. Recommended for tables with 100+ rows.
- **MobileMenu keyboard support**: Escape closes the menu and returns focus to the toggle button. Opening moves focus to the first focusable child. Shared `tcMobileMenuSet(menu, btn, open)` helper handles visibility, icon swap, `aria-expanded`, and focus management for both click and keyboard paths.
- **Global accent-color CSS**: `templates/custom.css` sets `accent-color: blue-600` (light) / `blue-400` (dark) on checkboxes, radios, range inputs, and progress bars. Consumers override via `@theme { --color-blue-600: #custom; }`. No Go code changes needed — same theming model as all other components.

## Demo Infrastructure

- **Demo CSS**: `examples/demo/demo.css` compiles Tailwind CSS scanning ALL `.templ` files in the repo (`@source "../../**/*.templ"`). Path is relative to the CSS file location (`examples/demo/`), so `../../` reaches the repo root. The compiled CSS is embedded via `//go:embed static/app.css` and served at `/css/app.css`.
- **@source path gotcha**: Tailwind v4's `@source` directive resolves relative to the CSS file, not the CWD. Using `@source "./**/*.templ"` from `examples/demo/demo.css` only scans `examples/demo/*.templ` — missing ALL component classes from `display/`, `forms/`, `feedback/`, etc. Must use `@source "../../**/*.templ"` to scan the entire repo.
- **Custom CSS**: Component-specific CSS (dialog animations, stylable select, auto-grow textarea, scroll-snap, accordion chevron, accent-color) lives in `templates/custom.css`. Both `templates/app.css` (consumer template) and `examples/demo/demo.css` (demo) import it via `@import "./custom.css"` / `@import "../../templates/custom.css"`. Single source of truth.
- **Dockerfile 3-stage pipeline**: CSS (Node 22, compiles Tailwind) → Go binary (templ generate + go build, overwrites committed CSS with fresh compile from Stage 1) → Distroless runtime. CSS is always freshly compiled during Docker build — the committed `static/app.css` is never stale because it's overwritten. `.dockerignore` excludes `.git/`, `website/`, `docs/` from Docker context (~653MB → ~15MB).
- **Demo endpoints**: `/health` returns `{"status":"ok"}` for Cloud Run health checks. `/css/app.css` serves embedded CSS with `Cache-Control: public, max-age=31536000, immutable`. `/api/load-more` and `/api/delete` are mock HTMX endpoints for interactive demo components.

## CI & Tooling Gotchas

- **pnpm 11 build-script approval is load-bearing in 2 config files.** pnpm 11 hard-fails (`ERR_PNPM_IGNORED_BUILDS`) any install whose dependency ships a build script that was not explicitly approved. The `pnpm` field in `package.json` is dead (warned + ignored); approvals live in `pnpm-workspace.yaml` under `allowBuilds: <pkg>: true`. Current sites: `website/pnpm-workspace.yaml` (`esbuild`, `sharp`) and the demo Dockerfile CSS stage (`@parcel/watcher`, written inline via `printf` before `pnpm add`). `pnpm init` auto-writes `devEngines.packageManager: "^11.x"`, so pnpm self-resolves to the latest 11.x at run time — the Docker CSS stage tracks pnpm 11.x minors. `pnpm dlx` inherits the same policy.
- **Never `pnpm add -g` on GitHub runners.** It fails preflight ("configured global bin directory ... is not in PATH") because runners do not set `PNPM_HOME` and corepack does not either (verified by container replication on node:24 + corepack pnpm@11.20.0). Use `npm install -g <cli>` for one-off global tools — the deploy job does this for firebase-tools.
- **Docker builder + `replace` directives:** root `go.mod` replaces all 6 sub-modules with local paths, so `go mod download` needs every sub-module's `go.mod` copied into the stage BEFORE it runs (see `examples/demo/Dockerfile`). Copying only root manifests fails with `reading charts/echarts/go.mod: no such file or directory`.
- **Visual tests need the pure fontconfig pin.** `makeFontsConf` is impure by design (includes `/etc/fonts/conf.d`, `/usr/share/fonts`, profile fonts) — never use it for cross-machine determinism. The `#visual` flake app uses a hand-written `fonts.conf` (`pkgs.writeText`) with ONLY `<dir>` entries for `pkgs.inter` + `pkgs.dejavu_fonts` and a tmp cachedir. If goldens flake: `rm -rf /tmp/tc-visualtest-fontconfig-cache` (the cache goes stale) and read the fc-match diagnostics the app echoes at startup.
- **upload-artifact v4 hides dotfiles by default.** Visual-regression failure screenshots live under `testdata/.fail/` — without `include-hidden-files: true` the artifact uploads empty and the failure evidence is lost.
- **BuildFlow `eslint-fix` breaks commits touching any `.ts`/`.js` file** (TODO #108): ESLint 10 runs at the repo root, which has no eslint config → exit 2. When a changeset includes website JS/TS, run the full verify matrix manually and commit with `--no-verify`, noting why in the body.
- **Daemon commits can regress same-day fixes.** During the v1.9.0 cut, daemon auto-commits prettier-un-minified `examples/demo/static/app.css` (+4780 lines → CSS Freshness CI failure) and flipped `website/package.json` typescript back to `^7.0.2` (astro check crashes on TS 7 — needs 6.x until `astro check` supports the native compiler). After any daemon commit lands, re-check: `nix run .#css` byte-stability, the website typescript pin, and CI status. The daemon also pushes master (and once, the release tags) without being asked — always `git fetch` before assuming local-only state.

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

**Tag format:** annotated + SSH-signed, message `<version>: <one-line summary>` (same key as v0.5.0).

**Post-release commits** (backfilling tests, doc fixes, post-release regeneration) land
normally on `master` and roll into the next release. Never retag the same version.

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
4. Installs an `EXIT`-trap rollback (`release_cleanup`) that restores `utils/version.go`, all `go.mod` files, `CHANGELOG.md`, and `FEATURES.md` if any later step fails — so a failed verify never leaves a dirty tree
5. Bumps `utils.Version` via in-place sed
6. Moves the `[Unreleased]` body under a new `## [<version>] — YYYY-MM-DD` heading (inserts a fresh empty `[Unreleased]` above)
7. Bumps `FEATURES.md` `**Version:**` + `**Updated:**` date (the three version files must move together; `utils.TestVersionMatchesFeatures` enforces it)
8. Regenerates `*_templ.go` and runs the full verify suite (build + test + lint)
9. Asserts the version drift-guard (`TestVersionMatches(Changelog|Features)`)
10. Strips the local `replace` directives (tagged `go.mod` files must be consumer-clean) and re-parses every `go.mod` as a sanity check
11. Stages and commits as `release: <version> — <summary>` (one-commit convention; body carries the release notes, `Assisted-by: Crush:${CRUSH_MODEL}`)
12. Creates annotated, SSH-signed tags: root `v<version>` plus one `<sub-module>/v<version>` per published sub-module, in lockstep. Guard with `scripts/check-release-tags.sh` before pushing — a root-only release (like v1.8.3) breaks every consumer.

The script does **not** push. House rule: "NEVER PUSH TO REMOTE". Push manually
after reviewing the release commit and tag with `git show v<version>` and
`git show <commit>`.

**Verify-before-strip (v1.9.0 lesson):** the script's build/test/lint phase MUST
run while the local `replace` directives are still present. go1.26.5 workspace
mode does NOT preempt module-graph resolution of `require` entries at unpushed
versions — with replaces stripped and the new tags not yet on the proxy, every
build fails with `unknown revision <sub>/v<version>` (GOPRIVATE is not a factor;
proven with a 2x2 replaces/GOPRIVATE matrix during the v1.9.0 cut). Stripping
now happens after verification. Also fixed then: bash keeps only ONE `EXIT`
trap, so the script's second `trap` silently disabled the rollback trap — both
cleanups now share one hook.

## Lint Command

```bash
# Multi-module: golangci-lint does not support go.work workspace mode — lint per module.
# Root module (examples/ excluded via .golangci.yml paths exclusion):
golangci-lint run ./display/... ./feedback/... ./forms/... ./integration/... ./internal/... ./layout/... ./navigation/... ./recipes/... ./cmd/...
# Sub-modules:
for mod in utils icons errorpage charts/echarts htmx datastar; do (cd "$mod" && golangci-lint run ./...); done
```

**Disabled linters (do NOT re-enable — fundamentally incompatible with this codebase):**

- `ireturn` — every component returns `templ.Component` (an interface) by design; the linter's premise is antithetical to templ.
- `godoclint` — demands exactly one `// Package` godoc per package, but the repo intentionally documents per-file.
- `testableexamples` — `Example*` funcs render HTML that is verbose and version-dependent; output isn't asserted.

**Reconciled at v0.18.1:** commit 73395d9 expanded to 67 linters but left `golangci-lint run` failing (187 findings, CI would have gone red on first push). The config now uses an explicit depguard allow-list (the `$module` token did not resolve — use literal `github.com/larsartmann/templ-components` + the three runtime deps), extends `varnamelen`/`mnd` ignore lists, and excludes test files from `err113`/`makezero`/`varnamelen`/`gocheckcompilerdirectives`. If you add a linter, run `golangci-lint run` to 0 findings before committing.

**`.golangci.yml` regression root cause (identified 2026-07-28, T1):** The 3 disabled linters re-entered the enable list **5 times** across sessions. Root cause: the AI agent's working tree holds a stale `.golangci.yml`, and broad `git add` for docs work silently stages the stale file. BuildFlow's pre-commit hook has a 60s budget and does NOT run `go test ./...`, so `TestGolangciDisabledLinters` only fires in CI, not at commit time. **Prevention (3 layers):** (1) `TestGolangciDisabledLinters` in `utils/lint_config_test.go` catches in CI via `go test ./...`; (2) `scripts/check-lint-config.sh` is a <50ms standalone grep guard wired into `.git/hooks/pre-commit` BEFORE BuildFlow runs; (3) CI step "Lint-config guard" in `.github/workflows/ci.yaml` runs the script before `golangci-lint` even installs. If you must modify `.golangci.yml`, run `scripts/check-lint-config.sh` to verify.

## `encoding/json/v2` Adoption

This library uses `encoding/json/v2` + `encoding/json/jsontext` (Go 1.26+ with
`GOEXPERIMENT=jsonv2`). The pre-commit hook (`scripts/pre-commit.sh`) sets
`GOEXPERIMENT=jsonv2` automatically. The `.golangci.yml` enables the
`goexperiment.jsonv2` build tag. The `flake.nix` devShell exports
`GOEXPERIMENT=jsonv2` via `shellHook`. **`.envrc`** (direnv) sets
`GOEXPERIMENT=jsonv2` repo-wide for ALL tools. `go.work` is active by default (lists all 7 modules). Use `GOWORK=off` for per-module isolation testing.

**Consumers** must set `GOEXPERIMENT=jsonv2` when building (or wait for Go 1.27
where it becomes stable). The `errorpage` package uses `json.MarshalEncode` +
`jsontext.NewEncoder` for JSON error responses. The `navigation/breadcrumbs`
package also uses `encoding/json/v2`. Remaining packages (tests) still use
`encoding/json` v1 — both coexist fine under the experiment flag.

## Conventions

- **Naming hygiene:** `forms/radio_go.go` renamed to `forms/radio.go` (the `_go.go` suffix falsely implied generated code). `icons.Close` added as alias for `icons.X` (prefer `Close` in new code — `X` is a single-letter identifier with poor discoverability). `errMsg` → `errorMessage` (no abbreviations). `cleanMessage` → `sanitizeErrorMessage` (precise verb). `htmxMainSRIDefault` → `sriHTMXMainDefault` (consistent word order with `sriHTMXMainByVersion`).
- **RTL keyboard mapping:** `display.Tabs`, the shared menu-keyboard nav (`display.Dropdown` + `display.ContextMenu`), and `display.Carousel` JS handlers check `document.documentElement.getAttribute('dir') === 'rtl'` and swap ArrowLeft/Right mappings per WAI-ARIA APG. In LTR: ArrowRight=next, ArrowLeft=prev. In RTL: ArrowLeft=next, ArrowRight=prev.
- **Rating DOM order must stay forward (1→N):** the interactive `forms.Rating` renders radio inputs in forward DOM order (value 1 first) so radiogroup arrow keys (ArrowDown/Right) increase the value per WAI-ARIA. The container uses `flex-row-reverse` (interactive branch only — read-only uses normal flow) so the `peer-checked` `~` fill still renders the correct ★★★☆☆ visual. Critically, the fill/color/hover/focus classes live on the `<label>` (a sibling of the hidden `.peer` radio), NOT the nested `<svg>` — Tailwind's `peer-checked` `~` combinator only matches siblings, so classes on the nested SVG never resolve. The SVG inherits color via `currentColor`. Never revert to reverse DOM order (it breaks arrow-key direction) or move fill classes back to the SVG (it silently breaks the fill).
- **Sub-template extraction pattern (ADR 0010):** Extract when 2+ callers, 5+ lines, clear domain name. Do NOT extract for single caller, demo code, no clean name, 8+ parameters, or same-file callers. See `docs/adr/0010-sub-template-extraction-pattern.md`.
- **Motion-reduce compliance test:** `utils.TestMotionReduceCompliance` scans all `.templ` files for `transition-*`/`animate-*` classes without `motion-reduce:` fallback. Run via `go test ./utils/... -run TestMotionReduce`.
- **SKILL.md drift-guard:** `utils.TestSkillComponentCount` logs actual vs documented component count. Informational (not failing) — intended to surface drift during code review.
- **Fuzz tests:** `forms.FuzzInputType`, `forms.FuzzFormMethod`, `display.FuzzButtonHTMLType` verify enum validation never panics on arbitrary input. Run via `go test -fuzz=. -run=Fuzz ./...`.
- **Benchmark suites:** Now in 7 packages (display, feedback, navigation, forms, layout, htmx, icons, utils). Run via `go test -bench=. -benchmem ./...`.
- **goconst zero issues:** all repeated string literals are named constants — keep it that way.
- **CSRFTokenName:** `forms.FormProps` has a `CSRFTokenName` field (defaults to `"csrf_token"`) for framework compatibility.
- **ErrorPage/NotFound404 landmark:** Both use `<main>` (not `<div role="region">`) for WCAG 2.4.1 Bypass Blocks compliance.
- **FromError fallback:** Unknown errors return `FamilyCorruption` (→500), not `FamilyInfrastructure` (→503). An unrecognized error is a bug, not a transient outage.
- **Native SVG Charts (Tier 1):** `display.LineChart`, `display.AreaChart`, `display.PieChart` — pure server-side SVG, zero JavaScript. All compose shared geometry primitives from `chart_geometry.go`: `ScalePoints`, `BuildPolylinePath`, `BuildSmoothPath` (Catmull-Rom spline), `BuildAreaPath`, `ComputeNiceTicks`, `FormatTickValue`. The shared chart color palette constants (`chartColorBlue`, `chartColorEmerald`, etc.) live in `chart_geometry.go` — all chart types use the same palette. LineChart and AreaChart share the `LineChartSeries` type and `LineChartStyle` enum; AreaChart adds `FillOpacity` and a filled area path. PieChart uses SVG arc paths (`computeArcPath`) — not stroke-dasharray — so it supports donut holes, label positioning, and non-integer totals. Dark mode via Tailwind `text-*`/`stroke-*`/`fill-*` classes with `dark:` variants on SVG elements; `currentColor` inheritance means a single `text-blue-600 dark:text-blue-400` class controls both stroke and fill. Accessibility: `role="img"` + consumer-provided `AriaLabel`; without a label, `aria-hidden="true"` (decorative). When adding a new chart type, compose the geometry helpers — don't reimplement math.
- **ECharts Adapter (Tier 2):** `charts/echarts` package — opt-in wrapper for interactive charts. Follows the `datastar` precedent: does NOT import go-echarts. Consumer builds their chart with go-echarts, calls `RenderSnippet()`, and passes `Element` + `Script` strings to `EChartsProps`. The `EChart` component wraps output in a CSP-safe `<div>` + inline `<script nonce>`. The `SDKScript` component loads the ECharts runtime from CDN (configurable version + host, self-hostable via `Src`). Dark mode bridge (`darkModeBridgeJS`) uses a `MutationObserver` on `document.documentElement.class` to sync ECharts theme with the Tailwind `.dark` class — singleton guard `window.tcEChartsDarkBridge`. The bridge does a shallow `setOption({...}, {merge: true})` — consumer custom axis styling may be partially overridden (documented limitation). The `chartScriptComponent` helper writes raw JS directly to the buffer via `templ.ComponentFunc` — required because templ's `<script>` context sanitizes `{ }` interpolation. See ADR-0031 for the two-tier decision.
