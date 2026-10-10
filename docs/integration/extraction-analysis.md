# Component Extraction Analysis: Four Consumer Projects

**Date:** 2026-10-10 · **Status:** all library-side builds shipped
(2026-10-10); consumer adoption waves X1/X2 pending (**source of truth:**
`docs/planning/2026-10-10_06-06_extraction-backlog-pareto-master-plan.md`, TODO_LIST #388–400)

This document is the consumer-side runbook produced by the four-project
extraction analysis (DiscordSync, mr-sync, nsfw-classifier, dnsblockd). Every
load-bearing claim below was verified at source (E1–E7 in the plan); the
session report lives at
`docs/status/2026-10-10_05-51_extraction-analysis-four-projects.md`.

Two architectural decisions came out of this analysis and shape everything
below: [ADR-0044](../adr/0044-semantic-accent-tokens.md) (semantic accent
tokens — kills the consumer CSS color bridges) and
[ADR-0045](../adr/0045-css-class-delivery.md) (class inventory shipping —
kills the consumer CSS scan machinery).

## Per-project adoption map

Hand-rolled consumer code that a library component now (or soon) replaces.

### DiscordSync (`~/projects/DiscordSync`)

| Hand-rolled (source)                                                                                      | Library replacement                                           | Status                             |
| --------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------- | ---------------------------------- |
| `commandLine`/`copyBtn` (`cmd/mr-sync/...` shape mirrored here via `copyIDRow`, `filters.templ:363–383`)  | `display.CodeBlock` (CompactID variant)                       | ✅ shipped                         |
| `filterForm`/`filterSelect`/`selectOptionsWithCurrent`/`selectOptionsFromItems` (`filters.templ:391–434`) | `forms.FilterBar` + `forms.Select`                            | ✅ shipped — adoption pending (X2) |
| `listTableWithHeader` sticky thead                                                                        | `Table.StickyHeader` flag                                     | ✅ shipped — adoption pending      |
| `imageThumbnail` blur-up placeholders                                                                     | `ImageProps.Placeholder` (zero-JS blur-up)                    | ✅ shipped — adoption pending      |
| Hand-rolled charts                                                                                        | `display.LineChart`/`AreaChart`/`PieChart`                    | X2 (post-ADR-0044)                 |
| Vendor color bridge (`input.css:184` — comment admits overriding library blue → brand purple)             | `--color-accent-*` tokens (ADR-0044)                          | decided; T1                        |
| Vendor-dir `@source` (`input.css:5`) — **inert under Tailwind v4.3: `vendor/` is gitignored**             | `templates/templ-components-classes.txt` inventory (ADR-0045) | decided; T2                        |

**Pin:** templ-components v1.19.4 (stale).

### mr-sync (`/home/lars/projects/mr-sync`)

| Hand-rolled (source)                                                  | Library replacement                   | Status                                                                                                                                                    |
| --------------------------------------------------------------------- | ------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `commandLine`/`copyBtn` (`cmd/mr-sync/dashboard_components.templ:11`) | `display.CodeBlock`                   | ✅ shipped                                                                                                                                                |
| `cardValueBlock`/`statCard`                                           | `display.StatCard`                    | migration ⫱ owner                                                                                                                                         |
| `fetchErrorBanner`                                                    | `feedback.Alert` / `errorpage` family | migration ⫱ owner                                                                                                                                         |
| Hand-rolled semantic CSS dashboard (no Tailwind/HTMX)                 | full library adoption                 | **owner-gated** — recommendation: migrate (small dashboard, near-total duplication); final call is the owner's, like PapDashboard in adoption survey #156 |

**Pin:** not a consumer yet (templ + custom CSS only).

### nsfw-classifier (`/home/lars/projects/nsfw-classifier`)

| Hand-rolled (source)                                                                | Library replacement                       | Status                        |
| ----------------------------------------------------------------------------------- | ----------------------------------------- | ----------------------------- |
| `historyChip` zero-JS filter chips (`internal/server/views/history_page.templ:105`) | `forms.FilterChips`                       | ✅ shipped — adoption pending |
| Preset/duration buttons                                                             | `forms.SegmentedControl`                  | ✅ shipped — adoption pending |
| Evidence proportion bars                                                            | `display.SegmentBar`                      | ✅ shipped — adoption pending |
| Image viewer                                                                        | `display.Lightbox`                        | ✅ shipped — adoption pending |
| `scorePercent`/`formatDuration(ms)` (`views/types.go:176–190`)                      | `format.Percent`/`format.CompactDuration` | ✅ shipped                    |
| `third_party/css-scan` mirror + `scripts/sync-css-scan-sources.sh`                  | class inventory (ADR-0045)                | decided; T2                   |
| SurfaceNav                                                                          | library `Nav` (determinism check first)   | X2                            |

**Pin:** v1.21.0.

### dnsblockd (`/home/lars/projects/dnsblockd`)

| Hand-rolled (source)                                                                                                   | Library replacement                                               | Status                        |
| ---------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- | ----------------------------- |
| Honesty-ladder state duplication — 9 `art-dupl:accept` markers across 4 view files (e.g. `devices_page.templ:177–183`) | `display.DataState`                                               | ✅ shipped — adoption pending |
| Staleness/live pills                                                                                                   | `display.StatusDot`/LivePill                                      | ✅ shipped — adoption pending |
| 4 local status-badge mappers                                                                                           | `StatusBadgeWith` injectable mapper                               | ✅ shipped — adoption pending |
| Stale-session/retry meta tags (`dashboard_login.templ:63`, `allow.templ:21`)                                           | `layout.MetaRefresh`                                              | ✅ shipped                    |
| `scripts/gen-library-classes.sh` + `library-classes.txt`                                                               | class inventory (ADR-0045 — dnsblockd's mechanism, library-owned) | decided; T2                   |
| Hand-rolled dashboard sidebar                                                                                          | `navigation.SidebarNav`                                           | X1                            |

**Pin:** v1.19.4 (stale).

## Cross-cutting findings

### Theming (ADR-0044)

All three Tailwind consumers ship CSS override bridges because the library
hardcodes `blue-*` (212 sites, shades 50–900) as its accent. DiscordSync's
bridge (`input.css:184–190`) is the verbatim admission. The decision: a
dedicated `--color-accent-*` family (literal-blue defaults, independence by
construction) shipped in `templates/custom.css`; components emit `accent-*`
after T1; the ADR-0008 semantic alias layer is completed for feedback colors.
Consumers then delete their bridges and override one `@theme` block.

### CSS class delivery (ADR-0045)

Three mechanisms existed; the analysis found the default-looking one
(vendor-dir `@source`) is **inert under Tailwind v4.3** whenever `vendor/` is
gitignored (DiscordSync's is — `.gitignore:90`; v4.3 never scans gitignored
paths, per dnsblockd's investigation). The decision: the library ships a
release-generated class inventory file; consumers copy it with one line and
`@source` it.

### Library-side waves (from the same analysis)

- **Shipped (2026-10-10):** `utils/format` (six helpers), `display.CodeBlock`,
  `layout.MetaRefresh`, and the full C-wave — FilterBar, DataState, StatusDot,
  LivePill, FilterChips, SegmentedControl, SegmentBar, Lightbox,
  `Table.StickyHeader`, `StatusBadgeWith`, `Image.Placeholder` blur-up.
- **Pending:** T1/T2 implementations (version-gated per the ADRs), consumer
  migrations X1/X2 (each repo runs its own verify ritual).

## Using this document

When picking up a consumer migration (X1/X2): read the project's row above,
bump its pin, adopt top-down (helpers → components → bridges last, after
T1/T2), and delete the hand-rolled twin in the same commit — the half-state
where both exist is how duplication regrows.

### T2 CSS-delivery migrations (ADR-0045 — inventory shipped 2026-10-10)

The library now ships `templates/templ-components-classes.txt` (generated by
`scripts/gen-class-inventory.sh`, release-refreshed, freshness-guarded by
`utils.TestClassInventoryFreshness`). Per-consumer migration, one commit each:

- **nsfw-classifier:** delete `third_party/css-scan/` +
  `scripts/sync-css-scan-sources.sh`; copy the inventory next to `app.css`;
  `@source "./templ-components-classes.txt"`. Trade-off: the corpus carries
  ALL library classes (their mirror was curated) — accepted per ADR-0045
  unless the CSS size delta proves to matter.
- **dnsblockd:** delete `scripts/gen-library-classes.sh` (their mechanism,
  now library-owned); either keep their tracked `library-classes.txt` name
  (replace content with the shipped inventory) or switch the `@source` line
  to the new filename.
- **DiscordSync:** replace the inert vendor `@source` (`input.css:5` — their
  `vendor/` is gitignored, so library utilities ride a stale compiled
  artifact); copy the inventory + rescan. Their color bridge
  (`input.css:184`) goes separately, post-T1.
