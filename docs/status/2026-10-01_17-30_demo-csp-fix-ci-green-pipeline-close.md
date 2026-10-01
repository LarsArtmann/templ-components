# Status Report — Demo MPA Fix Pipeline Closure: CSP Nonce Regression Fixed End-to-End, CI Green, Live Verified

**Session date:** 2026-10-01, 17:30 CEST (third session today; closes the 07:08 MPA build and the 12:35 fix-pipeline sessions)
**Scope:** Finish the verification pipeline for the demo-MPA + Firebase Hosting→Cloud Run work: triage every red gate to root cause, fix forward, witness green locally, push, watch CI to fully green, and re-verify the fixed deploy live.
**Repo:** templ-components @ master, `7f282d71` pushed (all session work on origin). Post-report activity by other actors: 3 daemon commits (13:51–15:12) and a **currently-dirty `flake.nix`** (aarch64-darwin removal being reworked — not this session's edit; see §c/d).
**Predecessor reports:** `2026-10-01_07-08` (MPA build) → `2026-10-01_12-35` (regression found + fixes, + Addendum 2) → **this report** (pipeline closure + honest close-out).

---

## Executive summary

The continuation session's mandate — "keep going until everything works" — is met for the pipeline: **CI (ci.yaml) is fully green for the first time today** (run 36857280299: Build & Test ✓ Lint ✓ Visual Regression ✓ CSS Freshness ✓ HTML validation ✓), the fixed demo + site are **deployed and live-verified** (0 empty-nonce scripts on every probed page — was 50; fresh CSS served; POST fragments and dual-mount work; SSE still buffered as documented), and `ci-repro.sh --lint --website` printed **VERDICT: PASS** at the pushed tip.

Getting there surfaced and fixed six more gates **after** the 12:35 report: the daemon's go.mod flip-flop had broken directive equality AND workspace toolchain resolution (aligned everything at `1.26.0`); the demo CSS artifact was one compile behind the new utility shades (CSS Freshness gate); the ViewTransitions component violated two fresh-vnu rules (style-not-first-child — fixed forward; `::view-transition-*` pseudo scoping — impossible to satisfy in markup, became the 16th documented ignore class); the tc-scaffolder mirror drifted behind a daemon commit (self-healing guard synced it); and the committed index route goldens had baked a stale-CSS render (regenerated).

Stat line: **6 post-report gates triaged to root cause and fixed · CI red→green across 5 pushes · live deploy healed and verified · 2 permanent CSP guards landed · several of my own process fuckups (§d) — the CSS-freshness miss and the daemon-race mirror drift being the most instructive.**

---

## a) FULLY DONE

1. **Kanban e2e validation of the nonce fix** — `TestDemoKanbanMoveButtonsBothTransports`, `TestDemoRTLKanbanMoveWorks`, `TestDemoKanbanHTTPContracts`, `TestKanbanE2EHTTPContractParity` all green. The move-buttons pair was the original signal of the CSP regression; it now passes because the kanban script executes.
2. **Two remaining hardcoded probes fixed** — the shared contract driver still verified the added card via `probe.get("/")` in two places; both moved to `fx.pagePath` (demo `/kanban`, parity twin `/`).
3. **Permanent CSP guards landed** (`visualtest/demo_csp_nonce_test.go`): `TestDemoCSPNonceIntegrity` sweeps 17 demo pages over HTTP asserting (a) the CSP header carries the demo nonce source and (b) **zero `<script nonce="">`** — the exact shipped-broken shape. `TestDemoCSPJSExecutes` proves in a real browser that `window.tcKanbanAttached` becomes true (inline scripts actually execute). Lint findings on the new file (noctx, tparallel) fixed; temporary diagnostics deleted (`tmp_kanban_diag_test.go`, `tmp_nonce_inventory_test.go` — the inventory was promoted, the diag was scaffolding).
4. **Full local verification battery green** — workspace build, `nix run .#lint` (0 issues, all modules), `nix fmt` no-op, `nix flake check` pass, all 6 sub-modules green (after aligning go directives), root module green, website module green (sales golden regenerated for the scrollback shade flow-through), full clean visual suite green twice.
5. **ci-repro VERDICT: PASS** at the pushed tip (`nix develop -c scripts/ci-repro.sh --lint --website`; coverage 71.5% ≥ 70%; changelog warmth satisfied; art-dupl advisory noted, non-blocking).
6. **Pushed and CI fully green** — 5 pushes fought the pipeline red→green; final run 36857280299: **all five jobs ✓** (Build & Test, Lint, Visual Regression, CSS Freshness, HTML validation). This ended a day-long red streak (every push since 05:12Z had failed).
7. **Website + demo redeployed with fresh CSS** (run 36853607372 ✓) — live CSS carries emerald-700/amber-700/green-700.
8. **Live behavioral re-verification of the fixed deploy** — 0 empty-nonce scripts across all probed canonical pages (was 50 occurrences across 10 distinct scripts); fresh CSS served; POST `/demo/api/wire/form` 200 with fragment through the proxy; site `/` 200; all four standalone recipe screens 200 via `/demo/recipes/*`; `/demo/errors/404|500` return their intentional demo status codes; raw run.app dual-mount intact.
9. **Six late gates root-caused and fixed** (details in §d where I caused them): go directive alignment (`go 1.26.0` across root/go.work/visualtest + website tidy); demo CSS recompile (one-line minified artifact, exactly the new shades); ViewTransitions style-before-script reorder (+ goldens, + htmx tests green); the 16th documented vnu ignore class (`::view-transition-*` pseudo scoping — impossible in markup, AGENTS count note updated); tc-mirror resync for `view_transitions.templ`; route golden refresh (index had a stale-CSS capture).
10. **Docs closed out** — status reports chained (07-08 → 12:35+Addendum-2 → this); AGENTS.md updated (nonce-passing convention, SSE limitation, `-update` flag-order gotcha, ignore-class count); CHANGELOG `[Unreleased]` carries the CSP-regression fix, the contrast pass, and the smoke-probe addition; TODO_LIST gained #329 (SSE bypass ⫱) and #330 (latent tone-class contrast sweep).

## b) PARTIALLY DONE

1. **Post-fix CI on subsequent daemon pushes — unverified at report time.** After my green run, the daemon pushed documentation commits (13:51, 15:09) and a `flake.nix` change (15:12); the triggered CI + Website runs are **queued** (6+ min in queue, verdict pending). My verified-green tip is `7f282d71`; everything after it is other actors' work.
2. **`flake.nix` is dirty right now** (not my edit): the worktree removes `aarch64-darwin` from `systems` and rewrites the comment to "Linux-only … nixpkgs does not ship Chromium for darwin" — on top of the daemon's 15:12 commit which already touched flake.nix. Whomever it belongs to (user or another agent), it is unpushed and unverified by me; I deliberately did not touch or revert it.
3. **CHANGELOG/version state** — `[Unreleased]` is warm and accurate, but **v1.19.5 is not cut** (library-visible contrast + ViewTransitions-order changes are in it). The version triple-bump + tag await owner go-ahead per the release ritual.
4. **ci-repro lane coverage** — `ci-repro --lint --website` does NOT reproduce four ci.yaml-only jobs: CSS Freshness, HTML validation, the tc-mirror Go tests, and the Visual Regression suite. I ran the missing ones ad hoc during this session (that's how each was caught locally — but only after CI flagged them, because I did not know to run them). The gap is fixable (§e2, §f items 3–5).
5. **TODO_LIST/ROADMAP harvest** — the 07-08 report's ~50-item §f and this session's additions are only partially transferred (#329, #330 landed). The consolidated docs-health pass remains open.
6. **AGENTS.md size** — BuildFlow now warns it is 413 lines vs the 377 budget (this session added the nonce convention, SSE limitation, and `-update` gotcha). The content is right; the budget needs a split-out pass.

## c) NOT STARTED

1. Owner decisions queued: **SSE bypass** (TODO #329: direct Cloud Run domain vs App Hosting vs raw-run.app EventSource + CORS + `connect-src`), run.app canonical-redirect policy, persistent site-header "Demo" link, release timing for v1.19.5.
2. **TODO #330** — sweep of demo-unrendered sub-4.5:1 text tone classes (ScrollbackToneSuccess-style patterns) with golden cascade.
3. Post-deploy smoke hardening: JS-singleton probe + SSE canary in `website.yml` (currently my new probe covers status/shell/CSP only), and the same in `ci-repro --website`.
4. Axe baseline re-triage: the `-1` waivers (index, index_dark, forms_dark) predate the per-page strict standard; policy says budgets or deletion.
5. `templ.GetNonce(ctx)` fallback in `scriptComponent` — needs an ADR or an explicit NO with rationale (would have prevented the entire regression class at the library level).
6. Docs sweep for the old `go test -update <pkg>` form outside AGENTS.md (recipes/, docs/) — AGENTS.md is corrected; other docs unchecked.
7. Pre-existing deferred items untouched: stale OG card, bun shim, PR-wall-clock budget, awesome-templ coverage gate (#28), TODO #250 prerender determinism, TODO #236 errorpage presentation.

## d) TOTALLY FUCKED UP (all repaired; honest ledger — this session's new entries only)

1. **I forgot to recompile the demo CSS after changing component classes — despite reasoning my way OUT of it out loud.** I checked that `text-amber-700`/`text-green-700` existed in the artifact and concluded "class set unchanged, artifact stable." Wrong: `bg-emerald-700`, `text-gray-600` (new usage counts), and the CountBadge pair were NOT in the compiled output. The repo's own gotcha ("after ANY class-affecting change, recompile") existed; I even cited a variant of it. Cost: a red CI run. **Lesson: run the artifact gates (CSS freshness, HTML validation) locally BEFORE pushing — they are cheap; "I reasoned it" is not a check.**
2. **I pushed through the daemon's sweep and shipped an unsynced mirror.** My ViewTransitions fix was committed by the auto-commit daemon (which bypasses pre-commit), so the self-healing tc-mirror guard never ran; CI's `TestSourcesMatchPackageFiles` + the guard self-test failed. The guard is *designed* to prevent exactly this and *did* fix it in 90ms when finally run. **Lesson: anything left uncommitted when the daemon sweeps lands UNGUARDED — commit promptly through the hook, or run the guards before every push.**
3. **My `-update` invocation over-matched and I accepted it without understanding why.** `-run TestDemoRouteGoldens/index` rewrote forms/users goldens too (Go's pattern matched more subtests than my mental model). I eyeballed the diffs (legitimate changes, accepted correctly) but never explained the mechanism — "it worked but I don't know why" is a debt I papered over.
4. **The go.work sed fixed the wrong side of the skew.** I aligned go.work `1.26.0`→`1.26` to satisfy the equality guard, breaking workspace toolchain resolution entirely (visualtest requires ≥1.26.0; GOTOOLCHAIN=local refused). 20 failed BuildFlow steps mid-ci-repro. Correct fix: root go.mod UP to `1.26.0`. The daemon flip-flop warning ("go line changed 20 times in 20 commits") was flashing in the same log I read — I fixed my symptom and left the flipper armed (see §g2).
5. **I re-cut the witness-run verdict with `| tail -5` and lost the final line — the exact documented lesson, third strike.** The run's ok/FAIL verdict was cut off; load had spiked to 49; I had to re-run. Rule restated: a verdict is only witnessed when its summary line is IN the capture.
6. **The committed index route goldens contained a stale-CSS render and I don't fully know how it got there.** The 12:37-mtime golden shows the home filter row stacked (sm:flex-row not applied) while my source-correct render is side-by-side; the clean witness run had PASSED before it, and I cannot reconstruct which writer produced it (no -update run of mine was active at that minute). I fixed the state (regenerated, verified, pushed) but not the mystery — an unexplained writer in the golden path is exactly the daemon-class hazard this repo keeps documenting.
7. **ci-repro needs `nix develop` and I ran it bare twice** (actionlint missing; one false BUILD-OK from `go build | head && echo` — head's exit code masked a real compile failure with struct-literal errors). Both fixed (nix develop invocation; `if go build; then` + stderr file), but I burned cycles re-learning documented lessons instead of applying them.
8. **Struct-literal promoted-field mistakes, four files deep** (ErrorAlertProps, ToastProps, NavProps, SimpleNavProps, DirtyGuardProps — `Nonce:` set directly instead of nested `BaseProps`). The repo documents this exact trap; the compiler caught every instance after I'd already written them.

## e) WHAT WE SHOULD IMPROVE

1. **Close the ci-repro lane gap (highest leverage):** add `--ci-extras` (or fold into default) reproducing CSS Freshness (`nix run .#css` + `git diff --exit-code`), HTML validation (`check-html-valid.sh`), the tc-mirror Go tests, and optionally the visual suite. Today ci-repro says PASS while four CI jobs remain unproven locally — this session's red streak is directly attributable to that blind spot.
2. **Daemon-race defense for guarded trees:** a pre-push guard (or ci-repro step) that runs the self-healing checks (mirror sync, CSS compile diff, generated-file sync) immediately before `git push`, so a daemon sweep can never land unguarded content that CI then rejects.
3. **Answer the golden-writer mystery:** add mtime/producer logging or a lockfile note to the route-golden `-update` path; an unexplained 12:37 writer means some tool captures goldens without a human asking (BuildFlow visual step? a stray background job?). Until identified, `-update` runs should pin a clean `git status` first.
4. **AGENTS.md diet:** move the three 2026-10-01 additions into `docs/` (e.g. `docs/demo-csp.md`) and leave pointers, restoring the ≤377-line budget.
5. **Behavioral smoke in CI's deploy job:** extend the website.yml post-deploy smoke with one JS-execution assertion (singleton flag) and an SSE canary with expected-failure annotation — my session proved HTTP-shape checks miss dead-JS outages.
6. **`GetNonce(ctx)` fallback decision** (ADR) — library-level immunity to "forgot to pass nonce" regressions.
7. **Keep the go-mod flip-flop on a leash:** the BuildFlow warning says the `go ` line changed 20 times in 20 commits; my `1.26.0` alignment will un-stick CI today but the flipper remains. Identify the writer (a BuildFlow provider or the daemon's go-mod-update) and pin or skip it deliberately.
8. **Watching-pipeline discipline:** after push, WATCH the CI run to its verdict before opening new work streams — I chained fixes reactively; a single `gh run watch` + one triage pass would have collapsed the 5-push red→green ladder into 2–3 pushes.

## f) TOP 50 THINGS TO GET DONE NEXT

_Carried forward from the 12:35 report (items survive; the pipeline thread 1–8 is now done) plus this session's additions. TODO_LIST next free ID: 331._

**Close the loop on today's push (1–4)**
1. Watch the queued CI/Website runs for the daemon's 15:12Z push (flake.nix + docs) — first unverified code-adjacent change since green.
2. Resolve the dirty `flake.nix` (aarch64-darwin removal): confirm intent (user/other agent), verify `nix flake check` on all systems it claims, commit or discard — it currently blocks a clean tree for the next pipeline run.
3. Extend ci-repro with the four missing ci.yaml lanes (CSS freshness, HTML validation, tc-mirror tests, visual suite) — closes §e1; prevents the entire red-streak class.
4. Add the pre-push guard running self-healing checks before `git push` (§e2).

**Regression hardening (5–12)**
5. Post-deploy JS-singleton probe in website.yml smoke (§e5).
6. SSE canary with expected-failure annotation until #329 lands.
7. Re-triage the axe `-1` ledger entries into budgets or deletions; document in the a11y policy.
8. Execute TODO #330 (sub-4.5:1 text tone sweep + golden cascade).
9. `-update` no-op canary in utils/golden (warn when an update request rewrites nothing while passing).
10. Prerender↔live nonce parity in `TestPrerenderMatchesLiveServer`.
11. `GetNonce(ctx)` fallback ADR for `scriptComponent` (§e6).
12. Identify the golden-writer / pin clean-tree requirement for `-update` runs (§e3).

**Canonical URL & hosting (13–17)**
13. Owner decision + implementation: SSE bypass (TODO #329).
14. run.app → canonical `/demo/` redirect decision (currently dual-serve).
15. Persistent "Demo" link in the website header (owner).
16. Demo sitemap/robots/canonical-tag story under the site domain.
17. Byte-equality guard: `demoCSP` (Go) vs firebase.json `/demo/**` header.

**Release & docs (18–23)**
18. Cut v1.19.5 once the queued runs settle (contrast + ViewTransitions order + CSP demo fixes are consumer-visible; `[Unreleased]` is warm).
19. AGENTS.md diet: move today's additions to docs/, restore ≤377 lines.
20. Sweep docs/ + recipes/ for the old `go test -update <pkg>` flag order.
21. Update `docs/testing/a11y-gate-policy.md`: sweep now runs animation-settled; axe ledger re-triage notes.
22. Fold the 07-08 §f + 12:35 §f + this §f into ONE TODO_LIST backlog (docs-health pass).
23. Record the "daemon sweeps land unguarded — commit promptly through the hook" lesson in AGENTS.md BuildFlow section (it has six sibling entries; this is the newest shape of the same hazard).

**Testing & quality (24–33)**
24. TODO #330 execution detail: grep `text-green-600|text-amber-600|text-red-600` text-size pairs in style maps; fix TEXT usages; regen.
25. Mobile-menu open/close e2e on shell pages (both breakpoints) — the shell is now the single point of failure for 14 pages.
26. Touch-target audit: deliberate header/footer sweep across all pages (today's hit was incidental).
27. Coverage push toward the 80% awesome-templ gate (71.5% overall today; measure per-package delta).
28. Watch `waitAnimationsSettled()` runtime cost in the axe sweep (21 routes × settle budget) — trim if CI budgets grow.
29. Reusable `visualtest/probe` helper: script inventory + singleton flags + SSE canary as one tool (the pieces of items 5/6/9).
30. Document the Chrome nonce-masking quirk (`getAttribute('nonce')` lies in-page) in docs/javascript-guide.md.
31. Review other go1.26 flag-semantics changes against scripts/*.sh.
32. Verify no `| head`-gated exit codes remain in scripts/ (the false BUILD-OK pattern).
33. Vision-review pass over the refreshed route goldens (`scripts/vision-review-goldens.sh`, human-confirm SUSPECTs).

**Infrastructure & hygiene (34–40)**
34. Identify and leash the go.mod flip-flop writer (§e7, §g2).
35. Confirm/clean the `flake.nix` aarch64-darwin work (item 2) — it touches the fleet-standard systems list.
36. Check the visualtest `.fail/` artifacts from today's runs are gitignored (they are) and purge stale ones from disk.
37. Cold-start latency measurement on the demo (min-instances=0) — owner cost decision.
38. Demo e2e through the Firebase proxy in CI (live-URL lane post-deploy).
39. `nix run .#css` byte-stability re-check after any future class change — automate into item 3's lane.
40. Prune vnu ignore classes on the next nixpkgs html5validator bump (TODO #216 family — now 16 classes).

**Owner-decision queue (41–50)**
41. SSE bypass choice (see 13) — the only consumer-visible broken feature left.
42. Canonical redirect + header Demo link (see 14/15).
43. v1.19.5 release timing (see 18).
44. Standing authorization: may pipeline work commit promptly through the hook instead of waiting for daemon sweeps? (Today's mirror drift shipped through a sweep; §d2.)
45. May I chase the go.mod flip-flop into BuildFlow config (owner repo) or realign periodically?
46. The flake.nix aarch64-darwin removal: yours/intentional? (It changes which systems `nix flake check` enforces.)
47. GetNonce fallback: build the ADR or record a NO (see 11).
48. Coverage investment this quarter for the awesome-templ gate (TODO #28).
49. Errorpage demo presentation decision (TODO #236 ⫱).
50. Confirm the 07-08 report's three questions are superseded by 41–43 (close them there).

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **SSE bypass (TODO #329):** the canonical `/demo/**` can never stream (Firebase Hosting buffers every rewrite — verified live, platform-documented). Which fix do you want — (a) map a Cloud Run-direct domain (e.g. `demo.lars.software`, needs your DNS access) and point the Datastar SSE route at it, (b) migrate the demo to Firebase App Hosting (streaming-capable, new pipeline), or (c) keep Hosting for pages but point the EventSource at the raw `*.run.app` URL via a deploy-time env var (+ CORS + CSP `connect-src`)? (c) is fully in-repo; (a)/(b) need infra you own.
2. **The go.mod flip-flop:** BuildFlow warns the root `go ` line changed 20 times in 20 commits — some writer (a BuildFlow provider or the daemon) keeps re-flipping it, and today it re-broke the directive-equality guard and the workspace toolchain until I pinned everything at `1.26.0`. Do you want me to hunt the flipper in `larsartmann/buildflow`'s config and pin/skip it, or accept periodic re-alignment as daemon weather?
3. **Commit/push authorization + release:** (a) may I commit verified pipeline work promptly through the pre-commit hook instead of letting daemon sweeps land it unguarded (today's mirror drift shipped that way)? — this is standing authorization for THIS repo's pipeline work, not a general rule change; and (b) should I cut **v1.19.5** now that CI is green (the contrast/ViewTransitions changes are consumer-visible and `[Unreleased]` is warm), or batch with the next feature set?

---

*Report ends. Nothing in flight belongs to this session; the queued CI/Website runs for the 15:12Z daemon push and the dirty `flake.nix` are the open observations handed to the next session (§b1–b2, §f1–2).*
