# Docs-Health Status Index — `docs/status/`

**Updated:** 2026-11-04 | **Maintainer ritual:** docs-health AUDIT (see the 2026-11-04 corpus
close-out sweep + self-review reports for method and gates).

This directory holds point-in-time session status reports. Historical value, not living docs.

## Layout

| Path                                                                  | Contents                                                                                                                                                      |
| --------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `2026-10-*.md`, `2026-11-*.md`                                        | Open-era reports (Oct/Nov 2026) — not yet swept; newest first                                                                                                 |
| 5 named 2026-09 files                                                 | Kept-open reports with genuinely-open rows (see below)                                                                                                        |
| `2026-11-04_03-27_docs-health-corpus-close-out-sweep.md`              | 2026-0* corpus close-out sweep report + 204-file archive manifest + postmortem                                                                                |
| `2026-11-04_03-35_docs-health-corpus-sweep-status-and-self-review.md` | Sweep self-review: done/partial/gaps, defect list, next-50                                                                                                    |
| `archived/` (204 files)                                               | Fully-annotated reports — every numbered/table row carries an inline verdict. Closed history; links inside are quoted as written (lychee-excluded by policy). |

## Kept-open reports (do not archive until their rows resolve)

- `2026-09-17_17-28_sidebar-light-mode-session-status.md` — dedicated `/navigation` demo route + goldens ask; nav-sidebar light eyeball; adaptive-brand guard idea.
- `2026-09-17_18-48_wire-action-design-review-session-status.md` — debounce-without-event htmx invariant test + keep-or-fix decision.
- `2026-09-17_20-00_card-sharp-default-squircle-opt-in-status.md` — squircle browser-proof gap (b1–b6) + follow-ups, two owner gates.
- `2026-09-18_07-08_kanban-followthrough-verification-hardening-status.md` — LoadingButton flake fix beyond diagnosis.
- `2026-09-23_15-56_prerender-drift-flake-fix.md` — prerender routes-table single-sourcing, calendar month pin, flake-stress mode, RelativeTime pinning evaluation, prerender-safe marker convention.

## Convention

- New reports: `YYYY-MM-DD_HH-MM_name.md`, newest insight first, honest gaps section mandatory.
- Annotate, never rewrite: bare rows in struck tables are OPEN items by design; verdicts are
  `done — evidence` / `**Won't implement…**` / `**NOT-DO…**` / `**→ routed TODO_LIST #N**`.
- Fully-resolved reports move to `archived/` AFTER `annotate.go check-rows` passes AND a
  bare-data-row scan (wholly-unstruck tables are invisible to check-rows — the 2026-11-04 lesson).
- Forward-looking items are harvested into `TODO_LIST.md` with provenance; the report itself
  never becomes the backlog.
