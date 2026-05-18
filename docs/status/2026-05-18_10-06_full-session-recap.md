# Comprehensive Status Report — templ-components + go-error-family

**Date:** 2026-05-18 10:06  
**Author:** Crush (automated)  
**Session context:** Continuation from 24-commit session on 2026-05-17

---

## Executive Summary

**templ-components** is a public Go UI component library (Go 1.26, templ, Tailwind v4) in good shape. Build green, lint green, 151 tests passing, 69.1% coverage. The previous session fixed real bugs (badge type mapping, WCAG ErrorAttrs, accessibility), added meaningful defaults, and reconciled massively stale documentation.

**go-error-family** is a public structured error protocol library, clean and stable at 76.2% coverage, 165 tests passing.

Both repos are **clean, pushed, and green**. No uncommitted changes.

---

## A) Fully Done ✅

### templ-components (54 of 75 TODO items complete — 72%)

| Area           | What                         | Details                                                                                                           |
| -------------- | ---------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| Build          | All builds pass              | `go build ./...` — zero errors                                                                                    |
| Tests          | 151 tests passing            | 9 packages, all green                                                                                             |
| Lint           | 0 issues on library packages | `golangci-lint run ./...` — clean                                                                                 |
| Coverage       | 69.1% overall                | utils 89.5%, internal/svg 79.0%, htmx 77.3%                                                                       |
| CI             | GitHub Actions               | Go 1.26, lint+build+test                                                                                          |
| Critical bugs  | 4 P0 fixes (#5-8)            | Literal string rendering, version mismatch, stale docs                                                            |
| Security       | 3 fixes (#9-12)              | NavLink shadow, ID validation, XSS escaping, accordion coupling                                                   |
| Architecture   | 10 refactorings (#13-22)     | Shared SVG, feedbackStyleSet, MapEnum, BaseProps everywhere, icon path map                                        |
| Accessibility  | 9 items (#30-41)             | alt text, aria-required, html lang, table scope, aria-live, focus trap, keyboard nav, ARIA linkage                |
| Testing        | 8 items (#42-44, 52-54)      | BDD tests for nav/htmx/layout/icons, a11y/dark mode/benchmark tests                                               |
| Dead code      | 2 items (#55, 60)            | IconAttrs removed, demo syntax fixed                                                                              |
| Documentation  | 6 items (#65-70)             | FEATURES.md, TODO_LIST.md, CONTEXT.md, CHANGELOG, migration guide, demo                                           |
| Session 5 work | 8 real fixes                 | BadgeNeutral bug, ErrorAttrs WCAG, aria-live on InlineLoading, 4 DefaultProps, tooltip cache, SimpleCard refactor |

### go-error-family

| Area              | What                         | Details                                  |
| ----------------- | ---------------------------- | ---------------------------------------- |
| Core              | Family, Tone, Audience types | Clean, well-tested                       |
| Coverage          | 76.2% overall                | root 97.1%, agent 100.0%, diagnose 60.6% |
| Dead code removal | Recent cleanup               | Unified rendering, removed dead code     |
| CI                | GitHub Actions               | Tag-based release workflow               |
| Docs              | Domain language template     | DDD ubiquitous language                  |

---

## B) Partially Done 🔨

| Item                         | Status                          | What's Left                                                                                                                                    |
| ---------------------------- | ------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| **#57 DefaultXxxProps**      | 6 of 8 have meaningful defaults | `DefaultAccordionProps()` and `DefaultStatCardProps()` still return zero-value structs                                                         |
| **#56 svg.FillIcon**         | Used by 4 components            | Exists as a thin proxy through `display/helpers.templ:fillIcon()`. Not dead code, but the TODO was about auditing if the indirection is needed |
| **Coverage target**          | 69.1% overall                   | display 66.0%, forms 68.1%, icons 68.3% — all below 70%                                                                                        |
| **go-error-family diagnose** | 60.6% coverage                  | Lowest package, needs more tests                                                                                                               |

---

## C) Not Started ⬜

### templ-components — 8 genuinely open items

| #   | Task                                                   | Priority | Effort | Impact            |
| --- | ------------------------------------------------------ | -------- | ------ | ----------------- |
| 23  | Unify JS attachment pattern (Accordion/Dropdown/Modal) | P2       | M      | Code consistency  |
| 24  | Extract shared dismiss JS (Alert/Toast)                | P2       | S      | Dedup             |
| 25  | Make toast icon SVG paths single-source                | P2       | S      | Dedup             |
| 35  | Add `aria-live` to HTMX error handling                 | P2       | S      | A11y              |
| 46  | Tests for Table mismatched header/row lengths          | P2       | S      | Robustness        |
| 47  | Tests for Modal/Dropdown with empty ID (panic)         | P2       | S      | Robustness        |
| 58  | Move test helpers to `internal/testutil/`              | P3       | S      | Code organization |
| 62  | Release automation (goreleaser)                        | P3       | M      | DevOps            |

### go-error-family — integration not started

Cross-project integration (go-error-family ↔ templ-components) is **researched but not started**. 3 architectural questions block all code work (see Section G).

---

## D) Totally Fucked Up 💥

### Stale TODOs discovered (6 items marked ⬜ but already done)

| #        | Task                                    | Reality                                                                               |
| -------- | --------------------------------------- | ------------------------------------------------------------------------------------- |
| **45**   | Add BDD tests for icons                 | `icons/bdd_test.go` exists with 5 test functions (47 subtests including all 42 icons) |
| **48**   | mapStatusToBadgeType boundary tests     | `display/helpers_test.go` has case-insensitive tests (Active, ERROR, In_Progress)     |
| **50**   | Improve utils coverage (56→75%)         | utils is at **89.5%** — well above the 75% target                                     |
| **63**   | Fix pre-commit hook to be executable    | Already executable: `-rwx--x--x` permissions                                          |
| **64**   | Exclude examples/ from lint (23 issues) | **0 issues** now — the 23 issues were apparently fixed                                |
| **#26b** | Extract cardShellClass                  | Already extracted as `const cardShellClass` in `card_templ.go:13`, used 3×            |

### Doc debt from previous session

The TODO_LIST.md update on 2026-05-17 was incomplete — these 6 items should have been marked ✅. The list claims 54/75 done, but it's actually **60/75 (80%)** once staleness is corrected.

### Pattern: TODO list chronically stale

Across multiple sessions, ~40-50% of items marked ⬜ turn out to be already done. The TODO list was last accurately verified 2026-05-07, but code kept changing. **Root cause: no automated TODO verification.**

---

## E) What We Should Improve 🔧

### High Impact

1. **Automated TODO staleness check** — Write a script that verifies each TODO against actual code state. Run it in CI or before status reports. The manual verification process wastes 30+ minutes every session.

2. **Coverage floors per package** — Set minimum coverage targets (e.g., 70% per package). Current gaps: display 66.0%, icons 68.3%, forms 68.1%. These drag the average down.

3. **JS pattern consolidation (#23)** — Three different JS attachment patterns across components is the biggest architectural inconsistency. Standardize before adding more interactive components.

4. **Golden file testing (#51)** — Current substring assertions are fragile. Golden files would catch regressions in full HTML output and make tests more maintainable.

### Medium Impact

5. **`internal/testutil/` extraction (#58)** — `utils/test_helpers.go` exports testing utilities used by all packages. This pollutes the public API. Move to internal.

6. **Table input validation (#46)** — No validation on header/row length mismatches. Silent rendering bugs possible.

7. **Empty ID panic tests (#47)** — `validateModalID()` and `validateDropdownID()` panic on empty ID, but no tests verify this.

8. **go-error-family diagnose coverage** — At 60.6%, this is the weakest package across both projects.

### Process Improvements

9. **Update TODO_LIST.md after every commit** — Not "at end of session." The 20+ stale items from this session prove batch updates don't work.

10. **Semver discipline** — Breaking changes (SimpleCard, BadgeDefault removal) pushed without version bumps. In 0.x this is technically allowed, but consumers have no migration signal.

---

## F) Top 25 Things To Do Next

Ranked by impact × effort ratio (Pareto order):

| Rank | Task                                                    | Package  | Effort | Impact | Why                                   |
| ---- | ------------------------------------------------------- | -------- | ------ | ------ | ------------------------------------- |
| 1    | **Mark 6 stale TODOs as ✅** (#45, 48, 50, 63, 64, 26b) | docs     | XS     | H      | Corrects project status from 54→60/75 |
| 2    | **Add aria-live to HTMX error handling** (#35)          | htmx     | S      | H      | Last open a11y bug                    |
| 3    | **Raise display coverage to 70%+**                      | display  | M      | H      | Lowest package at 66.0%               |
| 4    | **Unify JS attachment pattern** (#23)                   | display  | M      | H      | Biggest architectural debt            |
| 5    | **Raise icons coverage to 70%+**                        | icons    | S      | M      | 68.3% — close to threshold            |
| 6    | **Raise forms coverage to 70%+**                        | forms    | S      | M      | 68.1% — close to threshold            |
| 7    | **Extract shared dismiss JS** (#24)                     | feedback | S      | M      | Nearly identical code in Alert/Toast  |
| 8    | **Single-source toast icon paths** (#25)                | feedback | S      | M      | Duplicated in Go and JS               |
| 9    | **Table header/row length validation** (#46)            | display  | S      | M      | Silent rendering bugs possible        |
| 10   | **Empty ID panic tests** (#47)                          | display  | S      | M      | Missing test coverage for guards      |
| 11   | **Accordion meaningful defaults** (#57 partial)         | display  | XS     | L      | Zero-value struct with Items slice    |
| 12   | **StatCard meaningful defaults** (#57 partial)          | display  | XS     | M      | Zero-value struct                     |
| 13   | **Move test helpers to internal/testutil** (#58)        | utils    | S      | M      | Public API pollution                  |
| 14   | **Golden file tests** (#51)                             | all      | M      | H      | Regression safety                     |
| 15   | **Write TODO staleness checker script**                 | tooling  | M      | H      | Prevents 30min waste per session      |
| 16   | **Raise go-error-family diagnose coverage**             | diagnose | M      | M      | 60.6% → 75%+                          |
| 17   | **Release automation with goreleaser** (#62)            | DevOps   | M      | M      | Enables semver workflow               |
| 18   | **PageProps convention documentation** (#72)            | docs     | XS     | L      | Only Props struct without BaseProps   |
| 19   | **Document SimpleCard breaking change** (#72 related)   | docs     | XS     | M      | Missing from migration guide          |
| 20   | **ProgressBar test location** (#59)                     | display  | XS     | L      | Wrong package, tests feedback         |
| 21   | **Exclude examples from lint config** (#64)             | tooling  | XS     | L      | Already 0 issues but config cleaner   |
| 22   | **Documentation site generation** (#71)                 | docs     | L      | M      | Auto-generated from source            |
| 23   | **svg.FillIcon proxy audit** (#56)                      | display  | XS     | S      | Verify indirection is needed          |
| 24   | **BDD test for go-error-family agent**                  | agent    | S      | M      | 100% coverage but no behavior tests   |
| 25   | **Add CONTRIBUTING.md**                                 | docs     | M      | M      | Public repo, no contribution guide    |

---

## G) Top #1 Question I Cannot Answer Myself

### How should go-error-family integrate with templ-components?

The research is done. I understand both codebases deeply. But **3 product decisions require your input** before any code:

1. **Dependency model** — Should templ-components import go-error-family directly (adds a dep), or should there be a separate bridge package (e.g., `templ-components/errors`), or should integration be example-only (no import, just pattern documentation)?

2. **HTTP error envelope** — go-error-family has no HTTP-specific layer. templ-components' HTMX error handling (`htmx/error_handling.templ`) currently renders raw errors. Should integration add structured HTTP error rendering (status code → tone → toast style mapping), or stay presentation-only?

3. **Scope** — Should integration cover just display (error → toast/alert/badge styling), or also include form field error mapping (go-error-family validation errors → `ErrorAttrs`), or go full-stack (HTMX request intercept → structured error parsing → auto-render)?

**My recommendation:** Start with option (1) bridge package + (3) display-only scope. This is the leanest path that delivers value without coupling. HTTP envelope can come later if needed.

---

## Metrics Dashboard

### templ-components

| Metric              | Value                                      |
| ------------------- | ------------------------------------------ |
| Packages            | 9 library + 1 examples                     |
| `.templ` files      | 31                                         |
| Lines of code       | ~7,684                                     |
| Tests               | 151 passing                                |
| Coverage            | 69.1% overall                              |
| Lint issues         | 0                                          |
| TODO progress       | 60/75 (80%) — after stale correction       |
| Commits (session 5) | 24                                         |
| Dependencies        | 0 runtime (templ + tailwind-merge-go only) |

### go-error-family

| Metric        | Value         |
| ------------- | ------------- |
| Packages      | 3             |
| Go files      | 17            |
| Lines of code | ~3,372        |
| Tests         | 165 passing   |
| Coverage      | 76.2% overall |
| Dependencies  | 0             |

### Coverage by Package (templ-components)

| Package       | Coverage | Status       |
| ------------- | -------- | ------------ |
| utils         | 89.5%    | ✅ Excellent |
| internal/svg  | 79.0%    | ✅ Good      |
| htmx          | 77.3%    | ✅ Good      |
| feedback      | 73.2%    | ✅ Good      |
| layout        | 72.9%    | ✅ Good      |
| navigation    | 72.1%    | ✅ Good      |
| icons         | 68.3%    | ⚠️ Below 70% |
| forms         | 68.1%    | ⚠️ Below 70% |
| display       | 66.0%    | ⚠️ Below 70% |
| examples/demo | 0.0%     | — (expected) |

---

## Git State

### templ-components

- **Branch:** master
- **Status:** clean, up to date with origin
- **Last commit:** `7a38a3f` — docs(todo): mark #34 aria-live and #36 ErrorAttrs as done, #37 dot scaling

### go-error-family

- **Branch:** master
- **Status:** clean, up to date with origin
- **Last commit:** `983a71b` — chore: add domain language template, gitignore \*.db, update status doc
