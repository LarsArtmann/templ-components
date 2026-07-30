# Status Report — 2026-07-30 22:19

## Session Scope

Fixed two BuildFlow failures (visualtest compile error + `.golangci.yml` linter
regression), committed, discovered the BuildFlow daemon re-introduces both
issues, traced the root causes, filed feedback upstream, and partially applied
a permanent fix for one of them. This report covers only this session's work.

---

## a) FULLY DONE

### 1. Diagnosed and fixed `visualtest/doc.go` compile error (session start)

**Problem:** `:=` (short declaration) on line 76 shadowed the package-level
`sharedAllocCtx`/`allocCancel` vars, creating unused locals → compile error
blocking the entire `visualtest` module.

**Fix applied:** Changed `:=` to `=`.

**Root cause discovered later:** Not a human typo — the fatcontext linter's
autofix (run by BuildFlow's `golangci-lint --fix` repair step) unconditionally
converts `=` to `:=`. Filed as fatcontext issue #100 (closed as dup of #43).
The maintainer suggested splitting the assignment to make the autofix a no-op.

### 2. Diagnosed and fixed `.golangci.yml` linter regression (session start)

**Problem:** Three documented-disabled linters (`godoclint`, `ireturn`,
`testableexamples`) had re-entered the `enable` list + dead `ireturn:` settings
block was present. `TestGolangciDisabledLinters` caught all 4 violations.

**Fix applied:** Removed all three from `enable` + deleted `ireturn:` settings.

### 3. Identified fatcontext #100 as root cause of recurring `:=` bug

**What:** Checked `git log` and found the `:=` regression has been fixed
**8+ times** across sessions. Searched fatcontext issues, found #100 (filed by
LarsArtmann, closed as dup of #43). The maintainer confirmed: fatcontext's
autofix unconditionally converts `=` to `:=` without checking for outer-scope
shadowing. Suggested fix: split the assignment so `:=` is already present.

### 4. Applied permanent fix for `:=` regression

**What:** Rewrote the allocator assignment to declare locals with `:=` then
assign to package-level vars separately:

```go
allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
sharedAllocCtx = allocCtx
allocCancel = cancel
```

This makes fatcontext's autofix a no-op (`:=` already present). The locals are
immediately consumed so no "declared and not used". Build verified.

### 5. Committed the fixes

**Commit:** `d2efac6` — "docs: add status report for visualtest compile +
golangci linter regression fix". The code fixes were already committed by the
BuildFlow daemon in `1b10484`; this commit captured the status report + other
accumulated doc changes. BuildFlow pre-commit passed (44/0/0).

### 6. Diagnosed `.golangci.yml` regression root cause (7th+ occurrence)

**What:** After committing, BuildFlow's `test-race` step failed again — the
same three linters were back in `enable`. Re-removed them. This is the 7th+
documented occurrence of this exact regression.

### 7. Filed upstream feedback to golangci-lint-auto-configure

**What:** Wrote comprehensive feedback at
`/home/lars/projects/golangci-lint-auto-configure/docs/feedback/new/2026-07-30_repair-re-adds-linters-removed-from-enable.md`.

Explained that the 2026-07-25 fix (respect `linters.disable`) doesn't cover the
case where a linter is simply absent from both `enable` and `disable`. Proposed
5 solutions. Identified that `godoclint`/`testableexamples` are in no tier at
all and `ireturn` is in `PragmaticNoiseLinters` (enabled by default).

### 8. Identified immediate workaround

The project can add the three linters to `linters.disable` (not just remove
from `enable`) to signal intent that `repair` respects per the 2026-07-25 fix.

---

## b) PARTIALLY DONE

### `.golangci.yml` permanent fix — NOT YET APPLIED

The workaround (add to `disable`) was identified and offered to the user but
**not yet applied**. The working tree currently has the symptom fix (removed
from `enable`) but not the root-cause workaround (added to `disable`). Without
the workaround, the next BuildFlow `repair` run will re-add them again.

### `visualtest/doc.go` permanent fix — APPLIED BUT NOT COMMITTED

The split-assignment fix is in the working tree, verified to build, but not
committed. The committed version (`d2efac6`) has the old `=` fix which
BuildFlow will revert to `:=` on the next commit that triggers `templ-generate`.

---

## c) NOT STARTED

1. **Applying the `disable:` workaround** to `.golangci.yml` — the user was
   asked "Want me to apply that?" but the status report request interrupted.
2. **Committing the split-assignment fix** for `visualtest/doc.go`.
3. **Full verification pipeline** — `go test ./...`, `golangci-lint run`,
   `nix run .#verify` were never run this session. Only targeted checks.
4. **Updating AGENTS.md** — regression counts ("5+ sessions", "8th time") and
   the fatcontext #100 root cause are not documented there yet.
5. **Checking if fatcontext #43 has a fix pending** that would make the
   split-assignment workaround unnecessary.

---

## d) TOTALLY FUCKED UP

### Committed without verifying the working tree was clean

After commit `d2efac6`, BuildFlow's pre-commit repair step **immediately
re-introduced both regressions**:

- `visualtest/doc.go`: `=` reverted to `:=` (fatcontext autofix)
- `.golangci.yml`: three linters re-added to `enable` (golangci-lint-auto-configure repair)

I discovered this only by checking `git diff` after the commit. The commit
itself has the correct versions (pre-hook state), but the working tree was
left dirty with broken code. A developer pulling or continuing work would hit
compile errors.

**Lesson:** After any commit that triggers BuildFlow hooks, always run
`git diff` to verify the hook didn't re-introduce regressions.

---

## What I Forgot / Could Have Done Better

### 1. Did not apply the `disable:` workaround before committing

I identified the workaround (add linters to `disable` so `repair` respects
them) but only offered it as a question after the commit. The commit went out
with the fragile symptom fix (remove from `enable` only), guaranteeing the
next `repair` run re-adds them. I should have applied the workaround first.

### 2. Did not commit the split-assignment fix for visualtest/doc.go

I applied the permanent fix (split assignment) and verified it builds, but
never committed it. It's sitting in the working tree. If another commit
happens (e.g., the daemon), the split-assignment fix could be lost or
interleaved with unrelated changes.

### 3. Did not run the full test suite after the `.golangci.yml` fix

I ran `TestGolangciDisabledLinters` in isolation but never ran
`go test ./...`. If the `.golangci.yml` changes broke anything else (unlikely
but possible), I wouldn't know.

### 4. Did not update AGENTS.md with the fatcontext #100 root cause

The `:=` regression has been documented as a "BuildFlow bug" for sessions.
Now I know the exact root cause (fatcontext autofix, issue #100/#43). This
should be in AGENTS.md so future sessions don't waste time re-diagnosing.

### 5. Filed feedback but didn't check if the fix already exists

fatcontext #43 was filed before #100. I didn't check whether #43 has been
fixed in a recent fatcontext release that we could upgrade to.

### 6. Did not verify the split-assignment fix survives BuildFlow

I verified it builds, but I didn't run a BuildFlow `repair` cycle to confirm
the autofix is truly a no-op on the split assignment. This is the whole point
of the fix — it should be tested against the actual tool.

---

## e) WHAT WE SHOULD IMPROVE

### The two regression loops are systemic

Both issues share the same pattern: a tool's autofix "corrects" code that is
already correct, breaking it, and the developer must manually undo the fix
on every commit.

| Regression       | Tool                                | Root Cause                                              | Permanent Fix                              |
| ---------------- | ----------------------------------- | ------------------------------------------------------- | ------------------------------------------ |
| `:=` shadowing   | fatcontext autofix                  | Converts `=` to `:=` unconditionally (#100/#43)         | Split assignment (applied, not committed)  |
| Linter re-enable | golangci-lint-auto-configure repair | Re-adds linters absent from both `enable` and `disable` | Add to `disable` (identified, not applied) |

**Both fixes are workarounds, not root-cause fixes.** The real fixes belong
upstream (fatcontext and golangci-lint-auto-configure). The feedback filed is
correct, but until those ship, the project needs defensive workarounds +
documentation.

### BuildFlow's repair pipeline needs a "verify after repair" step

The pre-commit hook runs repair tools that modify the working tree, but never
verifies the result compiles or passes tests. A post-repair `go build ./...`
gate would catch autofix-induced regressions before they reach the commit.

---

## f) Next Steps (Up to 50)

### Immediate (blocking — do now)

1. Apply the `disable:` workaround: add `godoclint`, `ireturn`,
   `testableexamples` to `linters.disable` in `.golangci.yml`
2. Commit the split-assignment fix for `visualtest/doc.go`
3. Run `go test ./...` to verify the full suite passes
4. Run `git diff` after commit to verify BuildFlow didn't revert anything
5. Run `nix run .#verify` for the full pipeline

### Documentation (high value)

6. Update AGENTS.md `:=` regression section: root cause is fatcontext #100/#43,
   not generic "BuildFlow bug"
7. Document the split-assignment workaround pattern in AGENTS.md
8. Document the `disable:` workaround for linter regressions in AGENTS.md
9. Bump regression count in AGENTS.md (7+ for linter, 8+ for `:=`)
10. Update the `visualtest/doc.go` comment to reference fatcontext #100

### Verification (should have done this session)

11. Run `golangci-lint run` to confirm 0 findings after `.golangci.yml` fix
12. Run `scripts/check-lint-config.sh` after `.golangci.yml` edit
13. Verify the split-assignment survives a BuildFlow `repair` cycle
14. Run `govalid-generate` BuildFlow step to confirm it succeeds
15. Run `nix flake check` for format verification

### Upstream (root-cause fixes)

16. Check if fatcontext #43 has been fixed in a recent release
17. Check if upgrading fatcontext would resolve the `:=` autofix
18. Monitor golangci-lint-auto-configure feedback for response
19. Consider contributing a PR to fatcontext (check `pass.TypesInfo.Uses`
    before converting `=` to `:=`)
20. Consider contributing a PR to golangci-lint-auto-configure (don't re-add
    linters absent from both `enable` and `disable`)

### BuildFlow improvements (different repo)

21. File feedback to BuildFlow: add post-repair `go build ./...` gate
22. File feedback to BuildFlow: daemon should run `go test` before committing
23. Consider a `.buildflowignore` or per-file repair exclusion mechanism

### Testing improvements

24. Add a test that verifies `visualtest/doc.go` uses the split-assignment
    pattern (regression guard against fatcontext reverting it)
25. Add a test that verifies `.golangci.yml` `disable:` list contains the three
    linters (stronger than the current "not in enable" check)
26. Consider golden-testing `.golangci.yml` structure against an expected schema
27. Add an integration test that runs a BuildFlow repair cycle and asserts no
    regression in `.golangci.yml` or `visualtest/doc.go`

### Process improvements

28. Add "run `git diff` after every BuildFlow commit" to the session checklist
29. Add "apply workarounds before committing, not after" to personal workflow
30. Add "check upstream issues before diagnosing tool-induced regressions"
    (would have found fatcontext #100 in 2 minutes instead of 8+ sessions)
31. Consider a pre-push hook that runs the full test suite
32. Document the "autofix regression loop" pattern in AGENTS.md as a known
    anti-pattern to watch for

### Broader cleanup

33. Audit all `sync.Once` closures for similar `:=` shadowing risks
34. Audit all `.golangci.yml` linters against the three tiers in
    golangci-lint-auto-configure to find other potential conflicts
35. Consider whether `godoclint` and `testableexamples` should be proposed for
    `NeverAutoEnableLinters` or a new tier in golangci-lint-auto-configure
36. Review whether any other linters in the enable list are incompatible with
    templ projects
37. Consider a templ-project preset for golangci-lint-auto-configure
38. Review the status report from earlier this session
    (`2026-07-30_21-49_fix-visualtest-compile-and-golangci-linter-regression.md`)
    and reconcile/update it with this report's findings

---

## g) Questions I Cannot Answer Myself

1. **Should I apply the `disable:` workaround now and commit both fixes
   together?** The workaround (add `godoclint`, `ireturn`, `testableexamples`
   to `linters.disable`) should stop the regression loop. The split-assignment
   fix for `visualtest/doc.go` is ready. I can commit both now, but I need to
   know if you want to review them first or if I should just go.

2. **Should I check out the fatcontext #43 issue and consider contributing a
   PR?** The root cause is clear (`getSuggestedFixes` in `analyzer.go`
   unconditionally sets `token.DEFINE`). The fix is ~10 lines (check
   `pass.TypesInfo.Uses` before converting). This would end the `:=` regression
   loop permanently across all projects, not just this one. But it's a
   different repo and scope.

3. **Should the `ireturn:` settings block be kept or deleted when moving
   `ireturn` to `disable`?** The current fix deletes it (orphaned settings for
   a disabled linter). But if a future developer re-enables `ireturn`, they'd
   need to re-create the settings. I cannot determine your preference for
   "clean config now" vs "preserve settings for potential future re-enable."
