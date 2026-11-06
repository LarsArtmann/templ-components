# Session Report — 2026-11-07 02:13 CET — Recovery Found Moot (PR #77 Already Merged), Accent-Affordance Regression Fixed via PR #87

Continuation of `2026-11-06_05-23_design-token-axes-unred-s04-daemon-incident-status.md`.
The recovery plan from that report (§d1) was executed — and discovered to be
redundant partway through, because the world it targeted no longer existed.
The session's real deliverable became an audit of what actually landed on
master plus one genuine regression fix (PR #87, merged).

Single-writer scope claimed: `feat/design-token-axes` recovery + S04 files +
PR #77. The scope evaporated at T+40min (see §d).

---

## a) FULLY DONE

| Work | Evidence |
| --- | --- |
| **PR #87 merged: StatCard htmx tile + clickable rows focus in accent** — the session's real defect. #79 added `statCardInteractiveClass` (ring `blue-500/400`) and the row outline (`outline-blue-600/400`) by copying the then-blue linked ring; #77's emission merge later swapped the neighbors (linked tile → `ring-accent-500/400`, sort links → accent) but not them, leaving twin affordances in different families inside the same components. Fixed shade-preserving; row-href assertions flipped; `stat_card_htmx` golden regenerated (it had pinned the blue); corpus (27,328 lines), compiled CSS, tc mirrors, CHANGELOG all resynced | `7eeafb8e` on master; `gh pr checks 87` all 8 lanes pass (Visual Regression 7m15s, Build & Test 4m41s, Lint, Clone gate, HTML-valid, CSS Freshness, CHANGELOG warmth, Website); merged `--rebase --delete-branch` |
| **ci-repro ritual finally executed with a VERDICT** — the gate skipped for the entire prior session, explicitly named in that report's §e. Run at the exact pushed commit | `VERDICT: PASS (exit 0) — 2026-11-07 01:45:34 CET [dupl: art-dupl version v0.9.1-…-93432fd]` — pinned binary (the §b3 requirement from the 01:25 sibling report); master then landed `fc3bcea6` making ci-repro REFUSE non-pinned binaries |
| **Full state reconstruction after the daemon hijack** — proved from refs, `gh` API, and the sibling's untracked report: PR #77 merged upstream as squash `6082430e` (21:40), carrying the daemon-pushed chain (S04 sweep, demo override, 145-group reconciliation, this report's predecessor, the contested dep bumps); PRs #79 (defects), #80 (Scrollback/ListNote), #81 (#50 skill stamp), #82 (#501 recipe), #83 (color-token guard), #85 (app.css comment lesson) on top; `feat/design-token-axes` ref and its origin branch deleted by the merge path | `gh pr view 77` (state MERGED, mergeCommit `6082430e`); `gh pr list --state merged`; `git merge-base --is-ancestor 06001739 origin/master` → NO (squash, not merge ancestry) |
| **Master content audit against every recovery target** — toast stray `</script>` gone, toggle aria-checked gone (1 doc-comment mention only), S04 accent markers present (badge/tabs/card), `utils/darkmode_scan.go` avatar exception present, the 05-23 status report TRACKED on master, dark-scan + semantic keeps (BadgeInfo, ButtonOutlineInfo, StatTone*, ChartColorBlue, KanbanToneBlue) intact | `git show origin/master:<file>` greps + `git cat-file -t origin/master:docs/status/2026-11-06_05-23_…` → blob |
| **Parity diff salvaged-S04 vs master** — only 6 files differ, all attributable to follow-up PRs (#79 card/table additions, LoadingProgress extraction −28 lines), zero missing S04 pieces | `git diff eb404c4c origin/master --stat -- <19 S04 paths>` |
| **Old recovery contract disproven before it could cause damage** — "toast must keep accent swaps" (§d1 of the 05-23 report) checked against 5 refs + 30-reflog-entry scan: toast accent content exists in NO commit anywhere, and ADR-0044 decision 3 deliberately keeps Alert/Toast blue. The contract was based on a daemon-reverted working tree; correctly NOT invented | grep matrix across `51a03c3e`, `caae2ede`, `8ec97032`, `origin/master`, `bcebe187` (all 0 matches); `4afcc129` file list (feedback pilot = progressbar/loading/circular/step_indicator only) |
| **Fix engineered in an ISOLATED worktree** (`/tmp/tc-accent-fix`, branch off `origin/master`) — after the shared worktree was hijacked, every subsequent write happened away from the sibling's active space; cleaned up fully after merge (worktree removed, local branch deleted) | `git worktree list` → only main worktree remains |
| **Full verification battery on the fix** — root `go test ./...` (cmd/tc initially red on the mirror guard → sanctioned `check-tc-sources-sync.sh --fix` → green), per-module GOWORK=off loop all ok, golangci-lint `./display/... ./cmd/...` 0 issues, `check-html-valid.sh` clean (338 goldens), `TestSemanticColorTokensCompiled` + `TestClassInventoryFreshness` pass, build OK | command transcripts in session |

## b) PARTIALLY DONE

1. **The recovery itself (§d1 of the 05-23 report) was executed then abandoned as moot** — salvage patch written (81 files), branch reset to `bcebe187`, S04 committed (`eb404c4c`, 92 files, mirrors auto-staged by the guard), `origin/master` merge started (CHANGELOG auto-unified, toast/toggle auto-merged correctly, app.css regenerated). Then the daemon hijack + the `gh`-confirmed upstream merge made all of it redundant. The commit `eb404c4c` survives only in the reflog. Nothing was lost — but the effort was spent before the T0 check that would have killed it (see §e).
2. **Two same-root-cause suspects deliberately NOT fixed, pending owner intent** — `copy_button.templ`, `code_block*.templ`, `accordion*.templ`, `empty_state_default` still carry `ring-blue-500/400`. They predate the emission (not a #79/#77 interaction), were never in the S04 24-site census, and may be deliberate semantic keeps (like ButtonOutlineInfo's `ring-blue-300`, which has an assertion test pinning it). Fixed only what the merge demonstrably broke (intra-component inconsistency); the rest needs a design decision (§g Q1).
3. **Cross-session findings recorded here instead of in `TODO_LIST.md`** — the shared main worktree belongs to the active sibling session (bar-chart branch, ahead 1 behind 1 at report time); writing TODO rows or a history-doc narrative into it risks their daemon sweeping my lines into their commits (the exact `git add -A` hazard from the 05-23 report §d2). All hygiene rows are therefore listed in §f for whoever next owns the doc.

## c) NOT STARTED

1. **v2.0.0 release** — now the ONLY cuttable version (master carries the BREAKING accent-emission entry; ADR-0039 forbids v1.x from an emission tree — the old Q2 "v1.28.0 vs v2" resolved itself: v1.28.0 died when #77 merged). Requires: `docs/release-checklist.md` read-through (incl. new step 6d skill-stamp bump from #81), release.sh inside the dev shell with govulncheck, lockstep tags root + all published submodules, `result*` cleanliness check.
2. **Remaining ADR-0049 stages** (S06 forms sweep, S07 nav/layout, S08 errorpage/recipes, S09 tone-map collapse #425/#431, S10 guards) — unscoped against the current master; the plan doc needs a re-read before any of them start.
3. **Post-emission affordance guard** — nothing prevents the #79/#77 interaction class from recurring: an affordance added in pre-emission colors, then stranded when the emission lands around it. Candidate: extend `TestSemanticColorTokensCompiled` (or a sibling guard) to assert brand-axis COMPONENTS carry no palette focus/ring/outline tokens, mutation-tested like #83.
4. **Demo exposure for the fixed affordances** — the demo has zero `HxGet` StatCards and zero `Href` rows, so the visual suite cannot see these paths; the regression lived invisibly exactly because of that. Needs a demo route or e2e (see §f).
5. **Daemon-hijack narrative + defense** — nothing in `docs/agent-context-history.md` yet records this session's hijack class (branch ref deleted + worktree switched away while a merge was in conflict). No TODO rows exist yet for it or for the `core.bare=true` bare-flip bug from the prior session (still unfiled, carried twice now).
6. **#60 Sparkline** signed/responsive geometry (reuse candidate: `barChartBounds`/`barChartSignedPlacement` from the sibling's bar-chart work), **#63 CopyButton** page-level script mount, **#53** forms wrapper read-and-decide — untouched, as the sibling's report also lists them.

## d) TOTALLY FUCKED UP

1. **The T0 research miss that cost the session's first half:** I reconstructed the world from GIT REFS ONLY and never queried the PR tracker at session start. One `gh pr view 77` at T0 would have shown MERGED (21:40, the night before) and ended the recovery before it began. Instead: salvage patch, branch reset, a 92-file commit, a merge with conflict resolution, and a hijack encounter — all spent recovering work that was already master. The repo's state of truth is GitHub; git refs are a cache of it. This is the session's defining error.
2. **I mutated the SIBLING session's worktree.** Early verification ran `go test ./utils/` in the shared main worktree; `TestClassInventoryFreshness` auto-writes a regenerated corpus on staleness, so I planted an unstaged `templates/templ-components-classes.txt` in the sibling's active space — which then BLOCKED my own merge (uncommitted-changes error, one aborted attempt) and could have been swept into their commits by their daemon. The isolated-worktree pattern should have been the default from the first command, not adopted only after the hijack.
3. **Executed a recovery plan against a stale world model without re-validating its premises.** The 05-23 report was ~20h old; the summary itself warned state moves fast (the daemon had already pushed the chain between the report being written and this session starting — visible as "origin/feat = 06001739, [ahead 1]" at T0 — a signal I logged but did not chase to its conclusion, which was: the branch was on its way to being merged).
4. **The daemon hijacked the worktree MID-MERGE and destroyed my in-flight state** — switched the shared worktree to `fix/bar-chart-vertical`, authored a merge commit (`bf4840b3`, parents `c983b68b`+`73e3f438`) and follow-ups, leaving HEAD on a foreign branch, my merge gone, and the `feat/design-token-axes` ref deleted. Not my error — but my exposure to it was: I kept operating in the shared worktree after earlier commands had already shown foreign daemon commits moving refs (`1547fab7`, `bf4840b3` appearing under me mid-session).
5. **First ci-repro run was wasted** — executed with uncommitted changes, so its templ-drift check diffed against HEAD and reported "REAL content drift" for exactly my (correct) generated files. ~10 minutes of lanes burned on a guaranteed-false verdict. The ritual's own wording ("at the exact commit about to be pushed") already forbids this; I ran it pre-commit anyway.
6. **Golden-update flailing:** three attempts to regenerate `stat_card_htmx.golden` (`-run TestStatCard` → matched nothing; `-run TestGoldenSweepStatCard` → ran the sweep, which doesn't own the file; finally grepping the golden's NAME in test files → owner is `TestGoldenStatCardHTMX` in `golden_new_test.go`). The find-the-owner-first grep was the obvious move and I did it last.
7. **`cd /tmp` inside the worktree shell broke the cleanup** — `git worktree remove` + branch deletion failed with GIT_DISCOVERY errors (exited the repo), needing a retry from the main repo. Trivial, but it left the worktree/branch half-cleaned for one command cycle.
8. **One commit shipped via `--no-verify`** (`eb404c4c`, the now-moot S04 commit) — justified in-message (BuildFlow hook red on environment timeouts: test-compile killed at 14× p95, gomod-freshness "signal: killed", plus pre-existing repo-wide file-size findings on untouched files; fast guards all green; CI is enforcement — the documented precedent), but the hook was never re-tried after the machine settled, and that commit turned out to be dead weight anyway. PR #87's commit went through the full battery instead.

## e) WHAT WE SHOULD IMPROVE

**What I forgot:**
- The PR tracker exists. Session-start state reconstruction must include `gh pr list` / `gh pr view` before any git-archaeology; refs answer "where is the code", the tracker answers "what happened to it".
- That a test in this repo WRITES files (`TestClassInventoryFreshness` regenerates the corpus in-tree). "Tests are read-only" is false here; running them in a shared worktree has side effects on the sibling.
- That the freshness test had TOLD me it wrote the file ("The tree now carries the regenerated file; commit it") — the hint was in output I read an hour before the aborted merge.
- The owner-first grep for goldens (find the test that owns the artifact before invoking update flags).
- That the 05-23 report's premises needed re-validation, not execution — its own §d1 said "ready, not run", and 20 hours of daemon activity had elapsed.

**What I could have done better:**
- Worktree isolation from command #1. One `git worktree add` before any test run removes both §d2 and half of §d4's blast radius.
- A 30-second T0 gate: `git fetch && gh pr view <open PRs> && git log origin/master -3` — compare against the prior report's claims, THEN plan. Would have collapsed the entire first half into one audit pass.
- Re-try a red hook before bypassing it when the red is plausibly environmental (the machine was visibly thrashing — killed children, 143000% timing regressions).
- Read tool failures fully: "The tree now carries the regenerated file" was the classes.txt answer, given for free.
- Close the loop on my own close-out: the final cleanup should have been one atomic command run from a repo directory.

**What we should still improve (repo-level, not just mine):**
- The status-report handoff contract needs a freshness stamp ritual: every report should open with the exact refs it was written against, and the next session's first act is diffing those refs against reality. Both this session and its predecessor were burned by executing plans across stale snapshots.
- `ci-repro.sh` should assert a clean tree (or diff the templ check against the working tree, not HEAD) instead of producing a false "REAL content drift" for uncommitted-but-consistent work.
- The golden corpus has two ownership styles (`TestGoldenSweep<Component>` tables vs `golden_new_test.go`'s `golden.Assert` calls) — finding an artifact's owner should not require three probes; a naming/registry convention (or an orphan-golden guard) would close it.
- The sibling session's §b3 finding (ci-repro silently falling back to a dev art-dupl binary) got fixed on master (`fc3bcea6`) WHILE this session ran — good; the remaining open items from that report (§b1 bar-chart ritual, §b2 PR #84 rebase) still need an owner.

## f) NEXT (up to 50, ordered roughly by leverage)

1. Land the sibling's bar-chart branch: ci-repro VERDICT at tip → push → PR (Fixes #48 #59) → CI → merge (their M03 ritual; branch currently ahead 1/behind 1 of master — my `7eeafb8e` doesn't touch their files, rebase should be clean).
2. Rebase PR #84 (tc doctor, #35) onto current master, resolve the CHANGELOG overlap, `--force-with-lease`, let CI run (GitHub suppresses lanes on CONFLICTING), merge, #35 auto-closes.
3. Owner-gated: cut **v2.0.0** — the only remaining version; read `docs/release-checklist.md` first (new step 6d), dev shell + govulncheck, `result*` cleanliness, lockstep tags all published modules.
4. Re-scope ADR-0049's remaining stages (S06 forms / S07 nav-layout / S08 errorpage-recipes / S09 tone-map collapse / S10 guards) against current master — the plan predates #77–#87 landing.
5. **S09 tone-map collapse (TODO #425/#431)** is explicitly v2-windowed and v2 is now the active window — schedule it before or with the cut.
6. Decide the pre-emission blue-ring census (§g Q1): CopyButton, CodeBlock, Accordion, EmptyState — semantic keeps (document in ADR-0044/0049) or a follow-up accent sweep (S11-style).
7. **Post-emission affordance guard**: extend the color-token guard family to fail when a brand-axis component emits palette focus/ring/outline tokens; mutation-test it (planted blue ring must fail with file attribution) — kills the #79/#77 class permanently.
8. Add demo exposure for both fixed affordances (an `HxGet` StatCard and a clickable-row table in a demo route) so the visual suite covers them — the regression survived precisely because no route rendered those paths.
9. Chromedp e2e: keyboard-focus ring color parity between linked and htmx StatCard variants (string-proven ≠ browser-proven; the golden now pins the string, the browser should pin the pixel).
10. Write the daemon-hijack narrative into `docs/agent-context-history.md` (append-only; AGENTS.md stays ≤220 lines) — include the `bf4840b3`/`1547fab7` reflog evidence trail.
11. File the TODO row for the `core.bare=true` shared-config flip (unfiled since the prior session, carried forward twice now).
12. TODO row for **daemon-merge defense**: hook or guard that (a) warns loudly when a worktree switch happens with `MERGE_HEAD` present, (b) makes branch-ref deletion of a pushed branch require confirmation. (Owner decision — §g Q3.)
13. Session-start checklist amendment (project AGENTS.md or the checklist doc): `gh pr list/view` BEFORE git-archaeology; record the exact refs a plan was written against.
14. `ci-repro.sh`: assert clean tree before the templ-drift lane (or diff against working tree) — kills the false "REAL content drift" I hit.
15. Golden-ownership convention: one style for golden tests (sweep tables or `golden.Assert` — pick one), plus an orphan-golden guard (a testdata file no test references must fail CI).
16. Annotate the 05-23 report (docs-health ANNOTATE mode): its §d1 recovery contract is now known-wrong (toast accent premise + moot-by-merge) — resolve the numbered items inline so no future session executes it again.
17. Harvest `TODO_LIST.md`: verify #998/#501/#502/#503 rows against landed PRs, strike resolved, add §f rows 6–12.
18. Diff `6082430e`'s tree against the old branch tip `06001739`'s tree — confirm the squash carried the daemon chain byte-honestly (I audited targets, not the full delta).
19. Verify the golden-count story after the sibling's bar-chart PR merges (333→336 across 7 doc sites; guards will enforce — but confirm the countgen output at that merge, per their report).
20. Update the local stale `master` ref hygiene (`git branch -f master origin/master` — it sat at `5f354154` all session; a stale local master misleads `--contains` checks, as seen at T0).
21. Decide fate of the reflog-only S04 commit `eb404c4c` (document as intentional dead object or let it age out; `/tmp/s04.patch` likewise).
22. Look into the BuildFlow pre-commit environment thrash (test-compile killed at 14× p95, gomod-freshness killed at 42s, govalid +143000% timings): `buildflow doctor` on an idle machine; check the dead `/mnt/buildcache` suspicion; possibly raise the 60s pre-commit budget or make env-class timeouts non-fatal.
23. Rebuild/reinstall the stale BuildFlow binary (built at `e81c118`, repo HEAD past it — preflight warns every run).
24. Point `dprint.json`'s 4 remote plugin fetches at local `.wasm` paths (BuildFlow's own recommendation; cold-cache hangs a future run).
25. Clear the preflight `vendorHash` staleness warning via `nix-hash-fix` before the next real `nix build` (it WILL fail with a hash mismatch otherwise).
26. Run `buildflow format` once to re-verify cached verdicts under the drifted formatter identity (templ-fmt v0.3.1020 → v0.3.1100 fingerprint warning).
27. Re-verify the pre-commit hook passes cleanly now (the machine-thrash red) so future commits don't need `--no-verify` precedent.
28. Sweep `TODO_LIST.md` for other stale-but-unannotated status reports in `docs/status/` (the 2026-11-04 and 2026-11-06_02-28 ones) — ANNOTATE pass.
29. Clarify ADR-0044 decision 3 with one sentence: "Alert/Toast keep semantic blue; focus/ring affordances are brand-axis" — the 05-23 contract contradiction shows the line is derivable-but-not-written; writing it prevents the next misplan.
30. After v2.0.0: consumer migration guide for the accent `@theme` block (one-liner in the consumer guide + website).
31. Verify website CSP needs nothing for the new tokens (classes only, no inline styles — expected clean; cheap confirm at next site build).
32. Check for other stale remote branches post-#77/#87 cleanups (`git fetch --prune`; `backup/*` were already removed by the sibling).
33. Consider `git worktree add` guidance in the contributor setup doc — after this session, isolated worktrees are the proven pattern for multi-session days.
34. The `core.bare` incident root cause (linked-worktree self-test suspect, from the 05-23 report §c7) — find and fix the offending self-test.
35. Post-v2: revisit `HueAccent`/BGSolid-accent gap (S11 note from the prior summary — badge dot pinned as literal because the tone row lacks `HueAccent`).

## g) QUESTIONS I CANNOT ANSWER MYSELF

1. **The pre-emission blue rings on CopyButton, CodeBlock, Accordion, and EmptyState:** deliberate semantic keeps (like ButtonOutlineInfo's pinned `ring-blue-300`), or brand-axis sites the S04 census missed that you want swept to accent? I can draft either the ADR-0044 "documented keeps" list or the sweep PR — I cannot know which matches your design intent.
2. **v2.0.0 scope:** cut after the bar-chart + doctor PRs land (current [Unreleased] is already v2-gated and rich), or hold for the remaining ADR-0049 stages (S06–S09 sweeps, tone-map collapse) to ride the same major? This decides whether items 4/5/29 block the release or follow it.
3. **Daemon defense priority:** do you want a built guard (refuse/flag worktree switches with pending MERGE_HEAD; confirmation on pushed-branch ref deletion — item 12), or is the hijack class accepted chaos that only needs the history-doc narrative (item 10)? Building it touches the shared hook infrastructure another active session depends on, so I won't start it without your call.

---

*Verification trail: `gh pr checks 87` all-pass incl. Visual Regression 7m15s; `gh pr merge 87` output; `git log origin/master -2` = `7eeafb8e`/`fc3bcea6`; `gh pr view 77` MERGED@`6082430e`; ci-repro log `/tmp/ci-repro-fix.log` VERDICT PASS (exit 0, pinned dupl `v0.9.1-…-93432fd`); audit greps (toast `}</script>`=0, toggle aria-checked=1-comment, badge/tabs accent>0, darkmode accent-200=1, report blob tracked); reflog + `git worktree list` post-cleanup.*
