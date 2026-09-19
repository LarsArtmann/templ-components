# Status Report — Pareto Plan Execution, Sales-Page Hardening (Session 4)

**Date:** 2026-09-19, 18:46 CEST
**Session:** Execution of `docs/planning/2026-09-19_15-36_sales-page-pareto-execution-plan.md` (24 tasks), started from tip `7ee07b3b` (plan committed + pushed in session 3).
**Headline:** 12 of 24 plan tasks fully done and verified (the entire 1%/4%/20% Pareto core: prod verification, claim derivation, no-framework assertion, credibility links, nav, OG card, anchors/FAQ/README, dense variant, site route goldens, site a11y sweep). T9 is **half-done and the `layout` test lane is RED right now** — one test file still asserts the removed security metas. T9 must be finished first thing next session. No push happened this session: M03 witness-before-push is still pending; local master is ahead 33 (mostly daemon heuristic chunks).

---

## 0. Repo state at report time

| Item | State |
|---|---|
| Branch | `master`, **ahead 33** of `origin/master`, nothing pushed this session |
| Working tree | `M layout/base_templ.go` (regenerated), `M scripts/check-html-valid.sh` (corpus extended) |
| Test lanes | website module **green** · visualtest (demo + site, full `nix run .#visual`) **green** (105.9s) · **layout RED** — `TestSecurityHeaders/security_headers_rendered_when_enabled` in `layout/a11y_test.go:19` still asserts the two removed metas |
| Commits with detailed messages this session | `eaf2e8c1` (OG card), `16f9015d` (SplitSeq lint fix), `8da6539c` (T7+T8 visual net) |
| Daemon | committed ~28 heuristic chunks, including complete snapshots of my in-flight work (`0161ab09` swallowed the T4–T6 delta) |
| Parallel session suspected | `1ced1afa fix(ogshot): clear visualtest lint lane (gosec, nolintlint, golines)` at 18:10:26 is NOT my commit (I never authored it; message style is mine-like but I did not run ci-repro at that point). Also the pre-session daemon chunks (forms/calendar, tags_input, icons, utils/svg, CHANGELOG) were not mine. See question Q1. |

---

## a) FULLY DONE (implemented, tested, verified)

| Plan task | What shipped | Proof |
|---|---|---|
| **T1** — record gates | §2 annotated: G1 = README+nav, no social; G2 = **DECLINED** (no history rewrite; force-push needs explicit approval); G3 = **FULL machinery** | plan file §2 "Gate decisions" |
| **T3** — deployed verification | Prod `/sales` 200 + correct title; **CSP header byte-identical to committed firebase.json** (all 5 hashes incl. CopyButton); sitemap lists `/sales`, sitemap-index + robots sane; **prod HTML == local dist byte-for-byte modulo per-build nonce**; live vnu findings triaged into T9 | fetch/download + diff, agentic_fetch header report |
| **T4** — derive claims | `build.Stats` gained `GoVersion` (root go.mod `go` directive parsed, patch-trimmed, fail-soft) + `LibraryVersion` (`utils.Version`); hero badge + FAQ "at v1" render from them; CountStats fixture test extended (go directive in fixture, both assertions) | website suite green; output byte-identical today, self-maintaining on version bumps |
| **T5** — no-framework assertion | `assertNoFrameworkScripts` over **every** site page (external script srcs + framework identifiers **inside inline script bodies only** — prose comparisons naming Alpine.js stay legal) + `TestFindFrameworkViolations` negative control | caught nothing on real pages; detector proven on 4 frame classes |
| **T6** — credibility links + stars | Stars badge (shared `StarsLabel`, fail-soft) on sales hero; MIT badge → LICENSE href; FAQ production-ready answer → derived version + link to `/guides/invariants` (documents the three regression layers); cost answer → LICENSE link; new `faqAnswerRich`/`faqLink`/`faqProductionReady`/`faqCost` sub-templates | goldens updated same commit; links pass `build.CheckLinks` |
| **T10** — header nav | "Why" → `/sales` added to site header (all pages, desktop + shared mobile panel), positioned first (Why → Docs → GitHub) | dist grep + goldens |
| **T12** — OG image | `visualtest/tools/ogshot` + flake app `nix run .#ogshot`; card is **self-contained CSS** (not Tailwind classes — compiled app.css only contains utilities seen in `.templ` files; first render attempt proved it); no derived counts (static PNG cannot track them); `website/public/og/sales.png` 1200×630 rendered + visually verified; `SalesMeta.OGImage` wired; `og:image` in dist confirmed | `nix run .#ogshot`, dist grep, image inspection |
| **T13** — discoverability | Section anchors `#problem` `#proof` `#faq`; first FAQ `Open: true`; README link row gained "Why templ-components"; redundant `role="banner"`/`role="contentinfo"` removed (flagged by live vnu) | dist grep; anchors pass CheckLinks |
| **T22** — section dense variant | `sectionOpts{id, dense}` + helpers; opt-in `py-12` rhythm; zero-value keeps landing byte-identical | templ generate clean; landing golden unchanged by the refactor itself |
| **T7** — site route goldens | 8 goldens: site-{sales,index} × light/dark × desktop(1440×900)/mobile(390×844), full-page, **scroll-through first** (site `[data-animate]` reveals leave below-fold sections at opacity 0 otherwise) + `waitAnimationsSettled`; `.#visual` flake app now **builds website/dist first** (CI fails on skips, so missing dist must hard-fail) | `nix run .#visual` green; goldens committed in `8da6539c`; capture visually inspected (all sections revealed, badges/nav/FAQ visible) |
| **T8** — site a11y sweep | `TestAxeSweepSiteRoutes` (default-fail gate, shared ledger, `site_*` keys) + `TestSiteTouchTargetAudit` + `TestSiteZoomReflowAudit`; audit scrolls + **waits for finite animations** (scroll-through fires reveals; axe sampling mid-animation produced ~38%-opacity blended false positives) | suite green after fixes; 4 real issues found and FIXED forward (below); 1 documented ledger entry |
| T8 findings — **fixed forward** | (1) `--color-text-muted` kept light value in dark: **1.93:1** on sales problem cards → stone-400 in site.css dark block. (2) FAQ links relied on color alone (WCAG 1.4.1) + failed dark contrast → persistent underline + `dark:text-accent-light`. (3) Newsletter input `min-width:auto` → 15px page overflow at 320px → `min-w-0`. (4) Quickstart code tokens unbreakable → 32px overflow at 320px → `wrap-anywhere` + `min-w-0` on step cards. | reflow/contrast probes before+after; audits green |
| Misc | BuildFlow lint autofix on my `goDisplayVersion` committed as `16f9015d` with real message; ogshot lint lane cleared (see Q1 — possibly parallel session) | — |

Verification runs this session: website `go test ./...` green (repeatedly), `nix run .#visual` full suite green (105.9s, demo + site), CountStats/goldens/negative-controls all exercised. **Not yet run:** ci-repro witness (push gate) — nothing was pushed, correctly.

## b) PARTIALLY DONE

| Task | Done | Missing |
|---|---|---|
| **T9** — HTML validation gate | Root cause fixed: `layout.Base.SecurityHeaders` emitted two **invalid, browser-no-op metas** (`http-equiv="X-Content-Type-Options"` is header-only per spec; `Referrer-Policy` is not a registered http-equiv). Field is now a documented deprecated no-op (v1-compatible, no API break); `base.templ` + `bdd_test.go` + `integration_test.go` updated; `check-html-valid.sh` corpus extended with `website/internal/pages` | **(1) `layout/a11y_test.go:19` still asserts the metas → layout lane RED.** (2) Library goldens not regenerated (layout package + every Base-using golden). (3) Website goldens not regenerated after meta removal. (4) `check-html-valid.sh` not yet executed locally. (5) Expected vnu-staleness ignores not yet added/verified (`media` on meta theme-color; possibly `type="search"` mystery from the live validator). (6) `ci.yaml` html-validation job path filters not yet checked to include `website/**` + `layout/**` |
| **T19** — lint cleanups | ogshot lint lane cleared (`1ced1afa`, provenance uncertain) | gopls writestring (main.go:334/337), QF1002 (docs.templ:197), QF1003 (chart_shared.templ:59), full visualtest lint run — untouched |
| **T23** — squash | Decision recorded (G2 = DECLINED, no history rewrite) | nothing to do by decision — task closed as declined |
| **T24** — harvest | Plan §2 annotated; this report | TODO_LIST/ROADMAP harvest, status-report annotation, final plan outcome annotations per task |

## c) NOT STARTED

T11 (search index scope: per G1 decision, docs-search stays clean → likely a short "declined, keep invariant" wiring + assertion update), T6b (CopyButton browser e2e on built site — harness now exists), T14 (structural page count — kill `staticPages` hand constant), T15 (dark-tint tuning — goldens now make it safe), T16 (per-page script audit), T17 (sitemap lastmod for /sales + data.go content convention), T18 (Lighthouse/perf on / and /sales), T20 (upstream Scrollback Prompt idea — draft only, never auto-file), T21 (dogfood-marketing recipe doc).

## d) TOTALLY FUCKED UP (own mistakes, honest list)

1. **Left the `layout` lane RED.** I edited `base.templ` + 2 of 3 test files, then stopped mid-T9 for this report. Violates the plan's prime directive ("every task ends green"). Fix is ~10 minutes (flip `a11y_test.go:19` assertions to AssertNotContains, regen goldens, run layout suite).
2. **Lost my best commit message to the daemon.** The T4–T6 delta was fully committed by the daemon as `0161ab09 "chore: auto-commit 6 changed file(s)"` while my detailed `git commit` failed on the ref race ("cannot lock ref 'HEAD'"). The work is safe; the *why* is only in this report and the plan.
3. **sed silently no-op'd on a templ-fmt-realigned line** (`OGImage:     SiteURL...` had realigned whitespace) → built, tested, and nearly committed the OG wiring against the stale `home.png` until a dist grep caught it. Lesson: after any BuildFlow run, re-read before editing; prefer exact-match `edit` over sed in templ-touched files.
4. **Two wasteful iteration loops:** (a) T5's first detector scanned whole-page prose and flagged the landing's legitimate "Alpine.js" comparison cells, and my negative-control frame (`ReactDOM`) was un-matchable by design — rewritten to script-body-only + catchable frames; (b) the site axe sweep failed twice on animation-blend false positives before I root-caused it to scroll-reveal mid-flight opacity — the `[data-animate]` gotcha was *already documented* for goldens; I should have predicted it for axe on day one.
5. **`serveDist` returned a scheme-less URL** (`127.0.0.1:PORT`) → all 8 golden captures failed with "invalid URL" — cost a full nix run round-trip.
6. **Two throwaway debug tests hit avoidable friction:** JS try/catch syntax error in a chromedp Evaluate, and `t.Logf` output invisible without `-v` (burned 3 background runs). Also wrote a `runtime.Goexit()` placeholder into `siteDistBase` on first draft and referenced 4 helpers that didn't exist — compile-first discipline slipped while drafting long files.
7. **Batching invited the daemon.** Working ~25 minutes without committing on the T4–T6 chunk guaranteed a daemon sweep; committing smaller and sooner would have kept authorship history clean.

## e) WHAT WE SHOULD IMPROVE

1. **Commit per task, immediately** (repo rule already exists for explicit-commit work): the daemon wins every race longer than ~2 minutes.
2. **Green-checkpoint discipline:** never edit a test's subject across multiple files without running that package's suite before touching the next file — the a11y_test.go miss would have been caught in 20 seconds.
3. **Reuse documented gotchas proactively:** the scroll-reveal/animation-settle lesson existed; both new harnesses (goldens AND axe) needed it — design audits against the full known-gotcha list, not just the one in front of you.
4. **Templ formatter quirks deserve a memory entry:** composite literals (`templ.Attributes{...}`, `map[string]any{...}`) inside templ component bodies break `templ generate`'s formatter; plain Go helpers + spread `{ expr... }` (three dots, `templ.Attributer`) are the working pattern.
5. **Parallel-session coordination:** if another session is active (Q1), agree on file ownership or commit-message conventions to avoid mutual daemon-style churn.
6. **Push cadence:** this session produced a large unpushed delta; on a normal M03 cadence, interim witnessed pushes would cap the risk of a bad tip blocking everything.

## f) NEXT 50 (ordered; 1–6 are the immediate finish line)

1. **Finish T9:** flip `layout/a11y_test.go` security-headers assertions to AssertNotContains; run layout suite.
2. Regen ALL affected goldens (layout pkg + Base-using goldens + website pages after meta removal); update in same commit.
3. Run `nix shell nixpkgs#html5validator -c scripts/check-html-valid.sh`; triage new findings; add only documented staleness ignores (`media` on meta; check `type="search"`).
4. Verify `ci.yaml` html-validation + website job path filters cover `website/**` and `layout/**` so the gate fires on both.
5. Run `nix develop -c bash scripts/ci-repro.sh --lint --website` at the exact tip; witness `VERDICT: PASS`.
6. **Push immediately** after witnessing (33 commits waiting); confirm the website lane deploys and re-verify prod `/sales` CSP + og:image.
7. T6b: CopyButton click-through e2e on the built site (siteshots-smoke pattern; the harness now exists).
8. T14: export `buildSitePages(...)`, derive `TestSiteBuildIntegrity` wantPages, delete `staticPages`.
9. T17: `/sales` lastmod from git (`lastUpdated` non-docs variant) + write the data.go vs per-page content convention.
10. T16: per-page script audit — grep which routes load which siteScripts (sales legitimately uses newsletter.js; check docs pages don't).
11. T15: dark-tint tuning pass on sales Card/Accordion/StatCard surfaces (site's warm dark palette); goldens now protect it.
12. T18: Lighthouse perf/a11y/SEO on `/` and `/sales`; record in plan annotation; apply only cheap wins (font-display, image dims).
13. T11: wire the G1 decision — keep search index docs-only, assert the invariant in `assertSearchIndex` with a comment naming the decision.
14. T21: `docs/recipes/dogfood-marketing-page.md` (derived counts + self-render + CSP-hash + ogshot workflow).
15. T20: draft the Scrollback `Prompt` idea doc in Lars's voice (verify vendored source first; park for review — never auto-file).
16. T19 remainder: fix gopls writestring (main.go:334/337), QF1002 (docs.templ:197), QF1003 (chart_shared.templ:59 + regen).
17. T19 remainder: full `visualtest` module `golangci-lint run ./...` and clear findings.
18. T24 harvest: accepted items → TODO_LIST (bounded, owned); deferred → ROADMAP raw ideas.
19. T24: annotate the plan §3 with per-task outcomes (annotate, never rewrite) + annotate the 15:10 status report as superseded-by plan outcomes.
20. Answer Q1–Q3 below, then update the plan's gate table with any decision changes.
21. Confirm the parallel-session question (Q1) — if real, coordinate before next big batch.
22. Sweep the 33-commit local backlog's messages: accept heuristic history per G2, but ensure CHANGELOG `[Unreleased]` describes the user-visible set (derived claims, no-framework gate, OG card, site a11y net, nav/README discoverability).
23. Verify `utils.TestDocsCountDrift` and `TestCompiledCSSInventory` still green after site changes (rule: same-edit counts; counts didn't change, but confirm).
24. Re-check `website/firebase.json` headers after the SecurityHeaders meta removal: confirm X-Content-Type-Options + Referrer-Policy headers exist for ALL routes (they did for `/sales` — verify sitewide config shape).
25. Decide whether `DefaultPageProps()` should stop defaulting `SecurityHeaders: true` now that it's a no-op (behavior-neutral, docs-honest) — or defer to v2.
26. Consider deprecating/removing `SecurityHeaders` entirely in the planned v2 migration doc.
27. Add the scroll-reveal + animation-settle gotcha to AGENTS.md's visual-testing section (it bit twice).
28. Add the templ-formatter composite-literal gotcha to AGENTS.md (templ section).
29. Document the `.#ogshot` app + OG-card regeneration in `docs/` (website docs or recipes) so it's discoverable outside the tool's doc comment.
30. Extend the site golden set to one representative docs page (code-block + sidebar shapes) — cheap, catches docs-shell regressions.
31. Consider siteshots `-route` flag parity for ogshot (arbitrary OG cards for docs pages reuse the same tool).
32. Sitemap: `/` and `/sales` are the only URLs without lastmod (docs pages have it) — T17 covers sales; decide whether `/` should get one too.
33. Zoom-reflow: re-run at 320px after T15 tint changes (visual changes can alter break points).
34. Touch-target: log-only 44px AAA gap on site header controls — decide if the site should meet AAA like the demo aspirations.
35. Live-vnu vs pinned-vnu drift: after T9's ignores settle, re-run the live `html5.validator.nu` check on prod and reconcile (the `type="search"` live finding is unexplained; pinned checker may not flag it — document either way).
36. Wire `docs-testing/a11y-gate-policy.md` (if present) to mention the site tier alongside the demo tier.
37. Consider naming the site tier in `visualtest/doc.go` so future sessions find both tiers.
38. OG image: decide if docs pages' OG set needs regeneration in the same style (they predate this card design).
39. OG image: add a CI/integrity check that every page's og:image target exists in dist (link-checker covers `href/src` but og meta is absolute URL — verify CheckLinks actually skips it; add explicit check if so).
40. FAQ first-item-open changed the visual default — confirm the siteshots light/dark captures still read well (done for sales; eyeball index untouched).
41. `faqLink` templ helper is currently unused after inlining (check) — either use it in both inline anchors or delete it (dead-code hygiene).
42. Sales page hero Scrollback: consider `wrap-anywhere` parity check (`.tc-log` already has `overflow-wrap: anywhere` — verify it covers the long go-get token at 320px on the SITE dist, not just demo CSS).
43. Newsletter form: `target="popupwindow"` without a `popupwindow` JS opener — check buttondown docs; may be a leftover Astro pattern (T16-adjacent).
44. Search smoke (siteshots) runs against dist — confirm it still passes after header nav change (it did in this session's runs; keep in the loop after future header edits).
45. Add `site-sales` to the axe sweep's **light** ledger only if future real findings appear — currently clean; do not pre-accept.
46. Version claims: when Go bumps to 1.27, verify the sales badge auto-updates (the derivation exists; add a one-line note to the release checklist).
47. Same for library version: the FAQ "at v1" flips to "at v2" automatically at the major bump — verify then.
48. Consider extracting the repeated "scroll + settle + pin" chromedp sequence into one shared site-audit helper (three call sites now).
49. Retire the session-2/3 "awaiting gate answers" language in older status docs (annotate as resolved-by-execution).
50. Final plan close-out: mark the plan file's mermaid `done` node achieved once T24 lands, and cut the next release with the warmed `[Unreleased]`.

## g) QUESTIONS I CANNOT ANSWER MYSELF

**Q1 — Is another agent/session working in this repo in parallel?** Commit `1ced1afa fix(ogshot): clear visualtest lint lane` (18:10:26) is not mine but carries a deliberate message, and the pre-session daemon chunks touched files (forms/calendar, tags_input, icons, utils/svg) I never went near. If a parallel session exists, I should coordinate (file ownership, push cadence) instead of assuming the daemon is the only other actor.

**Q2 — G1 "no social": final, or should I prepare (not post) launch artifacts?** I decided G1 = README + in-site nav, no social posting. If you want a prepared launch post (X/LinkedIn/HN/Reddit draft + OG image variants) as a parked artifact for manual posting, say so and I'll add it as a bounded task — otherwise the current decision stands and /sales distributes via site + README only.

**Q3 — White-on-blue-500 (3.76:1) dark CTA: keep accepting, or schedule the palette fix?** The site ledger entry (and the demo's identical entries) all document the library-wide -600(light)/-500(dark) shade convention as deliberate debt. Fixing it means re-shading semantic colors library-wide (a visual release, likely v2). Keep accepting per-entry, or put "dark palette contrast release" on the ROADMAP with an owner?

---

**Standing contract honored:** T2/T3 done; execution ran G3-full per recorded defaults; nothing pushed (M03 pending); annotate-don't-rewrite observed on the plan; the one red lane (layout, T9 step 1) is flagged at the top of section (d)/(f) with its exact fix.
