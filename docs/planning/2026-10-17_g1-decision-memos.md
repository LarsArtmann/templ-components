# G1 Decision Memos — Bulk Session Prep (2026-10-17)

**Purpose:** one-stop pre-read for the G1 owner-decision session. Thirteen ⫱G1 items,
each one page: decision, recommendation, facts, options, reopen condition. Recommendation
= the default that applies if the owner does not engage (deliberate-do-nothing is itself
a recorded verdict — no item here stays ambiguous).

**How to run the session:** read the summary table, challenge any recommendation,
then strike the TODO row per verdict. Final ratification of tool adoptions lands in
the same sitting (#411).

| #   | Decision                                            | Recommendation                         |
| --- | --------------------------------------------------- | -------------------------------------- |
| 233 | flake.lock daemon nudge: keep or revert             | REVERT to `ef34387d…`                  |
| 234 | recipes-identifier drift: CI guard vs manual sweep  | BUILD the guard                        |
| 235 | website.yml npm tailwind: floating vs pinned        | PIN 4.3.3 (F-20 lockstep)              |
| 236 | errorpage demo: bordered boxes vs standalone routes | STANDALONE `/errors/*` routes          |
| 237 | AppShell/SidebarNav chrome default                  | RATIFY adaptive (shipped)              |
| 238 | Kanban failure UX: auto-retry vs honest revert      | KEEP immediate revert                  |
| 239 | Kanban pending-timeout escape hatch                 | DOCUMENT budget, defer API             |
| 303 | website-module content clones                       | BLANKET-ACCEPT (ADR-0009)              |
| 312 | export `display.HeadingTag`?                        | HOLD until consumer demand             |
| 358 | GOEXPERIMENT=jsonv2 exports                         | DROP in one lockstep commit            |
| 411 | ratify tool verdicts (lanes → adopt/trial/reject)   | ADOPT-as-trial, chromedp stays primary |
| 428 | warning hue: amber vs yellow                        | AMBER (memo: 2026-10-15)               |
| 322 | status→family: 500→Corruption vs Infrastructure     | KEEP 500→Corruption                    |

Also gated on G1 but self-contained in TODO_LIST (no memo needed): #295 (art-dupl
blocking flip), #190 (touch-drag), #212 (annotation policy), #217 re-check cadence,
#329 (SSE hosting), #400 (consumer waves), and the v1.24.0 cut timing.

---

## #233 — flake.lock daemon nudge (`b1b875982b`): keep or revert

**Decision:** the daemon moved the nixpkgs flake input on 2026-09-17 — not a deliberate
bump. Keep the newer lock, or revert to `ef34387d…`.

**Recommendation: REVERT.** The repo's toolchain-isolation policy (dedicated
`nixpkgs-go` input, AGENTS.md) exists precisely so toolchain-bearing inputs move only
deliberately. A daemon-written lock normalizes accidental input drift.

**Facts:** dev shell verified unaffected at the time (templ v0.3.1020 unchanged);
functionally harmless today; revert = re-lock + re-verify the shell (minutes).

**Options:** (a) revert — policy-clean, small one-time cost; (b) keep — zero work,
but every future daemon nudge becomes precedent to keep.

**Reopen condition:** if `ef34387d…` lacks a package the shell needs, re-lock
deliberately instead and record the bump as intentional.

---

## #234 — recipes-identifier audit: permanent CI drift-guard vs manual sweep

**Decision:** the 2026-09-17 audit (#230) caught 3 real identifier drifts in recipes
that had shipped for unknown weeks. Guard permanently, or sweep periodically?

**Recommendation: BUILD the guard.** Every other silent-drift class in this repo is
machine-guarded (`TestDocsCountDrift`, `check-templ-sync.sh`, …); recipes are the last
unguarded doc surface with a proven drift history. A test-level guard beats a calendar
sweep because it cannot be forgotten across sessions.

**Facts:** guard shape already exists — extraction-regex lessons in
`docs/status/2026-09-17_09-28` §d2; allowlist defined (external go-datastar SDK
symbols, `Test*` guard names, underscore-suffixed constants). ~1–2h effort.

**Options:** (a) CI drift-guard (test in `recipes` or `internal/contract`); (b) manual
quarterly sweep with a named owner.

**Reopen condition:** if the extraction regex proves brittle against new recipe
styles (2+ false positives after allowlist tuning), downgrade to the sweep and
calendar it.

---

## #235 — npm tailwind version policy for `website.yml`

**Decision:** floating latest v4 (matches nothing) vs pinned, for the website CI lane.

**Recommendation: PIN (4.3.3).** F-20 already mandates lockstep across all tailwind
lanes (nix dev shell + demo Dockerfile + website.yml — currently all 4.3.3). Floating
latest in one lane silently breaks the lockstep invariant and the re-probe discipline
(`--alpha()` / slash-utility assumptions must be re-probed with the NEW binary).

**Facts:** the lockstep rule exists because a tailwind binary change flips generator
output and utility semantics; the demo Dockerfile precedent is pinned pnpm 4.3.3.

**Options:** (a) pin all three lanes (status quo, codified); (b) float all three
(synchronized bump burden, no pin drift, but non-reproducible site builds).

**Reopen condition:** only as part of a deliberate F-20 revision, never lane-locally.

---

## #236 — errorpage demo presentation: bordered boxes vs standalone routes

**Decision:** `ErrorPage`/`NotFound404` are full-page components (`min-h-screen`
`<main>`) but render inside bordered demo boxes — nested `<main>`, per-section
viewport height, visually odd. Move to standalone `/errors/*` routes through
`ErrorHandler`, or keep boxes?

**Recommendation: STANDALONE routes.** Honest presentation of what consumers actually
ship; removes the nested-`<main>` wart; the axe sweep auto-audits new demo routes, and
route goldens already cover the demo page set mechanically.

**Facts:** execution plan exists (`docs/planning/2026-09-17_18-20_ERRORPAGE-VISUAL-PARITY-PARETO-PLAN.md`
M03); demo keeps standalone `layout.Base` shells for recipe/error routes already
(AGENTS.md Demo Infrastructure).

**Options:** (a) standalone routes (new route-golden/axe surface, one-time);
(b) keep boxes + explanatory note (zero churn, wart stays).

**Reopen condition:** if the added routes make the visual suite's runtime or flake
budget noticeably worse, fall back with the note.

---

## #237 — AppShell/SidebarNav chrome default: adaptive (shipped) vs classic-dark

**Decision:** 2026-09-17 made SidebarNav/AppShell chrome theme-adaptive by default
(white chrome in light mode). Classic permanently-dark admin sidebar remains a
pure-CSS opt-out via `--tc-sidebar-*`. Ratify or flip?

**Recommendation: RATIFY adaptive.** It is shipped, goldens + visual suite cover it,
it matches the library-wide dark-mode convention (no mixed slate, `gray-*` only), and
the classic look costs consumers one CSS variable block — the opt-out is the cheap
direction, not the default.

**Facts:** SidebarNav is explicitly NOT a dark-mode exemption (AGENTS.md); the
`--tc-sidebar-*` token contract is documented.

**Options:** (a) keep adaptive default; (b) flip to classic-dark default (breaking-ish
visual change, needs golden re-baseline + changelog callout).

**Reopen condition:** a consumer majority reports light-mode chrome as wrong for
admin contexts — collect evidence, do not flip on assumption.

---

## #238 — Kanban failure UX policy: one auto-retry vs immediate honest revert

**Decision:** after a server rejects a move: keep immediate revert + manual re-drag,
or add a single auto-retry?

**Recommendation: KEEP immediate revert.** 4xx rejections (validation, conflict, 422)
are exactly the cases that should surface instantly; auto-retry masks server
decisions and complicates the optimistic-pending model (ADR-0041). The demo now has
a live 422 path proving the revert UX (`examples/demo/kanban_demo.go` moveErr).

**Facts:** Datastar retries network errors natively (~2 min) but not 4xx; htmx would
need a re-submit mechanism that does not exist today.

**Options:** (a) keep revert; (b) auto-retry once for network-class failures only
(5xx/timeout), never 4xx; (c) auto-retry everything (rejected — masks 4xx).

**Reopen condition:** a consumer reports lost drags on flaky networks → implement (b)
with the network/4xx split documented in the Kanban recipe.

---

## #239 — Kanban pending-timeout escape hatch for hung servers

**Decision:** a server that never responds leaves the card pending until the 30s poll
budget. Decide the escape-hatch API now, or reject until demanded?

**Recommendation: DOCUMENT, defer the API.** Write the 30s budget into the Kanban
docs (one line, no API change), and hold the escape hatch to the F109 demand gate —
zero consumers have asked. If it is ever forced: `PendingPolicy` enum beats
`Optimistic *bool` (extensible, no tri-state ambiguity, matches the typed-enum
convention).

**Facts:** no `KanbanBoardProps` opt-out exists; the fallback shapes were sketched in
`docs/status/2026-09-17_18-10` f9/e6.

**Options:** (a) document + defer; (b) ship `PendingPolicy` now (API surface without
demand); (c) ship `Optimistic *bool` (rejected — tri-state ambiguity).

**Reopen condition:** first consumer complaint about stuck pendings.

---

## #303 — website-module content clones: extract or blanket-accept

**Decision:** icon-box trio + clamp-typography clones in the website module — extract
shared helpers or blanket-accept into the ADR-0009 baseline?

**Recommendation: BLANKET-ACCEPT.** The website is repo-local, unpublished, and
content-shaped; the clones have exactly one consumer each, so extraction buys zero
reuse and costs an abstraction + baseline churn. This is the same policy call the
demo content already got.

**Facts:** ADR-0009 is the canonical accepted-clones gate; the 2026-09-25 re-record
holds 124 accepted groups.

**Options:** (a) blanket-accept (add hashes to the ADR-0009 baseline at the next
re-record); (b) extract (two helpers, no second consumer).

**Reopen condition:** the website gains a second visual theme or a second site —
then extraction has a consumer.

---

## #312 — export `display.HeadingTag` (public helper)?

**Decision:** promote the internal heading-level helper to public API?

**Recommendation: HOLD.** Exports are permanent support surface (semver + docs +
goldens); there is no consumer ask on record. Additive export is cheap LATER; removal
is breaking. Consistent with the F109 demand gate that governs #217.

**Facts:** API-surface flag from the wave that produced it (f44); v2 (ADR-0039) will
run a full API-surface review anyway — a better venue for a batch of promotion
decisions than a one-off now.

**Options:** (a) hold; (b) export now (docs + test + FEATURES row).

**Reopen condition:** a consumer hand-rolls heading-level logic, or the v2 API review
bundles it.

---

## #358 — GOEXPERIMENT=jsonv2 disposition

**Decision:** `encoding/json/v2` is verified stable WITHOUT the flag on the Go 1.27
floor. Drop the repo-wide exports (.envrc, flake shellHook, pre-commit.sh, ci.yaml
env blocks, AGENTS.md, website installation.md) in one lockstep commit, or keep until
a release boundary?

**Recommendation: DROP now, one lockstep commit.** The row's own criterion ("one
deliberate commit, not a drive-by") is satisfiable immediately post-push; keeping it
is dead weight on every tool invocation for no remaining benefit. Guards
(`TestEncodingJSONV2`-class drift guards and the map-key sorting gotcha in AGENTS.md)
stay regardless — they document behavior, not the flag.

**Facts:** stable without flag on go1.27.1+; flag-gated only on 1.26.x (which is below
the repo floor).

**Options:** (a) drop now; (b) drop at the v1.24.0 boundary (bundled with release
churn).

**Reopen condition:** none foreseen; if a <1.27 consumer appears, that is a consumer-
side export, not a repo concern.

---

## #411 — owner ratifies tool verdicts (Playwright / achecker / k6; chromedp stance)

**Decision:** ratify the accumulated lane verdicts; nothing integrates before this.

**Recommendation: ADOPT-as-trial, chromedp stays primary.** Concretely: crossbrowser
Playwright lane = TRIAL, findings-only, never gates (48 verified captures, 3 real CSP
bugs Chromium-only chromedp structurally missed; kill-criteria + budget row already
live in `docs/testing/tool-trial-protocol.md`); IBM achecker = keep as on-demand
ledger check (zero-hit re-scan 2026-10-17); chromedp = primary e2e driver, Playwright
additive for cross-engine only; k6 = reject for now (no load-testing demand on
record). One sitting, one verdict table update in `docs/testing/tooling-evaluation.md`.

**Facts:** wave-5 struck the underlying rows (#412/#413/#418/#419/#420) with
artifact-verified provenance; the final adopt/reject ratification is the only
remaining gate.

**Options:** per-lane adopt/trial/reject as above; the alternative (defer again)
re-opens re-research cost the tooling-evaluation doc exists to prevent.

**Reopen condition:** the crossbrowser TRIAL converts to full ADOPT only after its
kill-criteria window passes clean (budget table owns the date).

---

## #428 — warning hue unification: amber vs yellow

**Decision:** one hue for the warning semantic — amber (statusdot/scrollback/errorpage
precedent) or yellow (badge/kanban/stat precedent)? Today an Alert warning and a Badge
warning render different hues on one page.

**Recommendation: AMBER.** Full memo with contrast math and migration sketch already
written: `docs/planning/2026-10-15_warning-hue-decision-memo.md` (default rec AMBER,
amber-700 carries WCAG text headroom). Execution lands with the v2 tone-table
consolidation (TODO #425/#431) — adapters exist per enum today.

**Reopen condition:** only if the WCAG sweep (#330, closed clean) is re-opened with
contrasting data.

---

## #322 — status→family mapping: 500→Corruption vs Infrastructure

**Decision:** `errorpage.FamilyFromStatus` shipped 2026-10-17 with 500→Corruption
(following the in-repo precedent: `FromError` unknown → Corruption). The cqrs modules
disagree — adminui maps 500→Corruption, dashboardui maps 500→Infrastructure. Keep or
flip?

**Recommendation: KEEP 500→Corruption.** An unclassified 500 means the contract
itself failed — the library cannot know it was infrastructure; claiming
Infrastructure fabricates a diagnosis. Corruption is the honest family for "unknown
and therefore suspect". Internal consistency beats matching one consumer's legacy
choice.

**Facts:** flip cost if owner disagrees: one line in `FamilyFromStatus`'s default
branch + 3 test literals (documented in the TODO row). The shipped mapping is
exhaustive-tested over 100–599.

**Options:** (a) keep Corruption (library consistency); (b) flip to Infrastructure
(cqrs-dashboardui compatibility).

**Reopen condition:** if cqrs-htmx adopts the helper and its dashboardui callers
break on the family label — then flip deliberately with a changelog note.
