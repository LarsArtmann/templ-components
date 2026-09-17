# Status Report — TODO-sweep hardening session (12 items)

**Date:** 2026-09-17 09:28 CEST
**Scope:** Single-session execution of the open actionable TODO_LIST items (#193, #198-#202, #221, #223, #225-#228, #230, #231), prioritized #231-first because it blocked ALL manual commits. Self-review + harvest included.
**Method:** load skills → investigate before editing → fix forward → prove every new gate fires (break-it/restore) → verify per touched module → warm CHANGELOG/TODO_LIST/AGENTS in the same session.
**Verification at close:** root `go build` + `go test ./...` green; `utils`, `datastar`, `visualtest` per-module green (visualtest full browser suite 11.3s); root lint 0 issues; visualtest lint 0 issues (was 21 findings); `scripts/check-lint-config.sh` green; ci.yaml YAML-valid; gofmt clean; `go.mod` still `go 1.26.7`; git tree clean (daemon committed as usual).

---

## a) FULLY DONE (12 TODO items, all verified)

| # | Item | Proof |
|---|------|-------|
| 231 | **Toolchain-skew fix** — root `go.mod` reverted 1.27.1 → 1.26.7. Root cause correction: the bump was NOT from the v1.17.0 release (as TODO believed) but daemon auto-commit `7252433f` TODAY, which also nudged `flake.lock` (nixpkgs `b1b875982b`). Dev shell verified unaffected (go 1.26.7 + templ v0.3.1020, zero-diff invariant intact). New guard `utils.TestGoDirectiveSkew` pins module-directive ≤ go.work for every `use` entry (derived dynamically — new modules auto-covered); **proven to fire** by reintroducing the skew (direct run fails naming the file), restored, green. Fires in CI's per-module `GOWORK=off` loop — the only place it can (workspace mode: the toolchain error preempts `go test` entirely). AGENTS.md gotcha bullet added. | `utils/infra_guards_test.go`, CHANGELOG Fixed, AGENTS.md |
| 198+199 | **Kanban LSP triage** — no-op confirmed: `golangci-lint run ./display/...` = **0 issues**; the LSP warnings are the already-documented stale-gopls noise. #199 closed as documented (AGENTS line exists); upstream filing left as optional owner action. | lint output |
| 200 | **Datastar scripts-are-inert fact** in `datastar/doc.go` — new "Patched fragments carry no executable scripts" section sourced from `docs/datastar-runtime-facts.md` §"Fragments carry no executable scripts". | module build+vet+tests green |
| 202 | **`.fail/` hygiene** — `TestMain` in `visualtest/main_test.go` prunes `testdata/.fail/` at run START; stale cross-session dirs + empty subdirectories (carousel/polloedregion/routes were sitting there, 0 bytes) proved gone after one run. Failure artifacts are now guaranteed to be from THIS run; CI upload path unaffected (prune precedes writes). | live prune observed |
| 193 | **Poll helpers** — `visualtest/poll.go`: `pollBool` (wraps predicate in `Boolean(...)`, `*bool` capture), `pollTrue` (verdict-in-error), `pollText` (string-expr, empty keeps polling). The bool-into-string chromedp mistake is now unwritable in the helpers' shape. Migrated: `demo_flows_e2e_test.go` (5 sites), `focus_preservation_e2e_test.go` (4 sites, ternary noise dropped), `calendar_nav_e2e_test.go` (3 sites). Removed 2 dead symbols on sight (`uploadEchoMarker`, `regionExistsExpr`). | compile + full visualtest suite |
| 225 | **visualtest lint lane + findings** — CI Lint job runs `(cd visualtest && golangci-lint run --timeout=5m ./...)`. All **21 findings fixed forward**: siteshots (static `errSearchNoHits` sentinel, `ListenConfig.Listen`, `Fprintf(os.Stdout,…)`, checked `server.Close`, named `screenshotQuality`/`navigationSettle` consts, prealloc via literal-embedded action, gosec G703 nolint on the true anchor line after extracting `cleanPath`), `demo_kanban_http_test.go` (`NewRequestWithContext` ×2, `strings.Cut`, stale-nolint removal, wsl), golines module-wide via `golangci-lint fmt`, and a commented `.golangci.yml` exclusion scoping gocognit/unparam to `_e2e_test.go` harness fixtures. **Module now lints 0 issues.** | lint run ×4 iterations to zero |
| 226 | **Contract-marker cheat-sheet** — "Demo contract markers (cheat-sheet)" table in `docs/visual-testing.md`: PORT override, `/health` JSON, `data-tc-kanban*` selectors, hidden-form field names (card/column/index — verified from `display/kanban.go` constants), CSRF input markup, all 6 kanban endpoints, 403/404 semantics. One correction caught during authoring: `#kb-htmx`/`#kb-ds` are e2e-harness board ids, NOT demo ids — row fixed before commit. | per-row source verification |
| 227 | **Cross-binding kanban parity** — shared builder impossible (visualtest boots the demo as an external binary), so ONE probe table: `runKanbanContractProbes` in `demo_kanban_http_test.go` drives `TestDemoKanbanHTTPContracts` (real demo binary) + new `TestKanbanE2EHTTPContractParity` (harness server) through identical probes (CSRF 403 ×2, same-origin 403, unknown-column 404, add 200, move 200, reset, datastar variants, + new datastar-reset probe so shared global boards return to initial layout). **Drift-catch proven**: broke the harness same-origin guard (403→200), parity test failed naming `e2e-parity`, restored, green. Cross-binding ANTI-DRIFT comments added to BOTH implementations (`examples/demo/kanban_demo.go` header + `kanban_e2e_test.go` header). | break-it/restore proof, both tests green |
| 221 | **`ExampleKanbanBoard_columnTone`** — compiling godoc example (todo=blue, blocked=red, done=green) next to `columnAction`; documents zero-value-toneless + aria-hidden dot. | `go vet ./display/` + TestKanban green |
| 223 | **Vision-review flagged set** — `kanban/section_action_tone_{light,dark}.png` added to `scripts/vision-review-goldens.sh` FLAGGED (with TODO ref in the header comment). | `bash -n` + file-existence check |
| 228 | **`docs/planning/TEMPLATE.md`** — the plan-authoring checklist made structural (per-task gate checkboxes: goldens-cover-this / wired⇒e2e-or-waiver / counts-bump-same-edit / demo-smoke-gate, ⫱ owner-gate convention, Pareto phase ordering, whole-plan verification block, follow-through seeds). Cross-linked from `docs/plan-authoring-checklist.md`. | written |
| 230 | **Exhaustive recipes audit** — all 36 recipe files; **194 unique package-qualified identifiers** checked against live APIs via `go doc -all` per package. **3 real drifts fixed**: `polled-region.md` `display.Alert`/`AlertError` → `feedback.Alert`/`FeedbackError` (v2 alias removal, ADR-0022); `custom-404-page.md` `NotFound404Link` → `NotFoundLink`; `login.md` `FormMethodPost` → `FormPost`. All other apparent misses verified NOT drift (external go-datastar SDK symbols, guard-test names, underscore-constant truncations like `DatastarVersion1_0_3`). Website copy (`website/content/docs`, 71 idents) audited with the same method — clean (10 apparent misses are all `utils.Test*` guard references, legit). | audit loop + triage |

**Bookkeeping done in the same session:** CHANGELOG `[Unreleased]` +10 entries (9 Added, 1 Fixed-pair), TODO_LIST rows closed for all 12 items, AGENTS.md +1 daemon-gotcha bullet (go-directive bumps).

---

## b) PARTIALLY DONE

- **#193 migration breadth.** Helpers shipped + flow tests migrated (the item's literal scope), but **~48 raw `chromedp.Poll` sites remain** in `kanban_e2e_test.go`, `wire_e2e_test.go`, `wire_form_e2e_test.go`, `wire_forms_pack_e2e_test.go`, `loading_button_e2e_test.go`, `polled_region_e2e_test.go`, `datastar_runtime_e2e_test.go`, `datastar_synthetics_e2e_test.go`. The mistake is still writable in those files.
- **#225 lane coverage.** visualtest lane added; the **website module remains unlinted** (same gap shape, deliberately out of the TODO's scope).
- **#198/#199.** Verified-noise + documented locally; the **upstream gopls/templ-LSP issue was never filed** (needs owner GitHub voice + a minimal repro outside this repo).

## c) NOT STARTED (session-scoped, deliberately)

- #189 file-backed kanban demo state + Dashboard-recipe kanban section — feature-scale, needs its own plan.
- #222 website kanban guide page — feature-scale (SSG page + nav + goldens + CSP re-hash).
- #224 sorted-view demo board + 422 e2e — feature-scale.
- #194/#195/#196 route goldens (dark/mobile/RTL) and #197 keyboard-only traversal — need `nix run .#visual` capture sessions.
- All Blocked/owner-decision items untouched (#80, #28/#29, #93-family, #107/#108/#124-#126, #190-#192, #211, #212, #162, #150, #213-#217).

## d) TOTALLY FUCKED UP — nothing unrecovered, but three honest near-misses

1. **False-green audit loop (caught).** The first #230 audit used `read - ref` (invalid option) — the loop never iterated and printed `TOTAL_MISSING=0` from a dead gate. The `read:` warning + the prove-the-gate discipline caught it; the corrected loop found 12 candidates. This is the *same* class as the AGENTS "pipeline masking" lesson — it almost recurred in the same session that documented it.
2. **Three-iteration extraction bugs.** Even after fixing the loop: (i) `go doc` extraction missed grouped consts (indented in `const (...)` blocks) → 50 false "missing"; (ii) the refs regex `[A-Za-z0-9]*` truncated underscored constants (`DatastarVersion1_0_3` → `DatastarVersion1`) → false positives needing manual triage. A correctly-designed extractor (underscores + indented const parsing) would have been one pass.
3. **Mid-edit daemon snapshots.** Daemon commits captured half-finished states (e.g., `1f004206` committed 3 lines of the in-progress guard; later commits caught compile-broken intermediates). One multiedit misplacement (`var png` inserted in the wrong function) survived my "applied 10/11" glance and surfaced as `undefined: png` at build — I fixed it, but I should view-after-multiedit when any edit reports failure. Also left `infra_guards_test.go` unformatted until end-of-session (daemon committed the unformatted version mid-way).

**Flagged autonomous decision (not fuckup, but ratify):** the `.golangci.yml` gocognit/unparam waiver for `_e2e_test\.go` applies repo-wide, not just visualtest — currently zero root-module files match; acceptable, but it is broader than the TODO asked.

## e) WHAT WE SHOULD IMPROVE

1. **Make the recipes audit a CI guard** (mechanize #230 with an allowlist for external-SDK/test symbols/underscore suffixes) — doc drift now recurs silently between sweeps.
2. **Migrate the remaining ~48 Poll sites** onto `pollBool`/`pollTrue`/`pollText` — then ban raw `chromedp.Poll` in visualtest via forbidigo or review convention.
3. **Lint the website module** (same lane shape as #225) — it is the last Go module without CI lint.
4. **Fix BuildFlow's go-directive bump path** (#93-family, separate repo) — today's #231 is its 7th documented incident class; the new test only detects after the daemon lands damage.
5. **Ratify the daemon's nixpkgs lock nudge** (`flake.lock` → `b1b875982b`): harmless today (go/templ unchanged), but wholesale input bumps are supposed to be deliberate.
6. **TestMain ordering note:** the parity test leaves shared package-global kanban boards reset; browser e2e tests mutating them under `-shuffle` is a pre-existing coupling worth a per-test board instance eventually.
7. **AGENTS convention line** for the poll helpers (like the existing transition-constants bullet) so future sessions use them instead of raw Poll.
8. **gopls-only warnings** in `visualtest/options_test.go` (nilness ×2) and `datastar_runtime_e2e_test.go` (writestring ×3) — lint passes them; fix or waive consciously at next touch.

## f) NEXT UP TO 50 (ranked-ish: owner gates ⫱ first, then Pareto)

1. ⫱ Ratify/rollback the daemon's `flake.lock` nudge (`b1b875982b`).
2. ⫱ Decide v1.18.0 release timing — `[Unreleased]` now ~20 entries (see g-1).
3. ⫱ Ratify the `_e2e_test.go` gocognit/unparam config waiver breadth (d-flag).
4. Decide fate of `templates/styles.css` + theme `.out.css` (#211 — evidence complete: recommend delete).
5. Policy for ~2.1k unannotated name-keyed report items (#212 — recommend (b) or (c)).
6. #192 release-timing supersession bookkeeping.
7. Migrate remaining Poll sites in `kanban_e2e_test.go` (~12).
8. Migrate remaining Poll sites in `wire_forms_pack_e2e_test.go` (~10).
9. Migrate `wire_e2e_test.go` + `wire_form_e2e_test.go` Poll sites.
10. Migrate `datastar_runtime` + `datastar_synthetics` Poll sites.
11. Migrate `loading_button` + `polled_region` Poll sites.
12. Recipes-identifier drift-guard test (mechanized #230).
13. Website module CI lint lane (+ fix its findings, if any).
14. #189 file-backed kanban demo state.
15. #189 Dashboard-recipe kanban section.
16. #222 website kanban guide page (seed from recipe; CSP re-hash + site goldens).
17. #224 sorted-view demo board + 422-rejection e2e.
18. #194 dark route goldens (6 routes).
19. #195 375px mobile route goldens (4 routes).
20. #196 RTL route goldens.
21. #197 keyboard-only demo traversal (Tab-order UX).
22. #229 session-scoped CSRF store for the demo (conditional).
23. #80/#162/#150 run `scripts/vision-review-goldens.sh` (needs API key) + human confirm SUSPECTs — flagged set now includes kanban action/tone.
24. #28 awesome-templ PR + #29 templ.guide listing (one sitting).
25. #216 vnu ignore re-triage on next nixpkgs html5validator bump (event-gated).
26. #213 PR wall-clock budget comment (needs baseline artifact infra).
27. #214 PR benchstat comment (needs bench.old storage).
28. #215 gremlins mutation pilot on utils.
29. #217 M24/M25 demand-check re-run when next survey lands (MultiSelect, DateRangePicker, FileDrop, command palette, toast positions, TreeView).
30. #93 BuildFlow: honest daemon commit messages (separate repo).
31. #107 BuildFlow preflight jsonv2 scan fix (separate repo).
32. #108 BuildFlow eslint-fix scoping (separate repo).
33. #124 BuildFlow: stop re-appending `*_templ.go` to .gitignore (separate repo).
34. #125 BuildFlow: provider flag against CSS un-minification (separate repo).
35. #126 BuildFlow: commit classifier for vetted artifacts (separate repo).
36. #190 KanbanBoard touch-drag story (owner Q1; vibe-kanban handle pattern documented).
37. #155 SimpleNav Wire transport symmetry (D3: demand check first).
38. #157 Calendar month-nav Wire candidate (D3 gate).
39. #178 typed interval/intersect triggers in `wire.Event` (needs ADR).
40. #33 `Validate() error` on remaining props structs (only where invalid states are representable).
41. #34 test helpers → `internal/testutil/` (large mechanical, post-v1.0 deferred).
42. #154 keep prerendered wire view in sync if demo grows (event-gated).
43. #120 CSS-recompile false negative — re-open only if CI/local disagree again.
44. #119-note remove the bun shim at `~/.local/bin/node` (user-level pnpm fix).
45. File the stale-gopls/templ-LSP diagnostics upstream (#199 residue, owner voice).
46. Fix or waive the gopls nilness/writestring warnings in visualtest (e-8).
47. AGENTS.md convention bullet for the poll helpers (e-7).
48. Per-test kanban board instances instead of package globals (e-6).
49. #156 consumer-adoption follow-ups live in consumer repos (owner action).
50. Next docs-health/status harvest pass over this report's §f.

## g) THREE QUESTIONS I CANNOT ANSWER MYSELF

1. **Cut v1.18.0 now or batch more?** `[Unreleased]` holds ~20 entries incl. today's 12. The release script makes a cut ~15 min; or batch #189/#222/#224 into it. This is your release-cadence call (#192 lineage).
2. **Keep or revert the daemon's `flake.lock` nudge (`b1b875982b`)?** I verified the dev shell is unchanged (go 1.26.7, templ v0.3.1020) so it is currently harmless — but it was NOT a deliberate input bump. Reverting restores the pre-incident lock; keeping it accepts the daemon's choice. Only you can ratify which.
3. **Should the recipes identifier audit become a permanent CI drift-guard** (mechanized with its allowlist: external go-datastar SDK symbols, `Test*` guard names, underscore-suffixed constants), or stay a manual periodic sweep? Cost = maintaining the allowlist against false positives; benefit = doc drift dies at PR time instead of recurring (today: 3 stale references shipped for unknown weeks).
