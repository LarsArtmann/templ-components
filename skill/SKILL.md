---
name: templ-components
description: Authoring and adoption playbook for the templ-components Go UI library (templ + HTMX + Tailwind v4). Use this skill whenever working in the github.com/larsartmann/templ-components repo, or when adding, editing, reviewing, or debugging any templ component, props struct, typed enum, style lookup, icon path, wire Action, HTMX/Datastar helper, guard test, or component count here. Also trigger when the user asks to "add a component", "fix a component", extend the library, wire up a new icon, add a typed enum, set up Tailwind v4 scanning for templ, integrate HTMX loading/error handling, or asks how a component should be structured, themed, made accessible, CSP-safe, or dark-mode aware in this codebase. Load it BEFORE reading or editing any .templ file, SKILL.md, or count-bearing doc in this repo.
metadata:
  tags: templ, templ-components, htmx, tailwind, tailwind-v4, go, ui, components, accessibility, csp, dark-mode, server-rendered
---

# templ-components

A server-rendered, type-safe, accessible, CSP-ready Go component library built on
[templ](https://templ.guide), [HTMX](https://htmx.org), and Tailwind v4.

Two audiences, one repo:

- **Consumers** (most readers): what exists, how to adopt it, how to wire it. → **Part 1**
  plus `references/component-catalogue.md` for exact signatures.
- **Authors** (maintainers): how to add, change, and review a component without breaking the
  library's contracts. → **Part 2**, then `references/guard-tests.md` before touching CSS,
  enums, generated files, or docs.

Bundled references (read on demand, not upfront):

| File | Use it for |
| --- | --- |
| `references/component-catalogue.md` | Full per-package signature tables, icons, `utils`/`utils/format`/`utils/wire` APIs |
| `references/consumer-guide.md` | Quick start, head-tag suppression, dual-transport form, theming, compiled-CSS rule, dark-mode checklist, recipes, adoption tracking |
| `references/authoring-guide.md` | Shared-helper inventory, full per-component test checklist, decision trees, website SSG registration |
| `references/guard-tests.md` | Every repo-wide guard test, what it prevents, and where it lives |

---

# Part 1: Consumer Guide

## Component catalogue

136 components across 11 packages (nine counted templ packages, plus `utils` and `icons`)
+ 102 icons. Full signatures and one-liners live in `references/component-catalogue.md`; this
section is the index. If you're about to hand-roll something, check here first.

### By use case (start here)

| You're building...             | Reach for                                                                                                                                                                                                                                                                                                                                                                               |
| ------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Dashboard / metrics page**   | `AppShell`, `Container`, `Grid`, `StatCard`, `Card`, `ProgressBar`, `CircularProgress`, `SkeletonCardGrid`, `PageHeader`, `Sparkline`, `BarChart`, `LineChart`, `PieChart`, `AreaChart`, `PolledRegion`                                                                                                                                                                                  |
| **Kanban / pipeline board**    | `KanbanBoard` — HTML5 drag-and-drop + keyboard move buttons, `Wire` move exchange (htmx or Datastar), `ParseKanbanMove` server helper, per-column `Action` add-card slot and `Tone` status dot                                                                                                                                                                                          |
| **List / table page**          | `Table` (`Flush` for card nesting, `CellPadding` for compact rows, `Table.Body` for custom rows, `Row.Href` for clickable rows), `Badge`, `StatusBadge`/`StatusBadgeWith`, `Avatar`, `Pagination`, `LoadMore`, `EndOfList`, `EmptyState`, `DataState` (disabled→unavailable→empty→content ladder), `ListNote`                                                                             |
| **Detail page**                | `Split`, `Card`, `DefinitionList`, `DefinitionGrid`, `Tabs`, `PageHeader`, `SectionHeading`, `DateRange`, `Breadcrumbs`                                                                                                                                                                                                                                                                |
| **Settings / data-entry form** | `Form`, `Input`, `Select` (supports `Groups` for optgroups), `Textarea`, `Toggle`, `Checkbox`, `RadioGroup`, `ValidationSummary`, `DirtyGuard`                                                                                                                                                                                                                                         |
| **Filter bar (horizontal)**    | `FilterBar` (auto-submit GET form, composite change trigger, Reset, noscript Apply, sticky opt-in, `Wire`) — children are plain controls; recipe `docs/recipes/horizontal-filter-bar.md` documents the footguns; `FilterChips` (zero-JS link-based chips, `aria-current`, `FilterToggleHref` query builder); `SegmentedControl` (segmented single-select — radio or link mode, zero JS) |
| **Feedback / notifications**   | `Toast`, `ToastContainer`, `Alert`, `Spinner`, `ProgressBar`, `CircularProgress`, `GlobalErrorHandling`                                                                                                                                                                                                                                                                                 |
| **Navigation**                 | `Nav`, `SimpleNav`, `SidebarNav`, `Breadcrumbs`, `Pagination`, `MobileMenu`                                                                                                                                                                                                                                                                                                             |
| **Modal / overlay**            | `Modal`, `Drawer`, `Dropdown`, `Tooltip`, `Popover`, `Accordion`                                                                                                                                                                                                                                                                                                                        |
| **Error pages**                | `ErrorPage`, `NotFound404`, `ErrorDetail`, `ErrorAlert`, `ErrorHandler`                                                                                                                                                                                                                                                                                                                 |
| **Full page shell**            | `Base`, `Minimal`, `ThemeScript`, `ThemeToggle`, `Script`, `AppShell` (admin shell), `Container` (content width), `Split` (main+aside), `Stack` (vertical rhythm)                                                                                                                                                                                                                       |

### By package

#### `display` — 52 components

Accordion, AreaChart, Avatar, Badge, BarChart, Button, Card, Carousel, CodeBlock,
CollapsibleSection, CommandPalette, ContextMenu, CopyButton, CountBadge, DataState,
DataTable, DateRange, DefinitionGrid, DefinitionList, Drawer, Dropdown, EmptyState,
ExternalLink, Eyebrow, Grid, Heatmap, HoverCard, Image, KanbanBoard, Lightbox, LineChart,
ListNote, LivePill, Modal, PageHeader, PieChart, Popover, Prose, RelativeTime, Scrollback,
SectionHeading, SegmentBar, SimpleCard, SimpleEmptyState, Sparkline, StatCard, StatusBadge,
StatusBadgeWith, StatusDot, Table, Tabs, Tooltip.

#### `forms` — 26 components

Calendar, Checkbox, Combobox, DatePicker, DirtyGuard, FieldError, FileInput, FilterBar,
FilterChips, FilterDropdown, FilterInput, Form, FormFieldWrapper, Input, InputGroup, Label,
Radio, RadioGroup, Rating, SegmentedControl, Select, Slider, TagsInput, Textarea, Toggle,
ValidationSummary.

#### `feedback` — 14 components

Alert, CircularProgress, InlineError, InlineLoading, InlineSuccess, LoadingOverlay,
ProgressBar, Skeleton, SkeletonGroup, SkeletonCardGrid, Spinner, StepIndicator, Toast,
ToastContainer.

#### `layout` — 11 components

AppShell, Base, Container, MetaRefresh, Minimal, Script, Split, Stack, Stylesheet,
ThemeScript, ThemeToggle.

#### `navigation` — 12 components

Breadcrumbs, EndOfList, Footer, LoadMore, MobileMenu, MobileMenuToggle, MobileNavLink, Nav,
NavLink, Pagination, SidebarNav, SimpleNav.

#### `htmx` — 9 components

CSRFToken, ConfirmDelete, GlobalErrorHandling, InlineLoadingOverlay, LoadingButton,
LoadingIndicator, PolledRegion, SwapOOB, ViewTransitions.

#### `datastar` — 6 components

Indicator, LiveRegion, LoadingButton, PolledRegion, SDKScript, SSEErrorHandling.

#### `errorpage` — 4 components + 6 constructors + handler

ErrorAlert, ErrorDetail, ErrorPage, NotFound404 — plus `ErrorHandler`, `WriteError`,
`WriteNotFound404`, `FromError`, and the `NotFound`…`InternalError` constructors.

#### `charts/echarts` — 2 components (opt-in module)

EChart, SDKScript. No `go-echarts` import: you pass in a `RenderSnippet()` element and the
module supplies the CSP-safe wrapper plus the dark-mode bridge (ADR-0031, tier 2).

#### `icons` — 102 icons

`Icon`, `IconRTL`, `IconWithStrokeWidth`, `Render`, `IconPathData`, `IconPathJS`,
`AnimatedIcon` (+ RTL/animation variants), `DefaultAnimation`, `AllAnimations`,
`AllIconNames`, `NameIsValid`.

#### `utils` — shared helpers

`BaseProps` (embed in every props struct), `ComponentProps`, `Class` (tailwind-merge),
`Lookup`, `Ternary`, `EnsureID`, `ScriptAttrs`, `ScriptComponent`, `Version`. Sub-packages:
`utils/format` (`Bytes`, `CompactDuration`, `ClockDuration`, `Percent`, `StringOrDash`,
`CompactCount`), `utils/wire` (transport-agnostic wiring, ADR-0036), `utils/svg`,
`utils/tone`, `utils/cdn`, `utils/golden`.

## Consumer essentials

```templ
@layout.Base(layout.DefaultPageProps()) {
    @navigation.SimpleNav(navigation.SimpleNavProps{
        BrandText:   "MyApp",
        Links:       []navigation.NavLinkProps{{Href: "/", Text: "Home"}},
        RightItems:  layout.ThemeToggle("Toggle theme", ""),
        CurrentPath: "/",
    })
    <main class="mx-auto max-w-7xl px-4 py-8">
        @display.Grid(display.GridProps{Cols: display.GridCols3}) {
            for _, u := range users {
                @display.StatCard(display.StatCardProps{Label: u.Name, Value: u.Count, Href: "/users/" + u.ID})
            }
        }
    </main>
    @feedback.ToastContainer("")
    @htmx.GlobalErrorHandling(htmx.DefaultErrorHandlingConfig())
}
```

- `RightItems` is a `templ.Component` field: pass the function value, never `@`-prefixed.
- One `wire.Action` drives HTMX or Datastar. Datastar patch target/mode come ONLY from
  response headers (`wire.Handler`); `Action.Target`/`Swap` render for htmx only. Under
  Datastar, form submission sets `ContentTypeForm` so the form's fields serialize.
- Inline scripts render through `utils.ScriptComponent`; external scripts through
  `layout.Script`. An empty nonce omits the attribute (never `nonce=""`).
- Theme by overriding Tailwind variables in `@theme`; never add color props to Go.
- Consumers gitignore `*_templ.go` and generate at build time (the library commits them).

Deep dive: `references/consumer-guide.md` — head-tag suppression, dual-transport form,
theming, compiled-CSS rule, dark-mode checklist, recipes, and adoption tracking.

# Part 2: Authoring Playbook

## Principles (the why behind every rule)

- **Make invalid states unrepresentable.** Variants are typed enums, lookups are maps with
  explicit fallbacks; the only permitted panic is the `icons` path-integrity check. This is a
  _library_: a panic in a consumer's render path is a bug we shipped.
- **Verify integrations against the counterparty artifact, not its docs.** SSE wire format,
  event names, HTMX attributes, CSP requirements: the ground truth is the pinned counterparty
  source (e.g. the `go-datastar/static` bundle). The 2026-08 SSE audit found the docs wrong in
  three places; string-asserting our own output cannot catch that class. Pin the contract in a
  guard test, and degrade empty/unknown values gracefully.
- **HATEOAS-first ([htmx.org/essays/hateoas](https://htmx.org/essays/hateoas/)).** HTML is
  the source of truth. Prefer native HTML (`<details>`, forms, links) over scripts; JS reads
  state FROM HTML attributes and is singleton-guarded so HTMX re-renders stay idempotent.
- **Accessibility is not optional.** ARIA, roles, keyboard nav, `motion-reduce:*`, and
  screen-reader text are part of "done".
- **CSP-safe by construction.** No `eval()`, no inline event handlers, no `javascript:` URLs;
  every inline script renders through `utils.ScriptComponent`/`layout.Script`.
- **Single source of truth.** SVG path data in `utils/svg` + `icons/icon_paths.go`, feedback
  styles in `feedback/styles.go`, shared class constants (`cardShellClass`, `mutedTextClass`).
- **Composition over configuration.** Props embed `utils.BaseProps`; consumers override via
  `@theme` CSS variables, never by editing Go. Slots are `templ.Component` children.

## Multi-module workspace (ADR-0034)

Seven published Go modules: root, utils, icons, errorpage, charts/echarts, htmx, datastar.
Local dev uses `go.work` (gitignored); CI and consumers use `replace` directives. Two
repo-local modules (`visualtest`, `website`) are never published.

**Module DAG (strict acyclic):**

```
Layer 0: utils (leaf)
Layer 1: icons, charts/echarts, htmx, datastar   [depend on utils]
Layer 2: errorpage                               [depends on utils, icons]
Layer 3: root                                    [depends on all above]
```

### Build & test commands (the complete forms)

```bash
# Regenerate + build everything (templ pinned by the dev shell)
find . -name '*_templ.go' -print0 | xargs -0 rm
templ generate ./...        # submodules: run from INSIDE each module dir, never a subdir pattern
go build ./...

# COMPLETE test run — the per-module loop. Root `go test ./...` does NOT include
# submodules even with go.work active (verified: `go list ./...` from the root
# expands to 10 root packages and zero submodule packages).
for mod in utils icons errorpage charts/echarts datastar htmx; do
  (cd "$mod" && GOWORK=off go test ./...)
done
(cd visualtest && GOWORK=off go test ./...)   # compile check; `nix run .#visual` for real runs

# Pre-push CI reproduction (read the VERDICT line, then push immediately)
scripts/ci-repro.sh --lint --css --website
```

`nix run .#test` runs `go test ./... -count=1 -race` from the repo root: fast and useful, but
it covers the root module only — never treat it as the full matrix. `nix run .#verify` IS the
complete lane (root test + per-module `GOWORK=off` loop + visualtest compile + per-module lint).

**Per-module lint** (golangci-lint does not support `go.work`):

```bash
golangci-lint run ./display/... ./feedback/... ./forms/... ./integration/... ./internal/... \
  ./layout/... ./navigation/... ./recipes/... ./cmd/...
for mod in utils icons errorpage charts/echarts datastar htmx; do
  (cd "$mod" && golangci-lint run ./...)
done
```

**Drift-guard tests** live in `utils`: `TestVersionMatches*`, `TestDocsCountDrift`,
`TestTemplVersionPin`, `TestGoDirectiveSkew`, `TestFeaturesEnumTableExhaustive`,
`TestFeaturesEnumValuesExhaustive`. Run them:
`cd utils && GOWORK=off go test ./... -run 'TestVersionMatches|TestDocsCountDrift|TestFeaturesEnum' -count=1`.

**Guard scripts:** the tracked `.githooks/pre-commit` (activated once per clone by
`scripts/setup-hooks.sh`, which sets `git config core.hooksPath .githooks`; `buildflow
precommit install` writes to the dead `.git/hooks/` and must not be used). It runs the fast
guards (`check-hooks-path`, `check-lint-config`, `check-templ-sync`, `check-version-sync`,
`check-module-sync`, `check-module-layers`, `check-replace-directives`, `check-css-minified`,
`check-lint-modules`, `check-tc-sources-sync`) then BuildFlow. `scripts/pre-commit.sh` is the
separate FULL pre-push verify.

## Flake commands (canonical entry points)

Build automation lives in `flake.nix`, not a Makefile. These wrap the pinned `templ`
(v0.3.1020, matching every `go.mod`) so generation is reproducible.

| Command | What it does |
| --- | --- |
| `nix develop` | Dev shell with the pinned templ, Go, golangci-lint, tailwind |
| `nix run .#build` | Regenerate `*_templ.go` + `go build ./...` |
| `nix run .#test` | `go test ./... -count=1 -race` — ROOT module only (see above) |
| `nix run .#lint` | golangci-lint across all 7 modules + actionlint |
| `nix run .#verify` | generate + build + test + lint |
| `nix run .#coverage` | Coverage profile + per-module summary |
| `nix run .#css` | Recompile `examples/demo/static/app.css` (Tailwind v4, minified) |
| `nix run .#visual` | Pixel-level visual regression (headless Chromium) |
| `nix run .#visual -- -update` | Regenerate visual goldens (the ONLY sanctioned way) |
| `nix run .#shots` | Full-page demo screenshots for manual inspection |
| `nix run .#crossbrowser` | 3-engine Playwright diagnostics — findings only, never gates |
| `scripts/release.sh <ver> "<summary>"` | One-commit release cut (see `docs/release-checklist.md`) |

The single done-check for component work: **`nix run .#verify`**, then the per-module test
loop. A "tests pass" without the loop is not a verdict.

## Process (run this when touching any `.templ` file)

1. **Enter the pinned dev shell before generating.** `nix develop`, or use the `nix run .#*`
   apps (same pinned templ). Any other templ binary produces cosmetic diffs across all 137
   generated files; restore `v0.3.1020` rather than bumping.
2. **Regenerate and build** (`nix run .#build`). Go never sees the `.templ` change until you
   generate. The root form (`templ generate ./...`) is right for root-module packages; when you
   edit a SUBMODULE's `.templ`, regenerate from inside that module (`cd utils && templ generate
   ./...`) so its `FileName` metadata stays module-relative. Never pass a subdir pattern from
   the root (it rewrites `FileName` paths repo-wide).
3. **Commit `*_templ.go`.** This is a library, not an app: the proxy serves tagged source
   as-is. `.gitignore` carries `!*_templ.go`. **BuildFlow gotcha:** the pre-commit
   `templ-generate` step re-appends the ignore line, hiding NEW generated files — `git add -f`
   them, and remove the re-added ignore line after the commit.
4. **Test the packages you touched**, then the full per-module loop. For an intentional output
   change: `go test ./<pkg>/... -update` — the `-update` flag goes AFTER the package list
   (the go command consumes a leading `-update` as its own flag).
5. **Recompile the CSS when classes change.** Any `templates/custom.css` edit or new class
   string in a `.templ` requires `nix run .#css` (the demo CSS scans every templ file); a
   stale `examples/demo/static/app.css` fails `TestCSSFreshness` in CI.
6. **Lint** (`nix run .#lint`; `examples/` is excluded by `.golangci.yml`).
7. **Keep `[Unreleased]` warm** in `CHANGELOG.md` in the same commit; `release.sh` refuses an
   empty section.
8. **No new dependencies.** Allowed: `templ`, `tailwind-merge-go`, `go-error-family`
   (errorpage only). Anything else needs an explicit decision.

## Definition of done for a new component

Beyond the code, a component is not finished until all of these land in the SAME commit:

- Golden HTML (`utils/golden`), a11y/behaviour/edge-case tests, and a godoc `Example`.
- The Props type registered in `internal/contract/component_props_test.go` (`componentTypes()`),
  enforced by `TestAllComponentPropsSatisfyInterface`.
- **Counts updated**: the per-package line here (`#### \`display\` — 52 components`),
  the README catalogue heading, and FEATURES where applicable — `TestDocsCountDrift` reads
  `skill/SKILL.md` directly and fails CI on a stale number.
- If the component is wired (HTMX/Datastar/JS), an e2e task in `visualtest/` — string-proven
  ≠ browser-proven.
- If it adds a closed-set enum, `IsValid()` plus its test in the same commit.

## Component anatomy (the canonical shape)

```
<package>/
  component.templ        # types, props, render template, private helpers
  component_templ.go     # GENERATED — commit it (git add -f if new)
  component_test.go      # unit + behaviour tests
  testdata/*.golden      # golden HTML snapshots
```

Inside `component.templ`:

1. **Typed enums** for every closed variant set, with `IsValid()` and a test.
   ```go
   type BadgeType string
   const (
       BadgePrimary BadgeType = "primary"
       BadgeSuccess BadgeType = "success"
   )
   ```
2. **Size constants** use `ComponentSize[SM|MD|LG]`.
3. **Props struct** embeds `utils.BaseProps` (only known exception: `layout.PageProps`),
   which auto-satisfies `utils.ComponentProps` via promoted pointer-receiver methods.
   **Gotcha:** promoted fields can be set directly in struct literals ONLY when the module's
   `go` directive is 1.27+ (`BadgeProps{Nonce: "x"}`); on older directives consumers must
   nest (`BadgeProps{BaseProps: utils.BaseProps{Nonce: "x"}}`). The `modernize` embedlit
   analyzer rewrites the nested form to the simple one — never hand-write the simple form in
   code that must compile below 1.27.
4. **Default constructor** `DefaultComponentProps()` with meaningful non-zero defaults whose
   zero value still matches the documented default.
5. **Render template** with a godoc example comment.
6. **Root element** propagates `props.ID`, `props.Class` (via `utils.Class(...)`),
   `props.Attrs`, `props.AriaLabel`. See `display/badge.templ`.
7. **Private lookup helpers** (`xxxClass()`, `xxxSizeClass()`) over typed-key maps.
8. **Shared rendering** at 2+ call sites extracts to a private sub-template (ADR-0010):
   5+ lines, a clear domain name, not a single caller, not demo code.
9. **Register the Props type** in `internal/contract/component_props_test.go` — forgetting
   this is the #1 way a component slips in without the BaseProps embed.

Shared-helper inventory, the full per-component test checklist, and every decision tree live in
`references/authoring-guide.md`; the complete guard matrix is in `references/guard-tests.md`.

## Repo-wide guard tests

They live in `utils/`, `internal/contract/`, `integration/`, and `visualtest/`. The full
matrix is in `references/guard-tests.md`; the ones authors trip most often:

- `TestTemplGeneratedInSync` — a `.templ` edit without regenerating.
- `TestDarkModeCompliance` / `TestMotionReduceCompliance` / `TestRTLLogicalProperties` /
  `TestCoarsePointerCompliance` / `TestAriaLivePoliteness` — class-string invariants.
- `TestContainerQueryCompliance` / `TestCustomCSSUtilities` / `TestTailwindGoSourceScanning`
  — breakpoint opt-ins and compiled-CSS coverage.
- `TestDocsCountDrift` — every count printed in this skill, README, AGENTS, FEATURES, and
  the site docs.
- `TestAllComponentPropsSatisfyInterface` / `TestPropsIDRendersInOutput` / `TestAttrsPropagateToRoot`
  — the consumer-facing contract.
- `TestAllInlineScriptsHaveNonce` — CSP.
- `TestTouchTargetAudit` / `TestZoomReflowAudit` + the axe-core sweep (default-FAIL, debt
  ledgered in `visualtest/testdata/axe_baseline.json`) — browser-measured a11y.

**When you add a cross-cutting rule, add a guard test rather than a comment.** The
auto-commit daemon bypasses the hook and never runs `go test`, so CI is the only reliable
enforcement layer.

## Mandatory conventions

- **Zero value = documented default.** Consumers write `XxxProps{...}` literals. Prefer
  inverted flags (`NoPad`, `NoValidate`) when the common case is "on" (the `ContainerProps.Pad`
  trap). Constructor-only defaults must be documented as such.
- **Closed-set enums ship `IsValid()` + a test in the same commit**; FEATURES tables are
  machine-guarded.
- **Counts are machine-guarded.** Adding/removing a component, icon, or golden updates the
  count-bearing docs (SKILL.md, README, FEATURES, site prose) in the same commit.
- **Dark mode:** `gray-*` exclusively, every neutral + semantic class carries `dark:`.
- **Motion:** every transition gets `motion-reduce:transition-none motion-reduce:duration-0`,
  every animation `motion-reduce:animate-none`; use the shared constants.
- **Class merging:** always `utils.Class(...)`; the only exception is `templ.KV` conditionals,
  which need comma-joining. `utils.Class` is mutex-protected; do not bypass it.
- **Tailwind variant literals:** write complete variant strings, never concatenate them at
  runtime (the scanner cannot see the token), and put peer-dependent classes on a SIBLING of
  `.peer` (`peer-checked:*:inset-s-N` for nested targets).
- **CSP nonce:** every inline script renders through `utils.ScriptComponent` (omit-empty rule);
  external scripts through `layout.Script`. New script emitters join
  `integration.TestNoEmptyNonceAttribute` + `TestEmptyNonceStillRendersScripts`.
- **RTL:** logical properties only (`ms-/me-/ps-/pe-/inset-s-*/inset-e-*/text-start/border-s-/border-e-`);
  never `ml-/mr-/pl-/pr-/left-/right-/text-left/border-l-/border-r-`, and never the deprecated
  bare `start-/end-` aliases. Only exception: `left-1/2` centering. Motion on logical insets
  (`peer-checked:*:inset-s-N` + `transition-[inset-inline-start]`), never `translate-x`.
- **Container queries:** 8 components carry `ContainerAware` — `Grid`, `Card`, `Nav`, `Split`,
  `DefinitionGrid`, `Form`, `Pagination`, `SkeletonCardGrid`. ALL are opt-in at the struct
  literal level (zero value = viewport); `DefaultGridProps()`/`DefaultSplitProps()` set `true`.
  Do NOT expand to marginal candidates (ADR-0018, `docs/container-query-strategy.md`).
- **Web Components / Shadow DOM are permanently rejected** (ADR-0033) — never re-evaluate.
- **`Card` shell is sharp** (`cardShellClass`, no `rounded-*`); restore rounding via `Class`.
  `SimpleCard` composes `Card`; `StatCard` uses the same shell.
- **Muted text:** `mutedTextClass` + a margin, never a bespoke string.
- **SVG paths:** reference `utils/svg` / `icons`; never inline `d="M..."`.
- **ProgressBar clamp:** `max(0, min(100, v))` (Go builtins).
- **templ language traps:**
  - `templ.KV(cond, class)` is a two-arg conditional attribute, NOT a ternary. Use
    `utils.Ternary(cond, a, b)` for value selection.
  - `checked={ someBool }` renders `checked="false"` when false, and any `checked`
    attribute means CHECKED. Render conditionally: `if sel { checked }`. Same for
    `disabled`, `selected`, `open`, `readonly`, `multiple`.
  - Void elements render WITHOUT a trailing `/>`; assert `>` endings.
  - Attribute values escape `'` as `&#39;`.
  - No implicit int→string: use `strconv.Itoa` in Go helpers.
  - No bare `continue`/`break` inside template loops (v0.3.1020 leaks the word as text) —
    use if/else.
  - NEVER import `github.com/a-h/templ` explicitly in a `.templ` file — the generator
    auto-injects it; an explicit import is a `templ redeclared` build error.
  - The templ LSP reports stale cross-module diagnostics: `nix run .#build` plus a fresh
    source read are ground truth; restart the LSP when it misleads.
- **Go directives move in lockstep:** `go 1.27` in `go.work` and every module, plus the
  toolchain, golangci-lint, CI, and Dockerfile image pins in ONE commit (guarded by
  `TestGoDirectiveSkew` + `TestGoWorkDirectiveMatchesRootGoMod` + `TestGoDirectivesAlignAcrossWorkspace`).
- **`cmd/tc/_sources` is a generated mirror:** never hand-edit it, never copy `*_templ.go` into
  it, and keep its package list in sync (`scripts/check-tc-sources-sync.sh` does the rest).
- **Demo pages pass the nonce:** every nonce-taking component rendered in `examples/demo/`
  gets `BaseProps.Nonce: demoNonceConst`; `nonce=""` is rejected by the demo CSP and kills the
  script silently (`TestDemoCSPNonceIntegrity`).
- **Auto-commit daemon:** it commits mid-session with generic messages, no tests, and can
  revert edits. After any daemon commit, verify your edits survived, `[Unreleased]` is still
  warm, and commit visual-golden PNG churn manually.

## Anti-patterns to refuse on review

- `switch` for style lookups → map + `utils.Lookup`.
- Inline `<script>` without nonce → `utils.ScriptComponent` / `layout.Script`.
- `*_templ.go` deleted or untracked after a `.templ` edit → `git add -f` if new.
- Raw `class={ "a " + b }` → `utils.Class("a", b)`.
- A new panic in render code → documented fallback.
- A new dependency outside the allowed three.
- `slate-*` dark-mode colors → `gray-*`.
- Hardcoded `aria-label` ignoring `props.AriaLabel`.
- Physical Tailwind properties or deprecated `start-`/`end-` insets → logical forms.
- Inline transition strings → the shared motion constants.
- Duplicated icon paths or class strings that already have a constant.
- Shipping a component with only assertion tests.
- Hand-typing a count into a doc instead of updating the guarded source.

## Where to read deeper

Load only when the task needs it.

- **`AGENTS.md`** — the full catalogue of conventions, gotchas, and the one-commit release
  convention; consult it whenever a rule's _why_ is unclear.
- **`references/`** — `component-catalogue.md` (signatures), `consumer-guide.md` (adoption),
  `authoring-guide.md` (helpers, tests, decisions), `guard-tests.md` (guard matrix).
- **`CONTRIBUTING.md`** — contributor setup and commit-message format.
- **`docs/adr/`** — decision records: two icon systems, shared SVG helpers, committing
  generated templ files, filled-vs-stroke icons, JS attachment, feedback unification,
  container queries, module split, transport wiring, typed triggers.
- **`docs/javascript-guide.md`** — the JS decision ladder, CSP, templ JS features
  (`OnceHandle`, `JSFuncCall`, `JSONString`, `JSONScript`), TypeScript, View Transitions.
- **`docs/DOMAIN_LANGUAGE.md`** — ubiquitous-language glossary.
- **`docs/count-definitions.md`** — what every sold count means.
- **`docs/visual-testing.md`** — visual suite, goldens, axe gate, emulation probes.
- **`internal/contract/component_props_test.go`** — the compile-time Props inventory.
- **`integration/composition_test.go`** — cross-package composition proof; extend it when a
  change affects how packages combine.
- **`examples/demo/`** — the canonical "how a consumer assembles a page" reference.
- **`FEATURES.md`** / **`TODO_LIST.md`** — honest inventory and in-flight work; check before
  proposing a new component.
- **`docs/release-checklist.md`** / **`scripts/release.sh`** — read before cutting a release.

## Installing this skill into Crush

The skill lives in the repo so it versions with the code. Crush auto-loads only from
`~/.config/crush/skills/`, so install it once from the repo root — link the entry point AND
`references/`, because a link to `SKILL.md` alone leaves every reference path dangling:

```bash
mkdir -p ~/.config/crush/skills/templ-components
ln -sfn "$PWD/skill/SKILL.md"   ~/.config/crush/skills/templ-components/SKILL.md
ln -sfn "$PWD/skill/references" ~/.config/crush/skills/templ-components/references
```

`-f` overwrites, `-n` treats a directory link as a file, so re-running is idempotent. Repeat
after cloning on a new machine. Verify `templ-components` appears in `available_skills` AND
that `references/component-catalogue.md` opens through the fan-out path — a skill whose links
resolve in git but not on disk is the 2026-09 stale-link failure class. Symlinking the whole
`skill/` directory as a top-level skill link is equally valid.

Never point this link at the nix store (repo policy: an editable skill must resolve to a
working checkout). The authoring norm and its guards are documented in
`/home/lars/projects/crush-config/README.md` § "Skill fan-out: authoring norm and integrity".

**Editing this skill:** keep the guard-matched counts accurate (`cd utils && GOWORK=off go
test -run TestDocsCountDrift -count=1`), keep the references non-overlapping (catalogue =
signatures, consumer guide = adoption, authoring guide = depth, guard tests = invariants), and
re-read the file through the fan-out path after editing — a change that lands in git but not
on disk is the stale-link failure class.
