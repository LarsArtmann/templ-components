# Status Report — templ-components (Session 3 Continuation)

**Date:** 2026-05-07 06:37 | **Branch:** master | **Total Commits:** 28

---

## Executive Summary

Continued executing the improvement plan from the previous audit. Completed 5 more small, focused changes across a11y, class merge, and documentation. All changes build, test, and lint clean. 6 uncommitted files ready for commit.

**Overall health: 7.5/10** — Improving steadily. Core components are solid. Remaining gaps are mostly polish and completeness.

---

## A) FULLY DONE ✅

### Session 3 Cumulative (10 commits committed + 1 pending)

| #   | Commit    | Description                                                                   |
| --- | --------- | ----------------------------------------------------------------------------- |
| 1   | `39e237e` | fix: Tailwind class merge, a11y gaps, form help text IDs                      |
| 2   | `1c014ef` | feat: 17 DefaultXxxProps() constructors                                       |
| 3   | `3d5aded` | feat!: type-safe icons (DropdownItem.Icon, EmptyStateProps.Icon → icons.Name) |
| 4   | `c9981d7` | feat(a11y): modal focus trap + Escape key                                     |
| 5   | `a1a81b5` | feat(a11y): dropdown keyboard navigation                                      |
| 6   | `e8acc3f` | fix(a11y): tabs ARIA linkage                                                  |
| 7   | `e6b403c` | fix(a11y): tooltip aria-describedby id                                        |
| 8   | `9884910` | docs: package doc comments                                                    |
| 9   | `dbbbdfc` | docs: comprehensive audit + TODO update                                       |
| 10  | `10dec01` | fix(a11y): aria-required on form inputs                                       |

### Uncommitted (pending commit after this report)

| #   | File(s)                          | Description                                                |
| --- | -------------------------------- | ---------------------------------------------------------- |
| 11  | `display/avatar.templ`           | `aria-hidden="true"` on status indicator dot               |
| 12  | `feedback/alert.templ`           | `utils.Class()` + `props.Class` support                    |
| 13  | `feedback/toast.templ`           | `utils.Class()` + `props.Class` support                    |
| 14  | `feedback/progress.templ`        | `utils.Class()` + `props.Class` support on wrapper         |
| 15  | `htmx/loading.templ`             | `role="status"` + `aria-live="polite"` on LoadingIndicator |
| 16  | `docs/migration/v0.1-to-v0.2.md` | Added icons.Name breaking change section                   |

### Previously Verified Already Done (from prior sessions)

All items from `docs/status/2026-05-07_06-09_comprehensive-audit-and-roadmap.md` section A remain done. Key items confirmed already-existed:

- All 22 Props structs have `DefaultXxxProps()` constructors — the audit was wrong about 11 missing
- `<html lang>` on Base layout — `props.Locale` (default `"en"`) already used as `lang` attribute
- Skip-to-content link — already in Base layout at line 93
- Avatar `Alt` field — already on AvatarProps and rendered on `<img>`
- Table `Caption` — already on TableProps, rendered as `<caption class="sr-only">`
- Table header `scope="col"` — already on all `<th>` elements

---

## B) PARTIALLY DONE 🔨

### utils.Class() Migration

| Component              | Before                  | After            | Status                                     |
| ---------------------- | ----------------------- | ---------------- | ------------------------------------------ |
| forms/Input            | comma-join              | `utils.Class()`  | ✅ Done                                    |
| forms/Select           | comma-join              | `utils.Class()`  | ✅ Done                                    |
| forms/Textarea         | comma-join              | `utils.Class()`  | ✅ Done                                    |
| forms/Checkbox         | comma-join              | `utils.Class()`  | ✅ Done                                    |
| display/Badge          | comma-join              | `utils.Class()`  | ✅ Done                                    |
| display/Card           | comma-join              | `utils.Class()`  | ✅ Done                                    |
| display/Tabs           | comma-join              | `utils.Class()`  | ✅ Done                                    |
| display/Tooltip        | comma-join              | `utils.Class()`  | ✅ Done                                    |
| display/Table          | comma-join              | `utils.Class()`  | ✅ Done                                    |
| display/Avatar         | comma-join              | `utils.Class()`  | ✅ Done                                    |
| display/Accordion      | comma-join              | `utils.Class()`  | ✅ Done                                    |
| navigation/Nav         | comma-join              | `utils.Class()`  | ✅ Done                                    |
| navigation/Pagination  | comma-join              | `utils.Class()`  | ✅ Done                                    |
| feedback/Alert         | comma-join              | `utils.Class()`  | ✅ Done (pending commit)                   |
| feedback/Toast         | comma-join              | `utils.Class()`  | ✅ Done (pending commit)                   |
| feedback/Progress      | comma-join              | `utils.Class()`  | ✅ Done (pending commit)                   |
| **display/Modal**      | `templ.KV`              | stays comma-join | ⚠️ Can't fix — KV returns non-string type  |
| **display/Dropdown**   | `templ.KV` + positional | stays comma-join | ⚠️ Can't fix — KV used for conditionals    |
| **navigation/NavLink** | `templ.KV`              | stays comma-join | ⚠️ Can't fix — KV used for active state    |
| **feedback/Loading**   | positional params       | N/A              | ⚠️ No Props struct, no BaseProps           |
| **layout/Base**        | string concat           | stays as-is      | ⚠️ External script tags, no class conflict |

**Summary:** 16/21 components converted. 5 remain with valid technical reasons (KV limitation or no Props struct).

---

## C) NOT STARTED ⬜

| #   | Task                                        | Priority | Effort | Type    |
| --- | ------------------------------------------- | -------- | ------ | ------- |
| 1   | Release automation (goreleaser)             | P3       | Medium | infra   |
| 2   | Documentation site generation               | P3       | Large  | DX      |
| 3   | Enhance demo app to showcase all components | P2       | Medium | DX      |
| 4   | Nix flake migration                         | P3       | Medium | infra   |
| 5   | Golden file snapshot tests                  | P2       | Medium | testing |
| 6   | Component composition integration tests     | P2       | Medium | testing |
| 7   | Improve forms test coverage (58% → 75%+)    | P2       | Medium | testing |
| 8   | Improve utils test coverage (56% → 75%+)    | P2       | Medium | testing |
| 9   | Tooltip JS-based aria-describedby injection | P2       | Medium | a11y    |
| 10  | Exclude `examples/` from golangci-lint      | P3       | Medium | DX      |

---

## D) TOTALLY FUCKED UP 💥

| #   | Issue                                      | Severity | Details                                                                                                                                                                         |
| --- | ------------------------------------------ | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **golangci-lint cache corruption**         | Annoying | Every lint run shows 40+ "Failed to persist facts to cache" warnings. Cache dir needs cleanup: `rm -rf ~/.cache/golangci-lint/`. Not a real issue — 0 actual lint issues found. |
| 2   | **LSP shows 40+ false errors**             | DX pain  | gopls can't see templ-generated types. Only `go build` matters. This is a templ ecosystem limitation.                                                                           |
| 3   | **Pre-commit hook was from `buildflow`**   | Weird    | `.git/hooks/pre-commit` contained a `buildflow` hook, not our templ hook. Fixed by copying `scripts/pre-commit.sh` over it.                                                     |
| 4   | **Audit overcounted missing DefaultProps** | Minor    | Previous status report claimed 11 missing DefaultProps. Actually 0 were missing — the script was matching component file names to constructor names incorrectly.                |

---

## E) WHAT WE SHOULD IMPROVE

### Already Good

- ✅ All Props structs have DefaultXxxProps()
- ✅ All form elements have aria-required, aria-invalid, aria-describedby
- ✅ Modal has focus trap + Escape + focus management
- ✅ Dropdown has keyboard navigation
- ✅ Tabs have full ARIA linkage
- ✅ Base layout has `<html lang>`, skip-to-content, security headers
- ✅ Avatar has alt text, Table has caption + scope
- ✅ All packages have doc comments
- ✅ Migration guide covers all breaking changes
- ✅ 16/21 components use utils.Class() for proper class merging
- ✅ HTMX LoadingIndicator has aria-live

### Still Needs Work

1. **Test coverage gaps** — forms (58%) and utils (56%) are well below the ~70% average
2. **No integration/composition tests** — components tested in isolation only
3. **No visual regression testing** — no golden file or screenshot comparison
4. **Demo app is minimal** — 151 lines, doesn't showcase all components
5. **No release automation** — no goreleaser, no version tags
6. **No documentation site** — only README + FEATURES.md

---

## F) TOP 25 THINGS TO DO NEXT

Sorted by impact × effort:

| #   | Task                                                | Impact | Effort  | Type    |
| --- | --------------------------------------------------- | ------ | ------- | ------- |
| 1   | Improve forms test coverage (58% → 75%+)            | Medium | Medium  | testing |
| 2   | Improve utils test coverage (56% → 75%+)            | Medium | Medium  | testing |
| 3   | Add component composition integration tests         | Medium | Medium  | testing |
| 4   | Enhance demo app to showcase all 53 components      | Medium | Medium  | DX      |
| 5   | Golden file snapshot tests                          | Medium | Medium  | testing |
| 6   | Tooltip JS-based aria-describedby injection         | Medium | Medium  | a11y    |
| 7   | Documentation site (templ-rendered)                 | High   | Large   | DX      |
| 8   | Set up goreleaser for versioned releases            | High   | Medium  | infra   |
| 9   | Nix flake migration                                 | Medium | Medium  | infra   |
| 10  | Add more icons (social, brand, common UI)           | Low    | Small   | feature |
| 11  | Add Form component (wraps inputs + validation)      | Medium | Medium  | feature |
| 12  | Add Dialog/Drawer component variants                | Medium | Medium  | feature |
| 13  | Add FileUpload component                            | Medium | Medium  | feature |
| 14  | Add Date Picker component                           | Medium | Large   | feature |
| 15  | Add Combobox/Autocomplete component                 | Medium | Large   | feature |
| 16  | Add Breadcrumb structured data (JSON-LD)            | Low    | Small   | SEO     |
| 17  | Add Pagination SEO rel=prev/next                    | Low    | Trivial | SEO     |
| 18  | Add skeleton component variants (card, table, list) | Low    | Small   | feature |
| 19  | Add dark mode toggle persistence example            | Low    | Trivial | DX      |
| 20  | Add responsive navigation pattern examples          | Low    | Small   | DX      |
| 21  | Investigate templ playground integration            | Low    | Medium  | DX      |
| 22  | Add CONTRIBUTING.md                                 | Low    | Small   | docs    |
| 23  | Add GitHub issue/PR templates                       | Low    | Trivial | infra   |
| 24  | Clean up golangci-lint cache corruption             | Low    | Trivial | DX      |
| 25  | Exclude examples/ from golangci-lint properly       | Low    | Medium  | DX      |

---

## G) TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF

**Should this library grow beyond display/feedback/forms/navigation into higher-level components (Form with validation, DataTable with sorting/filtering, full page layouts)?**

The current scope is atomic components — individual building blocks. But consumers often need composed patterns:

- `Form` that wraps inputs + validation + submission
- `DataTable` that adds sorting, filtering, pagination to `Table`
- `DashboardLayout` that composes Nav + Sidebar + Content

Adding these would dramatically increase value but also complexity, coupling, and opinionatedness. This is a scope decision that affects architecture direction and can't be made without understanding your product plans.

---

## Build & Test Status

```
✅ go build ./...         — PASS
✅ go test ./...          — PASS (all 10 packages)
✅ golangci-lint run      — 0 issues
✅ templ generate ./...   — PASS
```

## Test Coverage

| Package      | Coverage  |
| ------------ | --------- |
| display      | 66.6%     |
| feedback     | 71.8%     |
| forms        | 58.0%     |
| htmx         | 77.3%     |
| icons        | 73.0%     |
| internal/svg | 79.0%     |
| layout       | 73.0%     |
| navigation   | 72.0%     |
| utils        | 56.4%     |
| **Average**  | **69.7%** |

## Session Stats

| Metric                      | Value                                         |
| --------------------------- | --------------------------------------------- |
| Commits this session        | 10 (committed) + 1 (pending)                  |
| Files changed (uncommitted) | 6 files, +21/-4 lines                         |
| Dependencies                | 2 (templ v0.3.1001, tailwind-merge-go v0.2.1) |
| Total components            | 53                                            |
| Total icons                 | 42                                            |
