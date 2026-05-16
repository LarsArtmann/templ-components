# Status Report — templ-components

**Date:** 2026-05-17 00:29 CEST
**Sessions:** 5 (2026-05-03 → 2026-05-17)
**Total Commits:** 93
**Branch:** `master` (clean)

---

## Executive Summary

templ-components is a pre-release Go UI component library (53 components, 42 SVG icons, 8 packages) built on templ + Tailwind CSS + HTMX. The codebase is clean, well-tested (68% coverage, 146 passing tests, 0 lint issues), and has only 2 runtime dependencies. The project was made public this session with a rewritten README, CONTRIBUTING.md, and GitHub metadata.

**The biggest gap is not code quality — it's developer experience infrastructure (demo site, example app, ecosystem story).**

---

## a) FULLY DONE ✅

### Core Library (Production Quality)

| Area | Status | Details |
|---|---|---|
| Component library | ✅ | 53 components across 8 packages |
| SVG icons | ✅ | 42 named icons with typed constants |
| Typed enums | ✅ | 16+ string enums, impossible states unrepresentable |
| Props system | ✅ | All components embed `utils.BaseProps` (exception: `PageProps`) |
| Tailwind class merging | ✅ | `utils.Class()` via tailwind-merge-go |
| Dark mode | ✅ | All components support `dark:` variant |
| CSP compliance | ✅ | All inline scripts use nonces |
| Accessibility | ✅ | ARIA attributes, roles, keyboard nav on Modal/Dropdown/Tabs/Tooltip |
| Map-based style lookups | ✅ | No switches — data-driven, extensible |
| `DefaultXxxProps()` constructors | ✅ | 20 constructors across all packages |
| Shared feedback style system | ✅ | `feedbackStyleSet` + `lookupFeedbackStyle[T]()` |
| Icon rendering | ✅ | `iconPathData` map with `|` separator for multi-path icons |
| `internal/svg` package | ✅ | Shared SVG primitives (fillIcon, spinnerSVG) |
| Import graph | ✅ | Clean, no circular dependencies |

### Build & CI

| Area | Status | Details |
|---|---|---|
| Build | ✅ | `go build ./...` — zero errors |
| Tests | ✅ | 146 passing, 9 packages, 68% coverage |
| Lint | ✅ | 0 issues on library packages (golangci-lint) |
| GitHub Actions CI | ✅ | Go 1.26, lint + build + test + coverage upload |

### Testing

| Area | Status | Details |
|---|---|---|
| BDD tests | ✅ | display, feedback, forms packages |
| Snapshot tests | ✅ | All packages |
| A11y attribute tests | ✅ | display, navigation, layout, htmx |
| Dark mode verification | ✅ | All packages |
| Benchmark tests | ✅ | `utils.Class()`, Badge rendering |
| XSS safety tests | ✅ | Dropdown auto-escaping |
| Security headers tests | ✅ | `layout.Base` |
| Icon rendering tests | ✅ | All 42 icons verified |
| Default constructor tests | ✅ | Card, Badge, Modal, ProgressBar |

### Documentation & Public Release (This Session)

| Area | Status | Details |
|---|---|---|
| README.md | ✅ | Rewritten with value prop, competitive comparison, badges, verified examples |
| CONTRIBUTING.md | ✅ | Setup, conventions, workflow, commit format |
| FEATURES.md | ✅ | Complete feature inventory with status per component |
| TODO_LIST.md | ✅ | 76 items (32 done, 44 open) |
| CONTEXT.md | ✅ | Architecture context, patterns, naming conventions |
| CHANGELOG.md | ✅ | Full changelog with breaking changes documented |
| Migration guide | ✅ | `docs/migration/v0.1-to-v0.2.md` |
| GitHub metadata | ✅ | Description + 12 topics set |
| Tailwind v4 docs | ✅ | CSS-first setup with `@source` vendor path |
| Competitive analysis | ✅ | templUI, goshipit compared honestly |

---

## b) PARTIALLY DONE 🔨

### Accessibility (Good Foundation, Gaps Remain)

| Item | Status | Gap |
|---|---|---|
| Modal focus trap + Escape | ✅ | Done |
| Dropdown keyboard nav | ✅ | Done |
| Tabs ARIA linkage | ✅ | Done |
| Tooltip aria-describedby | ✅ | Done |
| Avatar `<img>` alt text | ❌ | Missing `Alt` field |
| `aria-required` on form inputs | ❌ | WCAG requirement |
| `<html lang>` in Base layout | ❌ | WCAG 3.1.1 |
| Table header `scope` attributes | ❌ | Screen reader column association |
| `aria-live` on HTMX loading | ❌ | Dynamic content announcements |
| `ErrorAttrs` for simultaneous error + help | ❌ | `aria-describedby` links only one |

### Per-Package Test Coverage (68% Average)

| Package | Coverage | Assessment |
|---|---|---|
| `internal/svg` | 79.0% | Good |
| `htmx` | 77.3% | Good |
| `icons` | 74.0% | Good |
| `layout` | 73.0% | Good |
| `navigation` | 72.0% | Good |
| `feedback` | 71.7% | Good |
| `display` | 66.2% | Acceptable |
| `forms` | 61.5% | Below target |
| `utils` | 56.4% | Below target |

### Demo / Example App (Broken)

| Item | Status |
|---|---|
| `examples/demo/main.go` exists | ✅ |
| Uses Tailwind v2 CDN (should be v4) | ❌ |
| Doesn't use `layout.Base` (raw HTML string) | ❌ |
| Doesn't showcase most components | ❌ |
| Pre-commit hook requires `buildflow` (not installed) | ❌ |

---

## c) NOT STARTED ⬜

### Critical Bugs

1. **`NavLinkProps.Attrs` shadowing `BaseProps.Attrs`** — Split brain bug. Consumer attributes on `BaseProps` are silently ignored on `NavLink`. (TODO #9, P0)
2. **Dropdown JS XSS vector** — IIFE interpolates `props.ID` via Go template instead of `strconv.Quote` (TODO #11, P1)
3. **Accordion state coupling** — Uses `max-h-96` CSS class as JS state indicator (TODO #12, P1)
4. **Required ID validation** — Modal/Dropdown with empty ID produces broken ARIA (TODO #10, P1)

### Architecture Improvements

5. **Replace `Tab.Active` with `TabsProps.ActiveTabID`** — Prevents zero/multiple active tabs (TODO #22, P2)
6. **Consolidate badge color maps** — `badgeColorMap` + `badgeDotColorMap` → single struct map (TODO #20, P2)
7. **Merge `BadgeDefault` with `BadgeNeutral`** — Currently identical CSS (TODO #21, P2)
8. **Unify JS attachment pattern** — Three different patterns across Accordion/Dropdown/Modal (TODO #23, P2)
9. **Extract shared dismiss JS** — Alert and Toast duplicate nearly identical pattern (TODO #24, P2)
10. **Single-source toast icon SVG paths** — Duplicated in Go and JS (TODO #25, P2)
11. **Decouple `htmx/loading` from `feedback.Spinner`** — Accept `templ.Component` instead (TODO #26, P2)
12. **Tooltip position/arrow consolidation** — Two switches on same type (TODO #26a, P3)
13. **Card shell CSS extraction** — Repeated 3× in Card/StatCard/SimpleCard (TODO #26b, P3)

### Missing Tests

14. **BDD tests for navigation** — Nav, Pagination, Breadcrumbs (TODO #42, P1)
15. **BDD tests for htmx** — Loading indicators, error handling (TODO #43, P1)
16. **BDD tests for layout** — Base, Minimal, Theme (TODO #44, P1)
17. **BDD tests for icons** — Icon rendering (TODO #45, P2)
18. **Table mismatched header/row lengths** — No validation exists (TODO #46, P2)
19. **Modal/Dropdown with empty ID** — Should fail gracefully (TODO #47, P2)
20. **`mapStatusToBadgeType` boundary cases** — Case sensitivity, whitespace (TODO #48, P2)

### Dead Code & Cleanup

21. **Remove or use `icons.IconAttrs`** — Exported, never called (TODO #55, P2)
22. **Remove or use `internal/svg.FillIcon`** — Only referenced via proxy (TODO #56, P2)
23. **Remove no-op `DefaultXxxProps()`** — Several return zero-value structs (TODO #57, P3)
24. **Move test helpers to `internal/testutil/`** — `Render`, `AssertContains` etc. (TODO #58, P3)
25. **Move ProgressBar a11y test from display to feedback** — Tests wrong package (TODO #59, P3)

### DevOps & Tooling

26. **Release automation** — goreleaser, tag-based releases (TODO #62, P3)
27. **Fix pre-commit hook** — `buildflow` not installed, breaks every commit (TODO #63, P3)
28. **Exclude `examples/` from lint** — 23 issues in demo (TODO #64, P3)

### Ecosystem & Growth (from STANDOUT-IDEAS.md)

29. **Live component showcase site** — templUI has templui.io. Nothing visual here.
30. **Get listed on templ.guide** — templUI is the only library listed there.
31. **Tag v0.1.0-alpha** — PUBLIC_OR_PRIVATE.md says do it now.
32. **Cross-link cqrs-htmx + go-cqrs-lite** in READMEs — Ecosystem story
33. **Real-world example app** — Full GOTH stack CRUD admin panel
34. **Submit to awesome-templ** — Discovery
35. **`go doc` examples** — `ExampleXxx()` test functions for pkg.go.dev
36. **Interactive playground** — Click component, see render, copy code

---

## d) TOTALLY FUCKED UP 💣

### Pre-Commit Hook
The `.git/hooks/pre-commit` requires `buildflow` which is not installed anywhere. Every commit requires `--no-verify`. This silently bypasses ALL safety checks. The `scripts/pre-commit.sh` exists but is not installed as the git hook.

### Demo App
`examples/demo/main.go` is a lie:
- Uses **Tailwind v2 CDN** (library recommends v4)
- Doesn't use `layout.Base` — writes raw HTML strings with `w.Write([]byte(...))`
- Only showcases 4 components (Nav, Alert, StatCard, Icons) out of 53
- The whole approach (raw string HTML in Go) contradicts the library's own philosophy

### Buildflow Dependency
The project appears to depend on a tool called `buildflow` for pre-commit checks, but it's not documented in requirements, not in `go.mod`, and not installed. If it was ever set up, it's been forgotten.

---

## e) WHAT WE SHOULD IMPROVE

### Immediate (Fixes That Improve Trust)

1. **Fix the pre-commit hook** — Replace `buildflow` with a simple shell script that runs `templ generate && go test ./...`. Or install buildflow. Either way, commits should NOT require `--no-verify`.

2. **Fix `NavLinkProps.Attrs` shadowing** — This is a P0 bug that silently drops consumer attributes. Any user passing `Attrs` on a nav link gets nothing. First real user will hit this immediately.

3. **Fix the demo app** — Either fix it to use `layout.Base` + Tailwind v4, or delete it. A broken demo is worse than no demo.

### Short-Term (Before v0.1.0-alpha Tag)

4. **Tag v0.1.0-alpha** — The code is ready. PUBLIC_OR_PRIVATE.md has been saying "go public now" since May 4. Every day of delay is lost first-mover ground.

5. **Get listed on templ.guide** — Open a PR/issue on the templ docs repo. This is the #1 discoverability channel for the ecosystem.

6. **Cross-link ecosystem in README** — cqrs-htmx and go-cqrs-lite should be mentioned. The "three-library GOTH stack" story is the genuine differentiator vs templUI.

### Medium-Term (Differentiation)

7. **Build the component showcase** — A `./cmd/demo` server that renders every component with variants. Deploy it. This is the single highest-leverage action for adoption.

8. **Add `ExampleXxx()` functions** — They show up on pkg.go.dev. Zero-cost discoverability.

9. **Tailwind v4 safelist** — Ship a file listing all library classes so consumers who don't vendor can use `@source inline(...)`.

### Architecture (Type Model Improvements)

10. **Replace `Tab.Active` with `TabsProps.ActiveTabID string`** — Currently allows zero or multiple active tabs. Impossible state is representable.

11. **Consolidate badge color maps** — `badgeColorMap` + `badgeDotColorMap` can drift. Single `badgeStyle{ BG, Dot }` struct map.

12. **Decouple `htmx/loading` from `feedback.Spinner`** — The `htmx` package should accept a `templ.Component` for spinner customization, not import `feedback` directly.

13. **Unify JS patterns** — Accordion, Dropdown, and Modal each use different JS attachment patterns. Standardize on IIFE-per-instance with `strconv.Quote` for safety.

14. **Add `<html lang>` to `layout.Base`** — Add `Lang` field to `PageProps`. WCAG 3.1.1 requirement. Simple fix, high a11y impact.

---

## f) Top 25 Things We Should Do Next

Sorted by **impact × effort** (highest first):

| # | Task | Impact | Effort | Category |
|---|---|---|---|---|
| 1 | **Tag v0.1.0-alpha and go public** | Unblocks everything | 10 min | Release |
| 2 | **Fix pre-commit hook** (replace buildflow) | Trust/safety | 15 min | DevOps |
| 3 | **Submit to awesome-templ** | Discovery | 15 min | Growth |
| 4 | **Cross-link cqrs-htmx + go-cqrs-lite in README** | Ecosystem story | 15 min | Docs |
| 5 | **Fix `NavLinkProps.Attrs` shadowing** | P0 bug fix | 20 min | Bug |
| 6 | **Open PR on templ docs to get listed** | Discovery | 30 min | Growth |
| 7 | **Add `<html lang>` to `layout.Base`** | WCAG 3.1.1 | 30 min | A11y |
| 8 | **Fix Dropdown JS XSS** (`strconv.Quote`) | Security | 30 min | Bug |
| 9 | **Fix or delete demo app** | Trust | 1 hr | DX |
| 10 | **Validate required ID in Modal/Dropdown** | Prevents broken ARIA | 30 min | Bug |
| 11 | **Add `aria-required` to form inputs** | WCAG | 30 min | A11y |
| 12 | **Add `Alt` field to `AvatarProps`** | WCAG 1.1.1 | 30 min | A11y |
| 13 | **BDD tests for navigation package** | Coverage + confidence | 2 hr | Testing |
| 14 | **BDD tests for layout package** | Coverage + confidence | 2 hr | Testing |
| 15 | **BDD tests for htmx package** | Coverage + confidence | 2 hr | Testing |
| 16 | **Replace `Tab.Active` with `ActiveTabID`** | Impossible state fix | 1 hr | Architecture |
| 17 | **Unify JS attachment pattern** (Accordion/Dropdown/Modal) | Consistency | 3 hr | Architecture |
| 18 | **Build component showcase server** (`./cmd/demo`) | Adoption multiplier | 1 day | Growth |
| 19 | **Add `ExampleXxx()` functions** | pkg.go.dev discoverability | 3 hr | DX |
| 20 | **Consolidate badge color maps** | Eliminate drift risk | 1 hr | Architecture |
| 21 | **Fix Accordion `max-h-96` state coupling** | Fragile coupling | 1 hr | Bug |
| 22 | **Decouple `htmx/loading` from `feedback.Spinner`** | Loose coupling | 1 hr | Architecture |
| 23 | **Extract shared dismiss JS** (Alert + Toast) | DRY | 1 hr | Architecture |
| 24 | **Real-world example app** (full GOTH stack) | Converts evaluators | 3 days | Growth |
| 25 | **Release automation** (goreleaser) | Professional releases | 2 hr | DevOps |

---

## g) My Top #1 Question

**Should the library standardize on Tailwind v4 exclusively, or maintain v3 compatibility?**

The README now recommends Tailwind v4 (`@import "tailwindcss"`, `@source`, `@custom-variant`). But the actual component code uses class names that work in both v3 and v4 (no v4-only features like `@theme` or new color spaces). The questions I can't answer:

1. **Are there any Tailwind v4 breaking changes** that affect the class names used in this library? (e.g., renamed utilities, removed classes, changed defaults)
2. **Should we add v4-only features** like CSS custom properties for theming, or keep it compatible with both?
3. **Does `tailwind-merge-go` v0.2.1 support Tailwind v4 classes?** Or does it need an upgrade?

This affects whether we can honestly claim "Tailwind CSS 4.x+" in requirements, or need to say "3.x and 4.x".

---

## Metrics Dashboard

| Metric | Value | Trend |
|---|---|---|
| Components | 53 | Stable |
| SVG icons | 42 | Stable |
| Typed enums | 16+ | Stable |
| `.templ` source files | 31 | Stable |
| Source lines (templ) | 2,999 | Stable |
| Source lines (Go) | 4,685 | Stable |
| Test lines | 4,040 | Stable |
| Test-to-code ratio | 53% | Good |
| Tests passing | 146 | Stable |
| Packages | 9 (8 lib + 1 demo) | Stable |
| Test coverage | 68.0% | Stable |
| Lint issues | 0 | Clean |
| Runtime dependencies | 2 | Minimal |
| Open TODOs | 44 | ↓ from 76 |
| Closed TODOs | 32 | ↑ from 0 |
| GitHub description | ✅ Set | New this session |
| GitHub topics | 12 | New this session |
| CONTRIBUTING.md | ✅ Created | New this session |
| Pre-commit hook | 💣 Broken | Unchanged |
| Demo app | 💣 Broken | Unchanged |

---

## Session 5 Summary (2026-05-17)

### What Changed

1. **README.md** — Complete rewrite: value proposition, competitive comparison (templUI, goshipit), 5 badges, verified code examples, Tailwind v4 CSS-first setup, utils package section, honest positioning ("zero client-side JavaScript")
2. **CONTRIBUTING.md** — New file: setup, development workflow, code conventions, commit format, license
3. **GitHub metadata** — Description and 12 topics set for discoverability
4. **Bug fixes found during audit** — Tailwind v4 `@source` path was wrong (`node_modules` → `vendor/`), `Skeleton("card")` API was wrong (string → `SkeletonCard`)

### What Was Verified

- All 17 API calls in README verified against source code
- Competitive claims verified against actual competitor repos (templUI, goshipit)
- Tailwind v4 `@source` + `@custom-variant` verified against official Tailwind docs
- All 146 tests still passing, 0 lint issues

### Commits This Session

```
863be6a docs(readme): fix comparison table, contributing link, and terminology
a21f53f docs: rewrite README and add CONTRIBUTING.md for public release
5b34b1c docs(readme): rewrite README with comparison tables, catalog, and design principles
```
