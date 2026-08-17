# Status Report — Demo Page "Is It Superb?" Audit & Overhaul

**Date:** 2026-08-17 14:26
**Session scope:** Audit `examples/demo` against the question "Is our demo page superb?", then fix everything found. Report covers this session only.

---

## The Verdict That Drove This Session

**No — the demo was good, not superb.** The audit (full source read of all 13 demo templates + `main.go` + live HTTP audit of every route) found:

| #  | Finding                                                                                                                                                                                              | Severity |
| -- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| 1  | Eager `PolledRegion` polled `/api/demo-stats` → **404 → error toast fired on every page load**                                                                                                       | Critical |
| 2  | 5 demo'd endpoints didn't exist (`/api/items`, `/api/items/123`, `/api/users`, `/api/demo-stats`, `/users`); 2 working mocks (`/api/load-more`, `/api/delete`) were orphaned (referenced by nothing) | Critical |
| 3  | Hero badge said **v0.17.0**; actual version **1.8.4**                                                                                                                                                | High     |
| 4  | Hero claimed **107 components / 10 packages**; FEATURES.md canonical totals: **116 / 11**                                                                                                            | High     |
| 5  | Sticky-nav TOC "Datastar" pointed at `#datastar-sdk-section`; real anchor was `datastar-sdk-script` — dead link in the primary nav                                                                   | High     |
| 6  | "Layout" section not in TOC at all                                                                                                                                                                   | Medium   |
| 7  | 4 pages unreachable by any link: `/forms`, `/recipes/dashboard`, `/recipes/settings`, `/recipes/login`                                                                                               | High     |
| 8  | Shipped-but-never-demoed: `SectionHeading`, `DateRange`, `CircularProgress`, `LoadingOverlay`, `ViewTransitions`, `SSEErrorHandling`, `AuthLayout` — and the **entire `charts/echarts` module**      | High     |
| 9  | `recipes.LoginCard` FormBody + AuthLayout demo content was a giant hand-rolled `templ.Raw` HTML string — the component library not using its own components                                          | Medium   |
| 10 | Nav demo links pointed at nonexistent `/about`, `/contact`                                                                                                                                           | Medium   |
| 11 | Accordion item 3 said "active development (v0.x), API may change before v1.0" — false at v1.8.4; CollapsibleSection debug block hardcoded "v1.6.0"                                                   | Medium   |
| 12 | `/users` StatCard href + DataTable `SortBaseURL` + Pagination `BaseURL` all targeted a 404 page                                                                                                      | High     |
| 13 | No results target for the FilterDropdown demo (`#results` referenced, never existed in DOM)                                                                                                          | Medium   |

---

## a) FULLY DONE

1. **Full audit** — read every demo template + main.go; ran the server (port 8091) and probed 17 routes with a Go HTTP client (curl is banned); grep-verified component coverage vs FEATURES.md catalogue.
2. **`main.go` rewritten backend**:
   - `/api/items` (LoadMore pagination: batch → next button → EndOfList), `/api/items/123` (ConfirmDelete, live target `#item-123` in DOM), `/api/save` (600ms delay so LoadingButton state is visible), `/api/demo-stats` (self-replacing PolledRegion response that **settles after 3 ticks** instead of polling forever), `/api/users` (real filtered fragment via templ).
   - New route `/users` — full server-driven page: DataTable sort (Name/email, asc/desc via `?sort=&dir=`), Pagination round-trips (`?page=`), Breadcrumbs, Card(flush). This is the flagship "here's the pattern the library is built for" page.
   - New route `/recipes/auth` — the never-demoed `recipes.AuthLayout`.
   - All fragment handlers render via templ (`componentOr500`) instead of string literals.
3. **Hero fixed**: version now renders `utils.Version` dynamically (can never go stale); component count 107→**116**, packages 10→**11** (canonical FEATURES.md numbers); icon count computed from `icons.AllIconNames()` (self-maintaining).
4. **Drift-guard test added** (`examples/demo/demo_counts_test.go`): `TestHeroCountsMatchFeatures` parses FEATURES.md Totals line and fails if the demo consts or icon count drift; also asserts the hero renders the live `utils.Version`.
5. **TOC fixed**: broken Datastar anchor corrected; Layout + ECharts + Recipes entries added.
6. **8 missing component demos added**: SectionHeading, DateRange (both variants incl. "Present"), CircularProgress (3 sizes/colors/label), LoadingOverlay (pinned-inside-a-frame variant via Class override, so it demos without hijacking the page), htmx.ViewTransitions, datastar.SSEErrorHandling, echarts.SDKScript + echarts.EChart (hand-authored RenderSnippet-equivalent strings — demo stays go-echarts-free, matching the adapter's zero-dep design), recipes.AuthLayout.
7. **New "Recipes & Full Pages" section** on the home page: 6 link cards (Users, Forms, Dashboard, Settings, Login, Auth) — every demo page now reachable in ≤1 click.
8. **Dogfooding fixes**: LoginCard FormBody now uses `forms.Form` + `forms.Input` + `display.Button` instead of raw HTML; Dashboard "Revenue trend" card now uses `display.LineChart` (was a hand-rolled flex-bar chart), "Recent activity" uses Avatar + RelativeTime (was "No recent activity"); `/forms` page wrapped in `forms.Form` (was a bare div with a submit button outside any form — invalid HTML); filter bar results targets (`#results`, `#user-results`) now exist.
9. **Stale copy fixed**: accordion "v0.x" → dynamic version; debug block version/time dynamic; nav links → real pages; breadcrumbs → real hrefs.
10. **prerender.go** updated with both new pages (`/users` is NOT in prerender — flagged below).

## b) PARTIALLY DONE

1. **Final build verification NOT re-run.** The last executed `templ generate + go build` failed with 7 compile errors; I fixed all 7 (ViewTransitionsProps has no Nonce field → BaseProps; `icons.User` doesn't exist → UserCircle; ButtonProps has no Class field → BaseProps.Class; Pagination wants uint; users_demo.go/templ split to fix a redeclared-import issue) — but the session was interrupted **before re-running generate+build**. State: fixes applied, compilation unconfirmed. The BuildFlow daemon committed everything anyway (see d).
2. **CSS recompilation not done.** New demo sections add Tailwind classes (`recipes-*` cards, CircularProgress, polled-region styles, ECharts container). The committed `static/app.css` is stale until `nix run .#css` (or full build) runs — this repo has `TestCSSFreshness` which **fails in CI** on stale CSS. The daemon's commit 7c2e3c2 regenerated `demo.out.css` but that is not the embedded `static/app.css`.
3. **Full verification suite not run**: `nix run .#test` (includes my new drift-guard test), `nix run .#lint`, golden tests untouched-but-unproven, live re-probe of all routes (the old binary on :8091 still serves the pre-fix code).

## c) NOT STARTED

1. `/users` page missing from `prerender.go` page list (only server-served).
2. Visual/screenshot pass of the new sections (light+dark, mobile) — `nix run .#visual` not attempted.
3. `demo.out.css` is a committed artifact that appears to be dead weight (static.go embeds `static/app.css`; nothing references `demo.out.css`) — not investigated, not removed.
4. CHANGELOG `[Unreleased]` entry for the demo overhaul — **repo rule: every master commit keeps [Unreleased] warm; the daemon committed without one.**
5. AGENTS.md demo section not updated with the new routes/pages.
6. Icon-search and section-filter demo JS still use ad-hoc scripts — fine (nonce'd), but untested this session.
7. Hero "Documentation" link → `https://templcomponents.lars.software` — never verified live this session.

## d) TOTALLY FUCKED UP

1. **Three templ parse errors on first generate** — all mine: (a) the word "for" starting a text line inside a `<p>` in htmx_demo (templ treats it as a loop keyword — must be reworded or escaped); (b) indentation-sensitive brace miscount in forms_demo after inserting the Form wrapper (took 3 edit rounds; `edit` re-indent warnings made it worse before better); (c) same miscount pattern in recipes_demo. Lesson: when wrapping templ children in a new component block, re-indent the **entire** child block, not just the first line — or rewrite the file.
2. **BuildFlow daemon race**: mid-session, the auto-commit daemon committed my half-verified work (97e7bb0) and then a follow-up commit (7c2e3c2) **modified my file** (`users_demo.go`, "simplify pagination ceiling" — replaced my `if pages < 1` clamp with something shorter) and regenerated CSS. This is documented expected behavior, but it means: unverified code is on master, and my file now contains an edit I did not author and have not reviewed.
3. The very first demo-server attempt raced an existing process on :8080 (leftover listener from a previous session is still occupying 8080 — never killed, still running).
4. I initially wrote `users_demo.templ` as one file mixing heavy Go + templ, hit a redeclared-import error from generated code, then split Go out to `users_demo.go` — one wasted generate cycle that the daemon then committed as churn.

## e) WHAT WE SHOULD IMPROVE (session learnings)

1. **A demo is a test surface.** Every `hx-*`/`data-*` URL in demo templates should be probed by a test — a dead demo endpoint is a silent lie to every evaluator. The `TestHeroCountsMatchFeatures` pattern should extend to "every href/hx-get/hx-post in rendered demo pages resolves to a registered route" (one table-driven test, kills finding #2 permanently).
2. **Never hardcode version/counts in a demo.** v0.17.0-at-v1.8.4 proves copy rot. Anything numeric should come from `utils.Version`/`AllIconNames()`/FEATURES.md or a guard test fails.
3. **Dogfood or it didn't happen.** The library's own recipe demos containing raw HTML strings undermines the pitch. Every demo surface should render through the library.
4. **Demos must settle.** Eager infinite polling hammering a 404 was the single worst bug — demo'd interactivity should terminate (my 3-tick settle pattern) or be user-armed.
5. **Daemon commits unverified work.** Known issue (AGENTS.md T13/T1 family). This session is another data point: 2,200+ lines of demo changes hit master without `go build`. Mitigation for next session: run verify _before_ the daemon's 60s window, or accept CI as the gate.
6. **templ + `edit` tool friction**: surgical edits on indentation-sensitive templ blocks are error-prone; prefer `write` for restructured regions.

## f) NEXT — up to 50 things

**Verification (do these first):**

1. `nix develop -c 'templ generate ./... && go build ./...'` — confirm my 7 error-fixes compile
2. Review daemon's `users_demo.go` tweak in 7c2e3c2 (I didn't author it)
3. `nix run .#css` — recompile `static/app.css` with the new classes
4. `nix run .#test` — includes new `TestHeroCountsMatchFeatures` + `TestCSSFreshness`
5. `nix run .#lint`
6. Re-run live route probe (all 17 old paths + `/users` + `/recipes/auth` + new API endpoints)
7. `nix run .#verify` as the single done-check
8. Add CHANGELOG `[Unreleased]` entry (demo overhaul + new routes)
9. Kill the stale process squatting on :8080
10. Check `.gitignore` for BuildFlow's re-appended `*_templ.go` line (documented gotcha; new files `users_demo_templ.go`, `echarts_demo_templ.go` must be tracked — verify with `git ls-files`)

**Coverage gaps still open:**
11. Add `/users` to `prerender.go`
12. Demo `forms.Select{Stylable: true}` (customizable `<select>` API — shipped, not shown)
13. Demo `layout.Minimal` (static/PDF shell)
14. Demo `layout.Script`/`Stylesheet` CSP helpers explicitly
15. Demo `htmx.SwapOOB` with a live trigger (button that returns an OOB swap) — currently static display only
16. Demo `display.Table` `LazyRows` (content-visibility) with a 100+ row table
17. Demo `TableRow.Href` clickable rows on `/users`
18. Demo `Table.Caption`
19. Demo `StatCard.Href` already used — add `TrendWarn`/`TrendNone` legend explanation
20. Demo `display.Image` `SrcSet`/`Sizes` responsive delivery
21. Demo `feedback.Skeleton` remaining variants (7 exist, ~4 shown)
22. Demo `forms.Textarea` `AutoGrow` + `EnterKeyHint` explicitly labeled
23. Demo `Input.EnterKeyHint` variants
24. Demo `display.CollapsibleSection` `StorageKey` persistence
25. Demo `icons.AnimatedIcon` / all 11 animation presets (only referenced in docs, not shown!)
26. Demo `Badge` all 7 types (5 shown)
27. Demo `Button` loading state + icon-only variants
28. Demo `Modal` other 4 sizes + `Drawer` left side
29. Demo `Tooltip` left/right positions + Escape dismiss behavior note
30. Demo `errorpage.ErrorHandler` JSON mode + `WriteError` via a `/api/error` endpoint
31. Demo `datastar` Indicator with a real fetching action (currently wired, untested)
32. Demo `AppShell` `MobileNav` slot (drawer on mobile)
33. Demo `Nav.ContainerAware` hamburger-by-container-width
34. Demo fluid typography `.tc-fluid-*` remaining classes
35. Demo `Pagination` ellipsis edge cases (page 1/last)

**Quality/infra:**
36. Write the "every demo'd URL resolves" test (see e.1) — biggest ROI item on this list
37. Extend drift-guard: demo section count vs `demoSection(` call count vs TOC entries (anchor existence test)
38. Investigate/remove `demo.out.css` if dead
39. Add visual-test goldens for `/users` and `/recipes/auth`
40. Dark-mode + mobile screenshot sweep of new sections
41. Update AGENTS.md "Demo Infrastructure" with new routes/pages
42. Update SKILL.md/README if they reference demo routes
43. Lighthouse/a11y pass (focus order on new pages, heading hierarchy — `/users` has h1→Card title hierarchy worth checking)
44. Consider a `/api/error` 500 endpoint to _intentionally_ demo GlobalErrorHandling's toast (today it only fires on real bugs)
45. Add `<noscript>` notice? (library is server-rendered-first; HTMX bits degrade — worth stating on the demo)
46. Hero: link "116 Components" stat to the catalogue section; wire stat cards to anchors
47. Prerender output → serve as static preview on the website (docs synergy)
48. Dedupe: `componentOr500` in main.go vs error handling in renderPage — minor, but same pattern twice
49. Consider extracting demo mock-data (`demoUsers`) into a fixture shared with tests
50. Track down why the daemon's message for 97e7bb0 was accurate this time ("feat(demo): expand showcase...") — the T13 stale-message bug may be fixed; worth confirming before more sessions rely on it

## g) QUESTIONS I CANNOT ANSWER MYSELF

1. **Is `https://templcomponents.lars.software` the canonical docs URL the demo should link?** It's what the hero linked before I touched it, but I found no source of truth in-repo this session (website/ exists but I didn't audit its routing) — confirm or correct the hero's Documentation link target.
2. **Should the demo's mock endpoints live behind `/api/*` in `main.go` forever, or do you want a separate `demo_handlers.go`?** The file is now ~240 lines mixing routes + fragments; I kept it one-file for diff review, but it's at the threshold where you may have a preference.
3. **The BuildFlow daemon modified `users_demo.go` after my edit (commit 7c2e3c2).** I didn't author that change and haven't reviewed it — do you want me to review/revert it on the next run, or is daemon-tidying of in-flight files accepted behavior I should trust?

---

**Bottom line:** Audit complete and verdict delivered; ~90% of the fixes are implemented and committed (by the daemon), but the session ended **before the final generate+build+test verification pass** — that is step 1 of the next session, no exceptions.
