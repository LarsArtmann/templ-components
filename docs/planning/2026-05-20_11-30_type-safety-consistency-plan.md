# Session 12.5 — Execution Plan: Type Safety, Consistency, Quality

**Generated:** 2026-05-20 | **Scope:** Reflection-driven improvements
**Principle:** Make impossible states unrepresentable, eliminate inconsistency, validate everything

## Execution Order (sorted by impact × low effort first)

| # | Task | Impact | Effort | Type |
|---|------|--------|--------|------|
| 1 | NavLink + MobileNavLink ID propagation | 🟡 Bug | 10m | Consistency |
| 2 | Checkbox unconditional `id=""` → conditional | 🟡 Bug | 5m | Correctness |
| 3 | Modal panel `props.Class` → `utils.Class()` | 🟡 Bug | 5m | Tailwind merge |
| 4 | alertIconName switch → map lookup | 🟢 Consistency | 5m | Pattern alignment |
| 5 | toastIconName switch → map lookup | 🟢 Consistency | 5m | Pattern alignment |
| 6 | spinnerSizeClass switch → map lookup | 🟢 Consistency | 5m | Pattern alignment |
| 7 | progressHeightClass switch → map lookup | 🟢 Consistency | 5m | Pattern alignment |
| 8 | avatarSizeClass + avatarDotSizeClass → map lookups | 🟢 Consistency | 8m | Pattern alignment |
| 9 | InputType validation (prevent XSS) | 🔴 Security | 10m | Type safety |
| 10 | StatCard Trend → tagged switch (lint hint) | 🟢 Lint | 5m | Code quality |
| 11 | Add t.Parallel() to TestSecurityHeaders | 🟢 Quality | 2m | Test quality |
| 12 | Test coverage: Modal without Title | 🟡 Quality | 5m | Coverage |
| 13 | Test coverage: ProgressBar Total=0, negative | 🟡 Quality | 5m | Coverage |
| 14 | Test coverage: Pagination edge cases | 🟡 Quality | 5m | Coverage |
| 15 | Test coverage: Dropdown empty Items | 🟡 Quality | 5m | Coverage |
| 16 | Test coverage: Alert/Toast edge cases | 🟡 Quality | 8m | Coverage |
| 17 | Test coverage: StepIndicator edge cases | 🟡 Quality | 5m | Coverage |
| 18 | Test coverage: Nav empty Links | 🟡 Quality | 5m | Coverage |
| 19 | Update AGENTS.md + TODO_LIST.md | 🟢 Docs | 5m | Documentation |
| 20 | Final verification: build + test + lint + coverage | 🟢 Verify | 5m | Verification |
