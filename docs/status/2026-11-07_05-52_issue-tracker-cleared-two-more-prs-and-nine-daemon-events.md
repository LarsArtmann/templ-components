# Status Report — 2026-11-07 05:52 CET

**Session scope:** continuation of the standing issue/PR sweep. This session landed the
stalled #63 (PR #91), implemented and landed #53 (PR #92), **cleared the issue tracker to
ZERO open issues**, discovered a repo-wide git-identity contamination, and catalogued four
more daemon race events (7-10) including one where the daemon merged an open PR under the
session.

---

## a) FULLY DONE

| Work                                                       | Evidence                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| ---------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **#63 CopyButton page-shell script mount landed** (PR #91) | branch was found already pushed + folded (the continuation of the prior session completed the fold of daemon commit `abfc9dc5` into `76515b1a`); all 11 CI lanes green (incl. Visual Regression 7m20s); **the daemon ran the rebase-merge itself at 03:38:06Z** (mergeCommit `9b442b9d`) and deleted the branch while the session was inspecting the PR — **#63 auto-closed 03:38:07Z**                                                                                                                                                                                                                                                                   |
| **PR #92 — `Stacked` form fields (#53)**                   | merged `2b7a05f7`; opt-in `Stacked bool` on Input/Select/Textarea/FileInput/Combobox/DatePicker AND on `FormFieldProps`/`FormFieldWrapper` directly; bare `<div>` container (no classes → intra-field spacing identical to today); children pass-through via `fieldContents` sub-template (ADR-0010: two callers, clear name); **all 55 existing forms goldens byte-identical** through the `-update` run; 2 new goldens (`input_stacked_full`, `select_stacked`); `TestFormFieldWrapperStackedPropagatesFromConsumers` proves the prop lands through all six components; ci-repro VERDICT PASS 04:55:48 at the pushed tip; **#53 auto-closed 05:44:00Z** |
| **Issue tracker cleared**                                  | `gh issue list --state open` → **0 open issues** (#48/#59/#35/#60/#63/#53 all closed this cycle; nothing left in the tracker)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| **Master green after both merges**                         | CI `completed success` @ `2b7a05f7`; Master Red Alert + Tidy Probe success; Website workflow correctly absent (path-filtered, forms-only change)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| **Housekeeping**                                           | skill fan-out verified through the symlink path (`skill/SKILL.md` stamp v1.27.0 line 38); count docs same-commit (341 HTML goldens in AGENTS/README/count-definitions/ROADMAP/comparison; countgen rewrote FEATURES; `TestDocsCountDrift` green); mirror synced via `check-tc-sources-sync.sh --fix` (7 files); TODO rows 385/500/502 re-read — all correctly ledgered, owner/demand-gated, no action                                                                                                                                                                                                                                                     |

## b) Git-identity contamination (ACTION NEEDED)

`.git/config` carries `[user] name=test / email=test@example.com` (repo-local; global and
system UNSET). **Every commit made from this worktree since Nov 6 04:13 carries the "test"
identity — including the merged-to-master PR commits (#86/#89/#84/#90/#91/#92 and the
daemon's docs commits).** Master history is merged and therefore untouchable; the config
was fixed repo-local to `Lars Artmann <git@lars.software>` this session (future commits
are correct). Follow-ups for the owner:

1. Check SIBLING repos for the same repo-local config (it predates today; origin unknown —
   a test fixture, a leftover `git config` from some harness, or a hostile-daemon artifact).
2. Decide whether the ~19 "test"-authored commits on master should be noted anywhere
   (rewriting merged master history is off the table by policy).

## c) Daemon events 7-10 (continuing the catalog from the 04:06 report)

| #  | Event                                                                                                                                                                                                                                                                                         | Outcome                                                                                                            |
| -- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| 7  | **Branch switch mid-command:** local HEAD silently moved from `feat/copybutton-script-mount` to `master` between two session tool calls; the session's `git reset --soft f6891f9f` then moved MASTER instead of the feature branch (caught by the chained-verification habit before any push) | recovered: content proven byte-identical to the feature commit, then restored                                      |
| 8  | **Staged-work sweep ×2:** the daemon committed the session's staged #53 work as junk commits `7e0ee7ae` + `d36d0465` on the feature branch (generic messages)                                                                                                                                 | folded into the single proper commit `d8665a97` (33 files, diff-verified identical)                                |
| 9  | **Mid-air PR merge:** the daemon merged OPEN PR #91 (rebase + delete-branch) while the session held it open for inspection; it also pushed the #53 branch (`d8665a97`) between the ci-repro verdict and the session's own push ("Everything up-to-date")                                      | both pushes carried EXACTLY the verified tips — the evidence bar held despite the daemon doing the mechanical push |
| 10 | **Master junk commit mid-merge:** daemon pushed its own docs report (`a36f52a9`) to master while PR #92's rebase-merge was landing; the merge stacked cleanly on top                                                                                                                          | local master re-synced via `pull --rebase`; CI cancelled the superseded run and went green on `2b7a05f7`           |

The chained ritual (verdict → `git status -sb` re-read → push in ONE command) caught
event 7 before damage; events 8-10 were harmless because every daemon push carried
already-verified content.

## d) State

| Ref             | State                                                                  |
| --------------- | ---------------------------------------------------------------------- |
| `origin/master` | `2b7a05f7` (#92), CI success; all workflows green                      |
| Local branches  | `master` (in sync), `revert/templ-pin-1020` (kept, unmerged by design) |
| Open issues     | **0**                                                                  |
| Open PRs        | 0                                                                      |
| `[Unreleased]`  | warm: #63 CopyButton + #53 Stacked entries                             |

## e) Follow-ups (for the next session, owner-gated)

- Identity audit of sibling repos (§b.1).
- Swap-scenario e2e for CopyButton fragment mode + the sweep for other singleton
  emitters sharing the innerHTML delivery gap (waiver noted in PR #91).
- From the 04:06 report's 50-item list, still open: coverage baseline re-record on
  go1.27.0 (item 7), `scripts/gen-templ.sh` wrapper (item 9), transient ci-repro FAIL
  root cause if it recurs (item 6), lint-lane vs Go-guard disagreement on
  `check-tc-sources-sync` (item 8), v2.0.0 ritual (user-gated), TODO 502 pre-commit
  debt (51 findings), TODO 500 EnsureID audit, TODO 385 lint-auto-configure allowlist.
