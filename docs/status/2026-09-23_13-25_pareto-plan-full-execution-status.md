# Status — Dedup Follow-Through Pareto Plan: Full Execution (66 Tasks)

**Date:** 2026-09-23, 13:25 CEST
**Session:** Continuation of the 2026-09-23_01:50 DEDUP-FOLLOW-THROUGH-PARETO-PLAN execution.
**Final state:** `VERDICT: PASS (exit 0) [dupl: art-dupl version dev]` at the tip + full visual suite green (137.8s, `-parallel 4`). Tree clean. **24 commits ahead of origin/master — NOTHING PUSHED (house rule: push is the human's).**

---

## 0. Executive Summary

All four tiers of the plan were executed. Every task is either DONE, or explicitly open with the blocker/owner-gate recorded in TODO_LIST.md. Beyond the plan, four unplanned-but-necessary pieces of work shipped: a flaky e2e fix, a guard-bug fix, a templ-fmt repo canonicalization, and a demo rate-limit + session-CSRF hardening pair. Two verification failures at the end were caused by my own session's code (lint findings, rate limiter 429ing the audits) — both found by the repo's own gates and fixed, which is the system working as designed.

The one process wart: the auto-commit daemon raced me on every squash attempt; the original 8-commit dedup blob got pushed by the daemon before any squash could land, so the "clean history" goal for that blob is moot (rewriting pushed history was off-limits).

---

## a) FULLY DONE

### Tier 1 — trust spine (10/10)
- **T01** ogshot smoke: website dist built, real capture → `website/public/og/sales.png` (90KB), exit 0.
- **T02** shots smoke: index light+dark full-page captures (2.8MB each).
- **T03** siteshots smoke: 25 captures light/dark × desktop/mobile + search smoke PASS (8 hits).
- **T04** Full `nix run .#visual` witnessed green twice (137s final). Root-caused `TestSiteSalesCopyButton` flake: bare `chromedp.Text` read raced the clipboard-write `.then` label-swap → replaced with a falsy-until-swapped `pollText` (repo's own poll-helper rule). Learned/recorded: **load-sensitive class on this 32-core shared box** — signature is light AND dark site goldens failing ~100% while demo-route goldens pass; remedy `-parallel 4` (recorded in AGENTS.md).
- **T05** Advisory dup-gate lane in `scripts/ci-repro.sh`: resolves `$ART_DUPL_BIN` → `/tmp/art-dupl-fork` → PATH, echoes `art-dupl --version` beside the verdict, never fails the run.
- **T05b** Binary-supply decision recorded in ADR-0009 (documented-PATH; flake pin deliberately deferred until fork detector stabilizes / #292/#293).
- **T06** Baseline determinism: two fresh recordings byte-identical modulo `recordedAt`, both equal to the committed baseline (and to HEAD's entries).
- **T07** Squash: landed `146fbf94` "feat: dedup follow-through — typed heading enum, drift guard, ADR hygiene" (33 files) as ONE properly-messaged commit. (The original 8-commit daemon blob was already pushed by the daemon — force-push forbidden, so that history stands.)
- **T08** Green at tip: build + display/forms/utils + gate all OK.
- **T09** Ritual rehearsal: full `--lint --website` **VERDICT: PASS** multiple times, finally at the sealing tip.

### Tier 2 — lock in the gains (11/11)
- **T10** CHANGELOG exemption rule written into `docs/release-checklist.md` (internal-only changes exempt; user-visible MUST entry in same commit); `[Unreleased]` warmed and corrected (stale "39 groups" fixed → 122, fork-pinned).
- **T11** `TestDocsCountDrift` extended with `countBaselineEntries` — pins CHANGELOG/ADR/AGENTS "N accepted groups" claims to the parsed baseline. **Proved itself on first run** by catching the CHANGELOG 39-vs-122 drift.
- **T13a** `HeadingTagType` enum (h1–h6 constants + `HeadingTagTypeIsValid` + tests) on Card/EmptyState/CollapsibleSection `TitleTag`; defined type so string literals keep compiling (v1-safe).
- **T13b** Redundant `isValidHeadingTag` normalization removed; goldens untouched (behavior-preserving proven); enum counts bumped 62→63 across FEATURES/README/website; FEATURES enum-table row added (both exhaustive guards green).
- **T14** Calendar month-nav aria-labels → `calendarPrevMonthLabel`/`calendarNextMonthLabel` constants.
- **T15** skill/SKILL.md "Shared internal helpers — reuse these" table (headingTag, chartSeriesGroup, sparkline projection, calendarMonthNav*, browser/distserver).
- **T17** Plan-authoring checklist step 5: run the clone gate in dedup-adjacent plans (#291 fully closed).
- **T18** ADR-0009 "Accepted twin pairs" section: datastar↔htmx normalize twins + ListNote↔EndOfList (#300).
- **T19** Gate re-run green at t=1, t=2, AND t=3; ledger re-read; historical line numbers disclaimed (#306).
- **T20** Historical-counts table marked detector-generation-specific.
- **T21** Advisory-invocation verdict recorded: `-t 1`, default min-lines (raising it hides the 2-3-line clones that hid the evening-pass extractions) (#294 partial).

### Tier 3 — hardening (35/35 planned tasks)
**Mirror/guard cluster:**
- **T22** `TestMirroredPackagesListsMatch` — pins bash `MIRRORED_PKGS` == Go `mirroredPackages` (#283a).
- **T23** `TestTypesFilesHaveTemplTwin` — ghost types-file guard (#283b).
- **T24** `TestAddSmoke` — execs the real `tc` binary for eyebrow + auth_layout in `t.TempDir()` (#283c).
- **T25/T26** `scripts/test-tc-sources-guard.sh` — **16 scenarios** (clean/drift/new/orphan/rename/missing-pkg/idempotence ×2) over detached worktrees. **It caught a real bug:** `--fix` treated a missing package directory as removable orphans (destructive); the guard now blocks with `ORPHAN-SKIPPED` instead.
- **T27** CI Lint job runs the self-test (#308).
- **T28/T29** `TestPackageDepsCoverPackageFiles` — `--list-deps` honesty; added `heading_tag.go`, `kanban.go`, `chart_geometry.go`, `area_chart.go`, `line_chart.go`, `pie_chart.go` (display), `calendar_nav.go` (forms), `embed.go` (layout) that silently predated it (#285).

**Docs-truth + guard-UX cluster:**
- **T30** `tc new` → shipped-commands rewording (2 live refs) (#288a).
- **T31** Verified docs carry NO numeric tc-add coverage claims — nothing to fix (#288b, closed as verified-no-claims).
- **T32** CHANGELOG `### Fixed` entry for the 22-component scaffolder rescue (#288c).
- **T33** `TC_SKIP_SYNC=1` opt-out with a loud banner (#289a).
- **T34** Guard runtime printed beside every verdict (86ms measured) (#289b).
- **T35** Starter CSS re-synced from templates/ (it predated the `.go`-scanning `@source` lesson + several custom.css sections) and pinned byte-identical by `TestStarterCSSMatchesTemplates`; `tc-snap-*` kept as documented consumer utilities (#290).
- **T36** `tc ls` derived "N components addable" footer (#294a).
- **T37** Guard-failure output review — all paths print actionable pointers (#294b).

**Tooling UX + site-tier cluster:**
- **T38** `visualtest/tools/README.md` pointing at the shared `internal/` packages (#304).
- **T39** `-selftest` flag on all three capture tools, verified live (#305).
- **T40** Twin cross-ref comments for `calendarNavQuery` (demo ↔ e2e) (#302).
- **T41** `site_routes_test.go` now uses `internal/distserver.Handler` (was a verbatim copy); packages MOVED from `tools/internal/` → `internal/` so the whole module can import them (#273 fully closed).
- **T42** 575-vs-548 forensics closed: the `_sources` exclude is NOT the whole delta (today a no-exclude scan would be ~657); paste-era generation non-reproducible (#310).
- **T43** varnamelen verdict: keep linter + forced renames; one-time churn (#307).
- **T44** Demo smoke coverage confirmed: display_demo renders EmptyState ×2 + CollapsibleSection ×2 + Card ×16 (#311).
- **T45** `//art-dupl:accept` directives REJECTED vs hash baseline — reasoning recorded in ADR-0009 (#314).

**Sweep/hardening cluster:**
- **T46** `role="combobox"`: both sites on `<input type="text">` per ARIA APG — compliant (#274).
- **T47** Zero `WriteString(literal + literal)` in library sources (#277).
- **T48** Demo CSS freshness verified (`nix run .#css` + `TestCSSFreshness`) (#278).
- **T49** `TestPrerenderMatchesLiveServer` made retry-tolerant (3 attempts, #250).
- **T50** Shared `kanbanParseMoveWithCSRF` prelude extracted (was duplicated 3×) (#262).
- **T51** ogshot is SITE_SKIP_STARS-independent (self-contained HTML); documented (#281).
- **T52** build.sh star-resolution inverted: default = deterministic fallback badge; `SITE_LIVE_STARS=1` for production (CI Website job sets it) (#271).
- **T53** `docs/visual-testing.md` "The site route tier" section (theme pin, skip-stars, scroll-reveal load caveat) (#275).
- **T54** SKILL.md website-SSG knobs section (topLevelPages, lastmod sources, search index) (#276).
- **T55** Warm-dark `dark:`-pair override pattern documented in the Tailwind adoption guide (#279).
- **T56** /tmp gate logs purged before salvage → declared accepted-loss in ADR-0009 (#313).

### Tier 4 — features + gated
- **T59** **`layout.HTMXNone` shipped**: HTMXSrc sentinel rendering no htmx runtime (no script/preconnect/SRI); `TestBaseHTMXNoneProvesHTMXOff` + self-host regression guard; CHANGELOG Added entry; AGENTS.md bullet (#282).
- **T60** `scripts/lighthouse.sh` skeleton (static-dist target, 4 categories, advisory; CI wiring deliberately deferred until 5 stable runs) (#272).
- **T61** Demo per-IP token-bucket rate limiter (`examples/demo/rate_limit.go`; knobs tuned to 10/s + burst 100 after the first values 429'd the visual audits) (#264).
- **T62** **Session-scoped demo CSRF**: per-visitor cookie (`tc_demo_csrf`, HttpOnly, SameSite=Lax), `withDemoSession` middleware, threaded through `demoPage → demoContent → kanbanDemo → kanbanDemoBoardProps`; move validation against the cookie; prerender-diff test normalized + cookie-jar client (#229).
- **T65** `TestKanbanSingleTransportBoards` — htmx-only and Datastar-only board renders pinned as goldens; shared `assertGoldenMatch` helper extracted (#258).
- **Gated rows recorded with session notes:** #303 (website clones), #312 (public HeadingTag export — recommendation: keep private), #295 (dup-gate→blocking: advisory green run 1 witnessed; needs run 2 + ratification), #224/#189 remain open with dependencies noted.

### Unplanned work the session surfaced and shipped
- CopyButton e2e race fix (see T04).
- Guard `--fix` destructive-orphan bug fix (see T25/T26).
- **templ fmt repo-wide canonicalization** (39 .templ files): BuildFlow's templ-fmt repair and the committed compact styles were fighting, desyncing the `cmd/tc/_sources` mirror on every manual commit; canonicalized once, mirror synced, goldens untouched, website pages goldens refreshed, stale `node_modules` ignore dropped from website/go.mod.
- TODO_LIST.md rows annotated: #291 (closed), #294 (partial), #296–#300, #305–#314, #258, #262, #264, #271–#281 verdicts.

---

## b) PARTIALLY DONE
1. **#229 session CSRF** — fully wired for the kanban demo endpoints, but the middleware+minting is demo-grade (no HMAC binding, no expiry); by design, documented in `csrf_session.go`.
2. **#294** — three of four sub-items done (`tc ls` footer, guard-runtime print, --min-lines verdict); "AGENTS post-daemon re-verify note" not written as a separate item.
3. **#295** (dup-gate advisory→blocking) — advisory lane green (run 1); promotion needs a second consecutive green run + owner ratification. Cannot be finished inside one session by definition.
4. **T60 Lighthouse** — skeleton + decision record only; no budgets asserted, not CI-wired (deliberate, but genuinely not "done").
5. **T07's original goal** — the big dedup-pass blob is clean-ish history but as 8 daemon-titled commits already on origin; only my later waves got proper messages.

## c) NOT STARTED
1. **Push** — 24 commits local-only; house rule forbids me pushing.
2. **#162/#150/#80** vision-review pass — needs a provider API key + human confirmation.
3. **#280** post-deploy spot-checks — needs the next production deploy to exist.
4. **#216** vnu re-triage — event-gated on the next nixpkgs html5validator bump.
5. **#217** M24/M25 components — demand-gated NO (F109).
6. **#213/#214/#215** CI budget/benchstat/mutation — deferred 2026-09-13, still deferred.
7. **Blocked ledger** (#28/29/93/107/108/124/125/126/190/191/211/212/232–239/269/270) — owner gates / separate repos, untouched by design.
8. **CHANGELOG version cut** — `[Unreleased]` is warm and healthy but no release was requested.

## d) TOTALLY FUCKED UP (honest failures)
1. **Lost the daemon race on every squash attempt** (3 attempts, 2 failed commits, exit 128 ref-lock + one daemon mid-commit) until the third attempt in a calm clean-worktree window. The original "squash the 6 unpushed daemon commits" plan goal failed outright — the daemon PUSHED them before I acted. Net effect: no data loss, but the history-hygiene goal for the main blob is dead permanently.
2. **My own rate limiter broke the visual audits** — burst 20 was smaller than one audit sweep's page+asset fetches; the audits measured unstyled pages (22px targets) and 429 error pages (missing title/lang). Caught by the suite, fixed by raising knobs — but I shipped a feature that broke the verification layer and only found out at the final sweep.
3. **My own new code introduced 4 clone groups** — the three copy-pasted `-selftest` blocks across the capture tools were exactly the duplication this repo fights. The gate caught it; the extraction (`internal/browser` + `internal/distserver.Selftest`) took three sub-iterations of lint fixes (gosec/noctx/mnd/wsl/gci/gocritic/err113) to get clean.
4. **Read `$?` after a pipe AGAIN** — the first visual run printed `VISUAL_EXIT=0` from `tail`, not nix. My own AGENTS.md documents this exact trap. Wasted one diagnosis cycle.
5. **Ran a bare `go test` inside visualtest/** and briefly treated a 0.024s "ok" as a pass — the documented skip-vs-pass trap again. Second occurrence of a documented mistake in one session.
6. **Edit-tool fumbles under daemon interference**: two edits silently didn't persist (visual-testing.md, oghot block) requiring re-reads; one edit collapsed a templ structure causing a generate failure; the kanban helper was first inserted inside a function body (syntax error). All recovered, but ~6 wasted tool cycles.
7. **Background-job budget exhaustion (50)** — piling up ~50 background shells from hook runs + suites jammed the shell tool entirely mid-commit; recovered by probing foreground and re-running, but it interrupted the flow at the worst moment.

## e) WHAT WE SHOULD IMPROVE
1. **Daemon vs pre-commit race is structural.** Every manual commit runs a 1–4 min hook while the daemon commits within seconds. Proposal: a `daemon pause` file the hook checks (or the hook writes), so atomic squashes become reliable.
2. **Squash BEFORE accumulating** — future sessions should squash daemon blobs at each milestone, not at the end, while the window is clean.
3. **Load-aware test invocation should be automatic**: a tiny wrapper that reads load average and appends `-parallel 4` (or skips timing-sensitive tests) would remove a recurring human-judgment tax.
4. **Demo feature changes need an immediate demo-suite run** — T61/T62 broke the visual audits and that was only discovered at the final sweep; the demo tests passed because they don't cover asset-fetch bursts.
5. **`node_modules`-class staleness**: two more stale-copy desyncs (starter CSS, _sources) were found this session; audit remaining hand-maintained copies (presets, theme css) for the same drift class.
6. **PIPESTATUS discipline**: my env-check scripts should capture exit codes to files, never after pipes — make it a personal checklist line since AGENTS.md alone didn't prevent it.
7. **The BuildFlow hook's transient failures under load** (tool timeouts, OOM kills) produce red commits-without-exit-128 confusion; a retry-once-on-transient mode would remove false blocks.

## f) NEXT 50 (ordered by impact)
1. Review + push the 24 local commits (human).
2. Second consecutive green advisory dup-gate run → promote to blocking (#295) + ratify.
3. #303 website content clones: extract or blanket-accept (owner gate).
4. #312 decide public `HeadingTag` export (owner gate; recommendation: keep private).
5. #224 build the sorted-view 422 demo endpoint, then the e2e.
6. #189 file-backed kanban demo state + recipe section.
7. Daemon pause-file mechanism (hook ↔ daemon coordination) — kills the race class.
8. Daemon commit messages: fix at source (`larsartmann/buildflow`) — still hallucinated titles.
9. #292/#293 art-dupl upstream: stable fingerprints + actionable-vs-raw counts (separate repo).
10. Flake-pin the art-dupl fork (`nix run .#dupl`) once its detector stabilizes.
11. Wire the Lighthouse lane into CI after 5 stable local runs (budgets asserted).
12. HMAC-bind the demo session CSRF cookie (real-app shape) if the demo keeps evolving.
13. Exempt static assets from the demo limiter explicitly (currently safe via burst; make it structural).
14. Migrate remaining ~48 raw `chromedp.Poll` sites to the poll helpers (#240).
15. #250 sibling: audit other poll-less assertions in visualtest (grep `chromedp.Text` after clicks).
16. Load-aware visual wrapper (`nix run .#visual-auto`) reading /proc/loadavg.
17. Prune the 14 vnu ignore classes when nixpkgs html5validator bumps (#216).
18. #282 follow-up: docs page showing an HTMXNone + Datastar-only demo composition.
19. AGENTS.md is 389 lines (BuildFlow warns >377) — split recipes into docs/.
20. Re-verify the demo limiter under the axe sweep after any knob change (add a smoke assertion).
21. #271 follow-up: confirm production deploy actually carries `SITE_LIVE_STARS=1` after the next deploy.
22. #281 sibling: ogshot golden — pin the OG card PNG too, not just the live capture.
23. Add `-selftest` to the smoke tool for symmetry (#305 residue).
24. Sweep docs/ for more hand-typed counts not yet drift-guarded.
25. Extract the three tools' shared `serveDist` variants into internal/distserver (ogshot injects the OG card; partial duplication remains).
26. Resolve the `TestPrerenderMatchesLiveServer` structure — the retry loop helps, but per-route subtests re-fetch; consider one client, table-driven retries.
27. Record the daemon-race squash recipe that finally worked (clean tree + calm window + single command) in AGENTS.md.
28. Split `scripts/ci-repro.sh` (now ~350 lines) — one file per lane, shared lib.
29. T60 follow-up: nix app `.#lighthouse` (chromium + npx pinned) so the skeleton is reproducible.
30. #305 residue: `--selftest` for `smoke`.
31. Add the demo limiter + session middleware to `visualtest/demo_flows_e2e_test.go` coverage (429 path untested).
32. Commit-message lint: reject `chore: auto-commit` on non-daemon authors (git hook guard).
33. website/go.mod `ignore dist` — verify CI tidy-check still green after the directive edit.
34. Check whether the earlier identical-PNG result for the two single-transport goldens is expected (attribute-level differences only) and note it in the test comment.
35. Reconcile `docs/planning/TEMPLATE.md` with the plan-authoring checklist step 5.
36. #290 residue: sweep `templates/presets/*.css` against current custom.css for the same drift class.
37. Add `-parallel 4` hint to the CI visual job (matches the local lesson; CI runners are less loaded but cheap insurance).
38. Capture-tools integration test: run all three `-selftest`s in one CI step.
39. Skill: document the session-CSRF middleware as the reference recipe for consumers adding CSRF to wired forms.
40. ADR-0009: add the 2026-09-23 selftest-wiring entry (22) to the entry-count claims if any doc cites "21 entries".
41. Prune `website/site.out.css`-style tracked artifacts if any lost their consumer (guard currently pins the set; re-audit).
42. Run `buildflow` binary refresh (preflight warns HEAD moved past the installed binary).
43. VACUUM the BuildFlow cache DB (1.34GB, preflight warning).
44. Fix `go-licenses` "not in devShell" BuildFlow warning (add to devShell).
45. Deduplicate the three tools' remaining flag/usage scaffolding if a clean builder API emerges (else keep ADR entry 22).
46. `docs/DOMAIN_LANGUAGE.md`: add session-CSRF, rate-limit, HTMXNone vocabulary.
47. Confirm CI is green on the pushed tip once the human pushes (the 24 commits ran the full local ritual, but CI is the real gate).
48. Annotate #227's anti-drift tie: the demo CSRF change altered one side of the kanban contract-parity table — verify the table comment still matches.
49. Consider `SITE_LIVE_STARS` deprecation path: document that SITE_SKIP_STARS is now the legacy alias.
50. Session retro: turn the "edit didn't persist" incidents into a rule — re-read after ANY daemon-adjacent pause before editing.

## g) Questions I cannot answer myself
1. **Push?** 24 commits are local and fully verified (ritual + visual green at tip). Do you want them pushed as-is, squashed first (I can attempt a final clean-tree squash), or reviewed commit-by-commit?
2. **#295 promotion**: do you ratify promoting the dup-gate lane to BLOCKING in ci-repro after I witness one more green advisory run (the plan's ladder), or keep it advisory for a longer probation?
3. **Demo hardening scope**: the rate limiter + session CSRF were built to demo-grade (in-memory, no HMAC). Should I invest in the real-app shape (HMAC-bound cookie, expiry, X-Forwarded-For trust behind the LB), or is the demo's current posture the intended end state?
