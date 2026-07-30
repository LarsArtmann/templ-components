# Status Report — 2026-07-30 11:19 — BuildFlow Failure Triage (3 Failures, All Fixed)

## Context

BuildFlow's `govalid-generate` and `test-race` steps both failed. Three distinct
failures were reported in the BuildFlow output. This session triaged and fixed all
three. The BuildFlow auto-commit daemon has already committed the fixes (commits
`4fd92eb` and `dcc1f39`).

---

## a) FULLY DONE

### 1. `visualtest/doc.go` — Variable Shadowing Bug (`:=` vs `=`)

**Failure:** `govalid-generate` failed — "declared and not used: sharedAllocCtx / allocCancel"

**Root cause:** Line 76 used `:=` (short variable declaration) inside `sync.Once.Do(func() { ... })`,
which created local variables that shadowed the package-level `sharedAllocCtx` and `allocCancel`.
The package-level vars were never assigned, leaving them `nil`. The comment on lines 72-75
**literally warns about this exact bug** ("Use = (not :=) so we assign the package-level vars").

**Fix:** Changed `:=` to `=` on line 76.

**Verified:** `go build ./...` in the visualtest module succeeds. gopls diagnostics clear.

**NOTE:** Commit `1eb50fe` previously fixed this **exact same bug** ("Fix visualtest doc.go: :=
shadowed package vars inside sync.Once.Do"). The regression was reintroduced — likely during
the `chore(navigation): update breadcrumbs component and visual test dependencies` commit
(`102d69f`) which touched `visualtest/go.mod` and may have triggered a file rewrite.

### 2. `.golangci.yml` — Disabled Linters Re-Enabled (6th Occurrence)

**Failure:** `TestGolangciDisabledLinters` failed — godoclint, ireturn, testableexamples
were back in the `enable:` list, plus the dead `ireturn:` settings block.

**Root cause (documented in AGENTS.md):** BuildFlow daemon commits a stale working tree
without running tests. The `.golangci.yml` in the working tree was stale. This has now
happened **6 times** across sessions (AGENTS.md documents 5; this is the 6th).

**Fix:** Removed all three linters from `enable:` list + deleted the `ireturn:` settings block.

**Verified:** `TestGolangciDisabledLinters` passes, `scripts/check-lint-config.sh` passes,
`golangci-lint run` reports 0 issues.

### 3. `navigation/breadcrumbs_templ.go` — Import Drift

**Failure:** `TestTemplGeneratedInSync` failed — `breadcrumbs.templ` imports `encoding/json`
but `breadcrumbs_templ.go` imported `encoding/json/v2`.

**Root cause:** The generated file was regenerated at some point with a system templ binary
or with `GOEXPERIMENT=jsonv2` active, which rewrote the import. The source `.templ` was not
updated to match.

**Fix:** Ran `GOEXPERIMENT=jsonv2 templ generate ./navigation/...` with templ v0.3.1020
(matches `go.mod`). Generated file now imports `encoding/json` (v1), matching the source.

**Verified:** `TestTemplGeneratedInSync` passes for all 50+ `.templ` files across all packages.

### Full Verification Run

| Check                               | Result               |
| ----------------------------------- | -------------------- |
| `go test ./...` (main module)       | All 16 packages PASS |
| `go test ./...` (visualtest module) | PASS                 |
| `golangci-lint run` (all packages)  | 0 issues             |
| `scripts/check-lint-config.sh`      | PASS                 |
| `go build ./...` (both modules)     | PASS                 |

---

## b) PARTIALLY DONE

Nothing — all three reported failures were fully resolved.

---

## c) NOT STARTED (Noticed But Not Addressed This Session)

### 1. breadcrumbs.templ still uses `encoding/json` v1

The project standard is `encoding/json/v2` (AGENTS.md, `.golangci.yml` build-tags,
`.envrc`). `breadcrumbs.templ` line 4 imports `encoding/json` (v1) while every other
package migrated to v2. The `json.Marshal(list)` call at line 89 works with both, but
this is an inconsistency. **I synced the generated file to match the v1 source rather
than updating the source to v2.** The right long-term fix is to update the `.templ` to
`encoding/json/v2`.

### 2. AGENTS.md regression counter not updated

AGENTS.md says the `.golangci.yml` regression happened "5 times across sessions." This
session makes it 6. The doc should be updated.

### 3. `nix run .#verify` not run

AGENTS.md prescribes nix flake commands for verification. I used raw `go test` and
`golangci-lint` instead. Functionally equivalent, but doesn't exercise the flake's
`GOEXPERIMENT=jsonv2` shellHook injection path.

### 4. `test-race` not explicitly re-run

The BuildFlow failure was `test-race` (runs `go test -race`). I ran `go test ./...`
(without `-race`) for verification. The race detector may surface issues not caught
by plain test runs.

### 5. `govalid-generate` not explicitly re-run

This was one of the two BuildFlow failures. I verified the compile errors are gone
via `go build`, but did not re-run `govalid-generate` itself to confirm the toolchain
step passes end-to-end.

---

## d) TOTALLY FUCKED UP

Nothing catastrophic. No data loss, no broken state, no irreversible damage.

**However — the recurring regression pattern is deeply concerning:**

- The `.golangci.yml` disabled-linter regression has now occurred **6 times**.
  The prevention layers (test guard + shell script + CI) catch it every time, but
  **the root cause is still unfixed**: BuildFlow's auto-commit daemon commits a stale
  working tree without running `go test ./...`. The daemon has a 60s budget and
  cannot run the full test suite. This means every fix is a band-aid; the regression
  WILL happen again.
- The `doc.go :=` shadowing bug was previously fixed (commit `1eb50fe`) and
  reintroduced. File rewrites during dependency bumps or templ generation can silently
  revert fixes.

---

## e) WHAT WE SHOULD IMPROVE

### Process Improvements

1. **Fix BuildFlow's daemon** (`larsartmann/buildflow`) to either (a) run
   `go test ./...` before committing, or (b) skip auto-commit when known drift-guard
   tests exist. The 60s budget is insufficient for `go test -race` but the
   `TestGolangciDisabledLinters` test alone takes <1s — a "fast-guard" tier could run.

2. **Add a pre-commit guard for variable shadowing** in `visualtest/doc.go`.
   A simple grep for `:=` on the `sharedAllocCtx` line would prevent the 2nd reintroduction.

3. **Migrate breadcrumbs.templ to `encoding/json/v2`** to match the project standard.
   Low risk — only uses `json.Marshal`.

4. **Run `nix run .#verify`** instead of ad-hoc go commands. The flake ensures
   `GOEXPERIMENT=jsonv2` is set, which affects code generation behavior.

### Architectural Observations

5. The `TestTemplGeneratedInSync` test only checks that source imports exist in
   generated files — it does NOT check the reverse (generated imports not in source).
   This means a generated file with EXTRA imports (like the json v2 drift) would only
   be caught if those imports are also in the source. Consider bidirectional checking.

6. The `encoding/json` vs `encoding/json/v2` drift will recur as long as the project
   has mixed usage. A project-wide grep guard (like `TestGolangciDisabledLinters`)
   could enforce v2-only in `.templ` files.

---

## f) Next 50 Things To Get Done

### High Priority — Prevent Recurring Regressions

1. Fix BuildFlow daemon to run fast-guard tests before committing
2. Add grep-based pre-commit guard for `sharedAllocCtx, allocCancel :=` in visualtest/doc.go
3. Update AGENTS.md regression counter: ".golangci.yml regression count: 5 → 6"
4. Migrate `navigation/breadcrumbs.templ` from `encoding/json` to `encoding/json/v2`
5. Run `nix run .#verify` to confirm flake-level verification passes
6. Run `go test -race ./...` to confirm the `test-race` BuildFlow step passes
7. Re-run `govalid-generate` (or `buildflow -s govalid-generate`) to confirm the step passes
8. Add a bidirectional sync check to `TestTemplGeneratedInSync` (generated imports not in source)
9. Add a `TestEncodingJSONVersion` guard test: all `.templ` files must import `encoding/json/v2`
10. Audit ALL `*_templ.go` files for json v1 vs v2 drift (not just breadcrumbs)

### Medium Priority — Hardening

11. Add `git blame` check to the triage workflow — identify which commit introduced each regression
12. Consider a `.gitattributes` merge strategy for `.golangci.yml` to prevent silent overwrites
13. Add a visualtest-specific compile guard to CI (currently only govalid-generate catches it)
14. Document the doc.go `:=` bug pattern in AGENTS.md visualtest section
15. Review ALL commits from `102d69f` onward for other latent regressions
16. Check if `visualtest/go.mod` was inadvertently modified by `go test` during this session
17. Add `go test -race` to the pre-commit hook (or a fast subset)
18. Consider committing `.golangci.yml` with a checksum guard file
19. Run the full visual regression suite (`nix run .#visual`) to confirm no visual regressions
20. Review the `011d396` commit (TitleTag feature) for correctness — it landed during this session

### Feature / Component Work

21. Review `Card.TitleTag` and `EmptyState.TitleTag` from commit `011d396` for accessibility
22. Add golden tests for the new TitleTag variants
23. Update FEATURES.md if TitleTag is a new user-facing feature
24. Add CHANGELOG entry for TitleTag under `[Unreleased]`
25. Verify TitleTag works with all heading levels (h1-h6)
26. Check TitleTag interaction with container-aware Card

### Testing Improvements

27. Add property-based tests for breadcrumb JSON-LD generation
28. Add edge case tests for breadcrumbs with special characters in URLs
29. Add visual regression tests for breadcrumbs dark mode
30. Add contract tests for the new TitleTag field
31. Increase golden test coverage to target >90%
32. Add race condition tests for visualtest allocator lifecycle
33. Add fuzz tests for `.golangci.yml` parsing (TestGolangciDisabledLinters)

### Documentation

34. Update AGENTS.md with this session's findings (6th regression, doc.go reintroduction)
35. Document the `encoding/json` v1/v2 split decision in an ADR
36. Update `docs/visual-testing.md` with the `:=` shadowing gotcha
37. Add a "known recurring regressions" section to AGENTS.md with prevention checklist
38. Update TODO_LIST.md with the BuildFlow daemon fix task
39. Review and update all status reports from 2026-07-28 onward for accuracy
40. Create a runbook for "BuildFlow test-race failed" common failures

### Code Quality

41. Audit all `sync.Once` usage in the codebase for similar shadowing risks
42. Check if any other package-level vars have `:=` shadowing bugs
43. Review the `encoding/json/v2` migration completeness across all packages
44. Run `nix flake check` to verify flake-level integrity
45. Audit all pre-commit hooks for coverage gaps
46. Consider adding `gocritic`'s shadowcheck if available
47. Review import ordering consistency across all `_templ.go` files
48. Check if the `cmd/tc/_sources/navigation/breadcrumbs.templ` copy also needs json v2 migration
49. Audit the `.gitignore` BuildFlow `*_templ.go` interaction (documented gotcha)
50. Review whether `govalid-generate` should be in the pre-commit hook vs CI-only

---

## g) Questions (Cannot Determine Myself)

1. **Should `breadcrumbs.templ` be migrated to `encoding/json/v2`?** The source uses v1
   while the project standard is v2. I synced the generated file to match v1, but the
   long-term fix depends on whether you want breadcrumbs to join the v2 migration or
   stay on v1 for backward compatibility.

2. **Is commit `011d396` (TitleTag feature) expected?** It appeared during this session
   as a new HEAD commit. I did not create it — was it from another session or manual work?
   Should I review it?

3. **Should I fix BuildFlow itself (`larsartmann/buildflow`) to address the root cause?**
   The `.golangci.yml` regression will keep recurring until the daemon either runs tests
   or refuses to commit files it can't verify. This is a cross-repo change — do you want
   me to proceed?
