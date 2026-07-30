# Status Report — 2026-07-30 11:35 — BuildFlow Repeat Failure Triage (Round 2)

> **CRITICAL:** This is the **exact same set of 3 failures** from the 11:19 report, 16 minutes later.
> The BuildFlow daemon re-introduced all three bugs after the previous session's fixes.

## Context

BuildFlow's `govalid-generate` and `test-race` steps both failed with the **identical** three
failures triaged in `docs/status/2026-07-30_11-19_buildflow-failure-triage-three-fixes.md`.
This session re-fixed all three. The BuildFlow auto-commit daemon has already committed the
fixes (commit `cb6cf5e`).

**Timeline of the repeat regression:**

| Time   | Commit         | What happened                                                                                                                                                |
| ------ | -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| ~11:19 | (prev session) | Previous session fixed all 3 failures. Status report written.                                                                                                |
| 11:23  | `a7f63aa`      | "chore(ci): update golangci config" — **RE-INTRODUCED** all 3 disabled linters + `ireturn:` settings block (the exact regression documented in AGENTS.md T1) |
| 11:25  | `ffbd9e8`      | "regenerate breadcrumbs" — **RE-INTRODUCED** the `encoding/json/v2` drift in `breadcrumbs_templ.go`                                                          |
| 11:28  | `cb6cf5e`      | BuildFlow daemon committed this session's fixes with a generic hallucinated message                                                                          |
| 11:35  | (this report)  | Session verifying and documenting                                                                                                                            |

---

## a) FULLY DONE

### 1. `visualtest/doc.go:76` — `:=` Variable Shadow Bug (3rd occurrence)

**Failure:** `govalid-generate` — "declared and not used: sharedAllocCtx / allocCancel"

**Root cause:** Line 76 used `:=` inside `sync.Once.Do(func() { ... })`, shadowing the
package-level `sharedAllocCtx` and `allocCancel` with unused locals. The comment on
lines 72-75 **literally documents this exact trap** ("Use = (not :=)").

**Fix:** `:=` → `=` on line 76.

**This bug has now been fixed THREE times:** commit `1eb50fe`, the 11:19 session, and this
session. Each time it regresses because the visualtest doc gets rewritten by an agent or
daemon that doesn't preserve the `=` operator.

**Verified:** `go build ./...` and `go test -race ./...` in visualtest module both pass.

### 2. `.golangci.yml` — Disabled Linters Re-Enabled (7th occurrence)

**Failure:** `TestGolangciDisabledLinters` — godoclint, ireturn, testableexamples found
in the enable list + dead `ireturn:` settings block still present.

**Root cause:** Commit `a7f63aa` (the previous session's "fix" commit) re-added all three
disabled linters back to the enable list. This is the **exact T1 regression** documented in
AGENTS.md — it has now happened **at least 7 times** across sessions (AGENTS.md says 5;
this is the 6th and 7th).

**Fix:** Removed `godoclint`, `ireturn`, `testableexamples` from the enable list. Deleted
the dead `ireturn:` settings block (`allow: error/empty/anon/stdlib/generic`).

**Verified:** `TestGolangciDisabledLinters` PASS, `scripts/check-lint-config.sh` PASS,
`golangci-lint run ./...` = 0 issues.

### 3. `navigation/breadcrumbs_templ.go` — Import Drift (2nd occurrence)

**Failure:** `TestTemplGeneratedInSync` — source imports `encoding/json` (v1, intentional)
but generated file had `encoding/json/v2`.

**Root cause:** Commit `ffbd9e8` regenerated breadcrumbs with the wrong import. The source
`.templ` file deliberately uses `encoding/json` (v1) — only `errorpage` uses `json/v2`.
Running `templ generate` with the system binary or at the wrong time produced the drift.

**Fix:** Regenerated via `nix develop -c templ generate ./navigation/...` (v0.3.1020, matches
`go.mod`). Single-line diff: `encoding/json/v2` → `encoding/json`.

**Verified:** `TestTemplGeneratedInSync` PASS (all 81 subtests across 8 packages).

### Full Verification

```
go build ./...                          ✅ PASS
go test ./...                           ✅ ALL PASS (17 packages)
go test -race ./visualtest/...          ✅ PASS
golangci-lint run ./...                 ✅ 0 issues
scripts/check-lint-config.sh            ✅ PASS
TestGolangciDisabledLinters             ✅ PASS
TestTemplGeneratedInSync                ✅ PASS (81 subtests)
```

---

## b) PARTIALLY DONE

Nothing — all three fixes are complete and verified.

---

## c) NOT STARTED

- **AGENTS.md regression count not updated** — should bump from "5 times" to "7 times" for the
  `.golangci.yml` regression. The count is stale.
- **Previous status report not annotated** — `docs/status/2026-07-30_11-19_buildflow-failure-triage-three-fixes.md`
  should note that the same bugs recurred 16 minutes later.
- **Root cause of the repeat cycle not addressed** — the BuildFlow daemon is caught in a loop
  where it commits a stale working tree, re-introducing regressions. The 3-layer prevention
  (TestGolangciDisabledLinters + check-lint-config.sh + CI guard) catches it in CI/test-race
  but does NOT prevent the daemon from committing the broken state in the first place.

---

## d) TOTALLY FUCKED UP

**The BuildFlow daemon commit loop is a systemic problem.** Here's the cycle:

1. Session A fixes bugs → daemon commits fixes
2. Daemon or next agent touches files → stale `.golangci.yml` or wrong templ regen sneaks in
3. Daemon commits the broken state with a generic message
4. BuildFlow `test-race` / `govalid-generate` fails
5. Session B (this one) re-fixes the same bugs
6. **GOTO 1**

The 3 disabled-linters regression has happened **7 times**. The doc.go `:=` bug has happened
**3 times**. The breadcrumbs drift has happened **2 times**. This is not a code quality problem
— it's a **tooling problem**. The guard tests work (they catch it every time), but they fire
AFTER the daemon has already committed the broken state, creating an endless whack-a-mole loop.

**The daemon commit messages remain useless.** Commit `cb6cf5e` says "update lint config,
breadcrumbs component, and visualtest docs" — a hallucinated description that mentions nothing
about the actual fixes (shadow bug, disabled linter removal, import sync). `git log --grep` for
"shadow" or "disabled linter" finds nothing. This is the documented behavior from AGENTS.md T13.

---

## e) WHAT WE SHOULD IMPROVE

1. **Fix BuildFlow itself** (`larsartmann/buildflow`) — the daemon should generate commit messages
   from `git diff --stat`, not from a template. It should run `go test ./...` before committing.
   It should have a budget > 60s. This is the root cause of ALL THREE repeat regressions.

2. **Add a pre-commit guard for doc.go `:=`** — a test or grep check that fails if line 76 of
   `visualtest/doc.go` contains `:=` instead of `=`. The comment is not enough; agents keep
   ignoring it.

3. **Make the lint config self-healing** — instead of just testing that disabled linters are
   absent, have a pre-commit hook that STRIPS them automatically (like `check-lint-config.sh`
   but with `sed -i` to fix, not just report).

4. **Pin the templ import in breadcrumbs** — the source deliberately uses `encoding/json` (v1).
   Add a test that asserts `breadcrumbs.templ` imports `encoding/json` (not v2) to catch the
   drift at the SOURCE level, not just the generated level.

5. **Update AGENTS.md regression counts** — the "5 times" count is stale; should be "7 times".
   Future sessions need accurate counts to understand severity.

6. **Annotate old status reports** — the 11:19 report should note the bugs recurred within
   16 minutes, preventing false confidence that they're "fixed."

---

## f) Up to 50 Things We Should Get Done Next

### Critical (root cause fixes)

1. Fix BuildFlow daemon to generate commit messages from `git diff --stat` (repo: `larsartmann/buildflow`)
2. Fix BuildFlow daemon to run `go test ./...` before committing (currently 60s budget, no tests)
3. Fix BuildFlow daemon to increase budget beyond 60s for `go test ./...`
4. Add a `TestDocGoShadowGuard` test that asserts `visualtest/doc.go` line 76 uses `=` not `:=`
5. Update AGENTS.md to bump `.golangci.yml` regression count from 5 to 7
6. Annotate `docs/status/2026-07-30_11-19_*.md` with recurrence note
7. Make `scripts/check-lint-config.sh` auto-fix (strip disabled linters) instead of just reporting

### Lint config hardening

8. Add a pre-commit hook that STRIPS disabled linters from `.golangci.yml` automatically
9. Add a `.golangci.yml` golden file test — compare current config against a known-good canonical version
10. Consider moving the disabled-linter list to a separate `.golangci-disabled.txt` that the hook enforces
11. Add a CI step that fails if `.golangci.yml` changes without an updated regression count in AGENTS.md
12. Add `recvcheck` to the disabled list check (future-proof against new incompatible linters)

### Breadcrumbs / templ sync

13. Add a source-level test asserting `breadcrumbs.templ` imports `encoding/json` (not v2)
14. Add a CI step that runs `templ generate` and asserts zero diff (catches drift before merge)
15. Consider pinning all `.templ` imports in a test file (canonical import manifest)
16. Add a test that `encoding/json/v2` appears ONLY in `errorpage` package, nowhere else

### visualtest hardening

17. Add a test that `ShutdownBrowser()` actually calls `allocCancel()` (integration test)
18. Add a test that `newTab()` derives from a non-nil `sharedAllocCtx`
19. Consider replacing `sync.Once` + package vars with a lazy singleton struct (harder to shadow)
20. Add a linter rule (via `forbidigo`) that bans `:=` after a comment containing "Use ="

### BuildFlow daemon investigation

21. Audit all daemon commits from the last 7 days — count how many had generic messages
22. Check if the daemon is re-applying a stale stash or cached working tree
23. Check if the daemon runs `templ generate` with the system binary (v0.3.1036) instead of nix (v0.3.1020)
24. Consider disabling the daemon entirely until it's fixed (manual commits only)
25. Add a `BUILDFLOW_COMMIT_PREFIX` env var so daemon commits are identifiable in `git log`

### Documentation

26. Update AGENTS.md "BuildFlow gotcha" section with the 7th regression incident
27. Add a "Known Repeat Regressions" table to AGENTS.md with counts and last-occurrence dates
28. Update `docs/status/` README (if exists) with cross-references between related reports
29. Write an ADR for the BuildFlow daemon commit message problem and proposed fix
30. Update FEATURES.md if any feature status changed (unlikely this session)

### Code quality (unrelated to this session's bugs, but noticed)

31. Run `nix run .#verify` to confirm the full pipeline is green at HEAD
32. Run `nix run .#visual` to confirm visual regression tests pass
33. Check if any other `*_templ.go` files have import drift (run full `templ generate` and diff)
34. Run `go test -race ./...` on the FULL suite (not just visualtest) to catch race conditions
35. Audit `visualtest/go.mod` — the daemon touched it (commit `ffbd9e8`); verify no unwanted deps

### Test coverage gaps

36. Add a fuzz test for `ensureAllocator` — concurrent calls should all get the same allocator
37. Add a test for the case where `CHROMEDP_CHROME_PATH` points to a non-executable file
38. Add a golden test for the breadcrumbs JSON-LD output (currently only string-tested)
39. Add a test that `breadcrumbJSONLD` handles empty `items` slice gracefully
40. Add a test that `resolveBreadcrumbURL` handles protocol-relative URLs (`//host/path`)

### Process improvements

41. Add a `make verify` / `nix run .#verify` step to the daemon's pre-commit (if possible)
42. Consider a Git hook that rejects commits with "Unknown Author" or generic messages
43. Add a `CONTRIBUTING.md` section on the daemon commit loop and how to break it
44. Tag the BuildFlow issue in `larsartmann/buildflow` with all 7 regression incidents
45. Consider switching from daemon-based auto-commit to pre-push hooks only

### Nice-to-have

46. Add `git config commit.template` with a structured message format for this repo
47. Add a `.git-blame-ignore-revs` file to hide pure-formatting daemon commits from blame
48. Add a CI badge for `test-race` specifically (separate from general test status)
49. Consider a `renovate.json` or Dependabot config for templ version pin management
50. Write a postmortem ADR linking all 7 regression incidents and the proposed BuildFlow fix

---

## g) Questions I Cannot Answer Myself

1. **Should we disable the BuildFlow daemon entirely until it's fixed?** The daemon has
   re-introduced the same 3 bugs at least twice today. Disabling it means manual commits,
   but eliminates the regression loop. I don't know your preference on manual vs. auto-commit
   or whether the daemon provides value I'm not seeing (e.g., auto-committing during long
   sessions to prevent data loss).

2. **Is commit `a7f63aa` (the previous session's "fix") the regression source, or was it
   already broken before that session?** The git log shows `a7f63aa` ADDED the disabled
   linters back (10 lines added to `.golangci.yml`), but I can't tell if the previous session
   intentionally modified the lint config for a valid reason that backfired, or if the daemon
   just committed a stale file. Understanding this determines whether the fix is "don't touch
   `.golangci.yml`" or "fix the daemon's file handling."

3. **Should the `visualtest/doc.go` allocator pattern be rewritten entirely?** The `sync.Once`
   - package-level vars + `=` (not `:=`) pattern is fragile — it's been broken 3 times. An
     alternative is a lazy-init struct (e.g., `type browserAllocator struct { once sync.Once;
ctx context.Context; cancel context.CancelFunc; init func() (context.Context, context.CancelFunc) }`)
     which makes shadowing impossible. But this changes the package's public API (`ShutdownBrowser()`
     signature) and I don't know if external consumers depend on the current API shape.
