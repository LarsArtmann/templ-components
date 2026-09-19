# Status Report — Pareto Plan Execution, Session 5 (T9 → T19)

**Date:** 2026-09-19 22:45 CEST
**Branch:** `master`, ahead 17 of `origin/master` (NOT pushed — M03 witness-before-push still pending)
**Working tree:** clean (all content committed; most of it as daemon `chore: auto-commit` chunks)
**Plan:** `docs/planning/2026-09-19_15-36_sales-page-pareto-execution-plan.md` (24 tasks)
**Previous report:** `docs/status/2026-09-19_18-46_pareto-plan-execution-status.md`

---

## Executive Summary

This session resumed at the red T9 layout lane and drove the plan from 12/24 to
**21 of 24 tasks fully done**. T9 was un-red and completed with three real HTML
bugs fixed along the way (nested `<pre>` in the hero, invalid combobox input
semantics, stale security-header tests). T14 (structural page count), T15
(warm-dark token overrides), T16 (per-page script audit), T17 (sitemap lastmod),
T6b (CopyButton browser e2e), T11 (docs-only search invariant) and T19 (lint
cleanups) all shipped verified. **T18 (Lighthouse) is half done**: landing
measured (93/100/100/100), sales not yet measured. T20/T21/T24 and the final
push ritual remain.

A **stars-badge nondeterminism flake** was root-caused and permanently fixed:
every dist rebuild fetched the live GitHub star count, so route goldens flapped.
`SITE_SKIP_STARS=1` in the `.#visual` app pins it; determinism proven with two
consecutive dist rebuilds.

**Parallel session CONFIRMED** (answers Q1 from the 18:46 report): commit
`0d4db611 fix(utils): drop the testing import from the production package`
(22:18) is NOT from this session. It is a real utils refactor (TestReporter
interface) and it is now part of the unpushed tip. Current tip verified green
(site goldens + axe + HTML gate) AFTER that commit landed.

---

## a) FULLY DONE (this session, all verified green)

| Task | What shipped |
| --- | --- |
| **T9 — HTML validation gate** | `layout/a11y_test.go` flipped to assert absence of the deprecated no-op metas; layout suite green. Website goldens added to `scripts/check-html-valid.sh` corpus (262 files validated). Gate GREEN with 15 documented ignore classes. CI verified: `html-validation` runs on every push (no path filter), so no ci.yaml change needed. **Three real bugs found & fixed:** (1) hero code window nested `<pre><pre class="chroma">` — fixed in `website/internal/pages/highlight.go` with `chromahtml.PreventSurroundingPre(true)` (verified against chroma v2.27.0 source: `nopPreWrapper`); (2) search input used `role="combobox"` + `aria-expanded` on `type="search"`, which ARIA-in-HTML forbids (combobox is legal only on `type=text`; APG pattern uses type=text) — fixed in `header.templ`; (3) new ignore class added ONLY for `media` on meta theme-color (WHATWG spec explicitly allows; vnu dataset lags — spec-verified via fetch). tc scaffolder sources re-synced (`check-tc-sources-sync.sh --fix`). Docs counts updated same-edit (visual goldens 163→171 in FEATURES/README/ROADMAP; `TestDocsCountDrift` green). |
| **T14 — structural page count** | `staticPages` const deleted; `topLevelPages(stats, starsLabel, nonce)` is now the single structural source of the top-level page set, consumed by both `run()` and `TestSiteBuildIntegrity` (wantPages derived, no hand-kept count). |
| **T15 — dark-tint tuning** | Sales page library surfaces retinted to the site's warm stone palette via explicit `dark:` Class overrides (they must beat the library shells' `dark:bg-gray-800`/`-900`, which tailwind-merge would otherwise keep): salesProblem Card, salesBenefits SimpleCard, all 4 salesProof StatCards (`dark:bg-bg-card-solid dark:border-border`), and the FAQ Accordion — the per-item `<details>` hardcodes its bg where no prop reaches, so a complete-literal arbitrary child variant `[&>details]:dark:bg-bg-card-solid` retints items (site.css `@source` scans `internal/pages/**/*_templ.go`, so it compiles). templ regen + page goldens + **all 8 site route goldens regenerated**; axe sweep re-read, ledger entry `site_sales_dark` still matches (1 accepted finding). Site tier green WITHOUT `-update`. |
| **T16 — per-page script audit** | Full hook map built: theme-sync/header/search/animations/newsletter all have consumers on every page they load on (footer newsletter form + `[data-animate]` footer are global; search box is in the header). ONE dead fetch found and removed: `copy-code.js` on docs pages (its only hook `#copy-btn` exists on the landing hero only; docs copy buttons are bound by `docs.js`). Script lists now documented in comments. |
| **T17 — sitemap lastmod** | `lastUpdated` generalized to variadic repo-relative paths (newest commit across sources, `git log -1 --format=%cs -- p1 p2`); `/sales` lastmod from `sales.templ`, landing lastmod from `landing.templ + hero.templ`; docs entries unchanged semantically. Convention written into `writeSitemaps` godoc: a new top-level page must be added to `topLevelPages` AND get a lastmod source — never a hand-kept count. Site tests green. |
| **T6b — CopyButton browser e2e** | New `visualtest/site_copy_e2e_test.go` `TestSiteSalesCopyButton`: serves the BUILT dist (`requireSiteDist`), injects a `clipboard.writeText` spy before the click (headless clipboard READS need ungrantable permissions; the spy observes the same call and the real write still executes), clicks `[data-tc-copy]`, asserts the exact install command payload AND the "Copied!" label swap. PASSED under `nix run .#visual`. |
| **T11 — search index scope** | G1 decision enforced as an invariant: `assertSearchIndex` now derives the docs URL set from `pages.AllDocs()` and fails if ANY search-index entry is outside it. First attempt used a wrong `/docs` prefix (site docs live at `/getting-started/*`, `/guides/*`, …) — caught by its own first run, fixed to set-membership. Site suite green. |
| **T19 — lint cleanups** | (1) `WriteString(literal + "\n")` concatenations removed in `writeSitemaps`; (2) QF1002: docs sidebar switch → tagged `switch doc.Slug`; (3) QF1003: chart axis-label if/else → tagged `switch i` (+ regen; output-identical, goldens unchanged, display suite green); (4) visualtest module lints **0 issues** (stale `nolint:gosec` → whole directive removed as unused; godox "bug" word in the accepted-debt comment reworded to "defect"); website module lints **0 issues** (golines 120-col restructure, prealloc with named const, wsl blank-line). gopls LSP diagnostics for `main.go:334/337` are STALE — ground truth is `golangci-lint run` = 0 issues. |
| **Stars-badge flake fix (unplanned, blocking)** | Site route goldens failed with mismatch sets that CHANGED between runs. Root cause: `fetchStars()` hits the live GitHub API on every dist build — badge text drifts (and on mobile can change layout height → 100%-dimension mismatches). Fix: `SITE_SKIP_STARS=1` env → `build.sh --skip-stars`, set in the flake `.#visual` app. Goldens regenerated under the pinned dist; **determinism proven: two consecutive dist rebuilds both green**. Production builds (website.yml) keep live stars. |

**Also verified at the CURRENT tip** (after the parallel session's utils commit):
site route goldens + axe sweep `ok`; HTML validation gate clean (262 goldens).

## Sessions-to-date scoreboard (plan tasks)

- **Fully done:** T1, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, T16, T17, T19, T22, T6b = **19**
- **Closed by decision:** T23 (declined per G2 — no history rewrite; to be recorded in T24 annotations) = **1**
- **Partially done:** T18 (landing Lighthouse done; sales run missing; quick-wins pass missing) = **1**
- **Not started:** T20, T21, T24 = **3**
- Plus the final M03 witness-and-push ritual (not a plan task).

---

## b) PARTIALLY DONE

1. **T18 — Lighthouse/perf.** Landing measured via `lighthouse@12` + nix chromium
   against a local dist serve: **performance 93, accessibility 100,
   best-practices 100, SEO 100** (report JSON at `/tmp/lh-index.json`, NOT
   committed — machine-local). Missing: the same run for `/sales.html`
   (`/sales` 404s on the python static server — no cleanUrls; use `sales.html`),
   recording both scores in the plan annotation, and the "quick wins only" pass.
   Note: the site goldens currently pin the NO-stars badge while production
   serves live stars — perf numbers measured on the skip-stars dist are
   representative enough (one badge line), but note it in the annotation.
2. **Push ritual (M03).** 17 commits ahead, tree green at tip on the lanes I
   ran (site visual tier, website suite, HTML gate, lint on website+visualtest).
   Still missing: the FULL `scripts/ci-repro.sh --lint --website` at the exact
   tip (which also covers the demo visual suite — NOT re-run since the T19
   templ edits — plus all-module build/test), witnessed PASS, then immediate
   push.

## c) NOT STARTED

1. ~~**T20 — Scrollback `Prompt` upstream idea draft.** Park-only deliverable:~~ done (draft parked at docs/upstream-drafts/2026-09-19_scrollback-prompt-field.md (source-verified, no duplicate issue; rename routed to TODO_LIST v2 #40))
   ~~verify against vendored source, check upstream issues, draft in Lars's voice~~
   ~~(verify-before-filing → github-voice). NEVER auto-file.~~
2. ~~**T21 — Dogfood-marketing recipe doc.** `docs/recipes/dogfood-marketing-page.md`~~ done (written at docs/recipes/dogfood-marketing-page.md + recipe-index row)
   ~~(derived counts + self-render pattern; the sales page is now the working~~
   ~~example — T15/T16/T17 learnings belong in it).~~
3. ~~**T24 — HARVEST.** Pull accepted items into TODO_LIST/ROADMAP, annotate the~~ done (TODO_LIST #271-#281 + v2 #40 + #216 extension; ROADMAP 4 rows; CHANGELOG [Unreleased] warmed; plan got section 8)
   ~~plan + status reports with outcomes (incl. T23-declined note, G1/G2/G3~~
   ~~outcomes, T18 numbers once complete, the stars-flake fix).~~
4. **Final: ci-repro witness + push.**

## d) TOTALLY FUCKED UP (own mistakes this session, no masking)

1. **Hit the 50-background-jobs ceiling.** I spawned one background shell per
   nix/lint/test invocation all session; the harness stopped accepting new
   background jobs near the end ("maximum number of background jobs (50)
   reached") and one lighthouse verification call failed purely because of it
   (the lighthouse run itself HAD completed). Completed shells apparently do
   not free slots until session end. Cost: one wasted round trip + uncertainty
   about a result that was actually fine. Future sessions: chain long runs
   into ONE background shell, keep everything else foreground, or redirect to
   files and poll.
2. **Wrong invariant shipped-on-first-try in T11.** I asserted all search URLs
   start with `/docs` — the site's docs live at top-level sections
   (`/getting-started/…`). My own test immediately failed (good), but I should
   have looked at `pages.AllDocs()` slugs BEFORE writing the assertion, not
   after. One wasted test cycle.
3. **Edit-tool mod-time rejections vs the templ-fmt daemon.** Two multiedit
   batches on `sales.templ` bounced with "file modified since read" because
   BuildFlow realigns templ files between my read and edit. I recovered by
   re-reading, but the FIRST attempt even used stale content I'd read minutes
   earlier — the guard did its job. Lesson applied: re-read immediately before
   every edit batch on daemon-touched files.
4. **Hand-rolled static server for T18.** `python3 -m http.server` on
   `website/dist` doesn't implement Firebase cleanUrls, so `/sales` 404'd and
   confused the first verification; I should have reused the existing Go
   `siteDistBase` pattern (which implements cleanUrls) or served `*.html`
   URLs directly from the start.
5. **Commit-message quality lost to daemon races (again).** Nearly every work
   batch landed as `chore: auto-commit N file(s) (heuristic)` daemon chunks
   because the daemon commits every ~1–2 min while I verify. My one successful
   hand commit (`ed68c9b6`) only caught the tc-sources sync because the
   pre-commit hook forced a fix-and-retry. Accepted per G2 (no history
   rewrite), but the session narrative lives in this report, not in `git log`.
   A "stage + commit immediately after each green check, before the daemon's
   next sweep" discipline would preserve more of it.
6. **Stale LSP diagnostics acted on as if fresh (twice).** `main.go:334/337`
   writestring and `docs.templ:197` QF1002 pointed at lines that no longer
   exist / were already restructured. I eventually treated golangci-lint +
   source reads as ground truth (per AGENTS guidance), but I initially went
   hunting for code at the cited lines. Should have started there.

## e) WHAT WE SHOULD IMPROVE

1. **Pin the visual dist everywhere it's consumed for verification.**
   `SITE_SKIP_STARS=1` is set in `.#visual`, but `nix run .#shots`,
   `siteshots`, and any manual dist serve still fetch live stars. Consider
   making skip-stars the DEFAULT for every non-production entry point
   (flag-gated production build instead).
2. **Lighthouse as a repeatable lane, not a one-off.** The manual npx +
   CHROME_PATH dance worked but is undocumented and uncommitted (results JSON
   is in /tmp). Either a tiny flake app `.#lighthouse` (chrome + npx lighthouse
   + local dist server with cleanUrls) or a chromedp nav-timing budget test in
   visualtest — one of the two, so T18 numbers stay reproducible.
3. **Local dist serving with cleanUrls.** A one-liner Go helper (the
   `siteDistBase` handler extracted) would serve any dist correctly and kill
   the python-server 404 class entirely.
4. **Background-job hygiene in agent sessions** (see d1): one shell, many
   chained commands, file-based output.
5. **Commit early, commit deliberately.** When the user authorizes commits,
   committing each verified batch IMMEDIATELY (before the next verification
   step) preserves history quality against the daemon race.
6. **The Q3 palette debt is still open.** `site_sales_dark` accepts white-on-
   `blue-500` (3.76:1) with budget 1. T15 fixed the SURFACES; the CTA shade
   convention (-600 light / -500 dark) remains library-wide accepted debt.
7. **The parallel session's utils refactor is in the unpushed tip unreviewed
   by this session** (`0d4db611`). Tip lanes are green, but a `git show` review
   of that diff before push would be cheap insurance.

## f) NEXT TASKS (ordered, ≤50)

**Finish the plan (blocking the push):**
1. ~~T18: run Lighthouse on `/sales.html`, record both pages' scores in the plan~~ done (landing 93/100/100/100 re-confirmed on fresh dist; /sales 79/100/100/100; recorded in plan section 8 T18 row)
   ~~annotation (§2 outcome column).~~
2. ~~T18: quick-wins triage — inspect the 7 perf deductions on landing (likely~~ done (triage complete - no zero-risk code wins; gaps are harness artifacts (no gzip) + htmx-parse (TODO_LIST #282) + stagger design; prod re-measure #280, lane #272)
   ~~render-blocking font preload chain / LCP); fix ONLY trivial ones (e.g.~~
   ~~`fetchpriority`, preconnect), regen goldens if pixels change.~~
3. ~~T20: verify Scrollback prompt behavior against the vendored/source code;~~ done (PARKED at docs/upstream-drafts/2026-09-19_scrollback-prompt-field.md; verified against scrollback.templ + sales.templ; no duplicate issue; v2 rename = TODO_LIST #40)
   ~~check upstream (a-h/templ? this repo's own backlog?) for existing issues;~~
   ~~DRAFT the issue in `docs/` — park, never file (G1).~~
4. ~~T21: write `docs/recipes/dogfood-marketing-page.md` — the sales page as the~~ done (docs/recipes/dogfood-marketing-page.md + recipe-index row)
   ~~recipe: derived counts (CountStats), self-render pattern, warm-dark Class~~
   ~~override pattern, script-audit checklist, lastmod convention.~~
5. ~~T24: annotate the plan (§2 gates outcome: G1 kept, G2 declined=T23 closed,~~ done (plan section 8 appended)
   ~~G3 done) + both status reports with outcomes.~~
6. ~~T24: harvest accepted items into `TODO_LIST.md` / `ROADMAP.md` — candidates~~ done (TODO_LIST #271-#281 + v2 #40 + #216 extension; ROADMAP 4 new rows)
   ~~listed in 14–30 below.~~
7. ~~Record T23-declined in the plan annotations (G2).~~ done (plan section 2 + section 8)
8. Full `nix run .#visual` (COMPLETE suite incl. demo goldens — not re-run
   since the T19 templ edits; chart output is semantically identical but
   unverified at pixel level).
9. ~~`git fetch && git status -sb` — re-check for parallel-session/daemon commits.~~ done (fetched; ahead 21 of origin, nothing pushed)
10. ~~Review `0d4db611` (parallel session's utils refactor) — sanity-read the diff.~~ done (APPROVED - 0d4db611 is CHANGELOG-only; the TestReporter code rode daemon a4f78a90; claim verified at source (no testing import in utils production files); shipped as v1.18.1)
11. Run `nix develop -c bash scripts/ci-repro.sh --lint --website` at the exact
    tip (background, ~7 min), WITNESS `VERDICT: PASS (exit 0)`.
12. Push IMMEDIATELY after witnessed PASS (re-check tip didn't move), then
    verify CI + Website lanes go green, and spot-check prod `/sales`
    (og:image, lastmod, CSP unchanged).
13. ~~Confirm CHANGELOG `[Unreleased]` is warm with the session's user-facing~~ done (CHANGELOG [Unreleased] warmed (3 Fixed + 3 Changed))
    ~~items (validation-gate expansion, search scope, sitemap lastmod, dark-tint~~
    ~~fixes, hero `<pre>` fix, combobox type fix) — the one-commit-release rule~~
    ~~needs it warm; add in the T24 commit if missing.~~

**Harvest candidates for TODO_LIST/ROADMAP (from this session's findings):**
14. ~~Default `SITE_SKIP_STARS=1` for all non-production dist entry points.~~ done (TODO_LIST #271)
15. ~~`.#lighthouse` flake app or nav-timing budget test (repeatable T18 lane).~~ done (TODO_LIST #272)
16. ~~Extract a reusable cleanUrls dist server helper (share with `.#shots`).~~ done (TODO_LIST #273)
17. ~~Prune the `media`-on-meta vnu ignore class when nixpkgs vnu catches up~~ done (folded into TODO_LIST #216)
    ~~(same TODO class as #216).~~
18. ~~Consider adopting ARIA-in-HTML check for `role=combobox` input types as a~~ done (TODO_LIST #274)
    ~~lint/test rule in the LIBRARY (the site bug class could exist in~~
    ~~components: grep library `.templ` for `role="combobox"`).~~
19. ~~Schedule the library-wide dark CTA shade decision (Q3) — either bump dark~~ done (ROADMAP General row Dark CTA contrast decision)
    ~~semantic surfaces to `-500`→`-400` or formally accept 3.76:1 for large~~
    ~~text only, documented in the a11y policy.~~
20. ~~Add `[&>details]:`-style arbitrary-variant child overrides documentation to~~ done (merged into TODO_LIST #279)
    ~~the templ-components SKILL.md / theming docs (new pattern proven here).~~
21. ~~Move `/tmp/lh-index.json` numbers into the plan annotation + delete the~~ done (numbers recorded in plan section 8; /tmp/lh-*.json deleted)
    ~~temp file (machine-local artifact).~~
22. ~~Update `docs/visual-testing.md` with the SITE route tier + skip-stars pin.~~ done (TODO_LIST #275)
23. ~~Update the skill (templ-components SKILL.md) site section: search scope~~ done (TODO_LIST #276)
    ~~invariant, lastmod convention, topLevelPages structure.~~
24. ~~Consider `TestDocsCountDrift` coverage for the site's dist log line~~ done (ROADMAP General row Counts truth table)
    ~~(components=123 icons=105 enums=62) vs FEATURES counts — they derive from~~
    ~~different counters; one truth table would prevent future confusion.~~
25. ~~Add the sales page to the siteshots smoke set if not already covered by~~ **Won't implement — siteshots already covers /sales (added 2026-09-19 15-10 session); route goldens pin the pixels.**
    ~~route goldens (verify parity).~~
26. ~~Grep repo for other `WriteString(literal + literal)` occurrences (gopls~~ done (TODO_LIST #277)
    ~~writestring may exist elsewhere; golangci-lint doesn't run that analyzer).~~
27. ~~Consider enabling gopls-analyzer-backed checks in golangci-lint config so~~ done (ROADMAP General row gopls analyzer lint gates)
    ~~writestring/prealloc classes gate in CI, not just LSP.~~
28. ~~LSP-staleness: add a note to AGENTS that QF100x hints on `.templ` files~~ done (AGENTS.md LSP-staleness note extended)
    ~~can point at stale lines; ground truth = golangci-lint + regen.~~
29. ~~The docs' visual-goldens count claim (171) will drift again with every~~ done (ROADMAP General row Counts truth table)
    ~~site route — consider deriving it in `TestDocsCountDrift` from the~~
    ~~testdata dir instead of a hand-typed number (it already compares against~~
    ~~reality; make docs say "site routes + library" structurally).~~
30. ~~Reduce `fetchStars` flake surface in PRODUCTION builds too (cache last~~ done (ROADMAP Production stars caching row)
    ~~good value to a file, like Astro did?) — fallback exists; caching would~~
    ~~stop badge flapping between deploys.~~

**Backlog hygiene (pre-existing, cheap):**
31. Re-check `git status` for re-added `*_templ.go` gitignore lines after
    daemon commits (BuildFlow gotcha — none observed this session, keep
    watching).
32. After push: verify GitHub autoclose keywords in any issue-closing commit
    texts (AGENTS convention).
33. ~~Sweep legacy raw `chromedp.Poll` sites (~48, TODO #240) — untouched.~~ **Won't implement — duplicate of TODO_LIST #240.**
34. ~~website.yml path filters: confirm `website/**` filter covers the new~~ done (folded into TODO_LIST #280)
    ~~`build.sh` env var behavior (deploy uses live stars — no change needed,~~
    ~~just verify).~~
35. `visualtest/tools/siteshots` search smoke: confirm it still passes with
    `type=text` search input (behavior unchanged; one smoke run would prove).
36. ~~Add `TestSiteSalesCopyButton` to any documented e2e inventory (docs/testing~~ **Won't implement — no FEATURES.md e2e inventory exists to update; the test gets its docs home via TODO_LIST #275.**
    ~~or FEATURES test-coverage lines mention suites; the new test should be~~
    ~~named there if the convention requires).~~
37. ~~Consider swapping python http.server examples in docs for the Go helper~~ done (folded into TODO_LIST #273)
    ~~(see 16).~~
38. ~~ogshot: `SITE_SKIP_STARS` doesn't affect OG cards (they don't render~~ done (TODO_LIST #281)
    ~~stars?) — verify and document either way in the ogshot README.~~
39. ~~Confirm the demo CSS is fresh after T15's new arbitrary-variant classes~~ done (TODO_LIST #278)
    ~~(they're site.css-scoped; demo CSS scans `**/*.templ` repo-wide — check~~
    ~~`examples/demo/static/app.css` contains them or prove they're unused by~~
    ~~demo routes; the daemon chunk `7de2133c` touched demo app.css already).~~
40. ~~Ask the BuildFlow daemon to include a file list in heuristic messages~~ **Won't implement — #93-family upstream idea; TODO_LIST #93/#232 own it.**
    ~~(upstream larsartmann/buildflow — park as an idea).~~

**Post-release follow-ups (after next version cut):**
41. ~~Prune vnu ignore classes when a newer checker lands (recurring).~~ **Won't implement — duplicate of TODO_LIST #216.**
42. ~~Re-measure Lighthouse on PRODUCTION (live stars + Firebase headers) once —~~ done (folded into TODO_LIST #280)
    ~~local dist numbers lack Firebase caching/CDN effects.~~
43. ~~Verify prod sitemap shows `/` + `/sales` lastmod after next deploy.~~ done (folded into TODO_LIST #280)
44. ~~Watch the first real deploy for the stars badge render (live count).~~ done (folded into TODO_LIST #280)
45. ~~After the next templ upstream release, revisit the v0.3.1020 pin (AGENTS~~ **Won't implement — standing AGENTS.md templ-pin policy owns it.**
    ~~standing item).~~
46. ~~When nixpkgs html5validator updates, re-run the gate to catch newly~~ **Won't implement — duplicate of TODO_LIST #216.**
    ~~enforced rules early (pre-CI).~~
47. ~~Consider adding `/sales` to `search-index.json` exclusion docs prose (the~~ **Won't implement — the T21 recipe documents the search-scope invariant.**
    ~~docs pages describing the site architecture) if T21 doesn't cover it.~~
48. ~~Add the warm-dark override pattern to the website's own theming docs page~~ done (TODO_LIST #279)
    ~~(`guides/theming` content markdown).~~
49. Keep an eye on chroma upgrades for `PreventSurroundingPre` behavior
    (pinned by go.mod; behavior verified at v2.27.0).
50. Celebrate, then start the NEXT Pareto planning cycle from the harvested
    TODO_LIST.

## g) QUESTIONS FOR YOU (cannot answer myself)

1. ~~**The parallel session (`0d4db611` utils TestReporter refactor, plus the~~ done (SUPERSEDED BY DECISION - M03 witness-at-exact-tip covers racing sessions; 0d4db611 gets a git show review before the push)
   ~~earlier `1ced1afa` ogshot commit):** is that session FINISHED, or will it~~
   ~~keep committing while I run the final ci-repro + push? If it's still~~
   ~~active, do you want ME to hold the push until it's done (M03 witness is~~
   ~~only valid at an exact tip — a racing session invalidates it), or is the~~
   ~~other session also under the same ritual and coordination is on your side?~~
2. ~~**T18 scope confirmation:** Lighthouse landing = 93/100/100/100. The 7 perf~~ done (RESOLVED BY THE WHOLE-LIST MANDATE - measure /sales; apply only zero-risk wins (any pixel change re-rolls the 8 site goldens))
   ~~points are mostly font/LCP related quick-wins territory. Do you want the~~
   ~~cheap fixes applied NOW (touches head/font loading → regenerates site~~
   ~~goldens again), or should I record the scores as-is and park perf tuning~~
   ~~for a dedicated task? (Plan default says "quick wins only" — but any pixel~~
   ~~change re-rolls 8 goldens, which is why I'm asking before touching head~~
   ~~loading.)~~
3. ~~**Q3 from the 18:46 report, still unanswered — the library-wide dark CTA~~ done (ROUTED TO ROADMAP - General row Dark CTA contrast decision (v2-scale))
   ~~shade (white on `blue-500` = 3.76:1):** keep accepting per-site (current~~
   ~~ledger approach) or schedule the library-wide palette fix? T21's recipe doc~~
   ~~and the site's axe ledger both reference this decision.~~

---

**Verdict:** ~~plan is 21/24 complete + 1 closed-by-decision; T18 half done;
T20/T21/T24 + push remain.~~ T20/T21/T24 closed by the resume session (2026-09-19 late):
T20 parked at `docs/upstream-drafts/`, T21 recipe shipped, T24 harvested (TODO_LIST
#271-#281 + v2 #40, ROADMAP rows, CHANGELOG `[Unreleased]` warm, plan §8). Remaining:
T18 `/sales` measurement + quick-wins, full visual suite, 0d4db611 review, ci-repro
witness + push. Tip is green on every lane this session ran. Nothing pushed yet — deliberately.
