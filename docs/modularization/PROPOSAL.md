# Modularization Proposal — templ-components

**Created:** 2026-05-14
**Status:** Self-reviewed

---

## 1. Executive Summary

templ-components is a Go component library with a single `go.mod` containing 8 public packages and 1 internal package. The project has clean, acyclic dependencies and well-defined boundaries. Modularizing into sub-modules will:

- **Enable selective imports** — consumers import only the components they need without pulling transitive deps
- **Enforce compile-time boundaries** — `internal/svg` becomes a proper module boundary, not a Go `internal/` trick
- **Enable independent versioning** — icons, utils, and layout can version independently from display/feedback
- **Speed up CI** — only changed modules need rebuilding/testing

**What changes:** One `go.mod` → ten `go.mod` files coordinated by `go.work`.

**What stays the same:** All public APIs, all package paths (except `internal/svg` → `svg`), all import paths for consumers.

---

## 2. Current State Analysis

### 2.1 Starting State

| Indicator            | Value                                                     |
| -------------------- | --------------------------------------------------------- |
| go.mod files         | 1 (root)                                                  |
| go.work              | No                                                        |
| Go version           | 1.26.2                                                    |
| External deps        | `a-h/templ` v0.3.1001, `Oudwins/tailwind-merge-go` v0.2.1 |
| Packages             | 8 public + 1 internal                                     |
| Circular imports     | None                                                      |
| State classification | **Monolith** (single go.mod, all packages in one tree)    |

### 2.2 Package Dependency Graph (Production)

```
                    ┌─────────┐
                    │  utils   │  (leaf — zero internal deps)
                    └────┬─────┘
                         │
              ┌──────────┼──────────┐
              │          │          │
         ┌────┴───┐ ┌────┴───┐ ┌───┴────┐
         │ forms  │ │feedback│ │navigation│
         └────────┘ └───┬────┘ └────────┘
                       │
                  ┌────┴─────┐
                  │  htmx    │
                  └──────────┘

         ┌──────────────────────────────┐
         │          internal/svg         │  (leaf)
         └──────┬───────────┬───────────┘
                │           │
         ┌──────┴──┐  ┌─────┴────┐
         │feedback │  │  icons   │
         └─────────┘  └──────────┘

         ┌──────────────────────────────┐
         │          icons               │
         └──────────┬───────────────────┘
                    │
              ┌─────┴──────┐
              │  display   │
              └────────────┘
```

### 2.3 Detailed Import Table (Production)

| Package         | Imports                                                |
| --------------- | ------------------------------------------------------ |
| `utils`         | _(none — leaf)_                                        |
| `internal/svg`  | _(none — leaf)_                                        |
| `layout`        | _(none — leaf, only uses `templ` directly)_            |
| `icons`         | `internal/svg`                                         |
| `forms`         | `utils`                                                |
| `feedback`      | `utils`, `internal/svg`                                |
| `navigation`    | `utils`, `internal/svg`                                |
| `display`       | `utils`, `icons`, `internal/svg`                       |
| `htmx`          | `feedback`                                             |
| `examples/demo` | `display`, `feedback`, `icons`, `layout`, `navigation` |

### 2.4 Test-Only Internal Imports

| Package        | Test-Only Imports                     |
| -------------- | ------------------------------------- |
| `display`      | `feedback`, `icons` (already in prod) |
| `htmx`         | `utils`                               |
| `icons`        | `utils`                               |
| `internal/svg` | `utils`                               |
| All others     | _(none)_                              |

### 2.5 Coupling Analysis

**Low coupling overall.** The project is well-structured:

- `utils` is a true leaf — no internal deps, consumed by 5 of 8 packages
- `internal/svg` is a leaf — consumed by 3 packages for SVG primitives
- `layout` is completely standalone — zero internal deps
- `htmx` → `feedback` is the only package-to-package non-utils dependency
- `display` is the most coupled — depends on `utils`, `icons`, and `internal/svg`

**No god-packages detected.** The largest package (`display`) has 26 files but they are all coherent UI components within the same domain.

### 2.6 Key Constraint: `internal/svg`

The `internal/` directory is a Go convention that prevents external consumers from importing it directly. In the current monolith, `display`, `feedback`, `icons`, and `navigation` can all reach it. After modularization, if `internal/svg` stays in one module, only that module's consumers can use it.

**Resolution:** Promote `internal/svg` to a public `svg` sub-module. It contains only 2 functions (`FillIcon`, `SpinnerSVG`) — a thin SVG primitive layer. Making it public is safe and enables cross-module sharing.

---

## 3. Proposed Module Structure

### 3.1 Module Definitions

| Module            | Path             | Purpose                                             | Prod Deps (Internal)                                   | External Deps                |
| ----------------- | ---------------- | --------------------------------------------------- | ------------------------------------------------------ | ---------------------------- |
| **utils**         | `/utils`         | Base types, Tailwind class merging, generic helpers | _(none)_                                               | `templ`, `tailwind-merge-go` |
| **svg**           | `/svg`           | Shared SVG rendering primitives                     | _(none)_                                               | `templ`                      |
| **layout**        | `/layout`        | Page layouts, theme toggle, dark mode               | _(none)_                                               | `templ`                      |
| **icons**         | `/icons`         | Named SVG icon constants and rendering              | `svg`                                                  | `templ`                      |
| **feedback**      | `/feedback`      | Alerts, toasts, spinners, progress, skeletons       | `utils`, `svg`                                         | `templ`                      |
| **forms**         | `/forms`         | Form controls: input, select, textarea, checkbox    | `utils`                                                | `templ`                      |
| **display**       | `/display`       | Cards, badges, modals, tables, tabs, avatars, etc.  | `utils`, `icons`, `svg`                                | `templ`                      |
| **navigation**    | `/navigation`    | Navbars, breadcrumbs, pagination, mobile menu       | `utils`, `svg`                                         | `templ`                      |
| **htmx**          | `/htmx`          | HTMX helpers: loading, error handling, CSRF         | `feedback`                                             | `templ`                      |
| **examples/demo** | `/examples/demo` | Demo application                                    | `display`, `feedback`, `icons`, `layout`, `navigation` | `templ`                      |

### 3.2 Proposed Sub-Modules (go.mod groupings)

Rather than giving every package its own go.mod, we group by **dependency boundary and versioning cadence**:

| Sub-Module        | Contains                                         | go.mod Location        | Rationale                               |
| ----------------- | ------------------------------------------------ | ---------------------- | --------------------------------------- |
| **utils**         | `utils/`                                         | `utils/go.mod`         | Leaf module — foundation for everything |
| **svg**           | `svg/` (renamed from `internal/svg`)             | `svg/go.mod`           | Leaf module — SVG primitives            |
| **layout**        | `layout/`                                        | `layout/go.mod`        | Leaf module — zero internal deps        |
| **icons**         | `icons/`                                         | `icons/go.mod`         | Thin dependency on svg only             |
| **core**          | `feedback/`, `forms/`, `display/`, `navigation/` | `core/go.mod`          | UI components with shared deps          |
| **htmx**          | `htmx/`                                          | `htmx/go.mod`          | Optional HTMX integration layer         |
| **examples/demo** | `examples/demo/`                                 | `examples/demo/go.mod` | Demo app — depends on everything        |

Wait — this creates a `core` module with 4 packages, which is just a smaller monolith. Let me reconsider.

### 3.3 Revised: Per-Package Sub-Modules

Since the dependency graph is clean and acyclic, the most honest approach is **one go.mod per public package**:

| Sub-Module      | go.mod Path            | Prod Internal Deps                                     |
| --------------- | ---------------------- | ------------------------------------------------------ |
| `utils`         | `utils/go.mod`         | _(none)_                                               |
| `svg`           | `svg/go.mod`           | _(none)_                                               |
| `layout`        | `layout/go.mod`        | _(none)_                                               |
| `icons`         | `icons/go.mod`         | `svg`                                                  |
| `feedback`      | `feedback/go.mod`      | `utils`, `svg`                                         |
| `forms`         | `forms/go.mod`         | `utils`                                                |
| `display`       | `display/go.mod`       | `utils`, `icons`, `svg`                                |
| `navigation`    | `navigation/go.mod`    | `utils`, `svg`                                         |
| `htmx`          | `htmx/go.mod`          | `feedback`                                             |
| `examples/demo` | `examples/demo/go.mod` | `display`, `feedback`, `icons`, `layout`, `navigation` |

**Total: 10 sub-modules** (8 original packages + 1 promoted `svg` + 1 examples/demo)

### 3.4 DAG Verification

The proposed dependency graph is a strict DAG:

```
Layer 0 (leaves):   utils    svg    layout
                      \      /\
Layer 1:              forms  /  icons
                       \    /     \
Layer 2:          feedback  navigation  display
                     |
Layer 3:            htmx

Layer 4:          examples/demo  (depends on display, feedback, icons, layout, navigation)
```

**Cycle check:** Every edge points downward. No upward or lateral dependencies. ✅

### 3.5 Module Paths

Each sub-module gets its own Go module path:

| Sub-Module    | Module Path                                             |
| ------------- | ------------------------------------------------------- |
| utils         | `github.com/larsartmann/templ-components/utils`         |
| svg           | `github.com/larsartmann/templ-components/svg`           |
| layout        | `github.com/larsartmann/templ-components/layout`        |
| icons         | `github.com/larsartmann/templ-components/icons`         |
| feedback      | `github.com/larsartmann/templ-components/feedback`      |
| forms         | `github.com/larsartmann/templ-components/forms`         |
| display       | `github.com/larsartmann/templ-components/display`       |
| navigation    | `github.com/larsartmann/templ-components/navigation`    |
| htmx          | `github.com/larsartmann/templ-components/htmx`          |
| examples/demo | `github.com/larsartmann/templ-components/examples/demo` |

**Import paths for consumers remain identical.** A consumer importing `github.com/larsartmann/templ-components/display` will resolve it the same way — `go.work` handles local development, and the Go module proxy handles published versions.

---

## 4. Replace / Workspace Strategy

**Recommendation: `go.work` at repo root.**

```
go 1.26.2

use (
    ./utils
    ./svg
    ./layout
    ./icons
    ./feedback
    ./forms
    ./display
    ./navigation
    ./htmx
    ./examples/demo
)
```

**Why go.work over replace directives:**

- Cleaner — no `replace` directives polluting each go.mod
- `go work sync` keeps everything consistent
- Consumers of published modules never see `go.work` (ignored by proxy)
- Standard Go tooling handles workspace mode natively

**Rules:**

- No `replace` directives in any go.mod
- `go.work` is committed to the repo (not gitignored)
- Each go.mod is self-contained and clean

---

## 5. Test Dependency Isolation

| Module        | Production Deps                              | Test-Only Deps |
| ------------- | -------------------------------------------- | -------------- |
| utils         | _(none)_                                     | _(none)_       |
| svg           | _(none)_                                     | utils          |
| layout        | _(none)_                                     | utils          |
| icons         | svg                                          | utils          |
| feedback      | utils, svg                                   | _(none)_       |
| forms         | utils                                        | _(none)_       |
| display       | utils, icons, svg                            | feedback       |
| navigation    | utils, svg                                   | _(none)_       |
| htmx          | feedback                                     | utils          |
| examples/demo | display, feedback, icons, layout, navigation | _(none)_       |

**Test dependencies in Go are handled by the module system automatically** — `_test.go` files can import any module listed in the `require` block. Since all modules are in the workspace, test cross-imports work without explicit `require` directives.

**Strategy:** No special isolation needed. The workspace ensures all modules are available for testing. Each module's go.mod only lists its production dependencies.

---

## 6. Interface Extraction

**Not needed for this project.**

templ-components is a component library, not an application. Every package exposes concrete types and functions — there are no interfaces to extract, no implementations to separate. The module boundaries ARE the interfaces.

The only consideration is `utils.BaseProps` — it's embedded in every component's Props struct. After modularization, `utils` becomes a leaf module, and all component modules depend on it. This is the correct direction: leaf → consumer.

---

## 7. Versioning Strategy

**Recommendation: Shared version (single git tag).**

| Factor            | Assessment                                                 |
| ----------------- | ---------------------------------------------------------- |
| Single team       | Yes                                                        |
| Tight coupling    | `utils.BaseProps` is embedded everywhere — changes cascade |
| Consumer profile  | Go developers importing specific packages                  |
| Publishing target | Go module proxy                                            |

**Rationale:** The components share `utils.BaseProps` and `utils.Class()` — a breaking change in utils affects every module. Shared versioning avoids version matrix hell where `display@v1.3` requires `utils@v1.1+`.

**Tag format:** `v1.2.3` at repo root. All modules bump together.

**Future option:** If `icons` stabilizes and rarely changes, it could graduate to independent versioning later. But premature independent versioning adds complexity without benefit.

---

## 8. Migration Strategy (Ordered Steps)

1. **Create `go.work`** at repo root (empty use directive, populate as we go)
2. **Extract `utils`** — create `utils/go.mod`, keep `utils/` in place
3. **Promote `internal/svg` → `svg`** — rename directory, create `svg/go.mod`, update all importers (4 files: `display/helpers.templ`, `feedback/loading.templ`, `navigation/pagination.templ`, `icons/icon.templ`)
4. **Extract `layout`** — create `layout/go.mod` (zero internal deps)
5. **Extract `icons`** — create `icons/go.mod`, add `svg` dependency
6. **Extract `feedback`** — create `feedback/go.mod`, add `utils` + `svg` deps
7. **Extract `forms`** — create `forms/go.mod`, add `utils` dependency
8. **Extract `display`** — create `display/go.mod`, add `utils` + `icons` + `svg` deps
9. **Extract `navigation`** — create `navigation/go.mod`, add `utils` + `svg` deps
10. **Extract `htmx`** — create `htmx/go.mod`, add `feedback` dependency
11. **Extract `examples/demo`** — create `examples/demo/go.mod`, add all UI deps
12. **Convert root go.mod to thin aggregator** — keep root go.mod with no packages, only go.work uses it
13. **Update `go.work`** — add all modules to use directive
14. **Regenerate templ** — `find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./...`
15. **Verify** — `go work sync`, `go build ./...`, `go test ./...`

Each step is independently revertable (single commit).

### 8.1 Step 3 Detail: `internal/svg` → `svg` Rename

Files requiring import path update:

| File                          | Old Import                                             | New Import                                    |
| ----------------------------- | ------------------------------------------------------ | --------------------------------------------- |
| `display/helpers.templ`       | `github.com/larsartmann/templ-components/internal/svg` | `github.com/larsartmann/templ-components/svg` |
| `feedback/loading.templ`      | `github.com/larsartmann/templ-components/internal/svg` | `github.com/larsartmann/templ-components/svg` |
| `navigation/pagination.templ` | `github.com/larsartmann/templ-components/internal/svg` | `github.com/larsartmann/templ-components/svg` |
| `icons/icon.templ`            | `github.com/larsartmann/templ-components/internal/svg` | `github.com/larsartmann/templ-components/svg` |

Plus their generated `*_templ.go` counterparts (regenerated automatically).

### 8.2 Root go.mod Decision

**Keep root go.mod as a thin aggregator** — no Go source files at root, just:

```go
module github.com/larsartmann/templ-components

go 1.26.2
```

The `go.work` file handles module coordination. The root go.mod exists for:

- `go install` compatibility
- GitHub/Go module proxy discovery
- Backward compatibility with consumers who `go get github.com/larsartmann/templ-components`

**Alternative considered (rejected):** Remove root go.mod entirely. This breaks `go get github.com/larsartmann/templ-components` and complicates module proxy resolution.

---

## 9. Risk Assessment

| Risk                                    | Likelihood | Impact | Mitigation                                                              |
| --------------------------------------- | ---------- | ------ | ----------------------------------------------------------------------- |
| Import path breakage for consumers      | Medium     | High   | Module paths match current package paths exactly                        |
| `internal/svg` rename breaks consumers  | Low        | High   | `internal/` was unimportable anyway — zero external consumers           |
| Templ generate breaks with multi-module | Medium     | Medium | `templ generate` works per-package; each module generates independently |
| CI complexity increases                 | Medium     | Low    | `go.work` keeps single-command build/test working                       |
| Version drift between modules           | Low        | Medium | Shared version strategy prevents this                                   |
| Tailwind-merge-go version mismatch      | Low        | Low    | All modules pin same version via go.work                                |

---

## 10. Build System Impact

| Component      | Current                        | After                                     |
| -------------- | ------------------------------ | ----------------------------------------- |
| Build          | `go build ./...` (root)        | `go build ./...` (root via go.work)       |
| Test           | `go test ./...` (root)         | `go test ./...` (root via go.work)        |
| Templ generate | `templ generate ./...` (root)  | `templ generate ./...` (root via go.work) |
| Lint           | `golangci-lint run ./pkgs/...` | Same paths, works with go.work            |
| CI             | Single job                     | Single job (go.work handles multi-module) |
| flake.nix      | Not present                    | N/A                                       |

**Key insight:** `go.work` makes multi-module feel like single-module for all build/test commands. No CI changes needed.

---

## 11. Key Decisions

1. **Per-package sub-modules** over grouped modules — honest boundaries, maximum selectivity
2. **`go.work`** over `replace` directives — cleaner, standard tooling
3. **Shared versioning** over independent semver — `utils.BaseProps` coupling makes independent versioning fragile
4. **Promote `internal/svg` → `svg`** — necessary for cross-module access; zero external consumers affected
5. **Keep root go.mod** as thin aggregator — preserves `go get` and module proxy discovery
6. **10 sub-modules** (8 packages + promoted svg + examples/demo) — the natural boundary count

---

## 12. Self-Review (Phase 4)

### Questions Answered

| #   | Question                          | Answer                                                                                                                                                                          |
| --- | --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | What did I forget?                | Root go.mod should be kept (not removed) for module proxy compatibility                                                                                                         |
| 2   | What could be better?             | Consider grouping feedback+forms+display+navigation into a `core` module to reduce module count — rejected because it hides honest boundaries                                   |
| 3   | What could still improve?         | The `utils` module is small (3 files) — it could be inlined into each consumer — rejected because `BaseProps` and `Class()` are shared types that need a single source of truth |
| 4   | Split brains?                     | No. `utils.BaseProps` lives in one place and is embedded (not duplicated). No shared types need duplication.                                                                    |
| 5   | Right granularity?                | Yes. Each package has a clear domain purpose. 10 modules for a component library is appropriate — each is independently importable.                                             |
| 6   | Existing code reuse?              | The current `internal/svg` code moves directly — no new packages needed. All modules are existing packages promoted to modules.                                                 |
| 7   | Type model improvements?          | No changes needed. Current types are clean and well-structured.                                                                                                                 |
| 8   | Established libraries?            | `templ` and `tailwind-merge-go` are the right choices. No banned deps.                                                                                                          |
| 9   | Replace/workspace strategy works? | Yes. `go.work` with per-module go.mod is the standard Go multi-module pattern.                                                                                                  |
| 10  | Test deps isolated?               | Yes. Test-only imports (utils in icon/svg/htmx tests, feedback in display tests) are handled by workspace availability.                                                         |
| 11  | CI faster?                        | Marginal — `go.work` still builds all modules. True per-module CI requires separate GitHub Actions jobs per module. Future optimization, not blocking.                          |
| 12  | Versioning realistic?             | Shared versioning is correct given the BaseProps coupling. Independent versioning would create version matrix hell.                                                             |
