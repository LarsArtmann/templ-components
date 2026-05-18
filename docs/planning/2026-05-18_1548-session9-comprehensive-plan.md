# Session 9 Execution Plan — Coverage Push, JS Cleanup, Release Prep

**Created:** 2026-05-18 15:48 | **Scope:** All remaining TODOs + discovered gaps
**Current State:** 673 tests, 70.7% coverage, 0 lint issues, 5 open TODOs

---

## Pareto Breakdown

### The 1% that delivers 51% of the result

- **Tag v0.2.0 release** — everything else is polish; the library is already production-quality

### The 4% that delivers 64% of the result

- Close the 4 critical coverage gaps (Skeleton 53.7%, emptyStateAction 44.1%, Avatar 59.4%, card_test compile errors)
- Unify JS attachment pattern (TODO #23) — last architectural inconsistency

### The 20% that delivers 80% of the result

- Push coverage toward 75%+ across all packages
- Golden file infrastructure (TODO #51)
- DefaultProps test coverage (6 functions at 0%)
- Documentation site (TODO #71) — pkgsite works, but custom docs would be superior

### Everything else

- Test helper API move (TODO #58) — breaking change, defer to v1.0
- ProgressBar test move (TODO #59) — cosmetic
- Examples/demo improvements

---

## Comprehensive Task Table

Sorted by: Impact × Customer Value / Effort (highest first)

| #   | Task                                                                              | Source   | Impact | Effort | Customer Value | Est.  |
| --- | --------------------------------------------------------------------------------- | -------- | ------ | ------ | -------------- | ----- |
|     | **TIER 1 — RELEASE (51% impact)**                                                 |          |        |        |                |       |
| 1   | Verify goreleaser config is correct and dry-run works                             | TODO #62 | HIGH   | 5min   | HIGH           | 5min  |
| 2   | Update CHANGELOG.md with final v0.2.0 entry                                       | Release  | HIGH   | 8min   | HIGH           | 8min  |
| 3   | Verify README.md metrics are accurate (tests, coverage, packages)                 | Release  | MED    | 5min   | MED            | 5min  |
| 4   | Tag v0.2.0 and push tag                                                           | Release  | HIGH   | 3min   | HIGH           | 3min  |
|     | **TIER 2 — COVERAGE GAPS (64% impact)**                                           |          |        |        |                |       |
| 5   | Fix card_test.go compile errors — align with SimpleCard API                       | Bug      | HIGH   | 8min   | HIGH           | 8min  |
| 6   | Test Skeleton image variant                                                       | Gap      | MED    | 5min   | MED            | 5min  |
| 7   | Test Skeleton table-row variant                                                   | Gap      | MED    | 5min   | MED            | 5min  |
| 8   | Test Skeleton unknown variant (default case)                                      | Gap      | MED    | 5min   | MED            | 5min  |
| 9   | Test emptyStateAction button branch (no href)                                     | Gap      | MED    | 5min   | MED            | 5min  |
| 10  | Test Avatar fallback SVG (no Src, no Initials)                                    | Gap      | MED    | 5min   | MED            | 5min  |
| 11  | Test DefaultAccordionProps                                                        | Gap      | LOW    | 3min   | LOW            | 3min  |
| 12  | Test DefaultDropdownProps                                                         | Gap      | LOW    | 3min   | LOW            | 3min  |
| 13  | Test DefaultTableProps                                                            | Gap      | LOW    | 3min   | LOW            | 3min  |
| 14  | Test DefaultAlertProps                                                            | Gap      | LOW    | 3min   | LOW            | 3min  |
| 15  | Test DefaultStepIndicatorProps                                                    | Gap      | LOW    | 3min   | LOW            | 3min  |
| 16  | Test DefaultToastProps                                                            | Gap      | LOW    | 3min   | LOW            | 3min  |
| 17  | Test DefaultMinimalProps                                                          | Gap      | LOW    | 3min   | LOW            | 3min  |
| 18  | Test DefaultNavLinkProps                                                          | Gap      | LOW    | 3min   | LOW            | 3min  |
|     | **TIER 3 — JS UNIFICATION (TODO #23)**                                            |          |        |        |                |       |
| 19  | Audit all 3 JS patterns — document current state and target design                | TODO #23 | HIGH   | 10min  | MED            | 10min |
| 20  | Refactor Accordion JS to use data-attribute + window flag pattern                 | TODO #23 | MED    | 10min  | MED            | 10min |
| 21  | Refactor Dropdown JS to use data-attribute + window flag pattern                  | TODO #23 | MED    | 10min  | MED            | 10min |
| 22  | Refactor Modal JS to use data-attribute + window flag pattern                     | TODO #23 | MED    | 10min  | MED            | 10min |
| 23  | Extract shared JS attachment helper function                                      | TODO #23 | MED    | 8min   | MED            | 8min  |
| 24  | Update all JS-using tests for new pattern                                         | TODO #23 | MED    | 10min  | MED            | 10min |
| 25  | Verify all interactive components work (Accordion, Dropdown, Modal, Alert, Toast) | TODO #23 | HIGH   | 10min  | HIGH           | 10min |
|     | **TIER 4 — COVERAGE PUSH (80% impact)**                                           |          |        |        |                |       |
| 26  | Display: test Badge with icon + text combination                                  | Coverage | MED    | 5min   | LOW            | 5min  |
| 27  | Display: test EmptyState with custom icon                                         | Coverage | MED    | 5min   | LOW            | 5min  |
| 28  | Display: test EmptyState with all icon styles                                     | Coverage | MED    | 5min   | LOW            | 5min  |
| 29  | Display: test Tabs with no tabs provided                                          | Coverage | MED    | 5min   | LOW            | 5min  |
| 30  | Display: test Tabs underline variant                                              | Coverage | MED    | 5min   | LOW            | 5min  |
| 31  | Display: test Dropdown with divider items                                         | Coverage | MED    | 5min   | LOW            | 5min  |
| 32  | Display: test Tooltip all 4 positions                                             | Coverage | MED    | 8min   | LOW            | 8min  |
| 33  | Feedback: test Alert with custom icon                                             | Coverage | MED    | 5min   | LOW            | 5min  |
| 34  | Feedback: test Alert all 4 types                                                  | Coverage | MED    | 5min   | LOW            | 5min  |
| 35  | Feedback: test InlineMessage all types                                            | Coverage | MED    | 8min   | LOW            | 8min  |
| 36  | Feedback: test Spinner all sizes                                                  | Coverage | LOW    | 5min   | LOW            | 5min  |
| 37  | Feedback: test LoadingOverlay with custom spinner                                 | Coverage | LOW    | 5min   | LOW            | 5min  |
| 38  | Feedback: test ProgressBar with label and animation                               | Coverage | MED    | 5min   | LOW            | 5min  |
| 39  | Feedback: test StepIndicator with mixed states                                    | Coverage | MED    | 5min   | LOW            | 5min  |
| 40  | Forms: test Input with placeholder                                                | Coverage | LOW    | 3min   | LOW            | 3min  |
| 41  | Forms: test Input with addon (left/right)                                         | Coverage | MED    | 8min   | LOW            | 8min  |
| 42  | Forms: test Checkbox with help + error simultaneously                             | Coverage | MED    | 5min   | LOW            | 5min  |
| 43  | Forms: test Select with multiple options selected                                 | Coverage | LOW    | 5min   | LOW            | 5min  |
| 44  | Navigation: test NavLink active with children                                     | Coverage | MED    | 5min   | LOW            | 5min  |
| 45  | Navigation: test Pagination edge cases (1 page, 2 pages, 100+ pages)              | Coverage | MED    | 10min  | LOW            | 10min |
| 46  | Navigation: test MobileMenu with open/closed state                                | Coverage | MED    | 5min   | LOW            | 5min  |
| 47  | Icons: test strokeIcon path rendering                                             | Coverage | LOW    | 5min   | LOW            | 5min  |
| 48  | Layout: test Base with all optional features enabled                              | Coverage | MED    | 8min   | LOW            | 8min  |
| 49  | Layout: test Minimal with custom attrs                                            | Coverage | LOW    | 5min   | LOW            | 5min  |
|     | **TIER 5 — GOLDEN FILES (TODO #51)**                                              |          |        |        |                |       |
| 50  | Design golden file test infrastructure                                            | TODO #51 | MED    | 10min  | MED            | 10min |
| 51  | Implement golden file helper (write/read/compare)                                 | TODO #51 | MED    | 10min  | MED            | 10min |
| 52  | Convert display snapshot tests to golden files                                    | TODO #51 | MED    | 10min  | MED            | 10min |
| 53  | Convert feedback snapshot tests to golden files                                   | TODO #51 | MED    | 10min  | MED            | 10min |
| 54  | Convert forms snapshot tests to golden files                                      | TODO #51 | MED    | 8min   | MED            | 8min  |
| 55  | Convert navigation snapshot tests to golden files                                 | TODO #51 | MED    | 8min   | MED            | 8min  |
| 56  | Add CI golden file update command                                                 | TODO #51 | LOW    | 5min   | MED            | 5min  |
|     | **TIER 6 — DOCUMENTATION SITE (TODO #71)**                                        |          |        |        |                |       |
| 57  | Evaluate documentation site generators (pkgsite, doc2go, custom)                  | TODO #71 | MED    | 10min  | HIGH           | 10min |
| 58  | Set up chosen doc generator with CI integration                                   | TODO #71 | MED    | 10min  | HIGH           | 10min |
| 59  | Write per-component usage examples with live preview support                      | TODO #71 | HIGH   | 12min  | HIGH           | 12min |
| 60  | Configure GitHub Pages deployment                                                 | TODO #71 | MED    | 8min   | MED            | 8min  |
|     | **TIER 7 — CLEANUP & POLISH**                                                     |          |        |        |                |       |
| 61  | Move display/a11y_test.go ProgressBar test to feedback/ (TODO #59)                | TODO #59 | LOW    | 5min   | LOW            | 5min  |
| 62  | Move test helpers from utils/ to internal/testutils/ (TODO #58)                   | TODO #58 | MED    | 10min  | LOW            | 10min |
| 63  | Update all imports after test helper move                                         | TODO #58 | MED    | 10min  | LOW            | 10min |
| 64  | Verify examples/demo builds and renders correctly                                 | Cleanup  | LOW    | 5min   | MED            | 5min  |
| 65  | Update FEATURES.md with final component list                                      | Cleanup  | MED    | 8min   | MED            | 8min  |
| 66  | Update AGENTS.md with session 9 outcomes                                          | Cleanup  | LOW    | 5min   | LOW            | 5min  |
| 67  | Update TODO_LIST.md — mark completed items                                        | Cleanup  | LOW    | 5min   | LOW            | 5min  |
| 68  | Final verification: build + test + lint + coverage                                | Cleanup  | HIGH   | 5min   | HIGH           | 5min  |
| 69  | Commit and push all session 9 work                                                | Cleanup  | HIGH   | 5min   | HIGH           | 5min  |

---

## Execution Order (by tier)

### Phase 1: Release v0.2.0 (Tasks 1-4) — ~21min

Get the release out. Everything else can be post-release.

### Phase 2: Critical Fixes (Tasks 5-10) — ~33min

Fix the card_test compile errors and close the 4 biggest coverage gaps.

### Phase 3: Default Props Coverage (Tasks 11-18) — ~24min

Quick wins — 8 Default\*Props() functions at 0%. Just call them and assert non-nil.

### Phase 4: JS Unification (Tasks 19-25) — ~68min

TODO #23 — the last architectural inconsistency. Standardize all JS attachment patterns.

### Phase 5: Coverage Push (Tasks 26-49) — ~157min

Systematic coverage improvement across all packages. Target: 75%+ total.

### Phase 6: Golden Files (Tasks 50-56) — ~61min

TODO #51 — convert snapshot tests to golden file comparison.

### Phase 7: Documentation Site (Tasks 57-60) — ~40min

TODO #71 — set up auto-generated documentation.

### Phase 8: Cleanup (Tasks 61-69) — ~53min

TODO #58, #59, and final housekeeping.

---

## D2 Execution Graph

```d2
title: Session 9 Execution Plan — templ-components

direction: right

phase1: Release v0.2.0 {
  shape: rectangle
  style.fill: "#4CAF50"
  t1: Verify goreleaser
  t2: Update CHANGELOG
  t3: Verify README
  t4: Tag v0.2.0
  t1 -> t2 -> t3 -> t4
}

phase2: Critical Fixes {
  shape: rectangle
  style.fill: "#FF9800"
  t5: Fix card_test compile
  t6: Skeleton image
  t7: Skeleton table-row
  t8: Skeleton default
  t9: emptyState button
  t10: Avatar fallback
  t5 -> t6 -> t7 -> t8 -> t9 -> t10
}

phase3: Default Props {
  shape: rectangle
  style.fill: "#2196F3"
  t11: DefaultAccordion
  t12: DefaultDropdown
  t13: DefaultTable
  t14: DefaultAlert
  t15: DefaultStepIndicator
  t16: DefaultToast
  t17: DefaultMinimal
  t18: DefaultNavLink
}

phase4: JS Unification {
  shape: rectangle
  style.fill: "#9C27B0"
  t19: Audit JS patterns
  t20: Refactor Accordion JS
  t21: Refactor Dropdown JS
  t22: Refactor Modal JS
  t23: Extract shared helper
  t24: Update tests
  t25: Verify interactive
  t19 -> t20 -> t21 -> t22 -> t23 -> t24 -> t25
}

phase5: Coverage Push {
  shape: rectangle
  style.fill: "#00BCD4"
  t26_32: Display tests (7 tasks)
  t33_39: Feedback tests (7 tasks)
  t40_43: Forms tests (4 tasks)
  t44_46: Navigation tests (3 tasks)
  t47: Icons tests
  t48_49: Layout tests
}

phase6: Golden Files {
  shape: rectangle
  style.fill: "#FF5722"
  t50: Design infra
  t51: Implement helper
  t52: Display golden
  t53: Feedback golden
  t54: Forms golden
  t55: Nav golden
  t56: CI update cmd
  t50 -> t51 -> t52 -> t53 -> t54 -> t55 -> t56
}

phase7: Docs Site {
  shape: rectangle
  style.fill: "#607D8B"
  t57: Evaluate generators
  t58: Set up generator
  t59: Write examples
  t60: GitHub Pages
  t57 -> t58 -> t59 -> t60
}

phase8: Cleanup {
  shape: rectangle
  style.fill: "#795548"
  t61: Move ProgressBar test
  t62_63: Move test helpers
  t64: Verify demo
  t65: Update FEATURES
  t66_68: Update docs & verify
  t69: Commit & push
}

phase1 -> phase2 -> phase3 -> phase4 -> phase5 -> phase6 -> phase7 -> phase8
```

---

## Coverage Target by Package

| Package      | Current   | After Tier 2-3 | After Tier 5 | Target   |
| ------------ | --------- | -------------- | ------------ | -------- |
| display      | 68.9%     | ~72%           | ~76%         | 75%+     |
| feedback     | 73.6%     | ~76%           | ~80%         | 78%+     |
| forms        | 70.3%     | ~70%           | ~75%         | 75%+     |
| htmx         | 77.3%     | ~77%           | ~80%         | 78%+     |
| icons        | 75.0%     | ~75%           | ~78%         | 78%+     |
| internal/svg | 79.0%     | ~79%           | ~82%         | 80%+     |
| layout       | 72.9%     | ~73%           | ~78%         | 78%+     |
| navigation   | 72.1%     | ~72%           | ~77%         | 75%+     |
| utils        | 89.5%     | ~89%           | ~92%         | 90%+     |
| **Total**    | **70.7%** | **~72%**       | **~77%**     | **75%+** |

---

## Risk Assessment

| Risk                                       | Probability | Impact | Mitigation                                             |
| ------------------------------------------ | ----------- | ------ | ------------------------------------------------------ |
| JS unification breaks interactive behavior | MEDIUM      | HIGH   | Test each component individually before moving to next |
| Golden files produce flaky tests           | LOW         | MEDIUM | Use deterministic rendering, normalize whitespace      |
| Test helper move breaks consumers          | HIGH        | HIGH   | Defer to v1.0 (already planned)                        |
| goreleaser config is stale                 | LOW         | MEDIUM | Dry-run first                                          |
| Coverage push hits diminishing returns     | MEDIUM      | LOW    | Stop at 75%, don't chase 90%                           |

---

## Session Time Estimate

| Phase                   | Tasks  | Time        |
| ----------------------- | ------ | ----------- |
| Phase 1: Release        | 4      | 21min       |
| Phase 2: Critical Fixes | 6      | 33min       |
| Phase 3: Default Props  | 8      | 24min       |
| Phase 4: JS Unification | 7      | 68min       |
| Phase 5: Coverage Push  | 24     | 157min      |
| Phase 6: Golden Files   | 7      | 61min       |
| Phase 7: Docs Site      | 4      | 40min       |
| Phase 8: Cleanup        | 9      | 53min       |
| **Total**               | **69** | **~457min** |

---

## Decision Points

1. **Phase 1 → 2:** User confirms v0.2.0 tag
2. **Phase 4 scope:** Full JS unification vs. just documenting the 3 patterns
3. **Phase 6 scope:** Full golden file conversion vs. infrastructure only
4. **Phase 7 approach:** Custom docs site vs. pkgsite (may defer entirely)
5. **Phase 8, TODO #58:** Defer test helper move to v1.0 (breaking change)
