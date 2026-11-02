# Status: Generics Session — Coalesce Shipped, 26 Sites Adopted, Self-Review

**Session:** 2026-11-01 evening → 2026-11-02 00:31 CET (single session, master worktree)
**Scope:** "How can we better use generics?" — discovery, implementation, adoption, verification, docs.
**Sibling activity:** one concurrent session shipped the five-skill architecture audit
(data-model-review HTML, assessment, D2 diagrams) in the same tree. No file conflicts;
its files were left untouched per the single-writer protocol.

---

## a) FULLY DONE

| # | Item                                                                                                                                                                                                                                                                                                                                                 | Evidence                                                                                                           |
| - | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| 1 | **`utils.Coalesce[T comparable](values ...T) T`** — first-non-zero-value helper, the generic form of the override-or-default idiom                                                                                                                                                                                                                   | `utils/utils.go`; 7-subtest `TestCoalesce` PASS; `BenchmarkUtils/Coalesce` added                                   |
| 2 | **26 templ call sites converted** across 6 packages (display 14, forms 5, feedback 2, layout 1, navigation 3, errorpage 2 — in 19 files). Nested chains collapsed: 3-level aria-label chain in `display/loading_dialog.templ:47`, 2-level chains in `tags_input`/`filter_dropdown`/`rating`                                                          | `rg -c` recount: 26 sites / 19 files (re-verified after publication — see d1)                                      |
| 3 | **21 generated `*_templ.go` files regenerated** via the canonical ci-repro sequence (root generate + per-module from inside each module); zero FileName churn elsewhere                                                                                                                                                                              | daemon commit `13af2e67` diff = exactly the touched packages                                                       |
| 4 | **`utils.Lookup` adoption** in remaining hand-rolled lookups: `gridClass`/`gridContainerClass` (default-key fallback form) and `forms/input.templ` `enterKeyHintValue`                                                                                                                                                                               | `display/grid.templ:74`, `forms/input.templ:77`                                                                    |
| 5 | **`utils/format` twins extracted** — `Bytes`/`CompactCount` now 3-line wrappers over a shared `compactUnits` scaler (sign + unit ramp + boundary bump); TODO_LIST row **#445 dispositioned DONE** with the plain-helper-not-generic rationale                                                                                                        | `utils/format/format.go`; boundary test `1048575 → "1.0 MiB"` PASS; benches 67–176 ns/op (in the documented range) |
| 6 | **Verification complete:** `go build ./...` OK; per-module test loop green (utils, icons, errorpage, charts/echarts, datastar, htmx); root suite exit 0; lint **0 issues** on root packages + utils + errorpage; goldens byte-identical (no `-update` needed — proof of output invariance); class inventory regenerated + committed per guard design | session verify runs; daemon commits `a42ae953`, `0f89d2a1`, `13af2e67`, `cbb0603c`                                 |
| 7 | **Docs synced in the same change:** FEATURES utils table (+Coalesce row), `skill/SKILL.md` helper list, `references/component-catalogue.md` (+row), CHANGELOG `[Unreleased]` (Added + Changed), AGENTS.md aria bullet (Ternary → Coalesce), TODO_LIST #445                                                                                           | all committed                                                                                                      |
| 8 | **Rejected-with-rationale (decisions, not omissions):** `display.List[T]` (ROADMAP demand-gate, no second ask), tone-map collapse (v2, TODO #425/#431), `withDefaults` generic (idea #80 — zero-value rule already covers it), generic enum-set type (Go can't enumerate constants)                                                                  | recorded in final summary + this report                                                                            |
| 9 | ** templ pin integrity held:** `templ --version` = v0.3.1020 verified against go.mod before generating; system PATH binary matched the dev-shell pin, so generation was safe                                                                                                                                                                         | session log                                                                                                        |

## b) PARTIALLY DONE

| # | Item                                       | State                                                                           | What remains                                                                                                                                                                                                   |
| - | ------------------------------------------ | ------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | Browser-level confirmation of the refactor | Rendered HTML proven byte-identical via golden sweeps; visual suite **not run** | `nix run .#visual` once, belt-and-braces, in a daemon-free window. Justified to skip (identical HTML ⇒ identical pixels), but the repo's own culture prefers browser-proven for anything touching render paths |
| 2 | Generic-idiom documentation depth          | AGENTS aria bullet + FEATURES + catalogue updated                               | `references/authoring-guide.md` has no decision node yet ("override-or-default ⇒ `Coalesce`, bool ⇒ `Ternary`") — future contributors will rediscover the distinction                                          |
| 3 | Benchmark evidence                         | One 100x smoke run                                                              | No Coalesce-vs-Ternary comparison numbers recorded anywhere durable                                                                                                                                            |

## c) NOT STARTED (noticed this session, deliberately not touched)

1. Stale doc comment `utils/darkmode_scan.go:32` still references `lookupFeedbackStyle[T]()` — a function removed when AlertType/ToastType aliases died (v2, ADR-0022). One-line fix, left for a dedicated hygiene commit.
2. `internal/contract` / authoring-guide guard coverage for "override-or-default must use Coalesce" — no guard exists; arguably shouldn't (style, not invariant).
3. `withDefaults`-adjacent idea #80 and `Slot` type idea #43 — untouched, still idea-grade.
4. Art-dupl baseline re-record — this refactor _removed_ duplication; stale accepted-group hashes are harmless but the baseline could shrink. Not started.

## d) TOTALLY FUCKED UP

**Nothing shipped broken** — build, full test loop, lint, and goldens all green, and rendered output is byte-identical. But honest process fuckups, ranked:

1. **I published a wrong count.** CHANGELOG (and my summary) said "24 templ render sites"; the true number is **26** (19 files). Caught it while writing this report by re-running `rg -c` instead of trusting my own arithmetic — exactly the brutal-self-review Q12 failure class (a claim asserted from memory, not re-verified). **Fixed forward this turn** (CHANGELOG now says 26). The lesson is stale: I had the full grep output in-context and still hand-summed wrong.
2. **Six parallel multiedits failed in one round trip** ("you must read the file before editing") — I violated my own read-before-edit discipline under batching pressure. Cost: one wasted tool round, no damage.
3. **`bar_chart.templ` edit failed on an indentation guess** (two sites, different nesting depths, I assumed equal tabs). Fixed with a tab-free `replace_all` — but that was a second attempt caused by not copying whitespace exactly from the View output.
4. **Sloppy bench invocation:** `go test -bench BenchmarkUtils/Coalesce` without `-run` executed the _entire_ utils suite including the class-inventory guard, producing a scary FAIL mid-verify that was just the by-design regenerate-once behavior.
5. **Em dashes in new godoc** — violated the global no-em-dash rule in the first `Coalesce` doc draft; self-caught, but the fix triggered a second class-inventory regeneration cycle.
6. **Session work fragmented across 4 daemon commits** (`a42ae953`, `0f89d2a1`, `13af2e67`, `cbb0603c`) with generic messages — known repo condition, not my doing, but the history for this change is hard to read (mitigation below, question g3).

## e) WHAT WE SHOULD IMPROVE

1. **Never publish a number I haven't re-counted with a tool in the same turn.** The 24-vs-26 slip is the entire lesson of this session.
2. **Read files before batch edits, always** — batching multiedits across unread files optimizes the wrong step.
3. **Isolate bench runs** (`-run '^$' -bench …`) so guard suites don't fire and muddy the signal.
4. **Codify the two lookup shapes in the authoring guide:** value-with-fallback ⇒ `utils.Lookup`; _conditional attribute render_ (omit when absent, e.g. `inputmode`) ⇒ plain `if v, ok :=` — a future contributor could wrongly "fix" inputModeForType into a always-rendering `Lookup` and emit `inputmode=""`.
5. **Decide Ternary's documentation status** (see g1) so the codebase converges on one override idiom instead of three coexisting shapes (Ternary, Coalesce, inline if).
6. **Fix the `lookupFeedbackStyle` ghost reference** (darkmode_scan.go:32) — small, but it misleads grep-driven navigation exactly like the one that produced this session's task.
7. **Belt-and-braces visual run** after output-shaping refactors, even when goldens prove identity — cheap confidence, aligns with the string-proven ≠ browser-proven culture.
8. **Daemon history curation:** for multi-commit sessions, consider one interactive rebase (local-only) before push so a feature's story isn't smeared across `chore: auto-commit` messages (see g3).

## f) UP TO 50 THINGS TO GET DONE NEXT

_Brainstorm-grade (status-report skill rule): most of these are ROADMAP fuel, not commitments. Grouped, rough priority order._

**Direct follow-ups (this session's tail):**

1. Run `nix run .#visual` once over the refactor for browser-level confirmation.
2. Add the lookup-shape decision node to `references/authoring-guide.md` (Lookup vs Coalesce vs conditional-render vs Ternary).
3. Record Coalesce-vs-Ternary benchmark numbers in the authoring guide (or drop the claim entirely).
4. Fix the stale `lookupFeedbackStyle[T]` comment in `utils/darkmode_scan.go:32`.
5. Consider a `Coalesce` godoc Example (godoc examples are a repo convention for exported API; Coalesce has tests but no runnable Example).
6. Check whether `references/consumer-guide.md` should mention `Coalesce` for consumer-side templ authors (it's exported API).
7. Sweep FEATURES `utils/format` section: note the `compactUnits` extraction (its prose still describes the pre-extraction twins implicitly).
8. Re-record the art-dupl baseline after the duplication reduction (shrinks accepted groups; log binary version per ADR-0009).

**Docs hygiene noticed en route:**
9. `docs/agent-context-history.md:221` still lists `lookupFeedbackStyle[T]()` as a live helper — archive-accurate? It's a _history_ file; verify it reads as historical, annotate if it reads as current.
10. `docs/planning/2026-05-20_comprehensive-plan.md` T11.1 (`LookupStyle[T]`) — long-shipped; harmless, but plan files could get an "executed" annotation sweep via docs-health ANNOTATE.
11. FEATURES utils table lacks `DismissScript`'s new callers? (verify — Alert/Toast/ErrorAlert; confirm table completeness while there).
12. SKILL.md "utils — shared helpers" bullet and FEATURES table now both list Coalesce — add it to the `TestDocsCountDrift`-adjacent guard surface? (Only if drift ever bites; no evidence it will.)

**Testing depth:**
13. Property test for `Coalesce` (returns first non-zero; zero iff all zero) — small gopter-style test; current table covers the space informally.
14. Golden test explicitly covering the 3-chain aria-label collapse (`LoadingDialog` with only Title set) — chain semantics changed shape, not output; one targeted assert would pin intent.
15. Fuzz `Coalesce` with `FuzzInputType`-style harness? Probably overkill — note as rejected unless a bug ever surfaces.
16. Add `enterKeyHintValue`'s Lookup rewrite to the existing input branch-coverage test sweep (verify branch coverage didn't drop — coverage gate passed, but the specific branch map could be checked).
17. Run the full `nix run .#verify` (not just the loop) at the next natural commit point to exercise the complete lane end-to-end.

**Process/tooling friction observed:**
18. Master is **ahead 11 of origin, unpushed** (this session + sibling session). Push decision needed (g3).
19. Daemon commit fragmentation — curate or accept (g3).
20. The class-inventory guard's fail-once-regenerate UX burned 3 test runs this session — consider a `-regen.inventory` flag or doc note so the first FAIL reads as expected behavior.
21. `-bench` without `-run` runs guard suites — a repo note in AGENTS "Testing" could prevent repeat confusion.
22. `scripts/ci-repro.sh` generation block is copy-paste gold — nothing to fix; noted as the canonical source (the AGENTS one-liner differs subtly; maybe cross-link them).
23. The system PATH `templ` matched the pin this time by luck of the nix profile — consider a version assert in ci-repro.sh's generate step (belt: `templ --version | grep -q v0.3.1020`).
24. `_sources` mirror duplicated my 26 edits as self-healing sync — worked as designed; no action.

**Generics-adjacent ideas (idea-grade, demand-gated):**
25. `display.List[T]` — parked on ROADMAP demand gate; do not build without a second consumer ask.
26. Typed pagination/result-page generic (`Page[T]`) — same demand-gate class; note only.
27. `wire.DecodeForm[T]` adoption breadth: audit remaining hand-rolled `r.FormValue` parsing sites for DecodeForm candidates (ParseKanbanMove uses it; are there others?).
28. Generic `CoalesceOr[T]` (first non-zero else explicit fallback) — rejected for now: last-arg-literal convention makes it redundant.
29. Generic slice helpers (`Map`/`Filter`) — rejected: repo has no samber/lo by closed dependency budget; render-path shapes don't need them.
30. Enum-set generic type (`utils.EnumSet[T]`) — rejected (Go can't enumerate constants; per-enum tables already minimal); revisit only if IsValid boilerplate grows.
31. `withDefaults` generic (idea #80) — rejected; struct-literal zero-value rule is the better contract.
32. Evaluate whether `feedbackRole`'s if-branch and similar two-way policies could be table+Lookup — rejected: structural if-branch is the documented convention for structural variants.

**Bigger repo items I touched the edges of (for the roadmap, not this lane):**
33. TODO #443: golangci-lint pin alignment (dev-shell 2.15.x vs CI v2.14.0 skew) — referenced by AGENTS; still open upstream of my lint runs.
34. templ v0.3.1070 unpinned bump + golden re-baseline window (TODO #335) — untouched, decided-deferred.
35. `internal/contract` component count guards — my change added no components; counts unaffected (verified by green guard).
36. Coverage baseline: new Coalesce/compactUnits code shifts statement counts — guard passed, but the recorded baseline file may deserve a `--record` refresh at the next natural point (toolchain-pinned).
37. Website prose mentions utils helpers? (unverified this session — docs-health VERIFY candidate, not researched now per scope limit).
38. Demo pages: none of the converted sites render differently, so demo screenshots need no refresh — noted to prevent a wasted `nix run .#shots` run.
39. AGENTS.md size budget: my aria bullet grew the file slightly — still under the 220-line gate (verified by the passing pre-commit-size path in CI? unverified locally this session; `wc -l AGENTS.md` is a 5-second check).
40. The sibling session's data-model-review HTML may propose overlapping type-system work — read it before the next refactor to avoid double-planning (I did not open it; scope limit).

**Parked hygiene (no evidence of pain, keep visible):**
41. `Ternary` call-site census after Coalesce adoption — how many remain, and are they all genuinely bool/computed shapes? (5-minute audit.)
42. `inputModeForType` guard comment ("deliberately not Lookup") — cheap insurance against a future wrong refactor.
43. `gridClass` doc comments could name the Lookup default-key form explicitly.
44. CHANGELOG entry wording: link the ADR? (No ADR exists for Coalesce; probably fine — a helper, not an architecture decision.)
45. `BenchmarkUtils` naming: one monolithic benchmark func vs per-helper — cosmetic.
46. Review whether `DismissScript` and `Coalesce` belong in a shared "idioms" doc section for contributors.
47. idea #43 `Slot` type — still parked.
48. idea #43-adjacent: `Card.Body`/`Table.Body` slot docs — untouched.
49. ROADMAP scroll/focus-preservation row — untouched (kanban papercut, unrelated but still open).
50. This report's section f harvest into TODO_LIST/ROADMAP per docs-health HARVEST (the skill's own closing rule — the actionable items above should be routed, not entombed here).

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **Should `Ternary` be soft-deprecated in the docs for the override-or-default shape?** All 26 converted sites previously used it; `Ternary` remains correct for bool→value and computed-value shapes (`HelpTextID`). Options: (a) doc guidance only — "override-or-default ⇒ Coalesce"; (b) add `// Deprecated:`-style framing to Ternary's godoc (too strong — it's not deprecated); (c) leave as-is. I can implement any of these but cannot decide the API-doc policy.

2. **Is golden-identity sufficient evidence for output-identical refactors, or do you want `nix run .#visual` as a standing ritual after any render-path edit?** I skipped the visual suite with the reasoning "byte-identical HTML ⇒ identical pixels." If your bar is "browser-proven, always," I'll add it to the checklist (and item f1 becomes mandatory, not optional).

3. **Master is ahead 11 of origin, and my session's work is smeared across 4 daemon commits with generic messages** (`a42ae953` utils Go, `0f89d2a1` templ sources, `13af2e67` generated+docs, `cbb0603c` docs). Local-only interactive rebase to reword/merge them into one readable commit before push — or push as-is and preserve the daemon trail? Rewriting local-only history is safe, but it's your history-curation call (past sessions have gone both ways).

---

**Self-review Q&A (brutal-self-review, condensed):** Did I lie? No — but I published one wrong number (d1) and corrected it. Ghost systems? One found: the `lookupFeedbackStyle` doc ghost (c1) — stale reference, no capability lost. Split brains? None created; Coalesce/Ternary/Lookup now have distinct, documented jobs (pending g1's doc convergence). Scope creep? Held — four candidate ideas explicitly rejected with reasons (a8). Removed something useful? No — only duplication; behavior pinned by pre-existing boundary tests. Test honesty? Full loop + lint ran; visual suite consciously skipped (b1, g2).
