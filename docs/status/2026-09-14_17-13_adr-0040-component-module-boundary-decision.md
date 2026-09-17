# Status Report — 2026-09-14 17:13 CEST — ADR-0040 Component Module-Boundary Decision

**Session scope:** go-modularize 7-phase re-assessment answering "should kanban, heatmap, … be
their own modules?" — plus this self-review/status report. Nothing else was touched.

**Verdict delivered:** NO split. KanbanBoard, Heatmap, and the native chart family stay in the
root `display` package. Decision is trigger-gated (T1 unique-dep / T2 consumer-demand / T3
own-build-pipeline), recorded in ADR-0040, timed with v2 per ADR-0039.

**Session artifacts (as committed by the auto-daemon, all unpushed):**

| Commit     | Files | Content                                                                                                                  |
| ---------- | ----- | ------------------------------------------------------------------------------------------------------------------------ |
| `443b9905` | 1     | proposal HTML copied from template (1077 lines)                                                                          |
| `9958974c` | 3     | proposal content (327-line rewrite), `docs/adr/0040-component-module-extraction-triggers.md` (83 lines), ROADMAP pointer |
| `ad00d189` | 2     | AGENTS.md decision pointer, CHANGELOG `[Unreleased]` entry                                                               |

Working tree clean at report time. `master` is **ahead 3, unpushed** — the daemon did NOT push
this time (it has pushed without being asked before).

> Format note: status-report skill's canonical output is HTML; the user explicitly requested
> `.md` at `docs/status/<ts>.md`, so the override is honored (consistent with the existing
> `.md` files in this directory).

---

## a) FULLY DONE (verifiable, committed, green)

1. **Phase 1 — current state detected.** All 9 `go.mod` files read (7 published + visualtest +
   website), `go.work` mapped, DAG confirmed acyclic, per-module dependency/replace tables built.
   Evidence: module table in
   `docs/modularization/2026-09-14_component-module-boundaries.html` §01.
2. **Phase 2 — coupling research.** Kanban imports `utils`, `utils/svg`, `utils/wire`
   (`display/kanban.go:6-13`); Heatmap imports `utils` + same-package `chart_geometry`
   (`display/heatmap.go:1-8`, shared nice-max math at `chart_geometry.go:49,61`); display's only
   module deps are utils/icons/htmx; co-change counts pulled (93 display commits since Aug;
   kanban 10, heatmap 6); ADR-0020/0034/0039 read in full; 12+ new-module touchpoints enumerated
   from `release.sh`, `check-module-sync.sh`, CI yaml.
3. **Phases 3–5 — proposal, self-review, plan.** Bauhaus-template HTML proposal with scoring
   matrix (7 existing modules all Keep; 3 candidates all trigger-gated NO), DAG ASCII diagram,
   T1/T2/T3 trigger table, 15-question self-review table, 3-step execution plan.
4. **ADR-0040 written and integrated.** Decision + evidence + triggers; linked from ROADMAP
   ("Per-package modules split" row), AGENTS.md module-structure section, and CHANGELOG
   `[Unreleased]`. No ADR index exists (checked — nothing to update there).
5. **Phase 6 — executed and verified.** `go build ./...` OK (workspace). `TestDocsCountDrift`
   and `TestVersionMatches*` PASS (run AFTER the last AGENTS.md edit). HTML tag-balance check
   clean (only self-closing `<meta/>`/`<br/>` parser false positives).
6. **Phase 7 — reflection + docs.** AGENTS.md pointer added so future sessions don't
   re-litigate; CHANGELOG warmed per the one-commit-release doctrine.
7. **Skills loaded per contract:** go-modularize (SKILL + phases + report-kit), how-to-golang,
   brutal-self-review, status-report (+ section-quality guide).

## b) PARTIALLY DONE (works, with named gaps)

1. **Native chart family as a third candidate** — named in the §03 matrix with depth
   "arguable" and zero dedicated evidence (no per-file import audit of
   line/area/pie/bar/sparkline/heatmap). The word "arguable" is now in a permanent artifact.
   Remaining: file-by-file import + substrate audit. Effort M. No blocker — I stopped because
   the user scoped the session to kanban/heatmap.
2. **F5 co-change evidence is inference, not proof.** Counts come from `git log` whose messages
   are daemon heuristics ("chore: auto-commit N files"); I inferred "repo-wide sweeps" from
   file counts, not content. Remaining: `git log --numstat` pair analysis. Effort M.
3. **Skill conformance deviations (accepted, undocumented in the artifacts):** skill's commit
   steps skipped (harness rule: never commit without explicit user request — the daemon
   committed instead, with heuristic messages); D2 diagram replaced by ASCII dep-tree;
   separate `<date>_EXECUTION_PLAN.html` folded into proposal §07 (skill allows this).
4. **DAG-enforcement claim not independently verified.** Proposal cites
   `scripts/check-module-layers.sh` as the enforcement mechanism; I never read the script.
   Effort S.
5. **how-to-golang loaded shallowly** (main SKILL.md only, no banned-libraries reference pass).
   Zero risk this session (docs-only, no deps changed), but the Phase 4.2 mandate was only
   partially met.

## c) NOT STARTED (planned/needed, zero work done this session)

1. ~~**HARVEST of section (f) into TODO_LIST.md/ROADMAP.md** (status-report skill: section (f)~~ done (DONE - harvested across the 09-14/09-16/09-17 sessions + the 2026-09-17 evening pass)
   ~~must not die in this timestamped file). Blocked on user instruction — user said "wait".~~
2. ~~**Doc drift found and deliberately deferred** (user: "do not research/fix unrelated stuff"):~~ done (DONE 2026-09-17 evening - AGENTS.md now lists display -> icons,utils,htmx; modularization README go-work-use now includes website (proposal-HTML row left open))
   ~~- AGENTS.md import-graph line omits `display → htmx` (card.templ imports `htmx.SwapStyle`~~
     ~~at `display/card.templ:6,289`) — real drift in the Production-deps enumeration.~~
   ~~- `docs/modularization/README.md` Contributing `go work use` list is missing `website`~~
     ~~(actual `go.work` has 9 entries incl. `./website`).~~
   ~~- Proposal HTML not added to that README's "Files" list (miss I created).~~
3. **F058 owner decision** (ADR-0039 "Proposed — owner decision pending"): ADR-0040's
   "timed with v2" clause silently assumes option 1. Not mine to decide.
4. ~~**Post-push CI watch** for the 3 unpushed daemon commits (documented daemon same-day~~ done (SUPERSEDED - pushed and released (v1.18.0, 511d3ed6); daemon drifts caught by the evening pass)
   ~~regression pattern: CSS un-minify, website typescript pin).~~
5. **Remainders from (b):** chart-family audit, content-level co-change, compile-time
   measurement, check-module-layers.sh read.
6. ~~**BuildFlow upstream fixes** (heuristic commit messages; templ-generate .gitignore~~ done (duplicate of TODO_LIST #93/#124 family)
   ~~re-append) — out of repo, documented 5+ sessions, untouched here.~~

## d) TOTALLY FUCKED UP (radical honesty — nothing is on-fire, three things stink)

Nothing shipped is broken: build green, guards green, no code touched, consumers unaffected.
The genuinely fucked-up parts are trust and discoverability:

1. **I shipped soft claims into permanent artifacts.** ADRs/proposals live for years. §03's
   "arguable" chart-family row and F5's daemon-inferred co-change are unevidenced-by-my-own
   standard. A top-tier engineer reading ADR-0040's evidence section would ask "how do you
   know?" for exactly these two rows — and the answer today is "I don't, fully." Severity:
   medium (credibility, not correctness — the verdict survives without them).
2. **Decision record committed with garbage messages.** ADR-0040 is a decision of record and
   its three commits say "chore: auto-commit N changed file(s) (heuristic)". `git log --grep`
   — the repo's own documented search practice — will never find this decision. Severity:
   low-medium; fix needs history rewrite (unpushed, so cheap) or acceptance.
3. **Trigger-list divergence risk (small split brain).** ADR-0040 generalizes ADR-0020's
   triggers (I added T3, dropped the "go.sum >50 entries" wording) without an explicit
   "supersedes-for-component-granularity" sentence in ADR-0020 itself. Two trigger lists now
   coexist; a future reader may not know which governs. Severity: low.
4. **Fix-on-sight tension resolved by silence.** AGENTS.md mandates fixing noticed drift
   immediately; the user's instruction said report-only. I followed the user (correct), but
   item 2/c is only captured here — if this report dies unharvested, the drift stays invisible
   to the next session. That's a process hole, not an error.

## e) WHAT WE SHOULD IMPROVE (process/design, concrete)

1. **Evidence bar for decision records:** no "arguable" rows in ADRs — measure it or omit it.
   One extra hour would have made ADR-0040 bulletproof (chart-family audit + numstat).
2. **Commit authorization for decision records, upfront.** Decision artifacts deserve real
   commit messages; the harness/daemon combo guarantees heuristic messages unless the user
   explicitly authorizes commits. Ask at session start when the deliverable is a decision.
3. **Co-change analysis method:** never trust daemon commit messages; use `--numstat`
   file-pair co-occurrence. Write this into the repo's analysis runbook (AGENTS.md already
   documents daemon behavior — add the analysis caveat).
4. **Recurring-question runbook:** "should X be its own module?" now has a 3-line answer path
   (ADR-0040 trigger check). Add it to the templ-components skill so the next session answers
   in minutes instead of re-running 7 phases.
5. **Fix-on-sight vs scope discipline:** keep collecting noticed-but-deferred drift into a
   named list (this report §c2) so HARVEST can't lose it.
6. **Report format convention:** docs/status/ is all `.md`; the status-report skill wants
   HTML. Align one of them (ask user / patch skill) to stop shipping overrides.
7. **Trigger-list single-source:** make ADR-0020 point to ADR-0040 for component granularity
   (one sentence) so exactly one trigger list governs per granularity level.

## f) NEXT TASKS (up to 50 — brainstorm, HARVEST-routed; Impact/Effort/Category)

**A. Doc drift + consistency (found this session)**

| #  | Task                                                                                                 | Impact | Effort | Category      |
| -- | ---------------------------------------------------------------------------------------------------- | ------ | ------ | ------------- |
| 1  | Fix AGENTS.md Production-deps line: add `display → htmx` (card.templ imports htmx.SwapStyle)         | High   | S      | Documentation |
| 2  | Add `website` to docs/modularization/README.md `go work use` command                                 | High   | S      | Documentation |
| 3  | Add 2026-09-14 proposal HTML to docs/modularization/README.md "Files" list                           | Medium | S      | Documentation |
| 4  | Add ADR-0040 pointer + "component extraction policy" note to docs/modularization/README.md           | Medium | S      | Documentation |
| 5  | One sentence in ADR-0020: component-granularity triggers now live in ADR-0040                        | Medium | S      | Documentation |
| 6  | Re-align ROADMAP "Per-package modules split" table row pipes after text growth                       | Low    | S      | Cleanup       |
| 7  | Link-check all relative links in ADR-0040 + proposal HTML resolve on GitHub                          | Medium | S      | Quality       |
| 8  | Regenerate go.work from the README command and diff vs current to catch other drift                  | Medium | S      | Quality       |
| 9  | Verify `scripts/check-release-tags.sh` covers all 7 sub-module tag sets (incl. datastar)             | Medium | S      | Quality       |
| 10 | Verify AGENTS claim "git status re-added `*_templ.go` line" behavior still matches current BuildFlow | Low    | S      | Documentation |

**B. Evidence hardening for ADR-0040**

| #  | Task                                                                                                               | Impact | Effort | Category      |
| -- | ------------------------------------------------------------------------------------------------------------------ | ------ | ------ | ------------- |
| 11 | Audit native chart family (imports/substrate per file) and replace §03 "arguable" row with evidence or drop it     | High   | M      | Quality       |
| 12 | Replace co-change inference with `git log --numstat` pair analysis for F5                                          | Medium | M      | Quality       |
| 13 | Measure `display` package compile time; cite it in the proposal's package-granularity argument                     | Low    | S      | Quality       |
| 14 | Read `scripts/check-module-layers.sh`; confirm the enforcement claim cited in ADR-0040                             | High   | S      | Quality       |
| 15 | Note in ADR-0040 that `utils/wire` relocation (if ADR-0035 freeze ever lifts) is a mini-trigger to re-check kanban | Low    | S      | Documentation |

**C. Verification debt**

| #  | Task                                                                                                                                               | Impact | Effort | Category |
| -- | -------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | -------- |
| 16 | Run full per-module loop (`for mod in utils icons errorpage charts/echarts datastar htmx; GOWORK=off go test ./...`) before next push              | High   | M      | Quality  |
| 17 | Run `nix run .#lint` (docs-only change, but AGENTS.md is guard-scanned)                                                                            | Medium | S      | Quality  |
| 18 | Open proposal HTML in a real browser once (parser check ≠ render check)                                                                            | Low    | S      | Quality  |
| 19 | `git fetch`; decide push of the 3 unpushed commits; watch CI after push (css byte-stability, website typescript pin per daemon-regression pattern) | High   | S      | Cleanup  |
| 20 | Re-verify `TestDocsCountDrift` on CI after AGENTS/CHANGELOG land (guard scans prose counts)                                                        | Medium | S      | Quality  |
| 21 | Confirm `ad00d189` AGENTS.md sentence renders correctly (view the committed diff, not the worktree)                                                | Low    | S      | Quality  |

**D. Decision + policy**

| #  | Task                                                                                                                                                                | Impact   | Effort | Category      |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------ | ------------- |
| 22 | Owner decides F058 (ADR-0039 option 1) — record confirmation in ADR-0039 Status line                                                                                | Critical | S      | Decision      |
| 23 | When v2 fires: append "re-run ADR-0040 trigger evaluation per component family" to the ADR-0039 mechanics checklist                                                 | Medium   | S      | Documentation |
| 24 | Policy: commit-per-task with explicit messages when user authorizes (re-learned today; daemon raced twice)                                                          | High     | S      | Process       |
| 25 | Optional: rewrite the 3 unpushed heuristic commits into properly-messaged docs commits via detached temp worktree (AGENTS documented pattern) — needs user approval | Medium   | M      | Cleanup       |
| 26 | Document in AGENTS Git Workflow: skills' commit steps defer to the Crush no-commit rule; daemon owns commits otherwise                                              | Medium   | S      | Process       |
| 27 | Encode the "module-extraction question → run ADR-0040 trigger check first" shortcut in the templ-components skill                                                   | Medium   | S      | Process       |

**E. BuildFlow / daemon (upstream)**

| #  | Task                                                                                          | Impact | Effort | Category |
| -- | --------------------------------------------------------------------------------------------- | ------ | ------ | -------- |
| 28 | Implement diff-based daemon commit messages in larsartmann/buildflow (documented 5+ sessions) | High   | L      | Bug      |
| 29 | Fix BuildFlow templ-generate step re-appending `*_templ.go` to .gitignore each run            | Medium | M      | Bug      |

**F. Component-adjacent smells noticed (flagged only, not researched)**

| #  | Task                                                                                                                                                        | Impact | Effort | Category      |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | ------------- |
| 30 | Clarify `display/collapsible_heatmap_test.go` / `collapsible_persist_test.go` naming (heatmap tests named after collapsible_section)                        | Low    | S      | Cleanup       |
| 31 | Verify Heatmap has golden-sweep + visual goldens (kanban's coverage was confirmed this session; heatmap's was not)                                          | Medium | S      | Quality       |
| 32 | Document why Heatmap needs no e2e (pure CSS grid, unwired) — kanban doctrine says wired components get e2e; make the unwired case explicit                  | Low    | S      | Documentation |
| 33 | Decide/document `display/card.templ`'s typed `HxSwap htmx.SwapStyle` cross-module field (accept + document, or re-type)                                     | Low    | S      | Documentation |
| 34 | Check whether `website/go.mod`'s root-module require is actually used (no `display` imports found in website this session); slim to used sub-modules if not | Medium | S      | Cleanup       |
| 35 | Confirm visualtest's root-module require is used (kanban/wire e2e import display — believed yes; make it explicit)                                          | Low    | S      | Documentation |

**G. Meta / skills / tooling**

| #  | Task                                                                                                                           | Impact | Effort | Category |
| -- | ------------------------------------------------------------------------------------------------------------------------------ | ------ | ------ | -------- |
| 36 | Patch go-modularize skill: commit steps conditional on harness commit policy                                                   | Medium | S      | Process  |
| 37 | Align status-report skill output format with docs/status/ `.md` convention (or migrate dir to HTML)                            | Low    | S      | Process  |
| 38 | Next Go-dependency-touching session: run the how-to-golang banned-libraries reference pass fully                               | Medium | S      | Process  |
| 39 | Add "noticed-drift ledger" habit: deferred fixes go into the session's status report §c (as done here), never silently dropped | Medium | S      | Process  |
| 40 | Evaluate (likely reject) a `TestAGENTSImportGraphMatchesImports` guard — probably too brittle; document the rejection if so    | Low    | M      | Quality  |

**H. HARVEST / follow-through**

| #  | Task                                                                                                                                    | Impact | Effort | Category      |
| -- | --------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | ------------- |
| 41 | Run docs-health HARVEST on this (f) list → TODO_LIST.md rows / ROADMAP fuel                                                             | High   | S      | Process       |
| 42 | Route items 1–5 + 22 into TODO_LIST as near-term rows; 11–14 as quality rows                                                            | High   | S      | Process       |
| 43 | Decide fate of items 6, 10, 18, 21 (likely batch into one cleanup commit)                                                               | Low    | S      | Cleanup       |
| 44 | After F058 resolution, revisit item 23 and ADR-0040's v2 wording in one pass                                                            | Medium | S      | Documentation |
| 45 | Verify GitHub renders ADR-0040's relative link to the HTML blob (manual click)                                                          | Low    | S      | Quality       |
| 46 | Sweep other prose for stale "7 modules" phrasing that should cite ADR-0040 (README top-level)                                           | Low    | S      | Documentation |
| 47 | Check whether any TODO_LIST item already covers items 28–29 (BuildFlow) to avoid duplicate rows during HARVEST                          | Low    | S      | Process       |
| 48 | Add ADR-0040 to the next release's CHANGELOG highlight list (already in `[Unreleased]`; verify it survives release.sh's cut)            | Medium | S      | Documentation |
| 49 | Give `display/collapsible_persist_test.go`'s "persist" naming a look during item 30 (possible ghost naming from a removed API)          | Low    | S      | Cleanup       |
| 50 | Close the loop: after HARVEST, mark this report as harvested (docs-health ANNOTATE convention) so the next session trusts (f) is routed | Medium | S      | Process       |

---

## g) QUESTIONS (3 — none answerable by me from the repo)

1. **F058 / ADR-0039:** Do you confirm **option 1** — the module path migrates to `/v2` exactly
   at the first real breaking change? ADR-0040's "timed with v2" clause and every T1–T3
   extraction's mechanics depend on this. I read ADR-0039 (Status: "Proposed — owner decision
   pending") and cannot decide it.
2. **Consumer signals:** Do you know of private/off-repo interest (e.g., from the #156 survey
   channels, cqrs-htmx, or direct messages) in **kanban-only, heatmap-only, or charts-only
   adoption**? Repo docs say 0 confirmed consumers for family-only imports; if private signals
   exist, T2 fires and today's NO flips for that family.
3. **History policy:** Should I leave the three unpushed heuristic daemon commits
   (`443b9905`, `9958974c`, `ad00d189`) as-is, or — with your explicit commit authorization —
   rewrite them into properly-messaged commits (temp-worktree pattern, then
   `--force-with-lease`) and commit per-task going forward for decision records?

**WAITING FOR INSTRUCTIONS.**
