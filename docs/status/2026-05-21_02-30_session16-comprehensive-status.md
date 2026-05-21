# Comprehensive Status Report

**Date:** 2026-05-21 02:30  
**Session:** 16 — TODO cleanup, bug fixes, examples  
**Branch:** master  
**Coverage:** 71.8% | **Tests:** 182 | **Lint:** 0 issues

---

## a) FULLY DONE

### TODO Cleanup (~70 items removed)
- Reviewed entire TODO_LIST.md (243 items → 112 remaining)
- Marked ~70 items as completed based on code verification
- Reorganized into 7 priority categories: Bugs, Breaking Changes, Enhancements, Testing, Infrastructure, Documentation, Release
- Removed stale items (e.g., BoolString removal, Deref removal, BadgeDefault removal — all already done)

### Bug Fixes
1. **Demo app HTMX** — Removed `props.HTMXVersion = ""` override in `examples/demo/main.go` so demo now uses default HTMX 2.0.6
2. **Checkbox conditional ID** — Added `if props.ID != ""` guard before `id={ props.ID }` in `forms/input.templ` (matches Input component pattern)
3. **ThemeToggle IIFE removal** — Removed unnecessary IIFE wrapper from `layout/theme.templ` for consistency with other components' JS patterns
4. **JS re-attachment research** — Verified that document-level event delegation correctly handles HTMX-swapped elements. Global singleton guards are the correct pattern for delegation; marked 3 related TODO items as "by design"

### Infrastructure
- **Pre-commit hook** — Replaced `buildflow` dependency with `scripts/pre-commit.sh` that runs: `templ generate → go build → go test → golangci-lint`
- **AGENTS.md** — Updated coverage from 68.3% → 71.8%

### Documentation
- **Example functions** — Added `ExampleXxx()` functions for pkg.go.dev discoverability:
  - `display`: ExampleBadge, ExampleCard, ExampleStatCard
  - `feedback`: ExampleAlert, ExampleToast, ExampleSpinner
  - `icons`: ExampleIcon
  - `navigation`: ExamplePagination, ExampleNavLink, ExampleBreadcrumbs

---

## b) PARTIALLY DONE

- **templ CLI version mismatch** — CLI v0.3.1001 vs go.mod v0.3.1020. Warning on every generate but not blocking.

---

## c) NOT STARTED

112 items remaining in TODO_LIST.md. Top categories:
- **New Components** (14): Radio, Toggle/Switch, File input, Date Picker, Combobox, Dialog/Drawer, Form wrapper, Skeleton, Step indicator vertical, Badge click/href, more Heroicons, ProgressBar indeterminate, client-side tab switching, tab keyboard nav
- **Testing** (20): BDD tests for navigation/htmx/layout/icons, coverage improvements, golden file tests, composition tests, integration tests, nonce audit
- **Infrastructure** (10): GitHub Actions CI, goreleaser, coverage threshold, Go workspace modularization, demo deployment, nix flake
- **Documentation** (8): README updates, CONTEXT.md updates, ADRs, go doc examples (partially started), DOMAIN_LANGUAGE.md
- **A11y** (7): aria-live in HTMX error handling, inline JS consolidation, Table scope, EmptyState landmark, Breadcrumb JSON-LD, Pagination rel=prev/next
- **Validation** (8): TrendDirection, SelectOption, SwapOOB, SVG pipe separator, DropdownItem.Disabled, MaxLength fields

---

## d) TOTALLY FUCKED UP!

Nothing. Build passes, tests pass, lint is clean.

---

## e) WHAT WE SHOULD IMPROVE!

1. **README.md stale** — Doesn't mention `AvatarStatus`, `StatCardProps`, `BreadcrumbsProps`, or any v0.2 API changes. Component catalog needs updating.
2. **Generated file churn** — Every `templ generate` touches ~30 files due to import block normalization. Consider pinning templ CLI version in CI.
3. **No CI pipeline** — No GitHub Actions. Every check is manual.
4. **No release tags** — Public repo with no tags. `go get` pulls floating commits.
5. **Example functions incomplete** — Only 4 packages have examples. Missing: forms, htmx, layout, utils.
6. **Coverage gap** — Several packages below 75%: display (71.8%), feedback (72.8%), forms (70.8%), navigation (72.2%).

---

## f) Top #25 Things to Get Done Next

1. Set up GitHub Actions CI (build + test + lint + coverage threshold)
2. Update README.md with v0.2 API changes and component examples
3. Tag v0.2.0 release + CHANGELOG.md
4. Add remaining ExampleXxx() functions (forms, htmx, layout, utils)
5. Fix inline JS consolidation — shared init strategy in base layout
6. Add `DropdownItem.Disabled` field
7. Add `InputProps.MaxLength` / `TextareaProps.MaxLength` / `CheckboxProps.Value`
8. BDD tests for navigation package (Nav, Pagination, Breadcrumbs)
9. BDD tests for htmx package (Loading, ErrorHandling)
10. BDD tests for layout package (Base, Minimal, Theme)
11. Add component composition integration tests (Card+Badge, Nav+Avatar, etc.)
12. Improve forms coverage (70.8% → 75%+)
13. Improve display coverage (71.8% → 75%+)
14. Add Table `scope` attributes for screen readers
15. Add EmptyState `role="region"` landmark
16. Extract HTMX CDN URL constant (repeated in base.templ + sri.go)
17. Add `aria-live="polite"` directly to HTMX error handling script
18. Set up goreleaser for automated releases
19. Add Radio button component
20. Add Toggle/Switch component
21. Write ADR for JS attachment patterns (singleton vs IIFE)
22. Write ADR for FeedbackType unification
23. Document `utils.Class()` thread-safety requirement
24. Add dark mode output verification tests
25. Verify `go get` works from clean project

---

## g) Top Question I Cannot Figure Out

**Is there an actual bug with the global `window.tc*Attached` pattern?**

All 5 components (Accordion, Dropdown, Alert, Toast, ThemeToggle) use document-level event delegation. The global singleton prevents duplicate listeners when multiple instances exist. Event delegation automatically handles dynamically added elements (including HTMX swaps). I cannot identify a concrete failure scenario where the current pattern breaks. The TODO items calling this a "bug" may be based on a misunderstanding of event delegation. Should these be removed entirely, or is there a specific edge case I'm missing?
