# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

## [1.8.4] - 2026-08-16

### Fixed

- **`errorpage`: `FamilyOrchestration` now has an explicit `familyStatusCodeMap` entry (500).** The map entry was missing — orchestration errors only returned 500 because `FamilyStatusCode`'s unknown-family fallback happened to be 500. The 500 now comes from data, not accident, and an end-to-end test guards the map entry's presence. Same gap fixed in `htmx.GlobalErrorHandling`'s family-aware toast map (`tcFamilyToastMap`): `'orchestration'` is now mapped to `'error'` explicitly instead of relying on the `|| 'error'` fallback. Stale doc comment ("five defined constants") corrected to six.
- **`display.Card`: `ContainerAware` reverted to default `false`.** The v2.0 default flip (shipped in v1.8.2) made every `Card` emit an `@container` wrapper whose `container-type: inline-size` containment suppresses intrinsic width — inside shrink-to-fit parents (flex rows, `inline-block`, auto-sized grid columns) cards collapsed to **zero width**. Caught by the visual regression suite (`visualtest/testdata/card/*` — golden 229px vs actual 32px). `Grid` and `Split` keep the `true` default (layout primitives used in definite-width contexts). Set `ContainerAware: true` explicitly to restore v1.8.2 behavior; give the wrapper a definite width when the parent is shrink-to-fit.

### Added

- **`display.StatCardProps.ValueID` and `display.TableProps.BodyID`.** Optional ids on the StatCard value `<dd>` and the Table `<tbody>`, so live-updating dashboards can address those nodes directly instead of querying the components' internal structure — a markup refactor inside the component can no longer break consumer scripts (dnsblockd's dashboard JS contract, consumer TODO T86).
- **`feedback.CircularProgress` component.** SVG-based circular progress indicator with typed size enum (sm/md/lg), optional percentage label, and value clamping. Follows the library's standard props/enum/golden-test pattern.
- **`display.SectionHeading` and `display.DateRange` components.** SectionHeading provides typed heading level (H1-H6) and text alignment (Left/Center/Right) with `break-after-avoid` for print. DateRange renders a date range with a "Present" fallback and typed date format enum.
- **Print-friendly display components.** `Card` now includes `print:shadow-none print:border-0 print:bg-transparent`. `Modal` and `Drawer` include `print:hidden`. `PageHeader` and `SectionHeading` include `break-after-avoid` on headings. Print-safe vs screen-only matrix documented in `docs/theming.md`.
- **`emerald` brand preset** (`templates/presets/emerald.css`). Emerald-green primary palette matching the LarsArtmann/CV project theme.
- **`PageProps.SkipLinkText`, `PageProps.OGType`, `PageProps.BodyDataAttrs`.** `SkipLinkText` overrides the default "Skip to main content" for i18n. `OGType` sets the Open Graph type (defaults to "website"). `BodyDataAttrs` adds extra attributes to the `<body>` element (e.g. `data-language`).
- **`utils.Class()` output now sorted alphabetically.** Fixes non-determinism from `Oudwins/tailwind-merge-go` map iteration, making golden tests stable across runs.

- **4 new icons: Book, CircleStack, DevicePhoneMobile, ArrowTrendingUp.** Expands the icon catalogue from 102 to 106 icons (105 path icons + Spinner). All use official Heroicons v2 path data with appropriate animation mappings (Book/CircleStack/ArrowTrendingUp→nod, DevicePhoneMobile→pulse).

- **Animated icons (heroicons-animated inspired).** `icons.AnimatedIcon(name, class)` and `icons.AnimatedIconWithAnimation(name, anim, class)` render any icon with a hover-triggered CSS animation. 11 animation presets (`AnimPulse`, `AnimBeat`, `AnimBounce`, `AnimWiggle`, `AnimSpin`, `AnimJump`, `AnimNod`, `AnimShake`, `AnimBlink`, `AnimWobble`, `AnimDraw`) covering all heroicons-animated patterns. Pure CSS (zero JavaScript), `prefers-reduced-motion` support, triggers on `:hover` and `:focus-within`. Every icon has an explicit default via `DefaultAnimation()` — Heart→pulse, Bell→wiggle, Settings→spin, Eye→blink, Beaker→wobble, Bolt→draw (self-draw via stroke-dashoffset), Refresh→spin, etc. Aliases (ArrowPath, Bars3, MapPin, HandThumbUp) resolve to their canonical icon's animation. CSS lives in `templates/custom.css` under `.tc-anim-*` classes.

- **Fluid typography via container query units.** Six `.tc-fluid-*` utility classes (`tc-fluid-display`, `tc-fluid-h1`–`tc-fluid-h4`, `tc-fluid-lead`) in `templates/custom.css` size text with `clamp(min, Ncqi + base, max)` — text scales smoothly with its container's width, not the viewport. Composes directly inside all 8 container-aware components. Baseline 2023, zero JavaScript. Recipe: `docs/recipes/fluid-typography.md`.
- **Container query leveraging strategy.** `docs/container-query-strategy.md` maps the full landscape: shipped foundation, container query length units, style queries (deferred until Baseline), the v2.0 default flip, named containers, and honest evaluation of new component candidates (5 evaluated, all rejected) and the `containerAwareWrapper` consolidation (declined).
- **Visual regression tests are now a hard CI gate.** The Visual Regression job detects skipped tests (e.g. missing Chromium) and fails the pipeline, preventing "vacuously green" runs where tests silently skip without a browser.
- **`TestCustomCSSUtilities` drift-guard.** New scanner in `utils/custom_css_test.go` asserts every `tc-*` CSS class used in `.templ` files is defined in `templates/custom.css` — catches silent CSS deletions and missing definitions before consumers hit visual regressions.
- **Datastar runtime version pinned via `go-datastar/static`.** `DatastarVersion1_0_2` is now derived from `static.Version` at compile time (`DatastarVersion(static.Version)`), making the embedded JS bundle and CDN URL structurally incapable of drifting. Zero transitive dependencies (`go-datastar/static` is a zero-require module). `TestDatastarVersionMatchesStatic` drift-guard test enforces the linkage.
- **`htmx` and `datastar` extracted as standalone Go sub-modules.** Both packages now have their own `go.mod` (Layer 1 in the DAG, depending only on `utils`), matching the `charts/echarts` and `icons` precedent. Consumers who want only HTMX or Datastar components no longer pull the root module's dependency graph. Test spinners replaced `feedback.Spinner` imports (circular dep elimination).

### Changed (v2.0 breaking changes — ADR-0022)

- **HTMX is now self-hosted by default via `//go:embed`.** `DefaultPageProps()` sets `HTMXSrc: HTMXSelfHost`, which embeds HTMX 2.0.10 inline and renders it as `<script nonce="...">`. No external CDN request. Consumers who prefer CDN loading set `HTMXSrc: ""` and provide `HTMXVersion`. **CSP impact:** the inline script requires `script-src 'nonce-...'` (or `'unsafe-inline'`) in the consumer's CSP policy. See `docs/migration/v1-to-v2.md`.
- **Container-aware default flipped to `true` for `Grid`, `Card`, and `Split`.** These three components now emit `@container` wrappers and use container-query breakpoints (`@sm:`/`@md:`/`@lg:`) by default instead of viewport breakpoints (`sm:`/`md:`/`lg:`). Set `ContainerAware: false` to restore viewport-breakpoint behavior. The other 5 container-aware components (`Nav`, `Pagination`, `Form`, `DefinitionGrid`, `SkeletonCardGrid`) remain opt-in (`false`).
- **`Grid.ContainerResponsive` renamed to `Grid.ContainerAware`** for consistency with all other container-aware components. Mechanical rename: replace `ContainerResponsive` with `ContainerAware` in all `GridProps` literals.
- **`AlertType`/`ToastType` aliases removed.** The backward-compat type aliases `type AlertType = FeedbackType` and `type ToastType = FeedbackType`, plus their constants (`AlertSuccess`/`Error`/`Warning`/`Info`, `ToastSuccess`/`Error`/`Warning`/`Info`), are removed. Use `FeedbackType` and `FeedbackSuccess`/`FeedbackError`/`FeedbackWarning`/`FeedbackInfo` directly. `AlertProps.Type` and `ToastProps.Type` are now typed `FeedbackType` (was `AlertType`/`ToastType`).

### Architecture

- **ADR-0034: Targeted 7-module workspace split.** The single-module library is now 7 Go modules coordinated by `go.work` (local dev) + `replace` directives (CI/consumers): `utils` (leaf: BaseProps, Class(), EnsureID, svg, cdn, golden), `icons` (102 SVG icons, icons-only adoption), `errorpage` (isolates `go-error-family`), `charts/echarts` (opt-in adapter), `htmx` (HTMX loading/error/OOB components), `datastar` (Datastar runtime + SSE LiveRegion), and root (all core UI + recipes + integration + demo + CLI). `internal/svg`, `internal/cdn`, `internal/golden` promoted to `utils/` sub-packages (Go's `internal/` rule blocks cross-module access). Import paths are unchanged for all externally-importable packages. Shared versioning: `scripts/release.sh` tags all 7 modules at release time. Per-module CI isolation tests verify standalone builds. Supersedes ADR-0020 (proposed per-package ~12-module split, deferred). See `docs/adr/0034-targeted-module-split.md`.
- **ADR-0033: Web Components permanently rejected.** Shadow DOM breaks the Tailwind-utility theming model; Custom Elements require JavaScript (violating the zero-JS principle); the cross-framework distribution problem WC solve doesn't exist for a Go-source library. The library achieves "use the platform" (the actual goal) via native APIs (`<dialog>`, Popover, `<details>`, scroll-snap, `@container`). Added to ROADMAP "Explicitly NOT Planned".

## [1.8.1] — 2026-08-09

### Added

- **`display.BarChart` richer per-bar metadata.** New `BarChartBar.Tooltip` (emits a native `title` attribute for hover tooltips on dense charts), `BarChartBar.ValueLabel` (overrides the auto-formatted value for composite labels like "123 (45%)" or "1.2 GB"), `BarChartProps.MinBarWidth` (Tailwind class, tightens vertical bars), and `BarChartProps.Gap` (Tailwind class, configurable vertical-bar spacing down to `gap-px`). All fall back to current behavior when unset.
- **`display.BarChart` fixed-height vertical charts.** New `BarChartProps.Height` (CSS value string, e.g. "8rem", "200px") applied as an inline `style="height:..."` on the vertical chart container. Fixes bars collapsing to zero height when the parent has no definite height. Horizontal charts are unaffected.
- **`navigation.SidebarNav` collapsible sections + header slot.** New `SidebarNavItem.Section` groups consecutive items under a collapsible `<details>` header (empty `Section` stays flat, preserving existing behavior); the section containing the active item auto-expands. New `SidebarNavProps.Header` slot renders between Brand and nav links for search inputs or filters.
- **Drift-guard extension: README.md + ROADMAP.md.** `TestDocsCountDrift` now asserts component/enum/visual-golden counts in README.md and ROADMAP.md, closing the gap that let the version badge and golden counts drift unnoticed for multiple releases.
- **Version-sync pre-commit guard.** New `scripts/check-version-sync.sh` extracts the version from `utils/version.go`, `CHANGELOG.md`, and `FEATURES.md` and blocks commits if they disagree (<50ms shell, mirrors `check-templ-sync.sh`). Wired into `.git/hooks/pre-commit` (Guard 3) and CI.
- **Actionlint in CI.** GitHub Actions workflow files are now linted with `actionlint` on every push/PR.
- **Fuzz tests for chart geometry.** `FuzzBuildSmoothPath` and `FuzzBuildAreaPath` join the existing chart fuzz suite, verifying the Catmull-Rom spline and area-path builders never panic on adversarial inputs (NaN, Inf, empty slices, extreme coordinates).
- **Unit test for ordered-substring guard predicate.** `TestIsOrderedTailwindSubstring` provides 16 table-driven cases (positive violations + negative non-violations), closing the meta-test gap where the drift guard itself was untested.

### Fixed

- **Visual test harness: parallel tab isolation.** `TestWaitAnimationSettled` subtests no longer share a single browser tab — each gets its own via `newTab(t)`, eliminating context-cancellation failures when parallel navigations clobbered each other.
- **Visual test harness: two-phase transition detection.** `waitAnimationSettled` now waits for `@starting-style` transitions to *register* (appear in `getAnimations()`) before waiting for them to *finish*, fixing the 99%+ false mismatch on drawers/popovers captured mid-slide under parallel load.
- **Visual test thresholds.** `TestPolledRegion` raised to 1% MaxMismatch (sub-pixel font rendering noise); `TestSpinner` raised to 8% (continuously rotating CSS animation catches random frames).

## [1.8.0] — 2026-08-08

### Added

- `charts/echarts` package added to `countExportedTemplFunctions` drift guard — `EChart` and `SDKScript` are now counted in the component total (110→112). The drift-guard regex is now flexible (`across \d+ packages`) so adding future packages doesn't break the test.
- Visual regression CI lane — `nix run .#visual` now runs in CI on every push/PR with Chromium provided by Nix (bit-identical renderer to golden PNGs). Previously visual tests only ran locally; CI silently skipped without a browser ("vacuously green" risk).
- `ChartPadding.Sanitize()` method clamps all padding fields to non-negative values before rendering. Negative padding previously produced negative plot dimensions, corrupting SVG path math.
- `SanitizeInnerRadius()` function clamps the PieChart `InnerRadius` to `[0, 1]`. Values outside this range previously produced broken arc paths. NaN is clamped to 0.
- `scripts/check-templ-sync.sh` — fast (<1s) pre-commit guard that catches `*_templ.go` drift before commit. Mirrors `check-lint-config.sh`, wired into `.git/hooks/pre-commit` and CI. Prevents the BuildFlow daemon from committing stale generated files.
- CSS freshness CI lane — recompiles demo CSS via `nix run .#css` and diffs against the committed file. Prevents stale CSS shipping (as happened in v1.7.0).
- Fuzz tests for chart geometry math: `FuzzScalePoints`, `FuzzComputeNiceTicks`, `FuzzComputeArcPath` — verify no panics on NaN/Inf/negative/extreme inputs (2M+ executions, zero failures).
- Visual regression tests for charts: LineChart, PieChart, DonutChart, AreaChart now have golden PNG baselines.
- Dark-mode visual variants for Combobox, Tooltip, Carousel, Skeleton, ErrorPage, NotFound404.
- Visual regression tests for v1.5–v1.6 components: BarChart, Heatmap, Sparkline, CollapsibleSection, ExternalLink, PolledRegion, DataTable.
- `TestWaitAnimationSettled` — dedicated test for the visual test harness covering polling, empty-animations, and timeout paths.
- Benchmarks for PieChart arc computation (`BenchmarkComputeSliceAngles`, `BenchmarkComputeArcPath`) and full LineChart render (`BenchmarkLineChartRender`).

### Changed

- Shared chart sub-templates extracted: `chartAxes`, `chartLegend`, `chartEmptyStateMsg` in `display/chart_shared.templ`. Eliminates ~80% template duplication between LineChart and AreaChart. `ChartRenderData` struct + `computeChartRenderData()` centralize the shared setup logic.
- Bumped `nixpkgs` and `treefmt-nix` flake inputs to latest upstream revisions (security patches, formatter improvements). Updated `github.com/chromedp/cdproto` in `visualtest/go.mod` to pick up the latest Chrome DevTools Protocol type definitions.
- Corrected all stale counts in README.md: enums (51/47/43 → 52), visual goldens (31 → 49), HTML goldens (102 → 175), components (107/98 → 112). The drift-guard test did not cover README, allowing silent drift.
- Merged orphaned `### Enums` section (GridGap alone) into the main display Enums table in FEATURES.md.
- Archived 6 fully-resolved historical docs (2 planning + 4 status) to `docs/{status,planning}/archived/` with resolution banners.

### Fixed

- `breadcrumbs_templ.go` out of sync (recurrence) — the generated file had drifted back to `encoding/json/v2` while the source uses `encoding/json` (stable). Regenerated to match (`TestTemplGeneratedInSync` now passes).
- Removed unused constant `pieChartLegendCharW` from `display/pie_chart.go` (dead code, verified via grep — zero references).
- Fixed "Catull-Rom" → "Catmull-Rom" typo in v1.7.0 CHANGELOG entry.
- **CSP nonce handling:** `datastar.SDKScript`, `layout.ThemeScript`, and `layout.ThemeToggle` no longer render an empty `nonce=""` attribute when no nonce is configured. Empty nonces are CSP no-ops and broke strict-CSP consumers (e.g. go-health-dashboard). The `nonce` attribute is now emitted only when a non-empty value is supplied. Golden baselines updated.

## [1.7.0] — 2026-08-04

### Added

- **`display.LineChart`** — pure-SVG line chart with axes, gridlines, multi-series support, data-point dots, legend, linear/smooth styles, and ARIA. Zero JavaScript. Part of the two-tier chart architecture (ADR-0031).
- **`display.PieChart`** — pure-SVG pie/donut chart with arc paths, external labels, legend, donut center label, custom colors, and ARIA. Zero JavaScript.
- **`display.AreaChart`** — pure-SVG area chart (line chart with filled areas), configurable fill opacity, multi-series, smooth curves.
- **`display.chart_geometry.go`** — shared SVG chart geometry helpers: `ScalePoints`, `BuildPolylinePath`, `BuildSmoothPath` (Catmull-Rom), `BuildAreaPath`, `ComputeNiceTicks`, `FormatTickValue`. The foundation for all native SVG chart components.
- **`charts/echarts`** — opt-in ECharts adapter package (`EChart`, `SDKScript`) with dark mode bridge. Accepts go-echarts `RenderSnippet()` output as strings — zero dependency on go-echarts. Follows the datastar opt-in precedent.
- **`LineChartStyle`** and **`PieChartLabelMode`** typed enums with `IsValid()` methods.
- Recipe docs: `docs/recipes/line-chart.md`, `pie-chart.md`, `area-chart.md`, `echarts-adapter.md`.
- ADR-0031: Two-Tier Chart Architecture (Native SVG + Opt-in ECharts).
- Golden snapshot tests for LineChart (10 baselines), PieChart (8 baselines), AreaChart (7 baselines).
- Demo showcase: SVG Charts section in `examples/demo/display_demo.templ`.
- **`recipes.AuthLayout`** — split-screen authentication layout with a card panel and branding panel. Supports reversed layout, panel features list, and panel footer.
- **`visualtest.Bool()`** — tri-state option helper for `Dark`/`RTL` fields. `*bool` semantics: nil=unset (default), `Bool(true)`=explicit dark/RTL, `Bool(false)`=explicit light/LTR. Prevents zero-value conflation.
- **`visualtest.ViewportMobile`**, **`ViewportTablet`**, **`ViewportDesktop`** — viewport preset constants for visual regression tests.
- **`visualtest.InteractionState.String()`** — debug-friendly string representation of interaction state.
- **HTMX golden snapshot tests** — 13 new golden baselines covering all HTMX components (LoadingIndicator, InlineLoadingOverlay, LoadingButton, CSRFToken, ConfirmDelete, SwapOOB, GlobalErrorHandling, ViewTransitions).
- **`nix run .#css`** — flake app that compiles Tailwind CSS with `--minify` for quick CSS rebuilds.
- **`tc version`** command and **`tc add --list-deps`** flag — CLI ergonomics for dependency inspection.
- **Popover edge-flipping** — dropdowns and popovers automatically flip to the opposite side when the preferred side clips off-screen.
- **`docs/testing-guide.md`** — comprehensive three-tier testing strategy documentation.
- **`docs/migration/skeletoncardgrid-api-change.md`** — migration guide for the SkeletonCardGrid Props API change.
- **`utils.TestNoOrderedTailwindSubstringsInTests`** — drift-guard that scans all test files for brittle ordered Tailwind substring assertions (e.g., `strings.Contains(out, "flex flex-col")`). Catches latent flakes from `utils.Class()` nondeterministic reordering.

### Changed

- **`visualtest.Options.Dark`/`RTL`** changed from `bool` to `*bool` — tri-state semantics prevent zero-value conflation between "unset" and "explicitly false". All visual tests updated to use `visualtest.Bool(true)`.
- **`visualtest` hover targeting** — `hoverAction` now descends to the first interactive child element instead of targeting `#tc-root` center, fixing hover-state screenshots.
- **treefmt aligned with gofumpt** — switched from `gofmt` to `gofumpt` to eliminate the latent formatter conflict with golangci-lint.
- **Chromium pinned via separate nixpkgs input** — added `nixpkgs-chromium` input decoupled from main `nixpkgs` so `nix flake update` doesn't shift pixel output in visual regression tests.

### Fixed

- **Modal/Drawer visual regression tests** — `<dialog Open=true>` doesn't promote to top-layer without `showModal()` JS. Fixed by using `FullViewport: true` + `WaitSelector: "dialog"`.
- **SkeletonCardGrid test assertion** — replaced brittle `strings.Contains(output, "rounded-lg border")` with `AssertNotContains`.
- **34 golangci-lint findings resolved** — renamed short variables in chart geometry, removed unused function, fixed whitespace and style violations.
- **`.envrc` missing `GOEXPERIMENT=jsonv2`** — the export was dropped, so every tool outside `nix develop` (gopls, BuildFlow, IDE, bare `go`) silently misbuilt the module. Restored (`TestEnvrcConsistency` now passes).
- **`breadcrumbs_templ.go` out of sync** — `navigation/breadcrumbs.templ` was switched to `encoding/json` (stable) but the generated file still imported `encoding/json/v2`. Regenerated `*_templ.go` to match (`TestTemplGeneratedInSync` now passes).
- **Visual overlay `MaxMismatch` calibration** — ran a 10× serialized calibration of all 8 overlay goldens (Dropdown/Popover/ContextMenu via `overlayOpen`; Modal/Drawer via `dialogOpen`) under the pinned Chromium; confirmed 0.0000% run-to-run mismatch (fully deterministic). Regenerated the stale `dropdown/open_dark` golden (was at a stable 0.7442% systematic diff — a stale golden the prior comment misattributed to anti-aliasing noise). Updated the `overlayOpen`/`dialogOpen` comments and `docs/visual-testing.md` with the rigorous data; the 1% threshold is validated as pure headroom for Chromium-version drift.
- **Visual test overlay race condition** — `dialogOpen` captures (Modal/Drawer) raced the `@starting-style` slide-in transition under parallel load: `WaitVisible("dialog")` returned before the 200ms transition completed, occasionally capturing the drawer off-screen (~90% false mismatch, ~20% full-suite flake rate). Added `waitAnimationSettled` to the harness, which polls `getAnimations()` until all CSS transitions finish before capture. Full suite now passes 8/8 under parallel load with 0.0000% overlay mismatch.

## [1.6.0] — 2026-07-30

### Added

- **`display.CollapsibleSection`** — native `<details>/<summary>` collapsible region with configurable heading level (h1-h6), open/closed default, optional `StorageKey` for localStorage persistence, and chevron rotation with motion-reduce and dark-mode support. Inspired by DiscordSync's `collapsibleSection` helper.
- **`display.Heatmap`** — CSS grid heatmap with row/column labels, opacity-based cell coloring, peak highlighting with ring, clickable cells, tooltip support, and empty state. No JavaScript, no SVG. Inspired by DiscordSync's `activityHourlyHeatmap` and `heatmapMatrix`.
- **`htmx.PolledRegionProps.Trigger`** — first-class `Trigger string` field that overrides the auto-generated `hx-trigger`, enabling custom SSE/WebSocket triggers without clobbering `Attrs`. Eliminates the DiscordSync wrapper hack.
- Recipe docs: `docs/recipes/collapsible-section.md`, `docs/recipes/heatmap.md`.
- ADRs: `0027-collapsible-section-native-details.md`, `0028-heatmap-css-table-opacity.md`, `0029-polled-region-trigger-field.md`.
- Golden snapshot tests (10 baselines), unit tests, a11y tests, and benchmarks for both new components.
- Demo showcase: CollapsibleSection and Heatmap sections in `examples/demo/display_demo.templ`.

### Changed

- Component count: 33 → 35 display components. Updated SKILL.md, README.md, FEATURES.md, AGENTS.md, website sections.ts.
- `display/FEATURES.md` version updated to 1.5.0 (was stale at 1.4.0).

### Fixed

- 8th BuildFlow `:=` shadowing regression in `visualtest/doc.go:76` (changing `=` to `:=`, shadowing package-level `sharedAllocCtx`/`allocCancel`).

## [1.5.0] — 2026-07-30

### Added

- **`htmx.PolledRegion`** — auto-refreshing HTMX region for dashboard stats. Polls a URL at a configurable interval, swaps content in place, includes optional eager-load, `aria-live` for screen readers, and a "Updated HH:MM:SS" freshness timestamp. Inspired by DiscordSync's `polledRegion` helper (the #1 dashboard pattern).
- **`display.Sparkline`** — tiny inline SVG line chart for trend visualization. Pure SVG (no JS), `currentColor` stroke, optional filled area, auto min/max bounds. Inspired by DiscordSync's `sparklineSVG`.
- **`display.BarChart`** — CSS-based horizontal/vertical bar chart. Per-bar colors, link labels, custom value formatting, empty-state message. No JavaScript, no SVG. Inspired by DiscordSync's 8+ hand-rolled bar chart variants.
- **`display.ExternalLink`** — safe-by-default off-site link with `target="_blank" rel="noopener noreferrer"`, external-arrow icon, and URL sanitization (plain string href, not `templ.SafeURL`). Inspired by DiscordSync's `externalLink` helper.
- **Golden HTML snapshot tests** for Sparkline, BarChart, ExternalLink, and PolledRegion — every component now has golden baselines matching the library's three-tier testing standard.
- **BDD, benchmark, and a11y tests** for all 4 new components — PolledRegion BDD in `htmx/bdd_test.go`, benchmarks in both packages, accessibility assertions in `display/a11y_new_test.go`.
- **Card and EmptyState `TitleTag` now supports h1–h6.** Previously limited to h2/h3; the switch now covers all semantic heading levels for correct a11y heading hierarchy.
- **PolledRegion `TimeFormat` field.** Configurable Go time format string for the timestamp footer (default: `"15:04:05"`). No more hardcoded format.
- **Recipe docs** for sparkline, bar-chart, polled-region, external-link.
- **ADRs** 0024 (PolledRegion design), 0025 (BarChart CSS vs SVG), 0026 (ExternalLink sanitization).
- **Carousel keyboard navigation.** The carousel region is now focusable (`tabindex="0"`) and responds to ArrowLeft/ArrowRight (prev/next slide), Home (first slide), and End (last slide). RTL-aware: ArrowLeft/Right are swapped in `dir="rtl"`. Follows WAI-ARIA carousel pattern.
- **MobileMenu keyboard support.** Escape closes the menu and returns focus to the toggle button. Opening the menu moves focus to the first focusable child. Extracted `tcMobileMenuSet(menu, btn, open)` shared helper for consistent open/close + focus management.
- **Dropdown keyboard enhancements.** Home (first item), End (last item), PageDown/PageUp (jump by quarter of the list) now move focus inside the menu. Disabled items (`[disabled]`) are skipped during navigation. First menuitem is auto-focused when the menu opens (via `toggle` event listener).
- **Tooltip Escape-to-dismiss.** Pressing Escape while a tooltip trigger has focus hides the tooltip via a `data-tc-tooltip-dismissed` attribute + CSS rule, while keeping focus on the trigger (preserving tab position). Hover or re-focus clears the dismissed state.
- **ContextMenu keyboard accessibility.** The context menu is now operable without a mouse: Shift+F10 and the dedicated ContextMenu key open it (positioned at the trigger element), and the menu supports full WAI-ARIA menu keyboard navigation (ArrowUp/Down with RTL-aware Left/Right, Home, End, focus-first-on-open). Menuitems gained `tabindex="-1"` (roving focus) and disabled items expose `aria-disabled="true"`. Escape and click-outside remain native via the Popover API.
- **Shared menu keyboard-navigation helper** (`display/shared.go`). Extracted the WAI-ARIA menu keydown handler (Arrow/Home/End/PageUp/PageDown + focus-first-on-toggle) into a single singleton script (`menuKeyboardNavScriptComponent`) shared by Dropdown and ContextMenu, replacing the Dropdown's inline duplicate. The nav selector now also skips `aria-disabled` items, not just `[disabled]`.
- **Carousel focus-visible outline.** The focusable carousel region now shows a visible focus ring (`focus-visible:ring-2`) so keyboard users can see where focus landed before navigating with arrow keys.

### Fixed

- **Sparkline Min/Max zero-sentinel bug.** Changed `Min`/`Max` from `float64` to `*float64`. Previously, a legitimate min/max of 0.0 was treated as "auto" (sentinel value). Now `nil` means auto-compute from data; a pointer to 0.0 explicitly sets the bound to zero.
- **Dropdown Enter/Space no longer breaks HTMX.** The custom JS handler used `window.location.href = item.href` for `<a>` menuitems, which triggered a full page load and bypassed HTMX's AJAX swap. Removed the handler entirely — native browser behavior activates links/buttons correctly, and HTMX intercepts native events as expected.
- **Rating arrow-key direction now matches value.** The interactive Rating rendered its radio inputs in reverse DOM order (5→1) for a CSS-only fill trick, which inverted radiogroup arrow-key behavior (ArrowDown/Right decreased the value). Switched to forward DOM order (1→N) with `flex-row-reverse` for the visual, so arrows increase the value per WAI-ARIA while the `peer-checked` fill still renders ★★★☆☆ correctly. Forward order is used only for the interactive branch; read-only rendering is unchanged.
- **Rating star fill now renders.** The `peer-checked` fill classes lived on the nested `<svg>` (not a sibling of the hidden `.peer` radio), so Tailwind's `~` combinator never matched and the selected star did not visually fill. Moved the color/checked/hover/focus classes to the `<label>` (the radio's sibling); the SVG inherits color via `currentColor`.

## [1.3.0] — 2026-07-28

### Added

- **Lint-config regression prevention (3-layer guard).** Root-caused the recurring `.golangci.yml` regression (5th occurrence) to the BuildFlow daemon committing stale working trees. Prevention: (1) `scripts/check-lint-config.sh` — <50ms standalone grep guard wired into `.git/hooks/pre-commit` BEFORE BuildFlow runs; (2) CI step "Lint-config guard" in `.github/workflows/ci.yaml` runs the script before `golangci-lint` even installs; (3) `TestGolangciDisabledLinters` in `utils/lint_config_test.go` catches in CI via `go test ./...`. Root cause and prevention layers documented in AGENTS.md.
- **Generated-file sync guard** (`utils/templ_sync_test.go`). `TestTemplGeneratedInSync` verifies every `*_templ.go` file's imports match its `.templ` source. Prevents stale generated artifacts (breadcrumbs drift: source had `encoding/json`, generated had `encoding/json/v2`). Runs in ~15ms across all 74 `.templ` files.
- **Container-query compliance scanner** (`utils/container_query_compliance_test.go`). `TestContainerQueryCompliance` scans `.templ` files for structural viewport breakpoints (`sm:grid-cols`, `md:flex`, `lg:hidden`) without a corresponding `ContainerAware` flag. Mirrors the proven `TestDarkModeCompliance`/`TestMotionReduceCompliance` pattern with an explicit exemption list.
- **Tailwind CSS Go-source scanning** (`utils/tailwind_source_test.go`). Added `@source "**/*.go"` to both `templates/app.css` (consumer template) and `examples/demo/demo.css`. Tailwind v4's content scanner now reads Go files for class lookup maps (`buttonVariantLookup`, `feedbackStyleMap`, `familyStyleMap`). Previously, errorpage family classes (amber/orange/purple) were silently missing from compiled CSS — `bg-amber-50` had 0 matches in the demo CSS. After fix: all shades present. Enforced by `TestTailwindGoSourceScanning`.
- **Shared Chromium process** (`visualtest/doc.go` + `visualtest/main_test.go`). Replaced per-test browser launches with a `sync.Once`-initialized shared allocator. Each test gets a new tab (~10ms) from the shared browser process. 15 visual tests now complete in ~2s total (was ~10s+ with per-test browsers). `TestMain` ensures cleanup.
- **Visual regression goldens expanded.** Added 16 new goldens across two passes: Modal (light/dark), Drawer (left-light/right-dark), Input (text-light/dark/error/disabled), Select (light/dark), RTL (button/card), Dropdown (light/dark), Popover, ContextMenu, plus Pagination/Breadcrumbs/Alert/Input snapshots. Total: **31 goldens** across 11 component types + RTL coverage.
- **Visual coverage metric** (`visualtest/coverage_test.go`). Reports golden-to-component ratio. Currently 31 goldens / 74 components = 41.9%.
- **Visual harness interaction states.** `StateClick` (click first interactive descendant), `StateContext` (right-click trigger), `FullViewport` (capture full page, not viewport-clipped — for overlays promoted to the top layer), and `WaitSelector` (wait for a selector before capture, e.g. an open menu). Enable capturing open-state of Popover API and `<dialog>` components (`visualtest/harness.go`, `visualtest/render.go`).
- **CSS staleness detection** (`utils/css_freshness_test.go`). Warns when committed `app.css` is older than the newest source file change.
- **GOWORK=off in devShell** + **`.envrc` (direnv)** for repo-wide `GOEXPERIMENT=jsonv2` + `GOWORK=off`. The devShell `shellHook` now sets both vars. `.envrc` ensures all tools (go, gopls, BuildFlow, IDE) inherit them without `nix develop`. Run `direnv allow` once after cloning.
- **Tailwind CSS compile step** added to `scripts/release.sh`. Now runs `tailwindcss --minify` after `templ generate` so the committed `app.css` is never stale at release time.
- **v2.0 migration design** (`docs/adr/0022-v2-default-flip-migration.md`) and **compound overlay API design** (`docs/adr/0023-compound-overlay-component-api.md`).
- **Visual regression testing framework.** New `visualtest/` module — a **separate Go module** (`visualtest/go.mod` with a local `replace` directive) so `chromedp` never pollutes the library's consumer dependency graph. Renders components in headless Chromium and diffs pixels against committed golden PNGs via `orisano/pixelmatch` (YIQ perceptual distance + anti-alias skip). Public API: `visualtest.AssertScreenshot(t, name, component, opts...)` with `Options{Dark, RTL, Viewport, MaxMismatch, Threshold, State}` (states: `StateHover`, `StateFocus`). Mirrors `internal/golden` DX: `-update` flag regenerates goldens, `.fail/` artifacts (actual + diff PNGs) auto-clean on pass. Wired into Nix (`nix run .#visual`) and CI (`.github/workflows/ci.yaml` `visual` job — Nix-based for renderer bit-parity, uploads diff artifacts on failure). Tests skip cleanly when no browser is available. See [`docs/visual-testing.md`](docs/visual-testing.md).
- **Container-aware components (ADR-0018 extension).** Extended the `ContainerAware` opt-in pattern to 5 more components so they adapt to their parent container's width instead of the viewport:
  - `layout.Split.ContainerAware` — main+aside collapses to stacked by container width (`@md:` variants)
  - `display.DefinitionGrid.ContainerAware` — term-detail card grid column count adapts (`@sm:`/`@lg:`)
  - `forms.Form.ContainerAware` — Grid layout label/value columns adapt (`@sm:`), ideal for forms in modals/drawers
  - `navigation.Pagination.ContainerAware` — mobile/desktop controls switch by container width (`@sm:` on nav root)
  - `feedback.SkeletonCardGrid` — converted from `count int` to `SkeletonCardGridProps{Count, ContainerAware}` for API consistency; loading skeletons now match `Grid.ContainerResponsive` behavior
  - All follow the exact ADR-0018 contract: opt-in flag, `@container` wrapper (or root class), lookup-map class swap, default-off backward compat. Tested per the established pattern.
- **Rename-safety test** (`errorpage/fromerror_safety_test.go`) — verifies every `errorfamily.Family` maps to a valid, distinct `Family` output (totality, correctness, injectivity), so a missing switch case fails loudly rather than silently rendering every error as Transient.
- **Environment & hook guards.** `TestEnvrcConsistency` asserts the tracked `.envrc` matches the expected `GOEXPERIMENT=jsonv2` + `GOWORK=off` content (catches the daemon re-ignoring it). `TestPreCommitHookInstallsGuard` asserts the pre-commit hook installs the required guards. `TestCSSFreshness` now fails in CI (`t.Errorf` when `CI` env) so a stale committed `app.css` red-lines CI instead of silently warning.

### Changed

- **`feedback.SkeletonCardGrid` signature changed (breaking).** `SkeletonCardGrid(count int)` is now `SkeletonCardGrid(SkeletonCardGridProps{Count, ContainerAware})` for API consistency with the container-aware pattern and the rest of the library. Migration: `SkeletonCardGrid(6)` → `SkeletonCardGrid(feedback.SkeletonCardGridProps{Count: 6})`. `DefaultSkeletonCardGridProps()` returns sensible defaults. This is a source-level breaking change in a post-v1.0 minor; it is accepted because `SkeletonCardGrid` has no known external consumers yet and the props struct is the forward-compatible shape. Consumers who upgrade will get a clear compile error at the call site.
- **Type-erasure fix in `FromErrorFamily`** (`errorpage/fromerror.go`) — replaced the `ParseFamily(f.String())` string round-trip with a typed `switch` on `errorfamily.Family`. Any future rename or addition of a family constant in go-error-family is now caught at compile time instead of silently collapsing to `FamilyTransient`. Added the `Orchestration` case plus a `FamilyOrchestration` constant and purple visual style in `familyStyleMap`.

### Fixed

- **`flake.nix` devShell now exports `GOEXPERIMENT=jsonv2` via `shellHook`.** Previously the devShell provided the tools but not the env var, so BuildFlow's `test-coverage` step (and any `go` subprocess outside the individual nix apps) ran without the experiment flag. This caused 5 packages to fail with `build constraints exclude all Go files in encoding/json/v2`: `errorpage`, `navigation`, `integration`, `internal/contract`, `examples/demo`. All nix apps already set the flag individually; this fix covers the gap for ad-hoc commands and third-party tools (BuildFlow) that run as subprocesses of the devShell.
- **`.golangci.yml` lint gate repaired (regression of the v0.19.0 fix).** The three linters documented in AGENTS.md as fundamentally incompatible with a templ library and "do NOT re-enable" (`ireturn`, `godoclint`, `testableexamples`) had re-entered the `enable:` list (along with the dead `ireturn:` settings block), causing `golangci-lint run` to exit 1 with 71 findings despite AGENTS.md and the v0.19.0 CHANGELOG claiming it exits 0. Removed all three from the enable list and deleted the dead `ireturn:` block. `golangci-lint run` now exits 0 as documented. **Added `TestGolangciDisabledLinters`** drift-guard test (`utils/lint_config_test.go`) — this was the **third recurrence** of the regression (the auto-commit daemon kept reverting the manual fix); the test now fails CI the moment any of the three linters re-enters the enable list, so it cannot recur a fourth time.
- **`.golangci.yml` disabled-linter regression removed again (6th recurrence, pre-v1.3.0).** The auto-commit daemon's commit `d1dd3da` ("align lint configuration") reformatted `.golangci.yml` and silently re-added `ireturn`, `godoclint`, `testableexamples` plus the dead `ireturn:` settings block to HEAD — the exact T1 regression, now committed rather than just a stale working tree. Removed all three and the dead block again pre-release; `scripts/check-lint-config.sh`, `TestGolangciDisabledLinters`, and the CI "Lint-config guard" step all would have blocked the v1.3.0 cut otherwise. `golangci-lint run` exits 0.
- **`visualtest` shared-allocator compile error fixed.** `ensureAllocator` used `:=` at `doc.go:76`, shadowing the package-level `sharedAllocCtx`/`allocCancel` with unused locals — a compile error ("declared and not used") that broke the separate `visualtest` module and blocked the BuildFlow pre-commit `govalid-generate` step. Changed to `=` so the assignment targets the package vars consumed by `newTab`/`ShutdownBrowser`, exactly as the comment above it already instructed. Root cause: the auto-commit daemon landed the buggy `:=` without building the `visualtest` sub-module (it only builds the root module).
- **Flaky `TestStackWithFeedbackComponents` root cause fixed.** The integration test asserted `strings.Contains(output, "flex flex-col")`, but `utils.Class()` wraps tailwind-merge-go which reorders Tailwind classes **nondeterministically** (output depends on LRU cache state) — the test failed ~13% of the time under `-race`. Replaced the ordered-substring assertion with `utils.AssertContainsAll(output, "flex", "flex-col", "space-y-4")` in `integration/appshell_composition_test.go`. Verified 0/40 failures under `-race` (was 4/30). Ordered-substring assertions on `utils.Class` output are latent flakes repo-wide (see TODO #81).

## [1.2.0] — 2026-07-23

### Fixed

- **Popover/Dropdown top-layer positioning bug (ADR-0017 revision).** The original Popover API migration assumed CSS class-based positioning (`top-full left-1/2`) would continue to anchor to the trigger. **This was wrong:** `popover="auto"` promotes the panel to the top layer where the UA stylesheet forces `position: fixed; inset: 0`, detaching it from the trigger's DOM subtree. CSS classes therefore resolved against the viewport, placing panels at the wrong location. Fixed via a shared singleton `popoverPositionJS` (in `display/shared.go`) that reads `getBoundingClientRect()` on `toggle` open and sets `style.left/top` with viewport clamping. Used by both `Popover` and `Dropdown`. `ContextMenu` already positioned via JS `inset` and was unaffected. ADR-0017 revised with a full explanation of the three approaches considered (Anchor Positioning, JS rect, hybrid) and why JS rect was chosen.
- **Tooltip `aria-describedby` propagation restored.** The v0.20.0 Popover migration deleted the singleton script that propagated `aria-describedby` from the non-focusable wrapper `<div>` to the first focusable child (button/link/input). Screen readers therefore stopped announcing tooltip text on focus. Fixed via `tooltipAriaJS` singleton that re-runs on load and `htmx:afterSettle`. Tooltip show/hide remains pure CSS.
- **`HTMXSrc` now suppresses CDN response-targets extension.** Setting `PageProps.HTMXSrc` to self-host htmx previously still loaded the response-targets extension from the CDN (default `HTMXResponseTargets: true`), defeating half the purpose of self-hosting. The condition now checks `props.HTMXSrc == ""` — self-host implies you manage extensions. Godoc on both fields updated.
- **`tc add` now warns about companion `.go` dependencies.** A `.templ` file references package-level helpers (class lookups, enums, sub-templates) defined in sibling `.go` files that are not embedded or copied. The CLI now prints a clear note after copy pointing the consumer to `go get` the full package for a working component.
- **Popover entrance/exit animations** via `@starting-style` + `allow-discrete` in `templates/custom.css`. Panels fade + scale in/out gracefully. Browsers without `allow-discrete` support snap instantly.

### Added

- **Demo routes for recipes** — `/recipes/dashboard`, `/recipes/settings`, `/recipes/login` in `examples/demo`. The three recipes are now visually showcaseable, not just documented.
- **`navigation.SidebarNav` — vertical sidebar navigation for admin panels and dashboards.** Renders a brand slot (top), nav links (each with an optional `icons.Name` icon), and a footer slot (bottom). `CurrentPath` auto-detects the active item via the shared `IsActive` matcher (unified with `NavLink`); explicit `Active=true` on an item takes priority. Permanently-dark surface for the persistent-sidebar pattern.
- **`DOMAIN_LANGUAGE.md` platform terms** — ContainerAware, Recipe, Semantic Token, Theme Preset, HTMXSrc, Popover API, tc CLI added to the glossary.
- **`ROADMAP.md` reconciled** — v1.0 marked SHIPPED, v1.1+ platform work documented, headless variants moved to Explicitly NOT Planned.

### Changed

- **`go-error-family` bumped to v0.8.0.** Internal dependency update — no API changes affect this library.
- **`.golangci.yml` — `depguard` linter removed.** The allow-list (stdlib + the three runtime deps) duplicated the module's existing import discipline with no additional safety; maintaining it added churn on every dependency change.

## [1.1.0] — 2026-07-21

### Added

- **`tc` CLI scaffolding tool.** New binary at `cmd/tc/` — install with `go install github.com/larsartmann/templ-components/cmd/tc@latest`. Three commands: `tc init` (scaffold `app.css` + `custom.css`), `tc ls` (list every component by package), `tc add <component> [--out DIR]` (copy a component's `.templ` + `_types.go` source to your project for customization). Embeds the library's `.templ` sources at build time via `go:embed`. New `docs/cli.md` documents the full usage. Non-breaking addition.
- **Two new ADRs**: [`0020-per-package-modules-split.md`](docs/adr/0020-per-package-modules-split.md) (modules split — **deferred until consumer demand**, full design documented) and [`0021-headless-variants-defer.md`](docs/adr/0021-headless-variants-defer.md) (headless variants — **deferred indefinitely**, three options evaluated and rejected).
- **CI lint scope narrowed** — `golangci-lint run` now takes an explicit package list (`./display/... ./forms/...` etc.) excluding `cmd/tc/`. The CLI tool uses different conventions (intentional path traversal for `--out`, globals for embed registry, Print to stdout) than library code. Documented in AGENTS.md lint section.

## [1.0.0] — 2026-07-21

### Removed (breaking — v1.0)

- **`display.ModalSizeFull`** — deprecated alias for `ModalSize2XL` (same `max-w-4xl` width). Replace with `display.ModalSize2XL`.
- **`display.DrawerFull`** — deprecated alias for `DrawerSize2XL` (same `max-w-2xl` width). Replace with `display.DrawerSize2XL`.
- **`errorpage.FamilyFromErrorFamily`** — deprecated alias for `FromErrorFamily`. Replace with `errorpage.FromErrorFamily`.
- **`forms.FormProps.Inline bool`** — deprecated legacy field. Replace with `Layout: forms.FormLayoutInline`.

See [`docs/migration/v0.22-to-v1.0.md`](docs/migration/v0.22-to-v1.0.md) for the full migration guide.

### Added

- **`errorpage.ErrorPageProps.Validate() error`** — verifies Family is valid, StatusCode (when set) is in `[400, 599]`, and the page has at least one of Title/Message/CauseChain. Opt-in; the renderer does not call it automatically. Call it explicitly in handlers where you want to catch programmer errors.
- **CI docs-health drift guard** — `.github/workflows/ci.yaml` now runs `go test ./utils/... -run TestDocsCountDrift` on every push and PR. Catches drift between component/generated/enum counts in FEATURES.md, AGENTS.md, SKILL.md, and website sections.ts.

### Changed

- **Default flip deferred to v2.0.** The planned v1.0 flip of HTMX self-host (CDN → self-host default) and semantic tokens (opt-in → default) is deferred to v2.0. Rationale: both shipped opt-in in v0.22.0 and a single minor cycle is insufficient deprecation time. Consumers upgrading to v1.0 see no change to default behavior.

## [0.22.0] — 2026-07-21

### Added

- **Semantic token layer (ADR-0008 implementation).** New `templates/templ-components-theme.css` aliases every Tailwind palette color used by the library (`blue-600`, `red-600`, `green-600`, etc.) to semantic tokens (`--color-tc-primary`, `--color-tc-danger`, `--color-tc-success`, etc.). `templates/app.css` now imports it by default. Consumers override one token (`--color-tc-primary: #4f46e5`) to re-skin every component — buttons, links, focus rings, toasts — without touching any `.templ source. New `docs/theming.md` documents the three-tier theming model (semantic tokens, direct palette override, per-component Class).
- **Theme presets.** Three starter CSS files in `templates/presets/`: `default.css` (library stock), `minimal.css` (muted slate/emerald palette, reduced saturation), `glass.css` (violet + frosted-glass surface helper). Import one from `app.css` to switch palette without redefining every token.
- **Self-host HTMX opt-in (ADR-0007 implementation).** New `PageProps.HTMXSrc string` field switches the htmx main script from the CDN URL to a self-hosted path (e.g. `/static/htmx.min.js`). When set: the CDN preconnect is skipped, no SRI hash is emitted (same-origin), and `HTMXVersion`/`HTMXCDN`/`HTMXUseSRI` are ignored for the main script. The response-targets extension still uses CDN unless `HTMXResponseTargets: false`. CDN remains the default — flip to self-host as default is planned for v1.0.

## [0.21.0] — 2026-07-21

### Added

- **`recipes/` composition package (ADR-0019).** Three screen-level components that compose existing primitives into complete layouts: `recipes.Dashboard` (AppShell + PageHeader + StatCard grid + chart slots), `recipes.SettingsLayout` (Container + Split + Card stack with section anchors), `recipes.LoginCard` (centered Card with form body + OAuth divider + footer). Every variable part is a `templ.Component` slot — recipes provide the scaffold only. Resolves TODO #31 ("blocks/composition examples," deferred since v0.12). New `docs/recipes/{dashboard,settings,login}.md` with copy-paste examples.
- **Container-aware components (ADR-0018).** `navigation.Nav.ContainerAware bool` makes the bar collapse to its hamburger menu based on its parent container width (`@container` + `@sm:`/`@lg:`) instead of the viewport. `display.Card.ContainerAware bool` swaps padding breakpoints `sm:` → `@sm:` the same way. Both default to `false` (viewport breakpoints) — zero behavior change unless opted in. Mirrors the existing `display.Grid.ContainerResponsive` pattern.
- Two new ADRs: [`0018-container-query-native-contract.md`](docs/adr/0018-container-query-native-contract.md) and [`0019-recipes-package.md`](docs/adr/0019-recipes-package.md).
- Component count: 98 → 101. Package count: 14 → 15 (recipes added). Generated file count: 87 → 90.

## [0.20.0] — 2026-07-21

### Changed

- **Popover API migration (ADR-0017).** `Dropdown`, `Popover`, and `ContextMenu` now use the native HTML `popover` attribute (`popover="auto"`) with declarative `popovertarget` invokers. This deletes ~70 lines of singleton JS across the three components and replaces custom click-outside / Escape / focus-management logic with native browser behavior (light-dismiss, top-layer rendering, focus restore). `Popover` is now **zero-JS** — entirely native. `Dropdown` keeps a thin singleton (~25 lines, down from ~50) for WAI-ARIA menu keyboard nav (ArrowUp/Down + RTL mapping) and first-item focus on open. `ContextMenu` keeps a thin singleton (~6 lines, down from ~12) for the `contextmenu` event → `showPopover()` call. Browser support: Popover API is Baseline 2024 (Chrome 114+, Safari 17+, Firefox 125+). See `docs/adr/0017-popover-api-migration.md` and `docs/research/popover-api.md`.

### Removed

- **`Tooltip` singleton script** — the `tooltipScriptComponent` / `tooltipJS` / `window.tcTooltipAttached` machinery in `display/shared.go` is deleted. `Tooltip` is now pure CSS (it already used `:hover` / `:focus-within`; the JS only added touch click-toggle, Escape-dismiss, and `aria-describedby` propagation). Trade-off: tooltips no longer toggle on touch-tap, and `aria-describedby` is no longer auto-propagated to focusable children — consumers should set `aria-describedby` directly on the focusable trigger element. Tooltips remain hover/focus-only on desktop, which matches their role as progressive enhancement. Documented in ADR-0017 and in the `TooltipProps` godoc.

## [0.19.0] — 2026-07-21

### Fixed

- **`.golangci.yml` lint gate repaired** — the three linters documented as "disabled" in AGENTS.md (`ireturn`, `godoclint`, `testableexamples`) were never actually removed from the `enable:` list, producing 68 findings that would have failed CI on the next push. Removed all three from the enable list and deleted the dead `ireturn:` settings block. `golangci-lint run` now exits 0 as the v0.18.1 CHANGELOG originally claimed.

### Added

- **`layout.AppShell` — the #1 admin dashboard pattern, now a first-class component.** Sidebar + sticky header + main content shell using `lg:grid lg:grid-cols-[var(--tc-sidebar-w)_minmax(0,1fr)] min-h-dvh`. The sidebar is `hidden lg:block`; an optional `MobileNav templ.Component` slot accepts any component (typically `display.Drawer`) for mobile navigation — keeping the `layout` package's deps minimal (`utils`, `icons` only; no `layout → display` import cycle). `SidebarWidth` typed enum (SM/MD/LG/Auto) drives the `--tc-sidebar-w` CSS variable. Renders _inside_ `Base`'s `<main>` — does NOT emit its own `<main>` or skip-link (nested `<main>` is invalid HTML; Base owns both per WCAG 2.4.1 Bypass Blocks). Recipe: `docs/recipes/appshell-dashboard-layout.md`.
- **`layout.Container` — centered max-width wrapper, replacing the `max-w-Nxl mx-auto px-4 sm:px-6 lg:px-8` snippet repeated in every consumer.** Typed `ContainerWidth` enum (SM/MD/LG/XL/Full/Prose) with `ContainerWidthIsValid()`. Optional responsive padding (`containerPadClass`). All class strings go through `utils.Class()` so consumer Tailwind overrides win.
- **`layout.Split` — 2-column content+aside for article+sidebar and detail+metadata layouts (the pattern flex cannot do well).** Typed `SplitRatio` (1To2/1To3/1To4) + `AsidePosition` (Start/End) with logical CSS positioning that mirrors automatically in RTL. Both columns get `min-w-0` as the flex/grid complement to `minmax(0,1fr)` (blowout guard). Aside uses `<aside>`; Main uses `<div>` (NOT `<main>` — Base owns the singleton main landmark).
- **`layout.Stack` — vertical rhythm with a typed `StackGap` enum (None/SM/MD/LG/XL).** Replaces the repeated `space-y-N` string pattern. Deliberately `flex flex-col`, NOT grid — Stack is a 1D layout and the grid-first rule (ADR-0016) reserves grid for genuine 2D layouts.
- **`navigation.Footer` multi-column mode** — `FooterProps.Columns []FooterColumn` renders `grid grid-cols-2 md:grid-cols-4 gap-8` with an optional brand block and `BottomBar templ.Component` slot (legal/copyright row). Fully backward-compatible: empty `Columns` renders the legacy single-row footer unchanged. `FooterLink` and `FooterColumn` types added.
- **`forms.Form` `Layout` enum (Stack/Inline/Grid)** — typed form layout selection. `FormLayoutGrid` emits `sm:grid-cols-[auto_minmax(0,1fr)]` for aligned settings forms (label column + input column). The legacy `Inline bool` field is preserved and soft-deprecated (`Layout` wins when both are set); removal targeted for v1.0. `FormLayoutIsValid()` included.
- **ADR-0016 — "grid = 2D, flex = 1D"** (`docs/adr/0016-grid-first-for-2d-layouts.md`) — codifies the library's layout-selection rule with a new-component checklist and the full flex-usage audit (48/48 existing flex usages retained as correctly 1D; 0 migrations). All flexible grid columns MUST use `minmax(0,1fr)`, never bare `1fr`.
- **`docs/recipes/grid-blowout-minmax.md`** — documents the grid-blowout bug (bare `1fr` + wide content = page-wide horizontal scroll), the `minmax(0,1fr)` fix, where the library enforces it, the `min-w-0` flex complement, and manual testing guidance.
- **`docs/research/css-subgrid.md`** — research note tracking CSS subgrid (Baseline since 2022-2023), what it would unlock (Card alignment, DefinitionList, Form fields, CSS tables), and the decision to track-only (no demand, API complexity, testing gap).
- **6 new typed enums + `IsValid()` methods + lookup maps**: `ContainerWidth`, `SidebarWidth`, `StackGap`, `SplitRatio`, `AsidePosition`, `FormLayout` — all following the existing `GridCols`/`GridGap` pattern. Component count 94→98; typed enum count 37→43; generated file count 82→87.

### Deprecated

- **`forms.FormProps.Inline bool`** — superseded by `FormProps.Layout` (`FormLayoutInline`). Still works (maps to `FormLayoutInline` when `Layout` is empty), but `Layout` wins when both are set. Removal targeted for v1.0. Pass `Layout: forms.FormLayoutInline` in new code.

- **`docs/recipes/horizontal-filter-bar.md` — two HTMX filter-bar footguns documented**, plus a production-grade pattern. Both bugs are silent (the filter bar looks correct until a user combines filters in a specific order) and were hit in production by DiscordSync (the heaviest `templ-components` consumer): (1) **Silent query-param drop** — HTMX serializes only controls that exist in the form, so a param your handler reads but that has no form control (e.g. `author_id` arriving via a row-click link rather than a visible `<select>`) is silently dropped the moment the user changes _any_ filter; fix is a hidden `<input>`. (2) **Checkbox trigger bug** — the natural `hx-trigger="change"` (or `change from:find select`) fires only for `<select>` changes, so a `filterCheckbox` toggle does nothing; fix is `change from:find select, change from:find input[type=checkbox]`. The new "Production-grade pattern" section shows the hardened `filterBar` shape (robust `hx-trigger`, `hx-disabled-elt` to prevent stacking changes on a pending request, a Reset link, and an `htmx-indicator` `Spinner`). Footgun #1 is fundamentally **not enforceable in a library component** — only the consumer knows which params their handler reads — so this is a doc fix rather than a new `forms.FilterForm` component (which would also duplicate the existing `forms.FilterDropdown` and contradict the recipe's "thin custom helper" guidance).

## [0.18.1] — 2026-07-18

### Fixed

- **`scripts/release.sh` — 3 release-integrity defects repaired** (identified in the v0.18.0 release postmortem, same-day fix). (1) Release commit body no longer duplicates the one-line summary as its first line — the body now carries the multi-paragraph release notes (`${RELEASE_NOTES}`), while the subject carries the summary. (2) Model attribution is no longer hardcoded to `MiniMax-M3` — now reads `${CRUSH_MODEL:-unknown}` so the `Assisted-by:` trailer is honest under any model. (3) The hostile stdin read loop (`while IFS= read -r line`) is gone — release notes now auto-extract from `CHANGELOG.md [Unreleased]` (the project's canonical source per the "[Unreleased] must be warm at all times" rule), - **`scripts/release.sh` — 3 more release defects repaired** (surfaced when cutting v0.18.1; the script failed verify and left the tree dirty). (4) `FEATURES.md` version was never bumped — the script bumped `utils.Version` and `CHANGELOG.md` but not `FEATURES.md`, despite AGENTS.md mandating all three move together and `utils.TestVersionMatchesFeatures` enforcing it; now bumps both `**Version:**` and `**Updated:**` date. (5) A failed `go test ./...` aborted via `set -e` with **no rollback**, leaving a dirty tree (only the changelog drift-guard had rollback); now an `EXIT` trap (`release_rollback`) restores all three version files on any failure between bump and commit. (6) Replaced the banned `git checkout --` with `git restore` (AGENTS.md prohibits `git checkout`). The drift-guard test now asserts 6 invariants (the original 3 plus these).
- **golangci-lint gate repaired** — commit 73395d9 expanded the linter set to 67 but never reconciled the config with the codebase, leaving `golangci-lint run` failing with 187 findings (CI would have gone red on the next push). Reconciled by: replacing the non-resolving `$module` depguard token with an explicit module + dependency allow-list (`a-h/templ`, `Oudwins/tailwind-merge-go`, `larsartmann/go-error-family`); extending `varnamelen` ignore-names with idiomatic Go short names (`w`/`r`/`tc`/`id`/`s`/`p`/`bp`/`cp` etc.) and excluding test files from the linter; adding HTTP status codes + unix file permissions to `mnd` ignored-numbers; disabling three linters fundamentally incompatible with a templ library (`ireturn` — every component returns the `templ.Component` interface by design; `godoclint` — demands a single godoc per package but the repo documents per-file; `testableexamples` — HTML-rendering `Example` funcs aren't deterministic); excluding `err113`/`makezero`/`gocheckcompilerdirectives` from test files; naming the `ComponentProps.SetBaseProps(props BaseProps)` interface parameter (`inamedparam`); and bumping `funlen` line limit 60→65 for the `ErrorHandler` dispatch. Result: `golangci-lint run` exits 0; the remaining 64 linters enforce cleanly.

### Added

- `scripts/release.sh --notes-file FILE` flag — accepts a path to a markdown file as the release-notes source, overriding the default CHANGELOG `[Unreleased]` extraction. Useful for scripted releases or when notes are authored separately.
- `utils/release_script_test.go` (`TestReleaseScriptInvariants`) — static-analysis drift guard that reads `scripts/release.sh` as text and asserts the 6 fixed invariants hold (body uses `${RELEASE_NOTES}` not `${RELEASE_SUMMARY}`; no hardcoded `MiniMax-M3`; no stdin read loop; `--notes-file` and `[Unreleased]` extraction present; `FEATURES.md` version bump present; `trap release_rollback EXIT` rollback installed; no banned `git checkout --`). Fails CI if any defect creeps back.
- `flake.nix` adopted `treefmt-nix` (mirrors `website/flake.nix`) — `nix fmt` formats `.nix` (nixfmt) + `.go` (gofmt + goimports); `nix flake check` runs the format verification as a CI gate. Generated `*_templ.go` files, `website/`, and `examples/demo/static/` are excluded to avoid churn against the generator. Replaces the former bare `formatter = pkgs.nixfmt;`. BuildFlow still owns the pre-commit hook.
- README Quick Start now documents the `GOEXPERIMENT=jsonv2` build flag required until Go 1.27 stabilizes it.

### Changed

- `.art-dupl-baseline.json` regenerated: the stale baseline recorded 17 clone groups from 2026-06-28, but the codebase now has 0 (refactored away). Regenerated to 0 groups so the `art-dupl check` "no new clones" CI gate reflects reality.
- `docs/icons-only-adoption.md` icon count corrected from 101 to 102 (101 path-icon `Name` constants + 1 animated Spinner; 5 discoverability aliases like `Close`→`X` resolve to canonical paths). Broken markdown list (orphaned Spinner bullet) also fixed.
- `TODO_LIST` #62 rescoped: "Add `Validate()` to top 5 props structs" was over-engineering — now scoped to `errorpage.ErrorPageProps` only (the one struct where invalid `StatusCode`/`Family` combos produce wrong HTTP responses; other props use graceful `utils.Lookup` fallback and need no `Validate`).
- v0.18.0 release postmortem annotated with `## Resolution (2026-07-18)` appendix answering the 3 open questions (Q1: don't rewrite published tag history; Q2: keep `Assisted-by:` but parameterize; Q3: CHANGELOG `[Unreleased]` is the canonical notes source) and recording what shipped.
- `AGENTS.md` Build & Test section: added Nix flake commands subsection documenting `nix fmt`, `nix flake check`, and the flake apps.

## [0.18.0] — 2026-07-18

Native `<dialog>` for Modal/Drawer, stylable `<select>`, auto-growing Textarea, responsive Image srcset, semantic `<search>` landmark, `hx-validate` on Form, and `content-visibility: auto` for large tables. All progressive enhancement — zero breaking changes.

### Changed

- **Custom CSS extracted to `templates/custom.css`.** Component-specific styles (dialog animations, stylable select, auto-grow textarea, scroll-snap, accordion chevron, accent-color) moved from inline in `templates/app.css` to a separate `templates/custom.css` file. Both `templates/app.css` (consumer entry point) and `examples/demo/demo.css` import it via `@import "./custom.css"`. Consumers must now copy both `app.css` and `custom.css` to their project.
- **Modal/Drawer migrated to native `<dialog>` element.** Replaces ~200 lines of custom focus-trap JS with browser-native `showModal()`/`close()`. Benefits: native focus trap, Escape-to-close, focus restore, top-layer rendering (no z-index wars), `::backdrop` dimming, automatic background `inert`. The `.tc-overlay`/`.tc-modal`/`.tc-drawer` CSS classes in `templates/custom.css` handle animations via `@starting-style` + `allow-discrete`. `tcOpenOverlay(id)`/`tcCloseOverlay(id)` wrappers kept for backward compatibility.

### Added

- `display.CardProps.TitleClass` and `HeaderClass` — consumer-supplied classes that override the default `<h3>` title and header wrapper classes respectively. Useful when the default header layout works but the typography or spacing needs adjustment without replacing the entire header via the `Header` slot.
- Demo site overhaul: new hero section with library stats, copy-to-clipboard install command, and documentation link; sticky section navigation with live section filtering; copy-paste Go code snippets for representative components; searchable icon gallery with empty-state feedback.
- `forms.SelectProps.Stylable` — opts into the modern customizable `<select>` API (`appearance: base-select`). Emits `<button><selectedcontent></selectedcontent></button>` structure for full CSS styling of the button, dropdown picker, options, checkmark, and arrow icon. Progressive enhancement: non-supporting browsers (Firefox, iOS Safari) ignore the structure and render a normal native `<select>`. Requires `.tc-select` CSS from `templates/custom.css`.
- `forms.TextareaProps.AutoGrow` (default `true`) — uses CSS `field-sizing: content` to auto-grow the textarea to fit content without JavaScript. Set `AutoGrow: false` to disable.
- `forms.TextareaProps.EnterKeyHint` + `forms.EnterKeyHintType` typed enum — sets the mobile keyboard's Enter key label. Constants: `EnterKeyHintSend` (chat), `EnterKeyHintDone`, `EnterKeyHintGo`, `EnterKeyHintNext`, `EnterKeyHintPrevious`, `EnterKeyHintSearch`, `EnterKeyHintEnter`. `EnterKeyHintTypeIsValid` included. `forms.InputProps.EnterKeyHint` uses the same enum — explicit override takes priority; otherwise auto-derived from `InputType` (email→next, search→search, etc.).
- `forms.FormProps.Validate` — when true, adds `hx-validate="true"` for HTML5 constraint validation before HTMX submit.
- `forms.Input` with `Type: InputSearch` now wraps the input in a `<search>` element — a semantic landmark for screen reader navigation. No API change needed; auto-detected from `InputType`.
- `display.ImageProps.SrcSet` and `Sizes` — typed fields for responsive image delivery (`srcset`/`sizes` attributes). Replaces the previous `Attrs` workaround documented in the old doc comment.
- `display.TableProps.LazyRows` — when true, applies `content-visibility: auto` to body rows. Browser skips rendering off-screen rows, giving 2-5x faster initial render for tables with 100+ rows.
- CSS for stylable `<select>` in `templates/app.css` — comprehensive `.tc-select` styles with light/dark mode, hover/focus states, custom dropdown arrow (`::picker-icon`), styled picker container (`::picker(select)`), option hover/selected states. Guarded by `@supports (appearance: base-select)`.
- Global `accent-color` in `templates/app.css` — all native form controls (checkboxes, radios, range inputs, progress bars) get the library's blue accent by default. Dark mode variant included. Override per-component with Tailwind `accent-*` utilities.

### Fixed

- `flake.nix` apps (`build`, `test`, `verify`, `coverage`) now export `GOEXPERIMENT=jsonv2` so the Nix-provided tooling matches the pre-commit hook and the library's `encoding/json/v2` usage.

## [0.17.0] — 2026-07-12

9 new components (DataTable, FilterDropdown, Slider, Rating, TagsInput, HoverCard, ContextMenu, Carousel, Calendar), 5 audit bug fixes, coverage tests, ADR 0013, demo /forms route.

### Added

- `display.DataTable` — data table with integrated sort management, optional pagination slot, and empty-state handling. Auto-generates sort-toggle URLs from `ActiveSortColumn`/`ActiveSortDir`/`SortBaseURL`. Composes `Table` internally.
- `forms.FilterDropdown` — compact select dropdown purpose-built for HTMX filter bars. Auto-submits via `hx-get`/`hx-target`/`hx-trigger` on change. `Value` pre-selects the current filter value.
- `forms.Slider` — labeled range input with Min/Max/Step, `ShowValue` display, error/help text. Dark mode compliant with `accent-blue` native styling.
- `forms.Rating` — star rating input using radio inputs for keyboard accessibility. `RatingSize` enum (SM/MD/LG), `ReadOnly` display mode, configurable `Max` stars. `RatingSizeIsValid` included.
- `forms.TagsInput` — tag input with add/remove via CSP-safe singleton JS. Hidden inputs for form submission. `MaxTags`/`AllowDuplicate` controls.
- `display.HoverCard` — CSS-only hover-activated card with 4 positions. Focus-within support, aria-describedby linkage. No JavaScript required.
- `display.ContextMenu` — right-click context menu with CSP-safe singleton JS. `role="menu"`/`role="menuitem"`, Escape and click-outside dismiss.
- `display.Carousel` — slide carousel with prev/next arrows and dot indicators. `aria-roledescription="carousel"`, singleton JS for slide control.
- `forms.Calendar` — month-view calendar with server-side navigation. Day links via `{year}`/`{month}`/`{day}` placeholder substitution. `MinDate`/`MaxDate` disabling.
- `display.HoverCardPositionIsValid` — typed enum validation for HoverCardPosition (top/bottom/start/end).
- `navigation.BreadcrumbsProps.CurrentPath` — auto-detects active breadcrumb from the current request path, matching the NavLink/SidebarNav pattern. Explicit `Active=true` takes priority.
- Demo `/forms` route — standalone form showcase at `http://localhost:8080/forms` with validation, all form components, and an HTMX filter bar example.
- ADR 0013: `encoding/json/v2` auto-formatter guard — documents why only `errorpage` uses v2 and how the pre-commit hook prevents accidental rewrites.

### Fixed

- `navigation.MobileMenu`/`MobileMenuToggle` double-prefix: `EnsureID("mobile-menu", ...)` already produces `tc-mobile-menu-<hex>`, but the template prepended `tc-mobile-menu-` again, producing `tc-mobile-menu-tc-mobile-menu-<hex>`. Removed the redundant prefix from both `aria-controls` and `id`.
- `layout.ThemeToggle` stale `aria-checked` after HTMX swap: the singleton guard prevented the `syncToggleAria` call from running on re-rendered buttons. Moved the sync call outside the guard so it runs on every render.
- `htmx.GlobalErrorHandling` retry used `.click()` instead of `htmx.trigger()`: non-click-triggered elements (e.g., `hx-trigger="change"`) wouldn't retry. Now reads the element's `hx-trigger` attribute and fires the matching event.
- `forms.RadioGroup` did not propagate `aria-invalid`/`aria-describedby` to individual radio `<input>` elements: screen readers couldn't announce the error state. Now passes error attributes through to each radio via `radioItemProps`.
- `navigation/end_of_list_test.go` ordered class assertion: `utils.Class()`/tailwind-merge reorders classes, causing a flaky test. Changed to `AssertContainsAll` for ordering-independent checks.

## [0.16.0] — 2026-07-12

Flush prop eliminates table-in-card double border; CellPadding typed enum adds compact density for admin dashboards; ADR 0012 documents the nesting convention.

### Added

- `display.TableProps.Flush` — when true, suppresses the Table wrapper div's border and rounded corners. Use when nesting a Table inside a `Card(CardPaddingNone)` to eliminate the double-border defect. The `overflow-x-auto` scroll wrapper is always retained.
- `display.TableCellPadding` typed enum (`TableCellPaddingComfortable` / `TableCellPaddingCompact`) + `TableCellPaddingIsValid`. Controls vertical density of header and body cells. Comfortable (px-4 py-3) is the default; Compact (px-4 py-2) suits data-heavy dashboards and admin panels.

### Fixed

- Table-in-Card double border: nesting `Table` inside `Card(CardPaddingNone)` previously rendered two concentric borders (one from each component). Set `Flush: true` on the Table to suppress its wrapper border — the card provides the outer border and the table sits flush against it.

## [0.15.0] — 2026-07-11

GridColsAutoFit + MinColWidth for container-responsive dashboard grids. Card.Header slot for custom card headers. CardPaddingNone fix for table-in-card. Dark mode packaging fixes: color-scheme: light dark, split @custom-variant into opt-in. Adoption guide rewrite with three dark mode paths, @theme palette override pattern, BaseProps.Class gotcha docs. Adopted encoding/json/v2 — consumers need GOEXPERIMENT=jsonv2. Also includes prior session features: TrendWarn, TableRow.Href, Select optgroups, EndOfList, theme-bridge recipe.

### Added

- `display.GridColsAutoFit` + `GridProps.MinColWidth` — CSS `auto-fit`/`minmax()` grid template for container-width-responsive layouts. Common dashboard pattern that previously required the `Class` escape hatch. Takes precedence over `ContainerResponsive` when both are set (auto-fit already responds to container width via CSS).
- `display.CardProps.Header` — `templ.Component` slot that replaces the entire default header section (title, subtitle, header action). Enables custom card headers without full `Body`-slot replacement. When nil, the default header renders as before.
- `docs/dark-mode-research.md` — first-principles analysis of dark mode mechanisms in Tailwind v4 + modern CSS. Documents three first-class consumer paths (OS-following, toggle, CSS-variable design system) with code examples and a comparison table.
- `docs/tailwind-v4-adoption-guide.md` — rewritten with "Setting Class on components" section (Go promoted-field struct literal gotcha), Go module cache `@source` path pattern, top-level "Theming" section with `@theme` palette override pattern, and "Default constructors" section.

### Fixed

- `templ-components-theme.css`: `color-scheme: light` → `color-scheme: light dark` on `:root` — native form controls (scrollbars, checkboxes, date pickers) were rendering light-only even in dark mode.
- `templ-components-theme.css`: `@custom-variant dark` moved to a clearly commented opt-in section. Previously, importing the theme file for color overrides silently forced class-based dark mode, overriding Tailwind v4's default `prefers-color-scheme` strategy.
- `display.CardProps.CardPaddingNone`: children/body now render directly inside the card shell without the wrapping padding `<div>`. Enables table-in-card layouts where `<table>` must be a direct child for `overflow-x-auto` to work correctly.
- `display.Grid` — `ContainerResponsive` no longer silently overrides `GridColsAutoFit`. When both are set, auto-fit takes precedence (it already responds to container width via CSS `minmax()`).
- `layout.ThemeScript` / `layout.ThemeToggle` godoc: clarified they are only needed for the toggle dark mode strategy, not for OS-following.

### Changed

- `display.GridColsIsValid` now recognizes `GridColsAutoFit` as a valid value.
- **Adopted `encoding/json/v2`** — the library now uses `encoding/json/v2` + `encoding/json/jsontext` in `errorpage`. Consumers must set `GOEXPERIMENT=jsonv2` (stable in Go 1.27). Deleted the `utils/jsonv2_guard_test.go` prohibition test and removed the pre-commit grep guard.
- `display.TrendWarn` — amber trend variant for StatCard (the #1 gap reported by DiscordSync, which previously mapped "warn" to `TrendDown` with an apology comment). Uses `text-amber-600 dark:text-amber-400` + right-pointing arrow icon, sr-only label "Holding at".
- `display.TableRow.Href` — clickable table rows. When set, the row gets `data-tc-row-href`, `role="link"`, `tabindex="0"`, `cursor-pointer`, and hover styling. A CSP-safe singleton script handles click navigation and keyboard support (Enter/Space). Clicks on interactive elements inside the row (links, buttons, inputs) are not hijacked. Replaces the `data-href` JS workaround used across 10+ pages in DiscordSync.
- `forms.SelectGroup` + `forms.SelectProps.Groups` — optgroup support for Select. When `Groups` is non-empty, options render as `<optgroup>` elements instead of a flat list. Each group's options go through the same `normalizeSelectOptions` normalization (Disabled+Selected clearing, single-Selected enforcement). Replaces the custom `channelGroupedSelect` in DiscordSync.
- `navigation.EndOfList` — "You've reached the end" indicator for the bottom of a list. Companion to `LoadMore` and `Pagination`. Customizable `Message`, `role="status"` for a11y, `text-gray-500 dark:text-gray-400` muted styling.
- `docs/recipes/theme-bridge.md` — documentation showing how consumers with custom semantic palettes (e.g. `bg-surface`, `bg-accent`) can remap standard Tailwind color tokens via `@theme` to use library components without forking. Replaces ~250 lines of reimplementation in SwettySwipperWeb.

## [0.14.0] — 2026-07-10

### Added

- `display.Popover` — new component: button-triggered floating panel with arbitrary content (children slot), 4 positions (`PopoverPositionTop/Bottom/Left/Right`), `role="dialog"`, `aria-haspopup`/`aria-controls`/`aria-labelledby` wiring, click-outside and Escape dismissal, CSP-safe singleton JS. The #1 most requested missing component across consumer feedback.
- `PopoverPosition` typed enum with `PopoverPositionIsValid()` + test coverage.

### Changed

- Icons coverage improved from 47.1% → 75.9% via `TestIconRTL` covering the shipped-but-untested `IconRTL` public API (all 100 path icons + spinner variant).

## [0.13.0] — 2026-07-10

### Fixed

- `htmx/InlineLoadingOverlay` — added `<span class="sr-only">Loading…</span>` for screen reader parity with `LoadingIndicator`.
- `forms/ValidationSummary` — error links now use raw `err.Field` instead of `SanitizeID(err.Field)`, so anchors match actual element IDs.
- `errorpage/FromError` — unknown errors now fall back to `FamilyCorruption` (→500) instead of `FamilyInfrastructure` (→503). An unknown error is more likely a bug than a temporary outage.
- `errorpage/ErrorPage` + `errorpage/NotFound404` — changed outer container from `<div role="region">` to `<main>` landmark (WCAG 2.4.1 Bypass Blocks).
- `forms/Form` — added `CSRFTokenName` field (defaults to `"csrf_token"`) so frameworks with different token names (Django `_csrf_token`, Rails `authenticity_token`, Spring `_csrf`) don't need to fork.

### Changed

- `navigation/Footer` — signature changed from `Footer(brandText string)` to `Footer(props FooterProps)`. `FooterProps` embeds `BaseProps` (Class/ID/Attrs/AriaLabel) matching every other component. **Breaking:** update callers to `Footer(FooterProps{BrandText: "MyApp"})`.
- Pre-commit hook (`scripts/pre-commit.sh`) — lint now uses `golangci-lint run ./...` (examples/ excluded via `.golangci.yml`), and includes a grep guard rejecting `encoding/json/v2` imports.
- AGENTS.md — renamed stale "Post-v0.9.0 Conventions" section to "Conventions"; fixed lint path typo (`./svg/...` → `./...`); corrected generated file count (61→62); added `encoding/json/v2` prohibition section.

## [0.12.1] — 2026-07-10

### Fixed

- `feedback/Toast` + `feedback/Alert` + `errorpage/ErrorAlert` — `tcToastColors`, `tcToastIcons`, and `DismissScript` were JSON-marshaled by templ's `ScriptContentOutsideStringLiteral` instead of emitted as raw JS. The `var` declarations became discarded string literals, so the variables were never declared. Fix: render entire `<script>` tags via `templ.Raw()` component, bypassing the JSON marshaling path entirely.

## [0.12.0] — 2026-07-09

### Fixed

- `navigation/Nav` — desktop link container now wraps gracefully when the number of links exceeds the viewport width. Changed from fixed `h-16` to `min-h-16 items-center` (grows to accommodate wrapped rows) and from `sm:space-x-8` to `sm:flex-wrap sm:gap-x-8 sm:gap-y-2` (wraps instead of overflowing off-screen). Works for any number of links on any screen width without JavaScript or horizontal scroll.

## [0.11.0] — 2026-07-09

### Added

### Added

- `templates/app.css` — ready-to-copy starter CSS entry-point with `@import "tailwindcss" source(none)`, project + vendored `@source` directives, `@custom-variant dark`, and commented `@theme`/`@import` blocks.
- BuildFlow `tailwind-build` provider — auto-discovers CSS entry-point files and compiles them via `tailwindcss` as part of BuildFlow's DAG, ordered after `go-mod-vendor` and `templ-generate`.
- `tailwindcss_4` added to templ-components devShell (`flake.nix`) so the binary is available without nix fallback.
- Documentation: README "Tailwind CSS Setup" section simplified to BuildFlow + starter template options. Adoption guide and migration guide updated.

### Fixed

- `errorpage/handler.go` — reverted accidental `encoding/json/v2` import to `encoding/json` (was introduced by an auto-formatter under GOEXPERIMENT=jsonv2; this repo does not use that flag).
- `navigation/breadcrumbs_templ.go` — corrected import from `encoding/json/v2` to `encoding/json` (same root cause).

### Changed

### Removed

## [0.10.0] — 2026-07-08

### Added

- Naming: `icons.Close` alias for `icons.X` (both map to `"x"`). Prefer `Close` in new code.
- RTL keyboard mapping for `display.Tabs` and `display.Dropdown` — ArrowLeft/Right swap when `dir="rtl"` is set on `<html>`.
- Demo: SkeletonCardGrid loading state showcase, anchor-linked TOC at top.
- Documentation: `docs/adr/0010-sub-template-extraction-pattern.md` (when to extract/when to keep duplication).
- Documentation: `docs/migration/v0.8-to-v0.9.md` migration guide.
- README: "Further reading" table cross-linking javascript-guide, motion-design, container-queries, recipes, and ADRs.
- README: `GridProps.Gap` + `ContainerResponsive` examples, `FormProps.Inline` filter bar example.
- `htmx.GlobalErrorHandling`: enhanced godoc example showing ToastContainer wiring.
- Benchmark suites for `forms` (Input, Select, Textarea, Combobox), `layout` (ThemeScript, ThemeToggle, Script, Minimal), `htmx` (LoadingIndicator, CSRFToken, SwapOOB), `icons` (Icon, IconWithStrokeWidth, IconPathData, IconPathJS), `utils` (Class, EnsureID, Ternary, Lookup).
- Fuzz tests for `forms.InputType`, `forms.FormMethod`, `display.ButtonHTMLType`.
- Motion-reduce compliance grep test (asserts every `transition-*`/`animate-*` in `.templ` files has `motion-reduce:` fallback).
- SKILL.md component count drift-guard test (informational, logs actual vs documented count).
- Golden package coverage boost: 70.5% → 81.8% (update-flag, mkdir, normalization edge cases, diff, lineAt).
- Dedicated sub-template tests for `errorHeader`, `actionLinkBody`, `goBackScript`, `skeletonContainer`, `definitionDetailContent`.
- `display.Tabs` auto-generates IDs for tabs that omit them (`ensureTabIDs`) and defaults `ActiveTabID` to the first tab when unset (`resolveActiveTabID`) — prevents invalid HTML and ensures WAI-ARIA keyboard-focus compliance.
- `display.Tooltip` JS propagates `aria-describedby` from wrapper to the focusable trigger element so screen readers announce tooltip text.
- `display.Accordion` uses CSS grid technique (`grid-rows-[1fr]`/`grid-rows-[0fr]`) instead of `max-h-96` — content of any height animates correctly without clipping.
- `errorpage.ErrorPageProps.StatusCode` — explicit HTTP status code override. When set (non-zero), takes precedence over the family-derived default. `NotFound()` sets 404, `Forbidden()` sets 403, `InternalError()` sets 500.
- `forms.RadioOption.Checked` — enables pre-selecting a radio option for edit forms.
- `forms.RadioProps.Required` — propagates `required` to individual radio inputs for native HTML5 validation.

### Fixed (htmx, errorpage, layout, forms, navigation)

- **`htmx.LoadingButton`**: `htmx-hide-during-request` was not a real CSS class — default text never hid during loading. Replaced with Tailwind arbitrary variant `[.htmx-request_&]:hidden`.
- **`htmx.InlineLoadingOverlay`**: static `aria-hidden="true"` was never toggled. Replaced with `role="status"` + `aria-live="polite"`.
- **`htmx` retry counter**: was set on `event.detail.elt` but cleared on `event.detail.target`. Now clears from the same element.
- **`htmx` error announcer**: `#tc-error-announcer` aria-live region was rendered but never populated. Now updated with error messages.
- **`htmx` missing catch-all**: no default `else` left `undefined` values for uncovered status codes. Added fallback.
- **`htmx.ConfirmDelete`**: `hx-confirm` was always rendered even when empty. Now conditional.
- **`htmx.SwapOOB`**: empty `Selector` produced malformed attribute. Now omits selector when empty.
- **`errorpage` status codes**: `NotFound()` returned 400 (should be 404), `Forbidden()` returned 400 (should be 403), `InternalError()` returned 503 (should be 500). Added `StatusCode` field.
- **`errorpage` a11y**: `role="region"` added to `ErrorPage`/`NotFound404` root divs. `ErrorAlert`: empty message guarded. `contextTable`: caption + `th scope`.
- **`layout.ThemeToggle`**: `querySelectorAll` syncs all instances. localStorage wrapped in try/catch.
- **`layout` FOUC**: `ThemeScript` moved before HTMX CDN scripts. Favicon type attribute removed. SRI integrity conditional.
- **`forms.RadioGroup`**: `Required` now propagates `required` to radio inputs for native validation.
- **`forms.InputGroup`**: right addon missing `pointer-events-none` — was blocking clicks.
- **`forms.FieldError`**: added `role="alert"`. Empty message guarded.
- **`navigation.LoadMore`**: `aria-label` moved from div to button.
- **`navigation` breadcrumb URL**: uses `net/url.Parse` instead of naive string check.

### Fixed (forms, feedback, display, navigation)

- **`forms.Toggle`**: `peer-checked:translate-*` classes were dynamically concatenated (`"peer-checked:" + translateClass`) at runtime, making them invisible to Tailwind's content scanner. The thumb did not slide when checked in production. Now stores complete variant-prefixed class literals (`peer-checked:translate-x-5`).
- **`navigation.Pagination`**: arrow button border-radius was dynamically concatenated (`"rounded-" + side + "-md"`), invisible to Tailwind's scanner. Now passes complete logical-property literals (`rounded-s-md`/`rounded-e-md`) that also auto-mirror in RTL.
- **`forms.Combobox` disabled hidden input**: the hidden submission input was not disabled when `Disabled: true`, so its value was still submitted (violating the HTML spec's disabled-exclusion contract). Now both visible and hidden inputs get `disabled`.
- **`forms.Combobox` stale hidden value**: typing in the text input without selecting an option left the hidden input's value stale (the pre-populated server value was silently submitted instead of the user's typed text). The `input` event handler now clears the hidden value when the user types.
- **`forms.Combobox` Enter blocking form submission**: `e.preventDefault()` was called unconditionally for Enter, even when no option was highlighted. This blocked form submission when pressing Enter in the combobox. Now Enter only prevents default when an option is actively highlighted.
- **`forms.Select` slice mutation**: `normalizeSelectOptions` mutated the caller's `[]SelectOption` in place, corrupting `Selected`/`Disabled` flags on re-render. Now returns a defensive copy.
- **`forms.Select` doc contradiction**: type comment said "Selected takes precedence (Disabled is cleared)" but code clears Selected. Documentation corrected to match the implementation (Selected is cleared).
- **`forms.Checkbox` invalid `for=""`**: a checkbox without an ID rendered `<label for="">` (invalid HTML that breaks label association). Now renders a `<span>` when no ID is present.
- **`feedback.Toast` auto-dismiss**: a toast with `Duration > 0` but no ID silently disabled auto-dismiss (the `setTimeout` was gated on `props.ID != ""`). Now auto-generates an ID via `EnsureID` so `DefaultToastProps()` (which sets Duration: 5000) auto-dismisses correctly.
- **`feedback.ProgressBar` aria-valuenow**: `aria-valuenow` used the raw `props.Current` value without clamping, producing values outside `[aria-valuemin, aria-valuemax]`. Now clamped to `[0, Total]`.
- **`display.Modal`/`display.Drawer` aria-hidden/inert sync**: the JS open/close functions only toggled CSS classes but never synced `aria-hidden` or `inert`. A JS-opened modal stayed `inert` (keyboard inaccessible) and `aria-hidden="true"` (screen reader invisible). Now `tcOpen` removes `inert` and sets `aria-hidden="false"`; `tcClose` adds `inert` and sets `aria-hidden="true"`.
- **`display.Dropdown` RTL dead code**: the RTL arrow-key ternary was inside a JS string literal (`e.key === '(document... ? ...)'`), making it dead code that never matched. Now computes `nextKey`/`prevKey` as variables outside the comparison.
- **`display.CopyButton` link navigation**: clicking the `<a>` variant followed the `href` before the "Copied!" feedback could show. Now calls `e.preventDefault()` so copy feedback is visible.

### Changed

- Deduplication sprint: 6 sub-template extractions across `errorpage`, `display`, and `feedback` packages — `errorHeader`, `goBackScript`, `actionLinkBody`, `skeletonContainer`, `definitionDetailContent`, and merged `overlayPanel` into `overlayShell`. Reduces production clone groups from t=8→4.
- Coverage boost: 152 new test functions across 6 packages (display, feedback, forms, navigation, errorpage, layout) targeting untested branches.
- Renamed `forms/radio_go.go` → `forms/radio.go` (misleading `_go.go` suffix).
- Renamed `forms/aria.go` parameter `errMsg` → `errorMessage` (descriptive, no abbreviation).
- Renamed `errorpage/fromerror.go` `cleanMessage` → `sanitizeErrorMessage` (precise verb).
- Standardized `layout/sri.go` naming: `htmxMainSRIDefault` → `sriHTMXMainDefault` (consistent with `sriHTMXMainByVersion`).
- Extracted `msgGoBack` constant in `errorpage/constructors.go` (goconst compliance — 0 lint issues).
- `goBackScript` and `overlayShellProps` reviewed for promotion/restructure — both documented with trigger conditions.
- ADR 0009 rewritten with rigorous per-clone justification for 6 remaining accepted clone groups.

### Removed

- Deleted stale `origin/modularize/strategic-split` remote branch (abandoned experiment, never merged).

## [0.9.1] — 2026-07-08

> **Note:** This version was never tagged. All changes were included in v0.10.0.
> Consumers should use `@v0.10.0` or later; `@v0.9.1` will return a 404 from the Go module proxy.

### Added

- Dark mode compliance tests: `utils.TestDarkModeCompliance` (neutral colors) and `utils.TestDarkModeSemanticColors` (semantic colors) — scanning all `.templ`/`.go` source files for missing `dark:` variants. Failing tests block CI.
- `color-scheme: light` on `:root` and `color-scheme: dark` on `.dark` in `templ-components-theme.css` — improves native form control rendering (scrollbars, checkboxes, radios, date pickers) in dark mode.
- Dark mode focus-ring and ring-offset variants on all interactive elements (`dark:focus:ring-*`, `dark:focus-visible:ring-*`, `dark:focus-visible:outline-*`, `dark:focus:ring-offset-gray-900`).
- Dark mode shadow variants on overlays and cards (`dark:shadow-black/20`).
- `progressbar.templ` modernized to use `max()`/`min()` builtins (Go 1.21+) instead of manual if-branch clamping.
- Doc comments updated with `dark:` variants in all example code.

### Fixed

- 30+ missing `dark:` variants fixed across all packages — buttons, avatars, badges, tabs, pagination, sidebar, breadcrumbs, mobile menu, theme toggle, toast dismiss, step indicator, error page families, form inputs, and more.
- `errorpage/handler.go` reverted from `encoding/json/v2` to `encoding/json` — `json/v2` requires `GOEXPERIMENT=jsonv2` which is not enabled in Go 1.26.4.

## [0.9.0] — 2026-07-06

### Added

- `display.GridProps.Gap`: typed `GridGap` enum (`GridGapSM`/`MD`/`LG`/`XL` → `gap-2`/`4`/`6`/`8`) with `gridGapLookup` map + `GridGapIsValid`. Replaces hardcoded `gap-4` in grid lookup maps — consumers can now control spacing. Defaults to `GridGapMD` (`gap-4`), backward compatible.
- `display.CopyButtonProps.Href`: when set, renders an `<a>` instead of a `<button>`. The link still copies to clipboard on click.
- `display.ImageProps.Rounded`: when `true`, adds `rounded-full` instead of `rounded-md`. Quick convenience for avatars and icons.
- `navigation.LoadMoreProps.InfiniteScroll`: when `true`, adds `hx-trigger="revealed"` for auto-loading when scrolled into view (infinite scroll pattern).
- `errorpage.NotFound404Props.LinksTitle`: configurable heading for the quick-links section. Defaults to "Popular pages".
- `errorpage.WriteNotFound404`: convenience HTTP handler that writes a `NotFound404` page with 404 status code.
- Demo app: 5 new showcase sections — sortable Table (`TypedHeaders`), inline Form filter bar (`FormProps.Inline`), container query Grid (`ContainerResponsive`), 404 page preview (`NotFound404`), Table.Body slot.
- Documentation: `ROADMAP.md` (v0.x/v1.0/v2.0+ direction), rewritten `CONTRIBUTING.md` (Nix setup, conventions, release flow), `docs/migration/v0.7-to-v0.8.md` (all changes with before/after examples).
- Benchmark suite for `errorpage` package (ErrorPage, NotFound404, ErrorDetail, ErrorAlert).

### Changed

- Contract test comment counts corrected (`display 18→23`, `navigation 6→7`).
- Grid lookup maps no longer include `gap-4` (gap is now a separate `Gap` field with its own lookup). `DefinitionGrid` updated to pass `gridGapClass(GridGapDefault)` explicitly.

## [0.8.0] — 2026-07-06

### Added

- `display.TableHeader` + `TypedHeaders []TableHeader`: sortable table columns with WAI-ARIA `aria-sort` (`ascending`/`descending`/`none`), clickable `<a>` sort links via `Href`, and visual ↑/↓ indicators. `TypedHeaders` takes precedence over `Headers []string` when set; backward compatible (empty `TypedHeaders` keeps existing header rendering). `SortDirection` enum (`None`/`Asc`/`Desc`) added to the display typed-enum set.
- `forms.FormProps.Inline`: horizontal form layout (`flex flex-wrap items-end gap-3`) instead of the default vertical stack (`space-y-6`). One-field toggle — useful for filter bars and compact toolbars. Follows the `RadioGroup.Inline` precedent.
- `navigation.Pagination`: `rel="canonical"` on the first-page link when ellipsis truncates it — tells search engines the first page is the canonical version of a paginated list. `activeSpanOrLink` sub-template gains an optional `rel` parameter.
- 14 new `IsValid()` methods across 5 packages, bringing every closed-set typed enum to full validation coverage (`AvatarStatus`, `DropdownItemKind`, `DropdownPosition`, `TabsVariant`, `OverlayKind`, `ButtonSize`, `ButtonHTMLType`, `StepIndicatorOrientation`, `ToggleSize`, `InputType`, `FormMethod`, `SwapStyle`, `icons.Name`, `SortDirection`). Every `IsValid` is now test-covered.
- `utils.TestVersionMatchesFeatures`: drift-guard test asserting `FEATURES.md` version matches `utils.Version` (mirrors the existing `TestVersionMatchesChangelog`).
- Recipe docs: [`docs/recipes/custom-table-rows.md`](docs/recipes/custom-table-rows.md) (Body slot + sortable `TypedHeaders`), [`docs/recipes/custom-404-page.md`](docs/recipes/custom-404-page.md) (`NotFound404` with custom links/search), [`docs/recipes/recipe-index.md`](docs/recipes/recipe-index.md) (index of all recipes).
- Recipe: [`docs/recipes/container-queries.md`](docs/recipes/container-queries.md) — when and how to use `ContainerResponsive` for parent-width-responsive grids.
- Reference: [`docs/motion-design.md`](docs/motion-design.md) — timing constants, duration guidelines, easing policy, and `motion-reduce` compliance rules.
- Reference: [`docs/javascript-guide.md`](docs/javascript-guide.md) — comprehensive JS patterns reference: decision ladder (native HTML → HTMX → singleton-guard → Alpine → Datastar → React islands), CSP compliance, templ's built-in JS features (`OnceHandle`, `JSFuncCall`, `JSONString`, `JSONScript`), TypeScript workflow, View Transitions API, event delegation, and anti-patterns.
- Audit: [`docs/audits/icon-rtl-mirroring.md`](docs/audits/icon-rtl-mirroring.md) — identifies 5 directional icons needing RTL mirroring, recommends `data-tc-dir-icon` + CSS approach, deferred to v1.0.
- ADR: [`docs/adr/0008-semantic-tokens.md`](docs/adr/0008-semantic-tokens.md) — semantic token layer (`bg-tc-primary`) migration plan, proposed and deferred to v1.0 with opt-in migration path.

### Changed

- **`display.StatCardProps.HxSwap` typed from `string` to `htmx.SwapStyle`** — consumers now pass typed constants (`htmx.SwapInnerHTML`) instead of raw strings, matching the pattern used by `SwapOOB`.
- **`ButtonHTMLType` converted from `map[X]bool` to `map[X]string` + `utils.Lookup`** — matches the convention used by all other enums (InputType, FormMethod, etc.).
- **`feedback.feedbackIconName` + `lookupFeedbackStyle` private helpers removed** — replaced with direct `utils.Lookup` calls, reducing custom boilerplate.
- **6 lookup maps converted from `map[string]string` to typed-key maps**: `cardPaddingLookup`, `avatarSizeLookup`, `avatarDotSizeLookup`, `badgeSizeLookup` (display); `spinnerSizeLookup`, `progressHeightLookup` (feedback). Eliminates all `string(v)` casts at call sites — invalid enum values are now caught at compile time rather than silently missing the map.
- `errorpage.CauseItem.Code` changed from raw `string` to the existing `Code` type (same package), closing a split brain where the `Code` type was defined but unused on this struct.
- `errorpage.FamilyStatusCode` simplified to use `utils.Lookup` instead of manual map+fallback.
- Motion constants (`transitionFast`, `transitionNormal`, `transitionColors`, `transitionTransform`) wired into `CopyButton` — previously only Modal and Drawer used them.
- SKILL.md authoring playbook updated with three new mandatory conventions: RTL logical properties, motion constants, and container queries.

### Fixed

- **Documentation/code split brain corrected**: AGENTS.md, flake.nix, and CHANGELOG v0.7.0 all falsely described a 6-module workspace with `go.work` — the modularization was prototyped on `modularize/strategic-split` but never merged. All three corrected to match the single-module reality.
- **`ModalSize2XL` and `DrawerSize2XL` both had value `"full"`** — identical to the deprecated `ModalSizeFull`/`DrawerFull` aliases. They resolved only by map-key accident (the alias's entry matched). Each now has its own value (`"2xl"`) with a dedicated map entry; the deprecated aliases keep `"full"` for backward compatibility.
- **Combobox WAI-ARIA compliance**: options now carry `aria-selected` (set to `"true"` on the active option alongside `data-selected`); Tab key now closes the listbox and clears `aria-activedescendant`/selection state instead of leaving stale focus. Extracted a shared `tcClearComboSelection()` helper across Escape/Enter/Tab/navigation paths.
- **Combobox `focusout` handler**: listbox now closes and `aria-activedescendant` clears when focus leaves the combobox container (mouse click outside, Tab away). Previously `aria-activedescendant` could remain stale if the outside-click handler didn't fire before blur.
- **Motion-reduce a11y gaps**: 7 `transition-colors` instances missing `motion-reduce` fallbacks fixed across `toast` (dismiss button ×2), `step_indicator`, `empty_state`, `file_input`, `errorpage` (action buttons ×2).
- `FEATURES.md` drift corrected: version `0.6.1` → `0.7.0`, removed phantom `BadgeType "Default"` value, removed already-fixed Tooltip "known issue", added `FeedbackType` to the feedback enum table (`AlertType`/`ToastType` are aliases).

## [0.7.0] — 2026-07-05

### Changed

- **RTL/i18n: all physical CSS properties migrated to logical**. Replaced `ml-`/`mr-` with `ms-`/`me-` (margin-inline-start/end), `pl-`/`pr-` with `ps-`/`pe-` (padding-inline-start/end), `left-0`/`right-0` with `start-0`/`end-0` (inset-inline-start/end), `text-left` with `text-start`, `border-l-`/`border-r-` with `border-s-`/`border-e-` across all `.templ` files. Zero behavioral change in LTR contexts (Tailwind logical utilities resolve identically). Makes the library RTL-ready for Arabic, Hebrew, Persian, and Urdu markets — consumers set `dir="rtl"` and components automatically mirror.
- **Multi-module workspace (prototyped, not shipped)**: a 6-module split was prototyped on `modularize/strategic-split` but not merged. `master` remains single-module. The split may be revisited post-v1.0.
- `go-error-family` remains a direct dependency of `errorpage` (not isolated to a sub-module since the split was not merged).

### Added

- `display.GridProps.ContainerResponsive`: when `true`, wraps the grid in an `@container` div and uses Tailwind container-query variants (`@sm:`, `@lg:`, etc.) instead of viewport breakpoints. The grid adapts to its parent container's width, not the browser viewport — useful for grids in sidebars, cards, or constrained layouts. Defaults to `false` (viewport breakpoints, current behavior). Requires Tailwind CSS v4 (container queries built into core).
- `display.CopyButton`: CSP-safe clipboard copy button with singleton event-delegation script. Copies text via `navigator.clipboard.writeText`, temporarily shows a "Copied!" label, reverts after 2s. Optional clipboard icon, fully accessible (type=button, focus ring, motion-reduce).
- `display.RelativeTime`: relative timestamp ("2 hours ago", "3 days ago") in a `<time datetime>` element. Server renders the initial text (pure Go formatting); `AutoRefresh` (defaults to `true`) injects a singleton script using native `Intl.RelativeTimeFormat` that live-updates every 30s and on `htmx:afterSettle`. Progressive enhancement — HTML carries the `datetime` attribute, JS just keeps the display fresh. Set `AutoRefresh: false` for static contexts (PDF, email).
- `display.CountBadge`: notification count overlay — renders children (e.g. a bell icon) with an absolutely-positioned count badge in the top-right corner. Overflow shows "N+" (default max 99). Zero count hides the badge entirely. Badge is `aria-hidden` (decorative — count is announced by the container's aria-label).
- `display.DefinitionGrid`: responsive grid of term-detail pairs in SimpleCard tiles. Composes through `Grid` + `SimpleCard` internally. Ideal for dashboard metrics and settings pages where many key-value pairs need to be scanned side by side.
- `display.Image`: `<img>` with lazy loading (`loading="lazy"` default), optional `width`/`height` for CLS prevention, and CSP-safe fallback source. The fallback swap uses a singleton error-capture listener (`data-tc-img-fallback` attribute) — no inline `onerror` handler.
- `navigation.LoadMore`: centered "Load more" button for cursor-based pagination. Uses `hx-get` + `hx-swap="outerHTML"` so the server response (next items + updated button) replaces this one in place. Cursor is appended as a query parameter.
- `display.CardProps.Body`: explicit `templ.Component` body slot for struct-based composition. When set, overrides children. Backward compatible — existing children-passing code is unaffected.
- `display.TableProps.Body`: explicit `templ.Component` body slot for custom `<tr>` rendering. When set, overrides `Rows` — ideal for templ loops where each row needs custom cell rendering. Follows the Card.Body pattern. Backward compatible.
- Recipe: [`docs/recipes/horizontal-filter-bar.md`](docs/recipes/horizontal-filter-bar.md) — horizontal HTMX filter bar pattern vs `forms.Form`, with copy-pasteable helper code.
- SKILL.md: "Components by use case" cross-reference table above the per-package catalogue. Consumer tip: track library component adoption in your project's AGENTS.md.
- `display.StatCardProps.HxGet`/`HxTarget`/`HxSwap`: typed HTMX fields on StatCard for HTMX-driven partial updates. When set, the corresponding `hx-*` attributes are rendered on both the `<a>` and `<div>` variants.
- Recipe: [`docs/recipes/cursor-pagination.md`](docs/recipes/cursor-pagination.md) — cursor-based pagination pattern with HTMX infinite scroll using `navigation.LoadMore`.
- ADR: [`docs/adr/0007-self-host-htmx-default.md`](docs/adr/0007-self-host-htmx-default.md) — decision to self-host htmx as default (CDN opt-in) in v1.0.

### Changed

- `display` package: 20 → 25 components (CopyButton, RelativeTime, CountBadge, DefinitionGrid, Image added).
- `navigation` package: 10 → 11 components (LoadMore added).
- README component count: 76 → 82. Display section updated with new component examples.
- Demo app: new "New Components (Session 7)" section showcasing all 6 new components + LoadMore.
- Registered all 7 new props types in `internal/contract/component_props_test.go` (the cross-package BaseProps contract inventory).
- `errorpage.NotFound404`: dedicated, visually striking 404 page with large gradient numeral hero, optional search form, quick-links card grid, and "Go home" / "Go back" actions. Welcoming navigation aid (not an alarming error). Accessible, CSP-safe, dark-mode aware. Composable via `NotFound404Props` + `NotFoundLink` types.

### Changed

- `layout.PageProps`: documented the two auto-injected `<head>` tags in godoc — `CSSPath` ("/app.css" default) and `HTMXVersion` (HTMXVersion2_0_10 default) — and how to suppress each by setting to "". `DefaultPageProps()` godoc now explicitly calls out these as the most common defaults to override when integrating with an existing asset pipeline. Addresses the "silent 404 / silent CDN tag" friction reported by two consumers.
- README layout section: new "Suppressing auto-injected `<head>` tags" subsection with copy-paste example for blanking `HTMXVersion` and `CSSPath`.
- README component catalogue: added `display.Grid` (count 19 → 20), `feedback.SkeletonCardGrid`, `StatCard.Href`, and `SimpleNav.RightItems` examples. Cross-linked the two new recipe docs.
- Registered `display.GridProps` in `internal/contract/component_props_test.go` (the cross-package BaseProps contract inventory).

### Internal

- Code review session 7: fixed stale count comments in component_props_test.go (display 18→23, nav 6→7), stale comment in fromerror.go, missing WayOutHref on 3 error constructors, extracted shared `scriptComponent()` helper (eliminates 4 near-identical functions), added `maxW2XL` named constant, fixed `allIconNames()` sorting, extracted `resolveCDNBase()` helper, removed competing package doc comments from 8 files
- Added `TestPinnedSRIMatchesCDN` network-gated test that fetches live CDN scripts and verifies pinned SRI hashes match the bytes (skips under `-short` and on network errors)
- Added `release.sh` pre-check: fails if `[Unreleased]` section body is empty
- Extracted `statCardInner` sub-template from `StatCard` so the linked (`<a>`) and unlinked (`<div>`) layouts share the icon/value/label body without duplication
- Added `TestSimpleNavForwardsRightItems`, `TestSimpleNavOmitsRightItemsWhenNil`, `TestStatCardRender/Href_*`, `TestGridResponsiveClasses`, `TestGridFallsBackForUnknownCols`, `TestGridRendersChildren`, `TestGridPropagatesBaseProps`, `TestScriptRender`, `TestSkeletonCardGridRender`
- Added golden tests for `Grid`, `StatCard.Href`, and `SkeletonCardGrid`
- Added BDD tests for `Grid` (responsive layout) and `StatCard.Href` (clickable filter)
- Added a11y tests for `Grid` (aria-label/ID propagation) and `StatCard.Href` (focus-visible ring)
- Added `ExampleGrid` godoc example
- Added `TestGridWithStatCards` and `TestStatCardWithHrefComposes` to integration composition suite
- Fixed `GridCols4`/`GridCols5` responsive ladders to include intermediate breakpoints (md) instead of jumping directly from 2 to the final count
- Modernized `ProgressBar` clamp from manual if-branch to `max(0, min(100, v))` (Go 1.21+ builtins)
- Updated `AGENTS.md`, `TODO_LIST.md`, `FEATURES.md` with session 6 conventions and component inventory
- Code review session 8: CopyButton `execCommand('copy')` fallback for non-secure HTTP contexts; `role="status"` + `aria-live="polite"` on label span; typed `OverlayKind` enum replaces untyped `closeKind`/`componentName` string fields on `overlayShellProps`; `formatRelativeTime` boundary tests; golden tests for StatCard HTMX (`div` + `a` variants) and Card.Body slot; 7 composition integration tests (CopyButton+Card, CountBadge overflow, Image+fallback, DefinitionGrid, Card.Body, Grid); benchmark tests for CopyButton, CountBadge, Image, RelativeTime, LoadMore; Image srcset/sizes documentation; replaced `formatInt` with `strconv.Itoa`; typed `Code` enum in errorpage; `IsValid()` methods for ButtonType, ModalSize, DrawerSize, DrawerSide, FeedbackType; SRI returns empty string for unknown HTMX versions
- Code review session 9: `templ.SafeURL()` XSS guard on Card.Href and Badge.Href; Image empty-src guard; `motion-reduce:animate-none` on Icon spinner, SkeletonGroup, ProgressBar indeterminate; HTMX nonce always rendered; stale FEATURES/README counts fixed; 9 more IsValid methods (BadgeType, BadgeSize, CardPadding, GridCols, TrendDirection, AvatarSize, SpinnerSize, ProgressBarSize, TooltipPosition, AvatarShape, SkeletonVariant); LoadMore uses `net/url` for cursor encoding (base64-safe); SimpleCard.Body slot; `layout.Stylesheet` helper; RTL rendering assertion tests; CSP nonce-presence assertion test across all inline-script components; dead transition constants removed; NotFound404Props registered in contract test

## [0.6.1] — 2026-07-04

### Added

- `PageProps.HTMXCDN`: overrides the CDN base URL for htmx scripts. Empty defaults to `https://cdn.jsdelivr.net/npm`. Both the htmx main script and the response-targets extension derive their URLs from this value, so consumers with a different CSP allow-list (e.g. `unpkg.com` or a self-hosted origin) no longer need to fork the library.

### Fixed

- htmx CDN switched from `unpkg.com` to `cdn.jsdelivr.net` — unpkg was not in any consumer's CSP allow-list, causing htmx scripts to be silently blocked by the browser
- `Favicon`: no `<link rel="icon">` tag is rendered when `Favicon` is empty, letting consumers provide their own favicon via `HeadContent` (e.g. a data URI that templ's URL sanitizer would otherwise reject)

### Internal

- Regenerated all `*_templ.go` files with standardized import grouping matching go.mod templ pin (v0.3.1020)
- Added cross-package `ComponentProps` contract test in `internal/contract`
- Added `scripts/release.sh` for automated one-commit releases

## [0.6.0] — 2026-06-29

### Added

- Tooltip touch-device support: click/tap toggles visibility, Escape and click-outside dismiss (idempotent JS body guarded by `window.tcTooltipAttached`, CSP-safe with nonce)
- Tooltip auto-generates an ID via `utils.EnsureID` when none is provided, so `aria-describedby` is always wired up
- Typed `HTMXVersion` enum (`HTMXVersion2_0_10`) replacing the bare string, matching the library's typed-constant convention
- `ThemeColor`/`DarkThemeColor` are now validated as CSS hex colors, falling back to `DefaultThemeColor`/`DefaultDarkThemeColor` for invalid values instead of emitting garbage into the `<meta>` tag
- Size constants (`AvatarSizeSM`/`MD`/`LG`, `BadgeSizeSM`, `SpinnerSM`, …) for programmatic size selection
- `Toggle`: `Required`, `Error`, and `HelpText` fields for form integration
- `ConfirmDelete` and `SwapOOB` now embed `BaseProps`, gaining `Class`/`ID`/`Attrs`/`AriaLabel`
- `ErrorHandlerConfig.Lang` to override the `<html lang>` attribute on error pages

### Changed

- **Breaking:** `forms.FormFieldWrapper` now takes a `FormFieldProps{ID, Label, Required, Error, HelpText}` struct instead of 5 positional parameters (affects `Input`, `Textarea`, `Select`, `FileInput`, `DatePicker`, `Combobox`)
- **Breaking:** `htmx.ConfirmDelete` now takes a `ConfirmDeleteProps{Delete, Target, Confirm}` struct instead of 3 positional strings
- **Breaking:** `htmx.SwapOOB` now takes a `SwapOOBProps{Selector, SwapStyle}` struct instead of positional parameters
- `errorpage` handler split into focused files; `WriteErrorPage` now derives the HTTP status from `props.Family` when `statusCode` is 0 (prevents status/family mismatch)
- `errorpage` renders to a buffer before writing the response, so a mid-stream templ failure can no longer emit a truncated HTML document at the wrong status code
- `Drawer` replaced inline `style="inset-y:0;left:0"` with Tailwind classes (`inset-y-0 left-0` / `right-0`) via `templ.KV` conditionals
- `PageProps.HTMXVersion` field type: `string` → `HTMXVersion`

### Fixed

- Library did not compile for consumers: four generated `*_templ.go` files (DefinitionList, ListNote, SidebarNav, PageHeader) were missing from the Git tag because a redundant `*_templ.go` line in `.gitignore` overrode the `!*_templ.go` unignore
- `Button`: invalid `aria-disabled:pointer-events-none` arbitrary variant (not real Tailwind) replaced with `pointer-events-none opacity-50` plus explicit `aria-disabled="true"` and `tabindex="-1"` for disabled links
- `Spinner`: now renders `role="img"` when `AriaLabel` is set, so the label is reachable (previously stayed `aria-hidden`)
- `Avatar`: status dot now renders in the initials/fallback branches, not just the image branch
- `errorpage.ExtractCauseChain` now handles `errors.Join` siblings (`Unwrap() []error`, Go 1.20+), not only single-error chains

### Internal

- Templ duplication reduced (19 → 17 clone groups at threshold 4) via shared `navLinkAnchor` sub-template, `emptyStateAction` helper, `mutedTextClass` constant, and `paginationPageItem`/`paginationEllipsisItem` sub-templates
- Duplicate default constants removed; `buttonVariantDefault`/`badgeStyleDefault` now derive from their lookup maps
- `internal/golden` test isolation fixed: package tests use `t.TempDir()` instead of a shared `testdata/` that raced under `t.Parallel`
- Accepted clones (`feedback/alert` ↔ `errorpage/erroralert` dismiss button; `Modal` ↔ `Drawer` panel body) documented with rationale comments

## [0.5.0] — 2026-06-28

### Added

- `display.ButtonHTMLType` enum: typed replacement for the raw `string` on `ButtonProps.Type` (button/submit/reset), with `buttonHTMLType()` normalizer that falls back to `"button"` for unknown values
- `forms.formMethod()` normalizer: validates `FormMethod` and falls back to `GET` (HTML spec default) for unknown values
- `utils.Version`: single source of truth for the library version string, with `TestVersionMatchesChangelog` drift-guard test
- GOTH stack ecosystem section in README (cross-links cqrs-htmx, go-cqrs-lite, go-error-family)

### Fixed

- `display.AvatarStatus`: unknown status values no longer render an invisible (colorless) dot — only `online` and `offline` render the status indicator
- `ButtonProps.Type`: previously a raw `string` emitted unvalidated to the DOM (`type="destroy-everything"` would render); now typed and validated
- `forms.Form`: invalid HTTP methods no longer render verbatim to the DOM
- CHANGELOG, FEATURES.md, CONTEXT.md, TODO_LIST.md: all metrics corrected to match actual code (73 components, 101 icons, 51 generated files)
- AGENTS.md: corrected false claims (generated file count 46→51, SanitizeID usage)

### Changed

- `ButtonProps.Type` field type: `string` → `ButtonHTMLType` (backward-compatible — untyped string constants still assign)
- `forms.FormProps.Method` rendering: now validated via `formMethod()` instead of raw `string()` cast
- Demo footer version: hardcoded string → `utils.Version` reference
- All 47 generated `*_templ.go` files: import grouping normalized by clean `templ generate` run

## [0.4.0] — 2026-06-27

### Added

- `display.PageHeader`: page header with Title, Subtitle, Breadcrumb, and Action component slots
- `display.DefinitionList`: two-column `<dl>` with typed `DefinitionItem` entries
- `display.ListNote`: "Showing N of M" truncation notice for truncated lists
- `navigation.SidebarNav`: vertical sidebar navigation with icons and active-route detection
- `display.StatCard.Icon` field: optional `icons.Name` rendered alongside the stat value
- `icons.IconPathData()`: exported function returning raw path data for consumers needing full `<svg>` wrapper control
- `icons.ArrowRightOnRectangle`, `icons.BuildingOffice2`, `icons.Key`: three new named icons
- `flake.nix`: reproducible devShell (go_1_26, gopls, golangci-lint, templ) and apps: `verify`, `test`, `lint`, `build`, `coverage`
- Golden snapshot tests for the `display` and `navigation` packages (`internal/golden.Assert`)
- `docs/adr-001-tailwind-v4-standard.md`, `docs/tailwind-v4-adoption-guide.md`, `docs/icons-only-adoption.md`: adoption and architecture docs

### Changed

- **Tailwind CSS v4+ adopted as the standard** for all LarsArtmann projects — no CSS-variable portability layer (see ADR-001)
- `display.Card` shell shadow: `shadow-sm` → `shadow-xs` (Tailwind v4 rename)
- `errorpage.ErrorPage` shadow: `shadow-sm` → `shadow-xs`
- `forms.Toggle` shadow: `shadow` → `shadow-sm`
- `display.Card` / `SimpleCard` share a `cardShellClass` constant for consistent styling

### Fixed

- README accuracy: corrected component count (69 → 73), icon count (99 → 101), CSS approach description ("Raw Tailwind only" → "Tailwind v4+ CSS-first"), test counts, and rewrote the theming section

## [0.3.0] — 2026-06-20

### Added

- `forms.DatePicker`: native HTML `<input type="date">` wrapper with min/max constraints, follows FormFieldWrapper pattern
- `forms.Combobox`: accessible autocomplete with client-side filtering, `role="combobox"` + `role="listbox"`, global singleton JS handler, auto-generated IDs via `utils.EnsureID`
- `utils.Lookup[K, V]`: generic map lookup with fallback — replaces the narrower `MapEnum`. Handles all map types including struct values and typed keys. Adopted at all 15 call sites, eliminating ~42 lines of duplicated boilerplate
- `utils.EnsureID(prefix, id)`: auto-generates unique IDs via `crypto/rand` (format `tc-<prefix>-<16hex>`) when a consumer omits `props.ID`
- `utils.RenderAll(t, components...)`: test helper for rendering multiple components into a concatenated string — supports integration tests
- `integration/composition_test.go`: cross-package composition tests verifying components render correctly together (full page, form with multiple inputs, modal with form content, CSP nonce consistency)
- Coverage boosters across all 10 packages — display, errorpage, feedback, forms, navigation, layout each gained dedicated coverage test files
- `display.overlayScriptComponent()`: shared overlay JS generator for Modal and Drawer — produces open/close functions, focus trap, focus save/restore, and CSP-safe `[data-tc-close]` click delegation from a single source of truth
- `navigation.SimpleNavProps` struct with `DefaultSimpleNavProps()` — replaces positional parameters, adds BaseProps embedding
- `forms.FormFieldWrapper()`: shared sub-template for Label + FieldError + helpText rendering, adopted by Input, Select, and Textarea
- `feedback.feedbackStyleMap` / `feedback.feedbackIconMap`: single source of truth for Alert and Toast styles — ensures identical appearance for the same severity
- `display.buttonVariantLookup` / `display.buttonSizeLookup`: map-based class lookups replacing switch statements
- `forms.toggleSizeMap` / `forms.toggleSizeSet`: map-based toggle size lookup replacing switch
- `errorpage.handler.go`: CSP-safe `data-tc-go-back` attribute replacing inline `onclick="history.back()"`

### Changed

- **BREAKING**: `utils.MapEnum[T ~string](m map[string]T, fallback T, key string) T` removed → replaced by `utils.Lookup[K, V](m map[K]V, key K, fallback V) V` — the old signature was too narrow, only handling string-keyed maps with string-like values. The new generic handles struct values and typed keys.
- **BREAKING**: `SimpleNav(brandText, brandHref, links, currentPath)` → `SimpleNav(SimpleNavProps)` — positional params replaced with props struct + BaseProps
- **BREAKING**: `forms.FormProps.Content templ.Component` removed — Form now uses `{ children... }` pattern matching Card, Modal, Drawer, InputGroup
- **BREAKING**: `navigation.PaginationProps.CurrentPage`, `TotalPages`, `MaxVisible` changed from `int` to `uint` — negative page numbers made unrepresentable at the type level
- **BREAKING**: `errorpage.BreadcrumbList` struct fields `Type` and `Context` swapped to match their JSON tags (`@type` and `@context`)
- Modal and Drawer: inline `onclick` handlers replaced with `data-tc-close` attribute + per-instance event delegation — CSP compliant (no `script-src-attr` needed)
- Alert and Toast: duplicate `alertStyleMap`/`alertIconMap` and `toastStyleMap`/`toastIconMap` consolidated into shared `feedbackStyleMap`/`feedbackIconMap`
- Input, Select, Textarea: now delegate field chrome rendering to `FormFieldWrapper` instead of manual Label+FieldError+helpText
- `errorpage.htmlEscape()` replaced with `html.EscapeString()` from stdlib
- `display.button_go.go`: two `switch` statements converted to map lookups with fallback constants
- `forms.toggle.templ`: `switch` converted to `toggleSizeMap` with `toggleSizeSet` struct
- `layout.ThemeToggle`: added `utils.Ternary` default for aria-label ("Toggle theme")
- `errorpage/styles.go`: `FamilyInfrastructure` changed from `slate-*` to `gray-*` for design system consistency
- `display/dropdown.templ`: stray leading whitespace on type declaration removed; `dark:hover:bg-slate-700` → `gray-700`
- `forms.InputType`: unknown types now fall back to "text" instead of panicking — matches HTML spec
- `icons.Name`: unknown icon names now fall back to the Question icon instead of crashing render
- `forms.RadioGroup`: `<fieldset>` now propagates `AriaLabel` from BaseProps (was silently dropped)
- `display.Avatar`: image branch wrapper `<div>` now propagates all BaseProps (ID, Class, AriaLabel, Attrs) — was only on inner `<img>`
- Modal, Drawer, Dropdown: empty `props.ID` now auto-generates a unique ID via `utils.EnsureID` instead of panicking
- `display.Accordion`: items with empty ID now auto-generate IDs instead of panicking
- `htmx.SwapOOB`: invalid swap styles now fall back to `outerHTML` instead of panicking
- `display.BadgeInfo`: changed from `indigo-*` to `blue-*` to match the library's primary color and `FeedbackInfo`

### Fixed

- Modal/Drawer CSP violations: 4 inline `onclick` handlers generated `script-src-attr 'unsafe-inline'` requirement — replaced with nonce'd event delegation
- Modal/Drawer HTMX regression: `data-tc-close` click listeners used per-element binding that broke on HTMX DOM swaps — replaced with event delegation on overlay container
- Toast icon split brain: server-rendered toasts showed XCircle for errors, client-side tcShowToast showed ExclamationTriangle — unified to use `feedbackIconMap` as single source of truth
- `navigation.BreadcrumbList` struct field naming lie: `Type`/`Context` were swapped relative to their JSON tags
- `forms/validation.templ`: pluralization `"error(s)"` → proper `"%d error%s"` with Ternary
- Removed dead code: `utils.AssertContainsClass` — identical to `AssertContains`, zero callers

## [0.2.0] — 2026-06-08

### Added

- `display.Drawer`: accessible side panel component with left/right slide, focus trap, Escape key, backdrop click, configurable size (SM/MD/LG/XL/Full)
- `forms.ValidationSummary`: accessible error summary with icon, error count, linked field errors, `role="alert"`
- 25 new Heroicons (98 path icons + 1 Spinner = 99 total): ArchiveBox, ArrowPath, Bars3, Beaker, Bolt, BugAnt, Calculator, Cube, FaceSmile, Fire, FolderOpen, Gift, HandThumbUp, Hashtag, PuzzlePiece, RocketLaunch, Server, Signal, Squares2x2, AcademicCap, ArrowDownOnSquare, ArrowUpOnSquare, BellSlash, Camera, NoSymbol
- `internal/golden`: golden file comparison package with CSS class normalization for deterministic snapshot testing
- Coverage tests for display (Drawer) and forms (ValidationSummary) packages
- CI coverage threshold raised from 60% to 70%
- `feedback/progress.templ` split into `progressbar.templ` + `step_indicator.templ` for code organization
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

- **BREAKING**: `Spinner(size SpinnerSize, colorClass string)` → `Spinner(SpinnerProps)` with BaseProps support (ID, Class, AriaLabel, Attrs), Size, Color fields
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
