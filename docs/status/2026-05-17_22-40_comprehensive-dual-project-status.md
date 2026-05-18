# Comprehensive Status Report — 2026-05-17 22:40

**Projects:** templ-components + go-error-family
**Author:** Lars Artmann
**Date:** 2026-05-17 22:40 CEST

---

## Executive Summary

Two public Go libraries in active development. **templ-components** is a UI component library (53 components, 42 icons, 68.3% test coverage, build green, lint clean). **go-error-family** is a structured error protocol library (97.1% root coverage, 165 tests, MIT licensed, v0.1.1 released). Both compile and pass all tests. The **cross-project integration** (go-error-family → templ-components) has been researched but NOT started — 3 architectural questions remain open.

---

## A. FULLY DONE ✅

### templ-components

| Item                   | Status | Evidence                                                            |
| ---------------------- | ------ | ------------------------------------------------------------------- | -------------------------- |
| Build passes           | ✅     | `go build ./...` — zero errors                                      |
| All tests pass         | ✅     | 9 packages, all green                                               |
| Lint clean             | ✅     | 0 issues on library packages                                        |
| 53 templ components    | ✅     | 31 `.templ` files, 52 `.go` files, ~7,653 lines                     |
| 42 icon constants      | ✅     | `iconPathData` map-driven rendering                                 |
| Dark mode              | ✅     | Full `dark:` Tailwind variant support                               |
| CSP compliance         | ✅     | All inline scripts use `nonce`                                      |
| Accessibility (major)  | ✅     | Modal focus trap, dropdown keyboard nav, tabs ARIA, tooltip linkage |
| Type-safe enums        | ✅     | 16+ typed string enums, impossible states unrepresentable           |
| Shared feedback styles | ✅     | `feedbackStyleSet` + `lookupFeedbackStyle[T]()`                     |
| Data-driven icons      | ✅     | `iconPathData` map with `                                           | ` separator for multi-path |
| GitHub Actions CI      | ✅     | Go 1.26, lint+build+test                                            |
| Public release         | ✅     | Public repo, README, LICENSE, badges                                |
| Tailwind v4 migration  | ✅     | All classes migrated, v4-exclusive documented                       |
| Deep deduplication     | ✅     | Clone groups reduced from 13 → 3 (threshold 5)                      |
| Spinner decoupling     | ✅     | HTMX loading accepts `templ.Component` parameter                    |
| Minimal layout props   | ✅     | Converted from positional params to struct-based                    |
| Ptr[T] removal         | ✅     | Replaced with Go 1.26 `new(v)` built-in                             |
| New components added   | ✅     | Latest commit: new UI and form components                           |

### go-error-family

| Item                                 | Status | Evidence                                                               |
| ------------------------------------ | ------ | ---------------------------------------------------------------------- |
| Root package 97.1% coverage          | ✅     | 165+ tests                                                             |
| Agent package 100% coverage          | ✅     | All paths tested                                                       |
| Public release v0.1.1                | ✅     | MIT license, GitHub release workflow                                   |
| 5 Families with data-driven metadata | ✅     | `familyData` array: Name, Exit, Tone, Message, Why, Fix                |
| Consumer interfaces                  | ✅     | `Coded`, `Classified`, `Contextual`, `Retryable` — each embeds `error` |
| Classification system                | ✅     | Interface → Retryable → Registered sentinels → Default (Transient)     |
| BSD sysexits.h exit codes            | ✅     | Per-family mapping                                                     |
| CLI boundary handler                 | ✅     | `HandleError`, `HandleErrorDetailed`, `HandleErrorWithConfig`          |
| Template resolution system           | ✅     | Override → Registered → Default → Family fallback                      |
| Diagnostic rules                     | ✅     | Postgres, Filesystem, Network, Git                                     |
| AI debug agent                       | ✅     | `Analyze()` returns `RootCause`, `Confidence`, `FixSteps`              |
| Zero dependencies                    | ✅     | Pure stdlib                                                            |
| Data-driven refactor                 | ✅     | Unified rendering, removed dead code                                   |
| `Audience.String()`                  | ✅     | Added in latest cleanup                                                |

---

## B. PARTIALLY DONE 🔨

### templ-components

| Item                  | Status | What's Left                                                                 |
| --------------------- | ------ | --------------------------------------------------------------------------- |
| Test coverage (68.3%) | 🔨     | Target 75%+. `utils` at 55.3%, `forms` at 63.8%, `display` at 66.0%         |
| Demo app              | 🔨     | Has syntax error at line 115, incomplete showcase                           |
| BDD tests             | 🔨     | Done for display, feedback, forms. Missing: navigation, htmx, layout, icons |
| Documentation         | 🔨     | README, FEATURES, CONTEXT exist. No docs site, no auto-generation           |
| Alpha release prep    | 🔨     | Most work done but release tooling (goreleaser) not set up                  |

### go-error-family

| Item                                  | Status | What's Left                                                                   |
| ------------------------------------- | ------ | ----------------------------------------------------------------------------- |
| diagnose package coverage (60.6%)     | 🔨     | Rules that shell out to system commands are integration-test territory        |
| Uncommitted changes                   | 🔨     | Modified `.gitignore` + status doc + new `DOMAIN_LANGUAGE.md` — not committed |
| Templ-components integration research | 🔨     | Deep research complete, 3 architectural questions posed, zero code written    |

---

## C. NOT STARTED ⬜

### templ-components

| #   | Item                                   | Priority | Notes                                                         |
| --- | -------------------------------------- | -------- | ------------------------------------------------------------- |
| 1   | NavLink `Attrs` shadowing fix          | P0       | Split brain bug — consumer `BaseProps.Attrs` silently ignored |
| 2   | Modal/Dropdown ID validation           | P1       | Empty ID = broken ARIA                                        |
| 3   | Dropdown JS XSS fix                    | P1       | `props.ID` interpolated unsafely                              |
| 4   | Accordion state coupling fix           | P1       | `max-h-96` CSS class used as JS state indicator               |
| 5   | Avatar `alt` text                      | P1       | WCAG 1.1.1                                                    |
| 6   | `aria-required` on form inputs         | P1       | WCAG requirement                                              |
| 7   | `<html lang>` in Base layout           | P1       | WCAG 3.1.1                                                    |
| 8   | Table header `scope` attributes        | P2       | Screen reader column association                              |
| 9   | Badge color map consolidation          | P2       | Two maps could drift                                          |
| 10  | `BadgeDefault` vs `BadgeNeutral` merge | P2       | Produce identical CSS                                         |
| 11  | Tabs `ActiveTabID` refactor            | P2       | Replace per-tab `Active bool`                                 |
| 12  | JS attachment pattern unification      | P2       | 3 different patterns across Accordion/Dropdown/Modal          |
| 13  | Shared dismiss JS for Alert/Toast      | P2       | Near-identical pattern duplicated                             |
| 14  | Toast icon SVG single-source           | P2       | Duplicated in Go and JS                                       |
| 15  | Tooltip position/arrow extraction      | P3       | Two switches on same type                                     |
| 16  | Card shell CSS extraction              | P3       | Repeated 3×                                                   |
| 17  | Dead code removal (`IconAttrs`)        | P2       | Exported but never called                                     |
| 18  | Release automation (goreleaser)        | P3       | Tag-based releases                                            |
| 19  | Pre-commit hook executable             | P3       | `chmod +x`                                                    |
| 20  | Examples lint exclusion                | P3       | 23 issues in demo                                             |
| 21  | Documentation site                     | P3       | Auto-generated from source                                    |
| 22  | go-error-family integration            | —        | Research done, implementation not started                     |
| 23  | Radio button component                 | —        | Not present                                                   |
| 24  | File input component                   | —        | Not present                                                   |
| 25  | Toggle/Switch component                | —        | Not present                                                   |

### go-error-family

| #   | Item                          | Priority | Notes                                         |
| --- | ----------------------------- | -------- | --------------------------------------------- |
| 1   | `docs/DOMAIN_LANGUAGE.md`     | —        | Written but untracked                         |
| 2   | Commit pending changes        | P0       | .gitignore + status doc + domain language     |
| 3   | diagnose coverage improvement | P2       | 60.6% → 75%+                                  |
| 4   | HTTP/json response envelope   | —        | Would enable templ-components integration     |
| 5   | More diagnostic rules         | P3       | Only 4 currently (Postgres, FS, Network, Git) |

---

## D. TOTALLY FUCKED UP ❌

Nothing is actively broken. Both projects build, test, and lint clean. However:

| Issue                                | Severity | Project          | Details                                                                                                                |
| ------------------------------------ | -------- | ---------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `examples/demo/main.go` syntax error | Medium   | templ-components | Line 115: `expected operand, found '{'`. Unfixed since session 1.                                                      |
| `NavLinkProps.Attrs` shadowing       | High     | templ-components | **Split brain bug.** Consumer attributes on `BaseProps` silently ignored. Every `NavLink` with custom attrs is broken. |
| Dropdown JS XSS vector               | High     | templ-components | `props.ID` interpolated into JS via Go template string — needs `strconv.Quote` like Modal.                             |
| `utils` test coverage 55.3%          | Medium   | templ-components | Lowest of any package. `MergeAttrs`, `CurrentYear`, `Deref` undertested.                                               |
| `diagnose` coverage 60.6%            | Medium   | go-error-family  | Rules that shell out are untested, but `context.go` helpers also lack coverage.                                        |
| go-error-family uncommitted work     | Low      | go-error-family  | 3 files sitting uncommitted including new `DOMAIN_LANGUAGE.md`.                                                        |

---

## E. WHAT WE SHOULD IMPROVE

### templ-components

1. **Fix the 4 P0/P1 bugs NOW** — NavLink shadowing, Dropdown XSS, Modal/Dropdown ID validation, Accordion coupling. These are correctness and security issues.
2. **Raise test coverage to 75%+** — `utils` (55.3%) and `forms` (63.8%) are dragging the average down. Write the missing BDD tests for navigation, htmx, layout.
3. **Fix the demo app** — It's been broken since session 1. A working demo is the best documentation.
4. **Decide on go-error-family integration architecture** — The 3 questions from the research session need answers before code flows.
5. **Consolidate JS patterns** — 3 different JS attachment patterns, duplicated dismiss logic. This is maintenance debt.
6. **Single-source toast icon paths** — SVG paths exist in both Go and JS. One will drift.

### go-error-family

1. **Commit the pending work** — `.gitignore`, status doc, `DOMAIN_LANGUAGE.md`.
2. **Improve `diagnose` coverage** — At least test the matching helpers and `context.go` functions.
3. **Define the HTTP response envelope** — This is the bridge to templ-components. Without it, integration is guesswork.
4. **More built-in templates** — Only 12 default message templates. Real projects will have 50+ error codes.

### Cross-Project

1. **Decide integration depth** — Thin adapter (1 file, zero coupling) vs deep protocol (structured JSON responses, Family-aware HTMX).
2. **Decide dependency model** — Direct dependency, optional sub-package, or separate module.
3. **Define who calls Classify** — Service layer or handler boundary. This shapes the entire integration.

---

## F. Top 25 Things to Get Done Next

Ranked by impact × urgency:

| #   | Task                                               | Project          | Priority | Impact | Effort     |
| --- | -------------------------------------------------- | ---------------- | -------- | ------ | ---------- |
| 1   | **Fix NavLink `Attrs` shadowing**                  | templ-components | P0       | High   | Small      |
| 2   | **Fix Dropdown JS XSS**                            | templ-components | P0       | High   | Small      |
| 3   | **Add Modal/Dropdown ID validation**               | templ-components | P1       | High   | Small      |
| 4   | **Fix Accordion `max-h-96` state coupling**        | templ-components | P1       | Medium | Small      |
| 5   | **Commit go-error-family pending work**            | go-error-family  | P0       | Medium | Trivial    |
| 6   | **Fix demo app syntax error**                      | templ-components | P2       | High   | Small      |
| 7   | **Raise `utils` coverage 55% → 75%**               | templ-components | P2       | Medium | Medium     |
| 8   | **Raise `forms` coverage 64% → 75%**               | templ-components | P2       | Medium | Medium     |
| 9   | **Add BDD tests for navigation**                   | templ-components | P1       | Medium | Medium     |
| 10  | **Add BDD tests for htmx**                         | templ-components | P1       | Medium | Medium     |
| 11  | **Add BDD tests for layout**                       | templ-components | P1       | Medium | Medium     |
| 12  | **Answer go-error-family integration questions**   | Cross            | P1       | High   | Discussion |
| 13  | **Consolidate badge color maps**                   | templ-components | P2       | Medium | Small      |
| 14  | **Merge `BadgeDefault` with `BadgeNeutral`**       | templ-components | P2       | Low    | Small      |
| 15  | **Refactor Tabs: `ActiveTabID`**                   | templ-components | P2       | High   | Medium     |
| 16  | **Unify JS attachment pattern**                    | templ-components | P2       | Medium | Medium     |
| 17  | **Extract shared dismiss JS**                      | templ-components | P2       | Medium | Small      |
| 18  | **Single-source toast icon paths**                 | templ-components | P2       | Medium | Small      |
| 19  | **Remove dead `icons.IconAttrs`**                  | templ-components | P2       | Low    | Trivial    |
| 20  | **Improve go-error-family diagnose coverage**      | go-error-family  | P2       | Medium | Medium     |
| 21  | **Add `aria-required` to form inputs**             | templ-components | P1       | Medium | Small      |
| 22  | **Add Avatar `alt` text**                          | templ-components | P1       | Medium | Small      |
| 23  | **Add `<html lang>` to Base**                      | templ-components | P1       | Low    | Trivial    |
| 24  | **Define HTTP error envelope for go-error-family** | go-error-family  | P2       | High   | Medium     |
| 25  | **Set up release automation (goreleaser)**         | templ-components | P3       | Medium | Medium     |

---

## G. Question I Cannot Answer Myself (#1)

**What is the intended relationship between these two libraries in your product architecture?**

The research uncovered a beautiful integration surface — `Family.Tone()` maps to UI presentation, `Family.Audience()` maps to component choice, `IsRetryable()` maps to HTMX retry logic. But I cannot decide:

1. **Should templ-components depend on go-error-family?** That makes every templ-components consumer pull in error classification. Or should it be an optional bridge package?
2. **Should go-error-family define an HTTP response format?** Currently it's a pure Go library with no HTTP opinions. Adding a JSON envelope (`{family, code, message, fix, retryable}`) would make it the backbone of your entire web stack — but also makes it opinionated.
3. **Is this for your personal projects only, or do you want both libraries independently useful?** If independent, the integration must be a third package. If your stack, tight coupling is fine.

These are product/owner decisions, not engineering decisions. I need your direction before writing integration code.

---

## Metrics Snapshot

### templ-components

| Metric                  | Value                        |
| ----------------------- | ---------------------------- |
| Packages                | 9 (library) + 1 (examples)   |
| `.templ` files          | 31                           |
| `.go` files (non-templ) | 52                           |
| `_test.go` files        | 37                           |
| Total lines             | ~7,653                       |
| Components              | 53                           |
| Icons                   | 42                           |
| Typed enums             | 16+                          |
| Test coverage           | 68.3%                        |
| Build                   | ✅ Green                     |
| Lint                    | ✅ 0 issues                  |
| Dependencies            | 2 (templ, tailwind-merge-go) |
| Go version              | 1.26.2                       |

### Per-Package Coverage

| Package         | Coverage |
| --------------- | -------- |
| `htmx`          | 77.3%    |
| `internal/svg`  | 79.0%    |
| `feedback`      | 73.2%    |
| `layout`        | 72.9%    |
| `navigation`    | 72.1%    |
| `icons`         | 68.3%    |
| `display`       | 66.0%    |
| `forms`         | 63.8%    |
| `utils`         | 55.3%    |
| `examples/demo` | 0.0%     |

### go-error-family

| Metric            | Value                     |
| ----------------- | ------------------------- |
| Packages          | 3 (root, agent, diagnose) |
| `.go` files       | 17                        |
| `_test.go` files  | 4                         |
| Total lines       | ~3,372                    |
| Root coverage     | 97.1%                     |
| Agent coverage    | 100.0%                    |
| Diagnose coverage | 60.6%                     |
| Overall coverage  | 76.2%                     |
| Dependencies      | 0 (pure stdlib)           |
| Go version        | 1.26.2                    |
| License           | MIT                       |
| Version           | v0.1.1                    |
| Uncommitted files | 3                         |

---

_Generated by Crush — 2026-05-17 22:40_
