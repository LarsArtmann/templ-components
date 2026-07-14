# Status Report — 2026-07-14 13:31

## CSS Source Scanning Bug Fix, Docker Hardening, Demo Polish

**Session:** 5 commits (`8e1d143` → `241cb4b`) on top of 14 from prior session.
**Total public-readiness effort:** 19 commits across 2 sessions.
**CI/CD:** Both workflows (CI + Website) green. All changes pushed and deployed.
**Working tree:** Clean.

---

## A. FULLY DONE ✅

### Critical: CSS @source Bug Fix

**The bug:** `examples/demo/demo.css` used `@source "./**/*.templ"` which resolved
relative to the CSS file location (`examples/demo/`), scanning ONLY demo files.
Every component class from library packages (`display/`, `forms/`, `feedback/`,
etc.) was missing from the compiled CSS. The entire live demo was rendering
**unstyled** — no Toggle animation, no Modal dialog, no Stylable Select, no
Textarea auto-grow, no Carousel scroll-snap, no Accordion chevron, no accent-color.

**The fix:**

1. Changed `@source` to `"../../**/*.templ"` to scan the entire repo from the
   CSS file's relative position.
2. Extracted component-specific CSS (dialog animations, stylable select,
   auto-grow, scroll-snap, accent-color, etc.) from `templates/app.css` into
   `templates/custom.css` — single source of truth imported by both the consumer
   template and the demo.
3. CSS went from 20KB (broken, demo-only) to 66KB (complete, all 97 components).

### Docker Hardening

| Item              | Before                               | After                                        |
| ----------------- | ------------------------------------ | -------------------------------------------- |
| `.dockerignore`   | Missing (653MB context)              | Added (37 rules, ~15MB context)              |
| CSS in Docker     | Committed static file (goes stale)   | Node 22 stage compiles fresh CSS every build |
| Dockerfile stages | 2 (Go → distroless)                  | 3 (CSS → Go → distroless)                    |
| Image size        | N/A                                  | 11.7MB                                       |
| CI workflow       | Manual templ generate + docker build | Simplified to just `docker build`            |

### Demo Polish

- `/health` endpoint — returns `{"status":"ok"}` for Cloud Run health checks
- `/favicon.svg` — indigo "tc" SVG served at the path `layout.Base` already links to
- All 26 inline scripts have CSP nonce
- 2 external CDN scripts (HTMX + response-targets extension)

### Documentation Updates

Updated all consumer-facing docs to reference `templates/custom.css`:

- AGENTS.md — 7 references updated + new "Demo Infrastructure" section
- CHANGELOG.md — extraction documented in `[Unreleased]`
- FEATURES.md — CSS Automation entry + accent-color location
- `docs/tailwind-v4-adoption-guide.md` — Option B starter template
- `docs/migration/play-cdn-to-tailwind-v4.md` — shortcut description

### Live Verification (all passing)

| Check                | Result                                      |
| -------------------- | ------------------------------------------- |
| Demo homepage        | 265KB HTML, 200 OK                          |
| Demo CSS             | 66KB, all 12 component CSS classes present  |
| Demo `/health`       | 200 `{"status":"ok"}`                       |
| Demo `/favicon.svg`  | 200 `image/svg+xml`                         |
| Demo dark mode       | 1,437 `dark:` class variants                |
| Demo CSP nonce       | 26/26 inline scripts have nonce             |
| Website              | 40KB, 200 OK                                |
| OG images            | `home.png` (46KB), `quick-start.png` (34KB) |
| Cloud Run cold start | 0.26s (instance was warm)                   |
| CI workflow          | Green (1m29s)                               |
| Website workflow     | Green (2m40s)                               |

---

## B. PARTIALLY DONE ⚠️

### Documentation Consistency for `custom.css` Extraction

Updated the primary consumer docs, but **stale references remain** in
historical/changelog sections that now point to `templates/app.css` for CSS
that lives in `templates/custom.css`:

| File                                   | Lines        | Issue                                                                                               |
| -------------------------------------- | ------------ | --------------------------------------------------------------------------------------------------- |
| `CHANGELOG.md`                         | 25-26        | `[Unreleased]` Added section references `templates/app.css` for `.tc-select` and `accent-color` CSS |
| `TODO_LIST.md`                         | 92, 137, 175 | Historical entries reference `templates/app.css` for RTL, accent-color, BuildFlow                   |
| `docs/adr/0014-dialog-migration.md`    | 76, 83, 86   | 3 references to `app.css` for `.tc-modal`/`.tc-drawer`/`.tc-overlay`                                |
| `docs/adr/0015-stylable-select-api.md` | 52, 76-77    | 3 references to `app.css` for `.tc-select` CSS                                                      |

These are **historical documents** — the CSS was in `app.css` when they were
written. Updating them is cosmetic, not functional, but they will confuse
consumers who follow the trail.

### CI/CD Pipeline

The Docker build is self-contained (CSS + templ + Go all inside Dockerfile),
but the CI workflow still uses `GOEXPERIMENT: jsonv2` env var manually. This
works but is a Go 1.26 workaround that becomes unnecessary in Go 1.27.

---

## C. NOT STARTED ❌

### Externally Blocked (require manual action outside the codebase)

1. **DNS propagation** — CNAME committed but `terraform apply` blocked
   (Namecheap API IP whitelist). `templcomponents.lars.software` doesn't resolve.
2. **Cloud Run custom domain** — Demo URL is ugly
   (`templcomponents-demo-132045829579.us-central1.run.app`). Needs
   `demo.templcomponents.lars.software` or similar.
3. **GCP budget alert** — No budget alert configured on `lars-software` project.
   Cloud Run min=0 max=1 + 256Mi is cheap, but an unbounded spike would be bad.

### Not Started (codebase work)

4. **Per-component documentation** — Only package-level overview exists on the
   website. No individual component pages with props tables, examples, live
   preview.
5. **Lighthouse audit** — No performance/accessibility/SEO audit run on either
   the website or demo.
6. **Analytics / Search Console** — No analytics installed, not submitted to
   Google Search Console.
7. **Website docs: `custom.css` references** — Website MDX files
   (`installation.mdx`, `theming.mdx`) still show only `app.css`, don't mention
   `custom.css`.
8. **Demo prerender for website** — Prerendered demo HTML exists in the codebase
   but isn't served on the website or linked from docs.
9. **CSS compilation in CI for website** — Website workflow compiles CSS in
   Docker, but doesn't run the CSS compilation step for the website itself
   (not needed — website uses its own Astro CSS, but worth noting).
10. **`robots.txt` on demo** — Demo has no `robots.txt`. Should disallow
    crawling or point to the main website.
11. **Sitemap** — Neither the website nor the demo has a `sitemap.xml`.
12. **Demo mobile menu** — Demo TOC navigation is a simple link list, not
    responsive on mobile.
13. **Rate limiting on demo HTMX endpoints** — `/api/load-more` and `/api/delete`
    have no rate limiting. Fine for a demo, but could be abused.

---

## D. TOTALLY FUCKED UP 💥

### The CSS @source Bug Itself

This bug was introduced in the prior session (commit `100eb60` — "eliminate
Tailwind browser CDN"). The CSS was compiled with `@source "./**/*.templ"` which
only scanned demo files. Nobody verified that component classes were actually
present in the compiled CSS before deploying. The demo was live for ~4 hours
rendering unstyled — every component missing its Tailwind utilities.

**Root cause:** Tailwind v4's `@source` directive resolves relative to the CSS
file's location, not the CWD where the compiler runs. This is documented but not
obvious. The prior session compiled CSS from the repo root
(`tailwindcss -i examples/demo/demo.css`) so it felt like `./` meant the repo
root — it doesn't.

**What I should have done:** After compiling CSS, verify that
component-specific classes (`peer-checked`, `tc-modal`, etc.) are present in the
output. I did this verification this session and caught it immediately.

### Favicon Data URI Attempt (commit `d2d7ce8`)

First favicon attempt used a `data:image/svg+xml,...` URI in a templ `href`
attribute. Templ's URL sanitizer rejected it (rendered as
`about:invalid#TemplFailedSanitizationURL`). I should have known templ
sanitizes URLs and tested before committing. Fixed in `241cb4b` by serving a
proper `/favicon.svg` endpoint.

---

## E. WHAT WE SHOULD IMPROVE 🔧

### Process Improvements

1. **Verify CSS output after compilation, every time.** The @source bug was
   invisible until someone loaded the demo. A simple `grep peer-checked
app.css` after compiling would have caught it. Add to BuildFlow or pre-commit.

2. **Test deployed demo, not just local build.** The bug was live for 4 hours
   because nobody fetched the live CSS URL and checked its contents. CI should
   fetch the deployed CSS and assert component classes are present.

3. **Don't trust `@source` relative paths.** Always verify by checking the
   compiled output for classes from files OUTSIDE the CSS file's directory.

4. **The CSS extraction was good, but the rollout was incomplete.** I updated
   5 primary docs but missed ADRs, CHANGELOG details, TODO_LIST, and website
   MDX. A grep-based sweep should have been the first step after extraction,
   not an afterthought.

### Architecture Improvements

5. **Single CSS entry point for the demo.** Currently `demo.css` imports
   `custom.css` and scans `../../**/*.templ`. If the directory structure changes,
   the relative paths break silently. Consider a build script or Makefile target
   that handles CSS compilation with explicit paths.

6. **Docker CSS stage could be cached better.** The Node stage copies the entire
   repo (`COPY . .`) before running `npm install`. If only Go files change, the
   Node stage rebuilds unnecessarily. Could split into `COPY` for CSS-relevant
   files only. (Minor — Docker layer caching helps.)

7. **No CSS regression test.** There's no test that asserts "if you compile
   `demo.css`, the output contains classes X, Y, Z." This would have caught the
   @source bug at build time.

---

## F. Next 50 Things to Get Done

### Priority 1: Unblock External Dependencies (do these first)

1. **Whitelist Namecheap API IP** and run `terraform apply` for CNAME — unblocks
   custom domain for both website and demo
2. **Map `demo.templcomponents.lars.software` to Cloud Run** —
   `gcloud run domain-mappings create`
3. **Set GCP budget alert** — `gcloud billing budgets create` for $5/mo
4. **Submit website to Google Search Console** — verify domain ownership
5. **Submit demo to Google Search Console** — or set canonical to website

### Priority 2: Documentation Consistency

6. **Update `CHANGELOG.md` lines 25-26** — change `templates/app.css` to
   `templates/custom.css` for `.tc-select` and `accent-color`
7. **Update `docs/adr/0014-dialog-migration.md`** — 3 references to `app.css`
   for dialog CSS
8. **Update `docs/adr/0015-stylable-select-api.md`** — 3 references to `app.css`
   for stylable select CSS
9. **Update `TODO_LIST.md`** — 3 historical entries reference `app.css` for CSS
   now in `custom.css`
10. **Update website `installation.mdx`** — mention `custom.css` alongside
    `app.css`
11. **Update website `theming.mdx`** — mention `custom.css`
12. **Add a CSS architecture doc** — explain the `app.css` (Tailwind directives)
    - `custom.css` (component CSS) split for consumers

### Priority 3: Demo Improvements

13. **Add `robots.txt` to demo** — point to main website or disallow
14. **Add `sitemap.xml` to website** — Astro has a sitemap plugin
15. **Run Lighthouse audit on demo** — performance, accessibility, SEO, best practices
16. **Run Lighthouse audit on website** — same
17. **Add mobile-responsive TOC** — current demo TOC is a flat link list
18. **Add dark mode screenshot to OG image** — current OG image is light-only
19. **Add "View Source" links** — link each demo section to the `.templ` source
    on GitHub
20. **Add interactive code examples** — let users edit props and see results
21. **Add search to demo** — filter components by name
22. **Add keyboard navigation** — `/` to focus search, arrow keys to navigate
23. **Add copy-button for component usage** — one-click import path copy
24. **Add loading states for HTMX demos** — show spinner during `/api/load-more`

### Priority 4: CI/CD and Infrastructure

25. **Add CSS regression test to CI** — assert compiled CSS contains
    `peer-checked`, `tc-modal`, etc.
26. **Add demo smoke test to CI** — after deploy, fetch `/health` and `/css/app.css`,
    verify 200
27. **Add dependency caching to Dockerfile Node stage** — cache `node_modules`
    between builds
28. **Add Dependabot/Renovate** — keep `tailwindcss`, `@tailwindcss/cli`,
    Astro, templ versions current
29. **Add `GOEXPERIMENT` comment in Dockerfile** — document that Go 1.27 removes
    the need for this flag
30. **Pin Node version in Dockerfile** — currently `node:22-slim`, could pin to
    `node:22.5.0-slim` for reproducibility
31. **Add Docker image vulnerability scanning** — `trivy` or `grype` in CI
32. **Set up Cloud Run min-instances=1** — eliminate cold start (~$5/mo) or
    accept the 2-3s cold start

### Priority 5: Website Improvements

33. **Add per-component documentation pages** — props tables, examples, live preview
34. **Add a component playground** — interactive props editor
35. **Add copy-paste installation snippets** — `go get` command per component
36. **Add a comparison page** — templ-components vs shadcn/ui vs templUI vs goshipit
37. **Add a "Theming" guide** — how to override colors, fonts, spacing
38. **Add an "HTMX Integration" deep-dive** — beyond the current guide
39. **Add a "CSP Compliance" guide** — how to set up Content-Security-Policy
40. **Add social sharing buttons** — Twitter/X, LinkedIn, Hacker News
41. **Add a GitHub star button** — prominent on the landing page
42. **Add analytics** — Plausible, Umami, or Google Analytics
43. **Add a changelog page** — render CHANGELOG.md as a webpage
44. **Add a migration guide** — for consumers upgrading between versions

### Priority 6: Code Quality

45. **Add a test for `@source` path resolution** — prevent the CSS bug from
    recurring
46. **Extract demo CSS compilation to a script** — `scripts/compile-demo-css.sh`
    with explicit path validation
47. **Add a `.nvmrc` or `engines` field** — specify Node version for CSS
    compilation
48. **Add `CONTRIBUTING.md` section on CSS** — how to add custom component CSS
49. **Review all `.tc-*` CSS classes** — ensure they're documented and tested
50. **Consider a CSS bundling strategy for consumers** — currently they must copy
    2 files (`app.css` + `custom.css`); could a single `go:generate` or
    `tc-css` tool help?

---

## G. Top 2 Questions I Cannot Answer Myself

### 1. Should the historical ADRs be updated or annotated?

`docs/adr/0014-dialog-migration.md` and `docs/adr/0015-stylable-select-api.md`
reference `templates/app.css` for CSS that now lives in `templates/custom.css`.
ADRs are supposed to be immutable historical records. Two options:

- **A) Update the references** — keeps them accurate for readers following links
- **B) Add an annotation** — "Note: CSS moved to `templates/custom.css` in [date]"

Which approach do you prefer? This affects how I handle all future ADR updates
when code locations change.

### 2. ~~Should we set Cloud Run min-instances=1?~~ RESOLVED

**Decision: min=0 (free, accept 2-3s cold start).** The demo is a showcase, not
a production app. Scale-to-zero keeps it free indefinitely. The `/health`
endpoint is already in place for when a warm-up strategy is needed later.

---

## Session Metrics

| Metric                         | Value                                                               |
| ------------------------------ | ------------------------------------------------------------------- |
| Commits this session           | 5                                                                   |
| Total public-readiness commits | 19 (across 2 sessions)                                              |
| Files changed this session     | 13 (+3,037 / -490 lines)                                            |
| Critical bugs fixed            | 1 (CSS @source — demo was fully unstyled)                           |
| New infrastructure             | .dockerignore, /health endpoint, /favicon.svg, Dockerfile CSS stage |
| CI status                      | ✅ Both workflows green                                             |
| Live demo                      | ✅ 66KB CSS, all component classes present                          |
| Working tree                   | Clean                                                               |
