# Status Report — Kanban Board, Dual-Transport (HTMX + Datastar)

**Date:** 2026-09-09 01:50
**Session scope:** `display.KanbanBoard` — design, implementation, tests, demo, e2e, docs, count-sync. Point-in-time snapshot of THIS session only.
**Branch:** `master` @ `da5e758` ("feat(component): add KanbanBoard component with drag-and-drop support" — daemon-committed) · working tree clean.

---

## a) FULLY DONE

| #  | Item                                                                                                                                                                                                                                                                                                                                                                                                                           | Proof                                        |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------------- |
| ~~1~~  | ~~**Research phase** — archived kanban project (`/home/lars/projects/archived/kanban`) read: research-only (3 docs, no code); extracted Board→Column→Card + stable-position-ordering model. Codebase patterns studied: `utils/wire` (Action/Handler/ADR-0036), `forms.FilterInput`/`forms.Form` wire defaults, Badge/Carousel/contract/custom.css/demo/e2e harness~~ done at `b0ef4e2` | ~~—~~ |
| ~~2~~  | ~~**`display.KanbanBoard` component** (`display/kanban.go` + `kanban.templ`): horizontally scrolling columns; cards draggable via HTML5 DnD; per-card keyboard move buttons (prev/next column); both paths funnel into ONE hidden move form (`card`/`column`/`index` + optional CSRF) and `requestSubmit()`~~ done at `b0ef4e2` | ~~`nix run .#verify` ALL PASSED~~ |
| ~~3~~  | ~~**Dual-transport wiring** — `KanbanBoardProps.Wire *wire.Action`: htmx renders `hx-post` + component-owned `hx-target="#<boardID>" hx-swap="outerHTML"`; Datastar renders `@post(url, {contentType: 'form'})` (response targeting via `wire.Handler` outer mode). Component applies POST/submit/form-encoding defaults on a COPY (non-mutation test-pinned). Nil/empty Wire = read-only board (no drag, no buttons, no script)~~ done at `b0ef4e2` | ~~`display/kanban_test.go` (7-case wire table)~~ |
| ~~4~~  | ~~**Server contract** — `display.ParseKanbanMove(r) (KanbanMove, error)`: typed decode, static wrapped errors, index default 0, rejects empty card/column and negative/non-numeric index~~ done at `b0ef4e2` | ~~`TestParseKanbanMove` (7 cases)~~ |
| ~~5~~  | ~~**DnD JS singleton** (`tcKanbanAttached`) — document-level delegation, insertion index from pointer midpoints, same-column removal adjustment, no-op drop suppression, drop indicator classes, `CSS.escape` on consumer IDs, live-region announcement ("Moving X to Y"), CSP nonce~~ done at `b0ef4e2` | ~~`TestKanbanJSSingletonGuard`; browser-proven~~ |
| ~~6~~  | ~~**Accessibility** — `role="region"` + label, `<section aria-labelledby>` columns, `ul/li` semantics, count badges with SR text ("2 cards"/"no cards"), hover/focus-revealed buttons with descriptive labels, `aria-live="polite"` region, motion-reduce via shared constants~~ done at `b0ef4e2` | ~~`kanban_a11y_test.go` (5 tests)~~ |
| ~~7~~  | ~~**Invalid states unrepresentable** — cards/columns WITHOUT a stable consumer ID render but don't participate in moves (no empty `column-body=""` identities, not draggable); ID-less zones excluded from keyboard navigation~~ done at `b0ef4e2` | ~~`kanban_edge_test.go` (7 tests)~~ |
| ~~8~~  | ~~**Test lenses** — golden sweep (5 goldens: readonly, wired htmx incl. CSRF, wired datastar, empty board, card Content), BDD (6 behavior specs), edge, example (2 godoc), helpers coverage~~ done at `b0ef4e2` | ~~all green~~ |
| ~~9~~  | ~~**Drag-feedback CSS** (`templates/custom.css`) — dragging opacity, zone highlight, insertion line (`::before`/`::after` with logical `inset-inline`), light+dark~~ done at `b0ef4e2` | ~~custom-CSS guard green~~ |
| ~~10~~ | ~~**Demo** — `examples/demo/kanban_demo.{go,templ}`: two boards side by side (htmx + Datastar), independent mutex-guarded server states, `POST /api/kanban/{htmx,datastar}` with `wire.Handler` on the Datastar side, sticky-nav entry, code snippet~~ done at `b0ef4e2` | ~~demo tests green; page live~~ |
| ~~11~~ | ~~**E2E under real Chromium** (`visualtest/kanban_e2e_test.go`) — BOTH transports × BOTH interaction paths: keyboard-button clicks AND synthetic HTML5 drag events (`new DataTransfer()` + `DragEvent` dispatches), incl. chained second move on the re-rendered board~~ done at `b0ef4e2` | ~~2 tests PASS (~0.6s)~~ |
| ~~12~~ | ~~**Count/doc drift surfaces all synced** — README (2× count tables + display heading + catalogue), FEATURES (totals 121/122 + row), SKILL.md (121/43 + by-use-case row + catalogue row), AGENTS (122 generated + display row), ROADMAP, website `sections.ts`, demo hero `componentCount`, contract inventory (`display (33)`)~~ done at `b0ef4e2` | ~~`TestDocsCountDrift` green~~ |
| ~~13~~ | ~~**Docs** — `docs/transport-wiring.md` new "Dual-transport kanban board" section (contract + handler recipe + facts); `docs/wire-gates-d1-d2-d3.md` post-round D3 adoption note; AGENTS.md conventions bullet (+ chromedp lessons); CHANGELOG `[Unreleased]` warm~~ done at `b0ef4e2` | ~~—~~ |
| ~~14~~ | ~~**Repo-wide verification** — `nix run .#verify` (generate+build+test+lint, all 7 modules, 0 lint issues), per-module `GOWORK=off` isolation loop, FULL visual suite in Chromium (23s, all pixel goldens + e2e), templ-sync/lint-config/css-minified guards, dark-mode/motion/RTL/container-query/tailwind-source guards~~ done at `b0ef4e2` | ~~all green~~ |
| ~~15~~ | ~~**Generated files landed** — `kanban_templ.go` + `kanban_demo_templ.go` force-added past the BuildFlow gitignore gotcha; 18 kanban files tracked; daemon-committed~~ done at `b0ef4e2` | ~~`git ls-files` = 18~~ |
| ~~16~~ | ~~**Incidental fix** — `nix run .#css` re-minified `examples/demo/static/app.css` (daemon had committed a 4,986-line un-minified file; TODO #125 recurrence)~~ done at `b0ef4e2` | ~~1-line minified, freshness guard green~~ |

## b) PARTIALLY DONE

| Item                                   | State                                                                                                                | Gap                                                                                                                                                                                                                                                                                                                                                 |
| -------------------------------------- | -------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Touch/mobile story**                 | Keyboard buttons exist and work, DnD is desktop-only                                                                 | **The move buttons are `opacity-0` until hover/focus — on coarse pointers (phones) hover never fires, so the ONLY touch-capable control is invisible.** Tappable but unseen. Needs a `@media (pointer: coarse)` always-visible rule (requires a `tc-` class hook on the button container — not present yet). This is the session's most honest miss |
| **RTL verification**                   | Logical properties everywhere (`ms/me/ps/pe/inset-inline`), arrow icons carry `rtl:rotate-180`, RTL guard test green | Never rendered under `dir="rtl"` in a browser — the arrow-flip + column-order flip are class-level claims, not pixel/browser-proven (fold into TODO #160 RTL sweep)                                                                                                                                                                                 |
| **Docs-health HARVEST of section (f)** | This report exists                                                                                                   | New next-tasks NOT yet routed into `TODO_LIST.md`/`ROADMAP.md` — waiting for instruction per the "WAIT" directive                                                                                                                                                                                                                                   |
| **Live-region announce-after-swap**    | Announcement fires at submit time ("Moving X to Y")                                                                  | The board outerHTML swap wipes the live region; a post-swap "moved" confirmation (htmx `afterSwap` / Datastar listener) is a polish item, not wired                                                                                                                                                                                                 |

## c) NOT STARTED (deliberate or deferred — nothing half-built)

1. ~~**Pixel visual goldens** for the board (`visualtest/testdata/kanban/*` light+dark) — deliberately skipped: e2e proves function, and agent-generated PNGs add to the #80 human-eyeball backlog. Layout regressions (column width, card overlap) currently have NO pixel pin.~~ done — kanban goldens M3
2. ~~**Demo screenshots** via `nix run .#shots` including the new kanban section.~~ done — shots all routes M9
3. **Within-column keyboard reorder** (move up/down) — scope decision; current keyboard path only crosses columns.
4. **Card create/delete, column management** (the archived research covers full-board CRUD) — the component is deliberately move-only; CRUD is consumer composition.
5. **WIP limits / column tone / card priority accents** — candidate props, no demand signal yet.
6. ~~**`docs/DOMAIN_LANGUAGE.md` entries** for board/column/card/move terms (file not consulted this session).~~ done — DOMAIN LANGUAGE M6
7. ~~**Fuzz test for `ParseKanbanMove`** and a **render benchmark** (other components have both; kanban has neither).~~ done — FuzzParseKanbanMove M5
8. ~~**Cross-board drag guard** — dragging a card from board A onto board B submits to B's endpoint (server no-ops on unknown card). Considered acceptable; NOT documented in the end.~~ done — cross board guard M4

## d) TOTALLY FUCKED UP! (all transient, all caught and fixed in-session — nothing shipped broken)

1. ~~**Careless first drafts.** `kanban_demo.go` v1: `[]KanbanColumn` where `[]KanbanCard` belonged (×3) and an invented `withBaseProps()` helper. e2e v1: an INCOMPLETE `apply()` (dangling loop referencing a nonexistent `s.columns`), a malformed `fmt.Sprintf` format string, `Poll`-into-`string` unmarshal bugs (twice).~~ done (docs-health pass 2026-09-08)
2. ~~**Test assertions written before reading rendered output** — three fix rounds: `<h3 id="...">` with a trailing `>` that never matches (templ emits more attributes), Datastar attribute apostrophes templ-escaped as `&#39;`, and substring assertions colliding with the component's own inline-script text (`data-tc-kanban-column-body` matched the JS source). Plus hand-rolled `indexOf` helpers written TWICE before remembering the `strings` package. Lesson recorded: render-and-inspect BEFORE writing string assertions.~~ done (docs-health pass 2026-09-08)
3. ~~**Dark-mode scanner violation** — class strings split across concatenated `.go` lines put light/dark pairs on different physical lines (scanner is line-based). Fixed by moving class helpers into the `.templ` file (where single long lines are the norm).~~ done (docs-health pass 2026-09-08)
4. ~~**e2e `NodeVisible` deadlock** — `chromedp.Click(sel, NodeVisible)` never resolves for the opacity-0-until-hover move buttons (45s timeout burn). Fixed: plain coordinate click. Lesson memorialized in AGENTS.md.~~ done (docs-health pass 2026-09-08)
5. ~~**Missed count on first pass** — forgot the demo's `kanban_demo_templ.go` counts toward the generated-files drift guard (121→122). Caught by `TestDocsCountDrift`, fixed.~~ done (docs-health pass 2026-09-08)
6. ~~**templ LSP stale all session** — "template closing brace not found" diagnostic persisted for a file that generated and built fine; LSP never recovered. Builds were ground truth (AGENTS rule). Process improvement: `lsp_restart` early instead of ignoring for hours.~~ done (docs-health pass 2026-09-08)

## e) WHAT WE SHOULD IMPROVE!

1. ~~**Render-before-assert discipline** — the 3-round test-fix cycle is this session's biggest waste. A quick `utils.Render` + inspect (or generating goldens FIRST) would have caught every assertion bug at authoring time.~~ done (docs-health pass 2026-09-08)
2. ~~**The per-file lint rule exists for a reason** — `kanban.go` accumulated funlen/wsl/varnamelen findings that a per-file `golangci-lint run` at authoring time would have shown (AGENTS 2026-09-08 lesson; I still batched the lint).~~ done (docs-health pass 2026-09-08)
3. ~~**Touch-first thinking** — the a11y path was designed for keyboards, not fingers. `pointer: coarse` belongs in the original design checklist for ANY hover-revealed control (this pattern may exist elsewhere in the library — worth an audit).~~ done (docs-health pass 2026-09-08)
4. ~~**Design the funlen-compliant shape upfront** — kanbanJS was written as one 90-line builder then split into 3; the split was mechanical and should have been the first draft.~~ done (docs-health pass 2026-09-08)
5. ~~**The e2e readiness/poll types** — chromedp `Poll` returns the expression's value; Go string receivers for boolean expressions failed twice. A tiny `pollTrue(ctx, expr)` helper would prevent recurrence.~~ done (docs-health pass 2026-09-08)
6. ~~**Verify daemon state BEFORE trusting clean git status** — the daemon committed 6+ times mid-session; integrity was verified after the fact (gitignore tail, CSS minified, lint guard, 18 files tracked — all fine this time), but the check should be proactive after every daemon commit landing during active work.~~ done (docs-health pass 2026-09-08)

## f) Up to 50 things we should get done next

**Kanban follow-ups (P0 first):**

1. ~~Coarse-pointer UX: always-visible move buttons under `@media (pointer: coarse)` — needs a `tc-` class hook on the button container (P0, shipped-UX gap)~~ done — tc-kanban-buttons M1
2. ~~Pixel goldens: `kanban/{light,dark}.png` via `AssertScreenshot`; human eyeball (#80 caveat)~~ done — kanban goldens M3
3. ~~RTL browser render of the board (fold into TODO #160 sweep)~~ done — rtl kanban N7
4. ~~Cross-board drag: suppress in the drop handler (src board === target board) or document in transport-wiring.md~~ done — cross board guard M4
5. ~~Post-swap live-region confirmation (htmx `afterSwap` / Datastar listener)~~ done — post swap announcement M4
6. ~~Document touch/DnD limitation (until #1 ships)~~ done — touch doc M6
7. ~~Fuzz `ParseKanbanMove`~~ done — FuzzParseKanbanMove
8. ~~Benchmark `KanbanBoard` render (match other components)~~ done — BenchmarkHotPaths
9. ~~Demo: add kanban clicks to the #168 demo click-through e2e~~ done — demo flows kanban N3
10. ~~Demo: transport-toggle integration for the kanban section (currently always-both)~~ **Won't implement — side by side judged sufficient.**
11. Demo: file-backed board state across restarts (nice-to-have)
12. Scope decision: within-column keyboard reorder (move up/down)
13. Scope decision: card create/delete wiring via `Wire` (add-card button per column)
14. Candidate prop: column WIP limit (count badge turns red over limit)
15. Candidate prop: card tone/priority accent border
16. Candidate prop: column accent tone
17. Add KanbanBoard to the Dashboard recipe (or a dedicated recipe)
18. ~~`docs/DOMAIN_LANGUAGE.md` entries for kanban terms~~ done — DOMAIN LANGUAGE M6
19. ~~Record the rejected ContainerAware-for-Kanban decision in `docs/container-query-strategy.md` (full-width scroller; preempts future re-litigation)~~ done — container query reject M6
20. ~~Audit OTHER hover-revealed controls for the same coarse-pointer gap (#1 generalization)~~ done — coarse pointer audit M7

**Repo health (noticed this session):**
21. ~~docs-health HARVEST this report's section (f) into `TODO_LIST.md` (pending instruction)~~ done — M2 harvest
22. ~~visualtest module: 68 pre-existing lint findings (not CI-linted) — triage or add the module to the lint matrix/exclusions formally~~ done — visualtest lint zero N4
23. ~~`examples/demo/main.go`: pre-existing unused `heroWireLine` (gopls) + writestring warnings — 5-minute nit cleanup~~ done — heroWireLine removed M9
24. Investigate templ LSP stale-diagnostics-after-regenerate (recurring; maybe upstream `a-h/templ` issue)
25. ~~Release decision: cut v1.16.0 with KanbanBoard or accumulate (see question c)~~ **Won't implement — answered v1.16.0 cut.**
26. ~~Re-check daemon push state (`git fetch`) before any release — daemon pushes master unasked~~ done — pre cut origin checks
27. TODO #157: Calendar month-nav Wire candidate (next D3 survey item)
28. ~~TODO #175: axe-core scan — kanban board is a prime target (drag a11y verdict)~~ done — axe sweep N2
29. ~~TODO #158/#159: overlay open-state captures + 375px mobile sweep (now incl. kanban)~~ done — overlay mobile N8 N6
30. ~~CSS freshness: consider a CI-local guard so daemon un-minified app.css commits (this session's incidental fix) fail BEFORE landing (currently only post-hoc)~~ done — check-css-minified.sh

_(30 items — quality over padding; the remaining backlog lives in TODO_LIST.md #28–#178.)_

## g) Questions I can NOT figure out myself

1. **Touch strategy:** fix the invisible-on-touch buttons purely in CSS (`pointer: coarse` always-visible) — or do you want a real touch drag story (long-press polyfill / pointer-events based DnD), which is a dependency-budget decision the library has deliberately kept closed?
2. **Feature scope for v1 kanban:** within-column keyboard reorder, per-column "add card" wiring, and WIP limits are all cheap follow-ups — which (if any) belong in the component vs. left as consumer composition?
3. ~~**Release:** cut v1.16.0 now (KanbanBoard is verified end-to-end and `[Unreleased]` is warm) or let it accumulate more work first?~~ **Won't implement — answered v1.16.0.**

---

_Point-in-time snapshot. Written by the 2026-09-09 kanban session. WAITING FOR INSTRUCTIONS._
