# F109 demand-gate re-check — 2026-10-08 (competitor evidence pass)

**Gate:** TODO_LIST #217 (F109 verdict, 2026-09-14): MultiSelect, DateRangePicker,
FileDrop, Command palette, Toast positions, TreeView had ZERO demand evidence in
the #156 22-repo survey — do not build speculatively; re-run when a consumer
asks or an adoption survey lands.

**Re-run trigger (logged 2026-10-08):** shadcn-templ (fka templUI, ~1.75k stars)
ships a set overlapping the gated list. This memo re-runs the check with that
evidence and renders a flip/hold verdict per component.

**Evidence standard:** the gate exists to prevent speculative builds. A
competitor shipping a component proves CATEGORY demand, not OUR-consumer
demand — so a flip requires either (a) ecosystem demand plus a named
in-repo consumer that would use it (dogfooding counts — the website and demo
are the library's largest consumers), or (b) a direct consumer ask. Evidence
below is from primary sources (shadcn-templ llms.txt + docs, verified
2026-10-08; see `docs/comparison.md` for the full competitor profile).

## Per-component verdicts

| Component | Competitor evidence (2026-10-08) | Our-consumer evidence | Verdict |
| --- | --- | --- | --- |
| Command palette | shadcn-templ ships `Command` ("Fast, composable command menu", cmdk-based) — the flagship shadcn/ui component ecosystem-wide | The website's docs search is a hand-rolled combobox (`website/assets/js/search.js`); the demo MPA has no cross-page navigation. A CSP-safe palette would replace the hand-rolled search UI and give the demo a jump-to-component surface — two named internal consumers | **FLIP — approved for build** (full ladder). The strongest candidate: ecosystem-confirmed category demand + named internal consumers |
| DateRangePicker | shadcn-templ ships `Date Picker` as a PATTERN composed of Calendar + Popover + Button — not a primitive | We already ship `forms.Calendar`; the composition is expressible today as a docs recipe without a new primitive | **HOLD as primitive — route as docs recipe.** If a consumer asks for range selection beyond the recipe, re-open |
| MultiSelect | NOT in shadcn-templ's registry (checked llms.txt: no multi-select; Toggle Group's "multiple" is a segmented control) | None | **HOLD** — no ecosystem signal, no ask |
| TreeView | NOT in shadcn-templ's registry | None | **HOLD** |
| FileDrop | NOT in shadcn-templ's registry (no file-upload/dropzone) | None | **HOLD** |
| Toast positions | shadcn-templ ships Toast (positions are table stakes there), but this gate row is an ENHANCEMENT to our existing `feedback.Toast`, not a new component | No ask; current single-position default has zero complaints on record | **HOLD** — revisit only with a consumer ask |

## Decision

- **Command palette** flips to build-eligible. Executed per subtask 4.3: routed
  to TODO_LIST as an open build item with the full testing ladder attached
  (spec → templ → props/enum → goldens → a11y sweep → e2e → docs counts →
  `tc` mirror), then built as plan task T20.
- Everything else **holds**. The gate stays armed for MultiSelect, TreeView,
  FileDrop, Toast positions, and DateRangePicker-as-primitive.
- Provenance: docs/status/2026-10-08_19-57 §f6 (trigger), #156 survey
  (original verdict), docs/comparison.md (competitor profile).
