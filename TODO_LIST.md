# TODO List — templ-components

**Updated:** 2026-05-20

Legend: ✅ Done | 🔨 In Progress | ⬜ Not Started | ❌ Blocked | ⏭️ Deferred

---

## Session 12 (2026-05-20) — Wave 1-2 Execution

### Completed This Session

| #   | Status | Task                                    | Notes                                                          |
| --- | ------ | --------------------------------------- | -------------------------------------------------------------- |
| 1   | ✅     | Fix Modal focus restore (WCAG)          | Save/restore activeElement via data-tc-prev-focus              |
| 2   | ✅     | Fix ID propagation (6 components)       | Alert, Toast, StatCard, Nav, Dropdown, ProgressBar             |
| 3   | ✅     | Add BaseProps to Breadcrumbs            | BreadcrumbsProps struct with Items + DefaultBreadcrumbsProps() |
| 4   | ✅     | Remove dead code                        | Deref, DerefOr, MergeAttrs, BoolString + tests                 |
| 5   | ✅     | Replace hardcoded SVGs with icon system | Alert dismiss, Toast JS dismiss, StepIndicator checkmark       |
| 6   | ✅     | ProgressBar negative clamp              | percent < 0 → 0                                                |
| 7   | ✅     | BoolString → strconv.FormatBool         | Accordion uses strconv directly                                |
| 8   | ✅     | Fix retry counter race                  | Per-element data-tc-retry attribute                            |
| 9   | ✅     | Deduplicate test data                   | testNavLinks in navigation tests                               |

### Explicitly Skipped (Low Value)

| #   | Task                                 | Reason                                                    |
| --- | ------------------------------------ | --------------------------------------------------------- |
| —   | Add BaseProps to Spinner             | Building block primitive, 12+ call sites, rarely needs ID |
| —   | Add BaseProps to SimpleNav           | Convenience wrapper, delegates to Nav                     |
| —   | Add BaseProps to SimpleEmptyState    | Trivial one-liner component                               |
| —   | SimpleCard compose through Card      | Current implementation is cleaner                         |
| —   | JS consolidation (shared tc-init.js) | Risky runtime behavior change                             |
| —   | ComponentProps interface             | Significant boilerplate, low current ROI                  |

---

## Session 10-11 (2026-05-19) — Previously Completed

| #   | Status                                                   | Task |
| --- | -------------------------------------------------------- | ---- |
| ✅  | Demo app rewrite with layout.Base + Tailwind v4          |
| ✅  | FeedbackType unification (AlertType/ToastType → aliases) |
| ✅  | LoadingOverlay → props struct                            |
| ✅  | StepIndicator BaseProps                                  |
| ✅  | FillIcon variadic → bool                                 |
| ✅  | ThemeToggle multi-instance fix                           |
| ✅  | Modal stable IDs                                         |
| ✅  | Tooltip aria-describedby                                 |
| ✅  | Breadcrumbs icon system                                  |
| ✅  | Pagination URL builder (net/url)                         |
| ✅  | Test cleanup (splitClasses, benchmarks)                  |
| ✅  | CONTRIBUTING.md fix                                      |
| ✅  | Icon validation (unknown names panic)                    |
| ✅  | IconPathJS stroke-width fix                              |
| ✅  | Exclamation icon removal                                 |

---

## Code Quality Baseline

|     | #   | Status                            | Notes                         |
| --- | --- | --------------------------------- | ----------------------------- |
| 1   | ✅  | Build passes (`go build ./...`)   | Clean. Zero issues.           |
| 2   | ✅  | All tests pass (`go test ./...`)  | 9 packages, all green.        |
| 3   | ✅  | Lint passes (`golangci-lint run`) | 0 issues on library packages. |
| 4   | ✅  | Coverage: ~68%                    | Range per package.            |

---

## Remaining Work (Prioritized)

### P1 — Should Do Soon

| #   | Status | Task                                       | Package    | Notes                                                                                           |
| --- | ------ | ------------------------------------------ | ---------- | ----------------------------------------------------------------------------------------------- |
| 1   | ⬜     | Fix JS re-attachment after HTMX swaps      | multi      | Global tc\*Attached guards prevent re-init after DOM swap. Use per-element data-tc-initialized. |
| 2   | ⬜     | Add test coverage gaps (10 areas)          | multi      | Alert/Toast/ProgressBar/StepIndicator/Nav/Dropdown/Modal/CSRF/Pagination edge cases             |
| 3   | ⬜     | Validate SelectOption Disabled+Selected    | forms      | Impossible state per HTML spec                                                                  |
| 4   | ⬜     | Validate Pagination CurrentPage > 0        | navigation | CurrentPage: 0 renders page-0 link                                                              |
| 5   | ⬜     | Extract shared test panic assertion helper | multi      | Repeated recover()+nil check pattern                                                            |

### P2 — Nice to Have

| #   | Status | Task                                                  | Package  | Notes                                             |
| --- | ------ | ----------------------------------------------------- | -------- | ------------------------------------------------- |
| 6   | ⬜     | Extract shared dismiss JS (Alert + Toast)             | feedback | Already partially shared via tcDismissAttached    |
| 7   | ⬜     | Document htmx→feedback runtime JS dependency          | htmx     | GlobalErrorHandling requires ToastContainer       |
| 8   | ⬜     | Document fill vs stroke 20×20/24×24 convention        | icons    | Code comment                                      |
| 9   | ⬜     | Add go doc examples (ExampleXxx functions)            | all      | For pkg.go.dev discoverability                    |
| 10  | ⬜     | Document thread-safety requirement in CONTRIBUTING.md | docs     | tailwind-merge-go NOT thread-safe, mutex required |

### P3 — Deferred (Post v1.0)

| #   | Status | Task                                            | Notes                              |
| --- | ------ | ----------------------------------------------- | ---------------------------------- |
| 11  | ⏭️     | Consolidate test files (37→15)                  | Post-v1.0                          |
| 12  | ⏭️     | Convert snapshot tests to golden files          | Post-v1.0                          |
| 13  | ⏭️     | Move test helpers out of utils/                 | Post-v1.0                          |
| 14  | ⏭️     | Add Radio, File input, Toggle/Switch components | New features                       |
| 15  | ⏭️     | Client-side JS tab switching                    | Enhancement                        |
| 16  | ⏭️     | PageProps zero-value safety                     | API change                         |
| 17  | ⏭️     | uint for Pagination fields                      | API change                         |
| 18  | ⏭️     | Icon list auto-gen from path map                | Build tooling                      |
| 19  | ⏭️     | ComponentProps interface                        | Significant boilerplate            |
| 20  | ⏭️     | DropdownItem typed variant                      | Breaking API change                |
| 21  | ⏭️     | Modularization (go.work)                        | Analysis concluded NOT recommended |
