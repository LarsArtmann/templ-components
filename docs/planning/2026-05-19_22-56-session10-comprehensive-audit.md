# Session 10 — Comprehensive 9-Skill Audit & Execution Plan

**Date:** 2026-05-19 | **Scope:** Full codebase audit with 9 skills

## Pareto Analysis

### 1% → 51% Impact (Do First)

| #   | Task                                   | Effort | Impact                                |
| --- | -------------------------------------- | ------ | ------------------------------------- |
| 6   | Remove/validate unknown icon names     | 15min  | Eliminates silent visual bugs         |
| 8   | Delete deprecated `Exclamation` icon   | 10min  | Removes dead code, fixes split brain  |
| 20  | Clamp ProgressBar percent to [0, 100]  | 5min   | Prevents visual overflow              |
| 19  | Validate Pagination CurrentPage > 0    | 5min   | Prevents broken pagination            |
| 37  | Remove unused `badgeTextLive` constant | 2min   | Cleans linter warning                 |
| 38  | Delete `TestPtr` in utils_test.go      | 2min   | Removes dead test                     |
| 7   | Fix IconPathJS stroke-width mismatch   | 10min  | Eliminates icon rendering split brain |

### 4% → 64% Impact (Do Second)

| #   | Task                                   | Effort | Impact                         |
| --- | -------------------------------------- | ------ | ------------------------------ |
| 9   | Unify AlertType/ToastType              | 30min  | Eliminates type duplication    |
| 10  | Merge alertStyleMap/toastStyleMap      | 20min  | Reduces style map duplication  |
| 27  | Fix FillIcon variadic bool             | 10min  | API quality improvement        |
| 16  | Add BaseProps to StepIndicatorProps    | 15min  | API consistency                |
| 17  | Convert LoadingOverlay to props struct | 20min  | API consistency                |
| 13  | Use stable IDs in modal JS             | 20min  | Robustness improvement         |
| 5   | Fix demo app to use layout.Base        | 45min  | Dogfooding, public-facing      |
| 15  | Use icon system in Breadcrumbs         | 15min  | Eliminates raw SVG duplication |
| 44  | Document htmx→feedback JS dependency   | 5min   | Consumer documentation         |

### 20% → 80% Impact (Do Third)

| #   | Task                                           | Effort | Impact                      |
| --- | ---------------------------------------------- | ------ | --------------------------- |
| 11  | Make SimpleCard compose through Card           | 30min  | Reduces shell duplication   |
| 12  | Add ComponentProps interface                   | 30min  | Enables generic handling    |
| 26  | Replace DropdownItem empty-Href discrimination | 45min  | Type safety                 |
| 28  | Audit tailwind-merge-go thread safety          | 30min  | Performance (mutex removal) |
| 36  | Consolidate test files                         | 120min | Reduces 37 → ~15 test files |
| 21  | Consolidate inline JS into shared init         | 90min  | JS architecture overhaul    |
| 33  | Eliminate icon list split brain                | 30min  | Maintenance burden          |

## Execution Order (Prioritized)

```
Phase 1 (Quick Wins — ~50min):
  #37, #38, #20, #19, #7, #8, #6

Phase 2 (Architecture — ~3hr):
  #9, #10, #27, #16, #17, #13, #15, #44

Phase 3 (Refactoring — ~4hr):
  #5, #11, #12, #26, #28, #33

Phase 4 (Consolidation — ~4hr):
  #36, #21, #22, #23, #24, #25

Phase 5 (Polish — ~2hr):
  #14, #18, #29, #30, #31, #32, #34, #35, #39, #40, #41, #42, #43, #45
```

## D2 Execution Graph

```d2
direction: right

phase1: {
  shape: rectangle
  style.fill: "#C8E6C9"
  "Quick Wins"
  -> "#37: Remove badgeTextLive"
  -> "#38: Delete TestPtr"
  -> "#20: Clamp ProgressBar%"
  -> "#19: Validate Pagination page"
  -> "#7: Fix IconPathJS width"
  -> "#8: Delete Exclamation icon"
  -> "#6: Validate icon names"
}

phase2: {
  shape: rectangle
  style.fill: "#BBDEFB"
  "Architecture"
  -> "#9: Unify AlertType/ToastType"
  -> "#10: Merge style maps"
  -> "#27: Fix FillIcon param"
  -> "#16: Add BaseProps to StepIndicator"
  -> "#17: LoadingOverlay props"
  -> "#13: Modal stable IDs"
  -> "#15: Breadcrumbs use icons"
}

phase3: {
  shape: rectangle
  style.fill: "#FFE0B2"
  "Refactoring"
  -> "#5: Fix demo app"
  -> "#11: SimpleCard compose"
  -> "#12: ComponentProps interface"
  -> "#26: DropdownItem type"
  -> "#28: Audit mutex"
  -> "#33: Icon list auto-gen"
}

phase4: {
  shape: rectangle
  style.fill: "#E1BEE7"
  "Consolidation"
  -> "#36: Consolidate test files"
  -> "#21: Shared JS init"
  -> "#22-25: JS fixes"
}

phase5: {
  shape: rectangle
  style.fill: "#F5F5F5"
  "Polish"
  -> "#14,18,29-32,34,35,39-43,45"
}

phase1 -> phase2 -> phase3 -> phase4 -> phase5
```
