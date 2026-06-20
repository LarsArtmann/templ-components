# Status Report — templ-components

**Date:** 2026-06-20 16:15
**Branch:** master
**HEAD:** `9eaf94e` — refactor: replace MapEnum with Lookup[K,V] generic

---

## Metrics Snapshot

| Metric           | Value                                                           |
| ---------------- | --------------------------------------------------------------- |
| Tests            | 1700+ (all passing)                                             |
| Packages         | 11 + demo + integration                                         |
| Source files     | 37 `.templ` + 12 `.go`                                          |
| Generated files  | 46 (`*_templ.go`)                                               |
| Icons            | 99 (98 path + 1 Spinner)                                        |
| Components       | 55                                                              |
| Lint issues      | 0                                                               |
| Average coverage | ~74.5%                                                          |
| Dependencies     | 3 direct (templ, tailwind-merge-go, go-error-family)            |
| Runtime panics   | 0 in component code (1 developer data integrity check in icons) |

### Per-Package Coverage

| Package         | Coverage        |
| --------------- | --------------- |
| utils           | 79.2%           |
| internal/svg    | 79.0%           |
| htmx            | 77.3%           |
| internal/golden | 76.9%           |
| icons           | 76.2%           |
| layout          | 73.6%           |
| errorpage       | 73.7%           |
| forms           | 73.0%           |
| feedback        | 73.2%           |
| display         | 72.6%           |
| navigation      | 72.6%           |
| examples/demo   | 0.0% (no tests) |

---

## A) COMPLETED THIS SESSION

### Panic Elimination

- [x] Modal, Drawer, Dropdown: auto-generate IDs via `utils.EnsureID()` (crypto/rand)
- [x] Accordion: auto-generate IDs for items missing them
- [x] SwapOOB: invalid swap styles fall back to `outerHTML`
- [x] **Zero runtime panics** in component code

### New Components

- [x] `forms.DatePicker` — native HTML date input with min/max constraints
- [x] `forms.Combobox` — accessible autocomplete with client-side filtering, hidden value for form submission, global singleton JS, CSP nonce

### Architecture Improvements

- [x] `utils.Lookup[K comparable, V any]` — generic map lookup replacing 14× duplicated if-ok-return boilerplate (old `MapEnum` was too narrow: only handled `map[string]T where T~string`)
- [x] `utils.EnsureID(prefix, id)` — crypto/rand-based ID auto-generation
- [x] `utils.RenderAll()` — test helper for multi-component integration tests
- [x] BadgeInfo color aligned from `indigo-*` to `blue-*` for consistency

### Testing

- [x] Coverage boosters for display, errorpage, feedback, forms, navigation, layout
- [x] Errorpage handler tests (HTMLShell, JSON, Override, WriteError/WriteErrorPage)
- [x] Integration test package verifying cross-package composition
- [x] In-package tests for `DismissScript`, `RenderAll`, `EnsureID` (fixed coverage regression)

### Documentation

- [x] AGENTS.md, CHANGELOG.md, README.md updated

---

## B) WHAT I FORGOT / COULD HAVE DONE BETTER

### Honest Self-Criticism

1. **utils coverage regression (68.2%)** — I added `DismissScript`, `RenderAll`, `EnsureID` but only tested them from OTHER packages. Go's `-cover` only instruments tests IN the package under test. Should have written in-package tests from the start.

2. **MapEnum was dead code for 14 of 15 call sites** — The helper existed but its signature `map[string]T where T~string` was too narrow. 14 sites duplicated the pattern because MapEnum couldn't handle struct values or typed keys. I should have caught this in the first self-review.

3. **Committed generated files were stale** — The pre-commit hook's templ generator (v0.3.1036) produced different import grouping than what was committed. I needed to commit the regenerated output.

4. **Combobox keyboard accessibility is incomplete** — No ArrowDown/Up navigation within listbox, no `aria-activedescendant` tracking. The current implementation handles basic filtering and click-to-select but doesn't meet the full WAI-ARIA combobox pattern.

5. **No Validate() error on props** — ErrorPageProps has 12 optional fields. `Message` is semantically required but nothing enforces it. Graceful fallbacks hide configuration errors.

---

## C) STILL TO IMPROVE

### High Impact, Low Effort

1. **Tag v0.3.0** — all features done, just needs `git tag`
2. **Submit to awesome-templ** — discoverability
3. **Update FEATURES.md** — stale (says 42 icons, 6 form components)

### Medium Impact, Medium Effort

4. **Combobox keyboard nav** — ArrowDown/Up, Enter, Home/End, `aria-activedescendant`
5. **Get coverage to 80%** — templ render functions have a ceiling around 75% due to boilerplate error handling
6. **Property-based tests** — verify invariants across random prop combinations

### High Impact, High Effort

7. **Deploy demo site** — GitHub Pages from examples/demo
8. **Documentation site** — pkgsite or doc2go
9. **Plan v1.0 API freeze** — define scope, set date

---

## D) ARCHITECTURE ASSESSMENT

### Type Model

- **BaseProps embedding**: sound — every props struct gets ID, Class, Attrs, AriaLabel, Nonce for free
- **ComponentProps interface**: exists but under-leveraged — no generic composition utilities use it yet
- **Enum types**: all use `type XxxType string` + const + `utils.Lookup` map pattern — consistent
- **ErrorPageProps**: 12 optional fields is a smell. Could benefit from required-field constructor or builder pattern. But graceful fallbacks (FamilyTransient default, empty Message renders without error) mean this is a DX issue, not a correctness issue.

### Library Hygiene

- **Zero runtime panics** in component code
- **CSP nonce propagation**: audited — all 14 inline scripts use `nonce={ props.Nonce }`
- **HTMX-safe event delegation**: all click handlers use document/container-level delegation
- **Generated files committed**: 46 `*_templ.go` files, required for library consumers
- **Import graph**: clean, no circular deps (`utils ← all`, `internal/svg ← display/feedback/icons`)

### Established Libraries

- **templ v0.3.1020**: core dependency, works well
- **tailwind-merge-go v0.2.1**: class conflict resolution, thread-safe via mutex
- **go-error-family v0.4.0**: error classification, integrated via errorpage
- **encoding/json v1**: stays — json/v2 requires GOEXPERIMENT flag in Go 1.26.3, can't require from a library until Go 1.27
