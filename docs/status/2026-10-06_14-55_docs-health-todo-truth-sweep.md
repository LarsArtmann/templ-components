# Docs-Health Pass — TODO_LIST Brought to True Current State

**Date:** 2026-10-06 ~14:55 CEST
**Branch:** `fix/nonce-omit-empty` (tree: TODO_LIST.md modified; daemon `4dbe67c9` already captured the bulk of the sweep — the last two ID-strike edits are the residual working-tree diff)
**Sources harvested:** `docs/status/2026-10-06_13-12_license-gate-…` §f, `docs/status/2026-10-06_13-13_github-issue-sweep-…` §b/§f, plus the 14:40 continuation report (§b ritual witness, §c open set)

## a) What was done

1. **Read all three 2026-10-06 status reports** (13:12, 13:13, 14:40) and diffed their
   forward-looking items against the existing TODO_LIST rows.
2. **Verified post-report reality against code before harvesting** (grep, not trust):
   - `scripts/session-tripwire.sh` exists → #337 struck correctly.
   - `.buildflow.yml` skips `go-auto-upgrade` → the 13:13 f12 disposition is DONE.
   - `utils/templ_pin_test.go` exists (`TestTemplVersionPin`) → the 13:13 f11 templ-pin
     drift guard is DONE.
   - `AGENTS.md` at 377 lines + `docs/agent-context-history.md` exists → the AGENTS diet
     (13-12 §f1, 13:13 f20) is DONE; the templ-pin sweep narrative lives in the history doc.
   - `cmd/tc/_sources/starter/` tracked set is {app.css, custom.css, templ-components-theme.css}
     → b7a27880 restored the theme file deliberately; `styles.css` + `templ-components-theme.out.css`
     are untracked+ignored daemon litter (tripwire class, guard refuses tracked resurrection only).
   - FEATURES.md has ZERO hits for `ScriptComponent`/`ScriptAttrs`/`HTMLDataAttrs` → the docs-surface gap is real (harvested as #355).
   - CHANGELOG `[Unreleased]` has exactly ONE PolledRegion entry → the 13:13 f19 dedupe task is already clean; not added.
   - `[&>details]` exists only in `website/internal/pages/sales.templ` → #278 settled, struck.
   - Living docs carry no stale `tc new`, no `tc add` coverage claims, CHANGELOG ships the
     22-file rescue → #288's three parts all settled, struck.
   - `docs/version-support.md` "What a floor bump looks like" still lacks the missed steps → #338 stays open.
   - PR #28 is **no longer 10/10**: 12 checks, `Deploy Website + Demo` +
     `Release tag completeness` FAILING (regressed since the 13:13 report) → captured in #353.
   - PR #30 open, base master, mergeable; branch ahead 9.
3. **Harvested 15 new rows (#353–#367)** into a new section, deduped against every existing
   row (the 13-12 table's #322-#352 overlaps and the already-struck/done items excluded).
4. **Corrected stale rows in place:** #270 (release state was v1.19-era; now names v1.21.0 +
   its loaded `[Unreleased]` contents + post-cut checks), #336 (now names PR #30, the
   follow-on commits it carries, the witnessed-green-at-`da8c75a9`-but-re-witness-at-tip state),
   #352 (go-licenses upstream watch), #290 (the struck verdict itself was falsified — see c).
5. **Struck settled rows:** #278, #288. Fixed the header's next-free-ID (said 352; 352 was taken).
6. **Ran the session tripwire** (`scripts/session-tripwire.sh da8c75a9`): only TODO_LIST.md
   dirty — no daemon churn, no resurrected files.

## b) True current state (one paragraph)

Master tip is still the daemon-poisoned red state (`eae4d3d1` templ bump); the complete
CI-unblock sits in PR #28 which has REGRESSED to 10/12 (Deploy Website + Demo and Release
tag completeness failing — diagnose before merge). PRs #30/#31/#32/#33 wait to rebase onto
post-#28 master. The nonce branch carries the whole follow-on wave (guards, diet, tripwire,
tc-init scaffold fix, templ-pin guard) with all three ritual lanes witnessed green at
`da8c75a9` — but 3+ commits landed after the witness, so the ritual must re-run at tip
before push/merge. `[Unreleased]` is heavily loaded; v1.21.0 is the next cut once the chain merges.

## c) The #290 correction (worth its own bullet)

The struck #290 verdict ("the live starter dir is exactly {app.css, custom.css}") was
falsified hours later by `b7a27880`: the deletion misdiagnosed `templ-components-theme.css`
as a zombie when the scaffold's `app.css` imports it — every `tc init` project failed its
first Tailwind compile. The row's verdict is now inline-corrected (annotation placement
rule): tc init scaffolds all three again, `TestStarterCSSMatchesTemplates` byte-parity is
the canon, and the inventory guard refuses only the two truly-dead files. Lesson recorded:
a deletion guarded by a test can still be a wrong deletion — the guard pinned the mistake.

## d) Left open deliberately

- **#340** (standing VERIFY sweep over open rows >7 days): this pass verified the RECENT
  rows plus spot-checks; the full older-row sweep remains its own task.
- Nothing was committed by hand (house rule: no commit unless asked); the daemon has
  already staged the bulk as `4dbe67c9`.
- The two failing PR #28 checks were NOT diagnosed — that is #353's first action.

## e) Next actions (now TODO_LIST #353+)

1. Diagnose PR #28's two failing checks, merge, then rebase/merge #31/#32/#33/#30 (#353).
2. Re-witness ritual at branch tip, push the 9 commits (#336).
3. Post-merge bundle (#354), then the v1.21.0 cut (#270).
