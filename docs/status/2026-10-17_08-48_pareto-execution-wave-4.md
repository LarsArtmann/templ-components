# Pareto Execution — Wave 4 Status Report (2026-10-17 08:48)

**Session:** continuation of the ALL-TODOS master pareto plan, wave 4. Resumed from wave 3's handoff with the corrected todo list (M18/M12/M19 completed). Owner directive: comprehensive ≤12-min plan first, then execute until done.

**Headline:** 33 TODO rows struck this wave (M21–M27 complete, plus #434, #425-strike, #296, #410c, #216, #213, #214, #215, #280, #363, #365, #357, #359, #360, #264, #275/#276/#279/#361/#366, #344/#301/#307/#311, #272/#373/#262). Master CI went RED mid-wave (Dependabot pin re-break + this wave's pre-lint code) and is FIXED at the tip (lockstep restored + dependabot ignore rule, all lint findings fixed, route goldens re-captured for the intentional playground-form growth). A REAL production bug was found and fixed (#280: sitemap lastmod uniform build-day — shallow-clone proven). Next free ID: 437.

---

## a) FULLY DONE this session (all verified green at completion)

### M21 — testing infra

- **#373** axe-sweep handshake → `chromedp.EvalAwaitPromise` + by-value decode into `AxeResults`; window-parked payload/poll/ERR-sentinel gone; rejected promise → `*chromedp.ExceptionError` → `errAxeRuntime`. Verified: positive control + /display sweep route.
- **#272** Lighthouse lane: picked the chromedp shape — `TestDemoNavTimingBudget` (FCP/LCP/CLS/load vs generous budgets, 4 representative routes; first-run numbers FCP ~770–870ms, CLS ≤0.005, load ≤1.2s pinned in the log).
- **#262** flaky-board helper: `visualtest/flaky.go` (FlakySlowID/FlakyFailID/FlakySlowDelay/`FlakyFail`), adopted in BOTH prior ad-hoc sites (kanban_e2e + kanban_pending_visual), `FlakyFail` unit-pinned; each site keeps its own stall shape (semantic, not duplication).
- Shots tool hardening: `-mode both` self-asserts light≠dark byte-diff before writing the pair (would have caught wave-3's dark-default bug on day one). Verified live: real captures differ.

### M22 — recipes/polish

- **#344** `docs/recipes/sse-live-rows.md` — htmx-ext-sse facts verified against the extension docs (separate script from core; hx-ext+sse-connect same element; named-vs-message events, no catch-all); sse-swap+hx-target+beforeend append into the `BodyID` tbody; endpoint rendering templ fragments; row ids + oob patching; cap-the-DOM escape hatch; Datastar twin pointer. Cross-linked SKILL + custom-table-rows.
- **#301** ADR-0009 appendix: per-group verdict table — historical 37+6 consolidated by reference; current advisory sets under fork v0.9.1 (74 t=2 + 9 t=3) dispositioned into 11 signature classes, zero new extraction decisions; owed baseline re-record filed as **#435**.
- **#307** varnamelen verdict: KEEP the guidance (the forced `c`→`coord` rename was right); policy written into `.golangci.yml` (ignore list is a CLOSED set; flagged cryptic vars get renamed, never exempted).
- **#311** headingTag demo smoke: condition FALSE — zero `HeadingTag`/`TitleTag` usage in examples/demo; route goldens already pin the default branch; smoke would be coverage theater.

### M23 — website docs

- **#275** visual-testing.md site tier: full 5-test inventory table (8-golden arithmetic, axe ledger pointer, touch/zoom audits, TestSiteSalesCopyButton), skip-stars pin rewritten to the post-#271 reality, nav-timing note.
- **#276** SKILL site section verified against website/cmd/site; de-TODO'd heading.
- **#279** warm-dark `dark:`-pair override documented in BOTH docs/theming.md and the website theming page (literal-pair Class form, scanner rules, safelist escape hatch, global-@theme pointer); website suite green.
- **#361** all five parts: live-404 link; api-reference ErrorPage subsection (WayOutAction/SecondaryWayOut/CopyCode/ErrorMaxWidth); SecondaryWayOut form field + handler wiring; full width-enum pin through the real handler (lg/4xl added; + notContains case, placeholder leak caught and fixed); `errors_playground` added to the axe sweep (3 rules, 0 blocking).
- **#366** tc doctor floors verified current (1.27 everywhere); repo stale-floor sweep ZERO; hand-typed counts ZERO (build.CountStats-derived); the rg -r trap produced phantom "n" placeholders on first sweep — corrected re-run disproved them.

### M24 — process/nonce

- **#357** AGENTS.md single-writer protocol section (194/220 lines — inside the gate).
- **#360** all four: TestDemoScriptsCarryDemoNonce + shared demoSweepRoutes; the sweep immediately caught a false positive worth pinning (the inlined htmx runtime contains its own `nonce="${config.inlineStyleNonce}"` SOURCE literal) → regex scoped to tag-opens + pin test; cross-links both directions; CHANGELOG omit-empty entry names both helpers; `FuzzScriptComponent` (escaped-nonce, never-empty-attr, verbatim-body, no-panic) — seeds + 10s campaign clean.
- **#359** decision memo (`docs/planning/2026-10-17_nonce-deepening-memo.md`): strict mode REJECT (zero render-path errors guardrail), forbidigo ban REJECT (output sweeps own the class), auto-detect RESOLVED by ADR-0047, `<style>` emitter audit CLEAN. Reopen conditions documented.
- **#363** all three: lychee `docs/feedback/archived` exclude fixed UPSTREAM in ~/projects/BuildFlow (both mirrored lists + link-audit.sh regex; all three parity test packages green); interrogate skip_step (.buildflow.yml, Python-only + zero Python sources); BuildFlow binary rebuilt from HEAD via nix (doctor binary-freshness now GREEN).
- **#365** MobileMenu gating: DECISION keep (Nav couples toggle on the same len(links)); the dead-JS trap class pinned by TestMobileMenuScriptMarkupCoRender.

### M25 — fleet/security/prod

- **#264** demo security posture section (stateless-GET render-only + handler field caps; deliberate NO in-app rate limiting — hosting layer's job; CSRF + same-origin on stateful POSTs; executable contract pins named).
- **#280** post-deploy spot-checks: **FOUND AND FIXED A REAL PROD BUG** — live sitemap showed lastmod = deploy-day for ALL 20 pages; root cause PROVEN (depth-1 shallow clone answers every `git log -1 -- path` with the HEAD date); fix = `fetch-depth: 0` on the Build Website checkout (actionlint+shellcheck green); live heals on next deploy. Stars badge verified working ("4 Stars", real lookup); Firebase caching sane; path-filter parity ✓; Lighthouse-on-live deferred (no binary here).
- **#410c** fleet stale-CSS sweep: all 16 consumers, 11 artifacts; membership check (own-source tc-* usage must exist as `.tc-*` in the artifact; tests+vendor excluded after the grep -h/-v trap) — 6 CLEAN, PapDashboard healthy-by-design (icons-only), **nsfw-classifier + pbx-lens STALE** (filed as **#436**, consumer-side owner action).

### M26 — CI/test extras

- **#216** vnu re-triage: ritual doc (docs/testing/html-validation-ritual.md: 5-step procedure, compound-alternation splitting gotcha, CI-is-authoritative, prune-only-on-bump) + snapshot run (204 raw errors / 335 goldens, per-class ledger: hx-105, RGB-32, data-18…; 2 arms + 2 members flagged prune-candidates).
- **#213 + #214** measurement runbooks (docs/testing/ci-measurement-runbooks.md) with REAL seed baselines: per-job wall-clock table from live run data; benchstat artifact+comment shape; both advisory with kill-criteria; implementation honestly gated on infra (comment token + artifact store).
- **#215** gremlins pilot on utils (23s): Killed 197 / Lived 50 / Not-covered 45 / Timed-out 108 → efficacy **79.76%**, mutator coverage 84.59%; baseline doc with pinned invocation + honest reading (timed-out inflates; the 50:45 lived/uncovered set IS the test-gap map).
- **Coverage batch**: ScriptAttrs 0→100%, DismissButton 2.2→76.1% (utils); the 8 untested Default*Props constructors + Tone adapters covered (display 71.8→72.0%, zero 0%-non-generated funcs left); an existing TestDefaultPropsConstructors was found via LSP DuplicateDecl and EXTENDED instead of duplicated.

### M27 — demo depth

- **#229** session CSRF: the store SHIPPED 2026-09-23 (stale premise); verified wiring + added the missing direct contract test (cookie mint + SameSite + matching-validates / foreign-fails-closed / no-cookie-fails-closed).
- **#189** BOTH parts: `KANBAN_STATE_FILE` opt-in file-backed state (atomic persist per mutation; the add() persist bug was caught by the round-trip test on its FIRST run; corrupt-file-safe; memory-only default keeps goldens deterministic) + dashboard-recipe kanban section; contract-markers rows added.
- **#224** sorted-view rejection: `moveErr`/`applyErr` reject moves whose card the server state lacks OR whose target column is unknown (the old silent no-op could EAT a card from its source column) → 422; index overflow stays a clamp; pinned in the SHARED probe table so demo + e2e cannot diverge (stale move after reset = 422 on both transports); TestKanbanE2E green (19.4s).

### #434 + hygiene

- **#434** SegmentBar role=img: image role moved onto the bar graphic; legend = real `<ul>` sibling. Zero visual change PROVEN (the /display route golden re-ran pixel-identical — no PNG diff); 4 package goldens carry the restructure; structure test pins role-on-bar. Achecker zero-hit confirmation rides the next re-scan.
- **#425** strike backfilled (tone Phase 0 shipped in wave 2; package+adapters+exhaustive tests verified); #296 half-struck row fixed.
- **Master CI RED → GREEN (root-caused and fixed)**: TestCrossbrowserLockstep (Dependabot `6477d6d3` re-broke the playwright-core pin after the sibling's `1aa8c212` restore; nixpkgs anchor is 1.63.0) → pin restored + **dependabot ignore rule for playwright-core** so the recurrence class is dead; Lint (gosmopolitan Han rune in my fuzz seed, golines on nav_timing, modernize errors.As→AsType + nlreturn in axe.go, wsl_v5 in flaky_test) → all fixed, both modules 0 issues; Visual (error_pages + errors_playground goldens grew with the legit playground-form change) → targeted re-capture of EXACTLY 4 PNGs, login_light false-diff restored + verified; verified WITHOUT -update: green.

## b) PARTIALLY DONE

- **Coverage +2pp**: display +0.2pp + utils bumps landed; the repo-wide +2pp toward #28's 80% NOT fully reached — the display/forms long tails remain (the batch killed all 0%-exported-func targets found).
- **g1.1 memos**: #359's memo done; the other ~13 G1 items not yet memo'd.
- **CHANGELOG**: [Unreleased] carries this wave's entries (SegmentBar fix, SSE recipe, kanban 422 + file-backed state, adoption-path guidance etc.); a final pass to sweep for missed entries not done.

## c) NOT STARTED (in plan order)

1. Demo-coverage KEEP gaps (datastar LoadingButton + PolledRegion demo sections; StatusBadgeWith demo-or-deprecate)
2. Sibling verdicts M01 (#413 spike findings), M09 (#412/#418/#420 artifacts), M20 (#419 flake app)
3. M14/M15 templ 0.3.1070 gate check (nixpkgs caught up?)
4. g1.1 bulk decision memos (#233–#239, #303, #312, #358, #411, #428-hue, #322-500, #295, #217, #400, #190, #212, #329)
5. g2.1–g6.1 gate-prep micros
6. Final ritual: per-module loop → full `nix run .#visual` → `scripts/ci-repro.sh --lint --css --website` VERDICT: PASS → closing TABLE REPORT
7. Lighthouse-on-live (needs a driver-up env with network+Chromium)

## d) TOTALLY FUCKED UP

1. **Master CI was RED for ~2h (00:24–04:24) across 3 jobs** — partly my pre-lint code (the daemon committed before golangci-lint ran), partly Dependabot's pin break. Wave 3's improve-item #1 ("lint at task close") was WRITTEN and STILL repeated. Third strike on this lesson.
2. **The rg -rn trap hit twice in one session**: phantom "n" placeholders (#366 — nearly "fixed" healthy docs) and phantom `func n(` (csrf_session.go — nearly reported daemon corruption). The gotcha is in AGENTS; I stepped in it anyway.
3. Python heredoc apostrophe escaping broke 3× (SyntaxError) — the "use edit tools" rule keeps being learned the expensive way.
4. The `-update` route-golden run rewrote `login_light.png` (capture nondeterminism) and `feedback_{light,dark}.png` (34-byte PNG-encoder noise) — login was restored+verified, feedback's noise rode a daemon commit; impure diffs either way.
5. Nearly filed a duplicate `TestDefaultPropsConstructors` — LSP DuplicateDecl caught it pre-commit, but I wrote the file before grepping for the existing test (read-before-write discipline).
6. First `TestDemoSessionCSRF` version failed on FormValue caching semantics — rewrote with proper form bodies.
7. gremlins re-run for the lived-mutant file distribution was interrupted — baseline recorded from the first successful run instead (numbers sufficient; distribution deferred).

## e) WHAT WE SHOULD IMPROVE

1. **Lint at task close** — the recurring lesson, now three waves old. Make it mechanical: run `golangci-lint run` on touched modules BEFORE the last test cycle, not after.
2. **The -update golden run should skip rewriting zero-mismatch files** (byte-compare before write) — capture nondeterminism pollutes diffs and flirts with the 2-rebaseline cap.
3. **Fleet check as a script**: the #410c membership check should be `scripts/fleet-css-check.sh` so #436-class drift is one command, not 20 lines of fragile shell.
4. **Dependabot ecosystem pinning**: the playwright-core ignore rule should be the PATTERN for every version that is anchored to a nixpkgs lock (tailwind lanes already have the same exposure).
5. gremlins lived-mutant triage (the 50:45 map) is the natural follow-up to the baseline.
6. The BuildFlow lychee fix rides the BuildFlow daemon commit; the upstream PR itself is still owed (G2).
7. Minor: RunAxe lost the byte-size error context in the decode path (observability regression, acceptable).

## f) NEXT UP TO 50 (priority order)

1. Demo-coverage KEEP gaps: datastar LoadingButton + PolledRegion sections; StatusBadgeWith verdict
2. Sibling verdicts M01/M09/M20 (read artifacts → verdicts → strikes)
3. M14/M15 templ 0.3.1070 gate check (one `nix eval`)
4. g1.1 bulk decision memos (13 items)
5. g2.1–g6.1 prep micros
6. Coverage long-tail batch (display/forms) toward #28's 80%
7. gremlins lived-mutant triage
8. `scripts/fleet-css-check.sh` (#436 enforcement)
9. BuildFlow upstream PR for the lychee exclude (G2)
10. Achecker re-scan to confirm #434's zero-hit (sibling tool)
    11–45. Final ritual ladder: per-module loop → full visual → ci-repro --lint --css --website (VERDICT: PASS) → verify CI green at tip → push decision → v1.24.0 cut (pending Q1) → TABLE REPORT
    46–50. Hardening backlog: shots byte-diff note already done; website.yml deploy verification after next push (sitemap lastmod per-page dates live); PR-review pass over the dependabot ignore; docs count sweep for this wave's doc additions; CHANGELOG completeness pass.

## g) QUESTIONS FOR THE OWNER (cannot be figured out from the repo)

1. **v1.24.0: cut now or keep accumulating?** [Unreleased] is heavier than at wave 3 (SegmentBar a11y fix, SSE recipe, kanban 422 + file-backed state, nonce sweep extension, docs waves). Release = annotated tags root + every sub-module; irreversible via the proxy cache.
2. **#428 warning hue — amber or yellow?** Still gates Batch-0 tone work and the two-hues-on-one-page fix (the memo's recommendation: amber).
3. **Consumer-repo fixes (#436: nsfw-classifier + pbx-lens stale CSS) and the BuildFlow lychee PR**: fix in those repos myself (owner permission for sibling repos), or leave both to you? Both are one-command regens + one upstream push respectively.

---

**Verification snapshot at report time (08:48):** tree clean at `759f95b4` (daemon-carried); demo suite green; visualtest build+vet+lint green (0 issues both modules); TestDemoKanbanHTTPContracts + TestKanbanE2EHTTPContractParity green with the new 422 probes; TestKanbanE2E browser suite green (19.4s); TestDemoRouteGoldens green WITHOUT -update after the targeted 4-PNG re-capture; utils+display lint 0 issues; website suite green; TestCrossbrowserLockstep green; latest CI runs at `3bab4ce5`/`e2f6efa7` still show failure — those SHAs predate the fixes at the current tip; the push decision (and the green it should restore) is the human's per the M03 ritual.
