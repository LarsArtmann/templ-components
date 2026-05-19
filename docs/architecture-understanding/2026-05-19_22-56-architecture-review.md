# Architecture Review — templ-components

**Date:** 2026-05-19 | **Reviewer:** Senior Staff Architect

---

## Overall Assessment: B+

A well-structured Go component library with strong type safety, good test coverage (71.8%), and clean package boundaries. The main architectural concerns are: JS architecture coupling, test suite duplication, and several shallow modules that could be deepened.

---

## Scalability & Modularity

### Package Dependency Graph (Current)

```
utils          ← all packages (0 deps)
internal/svg   ← display, feedback, icons, navigation
icons          ← display, feedback, layout, navigation
feedback       ← htmx (JS coupling only)
forms          ← utils only (clean boundary)
```

**Assessment:** The dependency graph is a clean DAG. No circular imports. `utils` is correctly at the bottom. `internal/svg` is properly internal. **This is good architecture.**

### Coupling Hotspots

| Hotspot                                              | Type                | Severity |
| ---------------------------------------------------- | ------------------- | -------- |
| `htmx` → `feedback` (via `tcShowToast` JS global)    | Runtime JS coupling | HIGH     |
| `display`, `feedback`, `navigation` → `internal/svg` | Compile-time, clean | LOW      |
| `AlertType` ≡ `ToastType` (same 4 values)            | Type duplication    | MEDIUM   |
| 7 files emit inline `<script>` blocks                | JS fragmentation    | MEDIUM   |

### Service Orientation Score: 7/10

The library is already well-decomposed into domain packages. Each package has a clear purpose. The `utils` package provides shared infrastructure. Improvements would come from:

1. Extracting a shared `FeedbackLevel` type for Alert/Toast/Badge
2. Creating a shared JS initialization strategy
3. Adding a `ComponentProps` interface for generic handling

---

## Composability Analysis

### Strong Points

1. **`utils.BaseProps` embedding** — All component props compose through a shared base. Excellent.
2. **`templ.Component` parameters** — Loading indicators accept arbitrary spinners. Cards accept footer/header content. This is the right pattern.
3. **Map-based style lookups** — Data-driven, extensible without code changes.
4. **`feedbackStyleSet` + generic lookup** — Best example of depth in the codebase.

### Weak Points

1. **`SimpleCard` doesn't compose with `Card`** — It duplicates the shell instead of calling `Card` with empty options.
2. **No `ComponentProps` interface** — 29 props structs share `BaseProps` but have no common interface for generic handling.
3. **`DropdownItem` uses empty-string discrimination** — Href="" means button, Href="url" means link. Should be a sum type or enum.
4. **Breadcrumb SVG is hardcoded** — Doesn't use the icon system, breaking the composition chain.

---

## Shallow Module Analysis

Using the "deletion test" (would deleting it concentrate complexity or just move it?):

| Module                      | Depth       | Deletion Test                    | Action               |
| --------------------------- | ----------- | -------------------------------- | -------------------- |
| `utils.Class()`             | **Deep**    | Yes — merge logic concentrates   | Keep                 |
| `feedbackStyleSet` + lookup | **Deep**    | Yes — generic style resolution   | Keep                 |
| `iconPathData` map          | **Deep**    | Yes — single source for 44 icons | Keep                 |
| `forms.SanitizeID()`        | **Shallow** | Just moves regex wrapping        | Keep (still useful)  |
| `layout.sri.go`             | **Medium**  | SRI hash management              | Keep                 |
| `feedback.Alert` vs `Toast` | **Shallow** | Nearly identical, separate maps  | **Merge style maps** |
| Test helpers in `utils/`    | **Medium**  | Useful for consumers             | Keep                 |

---

## JS Architecture Issues

The inline JS strategy creates maintenance problems:

| Component           | JS Lines | Global Guard            | Pattern          |
| ------------------- | -------- | ----------------------- | ---------------- |
| Accordion           | 27       | `tcAccordionAttached`   | Global singleton |
| Dropdown            | 48       | `tcDropdownAttached`    | Global singleton |
| Modal               | 45       | Per-instance IIFE       | Per-modal        |
| Toast+Alert dismiss | 10       | `tcDismissAttached`     | Shared           |
| Toast container     | 52       | —                       | DOM construction |
| Theme toggle        | 15       | `tcThemeToggleAttached` | Global singleton |
| Mobile menu         | 25       | `tcMobileMenuAttached`  | Global singleton |

**Total: ~222 lines of inline JS across 7 components.**

### Problems

1. **No shared init strategy** — Each component independently solves the "attach once" problem.
2. **HTMX swap breaks event listeners** — Global guards prevent re-attachment after HTMX swaps DOM elements.
3. **Theme toggle multi-instance bug** — `tcThemeToggleAttached` prevents second toggle from working.

### Recommendation

Create a `tc-init.js` module loaded once by `layout.Base` that provides:

- `tc.register(id, initFn)` — de-duplication per component ID
- `tc.dismiss(selector)` — shared dismiss handler
- `tc.onSwap(callback)` — re-attach after HTMX swaps

---

## Test Architecture Issues

Every package uses a 3-file pattern:

| File               | Purpose                         | Overlap with others               |
| ------------------ | ------------------------------- | --------------------------------- |
| `bdd_test.go`      | User-facing behavior assertions | 60-80% overlap with snapshot_test |
| `snapshot_test.go` | Render output assertions        | 60-80% overlap with bdd_test      |
| `a11y_test.go`     | Accessibility attribute checks  | 30-40% unique, rest overlaps      |

**Recommendation:** Consolidate into 1-2 files per package:

- `component_test.go` — All render + BDD + snapshot tests
- `a11y_test.go` — Only the unique a11y assertions (or merge into component_test)

---

## Key Metrics

| Metric           | Value                        | Assessment                         |
| ---------------- | ---------------------------- | ---------------------------------- |
| Packages         | 9                            | Clean separation                   |
| Circular imports | 0                            | Perfect                            |
| Test coverage    | 71.8%                        | Good for v0.x                      |
| External deps    | 2                            | Excellent (templ + tailwind-merge) |
| Props structs    | 29                           | Appropriate for 53 components      |
| Enums            | 17                           | Strong type safety                 |
| Inline JS        | 222 lines across 7 files     | Needs consolidation                |
| Test files       | 37                           | Too many, high overlap             |
| Max file size    | 174 lines (pagination.templ) | Under 350 threshold                |

---

## Top 5 Architectural Improvements

1. **Extract shared `FeedbackLevel` type** — Merge `AlertType`/`ToastType` into one type
2. **Consolidate inline JS** — Single shared init strategy, eliminate global guards
3. **Add `ComponentProps` interface** — Enable generic component handling
4. **Fix demo to use the library** — Currently uses raw HTML instead of layout.Base
5. **Consolidate test files** — Reduce 37 test files to ~15, eliminate 60%+ duplication
