# Status Report — Prerender Drift Flake Fix & Docs-Count Guard Reconciliation

**Date:** 2026-09-23 15:56 CEST
**Scope:** This session only — the failed `test-race` lane (`TestPrerenderMatchesLiveServer` drift on `/`), its root cause, the fix, and what was noticed along the way. No unrelated research was done.
**Tip at writing:** `0927826c` (working tree clean; daemon has committed everything)

---

## Session Summary

The `test-race` lane failed with `prerender drift on /: static snapshot differs from the live page beyond the by-design CSS link`. Root cause: the demo page renders **wall-clock time at second precision** in four places, and the prerender loop writes `index.html` before the live fetch happens — under `-race` on a loaded machine the two renders regularly straddle a second boundary. The normalizer (`normalizePrerender`) covered CSRF tokens and auto-generated IDs but was blind to time. Fixed by extending the normalizer; a pre-existing red `TestDocsCountDrift` (art-dupl 121-vs-122 prose drift) was found and fixed at the same time. Full lane verified green.

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| 1 | **Root-caused the flaky prerender test.** Diagnosed four wall-clock sources on `/`: `Build: {time.Now().UTC()...}` (examples/demo/display_demo.templ:513), `Updated {time.Now().Format("15:04:05")} UTC` (examples/demo/htmx_demo.templ:89), and `display.RelativeTime`'s `datetime` (RFC3339, second precision) + `title` (minute precision) attributes. | Drift reproduced deterministically with a 60-round prerender-vs-live loop; diff showed `datetime="...13:35:42+02:00"` vs `...13:35:43+02:00` — the exact signature. |
| 2 | **Fixed the normalizer.** `normalizePrerender` (examples/demo/prerender_diff_test.go:31) now strips wall-clock-derived values: two `<time>`-element attribute regexes (`datetime`, `title`, anchored so other elements' attributes are untouched) and two demo text-stamp regexes. Design unchanged: structure and all other content still compare byte-for-byte. | Committed (daemon) via `3ae5b6b9`/`cd921d49`; final file verified at HEAD (5 normalizer passes present). |
| 3 | **Fixed pre-existing red `TestDocsCountDrift` at HEAD.** Today's art-dupl second-pass re-recording set the baseline to 121 groups but ADR-0009 + CHANGELOG prose still said 122. Updated both citations (ADR now cites the 2026-09-23 re-recording, nix art-dupl 0.7.0-81ce00b, superseding 122/fork v0.7.0-74-ge7456139). | Guard failed with changes stashed (proving pre-existing), passes now: `go test ./utils -run 'TestVersionMatches\|TestDocsCountDrift'` ok. Committed `0927826c`. |
| 4 | **Verified the exact failed lane green.** `go test -race -shuffle=on -count=1 ./...` — all 10 root-module packages ok, including examples/demo. | Terminal output this session; repeated 3x (10 standalone stress runs of the target test before that, plus the 60-round drift-hunt loop). |
| 5 | **Lint clean.** `golangci-lint run ./examples/demo/...` — 0 issues. | Terminal output. |
| 6 | **CHANGELOG `[Unreleased]` warmed** with the fix entry (repo convention: every fix lands its entry immediately, never at release time). | CHANGELOG.md "Fixed" section; committed `0927826c`. |
| 7 | **Diagnostic scaffolding cleaned up.** Throwaway `zz_prerender_dump_test.go` (used to capture the drift) created, used, deleted. | git status clean at report time. |

## b) PARTIALLY DONE

| Item | What works | What remains | Effort |
|------|-----------|--------------|--------|
| **Flake-class elimination** | The observed failure mode (second-boundary crossings) is fixed and heavily verified. | The class is reduced, not zeroed: (1) RelativeTime **element text** ("12 minutes ago") on the recipes routes flips at minute granularity and is NOT normalized — only its `datetime`/`title` attributes are; (2) calendar Year/Month renders (wire demo, forms section) flip only at month boundaries — accepted risk, unnormalized; (3) the `Build:`/`Updated` text regexes are coupled to demo copy — a copy edit silently un-matches and the flake returns. | S |
| **art-dupl 122→121 reconciliation** | All prose counts now match the committed baseline; guards green. | The **reason** the second-pass recording dropped one group is not documented anywhere I could find — I fixed the numbers, not the story. | M |

## c) NOT STARTED

| Item | Why not started | Priority |
|------|-----------------|----------|
| Unit test for `normalizePrerender` itself (fixture HTML → expected normalization). It is a test-only helper with no direct test; a silent regex rot has no tripwire. | Out of scope while fixing the red lane; small follow-up. | High |
| Full pre-push ritual (`scripts/ci-repro.sh --lint --website`, per-module test loop, visual lane). Nothing was pushed — user did not ask, and the change is test-only + docs. | Waiting for instruction; per M03 ritual it must run at the exact tip before any push. | Critical (before next push) |
| `HARVEST` of this report's section (f) into `TODO_LIST.md`/`ROADMAP.md`. | Report first, wait-for-instructions next (user's explicit protocol). | High |
| gopls `writestring` hints at examples/demo/main.go:53,56,61 — visible in diagnostics all session, unaddressed (lint-clean, advisory only). | Cosmetic; not part of the failure. | Low |
| Design evaluation: pin the demo clock (`RelativeTimeProps.Now` already exists) vs. keep live time + normalize. | Product decision, flagged in question Q1. | Medium |

## d) TOTALLY FUCKED UP

*(radical honesty; items with current state)*

1. **The `test-race` lane was red and nobody knew why it was flaky.** The test's own error message ("re-cut with -prerender") was **misleading** — re-cutting would not have fixed anything, because the snapshot is cut fresh inside the test; the real problem was render-moment wall-clock. Anyone following the message would have wasted time. Severity: blocked the race lane intermittently. Mitigation: fixed this session.
2. **`TestDocsCountDrift` was red at HEAD for an unknown duration.** The daemon commits without running tests (documented in AGENTS.md, still unfixed), so a docs-count drift sat at the tip until this session's guard run tripped over it. Severity: every local root-module `go test ./...` failed; CI would have caught it, but local "why is my tree red" confusion is the actual cost. Mitigation: fixed; the daemon hole remains open.
3. **The daemon committed a throwaway diagnostic file into history** (`zz_prerender_dump_test.go` landed in one auto-commit, its deletion in another). Harmless history noise, but it proves the daemon snapshots the tree mid-session without judgment. Severity: cosmetic/history pollution. Workaround for next time: write diagnostic scratch outside the watched worktree (`/tmp` + a targeted test invocation).
4. **Daemon commit messages remain hallucinated/generic** ("chore: auto-commit N changed file(s) (heuristic)") — `git log --grep` for this fix finds nothing; only this status report and the CHANGELOG point at the commits. Known AGENTS.md issue, upstream (`larsartmann/buildflow`), still unfixed.

## e) WHAT WE SHOULD IMPROVE

1. **Normalize wall-clock by default in byte-compare tests of rendered pages.** This fix was per-incident; the third nondeterminism source was found only because I looped-to-failure. A shared normalizer contract (or a "no bare `time.Now` in demo render paths without a `// prerender-safe` marker" convention + scan) would stop the whack-a-mole.
2. **Loop-to-failure is the only honest flake diagnosis.** My first single-shot dump PASSED and nearly led to "not reproducible". The 60-round loop found the drift in seconds. Make the loop the standard first move for flaky render comparisons (candidate for a small helper or a note in `docs/testing/`).
3. **Diagnostic scratch must live outside the auto-committed worktree.** /tmp script + targeting existing test internals beats a committed `zz_*` file that the daemon snapshots.
4. **A test-only fix still deserves the ritual.** I verified the failed lane + lint but did not run `ci-repro.sh --lint --website`. Before the next push, the ritual must run at the exact tip (per M03) — queued in section (f).
5. **Misleading failure messages cost more than the bug.** "re-cut with -prerender" sent the reader in the wrong direction. Error messages in guards should name the actual invariant and its likely causes; worth an audit of similar "re-cut" hints in other guard tests.
6. **The docs-count guard caught the drift only because a human ran tests.** Guard coverage without a pre-commit/daemon execution point is theater for daemon-only commits — see (d)2.

## f) Top Things To Get Done Next (ranked; brainstorm for HARVEST — most items beyond the top few are ROADMAP fuel)

| # | Task | Impact | Effort | Category |
|---|------|--------|--------|----------|
| 1 | Run `scripts/ci-repro.sh --lint --website` at tip `0927826c` before any push (M03 ritual; includes visual + website lanes not run this session) | Critical | L | Quality |
| 2 | Normalize or pin the RelativeTime **text** on recipes routes ("12 minutes ago" minute-boundary flake remains) | Critical | S | Bug |
| 3 | Add unit test for `normalizePrerender` with fixture HTML (tripwire against silent regex rot) | High | S | Quality |
| 4 | Run the per-module test loop (`utils icons errorpage charts/echarts datastar htmx` + visualtest compile) — not run this session | High | M | Quality |
| 5 | Document WHY the art-dupl baseline went 122→121 (which group vanished; re-scan with pinned binary and triage the diff) | High | M | Documentation |
| 6 | HARVEST this report's section (f) into `TODO_LIST.md`/`ROADMAP.md` | High | S | Documentation |
| 7 | Single-source the prerender test's routes table from `prerender.go`'s page list (a new prerendered page silently escapes comparison today) | High | M | Quality |
| 8 | Pin or normalize calendar Year/Month in demo (month-boundary flake, wire + forms demos) | Medium | S | Bug |
| 9 | Add a flake-stress mode to `TestPrerenderMatchesLiveServer` (loop N rounds / fake clock) so drift-hunting doesn't need a throwaway file | Medium | M | Quality |
| 10 | Evaluate `RelativeTimeProps.Now` pinning for demo pages (decision pending Q1) | Medium | S | Feature |
| 11 | Audit all guard tests for misleading remediation hints ("re-cut with -prerender" class) | Medium | M | Quality |
| 12 | Daemon: cheap smoke (`go build ./...` or the drift guards) before auto-commit, so red-at-HEAD tips become impossible | High | L | Quality |
| 13 | Daemon: generate commit messages from `git diff --stat` (upstream buildflow fix; makes `git log --grep` useful) | Medium | L | Cleanup |
| 14 | AGENTS.md note: prerender normalizer now covers wall-clock; new demo time renders must be added to it (prevents regression-by-addition) | Medium | S | Documentation |
| 15 | Testing docs: write the "wall-clock in rendered HTML" rule into `docs/testing/` (generalize the lesson) | Low | S | Documentation |
| 16 | Fix gopls `writestring` hints at examples/demo/main.go:53,56,61 | Low | S | Cleanup |
| 17 | Verify `[Unreleased]` CHANGELOG coherence after daemon edits (entries accumulate under heuristic commits) | Low | S | Documentation |
| 18 | Decide patch-release now vs. ride next release for this fix (test-only; likely ride) | Medium | S | Decision |
| 19 | `prerenderCSRFToken` regex assumes `name` before `value` attribute order — harden or document | Low | S | Quality |
| 20 | Add backoff (or jitter) to the prerender test's 100ms fixed retry sleep | Low | S | Quality |
| 21 | Sweep for OTHER byte-compare tests (goldens, website build checks) that render bare `time.Now` content | Medium | M | Quality |
| 22 | Remove the bun shim at `~/.local/bin/node` (documented local pnpm blocker, still present) | Low | S | Cleanup |
| 23 | Upstream buildflow: templ-generate step re-appends `*_templ.go` to `.gitignore` each run (documented gotcha) | Low | L | Cleanup |
| 24 | Plan the deliberate lockstep Go toolchain bump (re-enables go-structure-linter + go-version-auto-configure; TODO #231 class) | Medium | L | Decision |
| 25 | Migrate remaining ~48 raw `chromedp.Poll` call sites (TODO_LIST #240, pre-existing) | Medium | L | Quality |
| 26 | Track art-dupl upstream truncation bug (TODO #292, pre-existing) | Low | M | Cleanup |
| 27 | Consider a root meta-guard that fails when `go list ./...` silently excludes workspace modules (the AGENTS-documented `go test ./...` blind spot) | Medium | S | Quality |
| 28 | Clean stale `/tmp/prerender-diff` artifacts from this session's diagnosis | Low | S | Cleanup |
| 29 | Add a marker convention (`// prerender-safe` or similar) + scan so new `time.Now` demo renders can't silently re-open the flake | Medium | S | Quality |
| 30 | Record the art-dupl binary version beside the baseline count wherever it is cited (ADR partially does; CHANGELOG entry now does) | Low | S | Documentation |

*(30 items — honest stop; the remaining "50" would be padding. Items 12–13 and 23–26 are pre-existing backlog observations from AGENTS.md, restated here because they intersected this session.)*

## g) Questions I Cannot Figure Out Myself

1. **Should the demo keep rendering live wall-clock time at all?** `RelativeTimeProps.Now` makes pinning the clock cheap, which would zero out this entire flake class — but the live "Updated HH:MM:SS" / "2 hours ago" stamps are demo *showmanship* (selling server-rendered freshness). Keep live + normalize, or pin for determinism? This is a product call about what the demo is for; I can argue both sides.
2. **Where did the `test-race` failure you pasted actually run — local buildflow or CI?** If local, nothing needs pushing and the fix just needs the pre-push ritual whenever you next push. If CI, master's tip was red for consumers of the race lane and I should run the full `ci-repro` + push ritual now (with your go-ahead). I cannot tell from the paste which environment produced it.
3. **Was today's art-dupl second-pass re-recording (122→121 groups) a deliberate re-triage with a rationale recorded somewhere I didn't look, or an accidental overwrite?** If deliberate, item 5 above is just "link the rationale". If accidental, the baseline may have silently dropped a group that should be accepted-clone-documented — a possible duplication-gate regression I cannot distinguish from a legitimate re-triage without redoing the binary-pinned scan.

---

*Report written per the status-report skill; format override honored (skill's canonical output is a styled HTML dashboard — user explicitly requested `.md`). Section (f) is the primary input for a future `docs-health` HARVEST run. **Waiting for instructions.***
