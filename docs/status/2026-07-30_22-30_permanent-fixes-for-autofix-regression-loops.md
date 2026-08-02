# Status Report — 2026-07-30 22:30

## Session Scope

Continuation of the BuildFlow autofix regression battle. This session covered
three full cycles of fix → BuildFlow reverts → diagnose → fix again, plus an
upstream feedback filing and a test rewrite. The core problem: two BuildFlow
repair tools (fatcontext autofix + golangci-lint-auto-configure repair) keep
breaking committed config that is already correct.

---

## Timeline (this session only)

| Time   | Event                                                                                                      |
| ------ | ---------------------------------------------------------------------------------------------------------- |
| ~21:44 | Fixed `visualtest/doc.go` `:=` → `=` and removed 3 linters from `.golangci.yml` enable list                |
| ~21:49 | Wrote first status report                                                                                  |
| ~21:59 | Committed (`d2efac6`). BuildFlow pre-commit immediately re-introduced BOTH regressions in the working tree |
| ~22:00 | Traced `:=` root cause to fatcontext issue #100 (dup of #43). Applied split-assignment fix.                |
| ~22:10 | `.golangci.yml` regression hit AGAIN (7th+ time). Re-removed linters.                                      |
| ~22:15 | Filed upstream feedback to golangci-lint-auto-configure repo                                               |
| ~22:19 | Wrote second status report                                                                                 |
| ~22:25 | BOTH regressions hit AGAIN. Plus new YAML corruption (indentation broken by repair tool).                  |
| ~22:30 | Applied three permanent fixes: `//nolint:fatcontext`, `disable:` list, section-aware test. All verified.   |

---

## a) FULLY DONE

### 1. `visualtest/doc.go` — `//nolint:fatcontext` (permanent fix, 3rd attempt)

**History of attempts this session:**

1. `:=` → `=` — fatcontext autofix reverted it to `:=` on next commit
2. Split assignment (`allocCtx, cancel := ...` then `sharedAllocCtx = allocCtx`) — fatcontext autofix converted `sharedAllocCtx = allocCtx` to `sharedAllocCtx := allocCtx`, shadowing the package var
3. `//nolint:fatcontext` comment on the `=` line — suppresses the diagnostic entirely so the autofix never fires

**Status:** Applied and verified (`go build ./...` in visualtest passes).

### 2. `.golangci.yml` — moved 3 linters to `disable:` (permanent fix)

**History:**

1. Removed from `enable` only — repair re-added them (7+ times)
2. Now added to `disable:` list — repair respects `disable:` per the 2026-07-25 fix in golangci-lint-auto-configure

**Status:** Applied. `golangci-lint linters` parses the config successfully.

### 3. `.golangci.yml` — YAML corruption fixed

**Problem:** BuildFlow's repair step broke YAML indentation on three entries:
`godot`, `loggercheck`, `testifylint` got nested-list indentation (6 spaces
instead of 4), which golangci-lint parsed as unknown linters
(`gocyclo - godot`, etc.).

**Fix:** Corrected indentation to 4 spaces for all list items.

### 4. `lint_config_test.go` — made section-aware

**Problem:** The `TestGolangciDisabledLinters` test matched any `- godoclint`
line anywhere in `.golangci.yml`, including the `disable:` section where the
linters now legitimately live.

**Fix:** Rewrote the test to track the current YAML section (`enable:` vs
`disable:`) and only flag disabled linters found in `enable:`.

**Status:** Verified — test passes with linters in `disable:`.

### 5. Filed upstream feedback to golangci-lint-auto-configure

Wrote comprehensive feedback at
`/home/lars/projects/golangci-lint-auto-configure/docs/feedback/new/2026-07-30_repair-re-adds-linters-removed-from-enable.md`
explaining the regression loop and proposing 5 solutions.

---

## b) PARTIALLY DONE

### Nothing is uncommitted-but-unverified

All three fixes are applied and verified in isolation. However:

- **Not committed** — the working tree has uncommitted changes across 3 files
- **Not tested against a full BuildFlow cycle** — the `//nolint:fatcontext` and
  `disable:` approaches are theoretically sound but have not been proven
  against an actual BuildFlow pre-commit run

---

## c) NOT STARTED

1. **Committing the three fixes** — waiting for user to say "commit"
2. **Full `go test ./...`** — only ran targeted tests
3. **`golangci-lint run`** — verified config parses, but did not run full lint
4. **`nix run .#verify`** — full pipeline not run
5. **Verifying fixes survive BuildFlow** — the whole point of these fixes is to
   survive the next `repair` cycle; untested
6. **Updating AGENTS.md** — fatcontext #100 root cause, `//nolint` workaround,
   `disable:` pattern, section-aware test, all undocumented
7. **Updating the two earlier status reports** from this session — they are now
   partially stale (their "next steps" were superseded)

---

## d) TOTALLY FUCKED UP

### Failed to apply permanent fixes on the first attempt

The session went through **three cycles** of the same pattern:

1. Apply symptom fix → BuildFlow reverts → diagnose deeper → apply better fix
2. Apply better fix → BuildFlow reverts differently → diagnose deeper → apply permanent fix

**What I should have done from the start:**

- For fatcontext: `//nolint:fatcontext` is the standard Go way to suppress a
  linter. I should have used it immediately instead of trying code restructuring
  (split assignment) that the autofix could still target.
- For `.golangci.yml`: Adding to `disable:` was identified as the workaround in
  the **first** status report (21:49) but I didn't apply it until the **third**
  regression cycle (22:30). Two regression cycles could have been avoided.

### The split-assignment "fix" was wrong

I spent time applying a split-assignment workaround that I claimed was
"permanent" — fatcontext immediately proved it wasn't by converting the second
`=` to `:=`. I should have known that fatcontext targets ALL assignment patterns,
not just function calls. The `//nolint` approach is the only reliable suppression.

---

## What I Forgot / Could Have Done Better

### 1. `//nolint:fatcontext` should have been the first attempt

The codebase already uses `//nolint` comments elsewhere. I diagnosed the
root cause (fatcontext #100) but didn't reach for the standard Go suppression
mechanism until the third attempt. Instead I tried two code restructurings that
the autofix defeated.

### 2. `disable:` workaround should have been applied immediately

I identified the `disable:` workaround in the first status report but offered
it as a question instead of applying it. Two more regression cycles happened
before I actually applied it.

### 3. Did not anticipate YAML corruption

The BuildFlow repair tool can corrupt YAML indentation when re-adding entries.
This was a new failure mode I hadn't seen before. The `.golangci.yml` now has
correct indentation, but there's no guarantee the repair tool won't corrupt it
again.

### 4. Test was not updated when the strategy changed

When I moved linters to `disable:`, the guard test immediately failed because
it was section-blind. I should have updated the test in the same edit as the
config change, not as a follow-up after discovering the test failure.

### 5. Still haven't run the full test suite or verification pipeline

Three status reports in this session, and none of them were backed by a full
`go test ./...` run. Only targeted tests.

### 6. Did not update AGENTS.md

Three permanent fixes applied (`//nolint`, `disable:` pattern, section-aware
test) and none are documented in AGENTS.md. Future sessions will not know why
these patterns exist.

---

## e) WHAT WE SHOULD IMPROVE

### The "autofix fighting" anti-pattern needs a playbook

This session spent ~45 minutes fighting two tools that "fix" already-correct
code. The pattern is now well-understood:

| Tool                                | Trigger                     | Fix it applies       | Why it's wrong here         | Permanent suppression |
| ----------------------------------- | --------------------------- | -------------------- | --------------------------- | --------------------- |
| fatcontext                          | `ctx, cancel = ...`         | Converts `=` to `:=` | Shadows package-level vars  | `//nolint:fatcontext` |
| golangci-lint-auto-configure repair | Linter absent from `enable` | Re-adds to `enable`  | User intentionally disabled | Add to `disable:`     |

This playbook should be in AGENTS.md so future sessions don't waste time.

### BuildFlow repair tools need a "don't touch what's already correct" principle

Both tools share a fundamental design flaw: they assume their autofix is always
helpful. There's no concept of "this was intentional, don't touch it." The
`//nolint` and `disable:` workarounds are the available signals, but they're
not obvious.

### Guard tests should be section-aware from the start

The `TestGolangciDisabledLinters` test was written to match any occurrence of
a linter name, regardless of YAML section. When the strategy changed from
"remove from enable" to "move to disable," the test broke. YAML-section-aware
parsing should have been in the original test.

---

## f) Next Steps (Up to 50)

### Immediate (blocking)

1. **Commit the three fixes** (`visualtest/doc.go`, `.golangci.yml`,
   `lint_config_test.go`) — waiting for user instruction
2. **Run `go test ./...`** to verify the full suite
3. **Run `git diff` after commit** to verify BuildFlow doesn't revert anything
4. **Run `nix run .#verify`** for the full pipeline

### Documentation (high value)

5. Update AGENTS.md with fatcontext #100 root cause and `//nolint` workaround
6. Document the `disable:` pattern for linter suppression in AGENTS.md
7. Document the section-aware test pattern in AGENTS.md
8. Add the "autofix fighting anti-pattern" playbook to AGENTS.md
9. Bump regression counts in AGENTS.md
10. Update/reconcile the two earlier status reports from this session

### Verification

11. Run `golangci-lint run` to confirm 0 findings
12. Run `scripts/check-lint-config.sh` after `.golangci.yml` edit
13. Verify `//nolint:fatcontext` survives a BuildFlow pre-commit cycle
14. Verify `disable:` list survives a BuildFlow repair cycle
15. Run `govalid-generate` BuildFlow step to confirm success
16. Run `nix flake check` for format verification

### Upstream

17. Monitor golangci-lint-auto-configure feedback for response
18. Check if fatcontext #43 has been fixed upstream
19. Consider contributing a PR to fatcontext (check `TypesInfo.Uses` before converting)
20. Consider adding `godoclint`/`testableexamples` to `NeverAutoEnableLinters` in golangci-lint-auto-configure
21. File feedback to BuildFlow about YAML corruption on repair

### Testing improvements

22. Add a test verifying `.golangci.yml` `disable:` list contains the 3 linters
23. Add a test verifying `//nolint:fatcontext` is present on the allocator line
24. Add an integration test that runs a simulated BuildFlow repair cycle
25. Consider a `.golangci.yml` schema validation test (catches YAML corruption)

### Process improvements

26. Add "use `//nolint` for linter autofix suppression, not code restructuring" to personal workflow
27. Add "apply workarounds immediately, don't offer as questions" to personal workflow
28. Add "run `git diff` after every BuildFlow commit" to session checklist
29. Add "update guard tests when changing suppression strategy" to personal workflow
30. Consider a pre-push hook that runs the full test suite

### Broader

31. Audit all `//nolint` comments in the codebase for correctness
32. Audit all `.golangci.yml` entries for potential repair conflicts
33. Consider whether other linters in the enable list are incompatible with templ
34. Review whether the repair tools can be configured to be less aggressive
35. Consider a templ-project preset for golangci-lint-auto-configure
36. Review if BuildFlow can exclude `.golangci.yml` from repair steps
37. Consider adding a `.golangci-lint-auto-configure.yml` sidecar for policy enforcement
38. Update the feedback file if the `disable:` workaround proves effective
39. Check if the YAML corruption is a known bug in golangci-lint-auto-configure
40. Consider whether the visualtest module needs its own lint config

---

## g) Questions I Cannot Answer Myself

1. **Should I commit now and verify against BuildFlow, or verify first and
   commit after?** The three fixes are verified in isolation but untested
   against a real BuildFlow cycle. Committing triggers BuildFlow's pre-commit
   hook, which is the real test — but if it fails, the working tree gets dirty
   again. I cannot determine whether you want to "commit and see" or "verify
   manually first."

2. **Should the `disable:` entries in `.golangci.yml` have inline comments
   explaining why each linter is disabled?** The current `disable:` list is
   bare (`- godoclint` with no reason). Inline YAML comments would help future
   developers, but I don't know if BuildFlow's repair tool would strip them.

3. **Should I update the two earlier status reports from this session to mark
   them as superseded?** Reports at 21:49 and 22:19 describe intermediate
   states that are now outdated. Leaving them as-is creates a misleading
   historical record, but editing them violates the "don't rewrite history"
   principle for point-in-time reports.
