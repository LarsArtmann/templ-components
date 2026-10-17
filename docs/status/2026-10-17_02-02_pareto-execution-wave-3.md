# Pareto Execution — Wave 3 Status Report (2026-10-17 02:02)

**Session:** continuation of the ALL-TODOS master pareto plan (`docs/planning/2026-10-15_20-00_ALL-TODOS-MASTER-PARETO-PLAN.md`), wave 3. Resumed from the wave-2 interruption (which left M18 mid-flight and claimed a red tree).

**Headline:** M18 fully closed (all four consumer asks), M12 verified+struck, M19 fully closed (including a REAL bug fix in the shots tool), M10 was already verified in wave 2. A sibling Crush session is CONFIRMED ACTIVE (commits `8ecccef9`, `6a41414a` — a11y re-scan doc + the comment-gotcha doc, both Crush-attributed, not mine).

---

## a) FULLY DONE this session (all verified green at time of writing)

### M18 #323 — AppShell sidebar-width class migration (the interrupted task, finished)

- Tree was NOT red on resume: the templ regen had completed before the wave-2 capture; daemon committed a coherent state. Verified `go build ./layout/` green first.
- `layout/appshell_test.go`: two stale subtests rewritten to the class contract (`tc-sidebar-w-{sm,md,lg,auto}`, `AssertNotContains "style="` CSP pin, no-sidebar no-leak). One iteration needed: `AssertNotContains "tc-sidebar-w-"` tripped on the component's own explanatory HTML comment — asserted the specific class instead (sibling later documented the gotcha in `6a41414a`).
- `utils/inline_style_compliance_test.go`: appshell exemption REMOVED (the guard's stale-exemption pruning now enforces the no-inline-style state).
- Starter mirror: `templates/custom.css` byte-copied over `cmd/tc/_sources/starter/custom.css` (`TestStarterCSSMatchesTemplates` green).
- Mirror sync via `check-tc-sources-sync.sh --fix` (appshell.templ re-copied; `appshell_types.go` already in packageDeps).
- All THREE compiled-CSS targets recompiled (`.#css` for demo + pinned tailwindcss for templates/styles.css + theme); all 4 classes verified present in both outputs.
- Class inventory regenerated; adoption-guide CSP caveat rewritten (AppShell off the inline-style list, SegmentBar added — it was missing — plus the class-switch note).
- CHANGELOG entry; TODO struck. Guards: utils/layout/integration/cmd-tc all green.

### M18 #324 — safelist.css shipped and wired

- LIVE-VERIFIED the mechanism first (2026-10-01 lesson: "a copied `@source inline(...)` is a hypothesis until compiled"): `/tmp` probe with pinned tailwindcss 4.3.3 — all 9 auto-fit classes emit through `@import "./safelist.css"`.
- `templates/app.css` imports it; `tc init` scaffolds it (file list + byte-synced starter copy + new `TestStarterCSSMatchesTemplates` pair).
- `templates/styles.css` recompiled — 9 `repeat(auto-fit` occurrences verified.
- Guide section + FEATURES Grid row updated; header de-TODO'd; CHANGELOG; TODO struck.

### M18 #322 — `errorpage.FamilyFromStatus`

- Implemented in `errorpage/styles.go` (net/http constants): 409→Conflict, 429/503→Transient, 502/504→Infrastructure, 500/501/505+→Corruption (in-repo precedent), other→Rejection.
- New `errorpage/family_from_status_test.go`: 21-case table, 100-599 sweep (validity + bucket), status-level round-trip for all 6 families + the two forced collapses pinned (Infrastructure→Transient@503, Orchestration→Corruption@500).
- `ExampleFamilyFromStatus`; FEATURES row; CHANGELOG with the G1-flip note; TODO struck with the one-line-flip residue.
- Post-commit: a hook/daemon edit swapped `fmt.Sprintf`→`strconv.Itoa` in the test — benign, re-verified green.

### M18 #326 — ErrorHandler adoption paths documented

- Decision: DOCUMENT the composition (wrapper stays minimal by design; not extended, no refusal ADR).
- "Choosing an adoption path" table added to `website/content/docs/guides/error-pages.md`: wrapper / `FromError`+own shell / per-status domain messages (+`FamilyFromStatus`) / HTMX fragment exchanges → real status + `htmx.GlobalErrorHandling` toast, never a full page.
- Website link checker caught my first link form (`/docs/guides/...` → `/guides/...`); fixed; full website suite green.

### M12 #417 — achecker spike verified (sibling work)

- Cross-verified `m08-achecker-final/report.json` against the §5.3 diff table in `docs/testing/tooling-evaluation.md`: post-fix per-route numbers (/kanban 2, /forms 24, /wire 34) and rule totals (input_label_visible 48, label_name_visible 16, style_color_misuse 16) all match. Seeds #421-#423 confirmed struck. Verdict TRIAL, ledger ⫱G2. Struck with provenance. (Minor sibling-owned drift noticed, NOT touched per single-writer: the eval doc claims the spike README documents the puppeteer gotchas; they live in the script instead.)

### M19 #289 — guard UX verified (was already implemented)

- Skip banner verified by run; measured runtime prints (`142ms`); `--fix` pointer in the failure path; no stale "<100ms" prose anywhere; `scripts/test-tc-sources-guard.sh` 19/19. Struck.

### M19 #287 — kanban e2e doctrine verify

- `nix run .#visual -- -run TestKanbanE2E` → ok 18.306s (all 8 e2e tests).
- `go build ./...` green on root + ALL 8 sub-modules (incl. visualtest + website). Struck.

### M19 #347 — playground pixels pinned + REAL BUG FIXED

- **Bug found:** `.#shots` light and dark captures were BYTE-IDENTICAL — headless Chromium defaults to `prefers-color-scheme: dark`, ThemeScript pre-applies `.dark`, so every "light" capture silently rendered DARK and the in-place toggle was a no-op. Every prior `-mode both` eyeball was half blind.
- **Fix:** the tool pins light (removes `.dark` before the first capture); captures now differ (verified byte-wise).
- `/errors/playground` added to the shots page list AND to the visual route goldens (`errors_playground_{light,dark}` — additive, not a re-baseline). First run wrote, full `TestDemoRouteGoldens` re-run green (137.9s). CHANGELOG entry. Struck.

### M19 #258 — kanban transport parity (premise corrected)

- Kanban IGNORES `?transport=` (renders both boards always; the param is wire-only). Parity now pinned three ways: route goldens light+dark (pixel), 8-test e2e suite on BOTH transports (behavior), corrected shots capture. Struck.

---

## b) PARTIALLY DONE

- **Nothing mid-flight.** M19's last two strikes + CHANGELOG entry were completed in the 60-second micro-cycle BEFORE this report (lesson from wave 2 applied). Tree is green; no red state anywhere I touched.
- **Caveat:** golangci-lint has NOT been run this session on the new code (errorpage/styles.go, family_from_status_test.go, shots/main.go edits). Build+tests are green; lint is unverified — first item of the final ritual.

## c) NOT STARTED (from the plan, in order)

- M21 (#272 #373 #262), M22 (#344 #301 #307 #311), M23 (#275 #276 #279 #361 #366), M24 (#357 #359 #360 #363 #365), M25 (#410c #264 #280), M26 (#216 #213 #214 #215 + coverage), M27 (#189 #224 #229)
- M01/M09/M20 sibling verdicts; M14/M15 templ gate check; g1.1-g6.1 gate memos
- Demo-coverage KEEP gaps (datastar LoadingButton/PolledRegion sections, StatusBadgeWith demo-or-deprecate)
- Final ritual: per-module loop, ci-repro, full visual suite, closing TABLE REPORT

## d) TOTALLY FUCKED UP

- **Nothing irreversible.** Honest incidents:
  1. The wave-2 handoff claimed a red tree; reality was green. Cost: one verification beat. (Better than blindly "fixing" a healed tree.)
  2. Two edit attempts without a prior View (tool refused) — friction only.
  3. I left the demo server on :8901 running ~50 min after the captures — killed at report time.
  4. Near-repeat of the wave-2 interruption failure (M19 strikes pending when the owner interrupted) — recovered because the missing bits were docs-only; the micro-cycle lesson is now written twice.
  5. The shots-tool dark bug was PRE-EXISTING (not introduced here) — but it invalidates the "light" half of every previous `-mode both` manual audit. Nobody noticed until byte-identical PNGs were diffed today.

## e) WHAT WE SHOULD IMPROVE

1. **Run the per-module lint at task close, not only at the end** — this session shipped ~4 new/edited Go files with zero golangci-lint reads.
2. **Shots tool should self-assert light≠dark** when `mode=both` (byte comparison in-tool; would have caught the dark-default bug on day one).
3. **Kill background servers immediately after the capture that needed them.**
4. **HTML comments are component output** — negative assertions must target specific literals, not family prefixes (now also documented by the sibling in `6a41414a`).
5. **Golden-count guards may now be stale** — I added 2 route goldens (`errors_playground_{light,dark}.png`); if `TestDocsCountDrift` counts testdata/routes PNGs, the "N golden files" prose needs the bump. NOT yet checked — immediate next step.
6. The daemon/hook mutates files post-commit (strconv.Itoa swap) — always re-run the affected package's tests after any daemon commit that touches your files.

## f) NEXT UP TO 50 (priority order)

1. Run `utils.TestDocsCountDrift` — bump golden-count prose if the 2 new route goldens moved it (SAME-COMMIT rule)
2. golangci-lint per module over this session's new code
3. M21 #272 Lighthouse lane (chromedp nav-timing shape preferred)
4. M21 #373 axe-sweep handshake → chromedp v0.20 EvalAwaitPromise
5. M21 #262 flaky-board (stall/fail cards) visualtest helper + one adoption
6. M22 #344 SSE live-row recipe (Table BodyID + Body + SSE append)
7. M22 #301 ADR-0009 appendix per-group verdict table
8. M22 #307 varnamelen config sweep verdict
9. M22 #311 headingTag demo smoke additions if demo-visible
10. M23 #275 docs/visual-testing.md site tier + skip-stars pin
11. M23 #276 SKILL.md website section (search scope, lastmod, topLevelPages)
12. M23 #279 warm-dark `dark:-pair` override pattern → theming docs
13. M23 #361 /errors/404-page "Try it live" + api-reference rows + width-enum pin
14. M23 #366 tc doctor floors text + website content floors sweep
15. M24 #357 AGENTS single-writer protocol section
16. M24 #359 nonce deepening decisions memo (⫱G1-adjacent)
17. M24 #360 nonce pins (demoNonceConst assert, cross-links, fuzz)
18. M24 #363 BuildFlow health warns (lychee exclude, interrogate, rebuild)
19. M24 #365 MobileMenu unconditional script render decision + change
20. M25 #410c stale-CSS sweep remaining fleet (SKILL grep-the-artifact per repo)
21. M25 #264 demo security posture note + HTTP contract tail
22. M25 #280 post-deploy prod spot-checks
23. M26 #216 vnu ignore re-triage ritual doc + snapshot run
24. M26 #213 PR wall-clock budget runbook
25. M26 #214 benchstat runbook
26. M26 #215 gremlins mutation pilot (utils) + baseline
27. M26 coverage +2pp batch toward #28's 80%
28. M27 #189 file-backed kanban demo + dashboard recipe
29. M27 #224 kanban 422 sorted-view rejection e2e
30. M27 #229 session-scoped CSRF store for demo
31. M01 sibling verdict (spike findings → ⫱G1 adopt/reject)
32. M09 sibling verdict (tooling persistence: #412/#418/#420)
33. M20 sibling verdict (#419 crossbrowser flake app)
34. M14/M15 templ 0.3.1070 gate check (nixpkgs caught up?)
35. g1.1 decision memos (incl. #428-hue, carried #322-500 flip note)
36. g2.1 BuildFlow sprint prep; 37. g3.1 art-dupl prep; 38. g4.1 coverage plan; 39. g5.1 v2 ADR roll-up; 40. g6.1 demand-recheck cadence
37. Demo-coverage KEEP gap: datastar LoadingButton demo section
38. Demo-coverage KEEP gap: PolledRegion demo section
39. Demo-coverage KEEP gap: StatusBadgeWith demo-or-deprecate
40. Shots tool self-assert light≠dark (hardening)
41. Installation docs: tc init now scaffolds FOUR files (safelist.css) — check getting-started lists
42. Full visual suite run (all tests, not subsets)
43. Complete per-module test loop
44. `scripts/ci-repro.sh --lint --css --website` → read VERDICT: PASS
45. Closing TABLE REPORT (user-required)
46. Release decision executed (see Q1)

## g) QUESTIONS FOR THE OWNER (cannot be figured out from the repo)

1. **v1.24.0: cut now or keep accumulating?** `[Unreleased]` is heavy (Position unification, StepState, FamilyFromStatus, safelist.css, CSP-safe sidebar width, readable golden diffs, theme-inline probe, CopyButton/Table/Split a11y fixes…). Release = annotated tags for root + every sub-module, irreversible via the proxy cache.
2. **#428 warning hue — amber or yellow?** The memo recommends amber; this gates any Batch-0-style tone-table work and the two-hues-on-one-page fix.
3. **Keep `FamilyFromStatus`'s 500→Corruption** (in-repo precedent, shipped + documented) **or flip to Infrastructure** (cqrs dashboardui flavor)? One line + 3 test literals; I shipped Corruption and noted the flip in the strike.

---

**Verification snapshot at report time:** root `go build ./...` + `go test ./...` green; errorpage/website/utils/layout/integration/cmd-tc module suites green; `TestKanbanE2E` green; full `TestDemoRouteGoldens` green; website suite green (incl. link check). NOT yet run this session: golangci-lint, full visual suite, ci-repro ritual. Tree tip: `6a41414a` (+ daemon commits carrying this session's work).
