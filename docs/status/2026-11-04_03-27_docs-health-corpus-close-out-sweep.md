# Docs-Health Corpus Close-Out Sweep — 2026-0* reports annotated + archived

**Date:** 2026-11-03 | **Mode:** docs-health AUDIT (VERIFY + HARVEST + ANNOTATE + ARCHIVE)
**Scope:** every non-archived `2026-0*` markdown doc: 185 `docs/status/`, 30 `docs/planning/`,
2 `docs/feedback/`, 1 `docs/architecture-understanding/` — 218 files.

## What ran

1. **VERIFY.** The 739 open table rows in the pre-2026-10 status corpus (plus 550 open rows in the
   planning corpus) were each verified against the current tree (v1.26.0) with grep-level evidence:
   ~40% resolved-done by later sessions, ~7% explicitly not-do/superseded, the rest routed.
2. **HARVEST.** Open residue that is bounded + actionable distilled into TODO_LIST **#488–#497**
   (one bundle per theme: release/infra hardening, keyboard/a11y, dark-mode scan extensions,
   props-API consolidation, table follow-ups, test/consistency, NotFound404 polish, visual
   tooling, docs tails, demo/site tails). Vague/long-term residue routed to ROADMAP
   ("2026-0* corpus residue" section: website-growth ideas, v2-leaning API ideas, testing ideas,
   docs/ADR ideas, M27 remainder). Owner-gated items routed to their existing ⫱ rows
   (#80/#93/#107/#108/#124-126/#157/#190/#211/#212/#232/#303/#312/#411/#443).
3. **ANNOTATE.** Every resolvable row now carries an inline verdict (`done — evidence`,
   `**Won't implement — superseded by X.**`, `**NOT-DO — reason.**`, or `**→ routed — TODO_LIST/ROADMAP #N**`).
   1548 verdicts applied across the corpus via `annotate.go repair` (739 + 116 placeholder-completed + 466 planning + 227 status placeholder rows + spot passes).
4. **ARCHIVE.** 213 fully-resolved files `git mv`'d into their `archived/` siblings
   (180 status, 30 planning, 2 feedback, 1 architecture-understanding). The check-rows gate ran
   BEFORE every move; the moved-file completeness gate (`grep -L '~~'`) reports **0 unmarked files**.

**Kept open (not archived) — 5 status files with genuinely-open rows:**
`2026-09-17_17-28_sidebar-light-mode-session-status.md` (dedicated `/navigation` route ask,
nav-sidebar light eyeball, adaptive-brand guard idea),
`2026-09-17_18-48_wire-action-design-review-session-status.md` (debounce-without-event invariant),
`2026-09-17_20-00_card-sharp-default-squircle-opt-in-status.md` (squircle browser-proof gap b1–b6

- follow-ups 1–9, two owner gates),
  `2026-09-18_07-08_kanban-followthrough-verification-hardening-status.md` (LoadingButton flake fix),
  `2026-09-23_15-56_prerender-drift-flake-fix.md` (prerender routes-table single-sourcing, calendar
  month pin, flake-stress mode, RelativeTime pinning evaluation, prerender-safe marker convention).

**TODO_LIST #212 resolved by this sweep:** the "policy for ~2.1k unannotated name-keyed items in
pre-2026-09-10 reports" owner gate — the numbered/table-row class is now fully annotated; the
remaining bare prose is accomplishment narrative (correctly unmarked), so the policy question is
moot. Struck in TODO_LIST with this report as evidence.

## Tooling notes (for the next corpus pass)

- The routed-marker regex (`**→[^*]+**`) breaks on `*` inside the verdict text — `2026-0*` was
  rewritten marker-free after the first check-rows run surfaced 28 false INCOMPLETE tables.
- The tool's already-annotated refusal caught an off-by-one verdict cluster (P-106/107/108 in the
  2026-06-28 remediation report) before it could land — dry-run first is load-bearing.
- Wholly-unstruck tables are NOT gate offenders (UNTOUCHED is only flagged as PARTIAL/INCOMPLETE
  when mixed) — one executed-but-unmarked plan (SELF-INTEGRATION 2026-09-19) slipped through the
  table gate and was caught by the moved-file `grep -L '~~'` completeness gate instead.

## Archive manifest (213 files)

| Moved to                                                                                              | Verdict | Reason                                                                     |
| ----------------------------------------------------------------------------------------------------- | ------- | -------------------------------------------------------------------------- |
| `docs/architecture-understanding/archived/2026-05-19_22-56-architecture-review.md`                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/feedback/archived/2026-07-05_DiscordSync.md`                                                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/feedback/archived/2026-07-05_swettyswipper-consumer-feedback.md`                                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-04-27_12-44-PARETO-EXECUTION-PLAN.md`                                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-05-17_01-25_alpha-release-critical-path.md`                              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-05-18_execution-plan.md`                                                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-05-20_comprehensive-plan.md`                                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-06-01_19-06_tailwind-v4-theming-pareto-plan.md`                          | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-06-08_execution-plan-all-todos.md`                                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-05_21-27_SUPERB-UI-LIBRARY-UPGRADES.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-07_18-10_comprehensive-execution-plan.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-08_17-15_dark-mode-compliance-plan.md`                                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-11_01-04_DARK-MODE-PACKAGING-AND-FEEDBACK-EXECUTION.md`               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-12_14-45_SUPERB-DOCS-HEALTH-AND-HARDENING.md`                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-12_19-00_MODERN-BROWSER-INTEGRATION.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-14_03-27_precompiled-css-sidebar.md`                                  | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-18_19-42_SUPERB-MULTI-SKILL-SELF-REVIEW-AND-RECOVERY.md`              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-20_22-51_SUPERB-GRID-LAYOUT-ADOPTION.md`                              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-21_03-54_NEXT-LEVEL-PLATFORM-FIRST-ROADMAP.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-07-27_22-52_pareto-hardening-prevention-coverage.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-08-14_16-29_CI-GREEN-AND-STANDOUT-MOAT-PARETO.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-07_16-24_forms-transport-supremacy.md`                                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-08_07-52_cv-templ-components-adoption-pareto-plan.md`                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-09_01-53_kanban-pareto-execution-plan.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-09_04-24_quality-tier-and-release-plan.md`                            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-13_11-14_SUPERB-PARETO-EXECUTION-PLAN.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-13_11-15_astro-to-templ-full-conversion.md`                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-17_06-00_RELEASE-FIRST-PARETO-MASTER-PLAN.md`                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-17_18-20_ERRORPAGE-VISUAL-PARITY-PARETO-PLAN.md`                      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-17_20-39_DRIFT-IMMUNITY-FIRST-PARETO-MASTER-PLAN.md`                  | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-19_15-36_sales-page-pareto-execution-plan.md`                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-19_17-48_SELF-INTEGRATION-PARETO-MASTER-PLAN.md`                      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/planning/archived/2026-09-22_22-42_TC-MIRROR-DUP-GATE-HARDENING-PARETO-PLAN.md`                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-01_21-14_comprehensive-review-and-execution-plan.md`                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-08_02-22_session-3-full-review-and-next-steps.md`                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-08_05-09_session-4-comprehensive-review.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-09_17-27_deduplication-zero-clones.md`                                  | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-09_17-50_deduplication-t30-zero-clones.md`                              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-20_06-26_v1.0-readiness-status-update.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-20_15-11_comprehensive-status-update.md`                                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-20_16-15_reflection-and-status.md`                                      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-20_23-06_post-dedup-and-type-safety-improvements.md`                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-27_22-04_cross-project-adminui-analysis-and-new-components.md`          | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-27_23-47_v0.4.0-shipped-and-adminui-tailwind-infra-ready.md`            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-28_03-34_full-code-review-remediation-status.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-28_04-05_phase-2-remediation-progress.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-28_05-03_dedup-cleanup-and-library-compile-fix.md`                      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-28_06-43_post-dedup-test-isolation-and-release-prep.md`                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-29_10-16_v0.6.0-release-and-pending-work.md`                            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-06-29_11-08_v0.6.0-released-and-post-release-hardening.md`                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-04_19-19_skill-authoring-playbook.md`                                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-05_16-56_final-session-report.md`                                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-05_17-36_session7-consumer-feedback-execution.md`                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-05_20-52_notfound404-component-status.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-05_20-54_session7-self-review-and-gaps.md`                              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-05_20-56_session8-comprehensive-status.md`                              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-05_21-05_discord-sync-feedback-actioned.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-05_22-58_v0.7.0-release-comprehensive-status.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-06_00-18_session10-brutal-review-hardening.md`                          | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-06_01-35_v0.8.0-release-typed-enums-split-brain-fix.md`                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-06_01-36_session11-comprehensive-hardening-js-research.md`              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-06_04-33_coverage-boost-sprint.md`                                      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-06_06-24_dedup-extraction-sprint.md`                                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-06_06-47_review-items-6-and-7.md`                                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-06_15-40_post-v090-hardening-sprint.md`                                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-07_18-02_bug-hunt-critical-fixes.md`                                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-08_04-32_comprehensive-bug-hunt-status.md`                              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-08_17-02_dark-mode-audit.md`                                            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-08_17-55_dark-mode-compliance-sprint.md`                                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-08_18-12_v0.10.0-release-cut.md`                                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-09_09-07_css-integration-cleanup-status.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-10_07-38_todo-list-builder-and-features-audit.md`                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-10_10-49_pareto-popover-and-coverage-gaps.md`                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-10_14-13_json-v2-recurrence-fix-and-prevention.md`                      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-11_09-31_consumer-driven-improvements.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-11_09-40_consumer-driven-dark-mode-and-feedback-execution.md`           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-12_02-34_table-flush-and-cellpadding-status.md`                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-12_15-49_SESSION-EXECUTION-STATUS.md`                                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-12_17-11_POST-V0.17.0-RETROSPECTIVE.md`                                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-12_21-28_DIALOG-MIGRATION-STATUS.md`                                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-12_21-53_MODERN-WEB-STANDARDS-STATUS.md`                                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-13_23-14_docs-health-audit-and-count-drift-fixes.md`                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-14_02-11_public-launch-status.md`                                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-14_02-46_post-commit-brutal-review.md`                                  | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-14_03-08_cloud-run-demo-live.md`                                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-14_04-21_precompiled-css-ci-green.md`                                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-14_13-31_css-source-fix-docker-hardening-demo-polish.md`                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-16_00-31_session-status.md`                                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-18_09-29_v0.18.0-release-postmortem.md`                                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-18_20-27_release-sh-repair-brutal-self-review.md`                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-21_01-24_GRID-LAYOUT-SESSION-BRUTAL-STATUS.md`                          | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-21_02-16_GRID-LAYOUT-POLISH-PASS-STATUS.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-21_07-38_PLATFORM-FIRST-ROADMAP-EXECUTION.md`                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-21_09-33_POST-V1-1-0-DEFECT-FIXES.md`                                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-22_18-51_pareto-execution-status.md`                                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-23_21-00_v1.2.0-release-cut.md`                                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-27_17-04_docs-health-and-update-old-docs-full-pass.md`                  | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-27_20-38_container-query-expansion-status.md`                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-27_20-38_visual-testing-comprehensive-status.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-27_21-15_buildflow-goexperiment-fix-brutal-review.md`                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-27_21-16_docs-health-self-review-and-honest-gaps.md`                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-27_22-49_docs-health-update-old-docs-3rd-pass-self-review.md`           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-27_22-49_flaky-stack-test-root-cause-fix.md`                            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-28_09-23_pareto-plan-execution-brutal-self-review.md`                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-28_10-14_hardening-pass-brutal-self-review.md`                          | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-28_14-59_docs-health-update-old-docs-brutal-self-review.md`             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-28_15-29_docs-health-completion-visualtest-shadowing-discovery.md`      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-28_16-31_v1.3.0-release-cut-chaotic-but-shipped.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_01-55_snapshot-test-infrastructure-overhaul.md`                      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_02-20_post-push-reality-check.md`                                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_11-19_buildflow-failure-triage-three-fixes.md`                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_11-35_buildflow-repeat-failure-triage-round-2.md`                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_15-53_keyboard-navigation-improvements.md`                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_16-22_keyboard-navigation-fixes-and-gaps.md`                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_16-59_keyboard-navigation-round3-fill-bug-and-gaps.md`               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_21-49_fix-visualtest-compile-and-golangci-linter-regression.md`      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_22-19_buildflow-autofix-regression-loops-root-caused.md`             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-07-30_22-30_permanent-fixes-for-autofix-regression-loops.md`               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-09_07-58_docs-health-close-out-execution-self-review.md`                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-09_08-36_docs-health-close-out-execution-and-self-review.md`            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_02-50_web-components-rejection-and-container-query-leveraging.md`    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_03-52_followup-drift-guards-doc-relocation-skill-update.md`          | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_04-22_multi-module-workspace-split.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_05-00_modularization-cleanup-and-self-review.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_06-31_v2-default-flip-and-modularization-completion.md`              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_06-32_go-datastar-static-integration.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_06-46_datastar-sub-module-extraction.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_07-05_v2-doc-cleanup-and-module-count-fix.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_07-06_htmx-sub-module-extraction-and-7-module-infra-update.md`       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-10_19-28_animated-icons-heroicons-animated-inspired.md`                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-11_03-57_animated-icons-rebuild-fix.md`                                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-11_06-57_animated-icons-polish-ship.md`                                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-11_07-24_full-source-verification-of-all-96-icon-mappings.md`           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-14_20-25_CI-RECOVERY-SESSION.md`                                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-15_00-50_CI-GREEN-DEPLOY-LAST-RED-AND-LESSONS.md`                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-15_01-33_DEPLOY-DOUBLE-GREEN-BRUTAL-SELF-REVIEW.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-16_03-44_errorpage-family-drift-review.md`                              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-17_14-26_demo-page-superb-audit-and-overhaul.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-17_15-14_demo-overhaul-verification-complete.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-21_08-36_sse-audit-and-sse-stack-leverage.md`                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-21_09-07_sse-followup-retrymode-and-harvest.md`                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-21_22-33_v1.9.0-release-cut.md`                                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-22_03-42_dnsblockd-block-page-learnings-shipped.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-22_04-23_v1.10.0-release-cut-daemon-race-recovery.md`                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-22_04-51_v1.10.0-github-release-ci-recovery.md`                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-31_16-12_branch-consolidation-master-ci-recovery.md`                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-08-31_17-22_fix-all-convention-repairs-and-ci-hardening.md`                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-02_20-00_todo-clearout-ci-hardening-datastar-lenses-go-toolchain.md`    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-03_15-49_go-datastar-static-v0.4.0-bump-and-daemon-regression-watch.md` | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-03_21-03_v1.12.0-release-cut-and-release-sh-sigpipe-recovery.md`        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-03_23-43_vuln-gates-release-fixture-test-master-green-again.md`         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-05_01-11_wire-pareto-execution-session.md`                              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-05_02-00_wire-release-v1130-cut-session.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-05_17-19_merge-open-branches-master-green.md`                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-07_16-20_dual-transport-forms-wire-content-type.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-07_18-01_forms-transport-phase0-trust-pinned.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-07_20-50_forms-transport-t01-t27-complete.md`                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-07_21-55_templ-components-cv-adoption-analysis.md`                      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_04-30_demo-visual-audit-7-pages-light-dark.md`                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_04-33_forms-pack-browser-e2e-goldens-and-ops-triage.md`              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_05-31_demo-fix-session.md`                                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_07-35_demo-fix-followthrough.md`                                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_15-31_todo-hardening-batch.md`                                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_15-35_plan-execution-cv-adoption.md`                                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_16-46_full-execution-complete-cv-adoption.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_16-47_followthrough-browser-proofs-route-goldens.md`                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_17-17_prs-merged-execution-wrap.md`                                  | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_21-52_release-1150-cut-cv-e2e-unblock.md`                            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_22-07_buildflow-failure-triage-green-sweep.md`                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-08_22-36_v1.15.1-release-cut-and-daemon-relapse.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-09_01-50_kanban-board-dual-transport.md`                                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-09_02-18_pareto-plan-and-m1-touch-fix.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-09_04-20_kanban-follow-through-execution.md`                            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-09_16-19_post-release-a11y-and-demo-e2e.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-09_23-21_n-plan-execution-complete-e2e-fixes-and-tail-tasks.md`         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-10_01-59_docs-health-annotate-archive-pass.md`                          | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-13_10-17_css-artifact-audit-and-dead-artifact-cleanup.md`               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-13_10-38_ghost-system-elimination-and-vision-integration.md`            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-13_12-12_docs-health-second-pass-annotations-resurrection-fix.md`       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-13_12-33_astro-to-templ-conversion-status.md`                           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-13_14-31_pareto-execution-wave2-v1.17.0-shipped.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-13_18-38_website-hardening-status.md`                                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-14_03-29_pareto-wave3-m14-m20-nine-bugs-shipped.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-14_05-56_pareto-wave4-followups-m21-m22-twelve-bugs-f109-gate.md`       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-14_09-33_mr-status-reconcile-branch-cleanup-ci-reds.md`                 | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-14_17-13_adr-0040-component-module-boundary-decision.md`                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-16_19-12_kanban-vibe-lessons-execution-status.md`                       | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_05-54_kanban-followthrough-hardening-status.md`                      | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_06-43_pareto-planning-session-status.md`                             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_09-28_todo-sweep-hardening-session.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_14-17_v1.18.0-release-session.md`                                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_15-27_errorpage-redesign-and-error-model-status.md`                  | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_18-09_website-session-self-review.md`                                | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_18-10_kanban-optimistic-pending-register.md`                         | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_18-49_demo-recipes-overhaul-status.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_19-47_kanban-test-lenses-png-evidence-status.md`                     | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_19-59_errorpage-pareto-plan-execution-kickoff-status.md`             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-17_20-17_docs-health-full-audit-status.md`                              | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-18_05-39_execution-session-status.md`                                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-18_07-06_errorpage-plan-waves-m01-m06-status.md`                        | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-18_09-36_errorpage-plan-m07-m25-close-out-status.md`                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-18_09-51_errorpage-plan-session-close-honest-retro-status.md`           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-19_15-10_website-sales-page-session-status.md`                          | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-19_17-27_svg-single-source-self-integration-status.md`                  | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-19_18-46_pareto-plan-execution-status.md`                               | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-19_19-20_plan-ship-blocked-by-concurrent-red-tip-status.md`             | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-19_22-45_pareto-session5-t9-t19-completion-status.md`                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-20_00-16_pareto-completion-t18-t24-visual-witness-status.md`            | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-22_22-38_tc-mirror-guard-art-dupl-baseline-status.md`                   | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-22_22-57_pareto-plan-ritual-push-two-counter-fixes-status.md`           | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-23_00-35_dedup-pass-artdupl-gate-binary-pin-status.md`                  | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-23_13-25_pareto-plan-full-execution-status.md`                          | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |
| `docs/status/archived/2026-09-23_15-18_dedup-t2-full-triage-chart-shell-status.md`                    | ARCHIVE | every numbered/table row resolved inline; check-rows gate green 2026-11-03 |

## Post-sweep defects & corrections (2026-11-04, self-review pass)

1. **CHANGELOG heading split-brain (fixed).** The sweep's CHANGELOG entry created a second
   `### Changed` in `[Unreleased]`; HEAD already carried duplicate `### Added` headings
   (daemon-era accumulation). All `[Unreleased]` headings consolidated to single
   Added/Changed/Fixed; `check-changelog-guard.sh` green.
2. **9 planning files were archived prematurely (fixed by un-archiving).** check-rows only
   flags MIXED tables — wholly-untouched tables pass silently, and the row-extraction pipeline
   keyed on the gate's detail output, so their item tables were never inventoried. The 9 files
   (2026-05-18, 2026-07-05_21-27, 2026-07-07_18-10, 2026-07-12_14-45, 2026-07-12_19-00,
   2026-07-21_03-54, 2026-09-13_11-15, 2026-09-19_15-36, 2026-09-19_17-48 — ~168 bare rows)
   moved back to `docs/planning/` and stay open for the next corpus pass. **The manifest above
   lists them as ARCHIVED; treat those 9 rows as superseded by this correction.**
3. **Stale prose path citations (fixed).** 64 `docs/status|planning/2026-0…` shorthand citations
   in TODO_LIST/ROADMAP rewritten to their `archived/` locations (kept-open files excluded).
   Markdown-link audit: 0 broken links in live docs; archived-internal links are quoted-as-written
   by policy and excluded from lychee (`lychee.toml exclude_path`).
4. **Gate blind spots recorded for the tooling backlog:** check-rows should flag wholly-unstruck
   tables in files being archived; the moved-file gate should require zero bare data rows, not
   merely ≥1 strikethrough (`docs/status/2026-11-04_03-35` self-review §e).
