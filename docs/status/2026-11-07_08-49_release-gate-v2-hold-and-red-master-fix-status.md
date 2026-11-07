# Status Report — Release Gate Investigation, v2 Hold, Red-Master Fix

**Point-in-time:** 2026-11-07 08:49 CET (push-state updated 08:53) · **Branch:** `master` @ `c10ffbe2` — the daemon PUSHED all 3 local commits to origin at ~07:46Z unasked; CI run 42352650488 on the fix commit was in progress at write time · **Session scope claimed:** release assessment + master CI un-red. No other scope touched.

## Scope line — what this report VERIFIED vs NOT

**Verified this session (command/artifact):**

- Tag/commit state: `git tag --sort=-v:refname`, `git log v1.27.0..HEAD --oneline` → 32 commits since v1.27.0 (8 daemon heuristic commits).
- `[Unreleased]` warmth + content: read the full CHANGELOG section, including both `BREAKING` entries.
- Emission state in tree: `grep -rn "accent-[0-9]" / "blue-[0-9]"` per package (counts below).
- ADR gating: read ADR-0039 (full), ADR-0049 (§Status/§Decision 1–2), ADR-0044 (referenced), release-checklist (full), master plan S04–S34 rows, TODO_LIST v2 rows (#397/#399/#425/#430/#431/#432/#484/#491).
- CI red root cause: `gh run view 42346078958 --log-failed` → `TestCrossbrowserLockstep`; `nix eval --raw nixpkgs#playwright-driver.version` → `1.64.1`; fix verified by `go test ./utils/ -run TestCrossbrowserLockstep -count=1` → ok.
- CI/Website status on tip `67a9537b` via background poll: CI=failure, Website=success.

**NOT verified / NOT attempted (and why):**

- `scripts/ci-repro.sh --lint --css --website` at `c10ffbe2` — NOT run (the M03 pre-push ritual; push is yours, and you chose hold before I reached the ritual step).
- CI on `c10ffbe2` — the daemon pushed at ~07:46Z unasked; run 42352650488 was IN PROGRESS (4m42s) at write time. Result not yet observed. Note: the intermediate daemon-only tip (`3befa4c3`) went through CI as run 42350615544 = failure — it still carried the broken pin before my fix landed.
- The 6 sub-module test lanes (icons/errorpage/charts/echarts/datastar/htmx) inside CI run 42346078958 — MASKED: the per-module loop runs under `bash -e` and aborted at `utils` FAIL; those lanes never executed in that run. They will run on the next push.
- Full audit of the daemon's 22-file commit `67a9537b` — I inspected the two npm files in detail and the file list only; I did NOT diff-audit the FEATURES.md / website version-support / recipe changes it carried.
- Crossbrowser spike runtime (`nix run .#crossbrowser`) — not run; only the lockstep pin guard.
- Sibling session's datastar-client research scope — untouched by protocol (single-writer).

---

## Session narrative (what actually happened)

1. **Release assessment.** Latest tag v1.27.0 (2026-11-04); 32 commits since; `[Unreleased]` warm with ~27 entries. MINOR bump (v1.28.0) looked mechanical — until the changelog audit surfaced **two `BREAKING` entries**, one carrying the explicit sentence: _"No v1.x release may be cut from a tree carrying this swap (ADR-0039 v2 gating)."_
2. **Tree verification vs the BREAKING claims.** Re-verified against source, not the changelog: the accent swap IS live (`feedback/*.templ` emits `accent-*`); `layout.ContainerProps.Pad` is gone (`NoPad` inverted flag landed, old spelling = compile error). Per-package emission counts: `feedback` 13 accent / 0 blue (fully swapped) · `display` 31/11 (partial) · `forms` 0/29 · `navigation` 0/10 · `layout` 0/2 · `errorpage` 0/16 · `recipes` 0/4 (not started).
3. **Gating verdict.** ADR-0049 (Accepted): the emission "ships only in the v2.0.0 window"; master deliberately carries the breaking state between the first sweep and the v2 cut. ADR-0039 (Proposed, owner decision **F058 pending**): the first breaking release requires the `/v2` module-path migration across all 7 modules. Master plan: v2 cut (S29) blocked on S04–S13 + S15–S18/S24/S27 + S33 + S34 ("v2 ships complete or not at all"). Remaining v2-scope checks: ADR-0022 alias removal already done (grep: no `AlertType`/`ToastType` left), TODO #425 done, #431 planned-not-executed, **#484 open with an owner call** ("v2 material alongside #425/#431").
4. **Key correction of an inherited claim.** The sibling status doc (2026-11-06 05:23) asserted "master itself is emission-free, so v1.28.0 is mechanically possible". That was true at write time and FALSE after `6082430e` (design token axes) landed on master later that day. The v1.28.0 escape hatch is closed as of now.
5. **Red master found and fixed.** CI on tip `67a9537b`: Build & Test FAILED, Website success. Root cause: the daemon swept an npm-install rewrite of `tools/crossbrowser-spike` playwright-core 1.64.1 → 1.65.0 into its 22-file commit; `TestCrossbrowserLockstep` (pins ⇄ nixpkgs `playwright-driver`) failed exactly as designed. Verified the lockstep direction (nixpkgs = 1.64.1, test constant = 1.64.1), restored both files to the pre-daemon content, guard green. Commit attempt blocked by the BuildFlow pre-commit findings gate (**54 pre-existing findings**: file-size-check 49 + docs-gates 5 — known debt, unrelated to the change). The daemon then swept my staged files into its own heuristic commit `841f657a` (3 files: my 2 + the sibling's untracked plan doc); I amended the message only → `c10ffbe2`, content untouched.
6. **Owner decision surfaced and taken:** hold — no release; `[Unreleased]` stays as the v2.0.0 staging area.

---

## a) FULLY DONE

| Item                                                                                 | Evidence                                                                                |
| ------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------- |
| Release-readiness verdict with source-level proof (v1.x forbidden, v2 gated)         | ADR-0049 §Status + grep audit + master plan S29/S34 blockers                            |
| Red master CI diagnosed to root cause (daemon-swept npm pin rewrite)                 | `gh run view --log-failed`, guard test output                                           |
| Lockstep direction verified before restore (nixpkgs 1.64.1 = test constant)          | `nix eval nixpkgs#playwright-driver.version`, test constant read                        |
| Playwright pin restored, guard test green                                            | `git restore --source=67a9537b^`, `go test -run TestCrossbrowserLockstep -count=1` → ok |
| Fix committed with a truthful message after daemon hijack (message-only amend)       | `c10ffbe2`, `git show --stat` = exactly 3 expected files                                |
| v2 deferred-list audit (aliases gone, #425 done, #431 planned, #484 open owner call) | TODO_LIST rows, feedback grep                                                           |
| Owner decision (hold) surfaced with complete tradeoff picture                        | question tool exchange                                                                  |
| Release-checklist + ADR-0039/0044/0049 + master plan read BEFORE any cut action      | session log order                                                                       |

## b) PARTIALLY DONE

| Item                             | Gap                                                                                                                                                                                                                                                                            |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| "Master is green" checklist gate | Was green-expected only up to `fdbcc9da` (06:00Z run). The daemon then pushed the broken-pin tip (CI failure, run 42350615544) AND my fix (`c10ffbe2`, run 42352650488 in flight at write time). Master red is fixed locally and CI verification is pending the in-flight run. |
| M03 pre-push ritual              | MOOT for c10ffbe2: the daemon pushed it without the ritual. Open question whether to run `scripts/ci-repro.sh --lint --css --website` retroactively at the pushed tip or accept CI-on-push (see §g2).                                                                          |
| Daemon 22-file commit audit      | File list reviewed, only the 2 npm files diff-audited; FEATURES.md/website/recipe diffs in that commit un-audited (CI's drift guards partially cover, but the Build&Test lane masked 6 module lanes).                                                                          |
| v2 emission sweep state mapping  | Counts per package captured, but the remaining `display` 11 blue sites were not classified (deliberate chart-hue keepers per plan vs leftover S04/S05 sites).                                                                                                                  |
| Crossbrowser lane health         | Pin guard green; runtime (`nix run .#crossbrowser`) unexercised — ADR-0048 marks that lane findings-only, ⫱G1 pending, so this matches plan state, not a new gap.                                                                                                              |

## c) NOT STARTED (within this session's natural reach)

- v2 campaign steps S04–S13 (display pt2, forms, navigation+layout, errorpage+recipes, tone-map collapse S09, guard upgrades S10, surface/ink S11–S13, theme css S13).
- S33 `/v2` migration + `docs/migration/v1-to-v2-module-path.md`; F058 decision record.
- S34 RC checklist; S29 release cut (blocked anyway).
- CHANGELOG structural fix: split v2-gated entries under a "Ships in v2.0.0" sub-heading (sibling Q2/item 39) so an emergency v1.28.x cut becomes mechanically safe.
- TODO_LIST harvest of section (f) below.
- `ci-repro` verification of `c10ffbe2`; the push itself.

## d) TOTALLY FUCKED UP (as found, with blame where due)

1. **The daemon committed onto master a guard-breaking pin rewrite** (npm 1.64.1→1.65.0) inside a 22-file docs sweep — master red ~08:25 window per CI (06:14Z run), undetected for ~2h because nothing pushes or gates daemon commits pre-push. The lockstep guard caught it in CI, but only AFTER landing. (Daemon behavior, known adversarial-by-accident class.)
2. **The pre-commit findings gate is red for every human/agent commit** (54 pre-existing findings) while the daemon bypasses hooks entirely — the hook currently gates only the actors least able to break things and no actor fixes findings under a 60s budget. Documented interim (--no-verify) is doing permanent duty.
3. **The daemon hijacked the fix commit** — my staged 2-file fix got committed under a hallucinated heuristic message together with an unrelated sibling plan doc, producing a mixed-purpose commit on master history. Amended message-only; content is correct; attribution and mixed scope are permanent blemishes in history.
4. **One CI lane failure masked six module lanes** (`bash -e` loop aborts at first failing module): the "Build & Test" red means only root + utils were actually proven failing that run; icons/errorpage/charts/datastar/htmx statuses are UNKNOWN in that run, not passing.
5. **Point-in-time status docs go stale within hours and lie convincingly**: the sibling's "master is emission-free → v1.28.0 mechanically possible" premise was invalidated the same day by the emission merge. I caught it by re-verifying against source; a session that trusted the doc would have cut a forbidden release.
6. **`[Unreleased]` mixes v1-safe and v2-gated BREAKING entries** — if anything ever forces an emergency v1.28.x, the notes mechanism as-is ships v2 BREAKING copy in a v1 tag body.
7. **The daemon pushed master unasked at 07:46Z**, including my fix commit, bypassing the M03 ritual the human-bound workflow requires (see §d6).

## e) WHAT WE SHOULD IMPROVE

1. **Gate the daemon, not just the humans.** The highest-leverage fix available: a post-daemon-commit hook/CI check that runs the cheap guard set (lockstep, replace-directives, version-sync) before the daemon's push step. Today every guard runs after the damage is on master.
2. **Split `[Unreleased]` into "Ships in v1.x" vs "Ships in v2.0.0" sub-headings NOW** (while no cut is imminent), so release mechanics stay safe under time pressure. (Sibling item 39 — endorse.)
3. **State single-writer scope at session start, always.** I stated it mid-report but not at session start per protocol; the datastar sibling appeared mid-session.
4. **Run the M03 ritual myself when I am the one producing the commit others will push** — I handed off the command instead of running ci-repro at `c10ffbe2`.
5. **Anticipate known-red gates**: read the newest status doc's BuildFlow debt row BEFORE the first commit attempt; one hook run (~60s+) was spent discovering a documented state.
6. **Never read "Build & Test: failure" as one fact** — enumerate the masked lanes (`bash -e` loop) before concluding blast radius.
7. **Status docs need a kill switch**: a header field like "superseded-by" or a freshness assertion when their core premise depends on a volatile state (emission-free, pin-X). The 05-23 doc actively misled for hours.
8. **Triage the 49 file-size + 5 docs-gates BuildFlow findings** into a dedicated task (fix or budget) so the pre-commit gate returns to actually gating. Until then, --no-verify remains the documented interim.
9. **Question-tool discipline**: two malformed calls before one landed — format checks cost nothing; round trips with the owner do.
10. **Keep re-verifying inherited claims against source** (the one thing this session did right — the emission-free claim re-check changed the entire verdict). Institutionalize: any "X is mechanically possible" claim from a status doc gets a source check in the same session that acts on it.

## f) Up to 50 things we should get done next (brainstorm — ROADMAP fuel, not commitments)

**Release / v2 campaign path:**

1. Confirm or defer F058 (`/v2` module-path migration) — owner decision, everything else queues behind it.
2. TODO #484 owner call: `ImageProps` eager/lazy zero-value fix shape (v2 material).
3. S04/S05: finish `display` accent sweep (11 remaining blue sites; classify chart-hue keepers first).
4. S06: `forms` accent sweep (29 sites: focus rings, checked states, filters).
5. S07: `navigation` + `layout` accent sweep (10 + 2 sites; `--tc-sidebar-*` interplay audit).
6. S08: `errorpage` + `recipes` accent sweep (16 + 4 sites).
7. S09: collapse remaining tone lookup maps into `utils/tone` (TODO #431 execution).
8. S10: guard-test upgrade — scanners recognize `tc-` pairs; docs-count bumps; ADR-0044 status → Phase 2 implemented.
9. S11: surface/ink tables + shell conversion (Base/AppShell/Card/Minimal).
10. S12: surface/ink sweep pt2 (Table/Dropdown/Popover/Modal/inputs/Toast/CodeBlock/Kanban).
11. S13: `templ-components-theme.css` tc-surface/ink literal ramps + presets re-emit + warm-dark recipe rewrite.
12. S15–S18: brand-as-data batch (master-plan items; not yet detailed here — read plan before scheduling).
13. S24 + S27: remaining release blockers per master plan.
14. S28: website/theming page updates after S11.
15. S33: `/v2` module-path migration per ADR-0039 runbook (7 go.mods + go.work + guard-script literals + website/visualtest same commit).
16. Write `docs/migration/v1-to-v2-module-path.md` (ADR-0039 step 4).
17. S34: release-blocker checklist + full light/dark visual at RC + compiled-CSS grep guard.
18. S29: the v2.0.0 cut itself (release.sh ritual, dev shell + govulncheck).
19. Restructure `[Unreleased]` into "Ships in v1.x" / "Ships in v2.0.0" sub-headings.
20. Post-cut: GitHub Release page from notes + 24h watch (checklist items).

**Master hygiene (unblocked now):**
21. ~~Run ci-repro, then push the 3 commits~~ MOOT — the daemon pushed them at 07:46Z. NEW: watch in-flight run 42352650488 on `c10ffbe2` to green; if red, the masked-lane audit (row 22) becomes urgent.
22. Verify the pushed CI run exercises all 6 sub-module lanes (masking audit from §d4).
23. Diff-audit the daemon's 22-file commit `67a9537b` (FEATURES.md, website version-support, recipe changes) — only the 2 npm files were ever reviewed.
24. Triage 49 file-size-check findings (fix / budget / raise limits deliberately).
25. Fix 5 docs-gates findings (docs-freshness class) so pre-commit stops demanding --no-verify.
26. Decide BuildFlow debt row (rebuild binary + triage vs permanent bypass) — sibling tier0 open question.
27. Daemon gating improvement (§e1): cheap guard-set run before daemon push.
28. Status-doc freshness convention (§e7): superseded-by header or premise-expiry note.
29. TODO_LIST harvest of this section (docs-health HARVEST mode).
30. Record F058 decision outcome in ADR-0039 §Status + decision log once made.

**Smaller items noticed en route:**
31. `TODO #358`: GOEXPERIMENT=jsonv2 export drop — one lockstep commit, owner call pending.
32. `TODO #80`: vision-review pass over eyeball PNGs — still blocked on an API key export.
33. `TODO #491`: props-API consolidation bundle (v2 window).
34. Document the `-parallel 4` guidance vs site-route blank-capture flake (sibling item 40, visual-testing.md).
35. Decide whether the deleted Toggle bfcache e2e scenario should be referenced from component docs (sibling item 38).
36. Sibling recovery leftovers: verify `/tmp/unred-master` worktree pruned + `/tmp/art-dupl-pinned` cleaned (sibling items 36/37).
37. Confirm compiled `app.css` contains `data-[tc-palette-active=true]:bg-accent-50` (sibling item 34).
38. Classify remaining `display` blue sites (deliberate chart hues vs sweep leftovers) — cheap, unblocks S04/S05 scoping.
39. Consider a guard test that fails if `[Unreleased]` contains the string "BREAKING" without a "ships in v2" marker (mechanizes §d6).
40. Write the session's daemon-hijack incident into `docs/agent-context-history.md` (append-only) — 4th occurrence class this week.

**Longer-horizon (from ADRs/plans touched this session):**
41. S30–S32: fleet sweeps post-v2 (bridge deletion, dnsblockd contrast hole fix).
42. ADR-0048 ⫱G1 crossbrowser findings lane disposition.
43. TODO #501 follow-through check (accessible clickable-rows recipe consumers).
44. Kanban/heatmap/charts module extraction — trigger-gated NO (ADR-0040); re-check triggers only at v2.
45. v2 tag-tree rehearsal: `check-release-tags.sh` against the `/v2` prefix set before the real cut (S33 validation).
46. Golden re-baseline strategy for v2 (one batched rebaseline per plan Batch 4) — schedule early, it is the long pole.
47. Decide S14 (brand-as-data) scope boundary before S15-S18 to avoid sweep scope creep.
48. Owner-facing: draft the v2 consumer announcement (22 repos, import-path churn) when F058 confirms.
49. Re-run `art-dupl` re-record if v2 emission changes clone fingerprints (ADR-0009 binary-pin discipline).
50. After v2 cut: post-propagation tidy sweep + pkg.go.dev verification (release-checklist 24h watch rows).

## g) Questions I cannot figure out myself (3)

1. **F058, the standing gate:** do you confirm ADR-0039 option 1 — the next release is v2.0.0 WITH the `/v2` module-path migration — and if yes, is the completion of the emission sweep (S04–S13 + blockers) the trigger, or do you want a different sequencing (e.g., sweep lands first, migration decision later)? Everything in section (f) rows 1–20 queues behind this.
2. **The daemon already pushed everything** (07:46Z, unasked), bypassing the M03 ritual; CI on the fix commit is in flight. Policy call: do you want the ritual retroactively enforced — I run `scripts/ci-repro.sh --lint --css --website` at the pushed tip now and report `VERDICT` — or is CI-on-push your accepted gate going forward? Related: should the daemon's push privilege be revoked/narrowed (§e1), given it pushed a red tip minutes earlier?
3. **Pre-commit findings gate policy:** the gate is red on every manual commit (54 pre-existing findings) while the daemon bypasses hooks entirely. Keep the documented `--no-verify` interim for agent commits, or should a dedicated task (rows 24–26) be scheduled NOW so the gate returns to gating before the v2 campaign adds commit volume?

---

_Report uncommitted per harness contract (Crush: no commit without explicit instruction; the auto-commit daemon picks the file up). Point-in-time snapshot; NOT verified after 08:53 CET. Push-state correction appended 08:53 (daemon pushed c10ffbe2; CI in flight)._
