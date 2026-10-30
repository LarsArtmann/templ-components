# Owner-Decision Packet — the gates that block templ-components work

**Created:** 2026-10-30 16:30 CET
**Author:** Crush session executing `docs/planning/2026-10-30_15-47_SUPERB-TOTAL-BACKLOG-PARETO-AND-EXECUTION-GRAPH.md` (C02)
**Purpose:** one memo carrying every owner-gated (⫱) row and open question, each with
options, evidence, blast radius, and a RECOMMENDED DEFAULT. Tasks gated on these rows
proceed on the default; an owner override later is a bounded re-run, not a redesign.

**How to answer:** reply per row with "accept default" or a different option letter.
Nothing here is irreversible except the tagged releases (#269's PR is branch-prepared).

---

## 1. The two blocking questions from the 15:40 report (§g)

### Q1 — Browser baseline target (gates C19, C22, C24 posture)

**Question:** which browser generation is the DEFAULT code path, with newer features as
opt-in layers vs the reverse?

| Option | Meaning | Consequence for the backlog |
| --- | --- | --- |
| **A. Widely-available 2024 (RECOMMENDED DEFAULT)** | `closedby="any"`, `::details-content`, popover, `interestfor`, anchor positioning, invoker commands STAY behind feature-detect/`@supports` fallbacks; the fallback path remains the shipped default for every component | C22/C23 ship the enhancement-layer shape (new platform path + JS/CSS fallback kept); C24 toasts keep the live region as primary; audit §5's warning honoured (Baseline claims age — do NOT promote on a guide's word) |
| B. 2025/2026-baseline-first | Newly-available features become the default path where the guide marks them Baseline; fallbacks demote to legacy engines only | Smaller payloads, but Firefox/Safari-before-N users get the degraded path; several audit "gap" rows (#473, #476…) become "skip" instead of "post-wave" |

**Evidence:** audit §5 (CLI per-browser tables contradict MDN on anchor positioning;
`interpolate-size`/`field-sizing`/`closedby`/invoker commands are all Newly-available);
the repo's own engine matrix (`docs/browser-support.md`) shows real Firefox/Safari
consumers; MWG guides themselves grade our fallbacks as the documented polyfill-equivalent.
**Blast radius of A:** none shipped today changes. **Unblocks:** C19 (theme model), C22
(HoverCard), C23 (tooltip layer), C24 (toasts), #473/#476–#482 disposition.
**Default applied in this execution:** A.

### Q2 — `LoadingOverlay` identity (gates C21)

**Question:** is the existing `LoadingOverlay` a page-blocking MODAL or a STATIC
in-place overlay?

| Option | Meaning | Consequence |
| --- | --- | --- |
| **A. Static by default, modal split out (RECOMMENDED DEFAULT)** | `LoadingOverlay` keeps its name and becomes an honest static region (`aria-busy`/`aria-live`, no `role="dialog"`); NEW `LoadingDialog` (on `overlayShell`) covers the page-blocking case; the demo's `Class: "absolute"` specimen migrates to the static variant, a modal specimen is added behind a trigger | The markup stops lying (today: `role="dialog" aria-modal` with the page still interactive behind it — the audit's "last legacy modal"); no consumer break for the demo shape |
| B. Modal by default | `LoadingOverlay` becomes `<dialog>` + `showModal()`; consumers wanting in-place spinners need a different component | Breaks the pinned demo specimen and every consumer using it inline; heavier migration |

**Evidence:** the demo pins `Class: "absolute"` inside a frame (in-place usage is the
shipped intent); the audit ranked this the weakest overlay (no focus containment, page
interactive behind a modal claim). **Blast radius:** A adds one component + rewords one;
B is a breaking change needing `docs/migration/`. **Unblocks:** C21, the `style=`
exemption list note in #327's guard. **Default applied in this execution:** A.

### Q3 — Push policy — RESOLVED

The 15:47 plan's instruction #8 said push after verification; pushes happen after a
`ci-repro VERDICT: PASS` at the exact tip (M03 ritual). No owner input outstanding.

---

## 2. The owner-gated (⫱) rows

Each row: evidence → options → recommendation → unblocks.

### #80 / #150 / #162 — Human eyeball of AI-captured goldens (vision-review first pass)

- **Evidence:** `scripts/vision-review-goldens.sh` exists and runs a vision-model first
  pass; it needs ONE exported API key (OPENAI/ANTHROPIC/GEMINI/OPENROUTER/XAI). Flagged
  sets: overlay PNGs (#80), wire dual-transport (#150), progressbar 45% fill (#162).
- **Options:** (a) export a key + run the script, human confirms only SUSPECT findings;
  (b) accept the 2026 pixel-golden suite as the eyeball (route+component goldens now
  number 331+53 and are human-reviewed at every re-baseline); (c) retire the rows.
- **Recommendation:** (a) once, at the next release eyeball window — the script makes
  the human cost small. Until then the rows stay ⫱ with no code impact.
- **Unblocks:** nothing downstream; pure risk-retirement.

### #211 — Fate of `templates/styles.css` + `templ-components-theme.out.css`

- **Evidence (complete since 2026-09-13):** no in-repo consumer outside `release.sh` and
  the guards; audit-session recommendation was DELETE.
- **Options:** (a) delete both + drop from `scripts/compiled-css-targets.txt` (RECOMMENDED);
  (b) keep + document a "download pre-built CSS" story in README.
- **Blast radius:** (a) touches the release script + `TestCompiledCSSInventory`; one
  commit. Consumers pulling the raw files from the repo would 404 — no evidence any does.
- **Unblocks:** #290-family cleanliness; shrinks the release artifact surface.

### #212 — Policy for ~2.1k unannotated name-keyed items in pre-2026-09-10 reports

- **Options:** (a) map IDs per file (M effort); (b) accept report-level resolution notes;
  (c) declare idea-residue open-by-design and stop tracking (RECOMMENDED — the reports
  are archived history; live work is tracked in TODO_LIST/FEATURES).
- **Unblocks:** closes the oldest meta-debt; no code impact.

### #233 — flake.lock daemon nudge: keep or revert

- **Evidence:** daemon moved nixpkgs 2026-09-17; dev shell verified unaffected
  (go 1.26.7 + templ v0.3.1020 — note: shell Go has since moved to 1.27.x via later
  deliberate bumps; templ pin unchanged).
- **Options:** (a) accept the new lock as fait accompli (RECOMMENDED — every later
  deliberate input bump has superseded it; reverting now rewrites history for nothing);
  (b) re-lock to `ef34387d…` + re-verify.
- **Unblocks:** closes a ⫱ that has been overtaken by events.

### #234 — Recipes-identifier audit: permanent CI drift-guard vs periodic sweep

- **Evidence:** the 2026-09-17 audit found 3 real drifts that shipped for weeks; guard
  shape is documented (extraction regex + allowlist). ~1–2h to build.
- **Options:** (a) build the guard (RECOMMENDED — same class as TestDocsCountDrift);
  (b) keep the manual periodic sweep with a calendar note.
- **Unblocks:** C25's docs-hygiene batch can fold the allowlist rules into the guard.

### #237 — AppShell/SidebarNav chrome default: adaptive vs classic-dark

- **Evidence:** adaptive shipped 2026-09-17 (white chrome in light mode); classic
  permanently-dark admin chrome stays a pure-CSS opt-out via `--tc-sidebar-*`;
  pixel goldens pin BOTH.
- **Options:** (a) ratify adaptive as default (RECOMMENDED — it is what shipped and what
  the goldens + docs now say); (b) flip default to classic-dark (breaking, needs
  migration note).
- **Unblocks:** row closes; ADR-0011 gets a "ratified" stamp.

### #238 / #239 — Kanban failure UX (auto-retry?) + pending-timeout escape hatch

- **Evidence:** current = immediate honest revert + manual re-drag; Datastar retries
  network errors natively (~2min) but not 4xx; no `KanbanBoardProps` opt-out exists.
- **Options:** (a) keep immediate revert, document as the contract (RECOMMENDED — an
  auto-retry on 4xx re-submits a mutation the server REJECTED; the pending register +
  revert is already the ADR-0041 behavior); (b) add one auto-retry; (c) add
  `PendingPolicy` enum now (speculative — no consumer ask).
- **Unblocks:** both rows close; the fallback API shape decision (#239) dies with (a).

### #269 — go-error-family upstream PR: `bridge.ClassifiedError.Message()`

- **Evidence:** branch `feat/bridge-classified-error-message` prepared in
  `/home/lars/projects/go-error-family` (red-green tests, lint clean, probes S1–S5).
- **Options:** (a) owner pushes + opens the PR (RECOMMENDED — work is done); (b) drop
  the branch and close the row.
- **Blast radius:** the PR kills the `[family]` message prefix — a consumer-visible
  string change in go-error-family; templ-components' errorpage does not depend on the
  prefix (verified by this repo's tests at the time of preparation).
- **Unblocks:** the errorpage message story stays as-is either way; this is upstream hygiene.

### #329 — SSE never streams through the Firebase Hosting `/demo/**` rewrite

- **Evidence (verified live 2026-10-01):** raw run.app streams at ~2.2s; canonical URL
  buffers forever; `X-Accel-Buffering: no` ignored by Firebase; the Datastar LiveRegion
  demo hangs on the canonical URL.
- **Options:** (a) map `demo.lars.software` (Cloud Run-direct) for SSE; (b) migrate demo
  to Firebase App Hosting; (c) env-var raw-run.app EventSource + CORS + CSP (RECOMMENDED
  as the CHEAPEST; (a) is the cleanest if a domain is available).
- **Unblocks:** the demo's Datastar page on the canonical URL; website docs can then
  demo LiveRegion honestly.

### #358 — Drop repo-wide `GOEXPERIMENT=jsonv2`

- **Evidence:** jsonv2 is stable WITHOUT the flag on the Go 1.27 floor; the flag is
  exported in `.envrc` + 3 flake places + pre-commit + CI env blocks + docs.
- **Options:** (a) one lockstep drop commit (RECOMMENDED — dead weight on every tool
  invocation; guards confirm behavior independently of the flag); (b) keep until the
  next release boundary.
- **Blast radius:** any toolchain < 1.27.1 consuming the repo loses the no-flag stability;
  the floor is 1.27 bare, so only patch-level old shells care.
- **Unblocks:** C18's toolchain-lockstep batch; removes a confusing env var from every
  probe (the "direnv contaminates ad-hoc probes" gotcha dies).

### #411 — Owner ratifies tool verdicts (Playwright / achecker / k6 / chromedp)

- **Evidence:** `docs/testing/tooling-evaluation.md` holds the verified-facts table;
  crossbrowser lane shipped findings-only + lockstep guard; achecker ledger at G2.
- **Options:** (a) ratify TRIAL dispositions as-is (RECOMMENDED — lanes are
  non-blocking by construction; nothing integrates further without this); (b) reject
  Playwright (kills the crossbrowser lane + its budget row).
- **Unblocks:** the tool-trial protocol's G1 gate; the budget table stops saying "pending".

---

## 3. What each answer unblocks (plan mapping)

| Answer | Unblocks |
| --- | --- |
| Q1 = A | C19, C22, C23, C24 execution as enhancement-layers; #473/#476–#482 stay post-wave |
| Q2 = A | C21 (split ships as: static keeps the name, `LoadingDialog` added) |
| #211 = a | release script + CSS inventory slimming |
| #234 = a | C25 folds allowlist into the guard |
| #237 = a | ADR-0011 ratified; SidebarNav docs stay |
| #238/#239 = a | Kanban failure contract documented; no new API |
| #269 = a | upstream PR opens (other repo) |
| #329 = c | demo Datastar page works on canonical URL |
| #358 = a | C18 lockstep batch; probe gotcha dies |
| #411 = a | tool-trial G1 closes; lanes stay findings-only |

**Defaults applied in this execution round:** Q1=A, Q2=A, #211 pending owner (not
executed — release-script change too close to C26), #234/#237/#238/#239/#358/#411
defaults NOTED in the relevant rows; only Q1/Q2 gate code work, and those defaults are
the conservative ones (everything stays an additive layer).

---

## 4. Decisions already answered by evidence (no owner input needed)

- **E2/#328** `display.CodeBlock` — SHIPPED (closed-by-verification 2026-10-30).
- **E7/#295** art-dupl promotion gate — two consecutive advisory greens exist; C15
  promotes the lane. The owner ratification is the CI job landing green; veto = revert
  one workflow file.
- **#428 warning hue** — evidence table lives in
  `docs/planning/2026-10-15_warning-hue-decision-memo.md`; default rec AMBER (G1
  pending). C20 executes with AMBER unless overridden (the memo's own recommendation).
