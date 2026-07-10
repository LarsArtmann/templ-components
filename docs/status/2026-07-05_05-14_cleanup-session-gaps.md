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

# Status Report — 2026-07-05 05:14 CEST

> **Updated:** 2026-07-06 (post-v0.8.0). Version at report: 0.6.x → **Current:** 0.8.0

**Session:** Feedback-driven improvements cleanup (session 6b) — closing gaps from self-review.
**Commits:** 3 (`985019a`, `2f20538`, `79c926c`)
**Done-check:** `nix run .#verify` → All checks passed (generate + build + test + lint, 0 issues)

> **UPDATE NOTE (2026-07-06):** Since this report, v0.7.0 and v0.8.0 were released.
> Most open items resolved. See status annotations below.

---

## a) FULLY DONE

| #   | What                                                                        | Evidence                                                    |
| --- | --------------------------------------------------------------------------- | ----------------------------------------------------------- |
| 1   | GridCols4/5 responsive ladders fixed (added intermediate `md` breakpoint)   | `display/grid.templ:27-28`                                  |
| 2   | ProgressBar clamp modernized (`max(0, min(100, v))`)                        | `feedback/loading.templ:118`                                |
| 3   | AGENTS.md: 8 new convention entries, BaseProps count 25→26                  | `AGENTS.md`                                                 |
| 4   | TODO_LIST.md: session 6 header added                                        | `TODO_LIST.md:3-6`                                          |
| 5   | FEATURES.md: Grid, SkeletonCardGrid, Script, GridCols enum                  | `FEATURES.md`                                               |
| 6   | Golden tests: Grid, StatCard.Href, SkeletonCardGrid (3 new `.golden` files) | `display/testdata/`, `feedback/testdata/`                   |
| 7   | BDD tests: Grid responsive, StatCard.Href clickable                         | `display/bdd_test.go`                                       |
| 8   | a11y tests: Grid aria-label/ID, StatCard.Href focus ring                    | `display/a11y_test.go`                                      |
| 9   | Example test: ExampleGrid godoc                                             | `display/example_test.go`                                   |
| 10  | Integration: Grid + StatCard.Href composition tests                         | `integration/composition_test.go`                           |
| 11  | CHANGELOG comprehensive [Unreleased] update                                 | `CHANGELOG.md`                                              |
| 12  | Generated files force-added (`grid_templ.go`, `script_templ.go`)            | tracked in git                                              |
| 13  | Planning doc with mermaid graph                                             | `docs/planning/2026-07-05_03-21_feedback-cleanup-pareto.md` |
| 14  | All pushed to `origin/master`                                               | `git push` confirmed                                        |

---

## b) PARTIALLY DONE

| Item                                    | Done                                                               | Missing                                                                                                                                            |
| --------------------------------------- | ------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| Test lens coverage for 3 new components | Display: golden+BDD+a11y+example+integration (Grid, StatCard.Href) | **Layout.Script**: zero golden, zero BDD, zero a11y, zero example tests. **feedback.SkeletonCardGrid**: golden only — no BDD, no a11y, no example. |
| README component count update           | Comparison table (76), display header (20), By the Numbers (76)    | **Feedback header still says "12 components"** — SkeletonCardGrid makes it 13. Missed in both commits.                                             |
| AGENTS.md metrics                       | Conventions added                                                  | Header metrics (components, tests, enums) in TODO_LIST header were updated, but AGENTS.md line 63 still says "25+" instead of "26+".               |

---

## c) NOT STARTED

| #   | What                                                  | Why it matters                                                                                                                                                                                                                               | Status (2026-07-06)                                                          |
| --- | ----------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| 1   | **SKILL.md update**                                   | New GridCols enum, Grid component, Script helper, SkeletonCardGrid are not in the skill's decision trees or conventions. The skill is the procedural counterpart to AGENTS.md — future sessions that load it will be missing these patterns. | ✅ Done — SKILL.md fully rewritten (Part 1 + Part 2)                         |
| 2   | **`examples/demo/` update**                           | Demo is the canonical "how a consumer assembles a page" reference. Grid + StatCard.Href not demonstrated.                                                                                                                                    | ✅ Done — demo updated                                                       |
| 3   | **Layout package test lenses for Script**             | Script has zero golden/BDD/a11y/example tests. The skill mandates all test lenses per component — Script got only assertion tests in `snapshot_test.go`.                                                                                     | ✅ Done — Script now has golden+BDD+a11y+example tests                       |
| 4   | **Feedback package test lenses for SkeletonCardGrid** | SkeletonCardGrid has golden + assertion tests but no BDD, no a11y, no example test.                                                                                                                                                          | ✅ Done — full test lens coverage added                                      |
| 5   | **CONTEXT.md update**                                 | Still says "Updated: 2026-06-27" with old metrics. Not updated this session.                                                                                                                                                                 | ⬠ Superseded — context now lives in AGENTS.md (CONTEXT.md no longer primary) |
| 6   | **Process 2 new untracked feedback files**            | `docs/feedback/2026-07-05_DiscordSync.md` and `docs/feedback/2026-07-05_swettyswipper-consumer-feedback.md` exist but were never read or processed. They may contain new actionable feedback.                                                | ✅ Done — all feedback processed in session 7                                |
| 7   | **StatCard golden with Href + Icon combined**         | The `stat_card_href.golden` tests Href without Icon. No golden covers the Href+Icon combination — a gap in regression coverage.                                                                                                              | ✅ Done — golden files cover Href + HTMX variants                            |

---

## d) TOTALLY FUCKED UP

| #   | What                                                                                                                                                                                                                                                                                                                                               | Severity                         | Status (2026-07-06)                                                      |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------- | ------------------------------------------------------------------------ |
| 1   | **README feedback count wrong — TWICE.** I bumped display (19→20) and the comparison table (73→76) but left `feedback — 12 components` unchanged when it should be 13 (SkeletonCardGrid). I literally identified this in the first self-review and still forgot to fix it.                                                                         | Low (cosmetic, but embarrassing) | ✅ Fixed in session 6b                                                   |
| 2   | **Layout.Script got the worst test coverage of any component in the library.** Every other component has golden+BDD+a11y+example. Script has only 3 assertion tests in snapshot_test.go. The skill explicitly says "each component package carries several complementary test lenses." I created the component, wrote minimal tests, and moved on. | Medium (consistency)             | ✅ Fixed — Script now has golden+BDD+a11y+example                        |
| 3   | **Feedback.SkeletonCardGrid is half-tested.** Got golden + assertion but no BDD/a11y/example. Same pattern — I hit the easy test types and skipped the rest.                                                                                                                                                                                       | Medium (consistency)             | ✅ Fixed — full test lens coverage                                       |
| 4   | **`.gitignore` still has `*_templ.go` on line 32.** This is the documented BuildFlow gotcha. I worked around it with `git add -f` instead of fixing the root cause. The next person who adds a component will hit the same trap. The fix is to remove line 32 (the trailing `*_templ.go` that overrides the `!*_templ.go` on line 2).              | Medium (recurring friction)      | ✅ Fixed — `.gitignore` now only has `!*_templ.go`, no trailing override |
| 5   | **SKILL.md forgotten entirely.** I listed it in the plan as T10 but never executed it. The skill is what teaches future sessions how to build components — missing GridCols, Script, SkeletonCardGrid from the decision trees.                                                                                                                     | Medium (knowledge debt)          | ✅ Fixed — SKILL.md fully rewritten                                      |

---

## e) WHAT WE SHOULD IMPROVE

1. ✅ **Fix README feedback count (12→13)** — Fixed.
2. ✅ **Remove `.gitignore` line 32** — Fixed. `.gitignore` now has only `!*_templ.go` on line 2, no trailing override.
3. ✅ **Add full test lens coverage for `layout.Script`** — Done.
4. ✅ **Add BDD + a11y + example for `feedback.SkeletonCardGrid`** — Done.
5. ✅ **Update SKILL.md** — Done. Full rewrite.
6. ✅ **Update `examples/demo/`** — Done.
7. ⬠ **Update CONTEXT.md** — Superseded. Context now lives in AGENTS.md.
8. ✅ **Read and process the 2 new feedback files** — Done in session 7.

---

## f) Up to 25 things to do next

| #   | Task                                                                                  | Impact | Effort | Status (2026-07-06)              |
| --- | ------------------------------------------------------------------------------------- | ------ | ------ | -------------------------------- |
| 1   | Fix README feedback count 12→13                                                       | Low    | 30s    | ✅ Done                          |
| 2   | Remove `.gitignore` line 32 (`*_templ.go`) — fix BuildFlow gotcha at root cause       | High   | 2m     | ✅ Done                          |
| 3   | Add golden test for `layout.Script`                                                   | Med    | 5m     | ✅ Done                          |
| 4   | Add BDD test for `layout.Script` (user includes CSP-safe script)                      | Med    | 5m     | ✅ Done                          |
| 5   | Add a11y test for `layout.Script` (nonce always present)                              | Med    | 5m     | ✅ Done                          |
| 6   | Add `ExampleScript` godoc example                                                     | Low    | 5m     | ✅ Done                          |
| 7   | Add BDD test for `feedback.SkeletonCardGrid` (user sees loading cards)                | Med    | 5m     | ✅ Done                          |
| 8   | Add a11y test for `feedback.SkeletonCardGrid` (role=status, motion-reduce, aria-busy) | Med    | 5m     | ✅ Done                          |
| 9   | Add `ExampleSkeletonCardGrid` godoc example                                           | Low    | 5m     | ✅ Done                          |
| 10  | Update `SKILL.md` — GridCols decision tree, Script pattern, SkeletonCardGrid          | Med    | 10m    | ✅ Done                          |
| 11  | Update `examples/demo/` — Grid + StatCard.Href showcase                               | Med    | 10m    | ✅ Done                          |
| 12  | Update `CONTEXT.md` with current metrics                                              | Low    | 5m     | ⬠ Superseded by AGENTS.md        |
| 13  | Read + process `docs/feedback/2026-07-05_DiscordSync.md`                              | High   | 10m    | ✅ Done                          |
| 14  | Read + process `docs/feedback/2026-07-05_swettyswipper-consumer-feedback.md`          | High   | 10m    | ✅ Done                          |
| 15  | Add golden test for StatCard with Href + Icon combined                                | Low    | 5m     | ✅ Done                          |
| 16  | Fix AGENTS.md "25+" → "26+" BaseProps count on line 63                                | Low    | 30s    | ✅ Done                          |
| 17  | Consider `GridProps.Gap` typed enum (gap-2/4/6/8)                                     | Low    | 10m    | ⬜ Not done                      |
| 18  | Consider `Card.Body` slot field (SEC feedback)                                        | Med    | 15m    | ✅ Done                          |
| 19  | Verify component count 76 by actual grep across all packages                          | Low    | 5m     | ✅ Done (82 components)          |
| 20  | Consider `layout.Stylesheet(nonce, href, attrs)` companion to `Script`                | Low    | 10m    | ✅ Done                          |
| 21  | Add `layout.Script` to the layout section of README component catalogue               | Low    | 5m     | ✅ Done                          |
| 22  | Cut v0.7.0 release once items 1–12 are done                                           | High   | 10m    | ✅ Done (v0.7.0 + v0.8.0)        |
| 23  | Audit all test files for `t.Parallel()` consistency                                   | Low    | 10m    | ⬜ Not done                      |
| 24  | Consider typed HTMX fields on StatCard (HxGet/HxTarget) vs Attrs workaround           | High   | 20m    | ✅ Done                          |
| 25  | Add CI check that `*_templ.go` files are tracked (prevent BuildFlow gotcha)           | Med    | 15m    | ⬜ Not needed — .gitignore fixed |

**Scorecard:** 21 of 25 complete (84%), 1 superseded, 3 not done.

---

## g) Top #1 question I cannot figure out myself

> ✅ **RESOLVED.** The `.gitignore` was fixed — line 32 was removed. Current `.gitignore`
> has only `!*_templ.go` on line 2. The BuildFlow managed block (lines 34+) no longer
> includes `*_templ.go`. The fix held through v0.8.0 — no re-additions observed.
