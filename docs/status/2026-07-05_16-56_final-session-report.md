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

| #     | Component / Feature                                                        | Files                                            | Tests                                                 |
| ----- | -------------------------------------------------------------------------- | ------------------------------------------------ | ----------------------------------------------------- |
| ~~1~~ | ~~`display.Grid` + `GridCols` enum (1–6)~~ done at `985019a`               | ~~`display/grid.templ`, `grid_templ.go`~~        | ~~golden, BDD, a11y, example, integration, coverage~~ |
| ~~2~~ | ~~`StatCardProps.Href` — renders `<a>` wrapper~~ done at `985019a`         | ~~`display/card.templ`, `card_templ.go`~~        | ~~golden, BDD, a11y, integration, coverage~~          |
| ~~3~~ | ~~`SimpleNavProps.RightItems` — forwarded to Nav~~ done at `985019a`       | ~~`navigation/nav.templ`, `nav_templ.go`~~       | ~~coverage, BDD~~                                     |
| ~~4~~ | ~~`layout.Script(nonce, src, attrs)` — CSP-safe helper~~ done at `985019a` | ~~`layout/script.templ`, `script_templ.go`~~     | ~~golden, BDD, a11y, example, snapshot~~              |
| ~~5~~ | ~~`feedback.SkeletonCardGrid(count)` — loading grid~~ done at `985019a`    | ~~`feedback/loading.templ`, `loading_templ.go`~~ | ~~golden, BDD, a11y, example, snapshot~~              |

### Recipes shipped (commit `985019a`)

| #     | Recipe                                                                     | Path                                                      |
| ----- | -------------------------------------------------------------------------- | --------------------------------------------------------- |
| ~~6~~ | ~~Play CDN → Tailwind v4 CSS-first (7-step migration)~~ done at `985019a`  | ~~`docs/migration/play-cdn-to-tailwind-v4.md`~~           |
| ~~7~~ | ~~Server-rendered HTMX error feedback (3 render modes)~~ done at `985019a` | ~~`docs/recipes/server-rendered-htmx-error-feedback.md`~~ |

### Documentation (commits `985019a`, `79c926c`, `865967a`, `cea5a66`, `29fccf1`)

| #      | What                                                                                         | Files                                |
| ------ | -------------------------------------------------------------------------------------------- | ------------------------------------ |
| ~~8~~  | ~~PageProps auto-inject godoc (HTMXVersion, CSSPath suppression)~~ done at `358d0b0`         | ~~`layout/base.templ`~~              |
| ~~9~~  | ~~README "Suppressing auto-injected `<head>` tags" subsection~~ done at `358d0b0`            | ~~`README.md`~~                      |
| ~~10~~ | ~~AGENTS.md: 8 new convention entries, count 25→26~~ done at `358d0b0`                       | ~~`AGENTS.md`~~                      |
| ~~11~~ | ~~TODO_LIST.md: session 6 header + Consumer Feedback Backlog (11 items)~~ done at `358d0b0`  | ~~`TODO_LIST.md`~~                   |
| ~~12~~ | ~~FEATURES.md: Grid, Script, SkeletonCardGrid, GridCols enum~~ done at `358d0b0`             | ~~`FEATURES.md`~~                    |
| ~~13~~ | ~~CONTEXT.md: updated metrics, package descriptions~~ done at `358d0b0`                      | ~~`CONTEXT.md`~~                     |
| ~~14~~ | ~~CHANGELOG: comprehensive `[Unreleased]`~~ done at `358d0b0`                                | ~~`CHANGELOG.md`~~                   |
| ~~15~~ | ~~README: component counts (73→76), examples (Grid, Href, RightItems)~~ done at `358d0b0`    | ~~`README.md`~~                      |
| ~~16~~ | ~~Feedback appendices: resolution status on all 5 feedback files~~ done at `358d0b0`         | ~~`docs/feedback/*.md`~~             |
| ~~17~~ | ~~SKILL.md rewritten to 10/10 — Part 1 Consumer Guide + Part 2 Authoring~~ done at `358d0b0` | ~~`skill/SKILL.md`~~                 |
| ~~18~~ | ~~3 status reports + 1 planning doc~~ done at `358d0b0`                                      | ~~`docs/status/`, `docs/planning/`~~ |

### Code fixes (commits `79c926c`, `29fccf1`)

| #      | What                                                                                                      | Why                                                 |
| ------ | --------------------------------------------------------------------------------------------------------- | --------------------------------------------------- |
| ~~19~~ | ~~Fixed `GridCols4`/`GridCols5` responsive ladders (added intermediate md breakpoint)~~ done at `79c926c` | ~~Design flaw — jumped 2→final~~                    |
| ~~20~~ | ~~Modernized ProgressBar clamp to `max(0, min(100, v))`~~ done at `358d0b0`                               | ~~templ minmax diagnostic~~                         |
| ~~21~~ | ~~Fixed stale `sidebar_nav.golden`~~ done at `358d0b0`                                                    | ~~Pre-existing failure from templ cosmetic change~~ |
| ~~22~~ | ~~Fixed 4 lint errors in `sri_net_test.go` (errcheck/noctx/paralleltest)~~ done at `358d0b0`              | ~~Pre-existing~~                                    |
| ~~23~~ | ~~Removed `*_templ.go` from `.gitignore` line 32~~ done at `29fccf1`                                      | ~~Root cause of BuildFlow gotcha~~                  |
| ~~24~~ | ~~Fixed README feedback count 12→13 (missed twice in self-reviews)~~ done at `29fccf1`                    | ~~Cosmetic but embarrassing~~                       |
| ~~25~~ | ~~Fixed AGENTS.md BaseProps count 25→26~~ done at `358d0b0`                                               | ~~Accuracy~~                                        |
| ~~26~~ | ~~Demo updated: StatCard section uses `display.Grid` + `StatCard.Href`~~ done at `358d0b0`                | ~~`examples/demo/demo.templ`~~                      |

### Test lens coverage

| Component                   | golden | BDD | a11y       | example    | integration | coverage |
| --------------------------- | ------ | --- | ---------- | ---------- | ----------- | -------- |
| `display.Grid`              | ✅     | ✅  | ✅         | ✅         | ✅          | ✅       |
| `StatCard.Href`             | ✅     | ✅  | ✅         | (existing) | ✅          | ✅       |
| `layout.Script`             | ✅     | ✅  | ✅         | ✅         | —           | ✅       |
| `feedback.SkeletonCardGrid` | ✅     | ✅  | ✅         | ✅         | —           | ✅       |
| `SimpleNav.RightItems`      | —      | ✅  | (existing) | —          | —           | ✅       |

### Contract / inventory

| #      | What                                                                                            |
| ------ | ----------------------------------------------------------------------------------------------- |
| ~~27~~ | ~~`GridProps` registered in `internal/contract/component_props_test.go`~~ done at `358d0b0`     |
| ~~28~~ | ~~`statCardInner` sub-template extracted (DRY for linked/unlinked StatCard)~~ done at `358d0b0` |

---

## b) PARTIALLY DONE

| Item                         | Done                                                                  | Missing                                                                                               |
| ---------------------------- | --------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| Consumer discoverability     | SKILL.md Part 1 has full catalogue; README has examples; demo updated | No standalone demo site; forms package still not "flagship" placement; no auto-generated catalog page |
| Consumer Feedback Backlog    | 11 items documented in TODO_LIST.md with sources                      | None implemented yet (CopyButton, RelativeTime, cursor pagination, etc.)                              |
| Tailwind v4 migration recipe | 7-step guide written                                                  | Not tested end-to-end by actually migrating a real consumer project                                   |

---

## c) NOT STARTED

| #     | What                                                                    | Why                                                         | Status (2026-07-06)                                  |
| ----- | ----------------------------------------------------------------------- | ----------------------------------------------------------- | ---------------------------------------------------- |
| ~~1~~ | ~~Implement CopyButton component~~ done — display/copy button.templ     | ~~DiscordSync feedback — in backlog~~                       | ~~✅ Done — shipped in session 7~~                   |
| ~~2~~ | ~~Implement RelativeTime component~~ done — display/relative time.templ | ~~DiscordSync feedback — in backlog~~                       | ~~✅ Done — shipped in session 7~~                   |
| ~~3~~ | ~~Implement cursor pagination~~ done — navigation/loadmore.templ        | ~~DiscordSync feedback — in backlog~~                       | ~~✅ Done — `navigation.LoadMore` shipped~~          |
| ~~4~~ | ~~Forms discoverability overhaul~~ done — skill/SKILL.md                | ~~SwettySwipper feedback — #1 gap, needs design decision~~  | ~~✅ Done — SKILL.md "by use case" table + recipes~~ |
| ~~5~~ | ~~Component catalog demo site~~ done — visualtest/demo.go               | ~~Multiple consumers — needs hosting decision~~             | ~~⬜ Not started~~                                   |
| ~~6~~ | ~~v0.7.0 release cut~~ done — CHANGELOG.md                              | ~~All `[Unreleased]` entries ready, release script tested~~ | ~~✅ Done (v0.7.0 + v0.8.0 released)~~               |

---

## d) TOTALLY FUCKED UP

Nothing is broken. Verify passes, git is clean, all pushed.

**Judgment failures this session (all fixed):**

| #     | What                                                                                                     | How fixed                                              |
| ----- | -------------------------------------------------------------------------------------------------------- | ------------------------------------------------------ |
| 1     | README feedback count "12" missed in 2 consecutive commits                                               | Fixed in `29fccf1`                                     |
| ~~2~~ | ~~`layout.Script` shipped with assertion-only tests (violating the skill's own rule)~~ done at `29fccf1` | ~~Fixed in `29fccf1` — added golden+BDD+a11y+example~~ |
| ~~3~~ | ~~`.gitignore` root cause worked around with `git add -f` instead of fixed~~ done at `29fccf1`           | ~~Fixed in `29fccf1` — removed line 32~~               |
| ~~4~~ | ~~SKILL.md forgotten entirely in first cleanup pass~~ done at `cea5a66`                                  | ~~Fixed in `cea5a66` — full rewrite~~                  |
| ~~5~~ | ~~`GridCols5` shipped with bad responsive ladder (2→5 jump)~~ done at `79c926c`                          | ~~Fixed in `79c926c`~~                                 |
| ~~6~~ | ~~AGENTS.md not updated in first pass~~ done at `79c926c`                                                | ~~Fixed in `79c926c`~~                                 |

**Remaining risk:**

| # | What                                                                                                                                          | Severity                                         | Status (2026-07-06)                                         |
| - | --------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------ | ----------------------------------------------------------- |
| 1 | BuildFlow's `gitignore-upserter` may re-add `*_templ.go` on next commit — the `.gitignore` fix may not be permanent if BuildFlow overrides it | Unknown — needs monitoring on next BuildFlow run | ✅ **Confirmed fixed** — `.gitignore` stable through v0.8.0 |

---

## e) WHAT WE SHOULD IMPROVE

1. ~~**Cut v0.7.0** — `[Unreleased]` is comprehensive, all tests pass, 5 new features + 2 recipes. No reason to wait.~~ done — CHANGELOG.md
2. ~~**Monitor the `.gitignore` fix** — if BuildFlow re-adds `*_templ.go`, we need to fix it in BuildFlow itself (`larsartmann/buildflow`).~~ done (docs-health pass 2026-09-08)
3. ~~**Forms discoverability** — SwettySwipper's #1 request. Needs a design decision: prominent README section? Separate forms demo page? Forms quickstart guide?~~ done — README.md
4. ~~**Component catalog** — auto-generated or manual, consumers need a single-page "what exists" reference.~~ done — visualtest/demo.go
5. **Consumer adoption testing** — actually migrate a real project (browser-history or DiscordSync) to validate the recipes work end-to-end.

---

## f) Up to 25 things to do next

| #      | Task                                                                                                                     | Impact   | Effort       | Status (2026-07-06)                                            |
| ------ | ------------------------------------------------------------------------------------------------------------------------ | -------- | ------------ | -------------------------------------------------------------- |
| ~~1~~  | ~~Cut v0.7.0 release via `scripts/release.sh`~~ done — CHANGELOG.md                                                      | ~~High~~ | ~~10m~~      | ~~✅ Done (v0.7.0 + v0.8.0)~~                                  |
| ~~2~~  | ~~Monitor `.gitignore` after next BuildFlow run — verify `*_templ.go` isn't re-added~~ done — .gitignore                 | ~~High~~ | ~~2m~~       | ~~✅ Done — stable through v0.8.0~~                            |
| ~~3~~  | ~~Forms discoverability: add prominent forms section to README with quickstart example~~ done — README.md                | ~~High~~ | ~~15m~~      | ~~✅ Done~~                                                    |
| ~~4~~  | ~~Forms demo page in `examples/demo/`~~ done — visualtest/demo.go                                                        | ~~Med~~  | ~~20m~~      | ~~⬜ Not started~~                                             |
| ~~5~~  | ~~Implement `display.CopyButton` (clipboard API + "Copied!" feedback)~~ done — display/copy button.templ                 | ~~Med~~  | ~~15m~~      | ~~✅ Done~~                                                    |
| ~~6~~  | ~~Implement `display.RelativeTime(timestamp)`~~ done — display/relative time.templ                                       | ~~Med~~  | ~~15m~~      | ~~✅ Done~~                                                    |
| ~~7~~  | ~~Implement cursor pagination pattern (document or `navigation.LoadMore`)~~ done — navigation/loadmore.templ             | ~~Med~~  | ~~20m~~      | ~~✅ Done~~                                                    |
| 8      | Auto-generate component catalog from source (script that greps `templ [A-Z]`)                                            | Med      | 30m          | ⬜ Not started                                                 |
| ~~9~~  | ~~Count badge overlay on icon (DiscordSync)~~ done — display/count badge.templ                                           | ~~Low~~  | ~~15m~~      | ~~✅ Done (`display.CountBadge`)~~                             |
| ~~10~~ | ~~`display.DefinitionGrid` wrapper (DiscordSync)~~ done — display/definition grid.templ                                  | ~~Low~~  | ~~10m~~      | ~~✅ Done~~                                                    |
| ~~11~~ | ~~`display.Image` with lazy loading + aspect ratio (SwettySwipper)~~ done — display/image.templ                          | ~~Low~~  | ~~20m~~      | ~~✅ Done~~                                                    |
| ~~12~~ | ~~Consider self-hosting htmx as default (v1.0 breaking change decision)~~ done — docs/adr/0007-self-host-htmx-default.md | ~~High~~ | ~~Decision~~ | ~~✅ Done — ADR 0007 written, deferred to v1.0~~               |
| ~~13~~ | ~~Consider typed HTMX fields on StatCard vs Attrs workaround~~ done — display/card.templ                                 | ~~Med~~  | ~~Decision~~ | ~~✅ Done~~                                                    |
| ~~14~~ | ~~Consider `Card.Body` explicit slot (SEC feedback)~~ done — display/card.templ                                          | ~~Low~~  | ~~15m~~      | ~~✅ Done~~                                                    |
| ~~15~~ | ~~StatCard golden with Href + Icon combined~~ done — display/testdata                                                    | ~~Low~~  | ~~5m~~       | ~~✅ Done~~                                                    |
| 16     | Test Play CDN migration recipe end-to-end on browser-history                                                             | Med      | 30m          | ⬜ Not started                                                 |
| 17     | Test HTMX error feedback recipe end-to-end on a real project                                                             | Med      | 30m          | ⬜ Not started                                                 |
| ~~18~~ | ~~Add `GridProps.Gap` typed enum (gap-2/4/6/8)~~ done — display/grid.templ                                               | ~~Low~~  | ~~10m~~      | ~~⬜ Not done~~                                                |
| ~~19~~ | ~~Consider `layout.Stylesheet(nonce, href, attrs)` companion to `Script`~~ done — layout/stylesheet.templ                | ~~Low~~  | ~~10m~~      | ~~✅ Done~~                                                    |
| ~~20~~ | ~~Audit component count 76 by actual grep across all packages~~ done — utils/skill count test.go                         | ~~Low~~  | ~~5m~~       | ~~✅ Done (82 components)~~                                    |
| ~~21~~ | ~~Add CI check that `*_templ.go` files are tracked (prevent future gotcha)~~ **Won't implement — root cause fixed.**     | ~~Med~~  | ~~15m~~      | ~~⬜ Not needed — .gitignore fixed~~                           |
| ~~22~~ | ~~Consider sortable `display.Table` (typed column definitions)~~ done — display/table.templ                              | ~~Med~~  | ~~30m~~      | ~~✅ Done (`TableHeader` + `TypedHeaders` shipped in v0.8.0)~~ |
| ~~23~~ | ~~`examples/demo/` add SkeletonCardGrid loading state showcase~~ done — visualtest/visual test.go                        | ~~Low~~  | ~~5m~~       | ~~⬜ Not done~~                                                |
| 24     | Consumer project: actually adopt templ-components in DiscordSync to validate discoverability fixes                       | High     | 60m          | ⬜ Not started                                                 |
| ~~25~~ | ~~v1.0 API freeze planning (move test helpers, Validate() error, freeze types)~~ done — docs/migration/v0.22-to-v1.0.md  | ~~High~~ | ~~60m~~      | ~~⬜ Not started~~                                             |

**Scorecard:** 14 of 25 complete (56%).

---

## g) Top #1 question I cannot figure out myself

> ✅ **RESOLVED.** The `.gitignore` fix held. Confirmed stable through v0.8.0 — the
> `gitignore-upserter` no longer re-adds `*_templ.go`. Current `.gitignore` has only
> `!*_templ.go` on line 2 with no trailing override.
