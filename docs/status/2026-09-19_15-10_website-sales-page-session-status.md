# Status Report — Website Sales Page Session (`/sales`)

**Date:** 2026-09-19 15:10 CEST
**Session scope:** Build a sales page for templ-components, rendered by templ-components itself, on the `website/` module (the library's largest dogfood consumer).
**Outcome:** Shipped, verified, and committed. One real pre-existing bug found and fixed along the way.

> **Format note:** The status-report skill's canonical output is a styled HTML
> dashboard. The user explicitly requested `.md` at this path — user
> instruction wins (per the skill's own override rule).

---

## Verification evidence (one glance)

| Check                                                                                            | Result                                                                       |
| ------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------- |
| `templ generate` (pinned binary, repo root)                                                      | zero unrelated drift; only session files                                     |
| Website module `go build` + `go test ./...`                                                      | PASS (19 pages: index, sales, 404, 16 docs)                                  |
| `TestSiteBuildIntegrity`                                                                         | PASS (page count, internal links, script nonces, search index, sitemap)      |
| CSP guard (`firebase.json` vs rendered pages)                                                    | PASS after `--update-csp` (new CopyButton script hash committed)             |
| Goldens                                                                                          | `sales.golden` created; landing goldens updated; `TestGoldenSweepPages` PASS |
| Full dist build (`bash website/build.sh`)                                                        | PASS (Go SSG + Tailwind v4 minify)                                           |
| vnu HTML validation (sales, index, 404)                                                          | exit 0, zero errors                                                          |
| Visual smoke (siteshots: light/dark × desktop/mobile, `/sales` added to routes)                  | reviewed manually — all sections render, no overflow, dark mode coherent     |
| Search smoke (siteshots)                                                                         | PASS (8 hits)                                                                |
| Lint                                                                                             | website `0 issues`; visualtest tools `0 issues`                              |
| Repo-wide drift guards (`TestTemplGeneratedInSync`, `TestDocsCountDrift`, `TestVersionMatches*`) | PASS                                                                         |
| `nix fmt`                                                                                        | 0 changed                                                                    |
| `git status` at end                                                                              | clean (all work committed via daemon)                                        |
| `website/dist/` tracking                                                                         | gitignored, 0 tracked files (verified — not accidentally committed)          |

**Session commits (BuildFlow daemon, heuristic messages):** `1e4089ef` (sales.templ + wiring + test), `cd21010b` (generated file + goldens + CSP), `413e1c17` (site.css fix), `b3bdbbc0` (siteshots routes), `d2614027` (CHANGELOG + `staticPages` const). Verified: the union of these 5 commits is exactly the 12-file session diff, nothing more, nothing missing.

---

## a) FULLY DONE

1. **Sales page authored and shipped** — `website/internal/pages/sales.templ` (314 lines; generated `sales_templ.go` 772 lines committed for the module proxy). Structure follows the conversion framework: hero (transformation headline, single primary CTA, risk-reversal badges, install-path terminal) → problem (frontend tax, 3 cards) → benefits (4 outcome cards with icon tiles) → proof ("This page is the demo": 4 derived StatCards + CopyButton) → objection-handling FAQ (5-item Accordion) → final CTA (try-one-component risk reversal). Evidence: `1e4089ef`, `cd21010b`; rendered at `dist/sales.html`.
2. **All on-page counts derived, not hand-typed** — components/icons/enums/modules come from `build.CountStats` at build time (rendered: 123/105/62/7). The hand-typed-count drift class that once shipped ("58 enums") cannot recur on this page.
3. **Route + sitemap wiring** — `sales.html` registered in `cmd/site/main.go`; `/sales` sitemap entry added; magic number replaced by a documented `staticPages` constant shared with the integrity test. Evidence: `1e4089ef`, `d2614027`; `TestSiteBuildIntegrity` PASS, sitemap contains `<loc>…/sales</loc>`.
4. **Landing CTA hands off to the pitch** — "Still evaluating? Read the full pitch — every argument on one page." added to the landing's closing section. Evidence: `sections.templ` in `1e4089ef`; landing goldens updated.
5. **Golden coverage** — `sales.golden` (47 lines) + snapshot entry in `golden_test.go:61`; both landing goldens regenerated. Evidence: `cd21010b`, `1e4089ef`.
6. **Real bug fixed: site CSS never imported `templates/custom.css`** — `.tc-log`/`.tc-log-line` (Scrollback) and all other `.tc-*` component CSS were missing from the site's compiled stylesheet; the sales page was the first page to use a custom-CSS component, which surfaced it. `site.css` now imports `../templates/custom.css` (same single-source-of-truth pattern as the demo). Evidence: `413e1c17`; `.tc-log` present in `dist/assets/app.css` after rebuild.
7. **CSP header kept honest** — the CopyButton's nonce'd inline script added a new body hash; `firebase.json` updated via the documented `--update-csp` flow and the guard now passes. Evidence: `cd21010b`.
8. **HTML validation** — vnu (html5validator) exit 0 on `sales.html`, `index.html`, `404.html`.
9. **Visual verification** — `/sales` added to `visualtest/tools/siteshots` routes; light+dark × desktop+mobile captures reviewed: badge pill, accent headline, Scrollback terminal (shell/done tags), tone-colored StatCards, CopyButton, Accordion all render correctly in both themes; mobile wraps cleanly with zero horizontal overflow. Search smoke still PASS. Evidence: `b3bdbbc0`, `/tmp/site-shots-sales/`.
10. **Lint + drift guards** — website and visualtest modules lint clean; repo-wide guards (templ sync, docs counts, version triples) pass; `nix fmt` clean.
11. **CHANGELOG `[Unreleased]` warmed** — sales-page entry added immediately (release convention), duplicate-heading slip fixed same-session. Evidence: `d2614027`.

## b) PARTIALLY DONE

1. **Automated regression coverage for the sales page's _look_** — works today via manual siteshots review; what's missing: theme-pinned route goldens (like `TestDemoRouteGoldens`) and an axe-core a11y sweep for site routes. Blocker: none — effort, not feasibility. Effort to finish: M (route goldens), M (site axe sweep).
2. **Discoverability of `/sales`** — reachable only via the landing CTA link (and direct URL/sitemap). Header nav deliberately untouched this session (would change every page's golden). Remaining: a placement decision + one line + golden regen. Effort: S.
3. **Dark-mode visual blending of library surfaces** — every library class carries its `dark:` variant (compliance tests own that), and the page is coherent, but library `gray-*` surfaces are blue-tinted vs the site's warm near-black tokens (visible on Card header bands and Accordion items). Accepted as dogfood-authentic; not tuned via `Class` token overrides. Effort if wanted: S.
4. **Commit hygiene** — all work is committed and complete, but as 5 daemon "heuristic" chunks with no semantic story (documented daemon behavior; I am not authorized to commit and the daemon races master). Content verified complete at tip. Remaining: optional squash with a real message — needs a user call (master is daemon-pushed; history rewrite has risk). Effort: S.
5. **Claims auditability** — the copy keeps numbers derived (done), but two claims are hand-typed: "Go 1.26+" (hero badge) and "v1 under semantic versioning" (FAQ). Both can drift from `go.mod`/`utils.Version`. Deriving them is small but not done. Effort: S.

## c) NOT STARTED

1. **Header nav link to `/sales`** — deferred (golden blast radius); still wanted? Needs a decision.
2. **Sales content in the search index** — `build.SearchDoc` entries are docs-only by design; the pitch/FAQ is not searchable. Not started; needs a scope decision (docs-search vs site-search).
3. **Dedicated OG image** (`/og/sales.png`) — page reuses `/og/home.png`. Not started; needs design asset.
4. **Sitemap `lastmod` for `/sales`** — empty (like index); deriving from git is trivial but not started.
5. **Pre-push CI reproduction** — `scripts/ci-repro.sh --lint --website` was NOT run this session (M03 ritual applies at push time; nothing was pushed this session). Must run at the exact tip before any push.
6. **CopyButton browser-proof on the built site** — clipboard click-through on `dist/sales.html` (siteshots-style assertion) not started; the component itself is browser-proven by the library's own e2e suite.

## d) TOTALLY FUCKED UP

Nothing remains broken — build, tests, lint, CSP, and HTML validation are all green at tip. This section is the radical-honesty ledger of what I got wrong this session:

1. **Violated a documented repo rule on first write** — I set `AriaLabel` as a promoted field in a `ScrollbackProps` literal; AGENTS.md explicitly documents that promoted BaseProps fields CANNOT be set in struct literals. `go build` caught it in seconds; cost was one extra generate/build cycle. Root cause: I wrote the whole 314-line file before the first compile instead of compiling a skeleton first.
2. **Duplicate `### Added` heading in CHANGELOG** — my first edit's `old_string` was too narrow, producing two `### Added` sections. Caught by re-reading immediately; fixed. Same lesson: narrower confidence, wider verification.
3. **Self-mangled search output** — an `rg -rn "257"` typo (stray `-r` replace flag) rewrote my own grep output, briefly making docs counts unreadable. No files harmed; caught at once.
4. **mnd lint failure on `3+len(docsPages)`** — first version used a magic number where the original `2` had been; fixed with the `staticPages` constant (arguably the constant is the better end state, but the failure was avoidable).
5. **Inherited-but-real: the missing `custom.css` import** (fixed this session, `413e1c17`) — severity was genuine: any current-or-future use of Scrollback, Accordion chevrons, Modal/Drawer animations, or the stylable select on the SITE was silently unstyled. Root cause: the site's Tailwind entry only `@source`-scanned `templates/*.css` for class candidates but never imported the rules. Caught only because the first custom-CSS component landed on a site page — meaning the gap would have kept growing invisibly.

## e) WHAT WE SHOULD IMPROVE

1. **Site pages have no automated a11y gate.** The axe sweep, touch-target, and zoom-reflow audits cover demo routes only. The sales page passed vnu and manual review, but color-contrast class findings (the 2026-09-17 theme-pin class of bug) would be invisible here. Fix: extend the visualtest audits to the built site's routes.
2. **Site pages have no committed pixel goldens.** `TestDemoRouteGoldens` pins demo routes with theme-pinned PNGs; the website (landing, sales, docs) relies on manual screenshots. The unpinned-theme lesson of 2026-09-14 applies verbatim. Fix: route goldens for index + sales (+ one docs page), light+dark.
3. **Site HTML isn't validated in CI.** `scripts/check-html-valid.sh` covers only the 9 library package dirs. I validated ad hoc this session. Fix: include `website/internal/pages/testdata/*.golden` (they're full documents) or run vnu over `dist/` in the Website workflow.
4. **Integrity-test page count is still a hand-maintained constant.** `staticPages` is better than a literal, but `run()` could export its page list so the test derives the count structurally — page-count drift becomes impossible rather than guarded.
5. **Two homes for site content data.** Landing content lives in `data.go`; sales content lives in `sales.templ` (`SalesPains`/`SalesBenefits`). Both patterns are fine; pick one convention and document it before the split-brain grows.
6. **Compile early when authoring large `.templ` files.** The one compile error of the session (promoted-field) would have cost 30 seconds on a skeleton instead of a full-file cycle. Personal process fix, zero code.
7. **Scrollback timestamps are semantically abused as prompt glyphs** (`$`, `→` in the Timestamp column). Works visually, `AriaLabel` compensates, but a `Prompt`/prefix field on `Scrollback` would be the honest API. Candidate upstream idea (goes through verify-before-filing first).
8. **`siteScripts` loads `newsletter.js` on every page** that uses the shared meta (sales included). If the footer has no newsletter form, that's a dead fetch per pageview. Audit per-page script lists.
9. **Daemon commit chunks fragment session history.** Five heuristic commits for one feature makes `git log --grep` useless (documented T13 pain, hit again). A post-session squash ritual (user-authorized) or a daemon message upgrade in BuildFlow would fix the class.

## f) NEXT TASKS (brainstorm — HARVEST fuel, not commitments)

Ranked by impact. Impact: Critical/High/Medium/Low. Effort: S <30min, M 30min–2h, L >2h.

| #  | Task                                                                                                                                                | Impact   | Effort | Category      |
| -- | --------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------ | ------------- |
| 1  | Run `scripts/ci-repro.sh --lint --website` at tip before the next push (M03 ritual; not yet run this session)                                       | Critical | S      | Quality       |
| 2  | Add theme-pinned route goldens for `/sales` + index (light+dark, desktop+mobile) to `visualtest`                                                    | High     | M      | Quality       |
| 3  | Extend the axe-core sweep (+ touch-target, zoom-reflow) to the built site's routes                                                                  | High     | M      | Quality       |
| 4  | Add website goldens or `dist/` to the HTML validation gate (script + CI Website workflow)                                                           | High     | S      | Quality       |
| 5  | Derive the two hand-typed claims on the sales page: "Go 1.26+" from `go.mod`, "v1" from `utils.Version`                                             | Medium   | S      | Quality       |
| 6  | Export the site page-list builder so `TestSiteBuildIntegrity` derives the page count structurally                                                   | Medium   | S      | Quality       |
| 7  | Decide `/sales` header-nav placement; if added, regenerate all page goldens in the same commit                                                      | Medium   | S      | Feature       |
| 8  | Squash the session's 5 daemon heuristic commits into one semantic commit (needs daemon coordination / user authorization)                           | Medium   | S      | Cleanup       |
| 9  | Assert the dogfood claim in the integrity test: `sales.html` must reference no framework script (react/vue/alpine/htmx-CDN)                         | Medium   | S      | Quality       |
| 10 | CopyButton click-through proof on the built site (clipboard assert, siteshots-style smoke)                                                          | Medium   | M      | Quality       |
| 11 | Include sales/pitch content in the search index (scope decision: docs-search vs site-search)                                                        | Medium   | M      | Feature       |
| 12 | Per-page script audit: drop `newsletter.js` from pages without a newsletter form                                                                    | Low      | S      | Cleanup       |
| 13 | Tune library card/Accordion surfaces on the site via `Class` token overrides (dark tint mismatch)                                                   | Low      | S      | Quality       |
| 14 | Dedicated OG image for `/sales` (`/og/sales.png`)                                                                                                   | Low      | M      | Feature       |
| 15 | First FAQ item `Open: true` (conversion best practice; CRO call)                                                                                    | Low      | S      | Feature       |
| 16 | Anchor IDs for sales sections (`#problem`, `#proof`, `#faq`) for shareable deep links                                                               | Low      | S      | Feature       |
| 17 | Derive sitemap `lastmod` for `/sales` from git (parity with docs pages)                                                                             | Low      | S      | Quality       |
| 18 | Document the site-content convention (data.go vs per-page tables); move sales tables if data.go wins                                                | Low      | S      | Cleanup       |
| 19 | Link the FAQ's "three regression layers" claim to the actual testing docs pages                                                                     | Low      | S      | Documentation |
| 20 | Add derived social-proof to the sales hero (GitHub stars badge, landing-style)                                                                      | Low      | S      | Feature       |
| 21 | Mention the sales page from README (badge/link row)                                                                                                 | Low      | S      | Documentation |
| 22 | Lighthouse/perf budget check for landing + sales (single CSS file, fonts preloaded — should be clean)                                               | Low      | M      | Quality       |
| 23 | Fix pre-existing gopls `writestring` warnings in `cmd/site/main.go:334/337`                                                                         | Low      | S      | Cleanup       |
| 24 | Apply `templ QF1002` tagged-switch hint in `docs.templ:197`                                                                                         | Low      | S      | Cleanup       |
| 25 | Apply `templ QF1003` tagged-switch hint in `display/chart_shared.templ:59`                                                                          | Low      | S      | Cleanup       |
| 26 | Run full `visualtest` module lint (this session only linted `./tools/...`)                                                                          | Low      | S      | Quality       |
| 27 | Upstream idea: `Scrollback` `Prompt` field so prompt glyphs don't abuse the timestamp column (verify-before-filing first)                           | Low      | S      | Feature       |
| 28 | Reusable recipe doc: "dogfood marketing page" pattern (derived counts + self-render proof) for other LarsArtmann projects                           | Low      | S      | Documentation |
| 29 | Cross-link the FAQ cost answer to the LICENSE file                                                                                                  | Low      | S      | Documentation |
| 30 | When nixpkgs' vnu updates, prune the check-html-valid.sh ignore list (documented standing task; I re-confirmed the local/CI vnu split still exists) | Low      | S      | Cleanup       |
| 31 | Consider `data-animate` section-padding variant (py-24 hardcoded) if more campaign-style pages follow                                               | Low      | S      | Feature       |
| 32 | Confirm the Website CI path filters still fire for `website/**` changes post-Astro-removal (my change assumed it; CI run will prove it)             | Medium   | S      | Quality       |
| 33 | Post-deploy: verify `https://templcomponents.lars.software/sales` serves with the committed CSP header (Firebase headers apply)                     | High     | S      | Quality       |
| 34 | Decide analytics/conversion-measurement stance for the CTA (privacy-preserving or none) — needs owner input                                         | Medium   | S      | Decision      |

_(Stopped at 34 honest items — the remaining gap to 50 would be padding. Per the status-report skill: items beyond the core are ROADMAP fuel and need HARVEST routing rigor, not automatic TODO_LIST entries.)_

## g) QUESTIONS I CANNOT ANSWER MYSELF

1. **What is the distribution plan for `/sales`?** (README link? Social/HN/Reddit posts? Nothing — it only serves direct links?) This single answer drives items 7, 11, 14, 20, 21, 34: whether it needs header nav, search indexing, an OG image, social proof, and whether conversion measurement matters at all. I tried inferring from the repo (no ads/analytics exist, landing structure is overview-first) — your intent is unknowable from code.
2. **Do you want the daemon's 5 heuristic commits squashed into one semantic commit for this feature — and are you willing to coordinate the daemon (it owns master pushes) for that?** I verified the union of chunks is exactly the session diff, so a squash is safe content-wise; the risk is purely the daemon racing a history-moving operation (AGENTS.md documents it fighting rebases).
3. **How much machinery should a marketing page carry?** Concretely: derive the two remaining hand-typed claims (#5), add site-route a11y goldens (#2/#3), and assert the no-framework claim (#9) — or keep the sales page deliberately light and let the library's existing gates own quality? I can build either; the cost/benefit taste call is yours.

---

**Per the status-report skill:** section (f) is the primary input for `docs-health` HARVEST (TODO_LIST/ROADMAP). Not run yet — the session instruction was to report and then wait.

**Per the harness:** no manual commit of this report (Crush forbids commits without explicit request); the auto-commit daemon will pick it up.
