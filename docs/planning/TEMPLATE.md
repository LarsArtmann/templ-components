# Plan: <SHORT-NAME>

Copy this file to `<YYYY-MM-DD>_<HH-MM>_<SHORT-NAME>.md` (see siblings for
examples) and fill every bracketed field. Delete sections you cannot justify —
an empty section is a lie about the plan. The four checklist gates below come
from `docs/plan-authoring-checklist.md` (rationale + post-mortem sources
there); AGENTS.md keeps the summary.

## Context

[What exists today, what hurts, and the evidence — report §, TODO #, survey,
or e2e finding. Link the source document.]

## Goal / non-goals

- Goal: [the one measurable outcome]
- Non-goals: [explicitly out of scope, with the reason]

## Phases (Pareto order: 1% → 4% → 20% → remaining)

### M<N> — <name> (medium 30–100 min | micro ≤12 min)

- What: [change, files]
- Proof: [the test that turns red→green]
- Checklist (answer EVERY gate per task — see plan-authoring-checklist.md):
  - [ ] **Goldens cover this** — grep `visualtest/` + `testdata/` for existing
        pixel coverage BEFORE assuming it; name the golden files this task
        adds/updates.
  - [ ] **Wired ⇒ e2e** — renders `wire.Action`/`hx-*`/`data-on:*`? Then the
        chromedp e2e task is IN THIS PLAN, or a written waiver lives here:
        [waiver text]. String-proven ≠ browser-proven.
  - [ ] **Counts bump in the same edit** — adds goldens/enums/components?
        List the README/FEATURES/ROADMAP/AGENTS number lines to touch
        alongside the CHANGELOG entry (`TestDocsCountDrift` +
        `TestFeaturesEnumTableExhaustive` enforce).
  - [ ] **Demo smoke is a gate** — demo endpoint/markup change? The
        `visualtest/tools/smoke` check (or `nix run .#visual`) runs before the
        task is called done.
- ⫱ [owner gate: the decision that must be Lars's, planned to the gate edge]

## Verification (whole plan)

- [ ] Per-module loop green: `for mod in utils icons errorpage charts/echarts datastar htmx; do (cd "$mod" && GOWORK=off go test ./...); done`
- [ ] Root build+test via `scripts/ci-repro.sh` (CI's Build & Test, step-for-step)
- [ ] Demo smoke over changed endpoints
- [ ] CHANGELOG `[Unreleased]` warm; drift counts bumped in the same edits

## Deferred / follow-through seeds

[Items discovered mid-plan that did not fit — each with a TODO_LIST ID when
harvested.]
