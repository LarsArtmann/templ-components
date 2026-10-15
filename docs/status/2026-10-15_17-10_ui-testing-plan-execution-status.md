# Status: UI-Testing-Tooling Pareto Plan Execution (2026-10-15 17:10)

**Session:** Execution of `docs/planning/2026-10-15_00-05_UI-TESTING-TOOLING-PARETO-PLAN.md`
(M01–M26, μ01–μ84). **Repo state at report time:** `master`, tree clean (auto-commit daemon
swept everything into `d0c2f47b`/`ee75d097`/`c480ac0f`), 3 commits ahead of origin, NOT
pushed (M03 ci-repro ritual deliberately not yet run — see §e).

**Critical context discovered mid-session:** a SECOND agent session is executing the
Tailwind-composability plan (`docs/planning/2026-10-14_20-45_*`) in this same checkout —
its `@variant dark` custom.css migration (compiled into app.css 07:13), its count-bumps
(absorbed my +2 goldens/+1 generated file), its leaked `tc-demo-route` processes, and its
website preview on :8080. Every piece of shared-tree evidence below is read in that light.

---

## a) FULLY DONE (verified)

| Task                                                 | Proof                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| ---------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Preflight: daemon pin-drift restore**              | Daemon had bumped templ v0.3.1020→v0.3.1070 (golden-flipping, forbidden by AGENTS.md) + go-error-family v0.11→v0.12 (no audit). Both restored offline from module cache; build + `TestTemplVersionPin` green.                                                                                                                                                                                                                                                                                                     |
| **M01 packaging claims (μ01–μ04)** — all four **GO** | (a) nixpkgs `playwright-driver.browsers` **empirically proven**: locked nixpkgs = 1.63.0 (chromium-1243, firefox-1543/FF 155.0, webkit-2359); `playwright-core@1.63.0` + `PLAYWRIGHT_BROWSERS_PATH` → all 3 engines launched + screenshotted. (b) rustywind v0.28.0 (2026-08-25, npm+nixpkgs+cargo). (c) stylelint+TWv4 feasible (ignoreAtRules or `@dreamsicle.io/stylelint-config-tailwindcss` ≥16.17). (d) cdproto has full `SetEmulatedVisionDeficiency` + `MediaFeature` surface.                            |
| **M02 probes (μ05–μ06)**                             | k6 v2.2.0: automatic extension resolution **works** (provisioned `explore` from registry) but **xk6-sse is NOT in the registry**; xk6 custom build blocked by this machine's Go-proxy restriction → M12 needs a Go-client stand-in locally (CI can build the custom binary). Bencher self-host/cloud cost note written.                                                                                                                                                                                           |
| **M03 close reads (μ07–μ09)**                        | SKILL.md read to EOF (961 lines); extend-not-duplicate harness inventory table written into eval doc §4.                                                                                                                                                                                                                                                                                                                                                                                                          |
| **M04 demo coverage map (μ10–μ13)**                  | `docs/testing/demo-coverage-map.md`: 136 = **124 demoed + 6 composition-covered + 6 true gaps** (datastar.LoadingButton, datastar.PolledRegion, StatusBadgeWith = KEEP; MetaRefresh, Minimal, Script = waive) + the Stylable prop gap.                                                                                                                                                                                                                                                                            |
| **M05 Stylable select demo route (μ14–μ18)**         | New `/stylable-select` route (pages.go entry, page .templ + generated, sidebar/cards/prerender auto-derived); axe sweep light+dark entries; route goldens light+dark captured **0.0000% mismatch**; smoke green on both mounts with markers (`tc-select`, `selectedcontent`, native comparison); axe finding (dark) = documented indigo-500 sidebar debt, ledgered in `axe_baseline.json` (+justifying comment); counts absorbed (FEATURES now 204 goldens / 137 generated), `TestDocsCountDrift -count=1` green. |
| **M06 3-engine spike — execution (μ19–μ24)**         | `tools/crossbrowser-spike/` (playwright-core pinned 1.63.0, version-lockstep documented, routes.json generated from pages.go, artifacts gitignored). **48/48 screenshots (16 routes × 3 engines), zero pageerrors, zero launch failures**, artifact sizes proportionally consistent (no layout blowout).                                                                                                                                                                                                          |
| **Two REAL CSP bugs found by M06 — FIXED**           | (1) `ThemeToggle("…","")` in navigation_demo + recipes_demo (3 sites) → nonce-less script dead under demo CSP → now pass `demoNonceConst`. (2) picsum.photos img-src violation (demo CSP listed only ui-avatars; plus fastly.picsum.photos redirect target) → fixed in `basepath.go` demoCSP + `website/firebase.json` `/demo/**` header (kept EQUAL per AGENTS.md). Re-shot: navigation clean in all engines.                                                                                                    |
| **M09 (μ34–μ36)**                                    | `docs/testing/tooling-evaluation.md` (verified-facts table, k6 probe, Bencher note, harness inventory) + **ADR 0048** drafted (additive-Playwright recommendation, ⫱G1–G3 placeholders, alternatives).                                                                                                                                                                                                                                                                                                            |

## b) PARTIALLY DONE

- **M06 verdict (μ25–μ26)** — run + triage data exists (console diffs per engine, size
  ratios); the written findings table (§5.1) and the "practical? effort?" verdict text are
  NOT yet in the eval doc.
- **Third CSP finding — inline `onclick="tcOpenOverlay(…)"` on the Modal/Drawer demo
  buttons** (`display_demo.templ:113,128`): CSP-dead in all engines (script-src-attr;
  Firefox reports it on /display). Root cause understood (hand-rolled demo buttons
  bypassing `display.Button` + a data-attribute pattern); the correct fix was still being
  pinned down when this report was requested (need to confirm where `tcOpenOverlay` is
  defined and whether Modal/Drawer expose a trigger wiring to bind instead).
- **Route-golden suite state** — MY route's goldens green; the SUITE still fails on
  index/layout/navigation/icons (+~0.1–1.4% pixel drift) from the CONCURRENT session's
  `@variant` migration hitting compiled CSS. Deliberately NOT rebaselined by me (their
  lane, their proof step). Shared-tree coordination question for the owner (§g).
- **Repo test surface** — visualtest route-golden + axe subset run; the FULL root-module
  `go test ./...` (incl. examples/demo tests iterating `newMux`) and the per-module loop
  have NOT been re-run since my demo additions. Required before push.

## c) NOT STARTED

M07 (Firefox base-select fallback assert — note: `firefox-light/stylable-select.png`
artifact already captured, assert itself not written) · M08 (achecker spike + axe diff) ·
M10 (aria snapshots) · M11 (CDP cheap wins into visualtest — the Go tests) · M12 (SSE load
spike, Go-client stand-in) · M13 (nix `.#crossbrowser` app + lane prep) · M14/M15 (verdict
pose + consolidation doc) · M16–M26 entirely (Bencher verdict, stylelint/rustywind trial,
trace viewer, engine hardening, achecker ledger, aria-golden class, k6 nightly, docs
sweep, trial protocol, HARVEST, ROADMAP seeds) · μ84 final checklist · CHANGELOG
`[Unreleased]` entries for the demo fixes (consumer-visible? demo-only → judge needed).

## d) TOTALLY FUCKED UP (honest list)

1. **Missed the concurrent session for ~40 minutes.** Its commits (M03/M04 labels from a
   DIFFERENT plan) were visible in git log from early on; I only understood the shared
   tree at 07:12 when its `@variant` edit collided with my generate runs. Earlier
   `git log --all` + process check would have saved several confusing loop-arounds
   (e.g. why existing goldens passed in run A but failed in run B).
2. **Hit the documented `rg -rn` footgun THREE times** (`-r` = replace, mangles output)
   despite AGENTS.md naming it verbatim. Cost: misread ThemeToggle call sites, wasted
   greps. Alias discipline needed.
3. **Port 8099 = ActivityWatch.** Burned ~4 debug cycles on a "404 that isn't the demo"
   before `ss -tlnp`. Should check listeners BEFORE starting servers.
4. **First `go get` hung 60 s+ on the blocked Go proxy and lost a background shell** —
   the module cache had everything; checking the cache first (as I did on retry) is the
   correct first move in this environment.
5. **Off-by-one in the coverage map** (17 routes vs actual 16) — caught and fixed, but it
   shipped in the first write of a doc whose entire purpose is counting precision.
6. **Let the daemon sweep my M05 route into a "chore: auto-commit 4 changed file(s)"**
   commit (82229eb7) — work is safe, but the canonical description lives only here and in
   the eval doc, not in history. Precise commits must beat the 60 s daemon budget.

## e) WHAT WE SHOULD IMPROVE

- **Two-sessions-one-tree protocol**: before executing any plan, check for concurrent
  authorship (`git log --since`, running demo/site servers, leaked test processes) and
  agree on lane ownership for shared artifacts (visualtest/testdata PNGs, compiled CSS,
  count-bump docs). The route-golden mismatch standoff is exactly this gap.
- **Demo CSP should be guarded by a test** — the picsum + ThemeToggle + onclick findings
  were all "dead under own CSP" classes invisible to string tests. A visualtest assertion
  ("zero CSP violations in console on every demo route", cross-checked per engine or at
  least Chromium) would have caught all three mechanically. Natural extension of
  `TestNoEmptyNonceAcrossDemoPages`.
- **`display.Button` needs a first-class non-wire trigger pattern** (data-attribute slot)
  so demo/consumers never hand-roll `onclick=` (the current bug's root cause).
- **Pin-drift tripwire at session start**: `rg 'templ v' **/go.mod` took 2 minutes and
  caught the daemon's forbidden bump — make it a standing first command (it already is in
  my mental checklist; formalize in AGENTS.md "session start" block).
- **Background-shell discipline**: long network commands → check module cache first,
  use explicit `timeout`, and prefer `run_in_background` + `job_output(wait=true)` over
  auto-backgrounding surprises.
- **Status of the suite**: `nix run .#visual` currently red for reasons outside my lane;
  if the other session stalls, its red will look like MY breakage in CI. Whoever finishes
  first should rebaseline with the diff-eyeball ritual and say so in `[Unreleased]`.

## f) Up to 50 things to do next (ordered by plan value)

**Immediate (finish the 1%–4%)**

1. Fix the Modal/Drawer demo buttons (replace inline `onclick` with data-attribute trigger + `display.Button`; confirm `tcOpenOverlay` definition site).
2. Write M06 findings table + verdict into eval doc §5.1 (μ25–μ26) → ⫱G1 material.
3. M07: write the Firefox base-select fallback assert (computed `appearance` + screenshot diff vs the native counterpart; already-captured artifact as input) + finding writeup (μ27–μ29).
4. Root-module test run (`go test ./...` incl. examples/demo) + per-module loop for touched modules.
5. CHANGELOG `[Unreleased]` entries (demo CSP fixes, ThemeToggle demo nonce, new route) + commit message repair notes (daemon-swept 82229eb7 described in eval doc).
6. M08: achecker spike over demo routes + finding-class diff vs `axe_baseline.json` (μ30–μ33) → §5.3.
7. M09 completion: fill §5.1/§5.2/§5.3 slots; ADR stays Proposed.
8. M12: Go-client SSE load spike vs demo binary /api/datastar/stream (steady/reconnect/client-drop; NEVER the canonical URL) → §5.5 (μ46–μ49).
9. M13 to ⫱G1 edge: `nix run .#crossbrowser` flake app wrapping the spike + kill-criteria row + CI lane yaml PREPARED (not installed) (μ50–μ53).
10. M10: aria-snapshot spike on Table/Modal/Dropdown/Tabs + regression-catch demo (μ37–μ40).

**Verdicts & consolidation (Phase 5)**
11. M14: verdict table (adopt/trial/reject × evidence links).
12. M14: additive-vs-replace recommendation section (chromedp primary, Playwright cross-engine-only).
13. M14/μ57: pose ⫱G1–G4 to Lars in ONE message.
14. M15 to ⫱G3 edge: boundary section in ADR (done in draft) + no zombie tooling sweep (μ58, μ60).
15. μ84-lite: tick plan Verification checklist items as they complete.

**Phase 6 long tail**
16. M16: Bencher keep/kill verdict (likely reject-for-now, trigger recorded) (μ61–μ62).
17. M17: stylelint run over custom.css + demo.css; triage (μ63).
18. M17: rustywind dry-run diff on .templ sources; adopt/reject note (μ64).
19. M18: trace-viewer repro of one flaky-prone e2e; keep/drop note (μ65–μ66).
20. M19: WebKit/Gecko fontconfig determinism in nix env (μ67).
21. M19: viewport policy row for engine captures (below-xl: caveat) (μ68).
22. M19: 3-consecutive-run determinism proof (μ69).
23. M20 to ⫱G2 edge: achecker ledger schema (Violation/Needs-review/Recommendation) (μ70–μ71).
24. M21 to verdict edge: aria-golden promotion docs + count wiring draft (μ72–μ73).
25. M22 to ⫱G4 edge: k6 nightly job + matrix doc (steady/reconnect/drop) (μ74).
26. M23: compat facts — `field-sizing` (all engines) + `base-select` (no Firefox) into docs surface (μ75).
27. M23: support-matrix adopt/waive decision note (μ76).
28. M23: Lost Pixel dead-tool note + rejected-tools log seed (μ77).
29. M24: tool-trial protocol doc (spike→verdict→ADR→integrate/reject) (μ78).
30. M24: CI tooling budget table (lane → owner → kill-criteria) (μ79).
31. M25: HARVEST pass 1 — TODO_LIST #411–#420 accuracy check (μ54).
32. M25: HARVEST — remaining open items → TODO_LIST/ROADMAP with IDs (μ80).
33. M25: AGENTS.md Testing section update — ONLY adopted tools (μ81).
34. M26: ROADMAP seeds — parity contracts, keyboard EXTENSION, Guidepup, contract tests (μ82).
35. M26: ROADMAP seeds — Lighthouse criteria, sampling economics, chaos, golden-churn, 6-month re-verify reminder (μ83).
36. μ84: full plan verification checklist ticked; stale notes annotated.

**Repo hygiene surfaced by this session**
37. Add CSP-violation console guard to visualtest (improvement §e-2; register as TODO item).
38. Register the 3 coverage-map KEEP gaps as TODO items (datastar twins, StatusBadgeWith).
39. Consider `nix run .#demo` app (pin-free port + no :8080/:8099 collisions) for local demo runs.
40. Clean leaked `tc-demo-route` processes from the visualtest suite (4 from the concurrent session; demoRouteBase should t.Cleanup-kill on interrupt).
41. Evaluate BuildFlow interaction of `tools/crossbrowser-spike/package.json` (node tools may start triggering on it) — verify pre-commit stays green.
42. Deconflict count-bump ownership in docs when two sessions bump simultaneously (drift test is the arbiter — note the pattern in AGENTS.md).
43. Push (after items 4–5): `scripts/ci-repro.sh --lint --css --website` at the exact commit, then push immediately (M03 ritual; covers my + concurrent session's landed work).
44. After push: verify CI green (lint/build+test/coverage + website lanes).
45. Update plan doc: tick executed tasks, annotate deviations (k6 stand-in, rebaseline ownership, M07 artifact reuse).

**Optional / deferred**
46. Dark-theme pass of the crossbrowser spike (only light ran; dark = where contrast bugs live).
47. Wire the spike's report.json into a tiny findings summarizer (md table) for the nightly lane artifact.
48. Re-check `xk6 x explore` after k6 registry updates (xk6-sse may get registered).
49. Consider `PLAYWRIGHT_HOST_PLATFORM_OVERRIDE=ubuntu-24.04` probe for WebKit determinism (not needed for launch, may matter for goldens).
50. Schedule the 6-month re-verification reminder (M26 seed) as an actual ROADMAP entry with date 2027-04-15.

## g) Questions for you (cannot answer myself)

1. **⫱G1 (crossbrowser lane):** ratify `tools/crossbrowser-spike` → `nix run .#crossbrowser` as a local diagnostic app + a NON-BLOCKING CI lane (artifacts-on-failure, kill-criteria in the M24 budget table)? And is the M21 aria-snapshot golden class approved to promote if the spike proves it catches regressions — or does that wait for a separate gate?
2. **Shared-tree coordination:** the concurrent Tailwind session's `@variant` migration currently leaves `nix run .#visual` red on index/layout/navigation/icons. Do I leave rebaselining 100% to that session's lane (my current stance), or — if it stalls — am I authorized to run `nix run .#visual -- -update -run TestDemoRouteGoldens`, eyeball the diff, and rebaseline their drift myself?
3. **Push authorization:** my work + the two CSP fixes are local-only (3 daemon commits ahead, unpushed). Push now after the root-module test + ci-repro ritual complete (docs + demo-only changes, no library API change), or hold until the whole plan's Phase 5 verdicts land?
