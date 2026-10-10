# Session Status: Extraction-Backlog Plan Execution (2026-10-10 13:53)

**Task:** Execute `docs/planning/2026-10-10_06-06_extraction-backlog-pareto-master-plan.md`
(the "GO" the previous session ended on). 15 of 16 todos completed; the final
verification ritual was interrupted mid-run. Everything is daemon-committed on
local `master` (28 commits ahead of origin, **NOT pushed**).

---

## a) FULLY DONE — with its full testing ladder verified

| ID | Deliverable | Proof |
|----|-------------|-------|
| D1 | **ADR-0044** semantic accent tokens DECIDED (`docs/adr/0044-semantic-accent-tokens.md`): `--color-accent-50..950` family (literal-blue defaults, independence by construction) in `templates/custom.css`; 212-site inventory; Phase 1 (complete ADR-0008 alias layer, v1.x, non-breaking) vs Phase 2 (class swap = T1, gated on v2/ADR-0039); token-ready rule effective NOW (no new hardcoded `blue-*` accent sites — all new components this session comply); verified against DiscordSync `input.css:184` needs line-by-line (ADR §Verification) | ADR file + TODO #397 annotation |
| D2 | **ADR-0045** CSS class delivery DECIDED (`docs/adr/0045-css-class-delivery.md`): release-generated class inventory `templates/templ-components-classes.txt` as the blessed default; tracked-source `@source` blessed for monorepos; compiled whole-library CSS rejected. **Deciding evidence found during D2a:** DiscordSync's vendor-dir `@source` is INERT under Tailwind v4.3 (their `vendor/` is gitignored; v4.3 never scans gitignored paths — dnsblockd-proven) | ADR file + TODO #397 annotation |
| F1/F2 | **`utils/format`** package: Bytes (IEC, boundary-bump fix), CompactDuration (operator style), ClockDuration, Percent (NaN→"—"), StringOrDash, CompactCount. Table tests + bench (~40–90ns/op) + examples + docs rows; lint 0 | `go test ./utils/format` green; golangci 0 |
| C1/C1T | **`display.CodeBlock`** (block + CompactID variants, CopyButton composition, `NoCopy` inverted flag, HTML-escaped) — 7 goldens, demo section, docs counts, tc mirror, contract registration | all tests green |
| M1 | **`layout.MetaRefresh`** (zero-JS reload/redirect, negative clamp) — 2 goldens | green |
| P1doc | **`docs/integration/extraction-analysis.md`** — per-project adoption maps (4 projects), cross-cutting findings, linked from TODO_LIST + the plan | written + linked |
| C2/C2T | **`forms.FilterBar`** — DiscordSync filterForm productized: composite change trigger (recipe footgun #2 encoded), Reset, noscript Apply, sticky shell, `NoPushURL`, dual-transport `Wire` (consumer action never mutated — pinned). 7 goldens + **browser e2e BOTH transports** (`TestWireE2EFilterBarAutoSubmits`, real Chromium via nix: open→change→result→second change, htmx + Datastar) | e2e 1.3s green |
| C3 | **`forms.FilterChips`** — zero-JS link chips, `aria-current`, labelled group, `FilterToggleHref` builder (sibling-preserving, escaped, non-mutating — all pinned) | 2 goldens + tests green |
| C4 | **`display.StatusDot` + `LivePill`** — `StatusTone` enum (5 tones + IsValid), pulse ring with motion-reduce fallback, sr-only Label semantics | 6 goldens + tests green |
| C5 | **`display.DataState`** — 4-rung honesty ladder composing EmptyState; `DataStateState` + IsValid; unknown→content fallback | 4 goldens + tests green |
| C6 | **`forms.SegmentedControl`** — radio mode (forward DOM order, peer-checked, bare `checked` attribute — caught and fixed the `checked="false"` all-checked bug class) + link mode; RTL scan test | 3 goldens + tests green |
| C7 | **`display.SegmentBar`** — flex-grow proportions (no rounding gaps), deterministic label-keyed palette, legend, composed aria-label; inline-style exemption allowlisted; **+2 pixel visual goldens (light/dark) baselined via nix** | green |
| C8 | **`display.Lightbox`** — native `<dialog>` viewer (thumbnails/trigger, wrap-around nav, captions, rotate, zoom), CSP singleton script, joined `integration` nonce render table; **browser e2e** (`TestLightboxE2E`: open, initial state, next+caption, rotate transform, wrap-around, close) | e2e green |
| C9–C11 | `Table.StickyHeader` (scroll-container thead pin), `StatusBadgeWith(mapper, status)` + exported `MapStatusToBadgeType`, `ImageProps.Placeholder` (zero-JS blur-up) — all zero-golden-regression (defaults byte-identical) | green |

**Cross-cutting, done per component:** contract registration, demo sections
(nonce-threaded), docs-count bumps in the SAME commit (README/FEATURES/AGENTS/
skill/ROADMAP/comparison/website ×3 surfaces), `cmd/tc` `_sources` mirrors +
`packageDeps` registration, CHANGELOG `[Unreleased]` warm (13 entries).

**Verification actually completed this session:** full root `go test ./...` →
0 FAIL; all 6 sub-modules `-short` green; `nix run .#lint` → **0 issues across
all modules** (after fixing 43 autofixable findings + 2 wrong nolint tags);
targeted browser runs (2 e2e + 1 visual baseline) green.

**Counts moved:** 124→**135 components** (+MetaRefresh, +StatusBadgeWith),
layout 10→12, display 44→51, forms 23→27; enums 65→**68** (65→67 with IsValid);
HTML goldens 279→**320**; visual goldens 200→202; generated files 126→135.

## b) PARTIALLY DONE

1. **Final CI ritual (`scripts/ci-repro.sh --lint --website`) INTERRUPTED** —
   background run was canceled mid-flight (tool context cancel); the visible
   output was only its pre-flight tree diff (the daemon hadn't committed the
   picsum-seed rename yet). **No VERDICT line was ever produced.** Per M03,
   nothing is "CI-green" until it prints `VERDICT: PASS (exit 0)`.
2. **Known-red CI lanes I can already name** (see d):
   CSS freshness + route goldens will fail until recompiled/re-baselined.

## c) NOT STARTED (deliberate, per plan scope)

- **T1** — ADR-0044 Phase 2 implementation (token class-swap batches, visual
  rebaseline, migration doc). Version-gated on v2 per the ADR.
- **T2** — ADR-0045 implementation (inventory generator + release.sh wiring +
  freshness guard + adoption-guide rewrite + consumer migration notes).
- **X1** — consumer wave 1 (dnsblockd pin bump v1.19.4→current, SidebarNav
  swap, DataState adoption; DiscordSync ExternalLink/CollapsibleSection).
- **X2** — consumer wave 2 (DiscordSync chart migration + bridge deletion
  post-T1; nsfw Lightbox + mirror deletion post-T2; mr-sync ⫱).
- **G1/G2** — JSONTree/Dropzone: parked by demand gate (unchanged, correct).

## d) TOTALLY FUCKED UP (recoverable, but CI-red as of THIS commit tip)

1. **Demo CSS is STALE — CI's CSS lane WILL fail.** Verified: 0 occurrences of
   `animate-ping` / `scale-[2]` / `blur-lg` (etc.) in
   `examples/demo/static/app.css`; every new component's utility classes are
   missing from the compiled artifact. `TestCSSFreshness` only warns locally
   (no CI env) — which is exactly why my green local runs hid it. I violated
   the plan's own anti-verschlimmbesser line ("demo CSS recompiled after …
   variant classes"). **Fix: `nix run .#css`, commit, re-verify.**
2. **Route goldens NOT re-baselined — Visual lane WILL fail.** `/display` and
   `/forms` are pinned light+dark (`visualtest/route_golden_test.go:241–274`);
   my demo sections guarantee pixel drift. **Fix: after CSS recompile,
   `nix run .#visual -- -update` (route tier), eyeball, commit.**
3. **Axe sweep + touch-target + zoom-reflow audits never ran on the new demo
   sections.** The sweeps audit live routes automatically in CI — new surfaces
   (FilterBar selects, chips, segmented radios, lightbox controls) are
   unaudited. Risk: `axe_baseline.json` entries needed or fix-forward.
4. **The ci-repro run was orphaned** (background shell canceled) — its partial
   state may include a poisoned go-build cache entry (I hit one earlier this
   session: `open /mnt/buildcache/... no such file`; retry self-healed).

## e) WHAT WE SHOULD IMPROVE (session lessons)

1. **Make CSS recompilation part of EVERY component's "tail" checklist** — I
   followed the ladder doc faithfully except this line; the guard that exists
   (TestCSSFreshness) is local-warning/CI-fail, i.e. invisible exactly when
   iterating. Propose: fail loud locally when `nix` is available, or move the
   recompile into `nix run .#build`.
2. **Route goldens are a per-demo-page obligation too** — same class of miss.
   The plan's ladder said "demo page" but not "route goldens + axe follow".
3. **Daemon races cost ~6 repair cycles** (packageCount reverted twice,
   FEATURES spacing, demo.templ, mirrors). The AGENTS warning exists; what's
   missing is a cheap post-write verification step for high-churn files
   (counts, mirrors) before moving on.
4. **Two lint-tag typos (`//nolint:exhaustruct` vs `exhaustruct_v5`)** — both
   caught by lint, but I wrote them from memory instead of copying the
   adjacent file's exact tag.
5. **`templ.KV` is not a ternary** — burned a build cycle on the thead sticky
   class; `utils.Ternary` is the tool. Also `checked={bool}` renders
   `checked="false"` (still ON in HTML) — boolean attributes need the
   if-block form. Both now demonstrated; candidates for a skill note.
6. **Synthetic Escape key events do not close `<dialog>` in headless
   Chromium** (probed: focus confirmed on in-dialog button; KeyEvent ignored
   by the cancel pipeline). Documented in the e2e; UA-native behavior stays
   untested there — consistent with Modal/Drawer precedent.
7. **Integration doc written before its sections shipped** —
   extraction-analysis.md still says "🔨 C2/C3/C4/C5/C7/C8 (in plan)" while
   all of them now ship. Needs one sweep pass (not done; listed below).

## f) NEXT — up to 50, in execution order

**Un-red the CI (blocking, ~30–45 min):**
1. `nix run .#css` — recompile demo CSS with all new utility classes; commit.
2. `nix run .#visual -- -update` for the route tier ONLY (`-run TestSiteRouteGoldens`… actual: `TestDemoRoutes`-class names — check `route_golden_test.go`); eyeball diffs (new sections expected).
3. Full `nix run .#visual` (all tiers incl. axe sweep, touch-target, zoom-reflow) — fix forward any critical/serious finding on new sections or ledger it.
4. `scripts/ci-repro.sh --lint --css --visual --website` → demand `VERDICT: PASS (exit 0)`.
5. `git status -sb` for daemon races; re-verify tip; **push** (needs your go — see q1).
6. Watch CI on GitHub; the Visual Regression + Website jobs are the risk lanes.

**Session hygiene (~20 min):**
7. Sweep `docs/integration/extraction-analysis.md` statuses (🔨→✅ shipped).
8. Mark plan Level-1 rows + TODO_LIST #388–400 done-shipped entries.
9. Add the two templ gotchas (KV-ternary, checked={bool}) to the skill's authoring rules.
10. Classify the art-dupl advisory note from the prior session (TODO backlog item).

**T2 — ADR-0045 implementation (next feature lane, ~1–2h):**
11. `scripts/gen-class-inventory.sh` (all 7 modules, deterministic order).
12. Wire into `scripts/release.sh` + freshness guard test (TestClassInventoryFreshness).
13. Add `templates/templ-components-classes.txt` to TestCompiledCSSInventory's known set (non-CSS note).
14. Rewrite `docs/tailwind-v4-adoption-guide.md` quickstart around the inventory.
15. `tc init` scaffold: emit the copy step.
16. Consumer migration notes (nsfw mirror deletion steps; dnsblockd switch; DiscordSync inert-line replacement).
17. Website docs page section for CSS delivery.

**T1 — ADR-0044 Phase 1 (v1.x-safe part, ~1–2h):**
18. Complete ADR-0008 alias layer shades (blue-50..900 minus already-aliased, amber-300/400/700, red/green tints).
19. `docs/recipes/scoped-theme-bridge.md` (the DiscordSync pattern, documented).
20. Re-audit `templates/custom.css` starter + `cmd/tc/_sources/starter/` mirror sync.

**T1 Phase 2 (only after q2 answered):** class-swap batches display→forms→
feedback→layout/navigation (21–24), guard-test regex updates (25), demo CSS
recompile + FULL visual rebaseline (26), migration doc + theme-bridge recipe
update (27), CHANGELOG + version plan (28).

**Consumer waves (other repos, each with its own verify ritual):**
29. X1a dnsblockd: bump templ-components v1.19.4→current, tidy, build, test.
30. X1b dnsblockd: SidebarNav swap; delete dashboardSidebar.
31. X1c dnsblockd: DataState adoption; delete ladder copies + 4 art-dupl markers.
32. X1d dnsblockd: adopt StatusBadgeWith; delete 4 local mappers.
33. X1e dnsblockd: adopt MetaRefresh; delete hand-rolled meta tags.
34. X1f DiscordSync: adopt ExternalLink + CollapsibleSection.
35. X2a DiscordSync: adopt CodeBlock (copyIDRow → CompactID variant).
36. X2b DiscordSync: replace inert vendor `@source` with the inventory file (post-T2).
37. X2c nsfw: adopt FilterChips (historyChip), SegmentedControl (presets), SegmentBar (evidence), format.Percent/CompactDuration.
38. X2d nsfw: delete `_css-scan` mirror + sync script (post-T2).
39. X2e mr-sync: decision ⫱ (q3) + migration kickoff if approved.
40. All three stale pins land current release; per-repo commits/PRs.

**Library follow-ups from this session's components:**
41. FilterBar + FilterInput composition doc (self-wired children nest forms — document the boundary).
42. SegmentedControl: consider `Description`/help text field if consumers ask.
43. Lightbox: keyboard arrow navigation between images (left/right when open) — native dialog gives Escape only.
44. CodeBlock: optional line-numbers variant if a second consumer asks (gate).
45. `utils/format`: adopt inside library demos (RelativeTime/CountBadge stay as-is; use format in demo stat cards).
46. Dark-viewport spot-check: run `nix run .#shots` on /display + /forms after CSS recompile (manual eyeball).
47. Benchmarks for the new components (repo convention: bench suites in 7 packages).
48. BDD lens for FilterBar/Lightbox (behavior specs; currently render+e2e only).
49. Consider `EnsureID` prefixes audit (tc-lightbox/tc-seg/tc-bar/tc-chips all golden-normalized — verify no collisions with consumer namespaces).
50. Update `docs/count-definitions.md` with the new baseline numbers (135/68/320/202).

## g) Three questions I cannot answer myself

1. **Push authorization.** Local master is 28 commits ahead (daemon-committed,
   content verified by me lane-by-lane, but the final ci-repro VERDICT is
   missing and two lanes are knowingly red until f.1–f.4). The previous
   session's explicit push grant was for that session. **Should I complete
   f.1–f.5 and push master now, or hold for your review of the diff first?**
2. **T1 Phase 2 timing.** ADR-0044 gates the accent class-swap on v2/ADR-0039
   (which is itself trigger-gated with no date). DiscordSync/nsfw/dnsblockd
   keep maintaining CSS bridges until then. **Do you accept shipping the swap
   as a v1.x minor with a migration note (breaking-ish for CSS-targeting
   consumers), or hold for a real v2.0 cut?**
3. **mr-sync (X2e, owner gate ⫱).** Recommendation stands: migrate to Tailwind
   + the library (small dashboard, near-total duplication). **Approve the
   migration, or keep mr-sync on its hand-rolled semantic CSS?**

---

**Bottom line:** all 13 planned builds shipped with their ladders (incl. two
browser e2e suites and pixel baselines), 0 lint findings, 0 test failures
locally — but the release-critical "compiled artifacts match sources" invariants
(demo CSS, route goldens) were NOT refreshed, so the tip is CI-red until f.1–f.4,
and nothing has been pushed. Awaiting instructions.
