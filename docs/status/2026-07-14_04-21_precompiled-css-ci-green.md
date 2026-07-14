# Status Report — 2026-07-14 04:21 — Pre-Compiled CSS, CI Green, Pushed

## Session Scope

Full public launch of templ-components: CI fixes, demo expansion (97 components), website (Astro + Starlight), Cloud Run deployment (distroless, scale-to-zero), pre-compiled CSS (eliminating Tailwind CDN), README rewrite, docs audit, OG images. 15 commits, all pushed.

---

## A) FULLY DONE

| #   | Item                                                                                                                    | Evidence                                            |
| --- | ----------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------- |
| 1   | CI jsonv2 fix — removed broken import guard, added GOEXPERIMENT to both jobs                                            | `ci.yaml` committed, **CI green on GitHub Actions** |
| 2   | CI golangci-lint fix — install from source (Go 1.26) instead of pre-built v1 binary                                     | **CI green** (was failing on every prior commit)    |
| 3   | templ version pinned — `@v0.3.1020` in both CI workflows                                                                | Matches `go.mod`                                    |
| 4   | Demo expansion — 7 per-package files, all 97 components across 9 packages                                               | HTTP verified: 265KB, all sections present          |
| 5   | Dark mode fix — carousel slide colors gained `dark:` variants                                                           | `TestDarkModeCompliance` passes                     |
| 6   | `icons.AllIconNames()` exported                                                                                         | Build passes, used in demo                          |
| 7   | README rewrite — badges, comparison table, full catalog                                                                 | Committed, pushed                                   |
| 8   | Website (Astro + Starlight) — 53 source files, 11 MDX docs, Firebase config                                             | Builds, deployed, 200 OK                            |
| 9   | Docs accuracy audit — 10 errors found and fixed (6 would have prevented compilation)                                    | All verified against source                         |
| 10  | Website CI/CD workflow — build + Firebase + Cloud Run deploy on push                                                    | **Green on GitHub Actions**                         |
| 11  | Firebase service account — `firebasehosting.admin` + `run.admin` + `artifactregistry.writer` + `iam.serviceAccountUser` | All IAM roles granted, secret in GitHub             |
| 12  | OG image generation — astro-og-canvas, 14 PNGs, `summary_large_image` card                                              | Generated, deployed                                 |
| 13  | Distroless Dockerfile — multi-stage, `distroless/static-debian12:nonroot`, 12.8MB                                       | Built, pushed, deployed                             |
| 14  | Cloud Run deployment — us-central1, min=0, max=1, 256Mi, scale-to-zero                                                  | **Live**: 200 OK on all endpoints                   |
| 15  | Mock HTMX endpoints — `/api/load-more`, `/api/delete`                                                                   | 200 OK                                              |
| 16  | Cloud Run CI/CD — Docker build + push + deploy in website workflow                                                      | **Green on GitHub Actions**                         |
| 17  | **Pre-compiled CSS** — Tailwind compiled at build time (15KB), embedded via `//go:embed`, served from `/css/app.css`    | **No Tailwind browser CDN**, instant first paint    |
| 18  | Theme customization — Inter/Space Grotesk fonts, stone neutrals, indigo accent                                          | Eliminates Bootstrap look                           |
| 19  | All work committed and pushed — 15 commits, clean working tree                                                          | `git status` clean, `git log` verified              |
| 20  | Both CI workflows green on GitHub Actions                                                                               | `gh run list` verified                              |
| 21  | GitHub repo metadata — description, homepage URL, 18 topics                                                             | `gh repo view` verified                             |
| 22  | Website live at `templcomponents.web.app`                                                                               | 200 OK, 40KB                                        |
| 23  | Demo live at `templcomponents-demo-132045829579.us-central1.run.app`                                                    | 200 OK, 265KB, CSS 200 OK                           |

---

## B) PARTIALLY DONE

| #   | Item                         | What's done                                            | What's missing                                                                                                                                             |
| --- | ---------------------------- | ------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Custom domain**            | CNAME committed in domains repo                        | `terraform apply` blocked (Namecheap API IP whitelist). Website on `.web.app`, demo on `.run.app`.                                                         |
| 2   | **Demo interactivity**       | Main page + forms + mock HTMX endpoints work           | HTMX components use mock static endpoints, not real backend logic. Sufficient for demo.                                                                    |
| 3   | **Dockerfile CSS build**     | CSS compiled locally and committed as `static/app.css` | CI/CD workflow doesn't recompile CSS — if component classes change, the committed CSS goes stale. Need a Node stage in Dockerfile or CI step to recompile. |
| 4   | **Firebase Hosting rewrite** | Tested, documented                                     | Requires Firebase Blaze plan. Demo uses direct Cloud Run URL.                                                                                              |

---

## C) NOT STARTED

| #   | Item                                                                          |
| --- | ----------------------------------------------------------------------------- |
| 1   | DNS propagation — Namecheap API IP whitelist → terraform apply                |
| 2   | SSL cert — Firebase cert `CERT_PENDING`, blocked on DNS                       |
| 3   | Custom domain for Cloud Run demo (e.g., `demo.templcomponents.lars.software`) |
| 4   | Dockerfile Node stage for CSS compilation (CSS is currently committed static) |
| 5   | Website Lighthouse audit                                                      |
| 6   | Per-component documentation pages                                             |
| 7   | Website analytics                                                             |
| 8   | Google Search Console submission                                              |
| 9   | `.dockerignore` for faster Docker builds                                      |
| 10  | Health check endpoint (`/health`) on demo binary                              |

---

## D) TOTALLY FUCKED UP / WHAT WENT WRONG

| #   | Issue                                                 | Impact                              | Root Cause                                                                                        | Fixed?                                |
| --- | ----------------------------------------------------- | ----------------------------------- | ------------------------------------------------------------------------------------------------- | ------------------------------------- |
| 1   | **golangci-lint v1 binary can't lint Go 1.26**        | CI was red on EVERY push for months | `golangci-lint-action` installs pre-built binary built with Go 1.24. Go 1.26 > 1.24 = hard error. | ✅ Fixed: `go install` from source    |
| 2   | **templ binary not in PATH in CI**                    | Website workflow failed             | `go install` puts binary in `GOPATH/bin` which isn't in PATH by default in GitHub Actions         | ✅ Fixed: `export PATH` inline        |
| 3   | **Cloud Run `actAs` permission denied**               | CI Cloud Run deploy failed          | Deploy service account needed `iam.serviceAccountUser` on the runtime account                     | ✅ Fixed: IAM binding                 |
| 4   | **Firebase Hosting rewrite → Cloud Run returned 404** | Wasted ~30 min                      | Firebase Spark (free) plan doesn't support Cloud Run rewrites. Billing ≠ Firebase plan.           | ⚠️ Pivoted to direct URL              |
| 5   | **Deployed to europe-west1 first**                    | Wasted ~15 min                      | Didn't check Firebase region constraint                                                           | ✅ Cleaned up, redeployed us-central1 |
| 6   | **Accidentally created root `node_modules/`**         | BuildFlow pre-commit failed         | `npm install tailwindcss` in root instead of temp dir                                             | ✅ Deleted                            |
| 7   | **Corrupted main.go with bad edit**                   | Build failed                        | Edit tool replaced wrong function boundary, created duplicate handler                             | ✅ Rewrote entire file                |
| 8   | **Demo URL is ugly**                                  | Bad first impression                | `templcomponents-demo-132045829579.us-central1.run.app` — no custom domain mapped                 | ❌ Not fixed                          |

---

## E) WHAT WE SHOULD IMPROVE

| #   | Area                                 | Current                                                 | Target                                                      |
| --- | ------------------------------------ | ------------------------------------------------------- | ----------------------------------------------------------- |
| 1   | **Dockerfile CSS compilation**       | CSS compiled locally, committed static                  | Node stage in Dockerfile compiles CSS → Go embeds fresh CSS |
| 2   | **Demo URL**                         | `templcomponents-demo-132045829579.us-central1.run.app` | Custom domain: `demo.templcomponents.lars.software`         |
| 3   | **Docker image versioning**          | `:latest` only                                          | `:sha-<git-sha>` for rollback                               |
| 4   | **`.dockerignore`**                  | None                                                    | Exclude `.git`, `node_modules`, `website/`, `docs/`         |
| 5   | **Cold start UX**                    | First visitor waits ~2-3s                               | Min-instances=1 or loading indicator                        |
| 6   | **Demo search/filter**               | Scroll through 97 components                            | Client-side filter by name                                  |
| 7   | **Per-component docs**               | Package-level overview only                             | Individual component pages with live examples               |
| 8   | **OG image verification**            | Generated, never visually inspected                     | Fetch and verify rendering                                  |
| 9   | **CSS goes stale if classes change** | Committed static file                                   | CI step to recompile and verify                             |
| 10  | **Cost monitoring**                  | No budget alerts                                        | GCP budget alert                                            |

---

## F) NEXT 50 THINGS TO DO

### Critical — Public credibility

| #   | Task                                                                  | Est.   |
| --- | --------------------------------------------------------------------- | ------ |
| 1   | Whitelist Namecheap API IP, `terraform apply` for DNS                 | 10 min |
| 2   | Verify DNS propagation + SSL cert + custom domain                     | 10 min |
| 3   | Map custom domain to Cloud Run (`demo.templcomponents.lars.software`) | 15 min |
| 4   | Add `.dockerignore` to speed up Docker builds                         | 5 min  |
| 5   | Add Node stage to Dockerfile for CSS compilation                      | 15 min |

### High — First impression polish

| #   | Task                                                 | Est.   |
| --- | ---------------------------------------------------- | ------ |
| 6   | Add CSS recompile step to CI website workflow        | 10 min |
| 7   | Verify OG images render correctly (fetch `home.png`) | 5 min  |
| 8   | Add Docker image tagging with git SHA in CI          | 10 min |
| 9   | Set GCP budget alert ($5/mo)                         | 5 min  |
| 10  | Submit sitemap to Google Search Console              | 5 min  |
| 11  | Add demo link to Starlight docs sidebar              | 5 min  |
| 12  | Add "Edit on GitHub" links to docs pages             | 10 min |
| 13  | Add `/health` endpoint to demo binary                | 5 min  |
| 14  | Add `Cache-Control` header to demo HTML responses    | 5 min  |
| 15  | Add reading time + last-updated dates to docs        | 15 min |

### Medium — Quality and robustness

| #   | Task                                                        | Est.   |
| --- | ----------------------------------------------------------- | ------ |
| 16  | Add Lighthouse CI to website workflow                       | 20 min |
| 17  | Make HTML validation a blocking CI step                     | 5 min  |
| 18  | Add broken link checker                                     | 15 min |
| 19  | Add CSP headers to Cloud Run responses                      | 10 min |
| 20  | Rename `forms_demo.templ` → `forms_standalone.templ`        | 5 min  |
| 21  | Delete or document `prerender.go` as alternative deployment | 5 min  |
| 22  | Add table of contents to long doc pages                     | 15 min |
| 23  | Add prev/next nav to docs                                   | 10 min |
| 24  | Add search/filter to demo page                              | 30 min |
| 25  | Add dark mode toggle to demo (not just system)              | 10 min |
| 26  | Add browser support matrix to docs                          | 15 min |
| 27  | Write migration guide from templUI                          | 30 min |
| 28  | Add "Copy to clipboard" on all code blocks                  | 15 min |
| 29  | Write announcement blog post                                | 30 min |
| 30  | Add structured data to demo page                            | 10 min |

### Future — Nice to have

| #   | Task                                                         | Est.    |
| --- | ------------------------------------------------------------ | ------- |
| 31  | Add per-component documentation pages                        | 2 hours |
| 32  | Interactive playground (edit templ, see result)              | 1 day   |
| 33  | Theme customizer (pick colors, see all components update)    | 4 hours |
| 34  | Versioned docs (v0.x, v1.x)                                  | 4 hours |
| 35  | Add analytics (Plausible/Fathom)                             | 15 min  |
| 36  | Newsletter signup                                            | 30 min  |
| 37  | Sponsors section                                             | 15 min  |
| 38  | Contributing guide for new components                        | 30 min  |
| 39  | Community link (Discord/Slack)                               | 10 min  |
| 40  | Visual regression testing for demo                           | 1 hour  |
| 41  | Bundle size monitoring                                       | 15 min  |
| 42  | Multi-region Cloud Run (us + eu)                             | 30 min  |
| 43  | Uptime monitoring (UptimeRobot)                              | 10 min  |
| 44  | Demo page performance optimization (inline critical CSS)     | 1 hour  |
| 45  | Add ETag support to demo server                              | 15 min  |
| 46  | Cloud Run logging to BigQuery                                | 20 min  |
| 47  | Add `/api/toast` mock endpoint for toast demo                | 5 min   |
| 48  | Firebase Blaze plan upgrade (for Hosting rewrite)            | 10 min  |
| 49  | Tailwind v4 `@source` scanning in CI for CSS freshness check | 15 min  |
| 50  | Export demo as static downloadable HTML for offline use      | 30 min  |

---

## G) TOP 2 QUESTIONS

### 1. Dockerfile CSS compilation: Node stage or CI step?

The pre-compiled CSS (`static/app.css`) is committed and goes stale if component Tailwind classes change. Two options:

- **A) Node stage in Dockerfile:** `FROM node:22-slim AS css → FROM golang:1.26 AS go → FROM distroless` — fully self-contained, CSS always fresh, but Docker image build takes longer (Node pull + npm install)
- **B) CI step in website workflow:** Recompile CSS before Docker build, commit if changed — simpler Dockerfile, but couples CSS freshness to CI

**Which approach?**

### 2. Should we map a custom domain to Cloud Run now, or wait for DNS?

The demo URL `templcomponents-demo-132045829579.us-central1.run.app` is terrible for public consumption. Cloud Run supports custom domain mapping via GCP (independent of Namecheap/Firebase). We could map `demo.templcomponents.lars.software` → Cloud Run now if the `lars.software` DNS is already managed somewhere we can add an A/AAAA record.

**Should I attempt the Cloud Run custom domain mapping now, or is this blocked by the same Namecheap DNS issue?**
