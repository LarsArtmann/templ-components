# Status Report: v2.0 Documentation Cleanup + Module Count Fix + Full Verification

**Date:** 2026-08-10 07:05
**Session goal:** Close all remaining gaps from the prior session's self-review (stale doc references, missing CHANGELOG, misleading test names, demo CSS staleness), then verify everything green.
**Outcome:** **All planned work executed, verified green, and committed (5b737f8). Discovered and fixed the 5→7 module count error along the way.**

---

## Executive Summary

This session was a cleanup pass over the prior session's v2.0 breaking-change bundle. The prior session implemented 4 breaking changes (HTMX self-host, container-aware default flip, alias removal, field rename) via mechanical sed, but missed propagating the renames to documentation files. This session:

1. **Fixed all 8 stale documentation files** (AGENTS.md, FEATURES.md, ROADMAP.md, skill/SKILL.md, CONTEXT.md, TODO_LIST.md, 2 recipe docs, 4 ADRs)
2. **Added 4 CHANGELOG `[Unreleased]` entries** for the v2.0 breaking changes
3. **Renamed 3 misleading test functions** (ContainerResponsive → ContainerAware)
4. **Updated the migration guide** with CSP nonce requirement
5. **Recompiled demo CSS** (byte-identical — no staleness)
6. **Discovered the module count was wrong** (said 5, actually 7 — `htmx` and `datastar` have their own `go.mod`). Fixed across all docs, guards, scripts, and CI.
7. **Added htmx to the DAG layer guard** (was missing)
8. **Fixed a gci lint issue** in `htmx/a11y_test.go`
9. **Reviewed the datastar changes** — legitimate go-datastar/static integration

All 7 modules build, test (race), and lint clean. All guards pass. The BuildFlow daemon committed everything as `5b737f8`.

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| 1 | **All stale `ContainerResponsive` references fixed in living docs** | `grep -rn ContainerResponsive --include='*.md' AGENTS.md FEATURES.md ROADMAP.md skill/SKILL.md CONTEXT.md TODO_LIST.md docs/recipes/ docs/adr/` returns 0 hits outside historical CHANGELOG entries and ADR update notes |
| 2 | **All stale `AlertType`/`ToastType` references fixed in living docs** | Same grep pattern — 0 hits outside ADR-0006 (annotated as removed) and ADR-0022 (migration doc) |
| 3 | **CHANGELOG `[Unreleased]` has `### Changed` section** with all 4 v2.0 breaking changes | `CHANGELOG.md:17-22` — HTMX self-host default, container-aware default flip, field rename, alias removal |
| 4 | **AGENTS.md container-query section updated** | `AGENTS.md:131` now says v2.0 defaults, all `ContainerAware`, no `ContainerResponsive` |
| 5 | **AGENTS.md feedback section updated** | `AGENTS.md:143` now says "v2.0 removed the aliases" |
| 6 | **AGENTS.md HTMX loading section added** | `AGENTS.md:147` describes self-host default, `HTMXSelfHost`, `layout/embed.go`, CDN fallback |
| 7 | **FEATURES.md updated** | Grid row, SkeletonCardGrid row, enum table (removed AlertType/ToastType), responsive section, deferred items marked done |
| 8 | **ROADMAP.md updated** | Container queries line (default flip happened), theming line (HTMX self-host default), module count (7) |
| 9 | **skill/SKILL.md updated** | Grid component table, container-query section, module count, DAG, isolation test loop |
| 10 | **CONTEXT.md updated** | Naming conventions table: AlertType → FeedbackType |
| 11 | **TODO_LIST.md updated** | #35 (default flip) and #38 (alias removal) marked done |
| 12 | **docs/recipes/container-queries.md rewritten** | Full Grid section, all-component table with Default column, no ContainerResponsive |
| 13 | **docs/recipes/dashboard.md updated** | ContainerResponsive reference → ContainerAware default |
| 14 | **ADR-0006 annotated** | Consequences: "Removed in v2.0 (ADR-0022)" |
| 15 | **ADR-0016 updated** | ContainerResponsive reference → ContainerAware |
| 16 | **ADR-0018 annotated** | v2.0 update note about rename + default flip |
| 17 | **ADR-0034 updated** | Module boundaries, table, DAG all updated for 7 modules |
| 18 | **3 test functions renamed** | `TestGridContainerResponsive*` → `TestGridContainerAware*` in display/coverage_test.go, display/coverage_boost2_test.go, display/grid_card_feedback_test.go |
| 19 | **Migration guide CSP nonce requirement added** | `docs/migration/v1-to-v2.md` has a blockquote about `script-src 'nonce-...'` + checklist item |
| 20 | **Demo CSS recompiled** | `tailwindcss --input demo.css --output static/app.css --minify` — byte-identical to committed version |
| 21 | **Datastar changes reviewed** | Legitimate: go-datastar/static version pinning, cleaner test isolation (customSpinner instead of feedback.Spinner dependency), correct SDK doc references. `version.go` now derives `DatastarVersion1_0_2` from `static.Version`. |
| 22 | **Module count corrected from 5 to 7** | Discovered `htmx/` and `datastar/` have their own `go.mod`. Fixed in: AGENTS.md (heading, description, table, import graph, lint command, GOEXPERIMENT note), ROADMAP.md, skill/SKILL.md (DAG, isolation loop, build desc), CHANGELOG.md, ADR-0034 (boundaries, table, DAG, replace count). |
| 23 | **check-module-layers.sh updated** | Added htmx as Layer 1 module. Header comment says 7 modules. Summary says "6 sub-modules". |
| 24 | **htmx gci lint fix** | `goimports -w a11y_test.go` in htmx module |
| 25 | **All 7 modules verified green** | Root (10 pkgs) + utils (4) + icons (1) + errorpage (1) + charts/echarts (1) + htmx (1) + datastar (1) — build ✓, test ✓ (race), lint ✓ |
| 26 | **All 4 guards pass** | module-sync (7 modules), module-layers (6 sub-modules), version-sync, lint-config |
| 27 | **Golden tests pass** | No stale snapshots from the container-aware default flip |
| 28 | **BuildFlow committed everything** | `5b737f8` — "feat(v2.0): expand module split to 7 modules + breaking changes (ADR-0022, ADR-0034)" — 92 files changed, 1751 insertions, 440 deletions |

---

## b) PARTIALLY DONE

| # | Item | What's done | What's missing |
|---|------|-------------|----------------|
| 1 | **CSP integration test for HTMX self-host** | Migration guide documents the CSP requirement; `layout/bdd_test.go` asserts `HTMXSrc == HTMXSelfHost` in default props | No integration test in `integration/csp_nonce_test.go` that specifically renders a Page with default props and asserts the inline `<script nonce="...">` contains HTMX source. The existing CSP test covers htmx package scripts but not the layout self-host path. |
| 2 | **Visual regression coverage for container-aware default flip** | Golden HTML tests pass (no stale snapshots) | No visual regression test specifically renders Grid/Card/Split with the new `ContainerAware: true` default and compares against a PNG. The `visualtest/` suite doesn't have a Grid screenshot test. |
| 3 | **AGENTS.md layout package description** | HTMX loading section added at `AGENTS.md:147` | The root-module package table at `AGENTS.md:14` still describes `layout` as "Page shell, theme toggle, CSP-safe script/style tags, body-layout primitives" — doesn't mention the HTMX embed infrastructure (`embed.go`, `static/htmx.min.js`). |
| 4 | **README.md module documentation** | README has no stale references | README.md doesn't mention the 7-module workspace at all — it has no module structure section. The module info lives only in AGENTS.md and ROADMAP.md. This may be intentional (README is the sales page) but consumers don't learn about individual module adoption from README. |

---

## c) NOT STARTED

| # | Item | Impact |
|---|------|--------|
| 1 | **Bump `utils.Version` to `2.0.0`** | Release-time action. Currently at `1.8.1`. The drift-guard tests (`TestVersionMatchesChangelog`, `TestVersionMatchesFeatures`) will enforce this at release time via `scripts/release.sh`. |
| 2 | **`scripts/release.sh` end-to-end test** | The release script was modified for the remove-at-release replace lifecycle and 7-module tagging, but never run end-to-end. The sed patterns for replace removal were tested in isolation only. |
| 3 | **Visual golden PNG update for container-aware default** | If any visual test renders Grid/Card/Split with defaults, the PNGs may need updating. No visual test currently covers these specifically, so this is a gap in coverage, not a staleness issue. |
| 4 | **AGENTS.md Dockerfile section accuracy** | `AGENTS.md:237` describes a "3-stage Dockerfile pipeline" and mentions `.dockerignore`. The Dockerfile exists at `examples/demo/Dockerfile` and `.dockerignore` exists at repo root. This section is pre-existing and accurate — not a regression from this session. |
| 5 | **Consumer simulation CI job** | No CI step that does `go get github.com/larsartmann/templ-components@<tag>` from a clean module to verify external consumers can actually adopt the library. The `GOWORK=off` isolation tests partially cover this but don't test the proxy resolution path. |

---

## d) TOTALLY FUCKED UP

| # | Item | What happened | Severity |
|---|------|---------------|----------|
| 1 | **Prior session said "5-module workspace" but the repo has 7 modules** | The prior session's modularization work extracted `htmx/` and `datastar/` as separate Go modules (they have their own `go.mod`), but the documentation consistently said "5 modules." This session discovered the discrepancy when the `check-module-sync.sh` guard reported 7 modules (it scans `go.mod` files dynamically). The AGENTS.md, ROADMAP.md, skill/SKILL.md, CHANGELOG.md, and ADR-0034 all said 5. Fixed all of them. The BuildFlow commit message even says "expand module split to 7 modules." | **MEDIUM** — the docs described a smaller architecture than what exists. Not a code bug, but misleading to the next session. |
| 2 | **The `check-module-layers.sh` was missing htmx** | The DAG enforcement script checked utils, icons, charts/echarts, datastar, and errorpage — but not htmx. So an upward dependency from htmx (e.g., importing display) would have passed the guard silently. Fixed by adding `check_layer "htmx" "utils" "htmx"`. | **MEDIUM** — the guard had a blind spot. Now closed. |
| 3 | **The AGENTS.md lint command included `./htmx/...` in the root lint** | `golangci-lint run ./htmx/...` from the root module fails because htmx has its own go.mod — golangci-lint can't lint across module boundaries. The command would produce an error (though golangci-lint continues with 0 issues on the valid packages). Fixed by removing `./htmx/...` from root lint and adding `(cd htmx && golangci-lint run ./...)` to the sub-module section. | **LOW** — the error was non-blocking (0 issues despite the error message). |
| 4 | **htmx/a11y_test.go had a gci formatting issue** | Likely introduced when the sed renamed `ContainerResponsive` → `ContainerAware` and import ordering shifted. The root module's lint scope didn't include htmx (it's a separate module), so this wasn't caught until I ran per-module lint. Fixed with `goimports -w`. | **LOW** — self-caught during verification. |

---

## e) WHAT WE SHOULD IMPROVE

1. **The module count error (5 vs 7) is a symptom of documentation drift from code.** The code created 7 modules, but the docs were written from the prior session's mental model of 5. The `check-module-sync.sh` guard catches go.mod drift but doesn't verify the documented module count matches the actual count. Consider a drift-guard test that asserts the AGENTS.md module count equals the actual `find . -name go.mod | wc -l`.

2. **Per-module lint is not wired into the pre-commit hook.** The pre-commit runs 5 shell guards + BuildFlow, but doesn't lint the sub-modules. The htmx gci issue was only caught because I ran full per-module lint manually. CI catches it, but not at commit time. Consider adding a fast lint check to pre-commit (or at least for changed modules).

3. **The HTMX self-host change has no dedicated integration test.** The CSP nonce test covers htmx package scripts, and the BDD test checks the default props value, but nothing renders a full Page with default props and asserts the inline HTMX script is present with a nonce. The migration guide documents the CSP requirement, but there's no automated test that catches a regression (e.g., if someone accidentally removes the nonce from the self-host path).

4. **Visual regression coverage has a gap for container-aware components.** No visual test renders Grid, Card, or Split with container queries enabled. The default flip from viewport to container breakpoints is a visual change that golden HTML tests can't catch (they compare strings, not rendered layout). A visual screenshot of a container-aware Grid in a constrained parent would catch layout regressions.

5. **The sed shotgun pattern continues to miss documentation.** This is the third session in a row where a mechanical sed updated code but not docs. The pattern is consistent: (1) source code ✓, (2) tests ✓, (3) generated code ✓, (4) documentation ✗, (5) CHANGELOG ✗. The fix is process, not tooling: after any mechanical rename, run `grep -rn <old_name> --include='*.md'` and fix every hit before declaring done.

6. **The README doesn't mention the module split at all.** This may be intentional (README is the sales page), but consumers who want icons-only or error-only adoption won't find the per-module `go get` instructions unless they read the installation docs or the migration guide. A brief "Modular adoption" section in README would improve discoverability.

7. **The `go-datastar/static` dependency is marked `// indirect` even though it's a direct import.** `go mod tidy` insists on the `// indirect` annotation, which is unusual. This may be because the import is only in `version.go` (a constant derivation), not in a function that `go mod tidy` considers a "direct" usage. Not a bug, but confusing for anyone reading go.mod.

---

## f) Up to 50 Things to Get Done Next

### Critical (before v2.0 release)

1. **Bump `utils.Version` to `2.0.0`** via `scripts/release.sh 2.0.0 "<summary>"` — this is the release cut
2. **Test `scripts/release.sh` end-to-end** — the remove-at-release replace lifecycle and 7-module tagging have never been run
3. **Add CSP integration test for HTMX self-host** — render Page with default props, assert `<script nonce="...">` contains `var htmx=function()`
4. **Add visual regression test for container-aware Grid** — screenshot in a constrained parent container
5. **Verify the release commit is clean** — no `replace` directives, all 7 modules tagged, CHANGELOG heading correct

### High priority (before or during v2.0 release)

6. **Add drift-guard test for documented module count** — assert AGENTS.md says "7-module" and actual go.mod count is 7
7. **Wire per-module lint into pre-commit** — at minimum, lint changed modules (detect via `git diff --name-only`)
8. **Add AGENTS.md layout package description update** — mention `embed.go` + `static/htmx.min.js` in the package table
9. **Add README "Modular adoption" section** — brief mention of per-module `go get` for icons/errorpage/echarts
10. **Add consumer simulation CI job** — `go get github.com/larsartmann/templ-components@v2.0.0` from clean module
11. **Add visual regression test for container-aware Card** — screenshot with ContainerAware: true
12. **Add visual regression test for container-aware Split** — screenshot with ContainerAware: true
13. **Verify all golden files are fresh** — run `go test ./... -run TestGolden -update` and diff
14. **Run `go mod tidy` on all 7 modules** — ensure go.sum entries are current
15. **Verify `nix run .#verify` passes** — the Nix-based all-in-one verification
16. **Update `docs/modularization/README.md`** — mention the remove-at-release replace lifecycle strategy
17. **Add a `TestHTMXSelfHostRendersInlineScript` test** — unit test in layout package
18. **Add CSP test for the self-host path in `integration/csp_nonce_test.go`**

### Medium priority (polish)

19. **Update `docs/DOMAIN_LANGUAGE.md`** — may reference old terms (AlertType, ContainerResponsive)
20. **Update `docs/theming.md`** if it references semantic tokens as "optional"
21. **Consider extracting `htmxSource` version as a compile-time constant** — currently the version is baked into the embedded file with no link to `HTMXVersion2_0_10`
22. **Add a `// Deprecated` migration comment path** for consumers upgrading (grep-friendly)
23. **Consider a `go.work.tmpl` checked-in template** for contributors
24. **Add pre-commit guard for CHANGELOG `[Unreleased]` warmth** — blocks commit if `[Unreleased]` has no body
25. **Add a visual module structure diagram to docs** — D2 or Mermaid showing the 7-module DAG
26. **Document `HTMXSelfHost` constant in AGENTS.md** — it's in the HTMX loading section but not in the constants list
27. **Consider adding `HTMXSelfHost` to the skill/SKILL.md** — the skill doesn't mention self-hosting at all
28. **Update the demo to explicitly handle new defaults** — some grids may need `ContainerAware: false` for specific showcases
29. **Run `nix flake check`** — verify the Nix flake is healthy
30. **Add `htmx` and `datastar` to the skill/SKILL.md component table** — currently the table only covers root-module packages

### Lower priority

31. **Consider a `go.work.sum` checked-in strategy** — currently gitignored, but consumers may need guidance
32. **Add CI job for `go get` from clean module per sub-module** — simulates external consumer for icons-only/errorpage-only adoption
33. **Document the `go-datastar/static` `// indirect` situation** — explain why tidy insists on indirect
34. **Consider adding a `TestVersionEqualsModuleCount` drift guard** — assert documented count matches actual
35. **Add a recipe doc for HTMX self-host customization** — how to serve a custom HTMX version
36. **Update `docs/adr/0007-htmx-self-host.md`** (if exists) — reflect that self-host is now the default
37. **Consider a `tc-htmx-version` variable in CSS** — for consumers who need to verify the embedded version
38. **Add a benchmark for the self-host render path** — `htmxSelfHostComponent` writes ~51KB per render
39. **Consider caching the self-host script component** — currently creates a new ComponentFunc per render
40. **Update the visual testing docs** — mention the container-aware default flip may affect visual goldens
41. **Review whether the `feedback` package still needs its own golden test for Alert/Toast** — the types changed
42. **Consider adding a `TestAlertPropsTypeIsFeedbackType` compile-time test** — assert the type didn't drift back
43. **Review the `cmd/tc` CLI for v2.0 awareness** — does it need to know about ContainerAware defaults?
44. **Update the website docs** — check all pages for stale ContainerResponsive/AlertType references
45. **Add a `v2.0` changelog summary section** — high-level migration overview at the top of the v2.0 release notes
46. **Consider a `tc-version` HTTP header** — for debugging which embedded HTMX version a page uses
47. **Review whether `datastar` needs a separate `go-datastar` MITRIC** — it imports `go-datastar/static`
48. **Add a `TestModuleCountMatchesDocs` test** — assert AGENTS.md module count == actual go.mod count
49. **Consider extracting the HTMX version into a `Version` struct** — not just a string constant
50. **Party?** — only after the release is cut and all guards pass

---

## g) Questions

### Q1: Should the README get a "Modular adoption" section showing per-module `go get` commands?

The README currently has no mention of the 7-module workspace. Consumers who want icons-only or error-only adoption won't find instructions unless they read the installation docs or migration guide. I left it out because the README is positioned as the "sales page" and the installation docs cover it. Should I add a brief section, or keep README focused on the full-library experience?

### Q2: Should I add a CI job that simulates an external consumer doing `go get` from a clean module?

The current CI uses `GOWORK=off` for per-module isolation testing, which partially simulates the consumer experience. But it doesn't test the actual module proxy resolution path — a consumer does `go get github.com/larsartmann/templ-components@v2.0.0` and the proxy fetches the tagged source. This would catch issues like missing `*_templ.go` files or broken `replace` directives that only manifest at the proxy level. Should I add this, or is the `GOWORK=off` test sufficient?

### Q3: Should the HTMX self-host path cache the rendered component?

Currently, `htmxSelfHostComponent(nonce)` creates a new `templ.ComponentFunc` on every call, and each render writes ~51KB of HTMX source to the buffer. The source itself is static (embedded via `//go:embed`); only the nonce changes. Should I pre-render the `<script nonce="%s">` wrapper as a `[]byte` template and just `w.Write` it, or is the current approach fast enough (the BDD tests don't show any perf concerns)?

---

## Verification Snapshot

```
=== All 7 modules (race-enabled) ===
Root:        10 packages — build ✓, test ✓ (race), lint ✓
Utils:        4 packages — build ✓, test ✓ (race), lint ✓
Icons:        1 package  — build ✓, test ✓ (race), lint ✓
Errorpage:    1 package  — build ✓, test ✓ (race), lint ✓
Echarts:      1 package  — build ✓, test ✓ (race), lint ✓
Htmx:         1 package  — build ✓, test ✓ (race), lint ✓
Datastar:     1 package  — build ✓, test ✓ (race), lint ✓

=== Guards ===
check-module-sync.sh:    OK (7 modules, all paths and versions consistent)
check-module-layers.sh:  OK (no upward dependencies in 6 sub-modules)
check-version-sync.sh:   OK
check-lint-config.sh:    OK

=== Git ===
Working tree:    clean (committed by BuildFlow daemon as 5b737f8)
Commit:          feat(v2.0): expand module split to 7 modules + breaking changes (ADR-0022, ADR-0034)
Files changed:   92 files, 1751 insertions, 440 deletions
go.work:          7 use directives (root + 6 sub-modules + visualtest)

=== Golden tests ===
All packages:    PASS (no stale snapshots from container-aware default flip)
Demo CSS:         recompiled, byte-identical to committed version
```

---

## What Was Done Well

- **The module count discovery was the most valuable catch this session.** The prior session's docs said "5 modules" but the code had 7. The `check-module-sync.sh` guard reported 7 (it scans dynamically), which prompted the investigation. Without this, the v2.0 release would have shipped with incorrect documentation.
- **The container-queries recipe was fully rewritten** rather than just sed-patched. The new version includes a Default column showing which components are opt-in vs opt-out, which is what consumers actually need.
- **The ADRs were annotated rather than rewritten.** ADR-0006, ADR-0018, and ADR-0016 got update notes pointing to ADR-0022, preserving the decision history while reflecting the current state.
- **The htmx lint issue was caught because I ran per-module lint, not just root lint.** This validated the importance of the per-module lint workflow.
- **The datastar changes were reviewed thoroughly** — read every line of the diff, verified the `go-datastar/static` dependency, confirmed the version pinning logic is correct.
