# Accessibility Gate Policy (axe-core)

**Status:** decided 2026-09-14 (plan item M05/F030). The harness this governs
lives in `visualtest/axe.go` + `visualtest/axe_sweep_test.go` and runs as part
of `nix run .#visual` and the CI Visual Regression job.

## The decision: default-fail, not opt-in

The axe sweep **fails by default** on any `critical`/`serious` violation found
on a demo route that is not explicitly accepted in the ledger
(`visualtest/testdata/axe_baseline.json`). There is no per-component opt-in
flag and there will not be one.

Why default-fail:

- **Opt-in gates rot.** An a11y check that a component author must remember
  to enable is a check that ships disabled — the moat argument (a11y is this
  library's stated differentiator) demands the inverse default.
- **The demo routes are the highest-leverage audit surface.** They are the
  composition patterns consumers copy; a violation there propagates into
  every consuming codebase. The sweep covers every route the demo registers,
  so coverage grows with the demo, not with author diligence.
- **The gate proves it can fail.** `TestAxeHarnessDetectsViolations`
  (deliberately invalid markup) pins the detector itself — the
  silently-green-gate failure mode is covered. The 2026-09-14 run is the
  proof the policy works: the sweep caught the Calendar MonthNav roleless
  `aria-label` arrows that every string/golden test passed.

## The severity line

| Impact   | Gate?  | Rationale                                                           |
| -------- | ------ | ------------------------------------------------------------------- |
| critical | BLOCK  | Unusable for the affected assistive-tech/mode.                      |
| serious  | BLOCK  | Blocks completion or understanding for real user groups.            |
| moderate | log    | Real but non-blocking; kept visible in output for pressure.         |
| minor    | log    | Stylistic/best-practice; noise-to-signal drops sharply below this.  |

Blocking on moderate/minor would bury regressions in churn; logging keeps the
signal without gating on taste.

## The baseline ledger is documented debt, not a pass

Format: `{"<route>": {"<rule>|<impact>": <budget>}}`

- A **positive budget** caps the accepted node count — an increase beyond it
  fails (debt cannot silently grow).
- **`-1` accepts unconditionally** — reserved for live pages where node counts
  fluctuate (polled regions); every `-1` is a candidate for conversion to a
  budget once the route stabilizes.
- Every accepted entry carries a justification in the ledger's header comment
  or the accepting commit. Today's set is palette-convention debt
  (`color-contrast` on muted captions, the `-600`/`-500` shade convention)
  tracked for owner review (TODO_LIST #175 follow-up) — fixing it means a
  deliberate visual release, not a drive-by shade tweak.
- **Pruning rule:** when an accepted rule stops producing findings (library
  re-shade, axe-core bump, markup fix), DELETE the ledger entry in the same
  change. An entry that protects nothing hides future regressions behind a
  stale accept.

## Change protocol

- **New demo route** → audited automatically on the next sweep run. Route
  authors fix forward or accept debt with justification; never weaken the
  sweep to land a route.
- **axe-core bump** → re-run the sweep, re-triage every ledger entry
  (same discipline as the vnu ignore classes, TODO_LIST #216).
- **A sweep failure** → fix the markup. Extending the ledger is the last
  resort and must state why the violation cannot be fixed now.
