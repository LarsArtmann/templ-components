# Session Report — 2026-11-07 01:25 CET — Five PRs Merged, BarChart Verified, Two Branches In Flight

Continuation of the standing issue/PR sweep. This session closed two more
issues, merged four PRs, and has one fully verified branch awaiting push plus
one conflicted branch awaiting rebase.

---

## a) FULLY DONE

| Work                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | Evidence                                                           |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| Master CI/Website verified green on `6082430e`, then every push gated by `ci-repro.sh --lint --css --website` VERDICT PASS                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | runs on `6082430e`, `a86cef37`, `c68c3424`, `6ee5db4a`, `59c6934b` |
| Housekeeping batch: 5 backup branches deleted after tree-diff/cherry content verification (`backup/*` ×4, `fix/ci-go-127`); `revert/templ-pin-1020` kept (unmerged by design); TODO **#998** re-verified green (already resolved in PR #78); TODO row **502** written (BuildFlow gate debt); merge-session lessons appended to history (§merge77)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | `a86cef37` on master                                               |
| **Issue #66 closed** — CI golangci-lint pin ALREADY at v2.15.0 (`482e3a93`, ancestor of master), dev shell matches, per-module lint green in CI + local ci-repro; `check-lint-config.sh` asserts the lockstep                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     | evidence comment on #66                                            |
| **Issue #50 fixed + auto-closed (PR #81, `c68c3424`)** — `skill/SKILL.md` now stamps `Released API: v1.27.0` + package-qualification note (`FilterBar` is `forms.FilterBar`); `utils.TestVersionMatchesSkillStamp` joins the version-guard family; `release.sh` gained step 6d (stamp bump + rollback list + targeted-guard regex). Root finding: the issue's symptom table probed `display/` for components living in `forms/` — all four "missing" components exist in v1.27.0                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | PR #81 merged, #50 closed                                          |
| **TODO #501 resolved (PR #82, `6ee5db4a`)** — `docs/recipes/accessible-clickable-rows.md` (contract + worked `TableCell.Content` example + stretched-link variation with tradeoffs + anti-patterns incl. the `role="link"` that v1.27.0 still renders — the removal is UNRELEASED, verified by tag probe); registered in recipe index; TODO row struck                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            | PR #82 merged                                                      |
| **Color-token CSS guard (PR #83, `59c6934b`)** — `utils.TestSemanticColorTokensCompiled`: extracts literal hue-shade tokens from every `.templ` class attribute (canonical tone hues + accent), asserts each resolves in the compiled demo CSS; family-anchored regex kills the z-50/duration-700 false-positive class; **mutation-tested** (planted `bg-emerald-950` failed with file attribution); `walkTemplFiles` lost its always-`".."` parameter (unparam)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  | PR #83 merged                                                      |
| **Issue #85-class root cause (PR #85, `73e3f438`)** — the "unreproducible" app.css drift SOLVED: `demo.css` scans `**/*.go` too, so my guard test's doc comment (`bg-accent-600/40`) was itself a Tailwind candidate; set-diffing the two blobs showed exactly ONE inserted rule. Comment reworded, AGENTS.md comment-hygiene rule extended to Go files, history doc §csscomment records the mechanism                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            | PR #85 merged                                                      |
| **Issue #35 implemented + locally verified (branch `fix/tc-doctor-false-fails`, `140e235f`, PR #84 OPEN)** — templ pin check is now a hand-parser (block + single-line require, `// indirect`, trailing comments, closed/reopened blocks, replace/exclude noise ignored; no `golang.org/x/mod` — dependency budget respected); Tailwind entry check discovers by content (any `*.css` carrying `@source` within 5 directory levels; fast-path candidates kept; `tc doctor --css PATH` override; per-file hint preserved when `**/*.templ` missing); `usage()` updated; 9 table-driven tests. Verified against the issue's repro layout (assets/css/input.css + block-form go.mod → one true pass, one true positive) and the repo itself (all checks pass where both used to false-FAIL)                                                                                                                                                                                          | local ci-repro VERDICT PASS at `140e235f`                          |
| **Issues #48+#59 implemented + browser-verified (branch `fix/bar-chart-vertical`, 3 commits, ahead of master)** — vertical BarChart rewritten: container stretches (dropped `items-end`), columns `h-full`, omitted `Height` defaults to `12rem` (was silent 0px), bars absolutely positioned; **signed rendering**: negatives hang below an auto-drawn zero baseline in the danger family (`NegativeBarColor`, default rose), `ShowBaseline` opts all-positive charts in, `MaxBarWidth` caps columns + centers the track; `barChartBounds`/`barChartSignedPlacement` helpers; 3 new goldens + updated vertical golden; 5 new unit tests incl. a placement table; **2 chromedp e2e tests pass in a real browser** (height contract: bars > 0px, ordered, 6x ratio; signed layout: positive bar < zero line < hanging negative by element centers); recipe sections; CHANGELOG Added+Fixed; mirror sync; golden-count 333→336 across 7 doc sites (countgen + hand-maintained rows) | final ci-repro **VERDICT: PASS (exit 0)** 01:19:46 CET             |

## b) PARTIALLY DONE

1. **`fix/bar-chart-vertical` is verified but NOT pushed.** VERDICT PASS at
   01:19:46, then two follow-up commits (mirror sync `96a5f076`, count bump
   `c983b68b`) landed after it. Per the M03 ritual the branch needs one more
   full ci-repro at the new tip (c983b68b) → push → PR (Fixes #48 #59) → CI →
   merge. Note: master moved to `73e3f438` (#85) after the branch was cut —
   rebase first; `examples/demo/static/app.css` WILL conflict (both sides
   recompiled it); resolve by regenerating at the tip, never by hand-picking.
2. **PR #84 (doctor) is OPEN and CONFLICTING** with master (CHANGELOG
   overlap with #81/#83 entries). GitHub never started its Actions lanes —
   `mergeable=CONFLICTING` suppresses them; I discovered this only after a
   close/reopen attempt fired nothing. Needs: rebase onto current master,
   resolve CHANGELOG, push with `--force-with-lease` (own feature branch),
   CI then runs, merge, #35 auto-closes.
3. **The clone gate in the final VERDICT line read `[dupl: art-dupl version dev]`**
   — the PINNED binary (ADR-0009 binary-pinning rule, `v0.9.1-…-93432fd`) was
   missing again (/tmp volatility) and ci-repro silently fell back to a dev
   build. The gate still reported zero new clones, but per the ADR a
   cross-binary compare can false-alarm or false-clear. Before pushing the
   BarChart branch: rebuild the pinned binary
   (`scripts/pinned-art-dupl.sh`) and re-run the clone gate with it; then
   make ci-repro fail loudly (or warn) when the pinned binary is absent.

## c) NOT STARTED

1. **#60 Sparkline** — responsive mode + zero-aware/signed bounds (+ optional
   zero line). The BarChart signed geometry (`barChartBounds` /
   `barChartSignedPlacement`) is a reuse candidate — compose, don't
   reimplement (also an art-dupl watch item).
2. **#63 CopyButton** — page-level `CopyButtonScript` mount + opt-out on the
   per-render singleton + fragment-contract docs (+ e2e for the swapped-in
   case).
3. **#53 forms.Input/Select layout wrapper** — read-fully-and-decide only;
   changing the DOM structure is breaking → likely a Deferred-v2 row with a
   design note, not an unilateral implementation.
4. **v2.0.0 release** — user-gated (the accent-emission BREAKING entry
   forbids any v1.x cut; ADR-0039 gating). `release.sh` ritual + lockstep
   tags + website pairing notes; ALSO dry-run the new 6d stamp-bump step
   (never executed in a real release yet).
5. TODO **500** (EnsureID render-random id audit sweep) — large.
6. TODO **502** (BuildFlow gate policy: 47 file-size + 4 docs-gate findings;
   fix sweep vs advisory) — decision pending, row written.
7. TODO **385** (BuildFlow allowlist for accepted-disabled linters).
8. Upstream watch (golangci-lint v2.16.0 will auto-file like #66 did).

## d) TOTALLY FUCKED UP (own failures, ranked)

1. **I caused the CSS "nondeterminism" I then spent ~30 minutes diagnosing.**
   My guard test's doc comment contained a complete class token; the daemon
   committed the correct recompile mid-flow and pushed it to the PR branch,
   so it landed on master inside the #83 merge. The fix (#85) is clean, but
   the whole detour — two failed ci-repro runs, a stat-cache goose chase, a
   phantom watcher hunt — was self-inflicted. The rule now exists (AGENTS.md,
   §csscomment); I should have thought of scanner candidates BEFORE writing
   the example.
2. **Missed the same-commit count bump for the 3 new goldens** — the rule is
   verbatim in AGENTS.md ("update the doc numbers in the SAME commit or CI
   fails"); I added the goldens, forgot countgen, and burned a full ci-repro
   cycle on TestDocsCountDrift.
3. **Missed the `_sources` mirror sync** — bar_chart.templ is mirrored;
   `--no-verify` commits bypass the self-staging guard, so the mirror drifted
   and burned another ci-repro cycle. Both misses are the same shape:
   known same-commit obligations, skipped in the moment.
4. **PR #84's dead CI**: I closed/reopened the PR hoping to re-trigger before
   checking `mergeable` — GitHub suppresses workflows for conflicting PRs.
   Checking mergeable FIRST (one API call) would have saved the detour.
5. **FileName-form flip-flop**: generated `bar_chart_templ.go` from inside
   display/ (module-relative), then from root (repo-root-relative), wasting
   two whole-tree regen cycles (44- and 72-file diffs) before settling on
   root-form to match origin/master. The committed canonical form per module
   should have been checked first.
6. **Daemon-interleave surprise ×7**: seven `chore: auto-commit` commits
   landed inside my work across five branches; every one needed a
   soft-reset fold. Documented since #77, still.costing attention each time.
7. **Hallucinated `utils.ParsePercent`** in the templ rewrite — caught
   immediately by the build, but it should never have been typed.
8. **Two test-expectation stumbles** in TestBarChartSignedPlacement (wrong
   height arithmetic) and the style-assertion trailing-semicolon — quick
   fixes, but the second one I "fixed" once via edit-tool before the file
   had settled, triggering a rewrite via heredoc (which itself hit an
   auto-formatter race). Net: three tool cycles for two one-line fixes.

## e) WHAT WE SHOULD IMPROVE

1. **Check `mergeable` immediately after `gh pr create`** — conflicting PRs
   silently get no CI; the check is one call and would have saved the #84
   detour.
2. **Same-commit obligations checklist as a pre-commit mental gate**: new
   goldens ⇒ `countgen -w` + TestDocsCountDrift; mirrored component touched ⇒
   `check-tc-sources-sync.sh --fix`; new prose in scanned trees ⇒ no complete
   class tokens. All three bit me tonight; all three are greppable.
3. **Pin the art-dupl binary non-volatilely** (nix derivation or
   repo-local cache path) and make ci-repro REFUSE the dev-build fallback —
   tonight's VERDICT line silently degraded gate quality (see b.3).
4. **Fold-the-daemon helper**: `scripts/fold-daemon.sh` doing
   fetch + reset-soft-to-base + recommit would turn a repeated 4-command
   dance into one deterministic step.
5. **CHANGELOG conflict pressure**: four PRs tonight each edited
   `[Unreleased] ### Added` head → every rebase conflicts. An append-at-top
   or per-PR anchor convention would remove the friction.
6. **`CHROMEDP_CHROME_PATH` in plain `nix develop`**: e2e tests SKIP without
   it (only `nix run .#visual` sets it). Exporting the pinned chromium in the
   devShell would let `go test ./...` in visualtest actually run e2e locally.
7. **Doctor's `templGeneratedWith` is a duplicated string** (`v0.3.1100` in
   cmd/tc/doctor.go vs utils' pin): the next templ bump can forget it. A
   cross-module sync guard (or reading the constant) would close it — small,
   worth doing with the next templ event.
8. **Kill the foreign watcher**: `tailwindcss -i assets/css/input.css` under
   `/home/lars/projects/ChastityAPI` (PID 2598845) has been running since
   Nov 6 at 48% CPU / 12 GB RAM in D-state. Another project's process — not
   mine to kill, but it will eat the machine until someone does.
9. **Stale demo servers**: eight `/tmp/tc-route-golden-*/tc-demo-route`
   processes from earlier visual runs are still alive — visualtest may be
   leaking servers on some path.
10. **ci-repro's line 15-17 note** ("never root-caused") should be updated to
    point at §csscomment — the mystery now HAS a root cause pattern.

## f) NEXT ITEMS (up to 50, ordered by leverage)

1. Rebuild pinned art-dupl, re-run clone gate with it (b.3) — BEFORE any push.
2. ci-repro: fail loudly when the pinned art-dupl binary is missing.
3. Rebase `fix/bar-chart-vertical` onto `73e3f438`, regenerate app.css at the
   tip, full ci-repro, PUSH, PR (Fixes #48 #59), CI green, merge.
4. Rebase `fix/tc-doctor-false-fails` onto current master (CHANGELOG),
   `--force-with-lease` push, verify CI fires (mergeable clean), merge → #35
   closes.
5. Confirm master CI/Website green after the BarChart + doctor merges.
6. Verify the fan-out picked up the SKILL.md stamp (check-skill-fanout).
7. Delete local branches after merges (feat/semantic-color-guard,
   fix/scanner-candidate-comment, fix/tc-doctor-false-fails,
   fix/bar-chart-vertical once merged).
8. `revert/templ-pin-1020` deletion (local + remote) — owner sign-off
   (unmerged by design, superseded by the v0.3.1100 migration).
9. #60 Sparkline (reuse the new signed geometry; goldens + e2e).
10. #63 CopyButtonScript page mount (+ fragment docs + e2e).
11. #53 read + decision (implement vs Deferred-v2 row).
12. Sparkline/BarChart geometry dedupe check (art-dupl watch).
13. BarChart: decide whether the negative bound also honors an override
    (currently max-side only).
14. BarChart horizontal signed rendering (currently clamps to 0-width; the
    NegativeBarColor doc says so) — v2 candidate note.
15. e2e for MaxBarWidth/ShowBaseline geometry (only signed + height have e2e).
16. Demo: add a vertical BarChart card (currently only horizontal; would sell
    the signed mode too — pixel re-baseline required, so batch with other
    visual changes).
17. Doctor constant ↔ utils pin sync guard ( Templ-bump drift risk).
18. Update ci-repro's stale "never root-caused" note → point to §csscomment.
19. Consider a repo lint for complete class tokens inside comments of scanned
    trees (cheap regex; would have caught the whole detour).
20. Kill the ChastityAPI tailwind watcher (user call, it's another project).
21. Investigate/clean the tc-demo-route server leak (8 stale processes).
22. CHANGELOG conflict-pressure convention (append-anchor or per-PR slots).
23. scripts/fold-daemon.sh helper.
24. TODO 502 decision (fix sweep vs advisory gates) — needs owner input.
25. TODO 500 EnsureID sweep (Modal/Drawer/Dropdown/… fleet).
26. TODO 385 BuildFlow allowlist ask.
27. v2.0.0 release prep checklist (after the remaining v2-shape work or now —
    owner decision).
28. release.sh dry-run exercising the new 6d stamp bump.
29. Website pairing notes for v2 (templ ≥ v0.3.1100 ↔ v2) on the install page.
30. Archive/annotate this + prior status reports (docs-health pass).
31. TestClassInventoryFreshness negative test? (family has none; consistent —
    decide once).
32. BarChart aria enhancement: announce the value range in the default label.
33. Confirm LoadingProgress's non-component (positional-args) catalogue
    semantics are intended.
34. docs/plan-authoring-checklist: add the two new same-commit gates (counts,
    scanner-token prose).
35. Upstream golangci-lint v2.16 watch (auto-issue will file).
36. Doctor: strict vs lenient unknown-arg handling (currently lenient) — note
    as intentional in --help or tighten.
37. Re-check TODO row 501's struck table renders post-dprint (pipes intact).
38. Verify TestVersionMatchesSkillStamp behaves when skill/SKILL.md is absent
    in odd checkouts (Skipf — matches family precedent ✓, just document).
39. Consider Make-less devshell export for CHROMEDP_CHROME_PATH (e.6).
40. After v2.0.0 cut: website install-page pairing + README badge + SKILL
    stamp all bump together (release.sh 6d/6c/6d verify).
41. Sweep docs/status for daemon-eaten sections (the §merge77 lesson applied
    retroactively).
42. Check whether any other _test.go/.go comments contain complete class
    tokens (grep-able hygiene sweep, related to 19).
43. Sparkline demo entry (if #60 adds responsive mode, demo it).
44. BarChart: MinBarWidth + MaxBarWidth interaction on tiny N — test the
    min>max edge or guard it.
45. Appraise #63's "opt-out" API shape against #53's outcome (both touch
    wrapper/emission contracts — sequence them).
46. Website build lints recipes? (bar-chart.md got a new section — website
    lane in ci-repro covers build only; check the site renders recipes or
    only curated pages — the sidebar does NOT list recipes, so no site work
    needed; confirm and note).
47. Re-run `nix run .#visual` locally once on master post-merges (pre-empt the
    CI Visual lane rather than discovering drift there).
48. Ask owner about the `art-dupl version dev` fallback hardening (e.3) if
    not done in this session.
49. Keep an eye on Actions minutes: tonight burned ~10 CI runs; concurrency
    group already queues; fine.
50. Sleep the daemon question: raise the "quiet flag" idea to the owner —
    auto-commit-on-worktrees is the recurring tax of every session.

## g) QUESTIONS (decisions I cannot make alone)

1. **v2.0.0 timing** (carried from the 22:45 report, still unanswered): the
   accent-emission BREAKING entry forbids any v1.x cut from this tree, and
   #48/#59 (just built) + #501/#50/#35/#66/#83/#85 are now on master. Cut
   v2.0.0 once the BarChart + doctor branches land, or wait for #60/#63
   (and the #53 decision)?
2. **BuildFlow gate policy** (TODO 502): 47 file-size + 4 docs-gate findings
   force `--no-verify` on every commit (tonight it directly caused the
   mirror drift). Dedicated fix sweep, or flip those two gates to advisory?
3. **ChastityAPI watcher**: it is eating 48% CPU / 12 GB since yesterday —
   kill it, or is it holding state you need?

---

## Evidence index

| Ref                         | State                                                                                                                                        |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| `origin/master`             | `73e3f438` (#85), CI+Website green; #81 `c68c3424`, #82 `6ee5db4a`, #83 `95aea41e` (incl. daemon app.css), #85 `73e3f438` all merged         |
| `fix/bar-chart-vertical`    | 3 commits (`273021df`, `96a5f076`, `c983b68b`), ci-repro VERDICT PASS 01:19:46 (dev-dupl caveat b.3), needs rebase onto `73e3f438` then push |
| PR #84 (doctor, `140e235f`) | OPEN, CONFLICTING with master, CI suppressed; needs rebase                                                                                   |
| Local ci-repro runs         | PASS at `a86cef37` 23:06, `e78687f4` 23:16, `7bb2ad7b` 23:28, `b1d62d3c` 00:19, doctor `140e235f` 00:32, `a30866f7` 00:35, BarChart 01:19    |
| Open issues remaining       | #63, #60, #53 (+#48/#59 close on PR)                                                                                                         |
| TODO rows touched           | 502 written, 501 struck, 998 confirmed resolved                                                                                              |
| Backup branches             | deleted (verified); `revert/templ-pin-1020` kept                                                                                             |
