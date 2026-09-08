# Status Report — Demo Visual Audit (7 pages × light + dark)

**When:** 2026-09-08 04:30 CEST
**Session scope:** Visual inspection of the `examples/demo` showcase — every route captured as a real-browser screenshot in light AND dark mode, every screenful actually viewed as an image, every suspect verified in source. No product code was changed; this was an audit-and-report session.
**Method:** Built the demo binary (`go build ./examples/demo`), served it on localhost:8901, drove headless Chromium via chromedp (same stack as `visualtest`) with a purpose-built capture tool, viewed all 28+ captures, tiled the 29,084px-tall index page into 28 viewable chunks, then root-caused each visual defect to exact `file:line`.

**Artifacts:** All raw captures retained at `/tmp/tc-shots/` (full-page PNGs, per-chunk JPEGs).

---

## Self-Critique (asked directly: what did I forget / could do better / still improve?)

**Forgot:**
- **Interactive states.** Every overlay component (Modal, Drawer, Dropdown, Popover, Tooltip, ContextMenu, Combobox, Carousel) was captured only in its *closed* state. `visualtest` already supports `State: Click/Hover` + `FullViewport` for opened top-layer elements — I reinvented a capture tool instead of reusing that harness and consequently have no opened-state screenshots.
- **Responsive + RTL.** Only one viewport (1280px). The library's core promises (container queries, MobileMenu hamburger, RTL logical-property mirroring) went visually unverified — one `EmulateViewport(375,…)` and one `dir="rtl"` JS call away.
- **Transport variants.** Only `?transport=both` captured; the `htmx`/`datastar` single-dialect index variants were never screenshotted.
- **A flagged item was dropped.** The 45% ProgressBar fill looked washed-out in both modes; I noted "check against goldens" and never resolved it. A finding without a verdict is a loose end.
- **Daemon exposure.** I created the temp capture tool *inside* the repo module. The BuildFlow daemon auto-committed it (then its deletion) into master history — 2 polluted commits. It should have lived in `/tmp` or a sanctioned tooling location from the first write.

**Could have done better:**
- The capture tool took 3 iterations (shared-browser multi-minute hangs → probe → fresh-browser-per-page). Reading chromedp's known gotchas (or the existing `visualtest` lifecycle code) before writing would have made it one iteration.
- Root-cause confidence was uneven: the Heatmap and LoadingButton findings are proven by code read; the AppShell SM-vs-`w-64` overlap mechanism is confirmed by code + screenshots but not DOM-measured (no `getBoundingClientRect` numbers).
- I confirmed the Heatmap root cause statically but never ran the 30-second dynamic check (set a defined var, re-screenshot) to close the loop experimentally.
- Demo server log was glanced at, not systematically grepped for 500s/warnings during the capture run.

**Could still improve:**
- Turn ad-hoc inspection into a repeatable, sanctioned workflow (see f-items 1–3) so the next "how is our demo?" is one command, not a session.
- Pair every found defect with a failing guard (visual golden or unit test) in the same fix commit, per repo convention.
- Add a11y (axe scan, keyboard-only pass) and offline/CDN-failure rendering to the audit checklist.

---

## a) FULLY DONE

| # | Item | Evidence |
|---|------|----------|
| 1 | Full visual audit executed: 7 routes (`/`, `/forms`, `/users`, `/recipes/dashboard`, `/recipes/settings`, `/recipes/login`, `/recipes/auth`) × light+dark = 14 full-page captures; index (29,084px tall) tiled into 14+14 chunks and **every chunk viewed** | `/tmp/tc-shots/` |
| 2 | All three component-level defects root-caused to exact lines | `layout/appshell.templ:10`, `display/heatmap.go:136`, `htmx/loading.templ:48-55` |
| 3 | All demo-content defects verified in source (not just visually) | `demo.templ:66-67`, `layout_demo.templ:72-80`, `display_demo.templ:245-252`, `forms_section.templ:91-107`, `recipes_demo.templ:69`, `navigation_demo.templ:26-33` |
| 4 | Icon count ground-truthed by an actual test run: `iconPathData=101`, `AllIconNames=102` — docs claiming 106 are stale | temp test in `icons` module, since removed |
| 5 | Capture tooling proven: fresh-browser-per-page pattern captures any demo page in ~3s (shared-browser sessions hang after very tall captures) | session-validated |
| 6 | Repo hygiene: temp tool removed, demo server stopped, unrelated in-flight edits from another session (docs/DOMAIN_LANGUAGE.md, docs/recipes/horizontal-filter-bar.md) detected and **left untouched** | `git status` |
| 7 | Findings delivered with severity ordering and mechanism explanations | previous message |

## b) PARTIALLY DONE

| # | Item | Gap |
|---|------|-----|
| 1 | Overlay/interactive inspection | closed states only; no opened Modal/Drawer/Dropdown/Popover/Tooltip/ContextMenu captures |
| 2 | Responsive audit | 1280px only; no 375px/768px captures; MobileMenu + ContainerAware collapse unverified visually |
| 3 | RTL audit | zero `dir="rtl"` captures despite the library's logical-property promise |
| 4 | Wire transport variants | only `?transport=both` screenshotted |
| 5 | E2E interaction | only `/health` + one raw HTML fetch; no browser click-throughs (LoadMore, ConfirmDelete, wire form roundtrip in *this* demo context) |
| 6 | AppShell SM-overlap finding | mechanism confirmed, but not DOM-measured for the exact overflow amount |
| 7 | ProgressBar 45% fill-color suspicion | flagged, never resolved against `visualtest/testdata/progressbar/*` goldens |
| 8 | Daemon pollution cleanup | temp tool deleted, but its add+delete daemon commits remain in master history |

## c) NOT STARTED

1. Fixing any of the ten findings — zero product code changed this session (report-only, by design).
2. a11y audit of the demo (axe, keyboard-only traversal, screen-reader pass).
3. Offline / CDN-failure rendering check (Google Fonts fallback behavior).
4. Prerendered static pages verification (`demo -prerender`; TODO #154 sync concern).
5. Live-server 404 route capture (skipped deliberately — default `http.NotFound`, not an errorpage render).
6. `TestCSSFreshness` / demo.css recompile state verification.
7. `docs-health` HARVEST of section (f) into `TODO_LIST.md` / `ROADMAP.md` (skill requires it after the report — **awaiting instructions**).

## d) TOTALLY FUCKED UP

**Product breakage found (severity-ordered):**

| # | What | Root cause | Severity |
|---|------|-----------|----------|
| 1 | `/recipes/dashboard` renders **collapsed** — stat cards clipped to slivers ("$4", "2.5"), charts/activity one word per line, ~75% of viewport empty. Both modes. | `layout/appshell.templ:10` — grid template `lg:grid-cols-[var(--tc-sidebar-w)_minmax(0,1fr)]` is **unconditional**; with `Sidebar == nil` the single content column lands in the 16rem sidebar track, the second column stays empty. The flagship recipe is unusable. | CRITICAL |
| 2 | `display.Heatmap` shows **no cells at all** — only row/column labels + the peak cell's ring outline. Both modes. | `display/heatmap.go:136` emits `rgba(var(--ds-brand-rgb), α)`; `--ds-brand-rgb` is defined **nowhere** (not in `templates/custom.css`, `templ-components-theme.css`, or demo CSS) → invalid CSS → declaration dropped → transparent cells. Ships broken by default; only `HighlightPeak` ring (different property) is visible. | CRITICAL |
| 3 | `htmx.LoadingButton` spinner **always visible at rest** — "⟳ Save changes", "⟳ Run job" (wire Busy-state HTMX side; the Datastar twin is correct). | `htmx/loading.templ:48-55`: the spinner `templ.Component` renders un-gated; only `loadingText` carries `htmx-indicator`. Bonus: `tc-btn-loading` class (loading.templ:49) is defined nowhere in `templates/custom.css`. No visual golden exists for this component — which is exactly why it slipped through. | HIGH |
| 4 | Demo hero says **"HELLO"** — leftover debug line, the first content visitors read. | `examples/demo/demo.templ:66-67` | HIGH (credibility) |
| 5 | AppShell inline demo doubly broken: sidebar nav overlaps content ("Ma\|in content area" clipped) **and** a ~900px black column bleeds into the Cards section. | `layout_demo.templ:72-80` passes `SidebarWidthSM` (12rem) but `SidebarNav` hardcodes `w-64` = 16rem (`navigation/sidebar_nav.templ:123`) → 64px overflow over the content column; AppShell's `min-h-dvh` (appshell.templ:10) inside a bounded demo card stretches the black column to viewport height. | HIGH |
| 6 | DateRange run-together: "March 2023 – Present**Jan 15, 2024** – Jun 30, 2025". | `display_demo.templ:245-252` — two inline DateRange spans adjacent; vertical margins don't apply to inline boxes. | MEDIUM |
| 7 | Index filter bar: `Status` select spans the full 1280px container while `Sort`+button sit compact below — looks broken next to the correct hand-rolled filter bar on `/forms`. | `forms_section.templ:91-107` — `forms.Select` inside `FormLayoutInline` ignores inline layout (renders full width). | MEDIUM |
| 8 | Component-count chaos: auth recipe claims **116** (`recipes_demo.templ:69`), hero const/FEATURES say **120**, SKILL.md says **118**. Three numbers in circulation. | hardcoded copy, never swept when counts moved | MEDIUM |
| 9 | Icon-count docs drift: AGENTS.md says 106, SKILL.md implies 105; reality (test-verified) **102**. The demo computes it correctly — docs lie. | stale docs after icon set changes | LOW |
| 10 | "Full Nav" demo renders **no brand**, no right slot — reads as unfinished next to SimpleNav two sections up. | `navigation_demo.templ:26-33` passes none | LOW |

**Session's own fuckup (owned):** temp capture tool placed inside `visualtest/tools/` → BuildFlow daemon auto-committed the file AND its deletion into master history. Tooling that touches the repo must be sanctioned or live outside it.

**Explicitly NOT fucked up (verified working, to prevent false alarms):** Datastar SSE LiveRegion live-patches ("SSE update #1 — streamed at 21:37:11" — the light capture's "Waiting for SSE stream…" was pure timing), PolledRegion auto-refreshes with its indicator, ECharts dark-mode bridge re-themes the chart, dark mode is excellent everywhere, icon gallery + search renders beautifully, `/forms` + Settings + Login + Auth recipes are polished.

## e) WHAT WE SHOULD IMPROVE

1. **Fix with guards, not patches.** Each of the three component bugs gets a failing test first: a visual golden with `Sidebar: nil` AppShell, a Heatmap default-var golden (currently NO golden pins the default render — that's the systemic hole), a LoadingButton rest-state golden.
2. **Close the no-golden gap.** Components with zero visual goldens today include LoadingButton, Heatmap default, AppShell variants — the exact three that broke. A coverage sweep "every demo section has at least one golden" would have caught all three before release.
3. **Sanction the inspection workflow.** A permanent `.#shots` flake app (fresh-browser-per-page, all routes, light+dark, optional `--mobile/--rtl/--open-overlays`) makes this audit a one-command regression sweep. Policy question below (g3).
4. **Demo copy hygiene.** The "HELLO" line and the 116/118/120 count chaos are the kind of drift the repo's own drift-guard philosophy exists to prevent — a demo copy test (like `demo_counts_test.go`) should own every user-visible number.
5. **Keep audit artifacts out of repo history.** Rule adopted for future sessions: throwaway tooling goes to `/tmp` with a module that reuses the module cache; only sanctioned tools enter the tree.
6. **Resolve every flag before closing.** The ProgressBar fill-color suspicion should never have survived the session unresolved — verdict or delete, same pass.

## f) NEXT — up to 50 things to get done

**P0 — broken user-visible behavior (fix now)**
1. AppShell: emit the single-column template (no sidebar track) when `Sidebar == nil`; add golden `appshell/nosidebar_light` + fix `/recipes/dashboard` render.
2. Heatmap: define `--ds-brand-rgb` (and dark variant) in `templates/custom.css`, or flip default `ColorVar` to a Tailwind-defined var — decide per g2; add default-render golden.
3. LoadingButton: wrap the spinner in `<span class="htmx-indicator">`; add rest-state + `.htmx-request` goldens; delete or define `tc-btn-loading`.
4. Remove the "HELLO" debug line from the demo hero (`demo.templ:66-67`) — keep or rewrite the tagline below it (g1).
5. AppShell demo: use `SidebarWidthMD` (matches `w-64` SidebarNav) and override `min-h-dvh` via `Class` for the inline context; re-screenshot.
6. Fix the two DateRange demo spans (wrap each in a block or render DateRange as block).
7. Fix `forms.Select` inside `FormLayoutInline` (component-level width constraint) or restructure the index filter bar like `/forms`.
8. Sweep every hardcoded user-visible count: "116 components" (recipes_demo.templ:69) → single source of truth; decide canonical count via `demo_counts_test.go` (FEATURES = 120).
9. Update icon counts in AGENTS.md (106 → 102) and SKILL.md ("104 path icons" → 101 + Spinner); grep for other stale counts.
10. Give "Full Nav" demo a brand + right slot so it actually showcases Nav vs SimpleNav.

**P1 — close this audit's gaps (verification debt)**
11. Capture opened overlays: Modal, Drawer, Dropdown, Popover, Tooltip, ContextMenu, Combobox, Carousel (reuse `visualtest` `State: Click` + `FullViewport`).
12. Mobile viewport sweep (375px): MobileMenu hamburger, ContainerAware collapse, form stacking, table overflow.
13. RTL sweep (`dir="rtl"`): verify logical-property mirroring on Nav, Split, Carousel, Drawer, Dropdown.
14. Capture `?transport=htmx` and `?transport=datastar` index variants.
15. Resolve the ProgressBar 45% fill-color suspicion against `visualtest/testdata/progressbar/half_light.png`.
16. Browser-verify LoadingButton *during* request (default text hides, loading text + spinner show) — the rest state is proven, the request state is not.
17. Dynamically confirm the Heatmap fix hypothesis (define var → cells appear) before cutting the fix.
18. DOM-measure the AppShell SM-vs-`w-64` overflow (getBoundingClientRect) to replace pixel estimates.
19. Click-through E2E on the demo itself: LoadMore → EndOfList, ConfirmDelete → row removal, wire form roundtrip, busy-state 800ms sleep, file upload echo.
20. Grep full server log for 500s/warnings during a capture run; assert zero.

**P2 — coverage & guard hardening**
21. Golden-coverage sweep: assert every demo section has ≥1 visual golden; write the missing ones (LoadingButton, Heatmap default, AppShell variants top the list).
22. Drift-guard for demo copy: test that user-visible numbers in the demo come from single sources (version, counts, icon totals).
23. Drift-guard for docs icon/component counts (AGENTS/SKILL vs computed reality) — extend the existing `TestDocsCountDrift` pattern.
24. Decide + implement AppShell's `min-h-dvh` opt-out for embedded contexts (Class override documented, demo uses it).
25. Guard test: `FormLayoutInline` children must not render full-width (or document the constraint loudly).
26. Add `DateRange` block-vs-inline semantics to its docs + a two-adjacent-ranges golden.
27. Add Heatmap dark-mode golden once the var exists (opacity on dark bg needs its own check).
28. Consider a `TestAllDemoRoutesRender` harness (httptest + every route + status 200 + zero-template-error).

**P3 — tooling & workflow**
29. Decide repo policy for inspection tooling (g3); if sanctioned: `visualtest/tools/shots` as its own module + `nix run .#shots` flake app with `--mode light|dark|both --mobile --rtl --open` flags.
30. Fold fresh-browser-per-page (or documented tab-reuse limits) into `visualtest` docs — the multi-minute hang is a real chromedp gotcha worth recording in AGENTS.md.
31. Commit the audit screenshots triage set (or a trimmed gallery) into `docs/` if the team wants persistent evidence; else document `/tmp` retention.
32. Prerender verification: run `demo -prerender` into a temp dir, diff against live server HTML for the 7 routes (TODO #154 sync).
33. Offline render check: block fonts.googleapis.com, screenshot, confirm acceptable fallback typography.
34. a11y pass: axe-core scan per route; keyboard-only traversal script; fix what it finds.
35. `TestCSSFreshness` state check + recompile demo CSS if stale before next release.

**P4 — docs & release hygiene**
36. HARVEST this report's (f) list into `TODO_LIST.md` (P0/P1) and `ROADMAP.md` (P2+ ideas) — per docs-health HARVEST mode. **Awaiting go-ahead.**
37. Warm `CHANGELOG [Unreleased]` when the P0 fixes land (repo convention: entry in the same commit).
38. AGENTS.md "Demo Infrastructure" section: add the known-issues list or link to this report until fixed.
39. SKILL.md component-count line (118) → align with canonical source.
40. Consider an "empty-slot contract" note for AppShell in `docs/` (what happens when Sidebar/Header/MobileNav are nil) — today's dashboard collapse proves the contract was implicit.

**P5 — roadmap fuel (bigger ideas surfaced by the audit)**
41. Interactive-state visual goldens as a first-class tier (click-to-open captures for all overlay components).
42. Multi-viewport golden matrix (1280/768/375) for the shell primitives (AppShell, Nav, Split, Grid).
43. Scheduled demo smoke job in CI (build → serve → screenshot → diff against tolerance) so demo rot is caught between releases.
44. Demo "theme override" toggle showcasing `@theme` remapping (sells the theming model visually).
45. Print/PDF rendering check for the demo (library claims `Minimal`/static-friendliness).
46. Performance pass: index page is 29k px tall / ~1.4MB HTML — consider lazy sections or pagination for the mega-page.
47. CSP negative test in the demo: serve with a strict CSP header and assert zero console violations during a capture run.
48. Icon gallery: add copy-to-clipboard feedback verification + keyboard operability check (click-to-copy is mouse-only today?).
49. Error pages: capture ErrorPage family matrix (rejection/conflict/transient/corruption/infrastructure) as goldens — currently only two exist.
50. Session retro item: codify "audit sessions produce reports + failing guards, never drive-by fixes" as a documented convention so fix scope stays reviewable.

## g) Questions I cannot figure out myself

1. **The "HELLO" line in the hero** (`demo.templ:66-67`) — is that leftover debug output you want deleted, or is "same component, either runtime — see the Wire section below" an intentional tagline you want kept (rewritten)? And was "HELLO" a deliberate canary you use to verify deploys?
2. **Heatmap theming contract**: should the default stay `--ds-brand` (and we define `--ds-brand-rgb` in `templates/custom.css` as the library-provided default, keeping the override hatch), or should the default flip to a Tailwind-native color (e.g. blue-600) with `ColorVar` documented as the opt-in? This touches the public theming model, so I won't pick unilaterally.
3. **Inspection tooling policy**: do you want a permanent, sanctioned demo-screenshot utility in the repo (own module + `nix run .#shots`), or should audit tooling always live outside the repo (as this session's daemon pollution demonstrated)? I can build either; I can't infer which you want maintained.

---

*Report generated per the status-report skill. Format override: user explicitly requested `.md`; the skill's canonical HTML dashboard was skipped this once. Section (f) is the primary input for a future `docs-health` HARVEST run.*
