# Status Report: Table-in-Card Double Border Fix + CellPadding

**Date:** 2026-07-12 02:34
**Session scope:** Fix `docs/feedback/2026-07-12_table-in-card-double-border.md`, self-review, gap closure, push.
**Version:** 0.15.0 (unreleased changes pending — `[Unreleased]` in CHANGELOG populated)
**Commits this session:** 7 (615ba99 → dfb2997)
**Files touched:** 19 (+657 / -184 lines)

---

## a) FULLY DONE

### Core Feature: `TableProps.Flush` (double-border fix)

- **`Flush bool`** field added to `TableProps` — suppresses the wrapper div's `border` + `rounded-lg` while keeping `overflow-x-auto`. Consumer usage: `Table(TableProps{Flush: true})` inside `Card(CardPaddingNone)`.
- `tableWrapperClass(flush bool)` helper — returns `"overflow-x-auto"` when flush, `"overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700"` when not.
- Wrapper div uses the helper directly (not `utils.Class()`) — it's an internal element with no consumer `props.Class`, so this preserves class ordering and avoids breaking existing tests that assert `"border border-gray-200"`.

### Core Feature: `TableCellPadding` typed enum (compact padding)

- **`TableCellPadding`** typed enum: `TableCellPaddingComfortable` (px-4 py-3, default) / `TableCellPaddingCompact` (px-4 py-2).
- Follows the library's typed-enum convention: `tableCellPaddingLookup` map + `tableCellPaddingClass()` via `utils.Lookup` fallback + `TableCellPaddingIsValid()` method.
- Wired into all `<th>` and `<td>` cells — both plain headers and `TypedHeaders` paths, and the `<td>` body cells.
- `CellPadding` field added to `TableProps` with godoc.

### Tests (all passing)

| Test                                                 | Type        | What it covers                                                          |
| ---------------------------------------------------- | ----------- | ----------------------------------------------------------------------- |
| `TestTableFlush` (2 subtests)                        | Unit        | Flush=true suppresses border; Flush=false has border                    |
| `TestTableCellPaddingOption` (5 subtests)            | Unit        | Default, compact, explicit comfortable, invalid fallback, typed headers |
| `TestGoldenTableFlush`                               | Golden      | Full HTML snapshot of flush table                                       |
| `TestGoldenTableCompact`                             | Golden      | Full HTML snapshot of compact table                                     |
| `TestTableFlush/default_table_has_wrapper_border`    | BDD         | Behavioral assertion on border presence                                 |
| `TestTableInCardNoDoubleBorder` (2 subtests)         | Integration | End-to-end: Table inside Card — flush=1 rounded-lg, non-flush=2         |
| `TestTableWrapperClassCoverage`                      | Coverage    | Both branches of `tableWrapperClass`                                    |
| `TestTableCellPaddingClassCoverage`                  | Coverage    | All paths: comfortable, compact, empty, invalid                         |
| `TestRTLRendering/table_uses_logical_text_alignment` | RTL         | `text-start` present, `text-left` absent                                |
| `TableCellPaddingIsValid` entries (3)                | Enum        | Comfortable, Compact, invalid                                           |
| `ExampleTable_flushInCard`                           | Example     | Godoc example for the flush + compact pattern                           |

### Docs Updated

| File                  | Change                                                                                 |
| --------------------- | -------------------------------------------------------------------------------------- |
| `CHANGELOG.md`        | `[Unreleased]` section: Added (Flush, CellPadding) + Fixed (double border)             |
| `AGENTS.md`           | Flush/CellPadding conventions added; enum counts updated (31 IsValid, 15 map+fallback) |
| `FEATURES.md`         | Table entry: `Flush` and `CellPadding` listed                                          |
| `skill/SKILL.md`      | Both catalogue tables (by use case + by package) updated                               |
| `display/table.templ` | Full godoc on `Flush`, `CellPadding`, `tableWrapperClass`, `tableCellPaddingClass`     |

### Demo

- New "Table in Card (Flush + Compact)" section in `examples/demo/demo.templ` — demonstrates the canonical dashboard pattern from the bug report.

### Verification

- `go build ./...` — clean
- `go test ./... -count=1 -race` — all packages pass
- `golangci-lint run` — 0 issues
- `*_templ.go` files regenerated and committed
- `.gitignore` checked after every commit (BuildFlow gotcha not triggered)
- All 7 commits pushed to `master`

---

## b) PARTIALLY DONE

Nothing — the feature is complete and all tests pass. The only "partial" aspect is that this is `[Unreleased]` — not yet cut as a version tag.

---

## c) NOT STARTED

- **Version cut** — `utils.Version` still says `0.15.0`. The `[Unreleased]` CHANGELOG section has content, ready for `scripts/release.sh 0.16.0 "..."`. Not started because user didn't ask for a release.
- **Consumer notification** — cqrs-htmx adminui still has the CSS workaround (`.overflow-hidden > .overflow-x-auto { border: 0 !important }`). They should remove it and use `Flush: true` once they upgrade.
- **README.md** — the main README component catalogue was not updated. The README has its own Table section that doesn't mention Flush/CellPadding.

---

## d) TOTALLY FUCKED UP

Nothing in this session. All 7 commits are clean, BuildFlow passed on every commit, no reverts needed.

**One near-miss:** the initial integration test used `strings.Count(output, "border border-gray-200")` to detect the double border, but `utils.Class()` reorders classes so the literal substring didn't match. Fixed immediately by switching to `rounded-lg` count (which is stable — `rounded-lg` only appears on wrapper-level elements, not internal cells).

---

## e) WHAT WE SHOULD IMPROVE

### Process improvements

1. **First-pass was incomplete.** The initial implementation (session 1) shipped the feature + unit/golden tests but missed: integration test, coverage tests, FEATURES.md, SKILL.md, demo, RTL test. The self-review (session 2) caught all of these. **Lesson:** always run a self-review checklist before declaring done — the per-component testing checklist in SKILL.md exists for this reason and wasn't followed.

2. **The wrapper div bypasses `utils.Class()`.** The `tableWrapperClass()` result is used directly in the template (`class={ tableWrapperClass(props.Flush) }`) instead of going through `utils.Class()`. This was necessary to preserve class ordering for existing tests, but it means consumer `props.Class` is NOT applied to the wrapper div — only to the `<table>` element inside. This is the pre-existing behavior (the old code also hardcoded the wrapper div class), so it's not a regression, but it's an inconsistency worth noting.

3. **Golden files are brittle to class reordering.** The golden file for the existing sortable table test (`table_sortable_headers.golden`) was not affected because we didn't change the non-flush path's classes, but any future `utils.Class()` change could reorder classes and break golden tests. The CSS normalization in `internal/golden` handles this partially.

### Code improvements

4. **`TableCellPadding` is a separate type from `CardPadding`.** Both control vertical density but have different value sets. This is correct (they're different domains), but a future `DensityLevel` shared enum could unify them if more components need it.

5. **The `Flush` name could be `FlushToParent`** for clarity, but `Flush` matches CSS terminology (`flush` = no gap/border) and is concise. No change needed.

6. **No `Compact` shorthand.** The feedback report suggested a `Compact` option. We implemented `CellPadding: TableCellPaddingCompact` instead, which is more flexible (allows future density levels). This is the right call — a boolean `Compact` would preclude a future `TableCellPaddingDense` or `TableCellPaddingSpacious`.

---

## f) Up to 50 things we should get done next

### Immediate (this feature)

| #   | Task                                                                            | Impact | Effort |
| --- | ------------------------------------------------------------------------------- | ------ | ------ |
| 1   | Cut v0.16.0 release with Flush + CellPadding                                    | High   | Low    |
| 2   | Update README.md Table section with Flush/CellPadding                           | Medium | Low    |
| 3   | Remove the CSS workaround in cqrs-htmx adminui and adopt `Flush: true`          | High   | Low    |
| 4   | Add a `docs/recipes/table-in-card.md` recipe showing the full dashboard pattern | Medium | Low    |

### Table component improvements

| #   | Task                                                                                                     | Impact | Effort |
| --- | -------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 5   | Add `TableSize` enum (sm/md/lg) controlling font-size + cell-padding as a unified control                | Medium | Medium |
| 6   | Add `StickyHeader bool` — `sticky top-0` on `<thead>` for long tables                                    | Medium | Low    |
| 7   | Add `ColumnAlign []TextAlign` — per-column text alignment (start/center/end)                             | Medium | Medium |
| 8   | Add `EmptyState templ.Component` slot — renders when `Rows` is empty (currently renders empty `<tbody>`) | High   | Low    |
| 9   | Add `Loading bool` — renders skeleton rows placeholder                                                   | Medium | Medium |
| 10  | Add `Selectable bool` + `Row.Selected` — checkbox column for bulk actions                                | High   | High   |
| 11  | `DataTable` wrapper (PLANNED in FEATURES.md) — sorting + filtering + pagination in one component         | High   | High   |
| 12  | Add `FooterRow templ.Component` — `<tfoot>` support for totals/summaries                                 | Medium | Low    |
| 13  | Add `StickyFirstColumn bool` — `sticky start-0` for row label columns                                    | Low    | Medium |
| 14  | Add `ZebraOdd bool` — control which row index gets striping (currently hardcoded to odd)                 | Low    | Low    |

### Cross-component improvements

| #   | Task                                                                                                                             | Impact | Effort |
| --- | -------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 15  | Unify `DensityLevel` enum across Card, Table, Badge, Avatar (compact/comfortable/spacious)                                       | Medium | High   |
| 16  | Add `SharedBorderClass` constant — the `border border-gray-200 dark:border-gray-700` string is duplicated between Card and Table | Low    | Low    |
| 17  | Extract `overflowXAutoClass` shared constant — currently inline in Table wrapper                                                 | Low    | Low    |
| 18  | Add `Card.BodyFlush` — alias for `CardPaddingNone` with clearer naming for table-in-card use case                                | Low    | Low    |

### Testing improvements

| #   | Task                                                                                       | Impact | Effort |
| --- | ------------------------------------------------------------------------------------------ | ------ | ------ |
| 19  | Add fuzz test for `TableCellPadding` — verify no panic on arbitrary string input           | Low    | Low    |
| 20  | Add benchmark for `tableWrapperClass` + `tableCellPaddingClass`                            | Low    | Low    |
| 21  | Add a11y test for Table — verify `aria-label` propagation, `scope="col"` presence          | Medium | Low    |
| 22  | Add dark-mode compliance test for Table — verify all `border-gray-*` have `dark:` variants | Medium | Low    |
| 23  | Add snapshot test for Table-in-Card composition (not just integration assertion)           | Low    | Low    |
| 24  | Test Table with `Flush=true` + `Bordered=true` interaction (both borders on table element) | Low    | Low    |

### Documentation improvements

| #   | Task                                                                                | Impact | Effort |
| --- | ----------------------------------------------------------------------------------- | ------ | ------ |
| 25  | Add ADR for the Flush pattern (component opts out of its own border when nested)    | Low    | Low    |
| 26  | Update `docs/recipes/horizontal-filter-bar.md` to mention compact table for results | Low    | Low    |
| 27  | Add Flush pattern to CONTRIBUTING.md conventions section                            | Low    | Low    |
| 28  | Document the wrapper-div-not-using-utils.Class decision in table.templ comment      | Low    | Low    |

### Broader library improvements

| #   | Task                                                                                                   | Impact | Effort |
| --- | ------------------------------------------------------------------------------------------------------ | ------ | ------ |
| 29  | Audit all components for hardcoded border classes that can't be suppressed when nested                 | Medium | Medium |
| 30  | Add `Flush` to SimpleCard (currently inherits via Card, but SimpleCard with Flush isn't documented)    | Low    | Low    |
| 31  | `StatCard` + `Card` nesting — does StatCard inside Card produce a double border? Audit needed          | Medium | Low    |
| 32  | Add `Tooltip.Flush` for tooltips inside bordered containers                                            | Low    | Medium |
| 33  | Review if `Badge` inside `Card` header has border conflicts                                            | Low    | Low    |
| 34  | `Form` component — does it have a border that conflicts when nested in Card?                           | Low    | Low    |
| 35  | Create a "nesting compatibility matrix" doc — which components can nest inside which without conflicts | Medium | Medium |

### Type model improvements

| #   | Task                                                                                                                                                  | Impact | Effort |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 36  | Make `TableHeader.SortDirection` use `SortDirection` type (currently uses the same type but no `SortDirectionIsValid` on the struct)                  | Low    | Low    |
| 37  | Add `TableColumnWidth` typed enum (auto/fixed/narrow/wide)                                                                                            | Low    | Medium |
| 38  | Consider `TableRowAction` struct for row-level action buttons (edit/delete)                                                                           | Medium | Medium |
| 39  | Add `TableColumn` struct with `Header`, `Width`, `Align`, `Sortable` — replace parallel `Headers`/`TypedHeaders` with unified `Columns []TableColumn` | High   | High   |

### Infrastructure

| #   | Task                                                                                                                  | Impact | Effort |
| --- | --------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 40  | Add visual regression testing (Playwright/screenshot diff) — golden files catch HTML changes but not visual rendering | High   | High   |
| 41  | Add `docs/recipes/dashboard-table.md` with sorting + pagination + compact + flush all together                        | Medium | Low    |
| 42  | Create interactive Storybook-like demo page with all Table options toggleable                                         | Medium | High   |
| 43  | Add CSP test for Table with clickable rows — verify nonce on the row-href script                                      | Low    | Low    |
| 44  | Add `utils.AssertCount` test helper — `strings.Count` assertions are inline and repeated                              | Low    | Low    |

### Consumer feedback items

| #   | Task                                                                                          | Impact | Effort |
| --- | --------------------------------------------------------------------------------------------- | ------ | ------ |
| 45  | Collect feedback from cqrs-htmx on whether Compact padding is compact enough (py-2 vs py-1.5) | Medium | Low    |
| 46  | Ask if consumers want a `TableVariant` enum (simple/bordered/flush) instead of separate bools | Low    | Low    |
| 47  | Survey: do consumers nest other components inside Card(CardPaddingNone) that need Flush?      | Medium | Low    |
| 48  | Document the consumer CSS workaround removal path in the CHANGELOG release notes              | Low    | Low    |
| 49  | Add migration guide: "From CSS workaround to Flush prop"                                      | Low    | Low    |
| 50  | Review all cqrs-htmx adminui table templates (5 reported) for Flush adoption                  | Medium | Low    |

---

## g) Top 2 Questions

### Q1: Should the `Flush` pattern be a reusable concept across other components?

**Context:** Table-in-Card is the first case, but `Form`, `Alert`, `DefinitionList`, and future components might also need to suppress their border when nested inside a `Card(CardPaddingNone)`. Should we:

- **(a)** Keep `Flush` as a Table-specific field (YAGNI for now), or
- **(b)** Create a shared `Flushable` concept (interface or shared BaseProps flag) that any component can opt into?

I lean toward **(a)** until a second component actually needs it — but the nesting compatibility audit (item #29 above) would surface this.

### Q2: Should we cut v0.16.0 now, or batch with more work?

**Context:** The `[Unreleased]` CHANGELOG section has two Added items (Flush, CellPadding) and one Fixed (double border). This is a meaningful consumer-facing fix (visible defect on every table page in adminui). Options:

- **(a)** Cut v0.16.0 now — small, focused release. Consumer can upgrade and remove CSS workaround immediately.
- **(b)** Batch with more Table improvements (EmptyState slot, StickyHeader, etc.) for a bigger v0.16.0.

I lean toward **(a)** — this is a bug fix, and making consumers wait for it behind unrelated features is worse than a small release.
