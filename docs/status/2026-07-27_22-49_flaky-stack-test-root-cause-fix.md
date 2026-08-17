# Status Report — 2026-07-27 22:49 CEST

## Session: Flaky `TestStackWithFeedbackComponents` fix

**Scope:** Single failing test from BuildFlow `test-race`. Root-caused, fixed, verified.
**Branch:** `master` (~~working tree dirty — one file changed, **not committed**~~ committed as `da156d6` on 2026-07-27).
**Commit:** ~~none yet (rule: never commit without explicit ask)~~ `da156d6` — "hell composition integration tests" (BuildFlow-daemon-authored message; the fix itself is correct).

---

## 1. What happened (the run)

BuildFlow `test-race` reported a single failure:

```
--- FAIL: TestStackWithFeedbackComponents (0.01s)
    appshell_composition_test.go:107: Stack must emit flex flex-col (1D, not grid)
FAIL    github.com/larsartmann/templ-components/integration
```

### Investigation path

1. Read failing test → asserted contiguous substring `"flex flex-col"`.
2. Read `layout/stack.templ` + `stack_types.go` → source clearly emits `flex flex-col`.
3. Rendered in isolation → always correct (`test-stack flex flex-col space-y-4`).
4. Ran suite 30× under `-race` → **4/30 failures (~13%)**, always same assertion.
5. Concurrency theories: stress-tested `utils.Class` with 3M varied calls → **0 corruption** (mutex is correct).
6. Stress-rendered `Stack` 100k× concurrently → **0 corruption**.
7. Instrumented the test to print output on failure → **captured the smoking gun**:

   ```
   <div class="test-stack flex-col space-y-4 flex"></div>
   ```

   `flex` was reordered to the end by `tailwind-merge-go`. Both tokens present — contiguous assertion wrong.

### Root cause

`utils.Class()` wraps `tailwind-merge-go`, which reorders Tailwind classes into canonical groups **non-deterministically** (output depends on LRU cache state, which depends on call history — hence shuffle/ordering-dependent). The race detector stays silent because access is mutex-guarded; the nondeterminism is logic-level, not a data race.

This is the **exact gotcha already documented in AGENTS.md**:

> Do NOT assert ordered class substrings in tests — `utils.Class`/tailwind-merge reorders classes; use `utils.AssertContainsAll` for multi-token checks.

### Fix applied

`integration/appshell_composition_test.go:106-109` — replaced two fragile `strings.Contains` checks with one `utils.AssertContainsAll(t, output, "flex", "flex-col", "space-y-4")`, matching the convention already used in `layout/stack_test.go`.

### Verification

- **0 failures / 40 iterations** under `go test -race ./integration/` (was 4/30).
- Full `go test ./...` green (no caches invalidated unexpectedly).
- Build clean (no orphaned `strings` import — still used elsewhere in file).

---

## 2. FULLY DONE

| # | Item                                | Evidence                                     |
| - | ----------------------------------- | -------------------------------------------- |
| 1 | Root-caused the flake               | captured `flex-col space-y-4 flex` output    |
| 2 | Proved `utils.Class` is thread-safe | 3M varied concurrent calls, 0 corruption     |
| 3 | Proved `Stack` component is correct | isolated render always emits `flex flex-col` |
| 4 | Fixed the broken assertion          | `AssertContainsAll` with independent tokens  |
| 5 | Verified fix kills the flake        | 0/40 under `-race`                           |
| 6 | No regression elsewhere             | full `go test ./...` green                   |

---

## 3. PARTIALLY DONE

| # | Item                          | Why partial                                                                                                                                                                                                                                                                |
| - | ----------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **Comment in the fixed test** | Added a 3-line `//` block explaining _why_ we use `AssertContainsAll`. Defensible (explains why, documents a non-obvious convention, AGENTS.md permits "why" comments) but AGENTS.md also says "NEVER ADD COMMENTS unless asked." Borderline — flagging for user judgment. |
| 2 | **Verification scope**        | Ran `go test ./...` (non-race) at the end + 40× `-race` on `integration/` only. Did **not** run 40× `-race` on every other package — other packages may have their own shuffle-dependent flakes that the original BuildFlow run would surface.                             |

---

## 4. NOT STARTED (gaps I noticed but did not touch)

| # | Gap                                                                                                                                                                                                                                                                                                       | Impact                       |
| - | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------- |
| 1 | **Did NOT scan repo for other tests with the same anti-pattern.** This is the biggest miss. AGENTS.md documents the rule; one test violated it; almost certainly more exist. Any `strings.Contains(output, "X Y")` where `X` and `Y` are Tailwind classes passed through `utils.Class` is a latent flake. | High — systematic fragility. |
| 2 | **No CHANGELOG `[Unreleased]` entry.** AGENTS.md release convention: every fix on `master` must add its changelog entry immediately, not defer.                                                                                                                                                           | Medium — drift risk.         |
| 3 | **No automated guard.** The convention lives in AGENTS.md prose only. No grep-test, no linter rule prevents re-introduction.                                                                                                                                                                              | Medium — recurrence certain. |
| 4 | **Did not escalate the tailwind-merge-go nondeterminism.** The library emits order-dependent-on-LRU output. Options (sort output in `utils.Class`, pin LRU, wrap with a stable layer) unexplored.                                                                                                         | Low (architectural).         |

---

## 5. TOTALLY FUCKED UP

Nothing destructive. No reverts, no force-pushes, no data loss, no unrelated files touched. Honest missteps:

| # | Misstep                                                                                                                                                                                                                                                                                                                           | Lesson                                                                                                                                                                      |
| - | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **Over-investigated exotic concurrency theories before capturing the output.** I wrote 3 separate stress tests (2 for `utils.Class`, 1 for `Stack`) before simply instrumenting the failing test to print its bad output — which instantly revealed the reordering.                                                               | When a test fails, **capture the actual output FIRST.** The answer was one `t.Errorf("...%q", output)` away. ~6 tool calls of speculative stress-testing could have been 1. |
| 2 | **Quoted the AGENTS.md rule early, then failed to apply it immediately.** I literally pasted _"Do NOT assert ordered class substrings"_ into my second message — then spent the next several steps proving the component and `Class` were innocent before fixing the _assertion_. The rule told me the answer; I didn't trust it. | Trust documented conventions over novel investigation.                                                                                                                      |
| 3 | **Did not run the exact BuildFlow `test-race` command.** Reproduced with `go test -race` (close, but BuildFlow may set flags/env differently). The original failure path is not bit-for-bit verified.                                                                                                                             | Reproduce via the real failing command, not a near-neighbor.                                                                                                                |

---

## 6. WHAT WE SHOULD IMPROVE

### A. Immediate (this fix's loose ends)

1. **Grep the whole `*_test.go` tree for `Contains(.* .*)` patterns on Tailwind classes** — find every fragile ordered-substring assertion. Likely 5–15 instances.
2. **Add CHANGELOG `[Unreleased]` entry** — `fix(tests): stop asserting ordered Tailwind class substrings (flaky under tailwind-merge reordering)`.
3. **Decide on the comment** in `appshell_composition_test.go:106-108` — keep, shorten, or drop per AGENTS.md "no comments" rule.

### B. Structural (prevent recurrence)

4. **Add `utils.TestNoOrderedTailwindSubstringsInTests`** — a repo-wide drift-guard test (like the existing `TestDarkModeCompliance`, `TestMotionReduceCompliance`) that scans `*_test.go` for `Contains("X Y")`/`AssertContains("X Y")` where both tokens look like Tailwind classes, and fails. One-time write, permanent protection.
5. **Add a `//nolint` or linter rule** if a static analyzer can catch this (e.g., a custom `golangci-lint` go/analysis pass — heavier lift).
6. **Document the pattern in a dedicated ADR** (`docs/adr/0NNN-test-class-assertion-convention.md`) so the _why_ lives next to the rule, not only in AGENTS.md prose.

### C. Deeper (the real source)

7. **Investigate making `utils.Class` output deterministic.** If tailwind-merge-go exposes a stable/deterministic mode (no LRU, or sorted output), flipping it on would make _all_ tests reproducible and let ordered assertions work again. Worth a 30-min spike.
8. **File/track upstream** if tailwind-merge-go's LRU nondeterminism is considered a bug by its maintainers.

### D. Process

9. **Always capture actual failing output before theorizing.** Make this a session-level habit — one `t.Errorf` with `%q` beats five stress tests.
10. **Trust documented gotchas.** When AGENTS.md says "X reorders classes," and a test asserts ordered classes, the diagnosis is done.
11. **Run the exact failing command** (BuildFlow `test-race`), not a manual approximation, before declaring "fixed."

---

## 7. Up to 50 things to do next (scoped, honest)

Ordered roughly by impact × ease.

### High impact, this fix's blast radius

1. Grep `*_test.go` repo-wide for `Contains("… …")` / `Assert(Contains).*"… …"` on Tailwind tokens → list all fragile assertions.
2. Fix each found instance → `AssertContainsAll` with independent tokens.
3. Add `utils.TestNoOrderedTailwindSubstringsInTests` drift-guard.
4. Add CHANGELOG `[Unreleased]` entry for this fix + the batch.
5. Decide comment fate in `appshell_composition_test.go:106-108`.
6. Run BuildFlow `test-race` (the real command) to confirm the original failure is gone.

### Verification hardening

7. Run `go test -race -count=20 ./...` per-package to surface any other shuffle flakes.
8. Add a CI matrix entry that runs `-race -count=N` nightly (flakes only show under repetition).
9. Spike: can tailwind-merge-go run deterministically? (sort output / disable LRU).
10. If yes → wrap in `utils.Class`, make ordered assertions safe again (or at least tests stable).

### Documentation / convention

11. ADR for the test-assertion convention.
12. Cross-link AGENTS.md rule ↔ the new ADR ↔ the new drift-guard test.
13. Add a "flaky-test playbook" note: _capture output → check for tailwind-merge reordering → use `AssertContainsAll`_.
14. Audit AGENTS.md for other "do not assert X" rules that lack automated enforcement.

### Broader test-quality sweep (discovered incidentally)

15. The repo has `TestDarkModeCompliance`, `TestMotionReduceCompliance`, `TestSkillComponentCount`, `TestVersionMatchesChangelog` — consider a `TestTestClassAssertionHygiene` in the same family.
16. Check `internal/golden` tests — do they normalize class order? If yes, they're immune; if no, they may also flake.
17. Check `visualtest` — pixel diffing is order-independent, so safe; confirm no string assertions sneak in.
18. Audit `integration/` for other contiguous-substring assertions on rendered HTML.
19. Audit `recipes/` tests (composition screens) — high chance of multi-class assertions.

### Tooling

20. Add a pre-commit hook snippet that greps for `Contains("[a-z]+-[a-z]+ [a-z]+-[a-z]+"` in `_test.go` (cheap heuristic).
21. Consider `rules_go` / custom analyzer if heuristic proves noisy.
22. Add `just`/flake target `nix run .#test-flake` = `go test -race -count=20 ./...` for quick flake hunting.

### Cleanup / housekeeping

23. Verify no other diagnostic edits left in tree (I added+removed several `_test.go` probes — confirm all cleaned).
24. `git status` review before any commit.
25. If committing: one commit, message per AGENTS.md release convention.

### Stretch (lower priority, mentioned for completeness)

26. Document tailwind-merge-go's nondeterministic behavior in `docs/research/`.
27. Benchmark: does the LRU actually buy speed worth the nondeterminism?
28. Pin tailwind-merge-go version in `go.mod` (check if a newer release fixes ordering).
29. Survey other templ UI libs (e.g., flowbite-go) for how they assert classes in tests.
30. Consider a `utils.AssertClassesEqual(t, got, want)` helper that sorts before comparing (canonical comparison).
31. If (30) added, migrate existing `AssertContainsAll` calls where exact-set comparison is intended.
32. Add fuzz test for `utils.Class` (random class subsets → output always contains all non-conflicting tokens).
33. Property test: `Class(a, b)` output ⊇ set(non-conflicting tokens from a, b).
34. Check whether the `//nolint:gochecknoglobals` on `classMu` is still needed after any refactor.
35. Doc: add "why is my test flaky?" FAQ entry pointing at tailwind-merge.
    36–50. _(Reserve for findings from the repo-wide grep in #1 — each new fragile assertion becomes its own fix item.)_

---

## 8. Questions I CANNOT figure out myself

1. **Scope of the batch fix:** Should I now grep the entire `*_test.go` tree and fix _every_ fragile ordered-substring Tailwind assertion in one pass (could be 5–15 files), or do you want to review the list first before I touch them?
2. **Comment policy for the fixed test:** Keep the 3-line explanatory `//` block I added (it documents _why_ `AssertContainsAll`), or strip it to honor the strict "NEVER ADD COMMENTS" rule now that a drift-guard test (#4) would carry the enforcement?
3. **CHANGELOG granularity:** Single `[Unreleased]` entry covering "this fix + the batch", or one entry per file touched?

---

## 9. TL;DR

- **The flake is dead** (0/40 under `-race`, was ~13%).
- **Root cause was a wrong test, not wrong code** — `Stack` and `utils.Class` were always correct.
- **The fix is one line of assertion logic** in `integration/appshell_composition_test.go`.
- **Not committed** (awaiting explicit go-ahead).
- **Biggest gap:** I did not check whether _other_ tests share the same anti-pattern — that's the likely-next-flake and the highest-value follow-up.

---

## Resolution (2026-07-28)

The fix shipped as commit `da156d6` (working tree clean). The CHANGELOG `[Unreleased]` now records it ("Flaky `TestStackWithFeedbackComponents` root cause fixed").

**Open items from §4 / §7 that remain open (harvested to TODO_LIST):**

- The repo-wide ordered-Tailwind-substring audit (§4.1 / §7.1–2) was **not** done — it is now **TODO #81** ("Audit repo-wide for ordered-Tailwind-substring test assertions" + `TestNoOrderedTailwindSubstringsInTests` drift-guard).
- The CHANGELOG entry (§4.2) is now added.
- Deeper investigations (deterministic `utils.Class`, upstream tailwind-merge-go LRU, ADR for the assertion convention) are deferred — not yet routed to a TODO; see TODO #81 if the audit surfaces a need.
