# Status Report — Kanban Optimistic Move + Honest Pending Register (ADR-0041)

**Date:** 2026-09-17 18:10 CEST
**Session scope:** design + implement + prove the "clear action register for client moves NOT yet on the server" (optimistic update, transparent to the human) for `display.KanbanBoard`, both transports.
**Repo state at end of session:** working tree clean; all work sits in daemon auto-commits (`07e43632`, `ab5a4414`, …) on local master. Full `nix run .#verify` GREEN, W3C HTML gate clean (250 goldens), all 9 kanban browser tests pass, lint 0 issues across all modules. **(2026-09-17 evening docs-health pass: §f harvested into TODO_LIST #238–#249/#250/#251/#258–#262; §b1 in progress by a concurrent session — `kanban/pending_state.png` + `failed_state.png` appeared in testdata; website CI drift found + repaired, see 18-09 report.)**

---

## a) FULLY DONE

1. **Research phase.** Read the full kanban implementation (`display/kanban.templ`, `kanban.go`, drag/announce JS, hidden-form contract), the demo endpoints, `visualtest/kanban_e2e_test.go`, and `templates/custom.css` kanban block before designing.
2. **Transport event contracts verified against ground truth (not docs):**
   - htmx 2.0.10 (vendored `layout/static/htmx.min.js`): `htmx:afterRequest` fires from `onload` for 2xx AND 4xx/5xx AND from `onerror`/`onabort`/`ontimeout`; `successful` is set only by response handling — undefined (never `true`) on network errors. `responseError` fires inside response handling; `sendError` only on network-level failures.
   - Datastar v1.0.3 (pinned bundle, sha-guarded): `datastar-fetch` dispatched on `document` (no bubbling needed), `detail = {type, el, argsRaw}`; sequences `started → error → finished` (HTTP ≥ 400) and `started → retrying… → retries-failed → finished` (network). `finished` after `error` is what makes the success-clear a safe no-op.
3. **ADR-0041 written** (`docs/adr/0041-kanban-optimistic-move-pending-register.md`) — decision, rejected alternatives (visible board-level chip), success/failure matrices, consequences.
4. **Library implementation:**
   - `display/kanban.templ`: `data-tc-kanban-count` hook on count badges, `data-tc-kanban-empty` hook on the empty-column placeholder, wired-only sr-only `role="alert"` region; godoc updated.
   - `display/kanban.go`: optimistic pipeline in `tcKbSubmit` (origin captured BEFORE placing — pinned), `tcKbPlace`, `tcKbSyncCounts`, pending registry keyed by board id, `tcKbSucceed` (3 independent clearers incl. announce poll re-budgeted 5s → 30s), `tcKbRevert` (exact-position restore, DOM recount, 4s `tc-kanban-move-failed` flash, alert announcement), transport listeners filtered on `data-tc-kanban-form`. Split into `kanbanPendingRegistryJS`/`kanbanPendingPipelineJS`/`kanbanPendingEventsJS` to satisfy funlen; emitted JS byte-identical (goldens prove it).
   - `templates/custom.css`: `.tc-kanban-pending` (dim + reduced-motion-safe spinner ring via `::after`, logical `inset-inline` for RTL) and `.tc-kanban-move-failed` (red border flash), both with dark variants.
5. **Unit tests** (`display/kanban_pending_test.go`): markup hooks (wired + read-only), ~25 pinned JS tokens, htmx/Datastar branch ordering, and the origin-before-place ordering pin (the regression that mattered).
6. **Goldens refreshed** (7 kanban goldens) — verified delta is exactly: count/empty hooks, alert region, new script.
7. **Demo:** move endpoints sleep 800ms (`kanbanDemoMoveDelay`) so the pending register is perceivable on the live demo; section description rewritten; demo contract tests (`TestDemoKanbanHTTPContracts` + parity) still pass.
8. **E2E browser proof** (`visualtest/kanban_e2e_test.go`): flaky board pair (`e-slow` stalls 1.2s, `e-fail` always 500s) under htmx AND real Datastar:
   - `TestKanbanE2EPendingStateBothTransports` — optimistic state (card in target column + pending class + `aria-busy` + both count badges) proven within a 500ms poll window, i.e. BEFORE the 1.2s response; then pending cleared + "Moved …" announced.
   - `TestKanbanE2EFailureRevertsBothTransports` — revert to original column, counts restored, alert text exact, flash self-clears after 4s.
   - All 7 pre-existing kanban e2e tests still pass.
9. **Docs:** CHANGELOG `[Unreleased]` entry; AGENTS.md kanban bullet extended (2026-09-17 section); FEATURES.md KanbanBoard line; `docs/transport-wiring.md` facts list.
10. **Full verification:** `nix run .#verify` green end-to-end (generate + workspace build + tests + per-module `GOWORK=off` tests + visualtest with real Chromium + per-module lint 0 issues); `scripts/check-html-valid.sh` clean (250 golden files); `node --check` on the emitted kanban JS.

## b) PARTIALLY DONE

1. **Visual state coverage:** the pending/failed states are proven via DOM/computed-style assertions, but no PNG evidence was captured (`nix run .#shots` on the delayed demo, or a visualtest golden of the pending card). Humans haven't SEEN the spinner ring/failed flash from this work.
2. **Test lenses:** the per-component testing checklist asks for golden (done), a11y (partially — alert region asserted in the markup test, nothing added to `kanban_a11y_test.go` proper), BDD (not in `kanban_bdd_test.go`), example (no `ExampleKanbanBoard` addition). The behavior is covered, the file conventions are not.
3. **Git history:** everything landed via the daemon as `chore: auto-commit N changed file(s)` blobs — including a mid-edit torn snapshot of ANOTHER session's work (`af96082a`, `bac46b97`, `b2b3413a`). The feature exists on local master but not as a reviewable, semantic commit, and nothing is pushed/PR'd.
4. **Concurrency story:** the ADR documents multi-move races (first swap clears a second in-flight move's indicator; failure after success is a no-op) as reasoning only — no test drives two overlapping moves on one board.
5. **README/website catalogue:** the `KanbanBoard` one-liners in README and the docs site were not touched to mention optimistic pending; only FEATURES/AGENTS/transport-wiring were.

## c) NOT STARTED (deliberately out of scope this session)

1. Retry affordance after a failed move (manual re-drag is currently the only path).
2. Escape-hatch prop (e.g., pending timeout / opt-out) — YAGNI'd on purpose, zero API change shipped.
3. Axe-core audit of the mid-move DOM state (the sweep audits static demo routes only).
4. Any `?transport=` single-transport view of the kanban demo (the param only affects the Wire section; kanban always renders both boards — unverified end-to-end on Cloud Run).
5. Page-weight measurement of the grown inline script (~+2.5KB unminified).

## d) TOTALLY FUCKED UP (honest list)

1. **My own revert-position bug:** first version of `tcKbOptimistic` captured `parent`/`next` AFTER `tcKbPlace` had already moved the card — a failed move would have "restored" to the wrong position. Caught by my own review BEFORE any e2e ran, fixed, and pinned with an ordering assertion. No user-visible damage, but it shipped into goldens first (see #2).
2. **Golden ordering mistake:** I regenerated goldens while the buggy JS was in place, then fixed the bug and had to refresh goldens again. Order should have been: implement → e2e → goldens last. Cost: one wasted cycle + a confusing golden diff mid-verify.
3. **Shared-worktree collision handled reactively, not proactively:** another session + the daemon were actively editing the same worktree (recipes demo refactor, ErrorPage redesign, chart label anchors, sidebar theming). I burned ~3 full verify cycles on their torn snapshots before recognizing the pattern. I repaired their breakage (sidebar golden, 13 chart goldens — anchor-only delta, 4 `cmd/tc/_sources/errorpage/*.templ` re-copies) but never paused to announce/coordinate, and my `TestPrerenderMatchesLiveServer` flake investigation only concluded after the tree stabilized.
4. **Minor sloppiness:** wrong chromedp poll option names twice (`WithPollTimeout` → `WithPollingTimeout`), one careless multiedit that glued a comment to a func signature in `kanban_e2e_test.go` (fixed immediately), and an abandoned bash/heredoc attempt at a Go helper (repo policy says use the edit tools; I self-corrected).

## e) WHAT WE SHOULD IMPROVE

1. **Stateful-JS work should be e2e-first:** the flaky-board harness would have caught the origin bug before goldens ever saw it.
2. **Treat goldens as the LAST step** of any JS pipeline change, not a mid-way checkpoint.
3. **Detect daemon/other-session activity before big gates:** a `git log --oneline -3` + `git status` before each verify would have saved cycles (AGENTS.md literally warns about this; I under-applied it).
4. **The daemon's 60s-budget snapshots keep producing broken commits that CI-only catches** — the existing guards caught everything this time (goldens, sources-sync, drift), which is the system working, but each catch cost a verify cycle.
5. **Accessibility of transient states** deserves a browser-level a11y assertion (aria-busy visible to AT mid-move), not just string checks.
6. **Escape-hatch surface:** if any consumer reports the optimistic default as wrong for them, the fallback is a flag — decide the API shape now (`KanbanBoardProps.Optimistic *bool` vs `PendingPolicy` enum) rather than under pressure later.

## f) UP TO 50 THINGS TO GET DONE NEXT (ordered by impact within this session's scope)

**Kanban / ADR-0041 follow-ups**

1. ~~Capture PNG evidence of pending + failed states (delayed demo + `nix run .#shots` or a visualtest golden with `WaitSelector`).~~ IN PROGRESS (2026-09-17 evening): a concurrent session captured `visualtest/testdata/kanban/pending_state.png` + `failed_state.png` (untracked at annotation time); TODO #248 tracks completion.
2. Add the BDD lens: `kanban_bdd_test.go` specs for "move looks instant", "pending is visible until confirmed", "failure restores".
3. Add `ExampleKanbanBoard_optimistic` godoc example.
4. Extend `kanban_a11y_test.go` with the `role="alert"` + `aria-busy` assertions where they conventionally live.
5. ~~Concurrency test: two overlapping moves on one board (first succeeds, second fails) — pin the documented no-op semantics.~~ harvested → TODO_LIST #249
6. ~~Concurrency test: two overlapping failures — both revert, both announced.~~ harvested → TODO_LIST #249
7. Browser-level a11y probe of the mid-move state (aria-busy queryable while in flight).
8. ~~Retry affordance decision: auto-retry once vs manual-only; if manual, consider a small "retry" hint in the alert region.~~ harvested → owner gate TODO_LIST #238; hint idea in ROADMAP
9. ~~Decide + document (or reject) a pending-timeout escape hatch for hung servers.~~ harvested → TODO_LIST #239
10. Measure + note the inline-script size delta; consider whether the events block justifies a minified form.
11. ~~Verify the kanban demo section renders correctly under `?transport=datastar` and `?transport=htmx` single-transport views (screenshots).~~ harvested → TODO_LIST #258
12. ~~Consider a visualtest golden for the pending card (deterministic with the flaky-board trick: stall the endpoint, WaitSelector the class, screenshot).~~ harvested → TODO_LIST #248 (merge with item 1)
13. ~~Document how consumers with `GlobalErrorHandling` get toast + inline revert together (complementary, not double-reporting) — one paragraph in transport-wiring.md.~~ harvested → TODO_LIST #260
14. Un-minified JS: run the kanban pipeline through the same lint/minify review as the htmx embed decision (one-line TODO).
15. Add `TestKanbanJSPendingSingletonIdempotence`-style guard if listeners could ever double-bind (currently guaranteed by the singleton guard; pin it explicitly for the new listeners).
16. Check pkg.go.dev rendering of the new godoc (KanbanBoardProps ADR reference).

**Repo hygiene noticed during this session**
17. Squash/polish the feature into a semantic commit (current history is daemon blobs); get `Fixes #N` linking if a TODO exists.
18. Push + open PR so CI (linux format check, html-validation, per-module lint) rules on the final tree.
19. ~~`nix run .#visual` FULL suite (not just `-run TestKanban`) before push — axe sweep + route goldens.~~ done — the 15:27 errorpage session ran the full visual suite incl. axe sweep (74s) green after the tree stabilized
20. ~~`scripts/ci-repro.sh --lint` once for the exact CI reproduction.~~ done — the 14:17 release session ran `ci-repro.sh --lint` ALL STEPS PASSED
21. ~~Investigate `TestPrerenderMatchesLiveServer` minute-boundary flake (prerendered timestamps vs live fetch under load) — make it deterministic or retry-tolerant.~~ harvested → TODO_LIST #250
22. ~~TODO_LIST.md: check for an existing optimistic-kanban entry; add/close items for the follow-ups above.~~ done — 2026-09-17 evening docs-health pass registered #238–#239, #247–#249, #250, #251, #258, #260–#262
23. ~~README + website docs-site `KanbanBoard` one-liners: mention optimistic pending register.~~ harvested → TODO_LIST #251
24. ~~Docs-count drift: AGENTS.md still says "248 golden files" somewhere in my memory of the corpus — utils guard passed, but re-check prose counts after the other session's additions (250 now).~~ done — verified 2026-09-17 evening: no stale 248/127 prose remains; `TestDocsCountDrift` green (goldens now 133 after the pending-state captures)
25. ~~Daemon-race hardening: consider a pre-commit guard that refuses to commit when `git status` shows files modified within the last N seconds (torn-snapshot tripwire) — BuildFlow-side fix.~~ harvested → TODO_LIST #232 (daemon commit-gate bundle)
26. ~~The `AGENTS.md` note "zero-value boards byte-identical" for `KanbanColumn.Action: nil` is now stale relative to the new count/empty hooks — reword to "visually identical".~~ done — AGENTS.md reworded 2026-09-17 evening
27. Consider gating `data-tc-kanban-count`/`-empty` hooks on `wired` (bytes-cleanliness for read-only boards) — deliberate no-op today, but document the choice. ← untouched = still open
28. ~~Post-landing: watch the concurrent ErrorPage/charts/sidebar session's CI run — its work mixed into the same commits.~~ done — 2026-09-17 evening: their 16:24 run green; the 16:28 red was NEW daemon drift (website go.mod/generated import), root-caused and repaired

**Pending register — product polish ideas (from the session, not yet decided)**
29. Pending ring color could follow `KanbanTone`/semantic tokens via `@theme` instead of hardcoded blues.
30. Failed flash duration (4s) as a CSS custom property so consumers can tune it.
31. Alert copy ("The board was restored.") — consider surfacing the server's error message when the response carries one (go-error-family `Why`/`Fix`).
32. Optimistic scroll/focus preservation when a swap lands mid-interaction (generic kanban+htmx papercut, observed while testing).
33. Consider `hx-sync` guidance for rapid multi-move users (drop vs queue) in the recipe docs.
34. Keyboard path: after a failed revert, move focus back to the move button that initiated it.
35. Datastar network-failure UX: ~2min of honest retrying before revert — consider surfacing "still retrying" in the live region on `retrying` events.
36. Add the pending register to the kanban card-anatomy recipe (`docs/recipes/kanban-card-anatomy.md`) so composed cards don't fight the spinner's top-right corner.
37. Reserve the card's top-right corner: document that `Content` slot authors should not place interactive elements at the top-right edge.
38. Mobile: verify the spinner ring is not clipped by `overflow-x-auto` on narrow viewports (pixel check at 375px).

**Testing infrastructure**
39. ~~Extract the flaky-board pattern into a reusable visualtest helper (stall/fail cards) for future dual-transport components.~~ harvested → TODO_LIST #262
40. ~~Add a `chromedp.WithPollingTimeout` primer to the chromedp-lessons section of AGENTS.md (option-name trap hit this session).~~ done — the AGENTS.md kanban bullet already documents the 500ms `WithPollingTimeout` pattern proven by the e2e
41. ~~Golden-update ordering: add one line to the skill/AGENTS ("goldens AFTER e2e for JS pipelines"). ← untouched = still open
42. ~~Consider a `go test`-level JS syntax gate (node --check) if node is guaranteed in CI — cheap tripwire for script edits.~~ harvested → TODO_LIST #259
43. Time-bound the announce poll's interval timer cleanup test (30s × 100ms per submit is fine, but a rapid-fire test would confirm no interval leak).

**Docs/writing**
44. ~~Ship a short blog/recipe: "Optimistic UI with htmx and Datastar — one markup, two runtimes" using this implementation as the case study.~~ harvested → ROADMAP (Optimistic-kanban product polish)
45. ~~Update `docs/javascript-guide.md` decision ladder with the kanban register as the worked example of transport-event-driven UI state.~~ harvested → TODO_LIST #261
46. ADR-0041: add the follow-up decisions (retry, timeout, tone-token) as a "Deferred" section once decided. ← event-gated, still open
47. ~~Website docs: kanban page (currently none) — the demo section prose is now long enough to deserve a real docs page.~~ duplicate → TODO_LIST #222

**Verification debt**
48. Run the full verify twice back-to-back on a quiet worktree to shake out the last flake (`TestPrerenderMatchesLiveServer`).
49. Re-run `TestDemoKanbanHTTPContracts` + parity after any future demo-delay change (they're the tripwire).
50. ~~After the next release cut: confirm the committed `*_templ.go` + goldens for kanban survive the release script's replace-strip dance (v1.17.0 lesson).~~ done — v1.18.0 (`511d3ed6`) shipped with all goldens + generated files intact (release session verified the tagged tree consumer-clean)

## g) QUESTIONS I CANNOT ANSWER MYSELF

1. ~~**Your link anchored `#wire-transport`** (the Wire demo section), but the quoted heading is the Kanban section. Was kanban the only target, or do you want the same optimistic/pending treatment for the Wire demo's buttons/forms (and eventually as a general `wire` pattern)?~~ routed → ROADMAP "Wire-demo optimistic pattern" (product direction, not a question blocking work)
2. ~~**Failure UX policy:** after a server rejection, should the library attempt ONE automatic retry before reverting (htmx would need a re-submit; Datastar retries network errors natively but not 4xx), or is immediate honest revert + manual retry the intended behavior?~~ routed → owner gate TODO_LIST #238
3. **Git/PR handling:** the feature currently lives inside daemon `chore: auto-commit` blobs on local master. Do you want me to leave history as-is, or squash into one semantic commit (and push / open a PR), accepting a force-with-lease rewrite of the daemon's tips?

---

_Prepared by Crush (GLM). Session: single-session design→implement→prove of ADR-0041; all gates green at time of writing._
