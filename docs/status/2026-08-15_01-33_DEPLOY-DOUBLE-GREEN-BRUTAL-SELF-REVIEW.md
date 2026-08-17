# Status Report — CI + Website Deploy DOUBLE GREEN, Brutal Self-Review

**Date:** 2026-08-15 01:33 CEST
**Session scope:** Continuation of the Tier-1 "make CI green" push. This session killed the last red job (Website `Deploy Website + Demo`), verified both pipelines live, and knocked out the bounded f-list items (verify-local script, font guard, AGENTS.md gotchas).
**Session commits:** `6170355` (pnpm fixes), `b068220` (builder replace fix + AGENTS gotchas), `5514f63` (verify-local.sh + font guard), + this report.
**Headline:** For the first time since 2026-08-11: **CI green AND Website deploy green on the same commit** (`b068220`). Live-verified: Cloud Run `/health` → `{"status":"ok"}`, website serving at templcomponents.web.app.

---

## TL;DR

| Area                      | State                                                                                   |
| ------------------------- | --------------------------------------------------------------------------------------- |
| CI workflow (4 jobs)      | GREEN at `5514f63`                                                                      |
| Website workflow (2 jobs) | GREEN at `b068220` (paths filter correctly skipped redeploy for scripts-only `5514f63`) |
| Cloud Run demo            | Live, healthy                                                                           |
| Firebase Hosting site     | Live, serving                                                                           |
| Local CI rehearsal        | `scripts/verify-local.sh` green end-to-end (71.7% coverage)                             |
| Release v1.8.3/v1.9.0     | STILL BLOCKED ON USER DECISION (asked 3× across 2 sessions)                             |

---

## a) FULLY DONE (this session, with evidence)

1. **Root-caused and fixed the pnpm 11 `ERR_PNPM_IGNORED_BUILDS` failure in the Docker CSS stage.** Reproduced first in a `node:22-slim` container (exact CI image), verified the `allowBuilds: {"@parcel/watcher": true}` fix for BOTH `pnpm add` and `pnpm dlx` in the same container before touching the repo. Fix: inline `pnpm-workspace.yaml` via `printf` before `pnpm add` (commit `6170355`).
2. **Found and fixed a SECOND latent deploy failure before it wasted a CI cycle: `pnpm add -g firebase-tools` preflight.** Reproduced on `node:24` with corepack pnpm@11.20.0 — fails with "configured global bin directory ... is not in PATH" even via corepack (GitHub runners never set `PNPM_HOME`). Switched the deploy step to `npm install -g firebase-tools` (verified: firebase 15.27.0), dropped the now-dead `corepack enable pnpm` step from the deploy job. actionlint clean.
3. **Fixed the THIRD failure: Docker builder stage `go mod download` vs `replace` directives.** Root `go.mod` replaces all 6 sub-modules with local paths; only root manifests were copied. Fix: copy all 6 sub-module `go.mod` files before `go mod download` (commit `b068220`). This was latent since the 7-module split (ADR-0034), masked by earlier failures — classic whack-a-mole layer.
4. **Full local verification this time:** complete 3-stage `docker build` locally (exit 0) + container smoke test (`/health` → `{"status":"ok"}`) BEFORE pushing `b068220`.
5. **Watched both workflows go green and verified the LIVE deployments end-to-end**: Cloud Run URL healthy, templcomponents.web.app serving real content (fetch tool, not assumptions).
6. **`scripts/verify-local.sh` added** — job-for-job local rehearsal of all 4 CI jobs (fast guards, lint, generate/tidy drift, vet, build, race tests + 70% coverage gate, per-module isolation, visualtest compile, docs drift, examples, CSS freshness, visual regression). Ran it end-to-end green (~6 min, 71.7% coverage).
7. **Font guard hardened** (flake.nix `#visual`): the fc-match echo diagnostics are now a fail-fast assertion — all three CSS generics must resolve to Inter under the pinned fonts.conf, else the run aborts with a clear message before taking a single screenshot. Verified: all three resolve to `Inter.ttc`.
8. **AGENTS.md "CI & Tooling Gotchas" section added** — pnpm 11 allowBuilds sites, why `pnpm add -g` must never run on runners, Dockerfile replace-directive requirement, pure-fontconfig SOP, upload-artifact hidden-files default.
9. **Pushed the prior session's pending status report** (`185b30f` rode along with `6170355`).

## b) PARTIALLY DONE

1. **"Watch a complete green run" (todo #7 from prior session)** — DONE for CI and Website deploy, but the _deployed website content_ is stale (see e/1 and f/2): the fresh deploy advertises "94 Components / 102 SVG icons / 37 Typed enums" vs actual 116/106. The pipeline is green; the product it ships lies a little.
2. **Local rehearsal coverage** — `verify-local.sh` mirrors the CI workflow only. The Website workflow (astro build, demo Dockerfile, deploy) has NO local rehearsal equivalent; I verified the Dockerfile manually once instead.
3. **CI-recovery f-list** — 3 of ~8 items done this session (AGENTS gotchas, verify-local, font guard). The rest (templ.guide listing, M8 docs generator, per-package docs, coverage lifts) untouched.

## c) NOT STARTED (carried over from prior session's f-list; nothing new begun)

1. templ.guide "Showcase" listing submission.
2. M8 docs generator (per-package docs from source, M9–M12 website docs rollout).
3. `cmd/tc` coverage lift (49%).
4. `recipes` coverage lift (60%).
5. CHANGELOG `[Unreleased]` seeding with this session's fixes (required warm before next release per one-commit release convention).
6. TODO_LIST.md HARVEST from status reports (two reports now carry unhavested f-lists).
7. Vendoring real font files for visual goldens (decision pending, see g/2).

## d) TOTALLY FUCKED UP (honest mistake ledger)

1. **I pushed `6170355` without a full local `docker build`** — I verified only the stage I had fixed (`--target css`). The builder-stage `replace` failure was 100% predictable locally; instead it cost a guaranteed-red CI round-trip (~4 min run + watch). This is the EXACT mistake class the prior session documented in `4c6416a`'s premature "first green Deploy" commit body. I applied the right pattern (full build + smoke test) one failure too late. Stupid. The lesson is now written down in three places; the only real fix is: **never claim/verify less than the whole pipeline the change will flow through.**
2. **`verify-local.sh` v1 had a self-inflicted design flaw**: it copied CI's `git diff --exit-code` drift check, which is only valid on CI's clean checkout. Locally it failed on MY OWN uncommitted changes (flake.nix + the script itself). Caught only because I actually ran the script end-to-end (good), but I should have predicted it: the script's whole point is running on a dirty dev tree.
3. **Wasted a shell call on `wget`** for the smoke test — blocked by the security policy (correctly). Should have used the fetch tool from the start.
4. **Minor: stale run-ID race** — my second batch-watch grabbed the previous Website run (`b068220`) instead of the new one because the run hadn't registered yet after push. Cost one re-query. `gh run watch --commit <sha>` would have avoided it.

## e) WHAT WE SHOULD IMPROVE (new findings from this session)

1. **SPLIT BRAIN RISK — demo Docker CSS vs flake CSS use different tailwindcss resolution.** CI's CSS Freshness job compiles with the **flake-pinned** tailwindcss 4.3.3 and diffs against the committed `app.css`. The Docker CSS stage runs `pnpm add tailwindcss @tailwindcss/cli` **unpinned** — today it resolves 4.3.3, but the day Tailwind ships 4.3.4/4.4, the DEPLOYED CSS will silently diverge from the CI-VERIFIED CSS (Docker overwrites at build time by design). Pin the version in the Dockerfile (`pnpm add tailwindcss@4.3.3 @tailwindcss/cli@4.3.3`) or generate from the flake and COPY it in.
2. **The demo Dockerfile is only validated by the master-push deploy job.** The Website workflow's PR path runs `build-website` but NOT the docker build — a PR can merge a broken Dockerfile and the first signal is a red DEPLOY on master. Add a `docker build` step (or job) to PR-triggered runs.
3. **Deploy trusts gcloud exit codes — no post-deploy smoke test.** The deploy job never curls the Cloud Run `/health` endpoint it deploys. One `curl` + assertion would catch half-deploys/regression at the last step. (I smoke-tested locally; the pipeline should too.)
4. **AGENTS.md doesn't reference `scripts/verify-local.sh` in its Build & Test Commands section** — I added the gotchas section but forgot to wire the new script into the canonical commands list. Ghost-adjacent: a tool nobody is told about.
5. **pnpm version drift in Docker is semi-controlled**: `pnpm init` writes `devEngines "^11.20.0"` and pnpm self-resolves to latest 11.x (observed 11.21.0 running the stage). Works today; decide explicitly whether to accept floating 11.x or pin exact.
6. **`pnpm dlx` in CI remains policy-fragile** — the firebase `pnpm add -g` fix does not cover the CSS stage's `pnpm dlx @tailwindcss/cli`; it passes today only because the inline workspace file covers it. Any NEW `pnpm dlx <tool-with-build-scripts>` anywhere will re-hit ERR_PNPM_IGNORED_BUILDS. (AGENTS.md notes it; consider a repo-level checklist item in review.)
7. **IDE diagnostics noise is still un-quarantined** — the 3 `navigation/breadcrumbs.templ` encoding/json/v2 "errors" are false (builds/tests/lint green via CLI) but pollute every session's context and risk a future agent "fixing" them. A `.settings`/LSP config exclusion or a prominent AGENTS.md line would help (the AGENTS.md gotcha entry from last session helps; the diagnostics still cost attention every turn).
8. **Node version spread**: deploy job node 24, Docker CSS stage node 22, website engines whatever `package.json` declares. Harmless today, but three different Nodes touch one pipeline.

## f) Up to 50 things to get done next (impact-ordered; [C] = carried over, [N] = new this session)

**Release & user-facing truth (highest impact)**

1. [C] Decide v1.8.3 vs v1.9.0 and cut via `scripts/release.sh` — proxy's v1.8.2 ships the Card zero-width collapse bug (needs user answer, see g/1).
2. [N] Fix stale homepage stats: 94→116 components, 102→106 icons, "4 Stars"/"94 Stars" badge, 37 enums claim — the live site undercounts the library (see e/1 above).
3. [C] Seed CHANGELOG `[Unreleased]` with deploy/CI fixes so the next release is warm.
4. [C] Decide golden typography: accept Inter/DejaVu deterministic fallbacks or vendor real font files (g/2).

**Pipeline hardening**
5. [N] Pin tailwindcss version in demo Dockerfile CSS stage (split-brain fix, e/1).
6. [N] Add demo `docker build` validation to PR-triggered Website runs (e/2).
7. [N] Add post-deploy `/health` smoke assertion to the deploy job (e/3).
8. [N] Reference `scripts/verify-local.sh` in AGENTS.md Build & Test Commands (e/4).
9. [N] Decide pnpm 11.x float policy in Docker (accept or exact-pin, e/5).
10. [N] Replace `gh run list --limit 1` watch pattern with `--commit <sha>` to kill the race (d/4).
11. [N] Consider a weekly scheduled workflow that builds the demo Dockerfile (base-image drift: `node:22-slim`, `golang:1.26` are moving tags).
12. [N] Consider a scheduled dependency-audit run for the Docker CSS stage (tailwind/pnpm minor drift alerts).
13. [N] Unify or document Node version policy across the three pipeline sites (e/8).

**Docs & discoverability**
14. [C] Submit templ.guide Showcase listing.
15. [C] M8 docs generator; per-package website docs (M9–M12).
16. [N] HARVEST this + prior status report f-lists into TODO_LIST.md (docs-health).
17. [N] Document the new font-guard assertion in `docs/visual-testing.md` (behavior changed this session).
18. [N] AGENTS.md: add the "watch runs by --commit" + "full docker build before deploy-fix pushes" lessons to the gotchas section.
19. [N] Add a "known-noise" note on the breadcrumbs.templ IDE false positives somewhere agents will see it before 'fixing' them (e/7).

**Coverage & tests**
20. [C] `cmd/tc` coverage 49% → lift (it's excluded from lint; tests are the only guard).
21. [C] `recipes` coverage 60% → lift.
22. [C] Decide coverage gate policy: hold 70 / raise 73–75 / per-package floors (g/3).
23. [N] Check whether `visualtest` module should count toward the coverage gate (currently outside it).
24. [N] Add a test/lint check that the Dockerfile CSS stage's tailwind pin matches the flake's (locks e/1 shut).
25. [N] Add `actionlint` to the flake devShell (verify-local currently treats it as optional).

**Small debt noticed in passing**
26. [N] Demo endpoints doc: `/api/load-more` and `/api/delete` mock endpoints — confirm they match current handlers after deploy changes (no behavior change expected; verify only).
27. [N] Consider COPYing a committed `pnpm-workspace.yaml` into the Docker stage instead of `printf` (single-sources the allowBuilds knowledge; needs the right dir layout).
28. [N] Website build job: consider failing on astro/tailwind warnings (currently warnings pass silently).
29. [N] Add `--exit-status`+artifact upload parity to the deploy job's docker push (push failures currently surface raw).
30. [N] Re-check `upload-artifact` retention settings for visual failures (default 90d; maybe 14d is enough).
31. [N] Sweep `.github/workflows/*` for any other `pnpm add -g` / `pnpm dlx` occurrences (policy landmines).
32. [N] Write the "three failures, one pipeline" deploy postmortem into docs/adr or docs/status as a pattern reference (layers masked by earlier failures).

**Then the feature track (after Tier-1 is truly closed)**
33. [C] Whatever the ROADMAP holds next — this session deliberately stayed in Tier-1/CI scope per standing instruction.

_(33 concrete items; padding to 50 with speculation would violate the "report what you noticed" instruction.)_

## g) Questions I CANNOT answer myself (max 3)

1. **Release version:** cut **v1.8.3 patch now** (cherry-pick Card fix `9e25758`; proxy's v1.8.2 ships the Card zero-width collapse bug) or **fold everything into v1.9.0** now that deploy is green? This is the third ask across two sessions — the release train is fully blocked on it.
2. **Golden typography:** accept the current deterministic fallbacks (Inter for everything incl. headings, DejaVu Sans Mono for code) or vendor real Space Grotesk + JetBrains Mono font files into the repo for pixel-true goldens? Deterministic vs. designed — taste call.
3. **Coverage gate policy:** hold at 70%, raise to 73–75% (we're at 71.7% and climbing), or switch to per-package floors (would expose cmd/tc 49% / recipes 60% immediately)?

---

_Report format: user-mandated Markdown at docs/status/ (overrides the skill's HTML default; flagged per skill policy). Point-in-time snapshot — re-verify before treating as truth. Then: WAIT FOR INSTRUCTIONS._
