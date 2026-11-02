# Owner-Decision Packet v2 — refresh (2026-11-02)

**Companion to:** `2026-10-30_16-30_SUPERB-OWNER-DECISION-PACKET.md` (v1 — full options,
evidence, and blast radius per gate; NOT superseded, this memo only refreshes evidence).
**Method:** every gate re-checked against the 2026-11-02 tree; deltas only. Where
evidence is unchanged, the v1 RECOMMENDED DEFAULT stands and nothing needs re-reading.

**How to answer:** reply per row with "accept default" or a different option letter.
Tasks gated on these rows keep proceeding on the defaults (plan T14 mandate).

## Evidence deltas since 2026-10-30

| Gate                          | Δ Evidence by 2026-11-02                                                                                                                                                                                                                                                                          | Default                                          |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------ |
| Q1 browser baseline           | none — `docs/browser-support.md` engine matrix unchanged; audit posture holds                                                                                                                                                                                                                     | **A** (2024-widely-available, fallbacks default) |
| Q2 LoadingOverlay identity    | none — demo specimen unchanged                                                                                                                                                                                                                                                                    | **A** (static default, modal split out)          |
| #80/#150/#162 vision review   | **API key still absent** from the session env (checked 2026-11-02); `vision-review-goldens.sh` ready and waiting; rows stay risk-retirement-only                                                                                                                                                  | (a) run once at next release eyeball window      |
| #211 templates fate           | none — `TestCompiledCSSInventory` still pins exactly the 3 tracked artifacts; no new consumer signal                                                                                                                                                                                              | v1 recommendation (keep as starter artifacts)    |
| #212 annotation policy        | **precedent landed:** the 2026-11-02 docs-health ANNOTATE pass annotated 4 stale reports inline (2026-05-19 architecture `.md`, 2026-07-18 modularity HTML, two data-model HTMLs — commit `127160e2`), including tag-balance checks; the inline-note pattern works and scales to the older corpus | agent-side inline ANNOTATE sweeps, batched       |
| #233 flake.lock daemon nudge  | none                                                                                                                                                                                                                                                                                              | v1 recommendation                                |
| #234 recipes audit guarding   | none                                                                                                                                                                                                                                                                                              | v1 recommendation                                |
| #237 AppShell/Sidebar chrome  | none — SidebarNav remains theme-adaptive via `--tc-sidebar-*` tokens; no consumer feedback either way                                                                                                                                                                                             | v1 recommendation                                |
| #238/#239 Kanban failure UX   | none — ADR-0041 pending register unchanged; no e2e flakes observed in this session's runs                                                                                                                                                                                                         | v1 recommendation                                |
| #269 go-error-family PR       | **branch still prepared locally** (`~/projects/go-error-family`, red-green tests, probe S1–S5 verified); needs owner go-ahead to push + open upstream PR; go-error-family v0.12.0 shipped here without it (nil-receiver Unwrap was the only delta)                                                | push branch + open PR on owner go-ahead          |
| #329 SSE via Firebase rewrite | none — still buffered through `/demo/**`; LiveRegion still hangs on the canonical URL; HTMX flows unaffected                                                                                                                                                                                      | v1 recommendation                                |
| #358 GOEXPERIMENT=jsonv2      | none — still exported repo-wide; Go floor unchanged at 1.27                                                                                                                                                                                                                                       | v1 recommendation                                |
| #411 tool verdicts ratify     | **new evidence, same posture:** T17 evaluations logged in `docs/testing/tooling-evaluation.md` §8 (erraudit LSP viable-not-wired; `--max` = audit surface, 13 generic-return findings all ACCEPTED per the SKILL ruling); `nix run .#crossbrowser` stays findings-lane                            | ratify at next tooling review                    |

## Adjacent state the owner should know (no decision required)

- **#28 awesome-templ#12:** still OPEN, awaiting maintainer review (re-verified 2026-11-02).
- **Dependabot:** templ + go-error-family now bot-ignored in all modules
  (commit `69d3a8ea`); the nine templ-0.3.1100 PRs closed with policy comments; zero
  dependabot refs remain. If the chromedp parts of former #67 are wanted, dependabot
  re-proposes them without templ.
- **#335 templ window:** full execution plan-of-record exists
  (`2026-11-02_TEMPL-MIGRATION-WINDOW-PLAN.md`); gate still CLOSED; nothing to decide
  until `check-templ-gate.sh` opens.
- **T15 demand re-check:** FileDrop FLIP-ready (2 consumers); MultiSelect WATCH;
  TreeView HOLD (row 217 has the sites).

## Still the two asks that unblock the most

1. **Q1 (browser baseline)** — unblocks C19/C22/C23/C24 + #473/#476–#482 dispositions.
2. **#80/#150/#162 API key export** — one script run retires three rows of risk.
