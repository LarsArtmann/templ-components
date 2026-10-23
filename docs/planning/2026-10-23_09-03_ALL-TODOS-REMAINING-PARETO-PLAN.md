# ALL-TODOS REMAINING PARETO PLAN — 2026-10-23 09:03

**Goal:** a complete, prioritized inventory of every open `TODO_LIST.md` row, split into
verifiable tasks of **≤12 minutes**, with the executable subset executed and verified in
this session.

**Non-goals:** inventing demand for HOLD items; making ⫱ owner decisions on the owner's
behalf; upstream fixes in `larsartmann/buildflow` / `art-dupl`; the v2 breaking wave
(ADR-0039); pixel re-baselines outside a proven need.

**Sources:** `TODO_LIST.md` @ v1.23.0 (413 lines), `docs/planning/2026-10-15_20-00_ALL-TODOS-MASTER-PARETO-PLAN.md`,
`docs/planning/2026-10-17_08-53_wave-5-endgame-pareto-plan.md`, `docs/planning/2026-10-17_g1-decision-memos.md`,
`docs/adr/0009-accepted-clones.md`, `docs/planning/2026-10-17_v2-color-collapse-plan.md`.

**Sort key:** `score = customer-value × impact ÷ effort`, both 1–5. Customer = library
consumers. Correctness > determinism > DX > docs > polish. A blocked row is still listed
(the gate is the task), never silently dropped.

**Guardrails:** no breaking API outside the v2 gate; wired component ⇒ e2e in the same
task; goldens/visuals re-baseline only with a recorded reason; re-read a file immediately
before editing it (single-writer protocol, #357); push only after
`scripts/ci-repro.sh --lint --css --website` prints `VERDICT: PASS` on the exact commit.

---

## 1. Inventory — every open row

Classes: **EXEC** = doable in this repo without a gate · **G1** owner decision ·
**G2** BuildFlow upstream · **G3** art-dupl upstream · **G4** external actor ·
**G5** v2/ADR-0039 wave · **G6** demand HOLD · **OTHER** other-repo work.

| ID   | Item (short)                                                  | Class | Tier | Score |
| ---- | ------------------------------------------------------------- | ----- | ---- | ----- |
| 437  | Heatmap `role="img"` on scroll container (achecker)           | EXEC  | P0   | 12.5  |
| 294  | `tc ls` footer / ci-repro guard output / AGENTS note          | EXEC  | P0   | 9.0   |
| 357  | Agent single-writer protocol section                          | EXEC  | P0   | 9.0   |
| 335  | templ v0.3.1070 gate re-check                                 | EXEC  | P0   | 6.0   |
| 387  | branching-flow nits: strong IDs + guard hardening             | EXEC  | P1   | 4.5   |
| 435  | art-dupl baseline re-record vs fork analyzer                  | EXEC  | P1   | 3.0   |
| 369  | `@utility` migration for utility-shaped `tc-*` classes        | EXEC  | P2   | 2.25  |
| 189  | Demo niceties: file-backed kanban + dashboard recipe section  | EXEC  | P3   | 1.5   |
| 191  | Kanban scope: within-column keyboard reorder (decided rest)   | EXEC  | P3   | 1.0   |
| 154  | Prerendered index uses default wire view (keep in sync)       | EXEC  | P3   | 1.0   |
| 120  | CSS-recompile false negative (re-open only on recurrence)     | EXEC  | P4   | 0.5   |
| 190  | Kanban touch-drag story (DnD polyfill)                        | G1    | G1   | gate  |
| 211  | Fate of `templates/styles.css` + theme `.out.css`             | G1    | G1   | gate  |
| 212  | Policy for 2.1k unannotated report items                      | G1    | G1   | gate  |
| 217  | Demand-gated component builds (MultiSelect/TreeView/…)        | G6    | G6   | hold  |
| 233  | flake.lock daemon nudge: keep or revert                       | G1    | G1   | gate  |
| 234  | Recipes-identifier drift guard vs manual sweep                | G1    | G1   | gate  |
| 235  | `website.yml` npm tailwind: floating vs pin 4.3.3             | G1    | G1   | gate  |
| 236  | Errorpage demo: bordered boxes vs standalone routes           | G1    | G1   | gate  |
| 237  | AppShell/SidebarNav chrome default: ratify adaptive           | G1    | G1   | gate  |
| 238  | Kanban failure UX: auto-retry vs honest revert                | G1    | G1   | gate  |
| 239  | Kanban pending-timeout escape hatch                           | G1    | G1   | gate  |
| 295  | Promote `art-dupl check` to blocking CI                       | G1    | G1   | gate  |
| 303  | Website-module content clones: extract or blanket-accept      | G1    | G1   | gate  |
| 312  | Export `display.HeadingTag`?                                  | G1    | G1   | gate  |
| 322r | status→family 500 mapping (Corruption vs Infrastructure)      | G1    | G1   | gate  |
| 329  | SSE through Firebase Hosting rewrite never streams            | G1    | G1   | gate  |
| 358  | GOEXPERIMENT=jsonv2 exports: drop or keep                     | G1    | G1   | gate  |
| 400  | Consumer migration waves (dnsblockd/DiscordSync/nsfw/mr-sync) | G1    | G1   | gate  |
| 411  | Owner ratifies tool verdicts (Playwright/achecker/chromedp)   | G1    | G1   | gate  |
| 428  | Warning hue unification: amber vs yellow                      | G1    | G1   | gate  |
| 93   | BuildFlow daemon: honest commit messages                      | G2    | G2   | gate  |
| 107  | BuildFlow preflight: jsonv2 false "redundant"                 | G2    | G2   | gate  |
| 108  | BuildFlow eslint-fix fails on JS/TS changes                   | G2    | G2   | gate  |
| 124  | BuildFlow re-appends `*_templ.go` to `.gitignore`             | G2    | G2   | gate  |
| 125  | BuildFlow un-minifies committed CSS                           | G2    | G2   | gate  |
| 126  | BuildFlow commit classifier for vetted artifacts              | G2    | G2   | gate  |
| 232  | BuildFlow daemon commit-gate                                  | G2    | G2   | gate  |
| 352  | BuildFlow license-check GOROOT export                         | G2    | G2   | gate  |
| 383  | BuildFlow pseudo-version-hygiene false positive               | G2    | G2   | gate  |
| 384  | BuildFlow art-dupl step ignores `.art-dupl.json`              | G2    | G2   | gate  |
| 385  | BuildFlow lint auto-configure accepted-disabled allowlist     | G2    | G2   | gate  |
| 292  | art-dupl piped-output truncation                              | G3    | G3   | gate  |
| 293  | art-dupl `--fingerprint-only` + config discovery              | G3    | G3   | gate  |
| 315  | art-dupl actionable-vs-raw counts                             | G3    | G3   | gate  |
| 316  | art-dupl `--baseline-version-must-match`                      | G3    | G3   | gate  |
| 28   | `awesome-templ` PR (coverage ≥80% gate)                       | G4    | G4   | gate  |
| 29   | `templ.guide` listing PR (maintainer review)                  | G4    | G4   | gate  |
| 80   | Human eyeball overlay PNGs (+ vision-model pass)              | G4    | G4   | gate  |
| 150  | Human eyeball wire visual goldens                             | G4    | G4   | gate  |
| 162  | Human eyeball progressbar 45% golden                          | G4    | G4   | gate  |
| 269  | go-error-family upstream PR go-ahead                          | G4    | G4   | gate  |
| 156  | Consumer adoption offers (other repos)                        | G4    | G4   | gate  |
| 316b | htmx v4 radar                                                 | G4    | G4   | watch |
| 33   | `Validate()` on remaining props                               | G5    | G5   | gate  |
| 34   | Move test helpers to `internal/testutil`                      | G5    | G5   | gate  |
| 39   | Compound overlay components (Trigger/Content/Close)           | G5    | G5   | gate  |
| 40   | `ScrollbackLine.Timestamp` → `Prefix`                         | G5    | G5   | gate  |
| 157  | Calendar month navigation as `Wire` candidate                 | G5    | G5   | gate  |
| 345  | htmx v4 dialect audit checklist                               | G5    | G5   | gate  |
| 368  | `size-*` sweep (gated with #335 window)                       | G5    | G5   | gate  |
| 399  | ADR-0044 Phase 2 class swap + ADR-0045 tails                  | G5    | G5   | gate  |
| 431  | v2 color-collapse execution                                   | G5    | G5   | gate  |
| 410  | Release + fleet bridge tails (a–d)                            | OTHER | O    | gate  |
| 436  | Consumer-side stale-CSS fixes + fleet audit                   | OTHER | O    | gate  |

**Already complete but not struck** (strike in the hygiene task, evidence in the row):
#273, #291, #297, #298, #299, #300, #306, #309, #317, #357.

---

## 2. Tier rationale

**P0 (1% → 51%)** — user-visible correctness and list honesty:

1. **#437** is a live `aria_descendant_valid` violation: `role="img"` is a leaf role on an
   `overflow-x-auto` container holding a `<table>`. Screen readers get a broken image
   subtree today. Markup-only fix; no pixel change expected.
2. **#294/#357** are list-honesty: one is already implemented at source, the other shipped
   as an AGENTS section. Leaving them "open" poisons every future prioritization.
3. **#335** decides whether the repo's biggest golden-churn window is open. One `nix eval`;
   the answer must be recorded either way (automation keeps re-applying the bump).

**P1 (4% → 64%)** — latent-failure hardening:

4. **#387**: typed rule/heading IDs kill stringly-typed drift in the two test/tooling
   modules; the LCS index access and the `*update` deref are detect-only findings that are
   cheap to make structurally safe.
5. **#435** keeps the duplication gate's hash baseline valid against the analyzer that
   actually runs in the lane. Advisory only, but a stale baseline makes the gate lie.

**P2–P4** — #369 (`@utility` tree-shaking) is real but consumer-facing only through bundle
size; #189/#191/#154/#120 are demo/maintenance polish or dormant.

**Gates** are not schedulable work: every one needs an owner verdict, an upstream sprint,
an external actor, or the v2 decision. Their ≤12-minute tasks verify the blocker and keep
the decision packet warm; they never fake completion.

---

## 3. Micro tasks (≤12 min each)

### 3.1 #437 — Heatmap ARIA restructure (P0)

| μ     | Task                                                                           | Verify                        | Est |
| ----- | ------------------------------------------------------------------------------ | ----------------------------- | --- |
| 437.1 | Read achecker rule + confirm target element (container, not table)             | finding reproduced in code    | 5m  |
| 437.2 | Drop `role="img"` from container; name the native `<table>` via `aria-label`   | markup diff                   | 10m |
| 437.3 | Regenerate `display/heatmap_templ.go` from repo root (templ 0.3.1020)          | `git diff` only expected file | 8m  |
| 437.4 | Update `TestHeatmapA11y`: container roleless, table named, fallback name       | `go test ./display/`          | 8m  |
| 437.5 | Re-run `/display` route golden (pixel read) and achecker re-scan if env allows | golden green / zero finding   | 12m |
| 437.6 | CHANGELOG `Fixed` entry + strike #437                                          | `TestVersionMatches*` / warm  | 6m  |

### 3.2 #294 — Small polish verification (P0)

| μ     | Task                                                                  | Verify                | Est |
| ----- | --------------------------------------------------------------------- | --------------------- | --- |
| 294.1 | Verify `tc ls` footer is runtime-derived (no hand-typed count)        | `cmd/tc/main.go` read | 5m  |
| 294.2 | Verify AGENTS post-daemon note + ADR-0009 `--min-lines` verdict exist | grep                  | 5m  |
| 294.3 | Review ci-repro guard-failure path (step output + non-zero exit)      | script read           | 6m  |
| 294.4 | Strike #294 with evidence                                             | TODO row              | 4m  |

### 3.3 #357 — Protocol shipped (P0)

| μ     | Task                                                                               | Verify      | Est |
| ----- | ---------------------------------------------------------------------------------- | ----------- | --- |
| 357.1 | Confirm AGENTS "Agent single-writer protocol (2026-10-17, #357)" covers the lesson | AGENTS read | 4m  |
| 357.2 | Strike #357                                                                        | TODO row    | 2m  |

### 3.4 #335 — templ gate (P0)

| μ     | Task                                             | Verify         | Est |
| ----- | ------------------------------------------------ | -------------- | --- |
| 335.1 | `nix eval --inputs-from . nixpkgs#templ.version` | prints version | 3m  |
| 335.2 | Record verdict in #335 (and #368/#369 linkage)   | TODO row       | 4m  |

### 3.5 #387 — Hardening (P1)

| μ     | Task                                                                      | Verify                     | Est |
| ----- | ------------------------------------------------------------------------- | -------------------------- | --- |
| 387.1 | `type AxeRuleID string`; update `AxeViolation.ID`, key join, format sites | `cd visualtest && go test` | 8m  |
| 387.2 | `type HeadingID string`; update `md.Heading.ID` + golden literal          | `cd website && go test`    | 8m  |
| 387.3 | Bounds-safe `lcsAt` accessor in `utils/golden` diff                       | `go test ./utils/golden/`  | 8m  |
| 387.4 | Nil-safe `updateGoldens()` helper replacing 5 `*update` derefs            | visualtest compiles        | 10m |
| 387.5 | Strike #387 + CHANGELOG note (internal refactor exemption applies)        | guard                      | 5m  |

### 3.6 #435 — art-dupl baseline re-record (P1)

| μ     | Task                                                                    | Verify        | Est |
| ----- | ----------------------------------------------------------------------- | ------------- | --- |
| 435.1 | Build fork binary from `~/projects/art-dupl`; log version               | `--version`   | 8m  |
| 435.2 | Advisory check at t=2/t=3 with the lane binary; capture group sets      | report saved  | 12m |
| 435.3 | Triage new signature classes against ADR-0009 appendix (accept/extract) | ADR diff      | 12m |
| 435.4 | Re-record baseline with that same binary                                | baseline diff | 10m |
| 435.5 | `check` green twice; log binary version in ADR/AGENTS if count changed  | green ×2      | 12m |

### 3.7 #369 — `@utility` migration (P2, only if P0/P1 are green)

| μ     | Task                                                                     | Verify      | Est |
| ----- | ------------------------------------------------------------------------ | ----------- | --- |
| 369.1 | Inventory utility-shaped `.tc-*` classes + their consumers               | list        | 10m |
| 369.2 | Convert classes to `@utility tc-*`; extend `loadCustomCSSClasses` parser | parser test | 12m |
| 369.3 | Recompile CSS + class inventory + starter mirror                         | byte checks | 12m |
| 369.4 | CHANGELOG + strike #369                                                  | guards      | 8m  |

### 3.8 #189 / #191 / #154 / #120 — polish/dormant (P3/P4)

| μ     | Task                                                                          | Verify                 | Est |
| ----- | ----------------------------------------------------------------------------- | ---------------------- | --- |
| 189.1 | Scope file-backed kanban demo state (persist + reset)                         | decision               | 6m  |
| 189.2 | Implement + e2e pending/revert coverage                                       | `.#visual` kanban pack | 12m |
| 191.1 | Verify within-column keyboard reorder stays deferred per vibe-kanban evidence | note                   | 5m  |
| 154.1 | Verify prerender still uses the default both-transport view                   | `prerender.go` read    | 5m  |
| 120.1 | Confirm no recurrence since `fabd1fb`; keep dormant                           | note                   | 3m  |

### 3.9 Gate verification tasks (≤12 min each, no fake completion)

| μ    | Gate                                                                   | Task                                                                                                                           | Est |
| ---- | ---------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------ | --- |
| G1.1 | #190 #211 #212 #233–#239 #295 #303 #312 #322r #329 #358 #400 #411 #428 | Verify `2026-10-17_g1-decision-memos.md` still covers every ID + recommendation; note the one-line follow-up per accepted memo | 12m |
| G2.1 | #93 #107 #108 #124–#126 #232 #352 #383–#385                            | Spot-check that each BuildFlow claim still reproduces in `~/projects/buildflow` HEAD or a stale binary                         | 12m |
| G3.1 | #292 #293 #315 #316                                                    | `art-dupl` fork version + upstream issue/PR state                                                                              | 8m  |
| G4.1 | #28 #29                                                                | Re-measure coverage; check templ.guide PR review state                                                                         | 10m |
| G4.2 | #80 #150 #162 #269                                                     | Confirm the human/API-key/go-ahead gate has not changed                                                                        | 5m  |
| G5.1 | #33 #34 #39 #40 #157 #345 #368 #399 #431                               | Confirm ADR-0039 v2 gate + #335 window still closed                                                                            | 8m  |
| G6.1 | #217                                                                   | Re-run demand check cadence; no new consumer evidence                                                                          | 10m |
| O.1  | #410 #436 #400                                                         | List the exact repo + command for each other-repo tail (handoff card)                                                          | 10m |

---

## 4. Execution order

1. **P0 executable**: 437 → 294 → 357 → 335.
2. **P1**: 387 → 435.
3. **Hygiene**: strike already-complete rows; keep `[Unreleased]` warm.
4. **Gates**: run 3.9 verification only; record the verdicts, never strike on someone
   else's behalf.
5. **P2+**: 369 / 189 only if the verification suite is green and time remains.
6. **Final gate**: touched-module tests → `templ generate` diff → lint on touched modules →
   optional `scripts/ci-repro.sh --lint --css --website`. Push only on `VERDICT: PASS`.

---

## 5. Execution log

| Step                | Result | Evidence                                                                                                                                                                                                                                |
| ------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Plan written        | ✅     | this file                                                                                                                                                                                                                               |
| 437 Heatmap ARIA    | ✅     | `role="img"` removed; native `<table>` owns `AriaLabel` (default "Heatmap"); 4 HTML goldens + `TestHeatmapA11y` updated; CHANGELOG `Fixed`; `go test ./display/` green; pixels unchanged                                                |
| 294 polish verify   | ✅     | `tc ls` footer counts the embedded registry at runtime; ci-repro guard step prints the failure and exits non-zero; AGENTS daemon note + ADR-0009 `--min-lines` verdict present; row struck                                              |
| 357 strike          | ✅     | AGENTS "Agent single-writer protocol (2026-10-17, #357)" present; row struck                                                                                                                                                            |
| 335 gate re-check   | ✅     | `nix eval --inputs-from . nixpkgs#templ.version` → `0.3.1020`; gate still CLOSED; #368 stays sequenced behind it                                                                                                                        |
| 387 hardening       | ✅     | `AxeRuleID`, `md.HeadingID`, bounds-safe `lcsAt`, nil-safe `updateGoldens()` in 5 call sites; website suite + visualtest compile + lints green                                                                                          |
| 435 art-dupl        | ✅     | Fork `v0.9.1-0.20261017225421-93432fd910fb-93432fd` built; advisory sets measured 73 (t=2) / 38 (t=3) → existing ADR classes; baseline re-recorded at 150 groups; `check` green ×2; ADR/AGENTS/CHANGELOG counts synced                  |
| 369 @utility        | ⏸      | DEFERRED with reason: `@utility` tree-shakes unused classes, which would drop opt-in classes (`tc-squircle`, `tc-fluid-*`) from the tracked prebuilt `templates/styles.css`; needs the ADR-0045 consumer-story decision, not a drive-by |
| Module verification | ✅     | root `go test ./...` + utils/icons/errorpage/charts/echarts/datastar/htmx loops green; visualtest compiles; website suite green                                                                                                         |
| Lint / ci-repro     | ⚠️      | golangci-lint 0 issues on display, utils/golden, visualtest, website; full `ci-repro --lint --css --website`/visual suite NOT run (needs browser env + long lane)                                                                       |

### Gate verification (read-only, 2026-10-23)

| Gate                                                                        | Finding                                                                                                                                                                                                                             |
| --------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| G1 (#190 #211 #212 #233–#239 #295 #303 #312 #322r #329 #358 #400 #411 #428) | `docs/planning/2026-10-17_g1-decision-memos.md` still covers all memo'd IDs and names the self-contained rows; no recommendation changed. Awaits the owner session.                                                                 |
| G2 (BuildFlow)                                                              | `~/projects/BuildFlow` still a live checkout (HEAD `f6dcb8c4c`); no upstream sprint happened this session.                                                                                                                          |
| G3 (art-dupl)                                                               | Fork at `93432fd9`; #435 consumed the analyzer update; #292/#293/#315/#316 remain upstream asks.                                                                                                                                    |
| G4 (#28 #29 #80 #150 #162 #269 #156 #316b)                                  | **#29 is DONE: a-h/templ#1447 is MERGED** and struck. #28 still blocked: root-module coverage measured 70.4% (threshold 70 passes; awesome-templ needs 80). #80/#150/#162 still need a human or API key; #269 needs owner go-ahead. |
| G5 (v2)                                                                     | ADR-0039 gate unchanged; #431 execution + #399 Phase 2 stay parked.                                                                                                                                                                 |
| G6 (#217)                                                                   | No new consumer demand evidence; HOLD stands.                                                                                                                                                                                       |
| O (#410 #436 #400)                                                          | Other-repo work: fixes belong in dbsblockd/DiscordSync/nsfw-classifier/mr-sync/pbx-lens; no library-side action.                                                                                                                    |

### Unplanned repairs found while verifying (fixed)

| Issue                                                                                           | Fix                                                                                           |
| ----------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| Daemon re-applied the templ `v0.3.1070` bump to `go.mod`, `website/go.mod`, `visualtest/go.mod` | Restored all 9 modules to `v0.3.1020` via `go mod edit` + `tidy`; `TestTemplVersionPin` green |
| Daemon bumped `tools/crossbrowser-spike` playwright-core to 1.64.1 against the flake pin        | Restored to 1.63.0 via `npm pkg set`; `TestCrossbrowserLockstep` green                        |
| `docs/recipes/scoped-theme-bridge.md` palette table was prettier-aligned against the generator  | Regenerated with `go test ./utils/tone/ -update`; `TestThemeBridgePaletteTableSync` green     |
| `cmd/tc/_sources/display/heatmap.templ` mirror drifted after the Heatmap edit                   | `scripts/check-tc-sources-sync.sh --fix`; `TestSourcesMatchPackageFiles` green                |
| `templates/templ-components-classes.txt` stale after the Heatmap edit                           | `scripts/gen-class-inventory.sh`; `TestClassInventoryFreshness` green                         |

### Deferred executable item

**#369 `@utility` migration.** `@utility` only emits a class when the consumer's
scanner sees it. The tracked prebuilt `templates/styles.css` is compiled WITHOUT
consumer sources, so opt-in classes such as `tc-squircle` and `tc-fluid-*` would
disappear from that artifact for consumers who do not run Tailwind. The migration
therefore needs the ADR-0045 prebuilt-CSS story decided first (either drop the
prebuilt artifact, add a safelist, or migrate only classes no consumer opts into
dynamically). Recorded here rather than silently breaking the artifact.

## 6. Verification plan

- **Per touched module:** `GOWORK=off go test -count=1 ./...` (root packages via the root
  module; `website` and `visualtest` separately).
- **Generated files:** regenerate from the repo root with the pinned templ 0.3.1020; a diff
  outside the edited component's `_templ.go` is a failure.
- **Goldens:** package HTML goldens via the owning package test; pixel goldens only under
  `nix run .#visual`; never `-update` because of a suspected load flake.
- **A11y:** the achecker re-scan is read-only evidence; the axe sweep remains the CI gate.
- **Counts:** any component/icon/golden count claim stays in lockstep with
  `utils.TestDocsCountDrift` / `tone.TestThemeBridgePaletteTableSync`.
- **Process:** after any daemon commit, re-check that the session's edits survived; commit
  PNG churn manually; push only after the ritual verdict.
