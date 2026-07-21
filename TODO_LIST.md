# TODO List — templ-components

**Updated:** 2026-07-18 | **Version:** 0.18.0

> Built from 42 `docs/**/2026-07-0*` files + code verification. Each item is verified against
> the actual codebase. Statuses: ✅ done, ⬜ deferred, ⚫ blocked (needs external resources).

---

## P0 — Current session: doc count drift and website newsletter

| #   | Task                                                                                       | Status  | Evidence                                                                                               | Source           |
| --- | ------------------------------------------------------------------------------------------ | ------- | ------------------------------------------------------------------------------------------------------ | ---------------- |
| 53  | Fix `docs/DOMAIN_LANGUAGE.md` broken glossary table (malformed separator, IconPath pipe)   | ✅ DONE | `docs/DOMAIN_LANGUAGE.md` now renders correctly; pipe removed from cell content                        | self-review      |
| 54  | Align component count across `FEATURES.md`, `SKILL.md`, `AGENTS.md`, website from 97 to 94 | ✅ DONE | All files now state 94 exported templ functions                                                        | self-review      |
| 55  | Align generated `*_templ.go` file count from 75 to 82 in `FEATURES.md` and `AGENTS.md`     | ✅ DONE | Both files now state 82 generated files                                                                | self-review      |
| 56  | Align website enum count from 34 to 37 in `website/src/data/sections.ts`                   | ✅ DONE | `sections.ts` now states 37 typed string enums                                                         | self-review      |
| 57  | Add newsletter signup component to website footer                                          | ✅ DONE | `website/src/components/Newsletter.astro` created and imported in `Footer.astro`                       | previous session |
| 58  | Bump nixpkgs in `flake.lock`                                                               | ✅ DONE | `flake.lock` updated to latest nixos-unstable                                                          | previous session |
| 59  | Add `TestDocsCountDrift` to guard component/generated-file/enum counts in docs             | ✅ DONE | `utils/docs_count_test.go` asserts counts in FEATURES.md, AGENTS.md, SKILL.md, and website sections.ts | self-review      |

---

## P0 — Real bugs & correctness gaps

| #   | Task                                                                                                                                                      | Status  | Evidence                                                                                                          | Source                  |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ----------------------------------------------------------------------------------------------------------------- | ----------------------- |
| 1   | `InlineLoadingOverlay` missing sr-only loading text (parity with `LoadingIndicator` which has `<span class="sr-only">Loading…</span>`)                    | ✅ DONE | `htmx/loading.templ:36` — added `<span class="sr-only">Loading…</span>`                                           | bug-hunt-status:143     |
| 2   | `SanitizeID` mismatch in `ValidationSummary` — error links point to `#user-email` (sanitized) but field IDs are whatever consumer set (e.g. `user_email`) | ✅ DONE | `forms/validation.templ:61` — link now uses raw `err.Field`; Field doc clarifies it should be the HTML element ID | bug-hunt-status:129     |
| 3   | `FromError` returns `FamilyInfrastructure` (→503) for unknown errors instead of `FamilyCorruption` (→500)                                                 | ✅ DONE | `errorpage/fromerror.go:71` — now returns `FamilyCorruption` (→500); tests updated                                | bug-hunt-status:127     |
| 4   | `Footer` doesn't accept `BaseProps` — can't set Class/ID/Attrs (API inconsistency; every other component embeds BaseProps)                                | ✅ DONE | `navigation/nav.templ:119` — now takes `FooterProps` with `BaseProps`; all callers/tests/README updated           | bug-hunt-status:180,188 |
| 5   | ErrorPage / NotFound404 missing `<main>` landmark — only `<div role="region">`, failing WCAG 2.4.1 (Bypass Blocks)                                        | ✅ DONE | `errorpage/errorpage.templ:7`, `errorpage/notfound404.templ:28` — changed to `<main>`; golden files updated       | bug-hunt-status:135     |
| 6   | `FormProps` CSRF token name hardcoded (`name="csrf_token"`) — frameworks use different names (Django, Rails, Spring)                                      | ✅ DONE | `forms/form.templ:71` — added `CSRFTokenName` field (defaults to `"csrf_token"`)                                  | bug-hunt-status:128     |
| 7   | `grid-rows-[0fr]` CSS output never verified against compiled Tailwind v4 — accordion collapse depends on it                                               | ✅ DONE | Verified: Tailwind v4.3.1 generates `grid-template-rows: 0fr` correctly                                           | bug-hunt-status:155     |

---

## P1 — Testing gaps

| #   | Task                                                                                                                                                                                                                                                  | Status     | Evidence                                                                                                                                                                                                                                                                                        | Source                 |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------- |
| 8   | Add regression tests for 18 untested Round-2 bug fixes (htmx LoadingButton/InlineLoadingOverlay/retry/announcer/catch-all/ConfirmDelete/SwapOOB, errorpage a11y, layout ThemeToggle/localStorage/FOUC/SRI, navigation LoadMore/breadcrumb/SidebarNav) | ✅ DONE    | 27 new regression tests across 4 packages (htmx, errorpage, layout, navigation)                                                                                                                                                                                                                 | bug-hunt-status:94-113 |
| 9   | `InlineLoadingOverlay` `role="status"` assertion — BDD and snapshot tests don't check for it                                                                                                                                                          | ✅ DONE    | `htmx/regression_test.go` — TestInlineLoadingOverlayAccessibility                                                                                                                                                                                                                               | code-verification      |
| 10  | Dark golden test variants — render components inside `<div class="dark">` wrapper and generate dark-mode golden files                                                                                                                                 | ✅ DONE    | `display/dark_golden_test.go` — badge_dark, card_dark, button_dark golden files                                                                                                                                                                                                                 | dark-mode-plan:213     |
| 11  | Toast JS-created toast golden test — `tcShowToast()` dynamically constructs HTML; only templ path is tested                                                                                                                                           | ✅ DONE    | `feedback/toast_regression_test.go` — 6 tests covering container + JS path                                                                                                                                                                                                                      | dark-mode-plan:212     |
| 12  | Coverage → 80%+ on packages still below: errorpage (~73%), feedback (~74%), forms (~73%), navigation (~73%)                                                                                                                                           | ✅ DONE    | **All non-templ Go code is at 90-100%.** The 72-74% ceiling is entirely from templ-generated `if err != nil` error branches that are structurally unreachable. Verified: feedback and navigation have **zero** non-templ files below 100%. 80% target was unrealistic for templ-heavy packages. | coverage-sprint        |
| 13  | Visual regression testing (Playwright screenshot diff light/dark)                                                                                                                                                                                     | ⚫ BLOCKED | Requires npm/playwright infra — no Node.js dependency allowed in this repo                                                                                                                                                                                                                      | bug-hunt-status:157    |

---

## P2 — Pre-commit / CI hardening

| #   | Task                                                                      | Status  | Evidence                                                                                | Source          |
| --- | ------------------------------------------------------------------------- | ------- | --------------------------------------------------------------------------------------- | --------------- |
| 14  | Add `encoding/json/v2` grep guard to pre-commit hook                      | ✅ DONE | `scripts/pre-commit.sh` — grep guard added; rejects `encoding/json/v2` imports          | css-cleanup:92  |
| 15  | Pre-commit lint uses hardcoded package paths instead of `./...`           | ✅ DONE | `scripts/pre-commit.sh:22` — now uses `./...`; `examples/` excluded via `.golangci.yml` | css-cleanup:100 |
| 16  | Document `encoding/json/v2` prohibition in AGENTS.md (blocked on Go 1.27) | ✅ DONE | New section in AGENTS.md documents the prohibition + pre-commit guard                   | css-cleanup:102 |

---

## P2 — Documentation accuracy

| #   | Task                                                                                               | Status  | Evidence                                                                  | Source            |
| --- | -------------------------------------------------------------------------------------------------- | ------- | ------------------------------------------------------------------------- | ----------------- |
| 17  | Fix AGENTS.md lint path typo: `./svg/...` → `./internal/svg/...`                                   | ✅ DONE | AGENTS.md lint command now uses `./...` (typo eliminated)                 | v0.10-release:108 |
| 18  | Add note to CHANGELOG `[0.9.1]` that it was never tagged (changes included in v0.10.0)             | ✅ DONE | `CHANGELOG.md:131` — note added warning consumers                         | v0.10-release:110 |
| 19  | ROADMAP.md doesn't mention dark mode compliance milestone                                          | ✅ DONE | `ROADMAP.md` — dark mode row added to v0.x table                          | v0.10-release:64  |
| 20  | Create `docs/migration/v0.9-to-v0.10.md` migration guide                                           | ✅ DONE | `docs/migration/v0.9-to-v0.10.md` — created with breaking changes + fixes | v0.10-release:61  |
| 21  | Update FEATURES.md with `templates/app.css` + BuildFlow `tailwind-build` provider entry            | ✅ DONE | Already present at `FEATURES.md:388` — verified                           | css-cleanup:58    |
| 22  | AGENTS.md "Post-v0.9.0 Conventions" section header is stale (shipped in v0.10.0) — rename or merge | ✅ DONE | Renamed to "Conventions" in AGENTS.md                                     | v0.10-release:66  |
| 23  | AGENTS.md claims "61 generated files" but actual count is 62                                       | ✅ DONE | AGENTS.md updated to 62 — matches actual count                            | code-verification |

---

## P2 — Code quality & consistency

| #   | Task                                                                                                                                                 | Status  | Evidence                                                                            | Source              |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ----------------------------------------------------------------------------------- | ------------------- |
| 24  | Wire shared motion constants into remaining ~19 transition-bearing components (only 3/22 use `transitionFast`/`transitionNormal`/`transitionColors`) | ✅ DONE | Constants promoted to `utils/motion.go`; wired into 13 components across 6 packages | session11:77        |
| 25  | `FamilyFromErrorFamily` stuttery name — consider `FromErrorFamily`                                                                                   | ✅ DONE | `errorpage/fromerror.go:12` — renamed; old name kept as deprecated alias            | naming-review       |
| 26  | Consolidate CHANGELOG "Round 1"/"Round 2" headings into single `### Fixed` section                                                                   | ✅ DONE | `CHANGELOG.md` — internal "Round 1/2" labels removed                                | bug-hunt-status:161 |

---

## P3 — Polish & community

| #   | Task                                                             | Status     | Evidence                                                                                                            | Source            |
| --- | ---------------------------------------------------------------- | ---------- | ------------------------------------------------------------------------------------------------------------------- | ----------------- |
| 27  | Demo site / showcase (live rendered components)                  | ✅ DONE    | `examples/demo` serves a live component showcase; Docker + Cloud Run deployment in place; README links to demo      | multiple-feedback |
| 28  | `awesome-templ` PR submission (updated component count)          | ⚫ BLOCKED | External repo submission — needs maintainer approval                                                                | session10:161     |
| 29  | `templ.guide` listing submission                                 | ⚫ BLOCKED | External repo submission — needs maintainer approval                                                                | session10:162     |
| 30  | Configure SSH tag signing (`gpg.ssh.allowedSignersFile`)         | ⚫ BLOCKED | Requires user's local git config + SSH key setup                                                                    | v0.8-release:82   |
| 31  | Blocks/composition examples (dashboard, login, settings layouts) | ✅ DONE    | `recipes/` package shipped in v0.21.0 — Dashboard, SettingsLayout, LoginCard; demo routes `/recipes/*` added v1.1.1 | v0.21.0           |
| 32  | Standalone `/forms` quickstart demo route                        | ✅ DONE    | `examples/demo` — `/forms` route shipped in v0.17.0 with all form components + HTMX filter bar                      | v0.17.0           |

---

## v1.0 — Deferred breaking changes

| #   | Task                                                                                  | Status      | Evidence                                             | Source       |
| --- | ------------------------------------------------------------------------------------- | ----------- | ---------------------------------------------------- | ------------ |
| 33  | `Validate() error` methods on props structs (83 components, design decision needed)   | ⬜ DEFERRED | No `Validate` methods exist — design decision needed | multiple     |
| 34  | Move test helpers to `internal/testutil/` (70+ test files depend on exported helpers) | ⬜ DEFERRED | No `internal/testutil/` directory — large migration  | multiple     |
| 35  | Self-host htmx as default, CDN opt-in (ADR 0007 written)                              | ⬜ DEFERRED | `layout/sri.go:69` — CDN is still default            | ADR-0007     |
| 36  | Semantic token layer `bg-tc-primary` (ADR 0008 written, 256 color refs)               | ⬜ DEFERRED | All components use hardcoded `bg-blue-600` etc.      | ADR-0008     |
| 37  | Icon RTL mirroring (`data-tc-dir-icon` + CSS `scaleX(-1)`) — 5 directional icons      | ✅ DONE     | `icons.IconRTL()` + CSS rule in `templates/app.css`  | session11:90 |
| 38  | Remove deprecated aliases (`AlertType`, `ToastType`, `FamilyFromErrorFamily`)         | ⬜ DEFERRED | Kept as aliases for backward compat until v1.0       | session8:120 |

---

## v2.0 — Architectural changes

| #   | Task                                                            | Status      | Evidence                                                                                                                     | Source         |
| --- | --------------------------------------------------------------- | ----------- | ---------------------------------------------------------------------------------------------------------------------------- | -------------- |
| 39  | Compound component pattern (Trigger/Content/Close) for overlays | ⬜ DEFERRED | Current Modal/Drawer are monolithic                                                                                          | research:70    |
| 40  | Native `<dialog>` element for Modal/Drawer                      | ✅ DONE     | `display/shared.templ` — `<dialog>` + `showModal()`/`close()`, ~200 lines of JS eliminated, `@starting-style` CSS            | session:dialog |
| 41  | Headless/unstyled component variants (Radix UI model)           | ❌ WONTFIX  | ADR-0021 evaluated 3 options; existing `Class` override accepted. Closed v1.1.0                                              | session8:126   |
| 42  | CLI tool (`templ-components add <component>`, shadcn-style)     | ✅ DONE     | `cmd/tc/` shipped in v1.1.0 — `tc init`, `tc ls`, `tc add <component>` with embedded `.templ` sources and dependency warning | v1.1.0         |

---

## New components — All shipped (v0.17.0)

| #   | Component                                                      | Status  | Evidence                                                                                                        |
| --- | -------------------------------------------------------------- | ------- | --------------------------------------------------------------------------------------------------------------- |
| 43  | `Popover`                                                      | ✅ DONE | `display/popover.templ` — button-triggered floating panel, 4 positions, click-outside/Escape dismiss            |
| 44  | `DataTable` (high-level sortable/filtering/pagination wrapper) | ✅ DONE | `display/table_data.templ` — integrated sort management, pagination slot, empty-state, 14 tests + golden + a11y |
| 45  | `FilterDropdown`                                               | ✅ DONE | `forms/filter_dropdown.templ` — HTMX auto-submit select, 12 tests + golden                                      |
| 46  | `Slider` (ARIA slider pattern)                                 | ✅ DONE | `forms/slider.templ` — labeled range input, ShowValue, help text, 8 tests + golden                              |
| 47  | `Rating` (star rating, keyboard support)                       | ✅ DONE | `forms/rating.templ` — star rating radio inputs, RatingSize enum, ReadOnly mode, 11 tests + golden              |
| 48  | `TagsInput`                                                    | ✅ DONE | `forms/tags_input.templ` — singleton JS add/remove, MaxTags/AllowDuplicate, 10 tests + golden                   |
| 49  | `ContextMenu` (right-click menu)                               | ✅ DONE | `display/context_menu.templ` — CSP-safe JS, role=menu, Escape/click-outside dismiss, 6 tests                    |
| 50  | `Carousel`                                                     | ✅ DONE | `display/carousel.templ` — prev/next arrows, dot indicators, singleton JS, 7 tests                              |
| 51  | `HoverCard`                                                    | ✅ DONE | `display/hover_card.templ` — CSS-only hover, 4 positions, focus-within support, 10 tests + golden               |
| 52  | `Calendar` (full calendar grid)                                | ✅ DONE | `forms/calendar.templ` — month-view, server-side nav, day links, MinDate/MaxDate, 10 tests + golden             |

---

## Modern web standards — Shipped in v0.18.0 (2026-07-18)

| Feature                                             | Component                 | Status  |
| --------------------------------------------------- | ------------------------- | ------- |
| Native `<dialog>` for Modal/Drawer                  | Modal, Drawer             | ✅ DONE |
| Stylable `<select>` API (`appearance: base-select`) | `SelectProps.Stylable`    | ✅ DONE |
| Auto-growing Textarea (`field-sizing: content`)     | `TextareaProps.AutoGrow`  | ✅ DONE |
| Unified `EnterKeyHintType` enum                     | Input, Textarea           | ✅ DONE |
| Form `hx-validate`                                  | `FormProps.Validate`      | ✅ DONE |
| Semantic `<search>` landmark wrapping               | Input (InputSearch)       | ✅ DONE |
| Image `SrcSet`/`Sizes` responsive delivery          | `ImageProps.SrcSet/Sizes` | ✅ DONE |
| Table `content-visibility: auto` (LazyRows)         | `TableProps.LazyRows`     | ✅ DONE |
| Global `accent-color` CSS                           | `templates/app.css`       | ✅ DONE |

---

## P1 — Next high-impact improvements (from Pareto plan)

Items #60, #63, #64 were removed as duplicates of #38, #35, #39 (see v1.0/v2.0 sections above).

| #   | Task                                                                                                                                                                                                           | Status  | Evidence                                                                         | Source |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | -------------------------------------------------------------------------------- | ------ |
| 61  | Add docs-health CI check                                                                                                                                                                                       | ✅ DONE | `.github/workflows/ci.yaml` runs `TestDocsCountDrift` since v1.0.0               | v1.0.0 |
| 62  | Add `Validate()` to `errorpage.ErrorPageProps` only (catches invalid `StatusCode`/`Family` combos that produce wrong HTTP responses). Other props use graceful `utils.Lookup` fallback — no `Validate` needed. | ✅ DONE | `ErrorPageProps.Validate()` shipped in v1.0.0 — 3 sentinel errors + 8 sub-tests. | v1.0.0 |

---

## P0 — release.sh repair follow-ups (from 2026-07-18 self-review)

Commit `f1a2592` fixed the 3 release.sh defects but introduced 3 new ones (see `docs/status/2026-07-18_20-27_release-sh-repair-brutal-self-review.md`). These track the recovery.

| #   | Task                                                                                     | Status      | Evidence                                                          | Source      |
| --- | ---------------------------------------------------------------------------------------- | ----------- | ----------------------------------------------------------------- | ----------- |
| 65  | Update AGENTS.md "Release Script" section — still describes the old hostile stdin prompt | ⬜ DEFERRED | AGENTS.md:339-360 references "Prompts for release notes on stdin" | self-review |
| 66  | Update AGENTS.md "Release Convention: One-Commit Release" section — references old flow  | ⬜ DEFERRED | Same section, stale step descriptions                             | self-review |
| 67  | Switch treefmt `gofmt` → `gofumpt` in flake.nix — latent conflict with golangci-lint     | ⬜ DEFERRED | flake.nix:142 `gofmt.enable` vs .golangci.yml:134 `gofumpt`       | self-review |

---

## Done — Verified complete (not actionable)

These were frequently listed as open in older reports but are confirmed DONE:

- ✅ Dark mode compliance (30+ `dark:` variants fixed, `TestDarkModeCompliance` + `TestDarkModeSemanticColors` failing tests block CI)
- ✅ All 30+ typed enums have `IsValid()` + tests (32 production methods verified)
- ✅ `Footer` — Nav links wrap gracefully (v0.12.0) + accepts `FooterProps` with `BaseProps` (v0.12.1)
- ✅ `NotFound404` component with search, links, `LinksTitle`, `WriteNotFound404`
- ✅ Sortable `TableHeader` + `TypedHeaders` with `aria-sort`
- ✅ `Form.Inline` horizontal layout + `CSRFTokenName` field
- ✅ `Grid.ContainerResponsive` + `GridProps.Gap` typed enum
- ✅ `layout.Script` + `layout.Stylesheet` CSP-safe helpers
- ✅ `feedback.SkeletonCardGrid`, `display.CopyButton`, `display.RelativeTime`, `display.CountBadge`, `display.Image`
- ✅ `navigation.LoadMore` with `InfiniteScroll` + cursor pagination recipe
- ✅ `display.Grid` responsive grid
- ✅ `SimpleNav.RightItems` slot
- ✅ `StatCard.Href` + typed HTMX fields (`HxGet`/`HxTarget`/`HxSwap: htmx.SwapStyle`)
- ✅ `Card.Body` / `SimpleCard.Body` / `Table.Body` slots
- ✅ RTL logical properties migration (0 physical properties remain)
- ✅ Motion-reduce compliance (0 gaps, grep test enforced; constants centralized in `utils/motion.go`)
- ✅ `OverlayKind` typed enum
- ✅ `icons.Close` alias for `icons.X`
- ✅ `icons.IconRTL()` for directional icons + CSS `[dir="rtl"] [data-tc-dir-icon] { transform: scaleX(-1) }`
- ✅ `color-scheme: light/dark` CSS
- ✅ `prefers-color-scheme` + `prefers-reduced-transparency` support
- ✅ Pre-commit has `go build ./...` + `encoding/json/v2` grep guard
- ✅ `errorpage.FromErrorFamily` (renamed from `FamilyFromErrorFamily`; old name is deprecated alias)
- ✅ `errorpage/handler.go` uses `encoding/json` (v1) — no json/v2 imports
- ✅ Dark golden test infrastructure (badge/card/button)
- ✅ Toast JS path regression tests
- ✅ Demo has ThemeToggle, TableHeader sortable, Form.Inline showcases
- ✅ Benchmark suites in 7 packages
- ✅ Fuzz tests for InputType, FormMethod, ButtonHTMLType
- ✅ ADRs 0007-0015 (self-host htmx, semantic tokens, accepted clones, sub-template extraction, dark mode convention, WCAG contrast, flush prop, json/v2 guard, dialog migration, stylable select)
- ✅ `templates/app.css` starter + BuildFlow `tailwind-build` provider (v0.11.0) + RTL icon mirroring CSS
- ✅ `tc-css` CLI deleted (over-engineered, YAGNI)
