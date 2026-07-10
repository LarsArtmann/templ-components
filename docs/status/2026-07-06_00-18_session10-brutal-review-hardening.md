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

# Status Report — Session 10: Brutal Self-Review & Comprehensive Hardening

> **Updated:** 2026-07-06 (post-v0.8.0 final review). Version at report: 0.7.0 → **Current:** 0.8.0

**Date:** 2026-07-06 00:18
**Version:** 0.7.0 (unreleased — changes target v0.8.0) → **Released as v0.8.0**
**Commits:** 7 (ced952b → 7778f95)
**Files changed:** 40 files, +1032 insertions, -171 deletions
**Tests:** 2,192 test cases, 13/13 packages green, 0 lint issues

> **UPDATE NOTE (2026-07-06):** This session's work was released as v0.8.0 (commit `2d2d127`).
> All "Not Started" items (#1–3: release, CHANGELOG, version bump) were completed in the
> v0.8.0 release session. The combobox `focusout` handler (partially done below) was shipped
> in `de8171c`. The TableHeader golden test was added in `cc88d41`. See status annotations below.

---

## Context

The user asked for a brutal self-review: "What did you forget? What could you have done better?" — then a comprehensive Pareto-sorted execution plan, then full implementation. This session delivered all three phases.

The brutal truth that surfaced: **the previous sessions lied in TODO_LIST.md** — claiming `ModalSize2XL` was fixed when it wasn't. 16 IsValid functions were dead code (0 callers, 0 tests). 6 lookup maps used `string` keys despite typed enums existing. FEATURES.md had a phantom BadgeType value and wrong version number.

---

## A) FULLY DONE ✅

### Bug Fixes

| #   | Fix                                                                                                                                                                                                | Severity | Commit  |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| 1   | **ModalSize2XL/DrawerSize2XL value `"full"`→`"2xl"`** — both had identical values to their deprecated aliases, working only by map-key accident. Now each has its own value + dedicated map entry. | Critical | ced952b |
| 2   | **FEATURES.md version 0.6.1→0.7.0** — was wrong since v0.7.0 release                                                                                                                               | High     | ced952b |
| 3   | **BadgeType phantom "Default" value** — listed in FEATURES.md but doesn't exist in code                                                                                                            | Medium   | ced952b |
| 4   | **Tooltip stale "Known Issue"** — claimed `tooltipLookupPosition()` called twice; already fixed (cached in `pos` variable)                                                                         | Low      | ced952b |
| 5   | **FeedbackType missing from FEATURES.md** — AlertType/ToastType are aliases for FeedbackType, not separate types. Doc now reflects this.                                                           | Medium   | ced952b |
| 6   | **TODO_LIST:184 lie** — claimed ModalSize2XL "FIXED: value changed from 'full' to '2xl'" when it wasn't. Now actually fixed + TODO updated.                                                        | Critical | ced952b |

### Type Safety Improvements

| #   | Change                                                                                                                                                                                                                                                                                                    | Files               | Commit  |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- | ------- |
| 7   | **6 lookup maps converted from `map[string]string` to typed-key maps** — `cardPaddingLookup`, `avatarSizeLookup`, `avatarDotSizeLookup`, `badgeSizeLookup` (display); `spinnerSizeLookup`, `progressHeightLookup` (feedback). Eliminated all `string(v)` casts.                                           | 8 files             | 766b754 |
| 8   | **`CauseItem.Code` changed from `string` to `Code` type** — the `Code` type existed in the same package but wasn't used on this struct.                                                                                                                                                                   | errorpage/styles.go | 766b754 |
| 9   | **14 missing IsValid methods added** — AvatarStatus, DropdownItemKind, DropdownPosition, TabsVariant, OverlayKind, ButtonSize, ButtonHTMLType (display); StepIndicatorOrientation (feedback); ToggleSize, InputType, FormMethod (forms); SwapStyle (htmx); Name (icons). Total enums with IsValid: 16→30. | 5 files             | 3e10d60 |
| 10  | **All 30 IsValid functions now tested** — table-driven tests with valid + invalid inputs across 5 packages. Eliminates the dead-code ghost system (16 functions with 0 callers → 30 functions with test coverage).                                                                                        | 5 new test files    | 3e10d60 |

### Accessibility

| #   | Change                                                                                                                                                                                   | Commit  |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| 11  | **Combobox `aria-selected` on active option** — was using non-standard `data-selected` only. Now sets both `data-selected` and `aria-selected="true"` for screen reader compliance.      | 0ee6bb1 |
| 12  | **Combobox Tab-to-close + cleanup** — Tab key now closes listbox and clears selection state (`data-selected`, `aria-activedescendant`). Previously Tab just moved focus without cleanup. | 0ee6bb1 |
| 13  | **Combobox `tcClearComboSelection()` helper** — extracted DRY cleanup used across Escape/Enter/Tab/navigation paths.                                                                     | 0ee6bb1 |

### New Features

| #   | Feature                                                                                                                                                                                                                                                                                                                                                                                | Commit  |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| 14  | **Sortable TableHeader** — `TableHeader` struct with `Sortable bool`, `SortDirection` enum (None/Asc/Desc), `Href` for server-side sort links. Renders `aria-sort="ascending/descending/none"` on sortable columns, ↑/↓ visual indicators, clickable `<a>` when Href is set. `TypedHeaders []TableHeader` on TableProps takes precedence over `Headers []string`. Backward compatible. | 74da41d |
| 15  | **Form.Inline horizontal layout** — `Inline bool` field on FormProps renders `flex flex-wrap items-end gap-3` instead of `space-y-6`. Follows the exact `RadioGroup.Inline` precedent.                                                                                                                                                                                                 | 74da41d |

### Test Coverage

| #   | Change                                                                                                                                                                      | Before → After            | Commit  |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------- | ------- |
| 16  | **htmx coverage boost** — ConfirmDelete full props, SwapOOB all 8 swap styles + invalid fallback, CSRFToken                                                                 | 68.5% → 75.7%             | b430980 |
| 17  | **layout coverage boost** — Stylesheet test (was 0%!), Script with attrs                                                                                                    | 69.6% → 74.5%             | b430980 |
| 18  | **All 13 packages now ≥70% coverage** — display 70.4%, feedback 72.5%, forms 72.3%, errorpage 73.0%, layout 74.5%, htmx 75.7%, utils 77.6%, icons 78.6%, internal/svg 79.0% | 3 packages were below 70% | b430980 |

### Documentation

| #   | Change                                                                                                                                                                                 | Commit  |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| 19  | **3 recipe docs** — custom-table-rows.md (Body slot + TypedHeaders sortable columns), custom-404-page.md (NotFound404 with custom links/search), recipe-index.md (links all 5 recipes) | 7778f95 |
| 20  | **errorpage/doc.go updated** — NotFound404 added to component list                                                                                                                     | 7778f95 |
| 21  | **FEATURES.md updated** — enum count 32→33 (+SortDirection), SortDirection added to display enums table                                                                                | 7778f95 |

---

## B) PARTIALLY DONE 🟡

| Item              | What's done                                                                                                                                               | What's missing (Status 2026-07-06)                                                                                                      |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| Combobox WAI-ARIA | aria-selected, Tab cleanup, ArrowDown/Up/Home/End/Enter/Escape all present. Type-ahead handled by input filter (correct for filterable combobox pattern). | ✅ **`focusout` handler shipped** (`de8171c`) — listbox closes and `aria-activedescendant` clears on blur. All WAI-ARIA gaps closed.    |
| Coverage          | All packages ≥70%.                                                                                                                                        | Unchanged: errorpage 72.9%, feedback 72.3%, forms 72.3%, navigation 72.6%. Could reach 75%+ with targeted tests.                        |
| IsValid system    | All 30 closed-set enums have IsValid methods + tests.                                                                                                     | `layout.HTMXVersion` enum has no IsValid (open-set, changes per release — arguably not needed). Unchanged.                              |
| TableHeader       | TypedHeaders with aria-sort + indicators fully working + tested.                                                                                          | ✅ **Golden test added** (`cc88d41`) — `TestGoldenTableSortableHeaders` renders 3-column table with aria-sort, sort links, ↑ indicator. |

---

## C) NOT STARTED ⬜

| #   | Item                                                                                                                                          | Status (2026-07-06)                                 |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------- |
| 1   | **v0.8.0 release** — all changes target this version but no release commit/tag cut yet                                                        | ✅ **Done** — v0.8.0 released (`2d2d127`)           |
| 2   | **CHANGELOG [Unreleased] entries** — none of the 7 commits this session added CHANGELOG entries                                               | ✅ **Done** — all entries added, released in v0.8.0 |
| 3   | **`utils.Version` bump to 0.8.0** — still says 0.7.0                                                                                          | ✅ **Done** — `utils.Version = "0.8.0"`             |
| 4   | **Sortable DataTable component** — TableHeader provides the type, but no high-level DataTable component that auto-manages sort state          | ⬜ Not started                                      |
| 5   | **Filter dropdown component** — recipe documents the manual pattern; no purpose-built component exists                                        | ⬜ Not started                                      |
| 6   | **`forms.InlineForm`** vs `Form.Inline` — the Inline field is done but a dedicated InlineForm constructor function might be more discoverable | ⬜ Not started (low priority)                       |
| 7   | **Demo app showcase** — new TableHeader and Form.Inline features not yet showcased in examples/demo                                           | ⬜ Not started                                      |

---

## D) TOTALLY FUCKED UP 💥

### Nothing critical this session — but documenting pre-existing damage found:

| #   | Issue                                                                                                                                                                                                                                                                                          | Severity     | Status                                |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------ | ------------------------------------- |
| 1   | **TODO_LIST.md lied about ModalSize2XL being fixed** — "FIXED: value changed from 'full' to '2xl'" was verifiably false. The previous session claimed to fix it but didn't. This is the kind of documentation lie that erodes trust in the entire TODO_LIST.                                   | **Critical** | Fixed this session (ced952b)          |
| 2   | **16 IsValid functions were dead code** — exported, documented, zero callers, zero tests. A ghost validation system that gave false confidence about type safety. Built across sessions 7-8 and never tested.                                                                                  | **High**     | Fixed this session (3e10d60)          |
| 3   | **6 lookup maps used `string` keys** despite typed enums existing — `badgeSizeLookup[string(v)]` instead of `badgeSizeLookup[v]`. The typed enums were created in earlier sessions but the maps were never updated, creating a split brain where the type system was bypassed at every lookup. | **High**     | Fixed this session (766b754)          |
| 4   | **FEATURES.md had 5 separate factual errors** — wrong version (0.6.1 vs 0.7.0), phantom BadgeType "Default" value, stale Tooltip known issue, missing FeedbackType enum, disagreeing coverage numbers. Documentation drift from sessions 7-9.                                                  | **Medium**   | Fixed this session (ced952b, 7778f95) |

---

## E) WHAT WE SHOULD IMPROVE 🔧

### Process Improvements

1. **Every IsValid function must ship with a test in the same commit** — the 16-function dead-code ghost system happened because IsValid methods were added without tests or callers. Enforce: no IsValid without a test.
2. **TODO_LIST claims must be verified before marking [x]** — the ModalSize2XL lie happened because someone trusted a previous claim. Rule: re-run the grep/test before marking done.
3. **Doc files (FEATURES.md, README.md) need a version-sync test** — like the existing `TestVersionMatchesChangelog`, add a test that FEATURES.md version matches `utils.Version`.
4. **Lookup maps should use typed keys from creation** — the `string`-key maps were created alongside typed enums but never connected. Establish: if a typed enum exists, its lookup map MUST use it as the key type.

### Code Quality Improvements Still Open

5. ✅ **`feedbackIconName()` and `FamilyStatusCode()` use manual map+fallback** — **Fixed** (`d3c8b88`): both replaced with `utils.Lookup()` calls.
6. ✅ **`StatCardProps.HxSwap` uses raw `string`** — **Fixed** (`cc88d41`): typed as `htmx.SwapStyle`.
7. ✅ **`ButtonHTMLType` uses `map[X]bool`** — **Fixed** (`cc88d41`): converted to `map[X]string` + `utils.Lookup`.

---

## F) Top 25 Things to Do Next

Sorted by impact × effort × customer value.

| #   | Task                                                                                                        | Status (2026-07-06)     |
| --- | ----------------------------------------------------------------------------------------------------------- | ----------------------- |
| 1   | **Cut v0.8.0 release** — bump version, CHANGELOG, tag, push                                                 | ✅ **Done** (`2d2d127`) |
| 2   | **Add CHANGELOG [Unreleased] entries** for all 7 commits this session                                       | ✅ **Done**             |
| 3   | **Add FEATURES.md version-sync test** (like TestVersionMatchesChangelog)                                    | ✅ **Done** (`6e94f93`) |
| 4   | **Demo app: showcase TableHeader sortable columns**                                                         | ⬜ Not started          |
| 5   | **Demo app: showcase Form.Inline**                                                                          | ⬜ Not started          |
| 6   | **Golden test for TableHeader sortable variant**                                                            | ✅ **Done** (`cc88d41`) |
| 7   | **StatCardProps.HxSwap: change `string` → `htmx.SwapStyle`**                                                | ✅ **Done** (`cc88d41`) |
| 8   | **ButtonHTMLType: convert `map[X]bool` → `map[X]string` + Lookup**                                          | ✅ **Done** (`cc88d41`) |
| 9   | **`feedbackIconName` + `FamilyStatusCode`: use `utils.Lookup`**                                             | ✅ **Done** (`d3c8b88`) |
| 10  | **Combobox `focusout` handler** — clear aria-activedescendant on blur                                       | ✅ **Done** (`de8171c`) |
| 11  | **Sortable DataTable component** — high-level wrapper around TableHeader                                    | ⬜ Not started          |
| 12  | **Filter dropdown component** — purpose-built for filter bars                                               | ⬜ Not started          |
| 13  | **Move test helpers to `internal/testutil/`** — deferred to v1.0 but plan it                                | ⬜ Deferred to v1.0     |
| 14  | **Add `Validate() error` to props structs** — v1.0 scope, but design now                                    | ⬜ Deferred to v1.0     |
| 15  | **errorpage coverage to 80%+** — handler edge paths, write failures                                         | ⬜ Not done (72.9%)     |
| 16  | **feedback coverage to 80%+** — StepIndicator branches, LoadingOverlay                                      | ⬜ Not done (72.3%)     |
| 17  | **forms coverage to 80%+** — Combobox rendering branches, RadioGroup                                        | ⬜ Not done (72.3%)     |
| 18  | **navigation coverage to 80%+** — SidebarNav, Breadcrumbs JSON-LD                                           | ⬜ Not done (72.6%)     |
| 19  | **AGENTS.md update** — document TableHeader, Form.Inline, typed-map convention, IsValid-test convention     | ✅ **Done** (`a0dbae7`) |
| 20  | **Icons-only adoption doc update** — mention new icons added since v0.7.0                                   | ⬜ Not done             |
| 21  | **awesome-templ PR submission** — component count updated, submit the prepared entry                        | ⬜ Not done             |
| 22  | **templ.guide listing submission** — prepared but never submitted                                           | ⬜ Not done             |
| 23  | **Tooltip: add `aria-describedby` via `props.ID`** — investigate if CSS-only tooltip needs JS for full a11y | ⬜ Not started          |
| 24  | **Pagination: add `rel="canonical"` for page 1** — SEO improvement                                          | ✅ **Done** (`098f7c3`) |
| 25  | **Add `TableHeader.IsValid` / `SortDirection.IsValid`** — complete the enum validation set                  | ✅ **Done** (`cc88d41`) |

**Scorecard:** 12 of 25 complete (48%).

---

## G) Top #1 Question I Cannot Answer

> ✅ **RESOLVED — DECISION: KEEP IT FLAT.** Table types stay in `display/` as `display.TableHeader`,
> `display.SortDirection`, etc. The breaking change cost of splitting into a sub-package isn't
> worth the organizational benefit until v1.0. This matches the established pattern — all other
> components live flat in `display/`. Confirmed by v0.8.0 release which kept types flat.

---

## Session Stats

| Metric                      | Value                                 |
| --------------------------- | ------------------------------------- |
| Commits                     | 7                                     |
| Files changed               | 40                                    |
| Lines added                 | +1,032                                |
| Lines removed               | -171                                  |
| New test files              | 7                                     |
| New IsValid methods         | 14                                    |
| Total IsValid methods       | 30 (was 16)                           |
| Bugs fixed                  | 6                                     |
| Features added              | 2 (TableHeader sortable, Form.Inline) |
| Ghost systems eliminated    | 1 (dead IsValid functions)            |
| Packages below 70% coverage | 0 (was 3)                             |
| Lint issues                 | 0                                     |
| Test cases                  | 2,192                                 |
| Typed enums                 | 33                                    |
| Recipe docs                 | 5 (was 2)                             |
