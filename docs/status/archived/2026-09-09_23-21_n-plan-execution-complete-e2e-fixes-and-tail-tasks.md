# Status Report — 2026-09-09 23:21 CEST — N-Plan Execution: E2E Suites, Defect Fixes, Tail Tasks

**Session scope:** resumption of the N1–N21 quality-tier plan. Entry state: N1–N2 complete,
N3 ~70% (suite written but never re-run). Exit state: **N1–N21 fully executed or explicitly
gated; final battery green** (`nix run .#verify` 21 packages, per-module loop 10 modules,
full visual suite 59.7s). 60 commits ahead of origin (daemon pushes on its own schedule).

**Verification at time of writing:** `nix run .#verify` PASS · per-module isolation PASS ·
full visual suite PASS (axe sweep + flows + smoke + mobile + RTL + route goldens +
synthetics + all component goldens) · `cmd/tc` sources sync PASS · docs count guards PASS.

---

## a) FULLY DONE

1. ~~**N3 — Demo click-through e2e suite is green.** All 5 flows pass in ~14s~~ done at `69880f6`
   ~~(LoadMore→EndOfList, ConfirmDelete, LoadingButton busy gate, multipart upload echo,~~
   ~~kanban move on both htmx and Datastar boards). Fixes that got it there: boolean~~
   ~~predicates polled into `string` (always errored → infinite retry) now poll into `bool`;~~
   ~~`demoClickUntil` wraps predicates in `Boolean(...)`; ConfirmDelete polls for the row's~~
   ~~REPLACEMENT (the old predicate could never be true); native `confirm()` is stubbed~~
   ~~(CDP trace proved the accept command is never sent on a dialog-paused target); every~~
   ~~flow tab is bounded at 120s (`newFlowTab`) so wedges fail instead of hanging the binary.~~
2. ~~**Real library bug fixed — tooltip JS TypeError.** The shared tooltip singleton called~~ done at `69880f6`
   ~~`e.target.closest(...)` unguarded in `keydown`/`mouseenter`/`focusin` document~~
   ~~listeners; a `mouseenter` fired on `document` (every pointer entry into the window on~~
   ~~any page with a tooltip) threw a console TypeError. All three handlers now route~~
   ~~through a target guard. Goldens updated; CHANGELOG `[Unreleased]` Fixed entry added.~~
3. ~~**N4 — visualtest lint triage: 65 findings → 0.** Real fixes: `resolveOptions`~~ done at `69880f6`
   ~~cyclomatic split (`mergeViewportOptions`), shared `e2ePageServer` helper replacing two~~
   ~~duplicated server fixtures, static base error for the wait-expression timeout,~~
   ~~blank-assign for best-effort artifact cleanup, golines via `golangci-lint fmt`.~~
   ~~Policy fix: one module-wide `visualtest/` waiver for `contextcheck`/`paralleltest`/~~
   ~~`wrapcheck` (serial-shared-browser design) replacing three per-file waivers; five stale~~
   ~~per-file `nolint` comments removed. tools/shots CLI brought to the same zero bar~~
   ~~(named constants, 0750/0600 permissions, `Fprintf(os.Stdout, ...)`).~~
4. ~~**N5 — CI demo smoke.** `visualtest/demo_smoke_test.go`: all 7 demo routes must render~~ done at `69880f6`
   ~~their unique `<title>` with zero 500s in the server log. Runs inside the existing~~
   ~~Visual Regression CI job (the flake's `.#visual` app already puts `go` on PATH, so the~~
   ~~demo-server harness works in CI without workflow changes). Commit `b8bad6d`.~~
5. ~~**N6 — 375px mobile sweep + real defect fixed.** Two new tests (no-horizontal-overflow~~ done at `69880f6`
   ~~across 4 routes; kanban reachable by internal panning). They caught a real defect:~~
   ~~KanbanBoard's `w-72 shrink-0` columns propagate ~1200px min-content through grid items~~
   ~~without `min-w-0`, widening the whole document by 875px on a phone. Fixed with~~
   ~~`min-w-0` on the board root (library) and the demo board wrappers (the actual grid~~
   ~~items); goldens updated; demo CSS recompiled; CHANGELOG Fixed entry. Commit `da2b3c8`.~~
6. ~~**N7 — RTL browser sweep.** `dir="rtl"` on all 4 demo routes: no overflow; kanban move~~ done at `69880f6`
   ~~buttons still resolve, submit, and land cards in the expected column on both~~
   ~~transports. Commit `6990fbb` (+ `4cf5007` golines).~~
7. ~~**N8 — Overlay open-state captures verified complete.** The plan was stale: Modal,~~ done at `69880f6`
   ~~Drawer, Tooltip, Combobox, Carousel, Dropdown, Popover, ContextMenu open-state captures~~
   ~~already existed. Verified by execution: 11 tests PASS, 0 SKIP.~~
8. ~~**N9 — Coverage margin.** recipes 60.3% → 63.6% via `coverage_matrix_test.go`~~ done at `69880f6`
   ~~(maximal/minimal slot permutations, BaseProps propagation) and~~
   ~~`TestRenderErrorPropagation` (failing writer + cancelled context exercise the generated~~
   ~~error branches). Library total 72.0% (71.98% exact) vs the 70% CI floor = the target~~
   ~~2pt headroom. Commit `0e396ff`.~~
9. ~~**N10 — `layout.Minimal` SEO support (#179).** New `SEO SEOMeta` field on MinimalProps;~~ done at `69880f6`
   ~~the four head tags (noindex/canonical/hreflang/JSON-LD) are emitted by a new shared~~
   ~~`seoHeadTags` sub-template also adopted by Base, so the two shells cannot drift.~~
   ~~Zero value emits nothing (Minimal stays dependency-free). Parity test asserts identical~~
   ~~tag sets across Base and Minimal. Commit `48f22e6`.~~
10. ~~**N11 — Page-level route goldens (#163).** Restored the 182-line suite from the~~ done at `69880f6`
    ~~`feat/layout-seo-meta` prototype with the missing `GOEXPERIMENT=jsonv2` build env,~~
    ~~checked type assertion, `NewRequestWithContext` health probe, gosec/wsl/nolint~~
    ~~compliance. 8 route captures (7 full-page + index above-the-fold), **0.0000%~~
    ~~run-to-run drift**; two goldens re-cut for legitimate v1.16 content drift (component~~
    ~~counter 118→120, Kanban TOC entry — diff PNG inspected before re-cutting).~~
11. ~~**N12 — Datastar JS synthetics (#147).** Three chromedp tests synthesize the runtime's~~ done at `69880f6`
    ~~document-level `datastar-fetch` CustomEvent: HTTP-error toast with status code,~~
    ~~retries-failed announcement, LiveRegion `aria-busy` clearing on first patch. No SSE~~
    ~~endpoint needed. Commit `afd0d43`.~~
12. ~~**N13 — FormLayoutInline width contract (#166).** Browser ground truth: fields render~~ done at `69880f6`
    ~~content-sized (218px inputs in a 780px form, same row) — the feared full-width~~
    ~~breakage does not exist. Pinned by `TestFormLayoutInlineWidthContract` (fails if a~~
    ~~wrapper gains `w-full` or fields stop sharing a row) and stated loudly in the enum~~
    ~~docs. Commit `aa4ebe2`.~~
13. ~~**N14 — DateRange docs + adjacent golden (#176).** Component doc now states the~~ done at `69880f6`
    ~~INLINE `<time>` semantics; `TestGoldenDateRangeAdjacent` pins the exact two-adjacent-~~
    ~~ranges output. Commit `4afebef`.~~
14. ~~**N15 — ErrorPage family matrix goldens (#177).** All five go-error-family values~~ done at `69880f6`
    ~~pinned through ErrorAlert plus the missing infrastructure family through ErrorDetail.~~
    ~~Commit `9c2256a`.~~
15. ~~**N16 — Prerender honesty (#167).** `examples/demo/prerender_diff_test.go`: prerender~~ done at `69880f6`
    ~~output vs `newMux()` live responses on all 7 routes are byte-identical after~~
    ~~normalizing exactly two by-design differences (live stylesheet link; fresh EnsureID~~
    ~~tokens). Commit `bf8a276`.~~
16. ~~**N17 — upstream-watch dry-run (#128).** Triggered via `gh workflow run -f~~ done at`69880f6`~~dry-run=true`; run 34391305391 completed **success** in 10s with real pin-extraction~~
    ~~and proxy/jq parsing in the log.~~
17. ~~**N18 — Docs mini-pack (#180).** innerHTML-no-scripts runtime fact appended to~~ done at `69880f6`
    ~~`docs/datastar-runtime-facts.md`; KanbanBoard FEATURES line now states touch-visible~~
    ~~move buttons; stale visual-golden counts 117→125 (README + ROADMAP) caught by~~
    ~~`TestDocsCountDrift`. Commit `3e162ae`.~~
18. ~~**N19 — Changelog policy (#133).** Policy decided and codified in~~ done at `69880f6`
    ~~`scripts/check-changelog-guard.sh` (library code requires CHANGELOG; test-only/~~
    ~~docs/examples/cmd/website exempt), CI job rewired to call it, shaken down locally~~
    ~~against 4 representative file lists (all behaved as specified). Commit `2f8d8bf`.~~
19. ~~**N21 — Flake input fold (#146).** `nixpkgs-go` folded into `nixpkgs` — both had~~ done at `69880f6`
    ~~locked the identical rev, so the split insulated nothing. Zero drift proven: templ~~
    ~~regenerates clean (v0.3.1020 pin intact), Go stays 1.26.7 (GO-2026 fixes), `nix flake~~
    ~~check` passes. Commit `3659d00`.~~
20. ~~**Final battery.** `nix run .#verify` green (21 packages); per-module loop green~~ done at `69880f6`
    ~~(utils/icons/errorpage/charts/echarts/datastar/htmx); full visual suite green~~
    ~~(59.7s); `cmd/tc` `_sources` re-synced for the form.templ/base.templ doc edits (the~~
    ~~sync guard caught them — working as designed); TODO_LIST pruned of 17 shipped rows +~~
    ~~#188 (v1.16.0 shipped) with #133 narrowed to its real-PR verification. Commit~~
    ~~`64fd128`.~~
21. ~~**Hygiene.** Killed 3 orphaned `tc-demo` + 13 orphaned `tc-demo-route` leaked server~~ done at `69880f6`
    ~~processes from timed-out runs; 3 scratch probes removed after diagnosis.~~

## b) PARTIALLY DONE

1. **#133 changelog policy — real-PR verification.** Works: script + CI wiring + local
   4-case shakedown. Missing: the throwaway-PR shakedown on real GitHub CI. Blocker: needs
   the 60 queued commits pushed (not ordered). Effort: S once pushed.
2. **v1.17.0 release.** `[Unreleased]` is warm (a11y batch, tooltip fix, kanban mobile
   fix, Minimal SEO, route goldens, synthetics). Not cut: the question "cut v1.17.0 after
   N3 lands?" was posed in the previous session's report and remains unanswered. Effort:
   S (release script + tags + push).
3. **Coverage headroom precision.** Displayed 72.0% but exact value is 71.98% — the 2pt
   headroom holds at CI's one-decimal granularity, with ~0.02pt to spare in exact math.
   recipes is still the weakest package at 63.6%. Effort to widen: M.
4. ~~**N8.** Counted as done, but it was VERIFICATION of existing coverage, not new work —~~ done (docs-health pass 2026-09-08)
   ~~the plan file was stale on this point. Lesson recorded in (e).~~
5. **AGENTS.md chromedp guidance.** The existing note "moved to goroutine with
   page.HandleJavaScriptDialog(true).Do(ctx)" documents a pattern this session PROVED
   never works (the accept command is never sent on a paused target; the working pattern
   is stubbing `window.confirm`). The doc correction is not yet written. Effort: S.

## c) NOT STARTED

1. **N20 — Demo niceties (file-backed kanban state, dashboard-recipe kanban section).**
   Owner-gated on Q1/Q2 (touch-drag JS + feature-scope props). Deliberately not started.
2. **Remaining owner decisions:** #190 (touch drag), #191 (kanban feature scope), #192
   (already answered for v1.16.0 by the "GET SHIT DONE" order — v1.17.0 timing is the new
   instance), #80/#162 human-eyeball PNG packs.
3. **BuildFlow upstream work:** #93/#107/#108/#123/#124/#125/#126 all live in
   `larsartmann/buildflow` or GitHub settings, untouched.
4. **`awesome-templ` / `templ.guide` submissions (#28/#29).** Awaiting upstream approval.
5. **Website sync.** New docs (route goldens, e2e tiers) not yet reflected in website/
   content.

## d) TOTALLY FUCKED UP

1. ~~**The previous session shipped a broken e2e suite and declared it 60% done.** The~~ done (docs-health pass 2026-09-08)
   ~~ConfirmDelete "fix" (goroutine + `page.HandleJavaScriptDialog(true).Do(ctx)`) NEVER~~
   ~~worked — this session's CDP trace proved the accept command is never sent once the~~
   ~~renderer pauses on the dialog. The LoadMore flow polled boolean predicates into~~
   ~~`string` (guaranteed unmarshal error → infinite retry), and the ConfirmDelete predicate~~
   ~~checked `#item-123` for its own replacement text AFTER the swap removes that element —~~
   ~~it could never be true. `go vet` passing was treated as "60% done". Severity: the~~
   ~~suite hung the go-test binary for its full 10-minute budget twice. Mitigation: all~~
   ~~fixed and green, but the pattern (declare done without running) is the root cause.~~
2. ~~**Two 10-minute binary timeouts burned before instrumenting.** The first hang should~~ done (docs-health pass 2026-09-08)
   ~~have triggered a CDP-trace probe immediately; instead one full re-run (10 min) was~~
   ~~spent before the scratch probe cracked it in 1.4s. Root cause of the hesitation:~~
   ~~trusting that "fixed last session" meant verified.~~
3. ~~**My own scratch probe contained the same bool-into-string bug** as the suite it was~~ done (docs-health pass 2026-09-08)
   ~~diagnosing — the mirror failure delayed the ConfirmDelete diagnosis by one round.~~
4. ~~**Daemon races fragmented history all session.** Seven+ `chore: auto-commit` snapshots~~ done (docs-health pass 2026-09-08)
   ~~landed mid-work (e.g. `5d7c513`, `6352db3`, `29b5a88`, `37ac64e`, `2f1002d`,~~
   ~~`fc28003`); one summary commit became a no-op (nothing staged — daemon swept first);~~
   ~~one amend lost a race against the daemon. Net effect: meaningful changes are all in,~~
   ~~but the narrative lives partly in heuristic-message commits.~~
5. ~~**Lint-after-commit violations (twice).** The RTL sweep landed with an `intrange`~~ done (docs-health pass 2026-09-08)
   ~~finding; the golines finding landed in a second commit. The repo rule is~~
   ~~lint-per-file-at-authoring; both were avoidable.~~
6. ~~**Shell flails:** a literal-garbage `rm` invocation, a `/dev/null` sed target, wrong~~ done (docs-health pass 2026-09-08)
   ~~grep regexes, and one `git -C ..` from the wrong CWD each cost a round trip. None~~
   ~~caused damage; all were noise.~~
7. ~~**Not fucked up, but embarrassing:** the N11 restore briefly wrote goldens via failed~~ done (docs-health pass 2026-09-08)
   ~~`git show > file` redirects into a nonexistent directory (8 "MISSING" messages) before~~
   ~~`mkdir -p` — the commit itself ended up correct.~~

## e) WHAT WE SHOULD IMPROVE

1. ~~**Run every test at authoring time.** The suite was declared functional twice without~~ done (docs-health pass 2026-09-08)
   ~~a single execution. Rule going forward: a test does not exist until it has RUN green~~
   ~~(or failed for the expected reason).~~
2. ~~**Commit inside the daemon's 60-second window.** Batched commits lost races repeatedly.~~ done (docs-health pass 2026-09-08)
   ~~Green → commit immediately, message included.~~
3. ~~**Lint at authoring time** (repo rule, violated twice this session): run~~ done (docs-health pass 2026-09-08)
   ~~`golangci-lint run` on the touched module before `git commit`, not after.~~
4. ~~**Probe-first for hangs.** Any chromedp action exceeding ~30s should trigger a~~ done (docs-health pass 2026-09-08)
   ~~WithDebugf probe, not a re-run. The probe cracked in 1.4s what two 10-minute runs~~
   ~~did not.~~
5. ~~**Shared chromedp poll helpers.** Add `pollBool(ctx, expr)` / `pollText(ctx, expr)` to~~ done (docs-health pass 2026-09-08)
   ~~the visualtest harness; the bool-into-string mistake then cannot be written.~~
6. ~~**Centralize the chromedp trap list.** AGENTS.md's dialog note is factually wrong~~ done (docs-health pass 2026-09-08)
   ~~(post-trace) and the Poll/unmarshal trap isn't recorded. One AGENTS.md section update~~
   ~~prevents every future session from re-learning these.~~
7. ~~**Check plan claims against reality before executing.** N8 was already done; 20 minutes~~ done (docs-health pass 2026-09-08)
   ~~of planning would have been wasted without a quick `ls testdata` + test run.~~
8. ~~**Exact-coverage tooling.** The awk one-liner for exact totals should be a flake app~~ done (docs-health pass 2026-09-08)
   ~~(`.#coverage-exact`) so precision never depends on shell improvisation.~~
9. ~~**Daemon posture.** Consider committing detailed messages IMMEDIATELY on green even~~ done (docs-health pass 2026-09-08)
   ~~for small units (done mostly), and treat every daemon snapshot near delicate~~
   ~~operations as suspect (worked this session; keep it).~~
10. ~~**Scratch files.** The diagnose-then-delete pattern worked, but probes should carry a~~ done (docs-health pass 2026-09-08)
    ~~`-run TestScratch` prefix convention from the start so they can never run in a full~~
    ~~sweep by accident.~~

## f) NEXT 50 (ranked; feeds docs-health HARVEST)

| #  | Task                                                                                                                | Impact   | Effort | Category      |
| -- | ------------------------------------------------------------------------------------------------------------------- | -------- | ------ | ------------- |
| 1  | Cut v1.17.0 via release script (CHANGELOG warm: a11y batch, tooltip fix, kanban mobile fix, Minimal SEO)            | Critical | S      | Release       |
| 2  | Push/drain the 60 queued commits (or confirm daemon coverage is acceptable)                                         | Critical | S      | Release       |
| 3  | Real-PR shakedown of `check-changelog-guard.sh` (2 throwaway PRs: docs-only pass, component-without-changelog fail) | High     | S      | Quality       |
| 4  | Correct AGENTS.md chromedp guidance: dialog-accept stall (confirm-stub works; goroutine accept never did)           | High     | S      | Documentation |
| 5  | Add `pollBool`/`pollText` helpers to visualtest harness; migrate flow tests                                         | High     | S      | Quality       |
| 6  | Route goldens: add dark variants for all 7 routes (only dashboard_dark exists)                                      | High     | M      | Quality       |
| 7  | Route goldens: 375px mobile captures for the 4 swept routes                                                         | High     | M      | Quality       |
| 8  | Route goldens: RTL captures (the sweep asserts overflow; pixels pin mirroring)                                      | High     | M      | Quality       |
| 9  | Keyboard-only demo traversal (residue of #175: axe sweep covers DOM, not Tab-order UX)                              | High     | M      | Quality       |
| 10 | Triage the 16 LSP warnings on display/kanban (unused funcs flagged — verify against real lint, delete or wire)      | Medium   | S      | Cleanup       |
| 11 | Document or file upstream the stale-gopls kanban diagnostics (false "closing brace" error)                          | Low      | S      | Cleanup       |
| 12 | Human-eyeball pack #80/#162 (overlay PNGs + progressbar 45% fill)                                                   | Medium   | S      | Bug           |
| 13 | Answer Q1/Q2, then execute N20 (file-backed kanban state, dashboard recipe section)                                 | Medium   | L      | Feature       |
| 14 | Kanban e2e: HTML5 drag-and-drop path (only click-move is browser-proven)                                            | Medium   | M      | Quality       |
| 15 | Datastar synthetics: `retrying` event console path                                                                  | Low      | S      | Quality       |
| 16 | Datastar: mirror the no-scripts fact into package doc.go                                                            | Low      | S      | Documentation |
| 17 | Upload e2e: Datastar-dialect upload (currently htmx-only)                                                           | Medium   | M      | Quality       |
| 18 | LoadingButton e2e: assert label swap ("Saving…"), not just htmx-request class                                       | Low      | S      | Quality       |
| 19 | Users route e2e: pagination click-through                                                                           | Medium   | S      | Quality       |
| 20 | Wizard multi-step click-through in flows suite                                                                      | Medium   | M      | Quality       |
| 21 | Coverage: lift recipes to 70%+ (settings/auth slot branches)                                                        | Medium   | M      | Quality       |
| 22 | Flake app `.#coverage-exact` (weighted exact %, no awk)                                                             | Low      | S      | Quality       |
| 23 | docs/visual-testing.md: document the route-golden tier + update command                                             | Medium   | S      | Documentation |
| 24 | CONTRIBUTING: "adding an e2e test" checklist (bounds, poll patterns, no native dialogs)                             | Medium   | S      | Documentation |
| 25 | Verify v1.16.0 pkg.go.dev propagation (24h-watch item from N1)                                                      | Medium   | S      | Release       |
| 26 | CHANGELOG [Unreleased]: link issue numbers for the a11y batch rows                                                  | Low      | S      | Documentation |
| 27 | SKILL.md drift check: component/golden counts vs reality (TestSkillComponentCount logs drift)                       | Low      | S      | Documentation |
| 28 | Website: sync new testing tiers into site content                                                                   | Low      | M      | Documentation |
| 29 | `ci-repro.sh --lint`: include the changelog-guard script in local parity                                            | Low      | S      | Quality       |
| 30 | visualtest: `.fail/` stale subdirectory disk hygiene (pre-session dirs linger)                                      | Low      | S      | Cleanup       |
| 31 | Consider `disable-dev-shm-usage` flag for CI Chromium stability                                                     | Low      | S      | Quality       |
| 32 | Consider parallel route-golden capture (separate tabs; measure against 14s serial baseline)                         | Low      | M      | Quality       |
| 33 | File BuildFlow issues upstream (#93/#124/#125/#126) if the repo accepts them                                        | Medium   | M      | Bug           |
| 34 | Demo: MobileMenu/hamburger showcase (index nav is a wrap-nav; the component is underdemoed)                         | Low      | S      | Feature       |
| 35 | shared page-builder options (dark/rtl/viewport) for e2e page fixtures                                               | Low      | S      | Quality       |
| 36 | `datastar` package: string-level aria-busy unit test as CI-fast complement                                          | Low      | S      | Quality       |
| 37 | Dedupe screenshot-quality constant between tools/shots and harness                                                  | Low      | S      | Cleanup       |
| 38 | Re-check website pins (typescript/html-validate) after next daemon sweep (recurring flip risk)                      | Medium   | S      | Bug           |
| 39 | Prune `.fail/` artifacts from disk after green runs (t.Cleanup sweep)                                               | Low      | S      | Cleanup       |
| 40 | Consider `-timeout 15m` on the flake visual app as a belt-and-suspenders bound                                      | Low      | S      | Quality       |
| 41 | AGENTS.md: add "poll result type" to the chromedp trap list                                                         | Medium   | S      | Documentation |
| 42 | Evaluate moving `firstCSV`/kanban helpers into a shared e2e helpers file (used by flows + RTL)                      | Low      | S      | Cleanup       |
| 43 | Sweep for remaining `http.Get` without context in test fixtures (guard found one class; grep for others)            | Low      | S      | Quality       |
| 44 | Add release-notes section to GitHub Release automation (N1 manual refill lesson)                                    | Low      | S      | Release       |
| 45 | Consider tagging visualtest module (it has go.mod but no tag; release script covers 7 modules)                      | Low      | S      | Release       |
| 46 | Docs: record the route-golden re-cut procedure (diff PNG inspection before -update)                                 | Low      | S      | Documentation |
| 47 | Check whether `TestSkillComponentCount` should escalate from log to failure now counts stabilized                   | Low      | S      | Quality       |
| 48 | Demo `?transport=` page: add route golden for both single-transport variants (shots has them; goldens don't)        | Low      | S      | Quality       |
| 49 | Consider `t.Chdir`-free temp-dir cleanup verification for route goldens (RemoveAll defer audited)                   | Low      | S      | Cleanup       |
| 50 | Harvest this list into TODO_LIST/ROADMAP via docs-health (the canonical next step for section f)                    | High     | S      | Documentation |

## g) QUESTIONS I CANNOT ANSWER MYSELF

1. **Cut v1.17.0 now, or batch more?** `[Unreleased]` carries the a11y batch, the tooltip
   TypeError fix, the kanban phone-overflow fix, Minimal SEO, route goldens and the
   synthetics. All verified green. Cutting is ~15 minutes via the release script plus a
   push — but pushing/cutting is yours to order (the previous session's same question
   was superseded by "GET SHIT DONE" for v1.16.0 only).
2. **Push policy for the 60 queued commits.** The daemon pushes master on its own
   schedule anyway, so the queue is draining itself — but unvetted. Do you want a
   deliberate manual push after my review, or is daemon-driven push acceptable as the
   norm?
3. **Kanban scope (Q1/Q2 for N20/#189):** is a real touch-drag story (JS/polyfill) wanted
   despite the closed dependency budget, and which feature-scope props (within-column
   keyboard reorder, add-card wiring, WIP limits, tone accents) belong IN the component
   vs consumer composition? Without this, N20 stays gated and #189/#190/#191 cannot move.

---

_Report written per the status-report skill; the user's explicit `.md` instruction
overrides the skill's HTML default (5th consecutive override). Section (f) is the input
for docs-health HARVEST._
