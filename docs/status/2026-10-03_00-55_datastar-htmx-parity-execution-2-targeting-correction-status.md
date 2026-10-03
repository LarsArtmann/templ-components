# Status — Datastar ↔ HTMX Parity, Execution Session 2 (2026-10-03 00:55 CEST)

Scope of this report: THIS session only (continuing `2026-10-02_12-43_…-execution-1-status.md`),
plus what that session left open. No unrelated research.

**Headline:** 15 plan tasks advanced to done, but the session's real product is a
**runtime-fact correction with library-code consequences**: the pinned Datastar
bundle ignores ALL client fetch options for patch targeting and merge mode —
`wire.Action.Swap` was rendering an inert `{mode}` option under Datastar and
`Action.Selector`'s documented targeting behavior does not exist in the runtime.
Found by the L1-06 browser e2e, root-caused (bundle decode + live probe), fixed in
code, tests, and 7 documentation surfaces. One demo test currently RED because of
the correction (fix known, not yet applied — you interrupted before I could).

---

## a) FULLY DONE (verified this session)

| Item                                                                                                                                                                                                                                                                                                                                                                                                                                                             | Evidence                                                                                                                                                                                                                                                       |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **L1-05 finished** (was 80%): ADR-0043 (`docs/adr/0043-typed-trigger-language.md`), ADR-0036/0042 cross-links, `docs/transport-wiring.md` trigger section rewritten (the stale "stays dialect-local" note replaced), FEATURES + skill wire rows for `Interval`/`Reveal`                                                                                                                                                                                          | utils suite green incl. wire tests                                                                                                                                                                                                                             |
| **L1-28 finished**: AGENTS.md daemon-flip note ("`1.26.0` flips harmless by construction")                                                                                                                                                                                                                                                                                                                                                                       | doc-only                                                                                                                                                                                                                                                       |
| **L1-07** `wire.Get/Post/Put/Patch/Delete(url)` constructors (`utils/wire/builders.go`)                                                                                                                                                                                                                                                                                                                                                                          | builds; wire suite green (dedicated table tests still MISSING — see c)                                                                                                                                                                                         |
| **L1-08** fluent `With*` builders (Event/Method/Target/Selector/Swap/ContentType/Debounce/Throttle/Interval/Reveal/PreventDefault — all value-receiver, copy-returning, non-mutating) + `WithFormDefaults(event)`; forms migrated (`formWireAttributes` → builder, `wireHosted` nil-boundary; filter_input, filter_dropdown); kanban rebuilt on the builder; the duplicated copy-and-default logic (`wireAttributesWithDefaults` + kanban's inline copy) DELETED | templ regenerated (pinned v0.3.1020); forms/display suites green                                                                                                                                                                                               |
| **L1-09** `ThrottleMS`: htmx `throttle:<n>ms` ↔ Datastar `__throttle.<n>ms` — **bundle-verified first** (decoded `de(e,t)` modifier parser: `throttle` set, `noleading`/`trailing` sub-mods, `te` duration parser ms/s/bare-ms)                                                                                                                                                                                                                                  | field + both renderers in; suite green                                                                                                                                                                                                                         |
| **L1-10** `utils/wire/doc.go`: one-object rule, `{year}`/`{month}` URL-template convention, common-subset + ADR extension map, server-side entry points                                                                                                                                                                                                                                                                                                          | doc-only                                                                                                                                                                                                                                                       |
| **L1-06 Swap/mode browser e2e** — `TestWireE2ESwapModes` GREEN, 5 cases: htmx outerHTML swap; Datastar header-driven outer patch; **server `Datastar-Mode: append` header beats a client `{mode:'inner'}`**; id-matched outer patching (header-less + fragment-root id); client `{selector}` does NOT target (probe folded in as permanent pin)                                                                                                                  | `nix run .#visual -- -run TestWireE2ESwapModes` PASS (witnessed)                                                                                                                                                                                               |
| **THE CORRECTION** (the big one): pinned `go-datastar/static` v0.6.1 patch target+mode are response-header-driven ONLY. Static decode (airtight: the merge branch reads fetch-state key `$` which is never assigned → dead code) + live browser probe (`regionGone=false`)                                                                                                                                                                                       | corrected in: `wire.go` (`datastarActionExpr` drops `{mode}`, renders `{selector}` only under `ContentTypeForm`; Swap/Selector godocs), `wire_test.go` (TestActionSwap, TestActionSelector), `triggers_test.go`, `invariants_test.go` — all suites green after |
| **Docs corrected** (7 surfaces): `docs/datastar-runtime-facts.md` (2 bullets rewritten as CORRECTED with evidence), ADR-0038 (inline correction annotations + a Correction section), `docs/transport-wiring.md` (6 sections: field table, dialect mapping, one-object para, form targeting row, Selector section, Swap section), FEATURES.md wire rows, AGENTS.md wire bullet, skill/SKILL.md wire row                                                           | grep sweep clean except 2 known sites (see d)                                                                                                                                                                                                                  |
| **L1-30 reassessment**: the daemon had already PUSHED the 9 snapshot commits before this session — folding them now would require force-push (forbidden without your approval). Marked MOOT; 8 NEW unpushed daemon snapshots exist (see f)                                                                                                                                                                                                                       | `git log origin/master..HEAD` = 8                                                                                                                                                                                                                              |
| L1-32/L1-31/L2-01.1/L1-01/02/03/04/L1-27/L1-29                                                                                                                                                                                                                                                                                                                                                                                                                   | done in session 1, re-verified loaded                                                                                                                                                                                                                          |

## b) PARTIALLY DONE

1. **L1-07/08/09 dedicated tests missing** — constructors/builders/ThrottleMS ship with only indirect coverage (the migrated component suites + corrected wire suite). The plan's L2-07.2 table tests, L2-08 migration assertions, L2-09 throttle rendering tests are NOT written. Violates the same-commit-tests rule; it happened because the e2e failure interrupted the batch.
2. **L1-11 (naming decision ADR addendum)**: the decision is recorded in plan §7 (session 1) and ADR-0043 documents `PatchMode` reuse, but the promised ADR-0038/0042 addendum paragraph about the naming is not written.
3. **Facts provenance table**: `docs/datastar-runtime-facts.md` has the corrected facts but the planned v0.6.1 provenance-table row (sha256, byte-identity note) is not added.
4. **CHANGELOG `[Unreleased]`**: NOTHING of sessions 1–2 is entered yet (PolledRegion, LoadingButton, typed triggers + ADR-0043, builders, ThrottleMS, the targeting correction). The release script would refuse to cut — correctly.
5. **L1-06 "complete"**: e2e green, but the demo-side collateral (next section) belongs to the same task and is not fixed.

## c) NOT STARTED (still open from the plan)

L1-13 (Tabs Wire), L1-14 (SimpleNav Wire), L1-15 (LoadMore → wire Swap + reveal trigger), L1-18/19/20 (demo cards + smoke), L1-25 (real fuzz runs), L1-26 (one-object demo-level guard), L1-12 (ViewTransition), L1-16 (Calendar settle doc), L1-17 (confirm eval), L1-21 (recipes + website api-reference/guide), L1-22 (DOMAIN_LANGUAGE Patch Mode), L1-23 (witnessed `nix run .#verify` + ci-repro at tip), L1-24 (full `nix run .#visual` pass), docs-count resync after demo cards.

## d) TOTALLY FUCKED UP (in-flight damage + process failures, honest list)

1. **RED: `TestWireDemoBusyCardRendersBothDialects`** (`examples/demo/wire_demo_test.go:701`) — the demo busy card was BUILT on the now-corrected false claim ("zero response-header routing, client `{selector}` targets"). My correction makes the test assert an expression that no longer renders. Worse than the test: **the demo's Datastar busy button can no longer patch its region** (endpoint `main.go:484` is a plain handler — no headers — and `wireBusyDone`'s root id doesn't match the region, so id-match can't save it). Fix (~15 min, planned, not applied because you called STOP): wrap the endpoint in `wire.Handler(PatchTarget{Selector: "#wire-busy-datastar-out", Mode: inner})`, drop `Selector` from the demo button's action, update the test's expected expression + the stale comment in `main.go:477-483`, and check the busy flow e2e in `visualtest/demo_flows_e2e_test.go`.
2. **Website guide still teaches the wrong model**: `website/content/docs/guides/transport-wiring.md` lines 71+73 (client-side `Action.Selector` targeting; "`{mode: 'inner'}` renders for Datastar") — will ship the correction's opposite to consumers if deployed.
3. **Process failure — I asserted before verifying**: I wrote the L1-06 e2e EXPECTATION from the session-1 note ("{mode} overrides the header") instead of bundle-verifying first. The repo's own quality bar (bundle-verify EVERY runtime claim) would have prevented one wasted e2e cycle and a wrong-test-detour. The e2e failing WAS the process working (that's why browser tests exist), but the detour was avoidable.
4. **Daemon races**: my e2e rewrite was rejected mid-edit twice because the daemon committed between read and edit; I also had to redo an edit after a stale-read rejection. Cost: ~4 wasted tool calls. Mitigation used after: re-read → edit; but I did NOT proactively check `git status` before long edit sequences.
5. **The stale `gopls` diagnostics** (phantom `datastarActionExpr` arg-count errors) persisted all session after the fix — known repo gotcha (LSP lag; build is ground truth), but I burned a few reads on them anyway.

## e) WHAT WE SHOULD IMPROVE

1. **Rule reinforcement**: "bundle-verify → then write the test expectation" must apply to e2e assertions too, not just component code. I'll encode the corrected targeting model as a checklist item in the skill's verification section.
2. **Browser-proof claims EARLY for any "client-side option" story**: string/golden tests cannot falsify runtime behavior — ADR-0038 shipped two claims (Selector-override, client-mode) that a single 2019-style e2e would have killed in September. New wire-surface claims get their e2e IN THE SAME session as the claim.
3. **Fold daemon commits at a quiet moment** (8 unpushed snapshots carry sessions 1–2 mixed together) — or better, get your explicit policy and stop triaging this every session.
4. **Dedicated tests in the same batch as new API** (builders/throttle) — no more deferring.
5. Consider a tiny **runtime-facts regression guard**: the corrected facts deserve a pinned contract test (e.g. golden asserting NO `mode:`/bare `selector:` in Datastar expressions beyond the form case) — partially exists in invariants_test now; could be tightened into `TestPinnedRuntimeBundleContract` token checks.
6. The docs sweep found stale claims in THREE layers (docs/, website/, demo/) — future corrections should include `grep -r` across all three as a fixed step.

## f) NEXT 50 (ranked; top 10 are the critical path)

1. Fix the demo busy card (endpoint → `wire.Handler`, drop Selector, fix test + comment) — unblocks RED suite.
2. Fix `website/content/docs/guides/transport-wiring.md` lines 71/73.
3. Re-run `go test ./examples/demo/...` + demo flow e2e to green.
4. Write dedicated tests: constructors table (L2-07.2), builders non-mutation/nil-path table (L2-08), ThrottleMS both dialects (L2-09.1/2).
5. CHANGELOG `[Unreleased]`: PolledRegion, LoadingButton, typed triggers + ADR-0043, builders/constructors, ThrottleMS, **targeting-model correction (breaking-ish, honest wording)**.
6. Facts provenance table row for v0.6.1 (sha256 + byte-identity to v0.5.0).
7. L1-11: naming addendum paragraph in ADR-0038/0042.
8. L1-13 Tabs Wire (+ goldens, both dialects).
9. L1-14 SimpleNav Wire (+ goldens).
10. L1-15 LoadMore → wire Swap + typed reveal trigger (+ goldens; e2e already covers infinite scroll? verify).
11. L1-12 wire ViewTransition (`useViewTransition` ↔ htmx `hx-view-transition`) — bundle-decode first.
12. L1-16 Calendar `settle:0s` exception doc.
13. L1-17 Datastar confirm evaluation (bundle fact first, then decide).
14. L1-18 demo card: Swap under both transports (uses the corrected model — handler owns Datastar mode).
15. L1-19 demo card: PatchModeRemove retraction.
16. L1-20 demo smoke route + goldens for the new surface.
17. `nix develop -c templ generate ./...` from repo root + demo CSS recompile after any `.templ` change.
18. L1-26 one-options-object demo-level render guard.
19. L1-25 real fuzz runs (-fuzztime=30s ×3) + seed-corpus note.
20. L1-21 recipes: transport-migration.md + website api-reference wire rows.
21. L1-22 DOMAIN_LANGUAGE.md Patch Mode entry + skill interop note.
22. Docs-count drift resync after demo cards (components/goldens counts unchanged unless demo pages count).
23. Update `docs/status` session-1 report's "3 open questions" section with today's resolutions.
24. L1-23: witnessed `nix run .#verify` + `scripts/ci-repro.sh --lint --website` at tip.
25. L1-24: full `nix run .#visual` pass; triage flakes vs real (theme-pin signature).
26. Fold the 8 unpushed daemon snapshots into semantic commits at the quiet moment (needs your push/force policy — see g1).
27. Re-check `TestDocsCountDrift` after every docs batch above.
28. Skill SKILL.md: add the corrected targeting model to the wire section's verification notes (done partially) + Part 2 checklist item.
29. Consider deprecating `Action.Selector` for non-form use entirely in v2 docs (it renders nothing there now).
30. Add `Reveal`/`Interval` rows to `docs/transport-wiring.md` dialect-mapping table (currently only in prose).
31. ADR-0043: append the correction cross-reference (mode claims inside it cite Swap rendering).
32. Check `cmd/tc/_sources` mirror sync after forms/display `.templ` edits (pre-commit guard does it; verify commit content).
33. Run `golangci-lint fmt` on the new files (builders.go, triggers.go) for golines/120 compliance.
34. `nix run .#lint` full pass before the witnessed verify.
35. website api-reference: PatchMode/Swap wording sync.
36. Golden sweep: confirm no golden anywhere contains `{mode:` (grep testdata).
37. `integration/csp_nonce_test.go` still green after demo changes.
38. Demo `wire_demo_test.go`: add a positive assertion that NO bare `{selector}` renders for the busy button (regression pin).
39. Consider a `wire.Handler` adoption in the demo validate endpoint (consistency).
40. Record the "dead-code `$` key" decode detail in the facts doc appendix (done in bullet; consider linking the minified snippet).
41. Update ROADMAP if the correction affects any parity item wording.
42. Check ADR-0042's surface-3 wording against the final trigger API (no drift).
43. DOMAIN_LANGUAGE: PatchMode vocabulary cross-link from Swap entry.
44. Verify `visualtest/wire_forms_pack_e2e_test.go` unaffected (it uses header targeting — should stay green; run it).
45. Run `go test ./integration/... ./internal/...` once before the witnessed verify.
46. Consider CHANGELOG "Fixed" entry phrasing review with github-voice before release.
47. Sweep `docs/recipes/*.md` for Selector/Swap claims (server-side-validation recipe likely fine; verify).
48. Confirm `tc add` mirror parity for any touched component sources.
49. After L1-15: re-check LoadMore's `InfiniteScroll` e2e coverage still passes with the typed reveal trigger.
50. Final: re-run the FULL plan checklist against AGENTS.md conventions (every new enum IsValid? — ThrottleMS is an int, no enum; verify no new enum slipped in untested).

## g) THREE QUESTIONS I CANNOT ANSWER MYSELF

1. **History policy**: the daemon pushed the 9 WIP snapshot commits to origin/master, and 8 more have accumulated unpushed. I will NOT force-push without your say-so. Fold the 8 at the tip via rebase (rewrites unpushed history only — safe), then push normally? Or leave snapshots and commit future work semantically on top?
2. **Release framing for the targeting correction**: `Action.Swap`/`Action.Selector` behavior changed (Datastar `{mode}` no longer rendered; `{selector}` only under form encoding). No in-repo consumer relied on the old behavior (it never actually worked), but downstream consumers MAY have copied the documented pattern. Ship as a bugfix in the next v1.x with a loud CHANGELOG/migration note, or hold the wire changes for v2?
3. **htmx.PolledRegion** (open since session 1): leave it htmx-native (my standing recommendation — its props are richer: eager-load, live-politeness plumbing) and let `datastar.PolledRegion` + `wire.Interval` serve the Datastar side, or migrate it onto wire triggers for internal consistency?

— Report ends. Per your instruction: WAITING FOR INSTRUCTIONS.
