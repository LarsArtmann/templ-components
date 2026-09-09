# Session 12.5 — Execution Plan: Type Safety, Consistency, Quality

**Generated:** 2026-05-20 | **Scope:** Reflection-driven improvements
**Principle:** Make impossible states unrepresentable, eliminate inconsistency, validate everything

## Execution Order (sorted by impact × low effort first)

| #  | Task                                               | Impact         | Effort | Type              |
| -- | -------------------------------------------------- | -------------- | ------ | ----------------- |
| ~~1~~  | ~~NavLink + MobileNavLink ID propagation~~ done — navigation/nav link.templ | ~~🟡 Bug~~ | ~~10m~~ | ~~Consistency~~ |
| ~~2~~  | ~~Checkbox unconditional `id=""` → conditional~~ done — CHANGELOG.md | ~~🟡 Bug~~ | ~~5m~~ | ~~Correctness~~ |
| ~~3~~  | ~~Modal panel `props.Class` → `utils.Class()`~~ done — display/modal.templ | ~~🟡 Bug~~ | ~~5m~~ | ~~Tailwind merge~~ |
| ~~4~~  | ~~alertIconName switch → map lookup~~ done — feedback/styles.go | ~~🟢 Consistency~~ | ~~5m~~ | ~~Pattern alignment~~ |
| ~~5~~  | ~~toastIconName switch → map lookup~~ done — feedback/styles.go | ~~🟢 Consistency~~ | ~~5m~~ | ~~Pattern alignment~~ |
| ~~6~~  | ~~spinnerSizeClass switch → map lookup~~ done — feedback/loading.templ | ~~🟢 Consistency~~ | ~~5m~~ | ~~Pattern alignment~~ |
| ~~7~~  | ~~progressHeightClass switch → map lookup~~ done — feedback/progressbar.templ | ~~🟢 Consistency~~ | ~~5m~~ | ~~Pattern alignment~~ |
| ~~8~~  | ~~avatarSizeClass + avatarDotSizeClass → map lookups~~ done — display/avatar.templ | ~~🟢 Consistency~~ | ~~8m~~ | ~~Pattern alignment~~ |
| ~~9~~  | ~~InputType validation (prevent XSS)~~ done — forms | ~~🔴 Security~~ | ~~10m~~ | ~~Type safety~~ |
| ~~10~~ | ~~StatCard Trend → tagged switch (lint hint)~~ done — display/card.templ | ~~🟢 Lint~~ | ~~5m~~ | ~~Code quality~~ |
| ~~11~~ | ~~Add t.Parallel() to TestSecurityHeaders~~ done — layout/a11y test.go | ~~🟢 Quality~~ | ~~2m~~ | ~~Test quality~~ |
| ~~12~~ | ~~Test coverage: Modal without Title~~ done — display/coverage test.go | ~~🟡 Quality~~ | ~~5m~~ | ~~Coverage~~ |
| ~~13~~ | ~~Test coverage: ProgressBar Total=0, negative~~ done — feedback/edge cases test.go | ~~🟡 Quality~~ | ~~5m~~ | ~~Coverage~~ |
| ~~14~~ | ~~Test coverage: Pagination edge cases~~ done — navigation/coverage test.go | ~~🟡 Quality~~ | ~~5m~~ | ~~Coverage~~ |
| ~~15~~ | ~~Test coverage: Dropdown empty Items~~ done — display/dropdown test.go | ~~🟡 Quality~~ | ~~5m~~ | ~~Coverage~~ |
| ~~16~~ | ~~Test coverage: Alert/Toast edge cases~~ done — feedback/edge cases test.go | ~~🟡 Quality~~ | ~~8m~~ | ~~Coverage~~ |
| ~~17~~ | ~~Test coverage: StepIndicator edge cases~~ done — feedback/snapshot test.go | ~~🟡 Quality~~ | ~~5m~~ | ~~Coverage~~ |
| ~~18~~ | ~~Test coverage: Nav empty Links~~ done — navigation/edge cases test.go | ~~🟡 Quality~~ | ~~5m~~ | ~~Coverage~~ |
| ~~19~~ | ~~Update AGENTS.md + TODO_LIST.md~~ done — AGENTS.md | ~~🟢 Docs~~ | ~~5m~~ | ~~Documentation~~ |
| ~~20~~ | ~~Final verification: build + test + lint + coverage~~ done — docs/status/2026-09-09 16-19 post-release-a11y-and-demo-e2e.md | ~~🟢 Verify~~ | ~~5m~~ | ~~Verification~~ |
