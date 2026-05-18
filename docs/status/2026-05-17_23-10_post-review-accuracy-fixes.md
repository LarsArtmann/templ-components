# Comprehensive Status Report — 2026-05-17 23:10

**Projects:** templ-components + go-error-family
**Author:** Lars Artmann
**Date:** 2026-05-17 23:10 CEST

---

## Executive Summary

Two public Go libraries. **templ-components** (UI components, 68.3% coverage, build/lint green, 8 sessions of work) and **go-error-family** (error protocol, 76.2% coverage, v0.1.1 released). This session focused on a deep self-review that uncovered **massive doc staleness** (10+ items marked undone that were already done), one real latent bug (`mapStatusToBadgeType`), and multiple convention violations (SimpleCard). All discovered issues fixed, docs reconciled, both repos clean and pushed.

---

## Session Summary (2026-05-17)

### Commits This Session

| Commit                      | Type         | Description                                                             |
| --------------------------- | ------------ | ----------------------------------------------------------------------- |
| `062ee1a`                   | docs         | Comprehensive dual-project status report                                |
| `ed9010f`                   | docs         | Mark TODO #20-22, #26 as done (stale status fix)                        |
| `0910413`                   | docs         | Fix stale known issues in FEATURES.md                                   |
| `1fcf7ae`                   | **fix**      | `mapStatusToBadgeType("default")` → `BadgeNeutral` (was `BadgePrimary`) |
| `d17cd6a`                   | docs         | Mark P0/P1 bugs #9-12 as done (already fixed)                           |
| `da4e32f`                   | **refactor** | SimpleCard → `SimpleCardProps` with BaseProps + utils.Class()           |
| `0e801eb`                   | **refactor** | Tooltip: cache position lookup, remove redundant default                |
| `f12c089`                   | docs         | Mark demo app #70 as done                                               |
| `f07d721`                   | docs         | AGENTS.md: add SimpleCard breaking change                               |
| (go-error-family) `983a71b` | chore        | Add domain language template, gitignore \*.db, update status            |

### What Changed in Code

- **1 real bug fix**: `mapStatusToBadgeType("default")` returned `BadgePrimary` instead of `BadgeNeutral` — leftover from `BadgeDefault` removal
- **1 convention fix**: `SimpleCard()` now uses `SimpleCardProps` with `BaseProps` + `CardPadding` + `utils.Class()` (was violating 3 conventions)
- **1 performance fix**: Tooltip caches `tooltipLookupPosition()` result (was calling twice per render)
- **10 docs commits**: TODO_LIST, FEATURES, AGENTS — reconciled with actual code state

---

## A. FULLY DONE ✅

### templ-components — Build & Quality

| Item              | Evidence                       |
| ----------------- | ------------------------------ |
| Build passes      | `go build ./...` — zero errors |
| All tests pass    | 9 library packages, all green  |
| Lint clean        | 0 issues on library packages   |
| GitHub Actions CI | Go 1.26, lint+build+test       |

### templ-components — Components (53 total, 42 icons)

| Package          | Components                                                                                                                                            | Status               |
| ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------- |
| `display` (14)   | Accordion, Avatar, Badge, StatusBadge, Card, SimpleCard, StatCard, Dropdown, EmptyState, SimpleEmptyState, Modal, Table, Tabs, Tooltip                | All FULLY_FUNCTIONAL |
| `feedback` (12)  | Alert, InlineError, InlineSuccess, Spinner, LoadingOverlay, InlineLoading, Skeleton, SkeletonGroup, ProgressBar, StepIndicator, ToastContainer, Toast | All FULLY_FUNCTIONAL |
| `forms` (6)      | Input, Checkbox, Select, Textarea, Label, FieldError                                                                                                  | All FULLY_FUNCTIONAL |
| `htmx` (7)       | LoadingIndicator, InlineLoadingOverlay, LoadingButton, ConfirmDelete, SwapOOB, CSRFToken, GlobalErrorHandling                                         | All FULLY_FUNCTIONAL |
| `icons` (1)      | Icon (42 named icons)                                                                                                                                 | FULLY_FUNCTIONAL     |
| `layout` (4)     | Base, Minimal, ThemeScript, ThemeToggle                                                                                                               | All FULLY_FUNCTIONAL |
| `navigation` (9) | Nav, SimpleNav, NavLink, MobileNavLink, Breadcrumbs, Pagination, MobileMenu, MobileMenuToggle, Footer                                                 | All FULLY_FUNCTIONAL |

### templ-components — Architecture Completed

| #      | Item                                                                                                     | Status   |
| ------ | -------------------------------------------------------------------------------------------------------- | -------- |
| 1-8    | All P0 bugs                                                                                              | ✅ Fixed |
| 9-12   | P0/P1 bugs (NavLink, Dropdown XSS, Accordion, ID validation)                                             | ✅ Fixed |
| 13-19  | Core architecture (SVG helpers, feedbackStyleSet, MapEnum, BaseProps, style maps, PageProps, icon paths) | ✅ Done  |
| 20-22  | Badge consolidation, BadgeDefault merge, ActiveTabID                                                     | ✅ Done  |
| 26     | HTMX loading decoupled                                                                                   | ✅ Done  |
| 26a    | Tooltip double lookup fixed                                                                              | ✅ Done  |
| 27-29  | AvatarStatus enum, TrendDirection, HTMXUseSRI                                                            | ✅ Done  |
| 38-41  | Modal focus trap, Dropdown keyboard, Tabs ARIA, Tooltip linkage                                          | ✅ Done  |
| 52-54  | a11y tests, dark mode tests, benchmarks                                                                  | ✅ Done  |
| 60, 70 | Demo app builds                                                                                          | ✅ Done  |

### go-error-family

| Item                                                           | Evidence                                                      |
| -------------------------------------------------------------- | ------------------------------------------------------------- |
| Root package 97.1% coverage                                    | 165+ tests                                                    |
| Agent package 100% coverage                                    | All paths tested                                              |
| v0.1.1 released                                                | MIT license, GitHub release workflow                          |
| 5 Families with data-driven metadata                           | `familyData` array                                            |
| Consumer interfaces (Coded, Classified, Contextual, Retryable) | Each embeds `error` for `errors.AsType`                       |
| CLI boundary handler                                           | `HandleError`, `HandleErrorDetailed`, `HandleErrorWithConfig` |
| Template resolution system                                     | Override → Registered → Default → Family fallback             |
| Diagnostic rules (4)                                           | Postgres, Filesystem, Network, Git                            |
| AI debug agent                                                 | Root cause analysis + FixSteps                                |
| Zero dependencies                                              | Pure stdlib                                                   |
| Domain language template                                       | `docs/DOMAIN_LANGUAGE.md`                                     |

---

## B. PARTIALLY DONE 🔨

### templ-components

| Item                | What's Left                                                         |
| ------------------- | ------------------------------------------------------------------- |
| Test coverage 68.3% | Target 75%+. `utils` 55.3%, `forms` 63.8%, `display` 66.1%          |
| Documentation       | README/FEATURES/CONTEXT exist. No docs site, no auto-generation     |
| Card shell CSS      | `cardShellClass` constant exists but StatCard still uses raw concat |

### go-error-family

| Item                      | What's Left                                                          |
| ------------------------- | -------------------------------------------------------------------- |
| diagnose coverage 60.6%   | Rules that shell out are integration-test territory                  |
| No JSON serialization     | `MarshalJSON`/`UnmarshalJSON` not implemented — blocks HTTP envelope |
| No HTTP response envelope | Would bridge to templ-components integration                         |

---

## C. NOT STARTED ⬜

### templ-components — Remaining TODO Items (25 items)

| #            | Task                                                        | Priority |
| ------------ | ----------------------------------------------------------- | -------- |
| 23           | Unify JS attachment pattern across Accordion/Dropdown/Modal | P2       |
| 24           | Extract shared dismiss JS for Alert and Toast               | P2       |
| 25           | Make toast icon SVG paths single-source                     | P2       |
| 26b          | Extract card shell CSS into consistent usage                | P3       |
| 30           | Add `alt` text to Avatar `<img>`                            | P1       |
| 31           | Add `aria-required` to form inputs                          | P1       |
| 32           | Add `<html lang>` to Base layout                            | P1       |
| 33           | Add Table header `scope` attributes                         | P2       |
| 34           | Add `aria-live` to HTMX loading indicators                  | P2       |
| 35           | Add `aria-live` to HTMX error handling                      | P2       |
| 36           | Fix ErrorAttrs for simultaneous error + help text           | P2       |
| 37           | Scale avatar status dot with avatar size                    | P3       |
| 42           | BDD tests for navigation package                            | P1       |
| 43           | BDD tests for htmx package                                  | P1       |
| 44           | BDD tests for layout package                                | P1       |
| 45           | BDD tests for icons package                                 | P2       |
| 46           | Tests for Table mismatched header/row lengths               | P2       |
| 47           | Tests for Modal/Dropdown with empty ID                      | P2       |
| 48           | Tests for mapStatusToBadgeType boundary cases               | P2       |
| 49           | Improve forms test coverage (64% → 75%+)                    | P2       |
| 50           | Improve utils test coverage (55% → 75%+)                    | P2       |
| 51           | Convert snapshot tests to golden file comparison            | P2       |
| 55           | Remove or use `icons.IconAttrs` (dead code)                 | P2       |
| 56           | Remove or use `internal/svg.FillIcon`                       | P2       |
| 57-64, 71-72 | DevOps, cleanup, docs (P3)                                  | P3       |

### go-error-family

| Task                                                    | Priority |
| ------------------------------------------------------- | -------- |
| diagnose coverage improvement (60.6% → 75%+)            | P2       |
| Add `MarshalJSON`/`UnmarshalJSON` for API serialization | P2       |
| Define HTTP response envelope                           | P2       |
| More diagnostic rules (beyond 4)                        | P3       |

### Cross-Project

| Task                                           | Status                                     |
| ---------------------------------------------- | ------------------------------------------ |
| go-error-family → templ-components integration | Research done, 3 questions open, zero code |

---

## D. TOTALLY FUCKED UP ❌

**Nothing actively broken.** Both projects build, test, and lint clean.

### Honest Assessment of Weaknesses

| Issue                                                  | Severity         | Details                                                                                                                               |
| ------------------------------------------------------ | ---------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| `utils` at 55.3% coverage                              | **Shameful**     | Core package, lowest coverage. `MergeAttrs`, `CurrentYear`, `Deref` undertested.                                                      |
| `forms` at 63.8% coverage                              | Low              | Second lowest package. Missing edge case tests.                                                                                       |
| Docs were 10+ items stale                              | **Embarrassing** | TODO_LIST had 8 items marked ⬜ that were ✅. FEATURES.md referenced deleted types. This means docs were wrong for multiple sessions. |
| 10 zero-value `DefaultXxxProps()`                      | Low              | Violate the stated convention "for every component with non-zero defaults".                                                           |
| `mapStatusToBadgeType` case-sensitive                  | Medium           | `"Active"`, `"SUCCESS"` fall through to `BadgeNeutral`. TODO #48 open.                                                                |
| No radio, file, toggle components                      | Low              | Forms package incomplete.                                                                                                             |
| `icons.IconAttrs` dead code                            | Low              | Exported, never called.                                                                                                               |
| go-error-family has no JSON serialization              | Medium           | Blocks the most valuable integration path (HTTP API responses).                                                                       |
| `SimpleCard` was a **breaking change** in this session | Medium           | Changed `SimpleCard()` → `SimpleCard(SimpleCardProps)`. No version bump.                                                              |

---

## E. WHAT WE SHOULD IMPROVE

### templ-components

1. **Raise test coverage to 75%+** — `utils` (55%) is embarrassing for a core package. Write targeted tests for `MergeAttrs`, `Deref`, `DerefOr`, `CurrentYear`, `BoolString`.
2. **Stop writing stale docs** — Every code change must update TODO_LIST immediately. The 10-item gap was pure negligence.
3. **Audit all `DefaultXxxProps()`** — 10 functions return zero-value structs. Either set meaningful defaults or remove them.
4. **Add `mapStatusToBadgeType` case-insensitive matching** — Current case sensitivity is a latent bug for any consumer passing mixed-case status strings.
5. **Add `MarshalJSON` to go-error-family types** — This is the single highest-leverage change for cross-project integration. Without it, HTTP handlers must hand-craft JSON responses.

### go-error-family

1. **Implement JSON serialization** — `Family.MarshalJSON()`, `Error.MarshalJSON()`. This unlocks the HTTP envelope.
2. **Define the HTTP error envelope** — `{family, code, message, fix, retryable, context}`. This bridges to templ-components.
3. **Improve diagnose coverage** — At least test the matching helpers in `diagnose.go`.

### Cross-Project

1. **Answer the 3 integration questions** — Integration depth, dependency model, Classify() caller. These are product decisions that block all implementation work.
2. **Build a thin adapter first** — Even without JSON, a Go-side `Family → AlertType` mapper in a separate package would prove the concept.

---

## F. Top 25 Things to Get Done Next

| #   | Task                                                      | Project | Impact | Effort     |
| --- | --------------------------------------------------------- | ------- | ------ | ---------- |
| 1   | **Answer go-error-family integration questions**          | Cross   | High   | Discussion |
| 2   | **Raise `utils` coverage 55% → 75%**                      | templ   | High   | Medium     |
| 3   | **Raise `forms` coverage 64% → 75%**                      | templ   | Medium | Medium     |
| 4   | **BDD tests for navigation**                              | templ   | Medium | Medium     |
| 5   | **BDD tests for layout**                                  | templ   | Medium | Medium     |
| 6   | **BDD tests for htmx**                                    | templ   | Medium | Medium     |
| 7   | **Add `aria-required` to form inputs**                    | templ   | Medium | Small      |
| 8   | **Add Avatar `alt` text**                                 | templ   | Medium | Small      |
| 9   | **Add `<html lang>` to Base**                             | templ   | Medium | Trivial    |
| 10  | **Fix ErrorAttrs for error + help text**                  | templ   | Medium | Small      |
| 11  | **Make `mapStatusToBadgeType` case-insensitive**          | templ   | Medium | Small      |
| 12  | **Unify JS attachment pattern**                           | templ   | Medium | Medium     |
| 13  | **Extract shared dismiss JS**                             | templ   | Medium | Small      |
| 14  | **Single-source toast icon SVG paths**                    | templ   | Medium | Small      |
| 15  | **Remove dead `icons.IconAttrs`**                         | templ   | Low    | Trivial    |
| 16  | **Add Table header `scope` attributes**                   | templ   | Medium | Small      |
| 17  | **Add `aria-live` to HTMX loading**                       | templ   | Medium | Small      |
| 18  | **Add `aria-live` to HTMX error handling**                | templ   | Medium | Small      |
| 19  | **Fix pre-commit hook executable**                        | templ   | Low    | Trivial    |
| 20  | **Add `MarshalJSON` to go-error-family**                  | error   | High   | Medium     |
| 21  | **Define HTTP error envelope**                            | error   | High   | Medium     |
| 22  | **Improve diagnose coverage**                             | error   | Medium | Medium     |
| 23  | **Remove zero-value DefaultXxxProps or add defaults**     | templ   | Low    | Small      |
| 24  | **Set up goreleaser**                                     | templ   | Medium | Medium     |
| 25  | **Build thin go-error-family → templ-components adapter** | Cross   | High   | Small      |

---

## G. Question I Cannot Answer Myself (#1)

**What is the intended relationship between go-error-family and templ-components in your product architecture?**

Three specific decisions I cannot make:

1. **Dependency model**: Should `templ-components` import `go-error-family` directly (every consumer gets error classification), or should integration live in a separate bridge package/module (optional)?

2. **HTTP envelope**: Should `go-error-family` define a standard JSON error response format (`{family, code, message, fix, retryable}`)? This makes it opinionated but enables the richest integration path. Or should the HTTP layer stay outside the library?

3. **Scope**: Are both libraries meant to be independently useful for other developers, or are they primarily tools for your own stack? If independent, tight coupling is wrong. If your stack, tight coupling is fine.

These are product/owner decisions. I need your direction.

---

## Metrics Snapshot

### templ-components

| Metric                  | Value                        |
| ----------------------- | ---------------------------- |
| Packages                | 9 library + 1 examples       |
| `.templ` files          | 31                           |
| `.go` files (non-templ) | 52                           |
| `_test.go` files        | 37                           |
| Total lines             | ~7,684                       |
| Components              | 53                           |
| Icons                   | 42                           |
| Typed enums             | 16+                          |
| Test coverage           | **68.3%**                    |
| Build                   | ✅ Green                     |
| Lint                    | ✅ 0 issues                  |
| Dependencies            | 2 (templ, tailwind-merge-go) |
| Go version              | 1.26.2                       |
| TODO items done         | 39 of 72 (54%)               |
| TODO items remaining    | 25 open + 8 P3               |

### Per-Package Coverage

| Package         | Coverage  | Trend                   |
| --------------- | --------- | ----------------------- |
| `internal/svg`  | 79.0%     | —                       |
| `htmx`          | 77.3%     | —                       |
| `feedback`      | 73.2%     | —                       |
| `layout`        | 72.9%     | —                       |
| `navigation`    | 72.1%     | —                       |
| `icons`         | 68.3%     | —                       |
| `display`       | 66.1%     | +0.1% (SimpleCard test) |
| `forms`         | 63.8%     | —                       |
| `utils`         | **55.3%** | —                       |
| `examples/demo` | 0.0%      | —                       |

### go-error-family

| Metric            | Value                     |
| ----------------- | ------------------------- |
| Packages          | 3 (root, agent, diagnose) |
| `.go` files       | 17                        |
| `_test.go` files  | 4                         |
| Total lines       | ~3,372                    |
| Root coverage     | **97.1%**                 |
| Agent coverage    | **100.0%**                |
| Diagnose coverage | 60.6%                     |
| Overall coverage  | 76.2%                     |
| Dependencies      | 0                         |
| Go version        | 1.26.2                    |
| License           | MIT                       |
| Version           | v0.1.1                    |

---

_Generated by Crush — 2026-05-17 23:10_
