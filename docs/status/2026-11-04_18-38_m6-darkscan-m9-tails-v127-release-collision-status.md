# Status — M5 Closure, M6 Dark-Scan, M9 Demo/Docs Tails, v1.27.0 Release Collision

**Written:** 2026-11-04 18:38 CET
**Session scope:** continuation of the CORPUS-CLOSE-OUT plan — M5 closure, M6 (#490), M9 (#497)
**Repo state at write time:** local tip `d6ba51c0`, **clean tree**, origin/master at `64f44638`, **7 commits unpushed**, tag `v1.27.0` exists LOCALLY (sibling cut it 18:36), no ritual yet at `d6ba51c0`.

## TL;DR

M5 is closed and pushed with two witnessed M03 rituals. M6 (dark-scan ring-offset/shadow color rules) is fully landed and pushed — zero library violations found, the scan extension was pure prevention. M9 content is 100% committed locally (`9d2675c7`): both undemonstrated features now have wired demo sections, the docs-site code chrome shipped with the md package's first tests, and both open decisions (forms_demo rename, ld+json) are recorded as NO. The session's defining event: **the v1.27.0 release collided with my work in the shared worktree** — the daemon ate the sibling's release staging (`9db26758`), the sibling aborted and re-cut (`128d76cc` revert → `97c49407` release → `d6ba51c0` replace-restore → tag v1.27.0). My commit survived beneath the release, the tree is clean again, but **the 7 unpushed commits + the local-only tag now need one coordinated ritual-and-push** that I have deliberately NOT started: it delivers the release to origin, and tag-push authority is the release session's per the checklist.

## a) FULLY DONE

1. **M5 closure end-to-end.** Diagnosed the "undiagnosed utils FAIL" from the morning: it was `TestDocsCountDrift` — the migration session had written a raw-git count (147 = 138 library + 9 website) into AGENTS.md where the guard counts library-only (fixed in `e1be762c`, verified with `-count=1`). Then two more `wsl_v5` findings in my own M5 test files (`display/flag_combinations_test.go`, `utils/utils_test.go`) surfaced lane-by-lane — both fixed (`6dcb93a9`, `279e5081`). Full ritual **VERDICT: PASS** at `279e5081` (17:36:14, log `~/var/log/tc-rituals/2026-11-04T1729-tip279e5081.log`), pushed immediately, tips matched.
2. **TODO #493 struck with line-numbered evidence** (all 8 items, the ScriptComponent nonce-fallback bug, the lint tails) plus #496 annotated (8/9 done, only the status-index decision open) and #497's Alert/LiveRegion sub-item struck — commit `750d5a28`, second witnessed ritual (PASS 17:44:01), pushed.
3. **Generated-file count drift fixed everywhere it lived**: AGENTS.md (guard-guarded claim → 138, narrative → "147 incl. website's 9") and `skill/SKILL.md:273` (unguarded prose the guard regex doesn't match — same drift, fixed in the same docs commit).
4. **M6 (#490) dark-scan extension, complete and pushed** (`8c0dee93` + CSS `64f44638`): color-carrying `ring-offset-*`/`shadow-*` forms joined the neutral+semantic regexes (hues derived from the tone table); size classes stay out of scope by design so the Toggle thumb needed **no** allowlist entry; `TestDarkModeRingShadowColorRules` (flagged/never-flagged tables + the real thumb line through the exception path); pre-flight scan found **zero** library gaps; AGENTS bullet + CHANGELOG entry; the freshness-guard `templates/templ-components-classes.txt` regenerated in the same commit; demo app.css recompiled when Tailwind's source scan picked up the new class-shaped tokens (CSS freshness lane caught it — the gate worked as designed). Ritual **VERDICT: PASS** at `64f44638` (18:01:20), pushed, tips matched.
5. **Post-migration follow-through probe**: `scripts/check-templ-gate.sh` → `GATE: CLOSED — pin and generator are in lockstep` (v0.3.1100 == v0.3.1100).
6. **M9 demo additions, committed** (`9d2675c7`): `display_demo.templ` gains a "Card Header Slot" section (`demoCardHeaderSlot()` named slot component + default-header comparison + code snippet) and a Grid auto-fit block (`GridColsAutoFit` + `MinColWidth: "190px"`); canonical rm-all + per-module generate sequence; workspace build green; **demo smoke green** (served fresh binary on :18081, all four content markers found on `/display` incl. `minmax(190px`); display route goldens re-recorded and verified.
7. **M9 docs-site code chrome, committed**: `md.go` wrapper renderer emits a chrome bar (fence title via goldmark brace attribute `{title="..."}`, language label, copy button moved into the bar); `codeBlockTitle` with nil-attribute guard and `string`/`[]byte` value handling; `site.css` rebuilt around the chrome (rounded join, mono labels, always-visible copy); `installation.md` ships the first titled fence (`​```bash {.title="terminal"}`); **`md_test.go` — the md package's FIRST tests** (lang label, HTML-escaped title, no-lang fence); pages HTML golden updated via `-update`; website module lint 0 issues + all 4 packages green.
8. **Both M9 decisions made and recorded in the TODO row**: `forms_demo.templ` rename → **NO** (sibling-convention file naming; the astro-conversion premise is obsolete); demo ld+json → **NO** (component gallery, not an SEO surface).
9. **Session-end repo state verified**: tree clean at `d6ba51c0`; the release collision resolved itself (see d/f); tag `v1.27.0` local-only.

## b) PARTIALLY DONE

1. **M9 closure**: content 100% committed, but the commit is 1 of **7 unpushed**; no ritual has run at `d6ba51c0`; TODO row #497 has every item struck/decided but the **row header is not yet fully struck** (`| 497 |` should become `| ~~497~~ |` like #493).
2. **The v1.27.0 release** (sibling's work, observed not owned): release commit `97c49407` + replace-restore `d6ba51c0` + tag exist locally; **nothing is pushed to origin**; per the release checklist the operator reviews `git show v1.27.0` then pushes master + tags manually — that step is pending and entangled with my 7 unpushed commits.
3. **Push authority question from the morning report**: partially answered empirically — the sibling pushed mid-session through `6dcb93a9`, I pushed `750d5a28`/`64f44638` — but the CURRENT 7-commit backlog (which contains their release commit) has no owner ruling yet.
4. **My CHANGELOG entries vs the release**: the dark-scan (Added) and code-chrome (Changed) entries were written while `[Unreleased]` was warm; the release commit consumed `[Unreleased]` into 1.27.0 minutes later. I have NOT yet verified my two entries survived the consumption verbatim (highly likely — release.sh moves the section — but unverified).
5. **Earlier owner questions** (§g of the 16:44 report): push ownership (partially resolved), daemon-eaten M5 commit message recovery (still open, though the #493 strike now documents the content), 626-table depth + `/mnt/buildcache` split (unanswered).

## c) NOT STARTED

1. **M11** (#492): Table `Flush+Bordered` interaction golden, pagination `current>total` case, demand-gated feature set (F109 pattern).
2. **M12** (#494): NotFound404 demand-gated polish bundle.
3. **M7** (#489): keyboard/a11y bundle (MobileMenu focus-trap, type-ahead, Combobox PageUp/PageDown, Tabs parity, Carousel aria-live, arrow-key e2e).
4. **M8** (#495): visual tooling (#tc-root padding config, 50× variance lane, overlay/blink diff generator).
5. **M10** (#488): release/infra hardening (release.sh --dry-run, step-0 preflight, CI shellcheck job, worktree lifecycle script, test-all-modules, CSS provenance stamp) — made MORE urgent by today's release collision.
6. **M13**: 2026-10/11 corpus sweep — depth still gated on the 626-wholly-unstruck-table ruling.
7. **M17**: ROADMAP residue go/no-go.
8. **#368 (size-*) / #369 (@utility)**: unblocked standalone templ events, each needing its own golden re-baseline window.
9. **erraudit sweep lane** and **coverage-gate ratchet check** — plan residue, unrun this session.

## d) TOTALLY FUCKED UP (what I forgot, broke, or mishandled)

1. **I walked into the known root-form `templ generate` trap AGAIN.** My first generate failed on a syntax error in my own comment (missing `//` on continuation lines) AND — because the root form ran with go.work active — rewrote FileNames workspace-relative across the submodules' generated files. I recovered with the canonical rm-all + per-module sequence (verified byte-identical afterwards, only my two files dirty), but this is a documented AGENTS.md gotcha I executed into anyway. The comment syntax error itself was careless: I wrote a three-line comment and only commented line one.
2. **I raced a sibling session's visual-suite run.** My scoped `-update` captured `display_*` goldens while the sibling was capturing `datastar_*`/`htmx_*`/`site-sales` in the same worktree; the deltas mixed, the daemon committed all of it (`03cc775f`), and my mid-flight `git restore` of the datastar goldens raced that commit (it lost — the capture versions landed). I had checked for active chromium processes but only AFTER the mismatch appeared, not before launching the run. Outcome verified green afterwards (scoped no-update run over display/datastar/htmx, 92.8s ok), but the sequence was luck-adjacent.
3. **My unpushed-commit backlog entangled the release cut.** The daemon ate the sibling's release staging as `9db26758` ("chore: auto-commit 15") — violating the one-commit-release convention (canonical message at tip) — and my `9d2675c7` landed beneath the re-cut release, so the release commit is no longer the tip. The sibling handled it (revert + re-cut + replace-restore), and `scripts/check-release-tags.sh` plus the checklist exist precisely for this, but my slower commit cadence (CSS-lane failure → recompile → second ritual round-trips) widened the collision window.
4. **The `wsl_v5` class is MY signature now.** Five findings across the session (2 in M5 files, 1 in the M6 test, 3 in md.go): I consistently under-whitespace new Go code, and because the daemon bypasses hooks when it commits under contention, my style slips reach master before any lint sees them. Every one was caught downstream — by the ritual, never by me.
5. **`rg -rn` — third strike**, the hook literally counts them now. Also a `multiedit` on AGENTS.md that failed on mtime staleness, and a display_demo multiedit that used a delete-then-reinsert pattern (momentarily deleting the `@demoSection("Cards", ...)` line before re-adding it). Sloppy edit hygiene under time pressure.
6. **Two wrong API assumptions in the chrome implementation**: I assumed goldmark brace attributes use `.title=` (`.x` means CLASS in goldmark; bare `title=` is the attribute) and that quoted values are `string` (they are `[]byte`). Both were caught by my own new md test — that's the test earning its keep — but each cost a build cycle. A third bug (`ctx.Attributes()` returning nil → panic in TestSiteBuildIntegrity) was caught by the existing suite.
7. **A test-design bug in my own M6 test**: I asserted the Toggle-thumb line with a synthetic path, tripping the pre-existing `bg-white` neutral rule the exception table never sees. Fixed by testing the size-only line synthetic AND the real line through the real path — should have thought of that first.
8. **Minor unfinished residue**: #497 row header not fully struck (see b1); SKILL.md's dark-mode section doesn't mention the new ring-offset/shadow coverage (AGENTS.md does); I restored the raced datastar goldens but not the equally-raced htmx ones (inconsistent reasoning in the moment — the subsequent verification made it moot).

## e) WHAT WE SHOULD IMPROVE

1. **Authoring-time lint, always.** Run targeted `golangci-lint run` on every newly written/edited Go file BEFORE staging. The wsl_v5 class disappears at zero cost; downstream ritual failures cost 8–20 minutes each.
2. **Canonical generate sequence only.** Under go.work, `rm-all → root generate → per-module generate` is the ONLY acceptable invocation. Bare root-form generate is banned for me from now on, regardless of how small the change seems.
3. **Sibling-activity check BEFORE any visual-suite or heavy nix run**: fresh `git fetch` + `git log origin/master..HEAD` + a chromium/go-test process peek + `git status` freshness. Today's collision was predictable from `5ff346ae` existing.
4. **Commit fast, push faster.** The M5 lesson (daemon eats staged work >5 min) repeated at release scale (daemon ate the release staging). Finished work should not sit staged; finished commits should not sit unpushed across a ritual failure — fix forward and re-ritual quickly.
5. **Goldem-capture hygiene for shared worktrees**: route-golden runs could fail-fast if `testdata/routes/**` changes under them (capture-hygiene guard). Proposal-worthy for M8 (#495).
6. **Pre-verify API surface before writing code against it**: one `rg` into the pinned module source (as I finally did for `CodeBlockContext` in the RIGHT version — the first read was against a stale v0.0.0 module while the build used v2.0.0-2023) would have saved two build cycles.
7. **The CSS-freshness lane taught a real lesson**: Go source strings containing class-shaped tokens (scanner regexes, comments) feed Tailwind's scan. Any future Go-side rule table that names classes must expect an app.css recompile in the same commit. Worth one AGENTS.md line.
8. **Keep using smoke-first verification**: the :18081 demo smoke (python urllib because curl is banned) caught content presence before goldens — cheap and effective; a small `scripts/demo-smoke.sh` would standardize it.

## f) 50 NEXT THINGS (priority-ordered, grouped)

**Immediate — unblock the push (1–6)**
1. Full M03 ritual at tip `d6ba51c0` (tree is clean NOW; durable log; `--lint --css --website`).
2. On VERDICT PASS: push the 7 commits (`5ff346ae`…`d6ba51c0`) to origin/master, re-fetch, verify tips match.
3. Coordinate the `v1.27.0` TAG push with the owner/release session (checklist: review `git show v1.27.0` first) — question g1.
4. Fully strike TODO row #497 (header strikethrough; all items resolved) — same docs commit as 24.
5. Verify my two CHANGELOG entries survived the release consumption verbatim in the 1.27.0 section; re-warm `[Unreleased]` with the next feature entry.
6. Run `scripts/check-release-tags.sh` post-push: root + every published submodule tagged in lockstep.

**Release follow-through (7–11)**
7. Post-push CI watch: Build & Test + Website workflows green at the pushed tip (the docs-layout HTML golden changed — website CI must re-render clean).
8. `check-templ-gate.sh` post-push probe (CLOSED expected) + `TestTemplVersionPin` green in CI.
9. Verify the v1.27.0 release notes match what actually shipped (my M6/M9 landed BEFORE the release cut — are they in the 1.27.0 section? They should be).
10. pkg.go.dev / module-proxy spot check after the tag lands publicly (consumer pairing: templ-components ≥ v1.27.0 ⇄ templ ≥ v0.3.1100).
11. Coverage ratchet re-check post-release (`coverage-library.sh --gate 94`) — M5/M6/M9 added tests; confirm no drift against `scripts/coverage-baseline.txt` prose.

**TODO residue from my own strikes (12–15)**
12. #496 residue: `docs/status/README.md` historical-index decision (accept archived/ layout vs new index) — the last open item in that row; smallest deliverable available.
13. SKILL.md dark-mode section: mention ring-offset/shadow coverage parity with AGENTS.md.
14. Add the SKILL.md-style "147 generated files" phrasing to `TestDocsCountDrift`'s guarded regexes so unguardable prose drift like today's gets caught (guard-gap I found and patched manually).
15. One AGENTS.md line: Go-source class-shaped tokens (scanner rule tables) feed the Tailwind scan → app.css recompile expected in the same commit.

**Plan milestones (16–23)**
16. M11 #492: Table Flush+Bordered interaction golden + pagination `current>total` case.
17. M12 #494: NotFound404 demand-gated items (each individually gated).
18. M7 #489: keyboard/a11y bundle (MobileMenu focus-trap, type-ahead, PageUp/PageDown, Tabs parity, Carousel aria-live, arrow-key e2e).
19. M8 #495: #tc-root render-shell padding config, 50× same-Chromium variance lane, overlay/blink diff generator (fold in the capture-hygiene guard idea, e5).
20. M10 #488: release.sh --dry-run + step-0 preflight (go.work clean, signing key, **daemon-quiesce**) — directly motivated by today's `9db26758`.
21. M10 continued: CI shellcheck job, worktree lifecycle script, `scripts/test-all-modules.sh` with sanitized log paths, compiled-CSS provenance stamp.
22. M13 corpus sweep 2026-10/11 (scope per the 626-table ruling, g3).
23. M17 ROADMAP residue go/no-go sweep.

**Unblocked templ events (24–25)**
24. #368 size-* utilities migration — standalone event with its own golden re-baseline (unblocked by the completed #335 window).
25. #369 @utility directive migration — same shape, separate event.

**Dark-scan follow-through (26–29)**
26. End-to-end `tc lint --dark` CLI run over the demo tree to prove the consumer path with the new rules (I proved the Go-test path only).
27. Scope decision: `ring-black/5` hairlines in popover/dropdown have no dark: twins (same-in-both-mode by design?) — either allowlist or extend the rule family; needs a deliberate ruling.
28. Consider `border-*`/`divide-*` semantic-hue coverage parity check (currently only gray in neutral; hues only via bg/text) — audit whether the asymmetry is principled.
29. Add `ring-offset-*`/`shadow-*` color rules to the `tc` CLI docs/help text if the rule table is user-documented anywhere.

**Demo/docs polish from today's work (30–35)**
30. xl-viewport probe of the two new display sections (route goldens render below xl:; the auto-fit block could overflow at 1280 — throwaway custom-viewport check per AGENTS.md).
31. Title the remaining getting-started fences (`{.title=...}`) where a filename/name is natural (content polish).
32. Chrome bar print/forced-colors audit on docs pages (site print CSS hides shadows; verify the chrome reads correctly when printed).
33. Verify the copy-button JS binding still fires from the chrome position (docs.js binds `.code-copy`; DOM position moved) — one CDP probe post-deploy.
34. `Card.HeaderClass` interplay with the Header slot — micro-test candidate (my demo exercises the slot, not the class override).
35. `scripts/demo-smoke.sh` helper (standardize the python-urllib smoke I hand-rolled twice).

**Tooling/guard proposals born this session (36–40)**
36. Daemon quiesce flag for release.sh runs (lockfile the script sets; daemon skips while present) — needs owner ruling (g2).
37. Capture-hygiene fail-fast in route-golden tests when testdata/routes changes mid-run (folded into M8, listed separately for visibility).
38. Extend the drift-guard family: assert SKILL.md/README count phrasings match the walk (per 14) — prevents the exact AGENTS/SKILL split-brain fixed today.
39. Pre-push ritual helper: a wrapper that refuses to start ci-repro while sibling git activity is <N minutes old (heuristic guard against visual/nix collisions).
40. Consider pinning ONE docs site route as a pixel golden (docs pages are currently HTML-golden-only; the chrome changed pixels with no pixel net) — new ask, ROADMAP row.

**Debt/infra (41–45)**
41. `/mnt/buildcache` split (shared Forgejo/dev caches: SQLITE_BUSY + gomod-freshness storms all session) — still needs the owner ruling (g3).
42. 626 wholly-unstruck-table archived-debt depth ruling (gates M13).
43. erraudit sweep lane (`scripts/erraudit-sweep.sh`) — plan residue, unrun.
44. Fleet CSS re-check (#436 consumers) — untouched this session, standing audit.
45. AGENTS.md daemon section: add the observed release-staging-eating incident class (9db26758) to the incident ledger in `docs/agent-context-history.md`.

**Housekeeping (46–50)**
46. Log today's ritual failures + fixes in the status-report index convention (this file is it).
47. Re-check AGENTS.md line budget (211/220; M6 bullet added zero lines — verify post-release docs churn didn't cross 220).
48. `git status --ignored` sweep for daemon-resurrected `.gitignore`/`result*` strays after the release churn.
49. Update my session protocol notes (AGENTS.md "Agent single-writer protocol" is holding; my additions are personal discipline items e1–e4, not repo edits).
50. Next session's literal first command: `nix develop -c scripts/ci-repro.sh --lint --css --website > ~/var/log/tc-rituals/<ts>-tip-d6ba51c0.log 2>&1` (assuming no new sibling commits — re-fetch first).

## g) THREE QUESTIONS I CANNOT ANSWER MYSELF

1. **Push sequencing for the 7 local commits + the local `v1.27.0` tag.** The tip now contains the sibling's release commit (`97c49407`) and replace-restore (`d6ba51c0`). Do you want me to run the ritual and push master AND the tag together (delivering the release to origin as one coordinated step), or does the release session own the tag push per its checklist (in which case I push master only after its go-ahead, or hold entirely)? I will not push until you rule.
2. **Should I build the daemon-quiesce mechanism for release runs** (f36: a lockfile `release.sh` sets and the auto-commit daemon honors — plus the same for ci-repro/visual lanes), or is daemon policy owned in the crush-config repo and out of this repo's scope? Today it ate BOTH my M5 commit and the release staging; the fix belongs somewhere, I just don't know where.
3. **The two still-open rulings from the 16:44 report**: (a) 626 wholly-unstruck tables in archived files — sweep them at what depth (annotate-only vs full re-triage)? (b) `/mnt/buildcache` split — do you want the Forgejo cache separated from dev caches, and is that mine to do? Both gate queued milestones (M13 scope; infra time).

---

### Evidence appendix

- Ritual logs (durable): `~/var/log/tc-rituals/2026-11-04T{1717-tip8790530d,1722-tipe1be762c,1729-tip6dcb93a9,1734-tip279e5081,1740-tip750d5a28,1754-tip8c0dee93,1800-tip64f44638}.log` — PASS at 279e5081 (17:36:14), 750d5a28 (17:44:01), 64f44638 (18:01:20); FAILs: count-drift (1717), wsl display (1722), wsl utils (1729), CSS freshness (1754).
- Pushes this session: `279e5081` and `750d5a28`, `64f44638` — all with immediate tip re-check (matched). Origin/master: `64f44638`.
- My commits this session: `e1be762c` (count drift), `6dcb93a9` (wsl display), `279e5081` (wsl utils), `750d5a28` (TODO strikes + SKILL), `8c0dee93` (M6 dark-scan), `9d2675c7` (M9 demo + chrome + decisions).
- Sibling/daemon commits interleaved: `5ff346ae` (sibling site-sales golden), `03cc775f` + `9db26758` (daemon), `128d76cc` (sibling release-abort revert), `97c49407` (**release: 1.27.0**), `d6ba51c0` (replace-restore), tag `v1.27.0` (local).
- My files touched this session: AGENTS.md, skill/SKILL.md, TODO_LIST.md, CHANGELOG.md, utils/darkmode_scan.go, utils/darkmode_compliance_test.go, templates/templ-components-classes.txt (freshness guard), examples/demo/display_demo.templ (+generated), visualtest/testdata/routes/display_{light,dark}.png, website/internal/md/md.go, website/internal/md/md_test.go (new), website/internal/pages/golden_test.go, website/internal/pages/testdata/docs-layout.golden, website/site.css, website/content/docs/getting-started/installation.md.
