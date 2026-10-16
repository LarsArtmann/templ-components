# Tool-trial protocol & CI tooling budget

Standing rule (M24, plan 2026-10-15): **no new tool or lane enters the repo's
test surface without a named owner and kill-criteria.** This file is the
register.

## The protocol

1. **Spike** — outside the Go test surface (e.g. `tools/`-side or a scratch
   module), time-boxed, artifacts + findings written down
   (`docs/testing/tooling-evaluation.md`).
2. **Verdict** — adopt / trial / reject, with evidence links, in the evaluation
   doc + ADR (gates ⫱G1–⫱G4 where the plan assigns them).
3. **ADR** — the adoption/rejection decision and its boundary recorded
   (`docs/adr/`).
4. **Integrate OR log-rejection** — adopted: wired as a lane with a budget-table
   row (below) BEFORE its first CI run; rejected: entered in the rejected-tools
   log (below) so it cannot be re-recommended innocently.
5. **Kill** — any lane can be killed by its kill-criteria without ceremony;
   delete the lane + ledger entries, keep the ADR trail.

## CI tooling budget table (μ79)

| Lane / tool                                                                                                  | Owner                       | Kill-criteria                                                                                                                                                                                                                                               |
| ------------------------------------------------------------------------------------------------------------ | --------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `nix run .#crossbrowser` (3-engine screenshots, dispatch-only CI draft `.github/workflows/crossbrowser.yml`) | Lars (⫱G1 pending)          | Zero engine-specific findings across 2 consecutive quarters AND no consumer reports of engine breakage → delete lane, keep spike scripts for on-demand use. Also kill if nixpkgs playwright lockstep breaks twice in a quarter (maintenance exceeds value). |
| IBM accessibility-checker (TRIAL, ⫱G2 pending)                                                               | Lars                        | Adopt permanently only if it finds real issues axe missed in 2 consecutive runs (§5.3 already banked one such class). Kill if its findings stay within noise classes for 2 quarters.                                                                        |
| SSE load (k6 nightly or `sseload`)                                                                           | Lars (⫱G4 = no nightly now) | Adopt a nightly if consumers report live-region disconnects, the stream handler grows real fan-out logic, or #329's canonical-URL buffering is fixed. Until then: on-demand `sseload` only.                                                                 |
| Bencher                                                                                                      | Lars                        | Revisit only when a benchmark regression reaches master unnoticed; then trial ONE package (M16 protocol).                                                                                                                                                   |
| stylelint                                                                                                    | Lars (deferred)             | Adopt when the CSS surface grows enough that a regression lands unnoticed; must stay zero-config with the TW4 config.                                                                                                                                       |
| Playwright trace viewer                                                                                      | nobody (skill, not lane)    | n/a — used ad hoc to debug crossbrowser-lane failures.                                                                                                                                                                                                      |

## Rejected tools (μ77) — dead, do not resurface

| Tool                                                            | Rejected                    | Reason (evidence)                                                                                                                                     |
| --------------------------------------------------------------- | --------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| Lost Pixel                                                      | permanently                 | Project archived 2026-04-22, team moved to Figma (verified 2026-10-14 research session).                                                              |
| rustywind                                                       | permanently (for this repo) | M17 dry-run: rewrites 183 files INCLUDING golden files whose class order is canonical per the library's own sorter — churn with zero behavioral gain. |
| k6 + xk6-sse (local)                                            | for this machine            | xk6-sse not in k6's extension registry; custom-binary build needs the Go proxy (blocked locally). CI remains viable.                                  |
| Second long-term pixel suite (any Playwright pixel-golden gate) | permanently                 | Plan guardrail 1 / ADR 0048: chromedp owns pixels; Playwright stays diagnostic.                                                                       |
| Guidepup screen-reader automation                               | parked (ROADMAP seed)       | Windows-runner cost exceeds current value; revisit only with a dedicated a11y milestone.                                                              |
