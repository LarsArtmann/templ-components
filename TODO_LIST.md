# TODO List — templ-components

**Updated:** 2026-07-12 | **Version:** 0.16.0

> Built from 42 `docs/**/2026-07-0*` files + code verification. Each item is verified against
> the actual codebase. Items marked ✅ are confirmed done; ⬜ are open.

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

| #   | Task                                                                                                                                                                                                                                                  | Status     | Evidence                                                                        | Source                 |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- | ------------------------------------------------------------------------------- | ---------------------- |
| 8   | Add regression tests for 18 untested Round-2 bug fixes (htmx LoadingButton/InlineLoadingOverlay/retry/announcer/catch-all/ConfirmDelete/SwapOOB, errorpage a11y, layout ThemeToggle/localStorage/FOUC/SRI, navigation LoadMore/breadcrumb/SidebarNav) | ✅ DONE    | 27 new regression tests across 4 packages (htmx, errorpage, layout, navigation) | bug-hunt-status:94-113 |
| 9   | `InlineLoadingOverlay` `role="status"` assertion — BDD and snapshot tests don't check for it                                                                                                                                                          | ✅ DONE    | `htmx/regression_test.go` — TestInlineLoadingOverlayAccessibility               | code-verification      |
| 10  | Dark golden test variants — render components inside `<div class="dark">` wrapper and generate dark-mode golden files                                                                                                                                 | ✅ DONE    | `display/dark_golden_test.go` — badge_dark, card_dark, button_dark golden files | dark-mode-plan:213     |
| 11  | Toast JS-created toast golden test — `tcShowToast()` dynamically constructs HTML; only templ path is tested                                                                                                                                           | ✅ DONE    | `feedback/toast_regression_test.go` — 6 tests covering container + JS path      | dark-mode-plan:212     |
| 12  | Coverage → 80%+ on packages still below: errorpage (~73%), feedback (~73%), forms (~73%), navigation (~73%)                                                                                                                                           | ⬜ OPEN    | `go tool cover -func` — gaps mostly in generated templ error branches           | coverage-sprint        |
| 13  | Visual regression testing (Playwright screenshot diff light/dark)                                                                                                                                                                                     | ⚫ BLOCKED | Requires npm/playwright infra — no Node.js dependency allowed in this repo      | bug-hunt-status:157    |

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

| #   | Task                                                             | Status      | Evidence                                                   | Source            |
| --- | ---------------------------------------------------------------- | ----------- | ---------------------------------------------------------- | ----------------- |
| 27  | Demo site / showcase (live rendered components)                  | ⚫ BLOCKED  | Requires deployment infra ( hosting, domain)               | multiple-feedback |
| 28  | `awesome-templ` PR submission (updated component count)          | ⚫ BLOCKED  | External repo submission — needs maintainer approval       | session10:161     |
| 29  | `templ.guide` listing submission                                 | ⚫ BLOCKED  | External repo submission — needs maintainer approval       | session10:162     |
| 30  | Configure SSH tag signing (`gpg.ssh.allowedSignersFile`)         | ⚫ BLOCKED  | Requires user's local git config + SSH key setup           | v0.8-release:82   |
| 31  | Blocks/composition examples (dashboard, login, settings layouts) | ⬜ DEFERRED | Deferred to v1.0 — needs design review for composition API | research:74       |
| 32  | Standalone `/forms` quickstart demo route                        | ⬜ DEFERRED | Deferred to v1.0 — demo restructuring needed               | session7:79       |

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

| #   | Task                                                            | Status      | Evidence                                          | Source       |
| --- | --------------------------------------------------------------- | ----------- | ------------------------------------------------- | ------------ |
| 39  | Compound component pattern (Trigger/Content/Close) for overlays | ⬜ DEFERRED | Current Modal/Drawer are monolithic               | research:70  |
| 40  | Native `<dialog>` element for Modal/Drawer                      | ⬜ DEFERRED | JS-based focus trap; browser `<dialog>` is better | research:72  |
| 41  | Headless/unstyled component variants (Radix UI model)           | ⬜ DEFERRED | All components ship with Tailwind classes         | session8:126 |
| 42  | CLI tool (`templ-components add <component>`, shadcn-style)     | ⬜ DEFERRED | No CLI exists                                     | session8:127 |

---

## New components — Not started

| #   | Component                                                      | Priority | Source                                                                                               |
| --- | -------------------------------------------------------------- | -------- | ---------------------------------------------------------------------------------------------------- |
| 43  | `Popover`                                                      | ✅ DONE  | `display/popover.templ` — button-triggered floating panel, 4 positions, click-outside/Escape dismiss |
| 44  | `DataTable` (high-level sortable/filtering/pagination wrapper) | High     | DiscordSync + Overview feedback                                                                      |
| 45  | `FilterDropdown`                                               | Medium   | Consumer-requested for filter bars                                                                   |
| 46  | `Slider` (ARIA slider pattern)                                 | Medium   | Research §5                                                                                          |
| 47  | `Rating` (star rating, keyboard support)                       | Low      | Research §5                                                                                          |
| 48  | `TagsInput`                                                    | Low      | Research §5                                                                                          |
| 49  | `ContextMenu` (right-click menu)                               | Low      | Research §5                                                                                          |
| 50  | `Carousel`                                                     | Low      | Research §5                                                                                          |
| 51  | `HoverCard`                                                    | Medium   | Research §5                                                                                          |
| 52  | `Calendar` (full calendar grid)                                | Medium   | Research §5                                                                                          |

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
- ✅ ADR 0007 (self-host htmx), ADR 0008 (semantic tokens), ADR 0009 (accepted clones), ADR 0010 (sub-template extraction), ADR 0011 (dark mode convention)
- ✅ `templates/app.css` starter + BuildFlow `tailwind-build` provider (v0.11.0) + RTL icon mirroring CSS
- ✅ `tc-css` CLI deleted (over-engineered, YAGNI)
