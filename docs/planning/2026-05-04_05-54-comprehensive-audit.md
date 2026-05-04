# Planning — 2026-05-04 Comprehensive Audit

## Pareto Analysis

### 1% → 51% Impact

| #   | Task                                                         | Impact                               | Effort |
| --- | ------------------------------------------------------------ | ------------------------------------ | ------ |
| 1   | Fix feedback split brain (unify alertStyleSet/toastStyleSet) | High — 60 lines, 2 modules deepened  | 30min  |
| 2   | Deepen icons (path-data map)                                 | High — 100 lines removed, extensible | 45min  |
| 3   | Add missing render tests (breadcrumbs, nav, mobile_menu)     | High — 3 uncovered components        | 60min  |

### 4% → 64% Impact

| #   | Task                       | Impact                           | Effort |
| --- | -------------------------- | -------------------------------- | ------ |
| 4   | Shared form error helper   | Medium — 30 lines deduped        | 30min  |
| 5   | AvatarStatus enum          | Medium — type safety             | 20min  |
| 6   | internal/svg tests         | Medium — foundation coverage     | 15min  |
| 7   | a11y validation tests      | Medium — accessibility guarantee | 45min  |
| 8   | Golden file snapshot tests | Medium — test maintainability    | 60min  |

### 20% → 80% Impact

| #   | Task                              | Impact | Effort |
| --- | --------------------------------- | ------ | ------ |
| 9   | TrendDirection enum for StatCard  | Low    | 15min  |
| 10  | HTMXSRI bool fix                  | Low    | 15min  |
| 11  | ProgressBar float precision       | Low    | 15min  |
| 12  | TableCell.Content templ.Component | Low    | 20min  |
| 13  | MapEnum direct test               | Low    | 10min  |
| 14  | Default\*Props constructor tests  | Low    | 20min  |
| 15  | Benchmarks for hot paths          | Low    | 30min  |
| 16  | CHANGELOG.md update               | Low    | 15min  |

## Execution Order

1. **Fix feedback split brain** (Architecture #6)
2. **Deepen icons** (Architecture #11)
3. **Add missing render tests** (Testing #22-25)
4. **Shared form error helper** (Architecture #12)
5. **AvatarStatus enum** (Architecture #13)
6. **internal/svg tests** (Testing #31)
7. **a11y tests** (Testing #27)
8. **Golden file tests** (Testing #26)
9. Then remaining P2/P3 items

## D2 Execution Graph

See: `docs/architecture-understanding/2026-05-04_05-54-current-state.d2`
See: `docs/architecture-understanding/2026-05-04_05-54-target-state-improved.d2`
