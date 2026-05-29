# Comprehensive Status Report — Deep Self-Review

**Date:** 2026-05-29 14:28 CEST | **Branch:** master | **Commits ahead:** 0 (all pushed) | **Working tree:** clean

---

## a) FULLY DONE

### Errorpage Package (3 sessions of work)

| Item | Status | Details |
|---|---|---|
| `ErrorPage` component | DONE | Full-page error view with What/Why/Fix/WayOut |
| `ErrorDetail` component | DONE | Inline error card with context table, cause chain |
| `ErrorAlert` component | DONE | Family-aware alert banner with dismiss |
| 5 Family enum constants | DONE | Rejection/Conflict/Transient/Corruption/Infrastructure |
| Visual style maps | DONE | Colors, icons, tones per family |
| Bridge helpers | DONE | `FamilyStatusCode()`, `ContextMap()`, `ExtractCauseChain()` |
| Shared sub-templates | DONE | 6 private templates in `shared.templ` |
| `http.Handler` integration | DONE | `ErrorHandler`, `WriteError`, `WriteErrorPage` |
| Pre-built constructors | DONE | 6 constructors with code constants |
| `FromError()` bridge | DONE | Extracts code/family/context/cause from any error |
| HTMX family-aware errors | DONE | `GlobalErrorHandling` parses structured JSON errors |
| BDD tests (25) | DONE | Behavior-driven test coverage |
| A11y tests (8) | DONE | Accessibility assertions |
| Edge case tests (17) | DONE | Nil/empty/invalid inputs |
| Example tests (6) | DONE | GoDoc examples |
| Handler tests (21) | DONE | HTTP handler, constructors, FromError |
| Demo page integration | DONE | ErrorAlert (all 5 families) + ErrorDetail (2 examples) |
| Documentation | DONE | AGENTS.md, CHANGELOG.md, FEATURES.md, doc.go |
| DismissScript extraction | DONE | Moved from feedback to `utils.DismissScript()` |
| Build/Test/Lint | DONE | 1027+ tests, 0 lint issues, clean build |

### Overall Project Health

| Metric | Value |
|---|---|
| Total tests | 1027+ |
| Packages | 10 + demo |
| Generated `*_templ.go` | 40 files |
| Lint issues | 0 |
| Build | GREEN |
| Dependencies | 2 direct (templ, tailwind-merge-go) |

---

## b) PARTIALLY DONE

### `FromError()` — Broken go-error-family Bridge

**THE BIGGEST ISSUE:** `FromError()` checks for `ErrorFamily() Family` where `Family` is `errorpage.Family` (string type). But go-error-family's `Classified` interface returns `errorfamily.Family` (int type). **A real go-error-family error will NEVER satisfy this assertion.** The bridge always falls back to `FamilyTransient`.

- The `ErrorCode()` and `ErrorContext()` bridges work correctly (they use string interfaces that both packages agree on)
- Only the `ErrorFamily()` bridge is broken due to type mismatch
- Tests pass because they use a mock that returns `errorpage.Family` — not a real go-error-family error

### HTMX Family-Aware Error Handling — Ghost System

The HTMX `GlobalErrorHandling` JS parses JSON error bodies with `family`/`code`/`message` fields. But `ErrorHandler()` renders HTML, not JSON. No code in the codebase produces the JSON format the HTMX handler expects. This is a **ghost system** — well-written JS code with no producer.

### Dark Mode Color Consistency

Errorpage uses mixed `gray-*` and `slate-*` in dark mode. The `familyStyleDefault` (unknown family fallback) uses `gray-*` while all 5 defined families correctly use `slate-*` for Infrastructure and `amber/orange/blue/red` for others. But the template files (errorpage.templ, errordetail.templ, erroralert.templ) use `gray-*` for neutral elements like "Go back" buttons and text colors.

---

## c) NOT STARTED

| Item | Impact | Effort |
|---|---|---|
| Fix `FromError()` to work with go-error-family's `int` Family | CRITICAL | Medium |
| Add `ParseFamily(string)` for string-based bridge | HIGH | Small |
| Add JSON error endpoint option to `ErrorHandler` | HIGH | Medium |
| HTML shell for standalone HTTP handler responses | HIGH | Small |
| Rename `ConflictError()` → `Conflict()` | LOW | Tiny |
| Test `DefaultErrorAlertProps()` | LOW | Tiny |
| Test `ConflictError("")` empty message default | LOW | Tiny |
| Test `ErrorHandler` with nil error | LOW | Tiny |
| Test nonce propagation in `ErrorHandler`/`WriteErrorPage` | LOW | Tiny |
| Test `Override` returning nil in `ErrorHandler` | LOW | Tiny |
| Add `FromError`/`ErrorHandler` examples to `doc.go` | MEDIUM | Small |
| Unify DismissScript call pattern (wrapper vs direct) | LOW | Tiny |
| Standardize `gray-*` vs `slate-*` dark mode palette | MEDIUM | Small |
| Log render errors in `ErrorHandler`/`WriteErrorPage` | MEDIUM | Tiny |
| Extract errorHeader shared sub-template from ErrorPage/ErrorDetail | LOW | Small |
| Consider `FamilyIsValid`/`FamilyIcon` unexporting | LOW | Tiny |
| `ExtractCauseChain` depth-capping test | LOW | Tiny |
| Combined-interface `FromError` test (all 3 interfaces on one error) | LOW | Tiny |

---

## d) TOTALLY FUCKED UP

### 1. `FromError()` Type Mismatch (CRITICAL)

```go
// handler.go:29 — checks for string-based Family
if c, ok := err.(interface{ ErrorFamily() Family }); ok {
    family = c.ErrorFamily()  // errorpage.Family (string)
}

// go-error-family returns int-based Family
// errorfamily.Family is int, not string
// → This assertion NEVER matches a real go-error-family error
```

The whole point of `FromError()` is to bridge go-error-family errors. Right now it **does not work** for the most important field (Family classification). Every real go-error-family error gets silently classified as Transient.

### 2. Handler Renders HTML Fragment, Not Valid Document

`ErrorHandler()` and `WriteErrorPage()` render `ErrorPage()` which is a `<div>`, not an HTML document. No `<!DOCTYPE html>`, no `<html>`, no `<head>`, no `<title>`, no charset meta. An HTTP handler returning this is serving invalid HTML.

### 3. HTMX Family-Aware Parsing Has No Producer

The JS code in `htmx/error_handling.templ` parses JSON error responses, but `errorpage/handler.go` only produces HTML. No JSON error endpoint exists. The family-aware feature is dead code in production.

---

## e) WHAT WE SHOULD IMPROVE

### Architecture

1. **Bridge pattern, not type matching:** `FromError()` should use go-error-family's `Classified` interface (returns `errorfamily.Family` int) and convert to our `errorpage.Family` string via `.String()`. This requires either importing go-error-family or using a string-based interface.

2. **HTML shell option:** `ErrorHandler` should optionally wrap ErrorPage in `layout.Minimal` or a minimal HTML shell with `<!DOCTYPE html>`, `<html>`, `<head>`, `<title>`.

3. **JSON error response:** `ErrorHandlerConfig` should support `JSON: true` to render JSON instead of HTML, enabling the HTMX family-aware parsing to actually work.

4. **Unified dismiss pattern:** Both `feedback` and `errorpage` should use the same call pattern for `utils.DismissScript()`.

### Type Model

5. **Consider importing go-error-family as an optional dependency.** The current "zero dependency" design is admirable but creates the broken bridge. Option: keep zero-dependency core, add a `errorpage/familybridge` sub-package that imports go-error-family.

6. **`ErrorPageProps` is 12 fields.** Consider splitting: `ErrorPageProps` (display) + `ErrorData` (data from error) so the display layer is separate from the error extraction layer.

### Testing

7. **`DefaultErrorAlertProps()` is untested.**
8. **No test for `FromError` with a real go-error-family error.**
9. **No test for `ErrorHandler` with nil error.**
10. **No test for `ConflictError("")` empty message.**
11. **Nonce propagation in handler tests not verified.**

### Polish

12. **`ConflictError()` naming inconsistency** — all other constructors are nouns (NotFound, Forbidden, BadRequest). ConflictError has an "Error" suffix.
13. **Dark mode `gray-*` vs `slate-*` inconsistency** across templates.
14. **Render errors silently discarded** in `ErrorHandler`/`WriteErrorPage`.

---

## f) Top 25 Things We Should Get Done Next

Sorted by **Impact × Ease** (Pareto ordering):

| # | Task | Impact | Effort | Category |
|---|---|---|---|---|
| 1 | Fix `FromError()` to detect go-error-family's string-based `ErrorFamily()` | CRITICAL | S | Bug fix |
| 2 | Add `ParseFamily(string) Family` for robust string→Family conversion | HIGH | S | API |
| 3 | Add test: `FromError` with go-error-family-style error | HIGH | S | Test |
| 4 | Rename `ConflictError()` → `Conflict()` | MED | XS | Naming |
| 5 | Test `DefaultErrorAlertProps()` returns `FamilyTransient` | LOW | XS | Test |
| 6 | Test `ConflictError("")` empty message default | LOW | XS | Test |
| 7 | Test `ErrorHandler` with nil error | LOW | XS | Test |
| 8 | Add HTML shell option to `ErrorHandlerConfig` | HIGH | S | Feature |
| 9 | Log render errors in handler instead of discarding | MED | XS | Fix |
| 10 | Add JSON error response mode to `ErrorHandlerConfig` | HIGH | M | Feature |
| 11 | Test nonce propagation in `ErrorHandler`/`WriteErrorPage` | LOW | XS | Test |
| 12 | Test `Override` returning nil in `ErrorHandler` | LOW | XS | Test |
| 13 | Update `doc.go` with `FromError`/`ErrorHandler` examples | MED | S | Docs |
| 14 | Unify DismissScript call pattern across feedback/errorpage | LOW | XS | Consistency |
| 15 | Standardize dark mode palette (`gray-*` → `slate-*` where appropriate) | MED | S | Polish |
| 16 | Extract errorHeader shared sub-template | LOW | S | Refactor |
| 17 | Consider unexporting `FamilyIsValid`/`FamilyIcon` if not consumer-facing | LOW | XS | API review |
| 18 | Test `ExtractCauseChain` depth capping | LOW | XS | Test |
| 19 | Test combined-interface `FromError` (all 3 on one error) | LOW | XS | Test |
| 20 | Consider optional go-error-family import in sub-package | MED | M | Architecture |
| 21 | Consider splitting `ErrorPageProps` into display + data layers | MED | M | Architecture |
| 22 | Verify `tailwind-merge-go` v0.2.1 concurrent safety (remove mutex?) | LOW | S | Perf |
| 23 | Consider `colorRoot` pattern for familyVisualStyle dedup | LOW | M | Refactor |
| 24 | Add `TestFromErrorWithRealGoErrorFamily` integration test | HIGH | S | Test |
| 25 | Wire HTMX JSON error response from handler to actual usage | HIGH | M | Integration |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should `templ-components` import `go-error-family` as a dependency?**

Arguments FOR:
- The broken `FromError()` bridge is caused by type mismatch — importing go-error-family would give us compile-time type safety
- We could use `errorfamily.Classified` interface directly instead of duck-typing
- `Family.String()` → `ParseFamily()` gives us the string bridge for free
- Both projects are by the same author (Lars)

Arguments AGAINST:
- `templ-components` currently has only 2 dependencies (templ, tailwind-merge-go) — adding go-error-family would break the "zero framework deps" principle
- Component library consumers who don't use go-error-family get an unnecessary dependency
- The current string-based approach is theoretically more flexible (works with any error classification system that returns "rejection"/"conflict"/etc strings)

**Proposed compromise:** Add `errorpage/familybridge/` sub-package that optionally imports go-error-family. The main `errorpage` package stays dependency-free. `FromError` in the main package checks for string-based `ErrorFamily() string` (not `Family`). The bridge sub-package converts `errorfamily.Family` → `errorpage.Family` using `ParseFamily(f.String())`.

**This requires your decision.** I cannot determine whether adding an optional sub-dependency is acceptable for this project's philosophy.
