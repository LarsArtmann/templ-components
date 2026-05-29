# Status Report — templ-components

**Date:** 2026-05-29 13:15 | **Branch:** master | **Since:** Last commit `7a34b02` (2026-05-27)

---

## Executive Summary

Added `errorpage` package — 3 new components for presenting structured errors on the web,
designed as a companion to the `go-error-family` library. All 940+ tests pass. Lint clean on
new code. Build green. Ready for commit.

---

## A) FULLY DONE

### New `errorpage` Package (2,105 lines total)

| File | Lines | Purpose |
|---|---|---|
| `errorpage/styles.go` | 188 | Family enum (5 families), visual style mappings, props structs, `ContextPair`, `CauseItem` types |
| `errorpage/errorpage.templ` | 128 | Full-page error view with Wix-style What/Why/Fix/WayOut layout |
| `errorpage/errordetail.templ` | 109 | Inline error detail card with context table, cause chain, suggested fix |
| `errorpage/erroralert.templ` | 75 | Family-aware alert banner with dismiss support |
| `errorpage/doc.go` | 26 | Package documentation with usage examples |
| `errorpage/bdd_test.go` | 318 | 25 BDD-style tests (all passing) |
| `errorpage/*_templ.go` (3 files) | 1,261 | Generated templ output (committed per library convention) |

**Family → Visual Mapping:**

| Family | Color | Icon | Tone |
|---|---|---|---|
| Rejection | Amber | ExclamationTriangle | Instructional |
| Conflict | Orange | ExclamationCircle | Explanatory |
| Transient | Blue | Refresh | Reassuring |
| Corruption | Red | ExclamationTriangle | Urgent |
| Infrastructure | Slate | Globe | Apologetic |

**Design decisions:**
- Zero dependency on go-error-family — bridge via `errorpage.Family(err.ErrorFamily().String())`
- Follows all project conventions: `BaseProps`, `utils.Class()`, `DefaultXxxProps()`, map lookups, CSP nonce
- Distinct 5-color palette (amber/orange/blue/red/slate) — each family is immediately distinguishable
- Supports go-error-family's What/Why/Fix/WayOut message template structure

### Regenerated `*_templ.go` Files

All 35 `*_templ.go` files regenerated with latest templ generator (v0.3.1036 vs v0.3.1020 in go.mod).
Minor formatting diffs in existing generated files — no behavioral changes.

### Updated AGENTS.md

Added errorpage to: import graph, generated file count, lint command, architecture section.

---

## B) PARTIALLY DONE

Nothing — all planned work for this session is complete.

---

## C) NOT STARTED

### Integration Examples
- No `examples/demo` entry for errorpage components yet
- No example showing go-error-family → errorpage bridge code

### Documentation
- No README update mentioning errorpage
- No CHANGELOG entry for the new package
- No ADR documenting the zero-dependency design decision

### Edge Case / Accessibility Testing
- No `a11y_test.go` for errorpage (existing packages have this)
- No `edge_cases_test.go` for errorpage (existing packages have this)
- No `snapshot_test.go` for errorpage (existing packages have this)

### Example Test File
- No `example_test.go` for errorpage (existing packages have this)

---

## D) TOTALLY FUCKED UP

### Pre-existing Lint Issues (NOT introduced this session)
- `feedback/styles.go:9` — unused `//nolint:revive` directive (nolintlint)
- `navigation/snapshot_test.go:205` — unused `//nolint:gochecknoglobals` directive (nolintlint)

These are in files we did NOT touch. Not our responsibility per project rules.

### Go Module Workspace Quirk
- `/home/lars/projects/go.work` exists at repo root level and does NOT include templ-components
- All `go build`/`go test` commands require `GOWORK=off` to work
- AGENTS.md build commands don't mention this — could confuse contributors

---

## E) WHAT WE SHOULD IMPROVE

### Critical
1. **AGENTS.md missing `GOWORK=off`** — Build commands in AGENTS.md don't mention the workspace issue. Anyone running `go build ./...` without `GOWORK=off` gets a confusing error.

### High Impact
2. **No example/demo for errorpage** — The demo is the best marketing. An interactive demo showing all 5 families would be compelling.
3. **Test coverage 65.6%** for errorpage — lowest in the project. Missing a11y tests, edge case tests, snapshot tests, and example tests that other packages have.
4. **No integration bridge example** — Consumers need to see the go-error-family → errorpage conversion pattern. A `// Example` in doc.go is not enough.

### Medium Impact
5. **Missing errorpage in FEATURES.md** — New package not documented in feature inventory.
6. **Missing errorpage in CHANGELOG.md** — No changelog entry for the new package.
7. **templ version drift** — Generator is v0.3.1036 but go.mod has v0.3.1020. AGENTS.md warns about this.

### Low Impact
8. **Pre-existing nolintlint warnings** — 2 unused nolint directives in untouched files.
9. **No ADR for zero-dependency design** — The decision to mirror go-error-family types instead of importing it deserves documentation.

---

## F) Top #25 Things We Should Get Done Next

### Tier 1: Ship-Blocking (do immediately)
1. **Add errorpage to the demo** — Interactive demo showing all 5 families with ErrorPage, ErrorDetail, ErrorAlert
2. **Write a11y_test.go for errorpage** — Role, aria-label, aria-live, focus management tests
3. **Write edge_cases_test.go for errorpage** — Unknown family, empty props, nil context, long strings
4. **Write snapshot_test.go for errorpage** — HTML snapshot tests for all 3 components × 5 families

### Tier 2: Polish (do before next release)
5. **Write example_test.go for errorpage** — GoDoc-ready examples for each component
6. **Add CHANGELOG.md entry** — Document the new errorpage package
7. **Update FEATURES.md** — Add errorpage to the feature inventory
8. **Fix AGENTS.md build commands** — Add `GOWORK=off` note or add templ-components to go.work
9. **Add integration example** — Show go-error-family → errorpage bridge in examples/
10. **Write ADR for zero-dependency design** — Document why errorpage mirrors types instead of importing go-error-family

### Tier 3: Quality (do soon)
11. **Upgrade templ in go.mod** — v0.3.1020 → v0.3.1036 to match generator
12. **Fix pre-existing nolintlint warnings** — 2 unused nolint directives
13. **Add errorpage to CI lint command** — Ensure `.github/workflows/ci.yaml` includes errorpage
14. **ErrorPage layout component integration** — ErrorPage should work inside layout.Base for full-page rendering
15. **ErrorPage HTTP status code mapping** — Add `FamilyStatusCode()` mapping (400/409/503/500/503) for HTTP handlers

### Tier 4: Enhancements
16. **ErrorPage dark mode visual verification** — Manual/snapshot check of all 5 families in dark mode
17. **ErrorAlert animation** — Add enter/exit transitions like Toast has
18. **ErrorDetail collapsible sections** — Make context/cause chain toggleable (currently always shown)
19. **ErrorPage background pattern** — Add subtle SVG pattern or gradient background for visual distinction
20. **Multi-error display** — Component to render `errors.Join` results as a list of ErrorDetails
21. **ErrorPage SEO meta** — Add `<title>` and `<meta name="robots">` support for HTTP error pages

### Tier 5: Ecosystem
22. **go-error-family `errorpage` bridge package** — A tiny package `github.com/larsartmann/go-error-family/errorpage` that converts `*Error` → `ErrorPageProps` in one call
23. **HTMX error page swap** — Component that uses HTMX to swap error pages without full reload
24. **ErrorPage analytics hooks** — `OnErrorShown` callback for tracking error impressions
25. **Error page templates** — Pre-built error pages for common HTTP codes (404, 500, 403, 503) with sensible defaults

---

## G) Top #1 Question I Cannot Figure Out Myself

**Should the `errorpage` package import `go-error-family` or remain zero-dependency?**

Arguments for zero-dependency (current design):
- No coupling, no version drift, works with any error classification system
- Parallel types are trivial to bridge
- Keeps templ-components dependency tree minimal (only templ + tailwind-merge-go)

Arguments for importing:
- Auto-conversion functions (`FromError(*errorfamily.Error) ErrorPageProps`)
- Type safety guarantees (no string-matching drift)
- Single source of truth for family names

A possible middle ground: create a **separate bridge package** `errorpage/bridge` or `errorpage/errorfamily` that has the go-error-family dependency, keeping the core zero-dep. But this requires a go.work multi-module setup or a separate repo.

This is a product/architecture decision only you can make.

---

## Metrics

| Metric | Value |
|---|---|
| Packages | 10 + demo |
| Source `.templ` files | 38 |
| Generated `*_templ.go` files | 35 |
| Test files | 44 |
| Total test cases (PASS) | 940+ |
| Overall coverage (weighted) | ~71% |
| New package coverage (errorpage) | 65.6% |
| Highest coverage (utils) | 83.3% |
| Lowest coverage (forms) | 64.3% |
| Lint issues (new code) | 0 |
| Lint issues (pre-existing) | 2 (nolintlint, untouched files) |
| Build status | GREEN |
| Test status | GREEN |
| Lines added this session | ~2,105 (errorpage) + ~100 (AGENTS.md, generated diffs) |

---

## Test Coverage Per Package

| Package | Coverage |
|---|---|
| utils | 83.3% |
| internal/svg | 79.0% |
| htmx | 77.3% |
| icons | 75.0% |
| layout | 73.2% |
| navigation | 72.1% |
| feedback | 70.3% |
| display | 68.6% |
| **errorpage** | **65.6%** |
| forms | 64.3% |
| examples/demo | 0.0% (no tests) |
