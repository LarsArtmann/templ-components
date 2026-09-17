# Plan-Authoring Checklist

Carried from the 2026-09-16 kanban-vibe execution post-mortem
(`docs/status/2026-09-16_19-12_kanban-vibe-lessons-execution-status.md` §e1/§e2,
`docs/status/2026-09-17_05-54_kanban-followthrough-hardening-status.md` §e).
AGENTS.md keeps the four-line summary; this file keeps the rationale. Every execution
plan that adds components, fields, or demo endpoints answers these per item:

1. **Do the goldens actually cover this?**
   Grep `visualtest/` route/section goldens and `testdata/` before assuming existing
   coverage. The kanban section sat below the fold with zero pixel coverage, and the
   2026-09-16 plan's "route goldens WILL change" premise was false — one grep at plan
   time would have caught it.

2. **Wired ⇒ e2e.**
   Anything rendering consumer-wired attributes (`wire.Action`, `hx-*`/`data-on:*`) gets
   a chromedp e2e task IN THE SAME plan, or an explicit written waiver in the plan.
   String-proven ≠ browser-proven: the demo add button's missing
   `hx-target`/`hx-swap="outerHTML"` survived regenerate+build+unit tests and only
   surfaced in a real click (`TestKanbanE2EAddAndResetBothTransports`).

3. **Drift counts are predictive, not reactive.**
   A task that adds goldens/enums/components bumps the README/FEATURES/ROADMAP/AGENTS
   counts in the SAME edit as its CHANGELOG entry — `TestDocsCountDrift` +
   `TestFeaturesEnumTableExhaustive` enforce, but the bump should never be a separate
   repair step.

4. **Demo smoke is a gate, not a bonus.**
   Run the live HTTP smoke (`visualtest/tools/smoke`) or `nix run .#visual` before
   finalizing any demo endpoint change. The unspecified-Method `hx-get` bug survived
   regenerate+build+unit tests; only the live smoke caught it.

Plan skeleton convention: phases ordered by Pareto tiers (1% → 4% → 20% → remaining),
medium tasks 30–100 min, micro-tasks ≤12 min, owner gates marked `⫱` and planned to the
gate edge. `docs/planning/TEMPLATE.md` carries this checklist structurally — copy it for
new plans. See
`docs/planning/2026-09-17_06-00_RELEASE-FIRST-PARETO-MASTER-PLAN.md` for
the reference shape.
