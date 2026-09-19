# Status Report — Kanban ADR-0041 Follow-Through: Test Lenses, Concurrency Pins, PNG Evidence

**Date:** 2026-09-17 19:47 CEST
**Session scope:** execute the follow-up list from `2026-09-17_18-10_kanban-optimistic-pending-register.md` — BDD/a11y/example test lenses, concurrency-semantics pins, and PNG evidence of the pending/failed states.
**Repo state at end of session:** all session work committed (daemon auto-commits `e9884c5c`…`b5043c85`); working tree has 3 modified docs files that are NOT this session's (see d6/g1). Filtered verification green: all new display tests pass, new visualtest golden test passes 3× consecutively (incl. verbose confirm), visualtest `go vet` clean, `golangci-lint run ./...` on display 0 issues.

---

## a) FULLY DONE

1. **Session bootstrap.** Loaded the templ-components + bdd-testing skills; recreated the todo list from the prior session's summary (9 tasks marked completed, follow-ups added); verified repo state first (`git log`/`git status` — AGENTS.md's own advice, applied this time).
2. **Scope correction discovered early:** `kanban_bdd_test.go`, `kanban_a11y_test.go`, `kanban_example_test.go`, `kanban_edge_test.go` already exist — the work was "extend with ADR-0041 coverage", not "create files". Read all of them before writing (convention: plain Go specs with `Spec:` comments in `package display`, NOT Ginkgo — repo convention won over the bdd-testing skill's Ginkgo default).
3. **New transport fact verified at source (not docs):** the vendored htmx 2.0.10 (`layout/static/htmx.min.js`) **serializes** same-element requests — with no `hx-sync`, a submit arriving while a request is in flight is queued as `"last"` and fires after the in-flight one completes (never aborts it). Datastar actions run concurrently. This materially sharpens the "two overlapping moves" race the ADR described.
4. **Concurrency pins** (`display/kanban_pending_test.go` → `TestKanbanJSConcurrentMoves`): optimistic placements stack (registry is a plain array, no dedup); revert is board-scoped and EXHAUSTIVE (loops every entry, `rest.push(e)` keeps other boards' entries); success clears the whole board (strip-all selector + `tcKbForget`); failure-after-success is a guarded no-op (`!tcKbPendingFor(b.id)` guard, ordering-pinned BEFORE the revert call). Extracted the shared `kanbanJSTokenAt` helper from the optimistic test.
5. **BDD lens** (`display/kanban_bdd_test.go`): `TestKanbanBehaviourMoveLooksInstant` — optimistic placement precedes `f.requestSubmit()` in the RENDERED board (guards a future refactor dropping the script component, asserted through the component output, not the bare JS string); `TestKanbanBehaviourCountLabelContract` — server-rendered count vocabulary (`aria-label="To do: 2 cards"`, `"no cards"`) matches the script's `tcKbCountLabel` vocabulary byte-for-byte, so optimistic recounts read identically to server renders.
6. **A11y lens** (`display/kanban_a11y_test.go`): `TestKanbanA11yFailureRegion` — wired-only sr-only `role="alert"` region, NEVER `aria-live="assertive"`, absent on read-only boards; `TestKanbanA11yBusyDuringFlight` — `aria-busy` set at submit and stripped on BOTH outcomes (exact count = 2 pin).
7. **Example lens** (`display/kanban_example_test.go`): `ExampleKanbanBoard_optimisticPending` godoc example (compile-only, matching the file's existing pattern) documenting that no extra props are needed.
8. **ADR-0041 updated:** new "Concurrency timings (verified in the vendored runtimes)" section — htmx serializes (queue `"last"` default), Datastar is concurrent; the register is conservative and correct under both; a result after the register cleared is a no-op by design.
9. **PNG evidence — the headline deliverable.** New internal test `visualtest/kanban_pending_visual_test.go` (`TestKanbanPendingRegisterVisualStates`): own stateless server (e-slow stalls until browser disconnect — no goroutine leak, clean `Server.Close`; e-fail 500s), one-board page, deterministic captures via (a) frozen spinner ring (inline style override = the exact rendering `prefers-reduced-motion` users get — a mid-rotation ring can never capture identically) and (b) the route-golden theme pin (localStorage + reload). Two goldens generated AND human-viewed:
   - `visualtest/testdata/kanban/pending_state.png` — Slow move in "In progress", dimmed, frozen blue ring top-right, counts 1/1.
   - `visualtest/testdata/kanban/failed_state.png` — Doomed move back in "To do" with red border flash, counts 2/0, "No cards" placeholder restored.
10. **Determinism proven:** 3 consecutive comparison-mode runs PASS (2 plain + 1 `-v` confirmed: `--- PASS: TestKanbanPendingRegisterVisualStates (2.22s)`).
11. **All 5 new display tests verified individually** (`-run` + `-v`: concurrent moves, move-looks-instant, count-contract, failure-region, busy-during-flight — all PASS; example compiles via package build).

## b) PARTIALLY DONE

1. **Full `nix run .#visual` (unfiltered)** — only kanban-filtered runs this session; the axe sweep + route goldens have not re-run on the current tree.
2. **`scripts/ci-repro.sh --lint`** — not run.
3. **Final `nix run .#verify`** — not run after this session's additions (partial gates green: display package tests + lint, visualtest vet + filtered visual).
4. **README + website KanbanBoard one-liners** — still untouched (carried over from the prior session).
5. **TODO_LIST.md reconciliation** — not done (carried over).
6. **AGENTS.md kanban bullet** — could absorb the new facts (2 PNG goldens, `TestKanbanJSConcurrentMoves`, htmx-serialization); not added. NOTE: daemon commit `b5043c85` touched AGENTS.md (6 lines) from another session — I did not check whether it already covers any of this.
7. **⚠️ Docs-count drift risk (unverified):** `utils.TestDocsCountDrift` counts visual goldens across README/FEATURES/AGENTS/skill — this session ADDED 2 visual goldens. If the guard pins an exact count, it now FAILS and prose counts need +2 in the same change. One command checks it: `cd utils && GOWORK=off go test ./... -run TestDocsCountDrift`. Not run before this report per your "report NOW" instruction — top of the next list.

## c) NOT STARTED (this session)

1. Browser-level a11y probe of the mid-move state (aria-busy queryable in-flight) — prior report item 7.
2. Retry-affordance decision + pending-timeout escape hatch (prior items 8–9; product decisions awaiting you).
3. Squash/push/PR of the daemon-blob history (prior §g3 — needs your call; force-with-lease requires approval).
4. All carried-over product-polish ideas (tone-token ring color, flash-duration custom property, alert copy surfacing server Why/Fix, focus restore after revert, Datastar "still retrying" announcements, hx-sync recipe guidance, card-anatomy corner reservation, 375px spinner clip check, flaky-board helper extraction, JS-size measurement, `?transport=` views, blog/recipe post, javascript-guide worked example, ADR Deferred section, website kanban docs page).

## d) TOTALLY FUCKED UP (honest list)

1. **First golden run captured DARK mode** — I hit the exact gotcha AGENTS.md documents (headless Chromium reports `prefers-color-scheme: dark` by default; the route goldens failed the same way before their fix). Saved by VIEWING the PNG before verifying: the pending state was perfect but the surfaces were dark. Fixed with the localStorage theme-pin + reload BEFORE any further gating. Cost: one regenerate cycle; lesson: eyeball every new golden at generation time, not after.
2. **Compile error from muscle memory:** used `r.Context()` inside the page-serve handler for `kanbanPendingVisualPage()` where `r` wasn't in scope (the page is request-independent). `go vet` caught it; one-line fix to `context.Background()`. Sloppy — the page never needed a request context.
3. **Stale LSP warning consumed a cycle:** `golangci_lint_ls` insisted `kanban_pending_test.go:125` violated gofumpt; `gofumpt -d` showed no diff and `golangci-lint run ./...` was 0 issues. I ran the formatter first and verified after — inverted order, though harmless (no-op).
4. **Mid-flight interrupt:** you requested this status report while the third determinism run was still in a background job. I had emitted no progress signal during the visual-evidence step (the longest of the session) — a one-line update before starting it would have set expectations.

## e) WHAT WE SHOULD IMPROVE

1. **"View the golden at generation time" is now a proven gate** — the image review caught the dark-mode trap instantly and cheaply. Make it a hard rule for NEW goldens (agent first-pass; human eyeball stays the second gate per the #80-family caveat in `kanban_visual_test.go`).
2. **The theme-pin is now duplicated in 3 places** (route goldens, siteshots, my visual test) — extract a harness helper (`pinTheme(ctx, "light")` or fold into capture) so the lesson stops being copy-paste.
3. **Frozen-animation capture deserves harness support:** my inline style override works, but a `ReducedMotion` option on `visualtest.Options` (raw cdproto `emulation.SetEmulatedMedia` — chromedp v0.16 lacks `EmulateMediaFeatures`) would make determinism declarative for future animation-state goldens.
4. **Test-lens completeness is convention, not enforcement:** the per-component checklist (golden/a11y/bdd/example) lives in the skill, but nothing fails when a new behavior ships without its lenses — I added them one session late. A drift-guard (component-with-pipeline ⇒ lens files mention it) would close this.
5. **Cross-session git noise continues:** the daemon auto-commits other sessions' edits (AGENTS.md, older status reports) while my work lands — the 3 modified docs files currently in the tree are not mine and I left them untouched (correct per the never-revert-others'-work rule, but the tree is never quiet for a clean verify).

## f) UP TO 50 THINGS TO GET DONE NEXT (ordered by impact)

**Immediate gates (this work's tail)**

1. Run `cd utils && GOWORK=off go test ./... -run TestDocsCountDrift` — the +2 visual goldens may break prose counts; fix counts in the same commit if so.
2. Full unfiltered `nix run .#visual` (axe sweep + route goldens on the current tree).
3. `scripts/ci-repro.sh --lint` — exact CI reproduction.
4. Final `nix run .#verify`.
5. W3C HTML gate re-run (`scripts/check-html-valid.sh`) — no HTML goldens changed, but cheap insurance on a shared tree.

**Docs increments for this session's work**
6. AGENTS.md kanban bullet: add the 2 PNG goldens, `TestKanbanJSConcurrentMoves`, and the htmx-serialization fact.
7. README KanbanBoard one-liner: "optimistic moves with pending register + failure revert" (+ website docs-site equivalent).
8. TODO_LIST.md: close the ADR-0041 follow-through items (test lenses, PNG evidence); open the remaining product decisions.
9. CHANGELOG `[Unreleased]`: extend the ADR-0041 entry with one clause (test lenses + visual evidence shipped).
10. ADR-0041 "Deferred" section: retry affordance + timeout hatch + tone-token ring (once decided, see g).
11. Check the daemon's AGENTS.md edit (`b5043c85`) for overlap with #6 before editing.

**Harness/test infrastructure**
12. Extract `pinTheme` theme-pin helper (3rd duplication site).
13. `ReducedMotion` capture option via cdproto `emulation.SetEmulatedMedia`.
14. Browser-level a11y probe: aria-busy queryable mid-move (e2e, not string-level).
15. Test-lens drift-guard idea: pipeline-bearing components must have bdd/a11y/example coverage.
16. Extract the flaky stall/fail board pattern into a reusable visualtest helper.
17. `chromedp.WithPollingTimeout` primer line in AGENTS.md's chromedp lessons (option-name trap, hit twice across sessions).
18. "Goldens AFTER e2e for JS pipelines" one-liner in the skill + AGENTS.md.
19. `node --check` JS syntax gate in CI if node is guaranteed there.
20. Announce-poll interval-leak test (rapid-fire submits).
21. Investigate `TestPrerenderMatchesLiveServer` minute-boundary flake (carryover).

**Git/release**
22. Squash daemon blobs into semantic commits; push; open PR (needs your approval, see g2).
23. After next release cut: confirm goldens + `*_templ.go` survive the replace-strip dance.
24. Watch the concurrent ErrorPage/charts session's CI run (its work shares this history).
25. Daemon torn-snapshot tripwire (pre-commit guard refusing commits with too-recent mtimes) — BuildFlow-side.

**Product polish (pending register)**
26. Retry affordance: auto-retry-once vs manual (g-pending).
27. Pending-timeout escape hatch for hung servers (API shape: `Optimistic *bool` vs `PendingPolicy` enum).
28. Pending ring color from `KanbanTone`/`@theme` tokens instead of hardcoded blues.
29. Failed-flash duration (4s) as a CSS custom property.
30. Alert copy surfacing the server's error message (go-error-family Why/Fix) when present.
31. Focus restore to the initiating move button after a failed revert.
32. Datastar "still retrying" live-region announcements on `retrying` events.
33. `hx-sync` guidance (drop vs queue) for rapid multi-move users in the recipe docs.
34. Document GlobalErrorHandling + inline revert as complementary (one paragraph, transport-wiring.md).
35. Card-anatomy recipe: reserve the card's top-right corner for the spinner ring.
36. Mobile 375px check: spinner ring not clipped by `overflow-x-auto`.
37. Consider gating `data-tc-kanban-count`/`-empty` hooks on `wired` (bytes-cleanliness; currently deliberate no-op).
38. Reword AGENTS.md "zero-value boards byte-identical" → "visually identical" (carryover, still open unless `b5043c85` did it).
39. Measure + note the inline-script size delta; minified-form decision for the events block.
40. Verify kanban demo section under `?transport=datastar` / `?transport=htmx` single-transport views (screenshots).

**Docs/writing**
41. Website kanban docs page (demo prose is long enough to deserve one; the new PNGs are ready-made figures).
42. Blog/recipe: "Optimistic UI with htmx and Datastar — one markup, two runtimes" using this as the case study.
43. `docs/javascript-guide.md` decision ladder: kanban register as the transport-event-driven worked example.
44. pkg.go.dev rendering check of the new godoc (ADR-0041 references + example).
45. Run the full verify twice back-to-back on a quiet tree (last-flake shakeout; carryover).

**Long-tail hygiene**
46. `TestCompiledCSSInventory`/golden-count guards: confirm no prose pins the visual-golden count in a second place beyond the drift guard.
47. Consider `-count=1` hygiene note: visual goldens passed 3×; record the determinism recipe (frozen animation + theme pin) in `docs/visual-testing.md`.
48. `TestDemoKanbanHTTPContracts` + parity re-run after any future demo-delay change (tripwire note, carryover).
49. Post-landing demo smoke: `nix run .#shots` on the kanban section for a fresh human-visible pass.
50. Archive/annotate this and the prior status report once the follow-ups land (docs-health ANNOTATE mode).

## g) QUESTIONS I CANNOT ANSWER MYSELF

1. **The 3 modified files in the working tree are not mine** (`docs/migration/v1-to-v2.md`, `docs/status/2026-09-17_05-54_…`, `docs/status/2026-09-17_09-28_…`) — is that another live session's work I should keep leaving untouched (my default), or do you want it committed/stashed before the next verify so the tree is quiet?
2. **Git/PR handling (carryover, still blocking push):** the ADR-0041 feature + this session's follow-through live in daemon `chore: auto-commit` blobs on local master. Squash into semantic commit(s) and push/PR (requires a force-with-lease rewrite of the daemon tips — your explicit approval), or leave history as-is?
3. **Should the two new PNGs stay test-only, or become website evidence?** `pending_state.png` / `failed_state.png` are clean, human-readable figures of the pending register and failure revert — I can build the website kanban docs page around them now, or keep them internal until that page is separately decided.

---

_Prepared by Crush (GLM). Session: test lenses + concurrency pins + PNG evidence for ADR-0041; all session-scoped gates green; repo tail (full visual + ci-repro + verify + docs-count check) queued._
