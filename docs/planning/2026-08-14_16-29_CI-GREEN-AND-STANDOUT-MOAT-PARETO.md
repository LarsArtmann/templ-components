# Pareto Plan — CI Green & Standout Moat

**Created:** 2026-08-14 16:29
**Status:** ACTIVE — Tier 1 executing now
**Sources:** CI run logs (31809023232, 31809023336), `TODO_LIST.md`, `docs/STANDOUT-IDEAS.md`, SKILL.md catalogue, working-tree audit 2026-08-14

---

## Why this plan exists

The last green CI run on `master` was **2026-08-10**. Since then **8+ consecutive runs failed**
across 5 distinct failure classes while feature commits kept landing. The BuildFlow daemon
commits without running `go test` (TODO #93), so nothing stopped the drift. Separately, the
library's biggest competitive gap (no per-component docs on the website — 112 components,
3 doc pages) and its genuine moats (HTMX-first package, ecosystem play) are untouched.

**Evidence snapshot (2026-08-14):**

| # | Failure | Where | Root cause |
|---|---------|-------|------------|
| F1 | `TestGoldenCopyButton`, `TestGoldenCardBodySlot` | CI Build & Test | Commit `177d88e` changed fixtures `npm install` → `pnpm add` without regenerating goldens |
| F2 | CSS Freshness: "Committed CSS is stale" | CI CSS job | `examples/demo/static/app.css` not recompiled after class changes (`nix run .#css` fixes) |
| F3 | Actionlint SC2044 | CI Lint job | `ci.yaml:85` — `for templ_file in $(find ...)` fragile loop |
| F4 | Website: "Unable to locate executable file: pnpm" | Website workflow | `setup-node` does not install pnpm; also `cache-dependency-path` points to npm `package-lock.json` that no longer exists after pnpm migration |
| F5 | Visual Regression exits 1 (~76s in) | CI Visual job | `nix run .#visual` fails inside `OUTPUT=$(...)` under `set -e` — output swallowed; first red run coincides with the CircularProgress feature (Aug 11) — likely PNG mismatches or missing snapshots; local repro running |
| F6 | (historic) "Verify no untracked changes" | 2 older runs | BuildFlow daemon commits stale trees — external, TODO #93 |

---

## Pareto breakdown

### The 1% that delivers 51% — **MAKE CI GREEN**

Everything else is worthless while the badge is red: consumers see a broken library, PRs
can't be trusted, and every subsequent task needs a green pipeline to verify itself.

1. Regenerate stale display goldens → `go test ./...` green locally
2. Recompile demo CSS (`nix run .#css`) → CSS Freshness green
3. Fix `ci.yaml` SC2044 (`find -print0` + `while read`) → Actionlint green
4. Fix `website.yml` pnpm setup (install pnpm, point cache at `pnpm-lock.yaml`) → Website green
5. Triage visual regression locally, fix snapshots or regressions, and stop swallowing `nix run` output in CI → Visual green
6. Watch a full CI run to completion on the pushed fix commit — **the only accepted proof**

### The 4% that delivers 64% — **DISCOVERABILITY + ONBOARDING**

These convert the existing (excellent) library into something a stranger can adopt:

7. Website per-component docs: generator (manifest → MDX) + `display`/`forms` pages
8. `GOEXPERIMENT=jsonv2` onboarding callout (README + installation.mdx) — my own first `go test` in this session failed without it
9. Datastar scope decision (ADR-0035: minimal complement, no parity pursuit — document it)
10. Ecosystem refresh: `STANDOUT-IDEAS.md` says "69 components, v0.3.0" (it's 112, v1.8.x); README GOTH-stack cross-links (cqrs-htmx, go-cqrs-lite)

### The 20% that delivers 80% — **QUALITY MOAT**

11. JS behavior test harness (chromedp key dispatch in `visualtest` module) + keyboard-nav tests for Tabs/Carousel/Dropdown/ContextMenu/Combobox — 17 hand-rolled singletons currently verified only by humans
12. Remaining docs pages (feedback/layout/navigation/htmx/datastar/errorpage/icons)
13. JS singleton consolidation: shared attach-once helper across 17 components
14. Compound overlay pattern (ADR-0023) — v2.0 epic, now planned in detail
15. Per-component docs ↔ live demo embed (examples/demo is deployed to Cloud Run)

### The other 20% (long tail to 100%)

16. `internal/testutil/` migration (TODO #34, 70+ test files)
17. `Validate()` assessment on remaining props (TODO #33)
18. Upstream listings: awesome-templ + templ.guide PRs (TODO #28/#29, blocked on maintainers)
19. Human PNG eyeball checklist (TODO #80, needs Lars)
20. BuildFlow external fixes: honest commit messages + `.gitignore` re-append bug (TODO #93, separate repo)
21. Version sync audit (root `1.8.1` vs sub-modules `1.8.2`, TODO_LIST `1.8.0`) + v1.9.0 release cut

---

## Comprehensive plan (tasks 30–100 min, sorted by impact / effort / customer value)

| ID   | Task                                                                                         | Tier | Impact | Effort | Customer value | Depends on |
| ---- | -------------------------------------------------------------------------------------------- | ---- | ------ | ------ | -------------- | ---------- |
| M1   | Make CI green: goldens + demo CSS recompile + verify locally                                  | 1%   | 10     | 45m    | Trust in library | —          |
| M2   | Fix workflow lint: `ci.yaml` SC2044 find-loop, `website.yml` pnpm install + lockfile path     | 1%   | 9      | 30m    | CI green, cache works | —      |
| M3   | Visual regression triage: local repro, fix PNG mismatches/missing snapshots                   | 1%   | 9      | 60m    | Layout regression safety | M1 (green suite) |
| M4   | CI visual job output hygiene: `tee` output, don't swallow `nix run` stderr under `set -e`     | 1%   | 6      | 30m    | Debuggable CI   | M3         |
| M5   | GOEXPERIMENT=jsonv2 onboarding: README callout + website `installation.mdx` troubleshooting   | 4%   | 7      | 30m    | First-run success | —        |
| M6   | Datastar scope ADR-0035 + website integration guide (minimal complement, revisit triggers)    | 4%   | 7      | 45m    | Clear adoption story | —       |
| M7   | Ecosystem refresh: STANDOUT-IDEAS.md stats, TODO_LIST version drift, README GOTH cross-links  | 4%   | 7      | 45m    | Discoverability | —          |
| M8   | Docs generator: component manifest (name/signature/one-liner/package) → MDX + sidebar wiring  | 4%   | 10     | 100m   | The #1 competitive gap | —    |
| M9   | Website docs pages: `display` (38 components)                                                 | 4%   | 10     | 100m   | " — "          | M8         |
| M10  | Website docs pages: `forms` (21 components)                                                    | 4%   | 9      | 100m   | " — "          | M8         |
| M11  | Website docs pages: `feedback` + `layout` + `navigation` (35 components)                       | 20%  | 8      | 100m   | " — "          | M8         |
| M12  | Website docs pages: `htmx` + `datastar` + `errorpage` + `icons` + `recipes`                    | 20%  | 7      | 90m    | " — "          | M8         |
| M13  | Docs copy-paste pattern: per-page runnable Go snippet + "Edit on GitHub"                       | 20%  | 7      | 60m    | shadcn-style DX | M9–M12 |
| M14  | JS behavior test harness: chromedp key-dispatch helper in `visualtest` module                  | 20%  | 8      | 90m    | Confidence in 17 JS singletons | — |
| M15  | JS tests: Tabs + Carousel keyboard nav (RTL included)                                         | 20%  | 7      | 60m    | WAI-ARIA proof | M14        |
| M16  | JS tests: Dropdown + ContextMenu shared menu nav                                              | 20%  | 7      | 60m    | " — "          | M14        |
| M17  | JS tests: Combobox + TagsInput                                                                | 20%  | 6      | 60m    | " — "          | M14        |
| M18  | JS singleton consolidation: shared attach-once helper, migrate display package                | 20%  | 6      | 90m    | Maintainability | M14       |
| M19  | JS singleton consolidation: migrate forms/layout/navigation                                    | 20%  | 5      | 60m    | " — "          | M18        |
| M20  | Compound overlays v2.0 part 1: Modal Trigger/Content/Close (ADR-0023)                         | 20%  | 7      | 100m   | Flexibility for real apps | M1 (green CI) |
| M21  | Compound overlays v2.0 part 2: Drawer + back-compat deprecations                              | 20%  | 7      | 100m   | " — "          | M20        |
| M22  | Docs ↔ demo embed: iframe/live-link per component page (demo runs on Cloud Run)               | 20%  | 6      | 60m    | See it before you `go get` | M9 |
| M23  | Upstream listings: awesome-templ + templ.guide PRs (external review)                          | tail | 8      | 30m    | Discoverability | M9 |
| M24  | Human PNG eyeball checklist doc + handoff (TODO #80)                                          | tail | 5      | 30m    | Visual QA honesty | M3 |
| M25  | `internal/testutil/` migration phase 1: shared render/assert helpers (TODO #34)               | tail | 4      | 60m    | Maintainability | — |
| M26  | `internal/testutil/` migration phases 2–3: golden + sweep helpers, 70+ files                  | tail | 4      | 60m    | " — "          | M25        |
| M27  | `Validate()` assessment: audit props where invalid states are representable (TODO #33)        | tail | 4      | 45m    | API honesty    | —          |
| M28  | Version sync audit (root 1.8.1 vs sub-modules 1.8.2 vs TODO_LIST 1.8.0) + v1.9.0 release      | tail | 6      | 60m    | Release hygiene | M1–M4 |
| M29  | BuildFlow external: honest commit messages + `.gitignore` re-append fix (TODO #93, own repo)  | blocked | 9    | 100m   | Stops the red-master class of rot | external |
| M30  | GOTH-stack example app epic (cqrs-htmx + templ-components + go-cqrs-lite, own repo)           | 20%  | 9      | multi-session | The single most convincing artifact | M9–M12 |

## Fine-grained breakdown (all tasks ≤ 12 min, sorted by execution order)

### Phase 0 — local verification baseline
| ID    | Micro-task                                                              | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G0.1  | `GOEXPERIMENT=jsonv2 go test ./display/ -run TestGolden -update`        | 2m | 2 failing tests now pass      |
| G0.2  | Full suite: `GOEXPERIMENT=jsonv2 go test ./...`                          | 5m | zero FAIL                     |
| G0.3  | Per-module isolation sweep (GOWORK=off loop from AGENTS.md)             | 8m | all 7 modules ok              |
| G0.4  | `nix run .#css` recompile demo CSS, inspect diff for sanity             | 8m | diff only expected classes    |
| G0.5  | Commit goldens + CSS: `fix(demo,display): regenerate stale goldens and demo CSS` | 3m | clean tree          |

### Phase 1 — workflow repairs (M2 + M4)
| ID    | Micro-task                                                              | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G1.1  | `ci.yaml:85`: rewrite find-loop as `find ... -print0 \| while IFS= read -r -d ''` | 5m | local shellcheck/actionlint clean |
| G1.2  | `website.yml`: add `corepack enable` + `pnpm/action-setup` (or corepack-managed pnpm) | 5m | `pnpm --version` resolves in step log |
| G1.3  | `website.yml`: `cache-dependency-path: website/pnpm-lock.yaml`          | 1m | no cache warning in CI        |
| G1.4  | `ci.yaml` visual job: `set +o pipefail`-safe capture — `nix run .#visual 2>&1 \| tee /tmp/vis.log`; grep the log; preserve exit code | 6m | failures show real output |
| G1.5  | Commit: `ci: fix actionlint SC2044, website pnpm bootstrap, visual output hygiene` | 2m | — |
| G1.6  | Push, `gh run watch` the CI + Website runs                              | 10m | both green or precise remaining failure |

### Phase 2 — visual regression triage (M3)
| ID    | Micro-task                                                              | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G2.1  | Read local `nix run .#visual` output (background job)                    | 2m | failure list                  |
| G2.2  | Classify each failure: real regression vs missing/stale PNG              | 5m | table of failures             |
| G2.3  | Real regressions → fix component code; stale PNGs → `-update` regen      | 12m | visual suite passes           |
| G2.4  | Commit + push, watch Visual job go green                                 | 5m | CI Visual green               |
| G2.5  | Write eyeball checklist for new/changed PNGs → `docs/visual-review-2026-08-14.md`, ping Lars (TODO #80) | 8m | checklist links each PNG |

### Phase 3 — onboarding + decisions (M5, M6, M7)
| ID    | Micro-task                                                              | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G3.1  | README: "Requirements" callout — `GOEXPERIMENT=jsonv2` until Go 1.27, exact error text shown otherwise | 6m | rendered README |
| G3.2  | `website/src/content/docs/getting-started/installation.mdx`: troubleshooting section with the same error | 6m | astro builds     |
| G3.3  | ADR-0035: Datastar = minimal opt-in complement; scope frozen at runtime/LiveRegion/Indicator/SSEErrorHandling; parity only on demonstrated demand | 10m | ADR file + index |
| G3.4  | `docs/STANDOUT-IDEAS.md`: refresh stats (112 components, 106 icons, v1.8.x), mark Tier-1 items done (demo deployed, v1.0 shipped) | 6m | no stale claims |
| G3.5  | TODO_LIST.md: bump Updated/Version, close out completed items            | 3m | dates consistent |
| G3.6  | README: GOTH-stack ecosystem section cross-linking cqrs-htmx + go-cqrs-lite | 6m | links valid     |
| G3.7  | Commit + push: `docs: onboarding callouts, ADR-0035 datastar scope, ecosystem refresh` | 2m | CI stays green  |

### Phase 4 — docs generator (M8) 
| ID    | Micro-task                                                              | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G4.1  | Define `docs/manifest.json` schema: package, component, signature, oneLiner, since | 8m | schema validates  |
| G4.2  | Seed manifest for `display` + `forms` from SKILL.md tables              | 12m | 59 entries        |
| G4.3  | Generator script (Go, `cmd/tc docs` or `scripts/gendocs`): manifest → MDX frontmatter + body skeleton | 12m | files emitted |
| G4.4  | Wire generated pages into Astro sidebar (`astro.config.mjs` sidebar)    | 8m | navigation lists all |
| G4.5  | Add one hand-written exemplar page (Button) with usage snippet as template | 10m | page renders     |
| G4.6  | Commit: `feat(website): component docs generator + display/forms skeletons` | 2m | — |

### Phase 5 — docs content fill (M9–M13, M22) — repeat per package
| ID    | Micro-task (per component page)                                         | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G5.1  | Fill page: props table (from Go doc), usage snippet, dark-mode note     | 10m | page builds       |
| G5.2  | Add live-demo link to Cloud Run demo route (iframe once demo routes exist) | 5m | link 200s     |
| G5.3  | Batch-commit per 5–8 components: `docs(website): Button, Badge, Card…`  | 2m | CI green          |
| G5.4  | Repeat G5.1–G5.3 for all 112 components (≈15 batches × 12m)             | ~180m total | all pages listed |

### Phase 6 — JS behavior tests (M14–M17)
| ID    | Micro-task                                                              | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G6.1  | `visualtest/interact.go`: chromedp helpers — `PressKey`, `Click`, `WaitVisible`, `AttributeOf` | 12m | compiles |
| G6.2  | Tabs test: arrow keys cycle focus, Home/End, `tabindex=0` invariant     | 12m | test passes  |
| G6.3  | Carousel test: arrows scroll track, dots sync, RTL key swap             | 12m | test passes  |
| G6.4  | Dropdown test: ArrowUp/Down skip disabled, Home/End, first-item focus on open | 12m | test passes |
| G6.5  | ContextMenu test: Shift+F10 opens at trigger, menu nav works            | 10m | test passes  |
| G6.6  | Combobox test: filter narrows options, Enter selects, hidden input syncs, Escape closes | 12m | test passes |
| G6.7  | Wire into `nix run .#visual` suite + CI                                 | 6m | CI green      |
| G6.8  | Commits per milestone: `test(visualtest): keyboard behavior for X`      | 2m×4 | — |

### Phase 7 — singleton consolidation (M18–M19)
| ID    | Micro-task                                                              | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G7.1  | Extract `tcAttachOnce(name, init)` shared JS emitter (display/shared.go) | 10m | goldens normalized |
| G7.2  | Migrate display singletons (carousel, ctx menu, tabs + shared.go)       | 12m | tests + goldens pass |
| G7.3  | Migrate forms/layout/navigation singletons (combobox, tags, theme, mobile menu) | 12m | tests pass |
| G7.4  | Regenerate goldens, full verify, commit                                 | 8m | suite green |

### Phase 8 — compound overlays v2.0 (M20–M21, only after release decision)
| ID    | Micro-task                                                              | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G8.1  | Design doc from ADR-0023: Modal.Trigger/Content/Close templ API + migration table | 12m | ADR updated |
| G8.2  | Implement Modal compound parts alongside monolithic (both work)         | 12m×4 | tests |
| G8.3  | Drawer compound parts + shared overlay internals                        | 12m×3 | tests |
| G8.4  | Deprecation notes + v2.0 migration guide page                           | 8m | docs build |

### Phase 9 — long tail (M23–M30)
| ID    | Micro-task                                                              | Est | Verify                        |
| ----- | ----------------------------------------------------------------------- | -- | ----------------------------- |
| G9.1  | awesome-templ PR: add repo entry (verify-before-filing: check listing criteria) | 12m | PR open |
| G9.2  | templ.guide PR/issue: submit library listing                            | 10m | issue open |
| G9.3  | Version sync: align root/sub-module/TODO_LIST versions, run drift-guard tests | 8m | tests pass |
| G9.4  | Cut v1.9.0 via `scripts/release.sh` once Tier 1+4 merged                | 12m | tag + CHANGELOG |
| G9.5  | testutil migration phase 1 (shared helpers into `internal/testutil/`, re-export shims) | 12m | suite green |
| G9.6  | testutil phases 2–3 (mechanical import rewrite, 70+ files)              | 12m×5 | suite green |
| G9.7  | `Validate()` audit list: which props have representable invalid states  | 10m | table in TODO_LIST |
| G9.8  | BuildFlow repo fixes (external): commit-message template from `git diff --stat`, drop `.gitignore` re-append | 100m | daemon commits honest |

---

## Mermaid execution graph

```mermaid
flowchart TD
    subgraph T1["1% — Make CI Green (51%)"]
        A1[G0.x Regenerate goldens + CSS] --> A2[G1.x Workflow fixes<br/>actionlint + pnpm + visual output]
        A2 --> A3[G1.6 Push + watch CI]
        A3 --> A4[G2.x Visual regression triage]
        A4 --> GREEN{{MASTER GREEN}}
    end

    subgraph T4["4% — Discoverability (64%)"]
        B1[G3.1-3.2 GOEXPERIMENT onboarding] --> GREEN2{{}}
        B2[G3.3 ADR-0035 Datastar scope] --> GREEN2
        B3[G3.4-3.6 Ecosystem refresh] --> GREEN2
        B4[G4.x Docs generator] --> B5[G5.1-5.4 display+forms pages]
        B5 --> B6[feedback/layout/navigation]
        B6 --> B7[htmx/datastar/errorpage/icons]
    end

    subgraph T20["20% — Quality Moat (80%)"]
        C1[G6.1 JS test harness] --> C2[G6.2-6.3 Tabs+Carousel]
        C1 --> C3[G6.4-6.5 Dropdown+ContextMenu]
        C1 --> C4[G6.6 Combobox+TagsInput]
        C5[G7.x Singleton consolidation] --> C2
        C6[G8.x Compound overlays v2.0] -.after v1.9.0. -> C7
        B7 --> C8[G5.2 demo embed per page]
    end

    subgraph TAIL["Other 20% — Long tail (100%)"]
        D1[G9.1-9.2 upstream listings]
        D2[G9.3-9.4 version sync + v1.9.0 release]
        D3[G9.5-9.6 testutil migration]
        D4[G9.7 Validate audit]
        D5[G9.8 BuildFlow external fixes — blocked]
        D6[M30 GOTH example app epic]
    end

    GREEN --> B1 & B2 & B3 & B4
    GREEN --> C1 & C5
    B5 --> D1
    C2 & C3 & C4 --> D2
    D2 --> C6
    D5 -.unblocks. -> A3
```

## Guardrails

- **No VERSCHLIMMBESSER:** every fix is verified by the actual CI run, not local assumptions. Goldens are only `-update`d when the *component output* is verified correct — never to silence an unexplained diff.
- Tier order is mandatory: nobody touches docs/tests while CI is red (that is how we got here).
- BuildFlow daemon may auto-commit mid-flight: `git status` before every commit; never stage unrelated changes.
- Visual `-update` regens require the human eyeball checklist (TODO #80) — pixel bugs AI cannot see.
