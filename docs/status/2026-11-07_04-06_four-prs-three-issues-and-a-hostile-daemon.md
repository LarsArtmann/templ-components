# Status Report — 2026-11-07 04:06 CET

**Session scope:** resume of the standing issue/PR sweep. Prior session left two open
branches (BarChart unpushed, tc-doctor PR #84 conflicting) and a broken dupl-pin ritual.
This session landed FOUR PRs, closed FIVE issues, survived six daemon sabotage events,
and ended mid-flight on issue #63 (implemented + verified, not pushed).

---

## a) FULLY DONE

| Work                                                                                             | Evidence                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| ------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **PR #86 — ci-repro refuses art-dupl results from non-pinned binaries**                          | merged `fc3bcea6`; the dupl advisory lane now ignores any binary whose version lacks the pinned rev, auto-builds via `scripts/pinned-art-dupl.sh`, and reports explicit `dupl-skip`; live-fire tested with a fake dev binary on PATH before the PR; every verdict since carries `[dupl: art-dupl version v0.9.1-…-93432fd]`                                                                                                                         |
| **Pinned art-dupl rebuilt**                                                                      | `/tmp` had wiped it; rebuilt at pinned rev `93432fd` via `scripts/pinned-art-dupl.sh`; clone gate GREEN — 0 new groups vs baseline 145 — re-verified at every branch tip                                                                                                                                                                                                                                                                            |
| **PR #89 — vertical BarChart measures, signs, scales** (was sitting unpushed from prior session) | merged `0bf45923`; rebased twice onto moving master (CHANGELOG + app.css conflicts resolved, artifacts regenerated at tip, NOT hand-picked); **#48 auto-closed**, **#59 closed manually** with evidence comment (PR body only carried `Fixes #48`); full ci-repro VERDICT PASS 01:57 + 02:04 at the pushed tips                                                                                                                                     |
| **PR #84 — tc doctor reads real go.mod/CSS layouts** (was CONFLICTING, CI suppressed)            | rebased onto new master, CHANGELOG conflict resolved, force-with-lease push (PR head was attached — deleting would have detached it, the #88 gotcha), all lanes green, merged `e75257a2`, **#35 auto-closed**                                                                                                                                                                                                                                       |
| **PR #90 — Sparkline responsive + zero-aware (#60)**                                             | merged `f6891f9f`; three opt-in props (`Responsive` viewBox+w-full+vector-effect, `ZeroAware` domain widening with Min/Max still winning, `ShowZeroLine` rule implying ZeroAware); 3 HTML goldens; **two chromedp browser proofs PASS** (responsive svg stretches into a 600px parent, fixed control stays 120px; zero line strictly separates signed series on screen); count docs same-commit (339 goldens, ~1,900 subtests); **#60 auto-closed** |
| **Daemon incident class discovered + recorded**                                                  | `docs/agent-context-history.md` §duplpin: the daemon created a MERGE commit carrying **app.css with literal conflict markers** (277,587 bytes = both sides concatenated) and **zeroed go.mod sibling pins** (`v0.0.0-00010101000000` — the exact state M18 + `check-replace-directives.sh` reject); recovery recipe documented (mixed reset → restore go.mod → regenerate both artifacts at tip → soft-fold)                                        |
| **Housekeeping**                                                                                 | fan-out verified through the symlink path (stamp line 38 correct); all merged branches deleted locally+remotely; master CI green on the BarChart/doctor pushes (doctor's master run was in_progress at last check — content identical to its green PR run)                                                                                                                                                                                          |
| **Issues closed this session**                                                                   | #48, #59, #35, #60, #66 (prior), plus #50/#81 via earlier merges                                                                                                                                                                                                                                                                                                                                                                                    |

## b) PARTIALLY DONE

1. **#63 CopyButton page-level script mount — implemented + verified, NOT pushed.**
   Branch `feat/copybutton-script-mount`, commit `605b2978` on master `f6891f9f`:
   `CopyButtonScript(nonce)` exported page-shell mount; `CopyButtonProps.NoScript`
   (inverted flag per repo convention — zero value = current behavior, non-breaking);
   default path **byte-identical goldens** (the first attempt used an `if` wrapper in
   markup, which leaked a space before `<script>` — goldens caught it; fixed with a
   `templ.NopComponent` selection in the code block, the repo's existing pattern);
   joined all three nonce guards (csp_nonce, no-empty-nonce, empty-nonce-still-renders);
   mirror synced; display+integration tests green; lint 0 issues; **full ci-repro VERDICT
   PASS at 02:55:48 CET with the pinned dupl binary**. Blocking remainder: one more
   daemon commit (`abfc9dc5`) landed on top carrying a **+27-line classes.txt delta I
   cannot yet explain** — must be regenerated-and-diffed before push; then fold, push,
   PR (`Fixes #63`), merge.
2. **Transient ci-repro FAIL on the CopyButton branch — root cause never proven.** First
   full run failed right after the utils module (output truncated, failing package not
   captured). Root tests re-ran green, website lane green, coverage gate green (96.5% ≥
   94), full re-run VERDICT PASS. Suspicion (unverified): toolchain skew — the run
   happened OUTSIDE nix develop and the coverage lane itself warned "baseline recorded
   with go1.27.0, this run used go1.27.3".
3. **Coverage baseline toolchain skew** (pre-existing, now visible): baseline is
   go1.27.0-recorded; non-dev-shell runs are go1.27.3. Numbers remain toolchain-specific
   per AGENTS.md; re-record on the canonical patch level.
4. **Master CI on the doctor merge**: last observed in_progress; the run content equals
   its all-green PR run, but the final master-run verdict was not re-read before the
   session pivoted.

## c) NOT STARTED

- **#53** (forms.Input/Select unwrapped-sibling layout hostility): assess only — prior
  plan says likely Deferred-v2; the issue itself offers a non-breaking additive
  `Stacked bool` option, so the call is genuinely open (question 3 below).
- **v2.0.0 release ritual** — still user-gated; ADR-0039/0049 policy, release.sh does
  not push.
- **TODO 502** (BuildFlow gate debt → the `--no-verify` standard that forces manual
  mirror-sync/artifact runs) — untouched.
- Demo-page variants for the new Sparkline/CopyButton modes (deliberately skipped to
  avoid visual-route golden churn — precedent: the BarChart PR also skipped demo
  changes).

## d) TOTALLY FUCKED UP (and what it cost)

1. **I violated the documented templ gotcha TWICE** (`templ generate ./display/` from
   repo root): first during sparkline, then again during CopyButton — each time
   rewriting ~29–31 submodule generated files with wrong workspace-relative FileNames,
   each time requiring a restore sweep. The gotcha is written down in AGENTS.md; there
   is no excuse except speed-pressure from the daemon races. Correct form (used for the
   final recovery every time): plain `templ generate ./...` from root + restore
   submodule files (they are already canonical).
2. **Forgot the self-healing mirror sync for `display/sparkline.templ`** — surfaced as
   `TestSourcesMatchPackageFiles` FAIL deep inside a ~8-minute ci-repro run, wasting a
   full ritual cycle. Under `--no-verify` the self-staging guards never run; the manual
   `check-tc-sources-sync.sh --fix` must be part of MY checklist, not the hook's.
3. **Pushed a non-verified tip once**: after the doctor VERDICT PASS at `b9908742`, the
   daemon raced commit `703baa92` (the sibling session's status report!) onto the
   branch in the seconds between verdict and push. I saw it in the same command's
   output and folded + re-pushed — but the first push DID go out with an unverified
   daemon commit as head. The ritual's "push IMMEDIATELY" cut both ways here.
4. **First ci-repro failure not diagnosed**: I re-ran instead of capturing the full
   log. The re-run passed, so the failure is now unexplainable except as transient.
   `>/tmp/cirepro-*.log` capture should have been the default from the start.
5. **The unexplained classes.txt +27 delta** sits committed on the CopyButton branch
   (daemon-authored). Not dangerous (inventory is a generated doc surface), but
   shipping an artifact delta nobody can explain is exactly the §csscomment failure
   class. Verify-before-push is mandatory.

## e) WHAT WE SHOULD IMPROVE

1. **The daemon is now actively hostile to this workflow.** This session it: created a
   merge commit with conflict markers; zeroed go.mod pins; **aborted a rebase I was
   mid-conflict-resolution on** (destroying in-progress resolution state); raced junk
   commits into the verdict→push window twice; committed sibling-session artifacts onto
   my branch; split my work into junk-message commits five times. The soft-reset fold
   ritual works but costs real time every hit. A session-scoped pause switch would pay
   for itself in one day (question 1).
2. **Verdict→push window**: the verified-SHA → push gap should be a single chained
   command (verdict tail + `git status -sb` re-read + push in ONE bash call), which I
   adopted late.
3. **Run the ritual inside `nix develop`** — the coverage lane's toolchain warning is a
   standing hazard for gates keyed to the recorded patch level.
4. **`check-tc-sources-sync.sh` could run as a cheap pre-push lane in ci-repro's lint
   step** — it caught nothing early because the mirror check inside the lint lane
   evidently passed while the Go guard failed later; worth investigating why the script
   and the test disagreed (the script found the drift only with `--fix`).
5. **`templ generate` wrapper**: a repo script that does generate + submodule-restore +
   mirror-sync in one call would eliminate my two most repeated mistakes.

## f) NEXT (ordered, ≤50)

1. Verify the `abfc9dc5` classes.txt delta: fresh `gen-class-inventory.sh` + diff; fold correct content.
2. Fold daemon commit into `605b2978`, quick guard re-run, push `feat/copybutton-script-mount`.
3. PR for #63 (`Fixes #63`) → watch CI → merge → verify #63 closes.
4. Confirm master CI + Website green after the CopyButton merge.
5. #53: implement `Stacked bool` opt-in OR ledger as Deferred-v2 (answer to question 3).
6. Root-cause the transient ci-repro FAIL if it recurs; keep full-log capture as default.
7. Re-record coverage baseline on go1.27.0 dev-shell toolchain (`--record`), clear the skew note.
8. Investigate lint-lane vs Go-guard disagreement on `check-tc-sources-sync` (item e.4).
9. Add `scripts/gen-templ.sh` wrapper (generate + submodule restore + mirror sync) — item e.5.
10. Answer v2.0.0 timing (question 2) → if GO: release ritual in dev shell (`release.sh <ver>`),
    tag tree check, `assert_release_tree`, manual push review.
11. TODO 502: BuildFlow gate debt — either fix the gates or document the standing --no-verify policy.
12. Sweep TODO_LIST.md rows made stale by this session's merges (#501 struck already; check 500/385/443/335-adjacent rows).
13. Delete or explicitly keep `revert/templ-pin-1020` (81579dae, unmerged by design — decide its fate).
14. Verify PR #88's closed state carries no dangling branch (it was delete+re-push replaced by #89).
15. Consider a swap-scenario e2e for CopyButton (htmx swap → dead button without shell script → alive with it).
16. Consider Sparkline demo-page variants behind a decision about visual-route golden churn.
17. Post-adoption note on #60/#63 for bank-sync (the consumer evidence behind both issues) — optional courtesy close-loop comment.
18. Audit remaining `docs/agent-context-history.md` anchors referenced by AGENTS.md still resolve.
19. Sweep for other singleton-script emitters with the same fragment-delivery gap (#63 pattern): menu key-nav, image fallback, go-back, kanban listeners, theme toggle — decide per-component whether a page-shell mount export is warranted.
20. Check `wc -l AGENTS.md` ≤ 220 budget (currently fine; re-check after any guard-list edit).
21. Item 19 output → possibly a tracking issue enumerating fragment-safe script delivery per emitter.
22. Verify FEATURES enum tables still exhaustive after this session (no enums added — TestFeaturesEnumTableExhaustive green in ritual).
23. Re-run `nix flake check` (only treefmt currently — sanity after tooling edits).
24. Confirm the dupl baseline's `recordedAt`/pin pair still matches after the hardening (visible in every verdict — yes).
25. Consider recording the daemon's rebase-abort behavior in §duplpin (it happened AFTER the section was written).
26. Sweep `docs/status/` untracked-or-daemon-committed reports for coherence (three landed this session).
27. `nix run .#visual` full pass on master post-merges (CI ran it green per-PR; one consolidated local pass closes the loop).
28. CopyButton: check `docs/recipes/` for a copy-pattern recipe worth adding the fragment contract to.
29. Sparkline: `docs/recipes/bar-chart.md` mentions the height contract; a one-line sparkline cross-reference could live in the same chart recipe family.
30. Check whether `skill/SKILL.md` needs the new exports mentioned (CopyButtonScript is a component-level API).
31. countgen: confirm component count stayed 138 (CopyButtonScript is a Go wrapper, not a templ declaration — drift guard green in ritual, but double-check FEATURES table).
32. Look at `docs/count-definitions.md` subtests rounding band (~1,900 vs actual 1,853 — nearest-hundred).
33. Investigate the +27 classes.txt delta root cause (which file fed new tokens) — item 1's deeper why.
34. Prune the 8 stale `/tmp/tc-demo-route` servers (session hygiene, machine-local).
35. ChastityAPI watcher PID 2598845 — still user-owned, still unasked (prior question, unanswered).
36. Verify `visualtest` compile lane in CI caught the e2e additions (it did — CI green on #90; note for future: compile-check runs without Chromium).
37. Consider `ResolveInclude`-style e2e for Sparkline Filled+Responsive combo at a fluid width (browser test exists for stretch, not for area path under scaling).
38. Check upstream templ for whitespace-leak semantics changes relevant to the `if`-wrapper space leak (filed knowledge only; no action yet).
39. Decide whether the `NoScript` name should be revisited before release (inverted-flag convention vs issue's `SelfScript` proposal — repo convention won; lock it in via CHANGELOG wording).
40. Update `references/component-catalogue.md` if it enumerates per-component script behavior (fan-out source is repo-side).
41. Sweep for lingering `TODO`/`FIXME` in files touched this session.
42. Re-read `docs/plan-authoring-checklist.md` against the #63 delivery (goldens-cover-this ✓, wired⇒e2e waived with reason ✓, count bumps ✓, demo smoke n/a) — formalize the waiver in the PR body at push time.
43. Confirm every `[Unreleased]` entry this session has an issue reference (✓ currently).
44. Housekeeping: `git worktree` list — `/tmp/tc-accent-fix` sibling checkout may still exist; leave alone (sibling session) but note it.
45. Decide the fate of `docs/status/2026-11-07_01-25_…` and this report after merge (daemon auto-commits them; they accumulate by design).
46. After #63 merge: close the loop on issue #22 (referenced in #63 as a different gap) — nothing to do, verify only.
47. After CopyButton push: verify the skill fan-out needs no stamp bump (no release).
48. Consider extracting the "fold daemon commit" ritual into a repo script (`scripts/fold-daemon.sh <base>`) given its frequency.
49. Check `docs/javascript-guide.md` mentions the fragment-script-delivery reality for CopyButton consumers.
50. Re-run the clone gate after all merges on final master (cheap; every verdict already did — do one on the release candidate tip).

## g) THREE QUESTIONS I CANNOT ANSWER MYSELF

1. **The daemon:** it sabotaged six events this session (merge commit with conflict
   markers, zeroed go.mod pins, aborting a rebase mid-resolution, two junk commits in
   the verdict→push window, one sibling-artifact transplant). May I pause/disable the
   auto-commit daemon for the remainder of the agent sweep (and if so, how do you
   prefer it stopped — service unit, env flag, cron remove?), or should I keep paying
   the fold tax?
2. **v2.0.0 + TODO 502 sequencing:** the deferred-v2 pile grows and every commit still
   rides `--no-verify` because of the BuildFlow gate debt (TODO 502). Do you want
   (a) fix the BuildFlow gates first, (b) cut v2.0.0 first (I run `scripts/release.sh`
   in the dev shell; it does NOT push), or (c) keep deferring both?
3. **#53 scope:** defer to the v2 pile per the prior plan, or implement the additive
   non-breaking `Stacked bool` wrapper prop now (opt-in, zero value unchanged, one
   golden + docs)? Both are legitimate; the issue's consumer pain is real but the
   documentation-only option also exists.

---

_Report ends. Awaiting instructions — nothing will be pushed, merged, or killed until
you answer._
