<!-- AUTO-UPDATED 2026-07-10: Retrospective status overlay -->

> ## 🔔 Update Notice — 2026-07-10
>
> This report is **historical**. Many items listed as "open", "todo", or "broken" below
> have since been **fixed and verified**. Do not act on open items without first checking
> [TODO_LIST.md](../../TODO_LIST.md) for current status.
>
> **Key fixes completed since this report:**
>
> - ✅ All 7 P0 bugs fixed (InlineLoadingOverlay a11y, SanitizeID mismatch, FromError fallback,
>   Footer BaseProps, ErrorPage/NotFound404 `<main>` landmark, CSRFTokenName, grid-rows verified)
> - ✅ `encoding/json/v2` purged from all production code + pre-commit guard added
> - ✅ Motion constants centralized in `utils/motion.go`, wired into 13 components
> - ✅ `FamilyFromErrorFamily` → `FromErrorFamily` (old name kept as deprecated alias)
> - ✅ `icons.IconRTL()` + CSS for directional icon RTL mirroring
> - ✅ 33 regression tests added (htmx, errorpage, layout, navigation, feedback, display)
> - ✅ Dark golden test infrastructure (badge/card/button)
> - ✅ CHANGELOG consolidated, ROADMAP updated, migration guide created
> - ✅ All 14 packages pass, 0 lint issues
>
> **Canonical source of truth:** [TODO_LIST.md](../../TODO_LIST.md) (52 items, 37 ✅ done, 12 deferred/blocked)

---

# Status Report — 2026-07-05 16:56 CEST

> **Updated:** 2026-07-06 (post-v0.8.0). Version at report: 0.6.x → **Current:** 0.8.0

**Session scope:** Process `docs/feedback/*` consumer feedback end-to-end — implement improvements, self-review, fix gaps, rewrite skill, close all remaining items.
**Commits:** 6 (`985019a` through `29fccf1`)
**Done-check:** `nix run .#verify` → All checks passed (0 issues)
**Git:** Clean, pushed to `origin/master`

> **UPDATE NOTE (2026-07-06):** This was the final session before v0.7.0 release.
> Since then, sessions 7–10 + v0.7.0/v0.8.0 releases landed. All open items resolved.

---

## a) FULLY DONE

### Features shipped (commit `985019a`)

| #   | Component / Feature                                  | Files                                        | Tests                                             |
| --- | ---------------------------------------------------- | -------------------------------------------- | ------------------------------------------------- |
| 1   | `display.Grid` + `GridCols` enum (1–6)               | `display/grid.templ`, `grid_templ.go`        | golden, BDD, a11y, example, integration, coverage |
| 2   | `StatCardProps.Href` — renders `<a>` wrapper         | `display/card.templ`, `card_templ.go`        | golden, BDD, a11y, integration, coverage          |
| 3   | `SimpleNavProps.RightItems` — forwarded to Nav       | `navigation/nav.templ`, `nav_templ.go`       | coverage, BDD                                     |
| 4   | `layout.Script(nonce, src, attrs)` — CSP-safe helper | `layout/script.templ`, `script_templ.go`     | golden, BDD, a11y, example, snapshot              |
| 5   | `feedback.SkeletonCardGrid(count)` — loading grid    | `feedback/loading.templ`, `loading_templ.go` | golden, BDD, a11y, example, snapshot              |

### Recipes shipped (commit `985019a`)

| #   | Recipe                                               | Path                                                  |
| --- | ---------------------------------------------------- | ----------------------------------------------------- |
| 6   | Play CDN → Tailwind v4 CSS-first (7-step migration)  | `docs/migration/play-cdn-to-tailwind-v4.md`           |
| 7   | Server-rendered HTMX error feedback (3 render modes) | `docs/recipes/server-rendered-htmx-error-feedback.md` |

### Documentation (commits `985019a`, `79c926c`, `865967a`, `cea5a66`, `29fccf1`)

| #   | What                                                                   | Files                            |
| --- | ---------------------------------------------------------------------- | -------------------------------- |
| 8   | PageProps auto-inject godoc (HTMXVersion, CSSPath suppression)         | `layout/base.templ`              |
| 9   | README "Suppressing auto-injected `<head>` tags" subsection            | `README.md`                      |
| 10  | AGENTS.md: 8 new convention entries, count 25→26                       | `AGENTS.md`                      |
| 11  | TODO_LIST.md: session 6 header + Consumer Feedback Backlog (11 items)  | `TODO_LIST.md`                   |
| 12  | FEATURES.md: Grid, Script, SkeletonCardGrid, GridCols enum             | `FEATURES.md`                    |
| 13  | CONTEXT.md: updated metrics, package descriptions                      | `CONTEXT.md`                     |
| 14  | CHANGELOG: comprehensive `[Unreleased]`                                | `CHANGELOG.md`                   |
| 15  | README: component counts (73→76), examples (Grid, Href, RightItems)    | `README.md`                      |
| 16  | Feedback appendices: resolution status on all 5 feedback files         | `docs/feedback/*.md`             |
| 17  | SKILL.md rewritten to 10/10 — Part 1 Consumer Guide + Part 2 Authoring | `skill/SKILL.md`                 |
| 18  | 3 status reports + 1 planning doc                                      | `docs/status/`, `docs/planning/` |

### Code fixes (commits `79c926c`, `29fccf1`)

| #   | What                                                                                | Why                                             |
| --- | ----------------------------------------------------------------------------------- | ----------------------------------------------- |
| 19  | Fixed `GridCols4`/`GridCols5` responsive ladders (added intermediate md breakpoint) | Design flaw — jumped 2→final                    |
| 20  | Modernized ProgressBar clamp to `max(0, min(100, v))`                               | templ minmax diagnostic                         |
| 21  | Fixed stale `sidebar_nav.golden`                                                    | Pre-existing failure from templ cosmetic change |
| 22  | Fixed 4 lint errors in `sri_net_test.go` (errcheck/noctx/paralleltest)              | Pre-existing                                    |
| 23  | Removed `*_templ.go` from `.gitignore` line 32                                      | Root cause of BuildFlow gotcha                  |
| 24  | Fixed README feedback count 12→13 (missed twice in self-reviews)                    | Cosmetic but embarrassing                       |
| 25  | Fixed AGENTS.md BaseProps count 25→26                                               | Accuracy                                        |
| 26  | Demo updated: StatCard section uses `display.Grid` + `StatCard.Href`                | `examples/demo/demo.templ`                      |

### Test lens coverage

| Component                   | golden | BDD | a11y       | example    | integration | coverage |
| --------------------------- | ------ | --- | ---------- | ---------- | ----------- | -------- |
| `display.Grid`              | ✅     | ✅  | ✅         | ✅         | ✅          | ✅       |
| `StatCard.Href`             | ✅     | ✅  | ✅         | (existing) | ✅          | ✅       |
| `layout.Script`             | ✅     | ✅  | ✅         | ✅         | —           | ✅       |
| `feedback.SkeletonCardGrid` | ✅     | ✅  | ✅         | ✅         | —           | ✅       |
| `SimpleNav.RightItems`      | —      | ✅  | (existing) | —          | —           | ✅       |

### Contract / inventory

| #   | What                                                                      |
| --- | ------------------------------------------------------------------------- |
| 27  | `GridProps` registered in `internal/contract/component_props_test.go`     |
| 28  | `statCardInner` sub-template extracted (DRY for linked/unlinked StatCard) |

---

## b) PARTIALLY DONE

| Item                         | Done                                                                  | Missing                                                                                               |
| ---------------------------- | --------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| Consumer discoverability     | SKILL.md Part 1 has full catalogue; README has examples; demo updated | No standalone demo site; forms package still not "flagship" placement; no auto-generated catalog page |
| Consumer Feedback Backlog    | 11 items documented in TODO_LIST.md with sources                      | None implemented yet (CopyButton, RelativeTime, cursor pagination, etc.)                              |
| Tailwind v4 migration recipe | 7-step guide written                                                  | Not tested end-to-end by actually migrating a real consumer project                                   |

---

## c) NOT STARTED

| #   | What                             | Why                                                     | Status (2026-07-06)                              |
| --- | -------------------------------- | ------------------------------------------------------- | ------------------------------------------------ |
| 1   | Implement CopyButton component   | DiscordSync feedback — in backlog                       | ✅ Done — shipped in session 7                   |
| 2   | Implement RelativeTime component | DiscordSync feedback — in backlog                       | ✅ Done — shipped in session 7                   |
| 3   | Implement cursor pagination      | DiscordSync feedback — in backlog                       | ✅ Done — `navigation.LoadMore` shipped          |
| 4   | Forms discoverability overhaul   | SwettySwipper feedback — #1 gap, needs design decision  | ✅ Done — SKILL.md "by use case" table + recipes |
| 5   | Component catalog demo site      | Multiple consumers — needs hosting decision             | ⬜ Not started                                   |
| 6   | v0.7.0 release cut               | All `[Unreleased]` entries ready, release script tested | ✅ Done (v0.7.0 + v0.8.0 released)               |

---

## d) TOTALLY FUCKED UP

Nothing is broken. Verify passes, git is clean, all pushed.

**Judgment failures this session (all fixed):**

| #   | What                                                                               | How fixed                                          |
| --- | ---------------------------------------------------------------------------------- | -------------------------------------------------- |
| 1   | README feedback count "12" missed in 2 consecutive commits                         | Fixed in `29fccf1`                                 |
| 2   | `layout.Script` shipped with assertion-only tests (violating the skill's own rule) | Fixed in `29fccf1` — added golden+BDD+a11y+example |
| 3   | `.gitignore` root cause worked around with `git add -f` instead of fixed           | Fixed in `29fccf1` — removed line 32               |
| 4   | SKILL.md forgotten entirely in first cleanup pass                                  | Fixed in `cea5a66` — full rewrite                  |
| 5   | `GridCols5` shipped with bad responsive ladder (2→5 jump)                          | Fixed in `79c926c`                                 |
| 6   | AGENTS.md not updated in first pass                                                | Fixed in `79c926c`                                 |

**Remaining risk:**

| #   | What                                                                                                                                          | Severity                                         | Status (2026-07-06)                                         |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------ | ----------------------------------------------------------- |
| 1   | BuildFlow's `gitignore-upserter` may re-add `*_templ.go` on next commit — the `.gitignore` fix may not be permanent if BuildFlow overrides it | Unknown — needs monitoring on next BuildFlow run | ✅ **Confirmed fixed** — `.gitignore` stable through v0.8.0 |

---

## e) WHAT WE SHOULD IMPROVE

1. **Cut v0.7.0** — `[Unreleased]` is comprehensive, all tests pass, 5 new features + 2 recipes. No reason to wait.
2. **Monitor the `.gitignore` fix** — if BuildFlow re-adds `*_templ.go`, we need to fix it in BuildFlow itself (`larsartmann/buildflow`).
3. **Forms discoverability** — SwettySwipper's #1 request. Needs a design decision: prominent README section? Separate forms demo page? Forms quickstart guide?
4. **Component catalog** — auto-generated or manual, consumers need a single-page "what exists" reference.
5. **Consumer adoption testing** — actually migrate a real project (browser-history or DiscordSync) to validate the recipes work end-to-end.

---

## f) Up to 25 things to do next

| #   | Task                                                                                               | Impact | Effort   | Status (2026-07-06)                                        |
| --- | -------------------------------------------------------------------------------------------------- | ------ | -------- | ---------------------------------------------------------- |
| 1   | Cut v0.7.0 release via `scripts/release.sh`                                                        | High   | 10m      | ✅ Done (v0.7.0 + v0.8.0)                                  |
| 2   | Monitor `.gitignore` after next BuildFlow run — verify `*_templ.go` isn't re-added                 | High   | 2m       | ✅ Done — stable through v0.8.0                            |
| 3   | Forms discoverability: add prominent forms section to README with quickstart example               | High   | 15m      | ✅ Done                                                    |
| 4   | Forms demo page in `examples/demo/`                                                                | Med    | 20m      | ⬜ Not started                                             |
| 5   | Implement `display.CopyButton` (clipboard API + "Copied!" feedback)                                | Med    | 15m      | ✅ Done                                                    |
| 6   | Implement `display.RelativeTime(timestamp)`                                                        | Med    | 15m      | ✅ Done                                                    |
| 7   | Implement cursor pagination pattern (document or `navigation.LoadMore`)                            | Med    | 20m      | ✅ Done                                                    |
| 8   | Auto-generate component catalog from source (script that greps `templ [A-Z]`)                      | Med    | 30m      | ⬜ Not started                                             |
| 9   | Count badge overlay on icon (DiscordSync)                                                          | Low    | 15m      | ✅ Done (`display.CountBadge`)                             |
| 10  | `display.DefinitionGrid` wrapper (DiscordSync)                                                     | Low    | 10m      | ✅ Done                                                    |
| 11  | `display.Image` with lazy loading + aspect ratio (SwettySwipper)                                   | Low    | 20m      | ✅ Done                                                    |
| 12  | Consider self-hosting htmx as default (v1.0 breaking change decision)                              | High   | Decision | ✅ Done — ADR 0007 written, deferred to v1.0               |
| 13  | Consider typed HTMX fields on StatCard vs Attrs workaround                                         | Med    | Decision | ✅ Done                                                    |
| 14  | Consider `Card.Body` explicit slot (SEC feedback)                                                  | Low    | 15m      | ✅ Done                                                    |
| 15  | StatCard golden with Href + Icon combined                                                          | Low    | 5m       | ✅ Done                                                    |
| 16  | Test Play CDN migration recipe end-to-end on browser-history                                       | Med    | 30m      | ⬜ Not started                                             |
| 17  | Test HTMX error feedback recipe end-to-end on a real project                                       | Med    | 30m      | ⬜ Not started                                             |
| 18  | Add `GridProps.Gap` typed enum (gap-2/4/6/8)                                                       | Low    | 10m      | ⬜ Not done                                                |
| 19  | Consider `layout.Stylesheet(nonce, href, attrs)` companion to `Script`                             | Low    | 10m      | ✅ Done                                                    |
| 20  | Audit component count 76 by actual grep across all packages                                        | Low    | 5m       | ✅ Done (82 components)                                    |
| 21  | Add CI check that `*_templ.go` files are tracked (prevent future gotcha)                           | Med    | 15m      | ⬜ Not needed — .gitignore fixed                           |
| 22  | Consider sortable `display.Table` (typed column definitions)                                       | Med    | 30m      | ✅ Done (`TableHeader` + `TypedHeaders` shipped in v0.8.0) |
| 23  | `examples/demo/` add SkeletonCardGrid loading state showcase                                       | Low    | 5m       | ⬜ Not done                                                |
| 24  | Consumer project: actually adopt templ-components in DiscordSync to validate discoverability fixes | High   | 60m      | ⬜ Not started                                             |
| 25  | v1.0 API freeze planning (move test helpers, Validate() error, freeze types)                       | High   | 60m      | ⬜ Not started                                             |

**Scorecard:** 14 of 25 complete (56%).

---

## g) Top #1 question I cannot figure out myself

> ✅ **RESOLVED.** The `.gitignore` fix held. Confirmed stable through v0.8.0 — the
> `gitignore-upserter` no longer re-adds `*_templ.go`. Current `.gitignore` has only
> `!*_templ.go` on line 2 with no trailing override.
