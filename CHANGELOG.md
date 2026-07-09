# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

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

### Fixed — Round 2 (htmx, errorpage, layout, forms, navigation)

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

### Fixed — Round 1 (forms, feedback, display, navigation)

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
