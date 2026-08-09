# Status Report — 2026-08-03 03:53

> **Resolution (2026-08-05):** All section A items shipped in v1.7.0.
> Section B/F doc-count drift items resolved — README, FEATURES, SKILL.md,
> sections.ts counts all corrected to actual values (52 enums, 49 visual
> goldens, 175 HTML goldens, 112 components). The charts/echarts drift-guard
> gap (B/F item 1) is now fixed — the package is included in
> `countExportedTemplFunctions`. Dark-mode visual variants for 6 newer
> components and chart visual tests are tracked as TODO_LIST.md #95–#96.
> The CHANGELOG "Catull-Rom" typo was fixed.

## Session: TODO Execution Sprint Cleanup & Lint/Docs/Visual Coverage Pass

---

## A. FULLY DONE (Verified Green)

| #   | Task                                        | Evidence                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| --- | ------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **34 golangci-lint findings → 0**           | Renamed short vars in chart_geometry.go (`n`→`valueCount`, `x`→`posX`, `p0/p1/p2`→`prev/curr/next`), pie_chart.go (`cx`→`centerX`, `x1/y1`→`startX/startY`), removed unused `lineChartFormatFloat`, fixed `modalSizeLookup` goconst (reused `maxWSM`/`maxWXL` from drawer), fixed nlreturn/wsl_v5 whitespace, fixed predeclared `min`→`minVal`, refactored `TestNoOrderedTailwindSubstringsInTests` from gocognit 58→<30 via `scanLinesForOrderedSubstrings` + `isOrderedTailwindSubstring` extraction. `golangci-lint run` → **0 issues**. |
| 2   | **CHANGELOG `[Unreleased]` populated**      | Added 14 entries: AuthLayout, visualtest Bool()/ViewportMobile/Tablet/Desktop, InteractionState.String(), htmx golden tests (13 baselines), `nix run .#css`, `tc version`/`tc add --list-deps`, popover edge-flipping, docs/migration + testing-guide, Chromium pin, gofumpt alignment, dialog visual test fix, SkeletonCardGrid test fix. `[Unreleased]` is warm for release.                                                                                                                                                              |
| 3   | **docs/visual-testing.md updated**          | Options table updated for `*bool` tri-state, added `Bool()` section, viewport presets section, `InteractionState.String()`, `<dialog>` overlay testing subsection with `WaitSelector: "dialog"` example. MaxMismatch docs updated from 2%→1%.                                                                                                                                                                                                                                                                                               |
| 4   | **docs/cli.md updated**                     | Added `tc version` command and `tc add --list-deps` section with output examples.                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| 5   | **6 new visual golden tests**               | Combobox, Tooltip (with button child via `templ.ComponentFunc`), Carousel (3 slides), Skeleton (SkeletonCardGrid), ErrorPage (full-page), NotFound404 (navigation page). 49 total PNGs (was 43). All pass against committed goldens.                                                                                                                                                                                                                                                                                                        |
| 6   | **MaxMismatch empirical calibration (#82)** | Measured mismatch rates across 5 runs. Spinner: 0–3.57% (5% threshold confirmed). Overlays: 0–0.74% (tightened from 2%→1%). Dialogs: 0% (tightened from 2%→1%). Updated comments with empirical justification.                                                                                                                                                                                                                                                                                                                              |
| 7   | **`.envrc` fix**                            | Restored missing `GOEXPERIMENT=jsonv2` export — pre-existing bug causing `TestEnvrcConsistency` failure.                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| 8   | **`nix flake check` fix**                   | treefmt was failing: `flake.nix` `perSystem` needed multi-line form. Ran `nix fmt` → 1 file changed → `nix flake check` → **all checks passed**.                                                                                                                                                                                                                                                                                                                                                                                            |

**Verification commands all green:**

```
go build ./... → ok
go test ./... -count=1 -buildvcs=false → 16/16 packages ok
golangci-lint run → 0 issues
nix flake check → all checks passed
visualtest → 49/49 PNGs pass (headless Chromium)
```

---

## B. PARTIALLY DONE

| Item                           | What's Done                                                                  | What's Missing                                                                                                                                                                                                                                                                                                                            |
| ------------------------------ | ---------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Visual test coverage (#79)** | 49 PNGs across 28 test functions covering 24 component families              | No dark-mode variants for the 6 new tests (Combobox, Tooltip, Carousel, Skeleton, ErrorPage, NotFound404). Existing pattern tests both light AND dark; new tests are light-only. Only 1 `Bool(true)` call in the entire test file (spinner).                                                                                              |
| **CHANGELOG**                  | All session work documented                                                  | Chart entries from the parallel process have a typo: "Catull-Rom" should be "Catmull-Rom" (missing 'm'). The entry was there before my session — I didn't introduce it but I should have caught it.                                                                                                                                       |
| **Doc count sync**             | FEATURES.md/AGENTS.md/SKILL.md component counts were synced by prior session | README.md and FEATURES.md still say **31 visual goldens** (actual: 49). README says **43 typed enums** (actual: 45 with IsValid). FEATURES says **51 typed enums, 49 with IsValid** (actual: 45). AGENTS says **102 golden files** (actual: 175). These drifted because chart components added golden files but nobody updated the count. |

---

## C. NOT STARTED

| Item                                            | Description                                                                                                                                                                                                                                                                      |
| ----------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **`nix fmt` verification**                      | I did NOT run `nix fmt` until the very end when `nix flake check` caught the flake.nix issue. I should have run it at the start of the lint cleanup as part of the standard verify cycle.                                                                                        |
| **Dark-mode visual variants for new tests**     | The 6 new visual tests (Combobox, Tooltip, Carousel, Skeleton, ErrorPage, NotFound404) are light-mode only. Every pre-existing component family has both light and dark goldens.                                                                                                 |
| **Ordered-substring drift-guard negative test** | `TestNoOrderedTailwindSubstringsInTests` scans for violations but has no test proving it actually catches a violation. A negative test (inject a known `strings.Contains(out, "flex flex-col")` and assert it's flagged) would verify the guard works.                           |
| **Tooltip test pattern review**                 | The Tooltip test uses `templ.ComponentFunc` with `io.Writer` to inject children. This works but is unusual — the established pattern in the repo is to test components without children. Worth reviewing whether this is the cleanest approach or if there's a templ-native way. |
| **Popover edge-flipping browser test**          | The popover edge-flipping JS (4 conditions added to `display/shared.go`) was only SSR-verified (golden test). No visual test verifies the flipping actually works in a browser when a popover clips.                                                                             |

---

## D. TOTALLY FUCKED UP

| Item                           | What Happened                                                                                                                                                                                                                                                                                                                  | Impact                                                                                                                                                                                             |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Ignored `nix flake check`**  | I ran `go test` and `golangci-lint` but never ran `nix flake check` until the report-gathering phase. It was **failing** the entire session — `flake.nix` had a treefmt formatting violation from the parallel process's `inputs'` migration. CI would have gone red.                                                          | **HIGH** — `nix flake check` is in the standard verify cycle (`nix run .#verify`). If this had been pushed, CI would fail. Fixed at the very end, but it should have been the first thing checked. |
| **gopls stale diagnostics**    | 12+ false `IncompatibleAssign` errors on `visualtest/visual_test.go` persisted throughout the entire session. The code compiled and tested fine, but I never restarted gopls to clear them. This made LSP-based diagnostics unreliable and noisy.                                                                              | **MEDIUM** — wasted context window tokens, couldn't trust diagnostics, had to verify everything via `go build`/`go test` instead of trusting the IDE.                                              |
| **BuildFlow blank commits**    | 3 commits (`b5ef189`, `73056c2`, `01b3917`) have completely blank messages. The BuildFlow daemon committed these during the session. This is the known #93 problem — the daemon generates messages from a template, not from `git diff --stat`.                                                                                | **LOW** (cosmetic) — git history has 3 unscannable commits. Not actionable without fixing BuildFlow itself.                                                                                        |
| **Doc count drift not caught** | I updated the CHANGELOG and docs/visual-testing.md but did NOT update the component/golden/enum counts in README.md (31→49 goldens, 43→45 enums), FEATURES.md (31→49 goldens, 51/49→45 enums), or AGENTS.md (102→175 golden files). The `TestDocsCountDrift` test doesn't check visual golden counts, so it didn't catch this. | **MEDIUM** — docs are out of sync. Consumers reading README see stale numbers.                                                                                                                     |

---

## E. WHAT WE SHOULD IMPROVE

1. ~~**Always run `nix flake check` as part of the verify cycle.**~~ lesson absorbed.

2. ~~**The verify cycle should be: `nix run .#verify && nix flake check`.**~~ lesson absorbed.

3. ~~**gopls stale diagnostics should be restarted immediately**~~ lesson absorbed.

4. ~~**Doc-count maintenance is a known brittle pattern**~~ → TODO_LIST #112 (extend drift guard to README + ROADMAP).

5. ~~**Dark-mode visual test coverage should be mandatory**~~ done — dark-mode variants shipped for all newer components (v1.8.0).

6. ~~**The CHANGELOG typo "Catull-Rom" → "Catmull-Rom"**~~ done — fixed.

7. ~~**MaxMismatch thresholds should be documented with measurement methodology.**~~ done — documented in `docs/visual-testing.md`.

---

## F. Up to 50 Things to Get Done Next

### High Priority (Release Blockers)

1. ~~Fix doc count drift: README visual goldens 31→49~~ done — updated to 66 (2026-08-09).
2. ~~Fix doc count drift: FEATURES visual goldens 31→49~~ done — updated to 66.
3. ~~Fix doc count drift: AGENTS golden files 102→175~~ done.
4. ~~Fix CHANGELOG typo: "Catull-Rom" → "Catmull-Rom"~~ done.
5. ~~Add dark-mode visual variants for the 6 new tests~~ done — v1.8.0.
6. ~~Run `nix fmt` as part of every verify cycle~~ done.
7. ~~Extend `TestDocsCountDrift` to check visual golden count~~ → TODO_LIST #112.

### Visual Test Coverage Gaps

8. ~~Add visual test for CollapsibleSection~~ done — v1.8.0.
9. ~~Add visual test for Heatmap~~ done — v1.8.0.
10. ~~Add visual test for PolledRegion~~ done — v1.8.0.
11. ~~Add visual test for Sparkline~~ done — v1.8.0.
12. ~~Add visual test for BarChart~~ done — v1.8.0.
13. ~~Add visual test for ExternalLink~~ done — v1.8.0.
14. ~~Add visual tests for LineChart, PieChart, AreaChart~~ done — v1.8.0.
15. ~~Add visual test for ECharts adapter~~ open — needs browser-based JS execution; → ROADMAP.
16. ~~Add visual test for DataTable~~ done — v1.8.0.
17. Add visual test for ContextMenu open state ← open.
18. Add visual test for Badge variants ← open.
19. ~~Add RTL variants for all new visual tests~~ partially done; expanding → ROADMAP.
20. Add hover/focus state tests for interactive components beyond buttons ← open.

### Quality Hardening

21. ~~Add negative test for `TestNoOrderedTailwindSubstringsInTests`~~ open (nice-to-have).
22. ~~Add browser test for popover edge-flipping~~ open.
23. ~~Review Tooltip test `templ.ComponentFunc` pattern~~ open.
24. ~~Verify chart components pass `TestMotionReduceCompliance`~~ done — drift-guard passes.
25. ~~Verify chart components pass container-query compliance~~ done — drift-guard passes.
26. ~~Add benchmarks for chart geometry helpers~~ done — #109.
27. ~~Add fuzz tests for chart geometry~~ done — #98.
28. ~~Review chart component ARIA compliance~~ done — charts have `role="img"` + ARIA.
29. ~~Add CSP nonce test coverage for chart components~~ done — charts emit no inline scripts.
30. ~~Review `packageDeps` map in `cmd/tc/main.go`~~ done — `tc` tracks package sources.

### Documentation

31. ~~Update `docs/recipes/` with AuthLayout~~ done.
32. ~~Add chart components to `examples/demo/recipes_demo.templ`~~ done.
33. ~~Add AuthLayout to demo showcase~~ done.
34. ~~Update SKILL.md with chart component patterns~~ done — AGENTS.md carries the patterns.
35. ~~Update `docs/adr/` — ADR-0031~~ done.
36. ~~Write migration guide for chart components~~ done — `docs/recipes/`.
37. ~~Update `docs/tailwind-v4-adoption-guide.md`~~ done — chart colors use standard Tailwind palette.
38. ~~Add visual testing section to `docs/testing-guide.md`~~ done.
39. ~~Document `nix run .#css` in README~~ done.

### Infrastructure

40. ~~Fix BuildFlow blank commit messages~~ → TODO_LIST #93 (blocked).
41. ~~Consider pinning `nixpkgs-chromium`~~ done — pinned via separate flake input.
42. ~~Add `nix fmt` to pre-commit hook~~ done — treefmt runs via `nix flake check`.
43. ~~Verify `.github/workflows/ci.yaml` runs `nix flake check`~~ done.
44. ~~Review `visualtest/go.mod` dependencies~~ done.
45. ~~Consider adding visual test count to CI output~~ open (nice-to-have).

### Release Preparation

46. ~~Verify `[Unreleased]` CHANGELOG section is complete~~ done — [Unreleased] has 3 entries.
47. ~~Review all unpushed commits~~ done — pushed.
48. ~~Consider cutting v1.7.0~~ done — v1.7.0 shipped; v1.8.0 shipped.
49. ~~Run `scripts/release.sh 1.7.0`~~ done.
50. ~~Verify `git show v1.7.0` tag~~ done.

---

## G. Questions (Genuinely Cannot Answer Myself)

### 1. ~~The chart components — are these ready for release?~~ **Resolved:** yes — shipped in v1.7.0 with visual tests (v1.8.0), benchmarks (#109), and fuzz tests (#98). All drift-guards pass.

### 2. ~~Should I push the 21 unpushed commits now, or wait for the doc-count fixes?~~ **Resolved:** pushed — v1.7.0 and v1.8.0 are both on remote.

### 3. ~~The `nixpkgs-chromium` pin and the `nix fmt` / `nix flake check` failure — did the parallel process's flake.nix changes get reviewed?~~ **Resolved:** yes — the `inputs'` migration was intentional and reviewed; `nix flake check` passes.
