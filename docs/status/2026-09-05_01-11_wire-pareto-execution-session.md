# Wire SDK Pareto Execution — Session Status Report

**Date:** 2026-09-05 01:11 CEST
**Session span:** 2026-09-04 ~22:20 → 2026-09-05 01:11 (single session, executing `docs/planning/2026-09-04_22-22_wire-sdk-pareto-execution-plan.md`)
**Repo tip at report time:** `origin/master` = `4d1308a` — **all 4 CI jobs green** (Lint, Build & Test, CSS Freshness, Visual Regression)
**Input plan:** 19 macro tasks (T1–T19), 91 micro tasks
**Headline:** 15 of 19 macro tasks fully done, master green with the entire wire SDK shipped and browser-proven; the release cut (T9) is staged but not executed; two gated WC tasks correctly stayed closed; the session was heavily contested by concurrent agent sessions and the auto-commit daemon, which cost several CI-red rounds that were all fixed forward.

---

## a) FULLY DONE (work verified complete and green)

| Task                                  | What shipped                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| ------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **T1 — Ship & CI green**              | `scripts/ci-repro.sh --lint` passed locally; root cause of the one CI red (CSS Freshness) found and fixed: Tailwind v4's automatic content detection scans `.go` files too, so the `"blur"` string literal in `utils/wire/wire.go:96` (`EventBlur`) generated a `.blur` utility in the demo CSS — fixed by committing the recompiled CSS. Pushed and watched to green.                                                                                                                                                                                                                    |
| **T2 — Browser-proof pack**           | `visualtest/wire_e2e_test.go`: two chromedp tests drive **real Chromium** against a page composed from production pieces — `layout.Base` (self-hosted htmx), `datastar.SDKScript` (pinned bundle served locally, readiness gated on the runtime's own `datastar-ready` event, discovered in the bundle), wired `display.Button`s. Both clicks verified: fragment lands in `#wire-htmx-out` (htmx) and `#wire-datastar-out` (Datastar via `wire.Handler` response headers). Both PASS in Chromium.                                                                                         |
| **T3 — Hardening pack**               | 6 wired-Button golden snapshots (htmx full/post, datastar, custom event, empty-URL inert, link+wire); `-race` green on utils/wire, display, examples/demo; `nix run .#css` byte-stable; `nix flake check` passed.                                                                                                                                                                                                                                                                                                                                                                         |
| **T4 — `wire.Handler`**               | `utils/wire/handler.go`: `Handler(PatchTarget, next)` middleware + typed `PatchMode` enum (**all 7 merge modes verified against the pinned bundle**: inner/outer/prepend/append/before/after/replace) + `IsDatastar`/`IsHTMX` predicates + empty-Selector-degrades-to-id-matching. Table tests for datastar/htmx/plain callers. Demo `/api/wire/fragment` refactored onto it.                                                                                                                                                                                                             |
| **T5 — Invariant pack**               | `utils/wire/invariants_test.go`: empty-URL inert across **every** enum combination; URL present in BOTH dialects for every Method×Event×Target; Target never renders under Datastar; wire rendering never emits `<script>`; `BenchmarkActionAttributes` (htmx vs datastar); `FuzzAction` extended with Target (1.7M execs clean).                                                                                                                                                                                                                                                         |
| **T6 — Docs adoption pack**           | README wire section; SKILL.md `utils/wire` catalogue block; website guide `guides/transport-wiring.mdx` + astro sidebar + api-reference row; ADR-0035 annotated (partially superseded); `wire.go` package comment split into `doc.go`; DOMAIN_LANGUAGE 6 new terms (Transport, Action, Wiring, Dialect, Patch Mode); modularization README Layer-0 updated; ADR-0036 fresh-eyes verified against code + annotated with the Handler addition; FEATURES row; AGENTS wire bullet extended; CHANGELOG `[Unreleased]` warmed with 4 entries.                                                   |
| **T7 — Decision gates**               | `docs/wire-gates-d1-d2-d3.md`: three memos with recorded, reversible outcomes — D1: WC module **not ratified**; D2: **fallback stays** (pinned by invariants); D3: **transport-symmetric components only**, case-by-case.                                                                                                                                                                                                                                                                                                                                                                 |
| **T8 — Demo transport toggle**        | `?transport=both\|htmx\|datastar` switches every wire button's dialect server-side; segmented control (pure links, CSP-safe, `aria-current`); typed `demoTransport` enum with graceful fallback; 4-case test pinning dialect presence/absence per param.                                                                                                                                                                                                                                                                                                                                  |
| **T10 — Count normalization**         | Drift guard now the single source of truth; fixed stale claims: IsValid 31→**56**, visual goldens 89→**91**, per-package counts (display 40→42, feedback 13→14, datastar 3→4 in README/SKILL). `TestDocsCountDrift` green.                                                                                                                                                                                                                                                                                                                                                                |
| **T11 — LoadMore Wire**               | `navigation.LoadMore` gained `Wire *wire.Action`: cursor appended to `Wire.URL`, self-replacement preserved (htmx `hx-target="this"` / datastar id-matching), `InfiniteScroll` documented htmx-only and silently ignored under datastar, empty URL inert. Existing golden byte-identical (backward compat proven). 4 new goldens + 5 behavior tests.                                                                                                                                                                                                                                      |
| **T12 — Forms wired validation demo** | `/api/wire/validate` endpoint via `wire.Handler`; verdict fragments (success/error/neutral); htmx dialect uses the typed contract (input carries its own name/value); **datastar uses the documented `Attrs` escape hatch** (`data-bind:value` + interpolated expression) because the static-URL common subset cannot carry the value — scope-honest, tested both.                                                                                                                                                                                                                        |
| **T13 — ConfirmDelete decision**      | Bundle grep: **zero** `confirm` in pinned datastar v1.0.2 → `hx-confirm` has no twin; `htmx.ConfirmDelete` stays htmx-only. Recorded in the D3 memo.                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| **T14 — Research pack**               | Signaling notes in `docs/transport-wiring.md`: htmx 2.0.10 self-host bundle contains **no aria at all** (loading is CSS-class based — use htmx module components); datastar signals via `datastar.Indicator` + LiveRegion's `aria-busy`; decision: docs-not-helper. Interval/intersect: `data-on-interval`/`data-on-intersect` ARE in the pinned bundle and htmx has `every`/`revealed` — symmetric concepts, dialect-specific syntax → ADR-sized future work (TODO #152). Form-submit parity documented. Upstream `go-datastar/static` **v0.5.0 exists** (repo pins v0.4.0) → TODO #151. |
| **T15 — Migration + misc docs**       | `docs/recipes/transport-migration.md` (3-move recipe + what-does-NOT-migrate table); links from transport-wiring.md; testing notes extended (entity-encoded assertions, E2E as template); demo hero wire snippet; CONTRIBUTING transport-wiring convention row.                                                                                                                                                                                                                                                                                                                           |
| **T16 — HARVEST**                     | TODO_LIST: #150–#154 added (wire golden eyeball, static v0.5.0 bump, wire trigger language ADR, catalogue Wire survey, prerender sync). ROADMAP: wire trigger language + busy-signaling ideas.                                                                                                                                                                                                                                                                                                                                                                                            |
| **T17.1 — WC draft ADR**              | `docs/adr/0037-light-dom-web-components-module.md` — **status Proposed, not ratified**: light-DOM hosts only, host-API-not-components shape, leaf module, explicit "what kills this ADR" section. The D1 decision is now one word.                                                                                                                                                                                                                                                                                                                                                        |
| **T19 — Keep-green rule**             | Applied at every pack boundary; master is green at report time.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |

**Incident repairs (forced by the environment, all shipped):**

- go-datastar/static floated v0.4.0→v0.5.0 (runtime 1.0.3) via daemon tidy sweep — **pinned back**, stale go.sum checksums removed (`f5aba52`, `7f4b56a`).
- `website/package.json` typescript flip to ^7.0.2 (AGENTS-documented astro-check breaker) — **restored to ^6.0.3**.
- Wire demo toggle dark-mode gap (`bg-blue-600` without `dark:bg-blue-500`) — fixed per convention.
- 8 lint findings in the new wire files (canonicalheader Hx-Request, noctx, gocognit, golines, wsl×3) — fixed; utils module lints at 0 issues.
- Statcard + wire visual goldens regenerated after component drift (9 PNGs, `4d1308a`).

---

## b) PARTIALLY DONE

| Item                        | State                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **T9 — Release 1.13.0**     | Pre-flight COMPLETE (checklist read, `[Unreleased]` warm with the 4 wire entries, version ordering 1.12.0→1.13.0 sane, worktree `/tmp/tc-master-wt` created at green master, govulncheck located in nix store at **v1.6.0** — the checklist says the script pins **v1.7.0** and fails loud on a missing binary; this mismatch was NOT yet resolved). `scripts/release.sh 1.13.0` **not executed**. Post-release steps (tag lockstep, push, proxy wait, tidy sweep, post-CI) obviously pending too. |
| **T17 — WC module phase 1** | Only the draft-ADR micro-task (T17.1) was in scope per the plan's own recommendation; scaffold/API/tests (T17.2+) remain gated behind D1 ratification.                                                                                                                                                                                                                                                                                                                                             |
| **Demo prerender**          | `prerender.go` updated for the new `demoPage` signature (default "both" view), but the prerendered `index.html` artifact itself was not regenerated/verified in this session (TODO #154 records the follow-up).                                                                                                                                                                                                                                                                                    |

## c) NOT STARTED (correctly out of scope or gated)

| Item                                           | Why                                                                                                                                                           |
| ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **T17.2+ / T18 — WC module build-out**         | Gated by D1 = not ratified. The draft ADR-0037 makes a future "yes" cheap.                                                                                    |
| **T9 sub-steps 9.3–9.9**                       | Blocked behind executing release.sh (see b).                                                                                                                  |
| **go-datastar/static v0.5.0 bump**             | Deliberately deferred to TODO #151 (needs full bundle re-audit per the bump protocol).                                                                        |
| **Wire trigger language (interval/intersect)** | TODO #152 — ADR-sized, deliberately not smuggled into the contract.                                                                                           |
| **Catalogue Wire rollout beyond LoadMore**     | TODO #153 — survey task under the D3 rule.                                                                                                                    |
| **html-report / website build verification**   | No `npx`/node toolchain on this machine (bun-shim issue, AGENTS-documented); website validated by CI instead — Website workflow green on the wire-guide push. |

## d) TOTALLY FUCKED UP (honest failure log)

1. **The shared-checkout war.** At least three other agent sessions plus the auto-commit daemon worked the same checkout simultaneously. Branch HEAD changed under me **at least 5 times** (`fix/live-region-nonce-guard`, `-pr`, `fix/statcard-dl-structure`, `-clean`, `fix/alert-nonce-fallback`). Consequences:
   - My uncommitted static-pin downgrade and dark-mode fix were **clobbered twice** (daemon restores the tree to HEAD); I had to re-apply and re-commit.
   - My lint-fix commit landed on the **wrong branch** (`fix/live-region-nonce-guard-pr`) because another session switched branches between my `git add` and `git commit`; recovery required a cherry-pick that then **conflicted** (stale hunks — the PR branch's parent had older file versions), and I initially committed a half-resolved state before catching it.
2. **I caused the static v0.5.0 drift trigger.** Adding `go-datastar/static` as a direct import in visualtest + the daemon's tidy sweep floated ALL modules to v0.5.0 (breaking the pinned-contract tests). My pin-back fix then left stale go.sum checksums → CI's tidy check went red again. Two avoidable CI rounds.
3. **Wrong lint scope at pack end.** My "full" lint pass covered root packages (`display/forms/layout`) but not the utils sub-module that held my new files — CI's per-module lint caught 8 findings (canonicalheader, noctx, gocognit>30, golines, wsl×3). Should have been caught locally.
4. **Stale visual goldens committed.** My wire dual-transport goldens were captured before the toggle/layout changes landed, producing a 100%-dimension-mismatch failure in CI Visual Regression; the statcard goldens from the concurrent session's DOM fix had the same problem (~1–1.8% drift). One more CI round to regenerate.
5. **Master was red for ~25 minutes** across three failed pushes (`f5aba52` tidy-diff, daemon's `8c87cec` sweep, `33924875928` lint) before `5e9ebef` + `4d1308a` went green.
6. **Minor templ-parser stumbling:** putting literal `{…}` code text and `//` comment text into demo hero HTML text broke the templ lexer twice; burned three generate cycles bisecting it (text must go through Go string expressions; `//` cannot appear raw in HTML text nodes).

## e) WHAT WE SHOULD IMPROVE

1. **Serialize agent sessions on one repo.** The branch-checkout war produced nearly every failure above. Either one session at a time, or every session works in its own `git worktree` from the start.
2. **Fix the daemon (BuildFlow) — TODO #125/#126/#93 keep paying interest.** It re-swept `typescript@^7` (a change AGENTS explicitly documents as broken) and floated dependency pins. Its working-tree-restore behavior also destroys uncommitted agent work.
3. **Pin go-datastar/static everywhere explicitly.** The indirect require in visualtest should never float; `go get module@v0.4.0` in all three go.mods is now done, but a version-sync guard for this dep (like `check-version-sync.sh`) would prevent recurrence.
4. **Run the FULL per-module lint loop after every pack**, not a subset — CI's lint job covers modules my local loop skipped.
5. **Regenerate visual goldens as the LAST step of any pack that touches component or demo markup**, immediately before pushing — never commit goldens captured from an intermediate state.
6. **Commit fixes within the daemon's commit window.** Uncommitted work in this checkout has a half-life measured in minutes.
7. **Resolve the govulncheck version mismatch before cutting T9** (nix store has 1.6.0; the release script expects the pinned 1.7.0 — verify `scripts/release.sh`'s actual lookup before running).
8. **The release should be cut from a dedicated worktree** (started this way — keep it) so main-checkout churn can't poison the cut; also confirm no other session pushes to master mid-cut.
9. **The `TestWireSection` capture should pin its own CSS state** — it depends on the compiled demo CSS containing the wire classes; a `tc-wire`-scoped smoke assertion (classes present in compiled CSS) would fail fast instead of 100%-diffing in CI.

## f) NEXT 50 (ordered, with sources)

**Release block (do first):**

1. Resolve the govulncheck 1.6.0-vs-1.7.0 expectation in `scripts/release.sh` (install 1.7.0 or accept the nix-provided 1.6.0 consciously).
2. Execute `scripts/release.sh 1.13.0 "<summary>"` in `/tmp/tc-master-wt` (notes auto-extract from the warm `[Unreleased]`).
3. Verify the release commit tree: version files agree, replaces stripped from tagged go.mods (script 8b does this — re-verify by hand).
4. `scripts/check-release-tags.sh` — 9 tags in lockstep.
5. Push master + all tags.
6. Wait for proxy propagation: `go list -m github.com/larsartmann/templ-components/utils@v1.13.0`.
7. Post-propagation `GOWORK=off go mod tidy` sweep across all 8 modules; commit go.sum refresh (the v1.11/v1.12 red class).
8. Confirm post-release CI green; 24h daemon watch for CSS/pin regressions.
9. Merge/close the three sibling PR branches after their sessions finish (`fix/live-region-nonce-guard-pr` CI is red — their session must fix before merge; master already carries the equivalent busy-script work? NO — verify whether 8c41463's change is fully represented on master, else merge conflicts land later).
10. Human-eyeball `wire/dual_transport_{light,dark}.png` + regenerated statcard PNGs (TODO #150, agent-capture caveat).

**Wire rollout (post-release, per D3 rule):**
11. Catalogue survey for next transport-symmetric `Wire` candidates (TODO #153).
12. Consider `Wire` on `SimpleNav`/`Nav` links if semantics are symmetric.
13. Evaluate `Form` + `EventSubmit` demo block (htmx carries fields natively; datastar needs bound signals — escape-hatch docs exist).
14. Wire-aware `EmptyState` action button?
15. Wire-aware `StatCard` `Href`-vs-`Wire` decision.
16. Consider `LoadMore` browser E2E (extend `wire_e2e_test.go` with a cursor flow).
17. Add `wire.Handler` to the cmd/tc scaffolder's generated endpoint template.
18.godoc example for `wire.Handler` (plan T4.5 leftover — doc-comment example).
19. `wire.Action` JSON marshaling decision (config-as-data use case) — needs ADR-0036 annotation if added.
20. aria-busy helper reconsideration trigger: document what consumer demand would look like (ROADMAP idea exists).

**Contract hardening:**
21. TODO #152: draft the wire trigger-language ADR (interval/intersect, both bundles verified).
22. TODO #151: go-datastar/static v0.5.0 deliberate bump + full bundle re-audit per the bump protocol (#129).
23. Version-sync guard extended: go-datastar/static appears in exactly 3 go.mods at the pinned version (extends #138).
24. Golden sweep: assert `datastarScriptURL` output for custom base + default (TODO #137, still open).
25. Fuzz `wire.Handler` header combinations (extend beyond Attribute fuzzing).
26. Benchmark `wire.Handler` middleware overhead.
27. Property test: Handler never sets Datastar headers for htmx/plain callers even under header-smuggling (duplicate Datastar-Request headers).
28. Document `PatchMode` → fragment-id-matching interplay in transport-wiring.md FAQ.
29. Add `TestWireDemoTransportToggle` E2E variant (click the toggle links in Chromium).
30. Prerendered index.html regeneration check after demo changes (TODO #154).

**Docs:**
31. Website: wire guide page — add the segmented-control screenshot once demo deploys.
32. transport-wiring.md: link the website guide page (currently only repo docs cross-link).
33. DOMAIN_LANGUAGE: add `PatchTarget`/`Escape hatch` terms (missed this session).
34. SKILL.md: add a wire E2E pattern note (visualtest/wire_e2e_test.go as the template).
35. ADR-0036: add the "Verification log" section listing browser-proof + lint pinning (annotation exists; log format would help future audits).
36. docs/DOMAIN_LANGUAGE.md: wire entry for `wire.Handler` term itself.
37. Migration recipe: add the LoadMore example (first catalogue component with Wire).
38. Record the `.blur` Tailwind-scans-.go finding in AGENTS.md gotchas (new gotcha discovered this session, not yet written down).
39. Record the templ-lexer findings (`//` and literal `{}` in HTML text break parsing) in AGENTS.md.
40. Record the concurrent-session protocol lesson in AGENTS.md (worktrees, commit-fast rule).

**Hygiene:**
41. TODO #133: changelog-guard policy for test-only PRs (still open, owner decision).
42. TODO #128: exercise upstream-watch.yml via workflow_dispatch.
43. TODO #134: GOWORK=off cheat sheet into AGENTS.md (the per-module loop bit me again this session).
44. TODO #139: golines max-width policy (golines flagged invariants_test.go twice).
45. TODO #142: go.work `go` directive version-sync guard.
46. TODO #143: model visualtest's datastar dependency explicitly in check-module-layers.sh (my E2E added the direct dep — the layer script should know).
47. Sweep stale branches (`fix/live-region-nonce-guard`, `-pr`, `fix/statcard-dl-structure`, `-clean`, `fix/alert-nonce-fallback`) once their sessions land — branch list is accumulating.
48. Clean up `/tmp/tc-master-wt` worktree after the release cut.
49. Demos: `/api/save` invisible response (TODO #149, untouched).
50. Update `docs/status/2026-09-04_22-15_transport-wiring-sdk-session.md` lineage note to point at this report (snapshot chain).

## g) QUESTIONS I CANNOT ANSWER MYSELF

1. **Release timing:** Master is green at `4d1308a` and the release is fully staged — do you want me to cut **v1.13.0 now** from the isolated worktree, or hold until the three concurrent sessions' in-flight work (live-region PR is red, statcard/alert branches still landing) has merged, so 1.13.0 contains everything and the tag lockstep runs once?
2. **The concurrent sessions:** Are the other agent sessions (live-region nonce guard, statcard dl structure, alert nonce fallback) yours and intentional? If yes, should I **pause my repo-mutating work** (and the release) until they finish, or are they stale leftovers whose branches I may close out after harvesting?
3. **D1 — the WC module:** I recorded "not ratified" as the D1 outcome (ADR-0037 stays _Proposed_, the consumer recipe remains the answer). Is that your actual position, or do you want to ratify ADR-0037 and have me build the light-DOM `wc` module (T17.2+/T18) in a follow-up session?

---

**Snapshot lineage:** `docs/status/2026-09-04_22-15_transport-wiring-sdk-session.md` → `docs/planning/2026-09-04_22-22_wire-sdk-pareto-execution-plan.md` → **this report**.
**Canonical living backlog:** `TODO_LIST.md` (this session added #150–#154).
**Format note:** user explicitly requested `.md`; markdown skill defaults overridden per instruction-wins rule.
