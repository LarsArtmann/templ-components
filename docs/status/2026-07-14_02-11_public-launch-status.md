# Status Report — templ-components Public Launch

**Date:** 2026-07-14 02:11
**Session goal:** Make the repo public-ready — README, wiki website, GitHub metadata, domains/Firebase hosting, comprehensive component demo.

---

## A. Fully Done (verified, not assumed)

1. **README.md rewritten** — Superb landing page with badges, comparison table (vs templUI/goshipit), quick start, full component catalog with code examples, design principles, theming, ecosystem links. 279 lines.
2. **Wiki website built and deployed** — Astro + Starlight at `https://templcomponents.web.app`. 13 pages (landing + 11 doc pages + 404). 1.6MB total build. Verified live via HTTP fetch.
3. **GitHub metadata configured** — Description updated, homepage set to `https://templcomponents.lars.software`, 18 topics (golang, templ, htmx, tailwindcss, component-library, design-system, goth-stack, etc.).
4. **Firebase hosting site created** — `templcomponents` site in `lars-software` project. Deployed successfully (65 files).
5. **Firebase custom domain added** — `templcomponents.lars.software` registered via Firebase Hosting REST API.
6. **DNS CNAME committed to domains repo** — `templcomponents.web.app` CNAME in `lars.software.tf`, already committed (commit 1090a52).
7. **Comprehensive component demo built** — All 97 components across 9 packages demonstrated. Build passes, lint clean (0 issues), all tests pass. Verified via HTTP that every component renders (200 OK, 260KB HTML, 15 presence checks passed).
8. **`icons.AllIconNames()` exported** — Was private `allIconNames()`, now public for demos/consumers.
9. **Website design** — Blue accent (matching library primary), dark/light toggle, JetBrains Mono code, Space Grotesk headings, scroll animations, responsive grid, favicon, manifest.json, robots.txt, sitemap.

---

## B. Partially Done

1. **Custom domain `templcomponents.lars.software`** — Firebase domain added (status: DOMAIN_ACTIVE, cert: CERT_PENDING, DNS: DNS_MISSING). The DNS CNAME is committed to terraform but **terraform apply has NOT been run** (blocked by Namecheap API IP whitelist — `api.ipify.org` is blocked, can't auto-detect public IP). Firebase cert won't provision until the CNAME propagates in DNS.
2. **Website documentation pages** — 11 doc pages exist but are concise first drafts. They cover the essentials but could be deeper (more code examples, more edge cases, screenshots/gifs).
3. **Forms demo** — The old `forms_demo.templ` (standalone `/forms` page) still exists alongside the new `forms_section.templ`. Some duplication. The old page isn't linked from the main demo TOC but is still routable at `/forms`.

---

## C. Not Started

1. **Terraform apply for DNS** — CNAME `templcomponents.lars.software → templcomponents.web.app` needs to actually be applied to Namecheap DNS.
2. **SSL certificate provisioning** — Firebase cert is pending. Won't start until DNS resolves.
3. **Website redeployment after DNS** — Once custom domain works, need to verify `https://templcomponents.lars.software` serves correctly with SSL.
4. **Nothing committed in the templ-components repo** — All 19 files (modified + untracked) are uncommitted. No git commit has been made.
5. **Website CI/CD** — No GitHub Actions workflow for auto-deploying the website on push. Manual `firebase deploy` only.
6. **OG image / social preview** — No Open Graph image generated for the website or GitHub repo. Just text meta tags.
7. **Demo deployment** — The demo binary runs locally only. Not deployed anywhere accessible publicly.
8. **Analytics** — No analytics on the website.

---

## D. Totally Fucked Up / Mistakes Made

1. **`forms_demo.templ` left orphaned** — The old standalone forms demo page (`/forms` route) still exists but is now superseded by `forms_section.templ` which is rendered inline in the main demo. This is confusing duplication. Should either delete the old file or consolidate.
2. **Button `Attrs` field mistake** — I tried `display.ButtonProps{Attrs: ...}` but `ButtonProps` doesn't expose `Attrs` directly (it embeds `BaseProps` but the Button component doesn't spread it the same way). Had to fall back to raw `<button>` HTML with `onclick` for the modal/drawer open buttons. This works but is inconsistent with the library's own patterns.
3. **`icons.IconLeft` field mistake** — I initially wrote `IconLeft: icons.Plus` but `ButtonProps.Icon` takes a `templ.Component`, not an `icons.Name`. Fixed to `Icon: icons.Icon(icons.Plus, "h-5 w-5")`. This was a read-before-write failure.
4. **Port 8080 conflict** — Another process was already on port 8080. Had to run the demo on 18080 to verify. The `main.go` default port is fine for consumers but it was annoying for testing.
5. **Domains terraform credentials** — The `terraform.tfvars` has a placeholder API key (`REPLACE_WITH_YOUR_API_KEY`). I couldn't apply DNS changes. This is a hard external blocker.
6. **`forms.MethodGet` typo** — The constant is `forms.FormGet`, not `forms.MethodGet`. I guessed the name instead of checking.
7. **`errorpage.FamilyNotFound` doesn't exist** — The family enum is `FamilyRejection` (which covers 400/403/404). I guessed a name that doesn't exist.
8. **HoverCard signature mistake** — I initially tried to pass two children blocks (`} {`) but HoverCard takes one children block (the trigger) and the content is a `Content` prop. Had to extract a helper template.

---

## E. What We Should Improve

1. **Deploy the demo publicly** — A live demo at `demo.templcomponents.lars.software` or `templcomponents.web.app/demo` would be far more impactful than a local binary. Could be a separate Firebase hosting site or a route on the Astro site.
2. **Integrate demo into the website** — The Astro docs site should embed live component previews, not just code snippets. Consider an iframe to a deployed demo, or static HTML screenshots.
3. **Add OG image generation** — Like gogenfilter's `src/pages/og/[...slug].ts`. Social sharing needs a visual preview.
4. **Website CI/CD pipeline** — GitHub Actions to auto-build and deploy the Astro site on push to `website/`.
5. **Consolidate forms demo** — Delete or merge `forms_demo.templ` into `forms_section.templ`.
6. **Add code copy buttons to demo sections** — Each demo section should show the code that produces the output, with a copy button.
7. **Component count in demo is hardcoded** — `const componentCount = 97` in demo.templ should be derived from actual counts to avoid drift.
8. **Website docs need depth** — The 11 doc pages are good starters but need more examples, edge cases, and visual references.
9. **Search engine submission** — Submit `templcomponents.lars.software` to Google Search Console once DNS is live.
10. **Performance audit** — The website loads Google Fonts; consider self-hosting for privacy and performance.

---

## F. Next 50 Things To Get Done

### Priority 1: Unblock the custom domain

1. Whitelist public IP on Namecheap API → run `terraform apply` for the CNAME
2. Wait for DNS propagation (check with `dig templcomponents.lars.software`)
3. Verify Firebase SSL cert transitions to `CERT_ACTIVE`
4. Verify `https://templcomponents.lars.software` loads with valid SSL
5. Update GitHub repo homepage if URL changes

### Priority 2: Commit all work

6. Commit the README rewrite
7. Commit the website directory (all config + source + docs)
8. Commit the demo expansion (7 new files + modified demo.templ)
9. Commit the `icons.AllIconNames()` export
10. Commit generated `*_templ.go` files alongside `.templ` sources
11. Clean up orphaned `forms_demo.templ` (delete or consolidate)

### Priority 3: Deploy the demo publicly

12. Create a Firebase hosting site for the demo (or use a sub-path)
13. Build a static HTML version of the demo page
14. Deploy demo HTML to Firebase
15. Link the live demo from the README and website

### Priority 4: Website improvements

16. Add OG image generation route
17. Add GitHub Actions CI for website auto-deploy
18. Add "Edit on GitHub" links to doc pages
19. Add component count badge that auto-updates
20. Add a component gallery page (visual grid with links to API docs)
21. Write deeper docs: each component gets its own page with props table + examples
22. Add migration guide from templUI
23. Add a "Quick Start" video or animated GIF
24. Add search functionality (Starlight Pagefind is already bundled)
25. Add social preview image for GitHub repo card

### Priority 5: Demo improvements

26. Add live code examples alongside each component demo section
27. Add a "copy code" button to each demo section
28. Make component count dynamic (not hardcoded `const componentCount = 97`)
29. Add responsive viewport toggle (mobile/tablet/desktop) for testing
30. Add dark/light toggle to the demo page header (currently only ThemeToggle in header)
31. Add a search/filter for the icons gallery
32. Add interactivity to the demo (working form submission, toast triggering)
33. Deploy demo with proper Tailwind CSS build (currently uses Tailwind CDN)
34. Add demo for the `layout.Script` and `layout.Stylesheet` helpers
35. Add DataTable with actual sort/pagination working (needs server-side handler)

### Priority 6: Polish and ecosystem

36. Submit to pkg.go.dev (ensure latest version is indexed)
37. Add the website to the GitHub repo's "About" section as the homepage
38. Write a blog post or announcement for the GOTH stack
39. Add the library to awesome-go / awesome-htmx lists
40. Create a GitHub Discussion template for Q&A
41. Add issue templates (bug report, feature request)
42. Add a CONTRIBUTING.md link from the website
43. Set up GitHub Pages as a fallback mirror for the docs
44. Add benchmarks to the website (performance section)
45. Create a visual changelog (screenshots for each version)

### Priority 7: Code quality

46. Add contract tests for the new `AllIconNames()` export
47. Add a test that verifies the demo builds without errors in CI
48. Add a fuzz test for `AllIconNames()` (ensure it never returns empty)
49. Consider exporting `allIconNames()` test helper for consumers
50. Update AGENTS.md with the new demo file structure

---

## G. Top 2 Questions I Cannot Answer Myself

### 1. Namecheap API credentials / IP whitelist

**The problem:** Terraform apply is blocked because the Namecheap provider can't auto-detect the public IP (`api.ipify.org` is blocked by a firewall/proxy). The `terraform.tfvars` has a placeholder API key.

**What I need:** Either (a) the real Namecheap API key in `terraform.tfvars`, or (b) the public IP of this machine whitelisted at https://ap.www.namecheap.com/settings/tools/apiaccess/whitelisted-ips, plus `NAMECHEAP_CLIENT_IP` environment variable set. Without this, I cannot apply DNS changes and the custom domain will never resolve.

### 2. Should the demo be deployed as a static site or a live Go server?

**The problem:** The current demo is a Go binary (`examples/demo/main.go`) that renders templ components server-side at runtime. Deploying it requires a Go runtime (Cloud Run, Fly.io, etc.). Alternatively, I could pre-render the demo pages to static HTML at build time and deploy to Firebase Hosting (which only serves static files).

**What I need to know:** Is there an existing deployment target for Go binaries in the LarsArtmann infrastructure? Or should I pre-render the demo to static HTML for Firebase? The static approach loses interactivity (modals won't open, HTMX won't work, forms won't submit) but is deployable to Firebase Hosting with zero additional infrastructure.
