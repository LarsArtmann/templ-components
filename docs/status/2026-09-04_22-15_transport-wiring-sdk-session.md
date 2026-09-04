# Status Report — Transport Wiring SDK Session (`utils/wire`)

**Date:** 2026-09-04 22:15 CEST
**Session scope:** "Smartly and composably support datastar, htmx, and WebComponents — does it make sense, and build it as a proper SDK."
**Verdict delivered:** htmx+Datastar composability = yes (built); Web Components as a third transport = category error; ADR-0033 reaffirmed, light-DOM WC = consumer recipe, library-level WC = pending owner ratification.
**Branch:** `master` (local only, **not pushed** — house rule). All session work auto-committed by the BuildFlow daemon across 6 heuristic commits (`07a51f7`, `1693783`, `cedc287`, `847b988`, `cafe728`, `7e2c13c`); working tree clean at report time.

---

## a) FULLY DONE

Each item: what + evidence + scope.

1. **ADR landscape research + runtime bundle verification.** Read ADR-0033 (WC rejection), ADR-0035 (Datastar scope freeze), ADR-0030, `docs/datastar-runtime-facts.md`; then verified the pinned `go-datastar/static@v0.4.0` bundle directly (Datastar v1.0.2): fetch actions have **no `target` option** (options are `selector, headers, contentType, filterSignals, openWhenHidden, payload, requestCancellation, retry*`), datastar sends `Datastar-Request: true`, non-SSE HTML responses are patched via `datastar-<key>` response headers, datastar on-event is `data-on:<event>` colon syntax. Evidence: greps/reads against `/mnt/buildcache/go-mod/.../static@v0.4.0/datastar.js`; facts match `docs/datastar-runtime-facts.md`. Scope: research only, no code.
2. **`utils/wire` package — the transport-agnostic wiring contract.** `Transport`/`Method`/`Event` typed string enums (3 new `*IsValid` funcs, package-level style matching datastar/enums_go.go), `Action{Transport, Method, URL, Event, Target}` with `Attributes() templ.Attributes` rendering `hx-*` (zero-value default per ADR-0030) or `data-on:<event>="@<method>('url')"`; empty URL → nil (inert); unknown enum → map/switch fallback; URL single-quote escaping mirrors `datastar.actionExpr`; four handler-contract constants (`HeaderDatastarRequest`, `HeaderDatastarSelector`, `HeaderDatastarMode`, `HeaderHXRequest`). Evidence: `utils/wire/wire.go`; `GOWORK=off golangci-lint run ./wire/...` = **0 issues**. Scope: new subpackage in the `utils` module (no go.mod dependency changes — templ was already a dep).
3. **Wire test suite.** 3 table-driven IsValid tests (all constants + unknowns), 10-case `TestActionAttributes` (both dialects, zero values, unknown fallbacks, quote escaping, nil-on-empty), `TestActionAttributesRender` through `templ.RenderAttributes` (proves `data-on:<event>` colon keys survive the attribute writer; pins the HTML-entity encoding `&#39;`), `FuzzAction` (3.8M execs, no panics, no empty-valued attrs). Evidence: `utils/wire/wire_test.go`, `go test ./wire/...` ok.
4. **`display.Button` pilot integration.** New `Wire *wire.Action` field on `ButtonProps`; attributes spread (after aria-label, before `Attrs`, so consumer `Attrs` win conflicts) in **both** `<a>` and `<button>` branches. Evidence: `display/button_go.go`, `display/button.templ`, regenerated `button_templ.go` with pinned templ v0.3.1020 (zero cosmetic drift); `TestButtonWire` (4 cases: nil-inert, htmx dialect, datastar dialect without target, Attrs-conflict precedence); full `go test ./display/` incl. all golden sweeps **pass** — nil-Wire rendering byte-identical to before.
5. **Demo: one endpoint, both transports, on the main page.** `examples/demo/wire_demo.templ` — `wireDemo()` section (TOC entry "Wire" → `#wire-transport`) rendering the *same* Action shape as an htmx Button (`hx-get` + `hx-target="#wire-htmx-out"`) and a Datastar Button (`data-on:click="@get(...)"`), with an explicit scope note about the Datastar target asymmetry; `/api/wire/fragment` handler branches on `wire.HeaderDatastarRequest` and echoes `Datastar-Selector: #wire-datastar-out` + `Datastar-Mode: inner`. Evidence: `examples/demo/wire_demo_test.go` — endpoint branching test (2 cases) + page-renders-both-dialects test, both pass; every static class token in the new markup cross-checked against committed `static/app.css` (13/13 present).
6. **ADR-0036: Transport-Agnostic Wiring Contract.** Documents the three-axis analysis (transport / element-model / encapsulation), formally supersedes ADR-0035's "no attribute-helper surface" clause **via that ADR's own revisit trigger** (owner request 2026-09-04), keeps the `datastar` module freeze intact (wire adds no Datastar components), reaffirms ADR-0033, and encodes the no-`target`-option bundle fact as the reason `Target` renders htmx-only. Evidence: `docs/adr/0036-transport-wiring-contract.md`.
7. **Consumer guide with the Web Components recipe.** `docs/transport-wiring.md`: axes model, quick start (Button `Wire` + `Attrs` composition for any component), zero-value table, dialect mapping table, "why Target is htmx-only" FAQ, one-handler-both-transports recipe with the constant table, deliberate scope-boundary table (what stays in htmx/datastar modules), and the complete light-DOM custom-element recipe with the ADR-0033 constraints. Evidence: file written; code snippets mirror the shipped demo/test code.
8. **Docs + drift-guard sync (all green).** FEATURES.md (new `utils/wire` section, Button row updated, Totals 57 typed enums / 54 IsValid / 117 generated files), CHANGELOG `[Unreleased]` Added entry (detailed), AGENTS.md (new convention bullet + 116→117 generated count), README.md both count mentions, `website/src/data/sections.ts` enum count 52→55. Evidence: `utils` module tests incl. `TestDocsCountDrift` **pass** (was failing red on exactly these 4 counts before the doc edits).
9. **Full verification sweep.** `templ generate ./...` → zero diff (dev-shell templ v0.3.1020); `go build ./...`; per-module GOWORK=off tests: utils, icons, errorpage, charts/echarts, datastar, htmx → all ok; root module `go test ./...` → all ok; `golangci-lint` root display/forms/layout → 0 issues; utils module lint → 0 issues (after fixing 14 findings: 10 forcetypeassert, 3 wrapcheck, 1 golines, 1 gci); `nix fmt` clean.

## b) PARTIALLY DONE

1. **SDK rollout breadth — pilot depth only.** Works: Button `Wire` field + universal `Attrs` composition. Remains: no other component has typed `Wire` sugar (LoadMore, ConfirmDelete, FilterDropdown, PolledRegion, forms all compose via `Attrs` only). Blocker: deliberate pilot scope; breadth is a design decision (see g3). Effort: M per component.
2. **Wire contract subset.** Works: method/URL/event/target common subset. Remains: no polling, reveal, `hx-swap`, confirm, indicators, OOB, SSE options — deliberately out (documented boundary table), but consumers wanting e.g. `hx-trigger="click delay:1s"` must drop to raw `Attrs`. Blocker: none; scope discipline (ADR-0035 history). Effort: S per opt-in extension, L if done as full dialect.
3. **Browser-level verification of the demo.** Works: HTTP-level tests assert exact rendered attributes + endpoint header branching. Remains: nothing clicked the buttons in a real browser — no chromedp/visualtest run proving htmx actually swaps into `#wire-htmx-out` and Datastar actually patches `#wire-datastar-out`. Blocker: none (nix provides Chromium); simply not done. Effort: M.
4. **CSS freshness assurance.** Works: all 13 new static class tokens + dynamic Grid/Button/InlineSuccess classes verified present in the *committed* CSS. Remains: did not run `nix run .#css` byte-stability recompile, so CI's CSS Freshness content-diff is unconfirmed locally. Effort: S.
5. **Race-mode verification.** Works: all tests pass in normal mode. Remains: `go test -race` (the flake `.#test` form) not run for wire/display this session. Effort: S.
6. **README / website / SKILL.md coverage of wire.** Works: counts updated everywhere; FEATURES/CHANGELOG/AGENTS document wire fully. Remains: README has no wire feature section; website (`sections.ts` site) has no transport-wiring guide page; the repo authoring SKILL.md (~/.config/crush) doesn't mention the wire package. Effort: S–M each.
7. **Commit hygiene.** Works: everything is committed and the tree is clean. Remains: the feature is smeared across 6 daemon commits with hallucinated messages (`chore: auto-commit N changed file(s) (heuristic)`) — invisible to `git log --grep`, and the docs/code boundary is arbitrary. Known repo-wide daemon behavior (documented in AGENTS.md), not fixable from this side. Effort: —.

## c) NOT STARTED

1. **Library-level Web Components module** — deliberately not started: requires owner ratification to supersede the Accepted ADR-0033 ("permanent architectural boundary"); the narrow shape (light-DOM hosts only, separate opt-in module, never Shadow DOM) is specified in the guide + ADR. Priority: waiting on g1.
2. **TODO_LIST.md / ROADMAP.md harvest** of this report's section (f) — per the status-report skill this belongs in docs-health HARVEST; not run (user instructed WAIT). Priority: next session opener.
3. **Wire adoption across the catalogue** (navigation.LoadMore, htmx.ConfirmDelete refactor, forms wired validation, FilterDropdown) — not started; awaiting g3.
4. **Server-side transport-branch helper** (`wire.Handler`-style wrapper encoding the demo's header-branching pattern as a library API) — not started; needs an API design pass.
5. **Golden snapshots for wired Button variants** — substring tests only; no golden files were added to `display/testdata`. Priority: medium (goldens are the repo's testing backbone).
6. **visualtest capture for the wire demo section** — none. Priority: low–medium.
7. **Benchmarks** for `wire.Action.Attributes` (repo convention: benchmark suites in 7 packages; wire has none). Effort: S.
8. **BDD spec** (`bdd_test.go` Ginkgo style exists in htmx/datastar modules) — wire has table tests only. Priority: optional per repo norms.
9. **Full `scripts/ci-repro.sh` pre-push reproduction** (incl. `--lint --css --visual`, coverage threshold, cold cache) — I ran the pieces manually; the complete CI step-for-step mirror was not executed this session.
10. **`nix flake check`** — not run this session.
11. **Release**: next version cut (utils v1.13.0 + root, 8-module tag lockstep) not started — `[Unreleased]` is warm and the release script will pick it up.
12. **ADR-0035 annotation** — ADR-0035's superseded clause doesn't yet point at ADR-0036 (docs-health ANNOTATE candidate). Effort: S.

## d) TOTALLY FUCKED UP

**Nothing shipped broken.** Master state at report time: all modules build, all tests pass, lint clean, count guards green. Four near-misses were caught and fixed *in-session* (listed for honesty, with root causes):

1. **Broken expression builder in first draft of `datastarActionExpr`** — produced `@'get'('...')` instead of `@get('...')` (string-concatenation confusion with escaped quotes). Severity if shipped: Datastar wiring silently dead. Root cause: hand-rolled quoting; fix: `fmt.Sprintf("@%s('%s')", method, escaped)` + the render test now pins the exact output. Caught: during self-review, before the first test run.
2. **`http.NewRequestWithContext` argument-order bug in the demo test** (`method` passed as `ctx`) plus a pointless `var _ = context.Background` import hack. Severity: compile error / smell. Caught: immediately after writing the file (build). Fixed properly (real `context.Background()` first arg, hack removed).
3. **Stray duplicated comment injected into `display/button_go.go`** by an imprecise edit (commented-out `External` line above the real one). Severity: cosmetic. Caught: right after the edit; removed.
4. **First wire test asserted the wrong HTML reality** — expected `data-on:click="@get('/api/items')"` literally, but templ HTML-entity-encodes attribute values (`'` → `&#39;`). This is correct HTML (decodes back in the DOM), so the *test expectation* was wrong, not the code. Caught by running the tests; fixed with an explanatory assertion + a testing note in the guide.

**Known-not-mine but noticed (pre-existing, untouched per "don't fix unrelated issues"):** README/FEATURES enum-count numbers were mutually inconsistent *before* this session (54 vs 53 vs 52 pairs across 4 locations) — I moved each by +3 and made the guard-tracked values exact, but the underlying inconsistency is only partially normalized (see f44). AGENTS.md's "31 enums have IsValid" line is stale vs. the actual 55.

## e) WHAT WE SHOULD IMPROVE

1. **String-assertions ≠ behavior proof.** I verified the demo via HTTP/HTML strings only; the whole point of the demo is runtime behavior (htmx swap + Datastar patch). Fix: one chromedp E2E (click both buttons, assert fragment content in both targets) — reusable pattern for all future transport demos. Impact: prevents "tests green, page broken" class.
2. **Prefer the repo's own CI mirror over hand-assembled verification.** I ran build/tests/lint manually; `scripts/ci-repro.sh --lint --css` exists precisely to catch what hand loops miss (per-module blind spots, CSS content diff). Fix: always finish sessions with the ci-repro form matching touched surfaces.
3. **Golden tests are the backbone; I shipped substring-only for new rendering.** Fix: add golden snapshots for wired Button variants (`-update`) alongside the substring invariants.
4. **ADR conflict surfaced too late.** I discovered the ADR-0033/0035 constraints mid-session and delivered the ratification question only in the final message; a 2-minute ADR glob at planning time would have front-loaded the decision. Fix: for any "should we support X" task, read `docs/adr/` first (pattern worth a line in the templ-components skill).
5. **Daemon-vs-feature commit hygiene.** 6 anonymous heuristic commits hide a coherent feature. Fix options (owner decision): short-lived feature branches even for sessions (PR flow is already the documented convention — I worked directly on master this time; the repo convention says `feat/*` branch + PR. That was a process miss on my part).
6. **`wire` silent fallback is a footgun for `Transport`.** A typo'd `wire.Transport("datastr")` silently renders htmx. Repo convention is map+fallback (no panics), but for a *dialect selector* (unlike a style variant) silent dialect-swap is semantically dangerous. Fix candidate: `Transport` unknown → keep fallback but expose `Action.Normalize() (Action, error)` or lint guard; needs a decision (see g2).
7. **Count-guard normalization debt.** Four doc files carry four different historical enum counts; guards only pin two of them. Fix: one normalization pass making all mentions derive from a single source (or from the guard itself).
8. **`wire` should document itself where component authors look.** SKILL.md + website guide pages are where future sessions/adoption will look first; AGENTS.md alone undersells it.

## f) 50 things to get done next

*(Brainstorm ranked roughly by impact; HARVEST should route TODO_LIST vs ROADMAP per docs-health. Impact/Effort/Category per the quality guide.)*

| #  | Task                                                                                                          | Impact  | Effort | Category      |
| -- | ------------------------------------------------------------------------------------------------------------- | ------- | ------ | ------------- |
| 1  | Owner decision: ratify/reject superseding ADR-0033 (library-level light-DOM WC module) — unblocks or kills #7  | High    | S      | Decision      |
| 2  | Harvest this list into TODO_LIST.md + ROADMAP.md (docs-health HARVEST)                                         | High    | S      | Documentation |
| 3  | Chromedp E2E: click both wire demo buttons, assert fragments land in both targets                              | High    | M      | Quality       |
| 4  | Run `scripts/ci-repro.sh --lint --css` pre-push reproduction                                                   | High    | M      | Quality       |
| 5  | Golden snapshots for wired Button variants (htmx + datastar) in `display/testdata`                             | High    | S      | Quality       |
| 6  | `go test -race` on utils/wire + display + examples/demo                                                        | High    | S      | Quality       |
| 7  | If #1 ratifies: scaffold `wc` module (light-DOM hosts, nonce-safe registration, zero templ-API coupling)       | High    | L      | Feature       |
| 8  | `wire.Handler` server helper: stdlib middleware that auto-branches Datastar/htmx headers (encode demo pattern)  | High    | M      | Feature       |
| 9  | Push `master` + watch CI (session ended local-only per house rule)                                             | High    | S      | Process       |
| 10 | ADR-0035 ANNOTATE: point superseded attribute-helper clause at ADR-0036                                        | Medium  | S      | Documentation |
| 11 | Decide unknown-`Transport` policy (fallback vs error path) — resolve g2, then pin with test                    | Medium  | S      | Decision      |
| 12 | `wire` doc.go split per repo convention                                                                        | Low     | S      | Cleanup       |
| 13 | Benchmarks: `Attributes()` htmx vs datastar paths                                                              | Low     | S      | Quality       |
| 14 | `navigation.LoadMore` gains `Wire` (transport-switchable pagination)                                           | Medium  | M      | Feature       |
| 15 | `htmx.ConfirmDelete` refactor onto `wire` internals or gains datastar twin (needs ADR check vs ADR-0035 freeze)| Medium  | M      | Feature       |
| 16 | `forms.Input`/`Select`/`Textarea` wired server-validation examples                                             | Medium  | M      | Feature       |
| 17 | visualtest capture for the wire demo section                                                                   | Medium  | S      | Quality       |
| 18 | Website (Astro) guide page: transport wiring + link from sections.ts                                           | Medium  | M      | Documentation |
| 19 | README feature section for `utils/wire` (beyond counts)                                                        | Medium  | S      | Documentation |
| 20 | SKILL.md (repo playbook) wire section + ADR-first trigger note                                                 | Medium  | S      | Documentation |
| 21 | `docs/DOMAIN_LANGUAGE.md` entries: Transport, Action, Wiring, dialect                                          | Low     | S      | Documentation |
| 22 | Normalize enum-count mentions across README/FEATURES/AGENTS to a single source                                 | Medium  | S      | Cleanup       |
| 23 | Update stale AGENTS.md "31 enums have IsValid" claim (actual: 55)                                              | Low     | S      | Cleanup       |
| 24 | Cut next release (utils v1.13.0 + root, 8-module tag lockstep) via release.sh                                  | High    | M      | Release       |
| 25 | Post-release go.sum tidy sweep (v1.11/v1.12 lesson) after tags propagate                                       | High    | S      | Release       |
| 26 | Cross-dialect property test: same Action must reference URL in both rendered dialects                          | Medium  | S      | Quality       |
| 27 | Fuzz `Target` selector input too (currently only transport/method/event/url fuzzed)                            | Low     | S      | Quality       |
| 28 | CSP invariant test: assert wire renders never emit `<script>` (integration/csp_nonce_test.go companion)        | Medium  | S      | Quality       |
| 29 | Demo: transport-switch toggle (query param flips Buttons between dialects) to sell the SDK pitch               | Medium  | M      | Feature       |
| 30 | Check upstream-watch issue for go-datastar/static pin drift (v0.5.0 seen in cache during research)             | Medium  | S      | Maintenance   |
| 31 | Migration recipe doc: converting an htmx-only page to dual-transport with wire                                 | Medium  | M      | Documentation |
| 32 | Evaluate `data-on-interval`/`data-on-intersect` support in pinned bundle before any Event-set expansion         | Medium  | S      | Research      |
| 33 | Wire Event coverage for form defaults (`EventSubmit` + `forms.Form` interplay doc)                             | Low     | S      | Documentation |
| 34 | Consider `aria-busy` helper for wired elements (htmx `hx-indicator` vs datastar indicator asymmetry doc first)  | Medium  | M      | Feature       |
| 35 | `docs/modularization/README.md`: mention `utils/wire` in the module DAG description                            | Low     | S      | Documentation |
| 36 | Demo CSS: run `nix run .#css` byte-stability + recompile if diff                                               | Medium  | S      | Quality       |
| 37 | `nix flake check` after all session changes                                                                    | Medium  | S      | Quality       |
| 38 | Consolidate session's 6 daemon commits story: note feature boundary in a docs commit (no history rewrite)      | Low     | S      | Cleanup       |
| 39 | Errorpage/datastar/htmx module lint loop after any further changes (kept green this session; keep it that way) | Low     | S      | Quality       |
| 40 | Property test: empty URL always renders nil across all enum combos (edge pin)                                  | Low     | S      | Quality       |
| 41 | Docs: testing note about HTML-entity-encoded attribute assertions promoted to a shared golden/testing doc      | Low     | S      | Documentation |
| 42 | Explore `hx-swap` opt-in field (`Swap string`) vs documented Attrs escape hatch — ADR note either way          | Low     | S      | Decision      |
| 43 | Add wire mention to `examples/demo` README/help if one exists                                                  | Low     | S      | Documentation |
| 44 | Pre-existing: AGENTS.md component-count table (display 42 vs skill 40) drift check                              | Low     | S      | Cleanup       |
| 45 | Pre-existing: FEATURES.md utils row says "0 components" though utils ships DismissScript/EnsureID etc.          | Low     | S      | Cleanup       |
| 46 | Review whether `wire` belongs in `internal/`-style docs sidebar of the website API reference                   | Low     | S      | Documentation |
| 47 | Consider committing a `wire` usage snippet to demo hero (visible SDK surface)                                  | Low     | S      | Feature       |
| 48 | Run docs-health VERIFY against ADR-0036 claims next session (fresh-eyes verification)                          | Medium  | S      | Quality       |
| 49 | Add `wire` to the contribution guide's "adding a package" checklist if such exists                              | Low     | S      | Documentation |
| 50 | Revisit ADR-0035 revisit-triggers after first external wire adoption data lands                                | Low     | S      | Decision      |

## g) Questions I cannot figure out myself

1. **Do you ratify superseding ADR-0033 with a light-DOM-only, opt-in `wc` wrapper module — or does the consumer-side recipe stay the final answer?** I read ADR-0033 in full (rationale + the "narrow exception" clause) and built everything that doesn't require overturning it; only you can decide whether "support Web Components" means library surface or recipe. This unblocks item f1/f7 (~L effort) or permanently closes the WC thread.
2. **For an unknown `Transport` value, do you want silent htmx fallback (current, repo-convention) or a loud failure path?** I tried to answer from the repo's conventions — they say map+fallback, zero runtime panics — but a typo'd `"datastr"` silently rendering the *wrong dialect* is qualitatively different from a wrong badge style, and the convention docs don't address that distinction.
3. **Should `Wire` become first-class across the catalogue next (LoadMore, ConfirmDelete, forms, FilterDropdown), or stay Button-pilot until a real consumer asks?** ADR-0035 exists because you value scope discipline ("nobody has asked for it" froze Datastar expansion) — but you just asked for this one, and I can't tell whether that makes the whole catalogue in-scope or only the contract + pilot.

---

*Point-in-time snapshot — goes stale. Section (f) is HARVEST input for TODO_LIST.md/ROADMAP.md. Format note: user explicitly requested `.md`; the status-report skill's canonical HTML dashboard format was overridden per its own rule.*
