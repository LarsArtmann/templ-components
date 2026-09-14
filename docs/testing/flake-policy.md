# Test Flake Policy (M19/F088)

How this repo classifies, handles, and eliminates test flakes. The policy
exists because a flaky gate is worse than no gate: teams learn to re-run
red builds reflexively, and real regressions ride the same reflex.

## Classification

| Class             | Example                                     | Handling                                                                                                                                                                                           |
| ----------------- | ------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Infra flake**   | proxy.golang.org 5xx during module fetch    | Retry at the SHELL level (bounded, loud `::warning::`) — see the visualtest compile step in ci.yaml (3 attempts, backoff)                                                                          |
| **Browser flake** | chromedp timing races, font cache staleness | Fix the TEST (readiness gates, retry-until-needle), never the gate. Known classes are documented in AGENTS.md (click-before-reattach, Poll bool-unmarshal, NodeVisible on opacity-0, dialog-pause) |
| **Real flake**    | Data race, time-dependence, map ordering    | Fix the CODE or the test to be deterministic. This library guarantees render determinism (`TestRenderDeterminism`) — a flaky golden is a bug                                                       |
| **Environment**   | missing Chromium, stale fontconfig cache    | Fail LOUD, never skip. Skipping guards protect nothing (AGENTS.md "guards must fail loud")                                                                                                         |

## Rules

1. **A flake is tracked from the moment it is observed.** Open a TODO_LIST
   entry with the failing test name, the CI run link class, and the
   hypothesis. Zero "probably fine" re-runs without a record.
2. **Retries are a LAST resort, bounded, and loud.** The only sanctioned
   retry wrapper is `visualtest.RetryOnce` (below) — it retries exactly
   once, logs both attempts, and fails if the retry also fails. Two
   consecutive failures are never "a flake"; they are a bug.
3. **Gates never weaken for flakes.** No `t.Skip`, no lowered thresholds,
   no `-count=1` removals to "stabilize" CI. If a gate is red for a known
   infra reason, revert the trigger or fix forward — do not blanket-retry.
4. **Three strikes = quarantine decision.** A test that flakes three times
   in 30 days gets an owner decision: fix, rewrite as non-flaky (determinism
   injection, readiness gates), or delete with rationale. Parking-lot
   flaky tests are not an option.
5. **Determinism injection is preferred over waiting.** Examples already in
   the codebase: `RelativeTime.Now` (injectable clock), golden EnsureID
   normalization, `TestRenderDeterminism` double-render pinning, the pure
   fontconfig pin for visual goldens.

## The retry helper

```go
// visualtest/flake_retry_test.go
visualtest.RetryOnce(t, "TestWireE2ESomething", func() error {
    // ... the flake-prone browser flow; return an error instead of using t,
    // so the retry starts from a clean slate
    return nil
})
```

- Retries the body ONCE on failure; both failures fail the test.
- The first-attempt failure is logged with `t.Logf`, never swallowed.
- Intended ONLY for browser-timing flakes whose root cause is documented;
  if a RetryOnce-wrapped test starts failing twice, treat it as class 3
  (real bug) and remove the wrapper while fixing the cause.
- **Dormant by design (status 2026-09-14): zero active call sites.** Every
  flake observed so far was root-caused and fixed instead (the Calendar
  MonthNav settle-window, the parallel-tab allocator contention, the stale
  fontconfig cache) — which is exactly what rule 5 demands. Wrapping a
  currently-green test just to give the helper a call site would weaken it
  (one free failure). The first qualifying browser-timing flake that
  SURVIVES a root-cause attempt gets the first call site; its AGENTS.md
  entry must link here.
