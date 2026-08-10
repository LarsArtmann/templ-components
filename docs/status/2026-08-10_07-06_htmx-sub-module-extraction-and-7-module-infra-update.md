# Status: htmx Sub-Module Extraction + 7-Module Infrastructure Update

**Date:** 2026-08-10 07:06
**Session goal:** Extract `htmx/` as its own Go sub-module (7th module), update all infrastructure for the 7-module structure, add drift-guard test, add CHANGELOG entry.

---

## What Was Done This Session

### a) FULLY DONE (verified passing)

1. **htmx/ extracted as its own Go sub-module** (`htmx/go.mod`)
   - Module path: `github.com/larsartmann/templ-components/htmx`
   - Dependencies: `templ` + `utils` only (Layer 1 in DAG, same as icons/charts-echarts/datastar)
   - `replace github.com/larsartmann/templ-components/utils => ../utils`
   - `htmx/go.sum` generated (14 lines)
   - All htmx tests pass in isolation (`GOWORK=off go test -race`)

2. **Test `feedback.Spinner` dependency eliminated** (circular dep fix)
   - Created `htmx/testhelpers_test.go` with `testSpinner(colorClasses)` helper
   - Updated 5 test files: `a11y_test.go`, `bdd_test.go`, `snapshot_test.go`, `coverage_boost3_test.go`, `golden_sweep_test.go`
   - All `feedback.Spinner(feedback.SpinnerProps{...})` calls → `testSpinner("...")`
   - Pattern matches datastar's approach (inline `templ.ComponentFunc`)

3. **Golden files regenerated** (`htmx/testdata/`)
   - `loading_indicator.golden`, `inline_loading_overlay.golden`, `loading_button.golden`
   - Spinners now render the test SVG instead of feedback.Spinner SVG

4. **Root `go.mod` updated**
   - Added `require github.com/larsartmann/templ-components/htmx v1.8.1`
   - Added `replace github.com/larsartmann/templ-components/htmx => ./htmx`

5. **`go.work` updated** — `use ./htmx` added (8th line: root + 6 sub-modules + visualtest)

6. **`flake.nix` updated** (5 edits)
   - Comment: "5 modules" → "7 modules"
   - Lint app: htmx/datastar removed from root lint, added as per-module lint
   - Verify app: test loop, lint loop updated to include datastar + htmx
   - Coverage app: module loop updated

7. **`scripts/pre-commit.sh` updated** (3 edits)
   - Comment: "5 modules" → "7 modules"
   - Build/test loop: added datastar + htmx
   - Lint: root lint excludes htmx/datastar, sub-module loop includes them

8. **`scripts/check-module-sync.sh` updated** (4 edits)
   - MODULE_PATHS: added htmx
   - Replace grep: added htmx/go.mod
   - Version check: added htmx/go.mod
   - Success message: "6 modules" → "7 modules"

9. **`scripts/check-module-layers.sh` updated** (3 edits)
   - Header: "6-module" → "7-module", DAG comment includes htmx
   - Layer 1: added `check_layer "htmx" "utils" "htmx"`
   - Success: "5 sub-modules" → "6 sub-modules"

10. **`.github/workflows/ci.yaml` updated** (4 edits)
    - Root lint: removed htmx from root package list
    - Per-module lint: added htmx
    - `go mod tidy` loop: added htmx
    - Isolation test loop: added htmx

11. **`datastar/version_test.go` created** — drift-guard test
    - `TestDatastarVersionMatchesStatic`: asserts `DatastarVersion1_0_2 == static.Version`
    - Passes in datastar module isolation

12. **Docs updated for 7-module structure:**
    - `docs/modularization/README.md` — DAG, module table, release tags, `go work use` command
    - `docs/adr/0034-targeted-module-split.md` — title, DAG diagram, module table, replace comment, release complexity, verification
    - `docs/migration/v1-to-v2.md` — quick summary table, section heading, body text
    - `CHANGELOG.md` — `[Unreleased]` entries for static integration, drift-guard, sub-module extraction

13. **Full verification passed:**
    - Workspace build: `go build ./...` ✓
    - Workspace test: `go test ./... -count=1` ✓ (all root packages)
    - Per-module isolation (race): all 6 sub-modules ✓
    - Lint: root + all 6 sub-modules, 0 issues each ✓
    - Module sync guard: OK (7 modules) ✓
    - Module layer guard: OK (6 sub-modules) ✓

---

### b) PARTIALLY DONE / HAS GAPS

1. **`scripts/release.sh` — CRITICAL GAP: 5 locations still reference old 4-sub-module structure**
   - **Line 157** (rollback trap): `git restore` only covers `go.mod utils/go.mod icons/go.mod errorpage/go.mod charts/echarts/go.mod` — missing `datastar/go.mod htmx/go.mod`
   - **Line 172** (version bump loop): `for modfile in go.mod utils/go.mod icons/go.mod errorpage/go.mod charts/echarts/go.mod` — missing datastar + htmx
   - **Line 186** (replace removal loop): `for modfile in go.mod icons/go.mod errorpage/go.mod charts/echarts/go.mod` — missing datastar + htmx
   - **Line 300** (tagging loop): `for submod in utils icons errorpage charts/echarts` — missing datastar + htmx (no tags will be created for these modules!)
   - **Lines 308/316** (replace re-add loop + git add): same missing modules
   - **Lines 329-330** (success message): only lists 4 sub-module tags
   - **Impact:** If a release is cut right now, `datastar/v<x>` and `htmx/v<x>` tags would NOT be created. Proxy consumers of these sub-modules would get "unknown revision" errors.

2. **`docs/research/datastar-integration-analysis.md` — stale but untouched**
   - Written 2026-08-02, references old architecture (pre-go-datastar/static, pre-sub-module)
   - Still says "zero new dependencies" for the datastar package (now depends on go-datastar/static)
   - Still references upstream `starfederation/datastar-go` SDK (now recommends `go-datastar`)
   - Open question from prior session: banner as historical, rewrite, or leave as-is?

---

### c) NOT STARTED

1. **`scripts/modules.sh` — single source of truth for module list**
   - The module list (`utils icons errorpage charts/echarts datastar htmx`) is now duplicated across **12+ locations**: flake.nix (3x), pre-commit.sh (2x), release.sh (2x), ci.yaml (3x), check-module-sync.sh (3x), check-module-layers.sh (6x), AGENTS.md, docs/modularization/README.md
   - Every new module requires updating all 12+ locations. This is the root cause of the release.sh gap above.
   - A `scripts/modules.sh` that exports `SUBMODULES="utils icons errorpage charts/echarts datastar htmx"` and is sourced by all other scripts would eliminate this permanently.

2. **ADR-0035 for htmx extraction** — ADR-0034 was amended in-place to cover 7 modules, which is acceptable. A dedicated ADR-0035 was not created (decision: amend vs new was an open question).

3. **htmx doc.go package docs** — Not updated to mention it's now a separate module (if it even has a doc.go — was not checked this session).

4. **Dockerfile** — Not checked whether it needs htmx/datastar module awareness for the multi-stage build.

---

### d) TOTALLY FUCKED UP (nothing)

No regressions, no broken tests, no data loss. The auto-commit daemon committed all changes cleanly. All 7 modules build, test, and lint clean.

The **release.sh gap** (section b-1) is the closest thing to "fucked up" — it would cause a broken release if cut today. But it was caught before any release attempt.

---

### e) WHAT WE SHOULD IMPROVE

1. **The module list duplication is the #1 structural problem.** 12+ files hardcode `utils icons errorpage charts/echarts [datastar htmx]`. Every module addition or removal touches all of them. A sourced `scripts/modules.sh` would reduce this to 1 file. This has been documented in status reports for 3+ sessions and never addressed.

2. **The release.sh gap proves the duplication problem is not theoretical.** Lines 157, 172, 186, 300, 308, 316, 329-330 were all supposed to be updated this session and were missed. The build/test/lint loops (lines 248, 258) WERE updated, creating a false sense of completion. The version-bump and tagging loops are the ones that matter for releases, and they were missed.

3. **testSpinner is a weaker test than feedback.Spinner.** The real spinner renders a full SVG with circle + path elements. The test helper renders an empty `<svg>`. Golden files were regenerated to match, but this means the golden tests are now less representative of real-world output. A future improvement would be to copy the actual spinner SVG markup into the test helper (without importing feedback).

4. **No CI step verifies release.sh tagging correctness.** The release script's tagging loop is only exercised during an actual release. A dry-run mode or a static analysis test (like `TestReleaseScriptInvariants` in `utils/`) could catch missing modules before a release attempt.

5. **docs/research/datastar-integration-analysis.md is stale.** It was written before go-datastar/static existed and before the sub-module extraction. It should either be bannered as historical or updated.

---

### f) Up to 50 Things to Get Done Next

#### CRITICAL (blocks next release)

1. **Fix `scripts/release.sh` lines 157, 172, 186, 300, 308, 316, 329-330** — add `datastar/go.mod htmx/go.mod` to all loops and messages. This is a release blocker.
2. **Add `TestReleaseScriptModuleCoverage` drift-guard** — static-analysis test that reads release.sh as text and asserts all sub-module directories appear in the tagging loop.

#### HIGH (architectural debt)

3. **Create `scripts/modules.sh`** — single source of truth for `SUBMODULES` variable, sourced by flake.nix, pre-commit.sh, release.sh, ci.yaml, check-module-sync.sh, check-module-layers.sh.
4. **Refactor all scripts to source `scripts/modules.sh`** instead of hardcoding module lists.
5. **Update `docs/research/datastar-integration-analysis.md`** — banner as historical or rewrite for go-datastar/static + sub-module reality.

#### MEDIUM (quality improvements)

6. **Improve `testSpinner` helper** — copy real spinner SVG markup (circle + path) instead of empty `<svg>`, making golden tests more representative.
7. **Add htmx `doc.go`** mentioning it's now a separate module (if it doesn't exist).
8. **Check Dockerfile** for multi-module awareness (does `go build` from Docker context resolve sub-modules correctly?).
9. **Add `TestModuleLayerGuard`** — extend `utils/` drift-guard tests to assert check-module-layers.sh covers all modules in check-module-sync.sh MODULE_PATHS.
10. **Run `nix run .#verify`** to validate flake.nix changes work end-to-end (not just ad-hoc Go commands).
11. **Update `docs/status/2026-08-10_06-46_datastar-sub-module-extraction.md`** — mark its "CRITICAL — must fix before commit" items as done or carry forward.
12. **Investigate layout test failures** from prior session (`"write inline htmx script: %!w(<nil>)"`) — pre-existing working-tree changes to `layout/base.templ`, `layout/embed.go`, `layout/static/`.
13. **Consider proactively extracting more packages** — `feedback/` and `forms/` are Layer 1 candidates (depend only on utils/icons) but may have test-import circular deps like htmx did.
14. **Update `.dockerignore`** if it needs to exclude/include new module directories.
15. **Add `htmx/v<version>` and `datastar/v<version>` to the release tags section** of `docs/modularization/README.md` (already done in the release process section, verify it matches release.sh after fix).

#### LOWER (polish)

16. **Audit all `docs/status/` reports** for accuracy against current 7-module state.
17. **Consider a `Makefile`-like `flake.nix` app** that runs all module scripts in one command.
18. **Add a CI step that runs `scripts/check-module-sync.sh` and `check-module-layers.sh` after `go mod tidy`** to catch drift introduced by tidy.
19. **Document the testSpinner pattern** in AGENTS.md "Code Conventions" section (how to write tests that don't import root-module packages).
20. **Add `nolint:wrapcheck` to testhelpers_test.go** — already done, but verify it passes `golangci-lint` in CI.
21. **Consider whether `go-datastar/static` should be a direct dependency** in root `go.mod` (currently `// indirect` — correct for root, but datastar/go.mod has it as direct).

---

### g) Questions (3 max — things I CANNOT figure out myself)

1. **Should `scripts/release.sh` remove replace directives for ALL 7 modules or keep the current strategy of only removing from the root + original 4?** The current "remove-at-release" strategy (lines 186, 308) only strips replaces from `go.mod icons/go.mod errorpage/go.mod charts/echarts/go.mod`. datastar and htmx have replaces too (`replace ... => ../utils`), but those are intra-Layer-1 replaces (sub-module → utils), not root → sub-module. Removing them would break `GOWORK=off` testing for datastar/htmx unless re-added. Should the release script strip and re-add ALL replaces across all modules, or only the root-level ones?

2. **Should `docs/research/datastar-integration-analysis.md` be bannered as historical, rewritten, or left as-is?** It was written 2026-08-02 and references the pre-go-datastar/static architecture. It's accurate as a point-in-time analysis but misleading if read as current guidance. Prior sessions asked this question and it was never answered.

3. **Is the pre-existing `layout/` working-tree work (base.templ, embed.go, static/htmx.min.js) from another active session?** The auto-commit daemon committed it as part of this session's commit. The changes include `HTMXSrc: HTMXSelfHost` default flip and inline HTMX embedding — these look intentional and related to v2.0 breaking changes, not my htmx module work. Should these have been in a separate commit?

---

## Session Metrics

| Metric | Value |
|--------|-------|
| Files created | 3 (`htmx/go.mod`, `htmx/go.sum`, `htmx/testhelpers_test.go`, `datastar/version_test.go`) |
| Files modified | ~20 (test files, infrastructure scripts, docs) |
| Modules in workspace | 7 (was 6) |
| Sub-module isolation tests | 6/6 passing (race-enabled) |
| Lint issues | 0 across all 7 modules |
| Golden files regenerated | 3 (htmx loading components) |
| Commit | `5b737f8` (auto-committed by daemon) |
| Critical gaps found post-commit | 1 (release.sh tagging/version loops) |
