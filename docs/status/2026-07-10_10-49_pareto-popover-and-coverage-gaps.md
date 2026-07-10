# Status Report — 2026-07-10 10:49

**Session goal:** Pareto analysis of TODO_LIST.md + FEATURES.md → execute the highest-value work.

**Version:** 0.14.0 (uncommitted) | **Branch:** master

---

## a) FULLY DONE

| #   | Work                                                                                     | Evidence                                                                                                                                                                                                                                                                                                              |
| --- | ---------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **`display.Popover` component** — the #1 most requested missing component                | `display/popover.templ` — button-triggered floating panel, 4 positions (`PopoverPositionTop/Bottom/Left/Right`), `role="dialog"`, `aria-haspopup`/`aria-controls`/`aria-labelledby` wiring, click-outside + Escape dismissal, CSP-safe singleton JS (`window.tcPopoverAttached`), arbitrary content via children slot |
| 2   | **`PopoverPosition` typed enum + `IsValid()` + contract registration**                   | `display/enums_go.go:51`, `display/enums_test.go` (3 cases), `internal/contract/component_props_test.go:53`                                                                                                                                                                                                           |
| 3   | **Popover tests — render, a11y, edge cases, golden, dark mode compliance, content slot** | `display/popover_test.go` — 22 sub-tests, all passing. Golden file at `display/testdata/popover_bottom.golden`                                                                                                                                                                                                        |
| 4   | **`icons.IconRTL` test coverage** — shipped public API was at 0%                         | `icons/snapshot_test.go` — 3 sub-tests covering regular icon, spinner variant, and all 100 path icons. Icons package coverage: **47.1% → 75.9%**                                                                                                                                                                      |
| 5   | **`utils.AssertContainsAll` test coverage** — last 0% function in utils                  | `utils/utils_test.go` — 3 sub-tests (all present, single, zero args)                                                                                                                                                                                                                                                  |
| 6   | **Documentation sync** — CHANGELOG, FEATURES, README, TODO_LIST, SKILL.md                | Version bumped 0.13.0 → 0.14.0, component count 83→84, enum count 33→34, generated files 62→63, Popover moved PLANNED→DONE in TODO_LIST                                                                                                                                                                               |
| 7   | **Full verify passes**                                                                   | `nix run .#verify` — generate + build + all tests + lint, 0 issues                                                                                                                                                                                                                                                    |

---

## b) PARTIALLY DONE

| #   | Work                            | What's done                                                                    | What's missing                                                                                                                                                                        |
| --- | ------------------------------- | ------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Popover JS behavior testing** | Server-rendered HTML is tested (a11y attributes, positions, dark mode, golden) | JS open/close/click-outside/Escape behavior is only testable via Playwright (blocked — no Node.js in repo). Singleton guard pattern matches Dropdown/Tooltip which have the same gap. |
| 2   | **Popover in `examples/demo`**  | Component is built and tested                                                  | Not wired into the demo binary — consumers discover it via FEATURES.md/SKILL.md but not via running demo                                                                              |

---

## c) NOT STARTTED

| #   | Work                                            | Why                                                                                        |
| --- | ----------------------------------------------- | ------------------------------------------------------------------------------------------ |
| 1   | BDD test for Popover (`bdd_test.go`)            | Skill checklist says every component needs a BDD spec. Not written.                        |
| 2   | Example test for Popover (`ExamplePopover`)     | Skill checklist says godoc example. Not written.                                           |
| 3   | Snapshot test for Popover in `snapshot_test.go` | Broader composition snapshot not added (golden covers the core render path).               |
| 4   | Coverage boost test for Popover private helpers | `popoverLookupPosition` branches partially covered by edge case test but not exhaustively. |
| 5   | Popover in integration CSP nonce test           | `integration/csp_nonce_test.go` should assert nonce on Popover's inline script. Not added. |
| 6   | Commit the work                                 | All changes are uncommitted. User hasn't asked to commit.                                  |

---

## d) TOTALLY FUCKED UP

| #   | What                                                    | Impact                                                                                                                                                                                                        | Fixed?                                                         |
| --- | ------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| 1   | **First draft of `AssertContainsAll` test was garbage** | Dead code, useless stub variable, explanatory comments for why I "couldn't" test it, `//nolint:contextcheck` directives that didn't apply. Violated 3 critical rules (no comments, no dead code, be concise). | **Yes** — rewrote clean with 3 focused sub-tests. Lint passes. |
| 2   | **Tried to edit AGENTS.md with precise counts**         | User correctly stopped this: AGENTS.md is enduring context, not a metrics file. Precise counts rot; `25+` is the honest representation.                                                                       | **Yes** — left AGENTS.md alone.                                |

---

## e) WHAT WE SHOULD IMPROVE

| #   | Issue                                                    | Recommendation                                                                                                                                                                                                                                                                                                                |
| --- | -------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Popover is missing BDD + example + integration tests** | The skill's per-component testing checklist requires 8 test lenses. Popover has 4 (render, a11y, edge, golden + dark mode). Missing: BDD, example, snapshot, coverage boost, CSP nonce integration. This is incomplete by the library's own standard.                                                                         |
| 2   | **Popover not in demo binary**                           | `examples/demo/main.go` should showcase Popover so consumers can see it live. Every other display component is there.                                                                                                                                                                                                         |
| 3   | **`utils` coverage is still 48.2%**                      | `EnsureID` (71.4%) and the `Assert*` helpers (66.7%) are the drag. `EnsureID`'s `crypto/rand` branch is hard to test without injecting a reader, but the helpers could reach 100% with a `testing.T` subtest-failure capture pattern. Low priority — these are test utilities, not production code.                           |
| 4   | **`display` coverage unchanged at 72.5%**                | Popover added more code than tests covered. The generated `popover_templ.go` error branches (context cancellation, buffer release) are the gap — same pattern as every other component. Chasing this is low-value busywork.                                                                                                   |
| 5   | **No Popover recipe doc**                                | Other interactive components have recipe docs. Popover is simple enough that FEATURES.md + godoc suffices, but a recipe showing Popover + form content (e.g., filter panel) would help consumers.                                                                                                                             |
| 6   | **Initial garbage test draft wasted a cycle**            | I should have written the test correctly the first time instead of copy-pasting a pattern from another LLM session. The `testing.T` subtest-failure pattern (`t.Run` returning `false`) is a known Go idiom — I should have used it immediately or accepted the happy-path-only pattern that the other `Assert*` helpers use. |

---

## f) Up to 50 things we should get done next

### Immediate (this session's unfinished work)

1. **Write Popover BDD test** — `display/bdd_test.go`, user-visible behavior spec
2. **Write `ExamplePopover`** — godoc example compiles + renders
3. **Add Popover to CSP nonce integration test** — `integration/csp_nonce_test.go`
4. **Add Popover to `examples/demo`** — wire into demo binary
5. **Write Popover snapshot composition test** — broader than golden
6. **Commit all changes** — user hasn't asked yet

### High value — new components

7. **`DataTable`** (#44) — high-level sortable/filtering/pagination wrapper around `Table` + `TableHeader` + `Pagination`. #2 most requested.
8. **`FilterDropdown`** (#45) — purpose-built for HTMX filter bars
9. **`Slider`** (#46) — ARIA slider pattern
10. **`HoverCard`** (#51) — like Popover but hover-triggered (can compose from Popover)

### High value — testing & quality

11. **BDD tests for all Popover-adjacent components** — Dropdown, Tooltip, Modal, Drawer all lack BDD specs despite the checklist requiring them
12. **Reassess the 80% coverage target** (TODO #12) — the 72-76% numbers are in generated templ error branches, not real logic. Write a targeted analysis rather than blindly adding tests.
13. **Add dark golden test variants for Popover** — `popover_dark.golden`
14. **Fuzz test for `PopoverPosition` validation** — match the `FuzzInputType`/`FuzzButtonHTMLType` pattern
15. **Benchmark for Popover render** — match the benchmark suite pattern in other packages

### Medium value — polish & DX

16. **Popover `Trigger` slot** — currently `TriggerText string`; a `templ.Component` trigger slot would allow icon-only triggers
17. **Popover `Open` prop** — server-controlled open state (like Modal's `Open` field) for HTMX-driven popovers
18. **Popover `DismissButton`** — explicit close button inside the panel
19. **Popover arrow/pointer** — visual arrow pointing to trigger (like Tooltip has)
20. **Recipe doc: Popover + filter form** — `docs/recipes/popover-filter-panel.md`

### Documentation

21. **SKILL.md: add Popover to the authoring playbook examples**
22. **SKILL.md: update component count in the process section**
23. **Add Popover to `docs/adr/` — decision: why Popover uses `role="dialog"` not `role="tooltip"`**
24. **CONTRIBUTING.md: mention Popover in the component list**

### v1.0 track (deferred, from TODO_LIST)

25. **`Validate() error` on props structs** (#33)
26. **Move test helpers to `internal/testutil/`** (#34)
27. **Self-host htmx as default** (#35, ADR 0007)
28. **Semantic token layer `bg-tc-primary`** (#36, ADR 0008)
29. **Remove deprecated aliases** (#38)

### v2.0 track (deferred)

30. **Compound component pattern** for overlays (#39)
31. **Native `<dialog>` element** (#40)
32. **Headless/unstyled variants** (#41)
33. **CLI tool `templ-components add <component>`** (#42)

### Remaining new components from TODO

34. **`Rating`** (#47)
35. **`TagsInput`** (#48)
36. **`ContextMenu`** (#49)
37. **`Carousel`** (#50)
38. **`Calendar`** (#52)

### Infrastructure (blocked)

39. **Visual regression testing** (#13) — Playwright screenshot diff
40. **Demo site deployment** (#27)
41. **`awesome-templ` PR** (#28)
42. **`templ.guide` listing** (#29)
43. **SSH tag signing config** (#30)

### Code quality

44. **icons package still at 75.9%** — the uncovered remainder is `strokeIconRTL` (68.3%) and `iconPaths()` panic path. Could add a targeted test for the stray `|` separator panic.
45. **`internal/golden.Assert` at 0%** — the golden comparison function itself is untested in the package it's defined in (tested via callers).
46. **Add `Popover` to `integration/composition_test.go`** — cross-package composition proof
47. **RTL test for Popover** — verify logical properties mirror correctly (no `left-`/`right-` except centering)
48. **Motion-reduce check for Popover** — the panel doesn't have transitions yet, but if added, must include `motion-reduce:*`
49. **`goconst` audit** — Popover JS string literals could be extracted to constants
50. **AGENTS.md `skill/` reference** — the `skill/SKILL.md` path in git status suggests the skill may be vendored in-repo; verify this is intentional vs. a symlink to `~/.config/crush/skills/`

---

## g) Top 2 questions I cannot answer myself

### 1. Should I commit this work?

All changes are uncommitted (11 modified + 4 new files). The repo convention is "NEVER COMMIT unless user explicitly says commit." You haven't said commit. But this is a complete, verified feature (Popover + tests + docs) at v0.14.0. Should I commit it as a single feature commit, or do you want to review the diff first?

### 2. Is the Popover complete enough by this library's standard, or should I finish the full test matrix first?

The skill's per-component testing checklist requires 8 test lenses (golden, a11y, BDD, edge cases, example, snapshot, coverage boost, CSP nonce). Popover currently has 4-5 of these. The other display components (Dropdown, Tooltip, Modal) also don't have all 8 — so Popover matches the _actual_ bar, not the _documented_ bar. Should I close the gap for Popover specifically, or is matching the existing components sufficient?
