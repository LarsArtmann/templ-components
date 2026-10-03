# Archive Manifest — 2026-10-03 docs-health sweep

Every archived report carries inline resolutions in the skill's marker format:
`~~<original item>~~ done at <evidence>` (variants: `Won't implement — <reason>`,
`NOT-DO — <reason>`). This manifest itself is an index, not a report.

One line per file archived by the 2026-10-03 docs-health pass (skill rule: bulk
archives carry a manifest; single-file archives are exempt). Classification per
the docs-health ANNOTATE/ARCHIVE table: every numbered item resolved inline
(strikethrough) before the move.

| File (relative to repo root)                                               | Classification | Deciding reason                                                                                                       |
| -------------------------------------------------------------------------- | -------------- | --------------------------------------------------------------------------------------------------------------------- |
| docs/status/archived/2026-06-08_00-51_todo-list-execution-comprehensive.md | ARCHIVE        | Last open item (consolidate inline JS into shared init) resolved by the ADR-0005 singleton-guard adoption             |
| docs/status/archived/2026-07-05_20-39_modularization-status-report.md      | ARCHIVE        | All 11 items struck (earlier pass); this pass verified zero unstruck items remain                                     |
| docs/status/archived/2026-09-04_22-15_transport-wiring-sdk-session.md      | ARCHIVE        | Last open item (wire BDD spec) closed Won't-implement — optional per repo norms; wire keeps table tests               |
| docs/planning/archived/2026-07-06_02-36_SUPERB-V0.9.0-HARDENING-SPRINT.md  | ARCHIVE        | Last two items (WriteNotFound404 handler + test) verified shipped (`errorpage/handler.go:190`, `handler_test.go:712`) |
| docs/feedback/archived/2026-07-10_cqrs-htmx-consumer-feedback.md           | ARCHIVE        | Last open ask (Card header overrides) shipped — Card has `TitleClass`/`HeaderClass`/`TitleComponent` + `Header` slot  |

Also annotated this pass (files moved earlier without inline markers, now
gate-compliant): `docs/feedback/archived/2026-07-05_overview-consumer-feedback.md`
(6 pain points struck inline, skeleton-demo residue verified shipped) and
`docs/feedback/archived/2026-07-12_table-in-card-double-border.md` (fixed by
`Table.Flush`, struck inline).
