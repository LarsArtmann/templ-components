# Status Report — Website Hardening Session (P0 Items from the Astro→templ Conversion)

**Date:** 2026-09-13 18:38
**Scope:** This session executed the P0 backlog from `docs/status/2026-09-13_12-33_astro-to-templ-conversion-status.md` — the six post-conversion blockers for the templ-components website — plus the follow-up work those blockers surfaced. Nothing else (library code untouched except doc-count drift fixes).
**End state:** website builds clean, 12 Go tests green, 0 lint findings, 0 vet findings, all 6 sub-modules pass, drift guards pass, browser-verified via screenshots + live search smoke test.

---

## a) FULLY DONE

| #  | Item                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | Evidence                                                                                                                                                                                |
| -- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1  | **CI demo-deploy regression fixed.** The Astro rewrite had dropped `examples/demo/**` from `website.yml` path filters, so demo Docker/Cloud-Run redeploys only fired when `website/` changed.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | `.github/workflows/website.yml:8,14` (both push + PR); YAML parsed clean                                                                                                                |
| 2  | **Hash-based CSP, synced to firebase.json.** `build.InlineScriptHashes` hashes every inline `<script>` body (nonce attrs change per build; bodies don't → header is build-invariant). `build.CSPHeader` emits `default-src 'self'; script-src 'self' 'sha256-…'×4; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'; manifest-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'`. Every build **checks** the committed header and fails with the fix command; `--update-csp` rewrites it. Design note: nonce-based CSP was rejected because the deploy job re-checks-out firebase.json and can't see the build job's random nonce.                                                           | `website/internal/build/csp.go`; `website/firebase.json` (10 headers incl. CSP); check-pass verified on a second run                                                                    |
| 3  | **404 empty-nonce bug fixed.** `<script nonce="">` (tcGoBackAttached) on 404.html — `siteNotFoundProps` never set `BaseProps.Nonce`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | `website/internal/pages/notfound.templ:26`; `rg 'nonce=""' dist/*.html` → 0 hits                                                                                                        |
| 4  | **Go-native docs search, zero dependencies.** Build emits `dist/search-index.json` (14 docs: title, description, heading anchors, plain-text bodies via `build.PlainText`). `assets/js/search.js` = CSP-safe vanilla-JS combobox: lazy fetch, AND-term scoring (title 25–40 / heading 15 / body 10), max 8 hits, `role="combobox"/"listbox"/"option"`, ArrowUp/Down/Enter/Escape, `<mark>` snippets, per-section anchor links, click-outside + pointerdown-safe navigation. One instance in the header serves desktop AND the mobile menu panel.                                                                                                                                                                                                                 | `website/internal/build/search.go`, `website/assets/js/search.js`, `website/internal/pages/header.templ` (docSearch); `site.css` search block; browser smoke: 8 hits (search-smoke.png) |
| 5  | **Website module test suite (was zero tests).** 12 tests: CountStats against a fixture repo (components/icons-dedupe/enums/`_test.go`-exclusion/modules), PlainText (script/style strip, entities, whitespace), WriteSearchIndex, CSP ×4 (hash extraction + dedupe + nonce-independence, header directives, sync round-trip idempotence, stale/missing/absent errors), link checker ×3 (clean site, 5 problem classes, resolveTarget table), page goldens (landing ±stars, 404, docs layout), and `TestSiteBuildIntegrity` — renders the WHOLE site into a temp dir and asserts page count (2+len(AllDocs)), zero broken internal links/anchors, script-nonce policy (JSON-LD exempt as non-executable data block), search-index validity, sitemap completeness. | `website/internal/build/{csp,build,links}_test.go`, `website/internal/pages/golden_test.go`, `website/cmd/site/main_test.go`; all green, lint 0 issues, vet clean                       |
| 6  | **Link/anchor checker.** `build.CheckLinks` validates every href/src: absolute-path enforcement, page resolution (clean URLs → `.html`), asset resolution, same-page + cross-page anchors against `id=` sets. Caught 205 broken asset refs in the test before the writeAssets fix; passes on real dist.                                                                                                                                                                                                                                                                                                                                                                                                                                                          | `website/internal/build/links.go`                                                                                                                                                       |
| 7  | **Visual pass tooling.** `visualtest/tools/siteshots`: serves dist with Firebase cleanUrls semantics, captures 5 routes × light/dark × desktop/mobile (20 full-page PNGs), fresh browser per page (per the #shots hang lesson), theme pinned via localStorage+reload (headless Chromium 152 reports `prefers-color-scheme: dark` by default!), scroll-through so IntersectionObserver reveals fire, plus the search smoke check.                                                                                                                                                                                                                                                                                                                                 | `visualtest/tools/siteshots/main.go`; `/tmp/site-shots/*` reviewed (landing light+dark, docs desktop, docs mobile, 404, search panel)                                                   |
| 8  | **Three real site bugs found & fixed by the visual pass:** (a) docs sidebar leaked a literal `continue` text node AND a broken `<a href="/">Full API on pkg.go.dev</a>` on every docs page — templ v0.3.1020 does not support bare `continue` in template loops (it emits it as text and falls through); restructured to if/else-if/else, goldens regenerated, dist greps clean. (b) Landing sections invisible without JS — `[data-animate]` hiding now gated on `html.js` set by animations.js. (c) Site build depended on CWD — writeAssets now resolves assets/public under `--repo-root`.                                                                                                                                                                   | `docs.templ` sidebar; `site.css` + `animations.js`; `main.go` writeAssets; golden diff shows only the sidebar fix                                                                       |
| 9  | **Docs stale-facts audit.** api-reference.md per-package table corrected (display 30→43, feedback 13→14, forms 21→23, layout 6→10, htmx 8→9, icons 102→105; added `datastar` 4 and `recipes` 4 rows — counts derived with the same regex the build uses, total 123 ✓); installation.md "5-module workspace"→7-module, "(not an pnpm package)" phrasing removed.                                                                                                                                                                                                                                                                                                                                                                                                  | `website/content/docs/api-reference.md`, `getting-started/installation.md`                                                                                                              |
| 10 | **CHANGELOG + AGENTS updated.** `[Unreleased]` warmed with 4 Added + 5 Fixed entries (release script will refuse an empty section). AGENTS.md website bullet expanded to cover the CSP guard, search, tests, CI filter rule, and four new gotchas; two stale references fixed (`website/pnpm-workspace.yaml`, `website/src/data/sections.ts` — deleted with Astro).                                                                                                                                                                                                                                                                                                                                                                                              | `CHANGELOG.md`, `AGENTS.md:148-152,278,359`                                                                                                                                             |
| 11 | **Full verification sweep.** Website: tests ✓, golangci-lint 0 issues, vet clean. All 6 sub-modules: OK. Drift guards: OK — which caught an unrelated 244→246 golden-count drift in FEATURES/ROADMAP/AGENTS (other sessions added goldens without updating claims); fixed to 246 and re-verified. Root-module spot tests (utils/layout/errorpage): OK.                                                                                                                                                                                                                                                                                                                                                                                                           | session logs; `utils.TestDocsCountDrift` green                                                                                                                                          |

## b) PARTIALLY DONE

| # | Item                                     | State                                                                                                                                                                                                                                                                                                                                                               |
| - | ---------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **Firebase cleanUrls live verification** | Config verified statically (`cleanUrls: true`, trailingSlash false, /docs redirect) + local server mimics cleanUrls semantics for screenshots. NOT verified on real Firebase serving — needs a preview-channel deploy (credentials + network).                                                                                                                      |
| 2 | **HTML validation**                      | The Go link checker covers dead links/anchors but not markup well-formedness. CI's `html-validate` step is still `continue-on-error: true` and was never run locally this session. Decision pending: keep (Node) vs replace with a Go validator vs make it blocking.                                                                                                |
| 3 | **Golden coverage breadth**              | 3 page shapes golden'd (landing ±stars, 404, docs layout). Individual landing sections, all 14 docs pages from real content, and 404 variants are not golden'd. Notably: my docs-layout golden used the REAL sidebar registry — it silently captured the `continue` bug on `-update`. Golden diffs must be eyeballed (AGENTS says so; I skipped it the first time). |
| 4 | **Search UX polish**                     | Core works (browser-proven at desktop). Missing: typo/fuzzy tolerance, result-count line, `/`-or-`Ctrl+K` focus shortcut, scroll-spy-style active hit announcement, and an actual browser check at mobile viewport (the panel lives inside the collapsible menu there).                                                                                             |
| 5 | **siteshots integration**                | The tool exists and works but is manual-only: not a flake app, not a CI step, screenshots not persisted as testdata, no pixel-diff regression mode (that's what `.#visual` does for the library).                                                                                                                                                                   |
| 6 | **CSP telemetry**                        | Strict CSP is now enforced but there is no `report-uri`/`report-to` — a real-world violation would fail silently for visitors.                                                                                                                                                                                                                                      |

## c) NOT STARTED

(Backlog items from the conversion status report and successors — untouched this session.)

1. ~~Go OG-image generator (static preserved PNGs still in use).~~ done (routed to ROADMAP Website-and-docs-ideas (2026-09-17 evening harvest))
2. ~~`nix run .#website` flake app.~~ done (routed to ROADMAP Website-and-docs-ideas)
3. ~~Lighthouse CI / performance budgets.~~ done (routed to ROADMAP Website-and-docs-ideas)
4. ~~Firebase preview-channel deploy workflow (credential-gated).~~ done (routed to ROADMAP Website-and-docs-ideas)
5. ~~RSS/Atom feed for the changelog.~~ done (routed to ROADMAP Website-and-docs-ideas)
6. ~~Docs pages' structured data (BreadcrumbList/TechArticle JSON-LD; only landing has SoftwareApplication).~~ done (routed to ROADMAP Website-and-docs-ideas)
7. ~~404 page's `SearchAction` wiring (NotFound404 search form unused — could point at the new search).~~ done (routed to ROADMAP Website-and-docs-ideas)
8. ~~Mobile docs sidebar (the docs nav is entirely hidden below `lg`; only header search/menu remain).~~ done (routed to ROADMAP Website-and-docs-ideas)
9. ~~Newsletter form hardening (buttondown embed only; no validation feedback).~~ done (routed to ROADMAP Website-and-docs-ideas)
10. `scripts/ci-repro.sh` full run (per-module loop + lint + website were run piecemeal instead).

## d) TOTALLY FUCKED UP

Nothing shipped broken — the final tree is verified (build, tests, lint, vet, browser smoke). But the _process_ had real faceplants, recorded here so they don't repeat:

1. **Three consecutive test-authoring failures.** Wrote the icon fixture in a format the regex can't match (regex parens are a capture group, not literals — didn't check the real `icon_names.go` first), asserted `ContainsAny(got, "<>&;")` which fails on legitimately decoded `&`/`>`, and asserted `ContainsAny(got, "<>")` which fails on decoded `>=` in text. Each was my assumption, not reality.
2. **chromedp API guessed from memory, wrong 3×** (`NewExecAllocator` returns 2 values; `FullScreenshot` takes `*[]byte` here; `SetValue` doesn't exist). The repo's own harness had every answer — I read it only after vet/runtime failures.
3. **Assumed `CHROMEDP_CHROME_PATH` is a chromedp feature.** It isn't — the repo's tools wrap it via `chromedp.ExecPath`. First run died with `exec: "google-chrome" not found`.
4. **The dark-capture bug nearly shipped as "verification".** First run's light/dark captures were byte-identical on 3 of 5 pages because headless defaults to `prefers-color-scheme: dark`. I noticed only because PNG file sizes matched exactly. A "visual pass" that captures the wrong thing is worse than none.
5. **Stale base-state assumptions.** Asserted 12 docs from the previous session's summary; reality is 14 (16 pages). The base branch had moved; I didn't re-derive counts before asserting.
6. **Edit-tool self-inflicted wound:** replaced `chromePath()` with `scrollRevealJS` instead of appending after it — vet caught the dangling references.
7. **Wrong-file edits + stale-read churn:** tried editing notfound content inside docs.templ, and hit "modified since read" twice (daemon formatter) after reading via `sed` instead of `view`.
8. **`lsp_replace_symbol` attempted first** despite the session-documented fact that gopls is broken for this module — method unsupported, one wasted call.
9. **The `continue` bug had survived since the conversion** because the golden test captured it as "expected output" on `-update`. The guard pattern itself failed until a human looked at pixels.

## e) WHAT WE SHOULD IMPROVE

1. **Read the existing harness before writing new tooling.** Three chromedp API mistakes and the env-var discovery would all have been avoided by one read of `visualtest/tools/shots/main.go` and `harness.go` first.
2. **Re-derive base-state numbers before writing assertions.** The "12 docs" failure class is exactly what AGENTS' "status reports are point-in-time" warning predicts.
3. **Always eyeball golden diffs on `-update`.** Goldens capture bugs as expectations; the `-update` step is a review step, not a formality (I skipped it; the `continue` leak rode in through it).
4. **Write the probe test first.** The icon-regex debug (temp test + real file) is what I should have done before writing the fixture — and the PlainText assertions should have been derived from what the function is ALLOWED to output (decoded text), not what looked "clean".
5. **Content-first defaults in design review.** `opacity: 0` until JS is a JS-dependency smell; the `html.js` gate should be the standing pattern for any scroll-reveal in this repo.
6. **Freshness-guard pattern applies to more than CSS.** The CSP check (fail build with fix command) is the TestCSSFreshness pattern generalized; consider it for any committed artifact derived from code (search-index is build-only, fine).
7. **Mechanical count updates are fine via script — with the guard test as the verifier.** The 244→246 sweep was safe precisely because `TestDocsCountDrift` re-verified it.
8. **Trust shell, not LSP, in this repo — consistently.** gopls reported 76 phantom errors all session; only the workspace build + per-module tests were truthful (as AGENTS already warns).

## f) NEXT (up to 50, brainstorm — harvest into TODO_LIST/ROADMAP, roughly impact-ordered)

**Website correctness & CI**

1. ~~Run the new website tests in CI: add `go test ./...` step to `website.yml` (currently only build+lint — the integrity test is not enforced yet).~~ done (routed to TODO_LIST #241 (website.yml go-test step))
2. Make the CSP firebase.json check fail the CI build explicitly (it already fails `build.sh` — verify the workflow surfaces it).
3. Deploy to a Firebase **preview channel**; verify cleanUrls (`/slug` → `slug.html`), CSP header delivery, and 404 handling live.
4. Decide + execute production deploy (site currently still serves the old Astro build until Firebase picks up the new dist).
5. ~~Add CSP `report-uri`/`report-to` telemetry (or at least document how to inspect violations).~~ done (routed to ROADMAP (CSP violation telemetry))
6. Replace the Node `html-validate` CI step (continue-on-error) with either a blocking run or a Go-based validator; delete the Node step if the Go checker suffices.
7. Wire `scripts/ci-repro.sh` website step to also run website tests (parity with CI).
8. Run `nix run .#visual` full suite + `scripts/ci-repro.sh --lint` once before the next push (not done this session).
9. Post-deploy smoke: hit `https://templcomponents.lars.software` + `/health` demo endpoint in CI after deploy.
10. ~~Re-check daemon same-day regressions after these commits land (CSS minification, config flips — historical pattern) + confirm CI green.~~ done (SUPERSEDED - the 09-14 mr-status session + 09-17 release session handled the daemon regressions)

**Website features**
11. Go OG-image generator (per-page PNGs) replacing the static preserved `public/og/*`.
12. Mobile docs navigation: expose the sidebar in the mobile menu panel (currently hidden below `lg`).
13. Search: `/` + `Ctrl+K` focus shortcut; Escape-to-blur; result count; "no results" suggestions.
14. Search: fuzzy/typo tolerance; show section breadcrumb in hit subtitle.
15. Search: browser-verify at mobile viewport (panel inside open menu) + debounce + index prefetch on focus.
16. Wire NotFound404's `SearchAction` to the new search (or remove the form).
17. Docs pages: emit BreadcrumbList JSON-LD + `rel="prev"/"next"` links (SEO + a11y).
18. Docs TOC scroll-spy active state.
19. Header "Docs" nav link active state on docs pages.
20. Newsletter form: inline validation + success/error feedback (currently bare buttondown embed).
21. RSS/Atom feed for the changelog page.
22. `search-index.json` Cache-Control (short TTL — fixed name, must not be immutable).
23. Content-hashed asset filenames so app.css/app.js changes don't fight the 1-year immutable cache.
24. Lighthouse CI with perf/a11y/SEO budgets.
25. Accessibility sweep of the site using the library's existing axe tooling (`visualtest/axe.go`) against real dist pages.
26. Contrast check for search dropdown + `mark` styling in both themes.
27. OG-image existence check in the integrity test (every `OGImage` path exists in dist).
28. Full read-through of remaining docs pages for stale claims (numeric grep ≠ content audit).
29. Sitemap: homepage `lastmod`, verify `robots.txt` in dist, decide `sitemap-index` vs direct reference.
30. Code blocks: language label / filename support in the md renderer.
31. `nix run .#website` flake app + `.#siteshots` flake app (mirror `.#shots`).
32. Document siteshots usage in `docs/visual-testing.md`.
33. Consider persisting siteshots captures as PR artifacts for visual review.

**Library/QA adjacent (noticed this session)**
34. Golden-test every docs page from real content (string tests + fixture goldens let the `continue` bug survive two sessions).
35. Guard test: repo-wide scan for bare `continue`/`break` inside `.templ` files (templ v0.3.1020 silently leaks them as text).
36. Guard test: `[data-animate]`-style hide-by-default CSS must be gated on `html.js` (scan site.css + custom.css).
37. templ upstream: verify whether `continue`/`break` support landed post-v0.3.1020; note in AGENTS until then (possible upstream issue to file).
38. Consider documenting the hash-CSP pattern as a recipe in `docs/` (it's reusable for any static consumer of the library).
39. Integration test: render a page with `SEO.JSONLD` set and assert hash-stability + nonce behavior (library-side complement to the site's guard).
40. `visualtest/go.mod`: confirm tools additions don't affect `go test ./...` runtime (vet only so far).

**Docs/housekeeping**
41. ~~Annotate `docs/status/2026-09-13_12-33_astro-to-templ-conversion-status.md` P0 items as resolved (docs-health ANNOTATE).~~ done (DONE - annotations executed in the 2026-09-17 evening docs-health pass)
42. Harvest this report's items into `TODO_LIST.md` / `ROADMAP.md` (docs-health HARVEST).
43. Update `docs/visual-testing.md` + README docs section for the website (search + CSP + tests are user-visible facts).
44. ~~Move `[Unreleased]` entries into the next release cut when ready (release.sh flow, tag set for all modules).~~ done (DONE - v1.18.0 shipped 2026-09-17 (511d3ed6))
45. Run `-race` on website tests once (repo standard for library tests).
46. Measure search-index.json size trend as docs grow (14 docs now; revisit client-side approach at ~50).
47. Font audit: confirm the 5 self-hosted woff2 cover all used weights (600 for headings?) and preloads match first paint.
48. Prune `docs/planning/2026-09-13_11-15` plan doc status line to reflect hardening completion (point-in-time doc, annotate only).
49. Demo repo hygiene: `examples/demo/main.go` + `wire_demo.templ` are dirty from the concurrent session — coordinate before any sweep.
50. Decide the html-validate/Node question once and remove the ambiguity from CI (ties to #6).

## g) QUESTIONS (cannot figure out myself)

1. **Deploy policy:** Should the next master push go straight to **production** Firebase Hosting (replacing the old Astro site), or do you want a **preview channel** first for your eyeball pass? This decides items #3/#4 and whether I need to wire any deploy gating.
2. **CI strictness for the website:** Do you want the website workflow to run the new Go tests + link checker as a **required** step, and should the Node `html-validate` step be kept (made blocking), or dropped in favor of the Go checker?
3. **Search scope & UX:** Should the index stay **docs-only** (current), or also include landing sections and (future) recipe pages? And do you want the full keyboard UX (`Ctrl+K`, modal-style) or the current minimal header combobox?

---

_Point-in-time snapshot. The auto-commit daemon picks this file up; per the harness contract no manual commit was made. Section (f) is HARVEST input for `TODO_LIST.md`/`ROADMAP.md`, not a commitment list._
