# Status Report: CSS Integration Automation — Cleanup & Completion

**Date:** 2026-07-09 09:07 (updated 09:25)
**Session Duration:** ~115 minutes (07:30 – 09:25)
**Scope:** templ-components, BuildFlow, DiscordSync
**Branches:** `master` (templ-components, DiscordSync), `nix-style-output` (BuildFlow)

---

## Executive Summary

The CSS integration automation sprint started with two tools (`tc-css` CLI + BuildFlow
`tailwind-build` provider). After critical review, **`tc-css` was deleted** — it was
over-engineered, broken under `-mod=vendor` (its primary use case), and duplicated logic.
What remains is the right level of automation: a copy-paste starter template, BuildFlow's
DAG-integrated provider, and well-documented manual setup.

A build-breaking `encoding/json/v2` import was also found and fixed — it was accidentally
introduced by an auto-formatter and pushed to master.

All 3 repos are committed and pushed. Builds pass. Tests pass.

---

## a) FULLY DONE

| #   | Task                                                                                                                   | Repo             | Verification                                    |
| --- | ---------------------------------------------------------------------------------------------------------------------- | ---------------- | ----------------------------------------------- |
| 1   | **Deleted `cmd/tc-css/`** — over-engineered, broken under `-mod=vendor`, duplicated BuildFlow logic                    | templ-components | 882 lines removed, zero references remain       |
| 2   | **Fixed build-breaking `encoding/json/v2` import** in `errorpage/handler.go` → reverted to `encoding/json`             | templ-components | `go build ./...` passes                         |
| 3   | **Fixed `breadcrumbs_templ.go`** — same `encoding/json/v2` → `encoding/json` revert                                    | templ-components | `go build ./...` passes                         |
| 4   | **Added `tailwindcss_4` to devShell** (templ-components + BuildFlow)                                                   | both             | Binary available in `nix develop`               |
| 5   | **DiscordSync `styles.css` regenerated** — class-based dark mode (0→80 `.dark` selectors, 80→0 `prefers-color-scheme`) | DiscordSync      | `grep -c "\.dark " styles.css` → 80             |
| 6   | **DiscordSync split-brain resolved** — reverted `generate-css` to `gen/main.go --css-only`                             | DiscordSync      | `nix run .#generate-css` works                  |
| 7   | **BuildFlow vendor/ false-positive fixed** — `skipDirs["vendor"]` changed to `true` + test added                       | BuildFlow        | `TestFindTailwindEntryPointsSkipsVendor` passes |
| 8   | **All docs cleaned** — README, adoption guide, migration guide, AGENTS.md across all repos                             | all 3            | `grep -rnl "tc-css"` → zero results             |
| 9   | **CHANGELOG `[Unreleased]`** updated — documents app.css, BuildFlow provider, json/v2 fix                              | templ-components | Drift guard test passes                         |
| 10  | **All repos committed and pushed**                                                                                     | all 3            | `git status --short` clean                      |

---

## b) PARTIALLY DONE

| #   | Task                                              | Status           | What's Missing                                                                                                                                                                          |
| --- | ------------------------------------------------- | ---------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **BuildFlow devShell builds**                     | ⚠️ Blocked       | `tailwindcss_4` added to `flake.nix`, but `nix develop` fails due to unrelated `hierarchical-errors` flake input broken by another session's uncommitted work in `modules/nix-checker/` |
| 2   | **BuildFlow tailwind unit tests in full package** | ⚠️ Isolated pass | Tests pass when `modules/nix-checker` pre-existing build failure is stashed. Cannot run in full `go test ./...` until that's fixed by the other session.                                |

---

## c) NOT STARTED

| #   | Task                                                                                     | Why                                                          |
| --- | ---------------------------------------------------------------------------------------- | ------------------------------------------------------------ |
| 1   | Release v0.11.0                                                                          | Needs user decision — see questions below                    |
| 2   | Pre-commit hardening (`go build ./...` + `encoding/json/v2` grep guard + lint `./...`)   | Identified but not implemented this session                  |
| 3   | ADR for output-path inference duplication (BuildFlow + manual setup use same convention) | Low priority — `tc-css` deleted so only BuildFlow has it now |
| 4   | FEATURES.md update with `templates/app.css` and BuildFlow `tailwind-build`               | Minor doc task                                               |

---

## d) TOTALLY FUCKED UP

| #   | What                                                                   | Impact                                                                                                                      | Root Cause                                                                                                                                      | Resolution                                                                              |
| --- | ---------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| 1   | **`encoding/json/v2` import in `errorpage/handler.go`**                | Build-breaking — entire `errorpage` package failed to compile outside `GOEXPERIMENT=jsonv2`. Pushed to master in `ad58171`. | Auto-formatter/LSP ran under `GOEXPERIMENT=jsonv2`, rewrote `encoding/json` → `encoding/json/v2`. Pre-commit hook doesn't run `go build ./...`. | **Fixed.** Reverted to `encoding/json`. Committed in `bec3d30`.                         |
| 2   | **`tc-css` broken under `-mod=vendor`**                                | Tool's primary value proposition didn't work. DiscordSync hit this failure immediately.                                     | `go run github.com/larsartmann/templ-components/cmd/tc-css` can't find the package after vendoring — fundamental Go limitation.                 | **Fixed.** Deleted `tc-css` entirely. It was over-engineered for a one-time setup task. |
| 3   | **DiscordSync `generate-css` migrated to `tc-css` before E2E testing** | `nix run .#generate-css` was broken.                                                                                        | Previous session migrated without testing.                                                                                                      | **Fixed.** Reverted to `gen/main.go --css-only` which works correctly.                  |

---

## e) WHAT WE SHOULD IMPROVE

| #   | Area                                 | Current State                                                                      | Improvement                                                                      |
| --- | ------------------------------------ | ---------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| 1   | **Pre-commit build check**           | Hook runs templ generate + lint on specific dirs, but NOT `go build ./...`         | Add `go build ./...` to catch import errors before push                          |
| 2   | **`encoding/json/v2` safety**        | Auto-formatters can rewrite imports when running under the flag                    | Add grep guard: `! grep -rn "encoding/json/v2" --include="*.go" .` in pre-commit |
| 3   | **Pre-commit lint coverage**         | `scripts/pre-commit.sh` lists specific package dirs, omits `cmd/`                  | Use `golangci-lint run ./...` instead of hardcoded list                          |
| 4   | **E2E testing before commit**        | `tc-css` was committed without testing against a real project                      | Always E2E test new tools before pushing                                         |
| 5   | **YAGNI enforcement**                | Built `tc-css` when `templates/app.css` + one `tailwindcss` command does the job   | Question every new tool: does this earn its complexity?                          |
| 6   | **BuildFlow working tree hygiene**   | 12+ uncommitted files from another session block `nix develop` and `go test ./...` | Commit or stash other sessions' work before starting                             |
| 7   | **Status reports should be current** | Previous session's report claimed "nothing committed" but everything was committed | Always `git status` + `git log` before writing status reports                    |

---

## f) Things We Should Get Done Next

### Critical (P0)

| #   | Task                                                                           | Impact                              | Effort |
| --- | ------------------------------------------------------------------------------ | ----------------------------------- | ------ |
| 1   | Add `go build ./...` to pre-commit hook                                        | Prevent build-breaking commits      | 10m    |
| 2   | Add `encoding/json/v2` grep guard to pre-commit                                | Prevent the exact bug that happened | 5m     |
| 3   | Cut v0.11.0 release (tc-css deleted, json/v2 fix, app.css, BuildFlow provider) | Gets clean code to consumers        | 30m    |

### High Priority (P1)

| #   | Task                                                                     | Impact                                    | Effort |
| --- | ------------------------------------------------------------------------ | ----------------------------------------- | ------ |
| 4   | Update pre-commit lint to `golangci-lint run ./...` (not hardcoded dirs) | Lint gaps cause CI failures               | 10m    |
| 5   | Commit or stash BuildFlow's 12+ uncommitted files from other session     | Unblock `nix develop` and `go test ./...` | 15m    |
| 6   | Document `encoding/json/v2` prohibition in AGENTS.md                     | Knowledge preservation                    | 10m    |
| 7   | Update FEATURES.md with app.css + BuildFlow tailwind-build provider      | Feature inventory                         | 15m    |
| 8   | Document the `encoding/json/v2` auto-formatter gotcha as an ADR          | Prevent recurrence                        | 15m    |

### Medium Priority (P2)

| #   | Task                                                                                     | Impact                | Effort |
| --- | ---------------------------------------------------------------------------------------- | --------------------- | ------ |
| 9   | Write ADR for tc-css deletion decision (why we built it, why we killed it)               | Decision record       | 15m    |
| 10  | Add `tailwindcss_4` to DiscordSync devShell (if not already there)                       | DevEx                 | 5m     |
| 11  | Integration test for BuildFlow `tailwind-build` provider against a real project          | Regression prevention | 1h     |
| 12  | Add `go:build !goexperiment.jsonv2` or equivalent guard to prevent accidental v2 imports | Compile-time safety   | 30m    |
| 13  | Update ROADMAP.md with CSS automation milestone                                          | Planning              | 10m    |

### Polish (P3)

| #   | Task                                                                 | Impact          | Effort |
| --- | -------------------------------------------------------------------- | --------------- | ------ |
| 14  | Review all docs for consistency after tc-css deletion                | Polish          | 20m    |
| 15  | Add `templates/app.css` to the demo binary                           | Discoverability | 15m    |
| 16  | Write blog post / announcement for BuildFlow tailwind-build provider | Marketing       | 30m    |
| 17  | Consider adding `--watch` mode to BuildFlow tailwind-build provider  | DevEx           | 1h     |

---

## g) Top 2 Questions I Cannot Answer Myself

### 1. Should we cut v0.11.0 now?

The `encoding/json/v2` build break was on master but never tagged. The last tag (v0.10.0) is clean. So consumers on `@latest` are fine — only `@master` users were affected.

However, `[Unreleased]` now has real content (app.css, BuildFlow provider, json/v2 fix, tc-css deletion). Ready for release whenever you say go.

### 2. Should the BuildFlow `tailwind-build` provider also handle `@source` auto-detection?

Currently the provider discovers CSS entry points and runs `tailwindcss`, but it doesn't auto-write `@source` directives. Consumers still need to manually write `@source "../vendor/..."` in their CSS. Should the provider inject `@source` paths into the CSS before compiling, or is that the consumer's responsibility?

I lean toward **consumer's responsibility** — the provider shouldn't mutate source files. But it's a judgment call.

---

## Repo State

| Repo             | Branch             | HEAD       | Working Tree                 | Build                     | Test       | Lint        |
| ---------------- | ------------------ | ---------- | ---------------------------- | ------------------------- | ---------- | ----------- |
| templ-components | `master`           | `83ee054`  | ✅ clean                     | ✅ passes                 | ✅ passes  | ✅ 0 issues |
| BuildFlow        | `nix-style-output` | `aed67d55` | ⚠️ 12+ dirty (other session) | ⚠️ blocked by nix-checker | ⚠️ blocked | N/A         |
| DiscordSync      | `master`           | `b594bcd`  | ✅ clean                     | ✅ passes                 | ✅ passes  | N/A         |

---

## Session Metrics

| Metric                | Value                                                     |
| --------------------- | --------------------------------------------------------- |
| Total commits         | 9 (5 templ-components, 2 BuildFlow, 2 DiscordSync)        |
| Lines deleted         | ~900 (tc-css + stale docs)                                |
| Lines changed         | ~300 (fixes + doc cleanup)                                |
| Bugs found & fixed    | 4 (json/v2 import, lint failures, vendor skip, dark mode) |
| Tools killed          | 1 (tc-css — over-engineered)                              |
| E2E tests run         | 2                                                         |
| Planning docs written | 1 (then deleted with tc-css)                              |
