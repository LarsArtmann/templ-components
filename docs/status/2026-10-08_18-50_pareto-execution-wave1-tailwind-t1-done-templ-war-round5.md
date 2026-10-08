# Status Report — Pareto Plan Execution Wave 1 (Tailwind T1 done + verified; templ-pin war round 5 live in HEAD)

**Written:** 2026-10-08 18:50 CEST
**Session scope:** Execution of `docs/planning/2026-10-08_15-18_tailwind-v4-audit-pareto-execution-plan.md` ("NOW GET SHIT DONE" dispatch). This report covers THIS execution session only — planning/audit history lives in the input artifacts.
**Branch state at writing:** `master` ahead of origin by 4 commits, ALL FOUR daemon-authored, ALL FOUR contaminated (details in §d). Working tree holds 4 uncommitted file restores that FIX the worst contamination.

---

## TL;DR

- **T1 (theming tokens — the plan's 1%→51% tier) is code-complete AND fully verified**, including the full visual regression suite passing GREEN with zero golden drift. The `@theme` re-skinning promise now holds for the entire `templates/custom.css` file, guarded forever by a new drift-guard test.
- **Verify-at-source paid off before a single site was migrated:** the probe against the pinned tailwindcss v4.3.3 binary corrected the audit's assumed `--alpha()` comma syntax (real: slash syntax), and confirmed all four migration targets (`inset-s-*`, `scrollbar-none`, `field-sizing-content`, `--alpha()`) exist.
- **The audit's finding 1 undercounted:** the real raw-literal population in custom.css was ~45 sites, not 12. All of them are gone.
- **The templ-pin war is in round 5:** while this session worked, the daemon re-applied the v0.3.1070 bump to ALL go.mod/go.sum files + flake.lock in 4 unpushed commits (02f1863f, ee2f10bd, a7f06496, f35bcd23-adjacent). The working tree has already been re-resolved back to v0.3.1020 (uncommitted). HEAD must not be pushed as-is.
- Everything else in the plan (T2 inset migration, lane pinning, T3 quick wins, T4 @utility/docs/tooling, release gates) is recon-complete but not yet edited — see §c.

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| a1 | **Verify-at-source probe of all four Tailwind claims against the pinned binary (v4.3.3)** | `/tmp/twprobe` compile probes: `.inset-s-0{inset-inline-start:0px}`, `.inset-s-full{…100%}`, `.inset-e-4{…calc(var(--spacing)*4)}`, `.scrollbar-none{scrollbar-width:none}`, `.field-sizing-content{field-sizing:content}` all EMIT. `--alpha()` works but ONLY with **slash syntax** `--alpha(var(--color-blue-500) / 30%)` — the audit's comma form is a hard build error ("requires a color and an alpha value"). Emission shape: static sRGB fallback (inlined value) + `@supports` branch with **var-preserving** `color-mix(in oklab, var(--color-x) N%, transparent)` — exactly the re-skinning behavior finding 1 needs. |
| a2 | **F-01..F-03 + F-49: full custom.css tokenization** — 19 edits, **~45 raw palette literals eliminated file-wide** (audit said 12; the `tc-select` stylable block, dialog backdrops, sidebar `#000`, reduced-transparency, heatmap tokens, kanban DnD + pending + failed states were all uncounted). All now `var(--color-*)`; alpha'd values use `--alpha(var(--color-x) / N%)`. Policy documented in the file header (per guardrail 7). `grep` confirms **0 remaining** raw hex/rgb/hsl in custom.css. |
| a3 | **F-49 partial (heatmap):** `--ds-brand` hex → `var(--color-violet-600/-500)`. The `--ds-brand-rgb` comma triplets deliberately KEPT: they are the ADR-0028 legacy alpha contract consumed by `heatmapCellBackground`'s `rgba(var(--x-rgb), α)` inline styles (documented public `ColorVar` override contract in `docs/recipes/heatmap.md`). Tokenizing them away requires a breaking component change → routed to v2 window (see f-item f-06). Guard exempts them correctly (bare number lists are not color literals). |
| a4 | **F-04: `utils.TestCustomCSSThemeTokens` written and PROVEN** — red-green sequence executed (planted `#123456` → FAIL with remediation message naming the line; reverted → PASS). Scans hex + `rgba?(` + `hsla?(` with `/* */` comment stripping. Lives in `utils/custom_css_test.go`. |
| a5 | **F-05: full guard family green** over the change: `TestCustomCSSThemeTokens`, `TestCustomCSSUtilities`, `TestDarkModeCompliance`, `TestDarkModeSemanticColors`, `TestMotionReduceCompliance`, `TestCoarsePointerCompliance`, `TestDocsCountDrift`, icons-module custom-CSS guard, `TestTemplVersionPin` (working tree). |
| a6 | **Demo CSS recompiled + CSS guards green:** `nix run .#css` → committed `examples/demo/static/app.css` refreshed (1 minified line); `TestCSSFreshness`, `TestCompiledCSSInventory`, `TestTailwindGoSourceScanning` all green. Compiled output spot-checked: `accent-color:var(--color-blue-600)` present; `.tc-kanban-over` carries BOTH the static fallback (`#e0e7ff99`) AND the var-preserving `color-mix` @supports branch; 0 legacy literals remain in the artifact. |
| a7 | **F-06..F-08 (T1 verification): FULL visual regression suite GREEN — 179.8s, `ok visualtest`, ZERO golden re-baseline needed.** The token swap is pixel-equivalent within the 0.1% tolerance. F-08's feared visual re-baseline turned out unnecessary for T1. |
| a8 | **Complete T2 recon (read-only):** all 11 `start-*`/`end-*` sites enumerated with exact replacement classes (`inset-s-0`, `inset-s-0.5`, `inset-e-0`, `inset-e-4`, `inset-s-full`, `inset-e-full`); workspace-wide grep proves NO additional sites hide in `.go` lookup maps (the audit only scanned `.templ`). |
| a9 | **Complete T3/T4 recon:** textarea AutoGrow emission site (`forms/textarea.templ:87` — `.tc-auto-grow` = `field-sizing + min-height:2.5rem + max-height:20rem`, exact stock-utility parity `field-sizing-content min-h-10 max-h-80`); `TestTextareaAutoGrowAddsClass` assertion located (must be updated); website `hero.templ:107` gradient + Dockerfile/website.yml unpinned installs + `site.css` scans GENERATED `*_templ.go` files (so post-regen classes flow automatically); `TestCustomCSSUtilities` parser needs a small extension for `@utility` syntax (F-37); `tc-snap-*` discovery — `tc-snap-x` ≡ stock `snap-x snap-mandatory` and `tc-snap-center` ≡ stock `snap-center snap-always`, so those two custom classes can be DELETED in favor of stock utilities (better than @utility-migrating them). |

## b) PARTIALLY DONE

| # | Item | Done | Missing |
|---|------|------|---------|
| b1 | **T1 lane (finding 1)** | Code + guards + visual verification (a2–a7) | Nothing functional — but the CHANGELOG `[Unreleased]` entry for it is NOT yet written (gated on F-50 ship decision), and the deep-dive HTML report still says "12 literals" (needs correction in F-28's appendix). |
| b2 | **Pin restoration (templ war round 5)** | Working tree ALREADY carries the fix (visualtest+website go.mod/go.sum back at v0.3.1020; `TestTemplVersionPin` green on tree) | Uncommitted; HEAD still carries the 1070 bump + drifted flake.lock in 4 unpushed daemon commits. Restore-commit + flake.lock revert + push still to do (f-01/f-02). |
| b3 | **T2 inset migration** | All 11 sites identified with verified replacement classes (a8); `TestRTLLogicalProperties` read and tightening design ready | Zero sites migrated; guard not yet tightened; AGENTS/skill RTL text not yet updated. |
| b4 | **T3 field-sizing** | Exact class-string parity established (`field-sizing-content min-h-10 max-h-80` ≡ `.tc-auto-grow`); test + docs reference map built | Component not edited; `.tc-auto-grow` not deleted; goldens not re-baselined. |
| b5 | **T4 @utility** | Candidate set classified during recon: `@utility` for `tc-squircle`, `tc-content-auto(-compact)`, `tc-log`, `tc-log-line-still`, `tc-fluid-*`; stock-utility DELETION for `tc-snap-*`; keep-plain for component-shaped classes (tc-select, tc-modal/drawer/overlay, tc-anim-*, tc-kanban-*) | No CSS edited; guard parser extension (a9) not implemented. |
| b6 | **Daemon-race containment this session** | Caught the race mid-session (f35bcd23), verified committed content integrity (77 token references present in daemon commit), caught round-5 re-bump at report time | Final containment (commit tree fix, push) blocked on user's "wait" instruction — correctly so. |

## c) NOT STARTED (plan tasks still open — all recon-ready)

| # | Plan ref | Task |
|---|----------|------|
| c1 | F-09..F-12 | Migrate the 11 inset sites; `nix run .#build`; regen from repo root |
| c2 | F-13 | Tighten `utils.TestRTLLogicalProperties` to ban `start-`/`end-` class forms (regex designed: require whitespace/quote before + Tailwind value after, so `text-start`/`items-start`/`data-tc-*` stay legal) |
| c3 | F-14/F-15 | RTL convention text in AGENTS.md + `skill/SKILL.md` (canonical examples swap to `inset-s-/e-`) |
| c4 | F-16/F-17 | Goldens for the 7 touched components + RTL visual spot check |
| c5 | F-18/F-19/F-20 | Pin `tailwindcss@4.3.3 @tailwindcss/cli@4.3.3` in Dockerfile CSS stage + website.yml npm install + AGENTS bump-checklist line |
| c6 | F-21/F-22 | Carousel `tc-no-scrollbar` → `scrollbar-none`; delete the custom class; sync guards |
| c7 | F-23/F-24 | `bg-gradient-to-br` → `bg-linear-to-br` (notfound404.templ:42), `bg-gradient-to-r` → `bg-linear-to-r` (website hero.templ:107); site goldens if affected |
| c8 | F-25..F-27 | Textarea field-sizing migration (see b4) |
| c9 | F-28/F-29 | Deep-dive report capability-ledger appendix + correction of the 12-literal and 21-vs-11 counts; chromedp screenshot pass of the HTML report (never yet rendered in a browser) |
| c10 | F-30..F-32 | HARVEST: TODO_LIST.md short-term entries, ROADMAP.md idea entries (incl. the new heatmap color-mix modernization item), annotate status + plan as harvested |
| c11 | F-33..F-37 | @utility audit table + migration batches + guard sync (design in b5) |
| c12 | F-38..F-42 | size-* sweep — **DEFERRED by plan guardrail 3** (golden re-baseline MUST ride the templ v0.3.1070 migration window; sweeping now = red suite by construction). Route to TODO_LIST as one gated batch. |
| c13 | F-43..F-46 | Docs: user-valid note, adoption-guide capability section, agent-context-history lessons (recount-rule, heredoc/sed relapse, probe-first wins), `TestDocsCountDrift` re-run |
| c14 | F-47/F-48 | `scripts/check-html-report-classes.sh` + screenshot-pass checklist line |
| c15 | F-50/F-51 | Q2 ship decision → CHANGELOG warm → version triple-lock → release ritual |
| c16 | F-52 | Q1 artifact keep/fold annotation |
| c17 | F-53 | `scripts/ci-repro.sh --lint --website` + push (currently BLOCKED by §d items) |

## d) TOTALLY FUCKED UP

| # | Item | Detail | Status |
|---|------|--------|--------|
| d1 | **templ-pin war ROUND 5 — HEAD is poisoned, 4 commits ahead of origin** | While this session worked, daemon commits re-applied `a-h/templ v0.3.1070` to go.mod/go.sum of ALL modules (02f1863f: root+utils+icons+htmx+errorpage+datastar+charts+website; ee2f10bd: visualtest) AND moved `flake.lock`'s nixpkgs input again (12 lines — the lock's templ binary is the generator-drift hazard). 5th documented recurrence; `TestTemplVersionPin` only guards the working tree, and the daemon does not run tests. | **Working tree ALREADY fixed** (go tooling re-resolved visualtest+website down to 1020 under workspace mode; all-module pin check green). Needs: verify root+utils+… go.mods in HEAD vs tree, revert flake.lock node (exact-string edit, jq lacks --slurpfile), single restore commit, push. DO NOT push the 4 daemon commits as-is. |
| d2 | **BuildFlow `modernize` override of a documented nolint (in daemon commit 02f1863f)** | `visualtest/render.go:86`: `//nolint:modernize // wrapper exists for API stability + tri-state clarity` was replaced by `//go:fix inline` + body `return &b` → `return new(b)`. This REVERSES the wrapper's documented intent: `//go:fix inline` actively REWRITES call sites away from `visualtest.Bool(...)` at next gopls fix pass. Not authored by this session; per the never-revert rule it is flagged, not touched. | Needs an operator decision: restore the nolint, or accept the inline directive (f-03). |
| d3 | **Self-inflicted: sed line-dance over-deleted custom.css's closing brace** | Proving the guard red (append 2 lines, remove 2) — I removed 3 `sed -i '$d'` times, silently chopping the file's final `}`. Caught only because I byte-compared the tail against HEAD (`git stash` + `od -c`). The guard test could NOT catch a structural brace deletion (it only scans color literals) — the file was one commit away from shipping broken CSS. | Fixed within the session; lesson recorded for f-40 (never add/remove lines via sed loops on files — use the edit tool or `git checkout`-free exact restores). |
| d4 | **Self-inflicted: `%` verb bug in the new guard's `t.Errorf` format string** | `--alpha(var(--color-blue-500) / 30%)` embedded raw in a format string → `%) unknown verb` BUILD FAILURE on first run (go vet catches this in CI, but it wasted a cycle locally). Fixed by passing the example as a `%s` argument. | Fixed. Same lesson family as the repo's "never patch code via heredocs": format strings with `%` need the same paranoia. |
| d5 | **Uncommitted work left on the tree at session pause** | 4 modified files (visualtest+website go.mod/go.sum — the CORRECT 1020 state) plus the earlier trailing-whitespace removal are uncommitted BY INSTRUCTION (user said report + wait). Risk: the next daemon cycle snapshots the tree into another heuristic commit, mixing my pin restore with whatever lands next. | Acceptable under the wait order; first action next session is the containment commit (f-01). |
| d6 | **Untracked CSS litter still breeding** | `examples/demo/demo.out.css`, `website/site.out.css`, `website/dist/**` present untracked (daemon tailwind-build cycles + this session's visual-suite website rebuild). `TestCompiledCSSInventory` reports them informationally only. | Known class, covered by guard design; no action needed beyond awareness. |

## e) WHAT WE SHOULD IMPROVE

1. **Probe-first is now proven policy for every Tailwind claim** — this session's single most valuable 5 minutes was the compile probe: it invalidated the audit's `--alpha()` comma syntax BEFORE it could ship a broken custom.css, and confirmed all three utility renames. The audit report (and the skill) should codify "no class/function migration without a binary probe."
2. **Audit counting needs a completeness gate** — finding 1's "12 literals" was ~45 in reality; last session's "11 sites" had been 21 a day earlier. Both came from prose-grepping instead of an exhaustive mechanical sweep. A tiny script (`rg` + strict class-boundary regex over `templates/*.css`) would make literal counts mechanical; same class as F-47's report-class checker.
3. **`TestTemplVersionPin` guards the tree, not HEAD** — round 5 proved the daemon can poison HEAD between the guard run and any push. The pre-push ritual should add a `git show HEAD:go.mod | grep templ` check (or the ritual's ci-repro already builds from HEAD — verify it would catch this; if yes, document that).
4. **flake.lock must join the pin-restore checklist** — the summary of the last restore explicitly reverted the flake.lock node; round 5's re-bump includes it again. Make a one-command `scripts/restore-templ-pin.sh`? (or at minimum a checklist line in AGENTS.md next to the existing paragraph — this is f-02.)
5. **Stock-utility replacement beats @utility where parity exists** — the `tc-snap-*` discovery (exact stock equivalents) deletes code instead of re-shelving it. The @utility audit (c11) should FIRST check "is this now a stock utility?" per class.
6. **Website sources drift from library conventions** — `website/internal/pages/hero.templ` carries physical `ml-3`/`text-left` (invisible to `TestRTLLogicalProperties`, which scans only library dirs). Two-character fixes riding F-24.
7. **Demo Dockerfile builds with `golang:1.26` while the repo floor is go 1.27** — pre-existing, out of plan scope, flagged for the next infra pass (can't be verified locally without a docker build).
8. **The plan's verification batching worked** — one full visual pass for T1 in isolation (green, zero drift) beats three partial passes; keep the consolidated-pass pattern for T2/T3 (one more full pass before push).

## f) Next things (up to 50)

**P0 — containment & gates (before any push):**
| # | Task | Est |
|---|------|-----|
| f-01 | Commit the working-tree pin restore (visualtest+website go.mod/go.sum at 1020) | 5m |
| f-02 | Revert flake.lock nixpkgs node to the pre-drift revision (exact-string edit; verify `templ` binary version in lock) | 10m |
| f-03 | Decide render.go: restore `//nolint:modernize` or accept `//go:fix inline` (operator input, §g/Q-adjacent) | 5m |
| f-04 | Squash/repair the 4 poisoned daemon commits (history rewrite on unpushed master: soft-reset to 65b74cbc + recommit coherent) or a corrective commit on top — pick per repo convention | 15m |
| f-05 | `git fetch` + re-verify remote tip immediately before push (daemon races pushes; ritual M03) | 2m |

**P1 — 4% tier (inset + pins):**
| # | Task | Est |
|---|------|-----|
| f-06 | ROADMAP: record Heatmap `color-mix` modernization (kill `--ds-brand-rgb` triplets + ADR-0028 amend) for the v2 window | 5m |
| f-07 | F-09: migrate notfound404.templ:119 + toast.templ:92 → `inset-s-0`/`inset-e-4` | 5m |
| f-08 | F-10: migrate toggle.templ:102 (`inset-s-0.5`), input_group.templ:60+68 | 5m |
| f-09 | F-11: migrate hover_card.templ:19-20 (`inset-e-full`/`inset-s-full`), avatar.templ:108+136, carousel.templ:86+96 | 10m |
| f-10 | F-12: `nix run .#build` (templ generate FROM REPO ROOT + build; check zero unexpected `_templ.go` diffs) | 10m |
| f-11 | F-13: tighten `TestRTLLogicalProperties` (ban-list regex; red→green proof) | 12m |
| f-12 | F-14: AGENTS.md RTL convention text + canonical example | 5m |
| f-13 | F-15: skill/SKILL.md RTL text (in-repo `skill/SKILL.md`) | 5m |
| f-14 | F-16: goldens `-update` for the 7 touched components (diff-review every file) | 12m |
| f-15 | F-17: RTL visual spot check (carousel/toggle/avatar RTL captures) | 12m |
| f-16 | F-18: Dockerfile `pnpm add tailwindcss@4.3.3 @tailwindcss/cli@4.3.3` | 5m |
| f-17 | F-19: website.yml `npm install --no-save tailwindcss@4.3.3 @tailwindcss/cli@4.3.3` | 5m |
| f-18 | F-20: AGENTS.md tailwind-version-bump checklist line next to the templ-pin paragraph | 5m |

**P2 — 20% tier (quick wins):**
| # | Task | Est |
|---|------|-----|
| f-19 | F-21: carousel.templ:64 → `scrollbar-none`; delete `.tc-no-scrollbar` (both rules) from custom.css | 5m |
| f-20 | F-22: guard sync + carousel goldens | 5m |
| f-21 | F-23: notfound404.templ:42 → `bg-linear-to-br` (keep dark: variants) | 5m |
| f-22 | F-24: website hero.templ:107 → `bg-linear-to-r`; site tests/goldens if affected | 10m |
| f-23 | F-24b (opportunistic): hero.templ `ml-3` → `ms-3`, `text-left` → `text-start` (same file, convention alignment) | 2m |
| f-24 | F-25: textarea.templ:87 → `field-sizing-content min-h-10 max-h-80` | 5m |
| f-25 | F-26: delete `.tc-auto-grow` from custom.css; sync `TestCustomCSSUtilities` expectations | 5m |
| f-26 | F-26b: update `TestTextareaAutoGrowAddsClass` (+ False-omits) to the new class set | 5m |
| f-27 | F-27: textarea goldens + demo spot check (grow behavior both themes) | 12m |
| f-28 | F-28: deep-dive report appendix — capability ledger + CORRECTED counts (12→~45 literals; include the --alpha slash-syntax correction) | 12m |
| f-29 | F-29: chromedp screenshot pass of the deep-dive HTML (hero/scorecard/tables/print) | 12m |
| f-30 | F-30: TODO_LIST.md short-term entries (tokens ✅-note, inset, pins, scrollbar, field-sizing, bg-linear) | 12m |
| f-31 | F-31: ROADMAP.md idea entries (@utility tree-shaking, size-* gated batch, masks/text-shadow/@container-size) | 10m |
| f-32 | F-32: annotate status report + plan as harvested | 10m |

**P3 — remaining 80% tier:**
| # | Task | Est |
|---|------|-----|
| f-33 | F-33/F-34: @utility audit — dead/alive table + current-vs-projected emitted bytes | 12m |
| f-34 | F-35: batch 1 — `@utility` for `tc-squircle`, `tc-content-auto(-compact)`; DELETE `tc-snap-*` → stock `snap-x snap-mandatory` / `snap-center snap-always` (probe first) | 12m |
| f-35 | F-36: batch 2 — `@utility` for `tc-fluid-*`, `tc-log`, `tc-log-line-still` | 12m |
| f-36 | F-37: extend `loadCustomCSSClasses` to parse `@utility tc-*` forms; full guard suite | 12m |
| f-37 | F-38..F-42: harvest size-* sweep as ONE gated TODO item (templ v0.3.1070 window) — do NOT sweep now (guardrail 3) | 5m |
| f-38 | F-43: user-valid docs note | 10m |
| f-39 | F-44: adoption-guide "Tailwind capability coverage" section | 12m |
| f-40 | F-45: agent-context-history lessons (probe-first win; recount rule; sed line-dance + `%`-format relapse) | 12m |
| f-41 | F-46: `TestDocsCountDrift` after all docs edits | 5m |
| f-42 | F-47: `scripts/check-html-report-classes.sh` + self-apply | 12m |
| f-43 | F-48: screenshot-pass checklist line (AGENTS docs-health section or script header) | 5m |

**P4 — release & ritual:**
| # | Task | Est |
|---|------|-----|
| f-44 | F-50: record Q2 decision (ship-now vs v2) | 5m |
| f-45 | F-51: warm CHANGELOG `[Unreleased]` (theming fix, inset migration, scrollbar, field-sizing, bg-linear, pins, guard) | 10m |
| f-46 | F-51b: pre-verify lint + touched packages BEFORE release.sh; then `nix shell nixpkgs#govulncheck -c nix develop -c scripts/release.sh` | ritual |
| f-47 | F-52: record Q1 decision; annotate artifacts | 5m |
| f-48 | Full `nix run .#verify` + consolidated full visual pass (T2/T3 changes) | 30m+ |
| f-49 | F-53: `scripts/ci-repro.sh --lint --website` at the exact commit → push IMMEDIATELY | ritual |
| f-50 | Post-push: verify referenced guards in CI; watch first daemon cycle for resurrection (`.out.css`, templ bump) | 10m |

## g) Questions I cannot figure out myself (plan Q1–Q3, still open after two asks)

1. **Q2 (blocks the release gate f-44..f-46):** Does finding 1 (theming tokens) ship as a **patch release now**, or wait for the **v2 window**? The change is consumer-visible bug-fix quality (rendered colors are pixel-equivalent; the fix restores a documented promise), which argues ship-now. If ship-now: I warm the CHANGELOG and run the release ritual; if v2: entries go to ROADMAP and only the code fix rides the next natural release.
2. **Q3 (blocks f-03/f-06 finality):** `--ds-brand` has no documented external consumer in this repo. (a) Is violet (per the 2026-09-08 user steer) still the intended brand default, and (b) do you know of any external consumer relying on the `--x-rgb` companion contract? If no → the v2-window modernization (color-mix, single token) is safe to roadmap; if yes → I document the contract instead.
3. **Q-render (blocks f-03):** The daemon/BuildFlow replaced the documented `//nolint:modernize` on `visualtest.Bool` with `//go:fix inline` (`return &b` → `return new(b)`), which will rewrite consumer call sites on the next gopls fix pass — the opposite of the wrapper's documented "API stability" intent. Restore the nolint (my recommendation: the wrapper exists precisely to be a stable API), or accept the inline directive?

---

## Verification log (commands → results, this session)

| Command | Result |
|---|---|
| `tailwindcss -i probe.css -o out.css` (pinned v4.3.3, 3 probe rounds) | inset-s-0/inset-s-full/inset-e-4/scrollbar-none/field-sizing-content all emitted; `--alpha` comma form = build ERROR; slash form emits fallback + var-preserving @supports branch |
| 19-edit multiedit on `templates/custom.css` | applied; `grep -E '#…|rgba?\(|hsla?\('` → CLEAN |
| guard red-green (planted `#123456`) | FAIL names custom.css line → PASS after removal |
| `go test ./utils -run 'TestCustomCSS…' -count=1` | ok |
| `go test ./utils -run 'TestDarkMode\|TestMotionReduce\|TestCoarsePointer\|TestDocsCountDrift' -count=1` | ok |
| `(cd icons && GOWORK=off go test ./...)` | ok |
| `nix run .#css` + `TestCSSFreshness` + `TestCompiledCSSInventory` + `TestTailwindGoSourceScanning` | CSS recompiled (1 minified line); all green |
| compiled-artifact spot check | `accent-color:var(--color-blue-600)`; `.tc-kanban-over` has static fallback + var color-mix @supports; 0 legacy literals |
| `nix run .#visual` (full suite) | **ok 179.8s — zero golden drift, zero re-baseline** |
| `go test ./utils -run TestTemplVersionPin -count=1` (working tree) | ok |

_Daemon-context note: 4 unpushed daemon commits (02f1863f, ee2f10bd, a7f06496 + f35bcd23) mix this session's T1 content with the round-5 templ re-bump, flake.lock drift, a ci.yaml action-annotation bump (benign), the render.go nolint override (§d2), and a parallel session's docs file. Containment plan in f-01..f-05._
