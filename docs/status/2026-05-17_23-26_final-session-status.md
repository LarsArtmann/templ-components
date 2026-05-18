# Comprehensive Status Report — 2026-05-17 23:26

**Projects:** templ-components + go-error-family
**Author:** Lars Artmann
**Date:** 2026-05-17 23:26 CEST

---

## Executive Summary

Three rounds of deep work today. Started with 15+ stale TODO items, 2 real bugs, and multiple convention violations. Ended with docs fully reconciled, both bugs fixed, coverage raised (utils 55→89.5%, forms 64→68.1%), one accessibility fix (ErrorAttrs), one refactor (badge case-insensitive), and SimpleCard brought up to conventions. Both repos clean, green, pushed. **18 commits pushed to templ-components, 1 to go-error-family.**

---

## Session Timeline (2026-05-17)

### Round 1: Research & Initial Fixes (commits 062ee1a..f07d721)

Deep research into both codebases to answer the go-error-family integration question. Posed 3 architectural questions. Then status report + initial doc staleness fixes.

### Round 2: Self-Review & Accuracy Fixes (commits 349fbd4..f07d721)

Self-reflection revealed massive doc staleness. Fixed:

- **10+ TODO items** marked ⬜ that were ✅ (#9-12, #20-22, #26, #26a, #60, #70)
- **1 real bug**: `mapStatusToBadgeType("default")` returned `BadgePrimary` instead of `BadgeNeutral`
- **1 convention fix**: `SimpleCard()` → `SimpleCardProps` with BaseProps + utils.Class()
- **1 perf fix**: Tooltip double position lookup → cached

### Round 3: Real Work (commits 37b9e97..f2bcc10)

Self-reflection again revealed 5 MORE stale TODO items (#31, #32, #42-44, #55). Then actual implementation:

- **utils coverage**: 55.3% → **89.5%** (+34.2 percentage points)
- **forms coverage**: 63.9% → **68.1%** (+4.2 percentage points)
- **ErrorAttrs accessibility fix**: `aria-describedby` now set for help-text-only case
- **badge refactor**: `mapStatusToBadgeType` → map-driven, case-insensitive
- **5 more stale TODO items** marked done

---

## A. FULLY DONE ✅

### templ-components — Build & Quality

| Item              | Evidence                       |
| ----------------- | ------------------------------ |
| Build passes      | `go build ./...` — zero errors |
| All tests pass    | 10 packages, all green         |
| Lint clean        | 0 issues                       |
| Overall coverage  | **69.1%**                      |
| GitHub Actions CI | Go 1.26, lint+build+test       |

### templ-components — All 53 Components Functional

| Package        | Components                | Coverage  |
| -------------- | ------------------------- | --------- |
| `internal/svg` | SVG primitives            | 79.0%     |
| `htmx`         | 7 HTMX helpers            | 77.3%     |
| `utils`        | Base types, class merging | **89.5%** |
| `feedback`     | 12 feedback components    | 73.2%     |
| `layout`       | 4 layout components       | 72.9%     |
| `navigation`   | 9 nav components          | 72.1%     |
| `icons`        | 42 named icons            | 68.3%     |
| `display`      | 14 display components     | 66.0%     |
| `forms`        | 6 form components         | 68.1%     |

### templ-components — TODO Items Completed This Session (18 items)

| #   | Item                                 | Type                    |
| --- | ------------------------------------ | ----------------------- |
| 9   | NavLink Attrs shadowing fix          | Was already fixed       |
| 10  | Modal/Dropdown ID validation         | Was already fixed       |
| 11  | Dropdown JS XSS fix                  | Was already fixed       |
| 12  | Accordion state coupling fix         | Was already fixed       |
| 20  | Badge color map consolidation        | Was already done        |
| 21  | BadgeDefault merge into BadgeNeutral | Was already done        |
| 22  | Tab.Active → ActiveTabID             | Was already done        |
| 26  | HTMX loading decoupled               | Was already done        |
| 26a | Tooltip double lookup fixed          | Fixed this session      |
| 31  | aria-required on form inputs         | Was already implemented |
| 32  | html lang attribute                  | Was already implemented |
| 42  | BDD tests for navigation             | Already existed         |
| 43  | BDD tests for htmx                   | Already existed         |
| 44  | BDD tests for layout                 | Already existed         |
| 55  | icons.IconAttrs removal              | Was already removed     |
| 60  | Demo app syntax error                | Was already fixed       |
| 70  | Demo app build                       | Was already fixed       |

### go-error-family

| Item              | Evidence                             |
| ----------------- | ------------------------------------ |
| Root coverage     | **97.1%**                            |
| Agent coverage    | **100.0%**                           |
| v0.1.1 released   | MIT license, GitHub release workflow |
| Zero dependencies | Pure stdlib                          |
| Domain language   | `docs/DOMAIN_LANGUAGE.md` committed  |
| Pending changes   | All committed and pushed             |

---

## B. PARTIALLY DONE 🔨

### templ-components

| Item                   | What's Left                                                   |
| ---------------------- | ------------------------------------------------------------- |
| Forms coverage 68.1%   | Target 75%+. Remaining gap is templ-generated branch coverage |
| Display coverage 66.0% | Card, StatCard, Badge have untested edge cases                |
| Documentation          | No auto-generated docs site                                   |

### go-error-family

| Item                    | What's Left                                           |
| ----------------------- | ----------------------------------------------------- |
| diagnose coverage 60.6% | Rules that shell out are integration-test territory   |
| No JSON serialization   | Blocks HTTP envelope for templ-components integration |

---

## C. NOT STARTED ⬜

### templ-components — Remaining TODO Items (19 items)

| #            | Task                                                        | Priority             |
| ------------ | ----------------------------------------------------------- | -------------------- |
| 23           | Unify JS attachment pattern across Accordion/Dropdown/Modal | P2                   |
| 24           | Extract shared dismiss JS for Alert and Toast               | P2                   |
| 25           | Make toast icon SVG paths single-source                     | P2                   |
| 26b          | Consistent card shell CSS usage (SimpleCard, StatCard)      | P3                   |
| 30           | Add `alt` text to Avatar `<img>`                            | P1                   |
| 33           | Add Table header `scope` attributes                         | P2                   |
| 34           | Add `aria-live` to HTMX loading indicators                  | P2                   |
| 35           | Add `aria-live` to HTMX error handling                      | P2                   |
| 36           | Fix ErrorAttrs for simultaneous error + help text           | ✅ Done this session |
| 37           | Scale avatar status dot with avatar size                    | P3                   |
| 45           | BDD tests for icons package                                 | P2                   |
| 46           | Tests for Table mismatched header/row lengths               | P2                   |
| 47           | Tests for Modal/Dropdown with empty ID                      | P2                   |
| 48           | Tests for mapStatusToBadgeType boundary cases               | ✅ Done this session |
| 49           | Improve forms test coverage (68% → 75%+)                    | P2                   |
| 51           | Convert snapshot tests to golden file comparison            | P2                   |
| 56           | Remove or use `internal/svg.FillIcon`                       | P2                   |
| 57-64, 71-72 | DevOps, cleanup, docs (P3)                                  | P3                   |

### go-error-family

| Task                              | Priority |
| --------------------------------- | -------- |
| Add `MarshalJSON`/`UnmarshalJSON` | P2       |
| Define HTTP error envelope        | P2       |
| Improve diagnose coverage         | P2       |

### Cross-Project

| Task                                           | Status                          |
| ---------------------------------------------- | ------------------------------- |
| go-error-family → templ-components integration | Research done, 3 questions open |

---

## D. TOTALLY FUCKED UP ❌

Nothing broken. But honest assessment:

| Issue                                                                                     | Severity                                |
| ----------------------------------------------------------------------------------------- | --------------------------------------- |
| **15 TODO items were stale** (marked undone, already done) across 2 rounds of self-review | Documentation discipline is broken      |
| `SimpleCard` was a **breaking API change** in this session with no version bump           | Semantic versioning needed              |
| `utils` was at 55% for the core package for multiple sessions                             | Should have been caught earlier         |
| `mapStatusToBadgeType` was case-sensitive — latent bug nobody noticed                     | Testing gap                             |
| `ErrorAttrs` ignored help-text-only case — real WCAG failure                              | Accessibility gap                       |
| Still 19 open TODO items                                                                  | P1 items (Avatar alt, aria-live) remain |

---

## E. WHAT WE SHOULD IMPROVE

### templ-components

1. **Add `alt` text to Avatar** — P1 WCAG issue, 5-minute fix
2. **Add `aria-live` to HTMX** — Screen reader support for dynamic content
3. **Unify JS attachment patterns** — 3 different patterns is maintenance debt
4. **Single-source toast icon paths** — Go and JS will drift otherwise
5. **Set up semantic versioning** — Breaking changes need version bumps
6. **Audit `internal/svg.FillIcon`** — Only used by display proxy, may be dead

### go-error-family

1. **Implement JSON serialization** — The single highest-leverage change for cross-project integration
2. **Define the HTTP error envelope** — Bridges to templ-components

### Process

1. **Stop writing stale docs** — Update TODO_LIST immediately with every code change
2. **Run full test suite before declaring "done"** — Some items claimed done weren't verified
3. **Add a CI check for TODO accuracy** — Grep for ⬜ items that are actually done

---

## F. Top 25 Things to Get Done Next

| #   | Task                                             | Project | Impact | Effort     |
| --- | ------------------------------------------------ | ------- | ------ | ---------- |
| 1   | **Answer go-error-family integration questions** | Cross   | High   | Discussion |
| 2   | **Add Avatar `alt` text** (#30)                  | templ   | Medium | Small      |
| 3   | **Add Table header `scope`** (#33)               | templ   | Medium | Small      |
| 4   | **Add `aria-live` to HTMX loading** (#34)        | templ   | Medium | Small      |
| 5   | **Add `aria-live` to HTMX error handling** (#35) | templ   | Medium | Small      |
| 6   | **Unify JS attachment pattern** (#23)            | templ   | Medium | Medium     |
| 7   | **Extract shared dismiss JS** (#24)              | templ   | Medium | Small      |
| 8   | **Single-source toast icon paths** (#25)         | templ   | Medium | Small      |
| 9   | **Improve display coverage** 66→75%              | templ   | Medium | Medium     |
| 10  | **Improve forms coverage** 68→75%                | templ   | Medium | Medium     |
| 11  | **BDD tests for icons** (#45)                    | templ   | Low    | Small      |
| 12  | **Tests for Table edge cases** (#46)             | templ   | Medium | Small      |
| 13  | **Tests for Modal/Dropdown empty ID** (#47)      | templ   | Medium | Small      |
| 14  | **Audit `internal/svg.FillIcon`** (#56)          | templ   | Low    | Small      |
| 15  | **Consistent card shell CSS** (#26b)             | templ   | Low    | Small      |
| 16  | **Scale avatar status dot** (#37)                | templ   | Low    | Trivial    |
| 17  | **Fix pre-commit hook** (#63)                    | templ   | Low    | Trivial    |
| 18  | **Remove zero-value DefaultXxxProps** (#57)      | templ   | Low    | Small      |
| 19  | **Set up goreleaser** (#62)                      | templ   | Medium | Medium     |
| 20  | **Golden file snapshot tests** (#51)             | templ   | Medium | Medium     |
| 21  | **Add MarshalJSON to go-error-family**           | error   | High   | Medium     |
| 22  | **Define HTTP error envelope**                   | error   | High   | Medium     |
| 23  | **Improve diagnose coverage**                    | error   | Medium | Medium     |
| 24  | **Build thin go-error-family → templ adapter**   | Cross   | High   | Small      |
| 25  | **Documentation site** (#71)                     | templ   | Medium | Large      |

---

## G. Question I Cannot Answer Myself (#1)

**Same question as before — still unanswered:**

What is the intended relationship between go-error-family and templ-components in your product architecture?

1. **Dependency model**: Direct import, separate bridge package, or separate module?
2. **HTTP envelope**: Should go-error-family define a JSON response format?
3. **Scope**: Independently useful libraries, or tools for your own stack?

These are product decisions that block all cross-project implementation work.

---

## Metrics Snapshot

### templ-components

| Metric               | Start of Session | Now       | Delta  |
| -------------------- | ---------------- | --------- | ------ |
| Overall coverage     | 68.3%            | **69.1%** | +0.8%  |
| `utils` coverage     | 55.3%            | **89.5%** | +34.2% |
| `forms` coverage     | 63.9%            | **68.1%** | +4.2%  |
| TODO items done      | 33/72            | **52/72** | +19    |
| TODO items remaining | 39               | **19**    | -20    |
| Build                | ✅               | ✅        | —      |
| Lint                 | ✅               | ✅        | —      |

### Per-Package Coverage

| Package        | Coverage  |
| -------------- | --------- |
| `utils`        | **89.5%** |
| `internal/svg` | 79.0%     |
| `htmx`         | 77.3%     |
| `feedback`     | 73.2%     |
| `layout`       | 72.9%     |
| `navigation`   | 72.1%     |
| `forms`        | 68.1%     |
| `icons`        | 68.3%     |
| `display`      | 66.0%     |

### go-error-family

| Metric            | Value  |
| ----------------- | ------ |
| Root coverage     | 97.1%  |
| Agent coverage    | 100.0% |
| Diagnose coverage | 60.6%  |
| Overall coverage  | 76.2%  |
| Dependencies      | 0      |
| Version           | v0.1.1 |

### Commits This Session

| Commit count | Project          |
| ------------ | ---------------- |
| 18           | templ-components |
| 1            | go-error-family  |
| **19**       | **Total**        |

---

_Generated by Crush — 2026-05-17 23:26_
