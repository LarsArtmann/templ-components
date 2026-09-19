# Status — Pareto plan COMPLETE (T18–T24 + full visual suite + witnessed ci-repro); ONLY THE PUSH REMAINS

**Date:** 2026-09-20 00:16 CEST · **Branch:** `master`, **ahead 25, nothing pushed** (deliberate — M03) · **Tip:** `bdb53316`
**Predecessor:** `docs/status/2026-09-19_22-45_pareto-session5-t9-t19-completion-status.md` (T9–T19 + the 50-task list this session harvested)
**Plan:** `docs/planning/2026-09-19_15-36_sales-page-pareto-execution-plan.md` — now 24/24 resolved (T23 closed-by-decision; §8 status table appended).

## a) TL;DR

The sales-page Pareto plan is DONE except the push. This resume session executed T18 (Lighthouse measured + triaged on BOTH pages), T20 (Scrollback `Prefix` rename draft — parked), T21 (dogfood recipe doc), T24 (full harvest: TODO_LIST #271–#282 + v2 item #40 + #216 extension, 4 ROADMAP additions + one structure repair, CHANGELOG `[Unreleased]` warmed, AGENTS LSP note extended, plan §8, inline annotation of the 22:45 report — 43 numbered items resolved). Reviewed the parallel session's `0d4db611` (APPROVED — it is CHANGELOG-only; the code rode daemon `a4f78a90`; claim source-verified). Ran the FULL visual suite (189s, everything green — first pixel verification since the T19 templ edits). Then `scripts/ci-repro.sh --lint --website` at exact tip `bdb53316`: **VERDICT: PASS (exit 0), witnessed 2026-09-20 00:15:35 CEST**. Tip has not moved since (verified 00:16).

## b) WHAT I DID THIS SESSION (resume at "WAIT FOR INSTRUCTIONS" → whole-list mandate)

| Task | Outcome |
|---|---|
| Todo sync | Synced stale list first (T9–T19 completed; T18 in_progress), then worked it top-down |
| T20 | Parked draft `docs/upstream-drafts/2026-09-19_scrollback-prompt-field.md`: claim verified at source (`display/scrollback.templ:22-30` godoc "Timestamp is preformatted" vs `sales.templ:91-94` passing `$`/`→`), no duplicate issue upstream (only open issue #18, unrelated), own backlog checked (nearest: the 2026-08-22 Trace/Timeline idea), drafted in voice-profile §8 register. **Never filed.** v2 rename routed to TODO_LIST Deferred #40 |
| T21 | `docs/recipes/dogfood-marketing-page.md` + recipe-index row: derived counts (`build.CountStats`), self-render-as-proof, warm-dark `dark:`-pair override pattern (incl. `[&>details]:` complete-literal rule), script-audit checklist, CSP-hash + lastmod + search-scope release checklist. Every cited fact re-verified in source before writing (`build.go:140`, `data.go:159`, `csp.go:133`, sales.templ lines) |
| T24 (harvest) | TODO_LIST: #271 skip-stars default, #272 lighthouse lane, #273 cleanUrls helper, #274 combobox-role audit, #275 visual-testing docs, #276 SKILL site section, #277 writestring sweep, #278 demo CSS freshness, #279 theming-docs pattern, #280 post-deploy spot-checks, #281 ogshot×skip-stars, #282 htmx gating (added after T18 triage), v2 #40 Scrollback rename, #216 extended (14→15 vnu classes). Header ID counter kept consistent (next free: 283) |
| T24 (ROADMAP) | 4 additions: General rows "Dark CTA contrast decision (v2-scale)" (Q3 routed here), "gopls analyzer lint gates", "Counts truth table"; Website row "Production stars caching". ALSO repaired pre-existing structure: the 09-18 Errorpage-remainder section had been inserted between the "Website & docs ideas" heading and its table — moved below it with a move-note |
| T24 (annotations) | CHANGELOG `[Unreleased]` warmed (3 Fixed: vnu-gate expansion + 3 real defects, sales dark surfaces, search scope; 3 Changed: sitemap lastmod, docs script trim, deterministic site goldens). AGENTS: LSP-staleness note extended (verified-twice QF100x/gopls stale-line lesson). Plan §8 appended (per-task outcome table + gate outcomes + T18 numbers). 22:45 report: §c 1–3, §f 37 items, §g 1–3 struck inline via the docs-health `annotate-prose.py` (dry-run first), verdict line corrected |
| T18 | Both pages measured on the FRESH dist (chromium 153 headless, lighthouse@12, local server): **landing 93/100/100/100** (reproduces the prior session's number), **/sales 79/100/100/100**. Triage: zero weighted a11y/bp/SEO audits fail anywhere; the perf gap is (a) harness artifact (`uses-text-compression` — python http.server sends no gzip; Firebase serves brotli in prod), (b) /sales TBT 590ms (wt 30) dominated by inline self-hosted htmx on a page with **zero `hx-` attributes** — library-contract decision → TODO #282, (c) Speed-Index 4.1s partly from the Scrollback stagger (design intent). `font-display:swap` already present; no zero-risk code wins — nothing patched for metrics' sake. Numbers recorded in plan §8; /tmp artifacts deleted |
| 0d4db611 review | APPROVED: the "hand commit" is CHANGELOG-only; the actual `utils/test_helpers.go` TestReporter change rode daemon `a4f78a90`; verified at source (`no testing import in utils production files`; interface in place); shipped as v1.18.1 (`c85235cb`) |
| Full visual suite | `nix run .#visual` → **ok 189.205s**: all library + demo + site-route goldens, axe sweep, touch/zoom audits, wire/kanban/forms e2e, new `TestSiteSalesCopyButton` — first full-suite pixel proof since the T19 templ edits |
| M03 witness | `nix develop -c bash scripts/ci-repro.sh --lint --website` at `bdb53316` → all module lanes 0 issues, website lane green + tidy-check green, **VERDICT: PASS (exit 0) at 00:15:35 CEST**, tip unchanged at 00:16 |

## c) WHAT I FORGOT / WHAT COULD BE BETTER

1. **The push — the only remaining plan item.** Witness is timestamped 00:15:35; if the daemon commits before you say go, the ritual demands a re-run at the new tip (~7 min). Say "push" and I re-check tip, re-witness if it moved, then push immediately.
2. **`check-draft.py` never ran on the parked T20 draft** (bash was ceiling-blocked during drafting; forgot to run it once bash returned). One command before ever filing it: `--kind body-issue --ai-drafted`.
3. **The 18:46 report's open questions were never annotated** (its §g differs from the 22:45 one). Deliberately left — superseded by the plan §8 table; say the word if you want that report annotated too.
4. **Lighthouse lane still manual** (TODO #272) — numbers die with /tmp again; nothing committed makes them reproducible.
5. **`annotate-prose.py` multi-line continuation handling produced `~~`-wrapped fragments on wrapped lines** — correct per the tool, but the 22:45 report's §f reads choppy in raw markdown. Rendering is fine; just noting.

## d) TOTALLY FUCKED UP (own mistakes this session, no masking)

1. **ROADMAP edit fiasco — 3 failed rounds + transient corruption.** I stacked multiedits on one region with old_strings built from assumed padding (divider-row dash counts) instead of exact bytes; then misread "Applied 1 of 2" WITHOUT determining which edit had applied — I assumed edit 1 (the dedup) when it was edit 2 (the insert), leaving a duplicated heading while I believed the opposite. One intermediate state even deleted the Errorpage bullets (recovered — they were in context). Fixed properly with `cat -A` byte inspection + a single targeted edit + `git diff` verification. Lesson applied: after ANY partial multiedit result, re-read the region BEFORE the next attempt; never construct table-row old_strings from eyeballed padding.
2. **Fabricated a crisis from a wrong-surface grep.** `grep chroma dist/assets/app.css` → 0 → I announced a "shipped regression: syntax highlighting gone". The stylesheet lives at `dist/assets/css/chroma.css` (separate generator output, `build.ChromaCSS()` at `main.go:163`). My own hasty check invented the bug — exactly the verify-before-claiming failure I police in others. Second surface check (`ls assets/css/`) killed it in one minute.
3. **`curl` in the chained Lighthouse script.** The bash tool BANS curl; my background chain died at the first curl with exit 1, wasting a full dist build + serve cycle. The banned list is in the tool documentation I read this session.
4. **Trusted a Lighthouse run served by a port-squatter.** A leftover `http.server 8901` from a PRIOR session survived in the background; my new server failed to bind (`Address already in use`), Lighthouse measured against the STALE dist — `/sales` 404'd (zeros) and landing showed a bogus 85. I had the 404 traceback in the same output as the numbers and still started triaging. Re-ran with explicit HTTP-200 gates before spending any Lighthouse run + killed the squatter. Lesson: gate measurements on a verified-fresh server, every time, no exceptions.
5. **Off-by-one on the TODO_LIST ID counter** (wrote "next free 283" for 11 added items → corrected to 282 → then added #282 → 283 is correct again). Cheap catch, but it happened because I counted rows instead of trusting the last-used ID.
6. **Daemon raced my deliberate commit — again.** I staged a real commit message for the T24 annotations; `bdb53316` landed the content before my `git add` executed. G2 accepted the noise; the e5 "commit immediately" discipline keeps losing to a 60–120s daemon cadence. Only a pre-commit-daemon lock (TODO #232 family) fixes this class.

## e) WHAT WE SHOULD IMPROVE

1. **A "port + server hygiene" helper for measurement lanes.** Squatters from prior sessions silently poison Lighthouse/golden measurements (d4). One `scripts/serve-dist.sh` (kill-stale + bind-fresh + 200-gate + cleanUrls) kills the whole class — folds TODO #273 + #272 into one shape.
2. **Lighthouse numbers belong in a committed artifact** (plan §8 is prose; a small `docs/perf/2026-09-19.json` + the #272 lane would make deltas diffable).
3. **`git commit` racing discipline:** when a deliberate commit matters, `git add` + `git commit` must be ONE `&&` chain issued within seconds of the green check — staging in a separate call is what lost this one.
4. **Q3 (dark CTA 3.76:1) now has ONE canonical home** (ROADMAP "Dark CTA contrast decision") — the per-site ledger entries reference it; next session touching a ledger should link the ROADMAP row explicitly.
5. **htmx-on-static-pages is the biggest honest perf lever on the site** (~0.5s TBT on /sales): #282 records it; it needs a small library decision (`HTMXOff`/`HTMXSrc` semantics), not site-side hacks.

## f) NEXT TASKS (ordered)

1. **PUSH** (only remaining plan item): re-check `git status -sb` — if the tip moved past `bdb53316`, re-run ci-repro; else push immediately on your go. Verify CI + Website lanes green, then T3 prod spot-checks (`/sales` 200 + og:image, sitemap lastmod, CSP unchanged).
2. TODO #280 post-deploy bundle: Lighthouse on LIVE prod (brotli + Firebase headers — expect the compression artifact to vanish; record real numbers), sitemap, stars badge.
3. Run `check-draft.py` on the parked Scrollback draft; then decide G-question 3 below.
4. TODO #282 decision (G-question 2 below), then implement if approved: library HTMX-off semantics + site adoption + goldens.
5. TODO #272 `.#lighthouse` flake app or chromedp nav-timing budget (make this session's numbers reproducible).
6. TODO #273 + e1: `scripts/serve-dist.sh` (kill-stale, cleanUrls, 200-gate) — retire the python-server class.
7. TODO #271: flip `SITE_SKIP_STARS=1` to default for `.#shots`/`siteshots` (flag-gate production).
8. TODO #275: `docs/visual-testing.md` site-tier section (8 goldens, sweep, copy e2e, skip-stars pin).
9. TODO #276: SKILL.md website section (search scope, lastmod convention, `topLevelPages`).
10. TODO #279: warm-dark override pattern → `docs/theming.md` + website theming page.
11. TODO #278: verify demo CSS vs the `[&>details]:` classes (or prove demo routes don't use them).
12. TODO #274: library grep `role="combobox"` input-type audit + policy call.
13. TODO #277: repo sweep for `WriteString(literal + literal)`.
14. TODO #281: ogshot × skip-stars verification + README note.
15. TODO #216 (extended): prune the `media`-on-meta ignore class on the next vnu bump.
16. ROADMAP "Counts truth table": one counter behind `build.CountStats` + `TestDocsCountDrift` + the docs' visual-goldens number.
17. ROADMAP "gopls analyzer lint gates": evaluate writestring/prealloc in `.golangci.yml`.
18. ROADMAP "Production stars caching": last-good badge value between deploys.
19. ROADMAP "Dark CTA contrast decision": make the v2-scale call (re-shade to `-400` vs accept-and-document).
20. Standing watches (open-by-design): `*_templ.go` gitignore re-adds (f31), autoclose keywords post-push (f32), v0.3.1020 templ pin (f49), chroma `PreventSurroundingPre` behavior on upgrades (f49-adjacent).
21. Celebrate → next Pareto planning cycle from the refreshed TODO_LIST (#271–#282 are all ≤1h bounded).

## g) QUESTIONS FOR YOU (cannot answer myself)

1. **Push now?** The witnessed PASS at tip `bdb53316` (00:15:35) is still exact — tip unchanged at 00:16, ahead 25. If the daemon has committed anything new by the time you read this, I re-witness first (~7 min). Say "push" (or tell me to hold for the parallel session).
2. **htmx gating (#282):** the /sales TBT is dominated by parsing htmx it never uses. Add an HTMX-off semantic to `layout.PageProps` (library change, rides next minor) — or accept static-page htmx until v2?
3. **The parked Scrollback draft:** file it as an issue on this repo now, or keep it parked for the v2 API-review wave (the rename is v2-shaped either way; TODO #40 holds the work item)?

---

**Verdict:** plan 24/24 resolved. Every lane green at `bdb53316`: full visual suite (189s), ci-repro `--lint --website` VERDICT: PASS (witnessed 00:15:35), working tree clean, 25 commits ahead, nothing pushed — deliberately, per M03. One instruction ("push") ends the plan.
