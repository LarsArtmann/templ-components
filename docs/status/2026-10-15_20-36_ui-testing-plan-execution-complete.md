# UI-Testing-Tooling Plan Execution — COMPLETE (status report §a–§g)

**When:** 2026-10-15 20:36 CEST · **Session:** continuation of the 17:10 report
(`2026-10-15_17-10_ui-testing-plan-execution-status.md`) — executed the remaining
~19 medium tasks (M06-finish → M26) + μ84 of
`docs/planning/2026-10-15_00-05_UI-TESTING-TOOLING-PARETO-PLAN.md`.

**Headline:** the plan is EXECUTED end-to-end. Verdicts + evidence in
`docs/testing/tooling-evaluation.md` §6; decision record ADR 0048; protocol +
budget `docs/testing/tool-trial-protocol.md`; engine matrix
`docs/browser-support.md`. **⫱G1–⫱G4 are now posed** (see §g). Master is ahead 5
commits, UNPUSHED — the pre-push ritual (ci-repro) has NOT run yet, and the
concurrent Tailwind-composability session is actively editing
`layout/appshell.templ` + `templates/custom.css` right now.

**What did I forget?** Three things, all caught and fixed in-session but each
cost real time: (1) the class-inventory self-heal — my `display/shared.go`
change staled `templates/templ-components-classes.txt`; the guard caught it on
the final battery run (I should have predicted it: ANY class-string change does).
(2) I accepted a "PASS" from the new M11 tests that was actually a silent SKIP
(no `CHROMEDP_CHROME_PATH` outside the nix env) — caught it only because I
requested `-v` output and saw "visual tests skipped". (3) The draft CI lane
pins `DeterminateSystems/nix-installer-action` at a SHA I have NOT verified
against upstream — it must be verified before anyone enables the lane
(verify-before-filing applies to my own lane too).

**What could I have done better?**
- The npm `accessibility-checker` install silently rewrote
  `tools/crossbrowser-spike/package.json`'s `playwright-core` pin 1.63.0 →
  1.64.1 and the daemon COMMITTED it — breaking the version-lockstep rule. I
  caught and restored it, but the right move was a scratch dir for the achecker
  trial + an immediate `git diff package.json` after any npm invocation.
- My "byte-determinism" check had a false positive (run 2 never executed; I
  compared stale files). Verify the thing RAN (mtimes/output) before trusting a
  match.
- `aria-snapshots.mjs` took 3 attempts (ariaSnapshot doesn't expose `aria-sort`;
  closed `<dialog>` snapshots empty; my helper re-navigated; /display has TWO
  sortable tables). I should have probed the component states before writing
  the script.
- Twice I used `edit` in a way that REPLACED a following section heading in the
  evaluation doc instead of inserting before it (M11, M12 headings) — both
  restored; sloppy anchor choice.

---

## a) FULLY DONE

| Task | Result |
| ---- | ------ |
| M06 finish (μ25–μ26) | 3rd CSP bug fixed at LIBRARY level: CSP-safe `data-tc-overlay-trigger` declarative attribute in `display/shared.go` overlay JS (the old `onclick="tcOpenOverlay(…)"` was doubly dead: undefined global + CSP-blocked). Demo buttons → `display.Button`; recipe doc re-taught; CHANGELOG entry; display goldens + `/display` route golden refreshed (targeted single-route refresh, temp test deleted after). Post-fix: 48/48 shots 0 errors; click-open verified in Chromium+Firefox+WebKit. Eval §5.1. |
| M07 (μ27–μ29) | `stylable-assert.mjs`: Firefox falls back to `appearance: auto`, native geometry (384×29 vs Chromium's styled 384×62), zero overflow, selection works → **degraded-but-usable, no bug**. Eval §5.2 + per-select crops. |
| M08 (μ30–μ33) | achecker over all 16 routes (initial → post-fix → final): caught REAL demo-wide duplicate ids (`status`, `name`, `q`, `email`, … reused across example forms), duplicate kanban/StepIndicator landmark names, unnamed regions. Fixed ALL demo-level findings (transport-suffixed ids, distinct labels, `role="group"` on code `pre`); repo-wide dup-id sweep now clean. **330 → 253 findings (−23%), 7 rule classes eliminated.** §5.3 diff table; component-level findings → TODO #421–#423. |
| M09 (μ34–μ36) | Eval doc fully populated (§5.1–§5.5 + §6 verdict table); ADR 0048 evidence bullets + ADDITIVE boundary. |
| M10 (μ37–μ40) | aria snapshots captured (table/modal-open/dropdown/tabs); μ40 regression demo PROVES ariaSnapshot cannot see `aria-sort` (YAML byte-identical after removal) → attribute asserts stay chromedp-owned. §5.4. |
| M11 (μ41–μ45) | `visualtest/visual_emulation_test.go` INTEGRATED (3 tests: vision-deficiency probe, reduced-motion kills spinner, print probe) — green under the nix browser env. Findings: no color-only info on probed routes; print has ZERO styling → TODO #424. |
| M12 (μ46–μ49) | `visualtest/tools/sseload` (Go stand-in for k6+xk6-sse): steady 10×20 s = 2000.05 ms ±8 ms, **0 continuity violations**; reconnect 10× = 2002–2012 ms each; 50-client synchronized drop = handlers reaped, service healthy (2005 ms). **Verdict: no nightly lane now;** kill-criteria drafted. §5.5. |
| M13 (μ50–μ53) | `nix run .#crossbrowser` flake app (pinned browsers, SAME pure fontconfig pin as visual, demo binary, 3 engines) — works end-to-end. Determinism: **0/13,949,440 pixels differ** across runs (byte diffs = PNG encoder chunks only). CI lane DRAFT committed (`workflow_dispatch`-only, inert). actionlint clean. |
| M14/M15 (μ55–μ60) | §6 per-tool verdict table + ADR 0048 boundary (ADDITIVE; chromedp primary). Gates prepared to the edge, not self-ratified. |
| M17 (μ63–μ64) | stylelint 17.16 + TW4 config: **0 findings** on custom.css + demo.css → DEFER. rustywind: would rewrite 183 files INCLUDental goldens → **REJECT**, logged. |
| M18 (μ65–μ66) | Trace recorded via playwright-core `context.tracing` (no test runner needed). KEEP as debugging skill, not a lane. |
| M19 (μ67–μ69) | fontconfig pin in the app; determinism proven; below-`xl:` viewport caveat documented. |
| M20/M21 | Gate-edge schemas documented (achecker ledger semantics; aria-golden promotion path) — wiring awaits ⫱G2/⫱G1. |
| M23 (μ75–μ77) | `docs/browser-support.md` support matrix ADOPTED (base-select/field-sizing/dialog/popover/CQ/content-visibility/SSE rows, evidence-linked). |
| M24 (μ78–μ79) | `docs/testing/tool-trial-protocol.md`: standing protocol + CI budget table (owner + kill-criteria per lane) + rejected-tools log (Lost Pixel, rustywind, k6-local, second pixel suite, Guidepup-parked). |
| M25 (μ54, μ80–μ81) | TODO_LIST: #421–#424 registered, execution-status note on #411–#420; AGENTS.md Testing bullet (189/220 lines). |
| M26 (μ82–μ83) | ROADMAP: 9 backlog seeds (parity contracts, keyboard extension, contract tests, sampling economics, chaos, golden-churn, re-verify reminder…). |
| μ84 (partial) | Plan verification checklist ticked with execution annotation; nix fmt clean; root + per-module tests green (after inventory self-heal). |

## b) PARTIALLY DONE

1. **μ84 final gate:** tests green, but `scripts/ci-repro.sh --lint --css --website`
   NOT yet run and **push NOT done** (master ahead 5). Deliberately paused: the
   concurrent session is mid-edit on `layout/appshell.templ`/`templates/custom.css`
   and their 8 route goldens still fail.
2. **M13 CI lane:** file committed but dispatch-only + one UNVERIFIED action SHA.
3. **M20/M21:** schemas/paths documented; zero wiring (gates).
4. **WebKit `base-select` support:** matrix says "verify at ⫱G1" — M07 assert only
   ran chromium+firefox (webkit leg not attempted; would have been cheap).
5. **Eval-doc cross-links:** the plan file's execution note points at the OLD
   status doc (17:10) as §g; this report supersedes it — the pointer text says so,
   but a follow-up edit after ratification should relink cleanly.

## c) NOT STARTED

- Push + gate report delivery to origin (blocked on ritual + coordination, §g).
- All ⫱G-gated execution: enabling the crossbrowser CI lane, achecker ledger
  wiring, aria-golden class wiring, k6 nightly (recommended: no).
- M15's alternative branch (chromedp→Playwright migration plan) — moot under the
  ADDITIVE recommendation unless ⫱G3 flips.
- ROADMAP seeds (parity contracts, chaos profiles, golden-churn tracking, …) —
  intentionally backlog.
- WebKit leg of M07 (above) and dark-theme cross-engine captures (spike ran
  light only).

## d) TOTALLY FUCKED UP

Nothing catastrophic — but one incident deserves the class of "would have bitten
a consumer": the npm-pinned `playwright-core` 1.64.1 got daemon-committed
(breaking the browser-revision lockstep; `nix run .#crossbrowser` would fail for
anyone cloning at that commit). Caught within the session, restored to 1.63.0,
verified the app still runs. Root cause: installing a spike dep INTO the pinned
manifest + trusting npm `--no-save` semantics + a daemon that commits fast. The
improvement (§e-4/§e-6) closes it permanently.

## e) WHAT WE SHOULD IMPROVE

1. Run the full ci-repro ritual and push once the tree stops moving (or now, if
   the concurrent session's lane can tolerate it).
2. Verify (or replace) the `nix-installer-action` SHA in the draft lane.
3. Add `TestCrossbrowserLockstep`: parse spike `package.json` `playwright-core`
   vs `nix eval --inputs-from . nixpkgs#playwright-driver.version` — cheap guard
   against a repeat of §d.
4. Run spike-foreign npm trials (achecker, stylelint) in /tmp scratch dirs, never
   inside `tools/crossbrowser-spike/`.
5. Fold `stylable-assert.mjs` + `aria-snapshots.mjs` into the `.#crossbrowser`
   app (one command = M06+M07+M10 evidence refresh).
6. Add a chromedp e2e click-test for `data-tc-overlay-trigger` in the PRIMARY
   suite — the 3-engine click proof currently lives only in the spike.
7. Re-check `-parallel` and load-flake behavior of the refreshed `/display`
   golden after the concurrent session lands (their CSS may shift it again).
8. Replace demo `picsum.photos` avatars with locally generated images — removes
   the `img-src` CSP exception AND a network flake source.
9. Quote the 3 CSP before/after reports (small JSON/text) into
   `docs/testing/evidence/` so the ADR trail survives without the gitignored
   PNGs.
10. Plan-authoring checklist: add "verify the env before trusting green (silent
    SKIP class)" — this session's cheapest lesson with the best payoff.

## f) NEXT 50 (ordered: unblock → ratify → pay down → explore)

1. Commit `templates/templ-components-classes.txt` if regenerated again at verify time.
2. `scripts/ci-repro.sh --lint --css --website` at the exact commit.
3. Push master (5 commits) — after #26 coordination answer.
4. Verify/replace `DeterminateSystems/nix-installer-action` SHA in `crossbrowser.yml`.
5. Add `TestCrossbrowserLockstep` (playwright-core ⇄ nixpkgs playwright-driver).
6. Pose/ratify ⫱G1 (crossbrowser lane adopt).
7. Pose/ratify ⫱G2 (achecker trial adopt).
8. Pose/ratify ⫱G3 (additive confirmed — likely no-op).
9. Pose/ratify ⫱G4 (SSE nightly: recommend NO).
10. On ⫱G1: enable `crossbrowser.yml`, add budget-row sign-off (M24 table).
11. On ⫱G1: fold aria-golden promotion (M21) — `testdata/aria/` class wiring.
12. On ⫱G2: wire `visualtest/testdata/achecker_baseline.json` + one documented finding flow (M20).
13. On ⫱G3: countersign the ADR boundary section.
14. #421 CopyButton `role="status"` outside the button (+ goldens + a11y test).
15. #422 `layout.Split` aside landmark decision (AsideLabel prop vs div).
16. #423 Table wrapper `role="region"` + accessible name (+ goldens).
17. #424 minimal `@media print` layer + consumer docs.
18. WebKit leg for `stylable-assert.mjs` (close the matrix's "verify at ⫱G1" cell).
19. Dark-theme 3-engine captures (spike ran light only).
20. chromedp e2e click-test for the Modal/Drawer triggers (browser-proof in-suite).
21. Replace picsum avatars with local generated images (CSP + flake removal).
22. Quote CSP evidence JSON into a committed docs/testing/evidence/ file.
23. Fold M07+M10 scripts into `.#crossbrowser` (`--all` mode).
24. Document `sseload` in a short README (usage, profiles, cadence flag).
25. sseload: fix gap attribution (currently round-robin per client) for precision.
26. sseload: add a server-RESTART reconnect profile (exercises `Retry:'always'` semantics).
27. Track Mozilla 1944403/1958443; re-run M07 assert when Firefox ships base-select.
28. Register the 6-month re-verification reminder with a concrete due date (2027-04).
29. Re-run the full visual suite after the concurrent session's lane lands; coordinate rebaseline of the 8 failing goldens (their lane).
30. Re-verify my `/display` golden after their CSS merge (may shift again).
31. Consider `data-tc-overlay-trigger` FEATURES.md row + docs snippet (new consumer feature).
32. Evaluate deprecating `tcOpenModal/tcOpenDrawer` wrappers in v2 notes (attribute supersedes).
33. Add a `.achecker.yml` config into the spike dir for repeatable policy (currently default policy).
34. Pin/document the puppeteer-chromium revision achecker uses (reproducibility row in README).
35. Check dependabot now sees `tools/crossbrowser-spike/package.json` (npm ecosystem block) — decide in/out.
36. Decide whether `package-lock.json` (untracked) should be committed for the spike (reproducibility vs drift noise).
37. Confirm the `/stylable-select` axe ledger entry survives the concurrent session's contrast work (re-read when their lane lands).
38. Extend `keyboard_traversal_test.go` routes (ROADMAP seed) — cheap first extension: `/stylable-select`.
39. htmx-fragment HTTP contract tests (ROADMAP seed) — start with `/demo/api/wire/*`.
40. Golden-churn tracking script (ROADMAP seed) — count golden commits per component.
41. Sampling-economics policy note before ⫱G1 integration (runtime budget).
42. Datastar chaos profile in sseload (ROADMAP seed) — only after kill-criteria trigger.
43. Lighthouse revisit criteria note (ROADMAP seed) — tie to website CWV work.
44. Clean /tmp scratch (pw-browsers, stylelint-trial, cb-det, tc-demo binaries) or document bootstrap in the spike README (partially there).
45. Add the "verify env before trusting green" lesson to `docs/plan-authoring-checklist.md`.
46. Update the plan file's §g pointer to THIS status doc after ratification.
47. Cross-check `docs/testing/demo-coverage-map.md` against any concurrent-session demo changes.
48. Re-run `TestDocsCountDrift -count=1` after the concurrent session merges (both lanes touch counts).
49. Consider a tiny `docs/testing/README.md` index linking evaluation/protocol/browser-support/coverage-map (discoverability).
50. Close the loop: after gates are ratified, annotate ADR 0048 Status → Accepted with the gate outcomes + dates.

## g) THREE QUESTIONS (owner-only)

1. **Push coordination:** master is 5 ahead with my work committed, the ritual
   hasn't run, and the concurrent session is mid-edit on
   `layout/appshell.templ`/`templates/custom.css` (their 8 route goldens still
   failing). Push my commits NOW (running the full ci-repro ritual first), or
   hold until their lane lands so the ritual runs once over the combined tree?
2. **⫱G1:** adopt the crossbrowser diagnostic lane as drafted (flake app now +
   dispatch-only CI lane, budget row per M24)? If yes, enable the lane
   immediately or after the sampling-economics note?
3. **Duplicate-id contract:** should the forms components auto-uniquify
   `name`-derived ids (EnsureID when no explicit ID — changes selector
   expectations, breaking-ish) or is "consumers own id uniqueness; demo must set
   explicit ids" (what I implemented) the intended long-term contract to document
   in the forms docs?
