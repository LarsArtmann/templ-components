# Status — Pareto plan full execution (21/21), corrections, and honest gaps

**Date:** 2026-11-02 16:27 CET · **Session:** pareto-plan full-execution mode over
`docs/planning/2026-11-02_00-23_PARETO-BACKLOG-TRUTH-AND-SAFETY-PLAN.md`
**Tip at report time:** `0c1d6bb5` (master == origin/master; the M03 ritual ran
`VERDICT: PASS (exit 0)` at this exact tip before it was pushed — the daemon pushed it,
the invariant "pushed tip == verified tip" held).
**⚠ Live caveat:** a `flake.lock` working-tree modification appeared AFTER the ritual
verdict (not authored this session — daemon/sibling race). Uncommitted, unverified by
the ritual. Never reverted (not mine); flagged here so the next session triages it
before any push.

---

## a) FULLY DONE (executed + verified)

| Task                                              | Deliverable                                                                                                                                                                                                                                                                                    | Verification                                                                                                                                                |
| ------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| T2 + T19 — dependency pin-train guard (TODO #995) | `utils/version_train_test.go`: `checkPinTrain` asserts pin SURFACE (exactly the known go.mods) + pin TRAIN (all versions equal) for go-error-family and tailwind-merge; negative-tested against the exact shipped-skew shape + surface/missing/replace-line cases via a recording `testing.TB` | mutation test (wrong display-count analog) proved both detection layers fire; full utils suite green; TODO #995 struck                                      |
| T3 — count-table codegen (TODO #996)              | `tools/countgen` (+ `cmd/countgen`): Derive/Check/Apply over 22 count sites in FEATURES.md + skill/SKILL.md; `TestDocsCountDrift` extended to the previously unguarded cells (4 SKILL headings, FEATURES table rows, recipes clause, "nine counted packages" word, errorpage constructors)     | check mode reported ZERO drift pre-switch (byte-parity proven); `-w` proven a no-op; both layers mutation-tested; suites green `-count=1`; TODO #996 struck |
| T4 — erraudit sweep                               | `scripts/erraudit-sweep.sh`: per-module fix dry-run + `legacy_as` GATE + `legacy_is` advisory; `--fix` opt-in                                                                                                                                                                                  | real tree PASS (matches known 0-findings state); scratch-copy injection test went RED naming `scratch_as.go:11:5`; AGENTS one-liner added                   |
| T13 — house-style doc                             | `docs/house-style.md`: intent-comment format for deliberate empty catches (3 canonical examples) + pipe-exit-code/PIPESTATUS rule with the 2026-08-14 misread citation                                                                                                                         | citation verified at source line                                                                                                                            |
| T5 — stale-review annotations                     | 2026-05-19 architecture `.md` (appendix), 2026-07-18 modularity HTML + 2026-07-18 data-model HTML + 2026-10-15 color-domain HTML (inline banners)                                                                                                                                              | cross-link claims verified against the 2026-11-01 sources ("60 typed enums", "90 IsValid guards"); tag-balance structural check OK on all 3 HTMLs           |
| T11 — provenance                                  | modularization README Files section indexes the 2026-11-01 assessment + ADR-0039/0040; `agent-context-history.md` gains §pinskew                                                                                                                                                               | ADR filenames verified before citing                                                                                                                        |
| T12 — bump protocol                               | `docs/external-dependency-bumps.md`: new step 6 (`go work sync` + go.work.sum + standalone per-module builds), step 3 cites the pin-train guard                                                                                                                                                | §pinskew anchor resolves both directions                                                                                                                    |
| T8 — dependabot hygiene                           | templ + go-error-family `ignore` in all 9 gomod blocks (yaml-validated); 9 open templ-0.3.1100 poison PRs closed with policy comments (#40/41/43/44/45/67/68/71/72); orphan refs deleted; zero dependabot refs remain                                                                          | `yaml.safe_load` semantic check; actionlint clean; comment posted on superseded PR #64                                                                      |
| T9 — erraudit SKILL fix (separate repo)           | `~/projects/SKILLS` go-error-modernization SKILL.md: `fix` takes `-w` not `--write`, no `--type-aware` on fix, table corrected, non-Go scanning note — commit `f7528e0`                                                                                                                        | fanout serves the new content (read through the symlink); `check-skill-fanout.sh`: 51 symlinks clean                                                        |
| T6 — templ window plan                            | `docs/planning/2026-11-02_TEMPL-MIGRATION-WINDOW-PLAN.md`: gate states, preconditions, ordered ONE-event recipe, #368/#369 bundling, per-stage rollback, exit checklist; NO preemptive nixpkgs-templ split (folded-back-input precedent cited)                                                 | claims cross-reviewed against `agent-context-history.md#templ-version`; TODO #335 row points at it                                                          |
| T7 — trigger dashboard                            | `scripts/display-split-triggers.sh` (counts 53/49, T1 unmet, afferent blast radius 3 packages) + ROADMAP v2 extraction order charts→overlays→tables                                                                                                                                            | two dashboard bugs found and fixed by its own first run (templ/runtime false positive; self-import in afferent count); chart_geometry.go existence verified |
| T17 — long-tail evaluations                       | contingent display-split D2 + compiled SVG; tooling-evaluation §8: erraudit LSP viable-not-wired; `--max` sweep = 13 generic-return findings all ACCEPTED (0 nolint suppressions repo-wide, verified)                                                                                          | SVG compiled (d2 v0.9.0); counts verified live                                                                                                              |
| T14 — owner packet v2                             | `2026-11-02_OWNER-DECISION-PACKET-V2-REFRESH.md`: 12 gates re-checked, deltas only, defaults standing                                                                                                                                                                                          | awesome-templ#12 re-verified OPEN; key absence verified in env                                                                                              |
| T15 — demand re-check                             | row 217: FileDrop 1→2 consumers (FLIP-ready), MultiSelect WATCH (2 mixed sites), TreeView HOLD                                                                                                                                                                                                 | hand-roll sites read at source in KeyHolderAI / DiscordSync / nsfw-classifier / standard-bug-tracking-schema                                                |
| T16 — backlog truth                               | #410 struck (7 v1.23.0 tags verified), #191 struck (decision record; v1.18.0 claims verified), #436 verified STILL PENDING with fresh evidence; 5 shipped-version claims spot-checked                                                                                                          | CHANGELOG headings + tags grepped                                                                                                                           |
| T21 — harvest                                     | row #997 harvested (release-time route-golden re-record); changelog warm                                                                                                                                                                                                                       | —                                                                                                                                                           |
| M03 push ritual                                   | `scripts/ci-repro.sh --lint --css --website` at `0c1d6bb5` → `VERDICT: PASS (exit 0)`; origin/master in sync with the verified tip                                                                                                                                                             | full log in `.crush` shell output; art-dupl ADVISORY noted (see d/e)                                                                                        |

## b) PARTIALLY DONE

1. **T18 vision review (#80/#150/#162)** — gate probe only: no vision API key in the
   session env; recorded in row 80 + packet v2. The script is ready; the run never
   happened. Caveat: I probed `env` in a non-direnv shell — if a key lives only in
   `.envrc`, my "absent" verdict could be false (unresolved).
2. **T15's registry leg** — the shadcn-templ llms.txt re-check was SKIPPED (no recorded
   URL in the 2026-10-08 doc; URL-guessing banned). Sibling-repo hand-roll search only;
   a GitHub-wide code search would strengthen the FileDrop/MultiSelect evidence.
3. **T14 packet v2 scope** — 12 gate sections refreshed with deltas; the plan's "26 ⫱
   rows" are covered via v1 pointers, not individually re-batched row-by-row. Thinner
   than the plan text; argued sufficient (evidence-unchanged rows), not proven sufficient.
4. **art-dupl advisory (dupl lane)** — final ritual printed
   `art-dupl: ADVISORY-FAIL (exit 2) — new clone groups outside baseline` (the new
   countgen error-wrapping similarity, most likely). Non-blocking, but UNCLASSIFIED:
   nobody has done the ADR-0009 extract-or-accept+rebaseline triage yet.
5. **SKILLS-repo fix** — committed (`f7528e0`), NOT pushed (never-push-unasked). The
   fanout serves the fix on this machine; other machines don't have it until pushed.
6. **countgen's own testing** — the generator has no Go unit test of its own; it is
   pinned indirectly (docs-count guard + mutation test done by hand this session).
   `tools/countgen/countgen_test.go` with synthetic-tree cases is still unwritten.

## c) NOT STARTED (placed, gated, or explicitly out of scope)

- **Q1 browser baseline + the rest of the owner packet asks** (C19/C22–C24 unblock
  chain) — owner decisions; packet v2 carries them.
- **#269 go-error-family upstream PR** — branch prepared in `~/projects/go-error-family`;
  needs owner go-ahead to push + open.
- **#997** — wire route-golden re-record into release.sh + release-checklist (harvested
  this session, not implemented).
- **#368 (size-*) / #369 (@utility)** — DECIDED-DEFERRED into the templ window
  (templ gate still CLOSED; re-probe after the next `nix flake update`).
- **FileDrop build row** (#217 FLIP-ready) — routed, not picked up.
- **Blocked rows 93/107/108/384/385** (BuildFlow repo), **#28** (upstream review) —
  separate repos by design.
- **Deferred-v2 catalog** — truth-passed, parked until v2.
- Anything requiring the still-unlanded sibling work: `display/rating.templ` never
  landed; if that session resumes, `TestClassInventoryFreshness` + the pixel suite
  need a re-run before trusting goldens.

## d) TOTALLY FUCKED UP (own failures, no masking)

1. **Shipped 3 commits carrying 10 lint findings.** My direct golangci-lint runs timed
   out under machine load and I gave up instead of retrying periodically — the findings
   (gocognit, thelper, nonamedreturns, unparam, wsl_v5, golines, gofumpt,
   embeddedstructfieldcheck) only surfaced at the M03 ritual, producing a red lane and
   a late fix commit (`0c1d6bb5`). The ritual did its job; my per-commit discipline
   didn't. House rule "fix so BOTH versions pass" was satisfied eventually, not early.
2. **countgen sank the coverage gate.** Adding `tools/countgen` (untested) without
   checking the coverage script's exclusion list caused a red `84.9% vs 94%` at the
   ritual — the scope fix + baseline re-record (`295b3318`) came after the fact. I
   should have checked "does my new package enter any measured scope?" BEFORE building
   under it.
3. **Two dashboard bugs shipped in T7's first version** (templ/runtime false positive;
   display counting itself as its own importer) — caught only by running the tool and
   reading its output. The initial inventory claim in my own guard comment
   ("charts/echarts legitimately absent") was also wrong — the pin-train guard itself
   corrected me. Net: the guards work, but my first-draft claims weren't verified
   before writing them into comments.
4. **A botched multiedit left orphaned code fragments** in `version_train_test.go`
   (my old_string missed the file's post-gofmt state), requiring a scripted region
   rewrite to recover. Compiling prevented damage, but the "read the file immediately
   before editing" discipline slipped after formatter interference.
5. **Missed the CHANGELOG/ritual ordering once:** the first ci-repro red (coverage) and
   the second red (lint) each produced a "done" narrative beat before the ritual
   actually passed. No masking shipped — every red got an honest verdict + fix — but
   twice I narrated completion ahead of verification.

## e) WHAT WE SHOULD IMPROVE

1. **Lint early, lint often:** run `golangci-lint` on touched packages right after each
   commit (queue retries under load) instead of letting the ritual be the first lint.
2. **Scope check before new top-level dirs:** any new directory (like `tools/`) must be
   checked against every measuring script's include/exclude lists (coverage, dupl,
   file-size, docs scans) in the same commit.
3. **countgen deserves real unit tests** (synthetic-tree parity, site-match failure
   modes, Apply idempotence) — currently only indirectly pinned.
4. **Hook-starvation protocol:** on `environment_starvation`, sleep-and-retry on a
   schedule instead of ad-hoc; consider `--max-concurrency` per the tool's own hint.
5. **The daemon doc-eating race** (planning/status docs deleted-then-restored inside
   `13af2e67..45da26d6`) deserves a §-entry in agent-context-history — noticed,
   unresolved, unrecorded.
6. **`flake.lock` race window:** the verdict-covered tree and the pushed tree can drift
   within minutes when the daemon/siblings commit; consider re-running a fast lane
   (lint-only) after any post-verdict movement, or making the daemon hold flake.lock.
7. **Demand-gate method:** record the registry URL in the recheck doc so future passes
   can re-fetch; add a GitHub code-search pass to the standard method.
8. **Packet v2 pattern:** next refresh should re-batch all rows individually (or state
   explicitly why unchanged rows don't need re-batching).
9. **Demo-server port hygiene:** `PORT` env + `shots -base` collision burned 3 attempts;
   document in docs/visual-testing.md (and consider shots auto-failing fast on 8901
   default when demo runs on 8080).
10. **AGENTS.md candidate facts** (not yet written there): tools/ coverage exclusion is
    only in the script; the shots/PORT gotcha is nowhere; both are ≤1-line adds within
    the 220-line budget if something else is demoted.

## f) Up to 50 things to get done next (impact-sorted, not all mine)

**Code/guards**

1. Triage the art-dupl advisory: classify countgen's new clone groups per ADR-0009
   (accept + re-record baseline with the pinned fork binary, or extract a shared
   error-wrap helper).
2. Write `tools/countgen/countgen_test.go` (synthetic-tree parity, multi-match
   structural-drift error, Apply idempotence, numberWord bounds).
3. Execute #997: add route-golden re-record to release.sh verify + release-checklist.
4. Extend countgen to README per-package headings (deletes the remaining hand-typed
   per-package tax; guard loop already covers 5 of them).
5. Wire `erraudit-sweep.sh` as a CI/ci-repro lane (currently a manual command).
6. Add a cheap countgen Check lane to pre-commit/CI for docs edits.
7. Decide a staleness policy for the erraudit LSP, then wire it into editor config (§8a).
8. Guard idea: fail when route goldens predate the last `utils.Version` bump
   (mtime-fragile — design carefully) or fold into release.sh (see #3).
9. Re-probe `scripts/check-templ-gate.sh` after the next `nix flake update`.
10. Triage the live uncommitted `flake.lock` modification before the next push.
    **Owner gates (packet v2)**
11. Q1 browser baseline decision (unblocks C19/C22–C24, #473/#476–#482).
12. Export vision API key → run `vision-review-goldens.sh` → human-confirm SUSPECTs
    (retires #80/#150/#162).
13. #269: go-ahead to push the prepared go-error-family branch + open the PR.
14. #211 templates fate; #212 annotation-policy scope; #233/#234/#237/#238/#239; #329;
    #358; #411 — packet v2 defaults stand until overridden.
15. FileDrop: green-light the build row (FLIP-ready, 2 consumers) → full ladder.
    **Docs/skills**
16. Push the SKILLS-repo commit `f7528e0` (flag corrections live for other machines).
17. Update the templ-components SKILL "add a component" ladder: the count step is now
    `go run ./tools/countgen/cmd/countgen -w` (check every mention of hand-bumped counts).
18. Record the doc-deletion daemon race (§ note) in agent-context-history.
19. Document the PORT/shots gotcha in docs/visual-testing.md.
20. Add the tools/-coverage-scope fact to AGENTS.md if the budget allows (or leave the
    script comment as canonical and link it).
21. Append a completion appendix to the 00-23 plan file (mark T1–T21 executed).
    **Follow-ups from this session's discoveries**
22. Verify dependabot's re-proposed PRs come back WITHOUT templ (validate the ignore).
23. Re-run `TestClassInventoryFreshness` + pixel suite if the rating.templ session lands.
24. MultiSelect: GitHub code-search pass to firm up (or refute) the WATCH posture.
25. Re-verify awesome-templ#12 at the next session start (was OPEN).
26. Reconcile the countgen "74 enums" hand-typed total IF a mechanical source ever
    emerges (documented as deliberately manual for now).
27. BuildFlow housekeeping (owner/external): restore the dead `/mnt/buildcache` mount;
    VACUUM the 0.69 GB state DB / 0.22 GB cache DB.
28. Re-record coverage baseline only at release boundaries (avoid sibling-churn drift);
    note the toolchain-patch rule is already documented.
29. Consider `--max-concurrency` for ci-repro when the box is loaded (tool hint).
30. Prune `visualtest/testdata/.fail/` litter (gitignored, but disk).
    **Backlog hygiene (from the session's TODO work)**
31. #436: chase the two consumers (nsfw-classifier, pbx-lens) to actually regen CSS —
    verified still stale 2026-11-02.
32. #217: re-check Toast positions (explicitly out of T15's scope this pass).
33. #425/#431 (tone collapse): stays v2-windowed; keep the pending cross-links current.
34. #444: per its own blocker text (untouched, unassessed this session).
35. Rows 124–126 / 190 / 232–239: individually re-batch in packet v3 if v2's pointer
    approach feels too thin.
    **Infrastructure/tooling candidates**
36. Evaluate wiring countgen Check into the BuildFlow fast lane (sub-second).
37. erraudit: consider `--enforce-go-error-family` posture decision (off today;
    repo does use the family — a deliberate choice to make, not drift).
38. Keep golangci-lint pin at 2.14.0; watch the dev-shell pin (TODO #443 resolved;
    stay alert).
39. Consider exporting `TC_ERRAUDIT_ROOT`-style target override in erraudit-sweep for
    scratch-tree testing (self-test currently copies the repo).
40. `nix run .#shots`: fail fast with a clear message when the demo server is down
    (today: connection-refused archaeology).
    **Nice-to-haves**
41. countgen: emit a `--json` summary for CI annotations.
42. house-style.md: add the "TRUST the hook, retry on starvation" operational note.
43. house-style.md: cite the T7 dashboard self-correction as the "run the tool before
    writing its claims" example.
44. SKILLS: generalize the flag-drift lesson into references/lessons.md (cross-project).
45. roadmap: the "Counts truth table" row is now HALF-shipped (countgen) — update it.
46. ROADMAP "since v2.0" phrasing on rows 27–28 reads as shipped although v2 is future —
    verify and reword (noticed, not fixed).
47. Dead-code sweep: `utils/dismiss.templ`'s exported templ function vs the "utils owns
    no components" claim — either document or fold (countgen documents it today).
48. Website landing: consider surfacing the new guarded-count provenance (CountStats +
    countgen) in a footer note.
49. Pin the shots tool's base-URL default to the demo's real port or read PORT.
50. Next session: start from packet v2 + this report; do NOT re-derive T1–T21.

## g) Three questions I cannot answer myself

1. **Vision API key:** is there an API key I should use for the T18 vision review
   (#80/#150/#162), and how should it reach the script (which env var name, exported
   where)? My "absent" probe ran in a non-direnv shell and could be a false negative.
2. **SKILLS push:** may I push the SKILLS repo (`f7528e0`, erraudit flag corrections)?
   It's committed and verified on the fanout here, but other machines don't have it
   until pushed — and I don't push without your say-so.
3. **FileDrop priority:** row #217 is FLIP-ready on two hand-rolled consumers
   (DiscordSync + nsfw-classifier). Do you want it routed to a build row now, or does
   it wait for the next demand-survey pass like the other gated candidates?
