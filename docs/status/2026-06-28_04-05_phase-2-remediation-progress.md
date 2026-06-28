# Status Update — 2026-06-28 04:05

## Context

**What triggered this work:** Continuation of the full-code-review remediation. The first phase completed all Tier 1 (Critical, 19/19) and Tier 2 (High, 22/22) tasks plus ~35 Tier 3/4 tasks. This phase (8 commits, 11 tasks) tackled the top remaining items from the Pareto plan.

**Baseline before this phase:** Build OK · tests green · lint 0 issues · working tree clean. 8 commits since last status report (`eed695a`).

**Current state after this phase:** Build OK · tests green (13/13 packages) · lint 0 issues · working tree clean.

---

## a) FULLY DONE — shipped and verified this phase

### Phase 2 — 11 tasks completed (8 commits)

| Task  | What was fixed                                                               | Package    |
| ----- | ----------------------------------------------------------------------------- | ---------- |
| P-074 | Tooltip auto-generates ID via EnsureID; aria-describedby always present       | display    |
| P-043 | FormMethod doc comment explains hx-put/hx-delete/hx-patch via Attrs            | forms      |
| P-091 | Removed buttonVariantDefault/badgeStyleDefault duplicate constants            | display    |
| P-092 | Drawer inline style → Tailwind classes (inset-y-0 left-0/right-0)              | display    |
| P-036 | WriteErrorPage derives status from FamilyStatusCode when statusCode is 0       | errorpage  |
| P-037 | renderToBuffer/renderShellToBuffer render to bytes.Buffer before WriteHeader   | errorpage  |
| P-038 | FromError extracts Title from errors implementing ErrorTitle()                 | errorpage  |
| P-041 | Split 354-line handler.go into fromerror.go + constructors.go + handler.go    | errorpage  |
| P-049 | FormFieldWrapper 5 positional params → FormFieldProps struct                  | forms      |
| P-047 | HTMXVersion typed string const (HTMXVersion2_0_10) replacing bare string      | layout     |
| P-048 | ThemeColor/DarkThemeColor hex validation with fallback to defaults            | layout     |
| P-105 | ExtractCauseChain handles errors.Join siblings (Unwrap() []error, Go 1.20+)   | errorpage  |
| P-073 | Tooltip touch fallback (click toggle) + Escape-to-dismiss via singleton JS     | display    |
| P-050 | ConfirmDelete/SwapOOB converted from positional params to Props structs       | htmx       |

### Also done (bonus, not in original plan)

| What                                                                  | Package   |
| -------------------------------------------------------------------- | --------- |
| ErrorHandlerConfig.Lang field for i18n html lang attribute override  | errorpage |
| writeFallbackError for graceful degradation when error page fails    | errorpage |
| ExtractCauseChain refactored: appendJoinSiblings + causeItemFromError helpers | errorpage |
| Tooltip: data-tc-tooltip attribute for JS targeting                  | display   |

### Cumulative totals (Phase 1 + Phase 2)

- **Fully done:** ~90 of 123 tasks (73%)
- **Tier 1 (Critical):** 19/19 (100%) ✅
- **Tier 2 (High):** 22/22 (100%) ✅
- **Tier 3 (Medium):** ~26/41 (63%)
- **Tier 4 (Polish):** ~23/41 (56%)

---

## b) PARTIALLY DONE

| Area               | What's done                                                                                        | What remains                                                                                                          |
| ------------------ | -------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| DropdownItem Kind  | P-046: DropdownItemKind type + Kind field exist; IsLink() uses Kind                                 | P-046: Href fallback not deprecated; backward-compat path still active                                                |
| errorpage safety   | P-036, P-037, P-038, P-041, P-105 all done                                                         | P-039 (contextTable → dl), P-064 (familyStyleMap builder), P-121 (i18n label props)                                   |
| forms architecture | P-049 (FormFieldWrapper Props), P-043 (FormMethod docs)                                            | P-051 (tabLink struct), P-052 (paginationArrow struct), P-053 (activeSpanOrLink struct)                              |
| htmx architecture  | P-050 (ConfirmDelete/SwapOOB Props), P-044 (SwapStyle enum), P-057 (singleton scripts)             | P-054 (diagnosticSection struct), P-056 (DismissScript move)                                                         |

---

## c) NOT STARTED

All remaining tasks are Tier 3 (Medium) or Tier 4 (Polish) — no remaining Critical or High issues:

| Task  | Description                                                      | Effort |
| ----- | ---------------------------------------------------------------- | ------ |
| P-039 | contextTable: use dl or add caption/th+scope                     | S      |
| P-042 | StatusBadge: typed status alias with constants                   | S      |
| P-045 | paginationArrow roundedSide: typed enum                          | S      |
| P-046 | DropdownItem: deprecate Href fallback, require Kind (partial)    | S      |
| P-051 | tabLink: pack active/inactive classes into struct                | S      |
| P-052 | paginationArrow (7 params) → props struct                        | S      |
| P-053 | activeSpanOrLink (5 params) → struct or split                     | S      |
| P-054 | diagnosticSection (5 params) → props struct                      | S      |
| P-056 | Move DismissScript out of utils (foundation leaks DOM/JS)        | S      |
| P-064 | errorpage familyStyleMap: builder to dedup 6 entries             | M      |
| P-076 | Icon accessible variant (role=img + title)                      | M      |
| P-080 | utils Class(): per-shard mutex or result cache                   | L      |
| P-111 | IconPathJS: accept strokeWidth arg (hardcoded 1.5)               | XS     |
| P-113 | layout: consolidate 7 overlapping test files into 2              | M      |
| P-114 | feedback: add toastJSStyles/toastJSIconPaths table tests          | S      |
| P-119 | paginationRange: add current>total test case                     | XS     |
| P-120 | DefinitionItem/TableCell: unify dual string+Component slot       | S      |
| P-121 | errorpage: accept label strings via props for i18n               | M      |
| P-122 | SidebarNavItem: embed BaseProps                                  | XS     |

---

## d) TOTALLY FUCKED UP

**Nothing.** No regressions introduced. Build, tests, and lint all pass. No files were accidentally destroyed.

**What I could have done better this phase:**

1. **Lint issues introduced:** The `contextcheck` and `nestif` linters flagged the new `ExtractCauseChain` and `renderToBuffer` code. Fixed immediately but should have anticipated them during implementation.
2. **The `nestif` fix** required extracting `appendJoinSiblings` and `causeItemFromError` helpers — this is actually better code, but it should have been the initial implementation, not a lint-driven refactor.
3. **The `multiedit` on handler.go** failed for one of two edits (file modified between view and edit). Had to re-read and retry. Should have used `view` immediately before `edit` on a file that was recently written.

---

## e) WHAT WE SHOULD IMPROVE

### Architecture

1. **Remaining positional params → Props structs:** 5 templ functions still use 5-7 positional string params (tabLink, paginationArrow, activeSpanOrLink, diagnosticSection, Label). These are the last holdouts from the library-wide Props struct convention. Each is a footgun where swapping adjacent strings compiles fine.

2. **utils.DismissScript leaks DOM concerns:** The foundation package returns raw browser JS. Moving to `htmx` or a `scripts` package would keep `utils` dependency-free.

3. **errorpage familyStyleMap:** 6 near-identical 8-field entries with no dedup. A builder pattern would reduce drift risk.

### Type Safety

4. **StatusBadge:** Accepts raw `string` — status values are known at compile time but get no type safety. A typed alias (like HTMXVersion) would catch typos.

5. **paginationArrow roundedSide:** Plain `string` with `"l"`/`"r"` literals — should be a typed enum.

### Testing

6. **layout test file consolidation:** 7 overlapping test files with massive duplication. Consolidating into 2 (unit + behaviour) would reduce maintenance burden.

7. **toastJSStyles/toastJSIconPaths:** Non-trivial string builders that serialize style/icon maps to JS have no direct unit tests.

8. **paginationRange:** No test case for `current > total` edge case.

### A11y

9. **Icon a11y:** Every icon hardcodes `aria-hidden="true"` with no API for meaningful icons (role="img" + title). Icons used as sole button content are invisible to AT.

10. **contextTable:** Uses bare `<table>/<td>` without `<dl>`, `<caption>`, or `<th scope>` — screen readers mis-announce definition data as table rows.

### Performance

11. **utils.Class():** Holds a global `sync.Mutex` across the entire `twmerge.Merge()` call. A per-shard mutex or `sync.Pool` of twmerge instances would reduce contention under concurrent SSR.

---

## f) Top 25 things to get done next

Sorted by impact × inverse effort (highest value first):

| #   | Task                                                      | Impact | Effort | Rationale                                                                 |
| --- | --------------------------------------------------------- | ------ | ------ | ------------------------------------------------------------------------- |
| 1   | P-052 paginationArrow → props struct                      | XS     | S      | 7 positional params; highest param count in the package                   |
| 2   | P-053 activeSpanOrLink → struct or split                  | XS     | S      | 5 positional params; used by NavLink and SidebarNav                       |
| 3   | P-051 tabLink pack active/inactive classes into struct    | XS     | S      | 4 string params; swapping active/inactive compiles silently              |
| 4   | P-054 diagnosticSection → props struct                    | XS     | S      | 5 positional params in errorpage shared template                          |
| 5   | P-042 StatusBadge typed alias                             | S      | S      | Magic string map is only validation; typed alias adds compile-time safety |
| 6   | P-045 paginationArrow roundedSide typed enum              | XS     | XS     | "l"/"r" literals → typed enum; 5-minute fix                               |
| 7   | P-046 DropdownItem deprecate Href fallback                | M      | S      | Dual discrimination drops href silently; Kind should be required          |
| 8   | P-039 contextTable → dl                                    | S      | S      | Screen readers mis-announce definition data as table                      |
| 9   | P-076 Icon accessible variant (role=img + title)          | L      | M      | Icons as sole button content invisible to AT                              |
| 10  | P-122 SidebarNavItem embed BaseProps                      | XS     | XS     | Missing Class/ID/Attrs/AriaLabel support                                  |
| 11  | P-111 IconPathJS accept strokeWidth arg                   | XS     | XS     | Hardcoded 1.5; 5-minute fix                                               |
| 12  | P-119 paginationRange current>total test case            | XS     | XS     | Edge case with no test coverage                                           |
| 13  | P-114 toastJSStyles/toastJSIconPaths tests                | S      | S      | Non-trivial untested string builders                                       |
| 14  | P-056 Move DismissScript out of utils                     | S      | S      | Foundation package leaks browser-specific DOM concerns                    |
| 15  | P-064 errorpage familyStyleMap builder                    | M      | M      | 6 near-identical 8-field entries → builder reduces drift                 |
| 16  | P-120 DefinitionItem/TableCell unify dual slot            | S      | S      | Two structs with same dual string+Component pattern                      |
| 17  | P-121 errorpage accept label strings via props for i18n   | M      | M      | "Go home", "Go back" etc. hardcoded; no i18n path                         |
| 18  | P-113 layout consolidate 7 test files → 2                 | M      | M      | Massive duplication; rename "coverage_boost" anti-pattern                 |
| 19  | P-080 utils Class() per-shard mutex                      | M      | L      | Real contention bottleneck under concurrent SSR                           |
| 20  | P-100 InputGroupPaddingClass: audit usage                  | XS     | XS     | Verify it's wired correctly or delete                                     |
| 21  | P-103 Remove unreachable props.Timestamp=="" branch      | XS     | XS     | Already done — verify and close                                           |
| 22  | P-096 pageURL preserve fragment across re-encode          | XS     | XS     | Already done — verify and close                                           |
| 23  | P-112 example_test // Output directives                   | S      | S      | Already done for forms + errorpage — audit remaining packages            |
| 24  | P-047 HTMXVersion typed const set                         | M      | S      | Already done — add more version constants as htmx releases                |
| 25  | P-048 ThemeColor hex validation                           | S      | S      | Already done — consider named color support (e.g. "indigo-600")           |

---

## g) Top #1 question I cannot figure out myself

**Should the remaining positional-param → Props-struct conversions (P-051 tabLink, P-052 paginationArrow, P-053 activeSpanOrLink) preserve backward compatibility, or is this a breaking change?**

These are private (lowercase) templ functions — `tabLink`, `paginationArrow`, `activeSpanOrLink`, `diagnosticSection` — so they're not part of the public API. Breaking them only affects internal code, not consumers. But `Label` is exported (`forms.Label(forID, text string, required bool)`) and IS part of the public API. Converting `Label` to a Props struct would be a breaking change for any consumer calling `forms.Label(...)` directly.

The question: **should `Label` stay as-is (3 positional params, exported, breaking to change), or should we add a `LabelProps` struct variant (`LabelWithProps(props LabelProps)`) alongside the existing function for backward compatibility?** The `FormFieldWrapper` conversion was safe because it's used internally, but `Label` is a public API used by consumers directly.

---

## Verification status

| Check          | Command                   | Result                          |
| -------------- | ------------------------- | ------------------------------- |
| Build          | `go build ./...`          | ✅ Clean                        |
| Tests          | `go test ./...`           | ✅ 13/13 packages green         |
| Lint           | `golangci-lint run ./...` | ✅ 0 issues                     |
| Git            | `git status`              | ✅ Clean working tree           |
