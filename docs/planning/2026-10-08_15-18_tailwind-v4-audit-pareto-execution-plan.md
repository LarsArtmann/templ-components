# Pareto Execution Plan — Tailwind v4 Audit Follow-Up

**Created:** 2026-10-08 15:18 CEST
**Scope:** execution plan for ALL todos from the 2026-10-08 Tailwind CSS v4 deep-dive session (8 findings, 35 status-report items). This is a PLANNING artifact — no source code changed for this plan.
**Format note:** pareto-planning skill's canonical output is HTML; the operator prompt explicitly demands `.md` with an inline mermaid.js/d2 graph — the explicit instruction wins (flagged per skill spec). Fine-granularity capped at 12min per task (operator) instead of the skill's 15min.
**Input artifacts:** `docs/research/2026-10-08_tailwindcss-deep-dive.html` (findings 1-8) · `docs/status/2026-10-08_15-08_tailwind-deep-dive-session-status.md` (items f1-f35, questions Q1-Q3).

---

## Guardrails (VERSCHLIMMBESSERN protection — read before executing)

1. **Never push on "should be green"** — `scripts/ci-repro.sh --lint --website` at the exact commit, then push immediately (repo RITUAL M03). Exception documented per-push if ever narrower.
2. **Golden updates:** `-update` goes AFTER the package list (`go test ./pkg/... -update`, never before — go1.26.7 flag-parsing trap). Eyeball every `.golden` diff.
3. **Batch re-baselines:** size-* sweep goldens MUST ride the planned templ v0.3.1070 migration window, not their own re-baseline.
4. **Daemon:** auto-commits every ~60s without running tests. After any edit: re-check your edits survived (`git log -1 -- <file>` moving backward = the 6th incident class). After pushing, re-verify remote tip.
5. **Version triple-lock:** any release work bumps `utils/version.go` + CHANGELOG heading + FEATURES `**Version:**` together (guards enforce).
6. **Do not touch:** templ pin (v0.3.1020), dependency budget, ContainerAware component list (ADR-0018 closed), Web Components (ADR-0033, binding), disabled linters.
7. **`--alpha()` caveat:** it is a BUILD-TIME function usable in Tailwind-processed CSS — custom.css IS processed (imported by app.css/demo.css), so it resolves; if a consumer imports custom.css without Tailwind, prefer `color-mix(in oklab, var(--color-x) N%, transparent)` which is plain CSS. Decide per-rule during T1.3 and document the choice in custom.css.
8. **RTL guard update ordering:** migrate sites BEFORE tightening `TestRTLLogicalProperties`, or the suite goes red mid-migration.

## Pareto breakdown (value = consumer-visible correctness + future-proofing + CI safety)

| Tier                      | Share of value | Tasks                                                                                       | Why                                                                                                                                                          |
| ------------------------- | -------------- | ------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **1%**                    | **51%**        | T1 theming tokens (+T14 --ds-brand)                                                         | The ONLY finding that breaks a documented consumer promise (`@theme` re-skinning). 12 literals, 6 a11y/state surfaces, one file. Everything else is hygiene. |
| **4%**                    | **64%**        | + T2 inset-s/e migration, T3 lane pinning                                                   | Kills the silent RTL-break exposure (deprecated aliases inside the RTL canon) and the version-skew CI risk (2 unpinned lanes). Both mechanical.              |
| **20%**                   | **80%**        | + T4 scrollbar-none, T5 bg-linear, T6 field-sizing, T7 report integrity, T8 harvest         | Quick wins (<1h each), each deletes hand-rolled code Tailwind now ships, plus closes this session's own accuracy debts.                                      |
| **remaining 80% of work** | **100%**       | T9 @utility, T10 size-*, T11 docs/lessons, T12 tooling, T13 ship window, T15 open questions | Real but deferrable: bundle-size win, idiom sweeps (gate on golden re-baseline windows), knowledge capture, and the release decision.                        |

---

## LEVEL 1 — Comprehensive plan (tasks 30-100 min, ALL todos, sorted by impact/effort/value)

| #     | Task                                                                                                                                                                              | Tier        | Finding | Est   | Depends on       | Covers status items |
| ----- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- | ------- | ----- | ---------------- | ------------------- |
| L1-01 | **Theming tokens:** replace 12 hardcoded color literals in custom.css with `var(--color-*)` / `--alpha()`; add `TestCustomCSSThemeTokens` guard                                   | 1%          | F1      | 60min | —                | f4-f6, f8           |
| L1-02 | **Themed-surface verification:** goldens + visual re-baseline for checkbox/toggle/radio/progress/kanban/validation                                                                | 1%          | F1      | 60min | L1-01            | f9                  |
| L1-03 | **Logical inset migration:** 11 `start-*`/`end-*` sites → `inset-s-*`/`inset-e-*`, regenerate, build                                                                              | 4%          | F2      | 45min | —                | f11                 |
| L1-04 | **RTL canon sync:** ban deprecated forms in `TestRTLLogicalProperties`; update AGENTS.md + skill/SKILL.md convention text                                                         | 4%          | F2      | 45min | L1-03            | f12, f13            |
| L1-05 | **Inset verification:** goldens for 7 touched components + RTL visual spot check                                                                                                  | 4%          | F2      | 45min | L1-04            | f14                 |
| L1-06 | **Lane pinning:** tailwindcss@4.3.3 in Dockerfile + website.yml; version-bump checklist in AGENTS.md                                                                              | 4%          | F3      | 30min | —                | f15-f17             |
| L1-07 | **scrollbar-none:** carousel swap, delete `.tc-no-scrollbar`, guard sync                                                                                                          | 20%         | F5      | 30min | —                | f22                 |
| L1-08 | **bg-linear rename:** notfound404 hero + website hero                                                                                                                             | 20%         | F8      | 30min | —                | f26                 |
| L1-09 | **field-sizing:** Textarea emits `field-sizing-content min-h-10 max-h-80`; retire `.tc-auto-grow`                                                                                 | 20%         | F6      | 45min | —                | f23                 |
| L1-10 | **Report integrity:** capability ledger appendix + screenshot pass on the deep-dive HTML (daemon-race check f1 already RESOLVED: HEAD carries the corrected 11-sites version)     | 20%         | —       | 30min | —                | f1✅, f2, f3        |
| L1-11 | **Harvest:** route all remaining items into TODO_LIST.md (short-term) / ROADMAP.md (ideas); annotate status + this plan as harvested                                              | 20%         | —       | 30min | plan approval    | f34, f35            |
| L1-12 | **@utility audit:** which `.tc-*` classes do components actually use vs docs-only; measure current vs projected emitted bytes                                                     | 80%-work    | F4      | 30min | L1-07            | f18, f19            |
| L1-13 | **@utility migration:** migrate utility-shaped classes, sync `TestCustomCSSUtilities`, full guard suite                                                                           | 80%-work    | F4      | 60min | L1-12            | f20, f21            |
| L1-14 | __size-_ sweep:_* 168 w/h token pairs → `size-N`, regenerate, verify artifact byte-stability                                                                                      | 80%-work    | F7      | 60min | —                | f25, f27            |
| L1-15 | __size-_ goldens:_* re-baseline — GATED on templ v0.3.1070 migration window                                                                                                       | 80%-work    | F7      | 30min | external gate    | f25                 |
| L1-16 | **Docs & lessons:** adoption-guide capability note, `user-valid` docs note, agent-context-history lessons, docs-count same-commit rule check                                      | 80%-work    | —       | 45min | L1-01..L1-09     | f24, f28, f29, f30  |
| L1-17 | **Verification tooling:** HTML report class-check script; screenshot-pass checklist for HTML deliverables                                                                         | 80%-work    | —       | 30min | —                | f31, f32            |
| L1-18 | **`--ds-brand` decision:** tokenize or document consumer (needs Q3 answer)                                                                                                        | 1%-adjacent | F1      | 15min | Q3               | f7                  |
| L1-19 | **Ship window:** resolve Q2 → if ship-now: CHANGELOG warm, version triple-lock, release ritual (pre-verify lint + touched packages FIRST, never let release.sh be the first gate) | gate        | F1-F6   | 60min | L1-01..L1-06, Q2 | f10                 |
| L1-20 | **Open questions:** resolve Q1 (artifact keep/fold) with operator                                                                                                                 | gate        | —       | 15min | —                | —                   |

## LEVEL 2 — Fine breakdown (max 12 min each, ALL todos, sorted by impact/effort/value)

| #    | Task (≤12min)                                                                                                                                                                                      | L1       | Est     | Verify                                                   |
| ---- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- | -------------------------------------------------------- |
| F-01 | custom.css:324,330 accent-color `rgb(37 99 235)`/`rgb(96 165 250)` → `var(--color-blue-600)`/`var(--color-blue-400)`                                                                               | L1-01    | 10min   | rg: 0 hits `rgb(37 99`                                   |
| F-02 | custom.css:341-348 user-valid/invalid borders → `var(--color-red-600)`/`var(--color-green-600)`                                                                                                    | L1-01    | 10min   | rg + git diff                                            |
| F-03 | custom.css:1045-1063 kanban literals → `var(--color-*)`; per guardrail 7 pick `--alpha()` vs `color-mix()` and comment the choice                                                                  | L1-01    | 12min   | rg: 0 raw hex in block                                   |
| F-04 | Write `utils.TestCustomCSSThemeTokens` (fail on raw hex/rgb in custom.css, comments exempt) + `-count=1` run                                                                                       | L1-01    | 12min   | `go test ./utils -run TestCustomCSSThemeTokens -count=1` |
| F-05 | Run dark-mode + custom-CSS + docs-count guard family over the change                                                                                                                               | L1-01    | 10min   | `go test ./utils -run 'TestDarkMode                      |
| F-06 | Goldens `-update` for themed surfaces (checkbox, toggle, radio, slider, progress, kanban) — diff-review each                                                                                       | L1-02    | 12min   | `git diff --stat '**/*.golden'`                          |
| F-07 | Remaining themed goldens (validation states, input_group)                                                                                                                                          | L1-02    | 12min   | same                                                     |
| F-08 | Visual re-baseline `nix run .#visual -- -update` for themed captures; eyeball PNG diffs                                                                                                            | L1-02    | 12min   | suite green                                              |
| F-09 | Migrate `errorpage/notfound404.templ:119` + `feedback/toast.templ:92` → `inset-s-0`/`inset-e-4`                                                                                                    | L1-03    | 10min   | rg `start-0                                              |
| F-10 | Migrate `forms/toggle.templ:102`, `forms/input_group.templ:60,68`                                                                                                                                  | L1-03    | 10min   | rg                                                       |
| F-11 | Migrate `display/hover_card.templ:19-20`, `display/avatar.templ:108,136`, `display/carousel.templ:86,96`                                                                                           | L1-03    | 12min   | rg: workspace `[" ](start\|end)-(0\|0.5\|4\|full)` = 0   |
| F-12 | `nix run .#build` (templ generate from repo root + go build)                                                                                                                                       | L1-03    | 10min   | build green, `*_templ.go` diffs only where expected      |
| F-13 | Tighten `utils.TestRTLLogicalProperties`: ban `start-`/`end-` class forms (allow `inset-s-/e-`)                                                                                                    | L1-04    | 12min   | suite red→green sequence proves guard bites              |
| F-14 | Update RTL convention in AGENTS.md (+ canonical example swap)                                                                                                                                      | L1-04    | 10min   | rg `start-` in AGENTS RTL section                        |
| F-15 | Update RTL convention in `skill/SKILL.md` (+ docs if referenced)                                                                                                                                   | L1-04    | 10min   | rg                                                       |
| F-16 | Goldens `-update` for the 7 touched components; review diffs                                                                                                                                       | L1-05    | 12min   | diff review                                              |
| F-17 | RTL visual spot check (`nix run .#visual` RTL captures of carousel/toggle/avatar)                                                                                                                  | L1-05    | 12min   | suite green                                              |
| F-18 | Pin `tailwindcss@4.3.3 @tailwindcss/cli@4.3.3` in `examples/demo/Dockerfile:18`                                                                                                                    | L1-06    | 10min   | rg pinned spec in file                                   |
| F-19 | Pin same versions in `.github/workflows/website.yml:68` npm install                                                                                                                                | L1-06    | 10min   | rg                                                       |
| F-20 | Add tailwind-version-bump checklist line next to the templ-pin paragraph in AGENTS.md                                                                                                              | L1-06    | 10min   | rg                                                       |
| F-21 | `display/carousel.templ:64` → `scrollbar-none`; delete `.tc-no-scrollbar` rules from custom.css                                                                                                    | L1-07    | 10min   | rg `tc-no-scrollbar` = 0                                 |
| F-22 | Sync `TestCustomCSSUtilities` expectations; run guards + carousel goldens                                                                                                                          | L1-07    | 10min   | guards green                                             |
| F-23 | `errorpage/notfound404.templ:42` → `bg-linear-to-br ...` (keep dark: variants)                                                                                                                     | L1-08    | 10min   | rg `bg-gradient` workspace = 0                           |
| F-24 | `website/internal/pages/hero.templ:107` → `bg-linear-to-r`; site goldens if affected                                                                                                               | L1-08    | 10min   | rg + site tests                                          |
| F-25 | Textarea: emit `field-sizing-content min-h-10 max-h-80` when AutoGrow; drop `tc-auto-grow` class emission                                                                                          | L1-09    | 12min   | golden shows utilities                                   |
| F-26 | Delete `.tc-auto-grow` from custom.css; sync `TestCustomCSSUtilities`                                                                                                                              | L1-09    | 10min   | rg `tc-auto-grow` = 0                                    |
| F-27 | Textarea goldens + manual e2e spot (grow behavior in demo, both themes)                                                                                                                            | L1-09    | 12min   | visual suite                                             |
| F-28 | Deep-dive report: add 21-group capability ledger appendix (or downgrade "14/21" to prose if ledger won't hold up)                                                                                  | L1-10    | 12min   | appendix present                                         |
| F-29 | Screenshot deep-dive HTML via visualtest chromedp one-off; eyeball hero/scorecard/tables/print                                                                                                     | L1-10    | 12min   | PNG reviewed                                             |
| F-30 | TODO_LIST.md: add short-term entries (tokens, inset, pins, scrollbar, field-sizing, bg-linear)                                                                                                     | L1-11    | 12min   | entries present                                          |
| F-31 | ROADMAP.md: idea entries (@utility tree-shaking, size-* sweep w/ templ window, masks/text-shadow/@container-size)                                                                                  | L1-11    | 10min   | entries present                                          |
| F-32 | Annotate status report + this plan as harvested (docs-health ANNOTATE style, inline note)                                                                                                          | L1-11    | 10min   | note present                                             |
| F-33 | `tc-*` usage audit: rg each class across component sources vs docs-only; write the dead/alive table                                                                                                | L1-12    | 12min   | table committed (report appendix or TODO note)           |
| F-34 | Measure current emitted custom.css bytes; project post-`@utility` (dead-class elimination)                                                                                                         | L1-12    | 12min   | numbers recorded                                         |
| F-35 | `@utility` migration batch 1: `tc-squircle`, `tc-content-auto(-compact)`, `tc-snap-*`                                                                                                              | L1-13    | 12min   | build + demo CSS fresh                                   |
| F-36 | `@utility` migration batch 2: `tc-fluid-*` (+ any audit-surviving utilities); keep component-shaped CSS as plain                                                                                   | L1-13    | 12min   | same                                                     |
| F-37 | Sync `TestCustomCSSUtilities` for `@utility` form; run full guard suite                                                                                                                            | L1-13    | 12min   | guards green                                             |
| F-38 | size-* sweep batch 1: display/                                                                                                                                                                     | L1-14    | 12min   | rg `w-N h-N` count drops                                 |
| F-39 | size-* sweep batch 2: forms/, feedback/, layout/, navigation/                                                                                                                                      | L1-14    | 12min   | same                                                     |
| F-40 | size-* sweep batch 3: htmx/, datastar/, errorpage/, utils/, recipes/ + `nix run .#build`                                                                                                           | L1-14    | 12min   | build green                                              |
| F-41 | Verify committed demo CSS byte-stability post-canonicalization (v4.2.2 collapsing)                                                                                                                 | L1-14    | 10min   | `nix run .#css` unchanged or refreshed deliberately      |
| F-42 | size-* goldens re-baseline — **GATED: execute only inside the templ v0.3.1070 migration window**                                                                                                   | L1-15    | 12min   | window checklist                                         |
| F-43 | Docs: note `user-valid:`/`user-invalid:` variant availability in validation docs                                                                                                                   | L1-16    | 10min   | prose present                                            |
| F-44 | Docs: "Tailwind capability coverage" section in adoption guide (v4.0 CSS-first, v4.2 logical insets, v4.3 scrollbar)                                                                               | L1-16    | 12min   | prose present                                            |
| F-45 | Lessons → `docs/agent-context-history.md`: recount-rule, heredoc-relapse, template-class-vs-guide verification                                                                                     | L1-16    | 12min   | entries present                                          |
| F-46 | Run `TestDocsCountDrift` after all docs edits; fix counts in same commit                                                                                                                           | L1-16    | 10min   | guard green                                              |
| F-47 | Write `scripts/check-html-report-classes.sh` (every class token in an HTML report exists in its CSS) + self-apply                                                                                  | L1-17    | 12min   | script green on both reports                             |
| F-48 | Add screenshot-pass checklist line for HTML deliverables (where: AGENTS docs-health section or script header)                                                                                      | L1-17    | 10min   | line present                                             |
| F-49 | Resolve `--ds-brand` per Q3: tokenize to `var(--color-*)` or document the external consumer in custom.css                                                                                          | L1-18    | 12min   | Q3 answer reflected                                      |
| F-50 | Resolve Q2 → mark ship-now or v2-bucket on finding 1 (and this plan's L1-19)                                                                                                                       | L1-19    | 5min    | decision recorded                                        |
| F-51 | If ship-now: warm CHANGELOG `[Unreleased]`, pre-verify lint + touched packages, then `nix shell nixpkgs#govulncheck -c nix develop -c scripts/release.sh` (never let release.sh be the first gate) | L1-19    | 12min×2 | release ritual green                                     |
| F-52 | Resolve Q1 → keep/fold artifacts; annotate accordingly                                                                                                                                             | L1-20    | 10min   | decision recorded                                        |
| F-53 | Pre-push ritual for any push containing F-output: `scripts/ci-repro.sh --lint --website` at the exact commit, then push immediately                                                                | standing | 12min   | `VERDICT: PASS`                                          |

_(All 35 status items are covered: f1 resolved during planning; f2-f35 map to F-01..F-52 above.)_

---

## Execution graph

```mermaid
flowchart TD
    subgraph PCT1["1% → 51% of value"]
        F01[F-01..F-03<br/>color literals → var(--color-*)] --> F04[F-04 theme-token guard]
        F04 --> F05[F-05 guard family]
        F05 --> F06[F-06..F-08<br/>goldens + visual re-baseline]
        F49{F-49 --ds-brand<br/>needs Q3} --> F01
    end

    subgraph PCT4["4% → 64% (cumulative)"]
        F09[F-09..F-11<br/>11 sites → inset-s-/e-] --> F12[F-12 generate + build]
        F12 --> F13[F-13 tighten RTL guard] --> F14[F-14..F-15 docs sync] --> F16[F-16..F-17 goldens + RTL visual]
        F18[F-18 Dockerfile pin] --> F20[F-20 AGENTS checklist]
        F19[F-19 website.yml pin] --> F20
    end

    subgraph PCT20["20% → 80% (cumulative)"]
        F21[F-21..F-22 scrollbar-none]
        F23[F-23..F-24 bg-linear]
        F25[F-25..F-27 field-sizing]
        F28[F-28..F-29 report integrity]
        F30[F-30..F-32 HARVEST]
    end

    subgraph REST["remaining 80% of work → 100%"]
        F33[F-33..F-34 @utility audit] --> F35[F-35..F-37 @utility migration]
        F38[F-38..F-41 size-* sweep] --> F42{F-42 goldens<br/>GATED on templ v0.3.1070}
        F43[F-43..F-46 docs + lessons]
        F47[F-47..F-48 tooling]
    end

    GATE2{F-50 Q2:<br/>ship now vs v2?}
    F06 --> GATE2
    F16 --> GATE2
    GATE2 -->|ship now| F51[F-51 CHANGELOG + release ritual]
    GATE2 -->|v2 bucket| DONE2[record decision]

    F05 --> RIT[F-53 ci-repro ritual] --> PUSH[push]
    F13 --> RIT
    F20 --> RIT
    F22 --> RIT
```

**Critical path:** F-49/F-01 → F-04 → F-06 → F-50 → F-51 (theming correctness ships first). Parallel lane with no dependencies: F-18/F-19 (lane pinning), F-23 (bg-linear), F-33 (audit). Longest-deferrable: F-42 (externally gated).

---

## Open questions blocking parts of the plan (from the status report)

1. **Q1 (F-52):** keep both session artifacts, or fold self-review into the deep-dive report and drop the status file?
2. **Q2 (F-50):** finding 1 ships as a patch release now, or waits for the v2 window (ADR-0039 convention)?
3. **Q3 (F-49):** does `--ds-brand` (custom.css:972-977) have a real external consumer, or can finding-1 tokenize it away?

---

_Point-in-time plan. When bringing it current later: docs-health ANNOTATE, never rewrite. If executed, harvest leftovers into TODO_LIST.md so this file can rest._
