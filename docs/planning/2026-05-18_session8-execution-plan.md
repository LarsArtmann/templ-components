# Comprehensive Execution Plan — Session 8

**Date:** 2026-05-18  
**Scope:** All remaining TODOs + all actionable improvements from audit  
**Method:** Each task ≤12 min, sorted by impact/effort/customer-value

## Source Data

### 7 Open TODO Items

- #23 JS pattern unification (P2, deferred — high risk)
- #24 Shared dismiss JS (P2, dedup)
- #49 Forms coverage 68.1→75%+ (P2)
- #51 Golden file tests (P2)
- #58 Move test helpers to internal/ (P3, breaking — defer)
- #59 Move ProgressBar test (P3)
- #71 Documentation site (P3, large effort)

### Coverage Gaps (packages below 70%)

- display: 67.2% (Avatar 59.4%, EmptyState 44-63%, Tooltip 57%, Badge 64%, Card 68%)
- icons: 68.3% (IconPathJS 0%)
- forms: 68.1% (Select 64.6%, Checkbox 66.7%, helpText 67.5%)

### Convention Violations

- `modalSizeClass` still uses switch (display/modal_go.go)
- `DefaultStepIndicatorProps` returns zero-value
- `DefaultNavLinkProps` returns zero-value
- `DefaultAlertProps` returns zero-value (line 36, coverage 0%)
- `DefaultToastProps` returns zero-value (line 82, coverage 0%)

### Untested Features

- Card: subtitle, footer, header action rendering
- Avatar: size variants, shape variants, status dot, status rendering
- Badge: pill shape, dot rendering, all types
- EmptyState: without icon, without action
- Tooltip: with custom class/ID
- Select: disabled options, no options
- Skeleton: custom dimensions

---

## Execution Plan — 35 Steps

### Tier 1: Real Bugs & Correctness (CUSTOMER VALUE: HIGH)

| #   | Task                                                   | Files                       | Effort | Impact        |
| --- | ------------------------------------------------------ | --------------------------- | ------ | ------------- |
| 1   | Convert `modalSizeClass` switch → map                  | `display/modal_go.go`       | 5min   | 🟡 Convention |
| 2   | Add meaningful defaults to `DefaultStepIndicatorProps` | `feedback/progress.templ`   | 3min   | 🟡 Convention |
| 3   | Add meaningful defaults to `DefaultNavLinkProps`       | `navigation/nav_link.templ` | 3min   | 🟡 Convention |
| 4   | Add meaningful defaults to `DefaultAlertProps`         | `feedback/alert.templ`      | 3min   | 🟡 Convention |
| 5   | Add meaningful defaults to `DefaultToastProps`         | `feedback/toast.templ`      | 3min   | 🟡 Convention |
| 6   | Test all 4 new defaults                                | Multiple test files         | 5min   | 🟢 Verify     |

### Tier 2: Coverage Push — Forms (68.1% → 70%+)

| #   | Task                                                       | Files                                         | Effort | Impact      |
| --- | ---------------------------------------------------------- | --------------------------------------------- | ------ | ----------- |
| 7   | Test Select with disabled options + no options             | `forms/select_test.go` or `forms/bdd_test.go` | 8min   | 🟡 Coverage |
| 8   | Test Checkbox with help text + checked state               | `forms/bdd_test.go`                           | 5min   | 🟡 Coverage |
| 9   | Test Label with help text only (no error) + without for ID | `forms/label_test.go`                         | 5min   | 🟡 Coverage |
| 10  | Test Textarea with rows + readonly + autofocus             | `forms/bdd_test.go`                           | 5min   | 🟡 Coverage |

### Tier 3: Coverage Push — Display (67.2% → 70%+)

| #   | Task                                                   | Files                                                | Effort | Impact      |
| --- | ------------------------------------------------------ | ---------------------------------------------------- | ------ | ----------- |
| 11  | Test Card with subtitle, footer, header action         | `display/card_test.go`                               | 8min   | 🟡 Coverage |
| 12  | Test Avatar size variants + shape + status             | `display/avatar_test.go`                             | 10min  | 🟡 Coverage |
| 13  | Test Badge pill shape + dot + all types                | `display/badge_test.go` or `display/helpers_test.go` | 8min   | 🟡 Coverage |
| 14  | Test EmptyState without icon + without action          | `display/card_test.go`                               | 5min   | 🟡 Coverage |
| 15  | Test Tooltip with custom class + ID + unknown position | `display/tooltip_test.go`                            | 5min   | 🟡 Coverage |
| 16  | Test Tabs with pills variant + no active tab           | `display/tabs_test.go`                               | 8min   | 🟡 Coverage |
| 17  | Test Dropdown with icon items + position right         | `display/dropdown_test.go`                           | 5min   | 🟡 Coverage |

### Tier 4: Coverage Push — Feedback (73.2% → 75%+)

| #   | Task                                         | Files                      | Effort | Impact      |
| --- | -------------------------------------------- | -------------------------- | ------ | ----------- |
| 18  | Test Skeleton with custom class + dimensions | `feedback/loading_test.go` | 5min   | 🟡 Coverage |
| 19  | Test Toast static render with all types      | `feedback/toast_test.go`   | 5min   | 🟡 Coverage |
| 20  | Test Alert dismiss script presence           | `feedback/alert_test.go`   | 3min   | 🟡 Coverage |

### Tier 5: Coverage Push — Icons (68.3% → 70%+)

| #   | Task                                               | Files                      | Effort | Impact      |
| --- | -------------------------------------------------- | -------------------------- | ------ | ----------- |
| 21  | Test `IconPathJS` function (currently 0% coverage) | `icons/icon_paths_test.go` | 5min   | 🟡 Coverage |

### Tier 6: Validation & Robustness

| #   | Task                                      | Files                   | Effort | Impact    |
| --- | ----------------------------------------- | ----------------------- | ------ | --------- |
| 22  | Test `SanitizeID` with special characters | `forms/helpers_test.go` | 3min   | 🟢 Verify |
| 23  | Test Table with Bordered=true             | `display/table_test.go` | 3min   | 🟢 Verify |
| 24  | Test Table with Hover=true                | `display/table_test.go` | 3min   | 🟢 Verify |

### Tier 7: JS Dedup (#24)

| #   | Task                                      | Files                                              | Effort | Impact    |
| --- | ----------------------------------------- | -------------------------------------------------- | ------ | --------- |
| 25  | Extract shared dismiss JS for Alert+Toast | `feedback/alert.templ`, `feedback/toast.templ`     | 10min  | 🟡 Dedup  |
| 26  | Test dismiss behavior for both components | `feedback/alert_test.go`, `feedback/toast_test.go` | 5min   | 🟢 Verify |

### Tier 8: Golden File Tests (#51)

| #   | Task                                   | Files                       | Effort | Impact            |
| --- | -------------------------------------- | --------------------------- | ------ | ----------------- |
| 27  | Create golden file test infrastructure | `internal/golden/golden.go` | 8min   | 🟡 Infrastructure |
| 28  | Convert Badge snapshot → golden file   | `display/badge_test.go`     | 5min   | 🟡 Migration      |
| 29  | Convert Card snapshot → golden file    | `display/card_test.go`      | 5min   | 🟡 Migration      |

### Tier 9: Code Organization (#59)

| #   | Task                                           | Files                                | Effort | Impact          |
| --- | ---------------------------------------------- | ------------------------------------ | ------ | --------------- |
| 30  | Move ProgressBar test from display to feedback | `display/a11y_test.go` → `feedback/` | 5min   | 🟢 Organization |

### Tier 10: Documentation (#71 partial)

| #   | Task                                            | Files                           | Effort | Impact  |
| --- | ----------------------------------------------- | ------------------------------- | ------ | ------- |
| 31  | Add package-level examples for top 5 components | `display/example_test.go`, etc. | 10min  | 🟢 Docs |
| 32  | Update README with current API surface          | `README.md`                     | 8min   | 🟢 Docs |

### Tier 11: Version & Release

| #   | Task                                      | Files          | Effort | Impact     |
| --- | ----------------------------------------- | -------------- | ------ | ---------- |
| 33  | Update CHANGELOG.md with v0.2.0 entries   | `CHANGELOG.md` | 5min   | 📝 Release |
| 34  | Tag v0.2.0 release (if user approves)     | git tag        | 2min   | 📝 Release |
| 35  | Update TODO_LIST.md date + mark new items | `TODO_LIST.md` | 2min   | 📝 Docs    |

---

## Items Explicitly Deferred

| Item                                             | Why                                                             |
| ------------------------------------------------ | --------------------------------------------------------------- |
| #23 JS pattern unification                       | High risk of breaking 3 working components. Low customer value. |
| #58 Move test helpers to internal/               | Breaking API change for public library consumers.               |
| #71 Full documentation site                      | Large effort (days), use pkgsite instead for now.               |
| Nonce parameter consistency                      | Breaking API change. Defer to v0.3.                             |
| Missing type fields (Disabled, MaxLength, Value) | Feature additions, not fixes. Separate PR after v0.2.0.         |

## Total Estimates

| Tier                     | Steps  | Time    | Cumulative |
| ------------------------ | ------ | ------- | ---------- |
| 1: Defaults & Convention | 6      | 22min   | 22min      |
| 2: Forms Coverage        | 4      | 23min   | 45min      |
| 3: Display Coverage      | 7      | 49min   | 94min      |
| 4: Feedback Coverage     | 3      | 13min   | 107min     |
| 5: Icons Coverage        | 1      | 5min    | 112min     |
| 6: Validation            | 3      | 9min    | 121min     |
| 7: JS Dedup              | 2      | 15min   | 136min     |
| 8: Golden Files          | 3      | 18min   | 154min     |
| 9: Organization          | 1      | 5min    | 159min     |
| 10: Documentation        | 2      | 18min   | 177min     |
| 11: Release              | 3      | 9min    | 186min     |
| **Total**                | **35** | **~3h** | —          |

## Expected Coverage Improvement

| Package   | Current   | Target   |
| --------- | --------- | -------- |
| display   | 67.2%     | ~72%     |
| forms     | 68.1%     | ~72%     |
| icons     | 68.3%     | ~71%     |
| feedback  | 73.2%     | ~76%     |
| **Total** | **69.5%** | **~72%** |
