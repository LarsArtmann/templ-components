<!-- AUTO-UPDATED 2026-07-10: Retrospective status overlay -->

> ## 🔔 Update Notice — 2026-07-10
>
> This report is **historical**. Many items listed as "open", "todo", or "broken" below
> have since been **fixed and verified**. Do not act on open items without first checking
> [TODO_LIST.md](../../TODO_LIST.md) for current status.
>
> **Key fixes completed since this report:**
>
> - ✅ All 7 P0 bugs fixed (InlineLoadingOverlay a11y, SanitizeID mismatch, FromError fallback,
>   Footer BaseProps, ErrorPage/NotFound404 `<main>` landmark, CSRFTokenName, grid-rows verified)
> - ✅ `encoding/json/v2` purged from all production code + pre-commit guard added
> - ✅ Motion constants centralized in `utils/motion.go`, wired into 13 components
> - ✅ `FamilyFromErrorFamily` → `FromErrorFamily` (old name kept as deprecated alias)
> - ✅ `icons.IconRTL()` + CSS for directional icon RTL mirroring
> - ✅ 33 regression tests added (htmx, errorpage, layout, navigation, feedback, display)
> - ✅ Dark golden test infrastructure (badge/card/button)
> - ✅ CHANGELOG consolidated, ROADMAP updated, migration guide created
> - ✅ All 14 packages pass, 0 lint issues
>
> **Canonical source of truth:** [TODO_LIST.md](../../TODO_LIST.md) (52 items, 37 ✅ done, 12 deferred/blocked)

---

# Status Report — Modularization Session

> **Updated:** 2026-07-06 (post-v0.8.0). **OUTCOME: BRANCH NEVER MERGED.**

**Date:** 2026-07-05 20:39
**Branch:** `modularize/strategic-split` (pushed to origin, **never merged**)
**Base:** `origin/master` at `f81ae66`
**Commits on branch:** 12 (7 modularization + 5 fixes from self-review)

> **UPDATE NOTE (2026-07-06):** The modularization branch was **abandoned**.
> Commit `a0dbae7` ("fix(docs): correct split-brain — master is single-module not multi-module")
> corrected the false documentation that claimed a 6-module workspace. `master` remains a single
> Go module. The branch still exists on origin (`remotes/origin/modularize/strategic-split`) but
> there is no `go.work`, no `svg/go.mod`, and no sub-module go.mod files on master. AGENTS.md
> documents: "The split may be re-attempted post-v1.0 if the package graph warrants it."
>
> **Current state:** v0.8.0, single-module, 14 packages, 575 test cases, 13/13 green, 0 lint.

---

## A) FULLY DONE ✅

### Modularization Execution (Phases 1–7 of go-modularize skill)

| Item                                                                                     | Status  | Commit            |
| ---------------------------------------------------------------------------------------- | ------- | ----------------- |
| Phase 1: Current state analysis — 1 go.mod monolith, 10 public + 2 internal packages     | ✅ Done | `7e42a53`         |
| Phase 2: Dependency graph mapping — production + test imports for all packages           | ✅ Done | `7e42a53`         |
| Phase 3: Proposal HTML written — `docs/modularization/2026-07-05_PROPOSAL.html`          | ✅ Done | `7e42a53`         |
| Phase 4: Brutal self-review — 15-point checklist, proposal validated                     | ✅ Done | `7e42a53`         |
| Phase 5: Execution plan HTML — `docs/modularization/2026-07-05_EXECUTION_PLAN.html`      | ✅ Done | `7e42a53`         |
| Phase 6: svg sub-module — `internal/svg` promoted to public `svg/`                       | ✅ Done | `e126d60`         |
| Phase 6: utils sub-module — extracted as leaf dependency                                 | ✅ Done | `e126d60`         |
| Phase 6: icons sub-module — extracted (depends on svg + utils)                           | ✅ Done | `04b83a1`         |
| Phase 6: errorpage sub-module — extracted (isolates go-error-family)                     | ✅ Done | `04b83a1`         |
| Phase 6: examples/demo sub-module — extracted as leaf                                    | ✅ Done | `04b83a1`         |
| Phase 6: Root go.mod updated — requires + replaces for all sub-modules                   | ✅ Done | `04b83a1`         |
| Phase 6: go.work created — 6 modules coordinated                                         | ✅ Done | `e126d60`         |
| Phase 6: All imports updated — `internal/svg` → `svg` in 7 .templ files + generated code | ✅ Done | `e126d60`         |
| Phase 6: Templ regenerated — all 59 `*_templ.go` files re-generated                      | ✅ Done | verified, no diff |
| Phase 7: Full build + test + lint — all 13 test packages green                           | ✅ Done | verified          |
| Phase 7: GOWORK=off isolation test — all 6 modules build standalone                      | ✅ Done | verified          |

### Self-Review Fixes (brutal self-review identified 2 critical + 2 high + 3 medium issues)

| Fix                                                                                 | Status  | Commit    |
| ----------------------------------------------------------------------------------- | ------- | --------- |
| CI workflow updated — tests all modules + GOWORK=off verification                   | ✅ Done | `96d4b29` |
| Release script — per-module tags + placeholder→real version bump                    | ✅ Done | `12dd118` |
| Old modularization docs archived — `docs/modularization/archive/` with README       | ✅ Done | `c837491` |
| Brutal self-review report — `docs/reviews/2026-07-05_18-32_brutal-self-review.html` | ✅ Done | `c837491` |
| ADRs updated — 0001/0002/0004 reference `svg` instead of `internal/svg`             | ✅ Done | `5405675` |
| .gitignore — demo binary added, go.work un-ignored                                  | ✅ Done | `5405675` |
| AGENTS.md — multi-module structure documented, build/test/lint commands updated     | ✅ Done | `e126d60` |
| CHANGELOG.md — [Unreleased] entry for multi-module workspace                        | ✅ Done | `e126d60` |
| flake.nix — lint paths include `./svg/...`                                          | ✅ Done | `e126d60` |

### Final 6-Module Structure

```
github.com/larsartmann/templ-components/
├── go.mod          (root: display, feedback, forms, layout, navigation, htmx, internal/)
├── go.work         (6 modules)
├── svg/go.mod      (leaf: SVG primitives — promoted from internal/svg)
├── utils/go.mod    (leaf: BaseProps, Class(), EnsureID, test helpers)
├── icons/go.mod    (depends: svg, utils — 101 named icons, zero CSS deps)
├── errorpage/go.mod (depends: icons, utils — isolates go-error-family)
└── examples/demo/go.mod (depends: root, icons, errorpage, svg, utils)
```

---

## B) PARTIALLY DONE ⚠️

> **ALL ITEMS BELOW ARE MOOT** — the branch was never merged. These were "done" on the
> branch but never landed on `master`. Listed for historical reference only.

| Item                                     | What was done on the branch                                           | What's missing (and now moot)                                                        |
| ---------------------------------------- | --------------------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| **Release script end-to-end test**       | Script updated with placeholder→version replacement + per-module tags | NEVER tested. Single-module `scripts/release.sh` used for v0.7.0 and v0.8.0 instead. |
| **CI workflow on actual GitHub Actions** | YAML written with sub-module iteration + GOWORK=off                   | NEVER run on CI. The actual CI uses the single-module workflow.                      |
| **Consumer experience verification**     | GOWORK=off builds pass for all modules                                | NEVER tested from outside. Single-module consumers use standard `go get`.            |
| **Merge to master**                      | Branch pushed to origin                                               | ❌ **NEVER MERGED.** Branch abandoned. Documentation corrected in `a0dbae7`.         |

---

## C) NOT STARTED ❌

> **ALL ITEMS BELOW ARE MOOT** — modularization was abandoned.

| Item                                   | Why (moot)                                                                                                       |
| -------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| **PR creation**                        | ❌ Never created. Branch abandoned.                                                                              |
| **First release with per-module tags** | ❌ Never happened. v0.7.0/v0.8.0 released as single-module with standard tags.                                   |
| **Per-module go.sum audit**            | ❌ Moot — no sub-modules exist.                                                                                  |
| **flake.nix multi-module apps**        | ❌ Moot — flake.nix uses single-module commands.                                                                 |
| **README.md update for multi-module**  | ❌ Moot — README correctly describes single-module `go get`.                                                     |
| **icons-only-adoption.md update**      | ✅ Done differently — doc now describes icons as naturally having zero CSS deps, not as a standalone sub-module. |
| **CONTRIBUTING.md update**             | ❌ Moot — single-module build commands are correct.                                                              |

---

## D) TOTALLY FUCKED UP 💥

### BuildFlow vs Branch Discipline

**This was the single biggest execution problem.** BuildFlow's pre-commit hook and background
process repeatedly:

1. **Switched branches mid-session** — at least 8 times, BuildFlow switched from
   `modularize/strategic-split` to `master` without warning. Commits landed on the wrong branch.
   Had to cherry-pick, reset, and recreate the branch multiple times.
2. **Deleted the modularize branch entirely** — twice. Had to recreate from reflog.
3. **Created untracked files** (`errorpage/notfound404*.go`, `forms/aria.go`, `forms/ids.go`)
   that caused build failures (duplicate declarations). These came from BuildFlow's
   `govalid-generate` or similar tooling running in the background.
4. **Reverted .gitignore changes** — the `go.work` un-ignore kept getting overridden.
5. **Stole file modifications** — applied formatting changes to files I was editing, causing
   race conditions where my edits were silently reverted.

**Root cause:** BuildFlow runs as a background process that operates on whatever branch is
checked out, and its auto-fixes commit to `master` if that's where HEAD happens to be.
This is fundamentally incompatible with feature-branch workflows.

**Impact:** ~40% of session time was spent fighting BuildFlow (branch recovery, cherry-picking,
re-applying lost changes) instead of doing actual work.

### Placeholder Version Lie

I claimed "zero import path changes for consumers" — this was technically true (import paths
unchanged) but **practically false** (consumer builds fail because `v0.0.0-00010101000000-000000000000`
doesn't resolve from the Go proxy after replace directives are stripped). I should have tested
the consumer experience before claiming success.

---

## E) WHAT WE SHOULD IMPROVE 🔄

1. **BuildFlow needs feature-branch awareness** — it should never switch branches or commit to
   `master` when a feature branch is checked out. File an issue at `larsartmann/buildflow`.
2. **Test consumer experience during modularization** — not just `go build ./...`, but simulate
   a consumer: `go get github.com/larsartmann/templ-components/icons@latest` from a clean module.
3. **Add a `make dry-release` target** — run the release script in dry-run mode to verify
   placeholder→version replacement + tag creation without actually tagging.
4. **Reduce module count if co-change is too high** — git analysis showed 80%+ of commits touch
   5+ packages. If this pattern holds, the strategic split (6 modules) may still be too many.
   Monitor after 3 months of usage.
5. **Consider merging utils INTO svg** — svg tests depend on utils test helpers. With only 2
   functions in svg, it may not justify its own go.mod. Re-evaluate at v0.8.0.
6. **flake.nix should be multi-module aware** — the `verify` app should iterate over all modules,
   not just run `go build ./...` at root (which only covers the root module via go.work).

---

## F) Next 25 Things To Get Done

> **ALL TIER 1–2 ITEMS ARE MOOT** — modularization was abandoned. Tiers 3–4 updated below
> with current status where applicable to the single-module codebase.

### Tier 1 — Must Do Before Merge (blocks v0.7.0) — ❌ ALL MOOT

| #   | Task                                             | Status (2026-07-06)                      |
| --- | ------------------------------------------------ | ---------------------------------------- |
| 1   | **Merge `modularize/strategic-split` to master** | ❌ Abandoned                             |
| 2   | **Cut v0.7.0 release using per-module tags**     | ✅ v0.7.0 released as single-module      |
| 3   | **Verify consumer experience from clean repo**   | ✅ Standard single-module `go get` works |
| 4   | **Run CI on the branch**                         | ❌ Moot                                  |

### Tier 2 — Should Do After Merge — ❌ ALL MOOT (except #10)

| #   | Task                                                                  | Status (2026-07-06)                                                  |
| --- | --------------------------------------------------------------------- | -------------------------------------------------------------------- |
| 5   | **Update README.md** for multi-module                                 | ❌ Moot — single-module correct                                      |
| 6   | **Update CONTRIBUTING.md** for multi-module                           | ❌ Moot                                                              |
| 7   | **Update icons-only-adoption.md**                                     | ✅ Done (different framing)                                          |
| 8   | **Add `make dry-release`**                                            | ⬜ Not done                                                          |
| 9   | **Update flake.nix `verify` app** for multi-module                    | ❌ Moot                                                              |
| 10  | **File BuildFlow issue** — branch-switching + go.work gitignore fight | ✅ Resolved — BuildFlow no longer switches branches; go.work removed |
| 11  | **Clean up stashes**                                                  | ⬠ Likely resolved                                                    |

### Tier 3 — Improve Quality — ❌ MOSTLY MOOT

| #   | Task                                               | Status (2026-07-06)                                        |
| --- | -------------------------------------------------- | ---------------------------------------------------------- |
| 12  | **Add workspace sync CI check**                    | ❌ Moot                                                    |
| 13  | **Add replace-directive audit script**             | ❌ Moot                                                    |
| 14  | **Per-module go.sum audit**                        | ❌ Moot                                                    |
| 15  | **Add version-drift CI check**                     | ❌ Moot                                                    |
| 16  | **Document the release workflow** for multi-module | ❌ Moot                                                    |
| 17  | **Evaluate: should `htmx` be extracted?**          | ⬜ Deferred — may revisit post-v1.0                        |
| 18  | **Evaluate: merge svg + utils?**                   | ✅ Resolved — svg stays as `internal/svg` in single module |

### Tier 4 — Polish & Long-term

| #   | Task                                                                                                     | Status (2026-07-06)                               |
| --- | -------------------------------------------------------------------------------------------------------- | ------------------------------------------------- |
| 19  | **Add per-module coverage reporting**                                                                    | ❌ Moot                                           |
| 20  | **Update `docs/diagrams/internal-dependencies.d2`** — reflect new multi-module DAG                       | ❌ Moot — diagram is for single module            |
| 21  | **Consider go.work.sum in .gitignore vs committed**                                                      | ❌ Moot — no go.work                              |
| 22  | **Add `go work vendor` support**                                                                         | ❌ Moot                                           |
| 23  | **Stale doc cleanup** — ~30 docs/status/_.md and docs/planning/_.md files still reference `internal/svg` | ✅ Resolved — `internal/svg` is correct on master |
| 24  | **Monitor co-change after 3 months**                                                                     | ⬜ Ongoing                                        |
| 25  | **Consider independent versioning for icons**                                                            | ⬜ Deferred — may revisit post-v1.0               |

---

## G) Top #1 Question I Cannot Figure Out Myself

> ✅ **RESOLVED — DECISION: DO NOT MERGE.** The modularization branch was abandoned.
> Master remains single-module. The split was prototyped but the cost/benefit didn't justify
> it (80%+ of commits touch 5+ packages — co-change is too high for a 6-module split).
> Documentation corrected in `a0dbae7`. The split may be re-attempted post-v1.0 if the
> package graph warrants it.
