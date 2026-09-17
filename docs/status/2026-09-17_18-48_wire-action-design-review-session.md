# Wire.Action Design Review Session — Status Report

**Date:** 2026-09-17 18:48
**Scope:** This session only — a design review of `utils/wire.Action` ("How could we
design `&wire.Action{` better?"), run under an explicit **IDEAS AND PLANS ONLY**
constraint (no code changes). Two research passes: an initial design review, then a
deeper pass after the owner challenged research completeness
("Is this all the research worth doing???").

**Deliverables produced:** 2 in-chat design reviews (10 weaknesses + A/B/C idea tiers +
phased plan; then a corrected review with new findings and 2 retracted claims).
**Zero files modified** in the library itself (by explicit instruction).

---

## a) FULLY DONE

1. **templ-components skill loaded** (Part 1, consumer guide) before touching the task —
   per the mandatory activation flow.
2. **Core wire package fully read:** `utils/wire/wire.go` (all 316 lines — Action struct,
   both dialect renderers, expression builder), `doc.go`, `handler.go` (PatchTarget,
   PatchMode, Handler, IsDatastar/IsHTMX), `form.go` (DecodeForm + setters map).
3. **Governing ADRs read:** ADR-0036 (transport-wiring contract, incl. the Handler
   annotation) and ADR-0038 (ContentType/DebounceMS extensions + the Selector second
   extension).
4. **Usage inventory taken:** ~180 `wire.Action{` / `Wire` match sites swept across
   display, forms, navigation, demo, visualtest, and `cmd/tc/_sources` — the 10 wired
   components (Button, Form, FilterInput, FilterDropdown, NavLink, LoadMore, KanbanBoard,
   Calendar MonthNav, KanbanColumn.Action) plus raw `Attrs` spread usage.
5. **First design review delivered** (chat): 10 numbered weaknesses with file:line
   evidence, 3 idea tiers (A non-breaking / B v2-breaking / C re-freeze), 6-step plan.
6. **Second (deeper) research pass delivered** after the owner's challenge:
   - `docs/datastar-runtime-facts.md` read in full — surfaced bundle-verified fetch
     options (`mode`, `retry`, `useViewTransition`) and event modifiers (`throttle`,
     `once`, `stop`, `viewtransition`) that `Action` does not express.
   - `docs/wire-gates-d1-d2-d3.md` read in full — the owner decision gates (D1 no WC
     module, D2 silent fallback pinned, D3 "both dialects, no new runtime facts" rule)
     that any redesign must argue through.
   - TODO_LIST/ROADMAP swept for wire items: #155 (SimpleNav next adoption), #157
     (Calendar), #178 (trigger language, ADR-gated), #156 (adoption survey), v2 rows
     (busy signaling docs-only verdict, SSE writer helper, wire fuzzing).
   - **`hx-swap` usage sweep** — the decisive finding: 6+ in-repo sites hand-roll swap
     mode outside `wire.Action` (below).
   - Existing test surface mapped: `FuzzAction`, `FuzzDecodeForm`, the 6 invariants,
     handler/form tests, benchmarks.
7. **Two of my own pass-1 claims retracted with corrections** (see section d) — the
   corrected review supersedes the first.
8. **Corrected plan delivered** (Phase 0 bundle verification / Phase 1 five non-breaking
   items / Phase 2 v2 ADR bundle) and framed inside the D3 gate.

## b) PARTIALLY DONE

1. **Research closure — still has unread spans** (design conclusions drawn from partial
   reads):
   - `utils/wire/wire_test.go` lines 120–486 unread (`TestActionAttributes`,
     `TestActionDebounce` bodies — coverage inferred from names + the read parts).
   - `docs/transport-wiring.md` only lines 1–150 read — the "signaling notes" section
     that TODO #178 cites lives beyond that.
   - `handler_test.go` / `form_test.go` never opened (asserted-by-name, not read).
   - ADR-0030 / ADR-0035 / ADR-0039 never read directly — relied on AGENTS.md summaries.
   - templ-components SKILL.md **Part 2 (author guidance) never read** — only Part 1.
   Effort to close: S each; matters before executing Phase 1.
2. **Prior-art verification — weakly evidenced.** Two sourcegraph queries returned
   zero/noise (bad repo-pattern syntax), one `agentic_fetch` errored. The "no prior art
   unifying htmx+Datastar attribute rendering" conclusion rests partly on domain
   knowledge, not verified searches. Directionally very likely right (the repo's own
   bundle-audit rule exists because no SDK abstracts this), but not proven. Effort: S
   (retry with plain queries, check the Datastar Go SDK docs directly).
3. **Swap-mode value mapping — claimed, not bundle-verified.** The pass-2 finding
   (`Action.Swap` ↔ `hx-swap` / `{mode: '…'}`) is grounded in the facts doc's line about
   the `mode` fetch option, but the exact htmx↔Datastar value mapping
   (`outerHTML↔outer`, `afterbegin↔prepend`, `beforeend↔append`, `beforebegin↔before`,
   `afterend↔after`; htmx-only: `none`/`delete`; Datastar-only: `replace`) still needs
   the decode-from-pinned-bundle pass (Phase 0) before any design is final. This is the
   plan's own declared remaining unknown.

## c) NOT STARTED

All execution work — correctly so, blocked by the owner's **IDEAS AND PLANS ONLY**
constraint, not by neglect:

1. Phase 0: bundle decode of the `mode` fetch option + facts-doc entry.
2. Phase 1 items: `Action.Swap` field; `wire.Get/Post/Put/Patch/Delete` constructors;
   promoted `WithEvent`/`WithContentType` builders in wire (deleting the 4 duplicated
   default-apply helpers + the test mirror); `FuzzAction` extension (Selector param +
   injection assertion); `ThrottleMS`; optional `ViewTransition`; URL-template doc note.
3. Phase 2 (v2 ADR): Target/Selector unification, value-typed `Wire` field,
   `PatchTarget` rename, TODO #178 trigger language.
4. HARVEST: the pass-2 findings (Swap gap, stale ROADMAP fuzz row, the 3 owner
   questions) exist only in this report — TODO_LIST/ROADMAP untouched.
5. Any CHANGELOG `[Unreleased]` entry (nothing shipped, so nothing owed yet).

## d) TOTALLY FUCKED UP

1. **Pass-1 factual error — claimed a missing fuzz test that exists.** I asserted "no
   fuzz on the expression builder" without grepping the test files; `FuzzAction` lives at
   `utils/wire/wire_test.go:561` and `FuzzDecodeForm` at `decode_fuzz_test.go:32`. Caught
   in pass 2 only because the owner challenged the research depth. Root cause: claimed
   absence from partial reads instead of verifying by search — the exact anti-pattern the
   repo's verify-external-claims/verify-before-filing rules exist for. Severity: medium
   (a wrong claim in a design review the owner might act on); mitigated by the pass-2
   retraction. **Lesson: grep for existing coverage before claiming anything is missing.**
2. **Pass-1 missed the single biggest design gap** — swap mode. A one-minute
   `hx-swap` grep (run only in pass 2) immediately surfaced 6+ hand-rolled sites
   (`display/kanban.go:134`, `forms/calendar_nav.go:40`,
   `navigation/loadmore.templ:118,125`, legacy `HxSwap` props on Card/FilterDropdown,
   demo + e2e repetition) and the bundle-verified `mode` option. Root cause: I reviewed
   the API's shape instead of sweeping for workarounds first — "how do call sites work
   around this abstraction?" is the highest-yield first question for any redesign.
3. **Wasted round trips on prior-art search:** two malformed sourcegraph repo-pattern
   queries + one failed agentic_fetch, then I proceeded on knowledge. Should have either
   retried with simpler syntax immediately or explicitly labeled the conclusion
   unverified (I only did the latter under pressure).

## e) WHAT WE SHOULD IMPROVE

1. **Workaround-sweep-first review order.** Before judging an API's design, grep its
   call sites for manual patching-around (`attrs["hx-…"] =`, hardcoded attributes,
   legacy twin props). Impact: pass 1 missed its top finding; this inversion finds gaps
   in minutes. (Candidate for the templ-components skill's Part 2 or the linter/review
   habits.)
2. **Verify-before-claim-absence.** Never state "X does not exist" from partial reads —
   search first (`grep func Fuzz`, `func Test`). This bit me once this session and is a
   repeatable failure mode for design reviews.
3. **Load the decision-gate docs before designing.** D1/D2/D3 in
   `docs/wire-gates-d1-d2-d3.md` reframed the whole plan (what looked like a pure
   API-taste question is governed by an owner-ratified adoption rule). AGENTS.md links
   ADRs, but gate memos deserve a "read before redesigning wire" pointer.
4. **Read skill files completely when the task is author-adjacent.** SKILL.md Part 2
   (author playbook) was skipped; a wire redesign is author work.
5. **Search-tool hygiene:** keep sourcegraph queries simple (`repo:owner/name term`);
   on failure, retry once with corrected syntax before falling back to knowledge, and
   label unverified conclusions as such in the same breath, not after challenge.
6. **Close the read spans before execution.** The unread wire_test/transport-wiring
   spans could contain invariants that constrain the design (e.g., debounce edge pins).
   S effort each; do them as Phase 1 step 0.

## f) Next tasks (session-scoped, ranked)

| #   | Task                                                                                                              | Impact   | Effort | Category      |
| --- | ----------------------------------------------------------------------------------------------------------------- | -------- | ------ | ------------- |
| 1   | Phase 0: decode `mode` fetch option + swap-value mapping from the pinned bundle; add facts-doc section            | Critical | S      | Research      |
| 2   | `Action.Swap` typed enum (hx-swap / `{mode}`); both-dialect tests + goldens per D3 rule                          | High     | M      | Feature       |
| 3   | Migrate kanban/loadmore/calendar_nav off manual `hx-swap`; deprecate legacy `HxSwap` props on Card/FilterDropdown | High     | M      | Cleanup       |
| 4   | `wire.Get/Post/Put/Patch/Delete(url)` constructors + table tests                                                | High     | S      | Feature       |
| 5   | Promote copy-and-default builders (`WithEvent`/`WithContentType`/`WithDebounce`) into wire; delete the 4 dupes + `formWireAttributesLike` mirror | High | S      | Quality       |
| 6   | Extend `FuzzAction`: add `Selector` param + no-unescaped-quote injection assertion                               | High     | S      | Quality       |
| 7   | HARVEST this report: Swap gap, constructors, builders, ThrottleMS, stale ROADMAP row → TODO_LIST/ROADMAP         | High     | S      | Documentation |
| 8   | Refresh stale ROADMAP "Wire fuzzing" row (both fuzzers already exist)                                            | Medium   | S      | Documentation |
| 9   | `ThrottleMS` sibling of `DebounceMS` (`throttle:Nms` ↔ `__throttle.Nms`, spelling already bundle-verified)       | Medium   | S      | Feature       |
| 10  | Read the unread spans (wire_test 120–486, transport-wiring 150+, handler/form tests, SKILL.md Part 2)            | Medium   | S      | Research      |
| 11  | Doc-drift guard: pin `docs/transport-wiring.md` dialect table against real `Attributes()` output                 | Medium   | M      | Quality       |
| 12  | Invariant test pinning debounce-without-event htmx behavior; then keep-or-fix decision                           | Medium   | S      | Quality       |
| 13  | Document the `{year}` URL-template convention in wire `doc.go`                                                    | Medium   | S      | Documentation |
| 14  | `ViewTransition bool` (`useViewTransition` ↔ `transition:true`) — optional, after Swap                            | Low      | S      | Feature       |
| 15  | Draft the v2 ADR: Target/Selector unification, value-typed `Wire`, `PatchTarget` rename (per ADR-0038 mechanism)  | High (v2)| M      | Documentation |
| 16  | Retry prior-art search with corrected queries; record verified result in the ADR                                  | Low      | S      | Research      |
| 17  | TODO #155 follow-through: SimpleNav `Wire` adoption under the D3 rule (named next candidate)                     | Medium   | M      | Feature       |
| 18  | Owner Q3 (below) resolution: wire Phase 1 lands now vs. v2 freeze — gates items 2–6                              | Critical | —      | Decision      |

(18 items — session-scoped. The repo-wide backlog lives in TODO_LIST.md/ROADMAP.md and
was deliberately not re-audited, per instruction.)

## g) Questions I cannot answer myself

1. **Ship Phase 1 now or hold for v2?** The D-gates are owner decisions; ADR-0036/0038
   allow ADR-amendment extensions any time, but you may want wire surface frozen until
2. **`Action.Swap` naming/semantics:** reuse the server-side `PatchMode` value set
   (`inner`/`outer`/`prepend`/…) across both dialects, or use htmx-style names
   (`innerHTML`/`outerHTML`/`afterbegin`/…)? Both are defensible; the first unifies with
   `wire.Handler`, the second is what htmx consumers already type. (Depends on Phase 0's
   verified mapping.)
3. **Report format:** the status-report skill's canonical output is a styled HTML
   dashboard; you asked for `.md` and got it. Keep `.md` as your standing override for
   future status reports, or was this a one-off?

---

**Bottom line:** a two-pass design review of `wire.Action` is complete and
self-corrected; the top actionable finding (missing `Swap` in the common subset, 6+
hand-rolled workarounds, bundle-verified twin) outranks everything from pass 1. All
execution is untouched by design. Research is closed except one bundle-decode task
(Phase 0) and small read spans. Waiting for instructions.
