# Status Report — Post-Fix Verification & Deep Audit

**Date:** 2026-05-29 15:08 CEST | **Branch:** master | **Working tree:** clean | **Tests:** 254 passing | **Lint:** 0 issues

---

## a) FULLY DONE

| Item                                | Details                                                                                      |
| ----------------------------------- | -------------------------------------------------------------------------------------------- |
| go-error-family integration         | `FromError()` uses `errors.AsType[errorfamily.Classified]()` for type-safe family extraction |
| `ParseFamily()`                     | Case-insensitive, trims whitespace, tested for uppercase/mixed/whitespace                    |
| `FamilyFromErrorFamily()`           | Converts `errorfamily.Family` (int) → `errorpage.Family` (string)                            |
| `ConflictError` → `Conflict` rename | Complete in all source, tests, docs                                                          |
| `ErrorHandlerConfig.HTMLShell`      | Wraps ErrorPage in valid HTML document                                                       |
| `ErrorHandlerConfig.JSON`           | Renders structured JSON for API/HTMX endpoints                                               |
| Render error logging                | `slog.Error` instead of silent discard                                                       |
| DismissScript unification           | Direct `utils.DismissScript()` calls everywhere, no private wrappers                         |
| `cleanMessage()` helper             | Uses `Message()` method for clean text (no `[family:code]` prefix)                           |
| `errorTimestamp()` helper           | Prefers error's own `Timestamp()` over `time.Now()`                                          |
| `doc.go` accuracy                   | Updated to reflect actual go-error-family integration                                        |
| `go.mod` correctness                | go-error-family correctly listed as direct dependency                                        |
| Tests                               | 30+ new tests for go-error-family integration, ParseFamily, HTMLShell, JSON                  |

---

## b) PARTIALLY DONE

| Item                        | Status                                                                                     | Gap                          |
| --------------------------- | ------------------------------------------------------------------------------------------ | ---------------------------- |
| Dark mode color consistency | errorpage uses `gray-*` in templates for neutral elements, rest of codebase uses `slate-*` | Cosmetic only, not blocking  |
| Templ version alignment     | go.mod has `v0.3.1020`, generator is `v0.3.1036`                                           | Works but generates warnings |

---

## c) NOT STARTED

| Item                                                                         | Impact | Effort |
| ---------------------------------------------------------------------------- | ------ | ------ |
| Upgrade templ to v0.3.1036 in go.mod                                         | LOW    | XS     |
| Remove dead `internal/svg` package                                           | MED    | S      |
| Consistently deprecate/remove `AlertType`/`ToastType` aliases                | LOW    | S      |
| Add `FromError` extraction of `IsRetryable()` for UI retry hints             | MED    | S      |
| Add `FromError` extraction of `Audience()` for show/hide technical details   | MED    | S      |
| Extract dismiss button pattern to shared sub-template                        | LOW    | M      |
| Wire HTMX `GlobalErrorHandling` to `ErrorHandlerConfig{JSON: true}` endpoint | HIGH   | M      |

---

## d) TOTALLY FUCKED UP — Nothing remaining

All previously identified critical issues have been fixed:

- ~~`FromError()` type mismatch~~ → Fixed with go-error-family integration
- ~~Handler renders HTML fragment~~ → Fixed with `HTMLShell` option
- ~~HTMX family-aware parsing ghost system~~ → Fixed with `JSON` response mode
- ~~`ParseFamily` not case-insensitive~~ → Fixed with `strings.ToLower` + `strings.TrimSpace`
- ~~`doc.go` false zero-dependency claim~~ → Fixed
- ~~go-error-family marked `// indirect`~~ → Fixed
- ~~Stale `ConflictError` test name~~ → Fixed
- ~~`FromError` uses `err.Error()` with `[family:code]` prefix~~ → Fixed with `cleanMessage()`
- ~~`FromError` ignores error's own timestamp~~ → Fixed with `errorTimestamp()`

---

## e) WHAT WE SHOULD IMPROVE

### Architecture

1. **Dead `internal/svg` package** — `FillIcon()` and `SpinnerSVG()` are never called. The `icons` package handles all SVG rendering independently. Should remove or integrate.

2. **Dismiss button duplication** — Alert, Toast, and ErrorAlert all have nearly identical dismiss button markup + JS injection. Could extract to a shared `dismissible()` sub-template.

3. **Feedback type aliases** — `AlertType` is deprecated but `ToastType` is not. Inconsistent. Both should either be deprecated or removed.

4. **Templ version mismatch** — go.mod `v0.3.1020`, generator `v0.3.1036`. Should upgrade.

### Type Model

5. **`ErrorPageProps` could gain `Retryable bool`** — extracted from go-error-family's `IsRetryable()` to drive "Try again" UI.

6. **`ErrorPageProps` could gain `Audience` field** — from go-error-family's `Audience()` to control whether to show/hide technical details (cause chain, context) for different audiences.

7. **`errorResponse` JSON struct** could include `retryable` and `timestamp` fields for richer API responses.

---

## f) Top 25 Things We Should Get Done Next

Sorted by **Impact × Ease** (Pareto):

| #   | Task                                                                                     | Impact | Effort | Category    |
| --- | ---------------------------------------------------------------------------------------- | ------ | ------ | ----------- |
| 1   | Upgrade templ to v0.3.1036 in go.mod                                                     | MED    | XS     | Deps        |
| 2   | Remove dead `internal/svg` package                                                       | MED    | S      | Cleanup     |
| 3   | Deprecate `ToastType` alias (or remove both `AlertType`/`ToastType`)                     | LOW    | XS     | Consistency |
| 4   | Add `Retryable bool` field to ErrorPageProps + extract in FromError                      | MED    | S      | Feature     |
| 5   | Add `errorResponse.Retryable` and `errorResponse.Timestamp` to JSON output               | MED    | XS     | Feature     |
| 6   | Demo: add ErrorHandler usage example (HTMLShell + JSON)                                  | MED    | S      | Docs        |
| 7   | Standardize dark mode `gray-*` → `slate-*` in errorpage templates                        | LOW    | S      | Polish      |
| 8   | Extract dismiss button to shared sub-template across alert/toast/erroralert              | LOW    | M      | Dedup       |
| 9   | Consider `Audience` field on ErrorPageProps for conditional detail display               | MED    | M      | Feature     |
| 10  | Add integration test: full HTMX error flow (JSON handler → HTMX toast)                   | HIGH   | M      | Test        |
| 11  | Verify `tailwind-merge-go` v0.2.1 concurrent safety (remove mutex?)                      | LOW    | S      | Perf        |
| 12  | Add `TestFromErrorCleanMessage` verifying no `[family:code]` prefix                      | MED    | XS     | Test        |
| 13  | Add `TestFromErrorTimestamp` verifying error's own timestamp is used                     | MED    | XS     | Test        |
| 14  | Remove `Exclamation` deprecated icon constant                                            | LOW    | XS     | Cleanup     |
| 15  | Fix FEATURES.md known issues (Toast icon duplication, Spinner SVG duplication)           | LOW    | M      | Docs        |
| 16  | Add `doc.go` to packages that lack them (display, forms, htmx)                           | LOW    | S      | Docs        |
| 17  | Consider `colorRoot` pattern for familyVisualStyle dedup                                 | LOW    | M      | Refactor    |
| 18  | Add `TestCombinedInterfaceFromError` (all interfaces on one error)                       | LOW    | XS     | Test        |
| 19  | Verify demo page renders all errorpage components correctly                              | LOW    | XS     | Test        |
| 20  | Add `TestExtractCauseChainDepthCapping`                                                  | LOW    | XS     | Test        |
| 21  | Consider removing `ExclamationCircle` duplicate icon                                     | LOW    | XS     | Cleanup     |
| 22  | Add `CHANGELOG.md` entry for `cleanMessage`/`errorTimestamp` improvements                | LOW    | XS     | Docs        |
| 23  | Check if `StatusBadge` still has correct mapping after all changes                       | LOW    | XS     | Test        |
| 24  | Consider adding `Layout` option to ErrorHandlerConfig (use layout.Base instead of shell) | LOW    | M      | Feature     |
| 25  | Review all `Default*Props()` constructors for consistency of defaults                    | LOW    | S      | Consistency |

---

## g) Top #1 Question

**Should we remove the `internal/svg` package?**

`FillIcon()` and `SpinnerSVG()` are defined there but never called from anywhere in the project. The `icons` package handles all SVG rendering via `icon.templ`. The package is dead code but removing it changes the package count and could affect any external consumers who import it (though `internal/` packages can't be imported externally by Go convention).

**Recommendation:** Remove it. It's internal (can't be imported), unused, and adds confusion.
