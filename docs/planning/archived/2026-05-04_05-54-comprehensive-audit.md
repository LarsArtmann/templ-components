# Planning — 2026-05-04 Comprehensive Audit

## Pareto Analysis

### 1% → 51% Impact

| #     | Task                                                                                       | Impact                                   | Effort    |
| ----- | ------------------------------------------------------------------------------------------ | ---------------------------------------- | --------- |
| ~~1~~ | ~~Fix feedback split brain (unify alertStyleSet/toastStyleSet)~~ done — feedback/styles.go | ~~High — 60 lines, 2 modules deepened~~  | ~~30min~~ |
| ~~2~~ | ~~Deepen icons (path-data map)~~ done — feedback/styles.go                                 | ~~High — 100 lines removed, extensible~~ | ~~45min~~ |
| ~~3~~ | ~~Add missing render tests (breadcrumbs, nav, mobile_menu)~~ done — feedback/styles.go     | ~~High — 3 uncovered components~~        | ~~60min~~ |

### 4% → 64% Impact

| #     | Task                                                     | Impact                               | Effort    |
| ----- | -------------------------------------------------------- | ------------------------------------ | --------- |
| ~~4~~ | ~~Shared form error helper~~ done — feedback/styles.go   | ~~Medium — 30 lines deduped~~        | ~~30min~~ |
| ~~5~~ | ~~AvatarStatus enum~~ done — feedback/styles.go          | ~~Medium — type safety~~             | ~~20min~~ |
| ~~6~~ | ~~internal/svg tests~~ done — feedback/styles.go         | ~~Medium — foundation coverage~~     | ~~15min~~ |
| ~~7~~ | ~~a11y validation tests~~ done — feedback/styles.go      | ~~Medium — accessibility guarantee~~ | ~~45min~~ |
| ~~8~~ | ~~Golden file snapshot tests~~ done — feedback/styles.go | ~~Medium — test maintainability~~    | ~~60min~~ |

### 20% → 80% Impact

| #      | Task                                                                             | Impact  | Effort    |
| ------ | -------------------------------------------------------------------------------- | ------- | --------- |
| ~~9~~  | ~~TrendDirection enum for StatCard~~ done — feedback/styles.go                   | ~~Low~~ | ~~15min~~ |
| ~~10~~ | ~~HTMXSRI bool fix~~ **Won't implement — superseded self hosted htmx adr-0007.** | ~~Low~~ | ~~15min~~ |
| ~~11~~ | ~~ProgressBar float precision~~ done — feedback/progressbar.templ                | ~~Low~~ | ~~15min~~ |
| ~~12~~ | ~~TableCell.Content templ.Component~~ done — display/table.templ                 | ~~Low~~ | ~~20min~~ |
| ~~13~~ | ~~MapEnum direct test~~ done — utils/utils test.go                               | ~~Low~~ | ~~10min~~ |
| ~~14~~ | ~~Default\*Props constructor tests~~ done — display/accordion test.go            | ~~Low~~ | ~~20min~~ |
| ~~15~~ | ~~Benchmarks for hot paths~~ done — display/benchmark test.go                    | ~~Low~~ | ~~30min~~ |
| ~~16~~ | ~~CHANGELOG.md update~~ done — CHANGELOG.md                                      | ~~Low~~ | ~~15min~~ |

## Execution Order

1. ~~**Fix feedback split brain** (Architecture #6)~~ done — feedback/styles.go
2. ~~**Deepen icons** (Architecture #11)~~ done — icons/icon paths.go
3. ~~**Add missing render tests** (Testing #22-25)~~ done — navigation/bdd test.go
4. ~~**Shared form error helper** (Architecture #12)~~ done — forms/aria.go
5. ~~**AvatarStatus enum** (Architecture #13)~~ done — display/avatar templ.go
6. ~~**internal/svg tests** (Testing #31)~~ done — utils/svg/svg test.go
7. ~~**a11y tests** (Testing #27)~~ done — forms/a11y test.go
8. ~~**Golden file tests** (Testing #26)~~ done — display/golden test.go
9. ~~Then remaining P2/P3 items~~ done — CHANGELOG.md

## D2 Execution Graph

See: `docs/architecture-understanding/2026-05-04_05-54-current-state.d2`
See: `docs/architecture-understanding/2026-05-04_05-54-target-state-improved.d2`
