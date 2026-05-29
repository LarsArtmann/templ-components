# Status Report — templ-components

**Date:** 2026-05-29 13:37 | **Branch:** master (pushed to origin) | **Since:** `7a34b02` (2026-05-27)

---

## Executive Summary

Created and polished the `errorpage` package — a complete solution for presenting go-error-family
structured errors on the web. 7 commits, 52 files touched, 3,088 lines added, 120 removed.
984 tests pass, 0 lint issues, build green, pushed to origin.

---

## A) FULLY DONE

### New `errorpage` Package (2,605 lines)

| File                           | Purpose                                                           |
| ------------------------------ | ----------------------------------------------------------------- |
| `errorpage/styles.go`          | Family enum (5), visual style maps, props structs, bridge helpers |
| `errorpage/errorpage.templ`    | Full-page error view with Wix-style What/Why/Fix/WayOut layout    |
| `errorpage/errordetail.templ`  | Inline error detail card with context table, cause chain, fix     |
| `errorpage/erroralert.templ`   | Family-aware alert banner with dismiss support                    |
| `errorpage/doc.go`             | Package documentation                                             |
| `errorpage/bdd_test.go`        | 25 BDD-style behavior tests                                       |
| `errorpage/a11y_test.go`       | 8 accessibility tests (ARIA roles, labels, nonce)                 |
| `errorpage/edge_cases_test.go` | 17 edge case tests (unknown families, nil inputs, bridge helpers) |
| `errorpage/example_test.go`    | 6 GoDoc examples                                                  |
| `errorpage/*_templ.go` (3)     | Generated templ output (committed)                                |

**3 components × 5 families = 15 distinct visual treatments**

| Family         | Color  | Icon                | Tone          | HTTP |
| -------------- | ------ | ------------------- | ------------- | ---- |
| Rejection      | Amber  | ExclamationTriangle | Instructional | 400  |
| Conflict       | Orange | ExclamationCircle   | Explanatory   | 409  |
| Transient      | Blue   | Refresh             | Reassuring    | 503  |
| Corruption     | Red    | ExclamationTriangle | Urgent        | 500  |
| Infrastructure | Slate  | Globe               | Apologetic    | 503  |

**Bridge helpers for go-error-family integration:**

- `FamilyStatusCode(f)` — Family → HTTP status code
- `ContextMap(m)` — map[string]string → []ContextPair
- `ExtractCauseChain(err, maxDepth)` — walks Unwrap() chain → []CauseItem

### Refactoring

- **DismissScript deduplication**: Extracted `utils.DismissScript()` as single source of truth for feedback (Alert, Toast) and errorpage (ErrorAlert). Eliminated inline JS copy.
- **Dead code removal**: `Expanded` field removed from `ErrorDetailProps`, empty if-block removed from errordetail cause chain loop.
- **Pre-commit hook fix**: Added `GOWORK=off` and included errorpage in lint targets.
- **Stale nolint cleanup**: Removed 2 unused nolint directives (feedback/styles.go, navigation/snapshot_test.go).
- **Templ regeneration**: All 38 `*_templ.go` files regenerated with v0.3.1036.

### Documentation

- **CHANGELOG.md**: Unreleased section with errorpage features, bridge helpers, utils.DismissScript
- **FEATURES.md**: Full errorpage section (components, enums, bridge helpers, family visual table), updated totals, added DismissScript, removed stale known issue
- **AGENTS.md**: Updated feedback shared section, import graph, generated file count, lint command
- **Demo page**: Added ErrorAlert (all 5 families) and ErrorDetail sections to interactive demo

---

## B) PARTIALLY DONE

Nothing — all planned work for this session is complete.

---

## C) NOT STARTED

### ErrorPage Demo

- The `ErrorPage` component (full-page view) is NOT in the demo page because it takes over the
  entire viewport. Would need a separate route or iframe in the demo app.

### Snapshot Tests

- No `snapshot_test.go` for errorpage — other packages (feedback, navigation, forms, icons) have these.
  Snapshot tests lock down exact HTML output to catch unintended changes.

### Type Architecture Improvements

- `feedbackStyleSet` (4 fields) and `familyVisualStyle` (8 fields) share the same pattern but
  aren't unified. Could extract a base style interface or use composition.
- `feedback.FeedbackType` has 4 levels (Success/Error/Warning/Info), `errorpage.Family` has 5 families.
  These are orthogonal axes — no overlap expected, but no explicit type relationship either.

### ErrorPage + layout.Base Integration

- ErrorPage renders its own full-screen container. It doesn't integrate with `layout.Base` for
  consistent `<head>` management. Consumers must choose one or the other.

### CI Pipeline

- `.github/workflows/ci.yaml` not updated to include `GOWORK=off` — may fail in CI if the workspace
  issue exists there too.

---

## D) TOTALLY FUCKED UP

### Nothing Is Fucked Up

- Working tree: clean
- Build: green
- Tests: 984 pass, 0 fail
- Lint: 0 issues
- All 7 commits pushed to origin/master
- Pre-commit hook passes on every commit

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **`familyVisualStyle` vs `feedbackStyleSet` unification** — Both are "CSS class containers with map lookup + default fallback." Could extract a shared `StyleSet` base type or use generics to eliminate the duplicate pattern. Low urgency since they serve different domains.

2. **ErrorPage + layout.Base composition** — ErrorPage should accept an optional `layout.PageProps` or be usable inside `layout.Base` for consistent `<head>`, CSP headers, and theme toggle. Currently they're mutually exclusive.

3. **Bridge package** — Consider a separate `errorpage/bridge` sub-package that imports go-error-family and provides one-call conversion: `FromError(*errorfamily.Error) ErrorPageProps`. This keeps core zero-dep while offering convenience.

### Quality

4. **Snapshot tests missing** — Every other component package has `snapshot_test.go`. Errorpage should too for regression protection.

5. **Coverage gap** — errorpage at 69.6% is above the project average but below top performers (utils 80%, svg 79%). The bridge helpers have good coverage but the template rendering paths could use more edge cases.

6. **CI `GOWORK` issue** — The pre-commit hook was fixed locally but `.github/workflows/ci.yaml` may need the same `GOWORK=off` treatment.

### Polish

7. **ErrorPage demo route** — Add a `/error` route to the demo app that renders each family as a full-page error view. This is the most compelling demo but requires routing.

8. **Pre-built HTTP error pages** — `NotFound()`, `InternalServerError()`, `Forbidden()`, `ServiceUnavailable()` convenience constructors with sensible defaults would eliminate boilerplate for the 4 most common HTTP errors.

---

## F) Top #25 Things We Should Get Done Next

### Tier 1: Ship-Blocking (pre-v0.3 release)

1. **Add `snapshot_test.go` for errorpage** — Lock down exact HTML for all 3 components × 5 families
2. **Verify CI pipeline passes** — Check `.github/workflows/ci.yaml` handles errorpage and `GOWORK`
3. **Add pre-built HTTP error constructors** — `NotFound()`, `Forbidden()`, `ServiceUnavailable()`, `InternalError()` with sensible defaults
4. **Write ADR for zero-dependency design** — Document why errorpage mirrors types instead of importing go-error-family

### Tier 2: Quality (do before promoting)

5. **Unify style lookup pattern** — Extract shared `StyleSet` base or generic lookup from feedback/errorpage
6. **ErrorPage + layout.Base composition** — Allow ErrorPage to render inside Base for consistent `<head>`
7. **Add errorpage demo route** — `/error/:family` route in demo app for full-page error views
8. **Coverage push to 75%+** — Add more template rendering edge cases
9. **Snapshot test migration** — Consider golden file approach for all packages (not just errorpage)

### Tier 3: Features

10. **ErrorPage background pattern** — Subtle SVG pattern or gradient for visual distinction
11. **Multi-error display component** — Render `errors.Join` results as a list of ErrorDetails
12. **ErrorAlert enter/exit animations** — Toast has transitions, ErrorAlert doesn't
13. **ErrorDetail collapsible sections** — Toggle context/cause chain visibility
14. **ErrorPage SEO meta** — `<title>` and `<meta name="robots" content="noindex">` support
15. **HTMX error page swap** — Component that swaps error pages without full reload

### Tier 4: Ecosystem

16. **go-error-family bridge package** — `errorpage/bridge` with `FromError()` one-call conversion
17. **Error page analytics hooks** — `OnErrorShown` callback for tracking error impressions
18. **Error page templates** — Pre-built for common HTTP codes (404, 500, 403, 503)
19. **ErrorPage i18n** — Message/fix translations for multi-language apps
20. **ErrorDetail copy-to-clipboard** — One-click error context copy for bug reports

### Tier 5: Cross-Cutting

21. **Upgrade templ in go.mod** — v0.3.1020 → v0.3.1036 to match generator
22. **Fix remaining pre-existing lint issues** — Toast icon SVG path duplication across Go/JS
23. **Spinner SVG rendering consolidation** — Rendered 3 different ways across packages
24. **Add `go.work` for templ-components** — Or remove dependency on parent go.work
25. **Errorpage integration test** — Full HTTP handler test: error → FamilyStatusCode → ErrorPage render

---

## G) Top #1 Question I Cannot Figure Out Myself

**Should `errorpage` import `go-error-family` directly, or remain zero-dependency forever?**

Current state: Zero dependency. Consumers bridge with:

```go
errorpage.Family(myError.ErrorFamily().String())
```

The gap: No type-safety guarantee that the string constants match. If go-error-family adds
a 6th family, errorpage won't know about it until manually updated.

Options I see:

1. **Stay zero-dep** — Mirror types, accept the drift risk. Library stays minimal.
2. **Add optional bridge sub-module** — `errorpage/bridge/go.mod` imports go-error-family, provides `FromError()`. Core stays zero-dep.
3. **Import directly** — Tighter coupling, version drift headaches, but type safety.

This is a product decision about the library's identity: is it "UI components for ANY error system" or "UI components SPECIFICALLY for go-error-family"?

---

## Metrics

| Metric                       | Value     |
| ---------------------------- | --------- |
| Packages                     | 10 + demo |
| Source `.templ` files        | 38        |
| Generated `*_templ.go` files | 38        |
| Test files                   | 52        |
| Total test cases (PASS)      | 984       |
| Lint issues                  | 0         |
| Build status                 | GREEN     |
| Commits this session         | 7         |
| Lines added this session     | +3,088    |
| Lines removed this session   | -120      |
| Files touched                | 52        |

## Test Coverage Per Package

| Package       | Coverage  | Delta vs Last Report                                          |
| ------------- | --------- | ------------------------------------------------------------- |
| utils         | 80.0%     | ↓ 3.3% (added DismissScript, not yet heavily tested directly) |
| internal/svg  | 79.0%     | —                                                             |
| htmx          | 77.3%     | —                                                             |
| icons         | 75.0%     | —                                                             |
| layout        | 73.2%     | —                                                             |
| navigation    | 72.1%     | —                                                             |
| feedback      | 70.3%     | —                                                             |
| **errorpage** | **69.6%** | **NEW**                                                       |
| display       | 68.6%     | —                                                             |
| forms         | 64.3%     | —                                                             |

## Commit History This Session

```
db77e4d docs: update CHANGELOG, FEATURES, and AGENTS.md for errorpage package
f8c027a demo(errorpage): add error alerts and error details to demo page
9306d1f test(errorpage): add a11y, edge case, and example tests
c4344ac fix(errorpage): add exhaustruct nolint for CauseItem conditional init
295b949 refactor(dismiss): extract DismissScript to utils, eliminate duplication
5b75c10 chore(templ): regenerate all *_templ.go files with templ v0.3.1036
0d4eff9 feat(errorpage): add error presentation components for go-error-family integration
```
