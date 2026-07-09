# TODO List — templ-components

**Updated:** 2026-07-09 | **Version:** 0.12.0

> Built from 42 `docs/**/2026-07-0*` files (5 feedback, 14 status, 8 planning, 4 HTML
> reviews, 11 older status/planning) + code verification. Each item is verified against
> the actual codebase. Items marked ✅ are confirmed done; ⬜ are open.

---

## P0 — Real bugs & correctness gaps

| #   | Task                                                                                                                                                      | Status  | Evidence                                                                                                              | Source                  |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | --------------------------------------------------------------------------------------------------------------------- | ----------------------- |
| 1   | `InlineLoadingOverlay` missing sr-only loading text (parity with `LoadingIndicator` which has `<span class="sr-only">Loading…</span>`)                    | ⬜ OPEN | `htmx/loading.templ:29` — has `role="status"` + `aria-live` but no sr-only text                                       | bug-hunt-status:143     |
| 2   | `SanitizeID` mismatch in `ValidationSummary` — error links point to `#user-email` (sanitized) but field IDs are whatever consumer set (e.g. `user_email`) | ⬜ OPEN | `forms/ids.go:9` + `forms/validation.templ:61` — `SanitizeID` transforms `_`→`-` but doesn't match actual element IDs | bug-hunt-status:129     |
| 3   | `FromError` returns `FamilyInfrastructure` (→503) for unknown errors instead of `FamilyCorruption` (→500)                                                 | ⬜ OPEN | `errorpage/fromerror.go:71` — generic `errors.New("nil pointer")` returns 503 (implies temporary outage)              | bug-hunt-status:127     |
| 4   | `Footer` doesn't accept `BaseProps` — can't set Class/ID/Attrs (API inconsistency; every other component embeds BaseProps)                                | ⬜ OPEN | `navigation/nav.templ:119` — `templ Footer(brandText string)`                                                         | bug-hunt-status:180,188 |
| 5   | ErrorPage / NotFound404 missing `<main>` landmark — only `<div role="region">`, failing WCAG 2.4.1 (Bypass Blocks)                                        | ⬜ OPEN | `errorpage/errorpage.templ:7`, `errorpage/notfound404.templ:28`                                                       | bug-hunt-status:135     |
| 6   | `FormProps` CSRF token name hardcoded (`name="csrf_token"`) — frameworks use different names (Django, Rails, Spring)                                      | ⬜ OPEN | `forms/form.templ:71` — needs `CSRFTokenName` field                                                                   | bug-hunt-status:128     |
| 7   | `grid-rows-[0fr]` CSS output never verified against compiled Tailwind v4 — accordion collapse depends on it                                               | ⬜ OPEN | `display/accordion.templ:79` — test asserts class string present, not that it generates correct CSS                   | bug-hunt-status:155     |

---

## P1 — Testing gaps

| #   | Task                                                                                                                                                                                                                                                  | Status  | Evidence                                                              | Source                 |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | --------------------------------------------------------------------- | ---------------------- |
| 8   | Add regression tests for 18 untested Round-2 bug fixes (htmx LoadingButton/InlineLoadingOverlay/retry/announcer/catch-all/ConfirmDelete/SwapOOB, errorpage a11y, layout ThemeToggle/localStorage/FOUC/SRI, navigation LoadMore/breadcrumb/SidebarNav) | ⬜ OPEN | ~20% coverage for Round-2 fixes                                       | bug-hunt-status:94-113 |
| 9   | `InlineLoadingOverlay` `role="status"` assertion — BDD and snapshot tests don't check for it                                                                                                                                                          | ⬜ OPEN | `htmx/bdd_test.go:40`, `htmx/snapshot_test.go:26`                     | code-verification      |
| 10  | Dark golden test variants — render components inside `<div class="dark">` wrapper and generate dark-mode golden files                                                                                                                                 | ⬜ OPEN | No `*dark*` golden test files exist                                   | dark-mode-plan:213     |
| 11  | Toast JS-created toast golden test — `tcShowToast()` dynamically constructs HTML; only templ path is tested                                                                                                                                           | ⬜ OPEN | `feedback/toast.templ` JS path untested                               | dark-mode-plan:212     |
| 12  | Coverage → 80%+ on packages still below: errorpage (~73%), feedback (~73%), forms (~73%), navigation (~73%)                                                                                                                                           | ⬜ OPEN | `go tool cover -func` — gaps mostly in generated templ error branches | coverage-sprint        |
| 13  | Visual regression testing (Playwright screenshot diff light/dark)                                                                                                                                                                                     | ⬜ OPEN | No browser testing infrastructure exists                              | bug-hunt-status:157    |

---

## P2 — Pre-commit / CI hardening

| #   | Task                                                                      | Status  | Evidence                                                         | Source          |
| --- | ------------------------------------------------------------------------- | ------- | ---------------------------------------------------------------- | --------------- |
| 14  | Add `encoding/json/v2` grep guard to pre-commit hook                      | ⬜ OPEN | `scripts/pre-commit.sh` — no guard; the import broke builds 3×   | css-cleanup:92  |
| 15  | Pre-commit lint uses hardcoded package paths instead of `./...`           | ⬜ OPEN | `scripts/pre-commit.sh:22` — omits `cmd/`, may miss new packages | css-cleanup:100 |
| 16  | Document `encoding/json/v2` prohibition in AGENTS.md (blocked on Go 1.27) | ⬜ OPEN | No mention in AGENTS.md                                          | css-cleanup:102 |

---

## P2 — Documentation accuracy

| #   | Task                                                                                               | Status  | Evidence                                                                 | Source            |
| --- | -------------------------------------------------------------------------------------------------- | ------- | ------------------------------------------------------------------------ | ----------------- |
| 17  | Fix AGENTS.md lint path typo: `./svg/...` → `./internal/svg/...`                                   | ⬜ OPEN | `AGENTS.md:317` — path doesn't exist, causes lint error                  | v0.10-release:108 |
| 18  | Add note to CHANGELOG `[0.9.1]` that it was never tagged (changes included in v0.10.0)             | ⬜ OPEN | `CHANGELOG.md:111-126` — no "untagged" note; consumers may try `@v0.9.1` | v0.10-release:110 |
| 19  | ROADMAP.md doesn't mention dark mode compliance milestone                                          | ⬜ OPEN | `ROADMAP.md` — no dark mode row                                          | v0.10-release:64  |
| 20  | Create `docs/migration/v0.9-to-v0.10.md` migration guide                                           | ⬜ OPEN | `docs/migration/` has v0.7→v0.8, v0.8→v0.9 only                          | v0.10-release:61  |
| 21  | Update FEATURES.md with `templates/app.css` + BuildFlow `tailwind-build` provider entry            | ⬜ OPEN | `FEATURES.md` — no mention of CSS automation                             | css-cleanup:58    |
| 22  | AGENTS.md "Post-v0.9.0 Conventions" section header is stale (shipped in v0.10.0) — rename or merge | ⬜ OPEN | `AGENTS.md` — section still named "Post-v0.9.0"                          | v0.10-release:66  |
| 23  | AGENTS.md claims "61 generated files" but actual count is 62                                       | ⬜ OPEN | `AGENTS.md:53` vs `find . -name '*_templ.go' \| wc -l` = 62              | code-verification |

---

## P2 — Code quality & consistency

| #   | Task                                                                                                                                                 | Status  | Evidence                                                                                 | Source              |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ---------------------------------------------------------------------------------------- | ------------------- |
| 24  | Wire shared motion constants into remaining ~19 transition-bearing components (only 3/22 use `transitionFast`/`transitionNormal`/`transitionColors`) | ⬜ OPEN | `display/shared.go` — Accordion, Tabs, Dropdown, Tooltip, Toast, etc. use inline strings | session11:77        |
| 25  | `FamilyFromErrorFamily` stuttery name — consider `FromErrorFamily`                                                                                   | ⬜ OPEN | `errorpage/fromerror.go:12`                                                              | naming-review       |
| 26  | Consolidate CHANGELOG "Round 1"/"Round 2" headings into single `### Fixed` section                                                                   | ⬜ OPEN | `CHANGELOG.md` — internal audit concepts in consumer-facing doc                          | bug-hunt-status:161 |

---

## P3 — Polish & community

| #   | Task                                                             | Status  | Evidence                                                      | Source            |
| --- | ---------------------------------------------------------------- | ------- | ------------------------------------------------------------- | ----------------- |
| 27  | Demo site / showcase (live rendered components)                  | ⬜ OPEN | `examples/demo` is a single-page binary, not a hosted catalog | multiple-feedback |
| 28  | `awesome-templ` PR submission (updated component count)          | ⬜ OPEN | Never submitted                                               | session10:161     |
| 29  | `templ.guide` listing submission                                 | ⬜ OPEN | Never submitted                                               | session10:162     |
| 30  | Configure SSH tag signing (`gpg.ssh.allowedSignersFile`)         | ⬜ OPEN | Tags are annotated but not SSH-signed                         | v0.8-release:82   |
| 31  | Blocks/composition examples (dashboard, login, settings layouts) | ⬜ OPEN | No formal "blocks" showing real-world composition             | research:74       |
| 32  | Standalone `/forms` quickstart demo route                        | ⬜ OPEN | Forms shown inline only, no dedicated showcase                | session7:79       |

---

## v1.0 — Deferred breaking changes

| #   | Task                                                                                  | Status  | Evidence                                          | Source       |
| --- | ------------------------------------------------------------------------------------- | ------- | ------------------------------------------------- | ------------ |
| 33  | `Validate() error` methods on props structs (83 components, design decision needed)   | ⬜ OPEN | No `Validate` methods exist                       | multiple     |
| 34  | Move test helpers to `internal/testutil/` (70+ test files depend on exported helpers) | ⬜ OPEN | No `internal/testutil/` directory                 | multiple     |
| 35  | Self-host htmx as default, CDN opt-in (ADR 0007 written)                              | ⬜ OPEN | `layout/sri.go:69` — CDN is still default         | ADR-0007     |
| 36  | Semantic token layer `bg-tc-primary` (ADR 0008 written, 256 color refs)               | ⬜ OPEN | All components use hardcoded `bg-blue-600` etc.   | ADR-0008     |
| 37  | Icon RTL mirroring (`data-tc-dir-icon` + CSS `scaleX(-1)`) — 5 directional icons      | ⬜ OPEN | Audit done at `docs/audits/icon-rtl-mirroring.md` | session11:90 |
| 38  | Remove deprecated aliases (`AlertType`, `ToastType`)                                  | ⬜ OPEN | Still kept as type aliases for backward compat    | session8:120 |

---

## v2.0 — Architectural changes

| #   | Task                                                            | Status  | Evidence                                          | Source       |
| --- | --------------------------------------------------------------- | ------- | ------------------------------------------------- | ------------ |
| 39  | Compound component pattern (Trigger/Content/Close) for overlays | ⬜ OPEN | Current Modal/Drawer are monolithic               | research:70  |
| 40  | Native `<dialog>` element for Modal/Drawer                      | ⬜ OPEN | JS-based focus trap; browser `<dialog>` is better | research:72  |
| 41  | Headless/unstyled component variants (Radix UI model)           | ⬜ OPEN | All components ship with Tailwind classes         | session8:126 |
| 42  | CLI tool (`templ-components add <component>`, shadcn-style)     | ⬜ OPEN | No CLI exists                                     | session8:127 |

---

## New components — Not started

| #   | Component                                                      | Priority | Source                                           |
| --- | -------------------------------------------------------------- | -------- | ------------------------------------------------ |
| 43  | `Popover`                                                      | High     | Most requested missing component across feedback |
| 44  | `DataTable` (high-level sortable/filtering/pagination wrapper) | High     | DiscordSync + Overview feedback                  |
| 45  | `FilterDropdown`                                               | Medium   | Consumer-requested for filter bars               |
| 46  | `Slider` (ARIA slider pattern)                                 | Medium   | Research §5                                      |
| 47  | `Rating` (star rating, keyboard support)                       | Low      | Research §5                                      |
| 48  | `TagsInput`                                                    | Low      | Research §5                                      |
| 49  | `ContextMenu` (right-click menu)                               | Low      | Research §5                                      |
| 50  | `Carousel`                                                     | Low      | Research §5                                      |
| 51  | `HoverCard`                                                    | Medium   | Research §5                                      |
| 52  | `Calendar` (full calendar grid)                                | Medium   | Research §5                                      |

---

## Done — Verified complete (not actionable)

These were frequently listed as open in older reports but are confirmed DONE in v0.12.0:

- ✅ Dark mode compliance (30+ `dark:` variants fixed, `TestDarkModeCompliance` + `TestDarkModeSemanticColors` failing tests block CI)
- ✅ All 30+ typed enums have `IsValid()` + tests (32 production methods verified)
- ✅ `Footer` — Nav links wrap gracefully (v0.12.0)
- ✅ `NotFound404` component with search, links, `LinksTitle`, `WriteNotFound404`
- ✅ Sortable `TableHeader` + `TypedHeaders` with `aria-sort`
- ✅ `Form.Inline` horizontal layout
- ✅ `Grid.ContainerResponsive` + `GridProps.Gap` typed enum
- ✅ `layout.Script` + `layout.Stylesheet` CSP-safe helpers
- ✅ `feedback.SkeletonCardGrid`, `display.CopyButton`, `display.RelativeTime`, `display.CountBadge`, `display.Image`
- ✅ `navigation.LoadMore` with `InfiniteScroll` + cursor pagination recipe
- ✅ `display.Grid` responsive grid
- ✅ `SimpleNav.RightItems` slot
- ✅ `StatCard.Href` + typed HTMX fields (`HxGet`/`HxTarget`/`HxSwap: htmx.SwapStyle`)
- ✅ `Card.Body` / `SimpleCard.Body` / `Table.Body` slots
- ✅ RTL logical properties migration (0 physical properties remain)
- ✅ Motion-reduce compliance (0 gaps, grep test enforced)
- ✅ `OverlayKind` typed enum
- ✅ `icons.Close` alias for `icons.X`
- ✅ `color-scheme: light/dark` CSS
- ✅ `prefers-color-scheme` + `prefers-reduced-transparency` support
- ✅ Pre-commit has `go build ./...`
- ✅ Demo has ThemeToggle, TableHeader sortable, Form.Inline showcases
- ✅ Benchmark suites in 7 packages
- ✅ Fuzz tests for InputType, FormMethod, ButtonHTMLType
- ✅ ADR 0007 (self-host htmx), ADR 0008 (semantic tokens), ADR 0009 (accepted clones), ADR 0010 (sub-template extraction), ADR 0011 (dark mode convention)
- ✅ `templates/app.css` starter + BuildFlow `tailwind-build` provider (v0.11.0)
- ✅ `tc-css` CLI deleted (over-engineered, YAGNI)
