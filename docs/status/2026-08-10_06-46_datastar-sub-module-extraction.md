# Status: datastar Sub-Module Extraction — Wiring & Self-Review

> **Date:** 2026-08-10 06:46
> **Session scope:** Extract `datastar/` into its own Go module (matching
> `charts/echarts`, `icons`, `errorpage`, `utils`), wire it into all
> build/CI/lint/release infrastructure, and self-review.
>
> **Predecessor:** `docs/status/2026-08-10_06-32_go-datastar-static-integration.md`
> (same session — wired `go-datastar/static` as the version source of truth).

---

## a) FULLY DONE

| #  | Task                                                                                                                   | Verification                                                |
| -- | ---------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------- |
| 1  | Created `datastar/go.mod` (requires: templ, go-datastar/static, utils; replaces utils → ../utils)                      | `go mod tidy` clean                                         |
| 2  | Removed `feedback` import from `bdd_test.go` (would be a circular dep: feedback is in root module)                     | Replaced with inline `templ.ComponentFunc`                  |
| 3  | Updated root `go.mod`: added datastar require + replace                                                                | `go build ./...` passes                                     |
| 4  | Updated `go.work`: added `use ./datastar`                                                                              | Workspace resolves correctly                                |
| 5  | Updated `scripts/check-module-sync.sh`: added datastar to 4 locations (paths, grep, version loop, count)               | Script passes: "6 modules"                                  |
| 6  | Updated `scripts/check-module-layers.sh`: added datastar to Layer 1 DAG check                                          | Script passes: "5 sub-modules"                              |
| 7  | Updated `.github/workflows/ci.yaml`: moved datastar from root lint to per-module lint; added to tidy + isolation loops | Verified by reading diff                                    |
| 8  | Updated `AGENTS.md`: module table, import graph (6 modules), lint command, build commands                              | Verified                                                    |
| 9  | Datastar sub-module isolation test passes (GOWORK=off, -race)                                                          | All tests pass                                              |
| 10 | All 6 sub-module isolation tests pass                                                                                  | utils, icons, errorpage, charts/echarts, datastar all green |
| 11 | Datastar lint passes (0 issues)                                                                                        | `golangci-lint run` clean                                   |
| 12 | Root build passes                                                                                                      | `go build ./...` clean                                      |

### Module DAG after this session

```
Layer 0: utils                          (leaf: BaseProps, Class, svg, cdn, golden)
Layer 1: icons → utils                   (102 SVG icons)
Layer 1: charts/echarts → utils          (ECharts adapter)
Layer 1: datastar → utils, go-datastar/static  ← NEW: Datastar components
Layer 2: errorpage → utils, icons        (go-error-family integration)
Layer 3: root → all above                (display, feedback, forms, layout, etc.)
```

### What the datastar sub-module looks like

```go
// datastar/go.mod
module github.com/larsartmann/templ-components/datastar

require (
    github.com/a-h/templ v0.3.1020
    github.com/larsartmann/go-datastar/static v0.1.0
    github.com/larsartmann/templ-components/utils v1.8.1
)
require github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
replace github.com/larsartmann/templ-components/utils => ../utils
```

---

## b) PARTIALLY DONE

| # | Task                        | What's missing                                                                                                                                                                                                                                                                                                                                                                                 |
| - | --------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **`flake.nix`**             | **FORGOT TO UPDATE.** 6 locations still reference the old 5-module structure. `nix run .#lint`, `nix run .#verify`, `nix run .#test` all miss the datastar module. Lines 56-57 (comment), 88 (root lint includes datastar), 89-96 (missing datastar lint step), 140/149/169 (module loops missing datastar). This is a **CI-critical miss** — the Nix lint/test commands won't cover datastar. |
| 2 | **`scripts/pre-commit.sh`** | **FORGOT TO UPDATE.** Lines 25, 34, 37 still reference old module list. Pre-commit lint won't cover datastar.                                                                                                                                                                                                                                                                                  |
| 3 | **`scripts/release.sh`**    | **FORGOT TO UPDATE.** Lines 248, 258 still reference old module list. Release verification won't cover datastar.                                                                                                                                                                                                                                                                               |

---

## c) NOT STARTED

| # | Task                                                                                     |
| - | ---------------------------------------------------------------------------------------- |
| 1 | **CHANGELOG entry** for the sub-module extraction                                        |
| 2 | **`docs/modularization/README.md`** — says "5 modules", needs "6 modules"                |
| 3 | **ADR-0034** — says "5 modules", needs amendment for the 6th                             |
| 4 | **`docs/migration/v1-to-v2.md`** — says "5-module workspace"                             |
| 5 | **Drift-guard test** — `TestDatastarVersionMatchesStatic` (from previous status)         |
| 6 | **Update status report** from 06:32 — it says AGENTS.md import graph is stale; now fixed |

---

## d) TOTALLY FUCKED UP

### Miss 1: Forgot `flake.nix` entirely

**What happened:** I updated CI, AGENTS.md, module scripts, and go.work — but
completely missed `flake.nix`. The Nix flake is the canonical build/lint/test
entry point for local development. It has 6 locations that reference the old
5-module structure and 3 module loops (`utils icons errorpage charts/echarts`)
that skip datastar entirely.

**Impact:** `nix run .#lint` won't lint datastar. `nix run .#verify` won't
verify it. `nix run .#test` won't test it. A developer using Nix (the
recommended workflow) will get false confidence that everything passes.

**Root cause:** Tunnel vision on Go-native tooling (`go build`, `golangci-lint`).
The flake.nix wraps these but I didn't trace the full dependency chain from
"what needs updating" to "what invokes lint/test."

**Severity:** High — must be fixed before commit.

### Miss 2: Forgot `scripts/pre-commit.sh`

**What happened:** Same pattern — the pre-commit script has two module loops
(lines 25, 37) that iterate over `utils icons errorpage charts/echarts` without
datastar.

**Impact:** Pre-commit lint/test won't cover datastar. A commit could pass
pre-commit with broken datastar code.

**Root cause:** Same as flake.nix — didn't enumerate all files that contain
module lists.

### Miss 3: Forgot `scripts/release.sh`

**What happened:** The release script has two module loops (lines 248, 258)
that skip datastar.

**Impact:** A release would ship without verifying datastar in isolation.

**Root cause:** Same — incomplete enumeration of build automation files.

### Pattern: "Module list" is duplicated in 8+ places

The root cause of all 3 misses is architectural: the list of sub-modules is
copy-pasted across `flake.nix`, `scripts/pre-commit.sh`, `scripts/release.sh`,
`.github/workflows/ci.yaml`, `scripts/check-module-sync.sh`, and
`scripts/check-module-layers.sh`. Adding a module means updating all of them.
There should be a single source of truth (e.g., a file listing modules, or a
`go work edit -json` invocation that auto-discovers).

---

## e) WHAT WE SHOULD IMPROVE

### Process

1. **Maintain a "files that reference module count" checklist.** When adding a
   module, grep the entire repo for `charts/echarts` (the last module added)
   and update every match. This is what I should have done before declaring done.

2. **Always run `nix run .#verify` as the final check, not `go build ./...`.**
   The Nix verify app is the closest thing to what CI runs. `go build` only
   checks workspace compilation — it doesn't lint or test. I would have caught
   the flake.nix miss immediately.

3. **Declare a single source of truth for the module list.** Every script and
   CI step reads from it. Eliminates this entire class of "forgot to update"
   bugs permanently.

### Architecture

4. **The `feedback` → `datastar` test dependency was a hidden coupling.** The
   bdd_test imported `feedback.Spinner` to test the Indicator's custom spinner
   slot. This worked only because both were in the root module. Extracting
   datastar exposed it. The fix (inline `templ.ComponentFunc`) is actually
   better — the test is now self-contained and doesn't depend on a specific
   spinner implementation.

5. **Sub-module extraction pattern is proven and repeatable.** The
   charts/echarts precedent made this straightforward. The only new wrinkle was
   the external dependency (`go-datastar/static`), which is zero-transitive-dep
   and therefore doesn't complicate the sub-module's go.sum.

---

## f) Next actions (up to 50)

### CRITICAL — must fix before commit

1. **Update `flake.nix`** — 6 locations: comment (56-57), root lint (88), lint sub-module list (89-96), test loop (140), verify lint (148-149), verify test loop (169)
2. **Update `scripts/pre-commit.sh`** — lines 25, 34, 37: add datastar to module loops, remove from root lint list
3. **Update `scripts/release.sh`** — lines 248, 258: add datastar to module loops
4. **Verify `nix run .#verify` passes** with all 6 modules
5. **Verify `nix run .#lint` passes** with all 6 modules

### Should fix soon

6. Update `docs/modularization/README.md` — "5 modules" → "6 modules" (3 locations)
7. Update `AGENTS.md` lines 3, 5 — "5-module" → "6-module"
8. Update `AGENTS.md` line 348 — "all 5 modules" → "all 6 modules"
9. Amend ADR-0034 or add a note: "datastar was extracted as the 6th module on 2026-08-10"
10. Update `docs/migration/v1-to-v2.md` — "5-module workspace" references
11. Add CHANGELOG `[Unreleased]` entry for sub-module extraction
12. Verify `nix run .#test` covers datastar

### From previous status report (still open)

13. Add drift-guard test: `TestDatastarVersionMatchesStatic`
14. Check FEATURES.md datastar section for stale SDK references
15. Add banner to research doc re: go-datastar (from Q1 of previous report)
16. Investigate pre-existing layout test failures (from Q2 of previous report)

### Architecture improvements

17. Create `scripts/modules.sh` — single source of truth for module list
18. Refactor all scripts to source `modules.sh` instead of hardcoding lists
19. Add a CI step: "verify module count matches go.work `use` directives"
20. Consider extracting `htmx/` as a sub-module too (same pattern, same benefits)
21. Consider a `go.work` CI check that verifies all `use` directives have corresponding go.mod

### Testing

22. Add contract test for datastar sub-module (like internal/contract for other packages)
23. Run `nix flake check` to verify flake formatting
24. Run `nix fmt` to format any new/changed files
25. Verify Dockerfile build still works with the new module
26. Add `nix run .#coverage` — verify datastar coverage is measured

### Documentation

27. Update ROADMAP.md if it references module count
28. Update `docs/modularization/README.md` contributor setup instructions
29. Add datastar sub-module to the module dependency diagram in docs
30. Consider whether ADR-0034 should be updated or a new ADR-0035 created

### Cleanup

31. Verify `.gitignore` handles `datastar/go.sum` correctly (should be committed)
32. Check if `go.work.sum` needs regeneration after adding datastar
33. Run `golangci-lint` on root module to verify datastar removal from root lint list didn't break anything
34. Update the previous status report (06:32) with a pointer to this one
35. Verify the `skill/SKILL.md` component catalogue doesn't reference module structure

---

## g) Questions I cannot answer myself

### Q1: Should I fix flake.nix + pre-commit + release.sh right now, or are those files mid-edit by another session?

The working tree already has uncommitted changes to many files (including
`scripts/release.sh` and `.github/workflows/ci.yaml`). If another session is
actively editing these, my changes could conflict. Should I proceed with the
fixes immediately, or wait?

### Q2: Should I create ADR-0035 for the datastar extraction, or amend ADR-0034?

ADR-0034 documents the original 5-module split. Datastar is a natural addition
(same pattern, same justification). Should I:

- **(A)** Amend ADR-0034 in-place (it's a living doc)
- **(B)** Create ADR-0035 "datastar sub-module extraction" referencing 0034
- **(C)** Just note it in CHANGELOG (the extraction is straightforward)

### Q3: Should `htmx/` follow the same sub-module extraction?

The `htmx` package is a close parallel to `datastar` — same "attribute emitter"
pattern, same "runtime injection via CDN" approach, no external deps today.
But if HTMX ever gets a `go-htmx/static` equivalent (like we just did for
datastar), extracting it as a sub-module would isolate that dependency. Should
I do it now proactively, or wait for an actual external dep?
