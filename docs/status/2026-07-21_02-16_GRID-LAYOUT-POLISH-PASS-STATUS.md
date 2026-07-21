# Grid-Layout Session — Polish Pass: Brutal Status Report

**Date:** 2026-07-21 02:16 CEST
**Scope:** Session 5 — the post-implementation polish pass after shipping the
grid-first 2D layout primitives (commit `2314f26`). This report covers ONLY
what happened in this polish session, not the broader project.
**Author:** Crush (autonomous), awaiting Lars review

---

## TL;DR

The previous session shipped 4 new layout primitives (AppShell, Container,
Split, Stack) + Footer multi-column + Form.Layout enum in one commit. This
session was the cleanup pass: CHANGELOG blocker resolved, 1 real lying-comment
bug fixed, 3 new failing CI guards added (RTL scanner, `<main>` contract,
cross-package composition), 1 invalid icon reference in a recipe caught, 2 RTL
physical-property violations in the demo fixed.

**Build OK · 13/13 packages pass `-race` · 0 lint issues · demo binary starts.**

But: the polish pass itself introduced 2 near-misses (a wrongly-deleted
function I caught only by luck; a recipe referencing a non-existent icon I
caught only because the status-report prompt forced re-verification). The
uncommitted surface is now 16 files. **Nothing is committed.**

---

## a) FULLY DONE (verified: build + test -race + lint = 0)

| #   | Item                                                                                                                                                                              | Evidence                                                                               |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| 1   | **CHANGELOG `[Unreleased]` updated** — 10 Added entries + 1 Deprecated entry covering every grid-layout change                                                                    | `CHANGELOG.md` — release script will no longer fail                                    |
| 2   | **`splitRatioCols` lying comment fixed** — the function's godoc wrongly said `splitRatioMainSpan returns...` (introduced in `2314f26`, caught this session)                       | `layout/split_types.go:38-42`                                                          |
| 3   | **`FormProps.Inline` `// Deprecated:` confirmed** — was already present from prior session (`form.templ:100`)                                                                     | verified, no change needed                                                             |
| 4   | **`--tc-sidebar-w` documented in custom.css** — `:root` fallback (16rem) + full comment explaining the SidebarWidth enum mapping + specificity note                               | `templates/custom.css:418-442`                                                         |
| 5   | **Demo dogfooding** — hero + sticky-nav sections replaced hand-rolled `max-w-6xl mx-auto px-4 sm:px-6 lg:px-8` with `@layout.Container`                                           | `examples/demo/demo.templ:56,111`                                                      |
| 6   | **Cross-package integration tests** — `TestAppShellCrossPackageComposition` (AppShell+SidebarNav+Nav+Grid), `TestSplitWithContentAndAside`, `TestStackWithFeedbackComponents`     | `integration/appshell_composition_test.go` — all PASS                                  |
| 7   | **`<main>` singleton contract test** — `TestBodyPrimitivesDoNotEmitMain` table-driven over AppShell/Split/Stack/Container                                                         | `layout/a11y_test.go:40-74` — would have caught the original Split `<main>` regression |
| 8   | **RTL logical-properties scanner** — `TestRTLLogicalProperties` scans all `.templ` for `ml-`/`mr-`/`pl-`/`pr-`/`text-left`/`text-right`/`border-l-`/`border-r-` (failing CI test) | `utils/rtl_compliance_test.go` — PASS, 0 library violations                            |
| 9   | **2 RTL demo violations fixed** — `ml-4`→`ms-4`, `pl-10`→`ps-10`                                                                                                                  | `examples/demo/display_demo.templ:278`, `examples/demo/forms_section.templ:21`         |
| 10  | **Mobile-drawer recipe example** — full `display.Drawer` + hamburger button wired to `tcOpenOverlay`, ready to drop into `MobileNav` slot                                         | `docs/recipes/appshell-dashboard-layout.md:99-160`                                     |
| 11  | **ROADMAP.md updated** — counts (97→98, 37→43), new Layout pillar, `FormProps.Inline` added to v1.0 removal list, RTL scanner mentioned                                           | `ROADMAP.md:13-24,40`                                                                  |
| 12  | **Full verify** — `templ generate` (87 files) + `go build ./...` + `go test -race ./...` (13/13 ok) + `golangci-lint run` (0 issues)                                              | ran twice, once after each significant change                                          |
| 13  | **Demo binary smoke test** — builds, starts on port 18923, stops cleanly                                                                                                          | `PORT=18923 /tmp/tc-demo` (no HTTP fetch — see section d)                              |
| 14  | **Invalid icon reference caught & fixed** — recipe used `icons.Cog6Tooth` which **does not exist**; replaced with `icons.Settings`                                                | `docs/recipes/appshell-dashboard-layout.md` — caught during status-report verification |

---

## b) PARTIALLY DONE

| Item                                               | What's done                                                        | What's missing                                                                                                                                                                                                                               |
| -------------------------------------------------- | ------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Demo binary smoke test**                         | Process starts, binds port, exits cleanly                          | **Never fetched a page** — did not verify the layout section renders, did not check for `AppShell`/`Container` output in the served HTML. A real smoke test would `fetch http://localhost:18923/` and assert the layout-demo section exists. |
| **form_layout_test.go whitespace normalization**   | multiedit reported "Applied 3 of 5 edits"; tests pass; lint passes | 2 edits silently failed (non-unique `old_string`). The file is consistent enough to pass lint, but I did not re-verify the file is style-identical to existing test files. Residual `wsl_v5` LSP warnings may or may not be stale.           |
| **CSS subgrid research note** (from prior session) | Tracked, decision documented (track-only)                          | No prototype, no Card/DefinitionList alignment proof                                                                                                                                                                                         |
| **`--tc-sidebar-w` documentation**                 | `:root` fallback + comment added                                   | **No test** asserts the variable is present; potential load-order conflict with consumer `:root` is mentioned but not stress-tested                                                                                                          |

---

## c) NOT STARTED (this session)

- **Committing any of the 16 changed files.** Nothing is committed. All work is in the working tree.
- **Pushing to remote.** (House rule: never push without explicit instruction.)
- **Cutting a release.** CHANGELOG is warm, but `utils.Version` is still `0.18.1` and no tag exists.
- **Recompiling demo CSS.** `examples/demo/demo.css` scans `.templ` files; I changed `demo.templ` (Container wrapping) but did **not** rebuild `examples/demo/static/app.css` (still dated Jul 14). The committed CSS is stale relative to the new Container wrappers. Docker rebuild is fine; `go run` is not.
- **Running `nix fmt`.** Changed `.go` files; `golangci-lint` covers Go, but `treefmt-nix` is the CI format gate (`nix flake check`). Did not run it.
- _\*Example_* compile-tests_* for the recipe code samples (plan F10.3, deferred).
- **Updating `website/src/content/`** with a dedicated layout-primitives page.
- **README "layout primitives" section** — README still describes the old 6-component layout package.
- **Visual RTL manual test** of AppShell + Split under `dir="rtl"`.
- **Accessibility manual test** of AppShell/Split with a screen reader.
- **Benchmark suite** for the 4 new primitives (display/feedback/navigation/forms/layout/htmx/icons/utils all have benchmarks; layout does not).

---

## d) TOTALLY FUCKED UP (near-misses, caught by luck or by re-verification)

### d.1 — I almost deleted `splitRatioCols` (a live function)

**What happened:** The LSP flagged `splitRatioCols` as `unused`. I trusted the
warning and deleted the function with an `edit`. Build still passed (because
the deletion removed both the function AND the generated-file caller was not
yet regenerated). I almost committed.

**How I caught it:** After deleting, I ran `grep splitRatioCols` across `.go`
AND `.templ` files. The function was referenced in `layout/split.templ:35` and
`layout/split_templ.go:59`. The LSP gives **false-positive `unused` warnings**
for any function called only from generated `*_templ.go` files, because
golangci-lint's path exclusions skip `_templ.go`.

**Root cause:** I trusted the LSP over a grep. The same false-positive pattern
affects `containerWidthClass`, `sidebarWidthValue`, `splitRatioMainSpan`,
`stackGapClass`, `formLayoutClass` — ALL flagged `unused`, ALL actually live.

**Lesson for AGENTS.md:** `unused` warnings on package-private functions in
this repo are unreliable for any function that might be called from a
`.templ` file. Always `grep -rn fnName *.templ *_templ.go` before deleting.
I restored `splitRatioCols` immediately with a corrected comment.

**Process failure grade:** 🔴 Critical. Without the post-delete grep, the next
`templ generate` would have broken the build, and if committed, the Go module
proxy would have served a non-compiling tag.

### d.2 — Recipe referenced a non-existent icon (`icons.Cog6Tooth`)

**What happened:** In the mobile-drawer recipe, I wrote
`{Label: "Settings", Href: "/settings", Icon: icons.Cog6Tooth}`. `Cog6Tooth`
is a real Heroicon name but **is not in this library's icon set** (only
`icons.Settings` exists).

**How I caught it:** Only because the status-report prompt said "report based
on what you noticed" and I decided to re-verify the recipe's icon names. The
recipe would have failed to compile if a consumer copied it verbatim.

**Root cause:** I wrote example code from Heroicons memory instead of
checking `icons/icon_names.go`.

**Lesson:** Every identifier in recipe code must be grep-verified against the
actual source. A recipe that doesn't compile is worse than no recipe.

**Process failure grade:** 🟡 Medium. Caught before commit, but only by luck
of the status-report prompt triggering re-verification.

### d.3 — I chased stale LSP `wsl_v5` ghosts for ~6 tool calls

**What happened:** The LSP showed 40 `wsl_v5` "missing whitespace" warnings
across my new test files. I spent multiple `multiedit` and `view` calls trying
to fix them. Then I checked existing clean test files
(`integration/composition_test.go`) and found the **exact same** `output := ...; if ...`
pattern with **0 warnings**. The LSP cache was stale.

**Root cause:** I didn't first check whether the warning was real by comparing
against an existing passing file. `golangci-lint run` (the CI gate) reported 0
issues the entire time.

**Lesson:** The CLI gate is authoritative. When LSP and CLI disagree, trust
the CLI, and restart the LSP rather than editing code.

**Process failure grade:** 🟢 Low (wasted ~5 minutes, no code damage).

### d.4 — `rtl_compliance_test.go` had a REAL `wsl_v5` failure I didn't catch until `golangci-lint run`

**What happened:** I wrote the RTL scanner, ran `go test` (passed), declared
it done. Then ran `golangci-lint run` for the final verify — it caught a
missing blank line before `t.Errorf`. I had to `lint --fix`.

**Root cause:** I ran the test but not the linter on the new file before
moving on. The test-passes-≠-lint-passes gap is exactly the D1 mistake from
the prior session's status report — and I repeated it.

**Lesson:** After writing ANY new `.go` file, run `golangci-lint run` on it
specifically before marking the task done. `go test` is not enough.

**Process failure grade:** 🟡 Medium (repeated a documented prior mistake).

---

## e) WHAT WE SHOULD IMPROVE (process, this session)

1. **Never trust LSP `unused` for package-private funcs in a templ repo.**
   Add to AGENTS.md: "LSP `unused` is a false positive for any function called
   only from `*_templ.go`. Grep `.templ` + `_templ.go` before deleting."
2. **Run `golangci-lint run` on every new file immediately**, not just at the
   end. `go test` passing ≠ lint passing (wsl_v5, gci, etc.).
3. **Verify every identifier in recipe/example code against source.** Icons,
   prop names, enum values — all must be grep-confirmed.
4. **Don't chase LSP warnings when the CLI gate is clean.** Restart the LSP
   instead. The CLI config is the source of truth.
5. **Commit incrementally.** 16 uncommitted files is too much surface for one
   commit. The lying-comment fix, the RTL scanner, and the CHANGELOG should be
   separate commits for bisect-ability.
6. **Recompile demo CSS after changing `.templ` files** if the demo is expected
   to run via `go run` (not just Docker).
7. **Fetch the actual HTTP response in a smoke test**, not just "process alive".
   A 200 with the right HTML is the only real proof.

---

## f) Up to 50 things to do next (prioritized)

### 🔴 Blockers / must-do before release

1. **Commit the 16-file polish surface** (after user picks commit granularity — see Q2).
2. **Recompile `examples/demo/static/app.css`** from `demo.css` so `go run` shows the new Container wrappers correctly (Docker rebuilds it; local dev does not).
3. **Decide release version** (v0.19.0 minor now, or hold for v1.0 — see Q1) and cut via `scripts/release.sh`.
4. **Run `nix fmt` + `nix flake check`** — the CI format gate. I changed `.go` files; treefmt-nix may want to reformat.
5. **Commit or trash the prior session's status report** at `docs/status/2026-07-21_01-24_GRID-LAYOUT-SESSION-BRUTAL-STATUS.md` (untracked).

### 🟠 Correctness / hardening

6. **Real HTTP smoke test**: `fetch http://localhost:PORT/` and assert the layout-demo section + `AppShell`/`Container` classes appear in served HTML.
7. **Verify `form_layout_test.go` is style-consistent** after the partial multiedit (3 of 5 applied). Re-read and align with existing test file style.
8. **Golden snapshot tests** for AppShell, Container, Split, Stack — the `internal/golden` package exists; new primitives have none.
9. **Benchmark suite for layout package** — every other package has `benchmark_test.go`; layout doesn't.
10. **`Validate()` on `AppShellProps`** (warn if `Content == nil`) and `SplitProps` (`Main == nil`) — graceful guard, not a panic.
11. **Stress-test `--tc-sidebar-w` `:root` fallback** for consumer-`:root` load-order conflicts.
12. **`TestDarkModeCompliance` + `TestMotionReduceCompliance`** explicitly run on the new primitives (they scan all `.templ`, so should pass — but verify, don't assume).
13. **Contract test for `ComponentProps` interface** — verify the 4 new props types appear in `componentTypes()` (prior session added them; verify no drift).
14. **CSP nonce check** on any inline scripts in new primitives (AppShell has none — verify the `integration/csp_nonce_test.go` still passes with them included).

### 🟡 Features / extensions

15. **AppShell desktop sidebar collapse** (icon-only collapse button) — common admin pattern.
16. **Split responsive collapse** — stack columns vertically below a breakpoint automatically.
17. **Stack direction prop** (horizontal flex-row variant) — though that's arguably just `flex flex-row gap`.
18. **Stack.Divider prop** — render `<hr>` or custom divider between children.
19. **Container.ContainerResponsive** — container-query variant (like `Grid.ContainerResponsive`).
20. **Container width "Between"** — e.g. `max-w-screen-xl` option.
21. **More `SidebarWidth` values** — XS (8rem), 2XL (24rem).
22. **Footer responsive column counts** — currently hardcoded `grid-cols-2 md:grid-cols-4`; make it a prop.
23. **Form Layout `TwoColumnGrid`** variant (beyond current single-column Grid).
24. **AppShell variant: top-bar only** (no sidebar) — or document that `Nav + Container` is the pattern.
25. **AppShell variant: dual sidebar** (left nav + right detail panel).
26. **Split.FlippedOnMobile** — swap source order on mobile.
27. **CSS subgrid prototype** on Card (align header/title/footer across sibling cards) — research note exists.

### 🟢 Docs / recipes / demo

28. **Dedicated layout-primitives page on the website** (`website/src/content/`).
29. **README "Layout primitives" section** — currently describes the old 6-component layout package.
30. **Recipe: Split with sticky aside** (common docs-layout pattern).
31. **Recipe: AppShell + HTMX** — sidebar links use `hx-get` for SPA-like nav.
32. **Recipe: Container + Grid dashboard** (the canonical dashboard composition).
33. **Recipe: multi-column Footer with real content** (links, newsletter, social).
34. **Demo: interactive AppShell** with collapsible sidebar.
35. **Demo: Split with real article + metadata.**
36. **Demo: Container width comparison** (all 6 widths side-by-side).
37. **Demo: Stack gap comparison.**
38. _\*Example_* compile-tests_* for recipe code (plan F10.3) — prevents another `Cog6Tooth` incident.
39. **Update `skill/Skill.md`** with the 3 new test names (`TestRTLLogicalProperties`, `TestBodyPrimitivesDoNotEmitMain`, `TestAppShellCrossPackageComposition`).
40. **ADR for the MobileNav slot decision** (avoiding `layout → display` import).
41. **ADR for the `<main>` singleton ownership** by Base.

### 🔵 Manual verification / a11y

42. **Screen-reader test** of AppShell (NVDA or VoiceOver) — verify sidebar + content landmark navigation.
43. **Screen-reader test** of Split `<aside>` (should be announced as complementary landmark).
44. **RTL manual test** — render AppShell + Split under `dir="rtl"`, verify sidebar on right, aside mirrors.
45. **Keyboard-only test** of AppShell — Tab through sidebar → header → content order.
46. **Reduced-motion test** — AppShell has no animations, but verify `StickyHeader` doesn't cause jank.
47. **Performance: AppShell grid vs flex** — benchmark `lg:grid` vs `flex` for the sidebar+main shell.

### ⚪ Meta / tooling

48. **Add an AGENTS.md note** about LSP `unused` false positives on templ-called functions (lesson d.1).
49. **Consider a pre-commit grep gate** that catches references to non-existent icons in `.md` recipe files (lesson d.2).
50. **Consider committing `examples/demo/static/app.css` regeneration** as part of the pre-commit hook when `.templ` files change.

---

## g) Questions I CANNOT figure out myself

### Q1 — Release version: v0.19.0 now, or hold for v1.0?

The CHANGELOG `[Unreleased]` is warm and ready. `scripts/release.sh` will cut
a clean release. But: `forms.FormProps.Inline` is now soft-deprecated with
removal targeted at v1.0. Options:

- **(a) v0.19.0 now** (minor bump, new features). Consumers get the layout
  primitives immediately. `Inline` stays deprecated through v0.x.
- **(b) Hold and batch into v1.0** where `Inline` is also removed. Fewer
  releases, but delays the layout primitives for consumers who need them now.

I cannot decide this — it's a product/scheduling call. My recommendation: **(a)**,
because the layout primitives are high-demand and the deprecation path is clean.

### Q2 — Commit granularity for the 16-file polish surface?

Options:

- **(a) One commit**: "chore: post-implementation polish — CHANGELOG, RTL
  scanner, contract tests, demo dogfooding, recipe" (simple, matches the
  one-commit-per-release convention).
- **(b) Split**: `fix(split): correct lying comment on splitRatioCols` /
  `test(rtl): add logical-properties scanner + fix 2 demo violations` /
  `docs: warm CHANGELOG + ROADMAP + mobile-drawer recipe` /
  `test(integration): cross-package composition + <main> contract` (better
  bisect-ability, more noise in `git log`).

I cannot pick — it's a style preference. My recommendation: **(b)**, because
the lying-comment fix and the RTL scanner are logically independent and the
commit messages are more informative.

### Q3 — The prior session's status report (`docs/status/2026-07-21_01-24_...`): commit or trash?

It's untracked. It contains the self-assessment that drove this polish pass.
Options:

- **(a) Commit it** alongside this report — preserves the full session history.
- **(b) Trash it** — this report supersedes it; keeping both is noise.
- **(c) Commit but mark superseded** with a one-line append pointing here.

I cannot decide — it's a record-keeping preference. My recommendation: **(c)**,
because the prior report has the original 50-item next-action list that this
report builds on, and deleting it loses that context.

---

## Verdict

**Ship quality: high.** The layout primitives are well-tested, the CHANGELOG
is honest, the new CI guards (RTL, `<main>`, composition) raise the floor.
**Process quality: mixed** — 2 near-misses (d.1, d.2) that I caught by luck
rather than by process, and 1 repeated mistake (d.4, test-passes-≠-lint-passes)
that was already documented in the prior session's status report.

The work is done. The commit is waiting on Q2.
