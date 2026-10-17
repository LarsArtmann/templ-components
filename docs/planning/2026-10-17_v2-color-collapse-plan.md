# v2 color collapse — ADR-ready batch plan (TODO #431)

Date: 2026-10-17 · Gate: **v2.0 only** (ADR-0039 timing, ADR-0044 T1).
Companion evidence: `2026-10-17_theme-inline-probe.md` (mechanism proven),
`utils/tone` Phase 0 (shipped, M07), the M16 Position alias precedent
(shipped), `feedback/tone_adapters.go` byte-equivalence tests (shipped).

## Preconditions (all must be green before Batch 1)

| #  | Condition                                                                    | Status                                                                  |
| -- | ---------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| P1 | `@theme inline` mode-flipping probe verdict GO                               | **Done** — computed-style-identical incl. color-mix path (pinned 4.3.3) |
| P2 | `utils/tone` Phase 0 shipped: Tone/Hue axes + adapters on every color enum   | **Done** (M07)                                                          |
| P3 | Warning hue decision (#428: amber recommended) settled by owner              | **G1 memo delivered 2026-10-15; awaiting decision**                     |
| P4 | tailwind-merge-go conflict handling for token classes unit-tested            | Open — small test at Batch 1 time                                       |
| P5 | v2 module-path/timing decision (ADR-0039) executed in the same release train | Owner gate                                                              |

## Batch order (each batch lands green independently; Batches 1–4 share ONE release)

### Batch 0 — decisions executed in the tone table (no component edits)

1. Apply the #428 hue decision: ONE warning hue in `tone.For`/`UtilitiesFor`
   (recommended amber); update the feedback-dialect anchors and their
   byte-equivalence tests; strike the two-hue note in
   `docs/recipes/scoped-theme-bridge.md` (regenerate via `go test ./utils/tone/ -update`).
2. `FeedbackError` → `FeedbackDanger` value rename (`"error"` → `"danger"`)
   in the canonical enum — v2 breaking, rides the same release; migration
   note in CHANGELOG (AlertType/ToastType aliases die here too, ADR-0022).

### Batch 1 — token emission (the class swap; the 669 `dark:` pairs die here)

- Define the role tokens in `templates/custom.css` + the ADR-0008 layer:
  `@theme inline { --color-accent: var(--tc-accent); … }` with `:root`/`.dark`
  literals EXACTLY equal to today's blue pairs (accent = blue-600/500 bg,
  600/400 text, per-slot from `tone.UtilitiesFor`), same shape for
  success/warning/danger/info/neutral roles.
- Swap component class strings mechanically, per package, tone-adapters
  first (they already name the mapping): `bg-blue-600 dark:bg-blue-500` →
  `bg-tc-solid`-style role classes; gray neutrals → neutral role tokens.
- What does NOT change: structural `dark:` variants (non-color), the literal
  legacy `left/right` position classes (unrelated), charts' categorical
  palette (Hue axis stays palette classes until demand says otherwise).
- Verification per package: HTML goldens re-baseline ONCE at the END of the
  batch (see Batch 4); computed-style equivalence spot-probes via the
  visualtest probe harness pattern for one component per package.

### Batch 2 — lookup maps become tone-table views

- `feedbackStyleSet`/`feedbackIconName`, BadgeType/KanbanTone/StatTone/
  StatusTone tables: replace per-component literal maps with views over
  `tone.For(t)` / `tone.UtilitiesFor(h)` (the M07 adapters already prove the
  equivalence byte-for-byte — this is deleting the duplicated literals).
- Keep the enum TYPES (consumer API); only their backing tables collapse.

### Batch 3 — enum collapse (alias-then-remove, M16 precedent)

- Remove the deprecated aliases: `AlertType`/`ToastType` (ADR-0022),
  physical `PositionLeft/Right` values + per-component position aliases,
  legacy `DropdownPositionLeft/Right`, `tc-primary-*` legacy named tokens.
- Rename/misc v2 enum items already ADR'd elsewhere (Scrollback rename,
  TrendWarn value, ButtonType split) — reference, not owned here.

### Batch 4 — the ONE re-baseline event (second and last allowed)

1. `go test ./... -update` for HTML goldens repo-wide; `nix run .#visual -- -update`.
2. Recompile the 3 tracked CSS artifacts; regenerate the class inventory.
3. Guard updates: `TestDarkModeSemanticColors` semantics shift from
   "every semantic class carries `dark:`" to "every role token has a `.dark`
   override + components emit role classes" (the scanner rule changes WITH
   the swap — write it in the same commit); counts (enums drop as aliases
   vanish — TestDocsCountDrift + FEATURES/README/website in lockstep).
4. CHANGELOG v2.0 entry with the consumer migration table
   (class renames + enum removals + theme-layer import requirement).

## Risks

| Risk                                              | Mitigation                                                                                 |
| ------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| Token class invisible to Tailwind scanner         | Classes stay complete static literals in templ (same contract as `Pair`)                   |
| tailwind-merge mis-merges role vs palette classes | P4 unit test before Batch 1                                                                |
| Consumer CSS bridges (attribute selectors) break  | Pre-announced in CHANGELOG + the theme-bridge recipe already documents the v2 escape hatch |
| Guard blindness during the swap                   | Swap is per-package with tests before daemon windows; re-baseline ONLY at Batch 4          |
| Scope creep into module-path v2 work              | Hard non-goal here; ADR-0039 owns that track                                               |

## Non-goals

Module path change, non-color v2 API breaks, new components, chart palette
tokenization, anything in G5 not listed above.
