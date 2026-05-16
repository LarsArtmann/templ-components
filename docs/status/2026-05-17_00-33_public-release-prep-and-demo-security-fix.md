# Status Report — templ-components

**Date:** 2026-05-17 00:33 UTC  
**Branch:** master (1 commit ahead of origin)  
**Session:** 6 — Public release prep, demo security fix  

---

## Project Snapshot

| Metric | Value |
|---|---|
| Version | 0.x (pre-release, no tag) |
| Packages | 8 component + 1 internal + 1 example |
| Components | 55 exported templ components |
| Icons | 43 named SVG icons |
| Typed enums | 16+ string enums |
| Test functions | 147 |
| Handwritten LOC | ~3,853 (Go + templ) |
| Test LOC | ~3,000+ |
| Runtime deps | 2 (`templ` + `tailwind-merge-go`) |
| Lint issues | 0 (library packages) |
| Build status | **PASSING** |
| CI | GitHub Actions (lint + build + test with coverage) |

---

## a) FULLY DONE ✅

### Core Library (Sessions 1–4)
- **55 templ components** across 8 packages: display (14), feedback (12), forms (6), htmx (7), icons (1+42), layout (4), navigation (9), utils (0, shared types)
- **All components embed `utils.BaseProps`** with ID, Class, Attrs, AriaLabel, Nonce (exception: `layout.PageProps`)
- **Map-based style lookups** everywhere — no switch statements
- **Shared `feedbackStyleSet` + `lookupFeedbackStyle[T]()`** for alert/toast consistency
- **`iconPathData` map** — data-driven icon rendering with `|` separator for multi-path icons
- **16+ typed string enums** — AvatarStatus, TrendDirection, AlertType, BadgeType, ModalSize, etc.
- **`DefaultXxxProps()` constructors** for every component with non-zero defaults
- **CSP compliance** — all inline scripts use `nonce={ props.Nonce }`
- **`utils.Class()`** — Tailwind class conflict resolution via tailwind-merge-go

### Accessibility (Session 3)
- Modal: focus trap, Escape key, focus management
- Dropdown: arrow key navigation, Escape to close
- Tabs: proper ARIA tablist/tab/tabpanel linkage
- Tooltip: deterministic ID for `aria-describedby`
- Accordion: `aria-expanded`, `aria-controls`, chevron rotation
- Forms: `ErrorAttrs()` for `aria-invalid` + `aria-describedby`
- Avatar: `aria-hidden` on status dot

### Testing (Sessions 1–4)
- **147 test functions** across 37 test files
- BDD-style tests for display, feedback, forms packages
- A11y attribute validation tests
- Dark mode class verification tests
- Benchmark tests for hot paths (`utils.Class()`, Badge render)
- XSS safety test for Dropdown
- Security headers test for `layout.Base`
- Layout `DefaultPageProps()` constructor test

### Documentation (Sessions 1–5)
- **README.md** — Quick start, comparison tables, component catalog, design principles
- **CONTRIBUTING.md** — Contribution guidelines
- **FEATURES.md** — Complete feature inventory with status per component
- **TODO_LIST.md** — 72-item tracked list with priorities
- **CONTEXT.md** — Architecture context, key patterns, naming conventions
- **CHANGELOG.md** — Full changelog with breaking changes documented
- **docs/migration/v0.1-to-v0.2.md** — Migration guide for breaking changes
- **docs/adr/0001-shared-svg-helpers.md** — Architecture decision record
- **docs/diagrams/** — D2 architecture diagrams (current + future)
- **docs/modularization/** — Modularization proposal, dependency graph, execution plan

### DevOps
- **GitHub Actions CI** — lint (golangci-lint v6), build, test with `-race`, coverage artifact upload
- **Pre-commit hook** — `scripts/pre-commit.sh` for auto `templ generate`
- **.golangci.yml** — 59 linters enabled, Go 1.26, 5min timeout

### Session 5 (2026-05-17)
- **README.md table formatting** — Fixed comparison table column alignment
- **"By the Numbers" table** — Fixed metric table alignment
- **CONTRIBUTING.md** — Added contribution guidelines for public release

### Session 6 (2026-05-17, this session)
- **XSS fix in demo app** — `examples/demo/main.go:96`: Added `templ.EscapeString(title)` to prevent HTML injection in `renderSection()`. Raw `%s` was writing unescaped user-controlled string into HTML output.

---

## b) PARTIALLY DONE 🔨

### Public Release Checklist (from PUBLIC_OR_PRIVATE.md)
| # | Condition | Status |
|---|---|---|
| 1 | Tag v0.1.0-alpha with disclaimer | ⬜ Not done |
| 2 | Fix empty CI — add working GitHub Actions | ✅ **DONE** (ci.yaml is complete) |
| 3 | Update CHANGELOG.md | ✅ **DONE** (comprehensive entries) |
| 4 | Add "Alpha / Pre-release" badge to README | ⬜ Not done |
| 5 | Add CONTRIBUTING.md | ✅ **DONE** |
| 6 | Submit to awesome-templ | ⬜ Not done |

### Demo App (`examples/demo/main.go`)
- **Partially converted** — Shows Nav, Alerts, StatCards, Icons
- **Missing showcase**: Accordion, Avatar, Badge, Breadcrumbs, Card, Checkbox, Dropdown, Footer, Modal, Pagination, ProgressBar, Select, Skeleton, Spinner, StepIndicator, Table, Tabs, Textarea, Tooltip
- **Still uses raw `w.Write`/`fmt.Fprintf`** for HTML instead of templ components
- **Security fixed** this session (XSS in section header)

---

## c) NOT STARTED ⬜

### Critical Bugs & Type Safety (from TODO_LIST.md)
| # | Task | Priority |
|---|---|---|
| 9 | Fix `NavLinkProps.Attrs` shadowing `BaseProps.Attrs` | P0 |
| 10 | Validate required `ID` in Modal and Dropdown | P1 |
| 11 | Fix Dropdown JS XSS vector (IIFE interpolates `props.ID`) | P1 |
| 12 | Fix Accordion state coupling with `max-h-96` | P1 |

### Architecture Improvements
| # | Task | Priority |
|---|---|---|
| 20 | Consolidate badge color maps into single struct map | P2 |
| 21 | Merge `BadgeDefault` with `BadgeNeutral` or differentiate | P2 |
| 22 | Replace `Tab.Active` with `TabsProps.ActiveTabID` | P2 |
| 23 | Unify JS attachment pattern across Accordion/Dropdown/Modal | P2 |
| 24 | Extract shared dismiss JS for Alert and Toast | P2 |
| 25 | Make toast icon SVG paths single-source | P2 |
| 26 | Decouple `htmx/loading` from `feedback.Spinner` | P2 |

### Accessibility
| # | Task | Priority |
|---|---|---|
| 30 | Add `alt` text to Avatar `<img>` | P1 |
| 31 | Add `aria-required` to form inputs when `Required: true` | P1 |
| 32 | Add `<html lang>` to Base layout | P1 |
| 33 | Add Table header `scope` attributes | P2 |
| 34 | Add `aria-live="polite"` to HTMX loading indicators | P2 |
| 35 | Add `aria-live="polite"` to HTMX error handling | P2 |
| 36 | Fix `ErrorAttrs` for simultaneous error + help text | P2 |

### Testing Gaps
| # | Task | Priority |
|---|---|---|
| 42 | Add BDD tests for navigation package | P1 |
| 43 | Add BDD tests for htmx package | P1 |
| 44 | Add BDD tests for layout package | P1 |
| 45 | Add BDD tests for icons package | P2 |
| 46 | Add tests for Table mismatched header/row lengths | P2 |
| 47 | Add tests for Modal/Dropdown with empty ID | P2 |
| 48 | Add test for `mapStatusToBadgeType` boundary cases | P2 |
| 49 | Improve forms test coverage (58% → 75%+) | P2 |
| 50 | Improve utils test coverage (56% → 75%+) | P2 |

### Dead Code & Cleanup
| # | Task | Priority |
|---|---|---|
| 55 | Remove or use `icons.IconAttrs` (dead code) | P2 |
| 56 | Remove or use `internal/svg.FillIcon` | P2 |
| 57 | Remove no-op `DefaultXxxProps()` functions | P3 |
| 58 | Move test helpers out of `utils/` to `internal/testutil/` | P3 |
| 59 | Move `display/a11y_test.go` ProgressBar test | P3 |

### DevOps & Tooling
| # | Task | Priority |
|---|---|---|
| 62 | Release automation (goreleaser) | P3 |
| 63 | Fix pre-commit hook to be executable | P3 |
| 64 | Exclude `examples/` from lint | P3 |

### Documentation
| # | Task | Priority |
|---|---|---|
| 71 | Documentation site generation | P3 |
| 72 | Document `PageProps` not embedding `BaseProps` | P3 |

---

## d) TOTALLY FUCKED UP 💣

### NavLinkProps.Attrs Shadowing (TODO #9 — P0)
**This is a real bug.** `NavLinkProps` defines its own `Attrs templ.Attributes` field that shadows the embedded `BaseProps.Attrs`. Consumers who set `BaseProps.Attrs` on a `NavLinkProps` will have their attributes silently ignored. This violates the core contract that "all components propagate BaseProps".

### Dropdown JS XSS Vector (TODO #11 — P1)
The Dropdown's IIFE interpolates `props.ID` via Go string template. If `props.ID` contains malicious content, it could inject JS. `modal_go.go` already uses `strconv.Quote` for this exact reason — Dropdown should do the same.

### Accordion State Coupling (TODO #12 — P1)
Accordion uses `max-h-96` CSS class as a JS state indicator. This is fragile coupling between CSS and JS — changing the max-height styling breaks the JS toggle logic. Should use `data-open` attribute or `aria-expanded` as the state source.

### Demo App Security (FIXED this session)
The `renderSection()` function was writing unescaped `title` directly into HTML via `fmt.Fprintf(..., "%s", title)`. **Fixed** by wrapping with `templ.EscapeString(title)`.

---

## e) WHAT WE SHOULD IMPROVE

### 1. Stop Writing Raw HTML in Go Code
The demo app (`examples/demo/main.go`) writes HTML with `fmt.Fprintf` and `w.Write([]byte(...))`. This is:
- **XSS-prone** (as proven this session)
- **Unmaintainable** (198-line Go file with embedded HTML strings)
- **Defeats the purpose** of using templ

**Fix:** Convert to `.templ` files. Split `main.go` into `main.go` (HTTP server setup) + `demo.templ` (page layout + sections).

### 2. Make Impossible States Unrepresentable (Consistently)
We did this for `AvatarStatus` and `TrendDirection`, but:
- `TabsProps` still uses `Tab.Active bool` — allows zero or multiple active tabs
- `NavLinkProps.Attrs` shadows `BaseProps.Attrs` — silent data loss
- `DropdownProps`/`ModalProps` don't validate `ID` — broken ARIA with empty string

### 3. Unify JS Patterns
Three different JS attachment patterns exist:
- Accordion: Global flag check
- Dropdown: IIFE with string interpolation
- Modal: Global functions registered by ID

**All should use IIFE-per-instance** for CSP compliance and isolation.

### 4. Single-Source SVG Paths
Toast icon SVG paths are duplicated: once in Go (`toastIconPath` map) and once in JavaScript (`tcToastIcons` object). They can drift apart silently.

### 5. Test Coverage Gaps
- `forms`: 58% coverage
- `utils`: 56% coverage  
- No BDD tests for navigation, htmx, layout, icons packages
- No tests for edge cases (empty ID, mismatched table dimensions, unknown status mappings)

### 6. Public Release Readiness
The repo is 80% ready for public release. Missing:
- Alpha tag (`v0.1.0-alpha`)
- Pre-release badge in README
- Submission to `awesome-templ`

---

## f) Top 25 Things We Should Get Done Next

| # | Task | Priority | Effort | Impact |
|---|---|---|---|---|
| 1 | Fix `NavLinkProps.Attrs` shadowing `BaseProps.Attrs` (P0 bug) | P0 | Small | Trust |
| 2 | Fix Dropdown JS XSS vector (`strconv.Quote`) | P1 | Small | Security |
| 3 | Add `alt` text to Avatar `<img>` (WCAG 1.1.1) | P1 | Small | A11y |
| 4 | Add `aria-required` to form inputs when `Required: true` | P1 | Small | A11y |
| 5 | Add `<html lang>` to Base layout (WCAG 3.1.1) | P1 | Trivial | A11y |
| 6 | Validate required `ID` in Modal and Dropdown | P1 | Small | Robustness |
| 7 | Fix Accordion state coupling (`max-h-96` → `data-open`) | P1 | Small | Maintainability |
| 8 | Tag `v0.1.0-alpha` with pre-release disclaimer | P0 | Trivial | Release |
| 9 | Add "Alpha / Pre-release" badge to README | P1 | Trivial | Communication |
| 10 | Submit to `awesome-templ` | P2 | Trivial | Visibility |
| 11 | Convert demo app to `.templ` files | P2 | Medium | Dogfooding |
| 12 | Replace `Tab.Active` with `TabsProps.ActiveTabID` | P2 | Small | Type safety |
| 13 | Unify JS attachment pattern (IIFE-per-instance) | P2 | Medium | Consistency |
| 14 | Extract shared dismiss JS for Alert/Toast | P2 | Small | DRY |
| 15 | Make toast icon SVG paths single-source | P2 | Small | DRY |
| 16 | Consolidate badge color maps into single struct map | P2 | Small | DRY |
| 17 | Decouple `htmx/loading` from `feedback.Spinner` | P2 | Small | Coupling |
| 18 | Add BDD tests for navigation package | P1 | Medium | Coverage |
| 19 | Add BDD tests for layout package | P1 | Medium | Coverage |
| 20 | Add BDD tests for htmx package | P1 | Medium | Coverage |
| 21 | Improve forms test coverage (58% → 75%+) | P2 | Medium | Coverage |
| 22 | Remove dead code (`icons.IconAttrs`, no-op constructors) | P2 | Small | Cleanliness |
| 23 | Fix `ErrorAttrs` for simultaneous error + help text | P2 | Small | Correctness |
| 24 | Add Table header `scope` attributes | P2 | Small | A11y |
| 25 | Build component showcase website | P3 | Large | Adoption |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should we release `v0.1.0-alpha` now with the known P0/P1 bugs unfixed, or fix them first?**

Arguments for releasing now:
- The repo is public-ready (CI, docs, README, CONTRIBUTING, MIT license)
- Alpha tag sets expectations that bugs exist
- First-mover advantage in the "pure Tailwind, pure Go" niche is time-sensitive
- Community feedback on what matters most is more valuable than our assumptions

Arguments for fixing first:
- `NavLinkProps.Attrs` shadowing is a silent data loss bug (P0)
- Dropdown JS XSS is a security issue (P1)
- First impressions matter — early adopters who hit these bugs may not return
- Only 4-5 items, estimated 1-2 hours of work

**I recommend fixing the P0 and P1 items (items 1–7 in the top 25 list) before tagging alpha.** But this is your call.
