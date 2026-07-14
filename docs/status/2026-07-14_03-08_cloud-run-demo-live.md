# Status Report — 2026-07-14 03:08 — Cloud Run Demo Live

## Session Scope

Two back-to-back work phases: (1) public launch prep (README, website, demo
expansion, CI/CD, docs audit), (2) Cloud Run deployment of the demo.

10 commits, 79 files changed (+14,200 / -1,681).

---

## A) FULLY DONE

| #   | Item                                                                                                                                    | Evidence                                   |
| --- | --------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------ |
| 1   | **CI jsonv2 fix** — removed broken import guard, added `GOEXPERIMENT=jsonv2` to both CI jobs                                            | `ci.yaml` committed, build passes locally  |
| 2   | **templ version pinned** — both CI workflows use `templ@v0.3.1020` (matches `go.mod`)                                                   | `grep templ@` verified                     |
| 3   | **Demo expansion** — 7 per-package files, all 97 components across 9 packages                                                           | HTTP verified: 260KB, all sections present |
| 4   | **Dark mode fix** — carousel slide colors gained `dark:` variants                                                                       | `TestDarkModeCompliance` passes            |
| 5   | **`icons.AllIconNames()` exported**                                                                                                     | Build passes, used in demo                 |
| 6   | **README rewrite** — badges, comparison table, full catalog                                                                             | Committed                                  |
| 7   | **Website** — 53 source files, Astro + Starlight, Firebase config                                                                       | Builds (13 pages + 14 OG images), deployed |
| 8   | **Docs accuracy audit** — 10 errors found and fixed (6 would have prevented compilation)                                                | All code examples verified against source  |
| 9   | **Website CI/CD workflow** — build + Firebase deploy on push                                                                            | `website.yml` committed                    |
| 10  | **Firebase service account** — `github-website-deploy@lars-software`, `firebasehosting.admin` + `run.admin` + `artifactregistry.writer` | IAM policy applied, secret in GitHub       |
| 11  | **OG image generation** — `astro-og-canvas`, 14 PNGs, `summary_large_image` Twitter card                                                | Generated, deployed                        |
| 12  | **Demo pre-rendering** — `-prerender` flag, static HTML output                                                                          | Works locally, used as fallback            |
| 13  | **Distroless Dockerfile** — multi-stage build, `distroless/static-debian12:nonroot`, 12.8MB                                             | Image built, pushed, deployed              |
| 14  | **Cloud Run deployment** — `us-central1`, min=0, max=1, 256Mi, scale-to-zero                                                            | Live: 200 OK on all 3 endpoints            |
| 15  | **Mock HTMX endpoints** — `/api/load-more`, `/api/delete`                                                                               | 200 OK, 222b HTML response                 |
| 16  | **Cloud Run CI/CD** — Docker build + push + deploy on push to master                                                                    | Workflow committed                         |
| 17  | **Artifact Registry** — `us-central1-docker.pkg.dev/lars-software/templcomponents`                                                      | Image pushed                               |
| 18  | **All work committed** — 10 commits, clean working tree                                                                                 | `git status` clean                         |
| 19  | **GitHub metadata** — description, homepage URL, 18 topics                                                                              | `gh repo view` verified                    |
| 20  | **Website deployed** — `templcomponents.web.app`                                                                                        | 200 OK, 40KB                               |
| 21  | **Demo deployed** — `templcomponents-demo-132045829579.us-central1.run.app`                                                             | 200 OK, 260KB                              |

---

## B) PARTIALLY DONE

| #   | Item                         | What's done                                  | What's missing                                                                                                                                                                                                                     |
| --- | ---------------------------- | -------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Custom domain**            | CNAME committed in domains repo              | `terraform apply` blocked (Namecheap API IP whitelist). `templcomponents.lars.software` doesn't resolve yet. Website + demo both on default `.web.app` / `.run.app` URLs.                                                          |
| 2   | **CI/CD end-to-end**         | Workflows written, secrets configured        | **Never pushed to remote.** GitHub Actions has never run any of these workflows. Could fail on first run.                                                                                                                          |
| 3   | **Demo interactivity**       | Main page + forms + mock HTMX endpoints work | HTMX components in the demo (LoadMore, ConfirmDelete, SwapOOB) point to mock endpoints, not real server state. The mock returns static HTML — functional for demo, not a real backend.                                             |
| 4   | **Firebase Hosting rewrite** | Configured, deployed, tested                 | **Requires Firebase Blaze plan.** Rewrite returned 404 — Firebase free tier doesn't proxy to Cloud Run. Pivoted to direct Cloud Run URL. Upgrading to Blaze would enable `templcomponents.web.app/demo/` path.                     |
| 5   | **Cloud Run cleanup**        | europe-west1 service + registry deleted      | The original Dockerfile used `golang:1.26-alpine` builder, changed to `golang:1.26` (Debian) because alpine doesn't ship `git` which `go mod download` needs. Image still works but is slightly larger than alpine-based would be. |

---

## C) NOT STARTED

| #   | Item                                                                 |
| --- | -------------------------------------------------------------------- |
| 1   | Push commits to remote (`git push`)                                  |
| 2   | DNS propagation — Namecheap API IP whitelist → terraform apply       |
| 3   | SSL cert — Firebase cert `CERT_PENDING`, blocked on DNS              |
| 4   | Custom domain mapping for Cloud Run demo URL                         |
| 5   | Firebase Blaze plan upgrade (for Hosting rewrite)                    |
| 6   | Website Lighthouse audit                                             |
| 7   | Per-component documentation pages                                    |
| 8   | Website analytics                                                    |
| 9   | Google Search Console submission                                     |
| 10  | OG image visual verification (generated but never fetched/inspected) |

---

## D) WHAT WENT WRONG / TOTALLY Fucked Up

| #   | Issue                                                                  | Impact                                       | Root Cause                                                                                                                                                                                                                                                                                                                                           |
| --- | ---------------------------------------------------------------------- | -------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Firebase Hosting rewrite → Cloud Run returned 404**                  | Wasted ~20 min debugging                     | Firebase **free tier (Spark)** doesn't support Cloud Run rewrites. I assumed it would work because the project has GCP billing enabled. Firebase plan ≠ GCP billing. Should have checked the plan first.                                                                                                                                             |
| 2   | **Deployed to europe-west1 first, then had to redo in us-central1**    | Wasted ~15 min                               | Firebase Hosting rewrites require same-region Cloud Run, and the convention in this project is actually flexible. But more importantly, I created resources in the wrong region, had to delete them, recreate in us-central1, re-tag the Docker image, re-push. Should have checked the Firebase Hosting rewrite region constraint before deploying. |
| 3   | **Used `golang:1.26-alpine` in Dockerfile initially**                  | Docker build would have failed in CI         | Alpine doesn't include `git`, which `go mod download` needs for some dependencies. Fixed to `golang:1.26` (Debian full), but this makes the builder stage larger. Should have known `go mod download` needs git.                                                                                                                                     |
| 4   | **Restricted Cloud Run ingress (`internal-and-cloud-load-balancing`)** | Service returned 404 to all external traffic | Firebase Hosting rewrite was supposed to go through Cloud Load Balancing, but since the rewrite didn't work (Spark plan), the service was unreachable. Had to change to `--ingress=all`. This is now publicly accessible without the Firebase Hosting proxy — meaning the demo URL bypasses any CDN caching Firebase provides.                       |
| 5   | **Cloud Run URL is ugly**                                              | Bad UX                                       | `templcomponents-demo-132045829579.us-central1.run.app` is a terrible URL for a public demo. Without Firebase Blaze (for rewrite) or custom domain mapping on Cloud Run, there's no clean URL.                                                                                                                                                       |
| 6   | **Port 8080 hardcoded initially**                                      | Would have failed on Cloud Run               | Cloud Run injects `PORT=8080` but I had it hardcoded. Fixed to read from env. Would have silently worked in this case, but fragile if Cloud Run ever changes the default.                                                                                                                                                                            |

---

## E) WHAT WE SHOULD IMPROVE

| #   | Area                         | Current State                                                              | Target                                                                                              |
| --- | ---------------------------- | -------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| 1   | **Push to remote**           | 10 commits local-only                                                      | Push immediately — untested CI is a liability                                                       |
| 2   | **Demo URL**                 | `templcomponents-demo-132045829579.us-central1.run.app`                    | Map a custom domain or use Firebase Blaze for clean `/demo/` path                                   |
| 3   | **Demo cold start**          | Cloud Run min-instances=0 means first visitor waits ~2-3s for cold start   | Either set min=1 (costs ~$5/mo) or add a loading indicator on the website link                      |
| 4   | **Dockerfile builder image** | `golang:1.26` (850MB Debian)                                               | Use `golang:1.26-bookworm-slim` or a two-stage approach with alpine + git installed                 |
| 5   | **Cloud Run region**         | `us-central1` (Iowa)                                                       | Consider `europe-west1` if audience is European — lower latency. But Firebase rewrite constraint... |
| 6   | **CI workflow testing**      | Never run                                                                  | Push and fix whatever breaks                                                                        |
| 7   | **Cost monitoring**          | No budget alerts                                                           | Set a GCP budget alert on the lars-software project for Cloud Run + Artifact Registry               |
| 8   | **Docker image versioning**  | `:latest` only                                                             | Use `:sha-<git-sha>` or `:v<version>` for rollback capability                                       |
| 9   | **Pre-render fallback**      | Code exists (`prerender.go`) but unused now that Cloud Run serves the demo | Either delete pre-render code or document it as an alternative deployment strategy                  |
| 10  | **Health check endpoint**    | Cloud Run uses TCP probe on :8080                                          | Add `/health` endpoint returning 200 for explicit liveness/readiness                                |

---

## F) NEXT 50 THINGS TO DO

### Critical — Do before declaring "public"

| #   | Task                                          | Est.   | Why                                                 |
| --- | --------------------------------------------- | ------ | --------------------------------------------------- |
| 1   | `git push origin master`                      | 1 min  | Nothing is live on GitHub. All work is local.       |
| 2   | Watch CI run, fix any failures                | 10 min | Workflows are completely untested.                  |
| 3   | Watch website workflow run, fix failures      | 10 min | Docker build + Cloud Run deploy never tested in CI. |
| 4   | Whitelist Namecheap API IP, `terraform apply` | 10 min | Custom domain is the #1 credibility signal.         |
| 5   | Verify DNS + SSL + custom domain loads        | 10 min | End-to-end verification.                            |

### High — Meaningful improvement to first impression

| #   | Task                                                          | Est.   | Why                                 |
| --- | ------------------------------------------------------------- | ------ | ----------------------------------- |
| 6   | Map custom domain to Cloud Run (or upgrade Firebase to Blaze) | 15 min | Cloud Run URL is ugly.              |
| 7   | Add cold-start warning or min-instances=1                     | 5 min  | First visitor gets 2-3s blank page. |
| 8   | Verify OG images render correctly (fetch `home.png`)          | 5 min  | Never visually inspected.           |
| 9   | Add `/health` endpoint to demo binary                         | 5 min  | Proper Cloud Run health checks.     |
| 10  | Set GCP budget alert ($5/mo)                                  | 5 min  | Avoid surprise costs.               |
| 11  | Add Docker image tagging with git SHA in CI                   | 10 min | Rollback capability.                |
| 12  | Submit sitemap to Google Search Console                       | 5 min  | SEO.                                |
| 13  | Add demo link to Starlight docs sidebar                       | 5 min  | Discoverability.                    |
| 14  | Add "Edit on GitHub" links to docs pages                      | 10 min | Community contribution friction.    |

### Medium — Polish and completeness

| #   | Task                                                            | Est.   |
| --- | --------------------------------------------------------------- | ------ |
| 15  | Optimize Dockerfile (slim builder, `.dockerignore`)             | 15 min |
| 16  | Add `COPY go.mod go.sum` cache layer (already there but verify) | 5 min  |
| 17  | Add Lighthouse CI to website workflow                           | 20 min |
| 18  | Make HTML validation a blocking CI step                         | 5 min  |
| 19  | Add broken link checker to website CI                           | 15 min |
| 20  | Add CSP headers to Cloud Run demo responses                     | 10 min |
| 21  | Rename `forms_demo.templ` → `forms_standalone.templ`            | 5 min  |
| 22  | Delete unused `prerender.go` or document it as alternative      | 5 min  |
| 23  | Add reading time + last-updated to doc pages                    | 15 min |
| 24  | Add table of contents to long doc pages                         | 15 min |
| 25  | Add prev/next nav to doc pages                                  | 10 min |
| 26  | Add search to demo page (filter by component name)              | 30 min |
| 27  | Add dark mode toggle to demo (not just system)                  | 10 min |
| 28  | Add browser support matrix to docs                              | 15 min |
| 29  | Write migration guide from templUI                              | 30 min |
| 30  | Add "Copy to clipboard" on all code blocks                      | 15 min |
| 31  | Add structured data to demo page                                | 10 min |
| 32  | Write announcement blog post                                    | 30 min |

### Future — Nice to have

| #   | Task                                                     | Est.    |
| --- | -------------------------------------------------------- | ------- |
| 33  | Add per-component documentation pages (top 10 first)     | 2 hours |
| 34  | Interactive playground (edit templ, see result)          | 1 day   |
| 35  | Theme customizer (pick colors, see all components)       | 4 hours |
| 36  | Versioned docs (v0.x, v1.x)                              | 4 hours |
| 37  | Add analytics (Plausible/Fathom)                         | 15 min  |
| 38  | Newsletter signup                                        | 30 min  |
| 39  | Sponsors section                                         | 15 min  |
| 40  | Contributing guide for new components                    | 30 min  |
| 41  | Community link (Discord/Slack)                           | 10 min  |
| 42  | Visual regression testing for demo                       | 1 hour  |
| 43  | Bundle size monitoring                                   | 15 min  |
| 44  | Cache-control headers for Cloud Run responses            | 10 min  |
| 45  | Add `Cache-Control: public, max-age=300` to demo HTML    | 5 min   |
| 46  | Add ETag support to demo server                          | 15 min  |
| 47  | Multi-region Cloud Run (us + eu)                         | 30 min  |
| 48  | Cloud Run service-level logging to BigQuery              | 20 min  |
| 49  | Uptime monitoring (UptimeRobot or equivalent)            | 10 min  |
| 50  | Demo page performance optimization (inline critical CSS) | 1 hour  |

---

## G) TOP 2 QUESTIONS

### 1. Should I push now, or is there something else to do first?

10 commits are sitting local. CI, website deploy, and Cloud Run deploy workflows have **never run**. Every minute we don't push, the remote falls further behind and the window for CI surprises grows. But if there's more work to batch, we could avoid multiple push/CI cycles.

**Do you want me to push now, or are there more changes to make first?**

### 2. Firebase Blaze or Cloud Run custom domain?

The demo URL (`templcomponents-demo-132045829579.us-central1.run.app`) is unusable for a public library. Two options:

- **A) Upgrade Firebase to Blaze (pay-as-you-go):** Enables Hosting rewrite so `/demo/` proxies to Cloud Run on the main domain. Near-zero cost at current traffic. Requires enabling billing on the Firebase project.
- **B) Map custom domain directly to Cloud Run:** e.g., `demo.templcomponents.lars.software` → Cloud Run. GCP handles SSL automatically. Doesn't require Firebase Blaze. But it's a different subdomain, not a path on the main domain.

**Which approach do you prefer?**
