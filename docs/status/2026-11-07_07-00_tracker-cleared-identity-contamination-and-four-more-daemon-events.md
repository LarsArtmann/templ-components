# Status Report — 2026-11-07 07:00 CET

**Session scope:** continuation window 04:06 → 07:00 CET of the standing issue/PR sweep.
Two PRs landed (#91 #63, #92 #53), the issue tracker reached ZERO open issues, a
repo-wide git-identity contamination was discovered and locally fixed, and four more
daemon race events (7-10) were catalogued.

**Verified here** (fresh commands, this window): `gh pr view 91` MERGED `9b442b9d`
03:38:06Z + `gh issue view 63` CLOSED 03:38:07Z; `gh pr checks 91` all 11 lanes pass
(read back 04:37); `gh pr checks 92 --watch` all lanes pass (Build & Test 4m55s, Visual
Regression 7m9s); `gh pr view 92` MERGED `2b7a05f7` + `gh issue view 53` CLOSED
05:44:00Z; `gh issue list --state open` → **0**; `gh run list --branch master` → CI
success @ `2b7a05f7`, Master Red Alert + Tidy Probe success; local
`scripts/ci-repro.sh --lint --css --website` → `VERDICT: PASS (exit 0)` 04:55:48 at
`d8665a97`; `go test ./forms/ ./integration/ ./utils/ -count=1` green;
`golangci-lint run ./forms/...` → 0 issues; `-update` golden run left all 55 existing
forms goldens unchanged (`git status` showed only 2 new files); `find . -name '*.golden'`
→ 341; `TestDocsCountDrift` green after the 5-file count bump; fold correctness proven
by `git diff 9b442b9d d36d0465` = exactly the 33-file change set; author of `d8665a97`
= `Lars Artmann <git@lars.software>` (config fixed first); master author distribution
read via `git log --format='%an' -40 origin/master`.

**NOT verified here:** the lost continuation's claimed ci-repro PASS at `79330718` (no
local log artifact exists; GitHub CI is the only confirmation); WHO authored the
04:12-04:17 continuation commits (git forensics says identity `test` — same as every
commit from this worktree since Nov 6, so it proves nothing); the mechanism of the
mid-session branch switch (inferred from state deltas, never observed); Website
workflow non-trigger for #92 (inferred from path filters, config not re-read); sibling
repos' git config (out of scope); the daemon's `a36f52a9` status-report content (stat
read only, text not proofread).

---

## a) FULLY DONE

| Work                                                        | Evidence                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| ----------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **#63 CopyButton page-shell script mount LANDED** (PR #91)  | MERGED `9b442b9d` (rebase-merge), merge at 03:38:06Z performed BY the daemon while the session held the PR open for inspection; all 11 CI lanes green incl. Visual Regression 7m20s; **#63 auto-closed** 03:38:07Z. Branch state at resume (folded commit `76515b1a`, pushed, PR open) was the work of a lost continuation 04:12-04:17 — reconstructed and verified via git forensics, not re-done                                                                                                               |
| **#53 `Stacked` form fields IMPLEMENTED + LANDED** (PR #92) | designed and built this window: `Stacked bool` on Input/Select/Textarea/FileInput/Combobox/DatePicker + `FormFieldProps`/`FormFieldWrapper`; bare `<div>` container via `fieldContents` pass-through sub-template (ADR-0010 qualified); **all 55 existing forms goldens byte-identical**; 2 new goldens; propagation test across all six consumers; 33-file single commit `d8665a97` (correct author); ci-repro VERDICT PASS at pushed tip; all PR lanes green; MERGED `2b7a05f7`; **#53 auto-closed** 05:44:00Z |
| **Git identity contamination DISCOVERED + locally FIXED**   | `.git/config` carried `user=test <test@example.com>` (repo-local; global/system UNSET); every worktree commit since Nov 6 04:13 carries it incl. merged master commits; config set to `Lars Artmann <git@lars.software>` BEFORE the #92 commit (verified: `d8665a97` author correct)                                                                                                                                                                                                                             |
| **Issue tracker CLEARED**                                   | `gh issue list --state open` → 0; no open PRs; `[Unreleased]` warm with 6 entries                                                                                                                                                                                                                                                                                                                                                                                                                                |
| **Master green after both merges**                          | CI success @ `2b7a05f7`; Master Red Alert + Tidy Probe success; local master in sync (`git status -sb` clean at 07:00)                                                                                                                                                                                                                                                                                                                                                                                           |
| **Daemon junk commits folded, content proven**              | `7e0ee7ae`+`d36d0465` (the daemon's sweep of my staged #53 work) folded into `d8665a97` ONLY AFTER `git diff 9b442b9d d36d0465` proved the combined content = the intended 33-file change set; earlier `bc38dad6` (daemon sweep of staged #63 work on master) proven redundant vs `7fccc348` before discard                                                                                                                                                                                                      |
| **Count docs + mirror hygiene, same commit**                | countgen `-w` rewrote FEATURES; 339→341 hand-fixed in AGENTS/README/count-definitions/ROADMAP/comparison; `TestDocsCountDrift` green; `check-tc-sources-sync.sh --fix` synced 7 mirror files; classes.txt inventory refreshed by its own guard and committed                                                                                                                                                                                                                                                     |
| **Housekeeping**                                            | skill fan-out verified through the symlink chain (`skill/SKILL.md`, stamp v1.27.0 line 38); TODO rows 385/500/502 re-read (all correctly ledgered, no action); wrap-up report written 05:52                                                                                                                                                                                                                                                                                                                      |

## b) PARTIALLY DONE

1. **Identity contamination — local fix done, audit NOT done.** Repo-local config fixed
   (evidence above). Remaining: audit SIBLING repos for the same repo-local config, and
   decide whether the ~19 `test`-authored master commits need any annotation. Blocker:
   none — effort S; simply not started (owner may want to lead the sibling-repo audit).
2. **AGENTS.md memory NOT updated.** The identity gotcha and the daemon-can-merge-PRs
   escalation are recorded in status reports but NOT in the project AGENTS.md — a fresh
   session reads AGENTS.md, not timestamped reports. Not done because the report-then-wait
   instruction landed first; effort S. (Same for a possible ADR on the page-shell script
   mount pattern from #63.)
3. **#53 demo discoverability.** The six components gained `Stacked`, but the demo forms
   page gained no example (opt-in prop, existing pages unaffected — CI green without it).
   Remaining: one demo section + smoke; effort S; deprioritized deliberately.
4. **Post-merge HARVEST not run.** The status-report skill says §f feeds docs-health
   HARVEST into TODO_LIST/ROADMAP; the session continued past the 05:52 report without
   running it. Remaining: one HARVEST pass; effort S.

## c) NOT STARTED

- **v2.0.0 release ritual** — user-gated (standing decision; `release.sh` does not push).
  `[Unreleased]` now holds SIX features; sequencing question in §g.
- **v1.28.0 cut** — the alternative to v2.0.0 batching; not started, owner decision.
- **Sibling-repo git-config audit** (§b.1) — waiting on owner input re: blast radius.
- **CopyButton swap-scenario e2e + singleton-emitter delivery-gap sweep** — waived in
  PR #91 with rationale (innerHTML non-execution is spec behavior; default path
  byte-identical); still the right follow-up work, demand-gated.
- **Coverage baseline re-record on go1.27.0** (prior report item 7) — statement counts
  are toolchain-specific; the committed baseline predates the skew warning.
- **`scripts/gen-templ.sh` wrapper** (item 9) — submodule-safe generate; still hand-run.
- **Transient ci-repro FAIL root cause** (item 6) — did not recur this window.
- **lint-lane vs Go-guard disagreement on `check-tc-sources-sync`** (item 8) — unchanged.
- **TODO 502** (51 pre-commit findings → `--no-verify` debt), **TODO 500** (EnsureID
  audit), **TODO 385** (lint-auto-configure allowlist) — all correctly ledgered, none
  started; 502 is the single highest-friction debt (bit THIS session: see d.1).

## d) TOTALLY FUCKED UP

1. **~19 commits on master carry author `test <test@example.com>` — permanently.**
   Severity: medium (cosmetic/attribution, not functional). Root cause: repo-local
   `.git/config [user] test` set sometime before Nov 6 04:13; EVERY session (mine, the
   daemon's sweeps, the continuation) committed under it; nobody — including me across
   two prior merges (#86/#89/#84/#90) — checked authorship before push/merge. Merged
   history is untouchable by policy. Mitigation: config fixed; but the residual risk is
   OTHER REPOS carrying the same config, unaudited.
2. **I violated the session's own `--no-verify` convention on one commit** (04:39, the
   #53 fold attempt). The BuildFlow hook ran, failed on the 54 documented pre-existing
   findings (TODO 502), and — because the hook held the index while the daemon was
   sweeping — my staged work landed inside daemon junk commits. Severity: low (recovered
   byte-identical), but it created the exact race window the ritual exists to avoid.
   Root cause: convention lives in the summary/conversation, not in AGENTS.md or a hook.
3. **A state-mutating git command ran without a same-command branch assertion** (04:39,
   `git reset --soft f6891f9f` while HEAD had been silently switched to master). My
   status check was 2 minutes stale; the daemon's branch switch landed inside the gap.
   Result: local master moved, feature changes staged on master, daemon swept them into
   `bc38dad6`. Recovered fully (content proven redundant before discard), but the ritual
   as practiced — "re-read `git status -sb` before pushing" — does NOT cover RESET-type
   commands, only pushes. Root cause: no enforced branch guard on mutating commands.
4. **The identity-fix window for PR #91 was missed because nobody was watching.** The
   continuation pushed the PR at ~04:17 under the `test` identity and went idle; I
   resumed at 04:37; the daemon merged at 04:38:06. The 21-minute unwatched gap is where
   clean history was lost. Root cause: resume checklist lacks an authorship check, and
   an open PR with a known defect sat unattended.
5. **Dead docs landed on master:** `docs/drafts/2026-11-07_issue63-copybutton-page-shell-script-pr.md`
   (`9b442b9d`) is a verbatim duplicate of the PR #91 body, and the #63 status report
   (`5fded784`) briefly existed as a duplicate branch commit before the merge collapsed
   them. My mid-flight cleanup (unstage + trash) was overwritten by the daemon's
   parallel push/merge. Severity: low (clutter), but it demonstrates the session and
   the daemon editing the same history without coordination.
6. **Session-resume context was stale and I initially trusted it.** The summary said
   "#63 NOT pushed"; the branch was in fact pushed with a PR open. Cost: ~10 minutes of
   forensic reconstruction (merge-base checks, ancestry probes) before acting. Correct
   instinct — re-derive from git/gh — but it should be the DEFAULT first step, not a
   reaction to surprise.

## e) WHAT WE SHOULD IMPROVE

1. **Same-command branch guard for EVERY mutating git command**, not just pushes:
   prefix `test "$(git branch --show-current)" = "<expected>" && ...`. Impact: would
   have prevented d.3 entirely. Effort: S (a shell habit + a line in AGENTS.md).
2. **`--no-verify` made structural** until TODO 502 lands: export `GIT_COMMIT_HOOK=skip`
   or wrap `git commit` in a repo alias, or fix the 51 findings so the hook stops
   blocking (the real fix). Impact: removes the recurring d.2 mistake class. Effort: S
   for the guard, M-L for 502.
3. **Authorship guard in the pre-push ritual:** `git log --format='%an %ae'
   origin/master..HEAD | grep -v 'Lars Artmann'` before any push; and an authorship
   check as the FIRST step of any session resume. Impact: prevents another d.1/d.4.
   Effort: S.
4. **Daemon authority decision (owner-gated):** it now merges PRs, pushes branches, and
   pushes master (events 9-10). Either constrain it (commits only, no push/merge) or
   formally accept its merge authority and drop the duplicated session-side merge
   effort. Impact: removes the entire race class or makes it harmless by contract.
   Effort: S (a config change on your side).
5. **Session-resume protocol: re-derive state BEFORE reading the summary's claims**
   (`git status -sb; git log --oneline -5; gh pr list; gh issue list` first). Impact:
   kills the d.6 confusion class. Effort: S.
6. **AGENTS.md as the durable memory:** identity gotcha, daemon-merges-PRs escalation,
   and the branch-guard rule belong there (size-budgeted, 2-3 lines), not only in
   timestamped reports. Impact: a fresh session stops re-deriving this. Effort: S.
7. **countgen scope:** the 341-golden bump needed hand-edits in 5 files;
   countgen could own AGENTS/README/count-definitions/ROADMAP/comparison the way it
   owns FEATURES. Impact: removes a recurring same-commit chore + drift risk. Effort: M.
8. **Drafts should not ride feature branches:** the PR-draft commit habit (continuation)
   creates daemon-swept clutter on master (d.5). Write PR bodies to /tmp or compose
   directly in `gh pr create`. Effort: S.

## f) Next tasks (up to 50; ranked, impact-first — HARVEST input for TODO_LIST/ROADMAP)

| #  | Task                                                                                                                                                                                                                                       | Impact | Effort | Category      |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------ | ------ | ------------- |
| 1  | Add identity guard (author check) to pre-push ritual + session-resume checklist                                                                                                                                                            | High   | S      | Quality       |
| 2  | Add same-command branch guard to all mutating git commands (AGENTS.md rule)                                                                                                                                                                | High   | S      | Quality       |
| 3  | Decide daemon authority: commits-only vs accepted-merger (owner)                                                                                                                                                                           | High   | S      | Decision      |
| 4  | Cut v1.28.0 OR batch into v2.0.0 (six warm `[Unreleased]` entries; owner)                                                                                                                                                                  | High   | S      | Decision      |
| 5  | TODO 502: fix 47 file-size + 4 docs-gates findings so hooks stop needing `--no-verify`                                                                                                                                                     | High   | L      | Quality       |
| 6  | Audit sibling repos for repo-local `user=test` git config                                                                                                                                                                                  | High   | S      | Cleanup       |
| 7  | Update AGENTS.md: identity gotcha, daemon-merge escalation, branch-guard rule (≤220-line budget)                                                                                                                                           | High   | S      | Documentation |
| 8  | Sweep singleton script emitters for the innerHTML delivery gap (#63 class): Alert/Toast dismiss, LiveRegion, Kanban, MenuKeyNav, MobileMenu, GoBack, table-row-href, Combobox, ImageFallback, Tooltip aria, ThemeToggle, CopyButton (done) | High   | M      | Bug           |
| 9  | Consider a `PageShellScripts` bundle component collecting all page-shell singleton scripts in one call                                                                                                                                     | Medium | M      | Feature       |
| 10 | CopyButton swap-scenario chromedp e2e (PR #91 waiver follow-through)                                                                                                                                                                       | Medium | S      | Quality       |
| 11 | ADR for the page-shell script-mount pattern if the sweep (item 8) confirms it generalizes                                                                                                                                                  | Medium | S      | Documentation |
| 12 | Delete dead `docs/drafts/2026-11-07_issue63-*.md` from master (verbatim PR-body duplicate)                                                                                                                                                 | Low    | S      | Cleanup       |
| 13 | Run docs-health HARVEST from this report's §f into TODO_LIST/ROADMAP                                                                                                                                                                       | Medium | S      | Documentation |
| 14 | Demo forms page: add a `Stacked` example + smoke (discoverability)                                                                                                                                                                         | Medium | S      | Feature       |
| 15 | Audit NON-wrapper labeled fields for the same grid-scatter problem (#53 scope was the 6 wrapper consumers; FilterInput/FilterDropdown/TagsInput/Slider/Checkbox/Radio/Toggle have own label handling)                                      | Medium | M      | Bug           |
| 16 | Decide whether the `Stacked` container should carry a canonical class (e.g. `tc-field`) for CSS hooking vs bare `<div>` (owner design call)                                                                                                | Low    | S      | Decision      |
| 17 | Coverage baseline re-record on go1.27.0 (toolchain-specific counts; prior item 7)                                                                                                                                                          | Medium | M      | Quality       |
| 18 | `scripts/gen-templ.sh` submodule-safe wrapper (prior item 9)                                                                                                                                                                               | Medium | S      | Quality       |
| 19 | Root-cause the transient ci-repro FAIL if it recurs (prior item 6)                                                                                                                                                                         | Low    | M      | Bug           |
| 20 | Resolve lint-lane vs Go-guard disagreement on `check-tc-sources-sync` (prior item 8)                                                                                                                                                       | Low    | M      | Bug           |
| 21 | TODO 500: EnsureID render-random-id audit                                                                                                                                                                                                  | Low    | M      | Quality       |
| 22 | TODO 385: golangci-lint-auto-configure accepted-disabled allowlist                                                                                                                                                                         | Low    | S      | Quality       |
| 23 | countgen: own the golden-count lines in AGENTS/README/count-definitions/ROADMAP/comparison (kills the 5-file hand chore)                                                                                                                   | Medium | M      | Quality       |
| 24 | Investigate the ORIGIN of the `test` git identity (shell history, harness configs, `.envrc`) — prevent recurrence repo-wide                                                                                                                | Medium | S      | Bug           |
| 25 | `git config --local include.path` / template check: ensure future clones of this repo can't inherit a bad identity (setup-hooks.sh assertion)                                                                                              | Medium | S      | Quality       |
| 26 | Stacked-prop parity check: struct-literal zero-value documentation claim vs actual default (verified this session, but add it to the convention checklist for future bool props)                                                           | Low    | S      | Documentation |
| 27 | Check `origin/feat/layout-seo-meta` staleness (sibling session branch; PR #87 merged long ago) — prune if dead                                                                                                                             | Low    | S      | Cleanup       |
| 28 | Post-release task: after next cut, verify `check-release-tags.sh` set includes all published sub-modules (standing ritual, next release)                                                                                                   | Medium | S      | Quality       |
| 29 | Consider `PageShellScripts` nonce handling doc in `docs/javascript-guide.md` once item 8/9 lands                                                                                                                                           | Low    | S      | Documentation |
| 30 | Re-run `scripts/coverage-library.sh --gate 94` after the two landed PRs (ratchet sanity; last recorded pre-#92)                                                                                                                            | Medium | S      | Quality       |
| 31 | Add a ci-repro lane printing `git branch --show-current` at start AND end (race detection for event-7-class switches)                                                                                                                      | Medium | S      | Quality       |
| 32 | File the `test`-identity master-history fact into `docs/agent-context-history.md` (append-only) so future attribution questions have a canonical record                                                                                    | Low    | S      | Documentation |

## g) Questions I cannot answer myself

1. **Daemon authority:** the auto-commit daemon has escalated again — it merged open PR
   #91 by itself (03:38Z) and pushed both branches and master (events 7-10). Do you want
   it (a) constrained to commits-only, (b) paused during agent sessions, or (c) formally
   accepted as merger-of-record so sessions stop racing it? I tried to find its config
   (the daemon rides your home-manager/crush-config, not this repo) and could not
   determine the right lever.
2. **Identity origin:** do you know when/why `git config user.name test` was set in this
   repo — a test harness, a sibling session, a tool install? I could not find the
   setting event from git alone (config files have no history), and the answer decides
   whether item 6 (sibling-repo audit) is a precaution or a fire drill.
3. **Release sequencing:** `[Unreleased]` carries six features (CopyButton fragments,
   Sparkline responsive/zero, BarChart signed series, semantic color guard, Stacked
   fields, tc doctor). Cut v1.28.0 now, or hold for the v2.0.0 window you flagged
   earlier? The v2.0.0 decision gates TODO 491-style breaking props too, so this
   sequencing is yours to call.
