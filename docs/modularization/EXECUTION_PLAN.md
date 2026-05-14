# Execution Plan — templ-components Modularization

**Created:** 2026-05-14
**Status:** Ready for execution

---

## Overview

Convert templ-components from a single-module monolith to a 10-module Go workspace.
Ordered by dependency layer: leaves first, consumers last.

**Total tasks:** 14
**Estimated total effort:** 3–4 hours

---

## Task List

### Task 1: Create `go.work` skeleton

**Layer:** Foundation
**Effort:** 5 min
**Dependencies:** None

Create an empty `go.work` at repo root:

```
go 1.26.2
```

(Populate `use` directive incrementally as modules are extracted.)

**Verify:** `go work sync` succeeds (no-op at this point)
**Rollback:** `git rm go.work`

---

### Task 2: Extract `utils` module

**Layer:** 0 — Leaf
**Impact:** 1% → 51% (foundational — everything depends on utils)
**Effort:** 10 min
**Dependencies:** Task 1

1. Create `utils/go.mod`:
   ```
   module github.com/larsartmann/templ-components/utils
   go 1.26.2
   require github.com/a-h/templ v0.3.1001
   require github.com/Oudwins/tailwind-merge-go v0.2.1
   ```
2. Run `go mod tidy` in `utils/`
3. Add `./utils` to `go.work` use directive

**Verify:** `go build ./utils/...` passes
**Rollback:** `git rm utils/go.mod utils/go.sum`, revert go.work change

---

### Task 3: Promote `internal/svg` → `svg` module

**Layer:** 0 — Leaf
**Impact:** 1% → 51% (foundational — 4 packages depend on svg)
**Effort:** 20 min
**Dependencies:** Task 1

1. `git mv internal/svg svg`
2. Create `svg/go.mod`:
   ```
   module github.com/larsartmann/templ-components/svg
   go 1.26.2
   require github.com/a-h/templ v0.3.1001
   ```
3. Update package declaration in `svg/*.go` from `package svg` (already correct — it's already `package svg`)
4. Update all importers (4 .templ files + their generated \_templ.go files):
   - `display/helpers.templ`: `internal/svg` → `svg`
   - `feedback/loading.templ`: `internal/svg` → `svg`
   - `navigation/pagination.templ`: `internal/svg` → `svg`
   - `icons/icon.templ`: `internal/svg` → `svg`
5. Delete old generated `*_templ.go` files
6. Run `templ generate ./...`
7. Run `go mod tidy` in `svg/`
8. Add `./svg` to `go.work` use directive

**Verify:** `go build ./svg/...` passes; `go build ./icons/...` passes (first consumer)
**Rollback:** `git mv svg internal/svg`, revert import changes, regenerate templ

---

### Task 4: Extract `layout` module

**Layer:** 0 — Leaf
**Impact:** 4% → 64%
**Effort:** 10 min
**Dependencies:** Task 1

1. Create `layout/go.mod`:
   ```
   module github.com/larsartmann/templ-components/layout
   go 1.26.2
   require github.com/a-h/templ v0.3.1001
   ```
2. Run `go mod tidy` in `layout/`
3. Add `./layout` to `go.work` use directive

**Verify:** `go build ./layout/...` passes
**Rollback:** `git rm layout/go.mod layout/go.sum`, revert go.work

---

### Task 5: Extract `icons` module

**Layer:** 1
**Impact:** 4% → 64%
**Effort:** 15 min
**Dependencies:** Task 3 (svg)

1. Create `icons/go.mod`:
   ```
   module github.com/larsartmann/templ-components/icons
   go 1.26.2
   require github.com/larsartmann/templ-components/svg v0.0.0
   require github.com/a-h/templ v0.3.1001
   ```
2. Run `go mod tidy` in `icons/` (go.work will resolve the svg dependency)
3. Add `./icons` to `go.work` use directive

**Verify:** `go build ./icons/...` passes; `go test ./icons/...` passes
**Rollback:** `git rm icons/go.mod icons/go.sum`, revert go.work

---

### Task 6: Extract `feedback` module

**Layer:** 2
**Impact:** 4% → 64%
**Effort:** 15 min
**Dependencies:** Task 2 (utils), Task 3 (svg)

1. Create `feedback/go.mod`:
   ```
   module github.com/larsartmann/templ-components/feedback
   go 1.26.2
   require github.com/larsartmann/templ-components/utils v0.0.0
   require github.com/larsartmann/templ-components/svg v0.0.0
   require github.com/a-h/templ v0.3.1001
   ```
2. Run `go mod tidy` in `feedback/`
3. Add `./feedback` to `go.work` use directive

**Verify:** `go build ./feedback/...` passes; `go test ./feedback/...` passes
**Rollback:** `git rm feedback/go.mod feedback/go.sum`, revert go.work

---

### Task 7: Extract `forms` module

**Layer:** 1
**Impact:** 4% → 64%
**Effort:** 10 min
**Dependencies:** Task 2 (utils)

1. Create `forms/go.mod`:
   ```
   module github.com/larsartmann/templ-components/forms
   go 1.26.2
   require github.com/larsartmann/templ-components/utils v0.0.0
   require github.com/a-h/templ v0.3.1001
   ```
2. Run `go mod tidy` in `forms/`
3. Add `./forms` to `go.work` use directive

**Verify:** `go build ./forms/...` passes; `go test ./forms/...` passes
**Rollback:** `git rm forms/go.mod forms/go.sum`, revert go.work

---

### Task 8: Extract `display` module

**Layer:** 2
**Impact:** 20% → 80% (largest package, most consumers)
**Effort:** 15 min
**Dependencies:** Task 2 (utils), Task 3 (svg), Task 5 (icons)

1. Create `display/go.mod`:
   ```
   module github.com/larsartmann/templ-components/display
   go 1.26.2
   require github.com/larsartmann/templ-components/utils v0.0.0
   require github.com/larsartmann/templ-components/svg v0.0.0
   require github.com/larsartmann/templ-components/icons v0.0.0
   require github.com/a-h/templ v0.3.1001
   ```
2. Run `go mod tidy` in `display/`
3. Add `./display` to `go.work` use directive

**Verify:** `go build ./display/...` passes; `go test ./display/...` passes
**Rollback:** `git rm display/go.mod display/go.sum`, revert go.work

---

### Task 9: Extract `navigation` module

**Layer:** 2
**Impact:** 4% → 64%
**Effort:** 10 min
**Dependencies:** Task 2 (utils), Task 3 (svg)

1. Create `navigation/go.mod`:
   ```
   module github.com/larsartmann/templ-components/navigation
   go 1.26.2
   require github.com/larsartmann/templ-components/utils v0.0.0
   require github.com/larsartmann/templ-components/svg v0.0.0
   require github.com/a-h/templ v0.3.1001
   ```
2. Run `go mod tidy` in `navigation/`
3. Add `./navigation` to `go.work` use directive

**Verify:** `go build ./navigation/...` passes; `go test ./navigation/...` passes
**Rollback:** `git rm navigation/go.mod navigation/go.sum`, revert go.work

---

### Task 10: Extract `htmx` module

**Layer:** 3
**Impact:** 4% → 64%
**Effort:** 10 min
**Dependencies:** Task 6 (feedback)

1. Create `htmx/go.mod`:
   ```
   module github.com/larsartmann/templ-components/htmx
   go 1.26.2
   require github.com/larsartmann/templ-components/feedback v0.0.0
   require github.com/a-h/templ v0.3.1001
   ```
2. Run `go mod tidy` in `htmx/`
3. Add `./htmx` to `go.work` use directive

**Verify:** `go build ./htmx/...` passes; `go test ./htmx/...` passes
**Rollback:** `git rm htmx/go.mod htmx/go.sum`, revert go.work

---

### Task 11: Extract `examples/demo` module

**Layer:** 4 — Consumer
**Impact:** Polish (examples only)
**Effort:** 15 min
**Dependencies:** Task 8 (display), Task 6 (feedback), Task 5 (icons), Task 4 (layout), Task 9 (navigation)

1. Create `examples/demo/go.mod`:
   ```
   module github.com/larsartmann/templ-components/examples/demo
   go 1.26.2
   require (
       github.com/larsartmann/templ-components/display v0.0.0
       github.com/larsartmann/templ-components/feedback v0.0.0
       github.com/larsartmann/templ-components/icons v0.0.0
       github.com/larsartmann/templ-components/layout v0.0.0
       github.com/larsartmann/templ-components/navigation v0.0.0
       github.com/a-h/templ v0.3.1001
   )
   ```
2. Run `go mod tidy` in `examples/demo/`
3. Add `./examples/demo` to `go.work` use directive

**Verify:** `go build ./examples/demo/...` passes
**Rollback:** `git rm examples/demo/go.mod examples/demo/go.sum`, revert go.work

---

### Task 12: Update root go.mod

**Layer:** Foundation
**Impact:** 1% → 51%
**Effort:** 10 min
**Dependencies:** All previous tasks

1. Remove all Go source files from root (there are none — all code is in sub-directories)
2. Keep root `go.mod` as thin aggregator:
   ```
   module github.com/larsartmann/templ-components
   go 1.26.2
   ```
3. Remove root `go.sum` (no dependencies at root level)
4. Ensure `go.work` has all modules listed

**Verify:** `go work sync` succeeds; `go build ./...` passes at root
**Rollback:** Restore original root go.mod/go.sum from git

---

### Task 13: Full regeneration and verification

**Layer:** Validation
**Impact:** Critical — ensures everything works end-to-end
**Effort:** 15 min
**Dependencies:** Task 12

1. Delete all generated files: `find . -name '*_templ.go' -print0 | xargs -0 rm`
2. Regenerate: `templ generate ./...`
3. `go work sync`
4. `go build ./...`
5. `go test ./...`
6. `golangci-lint run ./display/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...`
   (Note: `internal/` directory may be empty after svg promotion — check)

**Verify:** All commands pass with zero errors
**Rollback:** `git reset --soft HEAD~1` to undo this step only

---

### Task 14: Update documentation

**Layer:** Polish
**Impact:** 20% → 80% (developer experience)
**Effort:** 20 min
**Dependencies:** Task 13

Update the following files to reflect the new multi-module structure:

1. **README.md** — Update installation instructions, add per-module import examples
2. **AGENTS.md** — Update build/test commands for multi-module workflow
3. **CONTEXT.md** — Update package layout and import graph
4. **go.work** — Final state with all 10 modules

**Verify:** Read through all updated docs for accuracy
**Rollback:** `git checkout HEAD~1 -- README.md AGENTS.md CONTEXT.md`

---

## Dependency Graph Between Tasks

```
Task 1 (go.work) ──────────────────────────────────────────────┐
   │                                                            │
   ├── Task 2 (utils) ──┬── Task 6 (feedback) ── Task 10 (htmx)│
   │                     ├── Task 7 (forms)                     │
   │                     ├── Task 8 (display) ─────────────────┤
   │                     └── Task 9 (navigation) ──────────────┤
   │                                                           │
   ├── Task 3 (svg) ────┬── Task 5 (icons) ── Task 8 (display)│
   │                     ├── Task 6 (feedback)                  │
   │                     ├── Task 8 (display)                   │
   │                     └── Task 9 (navigation)                │
   │                                                           │
   └── Task 4 (layout) ──── Task 11 (examples/demo) ──────────┤
                                                               │
Task 11 (examples/demo) ── Task 12 (root go.mod) ── Task 13 ── Task 14
```

---

## Pareto Impact Summary

| Tier          | Tasks                                                              | Impact                                                |
| ------------- | ------------------------------------------------------------------ | ----------------------------------------------------- |
| **1% → 51%**  | Tasks 1–3 (go.work, utils, svg)                                    | Foundation — everything else depends on these         |
| **4% → 64%**  | Tasks 4–7, 9–10 (layout, icons, feedback, forms, navigation, htmx) | High leverage — each unlocks a clean module boundary  |
| **20% → 80%** | Tasks 8, 12–13 (display, root go.mod, full verification)           | Broad value — largest package + end-to-end validation |
| **Polish**    | Tasks 11, 14 (examples/demo, docs)                                 | Developer experience                                  |

---

## Safety Net

- **Branch:** Create `modularize/go-work` branch before starting
- **One commit per task** — each independently revertable
- **Build verification after every task** — never accumulate broken state
- **Final full verification in Task 13** — catch any integration issues

---

## Risk Mitigations During Execution

| Risk                                            | Mitigation                                                        |
| ----------------------------------------------- | ----------------------------------------------------------------- |
| `templ generate` fails with go.work             | Run templ generate from root — it walks all sub-directories       |
| Import cycles appear                            | DAG verified in proposal — if cycles appear, boundaries are wrong |
| `go mod tidy` removes needed deps               | Run `go work sync` before `go mod tidy`                           |
| Test failures after extraction                  | Fix immediately — don't proceed to next task                      |
| `internal/` directory empty after svg promotion | Remove empty `internal/` directory                                |
