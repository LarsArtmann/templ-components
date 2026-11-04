# Status Report — 2026-09-13 12:12 CEST — Docs-Health Second Full Pass: August Annotations, Resurrection Fix, Living-Doc Sync

**Session scope:** the user ordered a second full docs-health run: "View ALL *_/2026-0_ files! Execute the
docs-health SKILL! … TODO_LIST, CHANGELOG, AGENTS, README, ROADMAP, and FEATURES must be all SUPERB! …
Archive FULLY done and UPDATED (inline strikethrough) .md files!" This report covers this session only.
Entry state: master at `83f91c82`, v1.16.0, `[Unreleased]` warm, 09-10 annotation pass complete with
known residuals (19 skipped August files, ~2.1k unmatched items). A concurrent website session was
active the whole time (docs-feature work under `website/`) — read, judged, left alone throughout.

**Verification at time of writing:** `go build ./...` PASS (all 7 modules + visualtest compile) ·
utils module `GOWORK=off go test ./...` PASS ×5 packages (includes `TestDocsCountDrift`,
`TestVersionMatchesChangelog/Features/ReadmeBadge`, `TestCompiledCSSInventory`, `TestSkillComponentCount`).
The templ-LSP kanban diagnostics are again confirmed FALSE — build is ground truth.

**Commits (mine):** `2a57862a` (gitignore hardening + changelog) · `579c6e87` (guard asserts tracked set)
· `f8e8ea91` (living-doc sync + harvest) · `afa7157c` (TODO #212 policy row). The annotation edits and
the ghost-fix rode daemon commits (`0663a21c`, `a0f630ba`, `2131c7b9`) because my ghost-fix commit's
BuildFlow hook panicked (see d3). Working tree holds only the concurrent website session's files.

---

## a) FULLY DONE

| #   | Item                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | Evidence                                           |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------- |
| a1  | **Inventoried all 260 `2026-0*` files** (216 `.md`) across docs/{status,planning,feedback,reviews,proposals,research,architecture-understanding,modularization}; read both fresh 09-13 reports and the 09-10 pass report as entry context                                                                                                                                                                                                                                                                                                                                                                              | session transcript                                 |
| a2  | **19 skipped August status reports (2026-08-10…08-17) annotated** — 4 parallel classification agents produced per-item verdicts under strict evidence rules; I mechanically filtered them (h with validated hashes, w from the report itself, positive-v only, p→h upgrade inside completion-titled sections using each file's adding commit) and applied **629 inline verdicts**. **53 specs deliberately dropped** as unverifiable/negative — a false done is worse than an open item                                                                                                                                | daemon commits `0663a21c`, `a0f630ba`              |
| a3  | **Built + used a heading-aware annotator** for the two shapes the shipped scripts cannot address: `### N.` heading items and `###`-level section scoping (the shipped scripts bound only at `##`). Same strike semantics, atomic write, read-back check. All 19 files now carry inline verdicts                                                                                                                                                                                                                                                                                                                        | /tmp tool; diffs verified 18/18, 21/21 in-place    |
| a4  | **Caught + fixed a RED master — the predicted daemon resurrection happened.** Six deleted `.out.css` artifacts returned in daemon commit `504856b8`, ~90 minutes after `43f1522f` deleted them. `TestCompiledCSSInventory` failed exactly as designed with the remediation message. Fix: re-deleted, added `*.out.css` to `.gitignore` except the one tracked distribution artifact, and **rewrote the guard to assert the git-tracked set** (`git ls-files`, fail-loud on missing git, untracked worktree litter reported informationally) — the old disk-walk guard false-fails on every daemon tailwind-build cycle | `2a57862a`, `579c6e87`; guard PASS confirmed twice |
| a5  | **Corrected a false done-marker in a historical report**: the 2026-09-08_22-36 row claiming `scripts/daemon-regression-check.sh` was struck "done — visualtest lint zero N4", but the script never existed (verified `scripts/` empty of it). Re-marked **NOT-DO** citing what actually shipped: `check-css-minified.sh` (Guard 6), `ci-repro.sh --css`, TODO #120, root cause TODO #125. This closes the CSS-audit report's c10 item (whose pointer to the wrong 09-08 report I also corrected)                                                                                                                       | `2131c7b9`                                         |
| a6  | **TODO_LIST overhauled**: header date 09-09→09-13; next-free-ID 193→**212** (was stale with #193–#202 live); fixed the **broken markdown table** (rows #190–#192 floated as a headerless fragment) by merging them into the Blocked table; removed an empty harvested section (CV adoption, all rows already dropped); added **9 actionable rows (#203–#210)** and **2 owner decisions (#211 templates-CSS fate, #212 unmatched-annotation policy)** — every row citing source report + section                                                                                                                        | `f8e8ea91`, `afa7157c`                             |
| a7  | **ROADMAP gained a 9-direction depth-testing cluster** (harvested from the 100-idea review, verified absent before adding): axe-core gate, HTML validation over 242 goldens, determinism render-twice gate, golden-orphan detector, Firefox visual lane, forced-colors/contrast variants, chart-geometry property tests, wire fuzzing, go/analysis convention linter                                                                                                                                                                                                                                                   | `f8e8ea91`                                         |
| a8  | **Fix-on-sight trio from the CSS-audit report**: (b1) `TestCompiledCSSInventory` row added to the skill guard table; (b2) stale "`tc new` tests +34" row pruned from `docs/cross-project-analysis.md` after a live-docs audit found no other siblings; (b5) confirmed a NON-issue — `[Unreleased]` section order already matches Keep-a-Changelog                                                                                                                                                                                                                                                                      | `f8e8ea91`                                         |
| a9  | **AGENTS.md CSS-inventory bullet synced** to the new reality (tracked-set guard, gitignore rule, resurrection history) — the bullet previously promised the old disk-walk behavior                                                                                                                                                                                                                                                                                                                                                                                                                                     | `f8e8ea91`                                         |
| a10 | **Verified, not assumed**: my earlier session-scanner regex was wrong (reported 0 struck everywhere); re-verified against a known-annotated file before trusting any count. Agents produced contradictory evidence on `cmd/tc/_sources` existence — spec dropped rather than struck on contradiction. FEATURES confirmed release-pinned (unreleased features land at the v1.17.0 docs pass)                                                                                                                                                                                                                            | session transcript                                 |
| a11 | **Inline health report delivered** in-conversation (Accuracy 9.5 / Fitness 9.5, findings, visible math) per the docs-health AUDIT format; no report file written for it by design                                                                                                                                                                                                                                                                                                                                                                                                                                      | conversation                                       |

## b) PARTIALLY DONE

| #  | Item                                                                                                                                                                                                                                                                                                                                                 | Gap                                            |
| -- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| b1 | **The 19 August files are substantially but not fully annotated.** Every resolvable item is struck; the residue is f-list idea rows that are genuinely still open (correctly left unstruck) — so **none of the 19 became archive-eligible** under the "EVERY item resolved" rule. +0 archives this pass; the 26 from the 09-10 pass stand            | archive awaits genuine resolution, not markers |
| b2 | **Daemon litter is git-invisible but still regenerates.** `.gitignore` stops it entering master (guard green, tree clean), yet every commit hook's tailwind-build step re-creates the 6 zombie `.out.css` files on disk — confirmed re-grown at report time. Worktree hygiene only; upstream BuildFlow is the real fix (#125 family)                 | one `trash` pass per session until upstream    |
| b3 | **The annotation tooling lives in `/tmp`.** My heading-aware annotator + filter driver are throwaway; the next docs-health pass cannot reuse them. The shipped skill scripts still can't scope `###` sections or strike `### N.` heading items, and their read-back pattern mis-trips on `w` markers (my port fixed all three; fixes not upstreamed) | skill-asset PR or repo script decision         |
| b4 | **~600 agent verdict lines were hand-transcribed** into the driver's RAW block. Mechanical filters + fail-loud + dry-runs caught every defect (zero mis-strikes), but the risk class — re-typing machine-readable verdicts — is exactly what the 09-10 report's §e1 told us to eliminate. Still no durable agent→driver pipeline                     | process work, not heroics                      |
| b5 | **Verification was targeted, not canonical**: build + full utils suite ran; I did NOT run `scripts/ci-repro.sh --lint`, the full per-module test loop, or `nix run .#verify`. Justified mid-session (docs + one utils test file), but it's the pre-push gate and remains open                                                                        | CI will re-prove; run before push              |

## c) NOT STARTED

| #  | Item                                                                                                                                                                                                                                                                      | Why on the list                                                                                           |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| c1 | Second unmatched-annotation pass for the ~2.1k name-keyed items                                                                                                                                                                                                           | now explicitly owner-gated via TODO **#212** (map IDs / report-level resolution / declare open-by-design) |
| c2 | Upstreaming the annotate-script extensions (`###` scoping, heading items, w-marker read-back) into the docs-health skill assets                                                                                                                                           | skill repo decision; tool currently /tmp-only                                                             |
| c3 | `scripts/ci-repro.sh --lint` + full per-module loop + `nix run .#verify`                                                                                                                                                                                                  | pre-push gates, untouched this session                                                                    |
| c4 | v1.17.0 release cut (#192)                                                                                                                                                                                                                                                | owner decision; `[Unreleased]` is warm with wire.DecodeForm, hooks, guards, resurrection fix              |
| c5 | Everything else in the standing backlog — BuildFlow upstream family (#93/#107/#108/#124/#125/#126), vision review with a real key (#80/#150/#162), kanban Q1/Q2 (#190/#191), #189, adoption surveys (#155–#157), route-golden matrices (#194–#196), the 09-10 f-list tail | untouched by design — this session was docs-only plus one guard fix                                       |
| c6 | Website session's in-flight docs feature (many `website/` files churning all session)                                                                                                                                                                                     | another actor's work; coordination/verification deliberately not started                                  |

## d) TOTALLY FUCKED UP (honest ledger)

| #  | Item                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | Detail                          |
| -- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| d1 | **My first freshness scanner was buggy** (`~~`-regex missed `1. ~~` and table shapes) and reported 0 struck across ALL files — including files the 09-10 pass had provably annotated. Cost one wasted round; recovered by cross-checking a known-annotated file before acting on the number. Lesson already in AGENTS.md ("independently verify tool output") — applied, but the regex should never have shipped                                                                 | ~1 round trip                   |
| d2 | **The /tmp annotator's read-back guard false-tripped on `w` markers** (`"done" not in "Won't implement"`), aborting my batch chain mid-file. Worse: I first misattributed the failure to the wrong batch (index 108 belonged to a different file than I assumed) and burned several debugging rounds before instrumenting and finding the real cause. The WRITES were always correct (git-verified), but the debugging was slow and the chain-order confusion was self-inflicted | ~4 wasted rounds                |
| d3 | **My ghost-fix commit panicked the BuildFlow pre-commit hook** (stack trace out of `larsartmann/buildflow/cmd/buildflow/main.go`), so my commit aborted and the daemon swept the file instead — landing it mixed with the concurrent website session's unrelated files under a generic auto-commit message. Content is correct and verified in-tree, but attribution is the daemon's, and I did NOT investigate the panic (upstream binary, budget call)                         | attribution mess; panic unfiled |
| d4 | **Three edit-tool refusals** (CHANGELOG, AGENTS, cross-project-analysis) for editing before viewing in-session. The tool is right; I paid 3 round trips to save 3 quick views                                                                                                                                                                                                                                                                                                    | ~3 round trips                  |
| d5 | **Trashed litter BEFORE my last two commits instead of after** — each commit's hook re-ran tailwind-build and re-grew it. Confirmed re-grown at report time. Harmless (gitignored) but my cleanup ordering was backwards                                                                                                                                                                                                                                                         | worktree hygiene only           |
| d6 | **First draft of the driver carried a dead `run()` helper and a placeholder line** — sloppy first write, cleaned in the same pass before any run                                                                                                                                                                                                                                                                                                                                 | zero runtime impact             |
| d7 | **Used a python heredoc for the TODO #212 row insert** — sanctioned for docs tables by the AGENTS carve-out and verified with `git diff`, but `edit` was available and lower-risk; the carve-out is a seatbelt, not a preference                                                                                                                                                                                                                                                 | zero damage                     |

Net: zero fucked-up state survived in the repo; cost was ~10 wasted tool rounds, one misattributed debug session, and one daemon-mixed commit.

## e) WHAT WE SHOULD IMPROVE

1. **Durable agent-verdict pipeline.** Classification agents should write spec files the driver consumes verbatim — this session re-typed ~600 verdict lines by hand, the exact risk class the 09-10 report flagged. My filters + fail-loud + dry-runs made it safe this time; the process should make it safe by construction.
2. **Upstream the three annotate-script gaps**: `###`-level section scoping (boundary = section heading level, not "next `##`"), `### N.` heading-item support, and a read-back guard that accepts `w` markers. The shipped scripts failed safe on all three — that's why nothing was corrupted — but they can't address these shapes at all today.
3. **Archive-eligibility + freshness scanner as repeatable scripts** (09-10 f6/f7, still not built). This pass re-derived eligibility by hand; a `docs-health-check <file>` would have made the "none eligible" verdict mechanical.
4. **Resolve contradictions immediately.** Two agents disagreed on `cmd/tc/_sources` existence; I dropped the spec (safe) when a 5-second `ls` would have settled it. Dropping is the right default; verifying is a cheap upgrade.
5. **Commit-vs-daemon race discipline.** Twice my edits were swept by daemon commits; once a hook panic made the daemon the committer of record. Land commits immediately on green (I did), and re-check `git log` attribution after any hook failure (I didn't — the daemon had it).
6. **Cleanup AFTER the last mutating step.** Trash worktree litter post-commit, not pre-commit — hooks regenerate it.
7. **The "status-report `.md` override" keeps recurring** (09-10 §e8, count now higher). Record the standing order in the skill so sessions stop re-flagging it.

## f) NEXT 50 (ranked; this session's vantage — grouped)

**Decisions needed (owner)**

1. Cut **v1.17.0** via release script (#192; `[Unreleased]` warm: wire.DecodeForm, tracked hooks, count guards, resurrection fix).
2. **#211**: fate of `templates/styles.css` + `templates/templ-components-theme.out.css` — delete like the other ten, or document the consumer path (g1 below).
3. **#212**: policy for the ~2.1k unmatched items — per-file ID mapping, report-level resolution, or declare open-by-design.
4. **#190/#191**: kanban touch-drag + feature-scope Q1/Q2, then execute.
5. **#123**: branch protection — decide, then wontfix formally instead of leaving it parked.
6. Vision review provider key + cost approval for the 23-image flagged set (#80/#150/#162).
7. Upstream or accept: annotate-script extensions + the /tmp tooling (b3/e2).

**Direct follow-ups from THIS session (small, this week)**
8. **#203** lint baseline: drive the 8 named findings (gocognit ×6, golines, unused-nolint) to 0 or nolint-with-reason.
9. **#204** extend minification guard to the two `templates/` targets (fold into the Go guard).
10. **#205** single-source the compiled-CSS target list (release.sh ↔ guard test).
11. **#206** disambiguate the two same-named `templ-components-theme.css` files.
12. **#207** "committed artifact ⇒ named consumer" rule into CONTRIBUTING.md.
13. Re-trash the worktree `.out.css` litter (re-grown by hook cycles; confirmed at report time).
14. Watch one full daemon cycle end-to-end: confirm litter stays git-invisible and the guard stays green (b3 watch).
15. Run `scripts/ci-repro.sh --lint` before next push (c3).
16. Run `nix run .#verify` once before the next release cut (c3).
17. File/report the BuildFlow pre-commit **hook panic** observed in d3 (upstream issue material).
18. Correct the CSS-audit report's c10 pointer (it cites the wrong 09-08 report; the ghost lives in `_22-36`) — one-line annotation.

**Annotation/tooling continuation**
19. Second unmatched pass (per #212 outcome).
20. Archive-eligibility checker script (09-10 f6).
21. Freshness scanner as a repo/skill script (09-10 f7).
22. Migrate the 19-file verdict TSVs into a durable location for future re-verification.
23. Reconcile AGENTS "display=43" vs SKILL per-package sums (CSS report f31) — guards cover totals; per-package headings need one pass.
24. Audit remaining planning docs for `tc new` promises beyond cross-project-analysis (CSS report f32).

**Quality/CI backlog (top of the standing queues)**
25. **#208** post-tag consumer compile smoke in CI.
26. **#209** hook-adoption guard (`core.hooksPath` check).
27. **#210** `wire.DecodeForm` adoption: forms-pack + demo handlers.
28. **#193** pollBool/pollText visualtest helpers.
29. **#194–#196** route-golden dark/375px/RTL matrices.
30. **#197** keyboard-only demo traversal.
31. **#198/#199** kanban LSP-warning triage + gopls stale-diagnostics writeup.
32. **#200** Datastar innerHTML-no-scripts fact into datastar doc.go.
33. **#201** ci-repro changelog-guard parity.
34. **#202** visualtest `.fail/` disk hygiene.
35. BuildFlow upstream sprint: #93/#107/#108/#124/#125/#126 in one push (#40 of the 100-idea review).
36. `*.out.css` skip-list upstream in BuildFlow's tailwind-build provider (now proven twice).

**Feature/consumer backlog (already tracked — keep ranked)**
37. **#189** kanban demo niceties (file-backed state, dashboard section).
38. **#155** SimpleNav `Wire` candidate survey (D3 rule).
39. **#157** Calendar month-nav `Wire` survey.
40. **#156** consumer adoption gaps: AppShell CSS-var theming/breakpoint, Minimal head-content.
41. **#178** typed wire trigger language ADR.
42. **#133** changelog-guard real-PR shakedown.
43. **#39** compound overlay API (ADR-0023) at v2.
44. **#34/#33** testutil extraction + Validate() methods (deferred, keep parked).
45. Push queued commits / reaffirm daemon-push as the norm.

**Docs/site tail**
46. Re-derive the 100-idea HTML report's stat cards or stamp "at generation time" (10-38 b2).
47. Website: sync the new testing tiers into site content (09-10 f22).
48. Website docs-feature session: verify the concurrent in-flight work lands green (noticed all session).
49. **#119-note**: remove the bun shim blocking pnpm (user-level fix).
50. 30-day review: if the inventory guard never fires again, consider demoting severity (CSS report f50).

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **Do you know of anyone actually consuming the pre-compiled `templates/styles.css` / `templates/templ-components-theme.out.css`** (e.g. raw GitHub fetches, CDN-style usage), or should they join the ten deleted artifacts? release.sh treats them as distribution targets, but no in-repo consumer exists — this decides TODO #211 and whether #207's "named consumer" rule gets its first exception.
2. **Where should the annotation tooling live** — upstream the `###`-scoping/heading-item/w-guard fixes into the docs-health skill's shipped `annotate-*.py` (your skill repo), or land them as a repo script under `scripts/`? It's currently a `/tmp` throwaway; the next docs-health pass needs it to exist somewhere durable.
3. **The 09-13 ghost-elimination report claims you said you don't care about branch protection (#123).** Is that accurate? If yes I'll formally wontfix #123 with that reason instead of leaving it parked as "needs repo-owner decision".

---

_Point-in-time snapshot — 2026-09-13 12:12 CEST — master `b7f54d69` (daemon), tree clean except the concurrent website session's files, all module builds + utils guard suite green. Verify before acting on any claim (per repo policy)._
