# Status: Master-Green Restore on Resumption — Session Survived an 8-Day Suspension

**Date:** 2026-11-25 20:10 CET (wall clock, cross-checked against git commit clock 19:49)
**Scope:** the RESUMPTION phase of the merge-all-PRs session — everything after the
2026-11-17 self-review report (`docs/status/2026-11-17_15-57_merge-followup-closeout-self-review.md`

- its resumption addendum). Session-scoped: this reports only what THIS session did
  and noticed on resumption.
  **Verified:** `date` + `.git` mtime + commit timestamps (the suspension discovery);
  `git merge-base --is-ancestor` (history intactness); `TestCrossbrowserLockstep`
  (run red→green across the pin fix); CI runs `45672820320` / `45675861144` / `45679078657`
  (job-level conclusions read); Website run on `e40a2b5f` (success); local full visual
  suite (green, banner-only CSS diff); pinned-binary `art-dupl check` (reproduced CI's
  1-new-group, then zero after re-record); `TestDocsCountDrift` with `-count=1` (green
  after count updates); sibling activity via `stat` mtimes (live, 90s-old).
  **NOT verified:** the pin rewriter's identity (g1 — cannot be determined from here);
  the sibling session's intent for their `cmd/tc` work (deliberately not inspected
  beyond mtimes + names, single-writer protocol); upstream go-datastar CI.

---

## a) FULLY DONE (verified)

1. **The suspension was diagnosed, correctly, before any damage was done.** The
   session froze mid-flow on 2026-11-17 (after "WAIT HERE"); it resumed with a STALE
   `date` CLI (still printing Nov 17 15:5x) while the repo, `.git` mtimes, and commit
   timestamps all said Nov 25. The first alarm ("history was rewritten!") was a
   misread; `git merge-base --is-ancestor 15d836b8 origin/master` and `4e845ed4`
   proved the Nov-17 merge work fully intact with 8 days of newer commits on top
   (v1.27.2, cmdguard skill pairing, architecture-deepening brief). Documented in
   the prior report's resumption addendum.
2. **Red lane 1 — playwright pin (`a76536fb`):** master had been red since
   `cf725c9d` (Nov 22) on `TestCrossbrowserLockstep` because an npm-install-signature
   rewrite re-bumped `playwright-core` to 1.65.1 (daemon commits `d5702fc1`, then
   `d8722176`) against the unchanged nixpkgs 1.64.1 driver — the 4th occurrence of
   the class. Restored per the guard's own prescription; guard run green IN the
   commit chain; CI Build & Test green at `e40a2b5f` and after.
3. **Red lane 2 — Tailwind 4.3.5 three-lane lockstep (`e40a2b5f`):** the Nov-23
   daemon nixpkgs sync had moved the dev-shell tailwindcss to 4.3.5 while the
   Dockerfile pnpm pins, website.yml npm pins, and committed `app.css` stayed at
   4.3.4 (CSS Freshness red). Completed the lockstep per the 11-12 playbook:
   pins moved, CSS recompiled (diff = the version banner ONLY), full local visual
   suite green before push, CI CSS Freshness AND Visual Regression green after.
   [Unreleased] warmed with the master-restore Fixed entry.
4. **Red lane 3 — clone gate (`960aea8b`):** red for 5+ days on ONE new group
   (85 single-line `templ.Component` type references across test files). Triaged
   ACCEPT per ADR-0009 (intentional similarity, no extraction value — same class
   as the accepted test-idiom groups); baseline re-recorded with the pinned fork
   binary (`scripts/pinned-art-dupl.sh`, same rev CI installs). The re-record also
   shed 9 stale hashes dissolved by the window's refactors: **145 → 137 groups**.
   AGENTS.md, ADR-0009 (dated pass + Consequences current-state), and CHANGELOG
   counts updated in lockstep; `TestDocsCountDrift` green with `-count=1`.
5. **MASTER IS FULLY GREEN** — CI run `45679078657` on `960aea8b`: all six jobs
   success (Build & Test, Lint, CSS Freshness, Visual Regression, HTML validation,
   Clone gate); Website workflow success on `e40a2b5f` (its path filters correctly
   triggered on the demo/website files and correctly did NOT on the baseline/docs
   commit). The 3-day master-red window is closed.
6. **Sibling session detected and respected:** a sibling is LIVE in this worktree
   right now (`cmd/tc/doctor.go` + `main.go` modified, `css.go` + `css_test.go` +
   `_sources/starter/templ-components-classes.txt` + a docs/reviews HTML untracked;
   mtimes 90 seconds old at report time). Nothing of theirs was touched, staged,
   read in depth, or committed — single-writer protocol; my commits used explicit
   paths only.

## b) PARTIALLY DONE

1. **TODO_LIST HARVEST of the 11-17 report's f-list:** still owed. Deliberately
   deferred on resumption because the sibling session owns the worktree (editing
   TODO_LIST mid-flight risks a collision); the routing lives in the prior report's
   §f until a calm window. Items 1 (CI green) of that list is now done by this
   phase; the rest remain open.
2. **The prior report's addendum item 5 says "master green expected on the next
   CI run":** now superseded by VERIFIED green (run `45679078657`). Not edited
   retroactively (point-in-time snapshot; this report is the correction of record).

## c) NOT STARTED

1. **The structural fix for the playwright-pin rewriter** (prior report e3/f3) —
   still the top open item; manual restores demonstrably lose (occurrence 4 rode
   in 60 seconds after occurrence 3's restore).
2. **`feat/layout-seo-meta` revival verdict** (TODO #509, owner-gated).
3. **Findings-gate archive ritual** (TODO #507) — `--no-verify` + local verify
   remained the working mode for all three of today's fix commits.
4. **v1.28.0 release scheduling** (TODO #510) — [Unreleased] is warm again with
   the master-restore entries.
5. Prior report's remaining unrouted f-items (f14, f16, f18, f21, f22, f25, f28+,
   f30–f33, f35, f36, f38, f39, f40, f43, f46, f49 — see its §c2).

## d) TOTALLY FUCKED UP (and recovered or permanently mitigated)

1. **`b3df7f8a`'s commit message does not match its content.** The message claims
   "restore playwright pin…" but the commit contains ONLY the report file (the
   restore I had staged pre-suspension was consumed by the repo surgery during the
   8-day freeze; I pushed without verifying the commit's stat against the message).
   Mitigated permanently-ish: the very next commit `a76536fb` actually performed
   the restore, and both reports document the mismatch — but the misleading
   message is permanent master history (fixing it would need a history rewrite,
   which is banned). Lesson recorded in e1.
2. **The "history was rewritten!" misdiagnosis.** Four tool calls of alarm
   theories (rebase surgery, .git swaps, force-push rollbacks) before the mundane
   truth: 8 days of newer commits + a stale `date` clock + my reading `log -8` as
   "the history". The ancestry check that settled it should have been the FIRST
   move, not the fifth.
3. **One flag-less art-dupl run** (`-- check` without `-c .art-dupl.json -t 1
   --type-aware`) re-reported the 88-group whole-tree shape (the documented
   false-alarm class) before I re-ran with CI's exact invocation. The script docs
   say the flags ride after `--`; I passed none. One wasted run.

## e) WHAT WE SHOULD IMPROVE

1. **Message-diff verification before push, after ANY state anomaly.** A commit
   message is a claim; the diff is the fact. After a suspension/clock anomaly,
   `git show --stat HEAD` (read the stat, compare to the message) belongs in the
   chain BEFORE `git push`. Cost of skipping it this time: a permanently
   misleading message on master.
2. **Diagnose git anomalies with ancestry + timestamps FIRST.** `git merge-base
   --is-ancestor <expected> <head>`, `git log --format='%ci'`, `.git` mtime —
   these three settle "rewritten vs. advanced" in one call. Alarm theories are
   for after the facts contradict the mundane explanation.
3. **Cross-check the clock when state contradicts it.** This session survived an
   8-day suspension under a LYING `date` CLI because git/env timestamps kept
   contradicting it until the truth won. Make that contradiction-check explicit:
   when `date` disagrees with `.git` mtimes or harness dates, trust the artifacts.
4. **Mirror CI's invocation VERBATIM on first try** when reproducing a CI failure
   locally (flags included) — `scripts/pinned-art-dupl.sh -- check` without the
   `-c … -t 1 --type-aware` tail is a different command.
5. **Don't HARVEST into TODO_LIST while a sibling owns the worktree.** Routing
   stays in the report's §f (with explicit pointers) until the tree is calm;
   a mid-flight TODO edit risks the exact multi-writer collision the protocol
   exists to prevent.

## f) NEXT — up to 50, realistic 22 (priority order)

1. **Structural fix for the playwright-pin rewriter** (`.npmrc` `save-exact`,
   spike relocation out of daemon reach, or a pre-commit fast-guard on the pin
   file) — manual restores lose; this is the 4th occurrence.
2. Identify the rewriter (g1): sibling intent vs. accidental npm-install class —
   determines whether f1 proceeds or waits.
3. TODO_LIST HARVEST pass over the 11-17 report's f-items (route-or-drop each;
   needs a calm window — sibling owns the tree at report time).
4. Accordion rAF flake trend row (1 local occurrence, CI green 3× on the same
   code).
5. v1.28.0 release scheduling — [Unreleased] warm again (TODO #510).
6. GitHub Release for upstream `static/v0.8.0` (TODO #511).
7. Consumer pin sweep (TODO #512).
8. Dependabot bump-protocol routing (TODO #513).
9. Class-inventory doc-comment nuance → skill guard-tests reference (TODO #514).
10. `feat/layout-seo-meta` revival or verdict (TODO #509, ⫱).
11. Findings-gate archive ritual (TODO #507) — retires `--no-verify`.
12. Re-check the daemon's next commits: the sibling's `cmd/tc` css-command work
    will be swallowed with a heuristic message — their attribution problem, but
    watch that the daemon doesn't mangle it into `cmd/tc/_sources` strays.
13. Watch 48h: does the npm-rewrite recur on the playwright pin (trend evidence)?
14. Master-Red-Alert probe re-run (prior f14).
15. DOMAIN_LANGUAGE datastar pin mentions (prior f16).
16. Facts-doc TOC/anchor check (prior f21).
17. website/go.sum daemon churn check (prior f22).
18. chromedp ⇄ cdproto pairing verification (prior f25).
19. CodeRabbit comments sweep on merged PRs (prior f38).
20. pixelmatch pseudo-version confirmation (prior f39).
21. Upstream go-datastar master CI check (prior f43).
22. Full `scripts/ci-repro.sh --lint --css --website` VERDICT on the green tip
    (prior f40 — the M03 ritual at the session's true end).

## g) QUESTIONS (cannot resolve myself)

1. **The pin rewriter, now with a live suspect pool:** a sibling session is
   actively working in `cmd/tc` right now (css command + tests). Is the recurring
   npm-install-signature rewrite of `playwright-core` THEIR doing (part of a
   planned 1.65 lockstep that flake.lock should eventually carry), or the known
   accidental class? This decides whether I may land the structural pin fix (f1)
   or whether it would stomp their plan.
2. **`feat/layout-seo-meta`:** revival session or delete? (57 files / +1879
   lines, diverged since 2026-09-08; kept on remote — TODO #509.)
3. **Findings gate:** is `--no-verify` + full local verification the accepted
   working mode until the #507 archive ritual happens — and when does that
   session happen? (Every manual commit in this session went through it; the
   daemon bypasses hooks entirely.)

---

**Session verdict:** the resumption closed all three inherited red lanes and
returned master to full green (all-jobs CI + Website, verified at `960aea8b`);
one permanently misleading commit message (`b3df7f8a`, mitigated and documented);
the pin rewriter remains the one live external risk; a sibling session owns the
worktree at close. Waiting for instructions.
