<!-- AUTO-UPDATED 2026-07-10: Retrospective status overlay -->

> ## 🔔 Update Notice — 2026-07-10
>
> This report is **historical**. Many items listed as "open", "todo", or "broken" below
> have since been **fixed and verified**. Do not act on open items without first checking
> [TODO_LIST.md](../../TODO_LIST.md) for current status.
>
> **Key fixes completed since this report:**
>
> - ✅ All 7 P0 bugs fixed (InlineLoadingOverlay a11y, SanitizeID mismatch, FromError fallback,
>   Footer BaseProps, ErrorPage/NotFound404 `<main>` landmark, CSRFTokenName, grid-rows verified)
> - ✅ `encoding/json/v2` purged from all production code + pre-commit guard added
> - ✅ Motion constants centralized in `utils/motion.go`, wired into 13 components
> - ✅ `FamilyFromErrorFamily` → `FromErrorFamily` (old name kept as deprecated alias)
> - ✅ `icons.IconRTL()` + CSS for directional icon RTL mirroring
> - ✅ 33 regression tests added (htmx, errorpage, layout, navigation, feedback, display)
> - ✅ Dark golden test infrastructure (badge/card/button)
> - ✅ CHANGELOG consolidated, ROADMAP updated, migration guide created
> - ✅ All 14 packages pass, 0 lint issues
>
> **Canonical source of truth:** [TODO_LIST.md](../../TODO_LIST.md) (52 items, 37 ✅ done, 12 deferred/blocked)

---

# Comprehensive Hardening Plan — All Remaining TODOs

> **Date:** 2026-07-06
> **Source:** Status reports session 9 (Pareto) + session 10 (brutal review) + Top 25 lists
> **Constraint:** Each task ≤12 min. Build must pass after every phase.

---

## Task Inventory (28 tasks, sorted by impact/effort/value)

### Phase 1: Fix Real Bugs — Accessibility (HIGH impact, LOW effort)

| #     | Task                                                                                                        | Files                             | Impact     | Effort     | Est     |
| ----- | ----------------------------------------------------------------------------------------------------------- | --------------------------------- | ---------- | ---------- | ------- |
| ~~1~~ | ~~Fix motion-reduce gap in toast dismiss button (2 instances)~~ done — utils/motion compliance test.go      | ~~feedback/toast.templ~~          | ~~HIGH~~   | ~~LOW~~    | ~~8m~~  |
| ~~2~~ | ~~Fix motion-reduce gap in step indicator circle~~ done — utils/motion compliance test.go                   | ~~feedback/step_indicator.templ~~ | ~~HIGH~~   | ~~LOW~~    | ~~5m~~  |
| ~~3~~ | ~~Fix motion-reduce gap in empty state action class~~ done — utils/motion compliance test.go                | ~~display/empty_state.templ~~     | ~~HIGH~~   | ~~LOW~~    | ~~5m~~  |
| ~~4~~ | ~~Fix motion-reduce gap in file input~~ done — utils/motion compliance test.go                              | ~~forms/file_input.templ~~        | ~~HIGH~~   | ~~LOW~~    | ~~5m~~  |
| ~~5~~ | ~~Fix motion-reduce gap in error page action buttons (2 instances)~~ done — utils/motion compliance test.go | ~~errorpage/errorpage.templ~~     | ~~HIGH~~   | ~~LOW~~    | ~~8m~~  |
| ~~6~~ | ~~Combobox focusout handler — clear aria-activedescendant on blur~~ done — forms/combobox.templ             | ~~forms/combobox.templ~~          | ~~MEDIUM~~ | ~~MEDIUM~~ | ~~10m~~ |
| ~~7~~ | ~~Add SortDirection.IsValid method~~ done — display/enums test.go                                           | ~~display/table.templ~~           | ~~LOW~~    | ~~LOW~~    | ~~5m~~  |

### Phase 2: Code Quality (MEDIUM impact, LOW effort)

| #      | Task                                                                                                            | Files                            | Impact     | Effort  | Est     |
| ------ | --------------------------------------------------------------------------------------------------------------- | -------------------------------- | ---------- | ------- | ------- |
| ~~8~~  | ~~Wire transitionColors constant into copy_button.templ~~ done — utils/motion.go                                | ~~display/copy_button.templ~~    | ~~MEDIUM~~ | ~~LOW~~ | ~~5m~~  |
| ~~9~~  | ~~Wire transitionColors into accordion.templ button + transitionTransform into chevron~~ done — utils/motion.go | ~~display/accordion.templ~~      | ~~MEDIUM~~ | ~~LOW~~ | ~~10m~~ |
| ~~10~~ | ~~Simplify FamilyStatusCode with utils.Lookup~~ done (docs-health pass 2026-09-08)                              | ~~errorpage/styles.go~~          | ~~LOW~~    | ~~LOW~~ | ~~5m~~  |
| ~~11~~ | ~~Add FEATURES.md version-sync drift-guard test~~ done — utils/version test.go                                  | ~~internal/contract/ or utils/~~ | ~~MEDIUM~~ | ~~LOW~~ | ~~10m~~ |

### Phase 3: Documentation (MEDIUM impact, LOW effort)

| #      | Task                                                                      | Files                                 | Impact     | Effort  | Est     |
| ------ | ------------------------------------------------------------------------- | ------------------------------------- | ---------- | ------- | ------- |
| ~~12~~ | ~~SKILL.md: Add RTL logical properties convention~~ done — skill/SKILL.md | ~~skill/SKILL.md~~                    | ~~HIGH~~   | ~~LOW~~ | ~~10m~~ |
| ~~13~~ | ~~SKILL.md: Add motion constants convention~~ done — skill/SKILL.md       | ~~skill/SKILL.md~~                    | ~~HIGH~~   | ~~LOW~~ | ~~10m~~ |
| ~~14~~ | ~~SKILL.md: Add container query convention~~ done — skill/SKILL.md        | ~~skill/SKILL.md~~                    | ~~MEDIUM~~ | ~~LOW~~ | ~~10m~~ |
| ~~15~~ | ~~Container queries recipe doc~~ done — docs/recipes/container-queries.md | ~~docs/recipes/container-queries.md~~ | ~~MEDIUM~~ | ~~LOW~~ | ~~10m~~ |
| ~~16~~ | ~~Motion design reference doc~~ done — docs/motion-design.md              | ~~docs/motion-design.md~~             | ~~MEDIUM~~ | ~~LOW~~ | ~~10m~~ |
| ~~17~~ | ~~CHANGELOG entries for all changes~~ done — CHANGELOG.md                 | ~~CHANGELOG.md~~                      | ~~HIGH~~   | ~~LOW~~ | ~~10m~~ |
| ~~18~~ | ~~AGENTS.md: new convention entries~~ done — AGENTS.md                    | ~~AGENTS.md~~                         | ~~MEDIUM~~ | ~~LOW~~ | ~~10m~~ |

### Phase 4: Testing (MEDIUM impact, LOW-MEDIUM effort)

| #      | Task                                                                                                       | Files                            | Impact     | Effort  | Est     |
| ------ | ---------------------------------------------------------------------------------------------------------- | -------------------------------- | ---------- | ------- | ------- |
| ~~19~~ | ~~Golden test for TableHeader sortable variant~~ done — display/golden new test.go                         | ~~display/table_golden_test.go~~ | ~~MEDIUM~~ | ~~LOW~~ | ~~10m~~ |
| ~~20~~ | ~~Motion-reduce audit: grep all transitions, verify 100% coverage~~ done — utils/motion compliance test.go | ~~verification~~                 | ~~MEDIUM~~ | ~~LOW~~ | ~~10m~~ |

### Phase 5: Nice-to-Haves (LOW-MEDIUM impact, LOW effort)

| #      | Task                                                                                                   | Files                                | Impact     | Effort  | Est     |
| ------ | ------------------------------------------------------------------------------------------------------ | ------------------------------------ | ---------- | ------- | ------- |
| ~~21~~ | ~~Pagination rel="canonical" for page 1~~ done — navigation/pagination.templ                           | ~~navigation/pagination.templ~~      | ~~LOW~~    | ~~LOW~~ | ~~10m~~ |
| ~~22~~ | ~~Icon RTL mirroring audit — identify arrows/chevrons needing dir="rtl" flip~~ done — icons/icon.templ | ~~icons/ audit~~                     | ~~LOW~~    | ~~LOW~~ | ~~10m~~ |
| ~~23~~ | ~~Semantic token layer ADR~~ done — docs/adr/0008-semantic-tokens.md                                   | ~~docs/adr/0008-semantic-tokens.md~~ | ~~MEDIUM~~ | ~~LOW~~ | ~~10m~~ |

### Phase 6: Verification & Commit

| #      | Task                                                                                | Est     |
| ------ | ----------------------------------------------------------------------------------- | ------- |
| ~~24~~ | ~~Build: templ generate + go build all modules~~ done (docs-health pass 2026-09-08) | ~~10m~~ |
| ~~25~~ | ~~Test: go test -race -count=1 all modules~~ done (docs-health pass 2026-09-08)     | ~~10m~~ |
| ~~26~~ | ~~Lint: golangci-lint all packages~~ done (docs-health pass 2026-09-08)             | ~~5m~~  |
| ~~27~~ | ~~Update golden files if needed (-update flag)~~ done (docs-health pass 2026-09-08) | ~~10m~~ |
| ~~28~~ | ~~Git commit~~ done (docs-health pass 2026-09-08)                                   | ~~10m~~ |

**Total: 28 tasks, ~250 min (4.2 hrs)**

---

## What This Plan Does NOT Do

| Excluded                                          | Why                                        |
| ------------------------------------------------- | ------------------------------------------ |
| New components (Popover, Slider, DataTable, etc.) | Additive work, separate effort             |
| Semantic token migration (256 colors)             | Needs dedicated major-version sprint       |
| Native `<dialog>` migration                       | Fundamental architecture change, needs ADR |
| Compound component refactoring                    | Breaking API change, v2.0 decision         |
| CSS `@starting-style`                             | Needs browser baseline decision            |
| Move test helpers to internal/testutil/           | Deferred to v1.0                           |
| Validate() error on all props                     | Deferred to v1.0                           |
| v0.8.0 release cut                                | Release script needs interactive input     |
