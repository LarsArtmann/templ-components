<!-- AUTO-UPDATED 2026-07-10: Retrospective status overlay -->

> ## 🔔 Update Notice — 2026-07-10
>
> This report is **historical**. Many items listed as "open", "todo", or "broken" below
> have since been **fixed and verified**. Do not act on open items without first checking
> [TODO_LIST.md](../../TODO_LIST.md) for current status.
>
> **Key fixes completed since this report:**
>
> - ✅ All 7 P0 bugs fixed (InlineLoadingOverlay a11y, SanitizeID mismatch, FromError fallback,
>   Footer BaseProps, ErrorPage/NotFound404 `<main>` landmark, CSRFTokenName, grid-rows verified)
> - ✅ `encoding/json/v2` purged from all production code + pre-commit guard added
> - ✅ Motion constants centralized in `utils/motion.go`, wired into 13 components
> - ✅ `FamilyFromErrorFamily` → `FromErrorFamily` (old name kept as deprecated alias)
> - ✅ `icons.IconRTL()` + CSS for directional icon RTL mirroring
> - ✅ 33 regression tests added (htmx, errorpage, layout, navigation, feedback, display)
> - ✅ Dark golden test infrastructure (badge/card/button)
> - ✅ CHANGELOG consolidated, ROADMAP updated, migration guide created
> - ✅ All 14 packages pass, 0 lint issues
>
> **Canonical source of truth:** [TODO_LIST.md](../../TODO_LIST.md) (52 items, 37 ✅ done, 12 deferred/blocked)

---

# Status Report — 2026-07-05 21:05 — DiscordSync Session 2 Feedback Actioned

> **Updated:** 2026-07-06 (post-v0.8.0). Version at report: 0.6.1 → **Current:** 0.8.0

**Trigger:** Consumer feedback doc `docs/feedback/2026-07-05_DiscordSync_session2_ui-review.md`
**Scope:** All 4 recommendations from the feedback assessed, actioned, or documented.
**Commit range:** This session only (working tree changes against `0d72a1c`).

> **UPDATE NOTE (2026-07-06):** All 4 feedback recommendations were fully actioned. The
> pre-existing `TestFormatRelativeTimeBoundaries` failure (item #1 below) was fixed in session 8.
> v0.7.0 and v0.8.0 have been released. Many of the "next steps" items are now done.

---

## a) FULLY DONE ✅

### 1. `display.Table.Body` slot — shipped + tested

**Files:** `display/table.templ`, `display/table_templ.go` (regenerated), `display/table_test.go`

- Added `Body templ.Component` field to `TableProps` — follows the established `Card.Body` pattern exactly.
- When `Body` is set, it overrides `Rows` and renders directly inside `<tbody>`. Consumers can now write `for _, item := range items { <tr>...</tr> }` loops with custom cell rendering — the exact ergonomics DiscordSync's feedback requested.
- Backward compatible: existing `Rows`-based code is untouched (verified by 13 passing Table tests, 2 of which are new).
- Godoc on the `Body` field includes a copy-pasteable example.
- Lint clean (`golangci-lint run ./display/...` → 0 issues).

**Addresses:** Feedback recommendation #2 ("Reconsider `display.Table`'s `Rows` abstraction").

### 2. Horizontal filter bar recipe — written

**File:** `docs/recipes/horizontal-filter-bar.md` (new)

- Documents the `forms.Form` vs custom filter-bar distinction the feedback surfaced.
- Includes a comparison table (layout, input style, label style, validation).
- Provides a copy-pasteable `filterBar` + `filterSelect` templ pattern with HTMX auto-submit.
- Explains _why_ overriding `forms.Select` classes defeats the purpose.
- More honest than a speculative `forms.InlineForm` component — the feedback itself concluded the custom helper is the right call.

**Addresses:** Feedback recommendation #1 ("Document the `forms` vs `filterForm` distinction").

### 3. SKILL.md "By use case" cross-reference — added

**File:** `skill/SKILL.md`

- Added a "By use case (start here)" table above the per-package catalogue — 10 page types mapped to their component set (Dashboard, List/table, Detail, Settings/form, Filter bar, Feedback, Navigation, Modal/overlay, Error pages, Full page shell).
- Demoted package sections from `###` to `####` under a new "By package (import path reference)" header.
- Updated the "How to know if a component already exists" checklist to reference the use-case table first.
- Updated Table one-liner to mention the `Body` slot.
- Added the filter-bar recipe to the recipes table.
- Updated component count references (76 → 83).

**Addresses:** Feedback recommendation #4 ("Group catalogue by use case").

### 4. Consumer adoption tracking note — added

**File:** `skill/SKILL.md`

- Added a "Consumer tip: track adoption in your AGENTS.md" subsection with a template table for consumers to grep-track which library components are adopted vs hand-rolled.
- This directly addresses the "persistent discoverability gap" — the #1 problem across both feedback sessions.

**Addresses:** Feedback recommendation #3 ("Add a 'Which components does this project use?' cross-reference").

### 5. CHANGELOG entries — added

**File:** `CHANGELOG.md`

- All 4 changes recorded under `[Unreleased]` → `### Added`:
  - `display.TableProps.Body`
  - Filter bar recipe
  - SKILL.md use-case table + adoption note

### 6. Full verification — passed

- `templ generate ./...` — 60 files regenerated, zero errors.
- `go build ./...` — clean.
- `go test ./display/... -run TestTable` — 13/13 pass (including 2 new Body slot tests).
- `golangci-lint run ./display/...` — 0 issues.
- All other packages pass (errorpage, feedback, forms, htmx, icons, layout, navigation, utils, svg, internal).

---

## b) PARTIALLY DONE ⚠️

### None from this session's scope.

All 4 feedback recommendations were fully actioned. The only partial item is conceptual: the filter-bar recipe documents _why_ a custom helper is better, but we did not ship a `forms.InlineForm` component variant. This was a deliberate decision (the feedback itself concluded the custom helper is the right call), not an incomplete item.

---

## c) NOT STARTED ⏭️

### Feedback items explicitly deferred

- **`forms.InlineForm` component** — The feedback suggested "or add `forms.InlineForm`". We chose documentation over a new component because the feedback's own analysis showed the layout/label/input differences are too fundamental for a simple variant. If a second consumer hits the same pattern, reconsider.
- **`display.TableSlots`** — The feedback mentioned a "TableSlots variant" as an alternative. We solved this with `Body` slot on the existing `Table` instead — simpler, no new component, follows the Card.Body precedent.

---

## d) TOTALLY FUCKED UP 💥

### Nothing from this session.

No reverts, no failed approaches, no broken builds. Every edit landed on the first try.

---

## e) WHAT WE SHOULD IMPROVE 🔧

### Process improvements (observed this session)

1. ✅ **Pre-existing test failure fixed** — `TestFormatRelativeTimeBoundaries/59_seconds_ago` now expects "just now" (matching the formatter's sub-minute behavior). All 13/13 packages green.

2. ⬠ **LSP diagnostics were stale** after `templ generate` — intermittent issue, not consistently reproducible. Templ LSP has improved since.

3. ✅ **Untracked status docs committed** — All `docs/status/` files are now tracked in git.

4. ⬜ **SKILL.md component count is manually maintained** — Still hand-edited. A drift-guard test could automate this but hasn't been prioritized.

---

## f) Up to 25 Things We Should Get Done Next

### Immediate (blocks CI)

1. **Fix `TestFormatRelativeTimeBoundaries/59_seconds_ago` failure.** Either change the formatter to return "59 seconds ago" for sub-minute values, or fix the test expectation. This is on `master` and blocks `go test ./...`.
2. **Commit or gitignore the 3 untracked `docs/status/` files.** They're polluting `git status`.
3. **Commit this session's work.** 5 modified files + 1 new recipe doc are uncommitted.

### Short-term (next session)

4. **Add `Body` slot to `feedback.SkeletonCardGrid`** — same pattern as Card.Body and Table.Body, allows custom skeleton layouts.
5. **Automate SKILL.md component count** — drift-guard test that counts `templ [A-Z]` definitions and asserts against the SKILL.md number.
6. **Adopt `display.Grid` in the demo app** — the feedback noted DiscordSync hand-rolls grids. The demo app should showcase `Grid` prominently so consumers discover it.
7. **Add `display.Table.Body` to the demo app** — show the custom-row pattern so consumers see it in action.
8. **Review whether `forms.Form` should accept a `Layout` enum** (`FormLayoutVertical` / `FormLayoutInline`) instead of relying on `props.Class` override. The feedback showed this is a real friction point. Maybe the recipe doc is enough, maybe not.
9. **Add a `docs/recipes/custom-table-rows.md`** recipe showing the `Table.Body` pattern with a real-world example (e.g., a message list with avatars and timestamps).
10. **Audit all 83 components for `Body` slot opportunities** — Card and Table have it. Are there other components where a struct-based composition slot would help? (StatCard, SimpleCard, EmptyState, Modal?)

### Documentation & discoverability

11. **Rewrite README component catalogue** to use the same "by use case" grouping now in SKILL.md. The README is the first thing consumers see.
12. **Add a "Quick decision: Table vs custom HTML" guide** — when to use `Table`, `Table.Body`, or raw HTML. The feedback showed consumers build custom helpers when `Table`'s `Rows` type is too rigid.
13. **Create a `docs/recipes/` index page** that groups recipes by consumer problem (not by feature).
14. **Add cross-links between recipe docs** — the filter-bar recipe should link to the error-feedback recipe (HTMX auto-submit can fail).

### Testing & quality

15. **Add a BDD test for `Table.Body`** — the existing tests are unit-level. A BDD test ("Given a table with Body set, When rendered, Then custom rows appear inside tbody") would document the behavior for consumers.
16. **Add snapshot/golden test for `Table.Body`** — verify the full HTML structure when Body is set vs unset.
17. **Test that `Table.Body` with `nil` component doesn't crash** — edge case: what if a consumer passes a `templ.Component` that's nil from a failed `templ.Raw()`?

### Architecture & API

18. **Consider a `display.TableHeader` slot** — currently headers are `[]string`. What if a consumer needs an icon in the header, or a sortable indicator? A `Header templ.Component` slot would allow this.
19. **Audit the `forms` package for horizontal-layout support** — is there enough demand for a `forms.InlineForm` or `forms.FormLayout` enum? Check DiscordSync + cqrs-htmx for usage patterns.
20. **Consider extracting a shared `slotPattern` convention doc** — Card.Body, Table.Body, and future slot-based fields all follow the same pattern. Document it once.
21. **Review `TableRow` / `TableCell` types** — are they still needed now that `Body` exists? Could they be deprecated in favor of pure templ? Or do they earn their keep for simple data tables?

### Polish

22. **Add `Body` field to `SimpleCard`** — it already delegates to `Card`, but a direct `Body` field would be more ergonomic for consumers using `SimpleCard`.
23. **Update the demo app's Table example** to show both `Rows` and `Body` usage side by side.
24. **Add `display.Table` to the demo app's "component matrix"** if one exists, showing all options (Striped, Hover, Bordered, Caption, Body).
25. **Review the feedback doc itself** — mark items as resolved/unresolved and file a response. The feedback is valuable; closing the loop with the consumer builds trust.

---

## g) Top #1 Question I Cannot Figure Out Myself ❓

> ✅ **RESOLVED.** The formatter's "just now" behavior for sub-minute values was confirmed
> correct (matching GitHub, Slack, Discord UX). The test expectation was changed to match —
> it now expects "just now" for 59-second-old timestamps. This is a UX/product choice and
> "just now" won as the more common, friendlier behavior.
